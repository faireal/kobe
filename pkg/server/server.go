package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/faireal/kobe/api"
	"github.com/patrickmn/go-cache"
	"github.com/prometheus/common/log"
	uuid "github.com/satori/go.uuid"
	"time"
)

type Kobe struct {
	taskCache      *cache.Cache
	inventoryCache *cache.Cache
	chCache        *cache.Cache
	cancelCahce    *cache.Cache
	pool           *Pool
}

func NewKobe() *Kobe {
	return &Kobe{
		taskCache:      cache.New(24*time.Hour, 5*time.Minute),
		chCache:        cache.New(24*time.Hour, 5*time.Minute),
		inventoryCache: cache.New(10*time.Minute, 5*time.Minute),
		cancelCahce:    cache.New(10*time.Minute, 5*time.Minute),
		pool:           NewPool(),
	}
}

func (k *Kobe) CreateProject(ctx context.Context, req *api.CreateProjectRequest) (*api.CreateProjectResponse, error) {
	pm := ProjectManager{}
	p, err := pm.CreateProject(req.Name, req.Source)
	if err != nil {
		return nil, err
	}
	resp := &api.CreateProjectResponse{
		Item: p,
	}
	return resp, nil
}

func (k *Kobe) ListProject(ctx context.Context, req *api.ListProjectRequest) (*api.ListProjectResponse, error) {
	pm := ProjectManager{}
	ps, err := pm.SearchProjects()
	if err != nil {
		return nil, err
	}
	resp := &api.ListProjectResponse{
		Items: ps,
	}
	return resp, nil
}

func (k *Kobe) DeleteProject(ctx context.Context, req *api.DeleteProjectRequest) (*api.DeleteProjectResponse, error) {
	pm := ProjectManager{}
	err := pm.DeleteProject(req.Name)
	if err != nil {
		return nil, err
	}
	resp := &api.DeleteProjectResponse{
		Success: true,
	}
	return resp, nil
}

func (k *Kobe) GetInventory(ctx context.Context, req *api.GetInventoryRequest) (*api.GetInventoryResponse, error) {
	item, _ := k.inventoryCache.Get(req.Id)
	if item == nil {
		return nil, errors.New("inventory is expire")
	}
	resp := &api.GetInventoryResponse{
		Item: item.(*api.Inventory),
	}
	return resp, nil
}

func (k *Kobe) WatchResult(req *api.WatchRequest, server api.KobeApi_WatchResultServer) error {
	ch, found := k.chCache.Get(req.TaskId)
	if !found {
		return errors.New(fmt.Sprintf("can not find task: %s", req.TaskId))
	}
	t, found := k.taskCache.Get(req.TaskId)
	if !found {
		return errors.New(fmt.Sprintf("can not find task: %s", req.TaskId))
	}
	tv, ok := t.(*api.Result)
	if !ok {
		return errors.New(fmt.Sprintf("invalid cache"))
	}
	if tv.Finished {
		return errors.New(fmt.Sprintf("task: %s already finished", req.TaskId))
	}
	val, ok := ch.(chan []byte)
	if !ok {
		return errors.New(fmt.Sprintf("invalid cache"))
	}
	for buf := range val {
		_ = server.Send(&api.WatchStream{
			Stream: buf,
		})
	}
	return nil
}

func (k *Kobe) RunAdhoc(ctx context.Context, req *api.RunAdhocRequest) (*api.RunAdhocResult, error) {
	rm := RunnerManager{
		inventoryCache: k.inventoryCache,
	}
	ch := make(chan []byte)
	id := uuid.NewV4().String()
	result := &api.Result{
		Id:        id,
		StartTime: time.Now().Format("2006-01-02 15:04:05"),
		EndTime:   "",
		Message:   "",
		Success:   false,
		Finished:  false,
		Content:   "",
	}
	ctx, cancelFunc := context.WithCancel(context.Background())
	k.cancelCahce.Set(result.Id, cancelFunc, cache.DefaultExpiration)
	k.taskCache.Set(result.Id, result, cache.DefaultExpiration)
	k.chCache.Set(result.Id, ch, cache.DefaultExpiration)
	k.inventoryCache.Set(result.Id, req.Inventory, cache.DefaultExpiration)
	runner, err := rm.CreateAdhocRunner(req.Pattern, req.Module, req.Param)
	if err != nil {
		return nil, err
	}
	task := func() {
		runner.Run(ctx, ch, result)
		result.Finished = true
		result.EndTime = time.Now().Format("2006-01-02 15:04:05")
		defer func() {
			// 回收资源 防止内存泄漏
			fn, _ := k.cancelCahce.Get(result.Id)
			fn.(context.CancelFunc)()
		}()
	}
	k.taskCache.Set(result.Id, result, cache.DefaultExpiration)
	k.pool.Commit(task)
	return &api.RunAdhocResult{
		Result: result,
	}, nil
}

