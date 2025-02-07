package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"go_api/db"
	"go_api/models"
	"go_api/views"

	"sort"

	"github.com/tealeg/xlsx"
)

// getStocksFromSQLServer obtiene datos de SQL Server.
func getStocksFromSQLServer(db *sql.DB) ([]models.StockData, error) {
	query := `
        SELECT 
            s.ID_SUCURSAL,
            p.NOMBRE_PRODUCTO,
            p.CODIGO_INTERNO AS Codigo_Producto,
            s.ZETA,
            s.FECHA,
            p.PRECIO_VENTA,
            p.PRECIO_OFERTA,
            s.COSTO_UNITARIO,
            s.ANIO
        FROM STOCKS s
        INNER JOIN PRODUCTO p 
            ON s.ID_PRODUCTO = p.ID_PRODUCTO
        WHERE p.ACTIVO = 1 AND s.ID_SUCURSAL = 211
    `
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stocks []models.StockData
	for rows.Next() {
		var s models.StockData
		if err := rows.Scan(&s.IDSucursal, &s.NombreProducto, &s.CodigoProducto, &s.Zeta,
			&s.Fecha, &s.PrecioVenta, &s.PrecioOferta, &s.CostoUnitario, &s.Anio); err != nil {
			return nil, err
		}
		stocks = append(stocks, s)
	}
	return stocks, nil
}

// getSaldosFromMySQL obtiene saldos desde MySQL modificando el escaneo de Fecha_Ingreso.
func getSaldosFromMySQL(db *sql.DB, year string) ([]models.SaldoData, error) {
	query := `
        SELECT 
            COD_ART AS Codigo_Producto,
            ZET_ART AS Zeta,
            ANIO_PRO AS Año_Produccion,
            DES_INT AS Nombre_Producto, 
            UNI_CAJ AS Unidad_Caja,
            CIF_UNI AS Costo_CIF,
            cos_uni AS Costo_Real,
            MAX(FEC_ING) AS Fecha_Ingreso,    
            SUM(CAN_ING) AS Cantidad_Ingresada, 
            MAX(SAL_ANT) AS Saldo_Anterior,   
            DATEDIFF(CURDATE(), MAX(FEC_ING)) AS Dias_Desde_Ingreso
        FROM saldos 
        WHERE ANIO_PRO = ?  
        GROUP BY COD_ART, ZET_ART, ANIO_PRO, DES_INT, UNI_CAJ, CIF_UNI, cos_uni
        ORDER BY ANIO_PRO, COD_ART
    `

	rows, err := db.Query(query, year)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var saldos []models.SaldoData
	for rows.Next() {
		var s models.SaldoData
		var fechaIngresoBytes []byte
		var dias sql.NullInt64
		err := rows.Scan(
			&s.CodigoProducto,
			&s.Zeta,
			&s.AnioProduccion,
			&s.NombreProducto,
			&s.UnidadCaja,
			&s.CostoCIF,
			&s.CostoReal,
			&fechaIngresoBytes,
			&s.CantidadIngresada,
			&s.SaldoAnterior,
			&dias,
		)
		if err != nil {
			return nil, err
		}

		// Procesar fecha
		fechaStr := string(fechaIngresoBytes)
		if fechaStr == "" {
			s.FechaIngreso = time.Time{}
		} else {
			s.FechaIngreso, err = time.Parse("2006-01-02", fechaStr)
			if err != nil {
				s.FechaIngreso, err = time.Parse("2006-01-02 15:04:05", fechaStr)
				if err != nil {
					return nil, err
				}
			}
		}

		if dias.Valid {
			s.DiasDesdeIngreso = int(dias.Int64)
		}

		saldos = append(saldos, s)
	}

	return saldos, nil
}

// agruparStocksPorZeta agrupa la información de stocks por zeta.
func agruparStocksPorZeta(stocks []models.StockData) map[string]models.StockData {
	agrupados := make(map[string]models.StockData)
	for _, s := range stocks {
		if actual, existe := agrupados[s.Zeta]; !existe || s.Fecha.After(actual.Fecha) {
			agrupados[s.Zeta] = s
		}
	}
	return agrupados
}

