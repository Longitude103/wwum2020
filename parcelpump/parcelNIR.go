package parcelpump

import (
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
)

func parcelNIR(pgDB *sqlx.DB, db *sql.DB, sYear int, eYear int, parcels []Parcel) {
	for _, parcel := range parcels[:5] {
		fmt.Println(parcel)
	}
}
