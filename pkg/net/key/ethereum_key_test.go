package key

import (
	"crypto/ecdsa"
	"crypto/rand"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/secp256k1"
	"github.com/pborman/uuid"
)

// `geth` uses `go-ethereum` library to generate key with secp256k1 curve.
// `libp2p` does not recognize this curve and when it comes to creating peer's
// ID or deserializing the key, operation fails with unrecognized curve error.
//
// To overcome this limitation we rewrite ECDSA key referencing secp256k1 curve
// from `go-ethereum` library into a new key instance supported by `libP2P` and
// referencing secp256k1 curve from `btcsuite` used by `libp2p` under the hood.
// This happens in `toLibp2pKey` function.
//
// As long as all curve parameters are the same, this operation is valid.
// This test ensures that secp256k1 from `go-ethereum` and secp256k1 from
// `btcsuite` are the same. If this test starts to fails, we'll need to revisit
// how the key is ported from one instance to another in `toLibp2pKey` function.
func TestSameCurveAsEthereum(t *testing.T) {
	ethereumKey, err := generateEthereumKey()
	if err != nil {
		t.Fatal(err)
	}

	libp2pKey, _ := EthereumKeyToNetworkKey(ethereumKey)

	ethereumCurve := ethereumKey.PrivateKey.Curve.Params()
	libp2pCurve := libp2pKey.Curve.Params()

	if ethereumCurve.P.Cmp(libp2pCurve.P) != 0 {
		t.Errorf(
			"unexpected P\nexpected: %v\nactual: %v",
			ethereumCurve.P,
			libp2pCurve.P,
		)
	}

	if ethereumCurve.N.Cmp(libp2pCurve.N) != 0 {
		t.Errorf(
			"unexpected N\nexpected: %v\nactual: %v",
			ethereumCurve.N,
			libp2pCurve.N,
		)
	}

	if ethereumCurve.B.Cmp(libp2pCurve.B) != 0 {
		t.Errorf(
			"unexpected B\nexpected: %v\nactual: %v",
			ethereumCurve.B,
			libp2pCurve.B,
		)
	}

	if ethereumCurve.Gx.Cmp(libp2pCurve.Gx) != 0 {
		t.Errorf(
			"unexpected Gx\nexpected: %v\nactual: %v",
			ethereumCurve.Gx,
			libp2pCurve.Gx,
		)
	}

	if ethereumCurve.Gy.Cmp(libp2pCurve.Gy) != 0 {
		t.Errorf(
			"unexpected Gy\nexpected: %v\nactual: %v",
			ethereumCurve.Gy,
			libp2pCurve.Gy,
		)
	}

	if ethereumCurve.BitSize != libp2pCurve.BitSize {
		t.Errorf(
			"unexpected BitSize\nexpected: %v\nactual: %v",
			ethereumCurve.BitSize,
			libp2pCurve.BitSize,
		)
	}
}

func TestSameKeyAsEthereum(t *testing.T) {
	ethereumKey, err := generateEthereumKey()
	if err != nil {
		t.Fatal(err)
	}

	libp2pPrivKey, libp2pPubKey := EthereumKeyToNetworkKey(ethereumKey)

	if ethereumKey.PrivateKey.D.Cmp(libp2pPrivKey.D) != 0 {
		t.Errorf(
			"unexpected D\nexpected: %v\nactual: %v",
			ethereumKey.PrivateKey.D,
			libp2pPrivKey.D,
		)
	}

	if ethereumKey.PrivateKey.PublicKey.X.Cmp(libp2pPubKey.X) != 0 {
		t.Errorf(
			"unexpected X\nexpected: %v\nactual: %v",
			ethereumKey.PrivateKey.PublicKey.X,
			libp2pPubKey.X,
		)
	}

	if ethereumKey.PrivateKey.PublicKey.Y.Cmp(libp2pPubKey.Y) != 0 {
		t.Errorf(
			"unexpected Y\nexpected: %v\nactual: %v",
			ethereumKey.PrivateKey.PublicKey.Y,
			libp2pPubKey.Y,
		)
	}
}

func TestNetworkPubKeyToAddress(t *testing.T) {
	ethereumKey, err := generateEthereumKey()
	if err != nil {
		t.Fatal(err)
	}

	ethAddress := crypto.PubkeyToAddress(ethereumKey.PrivateKey.PublicKey).String()

	_, pubKey := EthereumKeyToNetworkKey(ethereumKey)
	libp2pAddress := NetworkPubKeyToEthAddress(pubKey)

	if ethAddress != libp2pAddress {
		t.Errorf(
			"unexpected address\nexpected: %v\nactual: %v",
			ethAddress,
			libp2pAddress,
		)
	}
}

func generateEthereumKey() (*keystore.Key, error) {
	ethCurve := secp256k1.S256()

	ethereumKey, err := ecdsa.GenerateKey(ethCurve, rand.Reader)
	if err != nil {
		return nil, err
	}

	id := uuid.NewRandom()

	return &keystore.Key{
		Id:         id,
		Address:    crypto.PubkeyToAddress(ethereumKey.PublicKey),
		PrivateKey: ethereumKey,
	}, nil
}