package views

import (
	"go_api/models" // Agregamos esta importación
	"html/template"
	"net/http"
	"time"
)

var combinedTemplate = `
{{define "title"}}Datos Combinados{{end}}

{{define "content"}}
    <div class="container mx-auto">
        <h1 class="text-3xl font-bold mb-6">Datos Combinados</h1>
        
        <div class="mb-4 flex justify-between items-center">
            <div class="flex items-center">
                <form method="GET" class="flex gap-4">
                    <input 
                        type="text" 
                        name="search"
                        value="{{.Search}}"
                        placeholder="Buscar..."
                        class="px-4 py-2 border rounded-lg">
                    
                    <select name="year" class="ml-4 px-4 py-2 border rounded-lg">
                        <option value="2024" {{if eq .Year "2024"}}selected{{end}}>2024</option>
                        <option value="2025" {{if eq .Year "2025"}}selected{{end}}>2025</option>
                        <option value="2026" {{if eq .Year "2026"}}selected{{end}}>2026</option>
                    </select>

                    <select name="pageSize" class="ml-4 px-4 py-2 border rounded-lg">
                        <option value="10" {{if eq .PageSize 10}}selected{{end}}>10 por página</option>
                        <option value="25" {{if eq .PageSize 25}}selected{{end}}>25 por página</option>
                        <option value="50" {{if eq .PageSize 50}}selected{{end}}>50 por página</option>
                    </select>
                    
                    <button type="submit" class="bg-blue-500 text-white px-4 py-2 rounded">
                        Filtrar
                    </button>
                </form>
            </div>
            
            <a href="/exportCombined?year={{.Year}}{{if .Search}}&search={{.Search}}{{end}}{{if .SortField}}&sort={{.SortField}}&dir={{.SortDir}}{{end}}" 
               class="bg-green-500 hover:bg-green-700 text-white font-bold py-2 px-4 rounded">
                Exportar Excel
            </a>
        </div>

        <div class="overflow-x-auto bg-white rounded-lg shadow">
            <table class="min-w-full">
                <thead class="bg-gray-800 text-white">
                    <tr>
                        <th class="px-4 py-2"><a href="?sort=CodigoProducto&dir={{.NextSort "CodigoProducto"}}&search={{.Search}}" class="text-white">Código {{.SortIndicator "CodigoProducto"}}</a></th>
                        <th class="px-4 py-2"><a href="?sort=Zeta&dir={{.NextSort "Zeta"}}&search={{.Search}}" class="text-white">Zeta {{.SortIndicator "Zeta"}}</a></th>
                        <th class="px-4 py-2"><a href="?sort=AnioProduccion&dir={{.NextSort "AnioProduccion"}}&search={{.Search}}" class="text-white">Año {{.SortIndicator "AnioProduccion"}}</a></th>
                        <th class="px-4 py-2"><a href="?sort=PrecioVenta&dir={{.NextSort "PrecioVenta"}}&search={{.Search}}" class="text-white">Precio Venta {{.SortIndicator "PrecioVenta"}}</a></th>
                        <th class="px-4 py-2">Precio Oferta</th>
                        <th class="px-4 py-2"><a href="?sort=NombreProducto&dir={{.NextSort "NombreProducto"}}&search={{.Search}}" class="text-white">Nombre {{.SortIndicator "NombreProducto"}}</a></th>
                        <th class="px-4 py-2">Fecha Ingreso</th>
                        <th class="px-4 py-2">CIF</th>
                        <th class="px-4 py-2">Real</th>
                        <th class="px-4 py-2">Cant.</th>
                        <th class="px-4 py-2">Saldo</th>
                        <th class="px-4 py-2">Días</th>
                    </tr>
                </thead>
                <tbody class="text-gray-700">
                    {{range .Data}}
                    <tr class="hover:bg-gray-50">
                        <td class="border px-4 py-2">{{.CodigoProducto}}</td>
                        <td class="border px-4 py-2">{{.Zeta}}</td>
                        <td class="border px-4 py-2">{{.AnioProduccion}}</td>
                        <td class="border px-4 py-2">{{.PrecioVenta}}</td>
                        <td class="border px-4 py-2">{{.PrecioOferta}}</td>
                        <td class="border px-4 py-2">{{.NombreProducto}}</td>
                        <td class="border px-4 py-2">{{formatDate .FechaIngreso}}</td>
                        <td class="border px-4 py-2">{{.CostoCIF}}</td>
                        <td class="border px-4 py-2">{{.CostoReal}}</td>
                        <td class="border px-4 py-2">{{.CantidadIngresada}}</td>
                        <td class="border px-4 py-2">{{.SaldoAnterior}}</td>
                        <td class="border px-4 py-2">{{.DiasDesdeIngreso}}</td>
                    </tr>
                    {{end}}
                </tbody>
            </table>
        </div>

        <div class="mt-4 flex items-center justify-between">
            {{template "pagination" .}}
        </div>

        {{if .Missing}}
        <div class="mt-8">
            <h2 class="text-2xl font-bold mb-4">Registros sin correspondencia en SQL Server</h2>
            <form method="GET" action="" class="mb-4">
                <input 
                    type="text" 
                    name="missingSearch"
                    value="{{.MissingSearch}}"
                    placeholder="Buscar en faltantes..."
                    class="px-4 py-2 border rounded-lg">
            </form>
            
            <table class="min-w-full bg-white">
                <thead class="bg-gray-800 text-white">
                    <tr>
                        <th class="px-4 py-2">Código</th>
                        <th class="px-4 py-2">Zeta</th>
                        <th class="px-4 py-2">Nombre Producto</th>
                        <th class="px-4 py-2">Fecha Ingreso</th>
                        <th class="px-4 py-2">Saldo</th>
                    </tr>
                </thead>
                <tbody class="text-gray-700">
                    {{range .Missing}}
                    <tr class="hover:bg-gray-50">
                        <td class="border px-4 py-2">{{.CodigoProducto}}</td>
                        <td class="border px-4 py-2">{{.Zeta}}</td>
                        <td class="border px-4 py-2">{{.NombreProducto}}</td>
                        <td class="border px-4 py-2">{{formatDate .FechaIngreso}}</td>
                        <td class="border px-4 py-2">{{.SaldoAnterior}}</td>
                    </tr>
                    {{end}}
                </tbody>
            </table>
        </div>
        {{end}}
    </div>
{{end}}

{{define "pagination"}}
    {{if gt .CurrentPage 1}}
    <a href="?page={{dec .CurrentPage}}&pageSize={{.PageSize}}{{if .Search}}&search={{.Search}}{{end}}{{if .SortField}}&sort={{.SortField}}&dir={{.SortDir}}{{end}}" 
       class="px-4 py-2 bg-gray-300 rounded">
        Anterior
    </a>
    {{else}}
    <span class="px-4 py-2 bg-gray-300 rounded opacity-50">Anterior</span>
    {{end}}
    
    <span>
        Página {{.CurrentPage}} de {{.TotalPages}}
    </span>
    
    {{if lt .CurrentPage .TotalPages}}
    <a href="?page={{inc .CurrentPage}}&pageSize={{.PageSize}}{{if .Search}}&search={{.Search}}{{end}}{{if .SortField}}&sort={{.SortField}}&dir={{.SortDir}}{{end}}" 
       class="px-4 py-2 bg-gray-300 rounded">
        Siguiente
    </a>
    {{else}}
    <span class="px-4 py-2 bg-gray-300 rounded opacity-50">Siguiente</span>
    {{end}}
{{end}}
`

