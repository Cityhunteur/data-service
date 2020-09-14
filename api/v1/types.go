// Package v1 contains types used by the API.
package v1

// Data represents a request for the API.
type Data struct {
	Title string `json:"title"`
}

// DataResponse represents the response for the API.
type DataResponse struct {
	ID        string   `json:"id"`
	Title     string   `json:"title"`
	Timestamp Time3339 `json:"timestamp"`
}
