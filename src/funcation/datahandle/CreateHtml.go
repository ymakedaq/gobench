package datahandle

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"funcation/core"
	"funcation/golog"
	"html/template"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
)

type Dreawhtml struct {
	Headtitle    string //html 的title
	Interval     int
	StartTime    []int
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

func (this *Dreawhtml) Newchart(rfile_name string) {
	tmpl := template.New("")
	tmpl.Parse(tpl)
	res, err := SysbenchResCut(rfile_name + ".txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	cpures, err := SysbenchCpucut(core.CPUFILE, 2000)
	if err != nil {
		fmt.Println(err)
		return
	}
	memres, err := SysbenchCpucut(core.MEMFILE, 2000)
	if err != nil {
		fmt.Println(err)
		return
	}

	ttime := time.Unix(cpures[0][0], 0).Format("2006,01,02,03,04,05")
	stime_arry, err := ConvertTime(ttime)
	if err != nil {
		return
	}
	this.Interval = core.AtouchTime * 1000
	this.StartTime = stime_arry
	this.MemXdata = TransintTotime(memres[0])
	this.MemYdata = memres[1:]
	this.MemHeadtitle = "Memory"
	this.CpuXdata = TransintTotime(cpures[0])
	this.CpuYdata = cpures[1:]
	this.CpuHeadtitle = "CPU idle"
	this.CpuYtitle = "cpu idle"
	this.Headtitle = res["headtitle"].(string)
	this.Xtitle = "tps/s"
	this.Ytitle = "tps/s"
	this.Xdata = res["thread"].([]string)
	this.Ydata = res["tps"].([]int)

	var html bytes.Buffer
	err = tmpl.Execute(&html, this)
	//fmt.Println(html.String())
	if err != nil {
		fmt.Println(err)
	}
	err = Createhmtl(rfile_name+".html", html.String())
	if err != nil {
		fmt.Println(err)
	}
}

func ConvertTime(s string) ([]int, error) {
	var k []int
	for _, v := range strings.Split(s, ",") {
		num, err := strconv.Atoi(v)
		if err != nil {
			golog.Error("datahandle", "Convertime", fmt.Sprint(err), 0)
			return []int{}, err
		}
		k = append(k, num)
	}
	return k, nil
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

func SysbenchResCut(filename string) (map[string]interface{}, error) {
	var pthread []string
	var ptps []int
	var pqps []int
	var pigrs []int
	mn := make(map[string]interface{})

	rows, err := Readfile(filename)
	if err != nil {
		fmt.Println(err)
		return nil, nil
	}
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
			res := strings.Replace(strings.Split(v, `run"`)[1], "}", "", -1)
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

func SysbenchCpucut(filname string, cut_limit int) ([][]int64, error) {

	cr, err := Readfile(filname)
	if err != nil {
		fmt.Println(err)
		return [][]int64{}, err
	}
	cc := *cr
	size := len(strings.Split(strings.TrimSpace(cc[len(cc)-1]), " "))
	t := make([][]int64, size)
	cut_step := len(*cr) / cut_limit
	for i := 0; i < len(cc); i += cut_step {
		if v := cc[i]; len(v) > 0 {

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
