package datahandle

import (
	"errors"
	"fmt"
	"funcation/golog"
	"io/ioutil"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type MysqlSysbenchResult struct {
	//rojectName           string
	Command               string
	Starttime             int64  //压测开始时间
	Endtime               int64  //压测结束时间
	Thread                int    //线程数
	Abtime                int    //压测时长
	Abtype                string //OLTP 类型
	Tread                 int
	Twrite                int
	Tohter                int
	Ttotal                int
	Transaction           int
	Tquery                int
	Tignore_errors        int
	Treconnect            int
	Exmin                 string
	Exavg                 string
	Exmax                 string
	Ex_95_percent         string
	Exsum                 string
	Event_avg             string
	Event_stddev          string
	Execution_time_avg    string
	Execution_time_stddev string
}

/*
*  获取得到执行的命令的时候，获取其中的一些相关参数
 */

func (this *MysqlSysbenchResult) CommandMarksomeFlag(commad string) error {

	this.Command = commad
	for _, ok := range strings.Split(this.Command, "--") {
		if strings.Split(ok, "=")[0] == "threads" {
			ts := strings.Split(ok, "=")[1]
			ts = strings.Replace(ts, " ", "", -1)
			ts = strings.Replace(ts, "\n", "", -1)
			key, err := strconv.Atoi(ts)
			if err != nil {
				fmt.Printf("Cover  %s  Fail!", ts)
				return err
			}
			this.Thread = key
		}

		if strings.Split(ok, "=")[0] == "time" {
			ts := strings.Split(ok, "=")[1]
			ts = strings.Replace(ts, " ", "", -1)
			ts = strings.Replace(ts, "\n", "", -1)
			key, err := strconv.Atoi(ts)
			if err != nil {
				golog.Error("Datahandler", "datahandler", fmt.Sprintf("Cover %s Fail!", ts), 0)
				fmt.Printf("Cover  %s  Fail!", ts)
				return err
			}
			this.Abtime = key
		}

	}
	return nil
}

func (this *MysqlSysbenchResult) DealBytes(barry []byte) {
	var Btmp []uint8
	var BRes [][]uint8
	for _, v := range barry {
		if v != 10 {
			Btmp = append(Btmp, v)
		} else {
			BRes = append(BRes, Btmp)
			Btmp = []uint8{}
		}
	}

	for _, v := range BRes {
		vt := strings.Replace(string(v), " ", "", -1)
		vt = strings.Replace(vt, "\n", "", -1)
		Rowscuts := strings.Split(vt, ":")
		for _, mes := range Rowscuts {

			switch mes {
			case "read":
				this.Tread, _ = strconv.Atoi(Rowscuts[1])
			case "write":
				this.Twrite, _ = strconv.Atoi(Rowscuts[1])
			case "other":
				this.Tohter, _ = strconv.Atoi(Rowscuts[1])
			case "total":
				this.Ttotal, _ = strconv.Atoi(Rowscuts[1])
			case "transactions":
				this.Transaction, _ = strconv.Atoi(strings.Split(string(Rowscuts[1]), "(")[0])
			case "queries":
				this.Tquery, _ = strconv.Atoi(string(strings.Split(Rowscuts[1], "(")[0]))
			case "ignorederrors":
				this.Tignore_errors, _ = strconv.Atoi(strings.Split(Rowscuts[1], "(")[0])
			case "reconnects":
				this.Treconnect, _ = strconv.Atoi(strings.Split(Rowscuts[1], "(")[0])
			case "min":
				this.Exmin = FloatString(Rowscuts[1])
			case "avg":
				this.Exavg = FloatString(Rowscuts[1])
			case "max":
				this.Exmax = FloatString(Rowscuts[1])
			case "95thpercentile":
				this.Ex_95_percent = FloatString(Rowscuts[1])
			case "sum":
				this.Exsum = FloatString(Rowscuts[1])
			case "events(avg/stddev)":
				this.Event_avg = FloatString(strings.Split(Rowscuts[1], "/")[0])
				this.Event_stddev = FloatString(strings.Split(Rowscuts[1], "/")[1])
			case "executiontime(avg/stddev)":
				this.Execution_time_avg = FloatString(strings.Split(Rowscuts[1], "/")[0])
				this.Execution_time_stddev = FloatString(strings.Split(Rowscuts[1], "/")[1])
			}
		}
	}
}

func FloatString(a string) string {
	b, _ := strconv.ParseFloat(a, 32)
	return strconv.FormatFloat(b, 'f', 2, 64)

}

func NewMysqlsysbenchRes(command string) (*MysqlSysbenchResult, error) {
	h := new(MysqlSysbenchResult)
	h.Starttime = time.Now().Unix()
	h.Abtime = 10
	err := h.CommandMarksomeFlag(command)
	if err != nil {
		golog.Error("datahandler", "datahandler", fmt.Sprintf("%s", err), 0)
	}
	err = h.CommandExec()
	if err != nil {
		golog.Error("datahandle", "datahandle", fmt.Sprintf("%s", err), 0)
		return nil, err
	}
	h.Endtime = time.Now().Unix()
	return h, nil
}

func (this *MysqlSysbenchResult) CommandExec() error {
	cmd := exec.Command("sh", "-c", fmt.Sprintf("%s", this.Command))
	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()
	defer stderr.Close()
	defer stdout.Close()
	if err := cmd.Start(); err != nil {
		golog.Error("datahandle", "datahandle", "执行命令出错", 0)
		return err
	}

	opBytes, err := ioutil.ReadAll(stdout)
	opError, err := ioutil.ReadAll(stderr)
	if err != nil {
		golog.Error("datahandler", "datahandler", "Read stdout Fail", 0)
		return err
	}
	if len(opError) > 0 {
		return errors.New(string(opError))
		golog.Error("datahandler", "datahandler", string(opError), 0)
	}
	this.DealBytes(opBytes)
	return nil

}

/*func (this *MysqlSysbenchResult) RunAbtest() {
	Rres, err := CommandExecResultBytes(this.Command)
	if err != nil {
		fmt.Println("Execute command Fail!")
		os.Exit(2)
	}
	this.DealBytes(Rres)
	//DealBytes(Rres)
}*/
