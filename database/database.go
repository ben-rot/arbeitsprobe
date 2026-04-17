package database

import (
	"embed"
	_ "embed"
	"fmt"
	"log"
	"time"

	"cwl/config"

	"github.com/jmoiron/sqlx"
	"github.com/qustavo/dotsql"
	_ "github.com/lib/pq"
)

//go:embed schema.sql
var schemaSQL string

//go:embed queries/*.sql
var queriesDir embed.FS



func loadQueries(fileName string) *dotsql.DotSql {

	path := "queries/" + fileName + ".sql"

	data, err := queriesDir.ReadFile(path)
	if err != nil {
		panic(fmt.Sprintf("CRITICAL: Failed to read embedded SQL file %s: %v", path, err))
	}

	dot, err := dotsql.LoadFromString(string(data))
    if err != nil {
        panic(fmt.Sprintf("CRITICAL: Failed to parse SQL syntax in %s: %v", path, err))
    }

    return dot
}



type Store struct {
	Users 		*UserRepo
	Families 	*FamilyRepo
	Clans 		*ClanRepo
	Accounts 	*AccountRepo
	DB 			*sqlx.DB
}



func NewStore(db *sqlx.DB) *Store {

	repo := func(name string) repository {
		return repository{db: db, dot: loadQueries(name)}
	}

	return &Store{
		Users: 		&UserRepo{repo("users")},
		Families: 	&FamilyRepo{repo("families")},
		Clans: 		&ClanRepo{repo("clans")},
		Accounts: 	&AccountRepo{repo("accounts")},
		DB: 		db,
	}
}




func ConnectDb(dbCfg *config.PostgresConfig) (*sqlx.DB, error) {

	connStr := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", 
		dbCfg.Host, dbCfg.Port, dbCfg.User, dbCfg.Password, dbCfg.Database,
	)


	db, err := sqlx.Open("postgres", connStr)
	if err != nil {
        return nil, fmt.Errorf("failed to open database: %w", err)
    }

	db.SetMaxOpenConns(25)
    db.SetMaxIdleConns(25)
    db.SetConnMaxLifetime(5 * time.Minute)

	if err = db.Ping(); err != nil {
        return nil, fmt.Errorf("failed to ping database: %w", err)
    }

	db.MustExec(schemaSQL)
	

	log.Println("Successfully connected to database")
	return db, nil
}