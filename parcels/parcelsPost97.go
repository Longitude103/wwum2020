package parcels

func ParcelsPost97(allParcels []Parcel, p97GWOParcels []Parcel) (p97Parcels []Parcel) {
	// delete GW Only parcels passed in
	parcelsWithOutGW := RemoveGWO(allParcels)

	// add GW Only parcels from 1997
	p97Parcels = append(parcelsWithOutGW, p97GWOParcels...)

	return p97Parcels
}

func RemoveGWO(allParcels []Parcel) (gwoRemoved []Parcel) {
	for _, parcel := range allParcels {
		if !parcel.IsGWO() {
			gwoRemoved = append(gwoRemoved, parcel)
		}
	}

	return
}
