package database

import (
	"database/sql"
	"errors"
)

type WellParcel struct {
	ParcelId int    `db:"parcel_id"`
	WellId   int    `db:"wellid"`
	Nrd      string `db:"nrd"`
	Yr       int    `db:"yr"`
}

type WellNode struct {
	WellId sql.NullInt64  `db:"wellid"`
	RegCd  sql.NullString `db:"regcd"`
	Node   sql.NullInt64  `db:"node"`
	Nrd    string         `db:"nrd"`
}

// GetWellParcels is a function that gets all the well parcel junction table values and creates one struct from them
// and also includes the year of the join as well as the nrd.
func GetWellParcels(v Setup) ([]WellParcel, error) {
	const query = "select parcel_id, wellid, nrd, yr from public.alljct();"
	qry := "select parcel_id, wellid, 'NP' nrd, 2014 yr from np.t2014_jct;"

	var wellParcels []WellParcel
	if err := v.PgDb.Select(&wellParcels, qry); err != nil {
		return wellParcels, errors.New("error getting parcel wells from db function\n")
	}

	if v.AppDebug {
		return wellParcels[:50], nil
	}

	return wellParcels, nil
}

// GetWellNode is a function that gets the wellid, regno and node number of the well so that we can add a location to
// the well when it is written out along with the nrd.
func GetWellNode(v Setup) (wellNodes []WellNode, err error) {
	const query = "select wellid, regcd, node, 'NP' nrd from np.npnrd_wells nw inner join model_cells mc on st_contains(mc.geom, nw.geom) union all select wellid, regcd, node, 'SP' nrd from sp.spnrd_wells sw inner join model_cells mc on st_contains(mc.geom, sw.geom)"

	if err := v.PgDb.Select(&wellNodes, query); err != nil {
		return wellNodes, errors.New("error getting well node locations from DB\n")
	}

	if v.AppDebug {
		return wellNodes[:50], nil
	}

	return wellNodes, nil
}
