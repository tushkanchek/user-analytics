package postgres

import (
	"database/sql"
	"fmt"
	"user-analytics/internal/config"
	_"github.com/lib/pq"
	
)



type Storage struct{
	db *sql.DB
}



func New(cfg config.DBConfig) (*Storage, error){
	const op = "storage.postgres.New"


	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode,
	)
	
	db, err := sql.Open("postgres", dsn)
	
	if err!=nil{
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	//table obtain info from /POST events
	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS events(
		id SERIAL PRIMARY KEY,
		user_id INT NOT NULL ,
		event_type TEXT NOT NULL,
		metadata JSONB NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT NOW()
	)`)
	if err!=nil{
		return nil, fmt.Errorf("%s: %w",op, err)
	}

	//index for fast search
	_, err = db.Exec(`
	CREATE INDEX IF NOT EXISTS idx_events_user_id ON events(user_id);
	CREATE INDEX IF NOT EXISTS idx_events_event_type ON events(event_type);
	CREATE INDEX IF NOT EXISTS idx_events_created_at ON events(created_at);
	`)
	if err!=nil{
		return nil, fmt.Errorf("%s: %w",op, err)
	}

	//db for 
	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS agregates(
		user_id INT NOT NULL,
		event_type TEXT NOT NULL,
		count INT,
		last_updated TIMESTAMP NOT NULL DEFAULT NOW()
	)`)
	if err!=nil{
		return nil, fmt.Errorf("%s: %w",op, err)
	}

	
	return &Storage{db:db}, nil
}


