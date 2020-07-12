
## defer学习

很多现代的变成语言中都会有defer关键字，Go语言的defer会在当前函数或是方法返回之前执行传入的函数，它会经常被用于
关闭文件描述符，关闭数据库链接和解锁资源。
作为一个编程语言中的关键字，defer 的实现一定是由编译器和运行时共同完成的，
不过在深入源码分析它的实现之前我们还是需要了解 defer 关键字的常见使用场景以及使用时的注意事项。
1. 使用defer的最常见的场景就是在函数调用结束的时候完成一些收尾工作，比如在defer中回滚数据库
```
func createPost(db *gorm.DB) error{
    tx := db.Begin()
    defer tx.Rollback()
    if err := tx.Create(&Post{Author: "Draveness"}).Error; err != nil{
        return err
    }
    return tx.Commit().Error

}
```
### 现象
我们在 Go 语言中使用 defer 时会遇到两个比较常见的问题，这里会介绍具体的场景并分析这两个现象背后的设计原理：
1. defer 关键字的调用时机以及多次调用 defer 时执行顺序是如何确定的；
2. defer 关键字使用传值的方式传递参数时会进行预计算，导致不符合预期的结果；

#### 作用域
向 defer 关键字传入的函数会在函数返回之前运行。假设我们在 for 循环中多次调用 defer 关键字：
defer的执行顺序是从后向前执行，先定义后执行
```
func main(){
    for i:=0;i<5;i++{
        defer fmt.Println(i)
    }
}
打印结果为：4，3，2，1

```
defer 传入的函数不是在退出代码块的作用域时执行的，它只会在当前函数和方法返回之前被调用。
```
func main(){
    {
        defer fmt.Println("defer runs")
        fmt.Println("block ends")
    }
    fmt.Println("main ends")
}
执行结果为：
block ends
main ends
defer runs           
```
#### 预计算参数
Go 语言中所有的函数调用都是传值的，defer 虽然是关键字，但是也继承了这个特性。假设我们想要计算 main 函数运行的时间，可能会写出以下的代码
```
func main() {
	startedAt := time.Now()
	defer fmt.Println(time.Since(startedAt))
	time.Sleep(time.Second)
}
执行结果为0s
```
调用 defer 关键字会立刻对函数中引用的外部参数进行拷贝，所以 time.Since(startedAt) 的结果不是在 main 函数退出之前计算的，
而是在 defer 关键字调用时计算的，最终导致上述代码输出 0s。解决这个方法只通过传参数的方式就可以解决。

### 数据结构
defer 关键字在 Go 语言源代码中对应的数据结构：
```
type _defer struct {
    siz     int32
    started bool
    sp      uintptr
    pc      uintptr
    fn      *funcval
    _panic  *_panic
    link    *_defer
}
```
runtime._defer 结构体是延迟调用链表上的一个元素，所有结构体都会通过link字段串链成链表
1. siz 是参数和结果的内存大小
2. sp和pc分别代表栈指针和调用方的程序技术器
3. fn是defer关键字中传入的函数
4. _panic 是触发延迟调用的结构体

### 编译过程
defer 关键字在运行期间会调用 runtime.deferproc 函数，这个函数接收了参数的大小和闭包所在的地址两个参数。
编译器不仅将 defer 关键字都转换成 runtime.deferproc 函数，它还会通过以下三个步骤为所有调用 
defer 的函数末尾插入 runtime.deferreturn 的函数调用：
1. cmd/compile/internal/gc.walkstmt 在遇到 ODEFER 节点时会执行 Curfn.Func.SetHasDefer(true) 设置当前函数的 hasdefer；
2. cmd/compile/internal/gc.buildssa 会执行 s.hasdefer = fn.Func.HasDefer() 更新 state 的 hasdefer；
3. cmd/compile/internal/gc.state.exit 会根据 state 的 hasdefer 在函数返回之前插入 runtime.deferreturn 的函数调用；

