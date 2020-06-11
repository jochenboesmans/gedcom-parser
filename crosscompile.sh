# clean
rm gedcom-parser.{darwin,linux,windows}.amd64.tar.gz

# build
env GOOS=darwin GOARCH=amd64 go build -o gedcom-parser.darwin.amd64
env GOOS=linux GOARCH=amd64 go build -o gedcom-parser.linux.amd64
env GOOS=windows GOARCH=amd64 go build -o gedcom-parser.windows.amd64

# tar
tar -czf gedcom-parser.darwin.amd64.tar.gz gedcom-parser.darwin.amd64
tar -czf gedcom-parser.linux.amd64.tar.gz gedcom-parser.linux.amd64
tar -czf gedcom-parser.windows.amd64.tar.gz gedcom-parser.windows.amd64

rm gedcom-parser.{darwin,linux,windows}.amd64
