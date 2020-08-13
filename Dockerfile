FROM golang:1.14-alpine

ENV TZ=Asia/Jerusalem
RUN apk add git tzdata

WORKDIR /go/src/app
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

ENV SLACK_SIGNING_SECRET ""
ENV SLACK_WEBHOOK_URL ""
ENV DEVOPS_ONDUTY_NAMESPACE ""
ENV DEVOPS_ONDUTY_CONFIGMAP ""

EXPOSE 8080
CMD ["app"]
