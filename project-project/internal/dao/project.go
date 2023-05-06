package dao

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"test.com/project-project/internal/data/pro"
	"test.com/project-project/internal/database"
	"test.com/project-project/internal/database/gorms"
)

type ProjectDao struct {
	conn *gorms.GormConn
}

func (p *ProjectDao) FindProjectById(ctx context.Context, projectCode int64) (pj *pro.Project, err error) {
	err = p.conn.Default(ctx).Where("id=?", projectCode).Find(&pj).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return
}

func (p *ProjectDao) FindProjectMemberByPid(ctx context.Context, projectCode int64) (list []*pro.ProjectMember, total int64, err error) {
	session := p.conn.Default(ctx)
	err = session.Model(&pro.ProjectMember{}).Where("project_code=?", projectCode).Find(&list).Error
	err = session.Model(&pro.ProjectMember{}).Where("project_code=?", projectCode).Count(&total).Error
	return
}

func (p *ProjectDao) UpdateProject(ctx context.Context, proj *pro.Project) error {
	return p.conn.Default(ctx).Updates(&proj).Error
}

func (p *ProjectDao) DeleteProjectCollect(ctx context.Context, memId int64, projectCode int64) error {
	return p.conn.Default(ctx).Where("member_code = ? and project_code = ?", memId, projectCode).Delete(&pro.ProjectCollection{}).Error
}

func (p *ProjectDao) SaveProjectCollect(ctx context.Context, pc *pro.ProjectCollection) error {
	return p.conn.Default(ctx).Save(&pc).Error
}

func (p *ProjectDao) UpdateDeletedProject(ctx context.Context, projectCode int64, deleted bool) error {
	session := p.conn.Default(ctx)
	if deleted {
		return session.Model(&pro.Project{}).Where("id=?", projectCode).Update("deleted", 1).Error
	}
	return session.Model(&pro.Project{}).Where("id=?", projectCode).Update("deleted", 0).Error
}

func (p *ProjectDao) FindCollectByPidAndMemId(ctx context.Context, projectCode int64, memberId int64) (bool, error) {
	var count int64
	session := p.conn.Default(ctx)
	sql := fmt.Sprintf("select count(*) from ms_project_collection where  member_code = ? and project_code = ?")
	raw := session.Raw(sql, memberId, projectCode)
	err := raw.Scan(&count).Error
	return count > 0, err
}

func (p *ProjectDao) FindProjectByPIdAndMemId(ctx context.Context, projectCode int64, memberId int64) (*pro.ProjectAndMember, error) {
	var pms *pro.ProjectAndMember
	session := p.conn.Default(ctx)
	sql := fmt.Sprintf("select a.*,b.project_code,b.member_code,b.join_time,b.is_owner,b.authorize from ms_project a, ms_project_member b where a.id = b.project_code and b.member_code = ? and b.project_code = ? limit 1")
	raw := session.Raw(sql, memberId, projectCode)
	fmt.Println("raw==>", raw)
	err := raw.Scan(&pms).Error
	fmt.Println("pms==>", pms)
	return pms, err
}

func (p *ProjectDao) SaveProject(conn database.DbConn, ctx context.Context, pr *pro.Project) error {
	p.conn = conn.(*gorms.GormConn)
	return p.conn.Tx(ctx).Save(&pr).Error
}

func (p *ProjectDao) SaveProjectMember(conn database.DbConn, ctx context.Context, pm *pro.ProjectMember) error {
	p.conn = conn.(*gorms.GormConn)
	return p.conn.Tx(ctx).Save(&pm).Error
}

func (p *ProjectDao) FindCollectProjectByMemId(ctx context.Context, memberId int64, page int64, size int64) ([]*pro.ProjectAndMember, int64, error) {
	var pms []*pro.ProjectAndMember
	session := p.conn.Default(ctx)
	index := (page - 1) * size
	sql := fmt.Sprintf("select * from ms_project where id in (select project_code from ms_project_collection where  member_code = ?) order by sort limit ?, ?")
	raw := session.Raw(sql, memberId, index, size)
	raw.Scan(&pms)
	var total int64
	query := fmt.Sprintf("member_code = ?")
	err := session.Model(&pro.ProjectCollection{}).Where(query, memberId).Count(&total).Error
	return pms, total, err
}

func (p *ProjectDao) FindProjectByMemId(ctx context.Context, memId int64, condition string, page int64, size int64) ([]*pro.ProjectAndMember, int64, error) {
	var pms []*pro.ProjectAndMember
	session := p.conn.Default(ctx)
	index := (page - 1) * size
	sql := fmt.Sprintf("select a.*,b.project_code,b.member_code,b.join_time,b.is_owner,b.authorize from ms_project a, ms_project_member b where a.id = b.project_code and b.member_code = ? %s order by sort limit ?, ?", condition)
	raw := session.Raw(sql, memId, index, size)
	raw.Scan(&pms)
	var total int64
	query := fmt.Sprintf("select count(*) from ms_project a, ms_project_member b where a.id = b.project_code and b.member_code = ? %s", condition)
	tx := session.Raw(query, memId)
	err := tx.Scan(&total).Error
	return pms, total, err
}

func NewProjectDao() *ProjectDao {
	return &ProjectDao{
		conn: gorms.New(),
	}
}
