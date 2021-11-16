package orm

import "strings"

type tags struct {
	Name      string
	Type      string
	Primary   bool
	Nullable  bool
	Indexname string
}

func newTags(tag string) tags {
	t := tags{
		Primary:  false,
		Nullable: false,
	}

	for _, s := range strings.Split(tag, ";") {
		if s != "" {
			kv := strings.Split(s, "=")
			if len(kv) == 2 {
				t.set(kv[0], kv[1])
			} else {
				t.set(kv[0], "")
			}
		}
	}

	return t
}

func (t *tags) set(key, value string) {
	k := strings.ToLower(strings.TrimSpace(key))
	v := strings.TrimSpace(value)

	if k == "" {
		return
	}

	switch k {
	case "name":
		t.Name = v
	case "type":
		t.Type = typestr(k)
	case "primary":
		t.Primary = true
	case "nullable":
		t.Nullable = true

	default:
		return
	}
}

func typestr(value string) string {
	split := strings.Split(value, "(")
	t := strings.ToLower(split[0])
	ret := ""

	switch t {
	case "int":
		ret = t
	case "nvarchar":
		ret = "nvarchar"
	default:
		break
	}

	return ret
}
