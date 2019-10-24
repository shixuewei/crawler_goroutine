package engine

//并发版的爬虫引擎
type ConcurrentEngine struct {
	//调度器
	Scheduler Scheduler
	//worker数量
	WorkerCount int
	//定义一个管道，类型是Item类型,用来保存相应的数据
	ItemChan chan Item
	//对url进行处理
	RequestProcessor Processor
}

//定义一个处理变量
type Processor func(Request) (ParseResult, error)

//定义一个调度器接口，在别处实现
type Scheduler interface {
	//提醒开始工作
	ReadyNotifier
	//用调度器处理请求（将engine发来的请求通过管道传输给一堆worker）
	Submit(Request)
	//提供方法用来告诉engine，我这有一个worker，你可以给我哪个engine.Request
	//输入
	WorkerChan() chan Request
	//运行
	Run()
}

//提醒开始工作（从外部告诉我们，有一个worker已经准备好了，可以接收request了）
//之所以放在外面单独调用，是因为如果直接在Scheduler中，则createWorker调用时传入
//Scheduler，显得太笨重，所以让它单独出来
type ReadyNotifier interface {
	WorkerReady(chan Request)
}

//引擎工作
func (c *ConcurrentEngine) Run(seeds ...Request) {
	//将ParseResult放在一个管道中进行传输（保证不出现数据竞争）{输出结果}
	out := make(chan ParseResult)
	//启动调度器
	//创建两个chan，然后就开始等待Request和worker的到来
	c.Scheduler.Run()
	//启动多个worker
	for i := 0; i < c.WorkerCount; i++ {
		//c.Scheduler.WorkerChan()输入
		c.createWorker(c.Scheduler.WorkerChan(), out, c.Scheduler)
	}

	//用调度器处理请求（将engine发来的请求通过管道传输给一堆worker）
	for _, r := range seeds {
		c.Scheduler.Submit(r)
	}

	//（处理输出结果）处理具体的item和request（有新的request就将其加到管道的后面）
	for {
		result := <-out
		for _, item := range result.Items {
			go func() { c.ItemChan <- item }()
		}

		//处理完Item后，将result中的Request送到Scheduler等待分发
		for _, request := range result.Requests {
			//分发Request
			c.Scheduler.Submit(request)
		}
	}
}

//创建一个worker（第一个参数为输入，第二个为输出，第三个为判断是否准备好）
//worker的输入与输出是一个channel，所以不存在循环等待的影响
func (c *ConcurrentEngine) createWorker(in chan Request, out chan ParseResult, ready ReadyNotifier) {
	//采用gofunc的原因就是因为是管道传输，用goroutine
	go func() {
		for {
			//提醒开始工作（从外部告诉我们，有一个worker已经准备好了，可以接收request了）
			ready.WorkerReady(in)
			//在request := <-in这一行，如果这行没有传入，那下面的out <- result就没有传出
			//然后就开始等待，最终卡死在这块
			request := <-in
			result, err := c.RequestProcessor(request)
			if err != nil {
				continue
			}
			out <- result
		}
	}()
}
