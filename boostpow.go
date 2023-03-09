package boostpow

// BoostArgs arguments needed to create a boost Locking Script
type BoostArgs struct {
	Category       int32
	Content        []byte
	Target         uint32
	Tag            []byte
	UserNonce      uint32
	AdditionalData []byte
}
