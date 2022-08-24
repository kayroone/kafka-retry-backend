FROM golang:1.17-buster as BUILDER

WORKDIR /app
COPY go.mod .
COPY go.sum .
COPY config.yaml /go

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
      -ldflags='-w -s -extldflags "-static"' -a \
      -o /go/bin/retry .

FROM scratch

COPY --from=BUILDER /go/bin/retry /go/config.yaml /go/bin/

ENTRYPOINT ["/go/bin/retry"]