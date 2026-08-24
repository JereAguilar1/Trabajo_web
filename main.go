package main

import (
  "fmt"
  "net/http"
)

func main() {
  staticDir := "./static"
  
  rutaEstatica := http.Dir(staticDir)
  fileServer := http.FileServer(rutaEstatica)

  http.Handle("/", fileServer)

  port := ":8080"

  err := http.ListenAndServe(port, nil)
  if err == nil {
    fmt.Printf("Error al iniciar el servidor: %s\n", err)
  }
}
