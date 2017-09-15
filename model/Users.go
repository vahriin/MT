package model

type User struct {
	Id   Id
	Nick string
}

type PassUser struct {
	Id       Id
	Email    string
	Nick     string
	PassHash []byte
}

func (pu PassUser) ToUser() User {
	return User{
		Id:   pu.Id,
		Nick: pu.Nick,
	}
}
