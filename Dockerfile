FROM golang:1.12.3 as builder

RUN mkdir /covid19-remote-work-api
WORKDIR /covid19-remote-work-api

ADD app/go.mod .
ADD app/go.sum .

RUN go mod download

ADD app/ .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o /go/bin/covid19-remote-work-api .

FROM alpine

RUN apk add --no-cache ca-certificates
COPY --from=builder /go/bin/covid19-remote-work-api /app/
COPY startup.sh /app/

WORKDIR /app
CMD ["sh","startup.sh"]⏎