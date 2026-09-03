package main

import (
	"fmt"
	"log"
	"net/http"

	"database/sql"
	"tp2/repository"

	_ "github.com/lib/pq"
)

func abrirDB() (*sql.DB, error) {
	dsn := "host=database port=5432 user=postgres password=supersecret dbname=tp2_db sslmode=disable"
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(25)
	return db, nil
}

func main() {
	db, err := abrirDB()

	if err != nil {
		log.Fatalf("Error al conectar la db: %v", err)
	}
	dir := "./static"

	nuevoUsuario := &repository.User{
		Name:  "Prueba modulos",
		Email: "prueba@ejemplo.com",
	}

	repo := repository.NewUserRepository(db)
	err = repo.CreateUser(nuevoUsuario)

	if err != nil {
		log.Fatalf("Error al crear usuario: %v", err)
	}

	fmt.Printf("Usuario creado correctamente con ID: %d\n", nuevoUsuario.ID)

	fileServer := http.FileServer(http.Dir(dir))
	http.Handle("/", fileServer)

	port := ":8080"

	fmt.Printf("Servidor escuchando en http://localhost%s\n", port)

	err = http.ListenAndServe(port, nil)

	if err != nil {
		fmt.Printf("Error al iniciar el servidor: %s\n", err)
	}

}
