package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Ullaakut/nmap/v3"
	"github.com/goburrow/modbus"
	"github.com/urfave/cli/v2"
)

func main() {

	// TO-DO feture
	// ============ basic  setup   ===============
	// set up target ip      -i     ip address or ip range
	// set up target port    -p     port number
	// ============ attack command ===============
	// modbus read			 -mR    regAddress
	// modbus write			 -mW    regAddress value
	// replay attack         -aR    FILE(*.pcap)
	// dos attack			 -aD
	// arp spoofing attack   -aA
	// ics protocol found    --ics

	ip, err := os.ReadFile("ip-config")
	if err != nil {
		fmt.Println("must set up ip first")
	}
	port, _ := os.ReadFile("port-config")
	if err != nil {
		fmt.Println("must set up port first")
	}
	var regAddr string = ""
	var value = 0

	// 印出ascii art
	fmt.Println(" ___   ____  ____    ____                     _                       ")
	fmt.Println("|_ _| / ___|/ ___|  / ___| _ __   __ _   ___ | | __  ___  _ __        ")
	fmt.Println(" | | | |    \\___ \\ | |    | '__| / _` | / __|| |/ / / _ \\| '__|    ")
	fmt.Println(" | | | |___  ___) || |___ | |   | (_| || (__ |   < |  __/| |          ")
	fmt.Println("|___| \\____||____/  \\____||_|    \\__,_| \\___||_|\\_\\ \\___||_|   ")
	fmt.Println("")

	app := &cli.App{
		Name:  "ICSCracker",
		Usage: "is a useful ICS attack tool",
		Flags: []cli.Flag{
			// --ip -i
			&cli.StringFlag{
				Name:    "ip",
				Aliases: []string{"i"},
				Usage:   "Set the target ip",
				Action: func(ctx *cli.Context, s string) error {
					ip := s
					fmt.Println("Set the target ip:", ip)
					os.WriteFile("ip-config", []byte(ip), 0644)
					return nil
				},
			},
			// --port -p
			&cli.StringFlag{
				Name:    "port",
				Aliases: []string{"p"},
				Usage:   "Set the target port",
				Action: func(ctx *cli.Context, s string) error {
					port := s
					fmt.Println("Set the target port:", port)
					os.WriteFile("port-config", []byte(port), 0644)
					return nil
				},
			},
		},
		Commands: []*cli.Command{
			// icsProtocolFound || ics
			{
				Name:    "icsProtocolFound",
				Aliases: []string{"ics"},
				Usage:   "Discover ICS protocol",
				Action: func(c *cli.Context) error {
					icsProtocolFound(string(ip))
					return nil
				},
			},
			// modbusRead || mR
			{
				Name:    "modbusRead",
				Aliases: []string{"mR"},
				Usage:   "Read modbus server register or coli value",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "regAddr",
						Usage: "Set you want to read register or coli address",
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
							client := connectModbus(string(ip), string(port))
							modbusRead(client, regAddr, "HR")
							return nil
						},
					},
					{
						Name:  "IR",
						Usage: "Read input registers",
						Action: func(c *cli.Context) error {
							client := connectModbus(string(ip), string(port))
							modbusRead(client, regAddr, "IR")
							return nil
						},
					},
					{
						Name:  "C",
						Usage: "Read coils",
						Action: func(c *cli.Context) error {
							client := connectModbus(string(ip), string(port))
							modbusRead(client, regAddr, "C")
							return nil
						},
					},
					{
						Name:  "IS",
						Usage: "Read input status",
						Action: func(c *cli.Context) error {
							client := connectModbus(string(ip), string(port))
							modbusRead(client, regAddr, "IS")
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
			// modbusWrite || mW
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
							client := connectModbus(string(ip), string(port))
							modbusWrite(client, regAddr, value, "HR")
							return nil
						},
					},
					{
						Name:  "IR",
						Usage: "Write input registers",
						Action: func(c *cli.Context) error {
							client := connectModbus(string(ip), string(port))
							modbusWrite(client, regAddr, value, "IR")
							return nil
						},
					},
					{
						Name:  "C",
						Usage: "Write coils",
						Action: func(c *cli.Context) error {
							client := connectModbus(string(ip), string(port))
							modbusWrite(client, regAddr, value, "C")
							return nil
						},
					},
					{
						Name:  "IS",
						Usage: "Write input status",
						Action: func(c *cli.Context) error {
							client := connectModbus(string(ip), string(port))
							modbusWrite(client, regAddr, value, "IS")
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
			// replay attack || aR
			{
				Name:    "replayAttack",
				Aliases: []string{"aR"},
				Usage:   "Replay attack packet",
				Action: func(c *cli.Context) error {
					fmt.Println("Replay attack")
					return nil
				},
			},
			// dos attack || aD
			{
				Name:    "dosAttack",
				Aliases: []string{"aD"},
				Usage:   "Denial of Service attack",
				Action: func(c *cli.Context) error {
					fmt.Println("Denial of Service attack")
					return nil
				},
			},
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func icsProtocolFound(ip string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()
	// nmap -Pn -sT --scan-delay 1s --max-parallelism 1 \
	//-p \
	//80,102,443,502,530,593,789,1089-1091,1911,1962,2222,2404,4000,4840,4843,4911,9600,19999,20000,20547,34962-34964,34980,44818,46823,46824,55000-55003 \
	//<ip>
	scanner, err := nmap.NewScanner(
		ctx,
		nmap.WithSkipHostDiscovery(),      // -Pn
		nmap.WithConnectScan(),            // -sT
		nmap.WithScanDelay(1*time.Second), // --scan-delay 1s
		nmap.WithMaxParallelism(1),        // --max-parallelism 1
		nmap.WithPorts("80,102,443,502,530,593,789,1089-1091,1911,1962,2222,2404,4000,4840,4843,4911,9600,19999,20000,20547,34962-34964,34980,44818,46823,46824,55000-55003"), // -p
		nmap.WithTargets(ip), // <ip>
	)
	if err != nil {
		log.Fatalf("unable to create nmap scanner: %v", err)
	}
	fmt.Print("Scanning ICS Protocol ...")
	result, warnings, err := scanner.Run()

	if len(*warnings) > 0 {
		log.Printf("run finished with warnings: %s\n", *warnings) // Warnings are non-critical errors from nmap.
	}
	if err != nil {
		log.Fatalf("unable to run nmap scan: %v", err)
	}

	// Use the results to print an example output
	for _, host := range result.Hosts {
		if len(host.Ports) == 0 || len(host.Addresses) == 0 {
			continue
		}

		fmt.Printf("Host %q:\n", host.Addresses[0])

		for _, port := range host.Ports {
			fmt.Printf("\tPort %d/%s %s %s\n", port.ID, port.Protocol, port.State, port.Service.Name)
		}
	}

	fmt.Printf("Finished: %d hosts up scanned in %.2f seconds\n", len(result.Hosts), result.Stats.Finished.Elapsed)
}

func intToBytes(n int, count uint16) []byte {
	byteCount := int(count) * 2
	bytes := make([]byte, byteCount)
	for i := 0; i < byteCount; i += 2 {
		bytes[i] = byte(n >> 8) // 高字節
		bytes[i+1] = byte(n)    // 低字節
	}
	return bytes
}

func connectModbus(ip string, port string) modbus.Client {
	connectionString := ip + ":" + port
	log.Printf("Connecting to %s\n", connectionString)
	handler := modbus.NewTCPClientHandler(connectionString)
	handler.Timeout = 10 * time.Second
	handler.SlaveId = 1
	if err := handler.Connect(); err != nil {
		log.Fatalf("Error connecting to Modbus server: %v", err)
	}
	client := modbus.NewClient(handler)
	return client
}

func modbusRead(client modbus.Client, regAddr string, readType string) {
	addr := strings.Split(regAddr, "-")
	if len(addr) == 1 {
		addrStart, _ := strconv.Atoi(addr[0])
		switch readType {
		case "HR":
			readResult, _ := client.ReadHoldingRegisters(uint16(addrStart), 1)
			fmt.Println(readResult)
		case "IR":
			readResult, _ := client.ReadInputRegisters(uint16(addrStart), 1)
			fmt.Println(readResult)
		case "C":
			readResult, _ := client.ReadCoils(uint16(addrStart), 1)
			fmt.Println(readResult)
		case "IS":
			readResult, _ := client.ReadDiscreteInputs(uint16(addrStart), 1)
			fmt.Println(readResult)
		default:
			fmt.Println("Invalid memory type")
		}
	} else if len(addr) == 2 {
		addrStart, _ := strconv.Atoi(addr[0])
		addrEnd, _ := strconv.Atoi(addr[1])
		switch readType {
		case "HR": //HoldingRegisters
			readResult, _ := client.ReadHoldingRegisters(uint16(addrStart), uint16(addrEnd)-uint16(addrStart)+1)
			fmt.Println(readResult)
		case "IR": //InputRegisters
			readResult, _ := client.ReadInputRegisters(uint16(addrStart), uint16(addrEnd)-uint16(addrStart)+1)
			fmt.Println(readResult)
		case "C": //Coils
			readResult, _ := client.ReadCoils(uint16(addrStart), uint16(addrEnd)-uint16(addrStart)+1)
			fmt.Println(readResult)
		case "IS": //InputStatus
			readResult, _ := client.ReadDiscreteInputs(uint16(addrStart), uint16(addrEnd)-uint16(addrStart)+1)
			fmt.Println(readResult)
		default:
			fmt.Println("Invalid memory type")
		}
	} else {
		fmt.Println("`--regAddr` format is [value]-[value] or [value]")
	}
}

func modbusWrite(client modbus.Client, regAddr string, value int, writeType string) {
	addr := strings.Split(regAddr, "-")
	if len(addr) == 1 {
		addrStart, _ := strconv.Atoi(addr[0])
		switch writeType {
		case "HR":
			writeResult, _ := client.WriteSingleRegister(uint16(addrStart), uint16(value))
			fmt.Println(writeResult)
		case "IR":
			writeResult, _ := client.WriteSingleRegister(uint16(addrStart), uint16(value))
			fmt.Println(writeResult)
		case "C":
			writeResult, _ := client.WriteSingleCoil(uint16(addrStart), uint16(value))
			fmt.Println(writeResult)
		case "IS":
			writeResult, _ := client.WriteSingleRegister(uint16(addrStart), uint16(value))
			fmt.Println(writeResult)
		default:
			fmt.Println("Invalid memory type")
		}
	} else if len(addr) == 2 {
		addrStart, _ := strconv.Atoi(addr[0])
		addrEnd, _ := strconv.Atoi(addr[1])
		switch writeType {
		case "HR": //HoldingRegisters
			fmt.Println("Write multiple registers")
			count := uint16(addrEnd) - uint16(addrStart) + 1
			client.WriteMultipleRegisters(uint16(addrStart), uint16(addrEnd)-uint16(addrStart)+1, intToBytes(value, count))

		case "IR": //InputRegisters
			count := uint16(addrEnd) - uint16(addrStart) + 1
			client.WriteMultipleRegisters(uint16(addrStart), uint16(addrEnd)-uint16(addrStart)+1, intToBytes(value, count))

		case "C": //Coils
			fmt.Println("Coils not support write multiple registers")

		case "IS": //InputStatus
			count := uint16(addrEnd) - uint16(addrStart) + 1
			client.WriteMultipleRegisters(uint16(addrStart), uint16(addrEnd)-uint16(addrStart)+1, intToBytes(value, count))

		default:
			fmt.Println("Invalid memory type")
		}
	} else {
		fmt.Println("`--regAddr` format is [value]-[value] or [value]")
	}
}

func dosAttack(ip string) {

}
