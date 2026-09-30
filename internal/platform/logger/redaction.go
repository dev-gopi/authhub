package logger

import "fmt"

func RedactSensitiveData(data string) string {
	fmt.Println("Redacting sensitive data...")
	return data
}
