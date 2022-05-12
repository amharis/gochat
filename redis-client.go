package main

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

func ExampleClient() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	err := rdb.Set(ctx, "key", "value", 0).Err()
	if err != nil {
		panic(err)
	}

	val, err := rdb.Get(ctx, "key").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("key", val)

	val2, err := rdb.Get(ctx, "key2").Result()
	if err == redis.Nil {
		fmt.Println("key2 does not exist")
	} else if err != nil {
		panic(err)
	} else {
		fmt.Println("key2", val2)
	}
	// Output: key value
	// key2 does not exist
	rdb.PubSubChannels()
}

/*
func main() {
	//init:connects to the default port 6379
	var client, client2 redis.Client

	//strings:set get del
	client.Set("a", []byte("hello"))
	client2.Set("b", []byte("world"))
	word, err := client.Get("a")
	if err == nil {
		fmt.Println("client get ", string(word))
	} else {
		fmt.Println("client get error ", err)
	}

	//lists
	data := []string{"a", "b", "c", "d", "e"}
	for _, v := range data {
		client.Rpush("list", []byte(v))
	}
	ret, err := client.Lrange("list", 0, -1)
	if err == nil {
		for index, val := range ret {
			fmt.Println(index, ":", string(val))
		}
	}
	client.Del("list")

	//pub/sub
	sub := make(chan string, 1)
	sub <- "channel"

	messages := make(chan redis.Message, 0)
	go client.Subscribe(sub, nil, nil, nil, messages)

	time.Sleep(10 * 1000 * 1000)
	client2.Publish("channel", []byte("cool"))
	msg := <-messages
	fmt.Println("received from:", msg.Channel, " message:", string(msg.Message))

	close(sub)
	close(messages)
}
*/
