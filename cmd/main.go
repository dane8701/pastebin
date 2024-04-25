package main

import (
	"context"
	"flag"
	"fmt"
	"pastebin/domain"
	"pastebin/store"
)

func main() {
	redis := flag.String("redis", "localhost:6379", "redis parameter")
	flag.Parse()

	svc, err := store.NewRedisDB(context.Background(), *redis)
	if err != nil {
		fmt.Printf("[error redis]: %v", err)
		return
	}

	err = domain.ServeAPI(svc)()
	if err != nil {
		fmt.Printf("[error server]: %v", err)
		return
	}
}
