package crypt

import (
	"github.com/bits-and-blooms/bitset"
	blst "github.com/supranational/blst/bindings/go"
)

type PublicKey = blst.P1Affine
type Signature = blst.P2Affine
type AggregateSignature = blst.P2Aggregate
type AggregatePublicKey = blst.P1Aggregate

type PartialSig Signature

type QuorumCert struct {
	view uint

	BlockHash    hash
	participants bitset.BitSet
	cert         *AggregateSignature
}

func (qc *QuorumCert) addPartialSig(sigma PartialSig, r uint) {
	qc.participants.Set(r)
	qc.cert.Add(sigma, false)
}

func (qc *QuorumCert) verifyCert(pks []PublicKey) {
	affine := qc.cert.ToAffine()

	memberKeys := make([]PublicKey, 0, qc.participants.Count())

	for i := range qc.participants.EachSet() {
		memberKeys = append(memberKeys, pks[i])
	}

	affine.FastAggregateVerify(true, pks, qc.BlockHash)
}
