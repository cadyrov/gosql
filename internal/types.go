package internal

type Column struct {
	Nullable     bool    `json:"nullable"`
	IsArray      bool    `json:"is_array"`
	IsPrimaryKey bool    `json:"is_primary_key"`
	SQLComment   *string `json:"sql_comment"`
	SQLName      string  `json:"sql_name"`
	SQLDataType  string  `json:"sql_data_type"`
	GoType       string  `json:"go_type"`
	GoName       string  `json:"go_name"`
}

type Table struct {
	SQLName string   `json:"sql_name"`
	GoName  string   `json:"go_name"`
	Scheme  string   `json:"scheme"`
	Columns []Column `json:"columns"`
}
