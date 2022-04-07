package main

import (
	driver "github.com/arangodb/go-driver"
	arangohttp "github.com/arangodb/go-driver/http"
)

func NewArangoDbConnector(url string, username string, password string, dbname string) (*ArangoConnector, error) {
	c := new(ArangoConnector)
	conn, err := arangohttp.NewConnection(arangohttp.ConnectionConfig{
		Endpoints: []string{url},
	})

	if err != nil {
		return nil, err
	}

	client, err := driver.NewClient(driver.ClientConfig{
		Connection:     conn,
		Authentication: driver.BasicAuthentication(username, password),
	})

	if err != nil {
		return nil, err
	}

	db, err := client.Database(nil, dbname)
	if err != nil {
		return nil, err
	}
	c.Db = db
	return c, nil
}

type ArangoConnector struct {
	Db driver.Database
}

func (c *ArangoConnector) ExecuteQuery(query string) (driver.Cursor, error) {
	cursor, err := c.Db.Query(nil, query, nil)
	if err != nil {
		return nil, err
	}
	defer cursor.Close()
	return cursor, err
}
