FROM golang:1.14 AS builder

WORKDIR /build
COPY . .
RUN go build

FROM alpine AS runner

WORKDIR /app
COPY --from=builder /build .
RUN mkdir -p ./io

CMD ["./gedcom-parser"]
