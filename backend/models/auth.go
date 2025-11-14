package models

import (
    "errors"
    "golang.org/x/crypto/bcrypt"
)

const HashCost = 12 // OWASP recommends 10–14

func HashPassword(password string) (string, error) {
    if len(password) < 8 {
        return "", errors.New("password too short")
    }
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), HashCost)
    return string(bytes), err
}

func CheckPassword(password, hashed string) bool {
    return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password)) == nil
}

func NeedsRehash(hashed string) bool {
    c, err := bcrypt.Cost([]byte(hashed))
    return err == nil && c < HashCost
}
