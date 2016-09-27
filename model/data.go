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
	ips   map[[4]byte]map[int]bool
}

func (m *Data) ReadLog() error {
	var (
		logFile *os.File
		ip      net.IP

		config = settings.GetInstance()

		line, lineUserID, lineIPAddress []byte
		ipBytes                         [4]byte
		userID, index                   int

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

		if _, ok := m.users[userID]; !ok {
			m.users[userID] = make(map[int]int8)
		}

		if _, ok := m.ips[ipBytes]; !ok {
			m.ips[ipBytes] = make(map[int]bool)
		}

		if _, ok := m.ips[ipBytes][userID]; !ok {
			m.ips[ipBytes][userID] = true
		}

		for uid := range m.ips[ipBytes] {
			if uid == userID {
				continue
			}

			if v, ok := m.users[uid]; ok {
				if counter, ok := v[userID]; !ok {
					v[userID] = 1
				} else {
					if counter < 2 {
						v[userID] += 1
					}
				}
			}
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
