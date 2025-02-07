package models

import "time"

// Datos extraídos de SQL Server (consulta 1)
type StockData struct {
	IDSucursal     int
	NombreProducto string
	CodigoProducto string
	Zeta           string
	Fecha          time.Time
	PrecioVenta    float64
	PrecioOferta   float64
	CostoUnitario  float64
	Anio           int
}

// Datos extraídos de MySQL (consulta 2)
type SaldoData struct {
	CodigoProducto     string // tal vez no lo uses, ya que se hará la fusión por Zeta
	Zeta               string
	AnioProduccion     int
	NombreProducto     string
	UnidadCaja         float64
	CostoCIF           float64
	CostoReal          float64
	FechaIngreso       time.Time
	CantidadIngresada  float64
	SaldoAnterior      float64
	DiasDesdeIngreso   int
	SaldoFinEnero      float64
	SaldoFinFebrero    float64
	SaldoFinMarzo      float64
	SaldoFinAbril      float64
	SaldoFinMayo       float64
	SaldoFinJunio      float64
	SaldoFinJulio      float64
	SaldoFinAgosto     float64
	SaldoFinSeptiembre float64
	SaldoFinOctubre    float64
	SaldoFinNoviembre  float64
	SaldoFinDiciembre  float64
}

// Estructura combinada final (puedes agregar o quitar campos según tus necesidades)
type CombinedData struct {
	// Datos provenientes de SQL Server:
	CodigoProducto string
	Zeta           string
	AnioProduccion int // o Anio, según corresponda
	PrecioVenta    float64
	PrecioOferta   float64

	// Datos provenientes de MySQL:
	NombreProducto     string
	UnidadCaja         float64
	CostoCIF           float64
	CostoReal          float64
	FechaIngreso       time.Time
	CantidadIngresada  float64
	SaldoAnterior      float64
	DiasDesdeIngreso   int
	SaldoFinEnero      float64
	SaldoFinFebrero    float64
	SaldoFinMarzo      float64
	SaldoFinAbril      float64
	SaldoFinMayo       float64
	SaldoFinJunio      float64
	SaldoFinJulio      float64
	SaldoFinAgosto     float64
	SaldoFinSeptiembre float64
	SaldoFinOctubre    float64
	SaldoFinNoviembre  float64
	SaldoFinDiciembre  float64
}
