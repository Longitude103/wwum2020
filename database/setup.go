package database

import (
	"errors"
	"fmt"
	"path/filepath"
	"time"

	"github.com/Longitude103/wwum2020/Utils"
	"github.com/Longitude103/wwum2020/logging"
	"github.com/jmoiron/sqlx"
)

type Setup struct {
	PgDb        *sqlx.DB
	SqliteDB    bool
	SlDb        *sqlx.DB
	SYear       int
	EYear       int
	Logger      *logging.TheLogger
	PNirDB      *DB
	RchDb       *RchDB
	AppDebug    bool
	ExcessFlow  bool
	Post97      bool
	OldGrid     bool
	MF640Grid   bool
	SteadyState bool
}

type Option func(*Setup)

// NewSetup is an initialization function for the Setup struct that sets the initial database connections, logger, and stores
// the flags for excess flow and debug.
func NewSetup(myEnv map[string]string, options ...Option) (*Setup, error) {
	s := &Setup{EYear: 1953, SYear: 2020, SqliteDB: true}
	path, fn, err := Utils.MakeOutputDir()
	if err != nil {
		return s, err
	}

	options = append(options, WithLogger(path, fn))

	for _, option := range options {
		option(s)
	}

	if s.AppDebug {
		s.Logger.Info("Debug is Set, limited records retrieved for speed, no output DB")
	}

	if s.ExcessFlow {
		s.Logger.Info("Using Excess Flows")
	}

	if s.Post97 {
		s.Logger.Info("Is in Post 97 Mode; 97 Land Use for Groundwater Only and MI Pumping is held constant")
	}

	if s.SqliteDB {
		s.Logger.Info("Setting Up Results database, getting postgres DB Connection.")
		var err error
		s.SlDb, err = GetSqlite(s.Logger, path, fn)
		if err != nil {
			return s, err
		}

		s.PNirDB, err = PNirDB(s.SlDb)
		if err != nil {
			return s, err
		}

		s.RchDb, err = ResultsRchDB(s.SlDb)
		if err != nil {
			return s, err
		}
	}

	s.PgDb, err = PgConnx(myEnv)
	if err != nil {
		return s, err
	}

	return s, nil
}

func WithLogger(path string, fileName string) Option {
	return func(s *Setup) {
		fileName := fmt.Sprintf("results%s.log", fileName)
		path := filepath.Join(path, fileName)

		l := logging.NewLogger(path)
		s.Logger = l
	}
}

func WithDebug() Option {
	return func(s *Setup) { s.AppDebug = true }
}

func WithExcessFlow() Option {
	return func(s *Setup) { s.ExcessFlow = true }
}

func WithNoSQLite() Option {
	return func(s *Setup) { s.SqliteDB = false }
}

func WithPost97() Option {
	return func(s *Setup) { s.Post97 = true }
}

func WithOldGrid() Option {
	return func(s *Setup) { s.OldGrid = true }
}

func WithMF640Grid() Option {
	return func(s *Setup) { s.MF640Grid = true }
}

// WithSteadyState sets the Steady State Bool to true and also sets the years to SYear to 1893 and EYear to 1952
func WithSteadyState(startYr, endYr int) Option {
	return func(s *Setup) {
		s.SteadyState = true
		s.SYear = startYr
		s.EYear = endYr
	}

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

func (s *Setup) CellType() int {
	ct := 2

	if s.OldGrid {
		ct = 1
	}

	if s.MF640Grid {
		ct = 3
	}

	return ct
}
