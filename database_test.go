package main

import (
	"encoding/json"
	"errors"
	"strings"
	"testing"
)

func TestDatabaseLoading(t *testing.T) {
	t.Run("returns an error if database does not exists on the provided path", func(t *testing.T) {
		_, err := loadDB("unvalid path")
		if err == nil {
			t.Error("does not return an error when provided with an existing json file")
		}
	})

	t.Run("returns an error if the file is not a json file", func(t *testing.T) {
		_, err := loadDB("./db.txt")
		if err == nil {
			t.Fatal("does not return an error when provided with a non json file")
		}
		if !errors.Is(err, errNotAJsonFile) {
			t.Error("does not return the proper error")
		}
	})

	t.Run("returns an error if database does not exists on the provided path", func(t *testing.T) {
		_, err := loadDB("./db.json")
		if err != nil {
			t.Error("should not return an error when provided with an existing json file")
		}
	})
}

func TestDatabaseReadAndWrite(t *testing.T) {
	db := SetupTestDb(t)
	t.Run("returns an error when a table does not exist", func(t *testing.T) {
		_, err := db.getTable("users")
		if err == nil {
			t.Fatal("should've returned an error when the table does not exist")
		}

		if !errors.Is(err, errTableNotFound) {
			t.Errorf("it should've returned errTableNotFound, instead got %s", err)
		}
	})
}

func SetupTestDb(t testing.TB) *database {
	t.Helper()
	content := `
{
	"posts": [
		{"id": 1, "content": "post 1"},
		{"id": 2, "content": "post 2"},
		{"id": 3, "content": "post 3"},
		{"id": 4, "content": "post 4"},
		{"id": 5, "content": "post 5"},
		{"id": 6, "content": "post 6"}
	]
}`

	database := &database{
		Tables: make(map[table][]row),
	}
	if err := json.NewDecoder(strings.NewReader(content)).Decode(&database.Tables); err != nil {
		t.Fatal("could not decode the testing database")
	}
	return database
}
