FROM ubuntu:18.04 AS runner

WORKDIR /app
COPY bin/mgr /app/bin/
COPY configs/ /app/configs/
EXPOSE 8001
EXPOSE 9001

CMD ["/app/bin/mgr", "-conf", "/app/configs/config.yaml"]