package benchmark

import (
	"reflect"
	"testing"

	"victimcacheproject/internal/model"
)

func suiteTestConfig() SuiteConfig {
	return SuiteConfig{
		BaseAddress:     0,
		BlockSizeBytes:  64,
		L1SizeBytes:     4 * 1024,
		L2SizeBytes:     64 * 1024,
		L2Associativity: 8,
		VictimEntries:   8,
		NumberOfBlocks:  4,
		Repetitions:     8,
		AccessSizeBytes: 8,
	}
}

func TestGenerateRepeatedScenario(t *testing.T) {
	scenario, err := GenerateScenario(TraceRepeated, suiteTestConfig())
	if err != nil {
		t.Fatal(err)
	}
	if len(scenario.Requests) != 32 {
		t.Fatalf("requests=%d, want 32", len(scenario.Requests))
	}
	for i, request := range scenario.Requests {
		if request.Address != 0 {
			t.Fatalf("request %d address=%d, want 0", i, request.Address)
		}
		if request.ID != uint64(i+1) {
			t.Fatalf("request %d ID=%d, want %d", i, request.ID, i+1)
		}
	}
}

func TestGenerateSequentialScenario(t *testing.T) {
	cfg := suiteTestConfig()
	scenario, err := GenerateScenario(TraceSequential, cfg)
	if err != nil {
		t.Fatal(err)
	}
	wantFirstRound := []uint64{0, 64, 128, 192}
	for i, want := range wantFirstRound {
		if scenario.Requests[i].Address != want {
			t.Fatalf("request %d address=%d, want %d", i, scenario.Requests[i].Address, want)
		}
	}
	if scenario.Requests[cfg.NumberOfBlocks].Address != 0 {
		t.Fatalf("second round did not restart at base address")
	}
}

func TestGenerateMixedScenarioIsDeterministicAndReadOnly(t *testing.T) {
	cfg := suiteTestConfig()
	first, err := GenerateScenario(TraceMixed, cfg)
	if err != nil {
		t.Fatal(err)
	}
	second, err := GenerateScenario(TraceMixed, cfg)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(first.Requests, second.Requests) {
		t.Fatal("mixed trace is not deterministic")
	}
	if len(first.Requests) == 0 {
		t.Fatal("mixed trace is empty")
	}
	for _, request := range first.Requests {
		if request.Op != model.Read {
			t.Fatalf("mixed trace contains non-read operation %v", request.Op)
		}
	}
}

func TestGenerateWritebackScenario(t *testing.T) {
	cfg := suiteTestConfig()
	scenario, err := GenerateScenario(TraceWriteback, cfg)
	if err != nil {
		t.Fatal(err)
	}
	if len(scenario.Requests) != cfg.L2Associativity+2 {
		t.Fatalf("requests=%d, want %d", len(scenario.Requests), cfg.L2Associativity+2)
	}
	wantStride := cfg.L2SizeBytes / uint64(cfg.L2Associativity)
	for i, request := range scenario.Requests {
		if request.Op != model.Write {
			t.Fatalf("request %d is not a write", i)
		}
		if i > 0 && request.Address-scenario.Requests[i-1].Address != wantStride {
			t.Fatalf("request spacing=%d, want %d", request.Address-scenario.Requests[i-1].Address, wantStride)
		}
	}
}

func TestGenerateSuiteContainsEveryTrace(t *testing.T) {
	scenarios, err := GenerateSuite(suiteTestConfig())
	if err != nil {
		t.Fatal(err)
	}
	if len(scenarios) != len(AllTraceKinds()) {
		t.Fatalf("scenarios=%d, want %d", len(scenarios), len(AllTraceKinds()))
	}
	for i, kind := range AllTraceKinds() {
		if scenarios[i].Kind != kind {
			t.Fatalf("scenario %d kind=%s, want %s", i, scenarios[i].Kind, kind)
		}
	}
}

func TestParseTraceKindRejectsUnknownValue(t *testing.T) {
	if _, err := ParseTraceKind("not-a-trace"); err == nil {
		t.Fatal("expected an error")
	}
}
