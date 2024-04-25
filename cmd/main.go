package main

import (
	"context"
	"flag"
	"fmt"
	"pastebin/domain"
	"pastebin/store"
	"math/rand"
	"time"
)

func generateSecretKey() []byte {
	rand.Seed(time.Now().UnixNano())
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890!@#$%^&*()-_=+"
	b := make([]byte, 32)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return b
}

func main() {
	redis := flag.String("redis", "localhost:6379", "redis parameter")
	flag.Parse()

	svc, err := store.NewRedisDB(context.Background(), *redis)
	if err != nil {
		fmt.Printf("[error redis]: %v", err)
		return
	}

	secretKey := generateSecretKey()

	err = domain.ServeAPI(svc, secretKey)()
	if err != nil {
		fmt.Printf("[error server]: %v", err)
		return
	}
}
