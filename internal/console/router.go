package console

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"lims/internal/audit"
	"lims/internal/instrument"
	"lims/internal/ns"
	"lims/internal/quota"
	"lims/internal/report"
	"lims/internal/result"
	"lims/internal/sample"
	"lims/internal/task"
)

// Deps wires every store into the console handlers.
type Deps struct {
	Namespaces  ns.Store
	Samples     sample.Store
	Tasks       task.Store
	Instruments instrument.Store
	Results     result.Store
	Reports     report.Store
	Audits      audit.Store
	Quotas      quota.LimitStore
	ReportFiles report.FileStore
}

// NewRouter builds the console HTTP router with embedded pages and JSON API.
func NewRouter(d Deps) http.Handler {
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Logger)

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/samples", http.StatusFound)
	})
	router.Get("/samples", pageHandler("samples.html"))
	router.Get("/tasks", pageHandler("tasks.html"))
	router.Get("/reports", pageHandler("reports.html"))
	router.Get("/audit", pageHandler("audit.html"))

	router.Route("/api", func(api chi.Router) {
		api.Get("/namespaces", listNamespaces(d))
		api.Post("/namespaces", createNamespace(d))
		api.Post("/namespaces/{id}/disable", disableNamespace(d))

		api.Get("/samples", listSamples(d))
		api.Post("/samples", registerSample(d))
		api.Get("/samples/{id}", getSample(d))
		api.Post("/samples/{id}/forward", forwardSample(d))

		api.Get("/tasks", listTasks(d))
		api.Post("/tasks", createTask(d))
		api.Post("/tasks/{id}/assign", assignTask(d))
		api.Post("/tasks/{id}/complete", completeTask(d))

		api.Get("/instruments", listInstruments(d))
		api.Post("/instruments", createInstrument(d))
		api.Post("/instruments/{id}/reconnect", reconnectInstrument(d))
		api.Post("/instruments/{id}/disconnect", disconnectInstrument(d))
		api.Post("/instruments/{id}/retry", retryInstrument(d))

		api.Get("/results", listResults(d))
		api.Post("/results", createResult(d))
		api.Get("/results/{id}", getResult(d))
		api.Post("/results/{id}/qc", runQC(d))

		api.Get("/reports", listReports(d))
		api.Post("/reports", createReport(d))
		api.Post("/reports/{id}/sign", signReport(d))
		api.Post("/reports/{id}/revise", reviseReport(d))
		api.Post("/reports/{id}/archive", archiveReport(d))
		api.Get("/reports/{id}/file", getReportFile(d))

		api.Get("/audit", listAudit(d))
		api.Get("/quota", listQuota(d))
		api.Put("/quota/{instrumentID}", updateQuota(d))
	})
	return router
}
