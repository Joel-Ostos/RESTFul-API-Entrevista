package models

import (
	"encoding/json"
	"net/http"
)

type UsersAPi struct {
  Results []User `json:"results"`
}
const url = "https://randomuser.me/api/?results=5000"

func(u *UsersAPi) GetUsers() ([]User, error) {
  resp, err := http.Get(url)
  if err != nil {
    return u.Results, err
  }
  json.NewDecoder(resp.Body).Decode(u)
  return u.Results, nil
}

