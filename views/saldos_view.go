package views

import (
	"go_api/models" // Agregamos esta importación
	"html/template"
	"net/http"
	"time"
)

var saldosTemplate = `
{{define "title"}}Tabla de Saldos{{end}}

{{define "content"}}
    <div class="container mx-auto">
        <h1 class="text-3xl font-bold mb-6">Saldos</h1>
        
        <div class="mb-4 flex justify-between items-center">
            <div class="flex items-center">
                <form method="GET" class="flex gap-4">
                    <input 
                        type="text" 
                        name="search"
                        value="{{.Search}}"
                        placeholder="Buscar..."
                        class="px-4 py-2 border rounded-lg">
                    
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
            
            <a href="/export{{if .Search}}?search={{.Search}}{{end}}" 
               class="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded">
                Descargar Excel
            </a>
        </div>

        <div class="overflow-x-auto bg-white rounded-lg shadow">
            <table class="min-w-full">
                <thead class="bg-gray-800 text-white">
                    <tr>
                        <th class="px-4 py-2"><a href="?sort=CodigoProducto&dir={{.NextSort "CodigoProducto"}}&search={{.Search}}" class="text-white">Código {{.SortIndicator "CodigoProducto"}}</a></th>
                        <th class="px-4 py-2"><a href="?sort=Zeta&dir={{.NextSort "Zeta"}}&search={{.Search}}" class="text-white">Zeta {{.SortIndicator "Zeta"}}</a></th>
                        <th class="px-4 py-2"><a href="?sort=AnioProduccion&dir={{.NextSort "AnioProduccion"}}&search={{.Search}}" class="text-white">Año Prod. {{.SortIndicator "AnioProduccion"}}</a></th>
                        <th class="px-4 py-2"><a href="?sort=NombreProducto&dir={{.NextSort "NombreProducto"}}&search={{.Search}}" class="text-white">Nombre {{.SortIndicator "NombreProducto"}}</a></th>
                        <th class="px-4 py-2">Unidad</th>
                        <th class="px-4 py-2">CIF</th>
                        <th class="px-4 py-2">Real</th>
                        <th class="px-4 py-2">Ingreso</th>
                        <th class="px-4 py-2">Cant.</th>
                        <th class="px-4 py-2">Saldo</th>
                        <th class="px-4 py-2">Días</th>
                    </tr>
                </thead>
                <tbody class="text-gray-700">
                    {{range .Items}}
                    <tr class="hover:bg-gray-50">
                        <td class="border px-4 py-2">{{.CodigoProducto}}</td>
                        <td class="border px-4 py-2">{{.Zeta}}</td>
                        <td class="border px-4 py-2">{{.AnioProduccion}}</td>
                        <td class="border px-4 py-2">{{.NombreProducto}}</td>
                        <td class="border px-4 py-2">{{.UnidadCaja}}</td>
                        <td class="border px-4 py-2">{{.CostoCIF}}</td>
                        <td class="border px-4 py-2">{{.CostoReal}}</td>
                        <td class="border px-4 py-2">{{formatDate .FechaIngreso}}</td>
                        <td class="border px-4 py-2">{{.CantidadIngresada}}</td>
                        <td class="border px-4 py-2">{{.SaldoAnterior}}</td>
                        <td class="border px-4 py-2">{{.DiasDesdeIngreso}}</td>
                    </tr>
                    {{end}}
                </tbody>
            </table>
        </div>

        <div class="mt-4 flex items-center justify-between">
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
        </div>
    </div>
{{end}}
`

type ViewData struct {
	Items       []models.Saldo
	CurrentPage int
	TotalPages  int
	PageSize    int
	Search      string
	SortField   string
	SortDir     string
}

func (d ViewData) SortIndicator(field string) string {
	if d.SortField == field {
		if d.SortDir == "asc" {
			return "↑"
		}
		return "↓"
	}
	return ""
}

func (d ViewData) NextSort(field string) string {
	if d.SortField == field && d.SortDir == "asc" {
		return "desc"
	}
	return "asc"
}

func RenderSaldos(w http.ResponseWriter, data interface{}) {
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

	if _, err = tmpl.Parse(saldosTemplate); err != nil {
		http.Error(w, "Error al cargar la plantilla de saldos", http.StatusInternalServerError)
		return
	}

	if err = tmpl.ExecuteTemplate(w, "layout.tmpl", data); err != nil {
		http.Error(w, "Error al renderizar la plantilla de saldos", http.StatusInternalServerError)
	}
}
