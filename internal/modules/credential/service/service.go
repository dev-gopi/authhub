package service

import "fmt"

type Service struct{}

func NewService() *Service {
	fmt.Println("Creating credential service...")
	return &Service{}
}
