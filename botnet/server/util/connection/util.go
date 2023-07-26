package connection

import (
	"sync/atomic"
	"math/rand"
	"time"
)
//*================================ ID GENERATOR===============================
type UniqueIDGenerator struct {
	counter int64
}

func NewUniqueIDGenerator() *UniqueIDGenerator {
	return &UniqueIDGenerator{}
}

// NextID generates the next unique ID
func (g *UniqueIDGenerator) NextID() int64 {
	return atomic.AddInt64(&g.counter, 1)
}
//*================================ END ID GENERATOR===============================

//*================================ Randomize list===============================
func ShuffleList(list []*Client){
	rand.Seed(time.Now().UnixNano())

	for i := len(list) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		list[i], list[j] = list[j], list[i]
	}
}
//*================================ END Randomize List===============================

//*================================ Text Animation===============================
func FillString(retunString string, toLength int) string {
	for {
		lengtString := len(retunString)
		if lengtString < toLength {
			retunString = retunString + ":"
			continue
		}
		break
	}
	return retunString
}