LABEL org.opencontainers.image.source=https://github.com/jochenboesmans/gedcom-parser

FROM alpine
COPY gedcom-parser.linux.amd64 /bin/gedcom-parser
EXPOSE 9000
ENTRYPOINT ["gedcom-parser"]
