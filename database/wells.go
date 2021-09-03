package database

import (
	"database/sql"
	"errors"
	"fmt"
)

type WellParcel struct {
	ParcelId int    `db:"parcel_id"`
	WellId   int    `db:"wellid"`
	Nrd      string `db:"nrd"`
	Yr       int    `db:"yr"`
}

type WellNode struct {
	WellId int            `db:"wellid"`
	RegCd  sql.NullString `db:"regcd"`
	Node   int            `db:"node"`
	Nrd    string         `db:"nrd"`
}

type SSWells struct {
	Id       int `db:"id"`
	WellName int `db:"wellname"`
	Rate     int `db:"defaultq"`
	Node     int `db:"node"`
	MVolume  [12]float64
}

// GetWellParcels is a function that gets all the well parcel junction table values and creates one struct from them
// and also includes the year of the join as well as the nrd.
func GetWellParcels(v Setup) ([]WellParcel, error) {
	const query = "select parcel_id, wellid, nrd, yr from public.alljct();"

	var wellParcels []WellParcel
	if err := v.PgDb.Select(&wellParcels, query); err != nil {
		fmt.Println("Err: ", err)
		return wellParcels, errors.New("error getting parcel wells from db function")
	}

	if v.AppDebug {
		return wellParcels[:50], nil
	}

	return wellParcels, nil
}

// GetWellNode is a function that gets the wellid, regno and node number of the well so that we can add a location to
// the well when it is written out along with the nrd.
func GetWellNode(v Setup) (wellNodes []WellNode, err error) {
	const query = "select wellid, regcd, node, 'np' nrd from np.npnrd_wells nw inner join model_cells mc " +
		"on st_contains(mc.geom, nw.geom) union all select wellid, regcd, node, 'sp' nrd from sp.spnrd_wells sw " +
		"inner join model_cells mc on st_contains(mc.geom, sw.geom)"

	if err := v.PgDb.Select(&wellNodes, query); err != nil {
		return wellNodes, errors.New("error getting well node locations from DB\n")
	}

	if v.AppDebug {
		return wellNodes[:50], nil
	}

	return wellNodes, nil
}

func GetSSWells(v Setup) (ssWells []SSWells, err error) {
	const ssQuery = "select ss_wells.id, wellname, defaultq, node from ss_wells inner join model_cells mc on " +
		"st_contains(mc.geom, st_translate(ss_wells.geom, 20, 20));"

	if err := v.PgDb.Select(&ssWells, ssQuery); err != nil {
		return ssWells, errors.New("error getting steady state wells from DB\n")
	}

	for i := 0; i < len(ssWells); i++ {
		if err := ssWells[i].monthlyVolume(); err != nil {
			return ssWells, errors.New("error setting monthly volumes\n")
		}
	}

	if v.AppDebug {
		return ssWells[:50], nil
	}

	return ssWells, nil
}

func (s *SSWells) monthlyVolume() (err error) {
	const daysInMonth = 30.436875
	annVolume := -1.0 * float64(s.Rate) * 365.25 / 43560

	for i := 0; i < 12; i++ {
		s.MVolume[i] = annVolume / daysInMonth
	}

	return nil
}