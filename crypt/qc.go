package crypt

import (
	"github.com/bits-and-blooms/bitset"
)

type QuorumCert struct {
	view uint

	BlockHash    hash
	participants bitset.BitSet
	cert         *AggregateSignature
}

func (qc *QuorumCert) addPartialSig(sigma *Signature, r uint) {
	qc.participants.Set(r)
	qc.cert.Add(sigma, false)
}

func (qc *QuorumCert) verifyCert(pks []*PublicKey) {
	affine := qc.cert.ToAffine()

	memberKeys := make([]*PublicKey, 0, qc.participants.Count())

	for i := range qc.participants.EachSet() {
		memberKeys = append(memberKeys, pks[i])
	}

	affine.FastAggregateVerify(false, pks, qc.BlockHash[:], []byte(dst))
}
