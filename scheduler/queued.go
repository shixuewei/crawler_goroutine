package scheduler

import "learning/crawler_goroutine/engine"

//队列调度器（将调度器通过队列进行控制）--用管道进行输入输出（用管道进行传输）
//将请求与工作两块分开，分别为一个goroutine，然后通过管道进行传输
type QueuedScheduler struct {
	//请求队列
	requestChan chan engine.Request
	//工作队列，用来处理worker
	workerChan chan chan engine.Request
}

//提供方法用来告诉engine，我这有一个worker，你可以给我哪个engine.Request
func (q *QueuedScheduler) WorkerChan() chan engine.Request {
	return make(chan engine.Request)
}

//提醒开始工作（从外部告诉我们，有一个worker已经准备好了，可以接收request了）
func (q *QueuedScheduler) WorkerReady(w chan engine.Request) {
	q.workerChan <- w
}

//提交（将请求数据传输到请求队列中）
func (q *QueuedScheduler) Submit(r engine.Request) {
	q.requestChan <- r
}

//总控
func (q *QueuedScheduler) Run() {
	//初始化workerChan
	q.workerChan = make(chan chan engine.Request)
	//初始化requestChan
	q.requestChan = make(chan engine.Request)
	//调用一个协程来进行调度器的运行
	go func() {
		//请求队列
		var requestQ []engine.Request
		//工作队列
		var workQ []chan engine.Request
		//循环等待
		for {
			//待激活的请求
			var activeRequest engine.Request
			//待工作的worker
			var activeWorker chan engine.Request
			//既有worker又有request时，准备好了，可以发了
			if len(requestQ) > 0 && len(workQ) > 0 {
				activeWorker = workQ[0]
				activeRequest = requestQ[0]
			}
			select {
			//发来的是请求，则将其加入到请求队列中
			case r := <-q.requestChan:
				requestQ = append(requestQ, r)
			//发来的是worker，则将其加入到worker队列中
			case w := <-q.workerChan:
				workQ = append(workQ, w)
			//如果已经将请求发送到worker中，则将他们从队列中拿掉
			case activeWorker <- activeRequest:
				requestQ = requestQ[1:]
				workQ = workQ[1:]
			}
		}
	}()
}
