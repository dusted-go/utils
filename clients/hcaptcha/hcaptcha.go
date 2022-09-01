package hcaptcha

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/dusted-go/fault/fault"
)

const (
	verificationEndpoint = "https://hcaptcha.com/siteverify"
)

// Verify checks the user provided captcha response with the hCaptcha server for correctness.
func Verify(ctx context.Context, siteKey, secret, captchaResponse string) (bool, error) {

	// Construct request:
	// ---
	reqBody := fmt.Sprintf(
		"response=%s&secret=%s&sitekey=%s",
		captchaResponse,
		secret,
		siteKey,
	)
	req, err := http.NewRequestWithContext(
		ctx,
		"POST",
		verificationEndpoint,
		bytes.NewBuffer([]byte(reqBody)),
	)
	if err != nil {
		return false,
			fault.SystemWrap(err, "hcaptcha", "Verify", "creating HTTP request failed")
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Send request:
	// ---
	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		return false,
			fault.SystemWrap(err, "hcaptcha", "Verify", "sending HTTP request to hCaptcha failed")
	}

	// Read response:
	// ---
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false,
			fault.SystemWrap(err, "hcaptcha", "Verify", "reading HTTP response body failed")
	}

	// Deserialize response:
	// ---
	captchaResult := struct {
		Success    bool     `json:"success"`
		ErrorCodes []string `json:"error-codes"`
	}{
		Success:    false,
		ErrorCodes: []string{},
	}
	err = json.Unmarshal(body, &captchaResult)
	if err != nil {
		return false,
			fault.SystemWrap(err, "hcaptcha", "Verify",
				"deserializing response body from JSON failed")
	}

	// Interpret result:
	// Docs: https://docs.hcaptcha.com/#siteverify-error-codes-table
	// ---
	if captchaResult.Success {
		return true, nil
	}
	for _, code := range captchaResult.ErrorCodes {

		if code == "missing-input-secret" {
			// Secret key is missing.
			return false,
				fault.SystemWrap(err, "hcaptcha", "Verify",
					"hCaptcha is misconfigured on the server")

		} else if code == "invalid-input-secret" {
			// Secret key is invalid or malformed.
			return false,
				fault.SystemWrap(err, "hcaptcha", "Verify",
					"hCaptcha is misconfigured on the server")

		} else if code == "not-using-dummy-passcode" {
			// You have used a testing sitekey but have not used its matching secret.
			return false,
				fault.SystemWrap(err, "hcaptcha", "Verify",
					"hCaptcha is misconfigured on the server")

		} else if code == "sitekey-secret-mismatch" {
			// The sitekey is not registered with the provided secret.
			return false,
				fault.SystemWrap(err, "hcaptcha", "Verify",
					"hCaptcha is misconfigured on the server")
		}
	}

	// No error, just badly solved captcha:
	return false, nil
}
