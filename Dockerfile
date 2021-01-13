FROM golang:1.15 AS getter
RUN go get -d github.com/jochenboesmans/gedcom-parser
RUN GOARCH=amd64 GOOS=linux go build github.com/jochenboesmans/gedcom-parser

FROM ubuntu AS runner
COPY --from=getter /go/gedcom-parser /bin/gedcom-parser
EXPOSE 9000
ENTRYPOINT ["gedcom-parser"]
