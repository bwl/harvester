package debug

import (
	"runtime"
	"sync"
	"time"
)

type Metrics struct {
	mu sync.RWMutex

	// Performance metrics
	FrameTime    time.Duration
	FrameRate    float64
	FrameHistory []time.Duration
	maxFrames    int

	// Memory metrics
	HeapAlloc     uint64
	HeapSys       uint64
	GCPauseTotal  time.Duration
	NumGoroutines int

	// Game specific metrics
	EntityCount int
	SystemTime  map[string]time.Duration
	RenderTime  time.Duration
	UpdateTime  time.Duration

	// Timing helpers
	lastFrameTime time.Time
	frameCount    uint64
}

var globalMetrics = &Metrics{
	maxFrames:  60,
	SystemTime: make(map[string]time.Duration),
}

func (m *Metrics) RecordFrame() {
	m.mu.Lock()
	defer m.mu.Unlock()

	now := time.Now()
	if !m.lastFrameTime.IsZero() {
		frameTime := now.Sub(m.lastFrameTime)
		m.FrameTime = frameTime

		// Update frame history
		if len(m.FrameHistory) >= m.maxFrames {
			copy(m.FrameHistory, m.FrameHistory[1:])
			m.FrameHistory[m.maxFrames-1] = frameTime
		} else {
			m.FrameHistory = append(m.FrameHistory, frameTime)
		}

		// Calculate average frame rate
		if len(m.FrameHistory) > 0 {
			var total time.Duration
			for _, ft := range m.FrameHistory {
				total += ft
			}
			avgFrameTime := total / time.Duration(len(m.FrameHistory))
			if avgFrameTime > 0 {
				m.FrameRate = 1.0 / avgFrameTime.Seconds()
			}
		}
	}

	m.lastFrameTime = now
	m.frameCount++
}

func (m *Metrics) RecordSystemTime(system string, duration time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.SystemTime[system] = duration
}

func (m *Metrics) RecordRenderTime(duration time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.RenderTime = duration
}

func (m *Metrics) RecordUpdateTime(duration time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.UpdateTime = duration
}

func (m *Metrics) SetEntityCount(count int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.EntityCount = count
}

func (m *Metrics) UpdateMemoryStats() {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	m.mu.Lock()
	defer m.mu.Unlock()

	m.HeapAlloc = memStats.HeapAlloc
	m.HeapSys = memStats.HeapSys
	m.GCPauseTotal = time.Duration(memStats.PauseTotalNs)
	m.NumGoroutines = runtime.NumGoroutine()
}

func (m *Metrics) GetSnapshot() MetricsSnapshot {
	m.mu.RLock()
	defer m.mu.RUnlock()

	systemTime := make(map[string]time.Duration)
	for k, v := range m.SystemTime {
		systemTime[k] = v
	}

	frameHistory := make([]time.Duration, len(m.FrameHistory))
	copy(frameHistory, m.FrameHistory)

	return MetricsSnapshot{
		FrameTime:     m.FrameTime,
		FrameRate:     m.FrameRate,
		FrameHistory:  frameHistory,
		HeapAlloc:     m.HeapAlloc,
		HeapSys:       m.HeapSys,
		GCPauseTotal:  m.GCPauseTotal,
		NumGoroutines: m.NumGoroutines,
		EntityCount:   m.EntityCount,
		SystemTime:    systemTime,
		RenderTime:    m.RenderTime,
		UpdateTime:    m.UpdateTime,
		FrameCount:    m.frameCount,
	}
}

type MetricsSnapshot struct {
	FrameTime     time.Duration
	FrameRate     float64
	FrameHistory  []time.Duration
	HeapAlloc     uint64
	HeapSys       uint64
	GCPauseTotal  time.Duration
	NumGoroutines int
	EntityCount   int
	SystemTime    map[string]time.Duration
	RenderTime    time.Duration
	UpdateTime    time.Duration
	FrameCount    uint64
}

// Global functions
func RecordFrame() {
	globalMetrics.RecordFrame()
}

func RecordSystemTime(system string, duration time.Duration) {
	globalMetrics.RecordSystemTime(system, duration)
}

func RecordRenderTime(duration time.Duration) {
	globalMetrics.RecordRenderTime(duration)
}

func RecordUpdateTime(duration time.Duration) {
	globalMetrics.RecordUpdateTime(duration)
}

func SetEntityCount(count int) {
	globalMetrics.SetEntityCount(count)
}

func UpdateMemoryStats() {
	globalMetrics.UpdateMemoryStats()
}

func GetMetricsSnapshot() MetricsSnapshot {
	return globalMetrics.GetSnapshot()
}

// SystemTimer helps measure system execution time
type SystemTimer struct {
	system string
	start  time.Time
}

func StartSystemTimer(system string) *SystemTimer {
	return &SystemTimer{
		system: system,
		start:  time.Now(),
	}
}

func (st *SystemTimer) Stop() {
	if st != nil {
		duration := time.Since(st.start)
		RecordSystemTime(st.system, duration)
		Debugf("performance", "System %s took %v", st.system, duration)
	}
}
