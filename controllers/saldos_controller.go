package controllers

import (
	"database/sql"
	"encoding/json"
	"go_api/db"
	"go_api/models"
	"go_api/views"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/tealeg/xlsx"
)

// getSaldos obtiene la lista de saldos desde la base de datos MySQL.
// Se le pasa la conexión a la BD (esto permite reutilizar la función con otras conexiones si es necesario).
func getSaldos(dbConn *sql.DB) ([]models.Saldo, error) {
	query := `SELECT 
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
ORDER BY ANIO_PRO, COD_ART;`

	rows, err := dbConn.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var saldos []models.Saldo
	for rows.Next() {
		var s models.Saldo
		var fechaIngresoBytes []byte // variable temporal para Fecha_Ingreso
		var dias sql.NullInt64       // variable temporal para Dias_Desde_Ingreso
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
			&dias, // se lee como sql.NullInt64
			&s.SaldoFinEnero,
			&s.SaldoFinFebrero,
			&s.SaldoFinMarzo,
			&s.SaldoFinAbril,
			&s.SaldoFinMayo,
			&s.SaldoFinJunio,
			&s.SaldoFinJulio,
			&s.SaldoFinAgosto,
			&s.SaldoFinSeptiembre,
			&s.SaldoFinOctubre,
			&s.SaldoFinNoviembre,
			&s.SaldoFinDiciembre,
		)
		if err != nil {
			return nil, err
		}
		// Convertir fecha de []byte a time.Time si no está vacía
		fechaStr := string(fechaIngresoBytes)
		if fechaStr == "" {
			s.FechaIngreso = time.Time{} // valor por defecto si está vacío
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
		} else {
			s.DiasDesdeIngreso = 0
		}
		saldos = append(saldos, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return saldos, nil
}

// Nueva función para paginación: obtiene 'limit' registros con 'offset'
func getSaldosPaginated(dbConn *sql.DB, offset, limit int, search, sortField, sortDir string) ([]models.Saldo, int, error) {
	// Construir la consulta base
	baseQuery := `SELECT 
        COD_ART AS Codigo_Producto,
        ZET_ART AS Zeta,
        ANIO_PRO AS Año_Produccion,
        DES_INT AS Nombre_Producto,
        UNI_CAJ AS Unidad_Caja,
        CIF_UNI AS Costo_CIF,
        cos_uni AS Costo_Real,
        FEC_ING AS Fecha_Ingreso,
        CAN_ING AS Cantidad_Ingresada,
        SAL_ANT AS Saldo_Anterior,
        DATEDIFF(CURDATE(), FEC_ING) AS Dias_Desde_Ingreso,
        FIN_ENE AS Saldo_Fin_Enero,
        FIN_FEB AS Saldo_Fin_Febrero,
        FIN_MAR AS Saldo_Fin_Marzo,
        FIN_ABR AS Saldo_Fin_Abril,
        FIN_MAY AS Saldo_Fin_Mayo,
        FIN_JUN AS Saldo_Fin_Junio,
        FIN_JUL AS Saldo_Fin_Julio,
        FIN_AGO AS Saldo_Fin_Agosto,
        FIN_SEP AS Saldo_Fin_Septiembre,
        FIN_OCT AS Saldo_Fin_Octubre,
        FIN_NOV AS Saldo_Fin_Noviembre,
        FIN_DIC AS Saldo_Fin_Diciembre
    FROM saldos`

	// Agregar condición WHERE si hay búsqueda
	whereClause := ""
	if search != "" {
		whereClause = " WHERE DES_INT LIKE ? OR COD_ART LIKE ? OR ZET_ART LIKE ?"
	}

	// Agregar ORDER BY si hay campo de ordenamiento
	orderClause := " ORDER BY "
	if sortField != "" {
		// Mapear nombres de campos
		switch sortField {
		case "CodigoProducto":
			orderClause += "COD_ART"
		case "Zeta":
			orderClause += "ZET_ART"
		case "AnioProduccion":
			orderClause += "ANIO_PRO"
		case "NombreProducto":
			orderClause += "DES_INT"
		default:
			orderClause += "COD_ART"
		}
		orderClause += " " + sortDir
	} else {
		orderClause += "COD_ART ASC"
	}

	// Construir consulta final con LIMIT y OFFSET
	query := baseQuery + whereClause + orderClause + " LIMIT ? OFFSET ?"

	log.Printf("Query ejecutada: %s", query)

	// Ejecutar consulta
	var rows *sql.Rows
	var err error
	if search != "" {
		searchPattern := "%" + search + "%"
		rows, err = dbConn.Query(query, searchPattern, searchPattern, searchPattern, limit, offset)
	} else {
		rows, err = dbConn.Query(query, limit, offset)
	}

	if err != nil {
		log.Printf("Error en la consulta: %v", err)
		return nil, 0, err
	}
	defer rows.Close()

	// Procesar resultados
	var saldos []models.Saldo
	for rows.Next() {
		var s models.Saldo
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
			&s.SaldoFinEnero,
			&s.SaldoFinFebrero,
			&s.SaldoFinMarzo,
			&s.SaldoFinAbril,
			&s.SaldoFinMayo,
			&s.SaldoFinJunio,
			&s.SaldoFinJulio,
			&s.SaldoFinAgosto,
			&s.SaldoFinSeptiembre,
			&s.SaldoFinOctubre,
			&s.SaldoFinNoviembre,
			&s.SaldoFinDiciembre,
		)
		if err != nil {
			log.Printf("Error al escanear fila: %v", err)
			return nil, 0, err
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
					log.Printf("Error al parsear fecha: %v", err)
					return nil, 0, err
				}
			}
		}

		if dias.Valid {
			s.DiasDesdeIngreso = int(dias.Int64)
		}

		saldos = append(saldos, s)
	}

	// Obtener total de registros
	var total int
	countQuery := "SELECT COUNT(*) FROM saldos" + whereClause
	if search != "" {
		searchPattern := "%" + search + "%"
		err = dbConn.QueryRow(countQuery, searchPattern, searchPattern, searchPattern).Scan(&total)
	} else {
		err = dbConn.QueryRow("SELECT COUNT(*) FROM saldos").Scan(&total)
	}
	if err != nil {
		log.Printf("Error al contar registros: %v", err)
		return nil, 0, err
	}

	return saldos, total, nil
}

