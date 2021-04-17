package distribution

import (
	"clibasic/color"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

const (
	host     = "long103-wwum.clmtjoquajav.us-east-2.rds.amazonaws.com"
	port     = 5432
	user     = "postgres"
	password = "rQ!461k&Rk8J"
	dbname   = "wwum"
)

type actCell struct {
	Rw        int
	Clm       int
	SoilCode  int
	TfgCellId int
}

func Distribution(debug *bool, startYr *int, endYr *int) {
	fmt.Println("Distribution")
	if *debug {
		fmt.Println(color.Red + "Debug Mode" + color.Reset)
	}

	fmt.Printf("Start Year: %d -> End Year %d\n", *startYr, *endYr)

	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	// open database
	db, err := sql.Open("postgres", psqlconn)
	CheckError(err)

	// close database
	defer db.Close()

	// check db
	err = db.Ping()
	CheckError(err)

	fmt.Println("Connected!")

	rows, err := db.Query(`SELECT rw, clm, soil_code, tfg_cellid FROM public.act_cells`)
	if err != nil {
		panic(err)
	}

	defer rows.Close()

	var Cells []actCell
	cell := actCell{}
	for rows.Next() {
		var rw, clm, soil, cellid int
		err = rows.Scan(&rw, &clm, &soil, &cellid)
		if err != nil {
			panic(err.Error())
		}
		cell.Rw = rw
		cell.Clm = clm
		cell.SoilCode = soil
		cell.TfgCellId = cellid
		Cells = append(Cells, cell)
	}

	for _, c := range Cells[:5] {
		fmt.Println(c)
	}
}

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}
