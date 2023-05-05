package dao

import (
	"context"
	"test.com/project-project/internal/data/task"
	"test.com/project-project/internal/database"
	"test.com/project-project/internal/database/gorms"
)

type TaskStagesDao struct {
	conn *gorms.GormConn
}

func (t TaskStagesDao) SaveTaskStages(ctx context.Context, conn database.DbConn, ts *task.TaskStages) error {
	t.conn = conn.(*gorms.GormConn)
	err := t.conn.Tx(ctx).Save(&ts).Error
	return err
}

func NewTaskStagesDao() *TaskStagesDao {
	return &TaskStagesDao{conn: gorms.New()}
}
