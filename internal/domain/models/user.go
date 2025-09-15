package models

type User struct {
	ID       int64
	Email    string
	PassHash []byte //salt чтобы не хранить пароли в открытом виде
}
