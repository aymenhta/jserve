package main

import (
	"cmp"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sort"
)

var (
	errTableNotFound  = errors.New("table does not exist")
	errColumnNotFound = errors.New("table does not exist")
	errRecordNotFound = errors.New("record does not exist")
)

type table string
type row map[string]any
type database map[table][]row

func loadDB(path string) (database, error) {
	// read from file
	content, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("could not read database: '%s'", err.Error())
	}

	// decode json
	db := make(database)
	if err = json.NewDecoder(content).Decode(&db); err != nil {
		return nil, fmt.Errorf("could not decode json: '%s'", err.Error())
	}

	return db, nil
}

func (db database) tableExists(name table) bool {
	_, ok := db[name]
	return ok
}

func (db database) getTable(name table) ([]row, error) {
	if !db.tableExists(name) {
		return nil, errTableNotFound
	}

	return db[name], nil
}

func (db database) DeleteRowById(name table, id float64) error {
	if !db.tableExists(name) {
		return errTableNotFound
	}

	for i, row := range db[name] {
		val, ok := row["id"]
		if !ok {
			return errColumnNotFound
		}

		if val == id {
			db[name] = append(db[name][:i], db[name][i+1:]...)
			return nil
		}
	}

	return errRecordNotFound
}

func (db database) EditRowById(name table, id float64, body row) (row, error) {
	if !db.tableExists(name) {
		return nil, errTableNotFound
	}

	for i, row := range db[name] {
		// check if key exist in row
		val, ok := row["id"]
		if !ok {
			return nil, errColumnNotFound
		}

		if val == id {
			db[name][i] = body
			return db[name][i], nil
		}
	}

	return nil, errRecordNotFound
}

func (db database) GetRowById(name table, id float64) (row, error) {
	if !db.tableExists(name) {
		return nil, errTableNotFound
	}

	for i, row := range db[name] {
		// check if key exist in row
		val, ok := row["id"]
		if !ok {
			return nil, errColumnNotFound
		}

		if val == id {
			return db[name][i], nil
		}
	}

	return nil, errRecordNotFound
}

func quickSort[T cmp.Ordered](rows []row, col string, descendingOrder bool) {
	sort.Slice(rows, func(i, j int) bool {
		if descendingOrder {
			return rows[i][col].(T) > rows[j][col].(T)
		}
		return rows[i][col].(T) < rows[j][col].(T)
	})
}
