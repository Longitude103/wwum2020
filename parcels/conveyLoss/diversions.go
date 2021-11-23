package conveyLoss

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Longitude103/wwum2020/database"
)

// efPeriod is a struct to hold the data from excess_flow_periods
type efPeriod struct {
	CanalId   int          `db:"canal_id"`
	StartDate sql.NullTime `db:"st_date"`
	EndDate   sql.NullTime `db:"end_date"`
}

// Diversion is a struct to hold the dailydiversions table data and also the results are this struct which is a monthly
// total using the first day of each month.
type Diversion struct {
	CanalId   int             `db:"canal_id"`
	DivDate   sql.NullTime    `db:"div_dt"`
	DivAmount sql.NullFloat64 `db:"div_amnt_cfs"`
}

// applyEffAcres is a method that multiplies the Diversion by the efficiency and acres passed in and converts the day cfs
// to acre-feet to yield a application in inches per acre from that surface water structure.
func (d *Diversion) applyEffAcres(eff float64, acres float64) {
	d.DivAmount.Float64 = d.DivAmount.Float64 * eff * 1.9835 / acres
}

func (d *Diversion) printString() string {
	return fmt.Sprintf("Canal ID: %d, DivDate: %d-%d, DivAmount: %.3f", d.CanalId, d.DivDate.Time.Month(), d.DivDate.Time.Year(), d.DivAmount.Float64)
}

// getDiversions retrieves the diversions from the pg database and returns a slice of Diversion struct for each canal
// during the year and also takes in a start year, end year and also excessFlow bool that if false will remove the
// excess flow from the daily diversions based on excess flow periods. Result diversions are in day cfs.
func getDiversions(v *database.Setup) (diversions []Diversion, err error) {
	// TODO: Change this to be a map[canalid][]Diversion

	if v.ExcessFlow {
		divQry := fmt.Sprintf(`select canal_id, make_timestamp(cast(extract(YEAR from div_dt) as int), 
cast(extract(MONTH from div_dt) as int), 1, 0, 0, 0) as div_dt, sum(div_amnt_cfs) as div_amnt_cfs 
from (select cdj.canal_id, div_dt, sum(div_amnt_cfs) as div_amnt_cfs from sw.dailydiversions 
inner join sw.canal_diversion_jct cdj on dailydiversions.ndnr_id = cdj.ndnr_id WHERE 
div_dt >= '%d-01-01' AND div_dt <= '%d-12-31' group by cdj.canal_id, div_dt) as daily_query group by canal_id, 
extract(MONTH from div_dt), extract(YEAR from div_dt) order by canal_id, div_dt;`, v.SYear, v.EYear)

		if err := v.PgDb.Select(&diversions, divQry); err != nil {
			v.Logger.Errorf("Error in getting diversion records with Excess Flow: %s", err)
			return nil, err
		}

	} else {
		// get all Diversion data
		divQry := fmt.Sprintf(`select cdj.canal_id, div_dt, sum(div_amnt_cfs) as div_amnt_cfs from sw.dailydiversions
inner join sw.canal_diversion_jct cdj on dailydiversions.ndnr_id = cdj.ndnr_id
WHERE div_dt between '%d-01-01' and '%d-12-31' group by cdj.canal_id, div_dt;`, v.SYear, v.EYear)

		var preDiversions []Diversion

		if err := v.PgDb.Select(&preDiversions, divQry); err != nil {
			v.Logger.Errorf("Error getting all diversions starting Excess Flow Limitation: %s", err)
			return nil, err
		}

		// remove the excess flows
		efQuery := "select canal_id, st_date, end_date from sw.excess_flow_periods;"

		var efPeriods []efPeriod
		if err = v.PgDb.Select(&efPeriods, efQuery); err != nil {
			v.Logger.Errorf("Error in getting Excess Flow Periods: %s", err)
		}

		// list of canals that have excess flows
		var efCanals []int
		for _, v := range efPeriods {
			if find(efCanals, v.CanalId) {
				efCanals = append(efCanals, v.CanalId)
			}
		}

		// loop through the list of canals with excess flows
		var canalMonthlyDiversion []Diversion
		for _, canal := range efCanals {

			// generate list of excess flow periods for this canal
			var efCanalPeriods []efPeriod
			for _, v := range efPeriods {
				if v.CanalId == canal {
					efCanalPeriods = append(efCanalPeriods, v)
				}
			}

			// filter daily canal diversions to this canal
			var allCanalDailyDiversion []Diversion
			for _, div := range preDiversions {
				if div.CanalId == canal {
					allCanalDailyDiversion = append(allCanalDailyDiversion, div)
				}
			}

			// filter out the times when there was excess flow and only return the times the flow wasn't during that
			// excess flow period
			var canalDailyDiversion []Diversion
			for _, div := range allCanalDailyDiversion {
				if findDiversion(div, efCanalPeriods) {
					canalDailyDiversion = append(canalDailyDiversion, div)
				}
			}

			// reduce the daily flows to monthly
			for y := v.SYear; y < v.EYear+1; y++ {
				for m := 1; m < 13; m++ {
					canalMonthlyDiversion = append(diversions, monthlyDiversion(canalDailyDiversion, m, y, canal))
				}
			}
		}

		// remaining diversions
		var strCanal []string
		for _, v := range efCanals {
			strCanal = append(strCanal, strconv.Itoa(v))
		}

		var canalLimit string
		if len(strCanal) > 0 {
			canalLimit = fmt.Sprintf("AND cdj.canal_id not in (%v)", strings.Join(strCanal, ", "))
		}

		// need to make an alternate if len(strCanal) == 0
		remainDiversionQry := fmt.Sprintf(`select canal_id, make_timestamp(cast(extract(YEAR from div_dt) as int), 
cast(extract(MONTH from div_dt) as int), 1, 0, 0, 0) as div_dt, sum(div_amnt_cfs) as div_amnt_cfs 
from (select cdj.canal_id, div_dt, sum(div_amnt_cfs) as div_amnt_cfs from sw.dailydiversions 
inner join sw.canal_diversion_jct cdj on dailydiversions.ndnr_id = cdj.ndnr_id WHERE div_dt >= '%d-01-01' 
AND div_dt <= '%d-12-31' %s group by cdj.canal_id, div_dt) as daily_query 
group by canal_id, extract(MONTH from div_dt), extract(YEAR from div_dt) order by canal_id, div_dt;`, v.SYear, v.EYear, canalLimit)

		var remainingDiversions []Diversion
		if err = v.PgDb.Select(&remainingDiversions, remainDiversionQry); err != nil {
			v.Logger.Errorf("Error in getting remaining diversions: %s", err)
		}
		diversions = append(canalMonthlyDiversion, remainingDiversions...)

	}

	amendedDivs := adjLaramie(diversions)

	return amendedDivs, nil
}

