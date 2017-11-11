package model

type User struct {
	Id   Id
	Nick string
}

type GoogleUser struct {
	User
	googleId []byte
}
