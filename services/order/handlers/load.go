package handlers

import (
	"crypto/sha256"
	"math/rand"
	"time"
)

var MemoryHolder [][]byte

func SimulateCPULoad(iterations int) {
	data := make([]byte, 1024)

	for i := 0; i < iterations; i++ {
		rand.Read(data)
		hash := sha256.Sum256(data)
		_ = hash
	}
}

func SimulateMemoryLoad(sizeMB int) {
	mem := make([][]byte, sizeMB)

	for i := 0; i < sizeMB; i++ {
		mem[i] = make([]byte, 1024*1024)
	}

	MemoryHolder = append(MemoryHolder, mem...)
}

func SimulateDelay(ms int) {
	time.Sleep(time.Duration(ms) * time.Millisecond)
}
