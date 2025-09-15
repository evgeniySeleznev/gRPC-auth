package models

type App struct {
	ID     int
	Name   string
	Secret string //чтобы подписывать и валидировать токены на стороне клиентского приложения
}
