package main

import (
	"testing"
)

// insert b
// insert c

func TestFirstGE(t *testing.T) {
	// hardcode instantiate a skiplist of form:
	// 			 		b		->		->		-> 		tail
	// `		 		b		-> 		c		->		tail
	// head		-> 		b		-> 		c		-> 		tail
	//
	var (
		tail = &skipListNode{}
		nodeC = &skipListNode{
			item: Item{
				Key:   "c",
				Value: "c's value",
			},
			nextNodes: []*skipListNode{tail, tail},
		}
		nodeB = &skipListNode{
			item: Item{
				Key:   "b",
				Value: "b's value",
			},
			nextNodes: []*skipListNode{nodeC, nodeC, tail},
		}
		skipList = &skipListOC{
			head: &skipListNode{
				item: Item{},
				nextNodes: []*skipListNode{nodeB},
			},
			tail: tail,
		}
	)

	var firstGEResult *skipListNode
	var updateVector []*nodeLevel

	firstGEResult, _, updateVector = skipList.getFirstGENodeLastLNodeAndUpdateVector("a")
	if firstGEResult.item.Key != "b" {
		t.Errorf("getFirstGENodeLastLNodeAndUpdateVector of 'a' should be 'b', instead it points to node with key: %s\n", firstGEResult.item.Key)
	}
	if updateVector != nil {
		t.Errorf("updateVector upon searching for 'a' should be nil, instead it's: %+v\n", updateVector)
	}
	firstGEResult, _, updateVector = skipList.getFirstGENodeLastLNodeAndUpdateVector("b")
	if firstGEResult.item.Key != "b" {
		t.Errorf("getFirstGENodeLastLNodeAndUpdateVector of 'b' should be 'b', instead it points to node with key: %s\n", firstGEResult.item.Key)
	}
	if updateVector != nil {
		t.Errorf("updateVector upon searching for 'b' should be nil, instead it's: %+v\n", updateVector)
	}
	firstGEResult, _, updateVector = skipList.getFirstGENodeLastLNodeAndUpdateVector("c")
	if firstGEResult.item.Key != "c" {
		t.Errorf("getFirstGENodeLastLNodeAndUpdateVector of 'c' should be 'c', instead it points to node with key: %s\n", firstGEResult.item.Key)
	}
	if updateVector[0].node.item.Key != "b" || updateVector[0].level != 2 {
		t.Errorf("updateVector upon searching for 'c' should be [{node: b, level: 2}], instead it's: %+v\n", updateVector)
		//fmt.Printf("loop length: %d\n", len(updateVector))
		//for i, nodeAndLevel := range updateVector {
		//	fmt.Printf("loop index: %d\n", i)
		//	fmt.Printf("nodeLevel: %+v\n", nodeAndLevel)
		//	fmt.Printf("node %s: %+v\n", nodeAndLevel.node.item.Key, nodeAndLevel.node)
		//	fmt.Printf("level: %d\n", nodeAndLevel.level)
		//}
		//fmt.Printf("address of B: %p\n", nodeB)
		//fmt.Printf("address of C: %p\n", nodeC)
		//fmt.Printf("address of tail: %p\n", tail)
	}
	firstGEResult, _, updateVector = skipList.getFirstGENodeLastLNodeAndUpdateVector("d")
	if firstGEResult != skipList.tail {
		t.Errorf("getFirstGENodeLastLNodeAndUpdateVector of 'd' should be tail, instead it's: %+v\n", firstGEResult)
	}
	if updateVector[0].node.item.Key != "b" || updateVector[0].level != 2 || updateVector[1].node.item.Key != "c" || updateVector[1].level != 1 {
		t.Errorf("updateVector upon searching for 'd' should be [{node: b, level: 2}, {node: c, level: 1}], instead it's: %+v\n", updateVector)
		//fmt.Printf("loop length: %d\n", len(updateVector))
		//for i, nodeAndLevel := range updateVector {
		//	fmt.Printf("loop index: %d\n", i)
		//	fmt.Printf("nodeLevel: %+v\n", nodeAndLevel)
		//	fmt.Printf("node %s: %+v\n", nodeAndLevel.node.item.Key, nodeAndLevel.node)
		//	fmt.Printf("level: %d\n", nodeAndLevel.level)
		//}
		//fmt.Printf("address of B: %p\n", nodeB)
		//fmt.Printf("address of C: %p\n", nodeC)
		//fmt.Printf("address of tail: %p\n", tail)
	}
}

