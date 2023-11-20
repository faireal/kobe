# kobe

kobe 是 Ansible 的轻量级封装，提供了 grpc api 和 golang cli。

## 主要功能

- playbook
- adhoc

## 部署
```shell
# 构建镜像
make docker
# 部署 
docker run -it -d  -p 8080:8080  trusfort/kobe:master
```

 
