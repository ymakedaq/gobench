package context

import (
	"fmt"
	"funcation/commandhandle"
	"funcation/core"
	"funcation/datahandle"
	"funcation/golog"
	"lib/cfg"
	"os"
	"strings"
	"time"
)

type BaseContext struct {
	SysConf     string
	SysClean    bool
	SysPrepare  bool
	SysLogfile  string
	SysLoglevel string
	SysLogo     string
}

func (this *BaseContext) Start() {
	logfile_time, err := golog.NewTimeRotatingFileHandler(this.SysLogfile, 3, 1)
	if err != nil {
		fmt.Println("Init Log Fail!", err)
		return
	}
	golog.GlobalSysLogger = golog.New(logfile_time, golog.Lfile|golog.Ltime|golog.Llevel)
	defer golog.GlobalSysLogger.Close()
	this.setLogLevel()
	gcmd := cfg.New_Gbh_cfg(this.SysConf)
	gcmd.Init_self()
	gblist := NewMysqlbenchResFromcfg(gcmd)
	go core.InitCpu()
	time.Sleep(1 * time.Second)
	go core.Collect()
	//core.DealChain()
	for map_index, _ := range gblist {
		//benchwork(map_index, gblist[map_index])
		this.benchwork(map_index, gblist[map_index])
		commandhandle.CommandExecResultBytes("> " + core.CPUFILE)
		commandhandle.CommandExecResultBytes("> " + core.MEMFILE)
	}
}

func (this *BaseContext) setLogLevel() {
	switch strings.ToLower(this.SysLoglevel) {
	case "debug":
		golog.GlobalSysLogger.SetLevel(golog.LevelDebug)
	case "info":
		golog.GlobalSysLogger.SetLevel(golog.LevelInfo)
	case "warn":
		golog.GlobalSysLogger.SetLevel(golog.LevelWarn)
	case "error":
		golog.GlobalSysLogger.SetLevel(golog.LevelError)
	default:
		golog.GlobalSysLogger.SetLevel(golog.LevelError)
	}
}

func (this *BaseContext) benchwork(res_flname string, cdl []*datahandle.MysqlSysbenchResult) {
	if len(cdl) <= 0 {
		fmt.Println("No Command list~")
		return
	}
	resfd, _ := golog.NewFileHandler(res_flname+".txt", os.O_CREATE|os.O_RDWR|os.O_APPEND)
	for index, v := range cdl {
		fmt.Println(index, "starting...")
		err := v.SysbenchRun()
		if err != nil {
			fmt.Println(err)
			return
		}
		if this.SysClean || this.SysPrepare {
			fmt.Println(index, "cleanup...")
			v.CommandCleanup(v.Command)
			v.CommandPrepare(v.Command)
		}
		resfd.Write([]byte(fmt.Sprintln(*v)))
		fmt.Println(*v)
	}
	resfd.Close()

	var mk datahandle.Dreawhtml
	mk.Newchart(res_flname)
}

func NewMysqlbenchResFromcfg(c *cfg.Gbh_cfg) map[string][]*datahandle.MysqlSysbenchResult {
	var rt map[string][]*datahandle.MysqlSysbenchResult
	rt = make(map[string][]*datahandle.MysqlSysbenchResult)
	for _, v := range c.Servers {
		for index, cmd := range v.Cmd_list {
			var t []*datahandle.MysqlSysbenchResult
			list_cmd := []string{}
			rt_name := v.Server_name + "_cmd" + fmt.Sprintf("%d", index+1)
			list_cmd = strings.Split(cmd, " ")[:1]
			list_cmd = append(list_cmd, "--mysql-host="+v.Mysql_host)
			list_cmd = append(list_cmd, "--mysql-user="+v.Mysql_user)
			list_cmd = append(list_cmd, "--mysql-password="+v.Mysql_password)
			list_cmd = append(list_cmd, "--mysql-db="+v.Mysql_db)
			list_cmd = append(list_cmd, "--mysql-port="+fmt.Sprintf("%d", v.Mysql_port))
			list_cmd = append(list_cmd, "--time="+fmt.Sprintf("%d", v.Bench_time))
			list_cmd = append(list_cmd, "--db-driver="+v.DB_Driver)
			list_cmd = append(list_cmd, "--tables="+fmt.Sprintf("%d", v.Tables))
			list_cmd = append(list_cmd, "--table_size="+fmt.Sprintf("%d", v.Table_size))
			for _, thread := range c.Thread_list {
				lt := new(datahandle.MysqlSysbenchResult)
				lt.Abtime = v.Bench_time
				lt.Thread = thread
				tmp_command := append(list_cmd, " --threads="+fmt.Sprintf("%d", thread))
				for _, other := range strings.Split(cmd, " ")[1:] {
					tmp_command = append(tmp_command, other)
				}
				lt.Command = strings.Join(tmp_command, " ")
				t = append(t, lt)
			}
			rt[rt_name] = t
		}
	}
	return rt
}
