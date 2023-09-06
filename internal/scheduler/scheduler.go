package scheduler

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	cipher_client "password-manager/internal/cipher"
	"password-manager/internal/client"
	"strings"
	"time"
)

/*
UserTokenScheduler this scheduler in cycle checks the user token and save it in file on user's pc
*/
func UserTokenScheduler(ctx context.Context, User *client.User) error {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	done := false
	for !done {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			if User.Token != "" {
				type claims struct {
					Exp    int    `json:"exp"`
					Iat    int    `json:"iat"`
					UserID string `json:"user_id"`
				}
				var c claims

				parts := strings.Split(User.Token, ".")
				if len(parts) != 3 {
					return fmt.Errorf("not valid token")
				}

				decodedPayload, err := base64.RawURLEncoding.DecodeString(parts[1])
				if err != nil {
					fmt.Println("Error decoding payload:", err)
					return err
				}

				err = json.Unmarshal(decodedPayload, &c)
				if err != nil {
					fmt.Println("Error unmarshaling payload:", err)
					return err
				}

				kmap, err := cipher_client.ReadEncryptedKey()
				if err != nil {
					return err
				}
				if _, exist := kmap[c.UserID]; exist {
					User.EncryptionKey = []byte(kmap[c.UserID])
					done = true
					return nil
				} else {
					err = cipher_client.SaveEncryptedKey(c.UserID, User.EncryptionKey)
					if err != nil {
						fmt.Println("Can not save user's secret key")
						return err
					}
					done = true
				}

			}
		}
	}
	return nil
}

func UserTokenServerSheduler(ctx context.Context, User *client.User) (string, error) {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()
	done := false
	for !done {
		select {
		case <-ctx.Done():
			return "", ctx.Err()
		case <-ticker.C:
			if User.Token != "" {

				type claims struct {
					Exp    int    `json:"exp"`
					Iat    int    `json:"iat"`
					UserID string `json:"user_id"`
				}
				var c claims

				type UserKey struct {
					Key string `json:"key"`
				}

				var k UserKey

				parts := strings.Split(User.Token, ".")
				if len(parts) != 3 {
					return "", fmt.Errorf("not valid token")
				}

				decodedPayload, err := base64.RawURLEncoding.DecodeString(parts[1])
				if err != nil {
					fmt.Println("Error decoding payload:", err)
					return "", err
				}

				err = json.Unmarshal(decodedPayload, &c)
				if err != nil {
					fmt.Println("Error unmarshaling payload:", err)
					return "", err
				}

				req, err := http.NewRequest(http.MethodGet, "https://localhost:8080/api/user/key", http.NoBody)
				if err != nil {
					return "", err
				}

				req.Header.Set("Authorization", "Bearer "+User.Token)
				resp, err := http.DefaultClient.Do(req)
				defer resp.Body.Close()
				if err != nil {
					return "", err
				}

				respBody, err := io.ReadAll(resp.Body)
				if err != nil {
					return "", err
				}

				if resp.StatusCode < 200 || resp.StatusCode >= 300 {
					return "", fmt.Errorf("Request failed with status: %d - %s", resp.StatusCode, string(respBody))
				}

				err = json.Unmarshal(respBody, &k)
				if err != nil {
					return "", err
				}

				if k.Key != "" {
					User.EncryptionKey = []byte(k.Key)
					done = true
					return "success", nil
				}
				done = true
			}
		}
	}
	return "Empty", nil
}

func SaveUserKey(User *client.User) error {

	type userKey struct {
		Key string `json:"key"`
	}

	var uk userKey

	uk.Key = string(User.EncryptionKey)

	body, err := json.Marshal(uk)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, "https://localhost:8080/api/user/key", bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+User.Token)
	resp, err := http.DefaultClient.Do(req)
	defer resp.Body.Close()
	if err != nil {
		return err
	}
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("Request failed with status: %d - %s", resp.StatusCode, string(respBody))
	}

	return nil
}
