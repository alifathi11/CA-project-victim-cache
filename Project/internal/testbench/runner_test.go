package testbench

import (
	"testing"

	"victimcacheproject/internal/benchmark"
	"victimcacheproject/internal/config"
)

func testSuiteConfig(base config.Config) benchmark.SuiteConfig {
	return benchmark.SuiteConfig{
		BlockSizeBytes:  base.BlockSizeBytes,
		L1SizeBytes:     base.L1SizeBytes,
		L2SizeBytes:     base.L2SizeBytes,
		L2Associativity: base.L2Associativity,
		VictimEntries:   base.VictimEntries,
		NumberOfBlocks:  4,
		Repetitions:     8,
		AccessSizeBytes: 8,
	}
}

func TestCompleteSuitePassesAllChecks(t *testing.T) {
	base := config.Default()
	scenarios, err := benchmark.GenerateSuite(testSuiteConfig(base))
	if err != nil {
		t.Fatal(err)
	}
	architectures, err := Architectures(PolicyBoth)
	if err != nil {
		t.Fatal(err)
	}
	results, err := RunSuite(base, scenarios, architectures)
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != len(scenarios)*len(architectures) {
		t.Fatalf("results=%d, want %d", len(results), len(scenarios)*len(architectures))
	}
	for _, check := range ValidateResults(results) {
		if !check.Passed {
			t.Errorf("failed check %q: %s", check.Name, check.Detail)
		}
	}
}

func TestMixedTraceExercisesEveryLevel(t *testing.T) {
	base := config.Default()
	scenario, err := benchmark.GenerateScenario(benchmark.TraceMixed, testSuiteConfig(base))
	if err != nil {
		t.Fatal(err)
	}
	result, err := RunCase(base, scenario, Architecture{
		Name:          "full-fifo",
		Topology:      config.TopologyFull,
		VictimEnabled: true,
		VictimPolicy:  config.ReplacementFIFO,
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.Stats.L1Hits == 0 || result.Stats.VictimHits == 0 || result.Stats.L2Hits == 0 || result.Stats.MemoryAccesses == 0 {
		t.Fatalf("mixed trace did not reach every level: L1=%d VC=%d L2=%d MEM=%d",
			result.Stats.L1Hits, result.Stats.VictimHits, result.Stats.L2Hits, result.Stats.MemoryAccesses)
	}
}

func TestConflictVictimImprovesBaseline(t *testing.T) {
	base := config.Default()
	scenario, err := benchmark.GenerateScenario(benchmark.TraceConflict, testSuiteConfig(base))
	if err != nil {
		t.Fatal(err)
	}
	baseline, err := RunCase(base, scenario, Architecture{Name: "l1-l2", Topology: config.TopologyL1L2})
	if err != nil {
		t.Fatal(err)
	}
	full, err := RunCase(base, scenario, Architecture{Name: "full-fifo", Topology: config.TopologyFull, VictimEnabled: true, VictimPolicy: config.ReplacementFIFO})
	if err != nil {
		t.Fatal(err)
	}
	if full.Stats.VictimHits == 0 {
		t.Fatal("expected victim hits")
	}
	if full.Stats.L2ReadRequests >= baseline.Stats.L2ReadRequests {
		t.Fatalf("L2 reads full=%d baseline=%d", full.Stats.L2ReadRequests, baseline.Stats.L2ReadRequests)
	}
	if full.Stats.TotalCycles >= baseline.Stats.TotalCycles {
		t.Fatalf("cycles full=%d baseline=%d", full.Stats.TotalCycles, baseline.Stats.TotalCycles)
	}
}

func TestWritebackTraceExercisesDirtyEviction(t *testing.T) {
	base := config.Default()
	scenario, err := benchmark.GenerateScenario(benchmark.TraceWriteback, testSuiteConfig(base))
	if err != nil {
		t.Fatal(err)
	}
	result, err := RunCase(base, scenario, Architecture{Name: "l1-l2", Topology: config.TopologyL1L2})
	if err != nil {
		t.Fatal(err)
	}
	if result.Stats.L2WriteRequests == 0 {
		t.Fatal("expected at least one dirty L2 writeback")
	}
}
