package database

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
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

type SSWell struct {
	Id       int `db:"id"`
	WellName int `db:"wellname"`
	Rate     int `db:"defaultq"`
	Node     int `db:"node"`
	MVolume  [12]float64
}

type MIWell struct {
	WellId     int    `db:"id"`
	WellName   string `db:"wellname"`
	Rate       int    `db:"defaultq"`
	MuniWell   bool   `db:"muni_well"`
	IndustWell bool   `db:"indust_well"`
	Stop97     bool   `db:"stop_97"`
	Start97    bool   `db:"start_97"`
	Node       int    `db:"node"`
	Pumping    []MIPumping
}

type MIPumping struct {
	WellId   int       `db:"well_id"`
	PumpDate time.Time `db:"dt"`
	Pump     float64   `db:"pumping"`
}

// GetWellParcels is a function that gets all the well parcel junction table values and creates one struct from them
// and also includes the year of the join as well as the nrd.
func GetWellParcels(v *Setup) ([]WellParcel, error) {
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
func GetWellNode(v *Setup) (wellNodes []WellNode, err error) {
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

// GetSSWells is a function that gets the data from the postgres DB and returns a slice of SSWell and also includes a call
// to the SSWell.monthlyVolume() method to set the monthly data from the annual data that is in the database.
func GetSSWells(v *Setup) (ssWells []SSWell, err error) {
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

// monthlyVolume is a method of SSWell that calculates the monthly volume of pumping from the rate that is included in the
// database records. It turns the value positive to make it uniform with the other results.
func (s *SSWell) monthlyVolume() (err error) {
	const daysInMonth = 30.436875
	annVolume := -1.0 * float64(s.Rate) * 365.25 / 43560

	for i := 0; i < 12; i++ {
		s.MVolume[i] = (annVolume / daysInMonth) * -1
	}

	return nil
}

func GetMIWells(v *Setup) (miWells []MIWell, err error) {
	const miQuery = "SELECT mi.id, mi.wellname, mi.defaultq, mi.muni_well, mi.indust_well, mi.stop_97, mi.start_97, " +
		"mc.node FROM mi_wells mi inner join model_cells mc on st_contains(mc.geom, st_translate(mi.geom, 20, 20));"

	if err = v.PgDb.Select(&miWells, miQuery); err != nil {
		return miWells, errors.New("error getting Municipal and industrial wells")
	}

	// use setup to get bounds of pumping (earliest is 1997)
	miPumpQuery := fmt.Sprintf("SELECT well_id, dt, pumping FROM mi_pumping where "+
		"extract(YEAR from dt) between %d and %d", v.SYear, v.EYear)

	var miPump []MIPumping

	if err = v.PgDb.Select(&miPump, miPumpQuery); err != nil {
		return miWells, errors.New("error getting Pumping for M&I wells")
	}

	for i := 0; i < len(miWells); i++ {
		for _, p := range miPump {
			if p.WellId == miWells[i].WellId {
				miWells[i].Pumping = append(miWells[i].Pumping, p)
			}
		}
	}

	return
}

func (w *MIWell) MIFileType() int {
	if w.MuniWell {
		return 210
	} else if w.IndustWell {
		return 211
	} else {
		return 212
	}
}
