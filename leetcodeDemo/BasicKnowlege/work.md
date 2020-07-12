### 面试笔记


## Golang的调度器
1. 谈到Golang的调度器，绕不开的是操作系统，进程和线程这些概念。多个线程是可以属于同一个进程的并共享内存空间，因为多线程
不需要创建新的虚拟空间，所以不需要内存管理单元处理的上下文的切换，线程之间的通信也是基于共享内存进行的，同重量级的进程相比
线程显得比较轻量
2. 虽然线程比较轻量，但是线程每一次的切换需要耗时1us左右的时间，但是Golang调度器对于goroutine的切换只要在0.2us
左右
3. Go 语言的调度器通过使用与 CPU 数量相等的线程减少线程频繁切换的内存开销，同时在每一个线程上执行额外开销更低的 Goroutine 来降低操作系统和硬件的负载。
#### 调度器种类
1. 单线程调度器：遵循如下调度过程
- 
2. 多线程调度器
3. 任务窃取调度器
4. 抢占式调度器
- 基于协作的抢占式调度器
- 基于信号的抢占式调度器：实现基于信号的抢占式调度器，垃圾回收在扫描栈时会触发抢占式调度；抢占的时间不够多，不能覆盖全部边缘情况；
  - 挂起goroutine的过程是在垃圾回收的栈扫描时候来完成的
  - 调用runtime.suspendG函数时会将处于运行状态的goroutine的preemptStop标记成为true
  - 调用runtime.preemptPark函数可以挂起当前的goroutine 将其状态更新为_Gpreemted并触发调度器的重新调度，该函数能够交出线程控制权
  - 在X86架构上增加异步抢占函数
  - 支持通过向线程发送信号的方式暂停运行的 Goroutine；
  - 在 runtime.sighandler 函数中注册 SIGURG 信号的处理函数 runtime.doSigPreempt
  - 实现 runtime.preemptM 函数，它可以通过 SIGURG 信号向线程发送抢占请求；
  - 修改 runtime.preemptone 函数的实现，加入异步抢占的逻辑；
- 目前的抢占式调度也只会在垃圾回收扫描任务时触发，我们可以梳理一下上述代码实现的抢占式调度过程
  - 程序启动时，在 runtime.sighandler 函数中注册 SIGURG 信号的处理函数 runtime.doSigPreempt；
  - 在触发垃圾回收的栈扫描时会调用 runtime.suspendG 挂起 Goroutine，该函数会执行下面的逻辑：
    1. 将 _Grunning 状态的 Goroutine 标记成可以被抢占，即将 preemptStop 设置成 true；
    2. 调用 runtime.preemptM 触发抢占；
  - runtime.preemptM 会调用 runtime.signalM 向线程发送信号 SIGURG；
  - 操作系统会中断正在运行的线程并执行预先注册的信号处理函数 runtime.doSigPreempt；
  - runtime.doSigPreempt 函数会处理抢占信号，获取当前的 SP 和 PC 寄存器并调用
  - runtime.sigctxt.pushCall 会修改寄存器并在程序回到用户态时执行
  - 汇编指令 runtime.asyncPreempt 会调用运行时函数 runtime.asyncPreempt2；
  - runtime.asyncPreempt2 会调用 runtime.preemptPark；
  - runtime.preemptPark 会修改当前 Goroutine 的状态到 _Gpreempted 并调用
  - runtime.schedule 让当前函数陷入休眠并让出线程，调度器会选择其它的 Goroutine 继续执行；

### 数据结构
##### G
1. 表示Goroutine 是一个等待执行的任务
2. 它只存在于Go语言的运行时，它是Go语言在用户态提供的线程，作为一种粒度更细的资源调度单元，如果使用得当能够在
高并发的场景下更加高效的利用机器CPU.
3. goroutine在运行的时候会使用私有结构体runtine.g表示，下面对具体的字段进行解释
   1. stack 字段描述了当前 Goroutine 的栈内存范围 [stack.lo, stack.hi)
   2. stackguard0 可以用于调度器抢占式调度
   3. m — 当前 Goroutine 占用的线程，可能为空
   4. atomicstatus — Goroutine 的状态；
   5. sched — 存储 Goroutine 的调度相关的数据；
      1. sched — 存储 Goroutine 的调度相关的数据；
      2. pc — 程序计数器（Program Counter）；
      3. g — 持有 runtime.gobuf 的 Goroutine
      4. ret — 系统调用的返回值
 4. goroutine的状态：主要有三种状态：等待中，可运行，运行中
    1. 等待中：Goroutine 正在等待某些条件满足，例如：系统调用结束等，包括 _Gwaiting、_Gsyscall 和 _Gpreempted 几个状态
    2. 可运行：Goroutine 已经准备就绪，可以在线程运行，如果当前程序中有非常多的 Goroutine，每个 Goroutine 就可能会等待更多的时间，即 _Grunnable；
    3. 运行中：Goroutine 正在某个线程上运行，即 _Grunning；

