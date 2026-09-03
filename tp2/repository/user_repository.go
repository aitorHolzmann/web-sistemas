package repository

import (
	"database/sql"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// User representa la entidad de base de datos
type User struct {
	ID        int
	Name      string
	Email     string
	CreatedAt time.Time
}

// Esto es para abstraer la logica de manejo de DB.
// De esta forma podemos cambiar de db sin problema
type UserRepository struct {
	db *sql.DB
}

// NewUserRepository inicializa el repositorio
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(user *User) error {
	query := `INSERT INTO users (name, email) VALUES ($1, $2) RETURNING id, created_at`
	err := r.db.QueryRow(query, user.Name, user.Email).Scan(&user.ID, &user.CreatedAt)

	return err
}

func (r *UserRepository) GetUserByID(id int) (*User, error) {
	user := &User{}

	query := `SELECT * FROM users WHERE id = $1;`

	//QueryRow se usa cuando esperamos como maximo una fila.
	err := r.db.QueryRow(query, id).Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt)

	if err != nil {
		return nil, err //Si no existe devuelve sql.ErrNoRows
	}
	return user, err
}

func (r *UserRepository) ListUsers() ([]*User, error) {
	query := `SELECT * FROM users`
	rows, err := r.db.Query(query)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*User

	for rows.Next() {
		user := &User{}
		err = rows.Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt)

		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	//Verificamos que no haya habido errores durante la iteracion
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (r *UserRepository) UpdateUser(user *User) error {
	query := `UPDATE users SET name = $1, email = $2 WHERE id = $3;`

	result, err := r.db.Exec(query, user.Name, user.Email, user.ID)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()

	if err != nil {
		return err
	}

	// Exec no falla si el WHERE no matcheó ninguna fila, así que
	// chequeamos RowsAffected para poder avisar si el ID no existía.
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *UserRepository) DeleteUser(id int) error {
	query := `DELETE FROM users WHERE id = $1;`
	result, err := r.db.Exec(query, id)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()

	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
