package service

import "fmt"

type Service struct{}

func NewService() *Service {
	fmt.Println("Creating tenant service...")
	return &Service{}
}