func TestGet (t *testing.T) {
	// hardcode instantiate a skiplist of form:
	// 			 		b		->		->		-> 		tail
	// `		 		b		-> 		c		->		tail
	// head		-> 		b		-> 		c		-> 		tail
	//
	var (
		tail = &skipListNode{}
		nodeC = &skipListNode{
			item: Item{
				Key:   "c",
				Value: "c's value",
			},
			nextNodes: []*skipListNode{tail},
		}
		nodeB = &skipListNode{
			item: Item{
				Key:   "b",
				Value: "b's value",
			},
			nextNodes: []*skipListNode{nodeC, nodeC, tail},
		}
		skipList = &skipListOC{
			head: &skipListNode{
				item: Item{},
				nextNodes: []*skipListNode{nodeB},
			},
			tail: tail,
		}
	)

	var getResultString string
	var getResultBool bool

	getResultString, getResultBool = skipList.Get("a")
	if getResultString != "" || getResultBool != false {
		t.Errorf("Get('a') should return empty string, false; instead got: %s, %t\n", getResultString, getResultBool)
	}
	getResultString, getResultBool = skipList.Get("b")
	if getResultString != "b's value" || getResultBool != true {
		t.Errorf("Get('b') should return 'b's value', true; instead got: %s, %t\n", getResultString, getResultBool)
	}
	getResultString, getResultBool = skipList.Get("c")
	if getResultString != "c's value" || getResultBool != true {
		t.Errorf("Get('c') should return 'c's value', true; instead got: %s, %t\n", getResultString, getResultBool)
	}
	getResultString, getResultBool = skipList.Get("d")
	if getResultString != "" || getResultBool != false {
		t.Errorf("Get('d') should return empty string, false; instead got: %s, %t\n", getResultString, getResultBool)
	}
}

func TestInsertNodeAndResetPointerAtLevel(t *testing.T) {
	var (
		skipList = skipListOC{}
		tail *skipListNode = &skipListNode{}
		existingNode *skipListNode = &skipListNode{item: Item{"Existing", "ExistingValue"}, nextNodes: []*skipListNode{tail}}
		level int = 0
		newNode *skipListNode = &skipListNode{item: Item{"Existing", "ExistingValue"}, nextNodes: []*skipListNode{tail}}
	)
	skipList.insertNodeAndResetPointerAtLevel(existingNode, level, newNode)
	if existingNode.nextNodes[0] != newNode {
		t.Errorf("Existing node at level %d should now point to newNode, instead points to: %+v\n", level, existingNode.nextNodes[0])
	}
	if newNode.nextNodes[0] != tail {
		t.Errorf("Existing node at level %d should now point to tail, instead points to: %+v\n", level, existingNode.nextNodes[0])
	}
}

// code taken and adapted from https://github.com/Bradfield-CSI-2020/rohan-coursework/blob/master/2_storage_retrieval/3_skip_lists/skip_list_oc_test.go

//func TestPut(t *testing.T) {
//
//	rand.Seed(0) // let's use predictable numbers for testing
//
//	list := newSkipListOC()
//
//	// insert two items in order
//	{
//		putResult := list.Put("b", "b")
//		if putResult != true {
//			t.Errorf("Successful insertion of a new key should return true. Got %t\n", putResult)
//		}
//
//		putResult = list.Put("c", "c")
//		if putResult != true {
//			t.Errorf("Successful insertion of a new key should return true. Got %t\n", putResult)
//		}
//	}
//
//	// and the bottom level should be ordered
//	{
//		firstItemKey := list.Header.Next[0].Key
//		if firstItemKey != "b" {
//			t.Errorf("'b' should be the first key at the base layer. Saw %s\n", firstItemKey)
//		}
//
//		secondItemKey := list.Header.Next[0].Next[0].Key
//		if secondItemKey != "c" {
//			t.Errorf("'c' should be the second key at the base layer. Saw %s\n", secondItemKey)
//		}
//	}
//
//	// insert another item out of order
//	{
//		// it should insert correctly
//		putResult := list.Put("a", "a")
//		if putResult != true {
//			t.Errorf("Successful insertion of a new key should return true. Got %t\n", putResult)
//		}
//
//		// and the keys should be in order, respecting the new insertion
//		itemKey0 := list.Header.Next[0].Key
//		itemKey1 := list.Header.Next[0].Next[0].Key
//		itemKey2 := list.Header.Next[0].Next[0].Next[0].Key
//		if itemKey0 != "a" || itemKey1 != "b" || itemKey2 != "c" {
//			order := fmt.Sprintf("(%s, %s, %s)", itemKey0, itemKey1, itemKey2)
//			t.Errorf("Expecting items in order (a, b, c) but found %s\n", order)
//		}
//	}
//
//	// insert an existing item
//	{
//		// it should insert correctly
//		putResult := list.Put("a", "x")
//		if putResult != false {
//			t.Errorf("Update of an extant key should return false. Got %t\n", putResult)
//		}
//
//		// and the keys should be in order, respecting the new values
//		itemVal0 := list.Header.Next[0].Val
//		itemVal1 := list.Header.Next[0].Next[0].Val
//		itemVal2 := list.Header.Next[0].Next[0].Next[0].Val
//		if itemVal0 != "x" || itemVal1 != "b" || itemVal2 != "c" {
//			order := fmt.Sprintf("(%s, %s, %s)", itemVal0, itemVal1, itemVal2)
//			t.Errorf("Expecting items in order (x, b, c) but found %s\n", order)
//		}
//	}
//}
