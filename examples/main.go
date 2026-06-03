package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/mytechcn/go-discovery/devicediscovery"
	"github.com/mytechcn/go-discovery/utils"
)

func main() {
	dev()
	// test()
}

func test() {
	hosts, err := devicediscovery.Discovery()
	if err != nil {
		log.Fatal(err)
	}
	for _, host := range hosts {
		log.Printf("IP: %s, Brand: %s\n", host.IP, host.Brand)
	}
	log.Println(hosts)
}

func dev() {
	tmpDir := "tmp"
	fmt.Printf("【主系统启动】正在读取本地临时缓存目录: %s", tmpDir)

	// 1. 获取目录下的文件
	files, err := os.ReadDir(tmpDir)
	if err != nil {
		fmt.Printf("打开 tmp 文件夹失败: %v\n", err)
		return
	}

	var totalAsset []devicediscovery.OnvifDevice

	// 2. 业务外围循环控制
	for _, file := range files {
		filename := utils.GetFileNameNoExt(file.Name())
		// 精准过滤非 XML 后缀文件
		if file.IsDir() || strings.ToLower(filepath.Ext(file.Name())) != ".xml" {
			continue
		}

		filePath := filepath.Join(tmpDir, file.Name())

		// 读取文件为字节流
		xmlBytes, err := os.ReadFile(filePath)
		if err != nil {
			fmt.Printf("跳过受损文件: %s, 错误: %v\n", file.Name(), err)
			continue
		}

		// 3. 【核心调用】：完全隔离独立的解析逻辑
		devices, err := devicediscovery.ParseDeviceXML(string(xmlBytes))
		if err != nil {
			fmt.Printf("文件 [%s] 格式不受支持，无法通过 ONVIF 校验: %v\n", file.Name(), err)
			continue
		}
		// 4. 写入 JSON 文件（可选步骤，展示如何将解析结果持久化）
		// 序列化
		buf, err := json.MarshalIndent(devices, "", "  ")
		if err != nil {
			fmt.Printf("序列化设备数据失败: %v\n", err)
			continue
		}
		fmt.Println(string(buf))
		utils.WriteJSONFile(filename+".json", string(buf))

		fmt.Printf("✅ 文件 [%s] 解析成功，发现设备%v\n", file.Name(), devices)

		// 丰富外围信息并汇总
		for _, d := range devices {
			totalAsset = append(totalAsset, d)

			// 打印当前文件的业务洗出日志
			fmt.Printf("➡️ [文件溯源: %s] 洗出 IP -> %s (%s - %s)\n", file.Name(), d.IP, d.Brand, d.Model)
		}
	}

	// 4. 汇总展示最终数组资产
	fmt.Println(strings.Repeat("=", 70))
	fmt.Printf("【处理完毕】tmp 目录下最终转化的结构体数组长度为: %d\n", len(totalAsset))
	fmt.Println(strings.Repeat("=", 70))
}
