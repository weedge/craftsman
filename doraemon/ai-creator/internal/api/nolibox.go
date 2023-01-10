package api

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func Img2acg(authVal string, imgUrl string) {
	url := "https://open.nolibox.com/prod-open-aigc/engine/img2acg"
	method := "POST"

	payload := strings.NewReader(fmt.Sprintf(`{
    "style": "Gorgeous_style",
    "image_url": %s,
    "use_seed": true,
    "keep_alpha": true,
    "w": 512,
    "h": 512,
    "fidelity": 0.3,
    "use_circular": false,
    "is_anime": true,
    "seed": 2180138447,
    "version": "v1.5",
    "promptInfo": {
        "lang": "zh",
        "original": "a cute European young lady and a handsome European young man, big eyes, beautiful detailed eyes, smile, ((style reference J.C. Leyendecker)), portrait, illustration, masterpiece, best quality, CG, HD, 8k, 4k, highly detailed, wallpaper",
        "originalNegative": "cropped, blurred, mutated, error, lowres, blurry, low quality, username, signature, watermark, text, nsfw, missing limb, fused hand, missing hand, extra limbs, malformed limbs, bad hands, extra fingers, fused fingers, missing fingers, bad breasts, deformed, mutilated, morbid, bad anatomy"
    },
    "timestamp": 1670569311577,
    "localModels": {},
    "negative_prompt": "cropped, blurred, mutated, error, lowres, blurry, low quality, username, signature, watermark, text, nsfw, missing limb, fused hand, missing hand, extra limbs, malformed limbs, bad hands, extra fingers, fused fingers, missing fingers, bad breasts, deformed, mutilated, morbid, bad anatomy",
    "sampler": "k_euler",
    "num_steps": 25,
    "guidance_scale": 7.5
}`, imgUrl))

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Authorization", fmt.Sprintf("Basic %s", authVal))
	req.Header.Add("User-Agent", "Apifox/1.0.0 (https://www.apifox.cn)")
	req.Header.Add("Content-Type", "application/json")

	fmt.Printf("%+v", req)
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(body))
}

type StatusResp struct {
	Data struct {
		Cdn          string  `json:"cdn"`
		Cos          string  `json:"cos"`
		CreateTime   float64 `json:"create_time"`
		Duration     float64 `json:"duration"`
		ElapsedTimes struct {
			AlgorithmLatencies struct {
				Total      float64 `json:"__total__"`
				Cleanup    float64 `json:"cleanup"`
				Download   float64 `json:"download"`
				GetModel   float64 `json:"get_model"`
				Inference  float64 `json:"inference"`
				Preprocess float64 `json:"preprocess"`
			} `json:"algorithm_latencies"`
			Audit        float64 `json:"audit"`
			Callback     float64 `json:"callback"`
			Pending      float64 `json:"pending"`
			Redis        float64 `json:"redis"`
			RunAlgorithm float64 `json:"run_algorithm"`
			Upload       float64 `json:"upload"`
		} `json:"elapsed_times"`
		EndTime float64  `json:"end_time"`
		ImgUrls []string `json:"img_urls"`
		Reason  string   `json:"reason"`
		Request struct {
			Model struct {
				CallbackURL      string `json:"callback_url"`
				ClipSkip         int    `json:"clip_skip"`
				CustomEmbeddings struct {
				} `json:"custom_embeddings"`
				Fidelity          float64       `json:"fidelity"`
				GuidanceScale     float64       `json:"guidance_scale"`
				IsAnime           bool          `json:"is_anime"`
				KeepAlpha         bool          `json:"keep_alpha"`
				MaxWh             int           `json:"max_wh"`
				NegativePrompt    string        `json:"negative_prompt"`
				NumSteps          int           `json:"num_steps"`
				Sampler           string        `json:"sampler"`
				Seed              int64         `json:"seed"`
				Text              string        `json:"text"`
				URL               string        `json:"url"`
				UseCircular       bool          `json:"use_circular"`
				VariationSeed     int           `json:"variation_seed"`
				VariationStrength float64       `json:"variation_strength"`
				Variations        []interface{} `json:"variations"`
				Version           string        `json:"version"`
				Wh                []int         `json:"wh"`
			} `json:"model"`
			Task string `json:"task"`
		} `json:"request"`
		Safe      bool    `json:"safe"`
		StartTime float64 `json:"start_time"`
		UID       string  `json:"uid"`
	} `json:"data"`
	Pending int    `json:"pending"`
	Status  string `json:"status"`
}

func Status(authVal string, uid string) (body []byte, err error) {
	url := "https://open.nolibox.com/prod-open-aigc/engine/status/" + uid
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	req.Header.Add("User-Agent", "Apifox/1.0.0 (https://www.apifox.cn)")
	req.Header.Add("Authorization", fmt.Sprintf("Basic %s", authVal))

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer res.Body.Close()
	body, err = ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(body))
	return
}

func GetAuth(ak, sk string) string {
	authVal := fmt.Sprintf("%s:%s", ak, sk)
	return base64.StdEncoding.EncodeToString([]byte(authVal))
}

type DraftResp struct {
	Code int `json:"code"`
	Data struct {
		UID string `json:"uid"`
	} `json:"data"`
	Msg string `json:"msg"`
}

// https://creator-nolibox.apifox.cn/api-56157562
func Draft(authVal string, task, text, imgUrl string) (body []byte, err error) {
	url := "https://open.nolibox.com/prod-open-aigc/engine/push/draft"
	method := "POST"

	payload := strings.NewReader(fmt.Sprintf(`{
    "task": "%s",
    "params": {
        "text": "%s",
        "url": "%s",
        "w": 512,
        "h": 512,
        "fidelity": 0.4
    }
	}`, task, text, imgUrl))

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		fmt.Println(err)
		return
	}

	req.Header.Add("User-Agent", "Apifox/1.0.0 (https://www.apifox.cn)")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Basic %s", authVal))

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer res.Body.Close()
	body, err = ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(body))

	return
}
