package utils

import (
	"encoding/json"
	"log"
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
	err := json.NewEncoder(w).Encode(error)
	if err != nil {
		log.Printf("Error encoding json %s", err)
	}
}

func Success(w http.ResponseWriter, success Response, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	success.Status = "OK"
	err := json.NewEncoder(w).Encode(success)
	if err != nil {
		log.Printf("Error encoding json %s", err)
	}
}
