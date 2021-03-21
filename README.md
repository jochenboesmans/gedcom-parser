# Gedcom parser
Lightweight, high performance GEDCOM 5.5.1 parser allowing for easy conversion between `.ged` and `.json` files representing lineage-linked family trees.

# Benchmarks
Benchmarks include full validation. CPU used for testing is Macbook Air M1 (ARM).

| | Harry Potter (16KB) | Gods from Wikipedia (1.4MB) | ITIS (95.7MB) |
| --- | --- | --- | --- |
Time | 0.001s | 0.015s | 1.42s |
Memory (heap) | N/A(?) | 2.61MB | 203.90MB

?: too low to profile

Time/space tradeoffs were made in favor of time to allow for p99 parsing times <1s.

## Installation
### Using binary
1. Download the appropriate release for your OS and architecture from GitHub releases
2. Unzip it
3. Add the resulting binary to your `$PATH`
### Using Go
Run `go get github.com/jochenboesmans/gedcom-parser`
## Usage
Please make sure to use the file extensions `.ged` and `.json` for respectively gedcom and json files and to include them in the filepaths.
### Parsing local files
* `gedcom-parser parse path/to/input/file path/to/output/file`
### gRPC service
* set up an S3 bucket and create a .env file with your `AWS_REGION`, `AWS_S3_BUCKET`, `AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY`
* `gedcom-parser serve` to launch server
* call `Parse(PathsToFiles)` from any gRPC client to trigger a parse (refer to `grpc/parse.proto` for the exact signature)

### Using Docker
Run `docker run -e AWS_REGION=... -e AWS_S3_BUCKET=... -e AWS_ACCESS_KEY_ID=... -e AWS_SECRET_ACCESS_KEY=... -p 9000:9000 jochenboesmans/gedcom-parser serve|parse`

Supplying AWS env variables is only necessary for running `serve`.
   
## Gedcom specification
The gedcom model used is based on a limited subset of GEDCOM 5.5.1 and is fully 5.5.1 spec extensible.
See `./gedcom/gedcom.proto` for the full specification.

## Validation
By default, the parser will validate the following:
1. Record id uniqueness
2. Cross-referential id integrity, i.e. references from within records to other records are valid

## Examples
See files in `./examples` and `./test-output`.
