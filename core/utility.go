package core

import (
	"crypto/sha3"
	"log"

	"github.com/gogo/protobuf/proto"
	pb "github.com/jamesh000/hotstuff/pb"
)

func createLeaf(p *Block, data string, qc *QuorumCert) *Block {
	b := new(Block)

	b.Parent = p
	b.Data = data
	b.Justify = qc
	b.Height = p.Height + 1

	b.hashPayload()

	return b
}

func createQC(b *Block) *QuorumCert {
	qc := new(QuorumCert)

	qc.Node = b

	qc.hashPayload()

	return qc
}

func sameChain(new, old *Block) bool {
	for b := new; b.Height >= old.Height; b = b.Parent {
		if b == old {
			return true
		}
	}
	return false
}

func (b *Block) hashPayload() {
	payload := pb.BlockPayload{
		Parent:  b.Parent.Id[:],
		Data:    b.Data,
		Justify: b.Justify,
	}

	plBytes, err := proto.Marshal(&payload)

	if err != nil {
		log.Println(err)
	}

	b.Id = sha3.Sum256(plBytes)
}

func (qc *QuorumCert) hashPayload() {
	payload := pb.QuorumCertPayload{
		Blockid: qc.Node.Id[:],
	}

	plBytes, err := proto.Marshal(&payload)

	if err != nil {
		log.Println(err)
	}

	qc.Id = sha3.Sum256(plBytes)
}