// Find is a filter function for slice of int in a int
func find(slice []int, val int) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}

// FindDiversion is a function to filter if a Diversion is within a period
func findDiversion(div Diversion, period []efPeriod) bool {
	for _, p := range period {
		if div.DivDate.Time.Before(p.StartDate.Time) || div.DivDate.Time.After(p.EndDate.Time) {
			return true
		}
	}

	return false
}

// monthlyDiversion is a function to aggregate a slice of dailyDiversion by months and returns a month total Diversion
func monthlyDiversion(dailyDiversion []Diversion, m int, y int, cID int) (mDiversion Diversion) {
	var totalDiversion float64
	for _, d := range dailyDiversion {
		if d.DivDate.Time.Year() == y && int(d.DivDate.Time.Month()) == m && d.DivAmount.Valid {
			totalDiversion += d.DivAmount.Float64
		}
	}

	mDiversion = Diversion{DivDate: sql.NullTime{Time: time.Date(y, time.Month(m), 1, 0, 0, 0, 0, time.UTC),
		Valid: true}, DivAmount: sql.NullFloat64{Float64: totalDiversion, Valid: true}, CanalId: cID}

	return mDiversion
}

// adjLaramie fills the diversions for the Laramie portion of the canal when there are no records for the stateline flow. Before this time
// there are only records for the river diversion in WY. This creates those records by converting the WY diversion amounts to stateline diversions
// based on the time period percentages from years of overlap.
func adjLaramie(diversions []Diversion) (amendDivs []Diversion) {
	for d := 0; d < len(diversions); d++ {
		if diversions[d].DivDate.Time.Before(time.Date(1981, 1, 1, 0, 0, 0, 0, time.UTC)) && diversions[d].CanalId == 272 {
			// record that needs to be used to create one for 52200 Laramie
			amendDivs = append(amendDivs, Diversion{CanalId: 26, DivDate: diversions[d].DivDate, DivAmount: sql.NullFloat64{Valid: true,
				Float64: laramieCanalFlow(int(diversions[d].DivDate.Time.Month()), diversions[d].DivAmount.Float64)}})
		} else {
			amendDivs = append(amendDivs, diversions[d])
		}
	}

	return amendDivs
}

// laramieCanalFlow is a function to recieve a month and flow amount and returns the adjusted diversion at the stateline. This was created by
// comparing the diversions at the state line with the river diversions by month and creating the percentage that flows into NE. Default is the average
// of the values in case there are some outliers.
func laramieCanalFlow(m int, div float64) float64 {
	remainingFlowPercent := 0.0

	switch m {
	case 5:
		remainingFlowPercent = .3597
	case 6:
		remainingFlowPercent = .3274
	case 7:
		remainingFlowPercent = .4052
	case 8:
		remainingFlowPercent = .4127
	case 9:
		remainingFlowPercent = 0.4645
	default:
		remainingFlowPercent = .3954
	}

	return remainingFlowPercent * div
}
