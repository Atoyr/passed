package orm

type Sqlconnection struct {
	database string
	instance string
	server   string
	user     string
	password string
}

func NewConn(database string, instance string, server string, user string, password string) Sqlconnection {
	conn := Sqlconnection{}
	conn.database = database
	conn.instance = instance
	conn.server = server
	conn.user = user
	conn.password = password

	return conn
}

func (sc *Sqlconnection) Connectionstring() string {
	var ret = make([]byte, 0, 512)
	ret = append(ret, "server="...)
	ret = append(ret, sc.server...)
	if sc.instance != "" {
		ret = append(ret, "\\"...)
		ret = append(ret, sc.instance...)
	}
	ret = append(ret, ";user id="...)
	ret = append(ret, sc.user...)
	ret = append(ret, ";password="...)
	ret = append(ret, sc.password...)
	ret = append(ret, ";database="...)
	ret = append(ret, sc.database...)
	return string(ret)
}
