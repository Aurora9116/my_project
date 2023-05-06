package dao

import (
	"context"
	"test.com/project-project/internal/data"
	"test.com/project-project/internal/database"
	"test.com/project-project/internal/database/gorms"
)

type TaskStagesDao struct {
	conn *gorms.GormConn
}

func (t TaskStagesDao) FindStagesByProjectId(ctx context.Context, projectCode int64, page, pageSize int64) (list []*data.TaskStages, total int64, err error) {
	session := t.conn.Default(ctx)
	err = session.Model(&data.TaskStages{}).Where("project_code=?", projectCode).Order("sort asc").Limit(int(pageSize)).Offset(int((page - 1) * pageSize)).Find(&list).Error
	err = session.Model(&data.TaskStages{}).Where("project_code=?", projectCode).Count(&total).Error
	return
}
func (t TaskStagesDao) SaveTaskStages(ctx context.Context, conn database.DbConn, ts *data.TaskStages) error {
	t.conn = conn.(*gorms.GormConn)
	err := t.conn.Tx(ctx).Save(&ts).Error
	return err
}

func NewTaskStagesDao() *TaskStagesDao {
	return &TaskStagesDao{conn: gorms.New()}
}
