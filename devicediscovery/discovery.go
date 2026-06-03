package devicediscovery

import (
	"fmt"

	"github.com/mytechcn/go-discovery/netdiscovery"
	"github.com/mytechcn/go-discovery/wsdiscovery"
)

func Discovery() ([]OnvifDevice, error) {
	var devList []OnvifDevice
	netInterfaces := netdiscovery.GetNetInterfaces()
	if netInterfaces == nil {
		return nil, fmt.Errorf("获取网络接口失败")
	}
	for _, netInterface := range netInterfaces {
		_devs, err := discoveryIPC(netInterface.Name)
		if err != nil {
			continue
		}
		devList = append(devList, _devs...)
	}
	return devList, nil
}

func discoveryIPC(interfaceName string) ([]OnvifDevice, error) {
	var devList []OnvifDevice
	devices, err := wsdiscovery.SendProbe(interfaceName, nil, []string{"dn:NetworkVideoTransmitter"}, map[string]string{"dn": "http://www.onvif.org/ver10/network/wsdl"})
	if err != nil {
		return nil, fmt.Errorf("搜索设备失败: %w", err)
	}
	for _, body := range devices {
		// utils.WriteFile(body)
		devs, err := ParseDeviceXML(body)
		if err != nil {
			continue
		}
		devList = append(devList, devs...)
	}
	return devList, nil
}
