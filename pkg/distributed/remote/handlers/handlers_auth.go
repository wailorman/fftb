package handlers

// import (
// 	"strings"
// 	"time"

// 	"github.com/dgrijalva/jwt-go"
// 	"github.com/labstack/echo/v4"
// 	"github.com/pkg/errors"
// 	"github.com/wailorman/fftb/pkg/distributed/models"
// )

// // SessionTimeout _
// var SessionTimeout = time.Duration(2 * time.Minute)

// // func authSkipperfunc(c echo.Context) bool {
// // 	if c.Path() == "/authorities" {
// // 		return true
// // 	}

// // 	return false
// // }

// // // JWTMiddleware _
// // func JWTMiddleware(sessionSecret []byte) echo.MiddlewareFunc {
// // 	return middleware.JWTWithConfig(middleware.JWTConfig{
// // 		Skipper:    authSkipperfunc,
// // 		SigningKey: sessionSecret,
// // 		ContextKey: "session",
// // 	})
// // }

// func extractToken(c echo.Context) (string, error) {
// 	if len(c.Request().Header["Authorization"]) == 0 {
// 		return "", models.ErrMissingAccessToken
// 	}

// 	authHeader := c.Request().Header["Authorization"][0]
// 	authParts := strings.Split(authHeader, " ")

// 	if len(authParts) < 2 {
// 		return "", models.ErrMissingAccessToken
// 	}

// 	return strings.ReplaceAll(authParts[1], "\"", ""), nil
// }

// func extractAuthor(c echo.Context) models.IAuthor {
// 	author, ok := c.Get("author").(models.IAuthor)

// 	if !ok {
// 		return nil
// 	}

// 	return author
// }

// // JWTMiddleware _
// func JWTMiddleware(sessionSecret []byte) echo.MiddlewareFunc {
// 	return func(next echo.HandlerFunc) echo.HandlerFunc {
// 		return func(c echo.Context) error {
// 			if c.Path() == "/authorities" || c.Path() == "/sessions" {
// 				return next(c)
// 			}

// 			token, err := extractToken(c)

// 			if err != nil {
// 				return c.JSON(newAPIError(err))
// 			}

// 			uid, err := ValidateToken(sessionSecret, token)

// 			if err != nil {
// 				unauthorizedErr := errors.Wrapf(models.ErrInvalidSessionKey, "Failed to validate session token `%s`: %s", token, err)
// 				return c.JSON(newAPIError(unauthorizedErr))
// 			}

// 			author := &models.Author{Name: uid}
// 			c.Set("author", author)

// 			return next(c)
// 		}
// 	}
// }

// // ValidateToken _
// func ValidateToken(secret []byte, key string) (name string, err error) {
// 	token, err := jwt.Parse(key, func(token *jwt.Token) (interface{}, error) {
// 		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
// 			return nil, errors.Wrapf(models.ErrUnknown, "Unexpected signing method: %v", token.Header["alg"])
// 		}

// 		return secret, nil
// 	})

// 	if err != nil {
// 		return "", errors.Wrap(err, "Parsing token")
// 	}

// 	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
// 		uidClaim, ok := claims["uid"].(string)

// 		if !ok {
// 			return "", errors.Wrapf(models.ErrUnknown, "Unexpected uid claim value: %#v", claims["uid"])
// 		}

// 		return uidClaim, nil
// 	}

// 	return "", errors.Wrap(models.ErrUnauthorized, "Invalid token")
// }

// // CreateAuthorityToken _
// func CreateAuthorityToken(authoritySecret []byte, name string) (string, error) {
// 	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
// 		"uid": name,
// 		"iat": time.Now().Unix(),
// 	})

// 	// Sign and get the complete encoded token as a string using the secret
// 	tokenString, err := token.SignedString(authoritySecret)

// 	if err != nil {
// 		return "", errors.Wrap(err, "Signing token error")
// 	}

// 	return tokenString, nil
// }

// // CreateSessionToken _
// func CreateSessionToken(authoritySecret, sessionSecret []byte, authorityKey string) (string, error) {
// 	uid, err := ValidateToken(authoritySecret, authorityKey)

// 	if err != nil {
// 		return "", errors.Wrapf(models.ErrInvalidAuthorityKey, "Validating authority key: `%s`", err)
// 	}

// 	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
// 		"uid": uid,
// 		"iat": time.Now().Unix(),
// 		"exp": time.Now().Add(SessionTimeout),
// 	})

// 	tokenString, err := token.SignedString(sessionSecret)

// 	if err != nil {
// 		return "", errors.Wrap(err, "Signing token error")
// 	}

// 	return tokenString, nil
// }
