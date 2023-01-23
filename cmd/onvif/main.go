package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"path"
	"regexp"
	"strings"

	"github.com/beevik/etree"
	discover "github.com/use-go/onvif/ws-discovery"
)

type Host struct {
	URL  string `json:"url"`
	Name string `json:"name"`
}

func main() {
	fmt.Println("ONVIF TEST START")
	// ctx := context.Background()

	// dev, err := goonvif.NewDevice(goonvif.DeviceParams{
	// 	Xaddr:      "192.168.0.101:80",
	// 	Username:   "",
	// 	Password:   "",
	// 	HttpClient: new(http.Client),
	// })
	// if err != nil {
	// 	panic(err)
	// }

	// systemDateAndTyme := device.GetSystemDateAndTime{}
	// getCapabilities := device.GetCapabilities{Category: "All"}

	// systemDateAndTymeResponse, err := sdk.Call_GetSystemDateAndTime(ctx, dev, systemDateAndTyme)
	// if err != nil {
	// 	log.Println(err)
	// } else {
	// 	fmt.Println(systemDateAndTymeResponse)
	// }
	// getCapabilitiesResponse, err := sdk.Call_GetCapabilities(ctx, dev, getCapabilities)
	// if err != nil {
	// 	log.Println(err)
	// } else {
	// 	fmt.Println(getCapabilitiesResponse)
	// }

	// 	On Unix-like systems (Linux, macOS):
	// "eth0"
	// "en0"
	// "wlan0"
	// On Windows:
	// "Ethernet"
	// "Wi-Fi"
	itfs, _ := net.Interfaces()
	fmt.Printf("%+v\n", itfs)

	interfaceName := "lo"

	var hosts []*Host
	devices, err := discover.SendProbe(interfaceName, nil, []string{"dn:NetworkVideoTransmitter"}, map[string]string{"dn": "http://www.onvif.org/ver10/network/wsdl"})
	if err != nil {
		log.Printf("error %s", err)
		return
	}
	for _, j := range devices {
		doc := etree.NewDocument()
		if err := doc.ReadFromString(j); err != nil {
			log.Printf("error %s", err)
		} else {

			endpoints := doc.Root().FindElements("./Body/ProbeMatches/ProbeMatch/XAddrs")
			scopes := doc.Root().FindElements("./Body/ProbeMatches/ProbeMatch/Scopes")

			flag := false

			host := &Host{}

			for _, xaddr := range endpoints {
				xaddr := strings.Split(strings.Split(xaddr.Text(), " ")[0], "/")[2]
				host.URL = xaddr
			}
			if flag {
				break
			}
			for _, scope := range scopes {
				re := regexp.MustCompile(`onvif:\/\/www\.onvif\.org\/name\/[A-Za-z0-9-]+`)
				match := re.FindStringSubmatch(scope.Text())
				host.Name = path.Base(match[0])
			}

			hosts = append(hosts, host)

		}

	}

	bys, _ := json.Marshal(hosts)
	log.Printf("done %s", bys)
}
