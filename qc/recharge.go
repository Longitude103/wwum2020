package qc

import (
	"bufio"
	"fmt"
	"github.com/paulmach/orb/geojson"
	"github.com/pterm/pterm"
	"os"
	"path/filepath"
	"time"
)

type rchResults struct {
	FileType int     `db:"file_type"`
	Desc     string  `db:"description"`
	Recharge float64 `db:"recharge"`
}

func (q *QC) rechargeBalance() error {
	var rchResultsSlice []rchResults

	spin, _ := pterm.DefaultSpinner.Start(fmt.Sprintf("Getting %d Data", q.Year))
	query := "select file_type, description, sum(result) recharge from results inner join file_keys fk on fk.file_key = results.file_type where strftime('%Y', dt) = "
	query += fmt.Sprintf("'%d' group by file_type, description;", q.Year)

	if err := q.v.SlDb.Select(&rchResultsSlice, query); err != nil {
		return err
	}
	d := pterm.TableData{{"File Type", "Description", "Recharge (AF)"}}
	total := 0.0

	for _, r := range rchResultsSlice {
		d = append(d, []string{fmt.Sprintf("%d", r.FileType), r.Desc, fmt.Sprintf("%.0f", r.Recharge)})
		total += r.Recharge
	}

	d = append(d, []string{"TOTAL", "All File Types", fmt.Sprintf("%.0f", total)})
	spin.Success()

	pterm.DefaultSection.Println("Recharge Summary")
	if err := pterm.DefaultTable.WithHasHeader().WithRightAlignment().WithData(d).Render(); err != nil {
		return err
	}

	return nil
}

type modelCells struct {
	Gjson []byte  `db:"geojson"`
	Ac    float64 `db:"area_ac"`
	Node  int     `db:"node"`
}

type resultData struct {
	Node int     `db:"node"`
	Rslt float64 `db:"rslt"`
}

func (q *QC) rechargeGeoJson() error {
	// get a slice of model cells in geojson
	pterm.DefaultSection.Println("GeoJSON Creation")
	spin, _ := pterm.DefaultSpinner.Start("Getting Model Cells from DB")
	var mCells []modelCells
	qry := "select st_asgeojson(q) geojson, area_ac, node from (select st_transform(geom, 4326), node, st_area(geom)/43560 area_ac from model_cells) q;"

	if err := q.v.PgDb.Select(&mCells, qry); err != nil {
		return err
	}
	spin.Success()

	spin, _ = pterm.DefaultSpinner.Start("Getting Result Data")
	rResMap := make(map[int][]resultData)
	for m := 1; m < 13; m++ {
		var rResults []resultData
		var mnString string
		if m < 10 {
			mnString = fmt.Sprintf("0%d", m)
		} else {
			mnString = fmt.Sprintf("%d", m)
		}

		rqry := fmt.Sprintf("select cell_node node, sum(result) rslt from results "+
			"where strftime('%%Y', dt) = '%d' and strftime('%%m', dt) = '%s' group by cell_node, strftime('%%m', dt);", q.Year, mnString)

		if err := q.v.SlDb.Select(&rResults, rqry); err != nil {
			return err
		}

		rResMap[m] = rResults
	}
	spin.Success()

	fn := fmt.Sprintf("%d_Recharge.geojson", q.Year)
	wd, _ := os.Getwd()
	path := filepath.Join(wd, q.fileName)
	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		return err
	}
	writeFile, err := os.Create(filepath.Join(path, fn))
	if err != nil {
		return err
	}

	// Follows format of https://datatracker.ietf.org/doc/html/rfc7946#section-1.5
	w := bufio.NewWriter(writeFile)

	header := `{"type":"FeatureCollection","features":[`

	_, _ = w.WriteString(header)
	cL := len(mCells)

	// unmarshal each item
	p, _ := pterm.DefaultProgressbar.WithTotal(cL).WithTitle("Model Cells").WithRemoveWhenDone(true).Start()
	for i := 0; i < cL; i++ {
		p.Increment()
		fc, err := geojson.UnmarshalFeature(mCells[i].Gjson)
		if err != nil {
			return err
		}

		// add property to them of the monthly result
		annTotal := 0.0
		for m := 1; m < 13; m++ {
			mn := time.Month(m)
			res := findResult(rResMap[m], mCells[i].Node)
			if q.Monthly {
				fc.Properties[mn.String()] = res
				fc.Properties[mn.String()+"_rate"] = res / mCells[i].Ac
			}
			annTotal += res
		}

		fc.Properties["Annual_Total"] = annTotal
		fc.Properties["AnnTotal_rate"] = annTotal / mCells[i].Ac

		// marshal that item back to json
		d, err := fc.MarshalJSON()
		if err != nil {
			return err
		}

		if _, err := w.WriteString(string(d)); err != nil {
			return err
		}

		if i < cL-1 {
			_, _ = w.WriteString(", ")
		}
	}

	_, _ = w.WriteString("]}")
	_ = w.Flush()
	_ = writeFile.Close()
	pterm.Success.Println("Check Output Files for GeoJson")

	return nil
}

func findResult(rData []resultData, node int) float64 {
	for i := 0; i < len(rData); i++ {
		if rData[i].Node == node {
			return rData[i].Rslt
		}
	}

	return 0
}
