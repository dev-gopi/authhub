package dto

import "fmt"

func MapToResponse(id string, email string) interface{} {
	fmt.Printf("Mapping platform user %s to response\n", id)
	return map[string]interface{}{"id": id, "email": email}
}
