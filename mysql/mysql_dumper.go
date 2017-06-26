package mysql

import (
	"context"
	"fmt"

	"github.com/bborbe/mysql_backup_cron/backup"
	"github.com/bborbe/mysql_backup_cron/model"
)

type Dumper struct {
	Database        model.MysqlDatabase
	Host            model.MysqlHost
	Port            model.MysqlPort
	User            model.MysqlUser
	Password        model.MysqlPassword
	Name            model.Name
	TargetDirectory model.TargetDirectory
}

func NewDumper(
	database model.MysqlDatabase,
	host model.MysqlHost,
	port model.MysqlPort,
	user model.MysqlUser,
	password model.MysqlPassword,
	name model.Name,
	targetDirectory model.TargetDirectory,
) *Dumper {
	d := new(Dumper)
	d.Database = database
	d.Host = host
	d.Port = port
	d.User = user
	d.Password = password
	d.Name = name
	d.TargetDirectory = targetDirectory
	return d
}

func (m *Dumper) Validate() error {
	if len(m.Host) == 0 {
		return fmt.Errorf("mysql host missing")
	}
	if m.Port <= 0 {
		return fmt.Errorf("mysql port missing")
	}
	if len(m.User) == 0 {
		return fmt.Errorf("mysql user missing")
	}
	if len(m.Password) == 0 {
		return fmt.Errorf("mysql password missing")
	}
	if len(m.Database) == 0 {
		return fmt.Errorf("mysql database missing")
	}
	if len(m.TargetDirectory) == 0 {
		return fmt.Errorf("mysql target dir missing")
	}
	if len(m.Name) == 0 {
		return fmt.Errorf("mysql name missing")
	}
	return nil
}

func (m *Dumper) Run(ctx context.Context) error {
	return backup.Create(
		m.Name,
		m.Host,
		m.Port,
		m.User,
		m.Password,
		m.Database,
		m.TargetDirectory,
	)
}
