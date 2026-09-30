package repository

type Interface interface {
	FindByID(id string) (interface{}, error)
}
