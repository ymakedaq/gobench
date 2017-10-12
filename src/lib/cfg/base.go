package cfg

import (
	"fmt"
	"funcation/golog"

	"funcation/common"
	"regexp"
	"strings"

	"github.com/robfig/config"
)

var validcmd = regexp.MustCompile(`cmd?[0-9]*[0-9]$`)

type Gbh_cfg struct {
	Cfg_name    string // 配置文件名称
	Thread_list []int  // 压测线程数
	Time_step   int    // cpu,mem 采样时间间隔
	Servers     []bench_service
}

type bench_service struct {
	Server_name    string
	Mysql_host     string
	Mysql_user     string
	Mysql_port     int
	Mysql_db       string
	Mysql_password string
	Tables         int
	Table_size     int
	Bench_time     int //压测时间 s
	DB_Driver      string
	Cmd_list       []string // 压测命令
}

func New_Gbh_cfg(filename string) *Gbh_cfg {
	k := new(Gbh_cfg)
	k.Cfg_name = filename
	return k
}

func (this *Gbh_cfg) Init_self() {
	cfg, err := config.ReadDefault(this.Cfg_name)
	if err != nil {
		golog.Error("cfg", "inits_serf", fmt.Sprintf("%s", err), 0)
		return
	}
	tl, err := cfg.RawStringDefault("thread_list")
	if err != nil {
		fmt.Println("thread_list 配置有误~")
		golog.Error("cfg", "inits_serf", fmt.Sprintf("%s", err), 0)
		return
	}
	tsp, err := cfg.Int("DEFAULT", "time_step")
	if err != nil {
		golog.Error("cfg", "inits_serf", fmt.Sprintf("%s", err), 0)
		return
	}
	thd_list, err := common.String_to_int(strings.Split(tl, ","))
	if err != nil {
		golog.Error("cfg", "inits_serf", fmt.Sprintf("%s", err), 0)
		return
	}
	this.Time_step = tsp
	this.Thread_list = thd_list
	sections := cfg.Sections()
	var t_servers []bench_service
	for _, s_name := range sections {
		if s_name != "DEFAULT" {
			var t_server bench_service
			options, err := cfg.Options(s_name)
			if err != nil {
				golog.Error("cfg", "inits_serf", fmt.Sprintf("%s", err), 0)
				return
			}
			t_server.Server_name = s_name
			for _, items := range options {
				switch {
				case items == "mysql-host":
					t_server.Mysql_host, _ = cfg.String(s_name, items)
				case items == "mysql-port":
					t_server.Mysql_port, _ = cfg.Int(s_name, items)
				case items == "mysql-user":
					t_server.Mysql_user, _ = cfg.String(s_name, items)
				case items == "mysql-password":
					t_server.Mysql_password, _ = cfg.String(s_name, items)
				case items == "mysql-db":
					t_server.Mysql_db, _ = cfg.String(s_name, items)
				case items == "db-driver":
					t_server.DB_Driver, _ = cfg.String(s_name, items)
				case items == "time":
					t_server.Bench_time, _ = cfg.Int(s_name, items)
				case items == "tables":
					t_server.Tables, _ = cfg.Int(s_name, items)
				case items == "table_size":
					t_server.Table_size, _ = cfg.Int(s_name, items)
				case validcmd.MatchString(items):
					cmd, _ := cfg.String(s_name, items)
					t_server.Cmd_list = append(t_server.Cmd_list, cmd)
				}
			}
			t_servers = append(t_servers, t_server)
		}
		this.Servers = t_servers
	}

}
