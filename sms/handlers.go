package sms

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type SMSPayload struct {
	Phone string `json:"phone"`
	Token string `json:"token"`
}

func SendSMSService(to, token string) error {
	apiURL := os.Getenv("SMS_API_URL")
	apiKey := os.Getenv("SMS_API_KEY")
	username := os.Getenv("SMS_USERNAME")
	sender := os.Getenv("SMS_SENDER_ID")

	if apiURL == "" || apiKey == "" {
		return fmt.Errorf("SMS_API_URL or SMS_API_KEY is not set")
	}
	if to == "" {
		return fmt.Errorf("no recipient phone number provided")
	}

	message := fmt.Sprintf("Your verification code is %s", token)

	form := url.Values{}
	form.Set("username", username)
	form.Set("to", to)
	form.Set("message", message)
	if sender != "" {
		form.Set("from", sender)
	}

	req, err := http.NewRequest(http.MethodPost, apiURL, strings.NewReader(form.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("apiKey", apiKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("sms provider returned status %d: %s", resp.StatusCode, string(body))
	}

	log.Printf("sms sent to %s", to)
	return nil
}
