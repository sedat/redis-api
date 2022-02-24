package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/sedat/redis-api/internal/item"
	"github.com/sedat/redis-api/internal/item/repository"
)

var ctx = context.Background()

func main() {
	redisAddress := fmt.Sprintf("%s:6379", os.Getenv("REDIS_URL"))
	client := redis.NewClient(&redis.Options{
		Addr:     redisAddress,
		Password: "",
		DB:       0,
	})
	_, err := client.Ping(ctx).Result()
	if err != nil {
		log.Fatal(err)
	}

	redisRepository := repository.NewRepository(client)

	http.HandleFunc("/get", item.GetItemHandler(ctx, redisRepository))
	http.HandleFunc("/get-keys", item.GetAllKeysHandler(ctx, redisRepository))
	http.HandleFunc("/get-items", item.GetAllItemsHandler(ctx, redisRepository))
	http.HandleFunc("/set", item.SetItemHandler(ctx, redisRepository))
	http.HandleFunc("/flush", item.FlushDBHandler(ctx, redisRepository))
	log.Println("Listing for requests at http://localhost:9000/")
	err = http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("PORT")), logRequest(http.DefaultServeMux))
	if err != nil {
		log.Fatal(err)
	}
	for range time.Tick(time.Minute * 1) {

	}
}

func logRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}
