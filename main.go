// Copyright 2022 Longitude103. All rights reserved.
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
package main

import "github.com/Longitude103/wwum2020/cmd"

// main function is the entry for the application that sets up the CLI and sets the flags needed for the application. This
// function also has an error checking to deal with flags not set correctly.
func main() {
	//	const help = `WWUM 2020 CLI for various tasks. At this point there are two main functions implemented.
	//1. runModel -> Run Full WWUMM 2020 Model
	//2. runSteadyState -> Run the Steady State Model Version
	//3. mfFiles -> Write ModFlow files from a results DB
	//4. qcResults -> Runs QC analysis on the results DB chosen, can output many things
	//-------------------------------------------------------------------------------------------------
	//Use this command: runModel
	//    Required Flags: --Desc: A description of the model being run, use " " around description
	//					--CSDir: The directory path to the CropSim Results text files, use " " if path has spaces
	//
	//    Optional Flags: --StartYr <year>: Will start the model at a specific year (default = 1997)
	//                    --EndYr <year>: Will end the model run at a specific year (default = 2020)
	//                    --noExcessFlow: sets the model to exclude excess flows (default = false)
	//                    --post97: sets the model to post 1997 mode where it holds 1997 acres and crop types constant, no excess flows
	//                    --oldGrid: sets the model use the "old 40 acre grid" and not the new unstructured grid
	//                    --mf6Grid40: sets the model to use the old 40 acre grid but with MF6 node numbers, not Row, Column
	//
	//Use this command: runSteadyState
	//    Required Flags: --Desc: A description of the model being run
	//					--CSDir: The directory path to the CropSim Results text files
	//
	//    Optionsal Flags: --StartYr <year>: Will start the model at a specific year (default = 1893)
	//                     --EndYr <year>: Will end the model run at a specific year (default = 1952)
	//                     --AvgStartYr <year>: Sets the start averaging year
	//                     --AvgEndYr <year>: Sets the end averaging year
	//                     --oldGrid: sets the model use the "old 40 acre grid" and not the new unstructured grid
	//                     --mf6Grid40: sets the model to use the old 40 acre grid but with MF6 node numbers, not Row, Column
	//
	//Use this command: mfFiles
	//	Required Flags: None, but will prompt for selection responses based on data
	//
	//Use this command: qcResults
	//	Required Flags: None, but will prompt for selection responses based on data
	//-------------------------------------------------------------------------------------------------
	//For help with those functions type: runModel -h or mfFiles -h`
	//
	//	var myEnv map[string]string
	//	myEnv, err := godotenv.Read()
	//	if err != nil {
	//		fmt.Println("Cannot load Env Variables:", err)
	//		os.Exit(1)
	//	}
	//	runModelCmd := flag.NewFlagSet("runModel", flag.ExitOnError)
	//	rModelDebug := runModelCmd.Bool("debug", false, "sets debugger for more log information")
	//	rModelStartY := runModelCmd.Int("StartYr", 1997, "Sets the start year of Command, default = 1997")
	//	rModelEndY := runModelCmd.Int("EndYr", 2020, "Sets the end year of Command, default = 2020")
	//	rModelCSDir := runModelCmd.String("CSDir", "", "REQUIRED! - CropSim Directory path")
	//	rModelEF := runModelCmd.Bool("noExcessFlow", false, "Sets to use Excess Flow or Not, default = false")
	//	rModelDesc := runModelCmd.String("Desc", "", "REQUIRED! - Model Description")
	//	rModelP97 := runModelCmd.Bool("post97", false, "If flag set, a post 97 run will be made")
	//	rModelGrid := runModelCmd.Bool("oldGrid", false, "If flag set, the model will use the 40 acre grid, not USG as default")
	//	rModelMF6Grid40 := runModelCmd.Bool("mf6Grid40", false, "If flag set, the model will use the 40 acre grid but in MF6 Node Numbers")
	//	runSSCmd := flag.NewFlagSet("runSteadyState", flag.ExitOnError)
	//	runSSAvgStart := runSSCmd.Int("AvgStartYr", 1953, "Sets the start year of Averaging, default = 1953")
	//	runSSAvgEnd := runSSCmd.Int("AvgEndYr", 2020, "Sets the end year of Averaging, default = 2020")
	//	runSSModelGrid := runSSCmd.Bool("oldGrid", false, "If flag set, the model will use the 40 acre grid, not USG as default")
	//	runSSModelMF6Grid40 := runSSCmd.Bool("mf6Grid40", false, "If flag set, the model will use the 40 acre grid but in MF6 Node Numbers")
	//	runSSModelCSDir := runSSCmd.String("CSDir", "", "REQUIRED! - CropSim Directory path")
	//	runSSModelDesc := runSSCmd.String("Desc", "", "REQUIRED! - Model Description")
	//	runSSStartY := runSSCmd.Int("StartYr", 1893, "Sets the start year of Command, default = 1893")
	//	runSSEndY := runSSCmd.Int("EndYr", 1952, "Sets the End year of Command, default = 1952")
	//
	//	if len(os.Args) < 2 {
	//		fmt.Println(help)
	//		os.Exit(0)
	//	}
	//
	//	switch os.Args[1] {
	//	case "runModel":
	//		err := runModelCmd.Parse(os.Args[2:])
	//		if err != nil {
	//			pterm.Error.Printf("Error in parsing arguments: %s", err)
	//			fmt.Println(help)
	//			os.Exit(1)
	//		}
	//
	//		if *rModelDesc == "" {
	//			pterm.Error.Println("Must include a model description before executing model run")
	//			fmt.Println(help)
	//			os.Exit(0)
	//		}
	//
	//		if *rModelCSDir == "" {
	//			pterm.Error.Println("Must include path to CropSim Files")
	//			fmt.Println(help)
	//			os.Exit(0)
	//		}
	//
	//		if *rModelGrid && *rModelMF6Grid40 {
	//			pterm.Error.Println("Error in Flags: Cannot have both --oldGrid and --mf6Grid40 flags present at the same time")
	//			os.Exit(1)
	//		}
	//
	//		if err := cmd.RunModel(*rModelDebug, rModelCSDir, *rModelDesc, *rModelStartY, *rModelEndY, *rModelEF, *rModelP97, *rModelGrid, *rModelMF6Grid40, myEnv); err != nil {
	//			pterm.Error.Printf("Error in Application: %s\n", err)
	//			os.Exit(1)
	//		}
	//	case "runSteadyState":
	//		err := runSSCmd.Parse(os.Args[2:])
	//		if err != nil {
	//			pterm.Error.Printf("Error Parsing Args for Steady State: %s", err)
	//			os.Exit(1)
	//		}
	//
	//		if *runSSModelDesc == "" {
	//			pterm.Error.Println("Must include a model description before executing model run")
	//			fmt.Println(help)
	//			os.Exit(0)
	//		}
	//
	//		if *runSSModelCSDir == "" {
	//			pterm.Error.Println("Must include path to CropSim Files")
	//			fmt.Println(help)
	//			os.Exit(0)
	//		}
	//
	//		if *runSSModelGrid && *runSSModelMF6Grid40 {
	//			pterm.Error.Println("Error in Flags: Cannot have both --oldGrid and --mf6Grid40 flags present at the same time")
	//			os.Exit(1)
	//		}
	//
	//		if err := cmd.RunSteadyState(*runSSModelDesc, *runSSModelCSDir, *runSSStartY, *runSSEndY, *runSSAvgStart, *runSSAvgEnd, *runSSModelGrid, *runSSModelMF6Grid40, myEnv); err != nil {
	//			pterm.Error.Printf("Error in Steady State Run: %s", err)
	//			os.Exit(1)
	//		}
	//	case "mfFiles":
	//		if err := cmd.MakeModflowFiles(); err != nil {
	//			fmt.Printf("Error in Application: %s\n", err)
	//			os.Exit(1)
	//		}
	//	case "qcResults":
	//		if err := cmd.QcResults(myEnv); err != nil {
	//			fmt.Printf("Error in Application: %s", err)
	//			os.Exit(1)
	//		}
	//	default:
	//		fmt.Println(help)
	//	}

	cmd.Execute()
}
