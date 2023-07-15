NAME=confluentacl
BINARY=terraform-provider-${NAME}
VERSION=0.1.0
GOOS=$(shell go env GOOS)
GOARCH=$(shell go env GOARCH)
OS_ARCH=${GOOS}_${GOARCH}

.PHONY: default build testacc .validate-testing-env-vars

default: build

.validate-testing-env-vars:
	ifndef TEST_ENV_ID
		$(error TEST_ENV_ID is undefined. You need to set to a real confluent cloud environment id)
	endif
	ifndef TEST_CLUSTER_ID
		$(error TEST_CLUSTER_ID is undefined. You need to set to a real confluent cloud cluster id that matches TEST_REST_ENDPOINT)
	endif
	ifndef TEST_SERVICE_ACCOUNT_NAME
		$(error TEST_SERVICE_ACCOUNT_NAME is undefined. You need to set to a real confluent cloud service account name)
	endif
	ifndef TEST_REST_ENDPOINT
		$(error TEST_REST_ENDPOINT is undefined. You need to set to a real confluent cloud rest endpoint that matches TEST_CLUSTER_ID)
	endif

build:
	go build -o ${BINARY}

testacc: .validate-testing-env-vars
	TF_ACC=1 go test -v ./...
