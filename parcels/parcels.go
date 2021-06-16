package parcels

import (
	"database/sql"
	"fmt"
	"github.com/Longitude103/wwum2020/database"
	"github.com/Longitude103/wwum2020/parcels/conveyLoss"
	"strconv"
	"strings"
	"time"
)

type Parcel struct {
	ParcelNo  int             `db:"parcel_id"`
	CertNum   sql.NullString  `db:"cert_num"`
	Area      float64         `db:"area"`
	Crop1     sql.NullInt64   `db:"crop1"`
	Crop2     sql.NullInt64   `db:"crop2"`
	Crop3     sql.NullInt64   `db:"crop3"`
	Crop4     sql.NullInt64   `db:"crop4"`
	IrrType   sql.NullString  `db:"irrig_type"`
	FirstIrr  sql.NullInt64   `db:"first_irr"`
	SwFac     sql.NullString  `db:"sw_fac"`
	ModelId   sql.NullString  `db:"model_id"`
	Crop1Cov  sql.NullFloat64 `db:"crop1_cov"`
	Crop2Cov  sql.NullFloat64 `db:"crop2_cov"`
	Crop3Cov  sql.NullFloat64 `db:"crop3_cov"`
	Crop4Cov  sql.NullFloat64 `db:"crop4_cov"`
	Sw        sql.NullBool    `db:"sw"`
	Gw        sql.NullBool    `db:"gw"`
	Nrd       string          `db:"nrd"`
	SwID      sql.NullInt64   `db:"sw_id"`
	PointX    float64         `db:"pointx"`
	PointY    float64         `db:"pointy"`
	SoilArea  float64         `db:"s_area"`
	SoilCode  int             `db:"soil_code"`
	CoeffZone int             `db:"coeff_zone"`
	Yr        int
	AppEff    float64
	Et        [12]float64
	DryEt     [12]float64
	Nir       [12]float64
	Ro        [12]float64
	Dp        [12]float64
	Pump      [12]float64
	SWDel     [12]float64
	Metered   bool
}

type IrrType int

const (
	Irrigated IrrType = 3
	DryLand   IrrType = 1
)

