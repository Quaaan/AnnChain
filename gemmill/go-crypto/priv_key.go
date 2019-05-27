// Copyright 2017 ZhongAn Information Technology Services Co.,Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package crypto

import (
	"bytes"

	"fmt"

	secp256k1 "github.com/btcsuite/btcd/btcec"
	"github.com/dappledger/AnnChain/gemmill/ed25519"
	"github.com/dappledger/AnnChain/gemmill/ed25519/extra25519"
	"github.com/dappledger/AnnChain/gemmill/go-wire"
	gcmn "github.com/dappledger/AnnChain/gemmill/modules/go-common"
)

const (
	// A series of combination of ciphers
	// ZA includes ed25519,ecdsa,ripemd160,keccak256,secretbox

	CryptoTypeZhongAn = "ZA"
)

const (
	PrivKeyLenEd25519   = 64
	PrivKeyLenSecp256k1 = 32
)

// PrivKey is part of PrivAccount and state.PrivValidator.
type PrivKey interface {
	Bytes() []byte
	Sign(msg []byte) Signature
	PubKey() PubKey
	Equals(PrivKey) bool
	KeyString() string
}

// Types of PrivKey implementations
const (
	PrivKeyTypeEd25519   = byte(0x01)
	PrivKeyTypeSecp256k1 = byte(0x02)
)

// for wire.readReflect
var _ = wire.RegisterInterface(
	struct{ PrivKey }{},
	wire.ConcreteType{PrivKeyEd25519{}, PrivKeyTypeEd25519},
	wire.ConcreteType{PrivKeySecp256k1{}, PrivKeyTypeSecp256k1},
)

func PrivKeyFromBytes(privKeyBytes []byte) (privKey PrivKey, err error) {
	err = wire.ReadBinaryBytes(privKeyBytes, &privKey)
	return
}

//-------------------------------------

// Implements PrivKey
type PrivKeyEd25519 [PrivKeyLenEd25519]byte

func (privKey PrivKeyEd25519) Bytes() []byte {
	return wire.BinaryBytes(struct{ PrivKey }{privKey})
}

func (privKey PrivKeyEd25519) Sign(msg []byte) Signature {
	privKeyBytes := [PrivKeyLenEd25519]byte(privKey)
	signatureBytes := ed25519.Sign(&privKeyBytes, msg)
	return SignatureEd25519(*signatureBytes)
}

func (privKey PrivKeyEd25519) PubKey() PubKey {
	privKeyBytes := [PrivKeyLenEd25519]byte(privKey)
	return PubKeyEd25519(*ed25519.MakePublicKey(&privKeyBytes))
}

func (privKey PrivKeyEd25519) Equals(other PrivKey) bool {
	if otherEd, ok := other.(PrivKeyEd25519); ok {
		return bytes.Equal(privKey[:], otherEd[:])
	} else {
		return false
	}
}

func (privKey PrivKeyEd25519) KeyString() string {
	return gcmn.Fmt("%X", privKey[:])
}

func (privKey PrivKeyEd25519) ToCurve25519() *[32]byte {
	keyCurve25519 := new([32]byte)
	privKeyBytes := [PrivKeyLenEd25519]byte(privKey)
	extra25519.PrivateKeyToCurve25519(keyCurve25519, &privKeyBytes)
	return keyCurve25519
}

func (privKey PrivKeyEd25519) String() string {
	return gcmn.Fmt("PrivKeyEd25519{*****}")
}

// Deterministically generates new priv-key bytes from key.
func (privKey PrivKeyEd25519) Generate(index int) PrivKeyEd25519 {
	newBytes := wire.BinarySha256(struct {
		PrivKey [PrivKeyLenEd25519]byte
		Index   int
	}{privKey, index})
	var newKey [PrivKeyLenEd25519]byte
	copy(newKey[:], newBytes)
	return PrivKeyEd25519(newKey)
}

