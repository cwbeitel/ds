package main 

import (
	"fmt"
	"math/rand"
	"time"
)

// The Hash object is an implementation of the hash map data
// structure using arrays, mapping string keys to string values.
type Hash struct {

	// The data array where pointers to hash nodes are stored,
	// into which the hash function indexes.
	data []*hashNode

	// The size of the allocated hash
	size int

}

// hashNode is a node in a Hash data structure, with multiple nodes
// linked in a chain for hash buckets where a collision has occurred.
type hashNode struct {

	// The data stored at this node in the hash
	data []string

	// A pointer to the next node in the nodes chained at this 
	// bucket of the hash
	next *hashNode

	// A poitner to the previous node in the nodes chained at
	// this bucket of the hash
	last *hashNode

}

// NewHash creates a new empty string to string hash map.
func NewHash() Hash {
	size := 1234 // TODO - determine this more carefully
	return Hash{make([]*hashNode, size), size}
}

// hash computes a hash key for a specified string and number of buckets
// using the classic rotating hash function.
//
// TODO: explain more about what is happening here.
func (h *Hash) hash(k string) int {
	key := []byte(k)
	l := len(key)
	v := byte(l)
	for i:=0; i<l; i++{
		v = (((v<<4)^(v >> 28))^key[i])
	}
	return int(v) % h.size
}


// Get obtains the entry in the hash corresponding to the provided
// key.
//
// Get will return a non-nil error if the key is not present in
// the hash, otherwise returning the value in the hash for that 
// key.
func (h *Hash) Get(key string) (string, error) {
	
	// Compute the hash index
	i := h.hash(key)
	n := h.data[i]
	if n == nil {
		return "", fmt.Errorf("no entry at hash bucket %d for key %s", i, key)
	}

	// For each node in the chain of all nodes that have collided at this
	// bucket, look for a node that matches the key.
	for {
		// If the key at this node matches the specified key, return the value
		if n.data[0] == key {
			return n.data[1], nil 
		}
		if n.next != nil {
			n = n.next // since there is a next node, proceed to it
		} else {
			return "", fmt.Errorf("hash: did not find value in hash with the specified key: %s", key)
		}
	}

}

// Set sets the provided value for the provided key in the hash.
func (h *Hash) Set(key, value string) error {

	// Compute the hash index
	i := h.hash(key)
	n := h.data[i]
	if n == nil {
		// there are not yet any entries for this bucket, insert a new one
		h.data[i] = &hashNode{[]string{key, value}, nil, nil}
		return nil
	}

	// For each node in the chain of all nodes that have collided at this
	// bucket, look through we either find that the node already exists or
	// we run out of nodes and insert this key,value as a new node
	for {

		// If the key at this node matches the specified key it already exists in the hash, return
		if (*n).data[0] == key {
			(*n).data[1] = value
			return nil
		}
		if n.next != nil {
			n = n.next // since there is a next node, proceed to it
		} else {
			// Does not exist, add it
			n.next = &hashNode{[]string{key, value}, n, nil}
		}
	}

}

// Delete removes an entry from a hash table given the specified key
// of that entry.
func (h *Hash) Delete(key string) error {

	// Compute the hash index
	i := h.hash(key)
	n := h.data[i]
	if n == nil {
		// there are not yet any entries for this bucket so there will not be one to delete
		return fmt.Errorf("hash: tried to delete entry from hash that was not found to be present in the hash: %s", key)
	}

	// For each node in the chain of all nodes that have collided at this
	// bucket, look for a node that matches the key.
	for {

		// If the key at this node matches the specified key, return the value
		if n.data[0] == key {
			if n.last == nil {
				h.data[i] = n.next
			} else {
				last := n.last
				next := n.next
				next.last = last
				last.next = next
			}
		}
		if n.next != nil {
			n = n.next // since there is a next node, proceed to it
		} else {
			return fmt.Errorf("hash: did not find value in hash with the specified key: %s", key)
		}
	}
}

// randStringArray generates an array of a specified number of strings of a 
// specified length.
func randStringArray(size, strlen int) []string {

	rand.Seed(time.Now().UnixNano())
	runes := []rune("abcdefghijklmnopqrstuvwxyz")
	rsa := make([]string, size)

	for i:=0; i<size; i++ {
		b := make([]rune, strlen)
		for j, _ := range b {
			b[j] = runes[rand.Intn(len(runes))]
		}
		rsa[i] = string(b)
	}

	return rsa

}

// scoreHashClustering computes a measure of the clustering of a set of
// hash keys generated by a hash function for a specified number of
// buckets and keys.
func scoreHashClustering(vals []int, numBuckets int) float64 {

	// Tally the number of things of each key type
	bucketCounts := make([]int, numBuckets)
	for i := 0; i < len(vals); i++ {
		bucketCounts[vals[i]] += 1.0
	}

	// Get the sum of the squares of bucket sizes
	tot := 0
	for i :=0; i<len(bucketCounts); i++ {
		tot += bucketCounts[i]*bucketCounts[i]
	}

	fmt.Println(tot)

	// return the clustering score as the bucket size mean minus
	// the expected number of keys per bucket (numKeys/numBuckets)
	return float64(tot)/float64(len(vals)) - float64(len(vals))*float64(1/numBuckets)

}