// getParcels returns a list of all parcels with crops irrigation types and areas. Returns data for both nrds. There
// can be multiples of the same parcels listed with different soil types. It sets the year into a field in the struct.
func getParcels(v database.Setup, Year int) []Parcel {
	query := fmt.Sprintf(`SELECT parcel_id, a.crop_int crop1, crop1_cov, b.crop_int crop2, crop2_cov, c.crop_int crop3, crop3_cov, d.crop_int crop4, crop4_cov, sw, gw,
       irrig_type, sw_fac, first_irr, cert_num::varchar, model_id, sw_id, st_area(i.geom)/43560 area, 'np' nrd,
       st_x(st_transform(st_centroid(i.geom), 4326)) pointx, st_y(st_transform(st_centroid(i.geom), 4326)) pointy,
       sum(st_area(st_intersection(m.geom, i.geom))/43560) s_area, m.soil_code, m.coeff_zone
FROM np.t%d_irr i inner join public.model_cells m on st_intersects(i.geom, m.geom)
                    LEFT join public.crops a on crop1 = a.crop_name
                    LEFT join public.crops b on crop2 = b.crop_name
                    LEFT join public.crops c on crop3 = c.crop_name
                    LEFT join public.crops d on crop4 = d.crop_name
GROUP BY parcel_id, a.crop_int, parcel_id, crop1_cov, b.crop_int, crop2_cov, c.crop_int, crop3_cov, d.crop_int, crop4_cov, sw, gw, irrig_type, sw_fac, cert_num::varchar, model_id, st_area(i.geom)/43560, st_x(st_transform(st_centroid(i.geom), 4326)), st_y(st_transform(st_centroid(i.geom), 4326)), m.soil_code, crop1_cov, crop2, crop2_cov, crop3, crop3_cov, crop4, crop4_cov, sw, gw, irrig_type,
         sw_fac, first_irr, cert_num::varchar, model_id, st_area(i.geom)/43560, st_x(st_transform(st_centroid(i.geom), 4326)),
         st_y(st_transform(st_centroid(i.geom), 4326)), nrd, m.soil_code, m.coeff_zone
UNION ALL
SELECT parcel_id, a.crop_int crop1, crop1_cov, b.crop_int crop2, crop2_cov, c.crop_int crop3, crop3_cov, d.crop_int crop4, crop4_cov, sw, gw,
       irr_type as irrig_type, sw_fac, first_irr, i.id as cert_num, null as model_id, sw_id, st_area(i.geom)/43560 area, 'sp' nrd,
       st_x(st_transform(st_centroid(i.geom), 4326)) pointx, st_y(st_transform(st_centroid(i.geom), 4326)) pointy,
       sum(st_area(st_intersection(m.geom, i.geom))/43560) s_area, m.soil_code, m.coeff_zone
FROM sp.t%d_irr i inner join public.model_cells m on st_intersects(i.geom, m.geom)
                    LEFT join public.crops a on crop1 = a.crop_name
                    LEFT join public.crops b on crop2 = b.crop_name
                    LEFT join public.crops c on crop3 = c.crop_name
                    LEFT join public.crops d on crop4 = d.crop_name
GROUP BY parcel_id, a.crop_int, parcel_id, crop1_cov, b.crop_int, crop2_cov, c.crop_int, crop3_cov, d.crop_int, crop4_cov, sw, gw, irrig_type,
         sw_fac, first_irr, i.id, model_id, st_area(i.geom)/43560, st_x(st_transform(st_centroid(i.geom), 4326)),
         st_y(st_transform(st_centroid(i.geom), 4326)), nrd, m.soil_code, m.coeff_zone;`,
		Year, Year)

	var parcels []Parcel
	err := v.PgDb.Select(&parcels, query)
	if err != nil {
		v.Logger.Errorf("Error in getting parcels for year %d, error: %s", Year, err)
	}

	for i := 0; i < len(parcels); i++ {
		parcels[i].Yr = Year
	}

	if v.AppDebug {
		return parcels[:10]
	}

	return parcels
}

// filterParcelByCert filters a slice of parcels by the CertNum and returns a slice of the parcels that have that CertNum.
func filterParcelByCert(p *[]Parcel, c string) (filteredParcels []Parcel) {
	for i := 0; i < len(*p); i++ {
		if (*p)[i].CertNum.String == c {
			filteredParcels = append(filteredParcels, (*p)[i])
		}
	}

	return filteredParcels
}

// parcelSWDelivery method uses the diversions to then calculate the total amount of surface water delivered to a parcel
// from those diversions. It returns nothing, but sets SWDel inside the Parcel
func (p *Parcel) parcelSWDelivery(diversions []conveyLoss.Diversion) {
	canalDivs := filterDivs(diversions, int(p.SwID.Int64))

	var swDelivery [12]float64
	for m := 0; m < 12; m++ {
		for _, d := range canalDivs {
			if int(d.DivDate.Time.Month()) == m+1 {
				swDelivery[m] = d.DivAmount.Float64 * p.Area
			}
		}
	}

	p.SWDel = swDelivery
}

// filterDivs function that receives a slice of divs and filters them to a canal that is given as an int
func filterDivs(divs []conveyLoss.Diversion, canal int) (d []conveyLoss.Diversion) {
	for _, v := range divs {
		if v.CanalId == canal {
			d = append(d, v)
		}
	}

	return d
}

