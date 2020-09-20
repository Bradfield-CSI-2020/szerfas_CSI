package main

import (
	"math/rand"
)

type skipListNode struct {
	item Item
	nextNodes []*skipListNode
}

type skipListOC struct {
	head *skipListNode
	tail *skipListNode
}

func newSkipListOC() *skipListOC {
	tail := &skipListNode{}
	head := &skipListNode{
		nextNodes: []*skipListNode{tail},
	}
	return &skipListOC{head, tail}
}

type nodeLevel struct {
	node *skipListNode
	level int
}

// todo: test with updateVector and lastL
// Find the first node such that node.item.Key >= key (firstGE)
// as we traverse, kep track of all nodes and their levels that need updating
// also keep track of the last node such that node.item.Key < key (lastL)
// return all three: the firstGE node, the lastL node, and the updateVector
// returns o.tail for firstGE node if no such node exists, and o.head for lastL if no such node exists
func (o *skipListOC) getFirstGENodeLastLNodeAndUpdateVector(key string) (firstGENode *skipListNode, lastLNode *skipListNode, updateVector []*nodeLevel) {
	//fmt.Println("inside getFirstGENodeLastLNodeAndUpdateVector")
	node := o.head.nextNodes[0]
	lastLNode = o.head
	if node == o.tail {
		return node, lastLNode, updateVector
	}

	offset := 0
	for node != o.tail && node.item.Key < key{
		// if next node is less than our key, we progress, but we can't compare against tail
		level := len(node.nextNodes) - 1
		//fmt.Printf("node: %+v\n", node)
		//fmt.Printf("level: %+v\n", level)
		if right := node.nextNodes[level - offset]; right != o.tail && right.item.Key < key {
			//fmt.Printf("moving right because key less than target\n")
			//fmt.Printf("nodeRight: %+v\n", right)
			lastLNode = node
			node = right
			// reset offset for next node
			offset = 0
		} else if level > offset {
			// every time we drop, we want to store in our updateVector, which is used in Put operations re-map pointers of various levels
			updateVector = append(updateVector, &nodeLevel{node, level - offset})
			// drop to next level and restart loop
			offset++
			//fmt.Printf("incrementing offset")
			//fmt.Printf("level of interest now: %d\n", level - offset)
		} else {
			lastLNode = node
			node = right
			// reset offset for next node
			offset = 0
			//fmt.Printf("moving right because at tail\n")
			//fmt.Printf("nodeRight: %+v\n", right)
		}
	}
	return node, lastLNode, updateVector
}

//// todo: test
//// Find the last node such that node.item.Key <= key
//// returns node one before o.tail if no such node exists
//func (o *skipListOC) lastLE(key string) (lastLENode *skipListNode, updateVector []*nodeLevel) {
//	fmt.Println("inside getFirstGENodeLastLNodeAndUpdateVector")
//	node := o.head.nextNodes[0]
//	if node == o.tail {
//		return o.tail, updateVector
//	}
//
//	offset := 0
//	for node != o.tail && node.item.Key < key{
//		level := len(node.nextNodes) - 1
//		// if next node is less than our key, we progress, but we can't compare against tail
//		if right := node.nextNodes[level - offset]; right != o.tail && right.item.Key < key {
//			lastLENode = node
//			node = right
//			// reset offset for next node
//			offset = 0
//		} else if level > offset {
//			// every time we drop, we want to store in our updateVector, which is used in Put operations re-map pointers of various levels
//			updateVector = append(updateVector, &nodeLevel{node, level - offset})
//			// drop to next level and restart loop
//			offset++
//		} else {
//			lastLENode = node
//			node = right
//			// reset offset for next node
//			offset = 0
//		}
//	}
//	return lastLENode, updateVector
//}

func (o *skipListOC) Get(key string) (string, bool) {
	node, _, _ := o.getFirstGENodeLastLNodeAndUpdateVector(key)
	if node != o.tail && node.item.Key == key {
		return node.item.Value, true
	}
	return "", false
}

func (o *skipListOC) insertNodeAndResetPointerAtLevel(existingNode *skipListNode, level int, newNode *skipListNode) {
	//fmt.Printf("\nlen of existingNode.nextNodes slice: %d\n", len(existingNode.nextNodes))
	//fmt.Printf("len of newNode.nextNodes slice: %d\n", len(newNode.nextNodes))
	//fmt.Printf("level: %d\n", level)
	newNode.nextNodes[level] = existingNode.nextNodes[level]
	existingNode.nextNodes[level] = newNode
}

func (o *skipListOC) Put(key, value string) bool {
	firstGENode, lastLNode, updateVector := o.getFirstGENodeLastLNodeAndUpdateVector(key)
	if firstGENode.item.Key == key {
		firstGENode.item.Value = value
		return false
	} else {
		levelNumber := o.getLevelNumber()
		newNode := &skipListNode{
			item: Item{key, value},
			nextNodes: make([]*skipListNode, levelNumber + 1),
		}
		// updateVector is a slice of the nodes and the levels at which their pointers were to a value greater than our target
		// they are listed in descending order of level
		// for every node in this list with a level <= to the level we're now creating for this node,
		// we want to re-map that node's pointers to this node, and point this node's pointers at that level to the next node.

		// For example, if we start with:
		// 			 		b		->		->		-> 		tail
		// `		 		b		-> 		c		->		tail
		// head		-> 		b		-> 		c		-> 		tail

		// and are inserting a key "d" with levelNumber 3, we'll want:
		// 			 		b		->		->		d		-> 		tail
		// `		 		b		-> 		c		d		->		tail
		// head		-> 		b		-> 		c		d		-> 		tail

		// whereas with levelNumber 2, we'll want:
		// 			 		b		->		->		->		-> 		tail
		// `		 		b		-> 		c		d		->		tail
		// head		-> 		b		-> 		c		d		-> 		tail
		var pointersToUpdate []*nodeLevel
		if len(updateVector) > 0 {
			pointersToUpdate = updateVector[len(updateVector) - levelNumber:]
		}
		for _, nodeLevel := range pointersToUpdate {
			o.insertNodeAndResetPointerAtLevel(nodeLevel.node, nodeLevel.level, newNode)
		}
		o.insertNodeAndResetPointerAtLevel(lastLNode, 0, newNode)
		return true
	}
	return false
}

// returns zero-index number of levels to set on a given node
func (o *skipListOC) getLevelNumber() (level int) {
	level = 0
	p := 0.5
	for rand.Float64() <= p {
		level++
	}
	return
}

func (o *skipListOC) Delete(key string) bool {
	return false
}

func (o *skipListOC) RangeScan(startKey, endKey string) Iterator {
	return &skipListOCIterator{}
}

type skipListOCIterator struct {
}

func (iter *skipListOCIterator) Next() {
}

func (iter *skipListOCIterator) Valid() bool {
	return false
}

func (iter *skipListOCIterator) Key() string {
	return ""
}

func (iter *skipListOCIterator) Value() string {
	return ""
}
