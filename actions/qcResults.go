package actions

import (
	"github.com/AlecAivazis/survey/v2"
	"github.com/Longitude103/wwum2020/database"
	"github.com/Longitude103/wwum2020/qc"
)

func QcResults(myEnv map[string]string) error {
	fileName, db, err := DbQuestion()
	if err != nil {
		return err
	}

	path, err := MakeOutputDir(fileName)
	if err != nil {
		return err
	}

	yr, err := yearQuestion()
	if err != nil {
		return err
	}

	graph, err := graphQuestion()
	if err != nil {
		return err
	}

	v := database.Setup{}

	if err := v.NewSetup(false, false, myEnv, true, ""); err != nil {
		return err
	}

	v.SlDb = db

	q := qc.NewQC(v, path, qc.WithGraph(graph), qc.WithYear(yr))

	if err := qc.QcRMain(q); err != nil {
		return err
	}

	return nil
}

func yearQuestion() (int, error) {
	var q = &survey.Input{
		Message: "Which YEAR should be analyzed?",
		Help:    "The year of the analysis that should be analyzed",
	}

	answer := 0

	// ask the question
	if err := survey.AskOne(q, &answer); err != nil {
		return 0, err
	}

	return answer, nil
}

func graphQuestion() (bool, error) {
	var q = &survey.Confirm{
		Message: "Want to graph the results?",
		Help:    "This will produce a .png graph of the results of the analysis",
	}

	g := false

	// ask the question
	if err := survey.AskOne(q, &g); err != nil {
		return false, err
	}

	return g, nil
}
