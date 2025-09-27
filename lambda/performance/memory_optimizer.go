package performance

import (
	"runtime"
	"sync"
	"time"
)

// MemoryOptimizer provides memory optimization utilities
type MemoryOptimizer struct {
	// Memory pools for reuse
	userPool    sync.Pool
	requestPool sync.Pool
	
	// Memory metrics
	allocCount    int64
	totalAllocs   int64
	gcCount       int64
	lastGCTime    time.Time
}

// NewMemoryOptimizer creates a new memory optimizer
func NewMemoryOptimizer() *MemoryOptimizer {
	mo := &MemoryOptimizer{
		userPool: sync.Pool{
			New: func() interface{} {
				return &User{}
			},
		},
		requestPool: sync.Pool{
			New: func() interface{} {
				return &Request{}
			},
		},
	}
	
	// Start memory monitoring
	go mo.monitorMemory()
	
	return mo
}

// User represents a user for memory pooling
type User struct {
	Username     string
	PasswordHash string
}

// Request represents a request for memory pooling
type Request struct {
	Path   string
	Body   string
	Headers map[string]string
}

// GetUserFromPool gets a user from the memory pool
func (mo *MemoryOptimizer) GetUserFromPool() *User {
	user := mo.userPool.Get().(*User)
	// Reset user fields
	user.Username = ""
	user.PasswordHash = ""
	return user
}

// ReturnUserToPool returns a user to the memory pool
func (mo *MemoryOptimizer) ReturnUserToPool(user *User) {
	mo.userPool.Put(user)
}

// GetRequestFromPool gets a request from the memory pool
func (mo *MemoryOptimizer) GetRequestFromPool() *Request {
	request := mo.requestPool.Get().(*Request)
	// Reset request fields
	request.Path = ""
	request.Body = ""
	request.Headers = nil
	return request
}

// ReturnRequestToPool returns a request to the memory pool
func (mo *MemoryOptimizer) ReturnRequestToPool(request *Request) {
	mo.requestPool.Put(request)
}

// monitorMemory monitors memory usage and triggers GC if needed
func (mo *MemoryOptimizer) monitorMemory() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()
	
	for range ticker.C {
		var m runtime.MemStats
		runtime.ReadMemStats(&m)
		
		// Check if memory usage is high
		if m.Alloc > 100*1024*1024 { // 100MB
			runtime.GC()
			mo.gcCount++
			mo.lastGCTime = time.Now()
		}
		
		mo.allocCount = int64(m.Alloc)
		mo.totalAllocs = int64(m.TotalAlloc)
	}
}

// GetMemoryStats returns current memory statistics
func (mo *MemoryOptimizer) GetMemoryStats() map[string]interface{} {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	
	return map[string]interface{}{
		"alloc_bytes":      m.Alloc,
		"total_alloc":      m.TotalAlloc,
		"num_gc":           m.NumGC,
		"gc_cpu_fraction":  m.GCCPUFraction,
		"heap_alloc":       m.HeapAlloc,
		"heap_sys":         m.HeapSys,
		"heap_idle":        m.HeapIdle,
		"heap_inuse":       m.HeapInuse,
		"heap_released":    m.HeapReleased,
		"heap_objects":     m.HeapObjects,
		"stack_inuse":      m.StackInuse,
		"stack_sys":        m.StackSys,
		"gc_count":         mo.gcCount,
		"last_gc_time":     mo.lastGCTime,
	}
}

// ForceGC forces garbage collection
func (mo *MemoryOptimizer) ForceGC() {
	runtime.GC()
	mo.gcCount++
	mo.lastGCTime = time.Now()
}

// OptimizeMemoryUsage optimizes memory usage
func (mo *MemoryOptimizer) OptimizeMemoryUsage() {
	// Force GC
	mo.ForceGC()
	
	// Set GC target percentage
	runtime.GC()
	
	// Set memory limit (if supported)
	// This is a placeholder for future memory limit features
}
