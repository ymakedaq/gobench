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

package main

import (
	//	"github.com/ivpusic/grpool"
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"funcation/golog"
	"html/template"
	"io"
	"os"
	"strconv"
	"strings"
)

const (
	CPUFILE    = `./cpufile.txt`
	MEMFILE    = `./memfile.txt`
	resultfile = `./res.txt`
)

type Dreawhtml struct {
	Headtitle string //html 的title
	Ytitle    string
	Xtitle    string
	Xdata     []string //x轴的数据值
	Ydata     []int    //y轴的数据值
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

	SysbenchCpucut()
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

func SysbenchCpucut() {
	//	var t [32][]float64
	cr, err := Readfile(CPUFILE)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(len(*cr))
	for _, v := range *cr {
		if len(v) > 0 {

			for index, value := range strings.Split(strings.TrimSpace(v), " ") {
				fmt.Println(index, value)
				single, err := strconv.ParseFloat(value, 32)
				if err != nil {
					fmt.Println(err)
					return
				}
				t[index] = append(t[index], single)
			}
		}
	}
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
