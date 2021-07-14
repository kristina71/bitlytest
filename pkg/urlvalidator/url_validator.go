package urlvalidator

import (
	"context"
	"net/http"
)

func ValidateUrl(ctx context.Context, url string) bool {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return false
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return false
	}

	defer resp.Body.Close()
	return resp.StatusCode == http.StatusOK
}
