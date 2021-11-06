package utils

import "fmt"

type Builder struct {
	withString   string
	withParams   []interface{}
	selectString string
	selectParams []interface{}
	order        string
	sqlString    string
	params       []interface{}
	groupString  string
	limit        int
	offset       int
}

func NewBuilder() *Builder {
	return &Builder{}
}

func (b *Builder) Order(statement string) {
	if len(statement) == 0 {
		b.order = ""

		return
	}

	b.order = "ORDER BY " + statement
}

func (b *Builder) Add(statement string, args ...interface{}) {
	b.sqlString += " " + statement
	b.params = append(b.params, args...)
}

func (b *Builder) Select(statement string, args ...interface{}) {
	b.selectString = statement

	b.selectParams = append(make([]interface{}, 0), args...)
}

func (b *Builder) GroupBy(statement string) {
	b.groupString = statement
}

func (b *Builder) With(statement string, args ...interface{}) {
	b.withString = statement

	b.withParams = append(make([]interface{}, 0), args...)
}

func (b *Builder) RawSQL() string {
	var sql string

	if b.withString != "" {
		sql += " WITH " + b.withString
	}

	if b.selectString != "" {
		sql += " SELECT " + b.selectString
	}

	sql += " " + b.sqlString

	sql += " " + b.order

	if b.groupString != "" {
		sql += " GROUP BY " + b.selectString
	}

	if b.limit > 0 {
		sql += fmt.Sprintf(" LIMIT %d", b.limit)
	}

	if b.offset > 0 {
		sql += fmt.Sprintf(" OFFSET %d", b.offset)
	}

	return sql
}

func (b *Builder) Values() []interface{} {
	res := make([]interface{}, 0)
	res = append(res, b.selectParams...)
	res = append(res, b.params...)

	return res
}

func (b *Builder) Pagination(limit int, offset int) {
	if limit >= 0 {
		b.limit = limit
	}

	if offset >= 0 {
		b.offset = offset
	}
}
