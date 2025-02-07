package models

import "time"

// Saldo representa la estructura de cada registro obtenido de la consulta.
type Saldo struct {
	CodigoProducto     string    `json:"Codigo_Producto"`
	Zeta               string    `json:"Zeta"`
	AnioProduccion     int       `json:"Año_Produccion"`
	NombreProducto     string    `json:"Nombre_Producto"`
	UnidadCaja         float64   `json:"Unidad_Caja"` // cambiado de int a float64
	CostoCIF           float64   `json:"Costo_CIF"`
	CostoReal          float64   `json:"Costo_Real"`
	FechaIngreso       time.Time `json:"Fecha_Ingreso"`
	CantidadIngresada  float64   `json:"Cantidad_Ingresada"` // cambiado de int a float64
	SaldoAnterior      float64   `json:"Saldo_Anterior"`     // cambiado de int a float64
	DiasDesdeIngreso   int       `json:"Dias_Desde_Ingreso"`
	SaldoFinEnero      float64   `json:"Saldo_Fin_Enero"`      // cambiado de int a float64
	SaldoFinFebrero    float64   `json:"Saldo_Fin_Febrero"`    // cambiado de int a float64
	SaldoFinMarzo      float64   `json:"Saldo_Fin_Marzo"`      // cambiado de int a float64
	SaldoFinAbril      float64   `json:"Saldo_Fin_Abril"`      // cambiado de int a float64
	SaldoFinMayo       float64   `json:"Saldo_Fin_Mayo"`       // cambiado de int a float64
	SaldoFinJunio      float64   `json:"Saldo_Fin_Junio"`      // cambiado de int a float64
	SaldoFinJulio      float64   `json:"Saldo_Fin_Julio"`      // cambiado de int a float64
	SaldoFinAgosto     float64   `json:"Saldo_Fin_Agosto"`     // cambiado de int a float64
	SaldoFinSeptiembre float64   `json:"Saldo_Fin_Septiembre"` // cambiado de int a float64
	SaldoFinOctubre    float64   `json:"Saldo_Fin_Octubre"`    // cambiado de int a float64
	SaldoFinNoviembre  float64   `json:"Saldo_Fin_Noviembre"`  // cambiado de int a float64
	SaldoFinDiciembre  float64   `json:"Saldo_Fin_Diciembre"`  // cambiado de int a float64
}
