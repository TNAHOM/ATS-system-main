package dto

type PaginationRequest struct {
	Page int `form:"page"`
	Size int `form:"size"`
}

type PaginationMeta struct {
	Total *int `json:"total,omitempty"` // pointer so we can omit when unknown
	Page  int  `json:"page"`
	Size  int  `json:"size"`
}

type PaginatedResponse[T any] struct {
	Items []T            `json:"items"`
	Meta  PaginationMeta `json:"meta"`
}

type PaginatedGetJobPostsResponse struct {
	Items []GetJobPostsResponse `json:"items"`
	Meta  PaginationMeta        `json:"meta"`
}