##### M
Go 语言并发模型中的 M 是操作系统线程。调度器最多可以创建 10000 个线程，但是其中大多数的线程都不会执行用户代码（可能陷入系统调用），最多只会有 GOMAXPROCS 个活跃线程能够正常运行。
在默认情况下，运行时会将 GOMAXPROCS 设置成当前机器的核数，我们也可以使用 runtime.GOMAXPROCS 来改变程序中最大的线程数。
在默认情况下，一个四核机器上会创建四个活跃的操作系统线程，每一个线程都对应一个运行时中的 runtime.m 结构体。
在大多数情况下，我们都会使用 Go 的默认设置，也就是线程数等于 CPU 个数，在这种情况下不会触发操作系统的线程调度和上下文切换，所有的调度都会发生在用户态，由 Go 语言调度器触发，能够减少非常多的额外开销。
1. g0 是持有调度栈的 Goroutine，curg 是在当前线程上运行的用户 Goroutine，这也是操作系统线程唯一关心的两个 Goroutine
   1. g0 是一个运行时中比较特殊的 Goroutine，它会深度参与运行时的调度过程，包括 Goroutine 的创建、大内存分配和 CGO 函数的执行
   
   
##### P
调度器中的处理器 P 是线程和 Goroutine 的中间层，它能提供线程需要的上下文环境，也会负责调度线程上的等待队列，通过处理器 P 的调度，每一个内核线程都能够执行多个 Goroutine，它能在 Goroutine 进行一些 I/O 操作时及时切换，提高线程的利用率。
因为调度器在启动时就会创建 GOMAXPROCS 个处理器，所以 Go 语言程序的处理器数量一定会等于 GOMAXPROCS，这些处理器会绑定到不同的内核线程上并利用线程的计算资源运行 Goroutine。

#### 调度器启动
1. 调度器通过 runtime.schedinit 函数初始化调度器：
2. 在调度器初始函数执行的过程中会将 maxmcount 设置成 10000，这也就是一个 Go 语言程序能够创建的最大线程数，虽然最多可以创建 10000 个线程，但是可以同时运行的线程还是由 GOMAXPROCS 变量控制。
3. 从环境变量 GOMAXPROCS 获取了程序能够同时运行的最大处理器数之后就会调用 runtime.procresize 更新程序中处理器的数量，在这时整个程序不会执行任何用户 Goroutine，调度器也会进入锁定状态，runtime.procresize 的执行过程如下：
   1. 如果全局变量 allp 切片中的处理器数量少于期望数量，就会对切片进行扩容；
   2. 使用 new 创建新的处理器结构体并调用 runtime.p.init 方法初始化刚刚扩容的处理器；
   3. 通过指针将线程m0同处理器allp[0]绑定到提起
   4. 调用runtime.p.destroy 方法释放不再使用的处理器结构；
   5. 通过截断改变全局变量 allp 的长度保证与期望处理器数量相等；
   6. 将除 allp[0] 之外的处理器 P 全部设置成 _Pidle 并加入到全局的空闲队列中；
   7. 调用 runtime.procresize 就是调度器启动的最后一步，在这一步过后调度器会完成相应数量处理器的启动，等待用户创建运行新的 Goroutine 并为 Goroutine 调度处理器资源。

