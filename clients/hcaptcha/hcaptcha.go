package hcaptcha

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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
		return false, fmt.Errorf("error creating HTTP request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Send request:
	// ---
	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		return false, fmt.Errorf("error sending HTTP request: %w", err)
	}

	// Read response:
	// ---
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, fmt.Errorf("error reading HTTP response body: %w", err)
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
		return false, fmt.Errorf("error deserializing JSON from HTTP response body: %w", err)
	}

	// Interpret result:
	// Docs: https://docs.hcaptcha.com/#siteverify-error-codes-table
	// ---
	if captchaResult.Success {
		return true, nil
	}
	for _, code := range captchaResult.ErrorCodes {

		if code == "missing-input-secret" || // Secret key is missing.
			code == "invalid-input-secret" || // Secret key is invalid or malformed.
			code == "not-using-dummy-passcode" || // You have used a testing sitekey but have not used its matching secret.
			code == "sitekey-secret-mismatch" { // The sitekey is not registered with the provided secret.
			return false, fmt.Errorf("configuration issue for hCaptcha: %s", code)
		}
	}

	// No error, just badly solved captcha:
	return false, nil
}
