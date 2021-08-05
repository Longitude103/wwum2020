package database

import "errors"

type wellParcel struct {
	ParcelId int    `db:"parcel_id"`
	WellId   int    `db:"wellid"`
	Nrd      string `db:"nrd"`
	Yr       int    `db:"yr"`
}

func GetWellParcels(v Setup) ([]wellParcel, error) {
	const query = "select parcel_id, wellid, nrd, yr from public.alljct();"

	var wellParcels []wellParcel
	if err := v.PgDb.Select(&wellParcels, query); err != nil {
		return wellParcels, errors.New("error getting parcel wells from db function")
	}

	if v.AppDebug {
		return wellParcels[:50], nil
	}

	return wellParcels, nil
}
