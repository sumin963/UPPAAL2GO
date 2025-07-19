# UPPAAL2GO
UPPAAL2GO consists of two components: UPPAAL2TADA and TADA2GO.


## UPPAAL2TADA
Converts a formally verified UPPAAL model into a TADA model, which serves as an intermediate representation to facilitate translation into Go code.


Usage:
Modify the internal variable read_file_path to point to the .xml file you wish to convert, then execute the program.

```
cd src
go run UPPAAL2TADA.go
```

## TADA2GO
Conversion of TADA Models to Go Language

```
go run TADA2GO.go
```

Usage:
Modify the internal variable input_path to specify the path of the XML file to be converted, then run the program.


Required Libraries
* https://github.com/beevik/etree
* https://github.com/dave/jennifer

