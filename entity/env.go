package entity

import "github.com/charmbracelet/bubbles/textinput"

type Env struct {
	Name string
	Vars map[string]string
}

type EnvTableRow struct {
	Key   textinput.Model
	Value textinput.Model
}
