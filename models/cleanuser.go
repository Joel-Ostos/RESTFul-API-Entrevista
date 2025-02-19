package models

type CleanUser struct {
  Gender 	string 	`json:"gender"`
  FirstName 	string 	`json:"first_name"`
  SecondName 	string 	`json:"second_name"`
  Country 	string 	`json:"country"`
  Email 	string 	`json:"email"`
  Uuid 		string 	`json:"uuid"`
  Age 		int 	`json:"age"`
}
