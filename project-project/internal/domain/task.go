package domain

import (
	"context"
	"test.com/project-common/errs"
	"test.com/project-project/internal/dao"
	"test.com/project-project/internal/repo"
	"test.com/project-project/pkg/model"
)

type TaskDomain struct {
	taskRepo repo.TaskRepo
}

func NewTaskDomain() *TaskDomain {
	return &TaskDomain{
		taskRepo: dao.NewTaskDao(),
	}
}

func (d *TaskDomain) FindProjectIdByTaskId(taskId int64) (int64, bool, *errs.BError) {
	task, err := d.taskRepo.FindTaskById(context.Background(), taskId)
	if err != nil {
		return 0, false, model.DbError
	}
	if task == nil {
		return 0, false, nil
	}
	return task.ProjectCode, true, nil
}
