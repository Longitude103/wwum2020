// Copyright 2020 Longitude103. All rights reserved.
// Use of this source code is governed by GPLv3
// license that can be found in the LICENSE file.

// Package wwum2020 implements the Western Water Use Management model, 2020 update.
//
// The model calculates the crop water use and subsequent
// soil water balance that is needed to operate the MODFLOW
// modeling used in administration decisions by the NRDs of
// Nebraska. The model is a CLI based model that uses two basic
// commands to operate it. The first is runModel that
// tells the model to operate and calculate all the required
// information and saves it in an output database locally.
// Once the output database and log files are created, the
// other command is run mfFiles that produces the MODFLOW
// WEL and RCH files needed to be imported into that model
// for operations. There are a series of questions that
// are asked as to the model options for the output files
// that allow some flexibility for the creation of the
// files.
//
package main

import (
	"flag"
	"fmt"
	"github.com/joho/godotenv"
	"os"

	"github.com/Longitude103/wwum2020/actions"
)

// main function is the entry for the application that sets up the CLI and sets the flags needed for the application. This
// function also has an error checking to deal with flags not set correctly.
func main() {
	const help = `WWUM 2020 CLI for various tasks. At this point there are two main functions implemented.
2. runModel -> Run Full WWUMM 2020 Model
-------------------------------------------------------------------------------------------------
Use this command: runModel
-------------------------------------------------------------------------------------------------
For help with those functions type: dist -h or rch -h`

	var myEnv map[string]string
	myEnv, err := godotenv.Read()
	if err != nil {
		fmt.Println("Cannot load Env Variables:", err)
		os.Exit(1)
	}
	runModelCmd := flag.NewFlagSet("runModel", flag.ExitOnError)
	rModelDebug := runModelCmd.Bool("debug", false, "sets debugger for more log information")
	rModelStartY := runModelCmd.Int("StartYr", 2014, "Sets the start year of Command, default = 1997")
	rModelEndY := runModelCmd.Int("EndYr", 2014, "Sets the end year of Command, default = 2020")
	rModelCSDir := runModelCmd.String("CSDir", "", "CropSim Directory path")
	rModelEF := runModelCmd.Bool("excessFlow", false, "Sets to use Excess Flow or Not, default = false")

	if len(os.Args) < 2 {
		fmt.Println(help)
		os.Exit(0)
	}

	switch os.Args[1] {
	case "runModel":
		err := runModelCmd.Parse(os.Args[2:])
		if err != nil {
			fmt.Println("Error", err)
			fmt.Println(help)
			os.Exit(1)
		}
		fmt.Println("Run Full Model")

		if err := actions.RunModel(*rModelDebug, rModelCSDir, *rModelStartY, *rModelEndY, *rModelEF, myEnv); err != nil {
			fmt.Printf("Error in Application: %s\n", err)
			os.Exit(1)
		}
	case "mfFiles":
		if err := actions.MakeModflowFiles(); err != nil {
			fmt.Printf("Error in Application: %s\n", err)
			os.Exit(1)
		}
	default:
		fmt.Println(help)
	}

}
