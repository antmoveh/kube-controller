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

##### 官方部署方式
```
# 基本操作
kubectl api-resources
kubectl api-versions
kubectl get crd

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

#### 自定义部署方案
```
# 由于官方模块使用了kube-proxy-rbac镜像，在此我们不使用内置镜像的方法，而是采用serviceaccount
kubectl yaml/deployment.yaml
```

#### adminwebhook 证书
```
# 官方提供的webhook证书方位为部署certmanager，但是观察其他CRD服务均未使用certmanager的方案
# 参考istio 使用自签名证书
./yaml/gen_webhookca.sh --service kunit-webhook-service --namespace kunit-system --secret kunit
# 配置webhook caBundle
CA_BUNDLE=$(kubectl config view --raw --minify --flatten -o jsonpath='{.clusters[].cluster.certificate-authority-data}')
sed -i "s#\${CA_BUNDLE}#${CA_BUNDLE}#g" yaml/webhook.yaml
kubectl apply -f yaml/webhook.yaml
# 会创建csr secret，将secret挂载到容器/tmp/k8s-webhook-server/serving-certs即可使adminwebhook正常工作
```

#### adminwebhook本地调试
```
yaml/webook.yaml中
   service:
        name: kunit-webhook-service
        namespace: kunit-system
        path: /mutate-custom-unit-crd-com-v1-unit
        port: 443
更改为url: https://192.168.56.102:9443/mutate-custom-unit-crd-com-v1-unit

yaml/gen_webhookca.sh生成签名部分增加Ip
[alt_names]
IP.1 = 192.168.56.102
DNS.1 = ${service}
DNS.2 = ${service}.${namespace}
DNS.3 = ${service}.${namespace}.svc
# CN=换成IP
openssl req -new -key "${tmpdir}"/server-key.pem -subj "/CN=${service}.${namespace}.svc" -out "${tmpdir}"/server.csr -config "${tmpdir}"/csr.conf

第二种方案
将namespace/kunit-system的service/kunit-webhook-service的endpoint更改为你的本地地址，这样证书无需修改。
```
##### 参考文章
- https://segmentfault.com/a/1190000020359577
- https://blog.hdls.me/15564491070483.html
- https://blog.hdls.me/15564491070483.html
- https://blog.upweto.top/gitbooks/kubebuilder/%E6%9C%AC%E5%9C%B0%E8%B0%83%E8%AF%95%E5%92%8C%E5%8F%91%E5%B8%83Controller.html
