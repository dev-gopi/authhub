package service

type Interface interface {
	CreateCredential(userID string) error
}
