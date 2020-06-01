package client

import (
	"database/sql"

	"github.com/matscus/Hamster/Package/Clients/client/httpsimpl"
	"github.com/matscus/Hamster/Package/Clients/client/postgres"
	"github.com/matscus/Hamster/Package/Clients/subset/pgclient"
)

//PGClient  struct for postgres client
type PGClient struct {
	DB *sql.DB
}

//New funct return client
func New(clientType string, config interface{}) interface{} {
	switch clientType {
	case "postgres":
		c := config.(postgres.Config)
		db, err := sql.Open(c.Driver, "user="+c.User+" password="+c.Password+" dbname="+c.DataBase+" sslmode="+c.SSLMode)
		if err != nil {
			return err
		}
		var client pgclient.PGClient
		client = postgres.PGClient{DB: db}
		return client
	case "https":
		client, err := httpsimpl.NewHTTPSClient()
		if err != nil {
			return nil
		}
		return client
	}
	return nil
}