#### 创建Goroutine
想要启动一个新的goroutine来执行任务，我们需要将Go语言中的go关键字，这个关键字会在编译期间通过下面方法cmd/compile/internal/gc.state.stmt 和 cmd/compile/internal/gc.state.call 两个方法将该关键字转换成 runtime.newproc 函数调用：
1. 编译器会将所有的go关键字转换为runtime.newproc 函数，该函数会接受大小和表示函数的指针funcval。在这个函数中我们还会
获取goroutine以及调用方的程序计数器，然后调用 runtime.newproc1 函数。runtime.newproc1 会根据传入参数初始化一个 g 结构体，我们可以将该函数分成以下几个部分介绍它的实现：
   1. 获取或者创建新的Groutine结构体
   2. 将传入的参数移植到Goroutine的栈上
   3. 更新Goroutine的调度相关性
   4. 将Goroutine加入处理器队列
    
##### 初始化结构体
- runtime.gfget通过两种不同的方式获取新的 runtime.g 结构体：
  - 从Goroutine所在的处理器的gFree列表或者调度器的sched.gFree 列表中获取 runtime.g 结构体；
  - 调用 runtime.malg 函数生成一个新的 runtime.g 函数并将当前结构体追加到全局的 Goroutine 列表 allgs 中。
-  runtime.gfget 中包含两部分逻辑，它会根据处理器中 gFree 列表中 Goroutine 的数量做出不同的决策：
   - 当处理器的 Goroutine 列表为空时，会将调度器持有的空闲 Goroutine 转移到当前处理器上，直到 gFree 列表中的 Goroutine 数量达到 32；
   - 当处理器的 Goroutine 数量充足时，会从列表头部返回一个新的 Goroutine；
- runtime.newproc1 会从处理器或者调度器的缓存中获取新的结构体，也可以调用 runtime.malg 函数创建新的结构体。

##### 运行队列
runtime.runqput 函数会将新创建的 Goroutine 运行队列上，这既可能是全局的运行队列，也可能是处理器本地的运行队列：
1. 当 next 为 true 时，将 Goroutine 设置到处理器的 runnext 上作为下一个处理器执行的任务；
2. 当 next 为 false 并且本地运行队列还有剩余空间时，将 Goroutine 加入处理器持有的本地运行队列；
3. 当处理器的本地运行队列已经没有剩余空间时就会把本地队列中的
一部分 Goroutine 和待加入的 Goroutine 通过 runqputslow 添加到调度器持有的全局运行队列上；
4. Go 语言中有两个运行队列，其中一个是处理器本地的运行队列，另一个是调度器持有的全局运行队列，只有在本地运行队列没有剩余空间时才会使用全局队列

##### 调度循环
调度器启动之后，Go 语言运行时会调用 runtime.mstart 以及 runtime.mstart1，前者会初始化 g0 的 stackguard0 和 stackguard1 字段，后者会初始化线程并调用 runtime.schedule 进入调度循环：
1. 为了保证公平，当全局运行队列中有待执行的 Goroutine 时，通过 schedtick 保证有一定几率会从全局的运行队列中查找对应的 Goroutine；
2. 从处理器本地的运行队列中查找待执行的 Goroutine；
3. 如果前两种方法都没有找到 Goroutine，就会通过 runtime.findrunnable 进行阻塞地查找 Goroutine；

##### 触发调度
运行时还会在线程启动 runtime.mstart 和 Goroutine 执行结束 runtime.goexit0 触发调度。我们在这里会重点介绍运行时触发调度的几个路径：
- 主动挂起 — runtime.gopark -> runtime.park_m
- 系统调用 — runtime.exitsyscall -> runtime.exitsyscall0
- 协作式调度 — runtime.Gosched -> runtime.gosched_m -> runtime.goschedImpl
- 系统监控 — runtime.sysmon -> runtime.retake -> runtime.preemptone

##### 线程管理
Go 语言的运行时会通过调度器改变线程的所有权，它也提供了 runtime.LockOSThread 和 runtime.UnlockOSThread 让我们有能力绑定 Goroutine 和线程完成一些比较特殊的操作。Goroutine 应该在调用操作系统服务或者依赖线程状态的非 Go 语言库时调用 runtime.LockOSThread 函数11，例如：C 语言图形库等。
1.  runtime.dolockOSThread 会分别设置线程的 lockedg 字段和 Goroutine 的 lockedm 字段，这两行代码会绑定线程和 Goroutine。
2.  当 Goroutine 完成了特定的操作之后，就会调用以下函数 runtime.UnlockOSThread 分离 Goroutine 和线程：









  
  


