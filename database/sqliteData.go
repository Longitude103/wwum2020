package database

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"time"
)

type FileKeys struct {
	FileKey     int    `db:"file_key"`
	Description string `db:"description"`
}

func (f FileKeys) Print() string {
	return fmt.Sprintf("%d - %s", f.FileKey, f.Description)
}

type MfResults struct {
	CellNode   int       `db:"cell_node"`
	ResultDate time.Time `db:"dt"`
	Rslt       float64   `db:"rslt"`
}

func (m MfResults) Date() time.Time {
	return m.ResultDate
}

func (m MfResults) Node() int {
	return m.CellNode
}

func (m MfResults) Value() float64 {
	return m.Rslt
}

func GetFileKeys(db *sqlx.DB, wel bool) ([]string, error) {
	var fKeys []FileKeys
	var resultFileKeys []string
	var query string
	if wel { // give the wel file_keys
		query = fmt.Sprint("SELECT file_key, description FROM file_keys WHERE file_key > 199;")
	} else { // give the rch file_keys
		query = fmt.Sprint("SELECT file_key, description FROM file_keys WHERE file_key < 200;")
	}

	if err := db.Select(&fKeys, query); err != nil {
		return resultFileKeys, err
	}

	for _, k := range fKeys {
		resultFileKeys = append(resultFileKeys, k.Print())
	}

	return resultFileKeys, nil
}

func GetAggResults(db *sqlx.DB, wel bool, excludeList []string) ([]MfResults, error) {
	var qry string
	var results []MfResults

	if wel { // is a wel file
		if len(excludeList) > 0 { // has an item in the to exclude list
			list := excludeList[0][0:3]
			for i := 1; i < len(excludeList); i++ {
				list += ", "
				list += excludeList[i][0:3]
			}

			qry = fmt.Sprintf("SELECT cell_node, dt, rslt from (SELECT cell_node, dt, sum(result) rslt "+
				"FROM wel_results WHERE file_type NOT IN (%s) group by cell_node, dt) where rslt > 0;", list)
		} else { // don't exclude anything
			qry = fmt.Sprint("SELECT cell_node, dt, rslt from (SELECT cell_node, dt, sum(result) rslt " +
				"FROM wel_results group by cell_node, dt) where rslt > 0;")
		}
	} else { // is a recharge file
		if len(excludeList) > 0 { // has an item in exclude list
			list := excludeList[0][0:3]
			for i := 1; i < len(excludeList); i++ {
				list += ", "
				list += excludeList[i][0:3]
			}

			qry = fmt.Sprintf("SELECT cell_node, dt, rslt from (SELECT cell_node, dt, sum(result) rslt "+
				"FROM results WHERE file_type NOT IN (%s) group by cell_node, dt) where rslt > 0;", list)
		} else { // don't exclude anything
			qry = fmt.Sprint("SELECT cell_node, dt, rslt from (SELECT cell_node, dt, sum(result) rslt " +
				"FROM results group by cell_node, dt) where rslt > 0;")
		}
	}

	if err := db.Select(&results, qry); err != nil {
		return results, err
	}

	return results, nil
}

func SingleResult(db *sqlx.DB, wel bool, fileKey string) ([]MfResults, error) {
	var results []MfResults
	var qry string

	if wel {
		qry = fmt.Sprintf("SELECT cell_node, dt, rslt from (SELECT cell_node, dt, sum(result) rslt FROM wel_results "+
			"WHERE file_type = %s group by cell_node, dt) where rslt > 0;", fileKey[0:3])
	} else {
		qry = fmt.Sprintf("SELECT cell_node, dt, rslt from (SELECT cell_node, dt, sum(result) rslt FROM results "+
			"WHERE file_type = %s group by cell_node, dt) where rslt > 0;", fileKey[0:3])
	}

	if err := db.Select(&results, qry); err != nil {
		return results, err
	}

	return results, nil
}