func (k *Kobe) RunPlaybook(ctx context.Context, req *api.RunPlaybookRequest) (*api.RunPlaybookResult, error) {
	rm := RunnerManager{
		inventoryCache: k.inventoryCache,
	}
	ch := make(chan []byte)
	id := uuid.NewV4().String()
	result := &api.Result{
		Id:        id,
		StartTime: time.Now().Format("2006-01-02 15:04:05"),
		EndTime:   "",
		Message:   "",
		Success:   false,
		Finished:  false,
		Content:   "",
		Project:   req.Project,
	}
	ctx, cancelFunc := context.WithCancel(context.Background())
	k.cancelCahce.Set(result.Id, cancelFunc, cache.DefaultExpiration)
	k.taskCache.Set(result.Id, result, cache.DefaultExpiration)
	k.chCache.Set(result.Id, ch, cache.DefaultExpiration)
	k.inventoryCache.Set(result.Id, req.Inventory, cache.DefaultExpiration)
	runner, err := rm.CreatePlaybookRunner(req.Project, req.Playbook, req.Tag)
	if err != nil {
		return nil, err
	}
	b := func() {
		runner.Run(ctx, ch, result)
		result.Finished = true
		result.EndTime = time.Now().Format("2006-01-02 15:04:05")
		defer func() {
			// 回收资源 防止内存泄漏
			fn, _ := k.cancelCahce.Get(result.Id)
			fn.(context.CancelFunc)()
		}()
	}
	k.taskCache.Set(result.Id, result, cache.DefaultExpiration)
	k.pool.Commit(b)
	return &api.RunPlaybookResult{
		Result: result,
	}, nil
}

func (k *Kobe) GetResult(ctx context.Context, req *api.GetResultRequest) (*api.GetResultResponse, error) {
	id := req.GetTaskId()
	result, found := k.taskCache.Get(id)
	if !found {
		return nil, errors.New(fmt.Sprintf("can not find task: %s result", id))
	}
	val, ok := result.(*api.Result)
	if !ok {
		return nil, errors.New("invalid result type")
	}
	if val.Project == "" {
		val.Project = "adhoc"
	}
	return &api.GetResultResponse{Item: val}, nil
}

func (k *Kobe) ListResult(ctx context.Context, req *api.ListResultRequest) (*api.ListResultResponse, error) {
	var results []*api.Result
	resultMap := k.taskCache.Items()
	for taskId := range resultMap {
		item := resultMap[taskId].Object
		val, ok := item.(*api.Result)
		if !ok {
			continue
		}
		results = append(results, val)
	}
	return &api.ListResultResponse{
		Items: results,
	}, nil
}

func (k *Kobe) CancelTask(ctx context.Context, req *api.CancelTaskRequest) (*api.CancelTaskResponse, error) {
	id := req.GetTaskId()
	result, found := k.taskCache.Get(id)
	if !found {
		return nil, errors.New(fmt.Sprintf("can not find task: %s result", id))
	}
	val, ok := result.(*api.Result)
	if !ok {
		return nil, errors.New("invalid result type")
	}
	// 如果任务还没结束，取消任务
	if !val.Finished {
		cancel, found := k.cancelCahce.Get(id)
		if !found {
			return nil, errors.New(fmt.Sprintf("can not find task: %s cancel func", id))
		}
		cancel.(context.CancelFunc)()
		log.Infof("cancel task: %s result: %v", id, val)
	}
	return &api.CancelTaskResponse{
		Success: true,
	}, nil
}
