package controllers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"go_api/db"
	"go_api/models"
	"go_api/views"

	"github.com/tealeg/xlsx"
	"sort"
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
func getSaldosFromMySQL(db *sql.DB) ([]models.SaldoData, error) {
	query := `
        SELECT 
            COD_ART AS Codigo_Producto,
            ZET_ART AS Zeta,
            ANIO_PRO AS Año_Produccion,
            MAX(DES_INT) AS Nombre_Producto,
            MAX(UNI_CAJ) AS Unidad_Caja,
            MAX(CIF_UNI) AS Costo_CIF,
            MAX(cos_uni) AS Costo_Real,
            MAX(FEC_ING) AS Fecha_Ingreso,
            SUM(CAN_ING) AS Cantidad_Ingresada,
            MAX(SAL_ANT) AS Saldo_Anterior,
            DATEDIFF(CURDATE(), MAX(FEC_ING)) AS Dias_Desde_Ingreso,
            MAX(FIN_ENE) AS Saldo_Fin_Enero,
            MAX(FIN_FEB) AS Saldo_Fin_Febrero,
            MAX(FIN_MAR) AS Saldo_Fin_Marzo,
            MAX(FIN_ABR) AS Saldo_Fin_Abril,
            MAX(FIN_MAY) AS Saldo_Fin_Mayo,
            MAX(FIN_JUN) AS Saldo_Fin_Junio,
            MAX(FIN_JUL) AS Saldo_Fin_Julio,
            MAX(FIN_AGO) AS Saldo_Fin_Agosto,
            MAX(FIN_SEP) AS Saldo_Fin_Septiembre,
            MAX(FIN_OCT) AS Saldo_Fin_Octubre,
            MAX(FIN_NOV) AS Saldo_Fin_Noviembre,
            MAX(FIN_DIC) AS Saldo_Fin_Diciembre
        FROM saldos
        GROUP BY COD_ART, ZET_ART, ANIO_PRO
        ORDER BY ANIO_PRO, COD_ART;
    `
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var saldos []models.SaldoData
	for rows.Next() {
		var s models.SaldoData
		var fechaIngresoBytes []byte // variable temporal para Fecha_Ingreso
		var dias sql.NullInt64       // variable temporal para Dias_Desde_Ingreso
		if err := rows.Scan(&s.CodigoProducto, &s.Zeta, &s.AnioProduccion, &s.NombreProducto,
			&s.UnidadCaja, &s.CostoCIF, &s.CostoReal, &fechaIngresoBytes,
			&s.CantidadIngresada, &s.SaldoAnterior, &dias, &s.SaldoFinEnero,
			&s.SaldoFinFebrero, &s.SaldoFinMarzo, &s.SaldoFinAbril, &s.SaldoFinMayo,
			&s.SaldoFinJunio, &s.SaldoFinJulio, &s.SaldoFinAgosto, &s.SaldoFinSeptiembre,
			&s.SaldoFinOctubre, &s.SaldoFinNoviembre, &s.SaldoFinDiciembre); err != nil {
			return nil, err
		}
		// Convertir el []byte recibido a time.Time
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
		// Convertir sql.NullInt64 a int
		if dias.Valid {
			s.DiasDesdeIngreso = int(dias.Int64)
		} else {
			s.DiasDesdeIngreso = 0
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
	saldos, err := getSaldosFromMySQL(db.MySQLDB)
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

	// Obtener datos...
	stocks, err := getStocksFromSQLServer(db.SQLServerDB)
	if err != nil {
		http.Error(w, "Error obteniendo stocks", http.StatusInternalServerError)
		return
	}

	saldos, err := getSaldosFromMySQL(db.MySQLDB)
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
	// Verificar conexiones
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

	if db.MySQLDB == nil {
		http.Error(w, "Conexión a MySQL no inicializada", http.StatusInternalServerError)
		log.Println("Conexión a MySQL no inicializada")
		return
	}
	saldos, err := getSaldosFromMySQL(db.MySQLDB)
	if err != nil {
		http.Error(w, "Error obteniendo saldos", http.StatusInternalServerError)
		log.Println("Error obteniendo saldos:", err)
		return
	}

	resultados := fusionarDatos(stocksMap, saldos)
	// Exportar TODOS los datos sin paginar.
	// Crear archivo Excel y exportar la totalidad de los registros.
	file := xlsx.NewFile()
	sheet, err := file.AddSheet("Datos Combinados")
	if err != nil {
		http.Error(w, "Error al crear el Excel", http.StatusInternalServerError)
		log.Println("Error al crear hoja en Excel:", err)
		return
	}
	// Encabezados con columnas adicionales
	row := sheet.AddRow()
	headers := []string{"Código", "Zeta", "Año Producción", "Precio Venta", "Precio Oferta", "Nombre Producto", "Fecha Ingreso", "Costo CIF", "Costo Real", "Cant. Ingresada", "Saldo Anterior", "Días Desde Ingreso"}
	for _, h := range headers {
		cell := row.AddCell()
		cell.Value = h
	}
	// Agregar todos los registros
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
	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	w.Header().Set("Content-Disposition", "attachment; filename=datos_combinados.xlsx")
	if err := file.Write(w); err != nil {
		http.Error(w, "Error al generar el Excel", http.StatusInternalServerError)
		log.Println("Error al escribir el Excel:", err)
	}
}
