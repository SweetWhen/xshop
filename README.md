# xshop
 一个练手的账号体系。
* 基于kratos，etcd服务注册与发现，openapi, otel链路追踪，入参validata, gorm；
* jwt/casbin账号认证与权限管理；
* envoy网关, idgen, 定时任务；
* es搜索；
* DDD
* IM、shopping...

整体架构如下
![jiagou](./doc/img/xshop架构2.png)


## quick start

> 目前仅支持用户增删查改、登录、搜索用户名（高亮）

### Docker
```bash
# build
docker build -t <your-docker-image-name> .

# run
docker run --rm -p 8000:8000 -p 9000:9000 -v </path/to/your/configs>:/data/conf <your-docker-image-name>
```
