package db

import (
	"time"

	"github.com/rs/zerolog/log"
)

type (

	dbi int

	CalendarProvider interface{
		Connection() db
	}

	//ClientOptions pass parameters for factory
	ClientOptions struct {
		ConnectionString string
		Collection string
		Database string
		Provider dbi
		Timeout time.Duration
	}



	mongoImpl struct {
		Options ClientOptions
	}

	sqlImpl struct {
		Options ClientOptions
	}
)

const (
	MongoProvider dbi = iota
	MssqlProvider
)

func NewDb(options ClientOptions) CalendarProvider {
	switch options.Provider {
	case MongoProvider:
		return &mongoImpl{Options: options}
	case MssqlProvider:
		return &sqlImpl{Options: options}
	default:
		log.Fatal().Msgf("not implemented %v",options.Provider)
		return nil
	}
}

func (c *sqlImpl) Connection() db {
	panic("implement me")
}
