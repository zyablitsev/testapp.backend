package model

import (
	"bufio"
	"bytes"
	"net"
	"os"
	"sort"
	"strconv"

	"github.com/zyablitsev/testapp.backend/settings"
)

type Data struct {
	seek    int64
	Records map[int]map[[4]byte]bool
}

func (m Data) ReadLog() error {
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

		if _, ok := m.Records[userID]; !ok {
			m.Records[userID] = make(map[[4]byte]bool)
		}

		if _, ok := m.Records[userID][ipBytes]; !ok {
			m.Records[userID][ipBytes] = true
		}
	}

	return nil
}

func (m Data) IsDupes(userID1, userID2 int) bool {
	// Return false if at least one of users not found
	if _, ok := m.Records[userID1]; !ok {
		return false
	}
	if _, ok := m.Records[userID2]; !ok {
		return false
	}

	type user struct {
		id  int
		ips SortedIps
	}

	var (
		user1 = &user{id: userID1, ips: make([]*[4]byte, len(m.Records[userID1]))}
		user2 = &user{id: userID2, ips: make([]*[4]byte, len(m.Records[userID2]))}
		users = [2]*user{user1, user2}

		userWithLessIps, userWithMoreIps *user = user1, user2

		i int
	)

	for _, x := range users {
		i = 0
		for j := range m.Records[x.id] {
			p := j
			x.ips[i] = &p
			i++
		}
	}

	sort.Sort(user1.ips)
	sort.Sort(user2.ips)

	if len(user1.ips) > len(user2.ips) {
		userWithLessIps, userWithMoreIps = userWithMoreIps, userWithLessIps
	}

	if bytes.Compare(
		(*userWithLessIps.ips[0])[:], (*userWithMoreIps.ips[0])[:]) < 0 {
		idx := sort.Search(
			len(userWithLessIps.ips),
			func(i int) bool {
				return bytes.Compare(
					(*userWithLessIps.ips[i])[:],
					(*userWithMoreIps.ips[0])[:]) >= 0
			})
		userWithLessIps.ips = userWithLessIps.ips[idx:]
	}

	if len(userWithLessIps.ips) == 0 {
		return false
	}

	if bytes.Compare(
		(*userWithMoreIps.ips[len(userWithMoreIps.ips)-1])[:],
		(*userWithLessIps.ips[len(userWithLessIps.ips)-1])[:]) < 0 {
		idx := sort.Search(
			len(userWithLessIps.ips),
			func(i int) bool {
				return bytes.Compare(
					(*userWithMoreIps.ips[i])[:],
					(*userWithLessIps.ips[len(userWithLessIps.ips)-1])[:]) <= 0
			})
		userWithLessIps.ips = userWithLessIps.ips[idx:]
	}

	if len(userWithLessIps.ips) == 0 {
		return false
	}

	i = 0
	for _, x := range userWithLessIps.ips {
		p := x
		if _, ok := m.Records[userID1][*p]; ok {
			i++
		}
		if i > 1 {
			break
		}
	}

	return i > 1
}

type SortedIps []*[4]byte

func (v SortedIps) Len() int {
	return len(v)
}

func (v SortedIps) Swap(i, j int) {
	v[i], v[j] = v[j], v[i]
	return
}

func (v SortedIps) Less(i, j int) bool {
	return bytes.Compare((*v[i])[:], (*v[j])[:]) < 0
}
