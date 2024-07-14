package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Request struct {
	UserID int `json:"userId"`
}

type Response struct {
	UserID int `json:"userId"`
}

func processHandler(w http.ResponseWriter, r *http.Request) {
	var req Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	processedUserID := req.UserID*2 + 1

	resp := Response{UserID: processedUserID}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func main() {
	http.HandleFunc("/api/process", processHandler)
	fmt.Println("Server is running on port 3000")
	log.Fatal(http.ListenAndServe(":3000", nil))
}
