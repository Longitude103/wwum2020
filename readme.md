# WWUMM 2020 CLI
[![Go Reference](https://pkg.go.dev/badge/github.com/Longitude103/wwum2020.svg)](https://pkg.go.dev/github.com/Longitude103/wwum2020)
[![GPLv3 license](https://img.shields.io/badge/License-GPLv3-blue.svg)](http://perso.crans.org/besson/LICENSE.html)
[![Release](https://img.shields.io/github/v/release/Longitude103/wwum2020?display_name=tag)](https://github.com/Longitude103/wwum2020/releases)

This CLI is the major starting point for the 2020 model operations. It's written in [![Go](https://img.shields.io/badge/--00ADD8?logo=go&logoColor=ffffff)](https://golang.org/)
and the binaries are compiled with [![GitHub go.mod Go version of a Go module](https://img.shields.io/github/go-mod/go-version/Longitude103/wwum2020)](https://github.com/Longitude103/wwum2020) and are available in the latest release.

[![Windows](https://svgshare.com/i/ZhY.svg)](https://svgshare.com/i/ZhY.svg) [![macOS](https://svgshare.com/i/ZjP.svg)](https://svgshare.com/i/ZjP.svg) [![Linux](https://svgshare.com/i/Zhy.svg)](https://svgshare.com/i/Zhy.svg)

## Database
The CLI uses the WWUM postgres database hosted on [Amazon Web Services RDS](https://aws.amazon.com/rds/?nc2=h_ql_prod_db_rds) using
a Postgresql database. The postgresql database uses the [PostGIS](https://postgis.net/) extension to calculate the spatial
components of the data.

## CLI functions
There are several CLI functions called with the CLI and this is a list of them. You can 
also get a list by running the program without anything. The "flag" package is used for the CLI which enables the use
of flags in the CLI. This is a core go package.

- runModel -> is the operation to start the main model run function, and it includes the following flags:
  - Required Flags: 
    - --Desc: A description of the model being run
    - --CSDir: The directory path to the CropSim Results text files
  - Optional Flags
    - --debug: run in debug mode with more log output and limited write operations
    - --StartYr: start year of the distribution, defaults to: 1997
    - --EndYr: end year of the distribution, defaults to: 2020


- mfFiles -> is the CLI function to use the results from the `runModel` function and results database to create
  new MODFLOW files for MODFLOW 6 in the form of RCH6 and WEL6 files. It creates an "OutputFiles" directory in the
  same location as the binary for the result files.
    
### Common CLI Examples
- Production runModel:
  - Windows -> In Powershell: `wwum2020-amd64.exe runModel --Desc "Test Run" --CSDir "<path>\WWUMM2020\CropSim\Run005_WWUM2020\Output"`
  - MacOS -> In a Terminal: `./wwum2020-amd64-darwin runModel --Desc "Test Run" --CSDir "<path>/WWUMM2020/CropSim/Run005_WWUM2020/Output"`
  - Linux -> In a Terminal: `./wwum2020-amd64-linux runModel --Desc "Test Run" --CSDir "<path>/WWUMM2020/CropSim/Run005_WWUM2020/Output"`
   
- Production mfFiles
  - Windows -> In Powershell: `wwum2020-amd64.exe mfFiles` -> Follow Prompts
  - MacOS -> In a Terminal: `./wwum2020-amd64-darwin mfFiles` -> Follow Prompts
  - Linux -> In a Terminal: `./wwum2020-amd64-linux mfFIles` -> Follow Prompts


- Development runs:
  - - `go run main.go runModel --CSDir "<path>"`
  - "--CSDir" is the path to the monthly output CropSim .txt files, you must qualify it in "<path>" if the path contains spaces
    - My example is `--CSDir "<path>/WWUMM2020/CropSim/Run005_WWUM2020/Output"`
  - `go run main.go runModel --CSDir "<path>" --debug`

### Output
To store the large volume of output, we're going to use SQLite3 as the storage container. 
This gives the ability to still be SQL enabled, but also allow compressing and store the data in Binary Format.
The only text files that are output are the rch and wel files. The package used is the [go-sqlite3](https://pkg.go.dev/github.com/mattn/go-sqlite3) package in the repo.
The file is named results.db that will be in the same dir as the executable.

### ModFlow Output Files
The second command that is used by the app is the `mfFiles` command that enables you to write out ModFlow 6 files needed for the 
further model runs.

- Run the model using `go run main.go mfFiles`