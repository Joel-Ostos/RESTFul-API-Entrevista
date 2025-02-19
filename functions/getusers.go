package functions

import (
  "time"
  "net/http"
  "sync"
  "example.com/models"
  "github.com/gin-gonic/gin"
  "fmt"
)

const ( 
  limit = 5000
  targetUsers = 15000
)

func GetUsersHandler(c *gin.Context) {
  var ResponseResult models.Response
  client := models.ClientRequest{
    UsersMap: make(map[string]bool, targetUsers),
    Wg: &sync.WaitGroup{},
    Mu: &sync.Mutex{},
    Time: time.Now(),
    Err: nil,
  }
  fmt.Println("Fetching users")

  FetchUsers(&client, &ResponseResult)

  if client.Err != nil {
    c.JSON(http.StatusInternalServerError, gin.H{"error": client.Err.Error()})
    return
  }

  ResponseResult.UsersPerCountry = make(map[string]int, 30)
  for _, user := range ResponseResult.Users {
    ResponseResult.UsersPerCountry[user.Country]++
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

func FetchUsers(client *models.ClientRequest, ResponseResult *models.Response) {
  for i := 0; i < targetUsers; i += limit {
    if client.Err != nil {
      return
    }
    client.Wg.Add(1)
    if targetUsers - len(client.UsersMap) < limit {
      go generate(client, ResponseResult,targetUsers - len(client.UsersMap))
      continue
    }
    go generate(client, ResponseResult, limit)
  }
  client.Wg.Wait()
}

func generate(client *models.ClientRequest,  ResponseResult *models.Response, limit int)  {
  defer client.Wg.Done()
  cant := 0
  for cant < limit {
    var ResponseAPI models.UsersAPi
    usersAPI, err := ResponseAPI.GetUsers()
    if err != nil {
      client.Err = err
      return 
    }
    for _, user := range usersAPI  {
      client.Mu.Lock()
      if _, exists := client.UsersMap[user.Login.Uuid]; exists {
	continue
      }
      ResponseResult.Users = append(ResponseResult.Users, user.GetCleanUser())
      client.UsersMap[user.Login.Uuid] = true
      client.Mu.Unlock()
      cant++
    }
  }
}
