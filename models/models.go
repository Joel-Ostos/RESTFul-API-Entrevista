package models

import (
	"encoding/json"
	"errors"
	"net/http"
)

type UsersAPi struct {
  Results 	[]User 	`json:"results"`
  Error 	string 	`json:"error"`
}

const url = "https://randomuser.me/api/?results=5000"

func(u *UsersAPi) GetUsers() ([]User, error) {
  resp, err := http.Get(url)
  if err != nil {
    return u.Results, err
  }
  json.NewDecoder(resp.Body).Decode(u)
  if u.Error != "" {
    return u.Results, errors.New(u.Error)
  }
  return u.Results, nil
}
