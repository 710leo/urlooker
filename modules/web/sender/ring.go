package sender

import (
	"stathat.com/c/consistent"
)

type ConsistentHashNodeRing struct {
	ring *consistent.Consistent
}

func NewConsistentHashNodeRing(NumberOfReplicas int, nodes []string) *ConsistentHashNodeRing {
	ret := &ConsistentHashNodeRing{ring: consistent.New()}
	ret.SetNumberOfReplicas(NumberOfReplicas)
	ret.SetNodes(nodes)
	return ret
}

func (this *ConsistentHashNodeRing) SetNumberOfReplicas(num int) {
	this.ring.NumberOfReplicas = num
}

func (this *ConsistentHashNodeRing) SetNodes(nodes []string) {
	for _, node := range nodes {
		this.ring.Add(node)
	}
}

func (this *ConsistentHashNodeRing) GetNode(pk string) (string, error) {
	return this.ring.Get(pk)
}
