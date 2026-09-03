package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	db "tp2/db/sqlc"
	//Estoy importando la carpeta db.
	//De esta forma puedo usar el codigo generado por sqlc.
	//Para poder usarlo primero debo ejecutar sqlc generate

	_ "github.com/jackc/pgx/v5/stdlib"
)

func abrirDB() (*sql.DB, error) { //Esto es igual que en el tp1. Es para poder ejecutar sql
	connStr := "host=database port=5432 user=postgres password=supersecret dbname=tp2_db sslmode=disable"
	//conn, err := sql.Open("postgres", connStr)
	conn, err := sql.Open("pgx", connStr)
	if err != nil {
		return nil, err
	}

	if err := conn.Ping(); err != nil {
		return nil, err
	}

	conn.SetMaxOpenConns(25)
	return conn, nil
}

func main() {
	conn, err := abrirDB()
	if err != nil {
		log.Fatalf("Error al conectar la db: %v", err)
	}
	defer conn.Close()

	ctx := context.Background() //Esto lo necesita sqlc
	queries := db.New(conn)     //Metodo de db.go. Crea objeto queries que conoce el nombre de los metodos

	//
	//HASTA ACA LA CONEXION
	//

	// 1. Crear
	//Resulta que sqlc genera un struct por cada consulta.
	//Por ejemplo db.CreateUserParams
	//Entonces lo llenamos y el lo usa para hacer la consulta

	nuevoUsuario, err := queries.CreateUser(ctx,
		db.CreateUserParams{
			Name:  "Prueba sqlc",
			Email: "prueba.sqlc@ejemplo.com",
		})
	if err != nil {
		log.Fatalf("Error al crear usuario: %v", err)
	}
	fmt.Printf("Usuario creado correctamente con ID: %d\n", nuevoUsuario.ID)

	// 2. Leer
	usuario, err := queries.GetUser(ctx, nuevoUsuario.ID)
	if err != nil {
		log.Fatalf("Error al leer usuario: %v", err)
	}
	fmt.Printf("Usuario leido: %+v\n", usuario)

	// 3. Actualizar
	err = queries.UpdateUser(ctx, db.UpdateUserParams{
		ID:    nuevoUsuario.ID,
		Name:  "Prueba sqlc actualizada",
		Email: nuevoUsuario.Email,
	})
	if err != nil {
		log.Fatalf("Error al actualizar usuario: %v", err)
	}
	fmt.Println("Usuario actualizado correctamente")

	// 4. Listar
	usuarios, err := queries.ListUsers(ctx)
	if err != nil {
		log.Fatalf("Error al listar usuarios: %v", err)
	}
	fmt.Printf("Hay %d usuario(s) en la tabla\n", len(usuarios))

	// 5. Eliminar
	err = queries.DeleteUser(ctx, nuevoUsuario.ID)
	if err != nil {
		log.Fatalf("Error al eliminar usuario: %v", err)
	}
	fmt.Println("Usuario eliminado correctamente")
}
