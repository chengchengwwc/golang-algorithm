
### 同步原语和锁
Golang作为一个原生支持用户态的语言，当提到并发进程，多线程的时候，是离不开锁的，锁是一种并发编程中的同步原语（Synchronization Primitives），它能保证多个 Goroutine 在访问同一片内存时不会出现竞争条件（Race condition）等问题。
#### 基于原语
go语言在sync包中提供了用于同步的一些基本原语，包括常见的sync.Mutex,sync.RWMutex,sync.WaitGroup,
sync.Once,sync.Cond.
这些基本原语提高了较为基础的同步功能，但是它们是一种相对原始的同步机制，在多数情况下，我们都应该使用抽象层级的更高的 Channel 实现同步。

##### Mutex
Mutex由两个字段：state,sema组成，其中state表示当前互斥锁的状态，而sema是用于控制锁的状态的信号量。上述两个加起来，只占用8个字节
1. 状态
   1. 在默认的情况下，互斥锁的所有的状态都是0，int32中不同位分别表示了不同的状态
      1. mutexLocked — 表示互斥锁的锁定状态；
      2. mutexWoken — 表示从正常模式被从唤醒；
      3. mutexStarving — 当前的互斥锁进入饥饿状态；
      4. waitersCount — 当前互斥锁上等待的 Goroutine 个数；
2. 正常模式和饥饿模
   1. 正常模式：锁的的等待者会按照新出的顺序获取锁，但是刚被唤起的goroutine和新创造的进程竞争的时候，大概率会获得锁，为了防止这种情况，一旦goroutine超过1ms没有获得锁，就会将当前的状态切换为饥饿模式，防止部分 Goroutine 被『饿死』。
   2. 在饥饿模式中，互斥锁会直接交给等待队列最前面的 Goroutine。新的 Goroutine 在该状态下不能获取锁、也不会进入自旋状态，它们只会在队列的末尾等待。如果一个 Goroutine 获得了互斥锁并且它在队列的末尾或者它等待的时间少于 1ms，那么当前的互斥锁就会被切换回正常模式。

3. 加锁和解锁
   1. 互斥锁的加锁是靠 sync.Mutex.Lock 完成的，最新的 Go 语言源代码中已经将 sync.Mutex.Lock 方法进行了简化，方法的主干只保留最常见、简单的情况 — 当锁的状态是 0 时，将 mutexLocked 位置成 1：如果互斥锁的状态不是0的时候就会调用sync.Mutex.lockSlow 尝试通过自旋（Spinnig）等方式等待锁的释放，该方法的主体是一个非常大 for 循环，这里将该方法分成几个部分介绍获取锁的过程：
       1. 判断当前goroutine是否进入自旋转
       2. 通过自旋等待互斥锁的释放；
       3. 计算互斥锁的最新状态；
       4. 更新互斥锁的状态并获取锁
   2. 自旋是一种多线程同步机制，当前的进程在进入自旋的过程中会一直保持 CPU 的占用，持续检查某个条件是否为真。在多核的 CPU 上，自旋可以避免 Goroutine 的切换，使用恰当会对性能带来很大的增益，但是使用的不恰当就会拖慢整个程序，所以 Goroutine 进入自旋的条件非常苛刻:
      1. 互斥锁只有在普通模式下才会进入自旋
      2. sync.runtime_canSpin 需要返回 true：
         1. 运行在多 CPU 的机器上；
         2. 当前 Goroutine 为了获取该锁进入自旋的次数小于四次；
         3. 当前机器上至少存在一个正在运行的处理器 P 并且处理的运行队列为空；
   3. 如果没有通过CAS 获得锁，会调用 sync.runtime_SemacquireMutex 使用信号量保证资源不会被两个 Goroutine 获取。sync.runtime_SemacquireMutex 会在方法中不断调用尝试获取锁并休眠当前 Goroutine 等待信号量的释放，一旦当前 Goroutine 可以获取信号量，它就会立刻返回。

#### RWMutex
读写互斥锁sync.RWMutex，是细粒度的互斥锁，她并不限制资源的并发读，但是读写，写写操作无法并行执行。一个常见的服务对资源的读写比例会非常高，因为大多数的读请求之间不会相互影响，所以我们可以读写资源操作的分离，在类似场景下提高服务的性能。
##### 结构体
sync.RWMutex 中总共包含以下 5 个字段：
```
type RWMUtex struct {
    w  Mutex
    writerSem   uint32
    readerSem   uint32
    readerCount int32
    readerWait  int32
}
```
- w 复用互斥锁提供的能力
- writerSem和readSem 分别用于写等待和读等待
- readerCount 存储了当前正在执行的读操作的数量
- readerWait 表示当写操作被阻塞时等待的读操作的个数

我们会依次分析获取写锁和读锁的实现能力，其中：
- 写操作使用 sync.RWMutex.Lock 和 sync.RWMutex.Unlock 方法；
- 读操作使用 sync.RWMutex.RLock 和 sync.RWMutex.RUnlock 方法；

