package orm

import "fmt"

type WhereComparisonType int

const (
	Equal WhereComparisonType = iota
	NotEqual
	MoreThan
	LessThan
	MoreThanOrEqual
	LessThanOrEqual
	IsNull
	IsNotNull
)

func (wct WhereComparisonType) createWherePhrase(value string, index int) (string, int) {
	s := fmt.Sprintf(" %s %s ", value, wct.String())
	if wct != IsNull && wct != IsNotNull {
		s = fmt.Sprintf("%s@p%d", s, index)
		index = index + 1
	}
	return s, index
}

func (wct WhereComparisonType) String() string {
	switch wct {
	case Equal:
		return " = "
	case NotEqual:
		return " != "
	case MoreThan:
		return " > "
	case LessThan:
		return " < "
	case MoreThanOrEqual:
		return " >= "
	case LessThanOrEqual:
		return " <= "
	case IsNull:
		return " is NULL "
	case IsNotNull:
		return " in not NULL "
	default:
		return ""
	}
}

type WherePhrase struct {
	Type  WhereComparisonType
	Key   string
	Value interface{}
}

type Where []WherePhrase

func (wps *Where) Append(wct WhereComparisonType, key string, value interface{}) {
	wp := WherePhrase{
		Type:  wct,
		Key:   key,
		Value: value,
	}
	*wps = append(*wps, wp)
}

func (wps *Where) CreateWherePhrase(startIndex int) (string, []interface{}) {
	var b = make([]byte, 0, 512)
	var ifs = make([]interface{}, 0, len(*wps))

	index := startIndex
	if index < 1 {
		index = 1
	}
	if len(*wps) > 0 {
		b = append(b, " where "...)
	}

	for _, v := range *wps {
		var s string
		s, index = v.Type.createWherePhrase(v.Key, index)
		b = append(b, s...)
		ifs = append(ifs, v.Value)
	}
	return string(b), ifs
}
