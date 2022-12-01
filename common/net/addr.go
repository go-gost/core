package net

import (
	"fmt"
	"net"
)

func ParseInterfaceAddr(ifceName, network string) (ifce string, addr []net.Addr, err error) {
	if ifceName == "" {
		addr = append(addr, nil)
		return
	}

	ip := net.ParseIP(ifceName)
	if ip == nil {
		var ife *net.Interface
		ife, err = net.InterfaceByName(ifceName)
		if err != nil {
			return
		}
		var addrs []net.Addr
		addrs, err = ife.Addrs()
		if err != nil {
			return
		}
		if len(addrs) == 0 {
			err = fmt.Errorf("addr not found for interface %s", ifceName)
			return
		}
		ifce = ifceName
		for _, addr_ := range addrs {
			if ipNet, ok := addr_.(*net.IPNet); ok {
				addr = append(addr, ipToAddr(ipNet.IP, network))
			}
		}
	} else {
		ifce, err = findInterfaceByIP(ip)
		if err != nil {
			return
		}
		addr = []net.Addr{ipToAddr(ip, network)}
	}

	return
}

func ipToAddr(ip net.IP, network string) (addr net.Addr) {
	port := 0
	switch network {
	case "tcp", "tcp4", "tcp6":
		addr = &net.TCPAddr{IP: ip, Port: port}
		return
	case "udp", "udp4", "udp6":
		addr = &net.UDPAddr{IP: ip, Port: port}
		return
	default:
		addr = &net.IPAddr{IP: ip}
		return
	}
}

func findInterfaceByIP(ip net.IP) (string, error) {
	ifces, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	for _, ifce := range ifces {
		addrs, _ := ifce.Addrs()
		if len(addrs) == 0 {
			continue
		}
		for _, addr := range addrs {
			ipAddr, _ := addr.(*net.IPNet)
			if ipAddr == nil {
				continue
			}
			// logger.Default().Infof("%s-%s", ipAddr, ip)
			if ipAddr.IP.Equal(ip) {
				return ifce.Name, nil
			}
		}
	}
	return "", nil
}
