// Package membership contains code that implements the Random Beacon Group
// Selection protocol as described in
// http://docs.keep.network/random-beacon.pdf
package membership

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"

	"github.com/btcsuite/btcd/btcec"
)

// Ticket is a message containing a pseudorandomly generated value, W_k, which is
// used to determine whether a given virtual staker is eligible for the group P
// (the lowest N tickets will be chosen) and a proof of the validity of the value
type Ticket struct {
	Value [sha256.Size]byte // W_k

	Proof *Proof // Proof(Q_j, vs)
}

// Proof consists of the components needed to construct the Ticket's value, and
// also acts as evidence for an accusing challenge against the Ticket's value.
type Proof struct {
	StakerValue        *btcec.PublicKey // Staker-specific value, Q_j
	VirtualStakerIndex uint64           // vs
}

// calculateTicket generates a Ticket from the previous beacon output, the
// staker's ECDSA public key, and the virtual staker index. This function is
// intended to be called in a loop, ranging over the list of virtual stakers.

// See Phase 1 of the Group Selection protocol specification.
func (s *Staker) calculateTicket(
	beaconOutput []byte,
	virtualStakerIndex uint64,
) (*Ticket, error) {
	if virtualStakerIndex > s.VirtualStakers || virtualStakerIndex < 1 {
		return nil, fmt.Errorf(
			"virtualStakerIndex not in range [1, %d]",
			s.VirtualStakers,
		)
	}

	var combinedValue []byte
	combinedValue = append(combinedValue, beaconOutput...)
	combinedValue = append(combinedValue, s.PubKey.SerializeCompressed()...)

	virtualStakerBytes := make([]byte, 64)
	binary.LittleEndian.PutUint64(virtualStakerBytes, virtualStakerIndex)
	combinedValue = append(combinedValue, virtualStakerBytes...)

	value := sha256.Sum256(combinedValue[:])

	return &Ticket{
		Value: value,
		Proof: &Proof{
			StakerValue:        s.PubKey,
			VirtualStakerIndex: virtualStakerIndex,
		},
	}, nil
}

// Tickets implements sort.Interface
type Tickets []*Ticket

// Len is the sort.Interface requirement for Tickets
func (ts Tickets) Len() int {
	return len(ts)
}

// Swap is the sort.Interface requirement for Tickets
func (ts Tickets) Swap(i, j int) {
	ts[i], ts[j] = ts[j], ts[i]
}

// Less is the sort.Interface requirement for Tickets
func (ts Tickets) Less(i, j int) bool {
	iBytes := ts[i].Value[:]
	jBytes := ts[j].Value[:]

	switch bytes.Compare(iBytes, jBytes) {
	case -1:
		return true
	case 1:
		return false
	}

	return true
}