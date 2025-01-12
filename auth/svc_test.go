package auth

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestCreateToken_Success(t *testing.T) {
    email := "test@example.com"
    password := "password123"
    role := "admin"

    // Act: Call CreateToken
    tokenString, err := CreateToken(email, password, role)

    // Assert: Check for errors and validate token
    assert.NoError(t, err, "Expected no error, but got one")
    assert.NotEmpty(t, tokenString, "Expected token string to be non-empty")
}