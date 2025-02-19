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
	"sync"

	"github.com/gin-gonic/gin"
)

type User struct {
	Gender  string `json:"gender"`
	Name    struct {
		First string `json:"first"`
		Last  string `json:"last"`
	} `json:"name"`
	Email   string `json:"email"`
	Country string `json:"country"`
	City    string `json:"city"`
	Login   struct {
		Uuid string `json:"uuid"`
	} `json:"login"`
}

type ResultsAPI struct {
	Results []User `json:"results"`
}

type Response struct {
	Users           []User         `json:"users"`
	UsersPerCountry map[string]int `json:"users_per_country"`
	Men             int            `json:"men"`
	Women           int            `json:"women"`
}

const url string = "https://randomuser.me/api/?results=5000"
const totalUsers int = 15000
const numGoroutines int = 3

func fetchUsers(result *Response, users map[string]User, wg *sync.WaitGroup, mu *sync.Mutex) {
	defer wg.Done()
	for len(result.Users) < totalUsers/numGoroutines {
		resp, err := http.Get(url)
		if err != nil {
		  log.Printf("Error en la solicitud HTTP: %v", err)
		  return
		}
		var tempUsers ResultsAPI
		if err := json.NewDecoder(resp.Body).Decode(&tempUsers); err != nil {
		  log.Printf("Error decodificando la respuesta: %v", err)
		  continue
		}
		fillMap(result, &tempUsers.Results, users, mu)
	}
}

func fillMap(result *Response, tempUsers *[]User, users map[string]User, mu *sync.Mutex) {
	mu.Lock()
	defer mu.Unlock()
	for _, user := range *tempUsers {
		if _, exists := users[user.Login.Uuid]; !exists {
			users[user.Login.Uuid] = user
			result.UsersPerCountry[user.Country]++
			if user.Gender == "female" {
				result.Women++
			} else {
				result.Men++
			}
			result.Users = append(result.Users, user)
		}
	}
}

func usersHandler(g *gin.Context) {
	var result Response
	result.UsersPerCountry = make(map[string]int)

	users := make(map[string]User)
	var wg sync.WaitGroup
	var mu sync.Mutex

	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go fetchUsers(&result, users, &wg, &mu)
	}
	wg.Wait()

	g.JSON(http.StatusOK, result)
}

func main() {
	r := gin.Default()
	r.GET("/users", usersHandler)
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
