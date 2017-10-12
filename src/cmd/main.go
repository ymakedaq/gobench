// work project main.go
package main

import (
	"flag"
	"fmt"
	"funcation/commandhandle"
	"funcation/core"
	"funcation/datahandle"
	"funcation/golog"
	"lib/cfg"
	"os"
	"path"
	"strings"
	"time"
)

func main() {

	var cdl []string
	var gblist map[string][]string
	gblist = make(map[string][]string)

	LOGO := `
		压测线程数  5, 10, 15, 20, 40, 60, 80, 100, 150, 200, 300, 400, 600, 800, 1000 递进压测
	`
	fmt.Println(LOGO)
	logFilePath := "/tmp/"
	logFile := path.Join(logFilePath, core.LogFilename)
	logfile_time, _ := golog.NewTimeRotatingFileHandler(logFile, 3, 1)
	golog.GlobalSysLogger = golog.New(logfile_time, golog.Lfile|golog.Ltime|golog.Llevel)
	setLogLevel(core.Loglevel)
	defer golog.GlobalSysLogger.Close()
	cmd := flag.String("c", "", "--sysbench command")
	conf_file := flag.String("f", "", "--Conf file")
	flag.Parse()

	if len(*cmd) >= 1 {
		cdl = Newcommand(*cmd)
		gblist["current"] = cdl
	} else {
		gcmd := cfg.New_Gbh_cfg(*conf_file)
		gcmd.Init_self()
		gblist = NewcommandFromcfg(gcmd)
	}
	go core.InitCpu()
	time.Sleep(1 * time.Second)
	go core.Collect()
	core.DealChain()

	for map_index, _ := range gblist {
		benchwork(map_index, gblist[map_index])
		commandhandle.CommandExecResultBytes("> " + core.CPUFILE)
		commandhandle.CommandExecResultBytes("> " + core.MEMFILE)
	}

}

func init() {
	CheckInstallsysbench()
}

func Newcommand(cmd string) []string {
	if len(cmd) <= 0 {
		return []string{}
	}
	step := []int{6, 10, 15, 20, 40, 60, 80, 100, 150, 200, 300, 400, 600, 800, 1000}
	fmt.Println("Fllow threads", step)
	var cmdlist []string
	for _, v := range step {
		cmdlist = append(cmdlist, strings.Replace(cmd, "--threads=2", fmt.Sprintf("--threads=%d", v), 1))
	}

	return cmdlist
}

func NewcommandFromcfg(c *cfg.Gbh_cfg) map[string][]string {
	var rt map[string][]string
	rt = make(map[string][]string)
	for _, v := range c.Servers {
		for index, cmd := range v.Cmd_list {
			var t []string
			rt_name := v.Server_name + "_cmd" + fmt.Sprintf("%d", index+1)
			list_cmd := strings.Split(cmd, " ")
			list_cmd[1] = "--mysql-host=" + v.Mysql_host
			list_cmd[2] = "--mysql-user=" + v.Mysql_user
			list_cmd[3] = "--mysql-password=" + v.Mysql_password
			list_cmd[4] = "--mysql-db=" + v.Mysql_db
			list_cmd[5] = "--mysql-port=" + fmt.Sprintf("%d", v.Mysql_port)
			list_cmd[6] = "--time=" + fmt.Sprintf("%d", v.Bench_time)
			list_cmd[7] = "--db-driver=" + v.DB_Driver
			for _, thread := range c.Thread_list {
				list_cmd[8] = " --threads=" + fmt.Sprintf("%d", thread)
				t = append(t, strings.Join(list_cmd, " "))
			}
			rt[rt_name] = t
		}
	}
	return rt
}

func benchwork(res_flname string, cdl []string) {
	if len(cdl) <= 0 {
		fmt.Println("No Command list~")
		return
	}
	resfd, _ := golog.NewFileHandler(res_flname+".txt", os.O_CREATE|os.O_RDWR|os.O_APPEND)
	for index, v := range cdl {
		fmt.Println(index, "--start--")
		fmt.Println(v)
		res, err := datahandle.NewMysqlsysbenchRes(`"` + v + `"`)
		if err != nil {
			fmt.Println(err)
			return
		}
		resfd.Write([]byte(fmt.Sprintln(*res)))
		fmt.Println(res)
	}
	resfd.Close()

	var mk datahandle.Dreawhtml
	mk.Newchart(res_flname)
}

func CheckInstallsysbench() {

	var wins string
	res, err := commandhandle.CommandExecResultBytes(`rpm -qa|grep  sysbench`)
	if len(res) == 0 || err != nil {
		fmt.Println("Not Found sysbench!,You Need Install? (Y|N)")
		fmt.Scanln(&wins)
		if wins == "Y" || wins == "y" {
			fmt.Println("Install ...")
			res, err := commandhandle.CommandExecResultBytes("yum -y install make automake libtool pkgconfig libaio-devel vim-common sudo -y && curl -s https://packagecloud.io/install/repositories/akopytov/sysbench/script.rpm.sh | sudo bash " +
				"&& yum-config-manager --save --setopt=akopytov_sysbench.skip_if_unavailable=true  && yum install  sysbench -y")
			fmt.Println(string(res))
			if err != nil {
				fmt.Println("Install Sysbench Fail!!")
				fmt.Println(time.Now().Format("2006-01-02 15:04:05"), err)
				return

			} else {
				fmt.Println("Installc Success!!")
				golog.Info("main", "main", "Installc Success!!", 0)
				///fmt.Println("Install  Success!!)
			}
		} else {
			fmt.Println("Chose Nothing!")
			os.Exit(1)
		}
	}
	fmt.Println(string(res) + "been  Install!")
}

func setLogLevel(level string) {
	switch strings.ToLower(level) {
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
