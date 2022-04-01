#### kubebuilder 创建自定义资源过程

##### 1. 安装kubebuilder

```shell
$ os=$(go env GOOS)
$ arch=$(go env GOARCH)

# download kubebuilder and extract it to tmp
$ curl -sL https://go.kubebuilder.io/dl/2.0.0-beta.0/${os}/${arch} | tar -xz -C /tmp/

# move to a long-term location and put it on your path
# (you'll need to set the KUBEBUILDER_ASSETS env var if you put it somewhere else)
$ mv /tmp/kubebuilder_2.0.0-beta.0_${os}_${arch} /usr/local/kubebuilder
$ export PATH=$PATH:/usr/local/kubebuilder/bin
```



##### 2. 项目初始化与CRD资源创建

```shell
# 创建项目
$ mkdir project
$ go mod init project
# 初始化域名
# k8s的资源由GVK三个关键字段定位，这里的`project.com`便是Group中后半部分
$ kubebuilder init --domain project.com 
# 创建一个crd资源
# 这里个group和domain共同组成k8s资源的Group:`custom.project.com`, kind首字母大写
$ kubebuilder create api --group custom --version v1beta1 --kind Abbc
```

##### 3. 创建不同Group的CRD资源

```shell
# 首先开启多组CRD支持
$ kubebuilder edit --multigroup=true
# 创建CRD，会创建到apis目录下,这样生成的Group为common.scope.cluster.domain.cn，注意在代码中删除domain.cn
$ kubebuilder create api --group common.scope.cluster --version v1beta1 --kind Cluster


```

##### 4. 生成CRD控制器及部署资源

```shell
# 修改CRD字段后，使用如下命令生成CRD部署yaml
$ make manifests
```

##### 5. 部署

```shell
# 查看集群内所有的资源版本
$ kubectl api-resources
$ kubectl api-versions
$ kubectl get crd
# 安装crd
$ make install
```

#### 参考文章

[kubebuilder学习笔记](https://segmentfault.com/a/1190000020359577 )