// fusionarDatos combina la información de stocks y saldos.
func fusionarDatos(stocksMap map[string]models.StockData, saldos []models.SaldoData) []models.CombinedData {
	var resultados []models.CombinedData
	for _, saldo := range saldos {
		if stock, ok := stocksMap[saldo.Zeta]; ok {
			combinado := models.CombinedData{
				// Datos de SQL Server:
				CodigoProducto: stock.CodigoProducto,
				Zeta:           stock.Zeta,
				AnioProduccion: stock.Anio,
				PrecioVenta:    stock.PrecioVenta,
				PrecioOferta:   stock.PrecioOferta,
				// Datos de MySQL:
				NombreProducto:     saldo.NombreProducto,
				UnidadCaja:         saldo.UnidadCaja,
				CostoCIF:           saldo.CostoCIF,
				CostoReal:          saldo.CostoReal,
				FechaIngreso:       saldo.FechaIngreso,
				CantidadIngresada:  saldo.CantidadIngresada,
				SaldoAnterior:      saldo.SaldoAnterior,
				DiasDesdeIngreso:   saldo.DiasDesdeIngreso,
				SaldoFinEnero:      saldo.SaldoFinEnero,
				SaldoFinFebrero:    saldo.SaldoFinFebrero,
				SaldoFinMarzo:      saldo.SaldoFinMarzo,
				SaldoFinAbril:      saldo.SaldoFinAbril,
				SaldoFinMayo:       saldo.SaldoFinMayo,
				SaldoFinJunio:      saldo.SaldoFinJunio,
				SaldoFinJulio:      saldo.SaldoFinJulio,
				SaldoFinAgosto:     saldo.SaldoFinAgosto,
				SaldoFinSeptiembre: saldo.SaldoFinSeptiembre,
				SaldoFinOctubre:    saldo.SaldoFinOctubre,
				SaldoFinNoviembre:  saldo.SaldoFinNoviembre,
				SaldoFinDiciembre:  saldo.SaldoFinDiciembre,
			}
			resultados = append(resultados, combinado)
		} else {
			log.Printf("No se encontró registro en SQL Server para zeta %s", saldo.Zeta)
		}
	}
	return resultados
}

// CombinedDataHandler utiliza las conexiones inicializadas en db/mysql.go y db/sqlserver.go.
func CombinedDataHandler(w http.ResponseWriter, r *http.Request) {
	// Utilizar conexión global a SQL Server
	if db.SQLServerDB == nil {
		http.Error(w, "Conexión a SQL Server no inicializada", http.StatusInternalServerError)
		log.Println("Conexión a SQL Server no inicializada")
		return
	}
	stocks, err := getStocksFromSQLServer(db.SQLServerDB)
	if err != nil {
		http.Error(w, "Error obteniendo stocks", http.StatusInternalServerError)
		log.Println("Error obteniendo stocks:", err)
		return
	}
	stocksMap := agruparStocksPorZeta(stocks)

	// Utilizar conexión global a MySQL
	if db.MySQLDB == nil {
		http.Error(w, "Conexión a MySQL no inicializada", http.StatusInternalServerError)
		log.Println("Conexión a MySQL no inicializada")
		return
	}
	saldos, err := getSaldosFromMySQL(db.MySQLDB, "2025")
	if err != nil {
		http.Error(w, "Error obteniendo saldos", http.StatusInternalServerError)
		log.Println("Error obteniendo saldos:", err)
		return
	}

	// Fusionar datos y devolver JSON
	resultados := fusionarDatos(stocksMap, saldos)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resultados); err != nil {
		http.Error(w, "Error codificando la respuesta", http.StatusInternalServerError)
		log.Println("Error codificando la respuesta:", err)
	}
}

// CombinedViewHandler ahora envuelve los datos paginados en una estructura con campos para la plantilla.
func CombinedViewHandler(w http.ResponseWriter, r *http.Request) {
	// Obtener parámetros de la URL
	query := r.URL.Query()
	page, _ := strconv.Atoi(query.Get("page"))
	if page < 1 {
		page = 1
	}

	pageSize, _ := strconv.Atoi(query.Get("pageSize"))
	if pageSize < 1 {
		pageSize = 25
	}

	search := query.Get("search")
	missingSearch := query.Get("missingSearch")
	sortField := query.Get("sort")
	sortDir := query.Get("dir")

	if sortDir != "desc" {
		sortDir = "asc"
	}

	// Obtener año de la URL, default 2025
	year := r.URL.Query().Get("year")
	if year == "" {
		year = "2025"
	}

	// Obtener datos...
	stocks, err := getStocksFromSQLServer(db.SQLServerDB)
	if err != nil {
		http.Error(w, "Error obteniendo stocks", http.StatusInternalServerError)
		return
	}

	// Pasar el año a la función getSaldosFromMySQL
	saldos, err := getSaldosFromMySQL(db.MySQLDB, year)
	if err != nil {
		http.Error(w, "Error obteniendo saldos", http.StatusInternalServerError)
		return
	}

	stocksMap := agruparStocksPorZeta(stocks)
	resultados := fusionarDatos(stocksMap, saldos)

	// Filtrar y ordenar resultados
	filteredResults := filterAndSortResults(resultados, search, sortField, sortDir)

	// Calcular paginación
	total := len(filteredResults)
	totalPages := (total + pageSize - 1) / pageSize
	start := (page - 1) * pageSize
	end := start + pageSize
	if end > total {
		end = total
	}

	// Obtener registros faltantes
	var missing []models.SaldoData
	for _, saldo := range saldos {
		if _, ok := stocksMap[saldo.Zeta]; !ok {
			if missingSearch == "" ||
				strings.Contains(strings.ToLower(saldo.Zeta), strings.ToLower(missingSearch)) ||
				strings.Contains(strings.ToLower(saldo.NombreProducto), strings.ToLower(missingSearch)) {
				missing = append(missing, saldo)
			}
		}
	}

	viewData := views.CombinedViewData{
		Data:          filteredResults[start:end],
		Missing:       missing,
		CurrentPage:   page,
		TotalPages:    totalPages,
		PageSize:      pageSize,
		Search:        search,
		MissingSearch: missingSearch,
		SortField:     sortField,
		SortDir:       sortDir,
		Year:          year, // Agregar el año a los datos de la vista
	}

	views.RenderCombined(w, viewData)
}

