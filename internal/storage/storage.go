package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/godovasik/tanki_docker_sql/internal/models"
	"github.com/godovasik/tanki_docker_sql/logger"
	"github.com/jackc/pgx/v5"
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
	cleanup := func() { pool.Close(); logger.Log.Infof("disconnected from db") }

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

// TODO - поменять заглавные на строчные буквы где нужно
type UserRepository interface { //интерфейс бд как я понял
	CreateUser(ctx context.Context, user models.User) error
	GetAllUsers(ctx context.Context) ([]models.User, error)
	GetUserById(ctx context.Context, id int) (*models.User, error)
	DeleteUser(ctx context.Context, id int) error
	// FindLastChangedDatastamp(ctx context.Context, user_id int) (*models.Datastamp, error)
	AddDatastamp(ctx context.Context, data *models.Datastamp, user_id int) (int, error)
	AddGearStats(ctx context.Context, datastamp_id, gear_key int, gearData models.GearData) error
	FindLastStampDate(ctx context.Context, user_id int) (time.Time, error)
	FindLastGearStats(ctx context.Context, user_id, gear_key int) (*models.GearData, error)
	UpdateDataForUser(ctx context.Context, data *models.Datastamp, user_id int) error
	loadGearMap(ctx context.Context, GearMap *map[string]models.GearData, user_id, datastamp_id int) (int, error)
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

func (r *userRepo) AddDatastamp(ctx context.Context, data *models.Datastamp, user_id int) (int, error) {
	timeHourRounded := time.Now().Truncate(time.Hour)
	query := `
		INSERT INTO datastamps (user_id, created_at, rank, kills, deaths, cry)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING datastamp_id
	`
	var datastampID int
	err := r.db.QueryRow(ctx, query, user_id, timeHourRounded, data.Rank, data.Kills, data.Deaths, data.EarnedCrystals).Scan(&datastampID)
	if err != nil {
		return 0, err
	}
	return datastampID, nil
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
		if err == pgx.ErrNoRows {
			return time.Unix(0, 0), nil
		}
		// logger.Log.Error(err)
		return time.Unix(0, 0), err
	}

	return *created_at, nil
}

// ТОДО - можно добавлять сразу несколько записей в одной квери, но мне лень
func (r *userRepo) AddGearStats(ctx context.Context, datastamp_id, gear_key int, gearData models.GearData) error {
	query := `
		INSERT INTO gear_stats (datastamp_id, gear_key, score_earned, seconds_played) 
		VALUES ($1, $2, $3, $4);
	`
	_, err := r.db.Exec(ctx, query, datastamp_id, gear_key, gearData.ScoreEarned, gearData.SecondsPlayed)
	return err
}

func (r *userRepo) FindLastGearStats(ctx context.Context, user_id, gear_key int) (*models.GearData, error) {
	query := ` SELECT gs.score_earned, gs.seconds_played
		FROM gear_stats gs
		JOIN datastamps ds ON gs.datastamp_id = ds.datastamp_id
		WHERE ds.user_id = $1 AND gs.gear_key = $2
		ORDER BY ds.datastamp_id DESC
		LIMIT 1;
	`
	var gearData models.GearData
	err := r.db.QueryRow(ctx, query, user_id, gear_key).Scan(&gearData.ScoreEarned, &gearData.SecondsPlayed)
	if err != nil {
		if err == pgx.ErrNoRows {
			return &gearData, nil
		}
		return nil, err
	}
	return &gearData, nil

}

func (r *userRepo) UpdateDataForUser(ctx context.Context, data *models.Datastamp, user_id int) error { // do fucking everything
	lastUpdate, err := r.FindLastStampDate(ctx, user_id)
	if err != nil {
		return err
	}

	// тут можно менять ограничение на то как часто могут писаться стампы. пишутся они все равно с truncate(Hour), независимо от значения снизу
	now := time.Now() //.Truncate(time.Hour)
	if now == lastUpdate {
		return fmt.Errorf("user %s is already up to date", data.Name)
	}

	datastamp_id, err := r.AddDatastamp(ctx, data, user_id)
	if err != nil {
		return err
	}

	hullsAdded, err := r.loadGearMap(ctx, &data.Hulls, user_id, datastamp_id)
	if err != nil {
		return err
	}
	turretsAdded, err := r.loadGearMap(ctx, &data.Turrets, user_id, datastamp_id)
	if err != nil {
		return err
	}

	logger.Log.Debugf("added %v hulls and %v turrets for user %v", hullsAdded, turretsAdded, data.Name)

	return nil

}

// он почемуто возвращает значение на 2 больше настоящего, хуй знает почему
func (r *userRepo) loadGearMap(ctx context.Context, GearMap *map[string]models.GearData, user_id, datastamp_id int) (int, error) {
	lines_added := 0
	for gearName, currentGearData := range *GearMap {
		gear_key, isthere := models.GetGearId(gearName)
		if !isthere {
			logger.Log.Errorf("bro there is no such gear: %s", gearName)
		} else {
			lastGearData, err := r.FindLastGearStats(ctx, user_id, gear_key)
			if err != nil {
				return lines_added, err
			}
			if currentGearData != *lastGearData {
				err = r.AddGearStats(ctx, datastamp_id, gear_key, currentGearData)
				if err != nil {
					return lines_added, fmt.Errorf("gear name: %v, err: %v", gearName, err)
				}
				lines_added++
			}
		}

	}
	return lines_added, nil
}

// эта штука не дописана и поке не используется. будем сохранять все датастампы пока что.
// func (r *userRepo) FindLastChangedDatastamp(ctx context.Context, user_id int) (*models.Datastamp, error) {
// 	query := `
// 	WITH last_values AS (
// 		SELECT
// 			(SELECT rank FROM datastamps d WHERE d.user_id = $1 AND rank IS NOT NULL ORDER BY created_at DESC LIMIT 1) AS last_rank,
// 			(SELECT kills FROM datastamps d WHERE d.user_id = $1 AND kills IS NOT NULL ORDER BY created_at DESC LIMIT 1) AS last_kills,
// 			(SELECT deaths FROM datastamps d WHERE d.user_id = $1 AND deaths IS NOT NULL ORDER BY created_at DESC LIMIT 1) AS last_deaths,
// 			(SELECT cry FROM datastamps d WHERE d.user_id = $1 AND cry IS NOT NULL ORDER BY created_at DESC LIMIT 1) AS last_cry
// 		FROM datastamps ds
// 		WHERE ds.user_id = $1
// 		LIMIT 1
// 	)
// 	SELECT * FROM last_values;`

// 	// data := models.Datastamp{}
// 	var lastRank, lastKills, lastDeaths, lastCry *int
// 	err := r.db.QueryRow(context.Background(), query, user_id).Scan(&lastRank, &lastKills, &lastDeaths, &lastCry)
// 	if err != nil {
// 		return nil, err
// 	}
// 	logger.Log.Debugf("rank: %d, kills: %d, deaths: %d, cry: %d", lastRank, lastKills, lastDeaths, lastCry)
// 	return nil, nil

// }
