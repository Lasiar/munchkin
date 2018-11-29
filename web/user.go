package web

import (
	"encoding/json"
	"fmt"
	"log"
	"munchkin/modal"
	"net/http"
)

type request struct {
	Password string `json:"password"`
	Login    string `json:"login"`
}

func UserRegistrator(w http.ResponseWriter, r *http.Request) {

	req := new(request)

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(req); err != nil {
		log.Printf("[Encode json] registration %v", err)
		return
	}

	if err := modal.New().SetUser(req.Login, req.Password); err != nil {
		log.Printf("[db] registration %v", err)
		return
	}
}

func UserAuthentications(w http.ResponseWriter, r *http.Request) {

	req := new(request)

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(req); err != nil {
		log.Printf("[Encode json] autentications %v", err)
	}

	valid, err := modal.New().Authentications(req.Login, req.Password)
	if err != nil {
		log.Printf("[DB]  Authentications %v", err)
	}

	fmt.Println(valid,err)

	encoder := json.NewEncoder(w)
	if err := encoder.Encode(struct {
		Aut bool `json:"auth"`
	}{valid}); err != nil {
		log.Printf("[web] encode %v", err)
	}
}
