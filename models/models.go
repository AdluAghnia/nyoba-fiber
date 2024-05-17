package models

type User struct {
	ID       int
	Username string
	Password string
	IsLogin  bool
}

type LoginResponse struct {
	Token string
}
