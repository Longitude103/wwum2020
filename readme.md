# WWUMM 2020 CLI
[![Go Reference](https://pkg.go.dev/badge/github.com/Longitude103/wwum2020.svg)](https://pkg.go.dev/github.com/Longitude103/wwum2020)

This CLI is the major starting point for the 2020 model operations. It's written mostly in Go or Python
and includes several operations required for the model function.

## Database
The CLI uses the WWUM postgres database hosted at `long103-wwum.clmtjoquajav.us-east-2.rds.amazonaws.com` using
the `postgres` username and the normal connection to the various schemas.

## CLI functions
There are several CLI functions called with the CLI and this is a list of them. You can 
also get a list by running the program without anything. The "flag" package is used for the CLI which enables the use
of flags in the CLI. This is a core go package.

- runModel -> is the operation to start the main model run function, and it includes the following flags:
  - --debug -> run in debug mode with more log output and limited write operations
  - --StartYr -> start year of the distribution, defaults to 2014
  - --EndYr -> end year of the distribution, defaults to 2015
  - --CSDir -> location of the CropSim results files directory  
    
### Common CLI Examples
- Production runs:
  - `go run main.go runModel --CSDir "/run/media/heath/part1/WWUMM2020/CropSim/Run005_WWUM2020/Output"`
- Development runs:
  - `go run main.go runModel --CSDir "/run/media/heath/part1/WWUMM2020/CropSim/Run005_WWUM2020/Output" --debug`

### Output
To store the large volume of output, we're going to use SQLite3 as the storage container. 
This gives the ability to still be SQL enabled, but also allow compressing and store the data in Binary Format.
The only text files that are output are the rch and wel files. The package used is the [go-sqlite3](https://pkg.go.dev/github.com/mattn/go-sqlite3) package in the repo.
The file is named results.db that will be in the same dir as the executable.