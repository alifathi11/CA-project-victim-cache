package metrics

type Stats struct {
	TotalRequests   uint64
	L1Hits          uint64
	L1Misses        uint64
	VictimHits      uint64
	VictimMisses    uint64
	L2Hits          uint64
	L2Misses        uint64
	MemoryAccesses  uint64
	L2ReadRequests  uint64
	L2WriteRequests uint64
	VictimSwaps     uint64
	TotalCycles     uint64
}

func (s Stats) L1HitRate() float64 {
	if s.TotalRequests == 0 {
		return 0
	}
	return float64(s.L1Hits) / float64(s.TotalRequests)
}

func (s Stats) VictimHitRate() float64 {
	accesses := s.VictimHits + s.VictimMisses
	if accesses == 0 {
		return 0
	}
	return float64(s.VictimHits) / float64(accesses)
}
