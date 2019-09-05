package cryptography

import (
	"errors"
	"strconv"

	"golang.org/x/crypto/bcrypt"
)

const (
	majorVersion = '2'
	minorVersion = 'a'
	//minorVersion       = 'b'
	maxSaltSize        = 16
	maxCryptedHashSize = 23
	encodedSaltSize    = 22
	encodedHashSize    = 31
	minHashSize        = 59
)

const (
	MinCost     int = 4  // the minimum allowable cost as passed in to GenerateFromPassword
	MaxCost     int = 31 // the maximum allowable cost as passed in to GenerateFromPassword
	DefaultCost int = 10 // the cost that will actually be set if a cost below MinCost is passed into GenerateFromPassword
)

var ErrMismatchedHashAndPassword = errors.New("crypto/bcrypt: hashedPassword is not the hash of the given password")
var ErrHashTooShort = errors.New("crypto/bcrypt: hashedSecret too short to be a bcrypted password")
var ErrInvalidHash = errors.New("crypto/bcrypt: invalid hashedPassword")
var ErrInvalidVersion = errors.New("crypto/bcrypt: invalid version")

type bcryptStruct struct{}

func BCrypt() *bcryptStruct {
	return &bcryptStruct{}
}

func (*bcryptStruct) Version(hashedBytes []byte) ([]byte, error) {
	if hashedBytes[0] != '$' {
		return nil, ErrInvalidHash
	}

	if hashedBytes[1] > majorVersion {
		return nil, ErrInvalidVersion
	}
	if hashedBytes[2] != '$' {
		return hashedBytes[1:3], nil
	}

	return hashedBytes[1:2], nil
}

func (*bcryptStruct) Cost(hashedBytes []byte) (int, error) {
	if len(hashedBytes) < minHashSize {
		return 0, ErrHashTooShort
	}

	if hashedBytes[0] != '$' {
		return 0, ErrInvalidHash
	}

	if hashedBytes[2] != '$' {
		cost, err := strconv.Atoi(string(hashedBytes[4:6]))
		if err != nil {
			return -1, err
		}
		return cost, nil
	}

	cost, err := strconv.Atoi(string(hashedBytes[5:7]))
	if err != nil {
		return -1, err
	}
	return cost, nil
}

func (*bcryptStruct) HashPassword(password string) (string, error) {
	if password == "" {
		return "", errors.New("password is empty")
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func (*bcryptStruct) IsCorrectPassword(hashedPassword, password string) (bool, error) {
	if hashedPassword == "" {
		return false, errors.New("hashedPassword is empty")
	}
	if password == "" {
		return false, errors.New("password is empty")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)); err != nil {
		return false, err
	}
	return true, nil
}
