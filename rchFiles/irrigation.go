package rchFiles

import (
	"errors"
	"github.com/heath140/wwum2020/database"
	"github.com/heath140/wwum2020/parcels"
	"time"
)

func IrrigationRCH(v database.Setup, AllParcels []parcels.Parcel) error {
	v.Logger.Info("Starting to write RCH information from Irrigated Parcels")
	for y := v.SYear; y < v.EYear+1; y++ {
		// filter all parcels to this year only
		parcelList, err := parcelFilterByYear(AllParcels, y)
		if err != nil {
			return err
		}

		irrCells, err := database.GetCellsIrr(v, y)
		if err != nil {
			return err
		}

		// use the RO + DP from parcel and split by acres to get recharge, will need to keep separate files for the various
		// distributions of scenarios.
		for i := 0; i < len(irrCells); i++ {
			p, err := parcelFilterById(parcelList, irrCells[i].ParcelId, irrCells[i].Nrd)
			if err != nil {
				return err
			}
			fileType, err := assignRCHType(p.Nrd, p.Sw.Bool, p.Gw.Bool, post97(p.FirstIrr.Int64))
			if err != nil {
				return err
			}

			cellRecharge := [12]float64{}
			for j := 0; j < 12; j++ {
				cellRecharge[j], err = cellRCH(p.Ro[j], p.Dp[j], p.Area, irrCells[i].IrrArea)
				if err != nil {
					return err
				}

				if cellRecharge[j] > 0 {
					v.Logger.Infow("recharge created and is:", cellRecharge[j])
					err = v.RchDb.Add(database.RchResult{Node: irrCells[i].Node,
						Dt:       time.Date(y, time.Month(j+1), 1, 0, 0, 0, 0, nil),
						FileType: fileType, Result: cellRecharge[j]})
					if err != nil {
						return err
					}
				}

			}
		}

	}

	return nil
}

func cellRCH(ro float64, dp float64, parcelArea float64, parcelInCellArea float64) (r float64, err error) {
	if parcelArea <= 0 {
		return 0.0, errors.New("total parcel area is zero, division by zero would occur")
	}
	r = ro*parcelInCellArea/parcelArea + dp*parcelInCellArea/parcelArea

	return r, nil
}

func assignRCHType(nrd string, sw bool, gw bool, post97 bool) (int, error) {
	switch {
	case sw && gw:
		// co-mingled
		if nrd == "sp" {
			if post97 {
				return 108, nil
			}
			return 107, nil
		}
		if post97 {
			return 106, nil
		}
		return 105, nil
	case sw && gw == false:
		// sw only
		if nrd == "sp" {
			return 104, nil
		}
		return 103, nil
	case sw == false && gw:
		// gw only
		if nrd == "sp" {
			if post97 {
				return 112, nil
			}
			return 111, nil
		}
		if post97 {
			return 110, nil
		}
		return 109, nil
	default:
		return 100, errors.New("cannot classify recharge")
	}
}

func post97(yr int64) bool {
	if yr > 1997 {
		return true
	}

	return false
}
