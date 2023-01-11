package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"text/template"
	"time"

	redis "github.com/go-redis/redis/v8"
	"github.com/gorilla/websocket"
	"github.com/weedge/craftsman/doraemon/ai-creator/internal/api"
)

const (
	// Time allowed to write the file to the client.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the client.
	pongWait = 60 * time.Second

	// Send pings to client with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10
)

var (
	rdb           redis.UniversalClient
	mapSub        = map[string]*redis.PubSub{}
	mapSubHandler = map[string]SubHandler{
		"draftTaskStatusCh": PullDraftTaskStatus,
	}

	templates = template.Must(template.New("").Funcs(createFuncMap()).ParseGlob(os.Getenv("TEMPLATE_DIR") + "*.gohtml"))

	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	statusResCh = make(chan []byte, 1)
)

func createFuncMap() template.FuncMap {
	// funcmap for go templates
	return template.FuncMap{
		/*

			"getTimeAgo": func(t time.Time) string {
				return timeago.English.Format(t)
			},
			"getTimeAgoForMillis": func(tUnix int64) string {
				return getTimeAgoForMillis(tUnix)
			},

		*/
	}
}

func reader(ws *websocket.Conn) {
	defer ws.Close()
	ws.SetReadLimit(512)
	ws.SetReadDeadline(time.Now().Add(pongWait))
	ws.SetPongHandler(func(string) error { ws.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, msg, err := ws.ReadMessage()
		if err != nil {
			log.Println("receive err:", err.Error())
			break
		}
		log.Println("receive:", msg)
	}
}

func writer(ws *websocket.Conn) {
	pingTicker := time.NewTicker(pingPeriod)
	defer func() {
		pingTicker.Stop()
		ws.Close()
	}()
	for {
		select {
		case res := <-statusResCh:
			ws.SetWriteDeadline(time.Now().Add(writeWait))
			if err := ws.WriteMessage(websocket.TextMessage, res); err != nil {
				log.Println("ws.WriteMessage err", err.Error())
				return
			}
			log.Printf("%s ws.WriteMessage ok", res)
		case <-pingTicker.C:
			ws.SetWriteDeadline(time.Now().Add(writeWait))
			if err := ws.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				log.Println("ws.WriteMessage err", err.Error())
				return
			}
			log.Println("Ping ws.WriteMessage ok")
		}
	}
}

func serveWsHandler(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		if _, ok := err.(websocket.HandshakeError); !ok {
			log.Println(err.Error())
		}
		return
	}

	go writer(ws)
	reader(ws)
}

func serveHomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	var v = struct {
		Host string
	}{
		r.Host,
	}
	templates.ExecuteTemplate(w, "home.gohtml", &v)
}

