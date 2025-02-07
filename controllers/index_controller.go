package controllers

import "net/http"

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	// Puedes modificar este contenido según lo que necesites mostrar o redirigir.
	w.Write([]byte("Bienvenido a la API de Saldos"))
}
