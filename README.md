# Gedcom parser
This application can be used as a CLI application for local gedcom parsing or as a gRPC service for gedcom parsing on AWS S3.
## Usage
Please make sure to use the file extensions `.ged`, `.json` and `.protobuf` for respectively gedcom, json and protobuf files and to include them in the filepaths.
### Parsing local files
* `gedcom-parser parse path/to/input/file path/to/output/file`
### gRPC service
* `gedcom-parser serve` to launch server
* call `Parse(PathsToFiles)` from any gRPC client to trigger a parse (refer to `grpc/parse.proto` for the exact signature)
   
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
