package routes

import (
	"Forum/handlers"
	"github.com/gorilla/mux"
)

func ReportRoutes(r *mux.Router) {
	r.HandleFunc("/send_report", handlers.SendReportHandler)
	r.HandleFunc("/delete_report", handlers.DeleteReportHandler)
}
