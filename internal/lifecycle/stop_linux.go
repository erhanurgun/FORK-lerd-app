package lifecycle

// BatchStopContainers is a no-op on Linux — systemd stops containers via unit
// deactivation so individual StopUnit calls are efficient and non-blocking.
func BatchStopContainers(_ []string) {}

// StopPodmanMachine is a no-op on Linux — Podman runs natively without a VM.
func StopPodmanMachine() {}
