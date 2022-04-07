


#### multi crd version 

##### 1. 创建CRD

```shell
# 选择已经存在v1beta1版本的，我们
$ kubebuilder create api --group study --version v1beta2 --kind World
```

##### 2. storage version

```shell
# 
```

##### 3. generate crd yaml

```shell
$ make kustomize
# multi group 
$ kustomize build config/crd > crd.yaml
```


##### 4. cert-manager

```shell
# cert
```


#### 注意事项

```shell
# 如果你是从CRD单分组升级到多分组，注意迁移api下所有文件到apis下

```


#### 参考文章
  - [kubebuilder教程](https://cloudnative.to/kubebuilder/multiversion-tutorial/tutorial.html)