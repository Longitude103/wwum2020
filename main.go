package main

import (
	"flag"
	"fmt"
	"github.com/joho/godotenv"
	"os"

	"github.com/Longitude103/wwum2020/actions"
)

func main() {
	const help = `WWUM 2020 CLI for various tasks. At this point there are two main functions implemented.
1. dist -> distribution of CropSim data by cell
2. rch -> create rch files
-------------------------------------------------------------------------------------------------
Use one of these two commands: dist or rch
-------------------------------------------------------------------------------------------------
For help with those functions type: dist -h or rch -h`

	var myEnv map[string]string
	myEnv, err := godotenv.Read()
	if err != nil {
		fmt.Println("Cannot load Env Variables:", err)
		os.Exit(1)
	}
	distCmd := flag.NewFlagSet("dist", flag.ExitOnError)
	distDebug := distCmd.Bool("debug", false, "sets debugger to true to not preform actual write")
	distStartY := distCmd.Int("StartYr", 1997, "Sets the start year of Command, default = 1997")
	distEndY := distCmd.Int("EndYr", 2020, "Sets the end year of Command, default = 2020")
	distCSDir := distCmd.String("CSDir", "", "CropSim Directory path")

	rchCmd := flag.NewFlagSet("rch", flag.ExitOnError)
	rchDebug := rchCmd.Bool("debug", false, "sets debugger for more log information")
	rchStartY := rchCmd.Int("StartYr", 2014, "Sets the start year of Command, default = 1997")
	rchEndY := rchCmd.Int("EndYr", 2014, "Sets the end year of Command, default = 2020")
	rchCSDir := rchCmd.String("CSDir", "", "CropSim Directory path")
	rchEF := rchCmd.Bool("excessFlow", false, "Sets to use Excess Flow or Not, default = false")

	if len(os.Args) < 2 {
		fmt.Println(help)
		os.Exit(0)
	}

	switch os.Args[1] {
	case "dist":
		err := distCmd.Parse(os.Args[2:])
		if err != nil {
			fmt.Println("Error", err)
			fmt.Println(help)
			os.Exit(1)
		}
		fmt.Println("Distribution of CropSim Data")
		_ = distDebug
		_ = distStartY
		_ = distEndY
		_ = distCSDir
		//distribution.Distribution(distDebug, distStartY, distEndY, *distCSDir)

	case "rch":
		err := rchCmd.Parse(os.Args[2:])
		if err != nil {
			fmt.Println("Error", err)
			fmt.Println(help)
			os.Exit(1)
		}
		fmt.Println("Run Recharge File Creation")

		actions.RechargeFiles(*rchDebug, rchCSDir, *rchStartY, *rchEndY, *rchEF, myEnv)
	default:
		fmt.Println(help)
	}

}
