package actions

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/Longitude103/Flogo"
	"github.com/Longitude103/wwum2020/Utils"
	"github.com/Longitude103/wwum2020/database"
	"github.com/jmoiron/sqlx"
)

func MakeModflowFiles() error {
	_, db, err := DbQuestion()
	if err != nil {
		return err
	}

	path, _, err := Utils.MakeOutputDir()
	if err != nil {
		return err
	}

	steadyState, err := database.GetSteadyState(db)
	if err != nil {
		return err
	}

	mDesc, err := database.GetDescription(db)
	if err != nil {
		return err
	}

	a, err := questions(db, steadyState)
	if err != nil {
		return err
	}

	var aggWel, aggRch []database.MfResults

	if steadyState {
		aggRch, err = database.GetAggResults(db, false, a.RchFK)
		if err != nil {
			return err
		}

		// process here before you make the file
		if err := MakeFiles(aggRch, false, true, a.RowCol, "AggregateRCH", path, mDesc); err != nil {
			return err
		}

		return nil
	}

	if !steadyState {
		aggWel, err = database.GetAggResults(db, true, a.WellFK)
		if err != nil {
			return err
		}

		var singleWELResults = make(map[string][]database.MfResults)
		var singleRCHResults = make(map[string][]database.MfResults)

		for _, w := range a.WellFK {
			singleWELResults[w], err = database.SingleResult(db, true, w)
			if err != nil {
				return err
			}
		}

		for _, r := range a.RchFK {
			singleRCHResults[r], err = database.SingleResult(db, false, r)
			if err != nil {
				return err
			}
		}

		if err := MakeFiles(aggWel, true, false, a.RowCol, "AggregateWEL", path, mDesc); err != nil {
			return err
		}

		for k := range singleWELResults {
			fn := fmt.Sprintf("%sWEL", k)
			if err := MakeFiles(singleWELResults[k], true, false, a.RowCol, fn, path, mDesc); err != nil {
				return err
			}
		}

		for k := range singleRCHResults {
			fn := fmt.Sprintf("%sRCH", k)
			if err := MakeFiles(singleRCHResults[k], false, true, a.RowCol, fn, path, mDesc); err != nil {
				return err
			}
		}

		aggRch, err = database.GetAggResults(db, false, a.RchFK)
		if err != nil {
			return err
		}

		if err := MakeFiles(aggRch, false, true, a.RowCol, "AggregateRCH", path, mDesc); err != nil {
			return err
		}
	}

	return nil
}

func suggestFiles(toComplete string) []string {
	files, _ := filepath.Glob("./OutputFiles/*/*.sqlite")
	return files
}

func DbQuestion() (string, *sqlx.DB, error) {
	var q = []*survey.Question{
		{
			Name: "file",
			Prompt: &survey.Input{
				Message: "Which results DB should be used?",
				Suggest: suggestFiles,
				Help:    "The SQLite file you want to use to build ModFlow files",
			},
			Validate: survey.Required,
		},
	}

	answers := struct {
		File string
	}{}

	// ask the question
	if err := survey.Ask(q, &answers); err != nil {
		return "", nil, err
	}

	// get the file_keys from sqlite and ask if any of these should be individual files
	sqliteDB, err := database.ConnectSqlite(answers.File)
	if err != nil {
		return "", nil, err
	}

	return answers.File, sqliteDB, nil
}

type Answers struct {
	WellFK []string
	RchFK  []string
	RowCol bool
}

func questions(sqliteDB *sqlx.DB, steadyState bool) (Answers, error) {
	var multiQs []*survey.Question
	if steadyState {
		multiQs = []*survey.Question{
			{
				Name: "RowCol",
				Prompt: &survey.Confirm{
					Message: "Do you want to Output Row-Column (node is default)?",
				},
			},
		}
	} else {
		wellFk, err := database.GetFileKeys(sqliteDB, true)
		if err != nil {
			return Answers{}, err
		}

		rchFK, err := database.GetFileKeys(sqliteDB, false)
		if err != nil {
			return Answers{}, err
		}

		// the questions to ask
		multiQs = []*survey.Question{
			{
				Name: "WellFK",
				Prompt: &survey.MultiSelect{
					Message: "Choose WEL files to exclude :",
					Options: wellFk,
				},
			},
			{
				Name: "RchFK",
				Prompt: &survey.MultiSelect{
					Message: "Choose which RCH files to exclude :",
					Options: rchFK,
				},
			},
			{
				Name: "RowCol",
				Prompt: &survey.Confirm{
					Message: "Do you want to Output Row-Column (node is default)?",
				},
			},
		}
	}

	var a Answers
	// ask the questions
	if err := survey.Ask(multiQs, &a); err != nil {
		return Answers{}, err
	}

	return a, nil
}

func MakeFiles(r []database.MfResults, wel bool, rch bool, Rc bool, fileName string, outputPath string, mDesc string) error {
	rInterface := make([]interface {
		Date() time.Time
		Node() int
		Value() float64
		RowCol() (int, int)
	}, len(r))
	for i, v := range r {
		if wel {
			// acre-feet / month -> (ft^3 / day) * -1
			v.Rslt = (v.Rslt * 43560) / float64(Utils.TimeExt{T: v.ResultDate}.DaysInMonth()) * -1
		}

		if rch {
			// acre-feet / month -> ft / day
			v.Rslt = (v.Rslt / v.CellSize.Float64) / float64(Utils.TimeExt{T: v.ResultDate}.DaysInMonth())
		}
		rInterface[i] = v
	}

	if err := Flogo.Input(wel, rch, Rc, fileName, rInterface, outputPath, mDesc); err != nil {
		return err
	}

	return nil
}

func processSSAggRCH(results []database.MfResults, fileName, outputPath, mDesc string) error {
	Rc := true
	if results[0].IsNodeResult() {
		Rc = false
	}

	var processResults []database.MfResults
	// needs to do a few things
	for _, r := range []int{1893, 1894} {
		// find 1893 and 1894 and process those into annual results
		// TODO: finish the first to years and make them annual values
		// TODO: convert into ft / day
		_ = r
	}

	// TODO: convert monthly values to ft / day

	rInterface := make([]interface {
		Date() time.Time
		Node() int
		Value() float64
		RowCol() (int, int)
	}, len(processResults))

	if err := Flogo.Input(false, true, Rc, fileName, rInterface, outputPath, mDesc); err != nil {
		return err
	}

	return nil
}

func filterResultsByYear(results []database.MfResults, yr int) []database.MfResults {
	// TODO: filter the results to each year to add them together by node

	return results
}
