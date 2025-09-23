package core

import (
	pb "github.com/jamesh000/hotstuff/pb"
)

func toProto(b *Block) *pb.BlockProto {
	s := new(pb.BlockProto)

	s.Id = b.Id[:]
	s.Payload = &pb.BlockPayload{
		Parent:  b.Parent.Id[:],
		Data:    b.Data,
		Justify: b.Justify.Id[:],
	}
	s.Height = uint64(b.Height)

	return s
}

func fromProto(bp *pb.BlockProto, storage BlockStorage) *Block {
	b := new(Block)

	parentId := hash(bp.Payload.Parent)

	b.Id = hash(bp.Id)
	b.Parent = storage[parentId]
	b.Justify = bp.Payload.Justify

}
