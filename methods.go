package updatingIdentifier

import (
	"math/rand"
	"strconv"
	"time"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

// watch runs in own goroutine
func (i *Data) watch() {
	for {
		time.Sleep(i.halfTimeout)
		i.generateValues()
	}
}

const (
	maxGenerateNextRawValueAttempts = 16
	valueFormatBase                 = 13
)

// generateValues should not be called concurrently
func (i *Data) generateValues() {
	currValue, currRawValue := i.valueAndRawValue()
	var nextRawValue uint64

	for attempts := 1; ; attempts++ {
		if nextRawValue = rand.Uint64(); nextRawValue != currRawValue {
			break
		}

		if attempts == maxGenerateNextRawValueAttempts {
			nextRawValue = currRawValue + uint64(rand.Uint32()) // fallback to guaranteed increment
			break
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