type CombinedViewData struct {
	Data          []models.CombinedData
	Missing       []models.SaldoData
	CurrentPage   int
	TotalPages    int
	PageSize      int
	Search        string
	MissingSearch string
	SortField     string
	SortDir       string
	Year          string
}

func (d CombinedViewData) SortIndicator(field string) string {
	if d.SortField == field {
		if d.SortDir == "asc" {
			return "↑"
		}
		return "↓"
	}
	return ""
}

func (d CombinedViewData) NextSort(field string) string {
	if d.SortField == field && d.SortDir == "asc" {
		return "desc"
	}
	return "asc"
}

func RenderCombined(w http.ResponseWriter, data CombinedViewData) {
	funcMap := template.FuncMap{
		"formatDate": func(t time.Time) string {
			return t.Format("2006-01-02")
		},
		"inc": func(i int) int { return i + 1 },
		"dec": func(i int) int {
			if i > 1 {
				return i - 1
			}
			return 1
		},
	}

	tmpl := template.New("layout.tmpl").Funcs(funcMap)
	tmpl, err := tmpl.ParseFiles("c:/Users/pc/Herd/go_api/views/layout.tmpl")
	if err != nil {
		http.Error(w, "Error al cargar el layout", http.StatusInternalServerError)
		return
	}

	if _, err = tmpl.Parse(combinedTemplate); err != nil {
		http.Error(w, "Error al cargar la plantilla", http.StatusInternalServerError)
		return
	}

	if err = tmpl.ExecuteTemplate(w, "layout.tmpl", data); err != nil {
		http.Error(w, "Error al renderizar la plantilla", http.StatusInternalServerError)
	}
}
