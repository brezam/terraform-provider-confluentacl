package client

import (
	"fmt"
	"terraform-provider-confluentacl/internal/client/request"
)

type LogicalCluster struct {
	ID   string  `json:"id"`
	Type *string `json:"type,omitempty"`
}

// Response Internal API
type ApiKeyResponseInternal struct {
	ApiKey ApiKeyInternal `json:"api_key"`
}
type ApiKeyInternal struct {
	Key             string           `json:"key"`
	Secret          string           `json:"secret"`
	HashedSecret    string           `json:"hashed_secret"`
	HashedFunction  string           `json:"hashed_function"`
	SASLMechanism   string           `json:"sasl_mechanism"`
	UserID          int              `json:"user_id"`
	Deactived       bool             `json:"deactived"`
	ID              int              `json:"id"` // This ID is an internal id. It's an integer (e.g.: 9712207)
	Description     string           `json:"description"`
	LogicalClusters []LogicalCluster `json:"logical_clusters"`
	AccountID       string           `json:"account_id"`
	ServiceAccount  bool             `json:"service_account"`
}

// Response IAM V2 Api
type ApiKeyIamV2 struct {
	ID   string          `json:"id"` // This ID is the same as the 'Key' of the internal object (e.g.: 7VNXC6EGOZ32SDAH)
	Spec ApiKeyIamV2Spec `json:"spec"`
}
type ApiKeyIamV2Spec struct {
	Description string              `json:"description"`
	Resource    ApiKeyIamV2Resource `json:"resource"`
	Owner       ApiKeyIamV2Owner    `json:"owner"`
}
type ApiKeyIamV2Resource struct {
	ID string `json:"id"`
}
type ApiKeyIamV2Owner struct {
	ID string `json:"id"`
}

// Create Request
type ApiKeyCreateRequestW struct {
	ApiKey *ApiKeyCreateRequest `json:"api_key"`
}
type ApiKeyCreateRequest struct {
	AccountID       string           `json:"account_id"`
	UserID          int              `json:"user_id,omitempty"`
	Description     string           `json:"description,omitempty"`
	LogicalClusters []LogicalCluster `json:"logical_clusters"`
}

// Update Request
type ApiKeyUpdateRequestW struct {
	ApiKey *ApiKeyUpdateRequest `json:"api_key"`
}
type ApiKeyUpdateRequest struct {
	ID              string           `json:"id"`
	AccountID       string           `json:"account_id"`
	Description     string           `json:"description"`
	LogicalClusters []LogicalCluster `json:"logical_clusters"`
}

// Delete Request
type ApiKeyDeleteRequestW struct {
	ApiKey *ApiKeyDeleteRequest `json:"api_key"`
}
type ApiKeyDeleteRequest struct {
	ID              string           `json:"id"`
	AccountID       string           `json:"account_id"`
	LogicalClusters []LogicalCluster `json:"logical_clusters"`
}

const (
	createApiKeyEndpoint = "api_keys"           // internal api
	readApiKeyEndpoint   = "iam/v2/api-keys/%s" // iam v2 api
	updateApiKeyEndpoint = "api_keys/%s"        // internal api
	deleteApiKeyEndpoint = "api_keys/%s"        // internal api
)

func (c *Client) CreateApiKey(userId int, envId, resourceId, description string) (*ApiKeyInternal, error) {
	body := &ApiKeyCreateRequestW{
		&ApiKeyCreateRequest{
			AccountID:       envId,
			UserID:          userId,
			LogicalClusters: []LogicalCluster{{ID: resourceId}},
			Description:     description,
		},
	}
	responseBody := &ApiKeyResponseInternal{}
	response, err := c.RequestBuilder().Endpoint(createApiKeyEndpoint).SetBody(&body).Post().ExecuteAndRetryOn429()
	if err != nil {
		return nil, err
	}
	if response.StatusCode != 200 {
		return nil, fmt.Errorf("api key creation failed. response status: %s", response.Status)
	}
	err = request.UnpackJSONResponse(response, &responseBody)
	if err != nil {
		return nil, err
	}
	return &responseBody.ApiKey, nil
}

func (c *Client) ReadApiKey(apiKey string) (*ApiKeyIamV2, error) {
	response, err := c.RequestBuilder().Endpoint(fmt.Sprintf(readApiKeyEndpoint, apiKey)).Get().ExecuteAndRetryOn429()
	if err != nil {
		return nil, err
	}
	if response.StatusCode == 403 { // Forbidden error happens when key isn't found
		return nil, nil
	}
	if response.StatusCode != 200 {
		return nil, fmt.Errorf("api key read failed. response status: %s", response.Status)
	}
	responseBody := &ApiKeyIamV2{}
	err = request.UnpackJSONResponse(response, &responseBody)
	if err != nil {
		return nil, err
	}
	return responseBody, nil
}

func (c *Client) UpdateApiKey(id, description, envId, resourceId string) error {
	body := &ApiKeyUpdateRequestW{
		&ApiKeyUpdateRequest{
			ID:              id,
			AccountID:       envId,
			LogicalClusters: []LogicalCluster{{ID: resourceId}},
			Description:     description,
		},
	}
	response, err := c.RequestBuilder().Endpoint(fmt.Sprintf(updateApiKeyEndpoint, id)).SetBody(&body).Put().ExecuteAndRetryOn429()
	if err != nil {
		return err
	}
	if response.StatusCode != 200 {
		return fmt.Errorf("delete request failed. response status: %s", response.Status)
	}
	return nil
}

func (c *Client) DeleteApiKey(id, envId, resourceId string) error {
	body := &ApiKeyDeleteRequestW{
		&ApiKeyDeleteRequest{
			ID:              id,
			AccountID:       envId,
			LogicalClusters: []LogicalCluster{{ID: resourceId}},
		},
	}
	response, err := c.RequestBuilder().Endpoint(fmt.Sprintf(deleteApiKeyEndpoint, id)).SetBody(&body).Delete().ExecuteAndRetryOn429()
	if err != nil {
		return err
	}
	if response.StatusCode == 403 { // Forbidden error happens when key isn't found
		return nil
	}
	if response.StatusCode != 200 {
		return fmt.Errorf("delete request failed. response status: %s", response.Status)
	}
	return nil
}
