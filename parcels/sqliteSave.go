package parcels

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"time"
)

// saveSqlite function saves the data for the parcel into local sqlite so that additional error checking can be preformed
// without loosing the data.
func saveSqlite(slDB *sqlx.DB, parcelID int, nrd string, pNIR [12]float64, yr int) {
	tx := slDB.MustBegin()

	for i, v := range pNIR {
		dt := time.Date(yr, time.Month(i+1), 1, 0, 0, 0, 0, time.UTC)
		tx.MustExec("INSERT INTO parcelNIR (parcelID, nrd, dt, nir) VALUES ($1, $2, $3, $4)", parcelID, nrd, dt.Format(time.RFC3339), v)
	}

	err := tx.Commit()
	if err != nil {
		fmt.Println("Error in SQLite Commit", err)
	}
}

func bulkSaveSqlite(slDB *sqlx.DB, values []Pumping, logger *zap.SugaredLogger) (err error) {
	bulks := len(values) / 500
	var subVals []Pumping

	for i := 0; i < bulks; i++ {
		for c := 0; c < 500; c++ {
			subVals = append(subVals, values[i*500+c])
		}

		_, err = slDB.NamedExec(`INSERT INTO parcelPumping (parcelID, nrd, dt, pump) 
										VALUES (:parcelID, :nrd, :dt, :pump)`, subVals)
		if err != nil {
			logger.Errorf("Error inserting parcel pumping into sqlite results, error: %s", err)
		}
	}

	return err
}
