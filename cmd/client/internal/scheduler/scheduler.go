package scheduler

import (
	cipher_client "client-password/internal/cipher"
	"client-password/internal/client"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

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
