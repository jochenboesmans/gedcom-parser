# clean previous zip
rm gedcom-parser.{darwin,linux,windows}.7z

# build
env GOOS=darwin GOARCH=amd64 go build -o gedcom-parser.darwin.amd64
env GOOS=linux GOARCH=amd64 go build -o gedcom-parser.linux.amd64
env GOOS=windows GOARCH=amd64 go build -o gedcom-parser.windows.amd64

# zip
7z a gedcom-parser.darwin.amd64.7z gedcom-parser.darwin.amd64
7z a gedcom-parser.linux.amd64.7z gedcom-parser.linux.amd64
7z a gedcom-parser.windows.amd64.7z gedcom-parser.windows.amd64

# clean build
rm gedcom-parser.{darwin,linux,windows}.amd64
