package devicediscovery

import (
	"net/url"
	"regexp"
	"strings"

	"github.com/beevik/etree"
)

// OnvifDevice 统一清洗后的摄像机设备资产结构
type OnvifDevice struct {
	IP           string `json:"ip"`            // 从 XAddrs 中提取出的干净 IPv4 地址
	XAddr        string `json:"xaddr"`         // 完整的 ONVIF 服务接入点 URL
	Brand        string `json:"brand"`         // 品牌/厂商 (如 HIKVISION, LTS)
	Model        string `json:"model"`         // 硬件型号 (如 DS-2CD1123G2-LIU)
	MAC          string `json:"mac"`           // 物理 MAC 地址
	Type         string `json:"type"`          // 设备类型 (IPC/NVR/DECODER/PC/UNKNOWN)
	RawType      string `json:"raw_type"`      // 原始的设备类型响应字符串
	LocationCity string `json:"location_city"` // 城市信息
}

var (
	// 在包级别预编译正则表达式，极大提升高并发/多文件批量解析时的执行效率
	//匹配一个或多个『既不是斜杠 /』，『也不是任意空白符 \s』的连续字符。
	//%s是制表符\t (Tab)、换行符\n (LF)、回车符\r (CR)等的统称，\s表示任意空白字符。
	reBrand = regexp.MustCompile(`onvif://www\.onvif\.org/name/([^/\s]+)`)
	reModel = regexp.MustCompile(`onvif://www\.onvif\.org/hardware/([^/\s]+)`)
	reMAC   = regexp.MustCompile(`onvif://www\.onvif\.org/MAC/([^/\s]+)`)
	reCity  = regexp.MustCompile(`onvif://www.onvif.org/location/city/([^/\s]+)`)
)

// ParseDeviceXML 保持独立性的核心纯净函数（解耦文件、网络与日志）
func ParseDeviceXML(xmlContent string) ([]OnvifDevice, error) {
	doc := etree.NewDocument()
	if err := doc.ReadFromString(xmlContent); err != nil {
		return nil, err
	}

	// 使用 XPath 的 local-name 模糊匹配通杀所有命名空间前缀(soap/s/env/wsd/d)
	matchNodes := doc.FindElements("//*[local-name()='ProbeMatch']")
	var devices []OnvifDevice

	for _, node := range matchNodes {
		var dev OnvifDevice

		// A. 提取原始设备类型并判断是否为 IPC
		if typesElem := node.FindElement(".//*[local-name()='Types']"); typesElem != nil {
			dev.RawType = typesElem.Text()
			dev.Type = extractType(dev.RawType)
		}

		// B. 提取 XAddrs 接入点，斩断 IPv6 混合地址，仅保留第一条 IPv4
		if xaddrsElem := node.FindElement(".//*[local-name()='XAddrs']"); xaddrsElem != nil {
			text := strings.TrimSpace(xaddrsElem.Text())
			if text != "" {
				dev.XAddr = strings.Split(text, " ")[0]
				//移除换行符和制表符等不可见字符，确保 URL 解析的干净输入
				dev.XAddr = strings.ReplaceAll(dev.XAddr, "\n", "")
				dev.XAddr = strings.ReplaceAll(dev.XAddr, "\t", "")
				if u, err := url.Parse(dev.XAddr); err == nil {
					dev.IP = u.Hostname()
				}
			}
		}

		// C. 深度清洗 Scopes 元数据字段
		if scopesElem := node.FindElement(".//*[local-name()='Scopes']"); scopesElem != nil {
			scopes := strings.Fields(scopesElem.Text())

			for _, s := range scopes {
				// 提取品牌/厂商名称
				if m := reBrand.FindStringSubmatch(s); len(m) > 1 {
					dev.Brand = extractBrand(m[1])
				} else if dev.Brand == "" && strings.Contains(s, "/name/") {
					dev.Brand = strings.Split(s, "/name/")[1]
					dev.Brand = extractBrand(dev.Brand)
				} else if dev.Brand == "" && strings.Contains(s, "LTS") {
					dev.Brand = "LTS" // 针对 LTS 某些缺省包的后备兜底
				}

				// 提取硬件型号
				if m := reModel.FindStringSubmatch(s); len(m) > 1 {
					dev.Model = m[1]
				} else if dev.Model == "" && strings.Contains(s, "hardware/") {
					dev.Model = strings.Split(s, "hardware/")[1]
				}

				// 提取 MAC 地址
				if m := reMAC.FindStringSubmatch(s); len(m) > 1 {
					dev.MAC = m[1]
				}

				// 提取城市信息
				if m := reCity.FindStringSubmatch(s); len(m) > 1 {
					dev.LocationCity = m[1]
				}
			}
		}
		devices = append(devices, dev)
	}

	return devices, nil
}

func extractBrand(rawBrand string) string {
	sp := "%20"
	if strings.Contains(rawBrand, sp) {
		return strings.Split(rawBrand, sp)[0]
	}
	return rawBrand
}

func extractType(rawTypes string) string {
	if strings.Contains(rawTypes, "NetworkVideoTransmitter") {
		return "IPC" // 确定是摄像机
	} else if strings.Contains(rawTypes, "NetworkVideoStorage") {
		return "NVR" // 确定是硬盘录像机
	} else if strings.Contains(rawTypes, "NetworkVideoDisplay") {
		return "DECODER" // 确定是解码器
	} else if strings.Contains(rawTypes, "pub:Computer") {
		return "PC" // 标记为 IT 干扰项，生产环境通常直接丢弃（不予入库）
	} else {
		return "UNKNOWN" // 其他未知节点
	}
}
