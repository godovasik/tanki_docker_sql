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

func ConnectToDb() (UserRepository, func(), error) {
	DSN := "postgres://tanki_enjoyer:r@172.24.125.42:5432/game_stats"
	logger.Log.Info("connecting to db...")

	cfg := Config{DSN: DSN}
	pool, err := NewPostgresPool(cfg)
	if err != nil {
		logger.Log.Error("db connection err:", err)
	}
	cleanup := func() { pool.Close() }

	logger.Log.Info("we connected to db!")

	userRepo := NewUserRepository(pool)
	return userRepo, cleanup, err
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
	FindLastChangedDatastamp(ctx context.Context, user_id int) (*models.Datastamp, error)
	AddDatastamp(ctx context.Context, data models.Datastamp, user_id int) error
	FindLastStampDate(ctx context.Context, user_id int) (time.Time, error)
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
	query := `select name from users where user_id = $1`
	user := &models.User{}
	err := r.db.QueryRow(ctx, query, id).Scan(&user.Name)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return nil, fmt.Errorf("user not found")
		} else {
			return nil, err
		}
	}
	user.Id = id
	return user, err
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

func (r *userRepo) DeleteUser(ctx context.Context, user_id int) error {
	query := `delete from users where user_id = $1`
	_, err := r.db.Exec(ctx, query, user_id)
	return err
}

func (r *userRepo) AddDatastamp(ctx context.Context, data models.Datastamp, user_id int) error {
	timeHourRounded := time.Now().Truncate(time.Hour)
	query := `
		insert into datastamps (user_id, created_at, rank, kills, deaths, cry)
		values ($1, $2, $3, $4, $5, $6)
	`
	_, err := r.db.Exec(ctx, query, user_id, timeHourRounded, data.Rank, data.Kills, data.Deaths, data.EarnedCrystals)
	return err
}

func (r *userRepo) FindLastStampDate(ctx context.Context, user_id int) (time.Time, error) {
	query := `SELECT created_at
				FROM datastamps
				WHERE datastamp_id = (
					SELECT MAX(datastamp_id) 
					FROM datastamps 
					WHERE user_id = $1
				);	
	`
	var created_at *time.Time
	err := r.db.QueryRow(ctx, query, user_id).Scan(&created_at)
	if err != nil {
		// logger.Log.Error(err)
		return time.Unix(0, 0), err
	}

	return *created_at, err

}

// эта штука не дописана и поке не используется. будем сохранять все датастампы пока что.
func (r *userRepo) FindLastChangedDatastamp(ctx context.Context, user_id int) (*models.Datastamp, error) {
	query := `
	WITH last_values AS (
		SELECT
			(SELECT rank FROM datastamps d WHERE d.user_id = $1 AND rank IS NOT NULL ORDER BY created_at DESC LIMIT 1) AS last_rank,
			(SELECT kills FROM datastamps d WHERE d.user_id = $1 AND kills IS NOT NULL ORDER BY created_at DESC LIMIT 1) AS last_kills,
			(SELECT deaths FROM datastamps d WHERE d.user_id = $1 AND deaths IS NOT NULL ORDER BY created_at DESC LIMIT 1) AS last_deaths,
			(SELECT cry FROM datastamps d WHERE d.user_id = $1 AND cry IS NOT NULL ORDER BY created_at DESC LIMIT 1) AS last_cry
		FROM datastamps ds
		WHERE ds.user_id = $1
		LIMIT 1
	)
	SELECT * FROM last_values;`

	// data := models.Datastamp{}
	var lastRank, lastKills, lastDeaths, lastCry *int
	err := r.db.QueryRow(context.Background(), query, user_id).Scan(&lastRank, &lastKills, &lastDeaths, &lastCry)
	if err != nil {
		return nil, err
	}
	logger.Log.Debugf("rank: %d, kills: %d, deaths: %d, cry: %d", lastRank, lastKills, lastDeaths, lastCry)
	return nil, nil

}
