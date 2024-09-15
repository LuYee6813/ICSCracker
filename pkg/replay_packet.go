package pkg

import (
	"errors"
	"log"
	"os"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
	"github.com/google/gopacket/pcapgo"
)

func ReplayPcap(pcapFile string, networkInterface string) error {
	handle, err := pcap.OpenLive(networkInterface, 1600, true, pcap.BlockForever)
	if err != nil {
		log.Fatalf("Error opening device: %v", err)
	}
	defer handle.Close()

	file, err := os.Open(pcapFile)
	if err != nil {
		log.Fatalf("Error opening pcap file: %v", err)
	}
	defer file.Close()

	var packetSource *gopacket.PacketSource
	pcapReader, err := pcapgo.NewNgReader(file, pcapgo.DefaultNgReaderOptions)
	if err != nil {
		file.Seek(0, os.SEEK_SET)
		pcapReaderLegacy, err := pcapgo.NewReader(file)
		if err != nil {
			return errors.New("unable to create either pcap or pcapng reader")
		}
		packetSource = gopacket.NewPacketSource(pcapReaderLegacy, pcapReaderLegacy.LinkType())
	} else {
		packetSource = gopacket.NewPacketSource(pcapReader, pcapReader.LinkType())
	}

	for packet := range packetSource.Packets() {
		// 重播流量
		err = handle.WritePacketData(packet.Data())
		if err != nil {
			log.Fatalf("Error sending packet: %v", err)
		}
		time.Sleep(time.Millisecond * 10)
	}

	return nil
}
