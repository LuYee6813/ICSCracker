package pkg

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Ullaakut/nmap/v3"
)

func IcsProtocolFound(ip string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	scanner, err := nmap.NewScanner(
		ctx,
		nmap.WithSkipHostDiscovery(),
		nmap.WithConnectScan(),
		nmap.WithScanDelay(1*time.Second),
		nmap.WithMaxParallelism(1),
		nmap.WithPorts("80,102,443,502,530,593,789,1089-1091,1911,1962,2222,2404,4000,4840,4843,4911,9600,19999,20000,20547,34962-34964,34980,44818,46823,46824,55000-55003"),
		nmap.WithTargets(ip),
	)
	if err != nil {
		log.Fatalf("Unable to create nmap scanner: %v", err)
	}

	fmt.Print("Scanning ICS Protocol ...")
	result, warnings, err := scanner.Run()

	if len(*warnings) > 0 {
		log.Printf("Run finished with warnings: %s\n", *warnings)
	}
	if err != nil {
		log.Fatalf("Unable to run nmap scan: %v", err)
	}

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
