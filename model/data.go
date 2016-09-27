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
	users map[int]map[int]int
	ips   map[[4]byte]map[int]int
}

func (m *Data) ReadLog() error {
	var (
		logFile *os.File
		ip      net.IP

		config = settings.GetInstance()

		line, lineUserID, lineIPAddress []byte
		ipBytes                         [4]byte
		userID, index                   int

		ips map[[4]byte]map[int]int = make(map[[4]byte]map[int]int)

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

		if _, ok := m.ips[ipBytes]; !ok {
			m.ips[ipBytes] = make(map[int]int)
			ips[ipBytes] = m.ips[ipBytes]
		}

		if _, ok := ips[ipBytes][userID]; !ok {
			ips[ipBytes][userID] = 1
			m.ips[ipBytes][userID] = 1
		} else {
			ips[ipBytes][userID] += 1
			m.ips[ipBytes][userID] += 1
		}
	}

	for _, v := range ips {
		for uid := range v {
			for uid2, count := range v {
				if uid2 == uid {
					continue
				}

				if _, ok := m.users[uid]; !ok {
					m.users[uid] = make(map[int]int)
				}

				if _, ok := m.users[uid][uid2]; !ok {
					m.users[uid][uid2] = count
				} else {
					m.users[uid][uid2] += count
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
