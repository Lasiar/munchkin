package web

import (
	"log"
	"munchkin/system"
	"net/http"
	"time"
)

func MiddlewareCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Header().Add("Access-Control-Allow-Headers", "*")
		w.Header().Add("Access-Control-Allow-Credentials", "true")
		w.Header().Add("Access-Control-Allow-Methods", "POST, OPTIONS")

		if r.Method == http.MethodOptions {
			return
		}

		next.ServeHTTP(w, r)
	})
}

func Run() {
	apiMux := http.NewServeMux()

	apiMux.HandleFunc("/api/registration", UserRegistrator)
	apiMux.HandleFunc("/api/authentications", UserAuthentications)

	webServer := &http.Server{
		Addr:           system.GetConfig().Port(),
		Handler:        MiddlewareCORS(apiMux),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	if err := webServer.ListenAndServe(); err != nil {
		log.Fatal(err)
	}

}
