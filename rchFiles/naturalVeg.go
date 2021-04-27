package rchFiles

import (
	"database/sql"
	"fmt"
)

func NaturalVeg(db *sql.DB, debug *bool) {
	fmt.Println("welcome to NaturalVeg")
	fmt.Println("Debug is set to", debug)
}
