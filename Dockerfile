FROM golang:1.19 as backend
WORKDIR /usr/src/app
COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o dist/monitor cmd/monitor/main.go

FROM alpine:latest
WORKDIR /root/
COPY --from=backend /usr/src/app/dist/monitor /root/monitor
CMD ["/root/monitor", "web"]