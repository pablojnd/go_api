package views

import (
	"html/template"
	"net/http"
)

var indexTemplate = `
{{define "title"}}Inicio - API de Saldos{{end}}

{{define "content"}}
    <div class="max-w-4xl mx-auto">
        <h1 class="text-4xl font-bold text-gray-800 mb-8">Bienvenido a la API de Saldos</h1>
        
        <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
            <div class="bg-white p-6 rounded-lg shadow-md hover:shadow-lg transition-shadow duration-300">
                <h2 class="text-2xl font-semibold text-blue-600 mb-4">Saldos</h2>
                <p class="text-gray-600 mb-4">Accede a la información detallada de saldos con funcionalidades de:</p>
                <ul class="list-disc list-inside text-gray-700 space-y-2">
                    <li>Búsqueda en tiempo real</li>
                    <li>Ordenamiento por columnas</li>
                    <li>Paginación dinámica</li>
                    <li>Exportación a Excel</li>
                </ul>
                <a href="/saldos" class="inline-block mt-4 px-6 py-2 bg-blue-500 text-white rounded hover:bg-blue-600 transition-colors duration-300">
                    Ver Saldos
                </a>
            </div>

            <div class="bg-white p-6 rounded-lg shadow-md hover:shadow-lg transition-shadow duration-300">
                <h2 class="text-2xl font-semibold text-green-600 mb-4">Datos Combinados</h2>
                <p class="text-gray-600 mb-4">Visualiza la información combinada con:</p>
                <ul class="list-disc list-inside text-gray-700 space-y-2">
                    <li>Datos integrados de múltiples fuentes</li>
                    <li>Filtrado avanzado</li>
                    <li>Identificación de registros faltantes</li>
                    <li>Exportación personalizada</li>
                </ul>
                <a href="/combined" class="inline-block mt-4 px-6 py-2 bg-green-500 text-white rounded hover:bg-green-600 transition-colors duration-300">
                    Ver Combinados
                </a>
            </div>
        </div>

        <div class="mt-8 bg-white p-6 rounded-lg shadow-md">
            <h2 class="text-2xl font-semibold text-gray-800 mb-4">API REST</h2>
            <p class="text-gray-600 mb-4">Accede a los datos programáticamente mediante nuestra API REST:</p>
            <div class="bg-gray-100 p-4 rounded">
                <code class="text-sm">
                    GET /api/saldos - Obtener lista de saldos<br>
                    GET /api/combined - Obtener datos combinados
                </code>
            </div>
            <a href="/api/saldos" class="inline-block mt-4 px-6 py-2 bg-gray-500 text-white rounded hover:bg-gray-600 transition-colors duration-300">
                Explorar API
            </a>
        </div>
    </div>
{{end}}
`

func RenderIndex(w http.ResponseWriter) {
	// Parseamos el layout y la plantilla específica del index
	tmpl, err := template.New("layout.tmpl").
		Funcs(template.FuncMap{"inc": func(i int) int { return i + 1 }, "dec": func(i int) int {
			if i > 1 {
				return i - 1
			}
			return 1
		}}).
		ParseFiles("c:/Users/pc/Herd/go_api/views/layout.tmpl")
	if err != nil {
		http.Error(w, "Error al cargar el layout", http.StatusInternalServerError)
		return
	}
	// Parseamos la plantilla de contenido
	if _, err = tmpl.Parse(indexTemplate); err != nil {
		http.Error(w, "Error al cargar la vista", http.StatusInternalServerError)
		return
	}
	if err = tmpl.Execute(w, nil); err != nil {
		http.Error(w, "Error al renderizar la vista", http.StatusInternalServerError)
	}
}
