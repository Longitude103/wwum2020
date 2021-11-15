package conveyLoss

import (
	"github.com/Longitude103/wwum2020/database"
	"github.com/pterm/pterm"
)

// Conveyance function finds the diversions and calculates the conveyance loss for all cells where there is a canal. This
// outputs to the results table in sqlite. Might update to return delivery by canal.
func Conveyance(v *database.Setup) (err error) {
	spRates := map[string]float64{"interstate": 0.4869, "highline": 0.2617, "lowline": 0.2513}

	spinner, _ := pterm.DefaultSpinner.Start("Getting Canal Cells and Diversions")
	canalCells, err := getCanalCells(v)
	if err != nil {
		spinner.Fail("Get Canals Failed")
		return err
	}

	diversions, err := getDiversions(v)
	if err != nil {
		spinner.Fail("Get Diversions Failed")
		return err
	}
	spinner.Success()

	p, _ := pterm.DefaultProgressbar.WithTotal(len(canalCells)).WithTitle("Process Canal Cells").WithRemoveWhenDone(true).Start()
	// loop over cells
	for _, cell := range canalCells {
		p.Increment()
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

		factor := getFactor(cell)

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

			ft := 113               // np by default
			if cell.CanalId == 54 { // western canal, only sp canal
				ft = 114
			}

			//if cell.Node == 51030 {
			//	fmt.Printf("Data for 51030: dt: %v, file: %d, st_loss: %g, factor: %g\n", div.DivDate.Time, ft, structureLoss, factor)
			//}

			if v.AppDebug {
				// if cell.CanalId == 26 {
				// 	fmt.Printf("Div: %+v\n", div)
				// 	fmt.Printf("StructureLoss: %f, Factor: %f\n", structureLoss, factor)
				// 	fmt.Printf("Cell Data: %+v\n", cell)
				// 	d := database.RchResult{Node: cell.Node, Size: cell.CellArea, Dt: div.DivDate.Time,
				// 		FileType: ft, Result: structureLoss * factor * 1.9835}
				// 	fmt.Printf("Cell Result: %+v\n", d)
				// }
			} else {
				if structureLoss > 0 {
					err := v.RchDb.Add(database.RchResult{Node: cell.Node, Size: cell.CellArea, Dt: div.DivDate.Time,
						FileType: ft, Result: structureLoss * factor * 1.9835})
					if err != nil {
						return err
					}
				}
			}
		}
	}

	pterm.Success.Println("Canal Loss Calculations")
	v.Logger.Info("Canal Loss Completed Successfully")
	return nil
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

func getFactor(cell CanalCell) (factor float64) {
	switch cell.CFlag {
	case 1: // DNR Factor
		if cell.DnrFact.Valid == false {
			factor = defaultFactor(cell)
		} else {
			factor = cell.DnrFact.Float64
		}
	case 4: // SatThick Factor
		if cell.SatFact.Valid == false {
			factor = defaultFactor(cell)
		} else {
			factor = cell.SatFact.Float64
		}
	case 2: // USGS Factor
		if cell.UsgsFact.Valid == false {
			factor = defaultFactor(cell)
		} else {
			factor = cell.UsgsFact.Float64
		}
	default: // default
		factor = defaultFactor(cell)
	}

	return factor
}

func defaultFactor(cell CanalCell) (factor float64) {
	if cell.CanalType == "Lateral" || cell.CanalType == "Spill" {
		factor = cell.StLength / cell.TotalLatLn.Float64
	} else {
		factor = cell.StLength / cell.TotalCanLn
	}

	return factor
}
