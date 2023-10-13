package cputemp

import (
	"strconv"
	"testing"
)

func TestBash(t *testing.T){
	WaitToCool(false)
	if temp := GetTemp(); temp > 64 {
		t.Error("CPU Too Hot!", strconv.Itoa(int(temp))+"°C")
	}
}
