// work project main.go
package main

import (
	"flag"
	"fmt"
	"funcation/commandhandle"
	"funcation/golog"
	"lib/context"
	"os"
	"time"
)

func main() {

	var BenchBase context.BaseContext
	conf_file := flag.String("f", "", "--Conf file")
	prepare_flag := flag.Bool("pr", false, "--Prepare for sysbench")
	clean_flag := flag.Bool("c", false, "--Cleanup for sysbench")
	flag.Parse()
	BenchBase.SysLogo = `########## START SYSBENCH #########`
	BenchBase.SysLogfile = "/tmp/gun.log"
	BenchBase.SysLoglevel = `info`
	BenchBase.SysPrepare = *prepare_flag
	BenchBase.SysClean = *clean_flag
	BenchBase.SysConf = *conf_file

	BenchBase.Start()
}

func init() {
	CheckInstallsysbench()
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
			}
		} else {
			fmt.Println("Chose Nothing!")
			os.Exit(1)
		}
	}
	fmt.Println(string(res) + "been  Install!")
}
