package main

import (
	"fmt"
)

import "github.com/go-redis/redis"
func main() {
	ExampleNewClient()
	ExampleClient()
}

var client *redis.Client
func ExampleNewClient() {
	client = redis.NewClient(&redis.Options{
		Addr:     "47.98.218.146:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	pong, err := client.Ping().Result()
	fmt.Println(pong, err)
	// Output: PONG <nil>
}

func ExampleClient() {
	//err := client.Set("key", "value", 0).Err()
	//if err != nil {
	//	panic(err)
	//}
	//
	//err = client.Set("key3", "value", time.Duration(time.Second*20)).Err()
	//if err != nil {
	//	panic(err)
	//}
	//
	//val, err := client.Get("key").Result()
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println("key", val)

	val2, err := client.Get("key3").Result()
	client.Del()
	if err == redis.Nil {
		fmt.Println("key2 does not exist")
	} else if err != nil {
		panic(err)
	} else {
		fmt.Println("key2", val2)
	}
	// Output: key value
	// key2 does not exist
}