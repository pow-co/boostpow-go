package boostpow

import (
	"encoding/binary"

	"github.com/libsv/go-bt/v2/bscript"
)

func CreateBoostOutputScript(ba BoostArgs) (*bscript.Script, error) {
	s := &bscript.Script{}

	categoryBytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(categoryBytes, uint32(ba.Category)) // convert the int32 to a byte array
	if err := s.AppendPushData(categoryBytes); err != nil {
		return nil, err
	}

	if err := s.AppendPushData(ba.Content); err != nil {
		return nil, err
	}

	if err := s.AppendPushData(ba.Tag); err != nil {
		return nil, err
	}

	userNonceBytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(userNonceBytes, uint32(ba.UserNonce)) // convert the int32 to a byte array
	if err := s.AppendPushData(userNonceBytes); err != nil {
		return nil, err
	}

	if err := s.AppendPushData(ba.AdditionalData); err != nil {
		return nil, err
	}

	s.AppendOpcodes(

		/* taken from boostpow WP - stopped halfway through
		bscript.OpCAT, bscript.OpSWAP,
		// copy mining pool’s pubkey hash to alt stack. A copy remains on the stack.
		bscript.Op5, bscript.OpROLL, bscript.OpDUP, bscript.OpTOALTSTACK, bscript.OpCAT,
		// expand compact form of target and push to altstack.
		bscript.Op2, bscript.OpPICK, bscript.OpTOALTSTACK,
		// check size of extra_nonce_1
		bscript.Op5, bscript.OpROLL, bscript.OpSIZE, bscript.Op4, bscript.OpEQUALVERIFY, bscript.OpCAT,
		// check size of extra_nonce_2
		bscript.Op6, bscript.OpROLL, bscript.OpSIZE, bscript.Op8,	bscript.OpLESSTHANOREQUAL, bscript.OpVERIFY, bscript.OpCAT,
		// create metadata document and hash it.
		OP_SWAP, OP_CAT, OP_HASH256,
		// target and content + merkleroot to altstack.
		OP_SWAP, OP_TOALTSTACK, OP_CAT, OP_TOALTSTACK,
		// general purpose bits
		push_hex("ff1f00e0"), OP_DUP, OP_INVERT, OP_TOALTSTACK, OP_AND,
		OP_SWAP, OP_FROMALTSTACK, OP_AND, OP_OR,
		OP_FROMALTSTACK, OP_CAT,                                // attach content + merkleroot
		OP_SWAP, OP_SIZE, OP_4, OP_EQUALVERIFY, OP_CAT,         // check size of timestamp.
		OP_FROMALTSTACK, OP_CAT,                                // attach target
		// check size of nonce. Boost POW string is constructed.
		OP_SWAP, OP_SIZE, OP_4, OP_EQUALVERIFY, OP_CAT,
		// Take hash of work string and ensure that it is positive and minimally encoded.
		OP_HASH256, ensure_positive,
		// Get target, transform to expanded form, and ensure that it is positive and minimally encoded.
		OP_FROMALTSTACK, expand_target, ensure_positive,
		// check that the hash of the Boost POW string is less than the target
		OP_LESSTHAN, OP_VERIFY,
		// check that the given address matches the pubkey and check signature.
		OP_DUP, OP_HASH160, OP_FROMALTSTACK, OP_EQUALVERIFY, OP_CHECKSIG
		*/

		// this is taken from the boostpow-js code:
		// CAT SWAP
		bscript.OpCAT,
		bscript.OpSWAP,

		// {5} ROLL DUP TOALTSTACK CAT                // copy mining pool’s pubkey hash to alt stack. A copy remains on the stack.
		bscript.Op5,
		bscript.OpROLL,
		bscript.OpDUP,
		bscript.OpTOALTSTACK,
		bscript.OpCAT,

		// {2} PICK TOALTSTACK                         // copy target and push to altstack.
		bscript.Op2,
		bscript.OpPICK,
		bscript.OpTOALTSTACK,

		// {5} ROLL SIZE {4} EQUALVERIFY CAT          // check size of extra_nonce_1
		bscript.Op5,
		bscript.OpROLL,
		bscript.OpSIZE,
		bscript.Op4,
		bscript.OpEQUALVERIFY,
		bscript.OpCAT,

		// {5} ROLL SIZE {8} EQUALVERIFY CAT          // check size of extra_nonce_2
		bscript.Op5,
		bscript.OpROLL,
		bscript.OpSIZE,
		bscript.Op8,
		bscript.OpEQUALVERIFY,
		bscript.OpCAT,

		// SWAP CAT HASH256                           // create metadata string and hash it.
		bscript.OpSWAP,
		bscript.OpCAT,
		bscript.OpHASH256,

		// SWAP TOALTSTACK CAT CAT                    // target to altstack.
		bscript.OpSWAP,
		bscript.OpTOALTSTACK,
		bscript.OpCAT,
		bscript.OpCAT,

		// SWAP SIZE {4} EQUALVERIFY CAT              // check size of timestamp.
		bscript.OpSWAP,
		bscript.OpSIZE,
		bscript.Op4,
		bscript.OpEQUALVERIFY,
		bscript.OpCAT,

		// FROMALTSTACK CAT                           // attach target
		bscript.OpFROMALTSTACK,
		bscript.OpCAT,

		// SWAP SIZE {4} EQUALVERIFY CAT             // check size of nonce. Boost POW string is constructed.
		bscript.OpSWAP,
		bscript.OpSIZE,
		bscript.Op4,
		bscript.OpEQUALVERIFY,
		bscript.OpCAT,

		// Take hash of work string and ensure that it is positive and minimally encoded.
		bscript.OpHASH256,
		// ensure positive:
		uint8([]byte{00}[0]), bscript.OpCAT, bscript.OpBIN2NUM,

		bscript.OpFROMALTSTACK,
		//expand target:
		bscript.OpSIZE,
		bscript.Op4,
		bscript.OpEQUALVERIFY,
		bscript.Op3,
		bscript.OpSPLIT,
		bscript.OpDUP,
		bscript.OpBIN2NUM,
		bscript.Op3,
		uint8([]byte{33}[0]), // in JS it is: Buffer.from('21', 'hex'),   // actually 33, but in hex
		bscript.OpWITHIN,
		bscript.OpVERIFY,
		bscript.OpTOALTSTACK,
		bscript.OpDUP,
		bscript.OpBIN2NUM,
		bscript.Op0,
		bscript.OpGREATERTHAN,
		bscript.OpVERIFY,
	)

	*s = append(*s, []byte{0000000000000000000000000000000000000000000000000000000000}...)

	s.AppendOpcodes(
		bscript.OpCAT,
		bscript.OpFROMALTSTACK,
		bscript.Op3,
		bscript.OpSUB,
		bscript.Op8,
		bscript.OpMUL,
		bscript.OpRSHIFT,

		// ensure positive:
		uint8([]byte{00}[0]), bscript.OpCAT, bscript.OpBIN2NUM,

		// check that the hash of the Boost POW string is less than the target
		bscript.OpLESSTHAN,
		bscript.OpVERIFY,

		// check that the given address matches the pubkey and check signature.
		// DUP HASH160 FROMALTSTACK EQUALVERIFY CHECKSIG
		bscript.OpDUP,
		bscript.OpHASH160,
		bscript.OpFROMALTSTACK,
		bscript.OpEQUALVERIFY,
		bscript.OpCHECKSIG,
	)

	return nil, nil
}
