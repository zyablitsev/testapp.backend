package model

import (
	"bufio"
	"bytes"
	"net"
	"os"
	"strconv"

	"github.com/zyablitsev/testapp.backend/settings"
)

type Data struct {
	seek  int64
	users map[int]map[int]int8
}

func (m *Data) ReadLog() error {
	var (
		logFile *os.File
		ip      net.IP

		config = settings.GetInstance()

		line, lineUserID, lineIPAddress []byte
		ipBytes                         [4]byte
		userID, index                   int

		ips map[[4]byte]map[int]bool = make(map[[4]byte]map[int]bool)

		err error
	)

	if logFile, err = os.OpenFile(config.LogPath, os.O_RDONLY, 0640); err != nil {
		return err
	}
	defer logFile.Close()
	logFile.Seek(m.seek, 0)

	stat, _ := os.Stat(config.LogPath)
	m.seek = stat.Size()

	scanner := bufio.NewScanner(logFile)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		line = scanner.Bytes()
		line = line[9:]

		index = bytes.IndexByte(line, 32)

		lineUserID = line[:index]
		lineIPAddress = line[index+1:]

		ip = net.ParseIP(string(lineIPAddress)).To4()
		copy(ipBytes[:], ip)

		if userID, err = strconv.Atoi(string(lineUserID)); err != nil {
			return err
		}

		if _, ok := ips[ipBytes]; !ok {
			ips[ipBytes] = make(map[int]bool)
		}
		ips[ipBytes][userID] = true

		for x := range ips[ipBytes] {
			if x == userID {
				continue
			}

			if _, ok := m.users[x]; !ok {
				m.users[x] = make(map[int]int8)
			}
			m.users[x][userID] += 1

			if _, ok := m.users[userID]; !ok {
				m.users[userID] = make(map[int]int8)
			}
			m.users[userID][x] += 1
		}
	}

	return nil
}

func (m Data) IsDupes(userID1, userID2 int) bool {
	// Return false if at least one of users not found
	if _, ok := m.users[userID1]; !ok {
		return false
	}
	if v, ok := m.users[userID1][userID2]; !ok {
		return false
	} else {
		return v > 1
	}

	return false
}
