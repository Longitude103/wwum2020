package cmd

import (
	"errors"
	"strconv"
	"strings"

	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var projectOutput = &cobra.Command{
	Use:   "project",
	Short: "Projection used for Robust Review Analysis",
	Long:  "Projection command used to create a file set for robust review analysis.",
	Run: func(cmd *cobra.Command, args []string) {
		yearsStr, _ := cmd.Flags().GetString("years")
		years, err := parseAndSetYears(yearsStr)
		if err != nil {
			pterm.Error.Println(err)
			return
		}

		projYears, _ := cmd.Flags().GetInt("ProjectYears")
		if projYears == 0 {
			pterm.Error.Println("Projection Years must be greater than 0")
		}

		fileName, _ := cmd.Flags().GetString("FileNames")

		if err := projectRR(myEnv, years, projYears, fileName); err != nil {
			pterm.Error.Println("Error in Application: ", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(projectOutput)
	projectOutput.Flags().String("Years", "2020", "A list of years to use for the projection, e.g. 2018, 2020, 2018 in that order")
	projectOutput.Flags().String("FileNames", "RR", "The name of the files to use for the projection output")
	projectOutput.Flags().Int("ProjectYears", 50, "The number of years to project forward")

	projectOutput.MarkFlagRequired("FileNames")
	projectOutput.MarkFlagRequired("Years")
}

func parseAndSetYears(yearsStr string) ([]int, error) {
	var years []int
	split := strings.Split(yearsStr, ",")
	for _, yr := range split {
		yearInt, err := strconv.Atoi(strings.TrimSpace(yr))
		if err != nil {
			return nil, errors.New("error parsing years: " + err.Error())
		}

		years = append(years, yearInt)
	}

	if len(years) == 0 {
		return nil, errors.New("must include at least one year")
	}

	return years, nil
}

func projectRR(myEnv map[string]string, years []int, projYears int, fileName string) error {

	return nil
}
