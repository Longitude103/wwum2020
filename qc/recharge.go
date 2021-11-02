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
	Mnth int     `db:"mnth"`
	Rslt float64 `db:"rslt"`
}

func (q *QC) rechargeGeoJson() error {
	// get a slice of model cells in geojson
	var mCells []modelCells
	qry := "select st_asgeojson(q) geojson, area_ac, node from (select geom, node, st_area(geom)/43560 area_ac from model_cells) q;"

	if err := q.v.PgDb.Select(&mCells, qry); err != nil {
		return err
	}

	var rResults []resultData
	rqry := "select cell_node node, strftime('%m', dt) mnth, sum(result) rslt from results where strftime('%Y', dt) = '1953' group by cell_node, strftime('%m', dt);"

	if err := q.v.SlDb.Select(&rResults, rqry); err != nil {
		return err
	}

	wfp := filepath.Join("./", "OutputFiles", "GeoJSON.geojson")
	writeFile, err := os.Create(wfp)
	if err != nil {
		return err
	}

	w := bufio.NewWriter(writeFile)
	// TODO: need to write a header file like: https://datatracker.ietf.org/doc/html/rfc7946#section-1.5

	// unmarshal each item
	for i := 0; i < len(mCells); i++ {
		fc, err := geojson.UnmarshalFeature(mCells[i].Gjson)
		if err != nil {
			return err
		}

		// add property to them of the monthly result
		for i := 0; i < 12; i++ {
			m := time.Month(i + 1)
			res := findResult(rResults, i+1, mCells[i].Node)
			fc.Properties[m.String()] = res
			fc.Properties[m.String()+"_rate"] = res / mCells[i].Ac
		}

		// marshal that item back to json
		d, err := fc.MarshalJSON()
		if err != nil {
			return err
		}

		if _, err := w.WriteString(string(d)); err != nil {
			return err
		}
		// TODO: need to write a comma between them

		// testing it
		//fmt.Println(string(d))
		if i == 5 {
			break
		}
	}

	// TODO: need to add the footer which is the close for the array "]", and close for the collection "}"

	_ = w.Flush()

	_ = writeFile.Close()
	//fc.Properties["Result"] = 5
	//
	//// just checking on it
	//for s, i := range fc.Properties {
	//	if s == "node" {
	//		fmt.Println("Node number is:", i)
	//	}
	//}

	// write to geojson file

	return nil
}

func findResult(rData []resultData, mon int, node int) float64 {
	for i := 0; i < len(rData); i++ {
		if rData[i].Node == node && rData[i].Mnth == mon {
			return rData[i].Rslt
		}
	}

	return 0
}
