package handler

import (
	"errors"
	"net/http"

	"forum/internal/entities"
	"forum/pkg/validator"
)

func (app *Application) administrationReportsView(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.render(w, http.StatusBadRequest, Errorpage, nil)
		return
	}

	page := 1
	pageSize := 10

	if p, err := validator.ValidateID(r.PostFormValue("page")); err == nil {
		page = p
	}

	paginationURL := "/administration/reports"
	reportsDTO, err := app.Service.Reaction.GetAllPaginatedPostReportsDTO(page, pageSize, paginationURL)
	if err != nil {
		app.Logger.Error("get post reports", "error", err)
		if errors.Is(err, entities.ErrNoRecord) {
		} else {
			app.render(w, http.StatusInternalServerError, Errorpage, nil)
			return
		}
	}

	data := app.newTemplateData(r)
	data.Reports = reportsDTO.Reports
	data.Header = "All reports from moderators"
	data.Pagination = pagination{
		CurrentPage:      reportsDTO.CurrentPage,
		HasNextPage:      reportsDTO.HasNextPage,
		PaginationAction: reportsDTO.PaginationURL,
	}
	app.render(w, http.StatusOK, "reports.html", data)
}
