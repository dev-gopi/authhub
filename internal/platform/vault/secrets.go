package vault

import "fmt"

func GetSecret(path string) map[string]interface{} {
	fmt.Printf("Getting secret from %s\n", path)
	return map[string]interface{}{}
}
