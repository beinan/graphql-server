package fastid

import (
	"errors"
	"net"
	"os"
	"strconv"
	"time"
)

const (
	StartTimeEnvName           = "FASTID_START_TIME"
	MachineIDEnvName           = "FASTID_MACHINE_ID"
	defaultStartTimeStr        = "2018-06-01T00:00:00.000Z"
	defaultStartTimeNano int64 = 1527811200000000000
	timeBits                   = 40
	seqBits                    = 7
	machineBits                = 16
)

var seq int64

var machineIDMask = ^(int64(-1) << machineBits)
var machineID = getMachineId() & machineIDMask

var timeMask = ^(int64(-1) << timeBits)
var startEpochNano = getStartEpochFromEnv()

var seqMask = ^(int64(-1) << seqBits)

var seqChan = make(chan int64, 32)

func init() {
	go func() {
		var lastGenTime int64 = 0
		var seq int64 = 0
		for {
			now := getTimestamp()
			if now > lastGenTime {
				lastGenTime = now
				seq = 0
			}
			if seq < seqMask {
				seq += 1
				seqChan <- seq
			} else {
				time.Sleep(5 * time.Microsecond)
			}
		}
	}()
}

func getTimestamp() int64 {
	//devided by 2^20 (10^6, nano to milliseconds)
	return (time.Now().UnixNano() - startEpochNano) >> 20 & timeMask
}

func GenInt64ID() int64 {
	timestamp := getTimestamp()
	seq := <-seqChan
	return timestamp<<(seqBits+machineBits) + seq<<machineBits + machineID
}

func getMachineId() int64 {
	//getting machine from env
	if machineIDStr, ok := os.LookupEnv(MachineIDEnvName); ok {
		if machineID, err := strconv.ParseInt(machineIDStr, 10, 64); err == nil {
			return machineID
		}
	}
	//take the lower 16bits of IP address as Machine ID
	if ip, err := getIP(); err == nil {
		return (int64(ip[2]) << 8) + int64(ip[3])
	}
	return 0
}

func getStartEpochFromEnv() int64 {
	startTimeStr := getEnv(StartTimeEnvName, defaultStartTimeStr)
	var startEpochTime, err = time.Parse(time.RFC3339, startTimeStr)

	if err == nil {
		return defaultStartTimeNano
	}

	return startEpochTime.UnixNano()
}

func getIP() (net.IP, error) {
	if addrs, err := net.InterfaceAddrs(); err == nil {
		for _, addr := range addrs {
			if ipNet, ok := addr.(*net.IPNet); ok {
				if !ipNet.IP.IsLoopback() && ipNet.IP.To4() != nil {
					ip := ipNet.IP.To4()

					if ip[0] == 10 || ip[0] == 172 && (ip[1] >= 16 && ip[1] < 32) || ip[0] == 192 && ip[1] == 168 {
						return ip, nil
					}
				}
			}
		}
	}
	return nil, errors.New("Failed to get ip address")
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
