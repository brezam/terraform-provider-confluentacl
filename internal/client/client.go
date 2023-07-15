package client

import (
	"sync"
	"terraform-provider-confluentacl/internal/client/request"
	"time"
)

type CachedAccessToken struct {
	token   string
	expires *time.Time
	mutex   sync.Mutex
}

type Client struct {
	cloudApiKey    string
	cloudApiSecret string
	accessToken    CachedAccessToken
	cache          map[string]interface{}
	cacheMutex     sync.RWMutex
}

const baseApiUrl = "https://confluent.cloud/api/"

func New(cloudApiKey, cloudApiSecret string) *Client {
	return &Client{
		cloudApiKey:    cloudApiKey,
		cloudApiSecret: cloudApiSecret,
		accessToken:    CachedAccessToken{},
		cache:          make(map[string]interface{}),
		cacheMutex:     sync.RWMutex{},
	}
}

func (c *Client) RequestBuilder() *request.Request {
	return request.NewRequestWithBasicAuth(baseApiUrl, c.cloudApiKey, c.cloudApiSecret)
}

func (c *Client) KafkaRestRequestBuilder(kafkaHttpEndpoint string) (*request.Request, error) {
	token, err := c.GetAccessToken()
	if err != nil {
		return nil, err
	}
	return request.NewRequestWithBearerAuth(kafkaHttpEndpoint, token), nil
}
