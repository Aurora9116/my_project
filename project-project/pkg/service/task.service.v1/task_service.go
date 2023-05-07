package task_service_v1

import (
	"context"
	"github.com/jinzhu/copier"
	"go.uber.org/zap"
	"test.com/project-common/encrypts"
	"test.com/project-common/errs"
	"test.com/project-common/tms"
	"test.com/project-grpc/task"
	"test.com/project-grpc/user/login"
	"test.com/project-project/internal/dao"
	"test.com/project-project/internal/data"
	"test.com/project-project/internal/data/pro"
	"test.com/project-project/internal/database"
	"test.com/project-project/internal/database/tran"
	"test.com/project-project/internal/repo"
	"test.com/project-project/internal/rpc"
	"test.com/project-project/pkg/model"
	"time"
)

type TaskService struct {
	task.UnimplementedTaskServiceServer
	cache                  repo.Cache
	transaction            tran.Transaction
	projectRepo            repo.ProjectRepo
	projectTemplateRepo    repo.ProjectTemplateRepo
	taskStagesTemplateRepo repo.TaskStagesTemplateRepo
	taskStagesRepo         repo.TaskStagesRepo
	taskRepo               repo.TaskRepo
}

func New() *TaskService {
	return &TaskService{
		cache:                  dao.Rc,
		transaction:            dao.NewTransactionImpl(),
		projectRepo:            dao.NewProjectDao(),
		projectTemplateRepo:    dao.NewProjectTemplateDao(),
		taskStagesTemplateRepo: dao.NewTaskStagesTemplateDao(),
		taskStagesRepo:         dao.NewTaskStagesDao(),
		taskRepo:               dao.NewTaskDao(),
	}
}
func (t *TaskService) TaskStages(c context.Context, msg *task.TaskReqMessage) (*task.TaskStagesResponse, error) {
	projectCode := encrypts.DecryptNotErr(msg.ProjectCode)
	page := msg.Page
	pageSize := msg.PageSize
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	stages, total, err := t.taskStagesRepo.FindStagesByProjectId(ctx, projectCode, page, pageSize)
	if err != nil {
		zap.L().Error("project SaveProject taskStagesRepo.FindStagesByProjectId error", zap.Error(err))
		return nil, errs.GrpcError(model.DbError)
	}
	var tsMessage []*task.TaskStagesMessage
	copier.Copy(&tsMessage, stages)
	if tsMessage == nil {
		return &task.TaskStagesResponse{Total: 0, List: tsMessage}, nil
	}
	stagesMap := data.ToTaskStagesMap(stages)
	for _, v := range tsMessage {
		taskStages := stagesMap[int(v.Id)]
		v.Code = encrypts.EncryptNoErr(int64(v.Id))
		v.CreateTime = tms.FormatByMill(taskStages.CreateTime)
		v.ProjectCode = msg.ProjectCode
	}
	return &task.TaskStagesResponse{Total: total, List: tsMessage}, nil
}
func (t *TaskService) MemberProjectList(c context.Context, msg *task.TaskReqMessage) (*task.MemberProjectResponse, error) {
	// 1.去project_member表中查询用户id列表
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	projectCode := encrypts.DecryptNotErr(msg.ProjectCode)
	projectMember, total, err := t.projectRepo.FindProjectMemberByPid(ctx, projectCode)
	if err != nil {
		zap.L().Error("project MemberProjectList projectRepo.FindProjectMemberByPid error", zap.Error(err))
		return nil, errs.GrpcError(model.DbError)
	}
	// 2.拿上用户id列表 去请求用户信息
	if projectMember == nil || len(projectMember) <= 0 {
		return &task.MemberProjectResponse{List: nil, Total: 0}, nil
	}
	var mIds []int64
	pmMap := make(map[int64]*pro.ProjectMember)
	for _, v := range projectMember {
		mIds = append(mIds, v.MemberCode)
		pmMap[v.MemberCode] = v
	}
	// 请求用户信息
	userMsg := &login.UserMessage{
		MIds: mIds,
	}
	memberMessageList, err := rpc.LoginServiceClient.FindMemInfoByIds(ctx, userMsg)
	if err != nil {
		zap.L().Error("project MemberProjectList LoginServiceClient.FindMemInfoByIds error", zap.Error(err))
		return nil, errs.GrpcError(model.DbError)
	}
	var list []*task.MemberProjectMessage
	for _, v := range memberMessageList.List {
		owner := pmMap[v.Id].IsOwner
		mpm := &task.MemberProjectMessage{
			Name:       v.Name,
			Avatar:     v.Avatar,
			MemberCode: v.Id,
			Code:       v.Code,
			Email:      v.Email,
		}
		if owner == v.Id {
			mpm.IsOwner = model.Owner
		}
		list = append(list, mpm)
	}
	return &task.MemberProjectResponse{Total: total, List: list}, nil
}
func (t *TaskService) TaskList(c context.Context, msg *task.TaskReqMessage) (*task.TaskListResponse, error) {
	stageCode := encrypts.DecryptNotErr(msg.StageCode)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	taskList, err := t.taskRepo.FindTaskByStageCode(c, int(stageCode))
	if err != nil {
		zap.L().Error("project task TaskList taskRepo.FindTaskByStageCode error", zap.Error(err))
		return nil, errs.GrpcError(model.DbError)
	}
	var taskDisplayList []*data.TaskDisplay
	var mIds []int64
	for _, v := range taskList {
		display := v.ToTaskDisplay()
		if v.Private == 1 {
			// 代表隐私模式
			taskMember, err := t.taskRepo.FindTaskMemberByTaskId(ctx, v.Id, msg.MemberId)
			if err != nil {
				zap.L().Error("project task TaskList taskRepo.FindTaskMemberByTaskId error", zap.Error(err))
				return nil, errs.GrpcError(model.DbError)
			}
			if taskMember != nil {
				display.CanRead = model.CanRead
			} else {
				display.CanRead = model.NoCanRead
			}
		}
		taskDisplayList = append(taskDisplayList, display)
		mIds = append(mIds, v.AssignTo)
	}
	if mIds == nil || len(mIds) <= 0 {
		return &task.TaskListResponse{List: nil}, nil
	}
	messageList, err := rpc.LoginServiceClient.FindMemInfoByIds(ctx, &login.UserMessage{MIds: mIds})
	if err != nil {
		zap.L().Error("project task TaskList LoginServiceClient.FindMemInfoByIds error", zap.Error(err))
		return nil, errs.GrpcError(model.DbError)
	}
	memberMap := make(map[int64]*login.MemberMessage)
	for _, v := range messageList.List {
		memberMap[v.Id] = v
	}
	for _, v := range taskDisplayList {
		message := memberMap[encrypts.DecryptNotErr(v.AssignTo)]
		e := data.Executor{
			Name:   message.Name,
			Avatar: message.Avatar,
		}
		v.Executor = e
	}
	var taskMessageList []*task.TaskMessage
	copier.Copy(&taskMessageList, taskDisplayList)
	return &task.TaskListResponse{List: taskMessageList}, nil
}
func (t *TaskService) SaveTask(c context.Context, msg *task.TaskReqMessage) (*task.TaskMessage, error) {
	// 1.检查业务逻辑
	if msg.Name == "" {
		return nil, errs.GrpcError(model.TaskNameNotNull)
	}
	stageCode := encrypts.DecryptNotErr(msg.StageCode)
	taskStages, err := t.taskStagesRepo.FindById(c, int(stageCode))
	if err != nil {
		zap.L().Error("project task SaveTask taskStagesRepo.FindById error", zap.Error(err))
		return nil, errs.GrpcError(model.DbError)
	}
	if taskStages == nil {
		return nil, errs.GrpcError(model.TaskStagesNotNull)
	}
	projectCode := encrypts.DecryptNotErr(msg.ProjectCode)
	project, err := t.projectRepo.FindProjectById(c, projectCode)
	if err != nil {
		zap.L().Error("project task SaveTask projectRepo.FindProjectById error", zap.Error(err))
		return nil, errs.GrpcError(model.DbError)
	}
	if project == nil || project.Deleted == model.Deleted {
		return nil, errs.GrpcError(model.ProjectAlreadyDeleted)
	}
	maxIdNum, err := t.taskRepo.FindTaskMaxIdNum(c, projectCode)
	if err != nil {
		zap.L().Error("project task SaveTask taskRepo.FindTaskMaxIdNum error", zap.Error(err))
		return nil, errs.GrpcError(model.DbError)
	}
	if maxIdNum == nil {
		a := 0
		maxIdNum = &a
	}
	maxSort, err := t.taskRepo.FindTaskSort(c, projectCode, stageCode)
	if err != nil {
		zap.L().Error("project task SaveTask taskRepo.FindTaskSort error", zap.Error(err))
		return nil, errs.GrpcError(model.DbError)
	}
	if maxSort == nil {
		a := 0
		maxSort = &a
	}
	assignTo := encrypts.DecryptNotErr(msg.AssignTo)
	ts := &data.Task{
		Name:        msg.Name,
		CreateTime:  time.Now().UnixMilli(),
		CreateBy:    msg.MemberId,
		AssignTo:    assignTo,
		ProjectCode: projectCode,
		StageCode:   int(stageCode),
		IdNum:       *maxIdNum + 1,
		Private:     project.OpenTaskPrivate,
		Sort:        *maxSort + 65536,
		BeginTime:   time.Now().UnixMilli(),
		EndTime:     time.Now().Add(2 * 24 * time.Hour).UnixMilli(),
	}
	err = t.transaction.Action(func(conn database.DbConn) error {
		err = t.taskRepo.SaveTask(c, conn, ts)
		if err != nil {
			zap.L().Error("project task SaveTask taskRepo.SaveTask error", zap.Error(err))
			return errs.GrpcError(model.DbError)
		}
		tm := &data.TaskMember{
			MemberCode: assignTo,
			TaskCode:   ts.Id,
			JoinTime:   time.Now().UnixMilli(),
			IsOwner:    model.Owner,
		}
		if assignTo == msg.MemberId {
			tm.IsExecutor = model.Executor
		}
		err = t.taskRepo.SaveTaskMember(c, conn, tm)
		if err != nil {
			zap.L().Error("project task SaveTask taskRepo.SaveTaskMember error", zap.Error(err))
			return errs.GrpcError(model.DbError)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	display := ts.ToTaskDisplay()
	member, err := rpc.LoginServiceClient.FindMemInfoById(c, &login.UserMessage{MemId: assignTo})
	if err != nil {
		return nil, err
	}
	display.Executor = data.Executor{
		Name:   member.Name,
		Avatar: member.Avatar,
		Code:   member.Code,
	}
	tm := &task.TaskMessage{}
	copier.Copy(tm, display)
	return tm, nil
}
func (t *TaskService) TaskSort(c context.Context, msg *task.TaskReqMessage) (*task.TaskSortResponse, error) {
	// 移动的任务id肯定有preTaskCode
	preTaskCode := encrypts.DecryptNotErr(msg.PreTaskCode)
	toStageCode := encrypts.DecryptNotErr(msg.ToStageCode)
	if msg.PreTaskCode == msg.NextTaskCode {
		return &task.TaskSortResponse{}, nil
	}
	err := t.sortTask(preTaskCode, msg.NextTaskCode, toStageCode)
	if err != nil {
		return nil, err
	}
	return &task.TaskSortResponse{}, nil
}

func (t *TaskService) sortTask(preTaskCode int64, nextTaskCode string, toStageCode int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	ts, err := t.taskRepo.FindTaskById(ctx, preTaskCode)
	if err != nil {
		zap.L().Error("project task TaskSort taskRepo.FindTaskById error", zap.Error(err))
		return errs.GrpcError(model.DbError)
	}
	t.transaction.Action(func(conn database.DbConn) error {
		ts.StageCode = int(toStageCode)

		if nextTaskCode != "" {
			// 意味着要进行替换
			nextTaskCode1 := encrypts.DecryptNotErr(nextTaskCode)
			next, err := t.taskRepo.FindTaskById(ctx, nextTaskCode1)
			if err != nil {
				zap.L().Error("project task TaskSort taskRepo.FindTaskById error", zap.Error(err))
				return errs.GrpcError(model.DbError)
			}
			prepre, err := t.taskRepo.FindTaskByStageCodeLtSort(ctx, next.StageCode, next.Sort)
			if err != nil {
				zap.L().Error("project task TaskSort taskRepo.FindTaskById error", zap.Error(err))
				return errs.GrpcError(model.DbError)
			}
			if prepre != nil {
				ts.Sort = (prepre.Sort + next.Sort) / 2
			}
			if prepre == nil {
				ts.Sort = 0
			}
			//ts.Sort, next.Sort = next.Sort, ts.Sort
			//err = t.taskRepo.UpdateTaskSort(ctx, conn, next)
			//if err != nil {
			//	zap.L().Error("project task TaskSort taskRepo.UpdateTaskSort error", zap.Error(err))
			//	return errs.GrpcError(model.DbError)
			//}
		} else {
			maxSort, err := t.taskRepo.FindTaskSort(ctx, ts.ProjectCode, int64(ts.StageCode))
			if err != nil {
				zap.L().Error("project task TaskSort taskRepo.FindTaskById error", zap.Error(err))
				return errs.GrpcError(model.DbError)
			}
			ts.Sort = *maxSort + 65536
		}
		err = t.taskRepo.UpdateTaskSort(ctx, conn, ts)
		if err != nil {
			zap.L().Error("project task TaskSort taskRepo.UpdateTaskSort error", zap.Error(err))
			return errs.GrpcError(model.DbError)
		}
		return nil
	})
	return err
}
