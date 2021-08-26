package main

import (
	"context"
	redisCache "github.com/go-redis/cache/v8"
	"gogs.iotopo.com/iotopo/iotopo-sdk-go/cache"
	"log"
	"time"
)

type Item struct {
	Foo string
}

type Object struct {
	Str string
	Num int
}

func main() {
	defer cache.Stop()
	//redis := cache.GetConn()
	//cache.Set("hello", )

	cache.SetEx("key1", "value1", 10)
	value1, _ := cache.Get("key1")
	log.Println("key1=", value1)

	// 对象缓存
	mycache := cache.NewCache(&cache.CacheOption{
		LocalCache: true, // 开启本地缓存
	})

	key := "mykey"
	obj := &Object{
		Str: "mystring",
		Num: 42,
	}

	if err := mycache.Set(&redisCache.Item{
		Ctx:   context.Background(),
		Key:   key,
		Value: obj,
		TTL:   time.Hour,
	}); err != nil {
		panic(err)
	}

	var wanted Object
	if err := mycache.Get(context.Background(), key, &wanted); err == nil {
		log.Println(wanted)
	}
	// Output: {mystring 42}


}
