FROM      alpine:latest
MAINTAINER dwj_wz@163.com

ENV GIN_MODE=release

WORKDIR /app
COPY smswebhook /app
COPY config.yml /app
COPY sms.tmpl /app

EXPOSE 8080
ENTRYPOINT ["/app/smswebhook"]