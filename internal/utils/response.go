package utils

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Status  string      `json:"status"`
	Message interface{} `json:"message"`
}

func Error(w http.ResponseWriter, error Response, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	error.Status = "ERR"
	json.NewEncoder(w).Encode(error)
}

func Success(w http.ResponseWriter, success Response, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	success.Status = "OK"
	json.NewEncoder(w).Encode(success)
}
