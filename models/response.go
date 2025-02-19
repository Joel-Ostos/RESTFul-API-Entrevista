package models

type Response struct { 
  Users 		[]CleanUser 	`json:"users"`
  UsersPerCountry 	map[string]int 	`json:"users_per_country"`
  AverageAge 		float64 	`json:"average_age"`
  MenUsers 		int 		`json:"men_users"`
  WomenUsers		int		`json:"women_users"`
}