func GenPrivKeyEd25519() PrivKeyEd25519 {
	privKeyBytes := new([PrivKeyLenEd25519]byte)
	copy(privKeyBytes[:32], CRandBytes(32))
	ed25519.MakePublicKey(privKeyBytes)
	return PrivKeyEd25519(*privKeyBytes)
}

// NOTE: secret should be the output of a KDF like bcrypt,
// if it's derived from user input.
func GenPrivKeyEd25519FromSecret(secret []byte) PrivKeyEd25519 {

	privKey32 := Sha256(secret) // Not Ripemd160 because we want 32 bytes.
	privKeyBytes := new([PrivKeyLenEd25519]byte)
	copy(privKeyBytes[:32], privKey32)
	ed25519.MakePublicKey(privKeyBytes)
	return PrivKeyEd25519(*privKeyBytes)
}

//-------------------------------------

// PrivKeySecp256k1 Implements PrivKey
type PrivKeySecp256k1 [PrivKeyLenSecp256k1]byte

func (privKey PrivKeySecp256k1) Bytes() []byte {
	return wire.BinaryBytes(struct{ PrivKey }{privKey})
}

func (privKey PrivKeySecp256k1) Sign(msg []byte) Signature {
	priv__, _ := secp256k1.PrivKeyFromBytes(secp256k1.S256(), privKey[:])
	sig__, err := priv__.Sign(Sha256(msg))
	if err != nil {
		gcmn.PanicSanity(err)
	}

	return SignatureSecp256k1(sig__.Serialize())
}

func (privKey PrivKeySecp256k1) PubKey() PubKey {
	_, pub__ := secp256k1.PrivKeyFromBytes(secp256k1.S256(), privKey[:])
	pub := [PubKeyLenSecp256k1]byte{}
	copy(pub[:], pub__.SerializeUncompressed()[1:])
	return PubKeySecp256k1(pub)
}

func (privKey PrivKeySecp256k1) Equals(other PrivKey) bool {
	if otherSecp, ok := other.(PrivKeySecp256k1); ok {
		return bytes.Equal(privKey[:], otherSecp[:])
	} else {
		return false
	}
}

func (privKey PrivKeySecp256k1) String() string {
	return gcmn.Fmt("PrivKeySecp256k1{*****}")
}

func (privKey PrivKeySecp256k1) KeyString() string {
	return gcmn.Fmt("%X", privKey[:])
}

func GenPrivKeySecp256k1() PrivKeySecp256k1 {
	privKeyBytes := [PrivKeyLenSecp256k1]byte{}
	copy(privKeyBytes[:], CRandBytes(32))
	priv, _ := secp256k1.PrivKeyFromBytes(secp256k1.S256(), privKeyBytes[:])
	copy(privKeyBytes[:], priv.Serialize())
	return PrivKeySecp256k1(privKeyBytes)
}

// NOTE: secret should be the output of a KDF like bcrypt,
// if it's derived from user input.
func GenPrivKeySecp256k1FromSecret(secret []byte) PrivKeySecp256k1 {
	privKey32 := Sha256(secret) // Not Ripemd160 because we want 32 bytes.
	priv, _ := secp256k1.PrivKeyFromBytes(secp256k1.S256(), privKey32)
	privKeyBytes := [PrivKeyLenSecp256k1]byte{}
	copy(privKeyBytes[:], priv.Serialize())
	return PrivKeySecp256k1(privKeyBytes)
}

func GenPrivkeyByBytes(cryptoType string, data []byte) (PrivKey, error) {
	var privkey PrivKey
	switch cryptoType {
	case CryptoTypeZhongAn:
		var ed PrivKeyEd25519
		copy(ed[:], data)
		privkey = ed
	default:
		return nil, fmt.Errorf("Unknow crypto type")
	}
	return privkey, nil
}

func GenPrivkeyByType(cryptoType string) (PrivKey, error) {
	var privkey PrivKey
	switch cryptoType {
	case CryptoTypeZhongAn:
		privkey = GenPrivKeyEd25519()
	default:
		return nil, fmt.Errorf("Unknow crypto type")
	}
	return privkey, nil
}
