#### 使用kubebuilder创建adminwebhook

##### 1. 创建adminwebhook

```

# 创建Adminssion Webhook
kubebuilder create webhook --group custom --version v1 --kind Unit --defaulting --programmatic-validation

Default() 用于修改，即对应mutating webhook
ValidateCreate()用于校验，对应validating webhook
ValidateUpdate()用于校验，对应validating webhook
ValidateDelete()用于校验，对应validating webhook
```

##### 2. 创建adminwebhook部署yaml

```shell

```

##### 3. 生成adminwebhook证书

```shell
# 官方提供的webhook证书方位为部署certmanager，但是观察其他CRD服务均未使用certmanager的方案
# 参考istio 使用自签名证书
./yaml/gen_webhookca.sh --service kunit-webhook-service --namespace kunit-system --secret kunit
# 配置webhook caBundle
CA_BUNDLE=$(kubectl config view --raw --minify --flatten -o jsonpath='{.clusters[].cluster.certificate-authority-data}')
sed -i "s#\${CA_BUNDLE}#${CA_BUNDLE}#g" yaml/webhook.yaml
kubectl apply -f yaml/webhook.yaml
# 会创建csr secret，将secret挂载到容器/tmp/k8s-webhook-server/serving-certs即可使adminwebhook正常工作
```

##### 4. adminwebhook本地调试

```shell
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

#### 参考文章

[深入剖析MutatingAdmissionWebhook](https://blog.hdls.me/15564491070483.html )

[使用kubebuilder创建自定义k8s AdmissionWebhooks](https://blog.hdls.me/15708754600835.html )