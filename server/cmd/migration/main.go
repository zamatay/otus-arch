package main

import (
	"context"
	"flag"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"os"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	_ "github.com/jackc/pgx/v5/stdlib"
)

var (
	flags = flag.NewFlagSet("goose", flag.ExitOnError)
)

type Person struct {
	Login     string    `fake:"{username}"`
	FirstName string    `fake:"{firstname}"` // Any available function all lowercase
	LastName  string    `fake:"{lastname}"`  // Any available function all lowercase
	Birthday  time.Time `fake:"-"`
	GenderId  string    `fake:"{gender}"` // Can call with parameters
	Interests []string  `fake:"{hobby}"`
	City      string    `fake:"{city}"` // Comma separated for multiple values
}

func main() {
	flags.Parse(os.Args[1:])
	args := flags.Args()

	if len(args) == 0 {
		flags.Usage()
		return
	}
	var dbstring string
	dbstring = args[0]

	//db, err := goose.OpenDBWithDriver(driver, dbstring)
	ctx := context.Background()
	db, err := pgxpool.New(ctx, dbstring)
	if err != nil {
		log.Fatalf("goose: failed to open DB: %v\n", err)
	}

	if _, err := db.Exec(ctx, "delete from public.users"); err != nil {
		log.Fatalf("goose: failed to truncate users: %v\n", err)
	}

	defer func() {
		db.Close()
		//if err := db.Close(); err != nil {
		//	log.Fatalf("goose: failed to close DB: %v\n", err)
		//}
	}()
	tx, err := db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return
	}

	startDate := time.Date(1950, 1, 1, 0, 0, 0, 0, time.UTC) // Начальная дата
	endDate := time.Date(2000, 12, 31, 0, 0, 0, 0, time.UTC) // Конечная дата
	// Создание экземпляра структуры

	var countColumn = 10
	var countData = 1_000_000
	type rowStruct []interface{}
	bulk := make([][]any, 0, countData)
	for i := range countData {
		var person Person
		row := make(rowStruct, countColumn)
		gofakeit.Seed(i)

		err := gofakeit.Struct(&person)
		if err != nil {
			return
		}

		idx := 0
		row[idx] = true
		row[idx+1] = person.Login
		row[idx+2] = person.FirstName
		row[idx+3] = person.LastName
		row[idx+4] = gofakeit.DateRange(startDate, endDate)
		row[idx+5] = GetGender(person.GenderId)
		row[idx+6] = person.Interests
		row[idx+7] = person.City
		row[idx+8] = time.Now()
		row[idx+9] = time.Now()
		bulk = append(bulk, row)

	}

	from, err := tx.CopyFrom(ctx, pgx.Identifier{"public", "users"}, []string{"enabled", "login", "first_name", "last_name", "birthday", "gender_id", "interests", "city", "created_at", "updated_at"}, pgx.CopyFromRows(bulk))
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Вставили %d строк", from)

	if err := tx.Commit(ctx); err != nil {
		log.Fatalf("goose: failed to commit transaction: %v\n", err)
	}
}

func GetGender(id string) int {
	if id == "male" {
		return 1
	}
	return 2
}
