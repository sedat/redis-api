package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-redis/redis/v8"
	"github.com/sedat/redis-api/internal/item"
	"github.com/sedat/redis-api/internal/item/model"
	"github.com/sedat/redis-api/internal/item/repository"
	"github.com/sedat/redis-api/internal/utils"
)

var ctx context.Context
var redisClient *redis.Client
var redisRepository repository.Repository

func init() {
	ctx = context.Background()
	redisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	redisRepository = repository.NewRepository(redisClient)
}

func TestRedisConnection(t *testing.T) {
	_, err := redisClient.Ping(ctx).Result()
	if err != nil {
		t.Errorf("Redis connection error")
	}
}

func TestSetItemHandler(t *testing.T) {
	var reqItem model.Item
	reqItem.Key = "test"
	reqItem.Value = "test"
	body, _ := json.Marshal(reqItem)
	req := httptest.NewRequest(http.MethodPost, "/set", bytes.NewReader(body))
	defer req.Body.Close()
	w := httptest.NewRecorder()
	handler := item.SetItemHandler(ctx, redisRepository)
	handler(w, req)
	res := w.Result()
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Errorf("Expected error to be nil but got %v", err)
	}
	var response utils.Response
	err = json.Unmarshal(data, &response)
	if err != nil {
		t.Errorf("Cannot conver data to json %v", err)
	}
	if response.Status != "OK" {
		t.Errorf("Status is not OK")
	}
	t.Logf("%v", response)
}

func TestGetItemHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/get?key=test", nil)
	w := httptest.NewRecorder()
	handler := item.GetItemHandler(ctx, redisRepository)
	handler(w, req)
	res := w.Result()
	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Errorf("Expected error to be nil but got %v", err)
	}
	var response utils.Response
	err = json.Unmarshal(data, &response)
	if err != nil {
		t.Errorf("Cannot convert data to json %v", err)
	}
	if response.Status != "OK" {
		t.Errorf("Status is not OK %v", response)
	}
	t.Logf("%v", response)
}

func TestFlushDBHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/flush", nil)
	w := httptest.NewRecorder()
	handler := item.FlushDBHandler(ctx, redisRepository)
	handler(w, req)
	res := w.Result()
	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Errorf("Expected error to be nil but got %v", err)
	}
	var response utils.Response
	err = json.Unmarshal(data, &response)
	if err != nil {
		t.Errorf("Cannot convert data to json %v", err)
	}
	if response.Status != "OK" {
		t.Errorf("Status is not OK %v", response)
	}
	t.Logf("%v", response)
}

func TestGetAllItemsHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/get-items", nil)
	w := httptest.NewRecorder()
	handler := item.GetAllItemsHandler(ctx, redisRepository)
	handler(w, req)
	res := w.Result()
	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Errorf("Expected error to be nil but got %v", err)
	}
	var response utils.Response
	err = json.Unmarshal(data, &response)
	if err != nil {
		t.Errorf("Cannot convert data to json %v", err)
	}
	if response.Status != "OK" {
		t.Errorf("Status is not OK %v", response)
	}
	t.Logf("%v", response)
}
