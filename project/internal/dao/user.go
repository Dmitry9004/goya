package dao

import (
	"goya/project/internal/models"

	"context"
	"log"
	"database/sql"
	//"sync"
	
	_ "github.com/mattn/go-sqlite3"
)



type UserDAO struct {
	//mu sync.Mutex
	db *sql.DB
	ctx context.Context
}

func NewUserDAO(ctx context.Context, db *sql.DB) *UserDAO{
	return &UserDAO {
		db:db,
		ctx: ctx,
	}
}

func (dao *UserDAO) CreateUsersTable() error {
	usersTable := `
		CREATE TABLE IF NOT EXISTS users(
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT UNIQUE,
			password TEXT
		);
	`
	if _, err := dao.db.ExecContext(dao.ctx, usersTable); err != nil {
		return err
	}
	log.Println("ALL OK!")
	return nil
}

func (dao *UserDAO) Save(user *models.User) (int, error) {
	query := `
		INSERT INTO users (username, password) VALUES ($1, $2)
	`
	
	result, err := dao.db.ExecContext(dao.ctx, query, user.Username, user.Password);
	if err != nil {
		log.Println(err)
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		log.Println(err)
		return 0, err
	}
	
	return int(id), err
}

func (dao *UserDAO) GetUserById(id int) (*models.User, error) {
	query := `
		SELECT * FROM users WHERE id = $1
	`
	var user models.User
	err := dao.db.QueryRowContext(dao.ctx, query, id).Scan(&user.Id, &user.Username, &user.Password)
	if err != nil {
		log.Println(err)
		return &user, err
	}
	
	return &user, nil
}


func (dao *UserDAO) GetUserByUsername(username string) (*models.User, error) {
	query := `
		SELECT * FROM users WHERE username = $1
	`
	var user models.User
	err := dao.db.QueryRowContext(dao.ctx, query, username).Scan(&user.Id, &user.Username, &user.Password)
	if err != nil {
		log.Println(err)
		return &user, err
	}
	
	return &user, nil
}

