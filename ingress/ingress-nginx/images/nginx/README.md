NGINX base image using [alpine](https://www.alpinelinux.org/)

This custom image contains:

- [nginx-http-auth-digest](https://github.com/atomx/nginx-http-auth-digest)
- [ngx_http_substitutions_filter_module](https://github.com/yaoweibin/ngx_http_substitutions_filter_module)
- [nginx-opentracing](https://github.com/opentracing-contrib/nginx-opentracing)
- [opentracing-cpp](https://github.com/opentracing/opentracing-cpp)
- [zipkin-cpp-opentracing](https://github.com/rnburn/zipkin-cpp-opentracing)
- [dd-opentracing-cpp](https://github.com/DataDog/dd-opentracing-cpp)
- [ModSecurity-nginx](https://github.com/SpiderLabs/ModSecurity-nginx) (only supported in x86_64)
- [brotli](https://github.com/google/brotli)
- [geoip2](https://github.com/leev/ngx_http_geoip2_module)

**How to use this image:**
This image provides a default configuration file with no backend servers.

_Using docker_

```console
docker run -v /some/nginx.con:/etc/nginx/nginx.conf:ro k8s.gcr.io/ingress-nginx/nginx:v20210104-g2254a9186@sha256:edd1d06bc6892b0dfb42de7d782ceb3c50eec843b09024abf3f95ba23f4feed5
```

_Creating a replication controller_

```console
kubectl create -f ./rc.yaml
```
