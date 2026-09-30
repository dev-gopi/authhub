package security

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"golang.org/x/crypto/argon2"
)

const (
	argon2Algorithm = "argon2id"

	defaultArgon2Time        uint32 = 3
	defaultArgon2Memory      uint32 = 64 * 1024
	defaultArgon2Parallelism uint8  = 4

	defaultArgon2SaltLength uint32 = 16
	defaultArgon2KeyLength  uint32 = 32
)

type Argon2Params struct {
	Time        uint32 `json:"time"`
	Memory      uint32 `json:"memory"`
	Parallelism uint8  `json:"parallelism"`
	SaltLength  uint32 `json:"salt_length"`
	KeyLength   uint32 `json:"key_length"`
	Version     int    `json:"version"`
}

type PasswordHasher struct {
	params Argon2Params
}

func NewPasswordHasher() *PasswordHasher {
	return &PasswordHasher{
		params: Argon2Params{
			Time:        defaultArgon2Time,
			Memory:      defaultArgon2Memory,
			Parallelism: defaultArgon2Parallelism,
			SaltLength:  defaultArgon2SaltLength,
			KeyLength:   defaultArgon2KeyLength,
			Version:     argon2.Version,
		},
	}
}

func (h *PasswordHasher) Params() Argon2Params {
	return h.params
}

func (h *PasswordHasher) Hash(password string) (string, error) {
	if password == "" {
		return "", errors.New("password cannot be empty")
	}

	salt := make([]byte, h.params.SaltLength)

	if _, err := rand.Read(salt); err != nil {
		return "", fmt.Errorf("generate password salt: %w", err)
	}

	hash := argon2.IDKey(
		[]byte(password),
		salt,
		h.params.Time,
		h.params.Memory,
		h.params.Parallelism,
		h.params.KeyLength,
	)

	encodedSalt := base64.RawStdEncoding.EncodeToString(salt)
	encodedHash := base64.RawStdEncoding.EncodeToString(hash)

	encoded := fmt.Sprintf(
		"$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		h.params.Version,
		h.params.Memory,
		h.params.Time,
		h.params.Parallelism,
		encodedSalt,
		encodedHash,
	)

	return encoded, nil
}

func (h *PasswordHasher) Verify(
	password string,
	encodedHash string,
) (bool, error) {
	params, salt, expectedHash, err := decodeArgon2Hash(encodedHash)
	if err != nil {
		return false, err
	}

	actualHash := argon2.IDKey(
		[]byte(password),
		salt,
		params.Time,
		params.Memory,
		params.Parallelism,
		uint32(len(expectedHash)),
	)

	return subtle.ConstantTimeCompare(
		actualHash,
		expectedHash,
	) == 1, nil
}

func (h *PasswordHasher) NeedsRehash(encodedHash string) (bool, error) {
	params, _, _, err := decodeArgon2Hash(encodedHash)
	if err != nil {
		return false, err
	}

	return params.Time != h.params.Time ||
		params.Memory != h.params.Memory ||
		params.Parallelism != h.params.Parallelism ||
		params.SaltLength != h.params.SaltLength ||
		params.KeyLength != h.params.KeyLength ||
		params.Version != h.params.Version, nil
}

func decodeArgon2Hash(
	encodedHash string,
) (Argon2Params, []byte, []byte, error) {
	var params Argon2Params

	parts := strings.Split(encodedHash, "$")

	if len(parts) != 6 {
		return params, nil, nil, errors.New(
			"invalid argon2 hash format",
		)
	}

	if parts[1] != argon2Algorithm {
		return params, nil, nil, errors.New(
			"unsupported password algorithm",
		)
	}

	versionString := strings.TrimPrefix(parts[2], "v=")

	version, err := strconv.Atoi(versionString)
	if err != nil {
		return params, nil, nil, fmt.Errorf(
			"parse argon2 version: %w",
			err,
		)
	}

	var memory uint32
	var iterations uint32
	var parallelism uint8

	if _, err := fmt.Sscanf(
		parts[3],
		"m=%d,t=%d,p=%d",
		&memory,
		&iterations,
		&parallelism,
	); err != nil {
		return params, nil, nil, fmt.Errorf(
			"parse argon2 parameters: %w",
			err,
		)
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return params, nil, nil, fmt.Errorf(
			"decode argon2 salt: %w",
			err,
		)
	}

	hash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return params, nil, nil, fmt.Errorf(
			"decode argon2 hash: %w",
			err,
		)
	}

	params = Argon2Params{
		Time:        iterations,
		Memory:      memory,
		Parallelism: parallelism,
		SaltLength:  uint32(len(salt)),
		KeyLength:   uint32(len(hash)),
		Version:     version,
	}

	return params, salt, hash, nil
}
