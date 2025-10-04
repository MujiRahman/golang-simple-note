package helper

import (
	"encoding/json"
	"net/http"
)

func RespondJSON(w http.ResponseWriter, code int, payload any) {
	if payload == nil {
		w.WriteHeader(code)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(payload)
}
