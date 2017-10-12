/*
type MysqlSysbenchResult struct {
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
*/

/*
import (
	"bufio"
	"bytes"
	//"encoding/json"
	"errors"
	"fmt"
	"funcation/golog"
	"html/template"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	CPUFILE    = `./cpufile.txt`
	MEMFILE    = `./memfile.txt`
	RESULTFILE = `./res.txt`
)

type Dreawhtml struct {
	Headtitle    string //html 的title
	Ytitle       string
	Xtitle       string
	Xdata        []string //x轴的数据值
	Ydata        []int    //y轴的数据值
	CpuHeadtitle string
	CpuYtitle    string
	CpuXtitle    string
	CpuXdata     []string
	CpuYdata     [][]int64
	MemHeadtitle string
	MemYtitle    string
	MemXtitle    string
	MemXdata     []string
	MemYdata     [][]int64
}

//渲染html  结构体的元素也需要首字母大写

func main() {
	var dat Dreawhtml
	tmpl := template.New("")
	tmpl.Parse(tpl)
	lines, err := Readfile("./res.txt")
	if err != nil {
		fmt.Println("ok")
	}
	res, err := SysbenchResCut(lines)
	if err != nil {
		fmt.Println(err)
		return
	}
	cpures, err := SysbenchCpucut(CPUFILE)
	if err != nil {
		fmt.Println(err)
		return
	}
	memres, err := SysbenchCpucut(MEMFILE)
	if err != nil {
		fmt.Println(err)
		return
	}
	dat.MemXdata = TransintTotime(memres[0])
	dat.MemYdata = memres[1:]
	dat.MemHeadtitle = "Memory"
	dat.CpuXdata = TransintTotime(cpures[0])
	dat.CpuYdata = cpures[1:]
	dat.CpuHeadtitle = "CPU idle"
	dat.CpuYtitle = "cpu idle"
	dat.Headtitle = res["headtitle"].(string)
	dat.Xtitle = "tps/s"
	dat.Ytitle = "tps/s"
	dat.Xdata = res["thread"].([]string)
	dat.Ydata = res["tps"].([]int)
	var html bytes.Buffer
	err = tmpl.Execute(&html, dat)
	fmt.Println(html.String())
	if err != nil {
		fmt.Println(err)
	}
	filname := "test.html"
	err = Createhmtl(filname, html.String())
	if err != nil {
		fmt.Println(err)
	}
	//list := []string{"dsasda", "xiaode", "xiaoke"}
	//tmpl.Execute(os.Stdout, list)
	//fmt.Println(doc.String())

}

func Readfile(f string) (*[]string, error) {
	rows := make([]string, 10)
	fd, err := os.OpenFile(f, os.O_RDONLY, 0660)
	if err != nil {
		fmt.Println("Openfile error:", err)
		return nil, err
	}
	defer fd.Close()
	rd := bufio.NewReader(fd)
	for {
		line, err := rd.ReadString('\n')
		if err != nil || io.EOF == err {
			break
		}
		rows = append(rows, line)

	}
	return &rows, nil
}

func SysbenchResCut(rows *[]string) (map[string]interface{}, error) {
	var pthread []string
	var ptps []int
	var pqps []int
	var pigrs []int
	mn := make(map[string]interface{})
	if len(*rows) < 0 {
		fmt.Println("NO DATA!")
		return mn, errors.New("no data!")
	}
	rrows := *rows
	b := strings.Split(strings.Split(rrows[len(*rows)-1], "run")[0], " ")
	c := strings.Split(b[len(b)-2], "/")
	mn["headtitle"] = c[len(c)-1]
	for _, v := range *rows {
		if len(v) > 1 {
			res := strings.Replace(strings.Split(v, "run")[1], "}", "", -1)
			res = strings.TrimSpace(res)
			time, _ := strconv.Atoi(strings.Split(res, " ")[3])
			pthread = append(pthread, strings.Split(res, " ")[2])
			tps, _ := strconv.Atoi(strings.Split(res, " ")[9])
			qps, _ := strconv.Atoi(strings.Split(res, " ")[10])
			igrs, _ := strconv.Atoi(strings.Split(res, " ")[11])
			ptps = append(ptps, tps/time)
			pqps = append(pqps, qps/time)
			pigrs = append(pigrs, igrs)

			mn["thread"] = pthread
			mn["tps"] = ptps
			mn["qps"] = pqps
			mn["igrs"] = pigrs

		}
	}
	return mn, nil
}

func SysbenchCpucut(filname string) ([][]int64, error) {

	cr, err := Readfile(filname)
	if err != nil {
		fmt.Println(err)
		return [][]int64{}, err
	}
	cc := *cr
	size := len(strings.Split(strings.TrimSpace(cc[len(cc)-1]), " "))
	t := make([][]int64, size)
	for _, v := range *cr {
		if len(v) > 0 {

			for index, value := range strings.Split(strings.TrimSpace(v), " ") {
				if index == 0 {
					single, err := strconv.Atoi(value)
					t[index] = append(t[index], int64(single))
					if err != nil {
						fmt.Println(err)
						return [][]int64{}, err
					}
				} else {
					single, err := strconv.ParseFloat(value, 32)
					t[index] = append(t[index], int64(single))
					if err != nil {
						fmt.Println(err)
						return [][]int64{}, err
					}
				}

			}
		}
	}

	return t, nil
}

func Createhmtl(filname string, htmldat string) error {
	fd, err := golog.NewFileHandler(filname, os.O_CREATE|os.O_RDWR|os.O_TRUNC)
	if err != nil {
		golog.Error("Createhtml", "Createhtml", fmt.Sprintf("%s", err), 0)
		return err
	}
	fd.Write([]byte(htmldat))
	fd.Close()
	return nil
}

func TransintTotime(unixt []int64) []string {
	var k []string
	for _, v := range unixt {
		k = append(k, time.Unix(int64(v), 0).Format("2006-01-02 03:04:05"))
	}
	return k
}

func CreateXdata(idle [][]int) string {
	return "ok"
}
*/
package main

import (
	"fmt"
	//"fmt"
	"lib/cfg"
	"regexp"
	"strings"
	//	"time"
)

var validcmd = regexp.MustCompile(`cmd?[0-9]*[0-9]$`)

func main() {
	b := [5]int{1, 2, 2, 3, 5}
	fmt.Println(b[1:])

}

func commandfromfile(c *cfg.Gbh_cfg) map[string][]string {
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
