// Package web contiene los handlers que devuelven HTML (vistas templ) y
// procesan formularios clásicos. Todos aplican Post/Redirect/Get: después de
// un POST redirigen a GET /, que vuelve a renderizar la página completa.
package web

import (
	"errors"
	"log"
	"net/http"
	"net/url"

	"github.com/lib/pq"
)

// Mensajes que puede mostrar la página tras un redirect. La URL solo lleva un
// código (/?error=codigo); el texto sale de acá, así que nadie puede armar un
// link que muestre un mensaje arbitrario.
var mensajes = map[string]string{
	"invalido":  "Datos inválidos: revisá los campos del formulario.",
	"no_existe": "El elemento referenciado no existe.",
	"error":     "Ocurrió un error inesperado.",
}

// Códigos SQLSTATE de PostgreSQL que nos interesan.
const (
	fkViolation = "23503"
)

// pgCode devuelve el código SQLSTATE si err es un error de PostgreSQL.
func pgCode(err error) string {
	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		return string(pqErr.Code)
	}
	return ""
}

// redirectOK vuelve a la página principal (303 See Other: el navegador
// repite la petición como GET).
func redirectOK(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// redirectError vuelve a la página principal mostrando el mensaje de `codigo`.
func redirectError(w http.ResponseWriter, r *http.Request, codigo string) {
	http.Redirect(w, r, "/?error="+url.QueryEscape(codigo), http.StatusSeeOther)
}

// fail registra un error inesperado y avisa al usuario.
func fail(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("%s %s: %v", r.Method, r.URL.Path, err)
	redirectError(w, r, "error")
}
