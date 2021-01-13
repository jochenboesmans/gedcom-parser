FROM alpine
COPY gedcom-parser.linux.amd64 /bin/gedcom-parser
EXPOSE 9000
ENTRYPOINT ["gedcom-parser"]
