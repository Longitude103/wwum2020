package parcelpump

import (
	"database/sql"
	"fmt"
	"github.com/heath140/wwum2020/parcelpump/conveyLoss"
	"github.com/jmoiron/sqlx"
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
	AppEff    float64
	Nir       [12]float64
	Ro        [12]float64
	Dp        [12]float64
	Usage     [12]float64
	SWDel     [12]float64
	Metered   bool
}

// getParcels returns a list of all parcels with crops irrigation types and areas. Returns data for both nrds. There
// can be multiples of the same parcels listed with different soil types.
// Need to implement multiple years
func getParcels(db *sqlx.DB, Year int) []Parcel {
	query := fmt.Sprintf(`SELECT parcel_id, a.crop_int crop1, crop1_cov, b.crop_int crop2, crop2_cov, c.crop_int crop3, crop3_cov, d.crop_int crop4, crop4_cov, sw, gw,
       irrig_type, sw_fac, cert_num::varchar, model_id, sw_id, st_area(i.geom)/43560 area, 'np' nrd,
       st_x(st_transform(st_centroid(i.geom), 4326)) pointx, st_y(st_transform(st_centroid(i.geom), 4326)) pointy,
       sum(st_area(st_intersection(m.geom, i.geom))/43560) s_area, m.soil_code, m.coeff_zone
FROM np.t%d_irr i inner join public.model_cells m on st_intersects(i.geom, m.geom)
    LEFT join public.crops a on crop1 = a.crop_name
    LEFT join public.crops b on crop2 = b.crop_name
    LEFT join public.crops c on crop3 = c.crop_name
    LEFT join public.crops d on crop4 = d.crop_name
GROUP BY parcel_id, a.crop_int, parcel_id, crop1_cov, b.crop_int, crop2_cov, c.crop_int, crop3_cov, d.crop_int, crop4_cov, sw, gw, irrig_type, sw_fac, cert_num::varchar, model_id, st_area(i.geom)/43560, st_x(st_transform(st_centroid(i.geom), 4326)), st_y(st_transform(st_centroid(i.geom), 4326)), m.soil_code, crop1_cov, crop2, crop2_cov, crop3, crop3_cov, crop4, crop4_cov, sw, gw, irrig_type,
         sw_fac, cert_num::varchar, model_id, st_area(i.geom)/43560, st_x(st_transform(st_centroid(i.geom), 4326)),
         st_y(st_transform(st_centroid(i.geom), 4326)), nrd, m.soil_code
UNION ALL
SELECT parcel_id, a.crop_int crop1, crop1_cov, b.crop_int crop2, crop2_cov, c.crop_int crop3, crop3_cov, d.crop_int crop4, crop4_cov, sw, gw,
       irr_type as irrig_type, sw_fac, i.id as cert_num, null as model_id, sw_id, st_area(i.geom)/43560 area, 'sp' nrd,
       st_x(st_transform(st_centroid(i.geom), 4326)) pointx, st_y(st_transform(st_centroid(i.geom), 4326)) pointy,
       sum(st_area(st_intersection(m.geom, i.geom))/43560) s_area, m.soil_code
FROM sp.t%d_irr i inner join public.model_cells m on st_intersects(i.geom, m.geom)
                    LEFT join public.crops a on crop1 = a.crop_name
                    LEFT join public.crops b on crop2 = b.crop_name
                    LEFT join public.crops c on crop3 = c.crop_name
                    LEFT join public.crops d on crop4 = d.crop_name
GROUP BY parcel_id, a.crop_int, parcel_id, crop1_cov, b.crop_int, crop2_cov, c.crop_int, crop3_cov, d.crop_int, crop4_cov, sw, gw, irrig_type,
         sw_fac, i.id, model_id, st_area(i.geom)/43560, st_x(st_transform(st_centroid(i.geom), 4326)),
         st_y(st_transform(st_centroid(i.geom), 4326)), nrd, m.soil_code;`,
		Year, Year)

	var parcels []Parcel
	err := db.Select(&parcels, query)
	if err != nil {
		fmt.Println("Error", err)
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

func filterDivs(divs []conveyLoss.Diversion, canal int) (d []conveyLoss.Diversion) {
	for _, v := range divs {
		if v.CanalId == canal {
			d = append(d, v)
		}
	}

	return d
}
