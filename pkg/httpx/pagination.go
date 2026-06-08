package httpx

import (
	"net/url"
	"strconv"
)

// Pagination default and bound constants for list endpoints.
const (
	DefaultLimit = 20
	MinLimit     = 1
	MaxLimit     = 100
)

// ParsePagination reads "limit" and "offset" from the query, applying defaults
// and clamping limit to [MinLimit, MaxLimit] and offset to >= 0. Invalid or
// missing values fall back to the defaults.
func ParsePagination(q url.Values) (limit, offset int) {
	limit = DefaultLimit
	if v := q.Get("limit"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			limit = n
		}
	}
	if limit < MinLimit {
		limit = MinLimit
	}
	if limit > MaxLimit {
		limit = MaxLimit
	}

	if v := q.Get("offset"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n >= 0 {
			offset = n
		}
	}
	return limit, offset
}

// Page is the envelope returned by paginated list endpoints.
type Page struct {
	Data   interface{} `json:"data"`
	Total  int         `json:"total"`
	Limit  int         `json:"limit"`
	Offset int         `json:"offset"`
}
