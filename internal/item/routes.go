package item

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/sedat/redis-api/internal/item/model"
	"github.com/sedat/redis-api/internal/item/repository"
	"github.com/sedat/redis-api/internal/utils"
)

func GetItemHandler(ctx context.Context, redisRepository repository.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		key := req.URL.Query().Get("key")
		value, err := redisRepository.Get(ctx, key)
		if err != nil {
			utils.Error(w, utils.Response{Message: err}, http.StatusBadGateway)
			return
		}
		item := model.Item{Key: key, Value: value}
		utils.Success(w, utils.Response{Message: item}, http.StatusAccepted)
	}
}

func GetAllKeysHandler(ctx context.Context, redisRepository repository.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		value, err := redisRepository.GetAllKeys(ctx)
		if err != nil {
			utils.Error(w, utils.Response{Message: err}, http.StatusBadGateway)
			return
		}
		utils.Success(w, utils.Response{Message: value}, http.StatusAccepted)
	}
}

func GetAllItemsHandler(ctx context.Context, redisRepository repository.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		value, err := redisRepository.GetAllItems(ctx)
		if err != nil {
			utils.Error(w, utils.Response{Message: err}, http.StatusBadGateway)
			return
		}
		utils.Success(w, utils.Response{Message: value}, http.StatusAccepted)
	}
}

func SetItemHandler(ctx context.Context, redisRepository repository.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		var item model.Item
		err := json.NewDecoder(req.Body).Decode(&item)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			utils.Error(w, utils.Response{Message: err}, http.StatusBadRequest)
			return
		}
		err = redisRepository.Set(ctx, item)
		if err != nil {
			w.WriteHeader(http.StatusBadGateway)
			utils.Error(w, utils.Response{Message: err}, http.StatusBadGateway)
			return
		}
		utils.Success(w, utils.Response{Message: item}, http.StatusCreated)
	}
}

func FlushDBHandler(ctx context.Context, redisRepository repository.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		err := redisRepository.Flush(ctx)
		if !err {
			w.WriteHeader(http.StatusBadGateway)
			utils.Error(w, utils.Response{Message: err}, http.StatusBadGateway)
			return
		}
		utils.Success(w, utils.Response{Message: "Redis flushed!"}, http.StatusCreated)
	}
}
