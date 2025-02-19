package models
import (
  "sync"
)

type ClientRequest struct {
  UsersMap map[string]CleanUser 
  Wg *sync.WaitGroup
  Mu *sync.Mutex
  Err error
}
