package server

import (
	"handlers"
	"log"
	"middleware"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func RunServer() {
	r := mux.NewRouter()

	// mfa
	r.HandleFunc("/", handlers.About)
	r.HandleFunc("/mfa/verify", middleware.VerifyMFAHandler).Methods("POST")
	r.HandleFunc("/generate-mfa-secret", middleware.GenerateMFASecretHandler).Methods("GET")

	//r.HandleFunc("/setcookie", setCookieHandler)
	//r.HandleFunc("/getcookie", getCookieHandler)

	c := cors.New(cors.Options{
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	handler := c.Handler(r)
	server := &http.Server{
		Addr:              ":8080",
		Handler:           handler,
		ReadHeaderTimeout: 10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       10 * time.Second,
		ReadTimeout:       20 * time.Second,
		MaxHeaderBytes:    1 << 20,
	}

	log.Printf("Server starting on http://localhost%s...\n", server.Addr)

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

//func setCookieHandler(w http.ResponseWriter, req *http.Request) {
//	//expire := time.Now().Add(12 * time.Hour)
//	cookie := middleware.JSONCookie{
//		Name:  "auth",
//		Value: "authenticated",
//		Path:  "/",
//		//Expires:  expire,
//		MaxAge:   12 * 60 * 60,
//		HttpOnly: true,
//	}
//
//	// Écrire la réponse JSON avec le cookie dans le corps
//	w.Header().Set("Content-Type", "application/json")
//	w.WriteHeader(http.StatusOK)
//	w.Write(cookie.StrucToJSON())
//}
//
//func getCookieHandler(w http.ResponseWriter, req *http.Request) {
//	var data map[string]string
//	err := json.NewDecoder(req.Body).Decode(&data)
//	if err != nil {
//		w.WriteHeader(http.StatusBadRequest)
//		json.NewEncoder(w).Encode(map[string]string{"message": "Invalid cookie format"})
//		return
//	}
//
//	cookie, err := middleware.JSONToStruct(data)
//
//	if err != nil {
//		fmt.Println("Cookie invalide")
//	}
//
//	if cookie.Value == "authenticated" {
//		fmt.Println("authentifié")
//	} else {
//		fmt.Println("non authentifié")
//	}
//
//	//fmt.Println(string(cookie.StrucToJSON()))
//}
