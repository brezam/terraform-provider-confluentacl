package request

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

type Request struct {
	UrlEndpoints []string
	authHeader   string

	body        interface{}
	method      string
	queryParams map[string]string
}

func NewRequestWithBasicAuth(baseUrl string, authUser string, authPassword string) *Request {
	basicToken := base64.StdEncoding.EncodeToString([]byte(authUser + ":" + authPassword))
	return &Request{UrlEndpoints: []string{baseUrl}, authHeader: "Basic " + basicToken}
}

func NewRequestWithBearerAuth(baseUrl string, jwtToken string) *Request {
	return &Request{UrlEndpoints: []string{baseUrl}, authHeader: "Bearer " + jwtToken}
}

func (r *Request) Endpoint(endpoints ...string) *Request {
	r.UrlEndpoints = append(r.UrlEndpoints, endpoints...)
	return r
}

func (r *Request) SetQueryParams(params map[string]string) *Request {
	r.queryParams = params
	return r
}

func (r *Request) SetBody(body interface{}) *Request {
	r.body = body
	return r
}

func (r *Request) Get() *Request {
	r.method = "GET"
	return r
}

func (r *Request) Post() *Request {
	r.method = "POST"
	return r
}

func (r *Request) Put() *Request {
	r.method = "PUT"
	return r
}

func (r *Request) Delete() *Request {
	r.method = "DELETE"
	return r
}

func UnpackJSONResponse(response *http.Response, container interface{}) error {
	defer response.Body.Close()
	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(bodyBytes, container)
	if err != nil {
		return err
	}
	return nil
}

func (r *Request) ExecuteAndRetryOn429() (*http.Response, error) {
	millisecondBackOffs := []int64{0, 100, 200, 400, 800, 1600, 3200}
	for _, wait := range millisecondBackOffs {
		time.Sleep(time.Duration(wait) * time.Millisecond)
		response, err := r.Execute()
		if err != nil {
			return nil, err
		}
		if response.StatusCode == 429 {
			continue
		}
		return response, nil
	}
	return nil, fmt.Errorf("exhausted retries")
}

func (r *Request) Execute() (*http.Response, error) {
	var bytesBody *bytes.Buffer
	if r.body != nil {
		jsonBody, err := json.Marshal(r.body)
		if err != nil {
			return nil, err
		}
		bytesBody = bytes.NewBuffer(jsonBody)
	}
	urlPath, err := r.resolveUrlEndpoints()
	if err != nil {
		return nil, err
	}
	var request *http.Request
	if bytesBody == nil {
		request, err = http.NewRequest(r.method, urlPath, nil)
	} else {
		request, err = http.NewRequest(r.method, urlPath, bytesBody)
	}
	if err != nil {
		return nil, err
	}
	request.Header.Add("Accept", "application/json")
	if bytesBody != nil {
		request.Header.Add("Content-Type", "application/json")
	}
	if len(r.queryParams) != 0 {
		q := request.URL.Query()
		for k, v := range r.queryParams {
			q.Add(k, v)
		}
		request.URL.RawQuery = q.Encode()
	}
	if r.authHeader != "" {
		request.Header.Add("Authorization", r.authHeader)
	}
	httpClient := &http.Client{}
	response, err := httpClient.Do(request)
	if response.StatusCode == 401 {
		return nil, errors.New("Unauthorized")
	}
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (r *Request) resolveUrlEndpoints() (string, error) {
	var endPointUrl, newPath *url.URL
	var err error
	for _, endpoint := range r.UrlEndpoints {
		newPath, err = url.Parse(endpoint)
		if err != nil {
			return "", err
		}
		endPointUrl = endPointUrl.ResolveReference(newPath)
	}
	return endPointUrl.String(), nil
}
