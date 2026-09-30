package service

type Interface interface {
	CreatePlatformUser(email string) error
}