// Modificación en SaldosHandler para paginación
func SaldosHandler(w http.ResponseWriter, r *http.Request) {
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
	sortField := query.Get("sort")
	sortDir := query.Get("dir")

	if sortDir != "desc" {
		sortDir = "asc"
	}

	// Obtener datos con los filtros aplicados
	offset := (page - 1) * pageSize
	saldos, total, err := getSaldosPaginated(db.MySQLDB, offset, pageSize, search, sortField, sortDir)
	if err != nil {
		http.Error(w, "Error al obtener los datos", http.StatusInternalServerError)
		return
	}

	totalPages := (total + pageSize - 1) / pageSize

	viewData := views.ViewData{
		Items:       saldos,
		CurrentPage: page,
		TotalPages:  totalPages,
		PageSize:    pageSize,
		Search:      search,
		SortField:   sortField,
		SortDir:     sortDir,
	}

	views.RenderSaldos(w, viewData)
}

// Actualización en ExportSaldosHandler para exportar datos paginados.
func ExportSaldosHandler(w http.ResponseWriter, r *http.Request) {
	// Leer parámetro "page" para paginación
	pageStr := r.URL.Query().Get("page")
	page := 1
	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}
	limit := 25
	offset := (page - 1) * limit

	// Obtener la página solicitada con paginación
	saldos, total, err := getSaldosPaginated(db.MySQLDB, offset, limit, "", "", "")
	if err != nil {
		http.Error(w, "Error al obtener los datos", http.StatusInternalServerError)
		log.Println("Error al exportar los saldos:", err)
		return
	}

	log.Printf("Exportando %d registros de un total de %d", len(saldos), total)

	// Crear un nuevo archivo Excel y exportar los datos de la página.
	file := xlsx.NewFile()
	sheet, err := file.AddSheet("Saldos")
	if err != nil {
		http.Error(w, "Error al crear el Excel", http.StatusInternalServerError)
		log.Println("Error al crear hoja en Excel:", err)
		return
	}
	// Agregar encabezado
	row := sheet.AddRow()
	headers := []string{"Código", "Zeta", "Año Producción", "Nombre Producto", "Unidad Caja", "Costo CIF", "Costo Real", "Fecha Ingreso", "Cant. Ingresada", "Saldo Anterior", "Días Desde Ingreso"}
	for _, h := range headers {
		cell := row.AddCell()
		cell.Value = h
	}
	// Agregar datos de la página actual
	for _, s := range saldos {
		row := sheet.AddRow()
		row.AddCell().Value = s.CodigoProducto
		row.AddCell().Value = s.Zeta
		row.AddCell().SetInt(s.AnioProduccion)
		row.AddCell().Value = s.NombreProducto
		row.AddCell().SetFloat(s.UnidadCaja)
		row.AddCell().SetFloat(s.CostoCIF)
		row.AddCell().SetFloat(s.CostoReal)
		row.AddCell().Value = s.FechaIngreso.Format("2006-01-02")
		row.AddCell().SetFloat(s.CantidadIngresada)
		row.AddCell().SetFloat(s.SaldoAnterior)
		row.AddCell().SetInt(s.DiasDesdeIngreso)
	}
	// Enviar archivo Excel como respuesta
	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	w.Header().Set("Content-Disposition", "attachment; filename=saldos.xlsx")
	if err := file.Write(w); err != nil {
		http.Error(w, "Error al generar el Excel", http.StatusInternalServerError)
		log.Println("Error al escribir el Excel:", err)
	}
}

// ApiSaldosHandler maneja la ruta /api/saldos y devuelve los datos en formato JSON.
func ApiSaldosHandler(w http.ResponseWriter, r *http.Request) {
	saldos, err := getSaldos(db.MySQLDB)
	if err != nil {
		http.Error(w, "Error al obtener los datos", http.StatusInternalServerError)
		log.Println("Error al obtener los saldos:", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(saldos)
}
