package domain

import (
	"context"
	"test.com/project-common/errs"
	"test.com/project-project/internal/dao"
	"test.com/project-project/internal/data"
	"test.com/project-project/internal/repo"
	"test.com/project-project/pkg/model"
)

type ProjectNodeDomain struct {
	projectNodeRepo repo.ProjectNodeRepo
}

func (d *ProjectNodeDomain) TreeList() ([]*data.ProjectNodeTree, *errs.BError) {
	list, err := d.projectNodeRepo.FindAll(context.Background())
	if err != nil {
		return nil, model.DbError
	}
	return data.ToNodeTreeList(list), nil
}

func (d *ProjectNodeDomain) NodeList() ([]*data.ProjectNode, *errs.BError) {
	list, err := d.projectNodeRepo.FindAll(context.Background())
	if err != nil {
		return nil, model.DbError
	}
	return list, nil
}

func NewProjectNodeDomain() *ProjectNodeDomain {
	return &ProjectNodeDomain{
		projectNodeRepo: dao.NewProjectNodeDao(),
	}
}
