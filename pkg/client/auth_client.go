package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"nftvc-profile/pkg/config"
	"nftvc-profile/pkg/logger"
	"time"
)

type AuthClient struct {
	log        logger.Logger
	baseURL    string
	httpClient *http.Client
}

func NewAuthClient(log logger.Logger, cfg *config.Config) *AuthClient {
	return &AuthClient{
		log:        log,
		baseURL:    cfg.AuthClientUrl,
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
}

type VerifyTokenRequest struct {
	AccessToken string `json:"access_token"`
}

type VerifyTokenResponse struct {
	AccountId string `json:"account_id,omitempty"`
}

type ChangeRoleRequest struct {
	Role string `json:"role"`
}

type ChangeRoleResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (a *AuthClient) VerifyToken(accessToken string) (*VerifyTokenResponse, error) {
	url := fmt.Sprintf("%s/api/auth/verify-token", a.baseURL)

	requestBody, err := json.Marshal(VerifyTokenRequest{
		AccessToken: accessToken,
	})
	if err != nil {
		a.log.Error("Error marshaling VerifyToken request: ", err)
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		a.log.Error("Error creating new request: ", err)
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := a.httpClient.Do(req)
	if err != nil {
		a.log.Debugf("Error sending request to %s: %v", url, err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		a.log.Error("Non-OK HTTP status: ", resp.Status)
		return nil, fmt.Errorf("non-OK HTTP status: %s", resp.Status)
	}

	var verifyResponse VerifyTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&verifyResponse); err != nil {
		a.log.Error("Error decoding response: ", err)
		return nil, err
	}

	return &verifyResponse, nil
}

func (a *AuthClient) ChangeRole(role string, accessToken string) (*ChangeRoleResponse, error) {
	url := fmt.Sprintf("%s/api/auth/change-role", a.baseURL)

	requestBody, err := json.Marshal(ChangeRoleRequest{
		Role: role,
	})
	if err != nil {
		a.log.Error("Error marshaling VerifyToken request: ", err)
		return nil, err
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		a.log.Error("Error creating new request: ", err)
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	req.Header.Set("Content-Type", "application/json")

	resp, err := a.httpClient.Do(req)
	if err != nil {
		a.log.Debugf("Error sending request to %s: %v", url, err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		a.log.Error("Non-OK HTTP status: ", resp.Status)
		return nil, fmt.Errorf("non-OK HTTP status: %s", resp.Status)
	}

	var changeRole ChangeRoleResponse
	if err := json.NewDecoder(resp.Body).Decode(&changeRole); err != nil {
		a.log.Error("Error decoding response: ", err)
		return nil, err
	}

	return &changeRole, nil
}
