package snowflake

import (
	"github.com/real-uangi/pu55y/config"
	"github.com/real-uangi/pu55y/plog"
	"github.com/real-uangi/pu55y/rdb"
	"strconv"
	"sync"
	"time"
)

// ID sealed for more operations
type ID int64

const (
	workerBits  uint8 = 10
	numberBits  uint8 = 12
	workerMax   int64 = -1 ^ (-1 << workerBits)
	numberMax   int64 = -1 ^ (-1 << numberBits)
	timeShift         = workerBits + numberBits
	workerShift       = numberBits
	epoch       int64 = 1684512000000
	redisRegKey       = "SNOWFLAKE:KEY:"
)

var (
	instance *Worker
	mu       sync.Mutex
	interval = 3600
)

type Worker struct {
	mu        sync.Mutex
	timestamp int64
	workerId  int64
	number    int64
}

func Init() {
	conf := config.GetConfig().Snowflake
	interval = conf.Interval
	getInstance()
}

func NextId() ID {
	return getInstance().nextId()
}

func getInstance() *Worker {
	mu.Lock()
	defer mu.Unlock()
	if instance == nil {
		instance = newWorker()
	}
	return instance
}

func newWorker() *Worker {
	var i int64
	for i = 0; i < workerMax; i++ {
		if rdb.TryLock(redisRegKey+strconv.Itoa(int(i)), strconv.FormatInt(time.Now().UnixMilli(), 10), interval) {
			plog.Info("Snowflake worker [" + strconv.Itoa(int(i)) + "] activating")
			return &Worker{
				timestamp: 0,
				workerId:  i,
				number:    0,
			}
		}
	}
	plog.Error("Failed to register Snowflake instance")
	panic("Failed to register Snowflake instance")
}

func keepInstanceOn() {
	time.Sleep(time.Duration(interval-60) * time.Second)
	mu.Lock()
	instance = newWorker()
	mu.Unlock()
	keepInstanceOn()
}

func (w *Worker) nextId() ID {
	w.mu.Lock()
	defer w.mu.Unlock()
	now := time.Now().UnixMilli()
	if w.timestamp == now {
		w.number = (w.number + 1) & numberMax
		if w.number == 0 {
			for now <= w.timestamp {
				now = time.Now().UnixMilli()
			}
		}
	} else {
		w.number = 0
	}
	w.timestamp = now
	id := (now-epoch)<<timeShift | (w.workerId << workerShift) | (w.number)
	return ID(id)
}

func (id ID) String() string {
	return strconv.FormatInt(int64(id), 10)
}

func (id ID) Int64() int64 {
	return int64(id)
}
