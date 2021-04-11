package remote

// import (
// 	"context"
// 	"net/http"
// )

// // Client _
// type Client struct {
// 	baseURL string

// 	httpClient HTTPClient
// }

// // HTTPClient _
// type HTTPClient interface {
// 	Do(req *http.Request) (*http.Response, error)
// }

// // NewClient _
// func NewClient() *Client {
// 	return &Client{
// 		baseURL: "http://localhost:8080",
// 	}
// }

// // NewRequest _
// func (rc *Client) NewRequest(
// 	method,
// 	path string,
// 	params map[string]string,
// 	reqBody interface{},
// 	needAuth bool) (*http.Request, error) {

// 	panic("not implemented")
// }

// // Do _
// func (rc *Client) Do(req *http.Request) (*http.Response, error) {
// 	return rc.httpClient.Do(req)
// }

// // Call _
// func (rc *Client) Call(
// 	ctx context.Context,
// 	method,
// 	path string,
// 	params map[string]string,
// 	reqBody interface{},
// 	needAuth bool,
// 	result interface{}) error {

// 	req, err := rc.NewRequest(method, path, params, reqBody, needAuth)

// 	if err != nil {
// 		return nil, errors.Wrap(err, "Building request")
// 	}

// 	req.WithContext(ctx)

// 	resp, err := rc.httpClient.Do(req)
// }

// // UnmarshalResponse _
// func (rc *Client) UnmarshalResponse(body []byte, val interface{}) error {
// 	err := json.Unmarshal(body, val)

// 	return err
// }
