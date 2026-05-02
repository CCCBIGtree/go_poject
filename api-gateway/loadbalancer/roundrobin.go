package loadbalancer

import "sync"

type RoundRobin struct {
	mu      sync.Mutex
	counter map[string]int
}

func NewRoundRobin() *RoundRobin {
	return &RoundRobin{counter: make(map[string]int)}
}

func (r *RoundRobin) Next(key string, nodes []string) string {
	if len(nodes) == 0 {
		return ""
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	idx := r.counter[key] % len(nodes)
	r.counter[key] = (r.counter[key] + 1) % len(nodes)
	return nodes[idx]
}
