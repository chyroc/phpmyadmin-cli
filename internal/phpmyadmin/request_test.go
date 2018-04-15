package phpmyadmin

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Chyroc/phpmyadmin-cli/internal/requests"
	"github.com/Chyroc/phpmyadmin-cli/internal/show"
)

func initClient() *phpMyAdmin {
	s, err := requests.NewSession()
	if err != nil {
		panic(err)
	}
	return &phpMyAdmin{
		session: s,
	}
}

func TestLogin(t *testing.T) {
	as := assert.New(t)
	p := initClient()
	p.SetURI("localhost:8000")

	as.Nil(p.Login("root", "pass"))
	err := p.Login("root", "error")
	as.NotNil(err)
	as.Equal("login err", err.Error())
}

func TestShowDatabases(t *testing.T) {
	p := initClient()
	as := assert.New(t)
	p.SetURI("localhost:8000")

	as.Nil(p.Login("root", "pass"))

	b, err := p.ExecSQL("root", "", "", "show databases;")
	as.Nil(err)
	as.NotNil(b)

	var buf = new(bytes.Buffer)
	show.TestSetOut(buf)
	show.ParseFromHTML(string(b))
	as.Equal(`+--------------------+
|      Database      |
+--------------------+
| information_schema |
| mysql              |
| performance_schema |
| sys                |
+--------------------+
`, buf.String())
}

func TestShowTables(t *testing.T) {
	p := initClient()
	as := assert.New(t)
	p.SetURI("localhost:8000")

	as.Nil(p.Login("root", "pass"))

	b, err := p.ExecSQL("root", "information_schema", "", "show tables;")
	as.Nil(err)
	as.NotNil(b)

	var buf = new(bytes.Buffer)
	show.TestSetOut(buf)
	show.ParseFromHTML(string(b))
	as.Equal(`+---------------------------------------+
|     Tables_in_information_schema      |
+---------------------------------------+
| CHARACTER_SETS                        |
| COLLATIONS                            |
| COLLATION_CHARACTER_SET_APPLICABILITY |
| COLUMNS                               |
| COLUMN_PRIVILEGES                     |
| ENGINES                               |
| EVENTS                                |
| FILES                                 |
| GLOBAL_STATUS                         |
| GLOBAL_VARIABLES                      |
| KEY_COLUMN_USAGE                      |
| OPTIMIZER_TRACE                       |
| PARAMETERS                            |
| PARTITIONS                            |
| PLUGINS                               |
| PROCESSLIST                           |
| PROFILING                             |
| REFERENTIAL_CONSTRAINTS               |
| ROUTINES                              |
| SCHEMATA                              |
| SCHEMA_PRIVILEGES                     |
| SESSION_STATUS                        |
| SESSION_VARIABLES                     |
| STATISTICS                            |
| TABLES                                |
| TABLESPACES                           |
| TABLE_CONSTRAINTS                     |
| TABLE_PRIVILEGES                      |
| TRIGGERS                              |
| USER_PRIVILEGES                       |
| VIEWS                                 |
| INNODB_LOCKS                          |
| INNODB_TRX                            |
| INNODB_SYS_DATAFILES                  |
| INNODB_FT_CONFIG                      |
| INNODB_SYS_VIRTUAL                    |
| INNODB_CMP                            |
| INNODB_FT_BEING_DELETED               |
| INNODB_CMP_RESET                      |
| INNODB_CMP_PER_INDEX                  |
| INNODB_CMPMEM_RESET                   |
| INNODB_FT_DELETED                     |
| INNODB_BUFFER_PAGE_LRU                |
| INNODB_LOCK_WAITS                     |
| INNODB_TEMP_TABLE_INFO                |
| INNODB_SYS_INDEXES                    |
| INNODB_SYS_TABLES                     |
| INNODB_SYS_FIELDS                     |
| INNODB_CMP_PER_INDEX_RESET            |
| INNODB_BUFFER_PAGE                    |
| INNODB_FT_DEFAULT_STOPWORD            |
| INNODB_FT_INDEX_TABLE                 |
| INNODB_FT_INDEX_CACHE                 |
| INNODB_SYS_TABLESPACES                |
| INNODB_METRICS                        |
| INNODB_SYS_FOREIGN_COLS               |
| INNODB_CMPMEM                         |
| INNODB_BUFFER_POOL_STATS              |
| INNODB_SYS_COLUMNS                    |
| INNODB_SYS_FOREIGN                    |
| INNODB_SYS_TABLESTATS                 |
+---------------------------------------+
`, buf.String())
}