### 编译过程
defer 关键字的运行时实现了分成两部分
1. runtime.deferproc 函数负责创建新的延迟调用
2. runtime.deferreturn 函数负责在函数调用结束时执行所有的延迟调用
#### 创建延迟调用
runtime.deferproc 会为 defer 创建一个新的 runtime._defer 结构体、设置它的函数指针 fn、
程序计数器 pc 和栈指针 sp 并将相关的参数拷贝到相邻的内存空间中：
```
func deferproc(siz int32,fn *funcval){
    sp := getcallersp()
    argp := uintptr(unsafe.Pointer(%fn)) + unsafe.Sizeof(fn)
    callerpc := getcallerpc()
    d := newdefer(siz)
    if d._panic != nil {
        throw("deferproc: d.panic != nil after newdefer")
    }
    d.fn =fn
    d.pc = callerpc
    d.sp = sp
    switch siz {
    case 0:
    case sys.PtrSize:
        *(*uintptr)(deferArgs(d)) = *(*uintptr)(unsafe.Pointer(argp))
    default:
        memmove(deferArgs(d),unsafe.Pointer(argp),uintptr(siz))
    }
    return0()
}

```
最后调用的 runtime.return0 函数的作用是避免无限递归调用 runtime.deferreturn，它是唯一一个不会触发由延迟调用的函数了。

runtime.deferproc 中 runtime.newdefer 的作用就是想尽办法获得一个 runtime._defer 结构体，办法总共有三个：
1. 从调度器的延迟调用缓存池 sched.deferpool 中取出结构体并将该结构体追加到当前的Goroutine的缓存池中
2. 从Goroutine的延迟调用缓存池的pp.deferpool中取出结构体
3. 通过runtime.mallocgc创建一个新的结构体
```
func newdefer(siz int32) *_defer {
	var d *_defer
	sc := deferclass(uintptr(siz))
	gp := getg()
	if sc < uintptr(len(p{}.deferpool)) {
		pp := gp.m.p.ptr()
		if len(pp.deferpool[sc]) == 0 && sched.deferpool[sc] != nil {
			for len(pp.deferpool[sc]) < cap(pp.deferpool[sc])/2 && sched.deferpool[sc] != nil {
				d := sched.deferpool[sc]
				sched.deferpool[sc] = d.link
				pp.deferpool[sc] = append(pp.deferpool[sc], d)
			}
		}
		if n := len(pp.deferpool[sc]); n > 0 {
			d = pp.deferpool[sc][n-1]
			pp.deferpool[sc][n-1] = nil
			pp.deferpool[sc] = pp.deferpool[sc][:n-1]
		}
	}
	if d == nil {
		total := roundupsize(totaldefersize(uintptr(siz)))
		d = (*_defer)(mallocgc(total, deferType, true))
	}
	d.siz = siz
	d.link = gp._defer
	gp._defer = d
	return d
}
```
无论使用哪种方式获取 runtime._defer，它都会被追加到所在的 Goroutine _defer 链表的最前面。

defer 关键字插入时是从后向前的，而 defer 关键字执行是从前向后的，而这就是后调用的 defer 会优先执行的原因。

#### 执行延迟调用

runtime.deferreturn 函数会多次判断当前 Goroutine 的 _defer 链表中是否有未执行的剩余结构，
在所有的延迟函数调用都执行完成之后，该函数才会返回。

### 小结
defer 关键字的实现主要是依靠编译器和运行时的协作
1. 编译期
   1. 将defer 关键字被替换成runtime.deferproc
   2. 在调用defer关键字的函数返回之前插入runtime.deferreturn
2. 运行时
   1. runtime.deferproc 会将一个新的 runtime._defer 结构体追加到当前 Goroutine 的链表头；
   2. runtime.deferreturn 会从 Goroutine 的链表中取出 runtime._defer 结构并依次执行；












