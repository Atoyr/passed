package database

import "fmt"

type UpdatePhrase struct {
	Scheme       string
	Table        string
	ColumnValue  map[string]interface{}
	WherePhrases WherePhrases
}

func NewUpdatePhrase(scheme, table string) *UpdatePhrase {
	up := new(UpdatePhrase)

	up.Scheme = scheme
	if up.Scheme == "" {
		up.Scheme = "dbo"
	}

	up.Table = table
	return up
}

func (up *UpdatePhrase) Query() (string, []interface{}) {
	var b = make([]byte, 0, 512)
	cap := len(up.ColumnValue) + len(up.WherePhrases)
	var ifs = make([]interface{}, 0, cap)

	index := 1
	b = append(b, fmt.Sprintf("update [%s].[%s] ", up.Scheme, up.Table)...)
	b = append(b, " set "...)
	isFirst := true

	for k, v := range up.ColumnValue {
		if isFirst {
			isFirst = false
		} else {
			b = append(b, " , "...)
		}
		b = append(b, fmt.Sprintf("%s = @p%d ", k, index)...)
		ifs = append(ifs, v)
		index++
	}
	w, param := up.WherePhrases.CreateWherePhrase(index)
	b = append(b, w...)
	ifs = append(ifs, param...)

	return string(b), ifs
}
