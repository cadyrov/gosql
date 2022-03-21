package main

import (
	"context"
	"fmt"
	"os"

	"github.com/cadyrov/gosql/internal"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jessevdk/go-flags"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func main() {
	cnf := internal.Config{}

	parser := flags.NewParser(&cnf, flags.Default)
	if _, err := parser.Parse(); err != nil {
		fmt.Printf("error parse env: %s\n", err.Error())

		os.Exit(1)
	}

	if cnf.Mode == internal.ModeMySQL {
		runMySQL(cnf)
	}

	if cnf.Mode == internal.ModePSQL {
		runPSQL(cnf)
	}
}

func runPSQL(cnf internal.Config) {
	url := "host=%s port=%d user=%s password=%s dbname=%s"

	if cnf.PSQL.SslMode != "" {
		url += " sslmode=" + cnf.PSQL.SslMode
	}

	if cnf.PSQL.Binary {
		url += " binary_parameters=yes"
	}

	dsn := fmt.Sprintf(url,
		cnf.SQL.Host, cnf.SQL.Port, cnf.SQL.UserName, cnf.SQL.Password, cnf.SQL.DBName)

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		panic(err)
	}

	gen := internal.NewGenerator(db, cnf.Scheme, cnf.Table, cnf.Template, cnf.Result)

	if err := gen.PSQL(context.Background()); err != nil {
		panic(err)
	}
}

func runMySQL(cnf internal.Config) {
	dsn := fmt.Sprintf(
		`%s:%s@tcp(%s:%d)/%s?parseTime=true`,
		cnf.SQL.UserName, cnf.SQL.Password, cnf.SQL.Host, cnf.SQL.Port, cnf.SQL.DBName)

	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		panic(err)
	}

	gen := internal.NewGenerator(db, cnf.Scheme, cnf.Table, cnf.Template, cnf.Result)

	if err := gen.MySQL(context.Background()); err != nil {
		panic(err)
	}
}
