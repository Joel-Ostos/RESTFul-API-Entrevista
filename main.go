/* Se debe desarrollar una API REST que consuma la API pública de RandomUser( https://randomuser.me/documentation ) para generar una lista de 15,000 usuarios únicos.

Requisitos
género, primer nombre, primer apellido,
Cada usuario debe incluir:  email, ciudad, país y UUID.
El UUID no debe repetirse.
La API debe responder en menos de 5.25 segundos.

Adicional se debe proporcionar un dato estadistico:
Distribución por género (cantidad de hombres y mujeres).
Cantidad de usuarios por país.
Edad promedio*/
package main

import (
	"encoding/json"
	"log"
	"net/http"
	"github.com/gin-gonic/gin"
)

type User struct {
  Gender 		string `json:"gender"`
  Name struct  {
    First 		string `json:"first"`
    Last 		string `json:"last"`
  } 			       `json:"name"`
  Email 		string `json:"email"`
  Country 		string `json:"country"`
  City 			string `json:"city"`
  Login struct {
    Uuid 		string `json:"uuid"`
  } 			       `json:"login"`
}

type ResultsAPI struct {
  Results []User `json:"results"`
}

type Response struct {
  Users 		[]User
  UsersPerCountry 	map[string]int
  Men 			int
  Women 		int
}

const url string = "https://randomuser.me/api/?results=5000"

func fetchUsers(result *Response, Users *map[string]User) {
  cant := 0
  for cant < 15000 {
    resp, err := http.Get(url)
    var tempUsers ResultsAPI
    if err != nil {
      return 
    }
    json.NewDecoder(resp.Body).Decode(&tempUsers)
    fillMap(&cant, result, &tempUsers.Results, Users)
  }
}

func fillMap(cant *int, result *Response, tempUsers *[]User, Users *map[string]User) {
  for i := 0; i < len((*tempUsers)); i++ {
    _, exists := (*Users)[(*tempUsers)[i].Login.Uuid] 
    if !exists {
      (*cant)++
      (*Users)[(*tempUsers)[i].Login.Uuid] = (*tempUsers)[i]
      result.UsersPerCountry[(*tempUsers)[i].Country]++
      if (*tempUsers)[i].Gender == "female" {
	result.Women++
      } else {
	result.Men++
      }
      result.Users = append(result.Users, (*tempUsers)[i])
    }
  }
}

func usersHandler(g *gin.Context) {
  var result Response
  result.UsersPerCountry = make(map[string]int, 100)
  var Users  = make(map[string]User, 15000)
  fetchUsers(&result, &Users)
  g.JSON(http.StatusOK, result)
}

func main() {
  r := gin.Default()
  r.GET("/users", usersHandler)
  if err := r.Run(":8080"); err != nil {
    log.Fatal(err)
  }
}
