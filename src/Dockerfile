FROM golang:1.14-alpine

ENV TZ=Asia/Jerusalem
RUN apk add git tzdata

WORKDIR /go/src/app
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

ENV SLACK_SIGNING_SECRET "xxx"
ENV LOG_LEVEL "debug"
ENV PORT "1012"

EXPOSE $PORT
CMD ["app"]
