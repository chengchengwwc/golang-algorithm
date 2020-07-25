## Golang channel
作为Go的核心的数据结构和Goroutine之间的通信，是支撑Go语言高并发的关键
### 设计原理
Go 语言提供了一种不同的并发模型，也就是通信顺序进程（Communicating sequential processes，CSP）1。Goroutine 和 Channel 分别对应 CSP 中的实体和传递信息的媒介，Go 语言中的 Goroutine 会通过 Channel 传递数据。
#### 先入先出
目前Channel收发操作先入先出的设计
1. 先从Channel读取数据的Goroutine会先收到数据
2. 先向 Channel 发送数据的 Goroutine 会得到先发送数据的权利

#### 无锁管道
1. 无锁（lock-free）队列更准确的描述是使用乐观并发控制的队列。乐观并发控制也叫乐观锁，但是它并不是真正的锁，很多人都会误以为乐观锁是一种真正的锁，然而它只是一种并发控制的思想.
Channel在运行时候，包换了一个用于保护成员变量的互斥锁，Channel本质上是一个用于同步和通信的有锁队列，使用互斥锁解决程序中可能存在的线程竞争问题。
2. 锁会导致休眠和唤醒带来的上下文切换
   - 同步channel 不需要缓冲区，发送数据直接到接收方
   - 异步channel 基于环形存储的传统生产者和消费者模型
   - chan struct{} 类型的异步 Channel — struct{} 类型不占用内存空间，不需要实现缓冲区和直接发送（Handoff）的语义

#### 数据结构
 Go 语言的 Channel 在运行时使用 runtime.hchan 结构体表示。我们在 Go 语言中创建新的 Channel 时，实际上创建的都是如下所示的结构体
 ```
 type hchan struct {
     qcount   uint
     dataqsiz uint
     buf      unsafe.Pointer
     elemsize uint16
     closed   uint32
     elemtype *_type
     sendx    uint
     recvx    uint
     recvq    waitq
     sendq    waitq
     lock     mutex
 }
 ```
 - qcount    Channel中元素的个数
 - dataqsiz  Channel中循环队列的长度
 - buf       Channel的缓冲区的数据指针
 - sendx     Channel 的发送操作处理到的位置
 - recvx     Channel 的接收操作处理到的位置

 #### 创建管道
 Go 语言中所有 Channel 的创建都会使用 make 关键字。编译器会将 make(chan int, 10) 表达式被转换成 OMAKE 类型的节点。
 - 如果当前channel不存在缓冲区，那么就只会为 runtime.hchan 分配一段内存空间。
 - 如果当前 Channel 中存储的类型不是指针类型，就会为当前的 Channel 和底层的数组分配一块连续的内存空间
 - 在默认情况下会单独为 runtime.hchan 和缓冲区分配内存；

 #### 发送数据
 当我们想要向 Channel 发送数据时，就需要使用 ch <- i 语句，编译器会将它解析成 OSEND 节点并在 cmd/compile/internal/gc.walkexpr 函数中转换成 runtime.chansend1。
 在发送数据的逻辑执行之前会先为当前 Channel 加锁，防止发生竞争条件。如果 Channel 已经关闭，那么向该 Channel 发送数据时就会报"send on closed channel" 错误并中止程序。
 流程分为下面三部
 1. 当存在等待的接收者时候，通过runtime.send直接将数据发送给阻塞的接收者
 2. 当缓冲区存在空余空间的时候，将发送数据写入Channel的缓冲区
 3. 当不存在缓冲区或是缓冲区已满的时候，等待其他Goroutine从Channel接收数据

 如果目标 Channel 没有被关闭并且已经有处于读等待的 Goroutine，那么 runtime.chansend 函数会从接收队列 recvq 中取出最先陷入等待的 Goroutine 并直接向它发送数据。

 #### 阻塞发送
 当Channel中没有接收者能够处理数据的时候，向Channel发送数据就会被下游阻塞，当前使用select关键字可以向Channel非阻塞的发送消息，向Channel阻塞的发送数据会执行下面代码
 ```
 func chansend(c *hchan,ep unsafe.Pointer,block bool,callerpc uintptr) bool {
     if !block {
         unlock(&c.lock)
         return false
     }
     gp := getg()
     mysg := acquireSudog()
     mysg.elem = ep
     mysg.g = gp
     mysg.c = c
     gp.waiting = mysq
     c.sendq.enqueue(mysg)
     goparkunlock(&c.lock, waitReasonChanSend, traceEvGoBlockSend, 3)
     gp.waiting = nil
     gp.param = nil
     mysg.c = nil
     releaseSudog(mysg)
     return true
 }
 ```
 1. 调用runtime.getg 获取发送数据使用的 Goroutine；
 2. 执行 runtime.acquireSudog 函数获取 runtime.sudog 结构体并设置这一次阻塞发送的相关信息，例如发送的 Channel、是否在 Select 控制结构中和待发送数据的内存地址等；
 3. 将刚刚创建并初始化的 runtime.sudog 加入发送等待队列，并设置到当前 Goroutine 的 waiting 上，表示 Goroutine 正在等待该 sudog 准备就绪
 4. 调用 runtime.goparkunlock 函数将当前的 Goroutine 陷入沉睡等待唤醒
 5. 被调度器唤醒后会执行一些收尾工作，将一些属性置零并且释放 runtime.sudog 结构体

#### 接收数据
通过下面两个方式来接收数据
```
<- ch
ok <- ch
```
- 当存在等待的发送者的时候，通过runtime.recv直接从阻塞的发送者或者缓冲区中获得数据
- 当缓冲区存在数据的时候，从channel的缓冲区中接收数据
- 当缓冲区中不存在数据，等待

1. 直接接收
   - 当 Channel 的 sendq 队列中包含处于等待状态的 Goroutine 时，该函数会取出队列头等待的 Goroutine，处理的逻辑和发送时相差无几，只是发送数据时调用的是 runtime.send 函数，而接收数据时使用 runtime.recv 函数
      1. 如果不存在缓冲区，则会将数据拷贝到目标的内存中
      2. 如果存在缓冲区，会将队列中的数据拷贝到接收方的内存地址中
2. 缓冲区
   - 当channel的缓冲区中已经包含数据的时候，从channel中接收数据会直接从缓冲区中的索引位置读取数据并处理。
3. 阻塞接收
   - 当 Channel 的发送队列中不存在等待的 Goroutine 并且缓冲区中也不存在任何数据时，从管道中接收数据的操作会变成阻塞操作，然而不是所有的接收操作都是阻塞的，与 select 语句结合使用时就可能会使用到非阻塞的接收操作。

#### 关闭管道
编译器会将用于关闭管道的 close 关键字转换成 OCLOSE 节点以及 runtime.closechan 的函数调用。当 Channel 是一个空指针或者已经被关闭时，Go 语言运行时都会直接 panic 并抛出异















