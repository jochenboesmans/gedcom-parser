FROM golang:1.14 AS builder

RUN go get github.com/jochenboesmans/gedcom-parser
WORKDIR /build
RUN go build

FROM alpine AS runner

ENV WDIR /go/src/github.com/jochenboesmans/gedcom-parser

WORKDIR ${WDIR}
COPY --from=builder /build ${WDIR}
RUN mkdir -p ./artifacts

CMD ["./gedcom-parser"]
