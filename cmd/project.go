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
		yearsStr, _ := cmd.Flags().GetString("Years")
		var years repeatYears
		err := years.parseAndSetYears(yearsStr)
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

func projectRR(myEnv map[string]string, years repeatYears, projYears int, fileName string) error {
	// load local database
	_, db, err := DbQuestion()
	if err != nil {
		return err
	}

	// sqlite database
	_ = db

	for i := 0; i < projYears; i++ {
		// which year to get
		year := years.NextYear()

		// query data for this year
		_ = year

	}

	return nil
}

type repeatYears struct {
	yearsList []int
	index     int
}

func (r *repeatYears) Len() int {
	return len(r.yearsList)
}

func (r *repeatYears) NextYear() int {
	if r.index < len(r.yearsList)-1 {
		yr := r.yearsList[r.index]
		r.index++
		return yr
	} else {
		r.index = 0
		return r.yearsList[0]
	}
}

func (r *repeatYears) parseAndSetYears(yearsStr string) error {
	var years []int
	split := strings.Split(yearsStr, ",")
	for _, yr := range split {
		yearInt, err := strconv.Atoi(strings.TrimSpace(yr))
		if err != nil {
			return errors.New("error parsing years: " + err.Error())
		}

		years = append(years, yearInt)
	}

	if len(years) == 0 {
		return errors.New("must include at least one year")
	}

	r.yearsList = years

	return nil
}
