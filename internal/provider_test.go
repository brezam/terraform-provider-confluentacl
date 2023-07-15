package internal

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

const (
	envVarTestEnvId        = "TEST_ENV_ID"
	envVarTestClusterId    = "TEST_CLUSTER_ID"
	envVarTestSaName       = "TEST_SERVICE_ACCOUNT_NAME"
	envVarTestRestEndpoint = "TEST_REST_ENDPOINT"
)

var (
	testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
		"example": providerserver.NewProtocol6WithError(NewProvider()),
	}
	testRealResource = &TestRealResources{
		EnvId:        os.Getenv(envVarTestEnvId),
		ClusterId:    os.Getenv(envVarTestClusterId),
		SaName:       os.Getenv(envVarTestSaName),
		RestEndpoint: os.Getenv(envVarTestRestEndpoint),
	}
)

type TestRealResources struct {
	EnvId        string
	ClusterId    string
	SaName       string
	RestEndpoint string
}

func testAccPreCheck(t *testing.T) {
	// Provider configuration env vars CONFLUENT_API_KEY and CONFLUENT_API_SECRET
	if v := os.Getenv(envVarCloudApiKey); v == "" {
		t.Fatal(envVarCloudApiKey + " must be set for acceptance tests")
	}
	if v := os.Getenv(envVarCloudApiSecret); v == "" {
		t.Fatal(envVarCloudApiSecret + " must be set for acceptance tests")
	}
	// Values of environment id, cluster id and service account name for real resources in Confluent.
	// This provider doesn't create them, so we already need them existing in order to test api key and acl creation
	if v := os.Getenv(envVarTestEnvId); v == "" {
		t.Fatal(envVarTestEnvId + " must be set for acceptance tests")
	}
	if v := os.Getenv(envVarTestClusterId); v == "" {
		t.Fatal(envVarTestClusterId + " must be set for acceptance tests")
	}
	if v := os.Getenv(envVarTestSaName); v == "" {
		t.Fatal(envVarTestSaName + " must be set for acceptance tests")
	}
	if v := os.Getenv(envVarTestRestEndpoint); v == "" {
		t.Fatal(envVarTestRestEndpoint + " must be set for acceptance tests")
	}
}
