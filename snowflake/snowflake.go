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

var enablePreGen *bool
var lazy *bool

var instanceLock sync.Mutex

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

// ID sealed for more operations
type ID int64

func (id ID) Int64() int64 {
	return int64(id)
}

func (id ID) String() string {
	return strconv.FormatInt(int64(id), 10)
}

type pool struct {
	ids      *queue.Queue
	min      int
	max      int
	step     int
	interval int
	lock     sync.Mutex
}

// no lock !!! for better pre-generation
func (w *worker) nextId() ID {
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
	id := (now-epoch)<<timerIdBitsMoveLen | (w.workerId << workerIdBitsMoveLen) | (w.number)
	return ID(id)
}

func getWorker() *worker {
	if thisWorker == nil {
		if instanceLock.TryLock() {
			defer instanceLock.Unlock()
			if thisWorker == nil {
				var i int64 = 0
				for i < maxWorkerId {
					i++
					if rdb.TryLock(redisRegKey+strconv.Itoa(int(i)), strconv.FormatInt(time.Now().UnixMilli(), 10), redisTtl) {
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

func NextId() ID {
	idPool.lock.Lock()
	defer idPool.lock.Unlock()
	if *enablePreGen {
		if idPool.ids.GetSize() <= idPool.min {
			go getWorker().preGen()
			if idPool.ids.IsEmpty() {
				return getWorker().nextId()
			}
		}
		id, _ := idPool.ids.Pop()
		return id.(ID)
	}
	return getWorker().nextId()
}

// lazy get worker id
func overTimeWatcher() {
	time.Sleep(syncTtl * time.Second)
	if *lazy {
		thisWorker = nil
		plog.Info("Snowflake worker has became inactive after " + strconv.Itoa(syncTtl) + " seconds for lazy mode")
		return
	}
	getWorker()
}

func Init() {
	// link to conf
	conf := config.GetConfig().Snowflake
	idPool = &pool{
		ids:      queue.New(),
		min:      conf.PreGen.Min,
		max:      conf.PreGen.Max,
		step:     conf.PreGen.Steps,
		interval: conf.PreGen.Interval,
	}
	lazy = &conf.Lazy
	enablePreGen = &conf.PreGen.Enable
	// pre handle
	if idPool.step > idPool.max-idPool.min {
		idPool.step = (idPool.max - idPool.min) / 2
	}
	if *lazy {
		plog.Info("Snowflake worker running in lazy mode")
		return
	}
	// init
	getWorker()
	if *enablePreGen {
		if conf.Lazy {
			getWorker().preGen()
		} else {
			go poolWatchDog()
		}
	}
}

func (w *worker) preGen() {
	idPool.lock.Lock()
	defer idPool.lock.Unlock()
	i := idPool.ids.GetSize()
	j := 0
	for i < idPool.max && j < idPool.step {
		idPool.ids.Push(w.nextId())
		j++
	}
	plog.Info("added " + strconv.Itoa(j) + " ids to pool, current storage : " + strconv.Itoa(idPool.ids.GetSize()))
}

func poolWatchDog() {
	for *enablePreGen && !*lazy {
		if idPool.ids.GetSize() < idPool.max {
			getWorker().preGen()
		} else {
			time.Sleep(time.Duration(idPool.interval) * time.Second)
		}
	}
}
