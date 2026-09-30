package service

import "fmt"

type Service struct{}

func NewService() *Service {
	fmt.Println("Creating platformuser service...")
	return &Service{}
}
