package conveyLoss

import (
	"github.com/heath140/wwum2020/database"
	"github.com/schollz/progressbar/v3"
)

// Conveyance function finds the diversions and calculates the conveyance loss for all cells where there is a canal. This
// outputs to the results table in sqlite. Might update to return delivery by canal.
func Conveyance(v database.Setup) (err error) {
	spRates := map[string]float64{"interstate": 0.4869, "highline": 0.2617, "lowline": 0.2513}

	clDB, err := database.ConveyLossDB(v.SlDb)
	if err != nil {
		return err
	}

	defer func(clDB *database.CLDB) {
		err := clDB.Close()
		if err != nil {
			return
		}
	}(clDB)

	canalCells := getCanalCells(v.PgDb)
	diversions := getDiversions(v)

	bar := progressbar.Default(int64(len(canalCells)), "Canal Cells")
	// loop over cells
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
				err := clDB.Add(database.CLResult{Node: cell.Node, Dt: div.DivDate.Time, FileType: ft, Result: structureLoss * factor * 1.9835})
				if err != nil {
					return err
				}
			}
		}

		_ = bar.Add(1)
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