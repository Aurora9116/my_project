package dao

import (
	"context"
	"test.com/project-project/internal/data"
	"test.com/project-project/internal/database/gorms"
)

type ProjectLogDao struct {
	conn *gorms.GormConn
}

func (p *ProjectLogDao) SaveProjectLog(pl *data.ProjectLog) {
	session := p.conn.Default(context.Background())
	session.Save(&pl)
}

func (p *ProjectLogDao) FindLogByTaskCode(ctx context.Context, taskCode int64, comment int) (list []*data.ProjectLog, total int64, err error) {
	session := p.conn.Default(ctx)
	model := session.Model(&data.ProjectLog{})
	if comment == 1 {
		model.Where("source_code=? and is_comment=?", taskCode, comment).Find(&list)
		model.Where("source_code=? and is_comment=?", taskCode, comment).Count(&total)
	} else {
		model.Where("source_code=?", taskCode).Find(&list)
		model.Where("source_code=?", taskCode).Count(&total)
	}
	return
}

func (p *ProjectLogDao) FindLogByTaskCodePage(ctx context.Context, taskCode int64, comment int, page int, pageSize int) (list []*data.ProjectLog, total int64, err error) {
	session := p.conn.Default(ctx)
	model := session.Model(&data.ProjectLog{})
	offset := (page - 1) * pageSize
	if comment == 1 {
		model.Where("source_code=? and is_comment=?", taskCode, comment).Limit(pageSize).Offset(offset).Find(&list)
		model.Where("source_code=? and is_comment=?", taskCode, comment).Count(&total)
	} else {
		model.Where("source_code=?", taskCode).Limit(pageSize).Offset(offset).Find(&list)
		model.Where("source_code=?", taskCode).Count(&total)
	}
	return
}

func NewProjectLogDao() *ProjectLogDao {
	return &ProjectLogDao{
		conn: gorms.New(),
	}
}
