package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("%s needs host name as argument\n", os.Args[0])
		os.Exit(1)
	}
	host := os.Args[1]
	addrs, err := net.LookupHost(host)
	if err != nil {
		fmt.Printf("failed to lookup %s: %v\n", host, err)
		os.Exit(1)
	}
	var errs []error
	for _, addr := range addrs {
		hostIP := net.ParseIP(addr)
		myIP, err := ipWithinSameSubnetAs(hostIP)
		if err != nil {
			errs = append(errs, fmt.Errorf("addr (%s): %v", addr, err))
			continue
		}
		fmt.Print(myIP.String())
		os.Exit(0)
	}
	for _, err := range errs {
		fmt.Println(err)
	}
	os.Exit(1)
}

func ipWithinSameSubnetAs(hostIP net.IP) (net.IP, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return hostIP, fmt.Errorf("failed to get interfaces: %v", err)
	}
	for _, iface := range ifaces {
		addrs, err := iface.Addrs()
		if err != nil {
			return hostIP, fmt.Errorf("failed to get addrs from interface (%s): %v", iface.Name, err)
		}
		for _, addr := range addrs {
			addrNet, ok := addr.(*net.IPNet)
			if !ok {
				continue
			}
			if addrNet.Contains(hostIP) {
				return addrNet.IP, nil
			}
		}
	}
	return hostIP, fmt.Errorf("couldn't find IP within same subnet as %s", hostIP.String())
}
