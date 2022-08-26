package domain

type Id string

type Organization struct {
	Id      Id
	Name    string
	Members []Member
}

type Member struct {
	Name string
}
