FROM golang:1.19.2 AS builder
WORKDIR /build
ADD go.mod .
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o confapp cmd/app/main.go

FROM alpine
WORKDIR /build
COPY --from=builder ./build/confapp .
COPY --from=builder ./build/config/ /build/config/
COPY --from=builder ./build/logs/ /build/logs/
CMD ["./confapp"]