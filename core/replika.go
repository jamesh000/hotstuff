package core

import (
	"encoding/hex"
	"log"
)

type BlockStorage map[hash]*Block

type Replika struct {
	votes     []int
	voteCount int
	vHeight   int
	lock      *Block
	exec      *Block
	qcHigh    *QuorumCert
	leaf      *Block

	n            int
	replikaCount int
	blocks       BlockStorage
}

func (r *Replika) update(bStar *Block) {
	b2 := bStar.Justify.Node
	b1 := b2.Justify.Node
	b := b1.Justify.Node

	r.updateQCHigh(bStar.Justify)

	if b1.Height > r.lock.Height {
		r.lock = b1
	}
	if b2.Parent.Id == b1.Id && b1.Parent.Id == b.Id {
		r.onCommit(b)
		r.exec = b
	}
}

func (r *Replika) onCommit(b *Block) {
	if r.exec.Height < b.Height {
		r.onCommit(b.Parent)
		log.Println("Committed block", hex.EncodeToString(b.Id[:]), "with data", b.Data)
	}
}

func (r *Replika) onReceiveProposal(msg Message) bool {
	new := msg.bNew
	if new.Height > r.vHeight && (sameChain(new, r.lock) || new.Justify.Node.Height > r.lock.Height) {
		r.vHeight = new.Height
		r.update(new)

		// vote
		return true
	}
	// no vote
	return false
}

func (r *Replika) onReceiveVote(msg VoteMessage) {
	if r.votes[msg.sender] != -1 {
		return
	}

	r.votes[msg.sender] = msg.partialSig
	r.voteCount++

	if r.voteCount > r.replikaCount-r.replikaCount/3 {
		// actual qc code here
	}
}

func (r *Replika) onPropose(leaf *Block, data string, qcHigh *QuorumCert) Message {
	new := createLeaf(leaf, data, qcHigh)

	r.blocks[new.Id] = new

	return Message{
		sender:  r.n,
		msgType: 1,
		bNew:    new,
	}
}

func (r *Replika) updateQCHigh(qc *QuorumCert) {
	if qc.Node.Height > r.qcHigh.Node.Height {
		r.qcHigh = qc
		r.leaf = qc.Node
	}
}
