package main

import (
	"flag"
	"fmt"
	"golang.org/x/sys/windows/registry"
	"log"
	"runtime"
	"strings"
)

type iList map[string]string

func getOSInfo() (nam, ver string) {
	ck, err := registry.OpenKey(registry.LOCAL_MACHINE, "SOFTWARE\\Microsoft\\Windows NT\\CurrentVersion", registry.READ)
	if err != nil {
		fmt.Println(err.Error())
		return "", ""
	}
	vn, err := ck.ReadValueNames(0)
	if err != nil {
		fmt.Println(err.Error())
		return "", ""
	}
	nam, _, err = ck.GetStringValue("ProductName")
	if err != nil {
		fmt.Println(err.Error())
		return "", ""
	}
	sb := strings.Builder{}
	if checkString("CSDVersion", vn) {
		cv, _, err := ck.GetStringValue("CSDVersion")
		if err != nil {
			fmt.Println(err.Error())
			return "", ""
		}
		sb.WriteString(cv)
		sb.WriteString(" ")
	}
	if checkString("DisplayVersion", vn) {
		dv, _, err := ck.GetStringValue("DisplayVersion")
		if err != nil {
			fmt.Println(err.Error())
			return "", ""
		}
		sb.WriteString(dv)
		sb.WriteString(" ")
	}
	cb, _, err := ck.GetStringValue("CurrentBuild")
	if err != nil {
		fmt.Println(err.Error())
		return "", ""
	}
	sb.WriteString(cb)
	ver = sb.String()
	return
}
func checkString(a string, b []string) bool {
	for _, x := range b {
		if a == x {
			return true
		}
	}
	return false
}
func getSoftwareInfo(key registry.Key, path string) (nam, ver string) {
	ck, err := registry.OpenKey(key, path, registry.READ)
	if err != nil {
		log.Fatal(err.Error())
	}
	vn, err := ck.ReadValueNames(0)
	if err != nil {
		log.Fatal(err.Error())
	}
	if checkString("DisplayName", vn) && checkString("DisplayVersion", vn) {
		nam, _, err = ck.GetStringValue("DisplayName")
		if err != nil {
			log.Fatal(err.Error())
		}
		ver, _, err = ck.GetStringValue("DisplayVersion")
		if err != nil {
			log.Fatal(err.Error())
		}
	}
	return
}
func getInstallSoftware(key registry.Key, path string) (instSoft iList, err error) {
	instSoft = iList{}
	k, err := registry.OpenKey(key, path, registry.READ)
	if err != nil {
		return nil, err
	}
	kns, _ := k.ReadSubKeyNames(0)
	for _, kn := range kns {
		ss := path + "\\" + kn
		nam, ver := getSoftwareInfo(key, ss)
		nam = strings.Split(nam, " - ")[0]
		instSoft[nam] = ver
	}
	return instSoft, nil
}
func main() {
	host := flag.String("h", "", "Server host ip address")
	flag.Parse()
	var SoftList = iList{}
	on, ov := getOSInfo()
	SoftList[on] = ov
	l1, err := getInstallSoftware(registry.LOCAL_MACHINE, "SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Uninstall")
	if err != nil {
		log.Fatal(err.Error())
		return
	}
	l2, err := getInstallSoftware(registry.CURRENT_USER, "SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Uninstall")
	if err != nil {
		log.Fatal(err.Error())
		return
	}
	for k, v := range l1 {
		if k != "" {
			SoftList[k] = v
		}
	}
	for k, v := range l2 {
		if k != "" {
			SoftList[k] = v
		}
	}
	if runtime.GOARCH == "amd64" {
		l3, err := getInstallSoftware(registry.LOCAL_MACHINE, "SOFTWARE\\Wow6432Node\\Microsoft\\Windows\\CurrentVersion\\Uninstall")
		if err != nil {
			log.Fatal(err.Error())
			return
		}
		for k, v := range l3 {
			if k != "" {
				SoftList[k] = v
			}
		}
	}
	err = sendData(*host, SoftList)
	if err != nil {
		fmt.Println(err.Error())
	}
}
