# kobe

kobe 是 Ansible 的轻量级封装，提供了 grpc api 和 golang cli。

## 主要功能

- playbook
- adhoc

## 部署
```shell
# 清理容器
docker ps -a | grep kobe | awk '{print $1}' | xargs docker rm -f
# 构建镜像
make docker
# 部署 
make docker-deploy
```

 
