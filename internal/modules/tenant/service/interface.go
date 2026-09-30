package service

type Interface interface {
	CreateTenant(name string) error
}