##### 写锁
1. 当资源的使用者想要获取写锁时，需要调用 sync.RWMutex.Lock 方法
2. 写锁的释放会调用 sync.RWMutex.Unlock 方法

与加锁的过程正好相反，写锁的释放分以下几个执行
1. 调用atomic.AddInt32 函数将变回正数，释放读锁；
2. 通过 for 循环触发所有由于获取读锁而陷入等待的 Goroutine
3. 调用 sync.Mutex.Unlock 方法释放写锁

##### 读锁
读锁的加锁方法 sync.RWMutex.RLock
```
func (rw *RWMutex) RLock(){
    if atomic.AddInt32(&rw.readerCount,1) < 0 {
        runtime_SemacquireMutex(&rw.readerSem,false,0)
    }
}
```
1. 如果该方法返回函数-其他 Goroutine 获得了写锁，当前 Goroutine 就会调用 sync.runtime_SemacquireMutex 陷入休眠等待锁的释放。
2. 如果该方法的结果为非负数 — 没有 Goroutine 获得写锁，当前方法就会成功返回.

当 Goroutine 想要释放读锁时，会调用如下所示的 sync.RWMutex.RUnlock 方法
```
func (rw *RWMutex) RUnlock() {
    if r := atomic.AddInt32(&rw.readerCount,-1);r<0{
        rw.rUnlockSlow(r)
    }
}
```

#### WaitGroup
sync.WaitGroup 可以等待一组 Goroutine 的返回，一个比较常见的使用场景是批量发出 RPC 或者 HTTP 请求：
```
reuqests := []*Requests{...}
wg := &sync.WaitGroup()
wg.Add(len(requests))
for _,request := range requests {
    go func(r *Request){
        defer wg.Done()
    }(request)
}
wg.Wait()

```
我们可以通过 sync.WaitGroup 将原本顺序执行的代码在多个 Goroutine 中并发执行，加快程序处理的速度。
##### 结构体
sync.WaitGroup 结构体中的成员变量非常简单，其中只包含两个成员变量
```
type WaitGroup struct {
    noCopy noCopy
    state1 [3]uint32
}

```
- noCopy 保证 sync.WaitGroup 不会被开发者通过再赋值的方式拷贝
- state1 存储着状态和信号量
##### 接口
其中的 sync.WaitGroup.Done 只是向 sync.WaitGroup.Add 方法传入了 -1，所以我们重点分析另外两个方法 sync.WaitGroup.Add 和 sync.WaitGroup.Wait
```
func (wg *WaitGroup) Add(delta int){
    statep,semap := wg.state()
    state := atomic.AddUint64(statep,uint64(delta)<<32)
    v := int32(state >>32)
    w := uint32(state)
    if v < 0 {
        panic("sync: negative WaitGroup counter")
    }

    if v > 0 || w == 0{
        return 
    }
    *statep = 0
    for ; w != 0; w-- {
        runtime_Semrelease(semap, false, 0)
    } 
}
```
另一个方法 sync.WaitGroup.Wait
```
func (wg *WaitGroup) Wait(){
    statep,semp := wg.state()
    for {
        state := atomic.LoadUint64(statep)
        v :=int32(state >> 32)
        if v == 0{
            return
        }
        if atomic.CompareAndSwapUint64(statep, state, state+1) {
			runtime_Semacquire(semap)
			if +statep != 0 {
				panic("sync: WaitGroup is reused before previous Wait has returned")
			}
            return 
        }
    }
}
```
当 sync.WaitGroup 的计数器归零时，当陷入睡眠状态的 Goroutine 就被唤醒

#### Once
Go 语言标准库中 sync.Once 可以保证在 Go 程序运行期间的某段代码只会执行一次。在运行如下所示的代码时，我们会看到如下所示的运行结果
```
func main() {
    o := &sync.Once{}
    for i:=0;i<10;i++{
        o.Do(func(){
            fmt.Println("ddd)
        })
    }
}
```
##### 结构体
每一个 sync.Once 结构体中都只包含一个用于标识代码块是否执行过的 done 以及一个互斥锁 sync.Mutex
```
type Once struct {
    done uint32
    m Mutex
}
```
##### 接口
sync.Once.Do 是 sync.Once 结构体对外唯一暴露的方法
- 如果传入的函数已经执行过了，就会直接返回
- 如果传入的函数没有执行过，就会调用sync.Once.doSlow执行传入函数
```
func (o *Once) Do(f func()){
    if atomic.LoadUint32(&o.done) == 0 {
        o.doSlow(f)
    }
}

func (o *Once) doSlow(f func()){
    o.m.Lock()
    defer o.m.Unlock()
    if o.done == 0 {
        defer atomic.StoreUinit32(&o.done,1)
        f()
    }
}
```

