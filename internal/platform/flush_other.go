//go:build !linux && !darwin

package platform

func flushTTYInput(int) {}
