package main

import (
	"bufio"
	"io"
	"log"
	"net/http"
	"os/exec"
	"sort"
	"strings"
	"sync"
	"time"
)

const (
	CsvEndPoint string = "https://www.vpngate.net/api/iphone/"
)

type Ip struct {
	ip string
	t  time.Duration
}

type FastVpn struct {
	vpnIps []Ip
	mu     sync.Mutex
	wg     sync.WaitGroup
}

func NewFastVpn() *FastVpn {
	return &FastVpn{
		vpnIps: make([]Ip, 0),
		mu:     sync.Mutex{},
		wg:     sync.WaitGroup{},
	}
}

func (f *FastVpn) SortByFaster() {
	// Sort the slice by values
	sort.Slice(f.vpnIps, func(i, j int) bool {
		return f.vpnIps[i].t < f.vpnIps[j].t
	})
}

func (f *FastVpn) Run() error {
	ic, err := f.Fetch(CsvEndPoint)
	if err != nil {
		return err
	}

	//Handle the data and get the ips list
	ips, err := f.HandleData(ic) //block until get the ips
	if err != nil {
		return err
	}

	//Ping the vpn servers using go rutine
	f.wg.Add(len(ips))
	for _, ip := range ips {
		go func(ip string, f *FastVpn, wg *sync.WaitGroup) {
			t, _ := f.pingVpnServer(ip)
			if t != 0 {
				//zero mean dead
				f.mu.Lock()
				f.vpnIps = append(f.vpnIps, Ip{
					ip: ip,
					t:  t,
				})
				f.mu.Unlock()
			}
			wg.Done()
		}(ip, f, &f.wg)
	}
	f.wg.Wait()
	return nil
}

func (f *FastVpn) Fetch(url string) (io.ReadCloser, error) {
	client := http.DefaultClient
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}

func (f *FastVpn) HandleData(d io.ReadCloser) ([]string, error) {
	ips := []string{}
	defer d.Close()
	log.Println("got the data and handling ...")
	sc := bufio.NewScanner(d)
	counter := 0
	for sc.Scan() {
		counter++
		if counter <= 2 {
			//avoid csv title and header
			continue
		}
		line := sc.Text()
		sL := strings.Split(line, ",") //Splited list
		if len(sL) > 1 {
			ips = append(ips, sL[1])
		}
	}
	log.Println("finished from handling")
	return ips, nil
}

func (f *FastVpn) pingVpnServer(ip string) (time.Duration, error) {
	log.Printf("Pinging %v ...", ip)

	startTime := time.Now()

	cmd := exec.Command("ping", "-n", "1", ip) // Adjust the arguments based on your OS, "-n" for Windows, "-c" for Unix/Linux
	output, err := cmd.CombinedOutput()

	pingTime := time.Since(startTime)

	if err != nil {
		log.Println("Error:", err)
		return 0, err
	}

	log.Printf("Ping time for %v: %v", ip, pingTime)
	log.Printf("Ping output for %v:\n%s", ip, output)

	return pingTime, nil
}

func main() {
	f := NewFastVpn()
	f.Run()
	f.SortByFaster()
	log.Println(f.vpnIps)
}
