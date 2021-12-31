FROM alpine


ENV APP_ROOT=/opt/ \
    PATH=${APP_ROOT}:$PATH:/usr/bin \
    TZ='Asia/Shanghai' \
    HOME=/opt/

RUN  mkdir -p ${APP_ROOT} \
     && sed -i 's/dl-cdn.alpinelinux.org/mirrors.tuna.tsinghua.edu.cn/g' /etc/apk/repositories \
     && apk update \
     && apk upgrade \
     && apk --no-cache add ca-certificates iputils\
     && apk add -U tzdata ttf-dejavu busybox-extras curl bash\
     && ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime

COPY http-server /usr/bin

WORKDIR /opt/

ENTRYPOINT ["/usr/bin/http-server"]
