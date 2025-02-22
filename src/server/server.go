package totpserver

import (
	"TOTP-telegram/src/data"
	totp "TOTP-telegram/src/totp-generator"
	"encoding/json"
	"net/http"
	"time"
)

type TOTPRequest struct {
	Name   string `json:"name"`
	Secret string `json:"secret"`
	UserID int64  `json:"user_id"`
}

func setupRoutes() {
	http.HandleFunc("/api/totp/generate", handleGenerateTOTP)
	http.HandleFunc("/api/totp/get", handleGetTOTP)
}

func StartServer(port string) error {
	setupRoutes()
	return http.ListenAndServe(":"+port, nil)
}

func handleGenerateTOTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req TOTPRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := data.SaveSecret(req.UserID, req.Name, req.Secret); err != nil {
		http.Error(w, "Failed to save secret", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"status": "Secret saved successfully",
	})
}

func handleGetTOTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req TOTPRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	secret, err := data.GetSecret(req.UserID, req.Name)
	if err != nil {
		http.Error(w, "Secret not found", http.StatusNotFound)
		return
	}

	otp := totp.GenerateTOTP(secret, time.Now().Unix())

	json.NewEncoder(w).Encode(map[string]interface{}{
		"otp": otp,
	})
}
