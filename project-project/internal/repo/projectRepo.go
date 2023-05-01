package repo

import (
	"context"
	"test.com/project-project/internal/data/pro"
	"test.com/project-project/internal/database"
)

type ProjectRepo interface {
	FindProjectByMemId(ctx context.Context, memId int64, condition string, page int64, size int64) ([]*pro.ProjectAndMember, int64, error)
	FindCollectProjectByMemId(ctx context.Context, memberId int64, page int64, size int64) ([]*pro.ProjectAndMember, int64, error)
	SaveProject(conn database.DbConn, ctx context.Context, pr *pro.Project) error
	SaveProjectMember(conn database.DbConn, ctx context.Context, pm *pro.ProjectMember) error
	FindProjectByPIdAndMemId(ctx context.Context, projectCode int64, memberId int64) (*pro.ProjectAndMember, error)
	FindCollectByPidAndMemId(ctx context.Context, projectCode int64, memberId int64) (bool, error)
}
type ProjectTemplateRepo interface {
	FindProjectTemplateSystem(ctx context.Context, page int64, size int64) ([]pro.ProjectTemplate, int64, error)
	FindProjectTemplateCustom(ctx context.Context, memId int64, organizationCode int64, page int64, size int64) ([]pro.ProjectTemplate, int64, error)
	FindProjectTemplateAll(ctx context.Context, organizationCode int64, page int64, size int64) ([]pro.ProjectTemplate, int64, error)
}
