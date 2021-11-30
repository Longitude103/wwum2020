package parcels

func parcelsPost97(allParcels []Parcel, p97GWOParcels []Parcel) (p97Parcels []Parcel) {
	// delete GW Only parcels passed in
	parcelsWithOutGW := removeGWO(allParcels)

	// add GW Only parcels from 1997
	p97Parcels = append(parcelsWithOutGW, p97GWOParcels...)

	return p97Parcels
}

func removeGWO(allParcels []Parcel) (gwoRemoved []Parcel) {
	for _, parcel := range allParcels {
		if !parcel.isGWO() {
			gwoRemoved = append(gwoRemoved, parcel)
		}
	}

	return
}
