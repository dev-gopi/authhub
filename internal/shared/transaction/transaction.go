package transaction

import "fmt"

func ExecuteInTransaction(fn func() error) error {
	fmt.Println("Executing in transaction...")
	return fn()
}
