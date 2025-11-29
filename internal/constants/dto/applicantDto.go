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

type GetMetaDataApplicationsResponse struct {
	TotalCount        int `json:"total_count"`
	AppliedCount      int `json:"applied_count"`
	InterviewingCount int `json:"interviewing_count"`
	RejectedCount     int `json:"rejected_count"`
	HiredCount        int `json:"hired_count"`
	ShortlistedCount  int `json:"shortlisted_count"`
}

type GetApplicationsResponseWithMetaData struct {
	Applications []GetApplicationsResponse       `json:"applications"`
	MetaData     GetMetaDataApplicationsResponse `json:"meta_data"`
}

type EnvelopeGetApplicationsResponse struct {
	Data  GetApplicationsResponseWithMetaData `json:"data"`
	Error *ErrorResponse                      `json:"error,omitempty"`
}

type UpdateApplicationResponse struct {
	ID             string `json:"id"`
	ProgressStatus string `json:"progress_status"`
}

type EnvelopeUpdateApplicationResponse struct {
	Data  UpdateApplicationResponse `json:"data"`
	Error *ErrorResponse            `json:"error,omitempty"`
}
