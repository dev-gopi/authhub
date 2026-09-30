package dto

import "fmt"

func MapToResponse(id string, name string) interface{} {
	fmt.Printf("Mapping tenant %s to response\n", id)
	return map[string]interface{}{"id": id, "name": name}
}
