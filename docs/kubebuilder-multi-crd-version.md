


#### 多CRD版本

##### 1. 创建CRD

```shell
# 选择已经存在v1beta1版本的，我们
$ kubebuilder create api --group study --version v1beta2 --kind World
```

##### 2. 多版本转换问题

-  当在相同分组下存在多个版本时候，需要提供一个存储版本，这个可以理解为需要提供一个基线版本，etcd中便存储这个基线版本。
- 比如world.study.example.cn存在v1beta1和v1beta2版本，此时我们规定v1beta1为存储版本，即在其api定义中增加 `// +kubebuilder:storageversion`注解

```shell
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:storageversion

type Job struct {
    metav1.TypeMeta   `json:",inline"`
    metav1.ObjectMeta `json:"metadata,omitempty"`
}
```

- 存储版本需要实现Hub方法`apis/study/v1beta1/world_conversion.go`，其他版本需要实现向存储版本的转换方法`apis/study/v1beta2/world_conversion.go`
- 假设有v1beta1、v1beta2和v1beta3它们之间的转换是通过存储版本转换的，并不会存在v1beta2和v1beta3直接转换方法，k8s资源多版本兼容也是如此实现的。这样引申出一个特性在我们的控制器中只要监听任意版本的资源均可，因为我们实现了CRD的转换方法
- CRD的转换方法是一个webhook，这个webhook的service定义在crd.yaml `deploy/crd.yaml`中

```yaml
  conversion:
    strategy: Webhook
    webhook:
      clientConfig:
        service:
          name: webhook-service
          namespace: system
          path: /convert
      conversionReviewVersions:
      - v1
```

##### 3.  实现convert的webhook

```shell
$ kubebuilder create webhook --group study --version v1beta1 --kind World --conversion
```

- 生成的webhook路径为`apis/study/v1beta1/world_webhook.go` 

##### 4. 其他方式实现convert的webhook

- 使用kubebuilder模板方法生成的总归有点绕，自己写方法实现convert方法

```shell
# 待补充
```

##### 5. 生成crd

Kubebuilder 在 `config` 目录下生成禁用 webhook bits 的 Kubernetes 清单。要启用它们，我们需要：

- 在 `config/crd/kustomization.yaml` 文件启用 `patches/webhook_in_<kind>.yaml` 和 `patches/cainjection_in_<kind>.yaml`。
- 在 `config/default/kustomization.yaml` 文件的 `bases` 部分下启用 `../certmanager` 和 `../webhook` 目录。
- 在 `config/default/kustomization.yaml` 文件的 `patches` 部分下启用 `manager_webhook_patch.yaml`。
- 在 `config/default/kustomization.yaml` 文件的 `CERTMANAGER` 部分下启用所有变量
- 因为我们不使用CERTMANGER可以不启动相关功能

```shell

# 生成crd，这里没有patch webhook服务
$ make manifests
# 安装
$ make kustomize
# 这里会patch webhook服务到crd yaml中
$ kustomize build config/crd > crd.yaml
```

##### 6. 部署测试

```shell
$ kubectl apply -f deploy/crd.yaml
$ kubectl apply -f world.v1beta1.study.example.cn.yaml
# 如果没有部署convert webhook该命令会执行失败
$ kubectl apply -f world.v1beta2.study.example.cn.yaml
# 通常情况下我们使用kubectl get kindname 查看所有版本资源，我们可以使用如下命令查看指定版本资源
$ kubectl get world.v1beta1.study.example.cn
$ kubectl get world.v1beta2.study.example.cn
# 这里只要你部署了v1beta2版本的，因为存在转换webhook，你查询v1beta1的也可以查询到创建的资源
```

#### 参考文章

  - [kubebuilder教程](https://cloudnative.to/kubebuilder/multiversion-tutorial/tutorial.html)