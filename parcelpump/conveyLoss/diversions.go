package conveyLoss

import (
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"strconv"
	"strings"
	"time"
)

// efPeriod is a struct to hold the data from excess_flow_periods
type efPeriod struct {
	CanalId   int          `db:"canal_id"`
	StartDate sql.NullTime `db:"st_date"`
	EndDate   sql.NullTime `db:"end_date"`
}

// diversion is a struct to hold the dailydiversions table data and also the results are this struct which is a monthly
// total using the first day of each month.
type diversion struct {
	CanalId   int             `db:"canal_id"`
	DivDate   sql.NullTime    `db:"div_dt"`
	DivAmount sql.NullFloat64 `db:"div_amnt_cfs"`
}

// applyEffAcres is a method that multiplies the diversion by the efficiency and acres passed in and converts the day cfs
// to acre-feet to yield a application in inches per acre from that surface water structure.
func (d *diversion) applyEffAcres(eff float64, acres float64) {
	d.DivAmount.Float64 = d.DivAmount.Float64 * eff * 1.9835 / acres
}

// getDiversions retrieves the diversions from the pg database and returns a slice of diversion struct for each canal
// during the year and also takes in a start year, end year and also excessFlow bool that if false will remove the
// excess flow from the daily diversions based on excess flow periods.
func getDiversions(pgDb *sqlx.DB, sYear int, eYear int, excessFlow bool) (diversions []diversion) {

	if excessFlow {
		divQry := fmt.Sprintf(`select canal_id, make_timestamp(cast(extract(YEAR from div_dt) as int), 
cast(extract(MONTH from div_dt) as int), 1, 0, 0, 0) as div_dt, sum(div_amnt_cfs) as div_amnt_cfs 
from (select cdj.canal_id, div_dt, sum(div_amnt_cfs) as div_amnt_cfs from wwnp.dailydiversions 
inner join wwnp.canal_diversion_jct cdj on dailydiversions.ndnr_id = cdj.ndnr_id WHERE 
div_dt >= '%d-01-01' AND div_dt <= '%d-12-31' group by cdj.canal_id, div_dt) as daily_query group by canal_id, 
extract(MONTH from div_dt), extract(YEAR from div_dt) order by canal_id, div_dt;`, sYear, eYear)

		err := pgDb.Select(&diversions, divQry)
		if err != nil {
			fmt.Println("Error in getting all diversions", err)
		}

	} else {
		// get all diversion data
		divQry := fmt.Sprintf(`select cdj.canal_id, div_dt, sum(div_amnt_cfs) as div_amnt_cfs from sw.dailydiversions
inner join sw.canal_diversion_jct cdj on dailydiversions.ndnr_id = cdj.ndnr_id
WHERE div_dt between '%d-01-01' and '%d-12-31' group by cdj.canal_id, div_dt;`, sYear, eYear)

		var preDiversions []diversion
		err := pgDb.Select(&preDiversions, divQry)
		if err != nil {
			fmt.Println("Error in getting all diversions starting Excess Flow Limitation", err)
		}

		// remove the excess flows
		efQuery := "select canal_id, st_date, end_date from sw.excess_flow_periods;"

		var efPeriods []efPeriod
		err = pgDb.Select(&efPeriods, efQuery)
		if err != nil {
			fmt.Println("Error in getting Excess Flow Periods", err)
		}

		// list of canals that have excess flows
		var efCanals []int
		for _, v := range efPeriods {
			if find(efCanals, v.CanalId) {
				efCanals = append(efCanals, v.CanalId)
			}
		}

		// loop through the list of canals with excess flows
		var canalMonthlyDiversion []diversion
		for _, canal := range efCanals {

			// generate list of excess flow periods for this canal
			var efCanalPeriods []efPeriod
			for _, v := range efPeriods {
				if v.CanalId == canal {
					efCanalPeriods = append(efCanalPeriods, v)
				}
			}

			// filter daily canal diversions to this canal
			var allCanalDailyDiversion []diversion
			for _, div := range preDiversions {
				if div.CanalId == canal {
					allCanalDailyDiversion = append(allCanalDailyDiversion, div)
				}
			}

			// filter out the times when there was excess flow and only return the times the flow wasn't during that
			// excess flow period
			var canalDailyDiversion []diversion
			for _, div := range allCanalDailyDiversion {
				if findDiversion(div, efCanalPeriods) {
					canalDailyDiversion = append(canalDailyDiversion, div)
				}
			}

			// reduce the daily flows to monthly
			for y := sYear; y < eYear+1; y++ {
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
group by canal_id, extract(MONTH from div_dt), extract(YEAR from div_dt) order by canal_id, div_dt;`, sYear, eYear, canalLimit)

		var remainingDiversions []diversion
		err = pgDb.Select(&remainingDiversions, remainDiversionQry)
		if err != nil {
			fmt.Println("Error in getting remaining diversions", err)
		}

		diversions = append(canalMonthlyDiversion, remainingDiversions...)
	}

	return diversions
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

// FindDiversion is a function to filter if a diversion is within a period
func findDiversion(div diversion, period []efPeriod) bool {
	for _, p := range period {
		if div.DivDate.Time.Before(p.StartDate.Time) || div.DivDate.Time.After(p.EndDate.Time) {
			return true
		}
	}

	return false
}

// monthlyDiversion is a function to aggregate a slice of dailyDiversion by months and returns a month total diversion
func monthlyDiversion(dailyDiversion []diversion, m int, y int, cID int) (mDiversion diversion) {
	var totalDiversion float64
	for _, d := range dailyDiversion {
		if d.DivDate.Time.Year() == y && int(d.DivDate.Time.Month()) == m && d.DivAmount.Valid {
			totalDiversion += d.DivAmount.Float64
		}
	}

	mDiversion = diversion{DivDate: sql.NullTime{Time: time.Date(y, time.Month(m), 1, 0, 0, 0, 0, time.UTC),
		Valid: true}, DivAmount: sql.NullFloat64{Float64: totalDiversion, Valid: true}, CanalId: cID}

	return mDiversion
}
