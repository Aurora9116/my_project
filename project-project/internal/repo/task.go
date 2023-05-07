package repo

import (
	"context"
	"test.com/project-project/internal/data"
	"test.com/project-project/internal/database"
)

type TaskStagesTemplateRepo interface {
	FindInProTemIds(ctx context.Context, id []int) ([]data.MsTaskStagesTemplate, error)
	FindByProjectTemplateId(ctx context.Context, projectTemplateCode int) (list []*data.MsTaskStagesTemplate, err error)
}
type TaskStagesRepo interface {
	SaveTaskStages(ctx context.Context, conn database.DbConn, ts *data.TaskStages) error
	FindStagesByProjectId(ctx context.Context, projectCode int64, page, pageSize int64) (list []*data.TaskStages, total int64, err error)
	FindById(ctx context.Context, id int) (ts *data.TaskStages, err error)
}
type TaskRepo interface {
	FindTaskByStageCode(ctx context.Context, stageCode int) (taskList []*data.Task, err error)
	FindTaskMemberByTaskId(ctx context.Context, taskCode int64, memberCode int64) (task *data.TaskMember, err error)
	FindTaskMaxIdNum(ctx context.Context, projectCode int64) (v *int, err error)
	FindTaskSort(ctx context.Context, projectCode int64, stageCode int64) (v *int, err error)
	SaveTask(ctx context.Context, conn database.DbConn, ts *data.Task) error
	SaveTaskMember(ctx context.Context, conn database.DbConn, tm *data.TaskMember) error
	FindTaskById(ctx context.Context, preTaskCode int64) (ts *data.Task, err error)
	UpdateTaskSort(ctx context.Context, conn database.DbConn, ts *data.Task) error
	FindTaskByStageCodeLtSort(ctx context.Context, stageCode int, sort int) (ts *data.Task, err error)
}
