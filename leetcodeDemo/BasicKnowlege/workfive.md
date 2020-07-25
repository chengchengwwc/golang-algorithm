### 网络轮询器

当前大多数的服务都是IO密集形的，应用程序需要花大量的时间进等待I/O操作的完成，网络轮询机制就是Go语言在运行的时候用来处理I/O操作的关键组件，它使用了操作系统提供的 I/O 多路复用机制增强程序的并发处理能力。

#### 设计原理
网络轮询器不仅仅只是用于监控网络I/O，还能用于监控文件的I/O，它利用了操作系统提供的I/O多路复用模型来提升设备的利用率，以及程序的性能
##### I/O模型
操作系统中包含阻塞 I/O、非阻塞 I/O、信号驱动 I/O 与异步 I/O 以及 I/O 多路复用五种 I/O 模型。我们在本节中会介绍上述五种模型中的三种
1. 阻塞 I/O 模型
2. 非阻塞 I/O 模型
3. I/O 多路复用模型
###### 阻塞I/O
阻塞 I/O 是最常见的 I/O 模型，对文件和网络的读写操作在默认情况下都是阻塞的。当我们通过 read 或者 write 等系统调用对文件进行读写时，应用程序就会被阻塞
###### 非阻塞 I/O
当进程把一个文件描述符设置成非阻塞时，执行 read 和 write 等 I/O 操作就会立刻返回。在 C 语言中，我们可以使用如下所示的代码片段将一个文件描述符设置成非阻塞的。第一次从文件描述符中读取数据会触发系统调用并返回 EAGAIN 错误，EAGAIN 意味着该文件描述符还在等待缓冲区中的数据；随后，应用程序会不断轮询调用 read 直到它的返回值大于 0，这时应用程序就可以对读取操作系统缓冲区中的数据并进行操作。进程使用非阻塞的 I/O 操作时，可以在等待过程中执行其他的任务，增加 CPU 资源的利用率。
###### I/O 多路复用
I/O 多路复用被用来处理同一个事件循环中的多个 I/O 事件。I/O 多路复用需要使用特定的系统调用，最常见的系统调用就是 select，该函数可以同时监听最多 1024 个文件描述符的可读或者可写状态；除了标准的 select 函数之外，操作系统中还提供了一个比较相似的 poll 函数，它使用链表存储文件描述符，摆脱了 1024 的数量上限。
多路复用函数会阻塞的监听一组文件描述符，当文件描述符的状态转变为可读或者可写时，select 会返回可读或者可写事件的个数，应用程序就可以在输入的文件描述符中查找哪些可读或者可写，然后执行相应的操作。
##### 多模块
Go 语言在网络轮询器中使用 I/O 多路复用模型处理 I/O 操作，为了提高 I/O 多路复用的性能，不同的操作系统也都实现了自己的 I/O 多路复用函数，例如：epoll、kqueue 和 evport 等。Go 语言为了提高在不同操作系统上的 I/O 操作性能，使用平台特定的函数实现了多个版本的网络轮询模块。
- src/runtime/netpoll_epoll.go
- src/runtime/netpoll_kqueue.go
- src/runtime/netpoll_solaris.go
- src/runtime/netpoll_windows.go
- src/runtime/netpoll_aix.go
- src/runtime/netpoll_fake.go
epoll、kqueue、solaries 等多路复用模块都要实现以下五个函数，这五个函数构成一个虚拟的接口：
```
func netpollinit(): 初始化网络轮询器，通过 sync.Once 和 netpollInited 变量保证函数只会调用一次；
func netpollopen(fd uintptr,pd *pollDesc) int32: 监听文件描述符上的边缘触发事件，创建事件并加入监听
func netpoll(delta int64) gList:轮询网络并返回一组已经准备就绪的 Goroutine，传入的参数会决定它的行为
func netollBreak() 唤醒网络轮询器，例如：计时器向前修改时间时会通过该函数中断网络轮询器
func netpollIsPollDescriptor(fd uintptr) bool 判断文件描述符是否被轮询器使用
```
#### 多路复用
网络轮询器实际上就是对 I/O 多路复用技术的封装
1. 网络轮询器的初始化
   - internal/poll.pollDesc.init — 通过 net.netFD.init 和 os.newFile 初始化网络 I/O 和文件 I/O 的轮询信息时；
   - runtime.doaddtimer — 向处理器中增加新的计时器时；
2. 向网络轮询器中加入待监控的任务
3. 从网络轮询器中获取触发的事件
   - 当我们在文件描述符上执行读写操作时，如果文件描述符不可读或者不可写，当前 Goroutine 就会执行 runtime.poll_runtime_pollWait 检查 runtime.pollDesc 的状态并调用 runtime.netpollblock 等待文件描述符的可读或者可写
   -  Go 语言的运行时会在调度或者系统监控中调用 runtime.netpoll 轮询网络，该函数的执行过程可以分成以下几个部分
      - 根据传入的 delay 计算 epoll 系统调用需要等待的时间；
      - 调用 epollwait 等待可读或者可写事件的发生；
      - 在循环中依次处理 epollevent 事件；
4. 截止日期
   - 截止日期在 I/O 操作中，尤其是网络调用中很关键，网络请求存在很高的不确定因素，我们需要设置一个截止日期保证程序的正常运行，这时就需要用到网络轮询器中的 runtime.poll_runtime_pollSetDeadline 函数。该函数会先使用截止日期计算出过期的时间点，然后根据 runtime.pollDesc 的状态做出以下不同的处理。
      - 如果结构体中的计时器没有设置执行的函数时，该函数会设置计时器到期后执行的函数、传入的参数并调用 runtime.resettimer 重置计时器。
      - 如果结构体的读截止日期已经被改变，我们会根据新的截止日期做出不同的处理
         - 如果新的截止日期大于 0，调用 runtime.modtimer 修改计时器
         - 如果新的截止日期小于 0，调用 runtime.deltimer 删除计时器
