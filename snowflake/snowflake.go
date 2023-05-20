package snowflake

import (
	"github.com/real-uangi/pu55y/colletion/queue"
	"github.com/real-uangi/pu55y/config"
	"github.com/real-uangi/pu55y/plog"
	"github.com/real-uangi/pu55y/rdb"
	"strconv"
	"sync"
	"time"
)

var thisWorker *worker
var idPool *pool

var enablePreGen bool = false

var instanceLock sync.Mutex
var poolLock sync.Mutex

// 支持 2 ^ 8 - 1 台机器
// 每一个毫秒支持 2 ^ 9 - 1 个不同的id
const (
	epoch               = int64(1684512000) //2023-05-20
	workerIdBitsMoveLen = uint(8)
	maxWorkerId         = int64(-1 ^ (-1 << workerIdBitsMoveLen))
	timerIdBitsMoveLen  = uint(17)
	maxNumId            = int64(-1 ^ (-1 << 9))
	redisRegKey         = "SNOWFLAKE:KEY:"
	syncTtl             = 600
	redisTtl            = syncTtl + 30
)

type worker struct {
	workerId  int64
	timestamp int64
	number    int64
}

type pool struct {
	ids      *queue.Queue
	min      int
	max      int
	step     int
	interval int
}

// no lock !!! for better pre-generation
func (w *worker) nextId() int64 {
	now := time.Now().UnixMilli()
	if now < w.timestamp {
		plog.Panic("Clock moved backwards")
	}
	if w.timestamp == now {
		w.number++
		if w.number > maxNumId {
			for now <= w.timestamp {
				now = time.Now().UnixMilli()
			}
		}
	} else {
		w.number = 0
		w.timestamp = now
	}
	ID := (now-epoch)<<timerIdBitsMoveLen | (w.workerId << workerIdBitsMoveLen) | (w.number)
	return ID
}

func getWorker() *worker {
	if thisWorker == nil {
		if instanceLock.TryLock() {
			defer instanceLock.Unlock()
			if thisWorker == nil {
				var i int64 = 0
				for i < maxWorkerId {
					if rdb.TryLock(redisRegKey+strconv.Itoa(int(i)), "snowflake", redisTtl) {
						thisWorker = &worker{workerId: i, timestamp: 0, number: 0}
						go overTimeWatcher()
						plog.Info("Snowflake worker [" + strconv.Itoa(int(i)) + "] activated")
						break
					}
				}
				if thisWorker == nil {
					plog.Panic("failed to assign worker an id")
				}
			}
		} else {
			for {
				time.Sleep(5 * time.Millisecond)
				return getWorker()
			}
		}
	}
	return thisWorker
}

func NextId() int64 {
	poolLock.Lock()
	defer poolLock.Unlock()
	if enablePreGen {
		if idPool.ids.GetSize() < idPool.min {
			go getWorker().preGen()
			if idPool.ids.IsEmpty() {
				return getWorker().nextId()
			}
		}
		id, _ := idPool.ids.Pop()
		return id.(int64)
	}
	return getWorker().nextId()
}

// lazy get worker id
func overTimeWatcher() {
	time.Sleep(syncTtl * time.Second)
	if config.GetConfig().Snowflake.Lazy {
		thisWorker = nil
		plog.Info("Snowflake worker has became inactive after " + strconv.Itoa(syncTtl) + " seconds for lazy mode")
		return
	}
	getWorker()
}

func Init() {
	conf := config.GetConfig().Snowflake
	idPool = &pool{
		ids:      queue.New(),
		min:      conf.PreGen.Min,
		max:      conf.PreGen.Max,
		step:     conf.PreGen.Steps,
		interval: conf.PreGen.Interval,
	}
	if conf.Lazy {
		plog.Info("Snowflake worker running in lazy mode")
		return
	}
	getWorker()
	enablePreGen = conf.PreGen.Enable
	if enablePreGen {
		getWorker().preGen()
	}
}

func (w *worker) preGen() {
	poolLock.Lock()
	defer poolLock.Unlock()
	i := idPool.ids.GetSize()
	j := 0
	for i < idPool.max && j < idPool.step {
		idPool.ids.Push(w.nextId())
		j++
	}
}
