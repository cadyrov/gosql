package internal

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
)

const (
	MysqlPrimary  = "auto_increment"
	MysqlNullable = "YES"
)

func FromMySQL(ctx context.Context, db *sqlx.DB, dbname string, table string) (Table, error) {
	tb := Table{}

	sql := `SELECT COLUMN_NAME, DATA_TYPE, IS_NULLABLE, COLUMN_COMMENT, EXTRA
			FROM information_schema.columns
			WHERE table_schema = ?
			  AND table_name = ?`

	rws, err := db.QueryContext(ctx, sql, dbname, table)
	if err != nil {
		return tb, err
	}

	defer rws.Close()

	tb.SQLName = table
	tb.Scheme = dbname

	for rws.Next() {
		cl := Column{}
		extra := ""
		nbl := ""

		if err := rws.Scan(&cl.SQLName, &cl.SQLDataType, &nbl, &cl.SQLComment, &extra); err != nil {
			return tb, err
		}

		if nbl == MysqlNullable {
			cl.Nullable = true
		}

		if extra == MysqlPrimary {
			cl.IsPrimaryKey = true
		}

		tb.Columns = append(tb.Columns, cl)
	}

	return tb.mySQLPrepare(), nil
}

func mySQLTypes() map[string]string {
	mp := map[string]string{
		"int":      "int64",
		"varchar":  "string",
		"text":     "string",
		"enum":     "string",
		"datetime": "time.Time",
	}

	return mp
}

func (tb Table) mySQLPrepare() Table {
	tb.GoName = sqlToGo(tb.SQLName)

	mp := mySQLTypes()

	for i := range tb.Columns {
		tb.Columns[i].GoType = mp[tb.Columns[i].SQLDataType]

		if tb.Columns[i].GoType == "" {
			tb.Columns[i].GoType = "string"
		}

		if tb.Columns[i].Nullable {
			tb.Columns[i].GoType = fmt.Sprintf("*%s", tb.Columns[i].GoType)
		}

		tb.Columns[i].GoName = sqlToGo(tb.Columns[i].SQLName)
	}

	return tb
}
