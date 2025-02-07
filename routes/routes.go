package routes

import (
	"go_api/controllers"
	"net/http"
)

func SetupRoutes() {
	// Ruta principal
	http.HandleFunc("/", controllers.IndexHandler)
	// Registrar rutas de API y vistas
	http.HandleFunc("/api/saldos", controllers.ApiSaldosHandler)
	http.HandleFunc("/saldos", controllers.SaldosHandler)
	// Ruta para exportar saldos paginados
	http.HandleFunc("/export", controllers.ExportSaldosHandler)
	// Ruta para visualizar datos combinados
	http.HandleFunc("/combined", controllers.CombinedViewHandler)
	// Nueva ruta para exportar datos combinados completos
	http.HandleFunc("/exportCombined", controllers.ExportCombinedHandler)
	// ...agregar más rutas si es necesario...
}
