# 11.4 Reference implementation: Casdoor

Casdoor is the reference layout for every integration. It mirrors the Laravel package while applying Go conventions.

## File mapping

| Laravel                                                                       | Go                                    |
| ----------------------------------------------------------------------------- | ------------------------------------- |
| `app/Infrastructure/Integrations/Casdoor/CasdoorClient.php`                   | `client.go`                           |
| `app/Infrastructure/Integrations/Casdoor/Contracts/ClientInterface.php`       | `contracts/client.go`                 |
| `app/Infrastructure/Integrations/Casdoor/Clients/UserClient.php`              | `clients/user_client.go`              |
| `app/Infrastructure/Integrations/Casdoor/DTOs/Requests/UpdateUserRequest.php` | `dto/requests/update_user_request.go` |
| `app/Infrastructure/Integrations/Casdoor/DTOs/Responses/GetUserResponse.php`  | `dto/responses/get_user_response.go`  |
| `app/Infrastructure/Integrations/Casdoor/Exceptions/CasdoorException.php`     | `exceptions/casdoor_exception.go`     |
| `app/Infrastructure/Integrations/Casdoor/Support/Endpoint.php`                | `support/endpoint.go`                 |

## Layout

```text
app/infrastructure/integrations/casdoor/
├── client.go
├── client_test.go
├── contracts/
│   └── client.go
├── clients/
│   ├── user_client.go
│   └── user_client_test.go
├── dto/
│   ├── requests/
│   │   └── update_user_request.go
│   └── responses/
│       └── get_user_response.go
├── exceptions/
│   └── casdoor_exception.go
└── support/
    └── endpoint.go
```

## `contracts/client.go`

```go
// app/infrastructure/integrations/casdoor/contracts/client.go
package contracts

import (
	"context"
	"net/url"
)

type Client interface {
	Get(ctx context.Context, path string, query url.Values) ([]byte, error)
	Post(ctx context.Context, path string, payload any) ([]byte, error)
	Put(ctx context.Context, path string, payload any) ([]byte, error)
	Patch(ctx context.Context, path string, payload any) ([]byte, error)
	Delete(ctx context.Context, path string) ([]byte, error)
}
```

The interface returns raw bytes because the transport adapter must not know response DTOs. Resource clients decode into their own DTOs.

## `client.go`

```go
// app/infrastructure/integrations/casdoor/client.go
package casdoor

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"go-fiber-vue/app/infrastructure/integrations/casdoor/exceptions"
)

type Config struct {
	BaseURL      string
	ClientID     string
	ClientSecret string
}

type Client struct {
	baseURL      string
	clientID     string
	clientSecret string
	httpClient   *http.Client
}

func NewClient(cfg Config, httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = &http.Client{Timeout: 15 * time.Second}
	}

	return &Client{
		baseURL:      strings.TrimRight(cfg.BaseURL, "/"),
		clientID:     cfg.ClientID,
		clientSecret: cfg.ClientSecret,
		httpClient:   httpClient,
	}
}

func (c *Client) Get(ctx context.Context, path string, query url.Values) ([]byte, error) {
	if query == nil {
		query = url.Values{}
	}
	query.Set("clientId", c.clientID)
	query.Set("clientSecret", c.clientSecret)

	return c.do(ctx, http.MethodGet, path, query, nil)
}

func (c *Client) Post(ctx context.Context, path string, payload any) ([]byte, error) {
	return c.do(ctx, http.MethodPost, path, nil, payload)
}

func (c *Client) Put(ctx context.Context, path string, payload any) ([]byte, error) {
	return c.do(ctx, http.MethodPut, path, nil, payload)
}

func (c *Client) Patch(ctx context.Context, path string, payload any) ([]byte, error) {
	return c.do(ctx, http.MethodPatch, path, nil, payload)
}

func (c *Client) Delete(ctx context.Context, path string) ([]byte, error) {
	return c.do(ctx, http.MethodDelete, path, nil, nil)
}

func (c *Client) do(ctx context.Context, method, path string, query url.Values, payload any) ([]byte, error) {
	endpoint := c.baseURL + "/" + strings.TrimLeft(path, "/")
	if len(query) > 0 {
		endpoint += "?" + query.Encode()
	}

	var body io.Reader
	if payload != nil {
		encoded, err := json.Marshal(payload)
		if err != nil {
			return nil, fmt.Errorf("casdoor: encode request: %w", err)
		}
		body = bytes.NewReader(encoded)
	}

	request, err := http.NewRequestWithContext(ctx, method, endpoint, body)
	if err != nil {
		return nil, fmt.Errorf("casdoor: build request: %w", err)
	}
	request.Header.Set("Accept", "application/json")
	if payload != nil {
		request.Header.Set("Content-Type", "application/json")
	}

	response, err := c.httpClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("casdoor: request failed: %w", err)
	}
	defer response.Body.Close()

	responseBody, err := io.ReadAll(io.LimitReader(response.Body, 1<<20))
	if err != nil {
		return nil, fmt.Errorf("casdoor: read response: %w", err)
	}

	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		return nil, exceptions.New("casdoor request failed", responseBody, response.StatusCode)
	}

	return responseBody, nil
}
```

