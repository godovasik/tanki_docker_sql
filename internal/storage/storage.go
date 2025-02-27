package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/godovasik/tanki_docker_sql/internal/models"
	"github.com/godovasik/tanki_docker_sql/logger"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Config struct {
	DSN string
}

func NewPostgresPool(cfg Config) (*pgxpool.Pool, error) { // эта функция создает пул соединений. говорят это хорошо.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, cfg.DSN)
	if err != nil {
		return nil, fmt.Errorf("cant create connection pool: %w", err)
	}
	if err := pool.Ping(ctx); err != nil {

		pool.Close()
		return nil, fmt.Errorf("ICANT ping db: %w", err)
	}
	logger.Log.Info("connected to db!")
	return pool, nil
}

// -----------------

type UserRepository interface { //интерфейс бд как я понял
	CreateUser(ctx context.Context, user models.User) error
	GetAllUsers(ctx context.Context) ([]models.User, error)
	GetUserById(ctx context.Context, id int) (*models.User, error)
	DeleteUser(ctx context.Context, id int) error
}

type userRepo struct { // ааааааааа это будет нашим интерфейсом наверно
	db *pgxpool.Pool
}

// -----------------

func NewUserRepository(db *pgxpool.Pool) UserRepository { // я хуй знает че это блядь
	return &userRepo{db: db}
}

func (r *userRepo) CreateUser(ctx context.Context, user models.User) error {
	query := `insert into users (username) values ($1)`
	_, err := r.db.Exec(ctx, query, user.Name)
	return err
}

func (r *userRepo) GetUserById(ctx context.Context, id int) (*models.User, error) {
	// query := `select user_id, name from users where id = $1`
	// user := models.User{}
	// err := r.db.QueryRow()
	// я заебался на сегодня, завтра доделаю Kappa
	return nil, nil
}

func (r *userRepo) GetAllUsers(ctx context.Context) ([]models.User, error) {
	query := `select user_id, name from users`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.Id, &user.Name); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return users, err
}

func (r *userRepo) DeleteUser(ctx context.Context, id int) error {
	query := `delete from users where id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}
