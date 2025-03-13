package models

import(
	"database/sql"
	"time"
	"errors"
	"strings"
	"golang.org/x/crypto/bcrypt"
	"github.com/jackc/pgx/v5/pgconn"
)

type User struct{
	ID int
	Name string
	Email string
	HashedPassword []byte
	Created time.Time
}

type UserModel struct{
	DB *sql.DB
}

func(m *UserModel) Insert(name, email, password string) error{
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil{
		return err
	}

	stmt := `INSERT INTO users (name, email, hashed_password, created)
			VALUES($1, $2, $3, CURRENT_TIMESTAMP AT TIME ZONE 'UTC')`
	
	_, err = m.DB.Exec(stmt, name, email, string(hashedPassword))
	if err != nil{
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr){
			if pgErr.Code == "23505" && strings.Contains(pgErr.Message, "users_uc_email"){
				return errors.New("dublicate email")
			}
		}
		return err
	}
	return nil
}

func (m *UserModel) Authenticate(email, password string) error{
	return nil
}

func (m *UserModel) Exist(id int) (bool, error){
	return false, nil
}