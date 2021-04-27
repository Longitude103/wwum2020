package database

import (
	"database/sql"
	"fmt"
)

// InitializeDb creates the database results table if it doesn't exist, so the records of the
// transaction can be stored properly, also creates file_keys table for result file_type integer
// and a foreign key restriction to results table
func InitializeDb(db *sql.DB) {
	stmt, err := db.Prepare(`
		create table if not exists file_keys
			(
				file_key integer not null
					constraint file_keys_pk
						primary key autoincrement,
				description TEXT
			);
			
		create unique index file_keys_file_key_uindex
				on file_keys (file_key);
	`)
	if err != nil {
		fmt.Println("Error", err)
	}

	// if table exists with records, then skip adding data.
	var count int
	err = db.QueryRow(`select count(*) from file_keys;`).Scan(&count)
	if err != nil {
		fmt.Println("Error", err)
	}

	if count == 0 {
		stmt, err = db.Prepare(`
		insert into file_keys (file_key, description) values (101, 'Dryland');
		insert into file_keys (file_key, description) values (102, 'Natural Vegetation');
		insert into file_keys (file_key, description) values (103, 'NP Surface Water Only');
		insert into file_keys (file_key, description) values (104, 'SP Surface Water Only');
		insert into file_keys (file_key, description) values (105, 'NP Comingled Parcel Pre 1998');
		insert into file_keys (file_key, description) values (106, 'NP Comingled Parcel Post 1997');
		insert into file_keys (file_key, description) values (107, 'SP Comingled Parcel Pre 1998');
		insert into file_keys (file_key, description) values (108, 'SP Comingled Parcel Post 1997');
		insert into file_keys (file_key, description) values (109, 'NP Groundwater Only Pre 1998');
		insert into file_keys (file_key, description) values (110, 'NP Groundwater Only Post 1997');
		insert into file_keys (file_key, description) values (111, 'SP Groundwater Only Pre 1998');
		insert into file_keys (file_key, description) values (112, 'SP Groundwater Only Post 1997');
		insert into file_keys (file_key, description) values (113, 'NP Canal Loss');
		insert into file_keys (file_key, description) values (114, 'SP Canal Loss');
		insert into file_keys (file_key, description) values (115, 'Outside NP and SP');
		insert into file_keys (file_key, description) values (117, 'NP Recharge Sites');
		insert into file_keys (file_key, description) values (118, 'SP Recharge Sites');
		insert into file_keys (file_key, description) values (201, 'NP Comingled Pre 1998');
		insert into file_keys (file_key, description) values (202, 'NP Groundwater Only Pre 1998');
		insert into file_keys (file_key, description) values (203, 'NP Comingled Post 1997');
		insert into file_keys (file_key, description) values (204, 'NP Groundwater Only Post 1997');
		insert into file_keys (file_key, description) values (205, 'SP Comingled Pre 1998');
		insert into file_keys (file_key, description) values (206, 'SP Groundwater Only Pre 1998');
		insert into file_keys (file_key, description) values (207, 'SP Comingled Post 1997');
		insert into file_keys (file_key, description) values (208, 'SP Groundwater Only Post 1997');
		insert into file_keys (file_key, description) values (209, 'Steady State');
		insert into file_keys (file_key, description) values (210, 'Municipal');
		insert into file_keys (file_key, description) values (211, 'Industrial');
		insert into file_keys (file_key, description) values (212, 'Other Wells');
		insert into file_keys (file_key, description) values (213, 'Western Canal Outside SP');
	`)
		if err != nil {
			fmt.Println("Error", err)
		}

		_, err = stmt.Exec()
		if err != nil {
			fmt.Println("Error", err)
		}
	}

	stmt, err = db.Prepare(`
		create table if not exists results
			(
				id integer not null
					constraint results_pk
						primary key autoincrement,
				cell_node int not null,
				file_type int not null
					constraint results_file_keys_file_key_fk
					references file_keys,
				result float,
				run_type int
			);
		
		create unique index results_id_uindex
			on results (id);
	`)
	if err != nil {
		fmt.Println("Error", err)
	}

	_, err = stmt.Exec()
	if err != nil {
		fmt.Println("Error", err)
	}

}
