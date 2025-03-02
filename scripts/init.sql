-- Создаем базу данных, если её нет
DO
$do$
BEGIN
    IF NOT EXISTS (SELECT FROM pg_database WHERE datname = 'game_stats') THEN
        CREATE DATABASE game_stats;
    END IF;
END
$do$;

-- Подключаемся к базе game_stats
\c game_stats

-- Создаем пользователя, если его нет
DO
$do$
BEGIN
    IF NOT EXISTS (SELECT FROM pg_roles WHERE rolname = 'tanki_enjoyer') THEN
        CREATE ROLE tanki_enjoyer WITH LOGIN PASSWORD 'rock';
    END IF;
END
$do$;

-- Выдаем права пользователю
ALTER DATABASE game_stats OWNER TO tanki_enjoyer;
GRANT ALL PRIVILEGES ON DATABASE game_stats TO tanki_enjoyer;

-- Создаем таблицу пользователей
CREATE TABLE IF NOT EXISTS users (
    user_id SERIAL PRIMARY KEY,
    name VARCHAR(128) NOT NULL UNIQUE
);

-- Создаем таблицу датастампов
CREATE TABLE IF NOT EXISTS datastamps (
    datastamp_id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL,
    rank SMALLINT,
    kills INTEGER,
    deaths INTEGER,
    cry INTEGER
);

-- Создаем индекс на поле created_at в таблице datastamps
CREATE INDEX IF NOT EXISTS idx_datastamps_created_at ON datastamps(created_at);

-- Создаем таблицу gear_stats
CREATE TABLE IF NOT EXISTS gear_stats (
    datastamp_id INTEGER NOT NULL REFERENCES datastamps(datastamp_id) ON DELETE CASCADE,
    hull_key SMALLINT NOT NULL,
    score_earned INTEGER NOT NULL,
    seconds_played INTEGER NOT NULL,
    PRIMARY KEY (datastamp_id, hull_key)
);

-- Выдаем права пользователю на все таблицы и последовательности
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO tanki_enjoyer;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO tanki_enjoyer;
