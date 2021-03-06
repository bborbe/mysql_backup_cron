package dumper

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"text/template"

	"github.com/bborbe/io/util"
	"github.com/bborbe/mysql-backup/model"
	"github.com/golang/glog"
)

type backup struct {
	name            model.Name
	host            model.MysqlHost
	port            model.MysqlPort
	user            model.MysqlUser
	pass            model.MysqlPassword
	targetDirectory model.TargetDirectory
}

func New(
	name model.Name,
	host model.MysqlHost,
	port model.MysqlPort,
	user model.MysqlUser,
	pass model.MysqlPassword,
	targetDirectory model.TargetDirectory,
) *backup {
	b := new(backup)
	b.name = name
	b.host = host
	b.port = port
	b.user = user
	b.pass = pass
	b.targetDirectory = targetDirectory
	return b
}

func (b *backup) Database(
	database model.MysqlDatabase,
	backupfile model.BackupFilename,
) error {
	return b.backup(database.String(), database.String(), backupfile)
}

func (b *backup) All(
	backupfile model.BackupFilename,
) error {
	return b.backup(
		"all",
		"--all-databases",
		backupfile,
	)
}

func (b *backup) backup(
	name string,
	database string,
	backupfile model.BackupFilename,
) error {
	path, err := util.NormalizePath("~/.my.cnf")
	if err != nil {
		return err
	}
	if err := writeMyCnfFile(path, b.user, b.pass); err != nil {
		return err
	}
	glog.V(1).Infof("mysqldump started")
	if err := runCommand("mysqldump", b.targetDirectory, "--defaults-file="+path, "--lock-tables=false", "--user", b.user.String(), "--host", b.host.String(), "--port", b.port.String(), "--result-file", backupfile.String(), "--all-databases"); err != nil {
		glog.V(2).Infof("mysqldump failed, delete incomplete backup: %v", err)
		if err := backupfile.Delete(); err != nil {
			glog.Warningf("delete incomplete backup failed: %v", err)
		}
		return err
	}
	glog.V(1).Infof("mysqldump finshed")
	return nil
}

const myCnfTemplaate = `
[mysqldump]
user={{.User}}
password={{.Pass}}
max_allowed_packet=1G
`

func writeMyCnfFile(path string, user model.MysqlUser, pass model.MysqlPassword) error {
	var data struct {
		User model.MysqlUser
		Pass model.MysqlPassword
	}
	data.Pass = pass
	data.User = user
	return writeTemplate(path, myCnfTemplaate, data, false)
}

func runCommand(command string, cwd model.TargetDirectory, args ...string) error {
	debug := fmt.Sprintf("%s %s", command, strings.Join(args, " "))
	glog.V(2).Infof("execute %s", debug)
	cmd := exec.Command(command, args...)
	if glog.V(4) {
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout
	}
	if cwd != "" {
		cmd.Dir = cwd.String()
	}
	if err := cmd.Start(); err != nil {
		return err
	}
	glog.V(2).Infof("%s started", debug)
	if err := cmd.Wait(); err != nil {
		glog.Warningf("%s failed: %v", debug, err)
		return fmt.Errorf("%s failed: %v", debug, err)
	}
	glog.V(2).Infof("%s finished", command)
	return nil
}

func writeTemplate(path string, templateContent string, data interface{}, executable bool) error {
	content, err := generateTemplate(path, templateContent, data)
	if err != nil {
		return err
	}
	return writeFile(path, content, executable)
}

func generateTemplate(name string, templateContent string, data interface{}) ([]byte, error) {
	tmpl, err := template.New(name).Parse(templateContent)
	if err != nil {
		return nil, err
	}
	content := bytes.NewBufferString("")
	if err := tmpl.Execute(content, data); err != nil {
		return nil, err
	}
	return content.Bytes(), nil
}

func writeFile(path string, content []byte, executable bool) error {
	var perm os.FileMode
	if executable {
		perm = 0755
	} else {
		perm = 0644
	}
	return ioutil.WriteFile(path, content, perm)
}
