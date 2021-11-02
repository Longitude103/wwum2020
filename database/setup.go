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
	SqliteDB   bool
	SlDb       *sqlx.DB
	SYear      int
	EYear      int
	Logger     *zap.SugaredLogger
	PNirDB     *DB
	RchDb      *RchDB
	AppDebug   bool
	ExcessFlow bool
	Desc       string
}

type Option func(*Setup)

// NewSetup is an initialization function for the Setup struct that sets the initial database connections, logger, and stores
// the flags for excess flow and debug.
func NewSetup(myEnv map[string]string, options ...Option) (*Setup, error) {
	s := &Setup{EYear: 1953, SYear: 2020, SqliteDB: true}
	for _, option := range options {
		option(s)
	}

	if s.AppDebug {
		s.Logger.Info("Debug is Set, limited records retrieved for speed.")
	}

	if s.ExcessFlow {
		s.Logger.Info("Using Excess Flows")
	}

	if s.SqliteDB {
		var err error
		s.SlDb, err = GetSqlite(s.Logger, s.Desc)
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

	var err error
	s.PgDb, err = PgConnx(myEnv)
	if err != nil {
		return s, err
	}

	return s, nil
}

func WithLogger() Option {
	return func(s *Setup) {
		l, _ := NewLogger()
		s.Logger = l.Sugar()
		s.Logger.Infow("Setting Up Results database, getting postgres DB Connection.")
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

func WithDescription(textDesc string) Option {
	return func(s *Setup) { s.Desc = textDesc }
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
