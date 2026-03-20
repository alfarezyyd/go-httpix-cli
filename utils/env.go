package utils

import "go-httpix-cli/entity"

func EnvToTableRows(e entity.Env) []entity.EnvTableRow {
	rows := make([]entity.EnvTableRow, 0, len(e.Vars))
	for k, v := range e.Vars {
		row := NewEnvTableRow()
		row.Key.SetValue(k)
		row.Value.SetValue(v)
		rows = append(rows, row)
	}
	return rows
}
