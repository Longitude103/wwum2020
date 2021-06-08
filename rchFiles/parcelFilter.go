package rchFiles

import (
	"errors"
	"github.com/heath140/wwum2020/parcels"
)

func parcelFilterByYear(parcels []parcels.Parcel, yr int) (p []parcels.Parcel, err error) {
	if len(parcels) < 1 {
		return nil, errors.New("no parcels in slice")
	}

	for i := 0; i < len(parcels); i++ {
		if parcels[i].Yr == yr {
			p = append(p, parcels[i])
		}
	}

	return
}

func parcelFilterById(p []parcels.Parcel, id int) (parcels.Parcel, error) {
	if len(p) < 1 {
		return parcels.Parcel{}, errors.New("no parcels in slice for id")
	}

	for i := 0; i < len(p); i++ {
		if p[i].ParcelNo == id {
			return p[i], nil
		}
	}

	return parcels.Parcel{}, errors.New("no parcel found with that id")
}
