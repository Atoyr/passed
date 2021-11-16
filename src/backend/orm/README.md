

# Tags
``` golang
type User struct {
  Name string `mssql-orm:"name=NAME;type=nvarchar(250);primary;"`
  SubName string `mssql-orm:"name=NAME;type=nvarchar(250);nullable;indexname=ix_foo"`
}
```


