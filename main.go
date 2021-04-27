package main

import (
	"clibasic/color"
	"flag"
	"fmt"
	"os"
	"wwum2020/actions"
	"wwum2020/actions/distribution"
)

func main() {
	const help = `WWUM 2020 CLI for various tasks. At this point there are two main functions implemented.
1. dist -> distribution of CropSim data by cell
2. convloss -> conveyance loss distribution
-------------------------------------------------------------------------------------------------
Use one of these two commands: dist or convloss
-------------------------------------------------------------------------------------------------
For help with those functions type: dist -h or convloss -h
`

	distCmd := flag.NewFlagSet("dist", flag.ExitOnError)
	distDebug := distCmd.Bool("debug", false, "sets debugger to true to not preform actual write")
	distStartY := distCmd.Int("StartYr", 1997, "Sets the start year of Command, default = 1997")
	distEndY := distCmd.Int("EndYr", 2020, "Sets the end year of Command, default = 2020")
	distCSDir := distCmd.String("CSDir", "", "CropSim Directory path")

	rchCmd := flag.NewFlagSet("rch", flag.ExitOnError)
	rchDebug := rchCmd.Bool("debug", false, "sets debugger for more log information")
	rchStartY := rchCmd.Int("StartYr", 1997, "Sets the start year of Command, default = 1997")
	rchEndY := rchCmd.Int("EndYr", 2020, "Sets the end year of Command, default = 2020")
	rchCSDir := rchCmd.String("CSDir", "", "CropSim Directory path")

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
		fmt.Println(color.Red + "Distribution of CropSim Data" + color.Reset)
		distribution.Distribution(distDebug, distStartY, distEndY, *distCSDir)
	case "rch":
		err := rchCmd.Parse(os.Args[2:])
		if err != nil {
			fmt.Println("Error", err)
			fmt.Println(help)
			os.Exit(1)
		}
		fmt.Println(color.Red + "Run Recharge File Creation" + color.Reset)
		actions.RechargeFiles(rchDebug, rchStartY, rchEndY, rchCSDir)
	default:
		fmt.Println(help)
	}

}
