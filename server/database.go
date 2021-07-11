package main

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io/ioutil"
	"os"
	"strconv"
	"time"
)

func utf2gbk(us []byte) (gs []byte, err error) {
	gs, err = ioutil.ReadAll(transform.NewReader(bytes.NewReader(us), simplifiedchinese.GBK.NewEncoder()))
	if err != nil {
		return nil, err
	}
	return
}
func countWrite(res <-chan softCount) {
	now := time.Now()
	filename := fmt.Sprintf("./softwareCount-%s.csv", now.Format("20060102150405"))
	for {
		select {
		case sc := <-res:
			countFile, err := os.Create(filename)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			cWriter := csv.NewWriter(countFile)
			for k, v := range sc {
				s, e := utf2gbk([]byte(k))
				if e != nil {
					fmt.Println(err.Error())
					continue
				}
				cWriter.Write([]string{string(s), strconv.Itoa(v)})
			}
			cWriter.Flush()
			countFile.Close()
		}
	}
}
func hostWrite(res <-chan iList) {
	now := time.Now()
	hostFile, err := os.Create(fmt.Sprintf("./hostSoftwareList-%s.csv", now.Format("20060102150405")))
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer hostFile.Close()
	cWriter := csv.NewWriter(hostFile)
	tick := time.Tick(time.Minute)
	for {
		select {
		case rs := <-res:
			host := rs["remote"]
			delete(rs, "remote")
			cWriter.Write([]string{host, "", ""})
			for k, v := range rs {
				s, e := utf2gbk([]byte(k))
				if e != nil {
					fmt.Println(k, " : ", err.Error())
					continue
				}
				cWriter.Write([]string{"", string(s), v})
			}
		case <-tick:
			cWriter.Flush()
		}
	}
}
