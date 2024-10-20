package metricpgrepo

import "github.com/Masterminds/squirrel"

const schemaDDL = `create table if not exists metric (
    id serial,
    key text unique,
    type text,
    name text,
    value text,
    datetime time
);`

const (
	metricTable = "metric"
)

const (
	idColumn       = "id"
	keyColumn      = "key"
	typeColumn     = "type"
	nameColumn     = "name"
	valueColumn    = "value"
	datetimeColumn = "datetime"
)

var psql = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

var baseSelectQuery = psql.Select(
	idColumn,
	typeColumn,
	nameColumn,
	valueColumn,
	datetimeColumn,
).From(metricTable)
