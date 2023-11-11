# 用户系统管理后台
有账号才有其他业务体系，先搞个账号服务，后续再加入其他业务，如IM、电商等

* 基于kratos，etcd服务注册与发现，openapi, otel链路追踪，入参validata, gorm；
* envoy网关, idgen, 定时任务；
* 基于Open Policy Agent的rbac权限体系；
* es搜索；
* DDD

## quick start

### Docker
```bash
# build
docker build -t <your-docker-image-name> .

# run
docker run --rm -p 8000:8000 -p 9000:9000 -v </path/to/your/configs>:/data/conf <your-docker-image-name>
```

