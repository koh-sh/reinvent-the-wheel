//go:build ignore

// https://github.com/koh-sh/etchosts-inventory
package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/netip"
	"os"
	"regexp"
	"strings"
)

type Inventory struct {
	Meta struct {
		Hostvars struct {
		} `json:"hostvars"`
	} `json:"_meta"`
	Targets struct {
		Hosts []string `json:"hosts"`
	} `json:"targets"`
}

func checkValidIp(ipaddress string) bool {
	ip, err := netip.ParseAddr(ipaddress)
	if err != nil {
		return false
	}
	switch {
	case ip.IsInterfaceLocalMulticast():
	case ip.IsLinkLocalMulticast():
	case ip.IsLinkLocalUnicast():
	case ip.IsLoopback():
	case ip.IsMulticast():
	case ip.IsUnspecified():
	case ip.Is6() && ip.IsGlobalUnicast():
	case ip.String() == "255.255.255.255":
	case ip.IsValid():
		return true
	}
	return false
}

func main() {
	hostsfile := "/etc/hosts"

	filePtr := flag.String("f", hostsfile, "Specify different hosts [for debug purpose]")
	// Ansible Dynamic Inventory must accept --list option
	_ = flag.Bool("list", false, "Option for Ansible execution.")

	flag.Parse()

	f, err := os.Open(*filePtr)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	hosts := []string{}
	reg := " +|	+"
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		l := scanner.Text()
		splitted := regexp.MustCompile(reg).Split(l, -1)
		switch {
		case strings.HasPrefix(l, "#"):
		case len(splitted) < 2:
		case splitted[1] == "":
		case checkValidIp(splitted[0]):
			hosts = append(hosts, splitted[1])
		}
	}

	body := `{
		"_meta": {
			"hostvars": {}
		},
		"targets": {
			"hosts": {}
		}
	}`
	inventory := Inventory{}

	json.Unmarshal([]byte(body), &inventory)
	inventory.Targets.Hosts = hosts

	bytes, err := json.MarshalIndent(inventory, "", "    ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(bytes))

}
