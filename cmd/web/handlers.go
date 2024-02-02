package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/jaked0626/snippetbox/internal/db/models"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}

	snippets, err := app.snippets.List(10)
	if err != nil {
		app.serverError(w, err)
		return
	}

	data := &templateData{
		Snippets: snippets,
	}

	app.render(w, http.StatusOK, "home.html", data)
	return
}

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	// validate input
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		app.badRequest(w)
		return
	} else if id < 1 {
		app.notFound(w)
		return
	}

	// get from db
	s, err := app.snippets.Get(id)
	if errors.Is(err, models.ErrNoRecord) {
		app.notFound(w)
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	}

	data := &templateData{
		Snippet: s,
	}

	app.render(w, http.StatusOK, "view.html", data)

	return
}

func (app *application) snippetList(w http.ResponseWriter, r *http.Request) {
	// validate input
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		app.badRequest(w)
		return
	} else if limit < 1 {
		app.notFound(w)
		return
	}

	// get from db
	snippets, err := app.snippets.List(limit)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// response
	res, err := json.Marshal(snippets)
	if err != nil {
		app.serverError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(res)
	return
}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}
	title := r.URL.Query().Get("title")
	content := r.URL.Query().Get("content")
	expires, err := strconv.Atoi(r.URL.Query().Get("expires"))
	if err != nil {
		app.badRequest(w)
		return
	}

	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/snippet/view?id=%d", id), http.StatusSeeOther)
	return
}
