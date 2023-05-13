package repo

import (
	"context"
	"test.com/project-project/internal/data"
)

type ProjectAuthRepo interface {
	FindAuthList(ctx context.Context, orgCode int64) (list []*data.ProjectAuth, err error)
}
