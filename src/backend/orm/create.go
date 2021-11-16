package orm

import "reflect"

type CreateConfig struct {
	Config
	ExistsAction func()
	StopIfExists bool
}

type Column struct {
	Name    string
	Type    string
	Optioin string
}

type CreateOption func(c *CreateConfig)

func (db *DB) Create(model interface{}, options ...CreateOption) error {
	t := reflect.TypeOf(model)
	config := &CreateConfig{
		Config:       Config{TableName: t.Name(), Scheme: "dbo"},
		ExistsAction: func() {},
		StopIfExists: true,
	}

	for _, option := range options {
		option(config)
	}

	query := "SELECT count(*) FROM sys.object WHERE object_id = OBJECT_ID(%p1) AND type in (N'U')"
	count := 0
	err := db.db.QueryRow(query, config.TableName).Scan(&count)
	if err != nil {
		return err
	}

	if count > 0 {
		config.ExistsAction()
		if config.StopIfExists {
			return nil
		}
	}

	buffer := make([]byte, 4192)

	buffer = append(buffer, "CREATE TABLE "...)
	buffer = append(buffer, ("[" + config.Scheme + "].[" + config.TableName + "] ")...)

	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		newTags(f.Tag.Get("mssql-orm"))

	}

	return nil
}