// Función auxiliar para filtrar y ordenar resultados
func filterAndSortResults(results []models.CombinedData, search, sortField, sortDir string) []models.CombinedData {
	// Filtrar primero
	filtered := make([]models.CombinedData, 0)
	searchLower := strings.ToLower(search)

	for _, item := range results {
		if search == "" ||
			strings.Contains(strings.ToLower(item.CodigoProducto), searchLower) ||
			strings.Contains(strings.ToLower(item.NombreProducto), searchLower) ||
			strings.Contains(strings.ToLower(item.Zeta), searchLower) {
			filtered = append(filtered, item)
		}
	}

	// Ordenar después
	sort.Slice(filtered, func(i, j int) bool {
		asc := sortDir != "desc"
		switch sortField {
		case "CodigoProducto":
			if asc {
				return filtered[i].CodigoProducto < filtered[j].CodigoProducto
			}
			return filtered[i].CodigoProducto > filtered[j].CodigoProducto
		case "Zeta":
			if asc {
				return filtered[i].Zeta < filtered[j].Zeta
			}
			return filtered[i].Zeta > filtered[j].Zeta
		case "AnioProduccion":
			if asc {
				return filtered[i].AnioProduccion < filtered[j].AnioProduccion
			}
			return filtered[i].AnioProduccion > filtered[j].AnioProduccion
		case "NombreProducto":
			if asc {
				return filtered[i].NombreProducto < filtered[j].NombreProducto
			}
			return filtered[i].NombreProducto > filtered[j].NombreProducto
		default:
			return filtered[i].CodigoProducto < filtered[j].CodigoProducto
		}
	})

	return filtered
}

// ExportCombinedHandler exporta los datos fusionados de la página solicitada a Excel.
func ExportCombinedHandler(w http.ResponseWriter, r *http.Request) {
	// Obtener los parámetros de filtrado de la URL
	year := r.URL.Query().Get("year")
	if year == "" {
		year = "2025"
	}
	search := r.URL.Query().Get("search")
	sortField := r.URL.Query().Get("sort")
	sortDir := r.URL.Query().Get("dir")

	// Obtener datos
	stocks, err := getStocksFromSQLServer(db.SQLServerDB)
	if err != nil {
		http.Error(w, "Error obteniendo stocks", http.StatusInternalServerError)
		log.Println("Error obteniendo stocks:", err)
		return
	}
	stocksMap := agruparStocksPorZeta(stocks)

	saldos, err := getSaldosFromMySQL(db.MySQLDB, year)
	if err != nil {
		http.Error(w, "Error obteniendo saldos", http.StatusInternalServerError)
		log.Println("Error obteniendo saldos:", err)
		return
	}

	resultados := fusionarDatos(stocksMap, saldos)

	// Aplicar filtros si existen
	if search != "" || sortField != "" {
		resultados = filterAndSortResults(resultados, search, sortField, sortDir)
	}

	// Crear archivo Excel
	file := xlsx.NewFile()
	sheet, err := file.AddSheet("Datos Combinados")
	if err != nil {
		http.Error(w, "Error al crear el Excel", http.StatusInternalServerError)
		return
	}

	// Agregar encabezados
	row := sheet.AddRow()
	headers := []string{
		"Código", "Zeta", "Año Producción", "Precio Venta", "Precio Oferta",
		"Nombre Producto", "Fecha Ingreso", "Costo CIF", "Costo Real",
		"Cant. Ingresada", "Saldo Anterior", "Días Desde Ingreso",
	}
	for _, h := range headers {
		cell := row.AddCell()
		cell.Value = h
	}

	// Agregar datos filtrados
	for _, c := range resultados {
		row := sheet.AddRow()
		row.AddCell().Value = c.CodigoProducto
		row.AddCell().Value = c.Zeta
		row.AddCell().SetInt(c.AnioProduccion)
		row.AddCell().SetFloat(c.PrecioVenta)
		row.AddCell().SetFloat(c.PrecioOferta)
		row.AddCell().Value = c.NombreProducto
		row.AddCell().Value = c.FechaIngreso.Format("2006-01-02")
		row.AddCell().SetFloat(c.CostoCIF)
		row.AddCell().SetFloat(c.CostoReal)
		row.AddCell().SetFloat(c.CantidadIngresada)
		row.AddCell().SetFloat(c.SaldoAnterior)
		row.AddCell().SetInt(c.DiasDesdeIngreso)
	}

	// Generar nombre del archivo con los filtros aplicados
	filename := "datos_combinados"
	if year != "" {
		filename += "_" + year
	}
	if search != "" {
		filename += "_filtrado"
	}
	filename += ".xlsx"

	// Enviar archivo
	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	if err := file.Write(w); err != nil {
		http.Error(w, "Error al generar el Excel", http.StatusInternalServerError)
		log.Println("Error al escribir el Excel:", err)
	}
}
