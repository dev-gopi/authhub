package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
	"unicode"

	platformmodel "github.com/dev-gopi/authhub/internal/modules/platformuser/model"
	rootmodel "github.com/dev-gopi/authhub/internal/modules/rootauth/model"
	"github.com/dev-gopi/authhub/internal/platform/database"
	sharedmodel "github.com/dev-gopi/authhub/internal/shared/model"
	"github.com/dev-gopi/authhub/internal/shared/security"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"golang.org/x/term"
	"gorm.io/gorm"
)

func main() {
	username := flag.String("username", "", "root administrator username")
	email := flag.String("email", "", "root administrator email address")
	flag.Parse()
	if flag.NArg() != 0 || strings.TrimSpace(*username) == "" || strings.TrimSpace(*email) == "" {
		log.Fatal("usage: root-admin --username USERNAME --email EMAIL")
	}

	password, err := readPassword("Password: ")
	if err != nil {
		log.Fatalf(
			"read password: %v",
			err,
		)
	}

	confirmPassword, err := readPassword("Confirm password: ")
	if err != nil {
		log.Fatalf(
			"read password confirmation: %v",
			err,
		)
	}

	defer clear(password)
	defer clear(confirmPassword)
	if !bytes.Equal(password, confirmPassword) {
		log.Fatal(
			"password confirmation does not match",
		)
	}

	if err := validatePassword(string(password)); err != nil {
		log.Fatalf("invalid password: %v", err)
	}
	// Load local development settings when present. Production deployments
	// should inject POSTGRES_DSN directly into the environment.
	_ = godotenv.Load()
	postgresDSN := os.Getenv("POSTGRES_DSN")
	if postgresDSN == "" {
		log.Fatal("POSTGRES_DSN is required")
	}

	postgresDB, err := database.NewPostgres(postgresDSN)
	if err != nil {
		log.Fatalf("connect postgres: %v", err)
	}
	defer postgresDB.Close()
	hasher := security.NewPasswordHasher()
	passwordHash, err := hasher.Hash(string(password))
	if err != nil {
		log.Fatalf("hash password: %v", err)
	}
	params, err := json.Marshal(hasher.Params())
	if err != nil {
		log.Fatalf("marshal password parameters: %v", err)
	}
	user, err := createRootAdmin(context.Background(), postgresDB.DB, strings.TrimSpace(*username), strings.ToLower(strings.TrimSpace(*email)), passwordHash, params)

	if err != nil {
		log.Fatalf(
			"create root admin: %v",
			err,
		)
	}

	fmt.Printf("Root admin created successfully (ID: %s, username: %s).\n", user.ID, user.Username)
}

func createRootAdmin(ctx context.Context, db *gorm.DB, username, email, passwordHash string, params []byte) (*platformmodel.PlatformUser, error) {
	now := time.Now().UTC()
	base := func() sharedmodel.BaseModel {
		return sharedmodel.BaseModel{ID: uuid.New(), CreatedAt: now, UpdatedAt: now, IsActive: true}
	}
	user := &platformmodel.PlatformUser{BaseModel: base(), Username: username, Email: email, DisplayName: username, Status: "active", IsRootAdmin: true, CredentialVersion: 1}
	password := &rootmodel.PlatformPassword{BaseModel: base(), PlatformUserID: user.ID, PasswordHash: passwordHash, PasswordAlgorithm: "argon2id", PasswordParams: params, PasswordVersion: 1, ChangedAt: now}
	if err := db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(user).Error; err != nil {
			return fmt.Errorf("create platform user: %w", err)
		}
		if err := tx.Create(password).Error; err != nil {
			return fmt.Errorf("create platform user password: %w", err)
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return user, nil
}

func readPassword(
	label string,
) ([]byte, error) {
	fmt.Print(label)

	value, err := term.ReadPassword(int(os.Stdin.Fd()))
	fmt.Fprintln(os.Stderr)
	return value, err
}

func clear(value []byte) {
	for i := range value {
		value[i] = 0
	}
}

func validatePassword(password string) error {
	if len(password) < 14 {
		return errors.New("must contain at least 14 characters")
	}
	if len(password) > 256 {
		return errors.New("exceeds 256 characters")
	}
	var upper, lower, number bool
	for _, ch := range password {
		switch {
		case unicode.IsUpper(ch):
			upper = true
		case unicode.IsLower(ch):
			lower = true
		case unicode.IsNumber(ch):
			number = true
		}
	}
	if !upper || !lower || !number {
		return errors.New("must contain uppercase, lowercase, and numeric characters")
	}
	if strings.Contains(strings.ToLower(password), "password") {
		return errors.New("is too weak")
	}
	return nil
}
