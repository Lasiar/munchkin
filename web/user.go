package web

import (
	"encoding/json"
	"fmt"
	"log"
	"munchkin/modal"
	"net/http"
	"time"
)

type request struct {
	Password string `json:"password"`
	Login    string `json:"login"`
}

func UserRegistrator(w http.ResponseWriter, r *http.Request) *webError {
	req := new(request)

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(req); err != nil {
		return &webError{fmt.Errorf("[web] decode json %v", err), errorJsonRead, 001}
	}

	if err := modal.New().SetUser(req.Login, req.Password); err != nil {
		return &webError{fmt.Errorf("[db] %v", err), internalServerError, 201}
	}
	return nil
}

func UserAuthentications(w http.ResponseWriter, r *http.Request) *webError {
	req := new(request)

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(req); err != nil {
		return &webError{fmt.Errorf("[web] decode json %v", err), errorJsonRead, 001}
	}

	hashCookie, auth, err := modal.New().Authentications(req.Login, req.Password)
	if err != nil {
		return &webError{fmt.Errorf("[db] %v", err), internalServerError, 201}
	}

	if !auth {
		return &webError{fmt.Errorf("[WEB] error auth %v", err), internalServerError, 403}
	}

	cookie := http.Cookie{Name: "userName", Value: hashCookie, Expires: time.Now().Add(365 * 24 * time.Hour)}

	http.SetCookie(w, &cookie)

	encoder := json.NewEncoder(w)
	if err := encoder.Encode(struct {
		Aut bool `json:"auth"`
	}{false}); err != nil {
		log.Printf("[web] encode %v", err)
	}
	return nil
}
