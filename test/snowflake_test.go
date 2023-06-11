package test

import (
	"github.com/real-uangi/pu55y/plog"
	"github.com/real-uangi/pu55y/runner"
	"github.com/real-uangi/pu55y/snowflake"
	"strconv"
	"sync"
	"testing"
	"time"
)

var idMap = make(map[int64]int)

var succeed int
var succeedMu sync.Mutex

var smu sync.Mutex

func TestSnowflake(t *testing.T) {
	succeed = 0
	runner.Prepare()
	a := 100
	for i := 0; i < a; i++ {
		go gen(i)
	}
	waitComplete(a)
	getCount()
}

func TestSnowflakeBench(t *testing.T) {
	runner.Prepare()
	start := float64(time.Now().UnixMilli())
	var dosage int = 1e6
	for i := 0; i < dosage; i++ {
		snowflake.NextId()
	}
	total := float64(time.Now().UnixMilli()) - start
	plog.Info("Totally takes " + strconv.Itoa(int(total)) + " ms for " + strconv.Itoa(dosage) + " samples")
	var rate float64
	rate = float64(dosage) / total / 1000
	plog.Info("Rate : " + strconv.FormatFloat(rate, 'f', 2, 64) + "k/ms")
	plog.Info(snowflake.NextId().String())
}

func gen(i int) {
	plog.Warn("group " + strconv.Itoa(i) + " started")
	for j := 0; j < 10000; j++ {
		id := snowflake.NextId().Int64()
		count(id)
	}
	countSucceed()
	plog.Info("group " + strconv.Itoa(i) + " finished!")
}

func countSucceed() {
	if succeedMu.TryLock() {
		defer succeedMu.Unlock()
		succeed++
	}
}

func count(id int64) {
	smu.Lock()
	defer smu.Unlock()
	idMap[id]++
}

func getCount() {
	smu.Lock()
	defer smu.Unlock()
	repeat := 0
	for _, v := range idMap {
		if v > 1 {
			repeat = repeat + v
		}
	}
	plog.Info("repeat " + strconv.Itoa(repeat))
}

func waitComplete(a int) {
	succeedMu.Lock()
	v := succeed
	succeedMu.Unlock()
	if v < a {
		plog.Info("Waiting ...")
		time.Sleep(time.Second)
		waitComplete(a)
	}
}
