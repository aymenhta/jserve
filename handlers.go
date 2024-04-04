package main

import (
	"net/http"
	"slices"
	"strconv"
	"strings"
)

func (app *application) listDBHandler(w http.ResponseWriter, r *http.Request) {
	app.writeJSON(w, http.StatusOK, app.db, nil)
}

func (app *application) showTableHandler(w http.ResponseWriter, r *http.Request) {
	table := table(r.PathValue("table"))
	// get the table
	data, err := app.db.getTable(table)
	if err != nil {
		app.writeJSON(
			w,
			http.StatusNotFound,
			envelope{"message": err.Error()}, nil)
		return
	}

	// *-----------------------------------------*
	// searching if available
	col := r.URL.Query().Get("col")
	val := r.URL.Query().Get("val")

	if col != "" && val != "" {
		if f, err := strconv.ParseFloat(val, 64); err == nil {
			data, err = searchRecords(data, col, f)
			if err != nil {
				app.writeJSON(
					w,
					http.StatusNotFound,
					envelope{"message": err.Error(), "name": table}, nil)
				return
			}
		} else if val == "true" || val == "false" {
			b, err := strconv.ParseBool(val)
			if err != nil {
				app.writeJSON(
					w,
					http.StatusInternalServerError,
					envelope{"message": err.Error()}, nil)
				return
			}
			data, err = searchRecords(data, col, b)
			if err != nil {
				app.writeJSON(
					w,
					http.StatusNotFound,
					envelope{"message": err.Error(), "name": table}, nil)
				return
			}
		} else {
			data, err = searchRecords(data, col, val)
			if err != nil {
				app.writeJSON(
					w,
					http.StatusNotFound,
					envelope{"message": err.Error(), "name": table}, nil)
				return
			}
		}
	}
	// *-----------------------------------------*

	// *-----------------------------------------*
	// sorting
	sortQ := r.URL.Query().Get("sort")
	if sortQ == "" {
		sortQ = "id"
	}

	// get sort direction
	sortColumn := sortQ
	descendingOrder := false
	if strings.HasPrefix(sortQ, "-") {
		sortColumn = strings.Split(sortQ, "-")[1]
		descendingOrder = true
	}

	// check if the sorting column exists
	for _, row := range data {
		// check if key exist in row
		_, ok := row[sortColumn]
		if !ok {
			app.writeJSON(
				w,
				http.StatusNotFound,
				envelope{"message": errTableNotFound.Error(), "name": table}, nil)
			return
		}
	}

	// sort
	result := slices.Clone(data) // so it doesn't effect the original data

	if _, ok := data[0][sortColumn].(int); ok {
		quickSort[int](result, sortColumn, descendingOrder)
	} else if _, ok := data[0][sortColumn].(float64); ok {
		quickSort[float64](result, sortColumn, descendingOrder)
	} else if _, ok := data[0][sortColumn].(string); ok {
		quickSort[string](result, sortColumn, descendingOrder)
	} else {
		app.writeJSON(
			w,
			http.StatusInternalServerError,
			envelope{"message": "column is not sortable"}, nil)
		return
	}
	// *-----------------------------------------*

	app.writeJSON(w, http.StatusOK, envelope{string(table): result}, nil)
}

func (app *application) addRecordHandler(w http.ResponseWriter, r *http.Request) {
	table := table(r.PathValue("table"))

	if !app.db.tableExists(table) {
		app.writeJSON(
			w,
			http.StatusNotFound,
			envelope{"message": errTableNotFound, "name": table}, nil)
		return
	}

	newRow := make(row)
	if err := app.readJSON(w, r, &newRow); err != nil {
		app.writeJSON(
			w,
			http.StatusInternalServerError,
			envelope{"message": err.Error()}, nil)
		return
	}

	app.db.Tables[table] = append(app.db.Tables[table], newRow)
	l := len(app.db.Tables[table])
	app.db.Tables[table][l-1]["id"] = app.db.Tables[table][l-2]["id"].(float64) + 1
	app.writeJSON(w, http.StatusOK, envelope{string(table): app.db.Tables[table][l-1]}, nil)
}

func (app *application) editRecordHandler(w http.ResponseWriter, r *http.Request) {
	table := table(r.PathValue("table"))
	pathID := r.PathValue("id")

	id, err := strconv.ParseFloat(pathID, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	newRow := make(row)
	if err := app.readJSON(w, r, &newRow); err != nil {
		app.writeJSON(
			w,
			http.StatusInternalServerError,
			envelope{"message": err.Error()}, nil)
		return
	}

	response, err := app.db.EditRowById(table, id, newRow)
	if err != nil {
		app.writeJSON(
			w,
			http.StatusNotFound,
			envelope{"message": err.Error()}, nil)
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{string(table): response}, nil)
}

func (app *application) deleteRecordHandler(w http.ResponseWriter, r *http.Request) {
	table := table(r.PathValue("table"))
	pathID := r.PathValue("id")

	id, err := strconv.ParseFloat(pathID, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err = app.db.DeleteRowById(table, id); err != nil {
		app.writeJSON(
			w,
			http.StatusNotFound,
			envelope{"message": err.Error()}, nil)
		return
	}
	app.writeJSON(w, http.StatusOK, envelope{"message": "row was deleted"}, nil)
}

func (app *application) showRecordHandler(w http.ResponseWriter, r *http.Request) {
	table := table(r.PathValue("table"))
	pathID := r.PathValue("id")

	id, err := strconv.ParseFloat(pathID, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	row, err := app.db.GetRowById(table, id)
	if err != nil {
		app.writeJSON(
			w,
			http.StatusNotFound,
			envelope{"message": err.Error()}, nil)
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{string(table): row}, nil)
}