#### Cond
Go标准库的中的sync.Cond是一个条件变量，它可以让一系列的goroutine都在满足特定条件下时候被唤醒，每一个 sync.Cond 结构体在初始化时都需要传入一个互斥锁，我们可以通过下面的例子了解它的使用方法
```
func main() {
    c := sync.NewCond(&sync.Mutex{})
    for i :=0;i<10;i++{
        go listen(c)
    }
    time.Sleep(1*time.Second)
    go broadcast(c)

    ch := make(chan os.Signal,1)
    signal.Notify(ch, os.Interrupt)
	<-ch
}

func broadcast(c *sync.Cond){
    c.l.Lock()
    c.Broadcast()
    c.l.Unlock()
}

func listen(c *sync.Cond) {
    c.l.Lock()
    c.wait()
    fmt.Println("ddd")
    c.l,Unlock()
}
```
上述代码同时运行了 11 个 Goroutine，这 11 个 Goroutine 分别做了不同事情：
- 10 个 Goroutine 通过 sync.Cond.Wait 等待特定条件的满足；
- 1 个 Goroutine 会调用 sync.Cond.Broadcast 方法通知所有陷入等待的 Goroutine；

sync.Cond.Signal 和 sync.Cond.Broadcast 方法就是用来唤醒调用 sync.Cond.Wait 陷入休眠的 Goroutine，它们两个的实现有一些细微差别：
- sync.Cond.Signal 方法会唤醒队列最前面的 Goroutine；
- sync.Cond.Broadcast 方法会唤醒队列中全部的 Goroutine；

在一般情况下，我们都会先调用 sync.Cond.Wait 陷入休眠等待满足期望条件，当满足唤醒条件时，就可以选择使用 sync.Cond.Signal 或者 sync.Cond.Broadcast 唤醒一个或者全部的 Goroutine。

#### ErrGroup
x/sync/errgroup.Group 就为我们在一组 Goroutine 中提供了同步、错误传播以及上下文取消的功能，我们可以使用如下所示的方式并行获取网页的数据
```
var g errgroup.Group
var urls = []string{
    "http://www.golang.org"
    "http://www.baidu.com"
}

for i := range urls {
    url := urls[i]
    g.Go(func() error {
        resp,err := http.Get(url)
        if err == nil{
            resp.Body.Close()
        }
        return err
    })
}

if err := g.Wait();err == nil{
    fmt.Println("Successfully fetched all URLs.")
}
```
x/sync/errgroup.Group.Go 方法能够创建一个 Goroutine 并在其中执行传入的函数，而 x/sync/errgroup.Group.Wait 会等待所有 Goroutine 全部返回，该方法的不同返回结果也有不同的含义：
- 如果返回错误 — 这一组 Goroutine 最少返回一个错误；
- 如果返回空值 — 所有 Goroutine 都成功执行

#### Semaphore
信号量是在并发编程中常见的一种同步机制，在需要控制访问资源的进程数量时就会用到信号量，它会保证持有的计数器在 0 到初始化的权重之间波动
- 每次获取资源时都会将信号量中的计数器减去对应的数值，在释放时重新加回来
- 当遇到计数器大于信号量大小时就会进入休眠等待其他线程释放信号

这个结构体对外也只暴露了四个方法：
- x/sync/semaphore.NewWeighted 用于创建新的信号量；
- x/sync/semaphore.Weighted.Acquire 阻塞地获取指定权重的资源，如果当前没有空闲资源，就会陷入休眠等待；
- x/sync/semaphore.Weighted.TryAcquire 非阻塞地获取指定权重的资源，如果当前没有空闲资源，就会直接返回 false；
- x/sync/semaphore.Weighted.Release 用于释放指定权重的资源；

在使用过程中需要注意以下几个问题
- x/sync/semaphore.Weighted.Acquire 和 x/sync/semaphore.Weighted.TryAcquire 方法都可以用于获取资源，前者会阻塞地获取信号量，后者会非阻塞地获取信号量；
- x/sync/semaphore.Weighted.Release 方法会按照 FIFO 的顺序唤醒可以被唤醒的 Goroutine；
- 如果一个 Goroutine 获取了较多地资源，由于 x/sync/semaphore.Weighted.Release 的释放策略可能会等待比较长的时间

#### SingleFlight
这个是Go语言的扩展包中提供的另外一个信号量，它能够在一个服务中抑制对下游的多次重复请求，比如在redis的缓存雪崩中，能够限制对同一个 Key 的多次重复请求，减少对下游的瞬时流量。

在资源获取非常昂贵的时候，就很适合使用x/sync/singleflight.Group
```
type service struct {
    requestGroup singleflight.Group
}

func (s *service) handleRequest(ctx context.Context, request Request) (Response error){
    v,err,_  := requestGroup.Do(request.Hash(),func() (interface{},error) {
        rows, err := // select * from tables
        if err != nil {
            return nil, err
        }
    })

    if err != nil{
        return nil,err
    }
    return Response {
        rows:rows,
    },nil

}
```
因为请求的哈希在业务上一般表示相同的请求，所以上述代码使用它作为请求的键。当然，我们也可以选择其他的唯一字段作为 x/sync/singleflight.Group.Do 方法的第一个参数减少重复的请求。









