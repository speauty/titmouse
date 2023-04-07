package repository

import (
	"sync"
	"titmouse/model"
)

var (
	apiTaskRepository  *TaskRepository
	onceTaskRepository sync.Once
)

func ApiTask() *TaskRepository {
	onceTaskRepository.Do(func() {
		apiTaskRepository = new(TaskRepository)
	})
	return apiTaskRepository
}

// TaskRepository 关于任务的相关操作仓库
type TaskRepository struct {
	*Repository
}

func (customTR TaskRepository) ActionTaskPusher(task *model.TaskModel) {
	// @todo 待完成
}

func (customTR TaskRepository) ActionQueryTask(task *model.TaskModel) {
	// @todo 待完成
}
