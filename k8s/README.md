

# **🚀 微服务部署指南（Kubernetes 版）**



> 所有微服务配置将从挂载的 config/config.yaml 中读取，相关逻辑请参考：[config/config.go](https://github.com/west2-online/fzuhelper-server/blob/main/config/config.go)



## **环境初始化步骤**

### **创建命名空间**

```
kubectl create namespace fzuhelper
```

### **创建阿里云镜像拉取 Secret**

```
kubectl create secret docker-registry aliyun-registry-secret \
  -n fzuhelper \
  --docker-server=xxx \
  --docker-username=xxx \
  --docker-password=xxx \
  --docker-email=xxx
```

### **部署基础组件**

#### **部署 etcd**

```
helm install etcd bitnami/etcd -n fzuhelper -f ./k8s/etcd/etcd-values.yaml

kubectl get pods -n fzuhelper -l app.kubernetes.io/name=etcd
```

#### **部署 MySQL**

```
helm install mysql bitnami/mysql -n fzuhelper -f ./k8s/mysql/mysql-values.yaml

kubectl get pods -n fzuhelper -l app.kubernetes.io/name=mysql
```

#### **部署 Redis**

```
helm install redis bitnami/redis -n fzuhelper -f ./k8s/redis/redis-values.yaml
kubectl get pods -n fzuhelper -l app.kubernetes.io/name=redis
```

### **创建配置 ConfigMap**

记得修改 configmap 的内容

```
kubectl apply -f k8s/config/configmap.yaml
```

### **部署微服务集群**

```
helm install fzuhelper-server ./k8s/fzuhelper-server -n fzuhelper
```

最后别忘了部署 ingress

## **更新服务配置**

```
helm upgrade etcd bitnami/etcd -n fzuhelper -f ./k8s/etcd/etcd-values.yaml
helm upgrade mysql bitnami/mysql -n fzuhelper -f ./k8s/mysql/mysql-values.yaml
helm upgrade redis bitnami/redis -n fzuhelper -f ./k8s/redis/redis-values.yaml
helm upgrade fzuhelper-server ./k8s/fzuhelper-server -n fzuhelper
```

## **卸载所有服务**

```
helm uninstall etcd -n fzuhelper
helm uninstall mysql -n fzuhelper
helm uninstall redis -n fzuhelper
helm uninstall fzuhelper-server -n fzuhelper
kubectl delete namespace fzuhelper
```



## **微服务资源配置汇总**



| **服务**      | **CPU**  | **Memory** |
| ------------- | -------- | ---------- |
| api           | 30m      | 128Mi      |
| user          | 20m      | 128Mi      |
| classroom     | 20m      | 128Mi      |
| course        | 20m      | 128Mi      |
| launch-screen | 20m      | 32Mi       |
| paper         | 20m      | 32Mi       |
| academic      | 20m      | 128Mi      |
| version       | 20m      | 32Mi       |
| common        | 20m      | 64Mi       |
| **小计**      | **190m** | **800Mi**  |

## **数据组件资源配置汇总**



| **组件** | **CPU**  | **Memory**  |
| -------- | -------- | ----------- |
| MySQL    | 100m     | 1Gi         |
| Redis    | 50m      | 512Mi       |
| Etcd     | 50m      | 64Mi        |
| **小计** | **200m** | **1.576Gi** |

------





## **总资源需求汇总（微服务 + 数据组件）**



| **类型** | **CPU**  | **Memory**  |
| -------- | -------- | ----------- |
| 微服务   | 190m     | 800Mi       |
| 数据组件 | 200m     | 1.576Gi     |
| **合计** | **390m** | **2.376Gi** |

