package skiplist

import (
	"math/rand"
)

const (
	maxLevels = 12
	branching = 4
)

type skipListOCNode struct {
	item  Item
	next  [maxLevels]*skipListOCNode
	level int
}

type skipListOC struct {
	head   *skipListOCNode
	levels int
}

func newSkipListOC() *skipListOC {
	head := &skipListOCNode{
		level: maxLevels,
	}
	return &skipListOC{
		head:   head,
		levels: 1,
	}
}

// Find the latest node in each level such that node.item.Key < key
func (o *skipListOC) findPrev(key string) [maxLevels]*skipListOCNode {
	var result [maxLevels]*skipListOCNode
	node := o.head
	for i := o.levels - 1; i >= 0; i-- {
		for node.next[i] != nil && node.next[i].item.Key < key {
			node = node.next[i]
		}
		result[i] = node
	}
	return result
}

func nextOrNil(prev [maxLevels]*skipListOCNode, key string) *skipListOCNode {
	if prev[0].next[0] != nil && prev[0].next[0].item.Key == key {
		return prev[0].next[0]
	}
	return nil
}

func randomLevel() int {
	result := 1
	for result < maxLevels && rand.Intn(branching) == 0 {
		result++
	}
	return result
}

func (o *skipListOC) Get(key string) (string, bool) {
	prev := o.findPrev(key)
	if node := nextOrNil(prev, key); node != nil {
		return node.item.Value, true
	}
	return "", false
}

func (o *skipListOC) Put(key, value string) bool {
	prev := o.findPrev(key)
	if node := nextOrNil(prev, key); node != nil {
		node.item.Value = value
		return false
	} else {
		level := randomLevel()
		node := &skipListOCNode{
			item:  Item{key, value},
			level: level,
		}
		if level > o.levels {
			for i := level - 1; i >= o.levels; i-- {
				o.head.next[i] = node
			}
			for i := o.levels - 1; i >= 0; i-- {
				node.next[i] = prev[i].next[i]
				prev[i].next[i] = node
			}
			o.levels = level
		} else {
			for i := level - 1; i >= 0; i-- {
				node.next[i] = prev[i].next[i]
				prev[i].next[i] = node
			}
		}
		return true
	}
}

func (o *skipListOC) Delete(key string) bool {
	prev := o.findPrev(key)
	if node := nextOrNil(prev, key); node != nil {
		for i := node.level - 1; i >= 0; i-- {
			prev[i].next[i] = node.next[i]
		}
		return true
	} else {
		return false
	}
}

func (o *skipListOC) RangeScan(startKey, endKey string) Iterator {
	prev := o.findPrev(startKey)
	node := prev[0].next[0]
	return &skipListOCIterator{o, node, startKey, endKey}
}

type skipListOCIterator struct {
	o                *skipListOC
	node             *skipListOCNode
	startKey, endKey string
}

func (iter *skipListOCIterator) Next() {
	iter.node = iter.node.next[0]
}

func (iter *skipListOCIterator) Valid() bool {
	return iter.node != nil && iter.node.item.Key <= iter.endKey
}

func (iter *skipListOCIterator) Key() string {
	return iter.node.item.Key
}

func (iter *skipListOCIterator) Value() string {
	return iter.node.item.Value
}
