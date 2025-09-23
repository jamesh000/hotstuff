package core

type hash [32]byte

type Block struct {
	Id      hash
	Parent  *Block
	Data    string
	Justify *QuorumCert
	Height  int
}

type Message struct {
	sender  int
	msgType int
	bNew    *Block
}

type VoteMessage struct {
	sender     int
	msgType    int
	b          *Block
	partialSig int
}
