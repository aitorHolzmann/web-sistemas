package repository

import (
	"database/sql"
	"fmt"
	"testing"
	"time"

	_ "github.com/lib/pq"
)

// setupTestDB abre una conexión a la misma base usada por la app.
// En un proyecto más grande esto normalmente apuntaría a una DB de test
// separada, pero para este TP alcanza con reusar la misma instancia.
func setupTestDB(t *testing.T) *sql.DB {
	// t.Helper() le dice a Go que, si algo falla acá adentro, reporte la
	// línea de quien LLAMÓ a esta función, no la línea de acá adentro.
	t.Helper()

	dsn := "host=database port=5432 user=postgres password=supersecret dbname=tp2_db sslmode=disable"
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Fatalf("no se pudo abrir la conexión: %v", err)
	}

	if err := db.Ping(); err != nil {
		t.Fatalf("no se pudo hacer ping a la db: %v", err)
	}

	return db
}

func TestUserRepository_CRUD(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close() // se ejecuta al final de la función, pase lo que pase

	repo := NewUserRepository(db)

	// Email único por corrida, para no chocar con el UNIQUE de la tabla
	// si corrés el test varias veces sin limpiar la DB a mano.
	email := fmt.Sprintf("test-%d@ejemplo.com", time.Now().UnixNano())

	user := &User{
		Name:  "Usuario de Prueba",
		Email: email,
	}

	// --- 1. Crear ---
	t.Run("Crear", func(t *testing.T) {
		err := repo.CreateUser(user)
		if err != nil {
			t.Fatalf("CreateUser devolvió error: %v", err)
		}
		if user.ID == 0 {
			t.Errorf("esperaba que CreateUser asignara un ID, quedó en 0")
		}
	})

	// --- 2. Leer ---
	t.Run("Leer", func(t *testing.T) {
		encontrado, err := repo.GetUserByID(user.ID)
		if err != nil {
			t.Fatalf("GetUserByID devolvió error: %v", err)
		}
		if encontrado.Name != user.Name {
			t.Errorf("nombre esperado %q, obtuve %q", user.Name, encontrado.Name)
		}
		if encontrado.Email != user.Email {
			t.Errorf("email esperado %q, obtuve %q", user.Email, encontrado.Email)
		}
	})

	// --- 3. Actualizar ---
	t.Run("Actualizar", func(t *testing.T) {
		user.Name = "Nombre Actualizado"
		err := repo.UpdateUser(user)
		if err != nil {
			t.Fatalf("UpdateUser devolvió error: %v", err)
		}

		// Volvemos a leer para confirmar que el cambio quedó guardado
		// en la DB, no solo en la variable local en memoria.
		actualizado, err := repo.GetUserByID(user.ID)
		if err != nil {
			t.Fatalf("GetUserByID después de actualizar devolvió error: %v", err)
		}
		if actualizado.Name != "Nombre Actualizado" {
			t.Errorf("esperaba nombre actualizado, obtuve %q", actualizado.Name)
		}
	})

	// --- 4. Listar ---
	t.Run("Listar", func(t *testing.T) {
		usuarios, err := repo.ListUsers()
		if err != nil {
			t.Fatalf("ListUsers devolvió error: %v", err)
		}

		encontrado := false
		for _, u := range usuarios {
			if u.ID == user.ID && u.Name == "Nombre Actualizado" {
				encontrado = true
				break
			}
		}
		if !encontrado {
			t.Errorf("el usuario actualizado no aparece en ListUsers")
		}
	})

	// --- 5. Eliminar ---
	t.Run("Eliminar", func(t *testing.T) {
		err := repo.DeleteUser(user.ID)
		if err != nil {
			t.Fatalf("DeleteUser devolvió error: %v", err)
		}

		_, err = repo.GetUserByID(user.ID)
		if err != sql.ErrNoRows {
			t.Errorf("esperaba sql.ErrNoRows después de borrar, obtuve: %v", err)
		}
	})
}
