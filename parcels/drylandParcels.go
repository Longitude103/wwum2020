package parcels

import (
	"github.com/heath140/wwum2020/database"
	"github.com/heath140/wwum2020/fileio"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

func DryLandParcels(pgDB *sqlx.DB, pNirDB *database.DB, sYear int, eYear int, csResults *map[string][]fileio.StationResults,
	wStations []database.WeatherStation, logger *zap.SugaredLogger) (dryParcels []Parcel, err error) {

	logger.Info("Getting parcels")
	for y := sYear; y < eYear+1; y++ {
		dryParcels = getDryParcels(pgDB, y, logger)

		for i := 0; i < len(dryParcels); i++ {
			err = (&dryParcels[i]).parcelNIR(pNirDB, y, wStations, *csResults, DryLand)
		}
		if err != nil {
			return nil, err
		}
	}

	return dryParcels, nil
}
