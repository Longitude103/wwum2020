package database

import (
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

// InitializeDb creates the database results table if it doesn't exist, so the records of the
// transaction can be stored properly, also creates file_keys table for result file_type integer
// and a foreign key restriction to results table
func InitializeDb(db *sqlx.DB, logger *zap.SugaredLogger) {
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
		logger.Errorf("Error in create statement for file keys: %s", err)
	}

	_, err = stmt.Exec()
	if err != nil {
		logger.Errorf("Error on create file_keys: %s", err)
	}

	// if table exists with records, then skip adding data.
	var count int
	err = db.QueryRow(`select count(*) from file_keys;`).Scan(&count)
	if err != nil {
		logger.Errorf("Error in count file_key records %s", err)
	}

	if count == 0 {
		stmt, err = db.Prepare(`
		insert into file_keys (file_key, description) values (101, 'Dryland'), (102, 'Natural Vegetation'), 
		  (103, 'NP Surface Water Only'), (104, 'SP Surface Water Only'), (105, 'NP Comingled Parcel Pre 1998'), 
		  (106, 'NP Comingled Parcel Post 1997'), (107, 'SP Comingled Parcel Pre 1998'), (108, 'SP Comingled Parcel Post 1997'), 
		  (109, 'NP Groundwater Only Pre 1998'), (110, 'NP Groundwater Only Post 1997'), (111, 'SP Groundwater Only Pre 1998'), 
		  (112, 'SP Groundwater Only Post 1997'), (113, 'NP Canal Loss'), (114, 'SP Canal Loss'), (115, 'Outside NP and SP'), 
		  (116, 'NP Recharge Sites'), (117, 'SP Recharge Sites'), (201, 'NP Comingled Pre 1998'), (202, 'NP Groundwater Only Pre 1998'), 
		  (203, 'NP Comingled Post 1997'), (204, 'NP Groundwater Only Post 1997'), (205, 'SP Comingled Pre 1998'), 
		  (206, 'SP Groundwater Only Pre 1998'), (207, 'SP Comingled Post 1997'), (208, 'SP Groundwater Only Post 1997'), 
          (209, 'Steady State'), (210, 'Municipal'), (211, 'Industrial'), (212, 'Other Wells'), (213, 'Western Canal Outside SP');
	`)
		if err != nil {
			logger.Errorf("Error in statement of key records: %s", err)
		}

		_, err = stmt.Exec()
		if err != nil {
			logger.Errorf("Error in insert of key records: %s", err)
		}
	}

	stmt, err = db.Prepare(`
		create table if not exists results
			(
				id integer not null
					constraint results_pk
						primary key autoincrement,
				cell_node int not null,
				dt text,
				file_type int not null
					constraint results_file_keys_file_key_fk
					references file_keys,
				result float
			);
		
		create unique index results_id_uindex
			on results (id);
	`)
	if err != nil {
		logger.Errorf("Error in creating results table statement: %s", err)
	}

	_, err = stmt.Exec()
	if err != nil {
		logger.Errorf("Error in creating results table: %s", err)
	}

	// add results table for parcelNIR
	stmt, err = db.Prepare(`create table if not exists parcelNIR
									(
										parcelID integer,
										nrd text,
										dt text,
										nir real
									);`)
	if err != nil {
		logger.Errorf("Error in statement of parcel nir table: %s", err)
	}

	_, err = stmt.Exec()
	if err != nil {
		logger.Errorf("Error in creating parcel nir table: %s", err)
	}

}
