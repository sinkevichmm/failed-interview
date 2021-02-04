package incrementor

import (
	"math"
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func TestNewIncrementor(t *testing.T) {
	inc := NewIncrementor()

	assert.Equal(t, uint32(math.MaxUint32), inc.GetMaximumValue())
	assert.Equal(t, uint32(0), inc.GetNumber())
}

func payloader(i *Incrementor, wg *sync.WaitGroup) {
	defer wg.Done()
	time.Sleep(time.Duration(rand.Int31n(100)) * time.Millisecond)
	i.IncrementNumber()
}

func TestIncrementNumber(t *testing.T) {
	inc := NewIncrementor()

	var wg sync.WaitGroup

	count := 10

	for i := 0; i < count; i++ {
		wg.Add(1)

		go payloader(inc, &wg)
	}

	wg.Wait()
	assert.Equal(t, uint32(count), inc.GetNumber())
}

func TestSetMaximumValue(t *testing.T) {
	inc := NewIncrementor()

	for i := 0; i < 100; i++ {
		inc.IncrementNumber()
	}

	maxValue := uint32(10)

	inc.SetMaximumValue(maxValue)

	assert.Equal(t, uint32(0), inc.GetNumber())

	for i := 0; i < int(maxValue); i++ {
		inc.IncrementNumber()
	}

	assert.Equal(t, maxValue, inc.GetNumber())

	inc.IncrementNumber()

	assert.Equal(t, uint32(0), inc.GetNumber())

	inc.SetMaximumValue(4294967295)

	inc.i = 4294967295

	inc.IncrementNumber()

	assert.Equal(t, uint32(0), inc.GetNumber())
}
