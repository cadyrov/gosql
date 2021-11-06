package internal

import (
	"context"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
)

func FromPSQL(ctx context.Context, db *sqlx.DB, dbname string, table string) (Table, error) {
	tb := Table{}

	sql := `SELECT a.attname                                                           AS column_name,
       format_type(a.atttypid, a.atttypmod)                                            AS data_type,
       CASE WHEN a.attnotnull THEN FALSE ELSE TRUE END                                 AS is_nullable,
       (SELECT EXISTS(SELECT i.indisprimary
                      FROM pg_index i
                      WHERE i.indrelid = a.attrelid
                        AND a.attnum = ANY (i.indkey)
                        AND i.indisprimary IS TRUE))                                   AS is_primary,
       col_description(t.oid, ic.ordinal_position)                                     AS description
FROM pg_attribute a
         JOIN pg_class t ON a.attrelid = t.oid
         JOIN pg_namespace s ON t.relnamespace = s.oid
         LEFT JOIN information_schema.columns AS ic
                   ON ic.column_name = a.attname AND ic.table_name = t.relname AND ic.table_schema = s.nspname
         LEFT JOIN information_schema.key_column_usage AS kcu
                   ON kcu.table_name = t.relname AND kcu.column_name = a.attname
         LEFT JOIN information_schema.table_constraints AS tc
                   ON tc.constraint_name = kcu.constraint_name AND tc.constraint_type = 'FOREIGN KEY'
         LEFT JOIN information_schema.constraint_column_usage AS ccu ON ccu.constraint_name = tc.constraint_name
WHERE a.attnum > 0
  AND NOT a.attisdropped
  AND s.nspname = $1
  AND t.relname = $2
GROUP BY a.attname, a.atttypid, a.attrelid, a.atttypmod, a.attnotnull, s.nspname, t.relname, ic.column_default,
         ic.table_schema, ic.table_name, ic.column_name, a.attnum, t.oid, ic.ordinal_position
ORDER BY a.attnum;`

	rws, err := db.QueryContext(ctx, sql, dbname, table)
	if err != nil {
		return tb, err
	}

	defer rws.Close()

	tb.SQLName = table
	tb.Scheme = dbname

	for rws.Next() {
		cl := Column{}

		if err := rws.Scan(&cl.SQLName, &cl.SQLDataType, &cl.Nullable, &cl.IsPrimaryKey, &cl.SQLComment); err != nil {
			return tb, err
		}

		tb.Columns = append(tb.Columns, cl)
	}

	return tb.pSQLPrepare(), nil
}

func pSQLTypes() map[string]string {
	mp := map[string]string{
		"bigint":           "int64",
		"integer":          "int",
		"text":             "string",
		"double precision": "float64",
		"boolean":          "bool",
		"ARRAY":            "[]interface{}",
		"json":             "json.RawMessage",
		"smallint":         "int16",
		"date":             "time.Time",
		"uuid":             "string",
		"jsonb":            "json.RawMessage",
	}

	return mp
}

func pSQLArrTypes() map[string]string {
	mp := map[string]string{
		"uuid[]":    "[]string",
		"integer[]": "[]int64",
		"bigint[]":  "[]int64",
		"text[]":    "[]string",
	}

	return mp
}

func containedTypes(in string) string {
	switch {
	case strings.Contains(in, "timestamp"):
		return "time.Time"
	case strings.Contains(in, "numeric"):
		return "float32"
	default:
		return "string"
	}
}

func (tb Table) pSQLPrepare() Table {
	tb.GoName = sqlToGo(tb.SQLName)

	mp := pSQLTypes()
	arrMp := pSQLArrTypes()

	for i := range tb.Columns {
		if arrKey, ok := arrMp[tb.Columns[i].SQLDataType]; ok {
			tb.Columns[i].GoType = arrKey

			tb.Columns[i].IsArray = true
		} else {
			tb.Columns[i].GoType = mp[tb.Columns[i].SQLDataType]
		}

		if tb.Columns[i].GoType == "" {
			tb.Columns[i].GoType = containedTypes(tb.Columns[i].SQLDataType)
		}

		if tb.Columns[i].Nullable && !tb.Columns[i].IsArray {
			tb.Columns[i].GoType = fmt.Sprintf("*%s", tb.Columns[i].GoType)
		}

		tb.Columns[i].GoName = sqlToGo(tb.Columns[i].SQLName)
	}

	return tb
}
