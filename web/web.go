package web

import (
	"encoding/json"
	"log"
	"munchkin/system"
	"net/http"
	"time"
)

const (
	alreadyRegistered 	string = "Вы уже зарегестрированы"
	errorJsonRead       string = "Ошибка чтения запроса"
	errorCookie         string = "Ошибка чтения cookie"
	internalServerError string = "Внутренняя ошибка"
)

type webError struct {
	Error   error
	Message string
	Code    int
}

type webHandler func(http.ResponseWriter, *http.Request) *webError

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

func (wh webHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if e := wh(w, r); e != nil {
		encoder := json.NewEncoder(w)

		log.Printf("[WEB] %v %v %v", e.Code, e.Message, e.Error)

		w.WriteHeader(http.StatusInternalServerError)

		if err := encoder.Encode(struct {
			Message   string
			ErrorCode int
		}{e.Message, e.Code}); err != nil {
			log.Printf("[WEB] %v", err)
		}

	}
}

func Run() {
	apiMux := http.NewServeMux()

	apiMux.Handle("/api/registration", webHandler(UserRegistrator))
	apiMux.Handle("/api/authentications", webHandler(UserAuthentications))

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
