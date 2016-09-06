package task

import "sync/atomic"

type Task interface {
	Run()
}
type TaskWaiter struct {
	doneChannel chan bool
	T           Task
}

func (waiter *TaskWaiter) runTask() {
	waiter.T.Run()
}
func (waiter *TaskWaiter) doneTask() {
	waiter.doneChannel <- true
}

func (waiter *TaskWaiter) Done() {
	<-waiter.doneChannel
}

const DEFAULT_TASK_BUFFER = 1024

type TaskExecutor struct {
	taskChannel        chan *TaskWaiter
	maxNumOfGoroutines int32 //正在运行的
	stopChannel        chan bool

	numOfRunningGoroutines int32
	numOfGoroutines        int32
}

func NewTaskExecutor(maxNumOfGoroutines int32) *TaskExecutor {
	executor := &TaskExecutor{
		taskChannel:            make(chan *TaskWaiter, DEFAULT_TASK_BUFFER),
		maxNumOfGoroutines:     maxNumOfGoroutines,
		stopChannel:            make(chan bool, 1),
		numOfGoroutines:        0,
		numOfRunningGoroutines: 0,
	}

	return executor
}
func (executor *TaskExecutor) isFull() bool {
	return atomic.LoadInt32(&executor.numOfGoroutines) ==
		atomic.LoadInt32(&executor.maxNumOfGoroutines)
}
func (executor *TaskExecutor) hasWaitingGoroutines() bool {
	return atomic.LoadInt32(&executor.numOfGoroutines) >
		atomic.LoadInt32(&executor.numOfRunningGoroutines)
}
func (executor *TaskExecutor) spawnGoroutine() {
	go func() {
		over := false
		for !over {
			select {
			case task := <-executor.taskChannel:
				atomic.AddInt32(&executor.numOfRunningGoroutines, 1)
				task.runTask()
				task.doneTask()
				atomic.AddInt32(&executor.numOfRunningGoroutines, -1)
			case <-executor.stopChannel:
				over = true
			}
		}
	}()
	atomic.AddInt32(&executor.numOfGoroutines, 1)
}

func (executor *TaskExecutor) Submit(t Task) *TaskWaiter {
	if !executor.hasWaitingGoroutines() && !executor.isFull() {
		executor.spawnGoroutine()
	}
	tw := &TaskWaiter{
		doneChannel: make(chan bool),
		T:           t,
	}
	executor.taskChannel <- tw
	return tw
}

func (executor *TaskExecutor) WaitToFinish() {

}

func (executor *TaskExecutor) Close() {
	executor.stopChannel <- true
}
