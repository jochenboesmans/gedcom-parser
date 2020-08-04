# Gedcom parser
Application for local gedcom parsing and gRPC service for gedcom parsing on AWS S3.

## Installation
### Using binary
1. Download the appropriate release for your OS and architecture
2. Unzip it
3. Add the resulting binary to your `$PATH`
### Using Go
Run `go get github.com/jochenboesmans/gedcom-parser`
## Usage
Please make sure to use the file extensions `.ged`, `.json` and `.protobuf` for respectively gedcom, json and protobuf files and to include them in the filepaths.
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
The gedcom model used is based on a limited subset of GEDCOM 5.5.1 as seen in the below proto spec:
```proto
message Gedcom {
    repeated Individual Individuals = 1;
    repeated Family Families = 2;

    message Individual {
        string Id = 1;
        repeated Name Names = 2;
        string Gender = 3;
        Date BirthDate = 4;
        Date DeathDate = 5;

        message Name {
            string GivenName = 1;
            string Surname = 2;
            bool Primary = 3;
        }
        message Date {
            uint32 Year = 1;
            uint32 Month = 2;
            uint32 Day = 3;
        }
    }
    message Family {
        string Id = 1;
        string FatherId = 2;
        string MotherId = 3;
        repeated string ChildIds = 4;
    }
}
```

Output GEDCOM files are fully 5.5.1 spec compliant, but there may be loss of data because of unsupported tags in the input GEDCOM.
## Examples
### GEDCOM -> JSON
#### Input
```
0 HEAD
...
0 @I50@ INDI
1 NAME Bacteria/Monera/
0 @I51@ INDI
1 NAME Schizomycetes Bacteria/Monera/
0 @I52@ INDI
1 NAME Archangiaceae Schizomycetes/Monera/
0 @I53@ INDI
1 NAME Pseudomonadales Schizomycetes/Monera/
...
0 @F51@ FAM
1 HUSB @I51@
1 CHIL @I53@
...
0 TRLR
```
#### Output
```json5
{
  "Individuals":[
    {
      "Id":"@I51@",
      "Names":[
        {
          "GivenName":"Schizomycetes Bacteria",
          "Surname":"Monera"
        }
      ]
    },
    {
      "Id":"@I52@",
      "Names":[
        {
          "GivenName":"Archangiaceae Schizomycetes",
          "Surname":"Monera"
        }
      ]
    },
    {
      "Id":"@I53@",
      "Names":[
        {
          "GivenName":"Pseudomonadales Schizomycetes",
          "Surname":"Monera"
        }
      ]
    }
  ],
  "Families":[
    {
      "Id":"@F51@",
      "FatherId":"@I51@",
      "ChildIds":[
        "@I53@"
      ]
    }
  ]
}
```