## `clients/user_client.go`

```go
// app/infrastructure/integrations/casdoor/clients/user_client.go
package clients

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"go-fiber-vue/app/infrastructure/integrations/casdoor/contracts"
	"go-fiber-vue/app/infrastructure/integrations/casdoor/dto/requests"
	"go-fiber-vue/app/infrastructure/integrations/casdoor/dto/responses"
	"go-fiber-vue/app/infrastructure/integrations/casdoor/exceptions"
	"go-fiber-vue/app/infrastructure/integrations/casdoor/support"
)

type UserClient struct {
	client contracts.Client
	owner  string
}

func NewUserClient(client contracts.Client, owner string) *UserClient {
	return &UserClient{client: client, owner: owner}
}

func (c *UserClient) GetUser(ctx context.Context, name string, query url.Values) (*responses.GetUserResponse, error) {
	if query == nil {
		query = url.Values{}
	}
	query.Set("id", c.owner+"/"+name)

	body, err := c.client.Get(ctx, support.EndpointGetUser, query)
	if err != nil {
		return nil, fmt.Errorf("casdoor: get user %q: %w", name, err)
	}

	var response responses.GetUserResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("casdoor: decode get-user response: %w", err)
	}

	return &response, nil
}

func (c *UserClient) UpdateUser(ctx context.Context, request requests.UpdateUserRequest) error {
	body, err := c.client.Post(ctx, support.EndpointUpdateUser, request)
	if err != nil {
		return fmt.Errorf("casdoor: update user %q: %w", request.Name, err)
	}

	var response struct {
		Status string `json:"status"`
		Msg    string `json:"msg"`
	}
	if err := json.Unmarshal(body, &response); err != nil {
		return fmt.Errorf("casdoor: decode update-user response: %w", err)
	}
	if response.Status != "ok" {
		return exceptions.New(
			"failed to update Casdoor user ["+request.Name+"]",
			body,
			http.StatusBadGateway,
		)
	}

	return nil
}
```

## `dto/requests/update_user_request.go`

```go
// app/infrastructure/integrations/casdoor/dto/requests/update_user_request.go
package requests

type UpdateUserRequest struct {
	Owner       string         `json:"owner"`
	Name        string         `json:"name"`
	ID          string         `json:"id"`
	DisplayName string         `json:"displayName"`
	Avatar      string         `json:"avatar"`
	Email       string         `json:"email"`
	Phone       string         `json:"phone"`
	Region      string         `json:"region"`
	Location    string         `json:"location"`
	Title       string         `json:"title"`
	Properties  map[string]any `json:"properties"`
}
```

## `dto/responses/get_user_response.go`

```go
// app/infrastructure/integrations/casdoor/dto/responses/get_user_response.go
package responses

type GetUserResponse struct {
	Data  map[string]any `json:"data"`
	Data2 map[string]any `json:"data2"`
}
```

## `exceptions/casdoor_exception.go`

```go
// app/infrastructure/integrations/casdoor/exceptions/casdoor_exception.go
package exceptions

type Exception struct {
	message    string
	response   []byte
	statusCode int
}

func New(message string, response []byte, statusCode int) *Exception {
	return &Exception{
		message:    message,
		response:   response,
		statusCode: statusCode,
	}
}

func (e *Exception) Error() string {
	return e.message
}

func (e *Exception) Response() []byte {
	return e.response
}

func (e *Exception) StatusCode() int {
	return e.statusCode
}
```

## `support/endpoint.go`

```go
// app/infrastructure/integrations/casdoor/support/endpoint.go
package support

const (
	EndpointGetUser    = "api/get-user"
	EndpointUpdateUser = "api/update-user"
)
```

