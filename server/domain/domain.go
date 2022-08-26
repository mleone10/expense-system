package domain

type OrgId string
type UserId string

type Organization struct {
	Id      OrgId
	Name    string
	Members []User
}

type User struct {
	Id   UserId
	Name string
}
