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
	AppDebug   bool
	ExcessFlow bool
}

func (s *Setup) NewSetup(debug, ef bool) error {
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

	s.SlDb, err = GetSqlite(s.Logger)
	if err != nil {
		return err
	}

	s.PgDb, err = PgConnx()
	if err != nil {
		return err
	}

	s.PNirDB, err = PNirDB(s.SlDb)
	if err != nil {
		return err
	}

	return nil
}

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

func NewLogger() (*zap.Logger, error) {
	cfg := zap.NewProductionConfig()
	path := fmt.Sprintf("./results%s.log", time.Now().Format(time.RFC3339))

	cfg.OutputPaths = []string{
		path,
	}

	return cfg.Build()
}
