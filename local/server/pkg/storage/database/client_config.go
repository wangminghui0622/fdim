package database

import (
	"context"

	"fdim/pkg/db/pagination"
	"fdim/pkg/storage/model"
)

type ClientConfig interface {
	Set(ctx context.Context, userID string, config map[string]string) error
	Get(ctx context.Context, userID string) (map[string]string, error)
	Del(ctx context.Context, userID string, keys []string) error
	GetPage(ctx context.Context, userID string, key string, pagination pagination.Pagination) (int64, []*model.ClientConfig, error)
}
