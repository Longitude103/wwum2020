package main

import (
	"clibasic/color"
	"flag"
	"fmt"
	"os"
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
	distStartY := distCmd.Int("StartYr", 1997, "Sets the start year of Dist Command, default = 1997")
	distEndY := distCmd.Int("EndYr", 2020, "Sets the end year of Dist Command, default = 2020")

	if len(os.Args) < 2 {
		fmt.Println(help)
		os.Exit(0)
	}

	switch os.Args[1] {
	case "dist":
		distCmd.Parse(os.Args[2:])
		fmt.Println(color.Red + "Distribution of CropSim Data" + color.Reset)
		distribution.Distribution(distDebug, distStartY, distEndY)
	default:
		fmt.Println(help)
	}

}
