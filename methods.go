package updatingIdentifier

import (
	"math/rand"
	"strconv"
	"time"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

// runs in own goroutine
func (i *Data) watch() {
	for {
		time.Sleep(i.halfTimeout)
		i.generateValues()
	}
}

const (
	maxGenerateNextRawValueAttempts = 4
	valueFormatBase                 = 13
)

func (i *Data) generateValues() {
	currValue, currRawValue := i.valueAndRawValue()
	attempts := 0
work:
	nextRawValue := rand.Uint64()

	if nextRawValue == currRawValue {
		if attempts++; attempts == maxGenerateNextRawValueAttempts {
			nextRawValue = currRawValue + uint64(rand.Uint32())
		} else {
			goto work
		}
	}

	nextValue := strconv.FormatUint(nextRawValue, valueFormatBase)

	defer i.mutex.Unlock()
	i.mutex.Lock()

	i.prevValue = currValue
	i.rawValue = nextRawValue
	i.value = nextValue
}

func (i *Data) valueAndRawValue() (value string, rawValue uint64) {
	defer i.mutex.RUnlock()
	i.mutex.RLock()

	value, rawValue = i.value, i.rawValue
	return
}

func (i *Data) Check(value string) bool {
	defer i.mutex.RUnlock()
	i.mutex.RLock()

	return value == i.value || value == i.prevValue
}

func (i *Data) String() string {
	defer i.mutex.RUnlock()
	i.mutex.RLock()

	return i.value
}
