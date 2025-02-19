package models

import (
	"sync"
	"time"
)

type ClientRequest struct {
  UsersMap map[string]bool 
  Wg *sync.WaitGroup
  Mu *sync.Mutex
  Time time.Time
  Err error
}
