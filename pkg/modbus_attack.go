package pkg

import (
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/goburrow/modbus"
)

func ConnectModbus(ip string, port string) modbus.Client {
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

func ModbusRead(client modbus.Client, regAddr string, readType string) {
	addr := strings.Split(regAddr, "-")
	if len(addr) == 1 {
		addrStart, _ := strconv.Atoi(addr[0])
		switch readType {
		case "HR":
			readResult, _ := client.ReadHoldingRegisters(uint16(addrStart), 1)
			log.Println(readResult)
		case "IR":
			readResult, _ := client.ReadInputRegisters(uint16(addrStart), 1)
			log.Println(readResult)
		case "C":
			readResult, _ := client.ReadCoils(uint16(addrStart), 1)
			log.Println(readResult)
		case "IS":
			readResult, _ := client.ReadDiscreteInputs(uint16(addrStart), 1)
			log.Println(readResult)
		default:
			log.Println("Invalid memory type")
		}
	} else if len(addr) == 2 {
		addrStart, _ := strconv.Atoi(addr[0])
		addrEnd, _ := strconv.Atoi(addr[1])
		switch readType {
		case "HR":
			readResult, _ := client.ReadHoldingRegisters(uint16(addrStart), uint16(addrEnd)-uint16(addrStart)+1)
			log.Println(readResult)
		case "IR":
			readResult, _ := client.ReadInputRegisters(uint16(addrStart), uint16(addrEnd)-uint16(addrStart)+1)
			log.Println(readResult)
		case "C":
			readResult, _ := client.ReadCoils(uint16(addrStart), uint16(addrEnd)-uint16(addrStart)+1)
			log.Println(readResult)
		case "IS":
			readResult, _ := client.ReadDiscreteInputs(uint16(addrStart), uint16(addrEnd)-uint16(addrStart)+1)
			log.Println(readResult)
		default:
			log.Println("Invalid memory type")
		}
	} else {
		log.Println("`--regAddr` format is [value]-[value] or [value]")
	}
}

func ModbusWrite(client modbus.Client, regAddr string, value int, writeType string) {
	addr := strings.Split(regAddr, "-")
	if len(addr) == 1 {
		addrStart, _ := strconv.Atoi(addr[0])
		switch writeType {
		case "HR":
			writeResult, _ := client.WriteSingleRegister(uint16(addrStart), uint16(value))
			log.Println(writeResult)
		case "IR":
			writeResult, _ := client.WriteSingleRegister(uint16(addrStart), uint16(value))
			log.Println(writeResult)
		case "C":
			writeResult, _ := client.WriteSingleCoil(uint16(addrStart), uint16(value))
			log.Println(writeResult)
		case "IS":
			writeResult, _ := client.WriteSingleRegister(uint16(addrStart), uint16(value))
			log.Println(writeResult)
		default:
			log.Println("Invalid memory type")
		}
	} else if len(addr) == 2 {
		addrStart, _ := strconv.Atoi(addr[0])
		addrEnd, _ := strconv.Atoi(addr[1])
		switch writeType {
		case "HR":
			count := uint16(addrEnd) - uint16(addrStart) + 1
			client.WriteMultipleRegisters(uint16(addrStart), uint16(addrEnd)-uint16(addrStart)+1, intToBytes(value, count))

		case "IR":
			count := uint16(addrEnd) - uint16(addrStart) + 1
			client.WriteMultipleRegisters(uint16(addrStart), uint16(addrEnd)-uint16(addrStart)+1, intToBytes(value, count))

		case "C":
			log.Println("Coils not support write multiple registers")

		case "IS":
			count := uint16(addrEnd) - uint16(addrStart) + 1
			client.WriteMultipleRegisters(uint16(addrStart), uint16(addrEnd)-uint16(addrStart)+1, intToBytes(value, count))

		default:
			log.Println("Invalid memory type")
		}
	} else {
		log.Println("`--regAddr` format is [value]-[value] or [value]")
	}
}

func intToBytes(n int, count uint16) []byte {
	byteCount := int(count) * 2
	bytes := make([]byte, byteCount)
	for i := 0; i < byteCount; i += 2 {
		bytes[i] = byte(n >> 8)
		bytes[i+1] = byte(n)
	}
	return bytes
}
