package models

type User struct {
  Gender string `json:"gender"`
  Name struct {
    First string `json:"first"`
    Last string `json:"last"`
  } `json:"name"`
  Country string `json:"country"`
  Email string `json:"email"`
  Login struct {
    Uuid string `json:"uuid"`
  } `json:"login"`
  Dob struct {
    Age int `json:"age"`
  } `json:"dob"`
}

func(u *User) GetCleanUser() CleanUser {
  return CleanUser{
    Gender 	: u.Gender,
    FirstName 	: u.Name.First,
    SecondName 	: u.Name.Last,
    Country 	: u.Country,
    Email 	: u.Email,
    Uuid 	: u.Login.Uuid,
  }
}
