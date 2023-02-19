FROM golang:1.19-alpine AS build
ENV GOPROXY=https://goproxy.cn
COPY . /app/account/
WORKDIR /app/account
RUN wget https://curl.haxx.se/ca/cacert.pem
RUN CGO_ENABLED=0 go build -o /bin/account cmd/account.go

FROM scratch
COPY --from=build /bin/account /app/account
COPY --from=build /app/account/configs/account.yaml /app/account.yaml
COPY --from=build /app/account/cacert.pem /etc/ssl/certs/

CMD ["/app/account", "--config", "/app/account.yaml"]