## Wiring

```go
import (
	"net"
	"net/http"
	"time"

	casdoor "go-fiber-vue/app/infrastructure/integrations/casdoor"
	casdoorclients "go-fiber-vue/app/infrastructure/integrations/casdoor/clients"
)

httpClient := &http.Client{
	Timeout: 15 * time.Second,
	Transport: &http.Transport{
		DialContext: (&net.Dialer{Timeout: 5 * time.Second}).DialContext,
	},
}

casdoorClient := casdoor.NewClient(casdoor.Config{
	BaseURL:      cfg.Casdoor.URL,
	ClientID:     cfg.Casdoor.ClientID,
	ClientSecret: cfg.Casdoor.ClientSecret,
}, httpClient)

casdoorUsers := casdoorclients.NewUserClient(casdoorClient, cfg.Casdoor.Owner)
```

Casdoor uses `CASDOOR_URL`, `CASDOOR_CLIENT_ID`, `CASDOOR_CLIENT_SECRET`, and `CASDOOR_OWNER`.

## Error mapping

```go
func writeCasdoorError(ctx *fiber.Ctx, err error) error {
	var casdoorErr *exceptions.Exception
	if !errors.As(err, &casdoorErr) {
		return helpers.Error(ctx, err.Error(), nil, fiber.StatusInternalServerError)
	}

	var payload any
	if len(casdoorErr.Response()) > 0 {
		_ = json.Unmarshal(casdoorErr.Response(), &payload)
	}

	return helpers.Error(ctx, casdoorErr.Error(), payload, casdoorErr.StatusCode())
}
```

## Test with a fake contract

```go
// app/infrastructure/integrations/casdoor/clients/user_client_test.go
package clients_test

import (
	"context"
	"errors"
	"net/url"
	"testing"

	"go-fiber-vue/app/infrastructure/integrations/casdoor/clients"
)

type fakeClient struct {
	get func(context.Context, string, url.Values) ([]byte, error)
}

func (f fakeClient) Get(ctx context.Context, path string, query url.Values) ([]byte, error) {
	return f.get(ctx, path, query)
}

func (f fakeClient) Post(context.Context, string, any) ([]byte, error) {
	return nil, errors.New("not implemented")
}

func (f fakeClient) Put(context.Context, string, any) ([]byte, error) {
	return nil, errors.New("not implemented")
}

func (f fakeClient) Patch(context.Context, string, any) ([]byte, error) {
	return nil, errors.New("not implemented")
}

func (f fakeClient) Delete(context.Context, string) ([]byte, error) {
	return nil, errors.New("not implemented")
}

func TestUserClientGetUserAddsOwnerToID(t *testing.T) {
	var gotQuery url.Values
	client := clients.NewUserClient(fakeClient{
		get: func(_ context.Context, _ string, query url.Values) ([]byte, error) {
			gotQuery = query
			return []byte(`{"data":{"name":"alice"}}`), nil
		},
	}, "sigma")

	response, err := client.GetUser(context.Background(), "alice", nil)
	if err != nil {
		t.Fatalf("GetUser() error = %v", err)
	}
	if got := gotQuery.Get("id"); got != "sigma/alice" {
		t.Fatalf("id = %q, want %q", got, "sigma/alice")
	}
	if response.Data["name"] != "alice" {
		t.Fatalf("data.name = %v, want alice", response.Data["name"])
	}
}
```

## Design notes

- `contracts.Client` is an adapter-local port. It keeps `clients.UserClient` independent of the transport and lets tests use a fake instead of an HTTP server.
- Laravel `UserClient` currently depends on the concrete `CasdoorClient`; the Go version injects `contracts.Client`.
- Laravel `getUser` and `updateUser` return raw JSON and do not use the DTOs; the Go version decodes into `GetUserResponse` and accepts `UpdateUserRequest`.
- `CasdoorException` becomes a typed error that preserves the upstream status and raw response body. The handler maps it to the same JSON error envelope as the rest of the API.
- Casdoor authenticates with `clientId` and `clientSecret` query parameters, so there is no `token_manager.go`. Add a token manager only for integrations that require bearer-token exchange.
- Do not log the full GET URL: it contains `clientSecret`. Do not log request bodies that contain credentials.
- The 1 MiB response cap is sufficient for Casdoor user payloads. Raise it per integration if an endpoint returns larger documents.
