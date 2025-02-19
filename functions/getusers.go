package functions

import (
	"net/http"
	"sync"
	"example.com/models"
	"github.com/gin-gonic/gin"
	"fmt"
)

const limit = 5000
const targetUsers = 15000

func GetUsersHandler(c *gin.Context) {
  var ResponseResult models.Response
  var client = models.ClientRequest{
    UsersMap: make(map[string]models.CleanUser, targetUsers),
    Wg: &sync.WaitGroup{},
    Mu: &sync.Mutex{},
    Err: nil,
  }
  fmt.Println("Fetching users")

  FetchUsers(&client)

  if client.Err != nil {
    c.JSON(http.StatusInternalServerError, gin.H{"error": client.Err.Error()})
    return
  }
  for _, user := range client.UsersMap {
    ResponseResult.Users = append(ResponseResult.Users, user)
    //ResponseResult.UsersPerCountry[user.Country]++
    ResponseResult.AverageAge += float64(user.Age)
    if user.Gender == "female" {
      ResponseResult.WomenUsers++ 
    } else {
      ResponseResult.MenUsers++
    }
  }
  ResponseResult.AverageAge /= float64(len(client.UsersMap))
  c.JSON(http.StatusOK, ResponseResult)
}

func FetchUsers(client *models.ClientRequest) {
  for i := 0; i < targetUsers; i += limit {
    fmt.Println("Dentro")
    client.Wg.Add(1)
    if targetUsers - len(client.UsersMap) < limit {
      go generate(client, targetUsers - len(client.UsersMap))
      continue
    }
    go generate(client, limit)
  }
  client.Wg.Wait()
}

func generate(client *models.ClientRequest, limit int)  {
  client.Mu.Lock()
  defer client.Wg.Done()
  cant := 0
  for cant < limit {
    var ResponseAPI models.UsersAPi
    var CleanUsers []models.CleanUser
    usersAPI, err := ResponseAPI.GetUsers()
    if err != nil {
      client.Err = err
      return 
    }
    for _, user := range usersAPI {
      CleanUsers = append(CleanUsers, user.GetCleanUser())
    }
    for _, user := range CleanUsers {
      if _, exists := client.UsersMap[user.Uuid]; exists {
	continue
      }
      client.UsersMap[user.Uuid] = user
      cant++
    }
  }
  client.Mu.Unlock()
}
