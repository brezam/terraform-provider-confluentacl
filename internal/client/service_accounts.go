package client

import (
	"fmt"
	"terraform-provider-confluentacl/internal/client/request"
)

type ServiceAccount struct {
	Id          string `json:"resource_id"`
	UserId      int    `json:"id"`
	ServiceName string `json:"service_name"`
}

type ServiceAccountResponse struct {
	Users []ServiceAccount `json:"users"`
}

const (
	serviceAccountsEndpoint = "service_accounts"
)

func (c *Client) ListServiceAccounts() ([]ServiceAccount, error) {
	// Check if data is in cache
	c.cacheMutex.RLock()
	serviceAccountsCache, ok := c.cache["serviceAccounts"]
	c.cacheMutex.RUnlock()

	if ok {
		// Return data from cache
		return serviceAccountsCache.([]ServiceAccount), nil
	}
	serviceAccounts := &ServiceAccountResponse{}
	response, err := c.RequestBuilder().Endpoint(serviceAccountsEndpoint).Get().ExecuteAndRetryOn429()
	if err != nil {
		return nil, err
	}
	if response.StatusCode != 200 {
		return nil, fmt.Errorf("listing service accounts failed. response status: %s", response.Status)
	}
	err = request.UnpackJSONResponse(response, &serviceAccounts)
	if err != nil {
		return nil, err
	}
	// Store data in cache
	c.cacheMutex.Lock()
	c.cache["serviceAccounts"] = serviceAccounts.Users
	c.cacheMutex.Unlock()

	return serviceAccounts.Users, nil
}

func (c *Client) GetSaNumericId(saName string) (int, error) {
	serviceAccountList, err := c.ListServiceAccounts()
	if err != nil {
		return 0, err
	}
	var numericUserId int
	for _, serviceAccount := range serviceAccountList {
		if serviceAccount.ServiceName == saName {
			numericUserId = serviceAccount.UserId
		}
	}
	return numericUserId, nil
}
