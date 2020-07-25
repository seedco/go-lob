package lob

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"strconv"
	"time"
)

// Errors returned by the webbook parsing
var (
	ErrInvalidHeader    = errors.New("webhook has invalid Lob-Signature header(s)")
	ErrNotTimestamped   = errors.New("no Lob-Signature-Timestamp header provided")
	ErrNotSigned        = errors.New("no Lob-Signature header provided")
	ErrNoValidSignature = errors.New("webhook has no valid signature")
	ErrTooOld           = errors.New("timestamp wasn't within tolerance")
)

//ValidateWebhookPayload valides the passed in payload and signature header
func ValidateWebhookPayload(payload []byte, timestampHeader, sigHeader, secret string, tolerance time.Duration) error {
	if sigHeader == "" {
		return ErrNotSigned
	}

	if err := hasValidTimestampWithTolerance(timestampHeader, tolerance); err != nil {
		return err
	}

	sig, err := hex.DecodeString(sigHeader)
	if err != nil {
		return ErrInvalidHeader
	}

	computed := computeSignature(timestampHeader, payload, secret)

	if hmac.Equal(computed, sig) {
		return nil
	}

	return ErrNoValidSignature
}

func hasValidTimestampWithTolerance(timestampHeader string, tolerance time.Duration) error {
	if timestampHeader == "" {
		return ErrNotTimestamped
	}

	t, err := strconv.ParseInt(timestampHeader, 10, 64)
	if err != nil {
		return ErrInvalidHeader
	}

	// lob.com sends the epoch time with milliseconds
	sentTime := time.Unix(0, t*int64(time.Millisecond))

	if time.Since(sentTime) > tolerance {
		return ErrTooOld
	}

	return nil
}

func computeSignature(timestamp string, payload []byte, secret string) []byte {
	mac := hmac.New(sha256.New, []byte(secret))

	mac.Write([]byte(timestamp))
	mac.Write([]byte("."))
	mac.Write(payload)

	return mac.Sum(nil)
}
