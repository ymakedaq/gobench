// work project main.go
package main

import (
	"flag"
	"fmt"
	"funcation/commandhandle"
	"funcation/core"
	"funcation/datahandle"
	"funcation/golog"
	"os"
	"path"
	"strings"
	"time"
)

func main() {

	//	var TheadWait sync.WaitGroup

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
	cmd := flag.String("c", "", "Please Input `sysbench  Commad`~")
	flag.Parse()
	if len(*cmd) <= 0 {
		fmt.Printf("Please Input Command!!")
		return
	}
	cdl := Newcommand(*cmd)

	go core.InitCpu()
	time.Sleep(1 * time.Second)
	go core.Collect()
	core.DealChain()

	resfd, _ := golog.NewFileHandler(core.RESFILE, os.O_CREATE|os.O_RDWR|os.O_APPEND)
	for index, v := range cdl {
		fmt.Println(index, "--start--", v)
		res, err := datahandle.NewMysqlsysbenchRes(v)
		if err != nil {
			fmt.Println(err)
			return
		}
		resfd.Write([]byte(fmt.Sprintln(*res)))
		fmt.Println(res)
	}
	resfd.Close()

	var mk datahandle.Dreawhtml
	mk.Newchart()

}

func init() {
	CheckInstallsysbench()
}

func Newcommand(cmd string) []string {

	step := []int{6, 10, 15, 20, 40, 60, 80, 100, 150, 200, 300, 400, 600, 800, 1000}
	fmt.Println("Fllow threads", step)
	var cmdlist []string
	for _, v := range step {
		cmdlist = append(cmdlist, strings.Replace(cmd, "--threads=2", fmt.Sprintf("--threads=%d", v), 1))
	}

	return cmdlist
}

func CheckInstallsysbench() {

	var wins string
	res, err := commandhandle.CommandExecResultBytes("rpm -qa|grep  sysbench")
	if len(res) == 0 || err != nil {
		fmt.Println("Not Found sysbench!,You Need Install? (Y|N)")
		fmt.Scanln(&wins)
		if wins == "Y" || wins == "y" {
			fmt.Println("Install ...")
			res, err := commandhandle.CommandExecResultBytes("yum install  sudo -y && curl -s https://packagecloud.io/install/repositories/akopytov/sysbench/script.rpm.sh | sudo bash "
			+"&& yum-config-manager --save --setopt=akopytov_sysbench.skip_if_unavailable=true  && yum install  sysbench -y")
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
