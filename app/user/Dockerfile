FROM ubuntu:18.04 AS runner

WORKDIR /app
COPY bin/user /app/bin/
COPY configs/ /app/configs/
EXPOSE 8000
EXPOSE 9000

CMD ["/app/bin/user", "-conf", "/app/configs/config.yaml"]
