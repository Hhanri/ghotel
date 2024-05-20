package api_util

type ResourceResponse[T any] struct {
	Data    []T `json:"data"`
	Results int `json:"results"`
	Page    int `json:"page"`
}

func NewResourceResponse[T any](data []T, page int) *ResourceResponse[T] {
	return &ResourceResponse[T]{
		Data:    data,
		Results: len(data),
		Page:    page,
	}
}
