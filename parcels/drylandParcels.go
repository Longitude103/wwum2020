package parcels

import (
	"github.com/heath140/wwum2020/fileio"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

func DryLandParcels(pgDB *sqlx.DB, slDB *sqlx.DB, sYear int, eYear int, csResults *map[string][]fileio.StationResults, logger *zap.SugaredLogger) (dryParcels []Parcel) {

	return dryParcels
}
