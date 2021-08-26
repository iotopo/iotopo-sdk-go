package cache

//
// TODO: rate limiting based on Redis. https://github.com/go-redis/redis_rate
// TODO: Simplified distributed locking implementation using Redis. https://github.com/bsm/redislock
import (
	"context"
	"github.com/go-redis/cache/v8"
	"github.com/go-redis/redis/v8"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

var conn redis.Cmdable
var client redis.UniversalClient
var ctx = context.Background()

func init() {
	addr := os.Getenv("REDIS_ADDR")
	if addr == "" {
		addr = "127.0.0.1:6379"
	}

	username := os.Getenv("REDIS_USERNAME")
	password := os.Getenv("REDIS_PASSWORD")

	db := 0
	if val := os.Getenv("REDIS_DB"); val != "" {
		i, err := strconv.Atoi(val)
		if err != nil {
			log.Fatalf("illegal envvar REDIS_DB: %s", val)
		}
		db = i
	}

	var minIdleConns int
	if val := os.Getenv("REDIS_MIN_IDLE_CONNS"); val != "" {
		i, err := strconv.Atoi(val)
		if err != nil {
			log.Fatalf("illegal envvar REDIS_MIN_IDLE_CONNS: %s", val)
		}
		if i > 0 {
			minIdleConns = i
		}
	}
	var poolSize int
	if val := os.Getenv("REDIS_POOL_SIZE"); val != "" {
		v, _ := strconv.Atoi(val)
		if v >= 0 {
			poolSize = v
		}
	}

	var poolTimeout int
	if val := os.Getenv("REDIS_POOL_TIMEOUT"); val != "" {
		v, _ := strconv.Atoi(val)
		if v >= 0 {
			poolTimeout = v
		}
	}

	var idleTimeout int
	if val := os.Getenv("REDIS_IDLE_TIMEOUT"); val != "" {
		v, _ := strconv.Atoi(val)
		if v >= 0 {
			idleTimeout = v
		}
	}

	//idleTimeout := time.Duration(idleTimeout) * time.Second
	//connectTimeout := time.Duration(connectTimeout) * time.Second

	adds := strings.Split(addr, ",")
	if len(adds) > 1 { //	集群模式
		client = redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:         adds,
			ReadOnly:      true,
			RouteRandomly: true,
			PoolSize:      poolSize,
			MinIdleConns:  minIdleConns,
			Username:      username,
			Password:      password,
			PoolTimeout:   time.Duration(poolTimeout) * time.Second,
			IdleTimeout:   time.Duration(idleTimeout) * time.Second,
		})
	} else {
		client = redis.NewClient(&redis.Options{
			Addr:         addr,
			Username:     username,
			Password:     password,
			DB:           db,
			PoolSize:     poolSize,
			MinIdleConns: minIdleConns,
			PoolTimeout:  time.Duration(poolTimeout) * time.Second,
			IdleTimeout:  time.Duration(idleTimeout) * time.Second,
		})
	}
	conn = client
	//_, err := conn.Ping(ctx).Result()
	//if err != nil {
	//	log.Fatal(err)
	//}
}

func GetConn() redis.Cmdable {
	return conn
}

func GetClient() redis.UniversalClient {
	return client
}

func Stop() {
	if client != nil {
		client.Close()
	}
}

func Del(key string) error {
	return conn.Del(ctx, key).Err()
}

func Set(key string, value interface{}) error {
	return conn.Set(ctx, key, value, 0).Err()
}

func SetEx(key string, value interface{}, expiration uint) error {
	return conn.SetEX(ctx, key, value, time.Duration(expiration)*time.Second).Err()
}

func Get(key string) (string, error) {
	return conn.Get(ctx, key).Result()
}

func HGet(key, childKey string) (string, error) {
	return conn.HGet(ctx, key, childKey).Result()
}

func HGetAll(key string) (map[string]string, error) {
	return conn.HGetAll(ctx, key).Result()
}

func HMGetAll(key string, fields ...string) ([]string, error) {
	result, err := conn.HMGet(ctx, key, fields...).Result()
	if err != nil {
		return nil, err
	}
	var out []string
	for i := range result {
		out = append(out, result[i].(string))
	}
	return out, nil
}

func HSet(key string, value map[string]interface{}) error {
	return HSetEx(key, value, 0)
}

func HSetEx(key string, value map[string]interface{}, seconds uint) error {
	var data []interface{}
	for k, v := range value {
		data = append(data, k)
		data = append(data, v)
	}
	pipe := conn.Pipeline()
	pipe.HSet(ctx, key, data)
	if seconds > 0 {
		pipe.Expire(ctx, key, time.Duration(seconds)*time.Second)
	}
	_, err := pipe.Exec(ctx)
	return err
}

type CacheOption struct {
	StatsEnabled   bool
	Marshal        cache.MarshalFunc
	Unmarshal      cache.UnmarshalFunc
	LocalCache     bool
	LocalCacheSize int
	LocalCacheTTL  time.Duration
}

func NewCache(opt *CacheOption) *cache.Cache {
	option := &cache.Options{
		Redis:        conn,
		StatsEnabled: opt.StatsEnabled,
		Marshal:      opt.Marshal,
		Unmarshal:    opt.Unmarshal,
	}
	if opt.LocalCache {
		ttl := opt.LocalCacheTTL
		if ttl == 0 {
			ttl = time.Minute
		}
		size := opt.LocalCacheSize
		if size <= 0 {
			size = 1000
		}
		option.LocalCache = cache.NewTinyLFU(size, ttl)
	}
	return cache.New(option)
}
