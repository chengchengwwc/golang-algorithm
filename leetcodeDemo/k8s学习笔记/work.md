#### Cluster 集群、
Cluster是计算，存储和网络资源的集合，kubernetes是利用这些资源运行各种基于容器的应用

#### master 控制节点
Master是Cluster的大脑，主要指责是调度，决定将应用放在哪里运行，为了实现高可用，可以运行多个master，调度应用程序，维护应用程序的所需状态，扩展应用程序和滚动更新都是master的工作

#### Node 节点
Node节点是运行容器的应用，Node是由master管理，Node负责监控并汇报容器的状态，同时根据master的要求管理容器的生命周期。
Node是kubernetes集群中的工作机器，可以是物理机器，也可以是虚拟机，每个工作节点都有一个kubelet，它是管理节点并同kubernetes master进行通信，节点上还应该具有处理容器操作的容器运行时。
一个 kubernetes工作集群至少包括三个节点，master管理集群，而node 用于托管正在运行的应用程序。
#### pod 资源对象
pod 是k8s中最小的工作单元，每一个pod包含一个或是多个容器，Pod中的容器会作为一个整体被master调度到一个node上运行。
k8s引入pod主要是基于下面两个目的
1. 可管理性能，有些容器天生需要紧密的联系在一起，一起工作，Pod提供了比容器更高层次的抽象，将相关容器封装到同一个部署单元中。
2. 通信和资源共享：pod中的容器使用同一个网络namespace,可以直接使用localhost进行通信，同样，这些容器可以共享存储， 当kubernetes挂载volume到pod上的时候，本质上是将volume挂载到Pod中的每一个容器中。

pod有下面两种运行方式
1. 运行单一容器
2. 运行多个容器

#### controller 控制器
kubernetes 通常不会直接创建pod，而是通过controller来管理pod ，controller中定义了pod的部署特性，比如有几个副本，在什么样的node下运行等，为了满足不同的业务场景，kubernetes提供了多种crontroler, 比如deployment,replicaSet daemonSet,StatefuleSet,job等
- deployment 用于管理pod的多个副本，并确保pod按照期望的状态运行。
- replicaSet 实现pod的多副本管理，使用deployment时候会自动创建replicaSet
- DaemonSet 用于每个node最多只运行一个副本的场景
- StatefuleSet可以保证pod的每个副本在生命周期里名称不会变
- job 用于运行结束后删除的应用。

#### service 服务
deployment 可以部署多个版本，外部通过service来对pod进行访问

#### namespace 命名空间
namespace 是对一组资源和对象的抽象组合，比如可以用来将系统内部的对象划分为不同的项目组或是产品组，如果有多个用户或是项目组使用同一个kubernetes cluster，这个主要是通过namespace来进行划分
- namespace 可以将一个物理的cluster逻辑上划分成多个cluster，每一个cluster就是一个namespace，不同的namespace里资源是完全隔离里的
- kubernetes 默认创建了两个namespace
- default: 创建资源时候不确定的时候，将被放到这个namespace中去
- kube-system: kubernetes自己创建的系统资源将放到这个namespace中去。





