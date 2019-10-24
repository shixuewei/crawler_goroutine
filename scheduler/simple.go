package scheduler

import "learning/crawler_goroutine/engine"

//简单版的调度器（用来对Request进行分发）所有Worker公用一个输入
//即，所有输入都在一个管道中，不同的worker进行抢夺下一个Request
type SimpleScheduler struct {
	//将请求数据存入管道中
	//输入
	workerChan chan engine.Request
}

//工作管道（用来将原来的Request放入管道中进行传输）
func (s *SimpleScheduler) WorkerChan() chan engine.Request {
	return s.workerChan
}

//提交（将请求数据传输到工作管道中）
func (s *SimpleScheduler) Submit(r engine.Request) {
	//将request发给worker chan
	//防止管道循环等待
	//（goroutine 可以防止数据在此卡住，底层会自动调用协程，会很快的将这个运行过去，
	//不会出现卡死的状态）
	go func() { s.workerChan <- r }()
}

func (s *SimpleScheduler) Run() {
	s.workerChan = make(chan engine.Request)
}

//提醒开始工作
func (s *SimpleScheduler) workerReady(chan engine.Request) {

}
