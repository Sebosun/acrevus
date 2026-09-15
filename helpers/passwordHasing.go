package helpers

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

const maxPasswordLength = 1024

var ErrPasswordTooLong = errors.New("password must not exceed 1024 bytes")

type Params struct {
	memory      uint32
	iterations  uint32
	parallelism uint8
	saltLength  uint32
	hashLength  uint32
}

var defaultArgon2Params = Params{
	memory:      19 * 1024,
	iterations:  2,
	parallelism: 1,
	saltLength:  16,
	hashLength:  32,
}

func HashPassword(password string) (string, error) {
	if len(password) > maxPasswordLength {
		return "", ErrPasswordTooLong
	}

	params := defaultArgon2Params
	salt, err := genCryptoSalt(params.saltLength)
	if err != nil {
		return "", err
	}

	hash := argon2.IDKey(
		[]byte(password), salt, params.iterations, params.memory, params.parallelism, params.hashLength,
	)

	return fmt.Sprintf(
		"$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version,
		params.memory,
		params.iterations,
		params.parallelism,
		base64.RawStdEncoding.EncodeToString(salt),
		base64.RawStdEncoding.EncodeToString(hash),
	), nil
}

func CheckPassword(password, encodedHash string) bool {
	if len(password) > maxPasswordLength {
		return false
	}

	params, salt, expectedHash, err := decodeArgon(encodedHash)
	if err != nil {
		return false
	}

	actualHash := argon2.IDKey(
		[]byte(password), salt, params.iterations, params.memory, params.parallelism, uint32(len(expectedHash)),
	)
	return subtle.ConstantTimeCompare(actualHash, expectedHash) == 1
}

func genCryptoSalt(saltSize uint32) ([]byte, error) {
	salt := make([]byte, saltSize)
	if _, err := rand.Read(salt); err != nil {
		return nil, fmt.Errorf("salt generation failed: %w", err)
	}
	return salt, nil
}

func decodeArgon(encodedHash string) (Params, []byte, []byte, error) {
	components := strings.Split(encodedHash, "$")
	if len(components) != 6 || components[0] != "" || components[1] != "argon2id" {
		return Params{}, nil, nil, errors.New("invalid Argon2id hash format")
	}

	version, err := parsePrefixedUint(components[2], "v=", 32)
	if err != nil || version != uint64(argon2.Version) {
		return Params{}, nil, nil, errors.New("unsupported Argon2id version")
	}

	parameterParts := strings.Split(components[3], ",")
	if len(parameterParts) != 3 {
		return Params{}, nil, nil, errors.New("invalid Argon2id parameters")
	}

	memory, memoryErr := parsePrefixedUint(parameterParts[0], "m=", 32)
	iterations, iterationErr := parsePrefixedUint(parameterParts[1], "t=", 32)
	parallelism, parallelismErr := parsePrefixedUint(parameterParts[2], "p=", 8)

	if memoryErr != nil || iterationErr != nil || parallelismErr != nil ||
		memory == 0 || memory > 512*1024 || iterations == 0 || iterations > 10 || parallelism == 0 || parallelism > 16 {
		return Params{}, nil, nil, errors.New("invalid Argon2id parameters")
	}

	salt, err := base64.RawStdEncoding.DecodeString(components[4])
	if err != nil || len(salt) < 8 {
		return Params{}, nil, nil, errors.New("invalid Argon2id salt")
	}

	hash, err := base64.RawStdEncoding.DecodeString(components[5])
	if err != nil || len(hash) < 16 || len(hash) > 64 {
		return Params{}, nil, nil, errors.New("invalid Argon2id hash")
	}

	return Params{
		memory:      uint32(memory),
		iterations:  uint32(iterations),
		parallelism: uint8(parallelism),
		saltLength:  uint32(len(salt)),
		hashLength:  uint32(len(hash)),
	}, salt, hash, nil
}

func parsePrefixedUint(value, prefix string, bitSize int) (uint64, error) {
	if !strings.HasPrefix(value, prefix) {
		return 0, errors.New("missing parameter prefix")
	}
	return strconv.ParseUint(strings.TrimPrefix(value, prefix), 10, bitSize)
}
