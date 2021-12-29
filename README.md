# static-server

静态资源下载 动态代理

主要目的解决

docker run -it -p 8080:8080 -v dist:/dist clarechu/http-server:v0.0.1
                  

## 本地安装http-server

```bash
# 下载 http-server
$ https://github.91chifun.workers.dev/https://github.com//clarechu/static-server/releases/download/v0.0.3/http-server-v0.3-win-x86_64.tar.gz

# 解压http-server
$ tar -xvf  http-server-v0.3-win-x86_64.tar.gz

# 配置环境变量引进http-server
```


### 变量配置

```bash
http-server 加载静态资源.

Usage:
  http-server [flags]

Flags:
  -f, --file string    static file path (default "./dist")
  -h, --help           help for http-server
  -i, --index string   static file path index.html (default "index.html")
  -P, --path string    url root path (default "/")
  -p, --port int32     static file server ports (default 8080)


```