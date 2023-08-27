package client

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

const (
	signUpEndpoint = `/api/user`
	BaseURL        = `https://localhost:8080`
	logInEndpoint  = `/api/login`
)

type User struct {
	client        http.Client
	Token         string
	Passwords     map[string]string
	EncryptionKey []byte
}

type UserData struct {
	Username string `json:"login"`
	Password string `json:"password"`
}

type UserPass struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type NewUserPass struct {
	NewName string `json:"new_name"`
	NewPass string `json:"new_pass"`
	OldName string `json:"old_name"`
}

func NewUser(EncryptionKey []byte) *User {

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	return &User{
		client:        http.Client{Transport: tr},
		EncryptionKey: EncryptionKey}
}

func (u *User) SignUp(username, password string) (string, error) {
	url := BaseURL + signUpEndpoint

	user := UserData{
		Username: username,
		Password: password,
	}
	body, err := json.Marshal(user)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := u.client.Do(req)
	if err != nil {
		return "", err
	}
	respBody, err := io.ReadAll(resp.Body)

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("Request failed with status: %d - %s", resp.StatusCode, string(respBody))
	}

	token := resp.Header.Get("Authorization")
	if token != "" {
		return token, nil
	}
	return "", errors.New("token is empty")
}

func (u *User) Login(username, password string) (string, error) {
	url := BaseURL + logInEndpoint

	user := UserData{
		Username: username,
		Password: password,
	}
	body, err := json.Marshal(user)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := u.client.Do(req)
	if err != nil {
		return "", err
	}
	respBody, err := io.ReadAll(resp.Body)

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("Request failed with status: %d - %s", resp.StatusCode, string(respBody))
	}

	token := resp.Header.Get("Authorization")
	if token != "" {
		return token, nil
	}
	return "", errors.New("token is empty")

}

func (u *User) Request(method, endpoint string, data interface{}) (interface{}, error) {
	url := BaseURL + endpoint

	req, err := http.NewRequest(method, url, http.NoBody)

	if method == http.MethodDelete || method == http.MethodGet {

		req.URL.Path = url + "/" + data.(string)

	} else {
		var requestBody []byte
		if data != nil {
			requestBody, _ = json.Marshal(data)
		}

		req, err = http.NewRequest(method, url, bytes.NewBuffer(requestBody))
		if err != nil {
			return nil, err
		}
	}

	req.Header.Set("Authorization", "Bearer "+u.Token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := u.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("Request failed with status: %d - %s", resp.StatusCode, string(respBody))
	}

	if method == http.MethodGet {
		var res UserPass
		err := json.Unmarshal(respBody, &res)

		if err != nil {
			return nil, err
		}
		return res.Password, nil
	}

	return nil, nil
}
