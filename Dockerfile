FROM golang:1.14.7 AS getter
RUN go get -d github.com/jochenboesmans/gedcom-parser
RUN GOARCH=amd64 GOOS=linux go build github.com/jochenboesmans/gedcom-parser

FROM ubuntu AS runner
RUN apt-get update && apt-get install -y ca-certificates && update-ca-certificates
COPY --from=getter /go/gedcom-parser /bin/gedcom-parser
EXPOSE 9000
ENTRYPOINT ["gedcom-parser"]
