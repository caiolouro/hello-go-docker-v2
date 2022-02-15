# Base stage
FROM golang:1.16 as base

# New stage called "dev"
FROM base as dev

RUN curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh -s -- -b $(go env GOPATH)/bin

WORKDIR /opt/app/api
CMD ["air"]

FROM base as built

WORKDIR /go/app/api
COPY . .

ENV CGO_ENABLED=0

RUN go get -d -v ./...
RUN go build -o /tmp/server ./*.go

FROM busybox

COPY --from=built /tmp/server /usr/bin/server
CMD ["server", "start"]