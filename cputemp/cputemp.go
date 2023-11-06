package cputemp

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/shirou/gopsutil/host"
)

// HighTemp is the highest temperature for the WaitToCool method to detect that the cpu is too hot
//
// note: if set to strict mode, LowTemp will take this variabes place
//
// default: 64
var HighTemp uint = 64

// LowTemp is the lowest temperature for the WaitToCool method to wait for, before deciding that the cpu has colled down,
// if it was previously throttled by HighTemp (except in strict mode, where LowTemp takes the place of HighTemp)
//
// default: 56
var LowTemp uint = 56

// when Logging is enabled, the WaitToCool method will log info to the console to report when it is waiting for the cpu to cool down
//
// you can set this var to false to disable this feature
//
// default: true
var Logging bool = true

// GetTemp returns the average cpu temperature in celsius
func GetTemp() uint {
	temps, err := host.SensorsTemperatures()
	if err != nil {
		return 0
	}

	var i float64
	var temp float64
	for _, t := range temps {
		if strings.HasSuffix(t.SensorKey, "_input") {
			i++
			temp += t.Temperature
		}
	}

	temp = math.Round(temp / i)
	if temp < 0 || temp > 1000 {
		return 0
	}
	return uint(temp)
}

// WaitToCool makes your function wait for the cpu to cool down
//
// by default, if the temperature > HighTemp, it will wait until the temperature <= LowTemp
//
// in strict mode, this will run if temperature > LowTemp
//
// HighTemp = 64
// LowTemp = 56
func WaitToCool(strict bool){
	flagTemp := HighTemp
	if strict {
		flagTemp = LowTemp
	}

	if temp := GetTemp(); temp > flagTemp {
		if Logging {
			fmt.Println("CPU Too Hot!")
			fmt.Println("Waiting for it to cool down...")
			fmt.Print("CPU Temp:", strconv.Itoa(int(temp))+"°C", "          \r")
		}
		for {
			time.Sleep(10 * time.Second)
			temp := GetTemp()
			if Logging {
				fmt.Print("CPU Temp:", strconv.Itoa(int(temp))+"°C", "          \r")
			}
			if temp <= LowTemp {
				break
			}
		}
		if Logging {
			fmt.Println("\nCPU Temperature Stable!")
		}
	}

	time.Sleep(300 * time.Millisecond)
}