func main() {
	flag.Parse()
	addr := flag.String("addr", ":8123", "http service address")

	setupRoutes()
	ctx := context.Background()
	connectToRedis(ctx, os.Getenv("REDIS_TYPE"))
	initSub(ctx)

	go func() {
		log.Printf("Go Web App Started on Port %s", *addr)
		log.Fatal(http.ListenAndServe(*addr, nil))
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	closeSub(ctx)
}

func setupRoutes() {
	http.HandleFunc("/", serveHomeHandler)
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/readiness", readinessHandler)
	http.HandleFunc("/ws", serveWsHandler)
	http.HandleFunc("/text2img", text2ImgHandler)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func readinessHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func text2ImgHandler(w http.ResponseWriter, r *http.Request) {
	keyword := r.FormValue("keyword")
	authVal := api.GetAuth(os.Getenv("NOLIBOX_API_AK"), os.Getenv("NOLIBOX_API_SK"))

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	respBody, err := api.Draft(authVal, "txt2img.sd", keyword, "")
	if err != nil {
		myMap := map[string]interface{}{}
		myMap["error"] = err.Error()
		tmp, _ := json.Marshal(myMap)
		w.Write(tmp)
		return
	}

	w.Write(respBody)
	data := api.DraftResp{}
	err = json.Unmarshal(respBody, &data)
	if err != nil {
		log.Println(respBody, "json.Unmarshal err", err.Error())
		return
	}
	pubSubChName := "draftTaskStatusCh"
	err = rdb.Publish(r.Context(), pubSubChName, data.Data.UID).Err()
	if err != nil {
		log.Println("pub err", err.Error(), data.Data.UID)
		return
	}
	log.Println(pubSubChName, "pub ok", data.Data.UID)
}

func connectToRedis(ctx context.Context, redisType string) {
	log.Println(os.Getenv("REDIS_URL"), os.Getenv("REDIS_PASSWORD"))

	switch redisType {
	case "cluster":
		rdb = redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:    []string{os.Getenv("REDIS_URL")},
			Password: os.Getenv("REDIS_PASSWORD"),
			Username: os.Getenv("REDIS_USERNAME"),

			MaxRetries:      3,
			MinRetryBackoff: 3 * time.Second,
			MaxRetryBackoff: 5 * time.Second,
			DialTimeout:     5 * time.Second,
			ReadTimeout:     3 * time.Second,
			WriteTimeout:    3 * time.Second,

			// connect pool
			PoolSize:           100,
			MinIdleConns:       10,
			MaxConnAge:         60 * time.Second,
			PoolTimeout:        5 * time.Second,
			IdleTimeout:        30 * time.Second,
			IdleCheckFrequency: 3 * time.Second,

			// To route commands by latency or randomly, enable one of the following.
			//RouteByLatency: true,
			RouteRandomly: true,
		})
	default:
		rdb = redis.NewClient(&redis.Options{
			Addr:     os.Getenv("REDIS_URL"),
			Password: os.Getenv("REDIS_PASSWORD"), // no password set
			Username: os.Getenv("REDIS_USERNAME"),
			DB:       0, // use default DB

			MaxRetries:      3,
			MinRetryBackoff: 3 * time.Second,
			MaxRetryBackoff: 5 * time.Second,
			DialTimeout:     5 * time.Second,
			ReadTimeout:     3 * time.Second,
			WriteTimeout:    3 * time.Second,

			// connect pool
			PoolSize:           100,
			MinIdleConns:       10,
			MaxConnAge:         60 * time.Second,
			PoolTimeout:        5 * time.Second,
			IdleTimeout:        30 * time.Second,
			IdleCheckFrequency: 3 * time.Second,
		})
	}

	pong, err := rdb.Ping(ctx).Result()
	if err == nil {
		log.Println(pong, "ok")
	} else {
		log.Fatalln(err.Error())
	}
}

type SubHandler func(ctx context.Context, msg *redis.Message) error

func PullDraftTaskStatus(ctx context.Context, msg *redis.Message) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	err := pullDraftTaskStatus(ctx, msg.Payload)
	return err
}

func initSub(ctx context.Context) {
	for subChName, handle := range mapSubHandler {
		mapSub[subChName] = rdb.Subscribe(ctx, subChName)
		go func(subChName string, handle SubHandler) {
			for {
				msg, err := mapSub[subChName].ReceiveMessage(ctx)
				if err != nil {
					log.Println(subChName, "receiveMsg err", err.Error())
					return
				}
				log.Printf("msg:%+v\n", msg)
				handle(ctx, msg)
			}
		}(subChName, handle)
		log.Println(subChName, "sub ok")
	}
}

func closeSub(ctx context.Context) {
	for name, item := range mapSub {
		err := item.Close()
		if err != nil {
			log.Println("sub", name, "err", err.Error())
		}
		log.Println("sub", name, "close ok")
	}
}

func pullDraftTaskStatus(ctx context.Context, uid string) error {
	authVal := api.GetAuth(os.Getenv("NOLIBOX_API_AK"), os.Getenv("NOLIBOX_API_SK"))
	ticker := time.NewTicker(1000 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			body, err := api.Status(authVal, uid)
			if err != nil {
				log.Println("status err", err.Error())
				continue
			}
			data := api.StatusResp{}
			err = json.Unmarshal(body, &data)
			if err != nil {
				log.Println("unmarshal err", err.Error())
				return err
			}
			switch data.Status {
			case "finished":
				statusResCh <- body
				return nil
			}
		case <-ctx.Done():
			log.Println("pullDraftTaskStatus done return")
			return nil
		}
	}
}
