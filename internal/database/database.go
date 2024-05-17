package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/joho/godotenv/autoload"
)

// Service represents a service that interacts with a database.
type Service interface {
	// Health returns a map of health status information.
	// The keys and values in the map are service-specific.
	Health() map[string]string

	// Close terminates the database connection.
	// It returns an error if the connection cannot be closed.
	Close() error

	AllAlbums() ([]Album, error)

	AlbumById(id int) (Album, error)

	AddAlbum(alb Album) (int64, error)

	DeleteAlbumByID(id int) (int, error)
}

type Album struct {
	ID     int
	Title  string
	Artist string
	Price  float32
}

type service struct {
	db *sql.DB
}

var (
	database   = os.Getenv("DB_DATABASE")
	password   = os.Getenv("DB_PASSWORD")
	username   = os.Getenv("DB_USERNAME")
	port       = os.Getenv("DB_PORT")
	host       = os.Getenv("DB_HOST")
	dbInstance *service
)

func New() Service {
	// Reuse Connection
	if dbInstance != nil {
		return dbInstance
	}
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", username, password, host, port, database)
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		log.Fatal(err)
	}
	dbInstance = &service{
		db: db,
	}
	return dbInstance
}

// Health checks the health of the database connection by pinging the database.
// It returns a map with keys indicating various health statistics.
func (s *service) Health() map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	stats := make(map[string]string)

	// Ping the database
	err := s.db.PingContext(ctx)
	if err != nil {
		stats["status"] = "down"
		stats["error"] = fmt.Sprintf("db down: %v", err)
		log.Fatalf(fmt.Sprintf("db down: %v", err)) // Log the error and terminate the program
		return stats
	}

	// Database is up, add more statistics
	stats["status"] = "up"
	stats["message"] = "It's healthy"

	// Get database stats (like open connections, in use, idle, etc.)
	dbStats := s.db.Stats()
	stats["open_connections"] = strconv.Itoa(dbStats.OpenConnections)
	stats["in_use"] = strconv.Itoa(dbStats.InUse)
	stats["idle"] = strconv.Itoa(dbStats.Idle)
	stats["wait_count"] = strconv.FormatInt(dbStats.WaitCount, 10)
	stats["wait_duration"] = dbStats.WaitDuration.String()
	stats["max_idle_closed"] = strconv.FormatInt(dbStats.MaxIdleClosed, 10)
	stats["max_lifetime_closed"] = strconv.FormatInt(dbStats.MaxLifetimeClosed, 10)

	// Evaluate stats to provide a health message
	if dbStats.OpenConnections > 40 { // Assuming 50 is the max for this example
		stats["message"] = "The database is experiencing heavy load."
	}

	if dbStats.WaitCount > 1000 {
		stats["message"] = "The database has a high number of wait events, indicating potential bottlenecks."
	}

	if dbStats.MaxIdleClosed > int64(dbStats.OpenConnections)/2 {
		stats["message"] = "Many idle connections are being closed, consider revising the connection pool settings."
	}

	if dbStats.MaxLifetimeClosed > int64(dbStats.OpenConnections)/2 {
		stats["message"] = "Many connections are being closed due to max lifetime, consider increasing max lifetime or revising the connection usage pattern."
	}

	return stats
}

func (s *service) AllAlbums() ([]Album, error) {
	// An albums slice to hold data from returned rows
	var albums []Album

	q := `
		SELECT id, title, artist, price FROM album
	`

	rows, err := s.db.Query(q)
	if err != nil {
		return albums, fmt.Errorf("[AllAlbums] Error: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var alb Album

		if err := rows.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
			return albums, fmt.Errorf("[AllAlbums] Error: %v", err)
		}

		albums = append(albums, alb)
	}

	if err := rows.Err(); err != nil {
		return albums, fmt.Errorf("[AllAlbums] Error: %v", err)
	}

	return albums, nil
}

func (s *service) AlbumById(id int) (Album, error) {
	var alb Album

	q := `
		SELECT id, title, artist, price FROM album WHERE id = $1 
	`
	row := s.db.QueryRow(q, id)

	if err := row.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
		if err == sql.ErrNoRows {
			// The error of no matching Album should be handled by caller of this method
			return alb, err
		}
		return alb, fmt.Errorf("[AlbumById] Error: %v", err)
	}

	return alb, nil
}

func (s *service) AddAlbum(alb Album) (int64, error) {
	var rowsEffected int64

	q := `
		INSERT INTO album (title, artist, price) 
			VALUES ($1, $2, $3)
	`

	res, err := s.db.Exec(q, alb.Title, alb.Artist, alb.Price)
	if err != nil {
		return rowsEffected, fmt.Errorf("[AddAlbum] Error: %v", err)
	}

	rowsEffected, err = res.RowsAffected()
	if err != nil {
		return rowsEffected, fmt.Errorf("[AlbumById] Error: %v", err)
	}

	return rowsEffected, nil
}

func (s *service) DeleteAlbumByID(id int) (int, error) {
	q := `
		DELETE FROM album WHERE id = $1 
	`

	_, err := s.db.Exec(q, id)
	if err != nil {
		return id, fmt.Errorf("[DeleteAlbumByID] Error: %v", err)
	}

	return id, nil
}

// Close closes the database connection.
// It logs a message indicating the disconnection from the specific database.
// If the connection is successfully closed, it returns nil.
// If an error occurs while closing the connection, it returns the error.
func (s *service) Close() error {
	log.Printf("Disconnected from database: %s", database)
	return s.db.Close()
}
