package internal

import (
	"context"
	"io/ioutil"
	"os"
	"text/template"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type Generator struct {
	db           *sqlx.DB
	scheme       string
	table        string
	templatePath string
	resultPath   string
}

func NewGenerator(db *sqlx.DB, scheme string, table string, templatePath string, resultPath string) Generator {
	return Generator{
		db:           db,
		scheme:       scheme,
		table:        table,
		templatePath: templatePath,
		resultPath:   resultPath,
	}
}

func (g *Generator) PSQL(ctx context.Context) error {
	tbl, err := FromPSQL(ctx, g.db, g.scheme, g.table)
	if err != nil {
		return err
	}

	return g.generate(tbl)
}

func (g *Generator) MySQL(ctx context.Context) error {
	tbl, err := FromMySQL(ctx, g.db, g.scheme, g.table)
	if err != nil {
		return err
	}

	return g.generate(tbl)
}

func (g *Generator) generate(tbl Table) error {
	tmpl := template.New("model")

	templateFile, err := os.Open(g.templatePath)
	if err != nil {
		return errors.Wrap(err, "generate")
	}

	data, err := ioutil.ReadAll(templateFile)
	if err != nil {
		return errors.Wrap(err, "generate")
	}

	tmpl = template.Must(tmpl.Parse(string(data)))

	file, err := os.Create(g.resultPath)
	if err != nil {
		return errors.Wrap(err, "generate")
	}

	primaryColumns := []Column{}
	nonPrimaryColumns := []Column{}

	for i := range tbl.Columns {
		if tbl.Columns[i].IsPrimaryKey {
			primaryColumns = append(primaryColumns, tbl.Columns[i])

			continue
		}

		nonPrimaryColumns = append(nonPrimaryColumns, tbl.Columns[i])
	}

	// Parse template to file
	if err = tmpl.Execute(file, struct {
		Table             Table
		PrimaryColumns    []Column
		NonPrimaryColumns []Column
		Columns           []Column
	}{
		Table:             tbl,
		PrimaryColumns:    primaryColumns,
		NonPrimaryColumns: nonPrimaryColumns,
		Columns:           tbl.Columns,
	}); err != nil {
		return errors.Wrap(err, "generate")
	}

	if err = file.Close(); err != nil {
		return errors.Wrap(err, "generate")
	}

	return nil
}
