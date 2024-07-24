package middleware

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pquerna/otp/totp"
	"github.com/skip2/go-qrcode"
)

func GenerateMFASecretHandler(w http.ResponseWriter, r *http.Request) {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "auth-app",
		AccountName: "pierre.caboor59@gmail.com",
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Générer le QR code
	qrCode, err := qrcode.Encode(key.URL(), qrcode.Medium, 256)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Écrire le QR code dans le corps de la réponse
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"secret": key.Secret(),
		"url":    key.URL(),
		"qrcode": "data:image/png;base64," + base64.StdEncoding.EncodeToString(qrCode),
	})
}

func VerifyMFAHandler(w http.ResponseWriter, r *http.Request) {
	var requestBody map[string]string
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	secret := requestBody["secret"]
	otp := requestBody["otp"]

	if totp.Validate(otp, secret) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "success"})
		fmt.Println("clé validé avec succès")
	} else {
		http.Error(w, "Invalid OTP", http.StatusUnauthorized)
	}

}
