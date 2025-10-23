package dto

type Envelope[T any] struct {
	Data  T      `json:"data,omitempty"`
	Error string `json:"error,omitempty"`
}

type ErrorResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

// --- Swagger helper envelope types (doc-only) ---

type EnvelopeAny struct {
	Data  interface{} `json:"data,omitempty"`
	Error string      `json:"error,omitempty"`
}

type EnvelopeCreateUserResponse struct {
	Data  CreateUserResponse `json:"data,omitempty"`
	Error string             `json:"error,omitempty"`
}

type EnvelopeLoginUserResponse struct {
	Data  LoginUserResponse `json:"data,omitempty"`
	Error string            `json:"error,omitempty"`
}

type EnvelopeGetAllUsersResponse struct {
	Data  []GetAllUsers `json:"data,omitempty"`
	Error string        `json:"error,omitempty"`
}

type EnvelopeCreateJobPostResponse struct {
	Data  CreateJobPostResponse `json:"data,omitempty"`
	Error string                `json:"error,omitempty"`
}

type EnvelopeGetAllJobPostsResponse struct {
	Data  []GetAllJobPostsResponse `json:"data,omitempty"`
	Error string                   `json:"error,omitempty"`
}

type EnvelopeUpdateJobPostResponse struct {
	Data  UpdateJobPostResponse `json:"data,omitempty"`
	Error string                `json:"error,omitempty"`
}
