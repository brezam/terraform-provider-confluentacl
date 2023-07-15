package client

import (
	"encoding/base64"
	"encoding/json"
	"strings"
	"terraform-provider-confluentacl/internal/client/request"
	"time"
)

const (
	accessTokenEndpoint = "access_tokens"
)

type AccessTokenResponse struct {
	Error string `json:"error"`
	Token string `json:"token"`
}

func (c *Client) GetAccessToken() (string, error) {
	c.accessToken.mutex.Lock()
	if c.accessToken.token != "" && c.accessToken.expires != nil && (time.Until(*c.accessToken.expires).Minutes() > 5) {
		c.accessToken.mutex.Unlock()
		return c.accessToken.token, nil
	}
	err := c.generateAccessToken()
	c.accessToken.mutex.Unlock()
	if err != nil {
		return "", err
	}
	return c.accessToken.token, nil
}

type JwtToken struct {
	Exp int64 `json:"exp"`
}

func (c *Client) generateAccessToken() error {
	response, err := c.RequestBuilder().
		Endpoint(accessTokenEndpoint).
		SetBody(struct{}{}).
		Post().
		ExecuteAndRetryOn429()
	if err != nil {
		return err
	}
	if err != nil {
		return err
	}
	var accessTokenResponse AccessTokenResponse
	err = request.UnpackJSONResponse(response, &accessTokenResponse)
	if err != nil {
		return err
	}
	fullToken := accessTokenResponse.Token
	middleTokenBit := strings.Split(fullToken, ".")[1]
	var jwtToken JwtToken
	result, err := base64.RawStdEncoding.DecodeString(middleTokenBit)
	if err != nil {
		panic(err)
	}
	json.Unmarshal(result, &jwtToken)
	c.accessToken.token = fullToken
	expTime := time.Unix(jwtToken.Exp, 0)
	c.accessToken.expires = &expTime
	return nil
}
