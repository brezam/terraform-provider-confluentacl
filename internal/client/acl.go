package client

import (
	"fmt"
	"terraform-provider-confluentacl/internal/client/request"
)

const (
	kafkaAclEndpoint = "kafka/v3/clusters/%s/acls"
	AclPatternTypeLiteral
)

type ACLRequest struct {
	ResourceType string `json:"resource_type"`
	ResourceName string `json:"resource_name"`
	PatternType  string `json:"pattern_type"`
	Principal    string `json:"principal"`
	Host         string `json:"host"`
	Operation    string `json:"operation"`
	Permission   string `json:"permission"`
}

type ACLListResponseWrapper struct {
	Kind string            `json:"kind"`
	Data []ACLListResponse `json:"data"`
}

type ACLListResponse struct {
	Kind         string `json:"kind"`
	ClusterId    string `json:"cluster_id"`
	ResourceType string `json:"resource_type"`
	ResourceName string `json:"resource_name"`
	PatternType  string `json:"pattern_type"`
	Principal    string `json:"principal"`
	Host         string `json:"host"`
	Operation    string `json:"operation"`
	Permission   string `json:"permission"`
}

func (c *Client) ListACLs(restEndpoint, clusterId string) ([]ACLListResponse, error) {
	return c.ListSpecificACLs(restEndpoint, clusterId, nil)
}

func (c *Client) ListSpecificACLs(restEndpoint, clusterId string, query *ACLRequest) ([]ACLListResponse, error) {
	var allAclsInCluster ACLListResponseWrapper
	endpoint := fmt.Sprintf(kafkaAclEndpoint, clusterId)
	requestBuilder, err := c.KafkaRestRequestBuilder(restEndpoint)
	if err != nil {
		return nil, err
	}
	var queryParams map[string]string
	if query == nil {
		queryParams = make(map[string]string, 0)
	} else {
		queryParams = map[string]string{
			"principal":     query.Principal,
			"resource_name": query.ResourceName,
			"resource_type": query.ResourceType,
			"pattern_type":  query.PatternType,
			"host":          query.Host,
			"operation":     query.Operation,
			"permission":    query.Permission,
		}
	}
	response, err := requestBuilder.
		Endpoint(endpoint).
		SetQueryParams(queryParams).
		Get().
		ExecuteAndRetryOn429()
	if err != nil {
		return nil, err
	}
	err = request.UnpackJSONResponse(response, &allAclsInCluster)
	if err != nil {
		return nil, err
	}
	return allAclsInCluster.Data, nil
}

func (c *Client) CreateACL(restEndpoint, clusterId string, request *ACLRequest) error {
	endpoint := fmt.Sprintf(kafkaAclEndpoint, clusterId)
	requestBuilder, err := c.KafkaRestRequestBuilder(restEndpoint)
	if err != nil {
		return err
	}
	response, err := requestBuilder.
		Endpoint(endpoint).
		SetBody(request).
		Post().
		ExecuteAndRetryOn429()
	if err != nil {
		return err
	}
	if response.StatusCode != 201 {
		return fmt.Errorf("Create ACL Failure. Response " + response.Status)
	}
	return nil
}

func (c *Client) DeleteAcl(restEndpoint, clusterId string, query *ACLRequest) error {
	endpoint := fmt.Sprintf(kafkaAclEndpoint, clusterId)
	requestBuilder, err := c.KafkaRestRequestBuilder(restEndpoint)
	if err != nil {
		return err
	}
	queryParams := map[string]string{
		"principal":     query.Principal,
		"resource_name": query.ResourceName,
		"resource_type": query.ResourceType,
		"pattern_type":  query.PatternType,
		"host":          query.Host,
		"operation":     query.Operation,
		"permission":    query.Permission,
	}
	response, err := requestBuilder.
		Endpoint(endpoint).
		SetQueryParams(queryParams).
		Delete().
		ExecuteAndRetryOn429()
	if err != nil {
		return err
	}
	if response.StatusCode != 201 {
		return fmt.Errorf("Create ACL Failure. Response " + response.Status)
	}
	return nil
}
