package main

import "fmt"

type Result struct {
	Context    Context     `json:"context"`
	Benchmarks []Benchmark `json:"benchmarks"`
}

type Context struct {
	Date              string  `json:"date"`
	Hostname          string  `json:"host_name"`
	Executable        string  `json:"executable"`
	NumCPUs           int     `json:"num_cpus"`
	MHzPerCPU         int     `json:"mhz_per_cpu"`
	CPUScalingEnabled bool    `json:"cpu_scaling_enabled"`
	Caches            []Cache `json:"caches"`
}

func (c Context) Equals(other Context) error {
	if len(c.Caches) != len(other.Caches) {
		return fmt.Errorf("different number of CPU caches: %d vs %d", len(c.Caches), len(other.Caches))
	}
	for i := range c.Caches {
		if err := c.Caches[i].Equals(other.Caches[i]); err != nil {
			return err
		}
	}
	if c.NumCPUs != other.NumCPUs {
		return fmt.Errorf("different number of CPUs: %d vs %d", c.NumCPUs, other.NumCPUs)
	}
	if c.MHzPerCPU != other.MHzPerCPU {
		return fmt.Errorf("different MHz/CPU: %d vs %d", c.MHzPerCPU, other.MHzPerCPU)
	}
	if c.CPUScalingEnabled != other.CPUScalingEnabled {
		return fmt.Errorf("different CPU scaling: %v vs %v", c.CPUScalingEnabled, other.CPUScalingEnabled)
	}
	return nil
}

type Cache struct {
	Type       string `json:"type"`
	Level      int    `json:"level"`
	Size       uint64 `json:"size"`
	NumSharing int    `json:"num_sharing"`
}

func (c Cache) Equals(other Cache) error {
	if c.Type != other.Type {
		return fmt.Errorf("different type of CPU cache: %v vs %v", c.Type, other.Type)
	}
	if c.Level != other.Level {
		return fmt.Errorf("different CPU cache level: %d vs %d", c.Level, other.Level)
	}
	if c.Size != other.Size {
		return fmt.Errorf("different CPU cache size: %d vs %d", c.Size, other.Size)
	}
	return nil
}

type Benchmark struct {
	Name            string  `json:"name"`
	RunName         string  `json:"run_name"`
	RunType         string  `json:"run_type"`
	Repetitions     uint64  `json:"repetitions"`
	RepetitionIndex uint64  `json:"repetition_index"`
	Threads         int     `json:"threads"`
	Iterations      uint64  `json:"iterations"`
	RealTime        float64 `json:"real_time"`
	CPUTime         float64 `json:"cpu_time"`
	TimeUnit        string  `json:"time_unit"`
}
