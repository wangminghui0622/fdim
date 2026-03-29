package admin

import (
	"context"
	"fdim/pkg/db/pagination"
	"time"
)

type RegisterAddGroup struct {
	GroupID    string    `bson:"group_id"`
	CreateTime time.Time `bson:"create_time"`
}

func (RegisterAddGroup) TableName() string {
	return "register_add_groups"
}

type RegisterAddGroupInterface interface {
	Add(ctx context.Context, registerAddGroups []*RegisterAddGroup) error
	Del(ctx context.Context, groupIDs []string) error
	FindGroupID(ctx context.Context, groupIDs []string) ([]string, error)
	Search(ctx context.Context, keyword string, pagination pagination.Pagination) (int64, []*RegisterAddGroup, error)
}
