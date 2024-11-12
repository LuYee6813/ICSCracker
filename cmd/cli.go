package cmd

import (
	"fmt"
	"net"

	"ICSCracker/pkg"

	"github.com/urfave/cli/v2"
)

// ValidateIP validates the IP address format
func ValidateIP(ip string) bool {
	return net.ParseIP(ip) != nil
}

func SetupCLI() *cli.App {
	var ip, port string
	var networkInterface string
	var pcapFile string
	var regAddr string
	var value int

	app := &cli.App{
		Name: "ICSCracker",
		// 是一個測試ICS場域的工具，可應用於所有IEC62443-2-2 ~ 4-1資安認證
		Usage: "A tool to test ICS area, applicable to all IEC62443-2-2 ~ 4-1 cybersecurity certifications",
		Flags: []cli.Flag{
			// ip
			&cli.StringFlag{
				Name:    "ip",
				Aliases: []string{"i"},
				Usage:   "Set the target IP address temporarily",
				Action: func(ctx *cli.Context, s string) error {
					if !ValidateIP(s) {
						return fmt.Errorf("invalid IP address format: %s", s)
					}
					fmt.Println("Setting target IP temporarily:", s)
					ip = s
					return nil
				},
			},
			// port
			&cli.StringFlag{
				Name:    "port",
				Aliases: []string{"p"},
				Usage:   "Set the target port temporarily",
				Action: func(ctx *cli.Context, s string) error {
					fmt.Println("Setting target port temporarily:", s)
					port = s
					return nil
				},
			},
		},
		Commands: []*cli.Command{
			// showNetworkInterface
			{
				Name:    "showNetworkInterface",
				Aliases: []string{"net"},
				Usage:   "Show all network interfaces on this computer",
				Action: func(c *cli.Context) error {
					pkg.ShowNetworkInterface()
					return nil
				},
			},
			{
				Name:    "icsProtocolFound",
				Aliases: []string{"ics"},
				Usage:   "Discover ICS protocol",
				Action: func(c *cli.Context) error {
					if ip == "" {
						return fmt.Errorf("no IP address configured. Please set the IP using --ip flag")
					}
					pkg.IcsProtocolFound(ip)
					return nil
				},
			},
			{
				Name:    "modbusRead",
				Aliases: []string{"mR"},
				Usage:   "Read modbus server register or coil value",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "regAddr",
						Usage: "Set the register or coil address you want to read",
						Action: func(ctx *cli.Context, s string) error {
							regAddr = s
							return nil
						},
					},
				},
				Subcommands: []*cli.Command{
					{
						Name:  "HR",
						Usage: "Read holding registers",
						Action: func(c *cli.Context) error {
							if ip == "" || port == "" {
								return fmt.Errorf("IP and port must be set using --ip and --port flags")
							}
							client := pkg.ConnectModbus(ip, port)
							pkg.ModbusRead(client, regAddr, "HR")
							return nil
						},
					},
					{
						Name:  "IR",
						Usage: "Read input registers",
						Action: func(c *cli.Context) error {
							if ip == "" || port == "" {
								return fmt.Errorf("IP and port must be set using --ip and --port flags")
							}
							client := pkg.ConnectModbus(ip, port)
							pkg.ModbusRead(client, regAddr, "IR")
							return nil
						},
					},
					{
						Name:  "C",
						Usage: "Read coils",
						Action: func(c *cli.Context) error {
							if ip == "" || port == "" {
								return fmt.Errorf("IP and port must be set using --ip and --port flags")
							}
							client := pkg.ConnectModbus(ip, port)
							pkg.ModbusRead(client, regAddr, "C")
							return nil
						},
					},
					{
						Name:  "IS",
						Usage: "Read input status",
						Action: func(c *cli.Context) error {
							if ip == "" || port == "" {
								return fmt.Errorf("IP and port must be set using --ip and --port flags")
							}
							client := pkg.ConnectModbus(ip, port)
							pkg.ModbusRead(client, regAddr, "IS")
							return nil
						},
					},
				},
				Action: func(cCtx *cli.Context) error {
					fmt.Println("EXAMPLES:")
					fmt.Println("   ICSCracker mR --regAddr 10 HR")
					fmt.Println("   ICSCracker mR --regAddr 50-60 C")
					return nil
				},
			},
			// modbusWrite
			{
				Name:    "modbusWrite",
				Aliases: []string{"mW"},
				Usage:   "Write modbus server register or coil value",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "regAddr",
						Usage: "Set the register or coil address you want to write",
						Action: func(ctx *cli.Context, s string) error {
							regAddr = s
							return nil
						},
						Required: true,
					},
					&cli.IntFlag{
						Name:  "value",
						Usage: "Value to be written to the register or coil",
						Action: func(ctx *cli.Context, i int) error {
							value = i
							return nil
						},
						Required: true,
					},
				},
				Subcommands: []*cli.Command{
					{
						Name:  "HR",
						Usage: "Write holding registers",
						Action: func(c *cli.Context) error {
							if ip == "" || port == "" {
								return fmt.Errorf("IP and port must be set using --ip and --port flags")
							}
							client := pkg.ConnectModbus(ip, port)
							pkg.ModbusWrite(client, regAddr, value, "HR")
							return nil
						},
					},
					{
						Name:  "C",
						Usage: "Write coils",
						Action: func(c *cli.Context) error {
							if ip == "" || port == "" {
								return fmt.Errorf("IP and port must be set using --ip and --port flags")
							}
							client := pkg.ConnectModbus(ip, port)
							pkg.ModbusWrite(client, regAddr, value, "C")
							return nil
						},
					},
				},
				Action: func(cCtx *cli.Context) error {
					fmt.Println("EXAMPLES:")
					fmt.Println("   ICSCracker mW --regAddr 10 --value 5 HR")
					fmt.Println("   ICSCracker mW --regAddr 50 --value 0xFF00 C")
					return nil
				},
			},
			// replayAttack
			{
				Name:    "replayAttack",
				Aliases: []string{"aR"},
				Usage:   "Replay attack packet",
				Flags: []cli.Flag{
					// pcap
					&cli.StringFlag{
						Name:    "pcap",
						Aliases: []string{"f"},
						Usage:   "Path to the pcap file to be replayed",
						Action: func(ctx *cli.Context, s string) error {
							pcapFile = s
							return nil
						},
						Required: true,
					},
					// networkInterface
					&cli.StringFlag{
						Name:    "networkInterface",
						Aliases: []string{"n"},
						Usage:   "Set the network interface to be used for replay attack",
						Action: func(ctx *cli.Context, s string) error {
							fmt.Println("Setting network interface:", s)
							networkInterface = s
							return nil
						},
						Required: true,
					},
				},
				Action: func(c *cli.Context) error {
					if pcapFile == "" {
						return fmt.Errorf("pcap file path must be provided using --pcap flag")
					}
					if networkInterface == "" {
						return fmt.Errorf("network interface must be set using --networkInterface flag")
					}
					fmt.Println("Replaying pcap file:", pcapFile)
					err := pkg.ReplayPcap(pcapFile, networkInterface)
					if err != nil {
						return fmt.Errorf("failed to replay pcap: %v", err)
					}
					fmt.Println("Replay attack completed successfully.")
					return nil
				},
			},
			// dosAttack
			{
				Name:    "dosAttack",
				Aliases: []string{"aD"},
				Usage:   "Denial of Service attack",
				Action: func(c *cli.Context) error {
					if ip == "" {
						return fmt.Errorf("IP address must be set using --ip flag")
					}
					pkg.DosAttack(ip)
					return nil
				},
			},
		},
	}
	return app
}

func PrintAsciiArt() {
	fmt.Println(" ___   ____  ____    ____                     _                       ")
	fmt.Println("|_ _| / ___|/ ___|  / ___| _ __   __ _   ___ | | __  ___  _ __        ")
	fmt.Println(" | | | |    \\___ \\ | |    | '__| / _` | / __|| |/ / / _ \\| '__|    ")
	fmt.Println(" | | | |___  ___) || |___ | |   | (_| || (__ |   < |  __/| |          ")
	fmt.Println("|___| \\____||____/  \\____||_|    \\__,_| \\___||_|\\_\\ \\___||_|   ")
	fmt.Println("")
}
