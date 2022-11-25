package apiKey

import "context"

type Repository interface {
	FindeApiKey(ctx context.Context,  key string) (ApiKeyResult, error)
	Delete(ctx context.Context, id string) error
}
