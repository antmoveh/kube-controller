#### kcontrolcrd
>- 实现kubernetes自定义资源

- 自定义资源Unit，将Ingress StatefulSet Service PVC 关系强绑定为应用

```
# 自定义
export CRD=Unit
export group=custom
export version=v1

mkdir -p CRD/${CRD} && cd CRD/${CRD}

export GO111MODULE=on

# 如果路径位于GOPATH/src下，go mod这一步可省略
go mod init ${CRD}

# domian可自定义
kubebuilder init --domain unit.crd.com

# 为CRD生成API groupVersion
# kubebuilder create api --group custom --version v1 --kind Unit
kubebuilder create api --group ${group} --version ${version} --kind ${CRD}
# 创建Adminssion Webhook
kubebuilder create webhook --group custom --version v1 --kind Unit --defaulting --programmatic-validation
```