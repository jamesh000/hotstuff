package crypt

import (
	blst "github.com/supranational/blst/bindings/go"
)

type hash [32]byte

const dst string = "BLS_SIG_BLS12381G1_XMD:SHA-256_SSWU_RO_POP_"

type PublicKey = blst.P1Affine
type Signature = blst.P2Affine
type AggregateSignature = blst.P2Aggregate
type AggregatePublicKey = blst.P1Aggregate

type SecretKey = blst.SecretKey
