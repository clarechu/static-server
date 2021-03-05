FROM alpine

RUN wget https://github.91chifun.workers.dev//https://github.com/ClareChu/static-server/releases/download/v0.0.2/http-server-amd-x86_64.tar.gz \
    && tar -xvf http-server-amd-x86_64.tar.gz \
    && mv /amd-x86_64/http-server /usr/bin/

CMD ["http-server"]
