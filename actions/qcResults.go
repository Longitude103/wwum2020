package actions

import (
	"github.com/AlecAivazis/survey/v2"
	"github.com/Longitude103/wwum2020/Utils"
	"github.com/Longitude103/wwum2020/database"
	"github.com/Longitude103/wwum2020/qc"
)

func QcResults(myEnv map[string]string) error {
	_, db, err := DbQuestion()
	if err != nil {
		return err
	}

	path, _, err := Utils.MakeOutputDir()
	if err != nil {
		return err
	}

	var opts []qc.Option

	aGJ, err := WellAnnGJQuestion()
	if err != nil {
		return err
	}

	if aGJ {
		opts = append(opts, qc.WithWellAnnGJson())
	}

	rechBalance, err := RechBalanceQuestion()
	if err != nil {
		return err
	}

	if rechBalance {
		opts = append(opts, qc.WithRechargeBalance())
	}

	graph, err := graphQuestion()
	if err != nil {
		return err
	}
	if graph {
		opts = append(opts, qc.WithGraph())
	}

	gj, err := GJQuestion()
	if err != nil {
		return err
	}

	if gj {
		mOrA, err := mOrYQuestion()
		if err != nil {
			return err
		}
		if mOrA == "Monthly" {
			opts = append(opts, qc.WithGJson(), qc.WithMonthly())
		} else {
			opts = append(opts, qc.WithGJson())
		}

	}

	wGJ, err := WellGJQuestion()
	if err != nil {
		return err
	}

	if wGJ {
		mOrA, err := mOrYQuestion()
		if err != nil {
			return err
		}

		if mOrA == "Monthly" {
			opts = append(opts, qc.WithWellGJson(), qc.WithMonthly())
		} else {
			opts = append(opts, qc.WithWellGJson())
		}
	}

	if gj || wGJ || rechBalance {
		yr, err := yearQuestion()
		if err != nil {
			return err
		}
		opts = append(opts, qc.WithYear(yr))
	}

	v, err := database.NewSetup(myEnv, database.WithNoSQLite())
	if err != nil {
		return err
	}

	v.SlDb = db

	q := qc.NewQC(v, path, opts...)

	if err := qc.StartQcRMain(q); err != nil {
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

func GJQuestion() (bool, error) {
	var q = &survey.Confirm{
		Message: "Want to Output a GeoJson file of one Year of RECHARGE results?",
		Help:    "This will produce a GeoJson file of the recharge results and saved to 'Output' Directory",
	}

	gj := false

	// ask the question
	if err := survey.AskOne(q, &gj); err != nil {
		return false, err
	}

	return gj, nil
}

func WellGJQuestion() (bool, error) {
	var q = &survey.Confirm{
		Message: "Want to Output a GeoJson file of one Year of WELL PUMPING results?",
		Help:    "This will produce a GeoJson file of the Well Pumping results and saved to 'Output' Directory",
	}

	gj := false

	// ask the question
	if err := survey.AskOne(q, &gj); err != nil {
		return false, err
	}

	return gj, nil
}

func RechBalanceQuestion() (bool, error) {
	var q = &survey.Confirm{
		Message: "Want to Output a Recharge Balance for a Given Year?",
		Help:    "This will produce a table in the CLI of the recharge totals for a single year",
	}

	rb := false

	// ask the question
	if err := survey.AskOne(q, &rb); err != nil {
		return false, err
	}

	return rb, nil
}

func WellAnnGJQuestion() (bool, error) {
	var q = &survey.Confirm{
		Message: "Want to Output a GeoJson file of the Annual WELL PUMPING of all results?",
		Help:    "This will produce a GeoJson file of the Well Pumping results and saved to 'Output' Directory",
	}

	gj := false

	// ask the question
	if err := survey.AskOne(q, &gj); err != nil {
		return false, err
	}

	return gj, nil
}

func mOrYQuestion() (string, error) {
	var q = &survey.Select{
		Message: "Want the output to include Monthly or just Annual Data",
		Help:    "Selecting Monthly will add the monthly data to the GeoJson file of the results and saved to 'Output' Directory",
		Options: []string{"Annual", "Monthly"},
	}

	aOrM := ""

	// ask the question
	if err := survey.AskOne(q, &aOrM); err != nil {
		return "", err
	}

	return aOrM, nil
}
