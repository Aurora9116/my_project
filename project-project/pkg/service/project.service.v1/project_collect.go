package project_service_v1

import (
	"context"
	"go.uber.org/zap"
	"strconv"
	"test.com/project-common/encrypts"
	"test.com/project-common/errs"
	"test.com/project-grpc/project"
	"test.com/project-project/internal/data"
	"test.com/project-project/pkg/model"
	"time"
)

func (p *ProjectService) UpdateCollectProject(ctx context.Context, msg *project.ProjectRpcMessage) (*project.CollectProjectResponse, error) {
	projectCodeStr, _ := encrypts.Decrypt(msg.ProjectCode, model.AESKey)
	projectCode, _ := strconv.ParseInt(projectCodeStr, 10, 64)
	c, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	var err error
	if msg.CollectType == "collect" {
		pc := &data.ProjectCollection{
			ProjectCode: projectCode,
			MemberCode:  msg.MemberId,
			CreateTime:  time.Now().UnixMilli(),
		}
		err = p.projectRepo.SaveProjectCollect(c, pc)
	}
	if msg.CollectType == "cancel" {

		err = p.projectRepo.DeleteProjectCollect(c, msg.MemberId, projectCode)
	}
	if err != nil {
		zap.L().Error("project UpdateCollectProject SaveProjectCollect error", zap.Error(err))
		return nil, errs.GrpcError(model.DbError)
	}
	return &project.CollectProjectResponse{}, nil
}
