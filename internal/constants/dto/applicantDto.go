package dto

import "encoding/json"

type JSONRawMessage = json.RawMessage

type GetApplicationsResponse struct {
	ID        string `json:"id"`
	JobPostID string `json:"jobPostID"`

	Name           string         `json:"name"`
	Email          string         `json:"email"`
	S3PATH         string         `json:"s3_path"`
	Status         string         `json:"status"`
	SeniorityLevel string         `json:"seniority_level"`
	ProgressStatus string         `json:"progress_status"`
	Analysis       JSONRawMessage `json:"analysis" swaggertype:"object"`

	CreatedAt string `json:"created_at"`
}

type EnvelopeGetApplicationsResponse struct {
	Data  []GetApplicationsResponse `json:"data"`
	Error *ErrorResponse            `json:"error,omitempty"`
}
