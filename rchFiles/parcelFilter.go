package rchFiles

import (
	"errors"
	"fmt"
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

func parcelFilterById(p []parcels.Parcel, id int, nrd string) (parcels.Parcel, error) {
	if len(p) < 1 {
		return parcels.Parcel{}, errors.New("no parcels in slice for id")
	}

	fmt.Println("parcels length", len(p))
	for i := 0; i < len(p); i++ {
		if p[i].ParcelNo == id && p[i].Nrd == nrd {
			return p[i], nil
		}
	}

	errMessage := fmt.Sprintf("no parcel with id: %d", id)
	return parcels.Parcel{}, errors.New(errMessage)
}
