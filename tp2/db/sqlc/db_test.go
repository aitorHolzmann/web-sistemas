package db

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// setupTestDB abre una conexion a la base de test. Usa el mismo host que
// el resto de la app ("database", el nombre del servicio en docker-compose).
func setupTestDB(t *testing.T) *sql.DB {
	t.Helper()

	dsn := "host=database port=5432 user=postgres password=supersecret dbname=tp2_db sslmode=disable"
	conn, err := sql.Open("pgx", dsn)
	if err != nil {
		t.Fatalf("no se pudo abrir la conexion: %v", err)
	}

	if err := conn.Ping(); err != nil {
		t.Fatalf("no se pudo hacer ping a la base: %v", err)
	}

	return conn
}

func TestQueries_CRUD(t *testing.T) {
	conn := setupTestDB(t)
	defer conn.Close()

	queries := New(conn)
	ctx := context.Background()

	var usuarioCreado User

	t.Run("CreateUser", func(t *testing.T) {
		email := fmt.Sprintf("test.sqlc.%d@ejemplo.com", time.Now().UnixNano())

		u, err := queries.CreateUser(ctx, CreateUserParams{
			Name:  "Test sqlc",
			Email: email,
		})
		if err != nil {
			t.Fatalf("CreateUser fallo: %v", err)
		}
		if u.ID == 0 {
			t.Errorf("se esperaba un ID > 0, se obtuvo %d", u.ID)
		}
		if u.Email != email {
			t.Errorf("email esperado %q, se obtuvo %q", email, u.Email)
		}

		usuarioCreado = u
	})

	t.Run("GetUser", func(t *testing.T) {
		u, err := queries.GetUser(ctx, usuarioCreado.ID)
		if err != nil {
			t.Fatalf("GetUser fallo: %v", err)
		}
		if u.Email != usuarioCreado.Email {
			t.Errorf("email esperado %q, se obtuvo %q", usuarioCreado.Email, u.Email)
		}
	})

	t.Run("UpdateUser", func(t *testing.T) {
		nuevoNombre := "Test sqlc actualizado"

		err := queries.UpdateUser(ctx, UpdateUserParams{
			ID:    usuarioCreado.ID,
			Name:  nuevoNombre,
			Email: usuarioCreado.Email,
		})
		if err != nil {
			t.Fatalf("UpdateUser fallo: %v", err)
		}

		u, err := queries.GetUser(ctx, usuarioCreado.ID)
		if err != nil {
			t.Fatalf("GetUser tras el update fallo: %v", err)
		}
		if u.Name != nuevoNombre {
			t.Errorf("nombre esperado %q, se obtuvo %q", nuevoNombre, u.Name)
		}

		usuarioCreado = u
	})

	t.Run("ListUsers", func(t *testing.T) {
		usuarios, err := queries.ListUsers(ctx)
		if err != nil {
			t.Fatalf("ListUsers fallo: %v", err)
		}

		encontrado := false
		for _, u := range usuarios {
			if u.ID == usuarioCreado.ID {
				encontrado = true
				break
			}
		}
		if !encontrado {
			t.Errorf("no se encontro el usuario id=%d en el listado", usuarioCreado.ID)
		}
	})

	t.Run("DeleteUser", func(t *testing.T) {
		if err := queries.DeleteUser(ctx, usuarioCreado.ID); err != nil {
			t.Fatalf("DeleteUser fallo: %v", err)
		}

		_, err := queries.GetUser(ctx, usuarioCreado.ID)
		if err != sql.ErrNoRows {
			t.Errorf("se esperaba sql.ErrNoRows tras borrar, se obtuvo: %v", err)
		}
	})
}