func getDryParcels(v database.Setup, Year int) []Parcel {
	query := fmt.Sprintf(`SELECT i.parcel_id, a.crop_int crop1, crop1_cov, b.crop_int crop2, crop2_cov, c.crop_int crop3, crop3_cov, d.crop_int crop4, crop4_cov,
       st_area(i.geom)/43560 area, 'np' nrd, st_x(st_transform(st_centroid(i.geom), 4326)) pointx,
       st_y(st_transform(st_centroid(i.geom), 4326)) pointy, sum(st_area(st_intersection(m.geom, i.geom))/43560) s_area,
       m.soil_code, m.coeff_zone
FROM np.t%d_dry i inner join public.model_cells m on st_intersects(i.geom, m.geom)
    LEFT join public.crops a on crop1 = a.crop_name
    LEFT join public.crops b on crop2 = b.crop_name
    LEFT join public.crops c on crop3 = c.crop_name
    LEFT join public.crops d on crop4 = d.crop_name
GROUP BY i.parcel_id, a.crop_int, parcel_id, crop1_cov, b.crop_int, crop2_cov, c.crop_int, crop3_cov, d.crop_int, crop4_cov,
    st_area(i.geom)/43560, st_x(st_transform(st_centroid(i.geom), 4326)), st_y(st_transform(st_centroid(i.geom), 4326)),
    m.soil_code, crop1_cov, crop2, crop2_cov, crop3, crop3_cov, crop4, crop4_cov, st_area(i.geom)/43560,
    st_x(st_transform(st_centroid(i.geom), 4326)), st_y(st_transform(st_centroid(i.geom), 4326)), nrd, m.soil_code, m.coeff_zone
UNION ALL
SELECT i.parcel_id, a.crop_int crop1, crop1_cov, b.crop_int crop2, crop2_cov, c.crop_int crop3, crop3_cov, d.crop_int crop4, crop4_cov,
       st_area(i.geom)/43560 area, 'sp' nrd, st_x(st_transform(st_centroid(i.geom), 4326)) pointx,
       st_y(st_transform(st_centroid(i.geom), 4326)) pointy, sum(st_area(st_intersection(m.geom, i.geom))/43560) s_area,
       m.soil_code, m.coeff_zone
FROM sp.t%d_dry i inner join public.model_cells m on st_intersects(i.geom, m.geom)
    LEFT join public.crops a on crop1 = a.crop_name
    LEFT join public.crops b on crop2 = b.crop_name
    LEFT join public.crops c on crop3 = c.crop_name
    LEFT join public.crops d on crop4 = d.crop_name
GROUP BY i.parcel_id, a.crop_int, parcel_id, crop1_cov, b.crop_int, crop2_cov, c.crop_int, crop3_cov, d.crop_int, crop4_cov,
    st_area(i.geom)/43560, st_x(st_transform(st_centroid(i.geom), 4326)), st_y(st_transform(st_centroid(i.geom), 4326)),
    m.soil_code, crop1_cov, crop2, crop2_cov, crop3, crop3_cov, crop4, crop4_cov, st_area(i.geom)/43560,
    st_x(st_transform(st_centroid(i.geom), 4326)), st_y(st_transform(st_centroid(i.geom), 4326)), nrd, m.soil_code, m.coeff_zone;`, Year, Year)

	var parcels []Parcel
	err := v.PgDb.Select(&parcels, query)
	if err != nil {
		v.Logger.Errorf("Error in getting dryland parcels for year %d", Year)
	}

	for i := 0; i < len(parcels); i++ {
		parcels[i].Yr = Year
	}

	if v.AppDebug {
		return parcels[:20]
	}

	return parcels
}

func (p *Parcel) String() string {
	return fmt.Sprintf("Parcel No: %d, NRD: %s, Year: %d", p.ParcelNo, p.Nrd, p.Yr)
}

func (p *Parcel) PrintNIR() string {
	str := strings.Builder{}

	for m, v := range p.Nir {
		str.WriteString(time.Month(m + 1).String())
		str.WriteString(": ")
		str.WriteString(strconv.FormatFloat(v, 'f', 2, 64))
		str.WriteRune('\n')
	}

	return str.String()
}

func (p *Parcel) GetXY() (x float64, y float64) {
	return p.PointX, p.PointY
}
