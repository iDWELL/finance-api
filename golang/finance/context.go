package finance

import (
	"context"
	"net/http"
)

type (
	ctxBaseURL    struct{}
	ctxHTTPClient struct{}
)

// WithBaseURL creates context.Context with baseURL param for making requests
func WithBaseURL(ctx context.Context, baseURL string) context.Context {
	return context.WithValue(ctx, ctxBaseURL{}, baseURL)
}

// WithHTTPClient creates context.Context with [*http.Client] for making requests
func WithHTTPClient(ctx context.Context, client *http.Client) context.Context {
	return context.WithValue(ctx, ctxHTTPClient{}, client)
}

func baseURLFromCtx(ctx context.Context) string {
	baseURL, ok := ctx.Value(ctxBaseURL{}).(string)
	if !ok || baseURL == "" {
		return BaseURL
	}

	return baseURL
}

func httpClientFromCtx(ctx context.Context) *http.Client {
	httpClient, ok := ctx.Value(ctxHTTPClient{}).(*http.Client)
	if !ok || httpClient == nil {
		return http.DefaultClient
	}

	return httpClient
}
