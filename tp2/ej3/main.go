package main

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	models "tp2/models"
)

func main() {
	// Cadena de conexión (DSN) reutilizada de tu Ejercicio 2
	dsn := "host=database port=5432 user=postgres password=supersecret dbname=tp2_db sslmode=disable"

	// Establecer la conexión usando GORM
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Falló la conexión a la base de datos: %v", err)
	}
	log.Println("Conexión exitosa con GORM.")

	// Ejecutar la migración automática
	err = db.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatalf("Falló la migración: %v", err)
	}
	log.Println("Migración de la tabla 'users' completada.")

	// A partir de acá irían las operaciones CRUD...

	// --- 1. Crear ---
	// Crea una instancia del struct y se pasa a db.Create().
	// GORM genera el INSERT automáticamente.
	nuevoUsuario := models.User{
		Name:  "GORM Test",
		Email: "gorm@ejemplo.com",
	}
	result := db.Create(&nuevoUsuario) // Pasa el puntero del objeto.
	if result.Error != nil {
		log.Fatalf("Falló al crear usuario: %v", result.Error)
	}
	log.Printf("Usuario creado con ID: %d\n", nuevoUsuario.ID)

	// --- 2. Leer (Uno) ---
	// Se usa db.First() para recuperar un registro por su clave primaria.
	var usuarioLeido models.User
	if err := db.First(&usuarioLeido, nuevoUsuario.ID).Error; err != nil {
		log.Fatalf("Falló al leer usuario: %v", err)
	}
	log.Printf("Usuario leído: %+v\n", usuarioLeido)

	// --- 3. Actualizar ---
	// Se puede actualizar un registro con db.Model().Update()[cite: 1, 2].
	if err := db.Model(&usuarioLeido).Update("Name", "GORM Test Actualizado").Error; err != nil {
		log.Fatalf("Falló al actualizar usuario: %v", err)
	}
	log.Printf("Usuario actualizado: %+v\n", usuarioLeido)

	// --- 4. Leer (Todos) ---
	// db.Find() recupera múltiples registros que coinciden con una condición.
	var usuarios []models.User
	if err := db.Find(&usuarios).Error; err != nil {
		log.Fatalf("Falló al listar usuarios: %v", err)
	}
	log.Printf("Hay %d usuario(s) en la tabla\n", len(usuarios))

	// --- 5. Eliminar ---
	// db.Delete() elimina un registro.
	if err := db.Delete(&models.User{}, usuarioLeido.ID).Error; err != nil {
		log.Fatalf("Falló al eliminar usuario: %v", err)
	}
	log.Println("Usuario eliminado correctamente")
}
