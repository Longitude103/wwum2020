package database

import (
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"time"
)

type Setup struct {
	PgDb       *sqlx.DB
	SlDb       *sqlx.DB
	SYear      int
	EYear      int
	Logger     *zap.SugaredLogger
	PNirDB     *DB
	RchDb      *RchDB
	AppDebug   bool
	ExcessFlow bool
}

// NewSetup is an initialization function for the Setup struct that sets the initial database connections, logger, and stores
// the flags for excess flow and debug.
func (s *Setup) NewSetup(debug, ef bool, myEnv map[string]string, noSqlite bool, mDesc string) error {
	l, err := NewLogger()
	if err != nil {
		return err
	}

	s.Logger = l.Sugar()
	s.Logger.Infow("Setting Up Results database, getting postgres DB Connection.")

	if debug {
		s.AppDebug = debug
		s.Logger.Info("Debug is Set, limited records retrieved for speed.")
	}

	if ef {
		s.ExcessFlow = ef
		s.Logger.Info("Using Excess Flows")
	}

	if !noSqlite {
		s.SlDb, err = GetSqlite(s.Logger, mDesc)
		if err != nil {
			return err
		}

		s.PNirDB, err = PNirDB(s.SlDb)
		if err != nil {
			return err
		}

		s.RchDb, err = ResultsRchDB(s.SlDb)
		if err != nil {
			return err
		}
	}

	s.PgDb, err = PgConnx(myEnv)
	if err != nil {
		return err
	}

	return nil
}

// SetYears is an initializer method for the Setup struct to set the start and end years of the application run.
func (s *Setup) SetYears(sYear, eYear int) error {
	if sYear > 1953 || sYear < time.Now().Year() {
		s.SYear = sYear
	} else {
		return errors.New("start year out of range")
	}

	if eYear > 1953 && eYear < time.Now().Year() {
		s.EYear = eYear
	} else {
		return errors.New("end year out of range")
	}

	return nil
}

// NewLogger is a function to setup the new zap.logger and set the path and file name.
func NewLogger() (*zap.Logger, error) {
	cfg := zap.NewProductionConfig()
	path := fmt.Sprintf("./results%s.log", time.Now().Format(time.RFC3339))

	cfg.OutputPaths = []string{
		path,
	}

	return cfg.Build()
}
