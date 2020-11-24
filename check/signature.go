package check

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
)

// SignatureChecker represents a payload signature Checker
//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . SignatureChecker
type (
	SignatureChecker interface {
		Check(payload, signature []byte) error
	}

	// Signature checker function type
	Signature func(payload, signature []byte) error
)

// Check checks the signature of a given payload
func (f Signature) Check(p, s []byte) error {
	return f(p, s)
}

// Sha256HMAC returns a hmac/sha256 base query validator
func Sha256HMAC(secret []byte) Signature {
	return func(p, s []byte) error {
		if len(p) == 0 {
			return fmt.Errorf("payload cannot be empty")
		}

		if len(s) == 0 {
			return fmt.Errorf("signature cannot be empty")
		}

		hash := hmac.New(sha256.New, secret)
		_, err := hash.Write(p)
		if err != nil {
			return fmt.Errorf("could not write to hash stream: %w", err)
		}
		ss := hash.Sum(nil)

		if !hmac.Equal(s, ss) {
			return fmt.Errorf("invalid signature")
		}

		return nil
	}
}
