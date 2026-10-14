package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
)

// WriteJSON serializa data como JSON, fija el status code y la
// cabecera Content-Type. La usan todos los handlers para responder
// de forma consistente.
func WriteJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// WriteError responde con un JSON de la forma {"error": "..."} y el
// status code indicado.
func WriteError(w http.ResponseWriter, status int, message string) {
	WriteJSON(w, status, map[string]string{"error": message})
}

// ParseID extrae y parsea el parámetro de ruta {name} como int32
// (todas las PKs de este dominio son SERIAL/int32).
func ParseID(r *http.Request, name string) (int32, error) {
	id, err := strconv.ParseInt(r.PathValue(name), 10, 32)
	if err != nil {
		return 0, err
	}
	return int32(id), nil
}
