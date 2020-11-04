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

Default() 用于修改，即对应mutating webhook
ValidateCreate()用于校验，对应validating webhook
ValidateUpdate()用于校验，对应validating webhook
ValidateDelete()用于校验，对应validating webhook
```

##### deploy and debug
```
# move kubectl and /root/.kube/config in this env
# install crd
make install 
# adminwebhooks need /tmp/k8s-webhook-server/serving-certs/tls.{crt,key}

# local debug
make run ENABLE_WEBHOOKS=false

# deploy
make docker-build docker-push IMG=docker.g.com/project-name:tag
make deploy IMG=docker.g.com/project-name:tag
```