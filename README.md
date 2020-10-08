# smart-dns-go
smart-dns-go  go version SmartDns 后端
![](./README/s1.png)

Smart DNS 不仅为域名解析系统 也可充当公共DNS服务器


### Deployment
``` 
docker pull redis:6.0.5-alpine
docker run -d  --name my_redis -p 6379:6379 redis:6.0.5-alpine
```
