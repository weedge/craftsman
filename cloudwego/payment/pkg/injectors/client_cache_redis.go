package injectors

import (
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/go-redsync/redsync/v4"
	redsyncredis "github.com/go-redsync/redsync/v4/redis"
	"github.com/go-redsync/redsync/v4/redis/goredis/v8"
)

type RedisClusterClientOptions struct {
	Addrs []string `mapstructure:"addrs"`

	RedisClientCommonOptions

	// To route commands by latency or randomly, enable one of the following.
	Route string `mapstructure:"route"`
}

type RedisClientCommonOptions struct {
	Password string `mapstructure:"password"`
	Username string `mapstructure:"username"`

	MaxRetries      int           `mapstructure:"maxRetries"`
	MinRetryBackoff time.Duration `mapstructure:"minRetryBackoff"`
	MaxRetryBackoff time.Duration `mapstructure:"maxRetryBackoff"`
	DialTimeout     time.Duration `mapstructure:"dialTimeout"`
	ReadTimeout     time.Duration `mapstructure:"readTimeout"`
	WriteTimeout    time.Duration `mapstructure:"writeTimeout"`

	// connect pool
	PoolSize           int           `mapstructure:"poolSize"`
	MinIdleConns       int           `mapstructure:"minIdleConns"`
	MaxConnAge         time.Duration `mapstructure:"maxConnAge"`
	PoolTimeout        time.Duration `mapstructure:"poolTimeout"`
	IdleTimeout        time.Duration `mapstructure:"idleTimeout"`
	IdleCheckFrequency time.Duration `mapstructure:"idleCheckFrequency"`
}

func DefaultRedisClusterClientOptions() *RedisClusterClientOptions {
	return &RedisClusterClientOptions{
		//Addrs:    []string{":26379"},
		Addrs: []string{":26379", ":26380", ":26381", ":26382", ":26383", ":26384"},
		RedisClientCommonOptions: RedisClientCommonOptions{
			Password: "",
			Username: "",

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
		},
		Route: "randomly",
	}
}

// InitRedisClusterDefaultClient init default redis cluster instance
func InitRedisClusterDefaultClient(opts *RedisClusterClientOptions) redis.UniversalClient {
	clusterOpts := &redis.ClusterOptions{
		Addrs:    opts.Addrs,
		Password: opts.Password,
		Username: opts.Username,

		MaxRetries:      opts.MaxRetries,
		MinRetryBackoff: opts.MinRetryBackoff,
		MaxRetryBackoff: opts.MaxRetryBackoff,
		DialTimeout:     opts.DialTimeout,
		ReadTimeout:     opts.ReadTimeout,
		WriteTimeout:    opts.WriteTimeout,

		// connect pool
		PoolSize:           opts.PoolSize,
		MinIdleConns:       opts.MinIdleConns,
		MaxConnAge:         opts.MaxConnAge,
		PoolTimeout:        opts.PoolTimeout,
		IdleTimeout:        opts.IdleTimeout,
		IdleCheckFrequency: opts.IdleCheckFrequency,

		// To route commands by latency or randomly, enable one of the following.
		//RouteByLatency: true,
		//RouteRandomly: true,
	}
	switch opts.Route {
	case "randomly":
		clusterOpts.RouteRandomly = true
	case "latency":
		clusterOpts.RouteByLatency = true
	default:
		clusterOpts.RouteRandomly = true
	}

	return redis.NewClusterClient(clusterOpts)
}

type RedisClientOptions struct {
	Addr string `mapstructure:"addr"`
	Db   int    `mapstructure:"db"`

	RedisClientCommonOptions
}

func DefaultRedisClientOptions() *RedisClientOptions {
	return &RedisClientOptions{
		Addr: "localhost:6379",
		RedisClientCommonOptions: RedisClientCommonOptions{
			Password: "",
			Username: "",

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
		},
	}
}

// InitRedisDefaultClient init default redis instance
func InitRedisDefaultClient(opts *RedisClientOptions, limiter redis.Limiter) redis.UniversalClient {
	return redis.NewClient(&redis.Options{
		Addr:               opts.Addr,
		Username:           opts.Username,
		Password:           opts.Password,
		DB:                 opts.Db,
		MaxRetries:         opts.MaxRetries,
		MinRetryBackoff:    opts.MinRetryBackoff,
		MaxRetryBackoff:    opts.MaxRetryBackoff,
		DialTimeout:        opts.DialTimeout,
		ReadTimeout:        opts.ReadTimeout,
		WriteTimeout:       opts.WriteTimeout,
		PoolSize:           opts.PoolSize,
		MinIdleConns:       opts.MinIdleConns,
		MaxConnAge:         opts.MaxConnAge,
		PoolTimeout:        opts.PoolTimeout,
		IdleTimeout:        opts.IdleTimeout,
		IdleCheckFrequency: opts.IdleCheckFrequency,
		//TLSConfig:          &tls.Config{},
		Limiter: limiter,
	})
}

// InitRedsync Create an instance of redisync to be used to obtain a mutual exclusion for distlock
// ref: https://redis.io/topics/distlock/
func InitRedsync(clients ...redis.UniversalClient) *redsync.Redsync {
	var pools []redsyncredis.Pool
	for _, client := range clients {
		pools = append(pools, goredis.NewPool(client))
	}

	return redsync.New(pools...)
}
