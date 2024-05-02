package content

import (
	"fmt"
	"sync"
	"time"
)

type LRUCache struct {
	capacity int
	// here key is of int type
	cache map[int]*Node
	head  *Node
	tail  *Node
	mutex sync.Mutex
}

type Node struct {
	key    int
	value  any
	expiry time.Time
	prev   *Node
	next   *Node
}

func (this *LRUCache) timerGoRoutine() {
	// single timer implementation.
	timer := time.NewTimer(1 * time.Second)

	for {
		select {
		case <-timer.C:
			timer.Reset(1 * time.Second)
			this.mutex.Lock()
			for key, node := range this.cache {
				if time.Now().After(node.expiry) {
					this.removeNode(node)
					delete(this.cache, key)
					fmt.Printf("[*] Removed node for the key %d\n", key)
				}
			}
			this.mutex.Unlock()
		}
	}
}

func Constructor(capacity int) LRUCache {
	return LRUCache{
		capacity: capacity,
		cache:    make(map[int]*Node),
		head:     nil,
		tail:     nil,
		mutex:    sync.Mutex{},
	}
}

func (this *LRUCache) Get(key int) any {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	if node, ok := this.cache[key]; ok {
		// Move the node to the front of the list
		this.moveAtop(node)
		return node.value
	}
	return -1
}

func (this *LRUCache) Set(key int, value any, expiry int) {
	go this.timerGoRoutine()
	this.mutex.Lock()
	defer this.mutex.Unlock()
	expiry_t := time.Now().Add(time.Duration(expiry) * time.Second)
	// if key already present, update the value.
	if node, ok := this.cache[key]; ok {
		node.value = value
		node.expiry = expiry_t
		this.addToLRUcache(node)
	} else {
		// If the key doesn't exist, create a new node
		node := &Node{key: key, value: value, expiry: expiry_t}
		// Add the new node to the front of the list
		this.addToLRUcache(node)
		this.cache[key] = node
		// If the cache exceeds capacity, remove the least recently used node
		if len(this.cache) > this.capacity {
			delete(this.cache, this.tail.key)
			this.removeLeastUsed()
		}
	}

	// commenting this out based on single routine execution
	// go func() {
	// 	// after the expiry seconds of value initialisation, the value will be deleted.
	// 	fmt.Printf("[*] Removing key %d after %d seconds.\n", key, expiry)
	// 	time.Sleep(time.Duration(expiry) * time.Second)
	// 	this.mutex.Lock()
	// 	defer this.mutex.Unlock()
	// 	if node, ok := this.cache[key]; ok {
	// 		this.removeNode(node)
	// 		delete(this.cache, key)
	// 		fmt.Printf("[*] Removed node for the key %d\n", key)
	// 	}
	// }()
}

func (this *LRUCache) moveAtop(node *Node) {
	if node == this.head {
		return
	}
	this.removeNode(node)
	this.addToLRUcache(node)
}

func (this *LRUCache) removeNode(node *Node) {
	if node.prev != nil {
		node.prev.next = node.next
	} else {
		this.head = node.next
	}
	if node.next != nil {
		node.next.prev = node.prev
	} else {
		this.tail = node.prev
	}
}

func (this *LRUCache) addToLRUcache(node *Node) {
	node.prev = nil
	node.next = this.head
	if this.head != nil {
		this.head.prev = node
	}
	this.head = node
	if this.tail == nil {
		this.tail = node
	}
}

func (this *LRUCache) removeLeastUsed() {
	if this.tail == nil {
		return
	}
	if this.tail.prev != nil {
		this.tail.prev.next = nil
	} else {
		this.head = nil
	}
	this.tail = this.tail.prev
}
