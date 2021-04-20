# WWUMM 2020 CLI
This CLI is the major starting point for the 2020 model operations. It's written mostly in Go or Python
and includes several operations required for the model function.

## Database
The CLI uses the WWUM postgres database hosted at `long103-wwum.clmtjoquajav.us-east-2.rds.amazonaws.com` using
the `postgres` username and the normal connection to the various schemas.

## CLI functions
There are serveral CLI functions called with the CLI and this is a list of them. You can 
also get a list by running the program without anything.

- dist -> is the operation to start the distribution of the Cropsim data files by cell and weather station.
 It includes the following flags:
  - --debug -> run in debug mode with more log output and no write operations
  - --StartYr -> start year of the distribution
  - --EndYr -> end year of the distribution 
    
### Common CLI Examples
- Development runs:
  - `go run main.go dist --debug --StartYr 1996 --EndYr 2019`
