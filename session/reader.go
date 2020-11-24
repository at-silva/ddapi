//Package session contains tools to introspect session scoped variables.
package session

import (
	"fmt"

	"github.com/dgrijalva/jwt-go"
)

type (
	// Reader reads the session data into a parameter map
	//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . Reader
	Reader interface {
		Copy(token string, pm map[string]interface{}) error
	}

	// Copy copies the params present in the session into a params map
	Copy func(token string, pm map[string]interface{}) error
)

// Copy all the params found in the session into the given map
func (f Copy) Copy(t string, pm map[string]interface{}) error {
	return f(t, pm)
}

// HS256JWT Copies all the claims in the given JWT into the given params map
func HS256JWT(secret []byte) Copy {
	return func(t string, pm map[string]interface{}) error {
		if t == "" {
			return fmt.Errorf("could not parse jwt: empty token")
		}

		if pm == nil {
			return fmt.Errorf("could not parse token: nil parameter map")
		}

		token, err := jwt.Parse(t, func(token *jwt.Token) (interface{}, error) {
			if token.Method != jwt.SigningMethodHS256 {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			return secret, nil
		})

		if err != nil {
			return fmt.Errorf("could not parse jwt: %w", err)
		}

		if !token.Valid {
			return fmt.Errorf("invalid jwt")
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			for k, v := range claims {
				pm[k] = v
			}
		}

		return nil
	}
}

// None no-op session reader
func None(_ []byte) Copy {
	return func(_ string, _ map[string]interface{}) error {
		return nil
	}
}
