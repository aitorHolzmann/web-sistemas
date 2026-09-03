package main // O el paquete donde decidas ubicarlo

import (
	"testing"
	"tp2/models" // Importamos tu modelo

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestGORM_CRUD(t *testing.T) {
	// 1. Configuración y Migración
	// Usamos el mismo DSN que en tu main.go
	dsn := "host=database port=5432 user=postgres password=supersecret dbname=tp2_db sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("Error al conectar a la BD de prueba con GORM: %v", err)
	}

	// Ejecuta AutoMigrate[cite: 2]
	err = db.AutoMigrate(&models.User{})
	if err != nil {
		t.Fatalf("Error en AutoMigrate: %v", err)
	}

	// Buena práctica: limpiar la tabla antes de cada ejecución[cite: 2].
	// TRUNCATE borra los datos y RESTART IDENTITY reinicia los IDs a 1 en PostgreSQL.
	db.Exec("TRUNCATE TABLE users RESTART IDENTITY CASCADE")

	// --- 1. Crear ---
	// Crea un nuevo User y usa db.Create()[cite: 2].
	nuevoUsuario := models.User{
		Name:  "Test User",
		Email: "test@ejemplo.com",
	}
	result := db.Create(&nuevoUsuario)

	// Comprueba que result.Error sea nil[cite: 2].
	if result.Error != nil {
		t.Fatalf("Falló la creación del usuario: %v", result.Error)
	}
	// Comprueba que el ID del objeto se haya actualizado[cite: 2].
	if nuevoUsuario.ID == 0 {
		t.Errorf("Se esperaba un ID asignado, pero es 0")
	}

	// --- 2. Leer ---
	// Usa db.First() para buscar el usuario por su ID[cite: 2].
	var usuarioLeido models.User
	if err := db.First(&usuarioLeido, nuevoUsuario.ID).Error; err != nil {
		t.Fatalf("Falló la lectura del usuario: %v", err)
	}
	// Verifica que los datos recuperados sean correctos[cite: 2].
	if usuarioLeido.Name != nuevoUsuario.Name {
		t.Errorf("Se esperaba el nombre '%s', pero se obtuvo '%s'", nuevoUsuario.Name, usuarioLeido.Name)
	}

	// --- 3. Actualizar ---
	// Modifica el objeto User y usa db.Model().Update() para persistir[cite: 2].
	nuevoNombre := "Nombre Modificado"
	if err := db.Model(&usuarioLeido).Update("Name", nuevoNombre).Error; err != nil {
		t.Fatalf("Falló la actualización: %v", err)
	}

	// Vuelve a leer el registro para confirmar la actualización[cite: 2].
	var usuarioActualizado models.User
	if err := db.First(&usuarioActualizado, nuevoUsuario.ID).Error; err != nil {
		t.Fatalf("Falló al volver a leer el usuario actualizado: %v", err)
	}
	if usuarioActualizado.Name != nuevoNombre {
		t.Errorf("La actualización no se aplicó. Se esperaba '%s', se obtuvo '%s'", nuevoNombre, usuarioActualizado.Name)
	}

	// --- 4. Eliminar ---
	// Usa db.Delete() para borrar el usuario[cite: 2].
	if err := db.Delete(&usuarioActualizado).Error; err != nil {
		t.Fatalf("Falló la eliminación: %v", err)
	}

	// Intenta buscarlo de nuevo y espera un error del tipo gorm.ErrRecordNotFound[cite: 2].
	var usuarioEliminado models.User
	err = db.First(&usuarioEliminado, nuevoUsuario.ID).Error
	if err != gorm.ErrRecordNotFound {
		t.Errorf("Se esperaba el error gorm.ErrRecordNotFound tras eliminar, pero se obtuvo: %v", err)
	}
}
