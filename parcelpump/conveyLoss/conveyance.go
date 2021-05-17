package conveyLoss

import (
	"fmt"
	"github.com/heath140/wwum2020/rchFiles"
	"github.com/jmoiron/sqlx"
	"github.com/schollz/progressbar/v3"
)

// Conveyance function finds the diversions and calculates the conveyance loss for all cells where there is a canal. This
// outputs to the results table in sqlite. Might update to return delivery by canal.
func Conveyance(pgDB *sqlx.DB, slDB *sqlx.DB, sYear int, eYear int, excessFlow bool) (err error) {
	spRates := map[string]float64{"interstate": 0.4869, "highline": 0.2617, "lowline": 0.2513}

	canalCells := getCanalCells(pgDB)
	//fmt.Println("First 10 Canal Cells")
	//tCells := 0
	//for _, v := range canalCells {
	//	if v.CanalId == 52 {
	//		tCells += 1
	//	}
	//}

	//fmt.Println("Total Canal Cells in 52", tCells)

	diversions := getDiversions(pgDB, sYear, eYear, excessFlow)
	//fmt.Println("First 10 Diversions")
	//for _, v := range diversions {
	//	if v.CanalId == 52 {
	//		fmt.Println(v)
	//	}
	//}
	//
	//fmt.Println("Total Canal Diversions", len(diversions))

	bar := progressbar.Default(int64(len(canalCells)), "Canal Cells")
	// loop over cells
	var cellResults []rchFiles.Result
	for _, cell := range canalCells {

		strLossPercent := 0.0
		cellIdDiv := 0

		// determine efficiency and get total loss factor
		if cell.CanalType == "Lateral" || cell.CanalType == "Spill" {
			strLossPercent = (1 - cell.CanalEff.Float64) * 1 / 6
			cellIdDiv = cell.CLinkId
		} else {
			cellIdDiv = cell.CanalId
			if cell.LatCount.Int64 == 0 {
				strLossPercent = 1 - cell.CanalEff.Float64
			} else {
				strLossPercent = (1 - cell.CanalEff.Float64) * 5 / 6
			}
		}

		//if cell.Node == 51030 {
		//	fmt.Printf("LatCount.Valid: %t, count is: %d\n", cell.LatCount.Valid, cell.LatCount.Int64)
		//	fmt.Printf("Math should be 5/6 * 1-canaleff: %g\n", (1-cell.CanalEff.Float64)*5/6)
		//	fmt.Printf("Canal Type: %s\n", cell.CanalType)
		//	fmt.Printf("StrLossPercent: %g, Cell Canal Eff: %g\n", strLossPercent, cell.CanalEff.Float64)
		//}

		// determine factor using DNR/USGS/SatThick/Length (Length is default)
		factor := 0.0
		switch cell.CFlag {
		case 1: // DNR Factor
			factor = cell.DnrFact.Float64
		case 4: // SatThick Factor
			factor = cell.SatFact.Float64
		case 2: // USGS Factor
			factor = cell.UsgsFact.Float64
		default: // default
			if cell.CanalType == "Lateral" || cell.CanalType == "Spill" {
				factor = cell.StLength / cell.TotalLatLn.Float64
			} else {
				factor = cell.StLength / cell.TotalCanLn
			}
		}

		// special cases for cells with Minatare, Mitchell Gering, and Highline and Lowline
		switch cell.CanalId {
		case 29: // Use the same diversions for north and south Minatare cells
			cellIdDiv = 30
		case 13: //  Mitchell Gering use Mitchell and Gering diversions for Gering Canal
			cellIdDiv = 13
		case 21, 25: // Highline and lowline canal use interstate canal diversions split by percentage of acerage
			cellIdDiv = 2
		}

		// if prev_id != cell_id_div save off and get new canal diversions -- not sure we need this...
		// filter diversions to the canal
		canalDiversions := filterCanal(diversions, cellIdDiv)

		structureLoss := 0.0
		for _, div := range canalDiversions {
			if cell.CanalId == 2 || cell.CanalId == 25 || cell.CanalId == 21 {
				switch div.DivDate.Time.Month() {
				case 1, 2, 3, 10, 11, 12:
					if cell.CanalId == 2 {
						structureLoss = strLossPercent * div.DivAmount.Float64
					} else {
						structureLoss = 0.0
					}
				default:
					switch cell.CanalId {
					case 2:
						structureLoss = strLossPercent * div.DivAmount.Float64 * spRates["interstate"]
					case 21:
						structureLoss = strLossPercent * div.DivAmount.Float64 * spRates["lowline"]
					case 25:
						structureLoss = strLossPercent * div.DivAmount.Float64 * spRates["highline"]
					}
				}
			} else {
				structureLoss = strLossPercent * div.DivAmount.Float64
			}

			ft := 114               // np by default
			if cell.CanalId == 54 { // western canal, only sp canal
				ft = 113
			}

			//if cell.Node == 51030 {
			//	fmt.Printf("Data for 51030: dt: %v, file: %d, st_loss: %g, factor: %g\n", div.DivDate.Time, ft, structureLoss, factor)
			//}

			if structureLoss > 0 {
				cellResults = append(cellResults,
					rchFiles.Result{Node: cell.Node, Dt: div.DivDate.Time, FileType: ft, Result: structureLoss * factor * 1.9835})
			}

			if len(cellResults) == 500 {
				// save the data off, then clear slice
				err := insertSql(slDB, cellResults)
				if err != nil {
					fmt.Println("Error in insert of SQL", err)
				}
				cellResults = nil
			}
		}

		_ = bar.Add(1)
	}

	err = insertSql(slDB, cellResults)
	if err != nil {
		fmt.Println("Error in insert of SQL", err)
	}

	_ = bar.Close()
	return err
}

// filterCanal filters the canal diversions to a specific canal and returns a slice of Diversion
func filterCanal(diversions []Diversion, canal int) (canalDiversion []Diversion) {
	for _, v := range diversions {
		if v.CanalId == canal {
			canalDiversion = append(canalDiversion, v)
		}
	}
	return canalDiversion
}

func insertSql(slDB *sqlx.DB, values []rchFiles.Result) (err error) {
	_, err = slDB.NamedExec(`INSERT INTO results (cell_node, dt, file_type, result) 
										VALUES (:cell_node, :dt, :file_type, :result)`, values)
	if err != nil {
		fmt.Println("Error in insert of Cell Loss", err)
	}

	return err
}
