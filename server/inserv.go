package main

import (
	"fmt"
	"strings"
	"time"
)

type softCount map[string]int
type iList map[string]string

func centerString(s, ins string, width int) string {
	return fmt.Sprintf("%*s  %s  %*s", (width-len(s)-2)/2, ins, s, (width-len(s)-2)/2, ins)
}
func (il *iList) keys() []string {
	s := make([]string, len(*il)/2)
	for k := range *il {
		s = append(s, k)
	}
	return s
}
func (il *iList) values() []string {
	vl := make([]string, len(*il)/2)
	for _, v := range *il {
		vl = append(vl, v)
	}
	return vl
}
func (il *iList) String() string {
	i := *il
	sl := make([]string, len(i))
	s := centerString(i["remote"], "-", 400) + "\n"
	delete(i, "remote")
	for k, v := range i {
		sl = append(sl, strings.Join([]string{k, v}, " : "))
	}
	s += strings.Join(sl, "\n")
	return s
}
func checkString(a string, b []string) bool {
	for _, x := range b {
		if a == x {
			return true
		}
	}
	return false
}
func main() {
	c := softCount{}
	var softwareList []string
	var hostList []string
	res := make(chan iList)
	wres := make(chan iList)
	wcou := make(chan softCount)
	tick := time.Tick(time.Minute)
	go recv(res)
	go countWrite(wcou)
	go hostWrite(wres)
	for {
		select {
		case rs := <-res:
			host := rs["remote"]
			if checkString(host, hostList) {
				continue
			} else {
				fmt.Printf("Recive %s SoftWare Data!\n", host)
				wres <- rs
				hostList = append(hostList, host)
				for _, k := range rs.keys() {
					if k == "remote" || k == "" {
						continue
					}
					if checkString(k, softwareList) {
						c[k] += 1
					} else {
						c[k] = 1
						softwareList = append(softwareList, k)
					}
				}
			}
		case <-tick:
			wcou <- c
		}
	}
}
