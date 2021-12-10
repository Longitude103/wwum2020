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
	fileName, db, err := DbQuestion()
	if err != nil {
		return err
	}

	path, err := Utils.MakeOutputDir(fileName)
	if err != nil {
		return err
	}

	mDesc, err := database.GetDescription(db)
	if err != nil {
		return err
	}

	welFK, rchFK, err := questions(db)
	if err != nil {
		return err
	}

	aggWel, err := database.GetAggResults(db, true, welFK)
	if err != nil {
		return err
	}

	aggRch, err := database.GetAggResults(db, false, rchFK)
	if err != nil {
		return err
	}

	var singleWELResults = make(map[string][]database.MfResults)
	var singleRCHResults = make(map[string][]database.MfResults)

	for _, w := range welFK {
		singleWELResults[w], err = database.SingleResult(db, true, w)
		if err != nil {
			return err
		}
	}

	for _, r := range rchFK {
		singleRCHResults[r], err = database.SingleResult(db, false, r)
		if err != nil {
			return err
		}
	}

	if err := MakeFiles(aggWel, true, false, "AggregateWEL", path, mDesc); err != nil {
		return err
	}

	for k := range singleWELResults {
		fn := fmt.Sprintf("%sWEL", k)
		if err := MakeFiles(singleWELResults[k], true, false, fn, path, mDesc); err != nil {
			return err
		}
	}

	if err := MakeFiles(aggRch, false, true, "AggregateRCH", path, mDesc); err != nil {
		return err
	}

	for k := range singleRCHResults {
		fn := fmt.Sprintf("%sRCH", k)
		if err := MakeFiles(singleRCHResults[k], false, true, fn, path, mDesc); err != nil {
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

func questions(sqliteDB *sqlx.DB) (a2 []string, a3 []string, err error) {
	wellFk, err := database.GetFileKeys(sqliteDB, true)
	if err != nil {
		return nil, nil, err
	}

	// the questions to ask
	var multiQs = []*survey.Question{
		{
			Name: "wellFK",
			Prompt: &survey.MultiSelect{
				Message: "Choose WEL files to exclude :",
				Options: wellFk,
			},
		},
	}

	var answers2 []string
	// ask the question
	if err := survey.Ask(multiQs, &answers2); err != nil {
		return nil, nil, err
	}

	rchFK, err := database.GetFileKeys(sqliteDB, false)
	if err != nil {
		return nil, nil, err
	}

	multiQs = []*survey.Question{
		{
			Name: "rchFK",
			Prompt: &survey.MultiSelect{
				Message: "Choose which RCH files to exclude :",
				Options: rchFK,
			},
		},
	}

	var answers3 []string
	// ask the question
	if err := survey.Ask(multiQs, &answers3); err != nil {
		return nil, nil, err
	}

	return answers2, answers3, nil
}

func MakeFiles(r []database.MfResults, wel bool, rch bool, fileName string, outputPath string, mDesc string) error {
	rInterface := make([]interface {
		Date() time.Time
		Node() int
		Value() float64
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

	if err := Flogo.Input(wel, rch, fileName, rInterface, outputPath, mDesc); err != nil {
		return err
	}

	return nil
}
