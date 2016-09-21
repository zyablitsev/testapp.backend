package testapp

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/zyablitsev/testapp.backend/settings"
)

var (
	logger *log.Logger
)

type dataGroup struct {
	usersIds      []int
	ips           []string
	requestsToLog int
}

type dataGroups []*dataGroup

func init() {
	rand.Seed(time.Now().Unix())
}

func random(min, max int) int {
	return rand.Intn(max-min) + min
}

func GenerateLogFile() error {
	const (
		cidrGroup1 string = "10.0.0.1/24"
		cidrGroup2 string = "10.0.1.1/24"
	)

	var (
		wg sync.WaitGroup

		logFile *os.File

		config     = settings.GetInstance()
		configPath = filepath.Dir(config.LogPath)

		users [6000]int
		ips   [508]string

		dataGroup1 = &dataGroup{
			usersIds:      users[:2000],
			ips:           ips[:254],
			requestsToLog: 30,
		}
		dataGroup2 = &dataGroup{
			usersIds:      users[2000:4000],
			ips:           ips[254:],
			requestsToLog: 30,
		}
		dataGroup3 = &dataGroup{
			usersIds:      users[4000:],
			ips:           ips[:],
			requestsToLog: 450}

		dataGroups = [3]*dataGroup{dataGroup1, dataGroup2, dataGroup3}

		err error
	)
	if _, err = os.Stat(filepath.Dir(config.LogPath)); err != nil {
		os.Mkdir(configPath, 0755)
	}

	if logFile, err = os.OpenFile(config.LogPath, os.O_WRONLY|os.O_CREATE, 0640); err != nil {
		return err
	}
	defer logFile.Close()
	logger = log.New(logFile, "", log.Ltime)

	// Populate users ids
	for i := 0; i < len(users); i++ {
		users[i] = i + 1000
	}

	// Populate dataGroup1.ips with cidrGroup1
	if err = populateIPsFromCIDR(cidrGroup1, dataGroup1.ips); err != nil {
		return err
	}

	// Populate dataGroup2.ips with cidrGroup2
	if err = populateIPsFromCIDR(cidrGroup2, dataGroup2.ips); err != nil {
		return err
	}

	// Write log file randomly
	for _, x := range dataGroups {
		wg.Add(1)
		go writeToLogFile(x, &wg)
	}

	wg.Wait()

	return nil
}

func writeToLogFile(dg *dataGroup, wg *sync.WaitGroup) {
	var (
		l = len(dg.ips)
	)
	defer wg.Done()

	for _, x := range dg.usersIds {
		for i := 0; i < dg.requestsToLog; i++ {
			logger.Println(
				strings.Join([]string{strconv.Itoa(x), dg.ips[random(0, l)]}, " "),
			)
		}
	}
}

func populateIPsFromCIDR(cidr string, ips []string) error {
	var (
		ip    net.IP
		ipnet *net.IPNet
		err   error
	)

	if ip, ipnet, err = net.ParseCIDR(cidr); err != nil {
		return err
	}

	for i := 0; i < len(ips); i++ {
		if !ipnet.Contains(ip.Mask(ipnet.Mask)) {
			err = fmt.Errorf("Can't populate IPs from CIDR: %s", cidr)
			return err
		}
		ips[i] = ip.String()
		incrementIP(ip)
	}

	return nil
}

func incrementIP(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}
