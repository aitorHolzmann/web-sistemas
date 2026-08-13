package main

import (
	"fmt"
	"html/template"
	"net/http"
)

func main(){
  
  //De esta forma "publicamos" la carpeta abierta. 
  // Si conozco el nombre del archivo lo puedo ir a buscar a mano
  //staticDir := "./static"
  // fileServer := http.FileServer(http.Dir(staticDir))
  // http.Handle("/", fileServer)
  // fmt.Println("Sirviendo archivos desde %s", staticDir)

  //Si quiero tener mayor control puedo hacer lo siguiente
  http.HandleFunc("/", router)
  port := ":8080"
  fmt.Println("Servidor escuchand en https://localhost%s", port)

  err := http.ListenAndServe(port, nil)
  if (err != nil){
    fmt.Println("Error al iniciar el servidor: %s", err)
  }

}

func router (w http.ResponseWriter, r *http.Request){
  userPath := r.URL.Path
  basePath := "static/"
  // 3. Establece la cabecera Content-Type
  w.Header().Set("Content-Type", "text/html; charset=utf-8")
		// 4. Escribe el HTML en la respuesta
		// fmt.Fprint(w, htmlContent)
  switch userPath {
    case "/": 
      path := basePath+"index.html"
		  http.ServeFile(w, r, path)
    case "/about": 
      path := basePath+"about.html"
		  http.ServeFile(w, r, path)
    case "/formulario": 
      path := basePath+"formulario.html"
		  http.ServeFile(w, r, path)

    case "/saludo":
      
      if r.Method != http.MethodPost {
        http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
        return
      }
      
      // 1. Parsear los datos del formulario (¡Crucial!)
      if err := r.ParseForm(); err != nil {
        http.Error(w, "Error al parsear", http.StatusBadRequest)
        return
      }
      
      // 2. Obtener el valor del campo 'user'
      username := r.FormValue("nombre")
      path := basePath + "saludar.html"

      // 3. Parsear el archivo HTML como un template
      tmpl, err := template.ParseFiles(path)
      if err != nil {
        http.Error(w, "Error al cargar la página", http.StatusInternalServerError)
        return
      }

      err = tmpl.Execute(w, username)
      // 4. Generar y enviar la respuesta HTML
      w.Header().Set("Content-Type", "text/html; charset=utf-8")
      fmt.Fprintf(w, `<!DOCTYPE html><html><head><title>Bienvenido</title></head> <body><h1>¡Hola, %s!</h1>
        <p>Recibimos tus datos.</p> <a href="/">Volver</a></body></html>`, username)

		  http.ServeFile(w, r, path)
    default:
      path := basePath+"error.html"
		  http.ServeFile(w, r, path)
  }
}