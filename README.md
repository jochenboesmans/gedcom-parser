# Gedcom parser
This CLI application allows you to convert GEDCOM, JSON and Protobuf files representing family trees.
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