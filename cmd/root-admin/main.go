package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	platformrepository "github.com/dev-gopi/authhub/internal/modules/platformuser/repository"
	rootrepository "github.com/dev-gopi/authhub/internal/modules/rootauth/repository"
	rootbootstrap "github.com/dev-gopi/authhub/internal/modules/rootauth/service"

	"github.com/dev-gopi/authhub/internal/config"
	"github.com/dev-gopi/authhub/internal/platform/database"
	"github.com/dev-gopi/authhub/internal/shared/security"

	"golang.org/x/term"
)

func main() {
	ctx := context.Background()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf(
			"load configuration: %v",
			err,
		)
	}

	postgresDB, err := database.NewPostgres(
		cfg.Postgres.DSN,
	)
	if err != nil {
		log.Fatalf(
			"connect postgres: %v",
			err,
		)
	}

	defer func() {
		if err := postgresDB.Close(); err != nil {
			log.Printf(
				"close postgres: %v",
				err,
			)
		}
	}()

	reader := bufio.NewReader(os.Stdin)

	username := readLine(
		reader,
		"Username: ",
	)

	email := readLine(
		reader,
		"Email: ",
	)

	displayName := readLine(
		reader,
		"Display name: ",
	)

	password, err := readPassword(
		"Password: ",
	)
	if err != nil {
		log.Fatalf(
			"read password: %v",
			err,
		)
	}

	confirmPassword, err := readPassword(
		"Confirm password: ",
	)
	if err != nil {
		log.Fatalf(
			"read password confirmation: %v",
			err,
		)
	}

	if password != confirmPassword {
		log.Fatal(
			"password confirmation does not match",
		)
	}

	platformUserRepository :=
		platformrepository.NewPostgresRepository(
			postgresDB.DB,
		)

	passwordRepository :=
		rootrepository.NewPasswordRepository(
			postgresDB.DB,
		)

	passwordHasher :=
		security.NewPasswordHasher()

	service := rootbootstrap.NewService(
		postgresDB,
		platformUserRepository,
		passwordRepository,
		passwordHasher,
	)

	user, err := service.CreateRootAdmin(
		ctx,
		rootbootstrap.CreateRootAdminRequest{
			Username:    username,
			Email:       email,
			DisplayName: displayName,
			Password:    password,
		},
	)

	// Reduce the amount of time plaintext passwords
	// remain referenced by this function.
	password = ""
	confirmPassword = ""

	if err != nil {
		log.Fatalf(
			"create root admin: %v",
			err,
		)
	}

	fmt.Println()
	fmt.Println("Root Admin created successfully.")
	fmt.Printf(
		"ID: %s\n",
		user.ID,
	)
	fmt.Printf(
		"Username: %s\n",
		user.Username,
	)
}

func readLine(
	reader *bufio.Reader,
	label string,
) string {
	fmt.Print(label)

	value, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalf(
			"read input: %v",
			err,
		)
	}

	return strings.TrimSpace(value)
}

func readPassword(
	label string,
) (string, error) {
	fmt.Print(label)

	value, err := term.ReadPassword(
		int(os.Stdin.Fd()),
	)

	fmt.Println()

	if err != nil {
		return "", err
	}

	return string(value), nil
}
