package handlers

// import (
// 	"testing"

// 	"github.com/stretchr/testify/assert"
// )

// var authoritySecret = []byte("authority secret")
// var sessionSecret = []byte("session secret")

// func Test__Success(t *testing.T) {
// 	authorityToken, err := CreateAuthorityToken(authoritySecret, "local_handlers")

// 	if err != nil {
// 		t.Fatalf("Failed to create authority token: %s", err)
// 	}

// 	sessionToken, err := CreateSessionToken(authoritySecret, sessionSecret, authorityToken)

// 	if err != nil {
// 		t.Fatalf("Failed to create session token: %s", err)
// 	}

// 	authorName, err := ValidateToken(sessionSecret, sessionToken)

// 	if err != nil {
// 		t.Fatalf("Failed to create session token: %s", err)
// 	}

// 	assert.Equal(t, "local_handlers", authorName)
// }
