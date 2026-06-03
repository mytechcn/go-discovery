package netdiscovery

import "net"

type NetInterface struct {
	Name string
	MAC  string
	IP   string
}

func GetNetInterfaces() []*NetInterface {
	var netInterfaces []*NetInterface
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil
	}
	for _, iface := range ifaces {
		netInterface := &NetInterface{
			Name: iface.Name,
			MAC:  iface.HardwareAddr.String(),
			IP:   "",
		}
		addrs, _ := iface.Addrs()
		for _, addr := range addrs {
			netInterface.IP = addr.String()
			break
		}
		netInterfaces = append(netInterfaces, netInterface)
	}
	return netInterfaces
}
