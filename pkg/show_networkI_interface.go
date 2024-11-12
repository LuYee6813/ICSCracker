package pkg

import (
	"fmt"
	"net"
)

// ShowNetworkInterface 显示这台电脑上的所有网络接口和它们的 ID
func ShowNetworkInterface() {
	interfaces, err := net.Interfaces()
	if err != nil {
		fmt.Printf("Error fetching network interfaces: %v\n", err)
		return
	}

	// 打印表头
	fmt.Println("┌─────┬─────────────────┐")
	fmt.Printf("│ %-3s │ %-15s │\n", "ID", "Interface Name")
	fmt.Println("├─────┼─────────────────┤")

	// 打印每个网络接口的信息
	for i, iface := range interfaces {
		fmt.Printf("│ %-3d │ %-15s │\n", i, iface.Name)
	}

	// 打印表格底部
	fmt.Println("└─────┴─────────────────┘")
}
