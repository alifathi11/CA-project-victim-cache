# Simulation adapter

`adapter.go` provides the complete synchronous reference runner used for deterministic correctness tests.

For an actual Akita event-driven implementation, keep the same `System`, request model, cache policies, and metrics, but wrap them with Akita components, ports, messages, and scheduled events. Detailed steps are in the repository root file `AKITA_INTEGRATION.md`.
