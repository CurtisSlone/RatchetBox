# crypto/hpke (Go standard library)

Package hpke implements Hybrid Public Key Encryption (HPKE) as defined in RFC

Import path: crypto/hpke   Toolchain: go1.26.4

package hpke // import "crypto/hpke"

Package hpke implements Hybrid Public Key Encryption (HPKE) as defined in RFC
9180.

[RFC 9180]: https://www.rfc-editor.org/rfc/rfc9180.html

FUNCTIONS

func Open(k PrivateKey, kdf KDF, aead AEAD, info, ciphertext []byte) ([]byte, error)
    Open instantiates a single-use HPKE receiving HPKE context like
    NewRecipient, and then decrypts the provided ciphertext like Recipient.Open
    (with no aad). ciphertext must be the concatenation of the encapsulated key
    and the actual ciphertext.

func Seal(pk PublicKey, kdf KDF, aead AEAD, info, plaintext []byte) ([]byte, error)
    Seal instantiates a single-use HPKE sending HPKE context like NewSender,
    and then encrypts the provided plaintext like Sender.Seal (with no aad).
    Seal returns the concatenation of the encapsulated key and the ciphertext.


TYPES

type AEAD interface {
	ID() uint16

	// Has unexported methods.
}
    The AEAD is one of the three components of an HPKE ciphersuite, implementing
    symmetric encryption.

func AES128GCM() AEAD
    AES128GCM returns an AES-128-GCM AEAD implementation.

func AES256GCM() AEAD
    AES256GCM returns an AES-256-GCM AEAD implementation.

func ChaCha20Poly1305() AEAD
    ChaCha20Poly1305 returns a ChaCha20Poly1305 AEAD implementation.

func ExportOnly() AEAD
    ExportOnly returns a placeholder AEAD implementation that cannot encrypt or
    decrypt, but only export secrets with Sender.Export or Recipient.Export.

    When this is used, Sender.Seal and Recipient.Open return errors.

func NewAEAD(id uint16) (AEAD, error)
    NewAEAD returns the AEAD implementation for the given AEAD ID.

    Applications are encouraged to use specific implementations like AES128GCM
    or ChaCha20Poly1305 instead, unless runtime agility is required.

type KDF interface {
	ID() uint16

	// Has unexported methods.
}
    The KDF is one of the three components of an HPKE ciphersuite, implementing
    key derivation.

func HKDFSHA256() KDF
    HKDFSHA256 returns an HKDF-SHA256 KDF implementation.

func HKDFSHA384() KDF
    HKDFSHA384 returns an HKDF-SHA384 KDF implementation.

func HKDFSHA512() KDF
    HKDFSHA512 returns an HKDF-SHA512 KDF implementation.

func NewKDF(id uint16) (KDF, error)
    NewKDF returns the KDF implementation for the given KDF ID.

    Applications are encouraged to use specific implementations like HKDFSHA256
    instead, unless runtime agility is required.

func SHAKE128() KDF
    SHAKE128 returns a SHAKE128 KDF implementation.

func SHAKE256() KDF
    SHAKE256 returns a SHAKE256 KDF implementation.

type KEM interface {
	// ID returns the HPKE KEM identifier.
	ID() uint16

	// GenerateKey generates a new key pair.
	GenerateKey() (PrivateKey, error)

	// NewPublicKey deserializes a public key from bytes.
	//
	// It implements DeserializePublicKey, as defined in RFC 9180.
	NewPublicKey([]byte) (PublicKey, error)

	// NewPrivateKey deserializes a private key from bytes.
	//
	// It implements DeserializePrivateKey, as defined in RFC 9180.
	NewPrivateKey([]byte) (PrivateKey, error)

	// DeriveKeyPair derives a key pair from the given input keying material.
	//
	// It implements DeriveKeyPair, as defined in RFC 9180.
	DeriveKeyPair(ikm []byte) (PrivateKey, error)

	// Has unexported methods.
}
    A KEM is a Key Encapsulation Mechanism, one of the three components of an
    HPKE ciphersuite.

func DHKEM(curve ecdh.Curve) KEM
    DHKEM returns a KEM implementing one of

      - DHKEM(P-256, HKDF-SHA256)
      - DHKEM(P-384, HKDF-SHA384)
      - DHKEM(P-521, HKDF-SHA512)
      - DHKEM(X25519, HKDF-SHA256)

    depending on curve.

func MLKEM1024() KEM
    MLKEM1024 returns a KEM implementing ML-KEM-1024 from draft-ietf-hpke-pq.

func MLKEM1024P384() KEM
    MLKEM1024P384 returns a KEM implementing MLKEM1024-P384 from
    draft-ietf-hpke-pq.

func MLKEM768() KEM
    MLKEM768 returns a KEM implementing ML-KEM-768 from draft-ietf-hpke-pq.

func MLKEM768P256() KEM
    MLKEM768P256 returns a KEM implementing MLKEM768-P256 from
    draft-ietf-hpke-pq.

func MLKEM768X25519() KEM
    MLKEM768X25519 returns a KEM implementing MLKEM768-X25519 (a.k.a. X-Wing)
    from draft-ietf-hpke-pq.

func NewKEM(id uint16) (KEM, error)
    NewKEM returns the KEM implementation for the given KEM ID.

    Applications are encouraged to use specific implementations like DHKEM or
    MLKEM768X25519 instead, unless runtime agility is required.

type PrivateKey interface {
	// KEM returns the instantiated KEM.
	KEM() KEM

	// Bytes returns the private key as the output of SerializePrivateKey, as
	// defined in RFC 9180.
	//
	// Note that for X25519 this might not match the input to NewPrivateKey.
	// This is a requirement of RFC 9180, Section 7.1.2.
	Bytes() ([]byte, error)

	// PublicKey returns the corresponding PublicKey.
	PublicKey() PublicKey

	// Has unexported methods.
}
    A PrivateKey is an instantiation of a KEM (one of the three components of an
    HPKE ciphersuite) with a decapsulation key (i.e. the secret key).

    A PrivateKey is usually obtained from a method of the corresponding KEM,
    such as KEM.GenerateKey or KEM.NewPrivateKey.

func NewDHKEMPrivateKey(priv ecdh.KeyExchanger) (PrivateKey, error)
    NewDHKEMPrivateKey returns a PrivateKey implementing

      - DHKEM(P-256, HKDF-SHA256)
      - DHKEM(P-384, HKDF-SHA384)
      - DHKEM(P-521, HKDF-SHA512)
      - DHKEM(X25519, HKDF-SHA256)

    depending on the underlying curve of priv (ecdh.X25519, ecdh.P256,
    ecdh.P384, or ecdh.P521).

    This function is meant for applications that already have an instantiated
    crypto/ecdh private key, or another implementation of a ecdh.KeyExchanger
    (e.g. a hardware key). Otherwise, applications should use the
    KEM.NewPrivateKey method of DHKEM.

func NewHybridPrivateKey(pq crypto.Decapsulator, t ecdh.KeyExchanger) (PrivateKey, error)
    NewHybridPrivateKey returns a PrivateKey implementing

      - MLKEM768-X25519 (a.k.a. X-Wing)
      - MLKEM768-P256
      - MLKEM1024-P384

    from draft-ietf-hpke-pq, depending on the underlying curve of t
    (ecdh.X25519, ecdh.P256, or ecdh.P384) and the type of pq.Encapsulator()
    (either *mlkem.EncapsulationKey768 or *mlkem.EncapsulationKey1024).

    This function is meant for applications that already have instantiated
    crypto/ecdh and crypto/mlkem private keys, or another implementation of a
    ecdh.KeyExchanger and crypto.Decapsulator (e.g. a hardware key). Otherwise,
    applications should use the KEM.NewPrivateKey method of e.g. MLKEM768X25519.

func NewMLKEMPrivateKey(priv crypto.Decapsulator) (PrivateKey, error)
    NewMLKEMPrivateKey returns a KEMPrivateKey implementing

      - ML-KEM-768
      - ML-KEM-1024

    from draft-ietf-hpke-pq, depending on the type of priv.Encapsulator()
    (either *mlkem.EncapsulationKey768 or *mlkem.EncapsulationKey1024).

    This function is meant for applications that already have an instantiated
    crypto/mlkem private key. Otherwise, applications should use the
    KEM.NewPrivateKey method of e.g. MLKEM768.

type PublicKey interface {
	// KEM returns the instantiated KEM.
	KEM() KEM

	// Bytes returns the public key as the output of SerializePublicKey.
	Bytes() []byte

	// Has unexported methods.
}
    A PublicKey is an instantiation of a KEM (one of the three components of an
    HPKE ciphersuite) with an encapsulation key (i.e. the public key).

    A PublicKey is usually obtained from a method of the corresponding KEM or
    PrivateKey, such as KEM.NewPublicKey or PrivateKey.PublicKey.

func NewDHKEMPublicKey(pub *ecdh.PublicKey) (PublicKey, error)
    NewDHKEMPublicKey returns a PublicKey implementing

      - DHKEM(P-256, HKDF-SHA256)
      - DHKEM(P-384, HKDF-SHA384)
      - DHKEM(P-521, HKDF-SHA512)
      - DHKEM(X25519, HKDF-SHA256)

    depending on the underlying curve of pub (ecdh.X25519, ecdh.P256, ecdh.P384,
    or ecdh.P521).

    This function is meant for applications that already have an instantiated
    crypto/ecdh public key. Otherwise, applications should use the
    KEM.NewPublicKey method of DHKEM.

func NewHybridPublicKey(pq crypto.Encapsulator, t *ecdh.PublicKey) (PublicKey, error)
    NewHybridPublicKey returns a PublicKey implementing one of

      - MLKEM768-X25519 (a.k.a. X-Wing)
      - MLKEM768-P256
      - MLKEM1024-P384

    from draft-ietf-hpke-pq, depending on the underlying curve of t
    (ecdh.X25519, ecdh.P256, or ecdh.P384) and the type of pq (either
    *mlkem.EncapsulationKey768 or *mlkem.EncapsulationKey1024).

    This function is meant for applications that already have instantiated
    crypto/ecdh and crypto/mlkem public keys. Otherwise, applications should use
    the KEM.NewPublicKey method of e.g. MLKEM768X25519.

func NewMLKEMPublicKey(pub crypto.Encapsulator) (PublicKey, error)
    NewMLKEMPublicKey returns a KEMPublicKey implementing

      - ML-KEM-768
      - ML-KEM-1024

    from draft-ietf-hpke-pq, depending on the type of pub
    (*mlkem.EncapsulationKey768 or *mlkem.EncapsulationKey1024).

    This function is meant for applications that already have an instantiated
    crypto/mlkem public key. Otherwise, applications should use the
    KEM.NewPublicKey method of e.g. MLKEM768.

type Recipient struct {
	// Has unexported fields.
}
    Recipient is a receiving HPKE context. It is instantiated with a
    specific KEM decapsulation key (i.e. the secret key), and it is stateful,
    incrementing the nonce counter for each successful Recipient.Open call.

func NewRecipient(enc []byte, k PrivateKey, kdf KDF, aead AEAD, info []byte) (*Recipient, error)
    NewRecipient returns a receiving HPKE context for the provided KEM
    decapsulation key (i.e. the secret key), and using the ciphersuite defined
    by the combination of KEM, KDF, and AEAD.

    The enc parameter must have been produced by a matching sending HPKE
    context with the corresponding KEM encapsulation key. The info parameter is
    additional public information that must match between sender and recipient.

func (r *Recipient) Export(exporterContext string, length int) ([]byte, error)
    Export produces a secret value derived from the shared key between sender
    and recipient. length must be at most 65,535.

func (r *Recipient) Open(aad, ciphertext []byte) ([]byte, error)
    Open decrypts the provided ciphertext, optionally binding to the additional
    public data aad, or returns an error if decryption fails.

    Open uses incrementing counters for each successful call, and must be called
    in the same order as Seal on the sending side.

type Sender struct {
	// Has unexported fields.
}
    Sender is a sending HPKE context. It is instantiated with a specific KEM
    encapsulation key (i.e. the public key), and it is stateful, incrementing
    the nonce counter for each Sender.Seal call.

func NewSender(pk PublicKey, kdf KDF, aead AEAD, info []byte) (enc []byte, s *Sender, err error)
    NewSender returns a sending HPKE context for the provided KEM encapsulation
    key (i.e. the public key), and using the ciphersuite defined by the
    combination of KEM, KDF, and AEAD.

    The info parameter is additional public information that must match between
    sender and recipient.

    The returned enc ciphertext can be used to instantiate a matching receiving
    HPKE context with the corresponding KEM decapsulation key.

func (s *Sender) Export(exporterContext string, length int) ([]byte, error)
    Export produces a secret value derived from the shared key between sender
    and recipient. length must be at most 65,535.

func (s *Sender) Seal(aad, plaintext []byte) ([]byte, error)
    Seal encrypts the provided plaintext, optionally binding to the additional
    public data aad.

    Seal uses incrementing counters for each call, and Open on the receiving
    side must be called in the same order as Seal.

## idiomatic usage

Encrypt a single message from a sender to a recipient with the HPKE one-shot API, choosing a KEM, KDF, and AEAD; the recipient generates a key pair and the sender Seals to the public key while the recipient Opens. Keywords: Seal Open GenerateKey NewPublicKey PublicKey Bytes MLKEM768X25519 HKDFSHA256 AES256GCM DHKEM KEM KDF AEAD hybrid public key encryption encapsulation associated data ciphertext plaintext.

```go
import (
	"crypto/hpke"
	"fmt"
)

func Example() {
	kem, kdf, aead := hpke.MLKEM768X25519(), hpke.HKDFSHA256(), hpke.AES256GCM()

	// Recipient generates a key pair and publishes the public key.
	recipientPrivateKey, err := kem.GenerateKey()
	if err != nil {
		panic(err)
	}
	publicKeyBytes := recipientPrivateKey.PublicKey().Bytes()

	// Sender seals a message to the recipient's public key.
	publicKey, err := kem.NewPublicKey(publicKeyBytes)
	if err != nil {
		panic(err)
	}
	ciphertext, err := hpke.Seal(publicKey, kdf, aead, []byte("example"), []byte("|-()-|"))
	if err != nil {
		panic(err)
	}

	// Recipient opens the ciphertext.
	plaintext, err := hpke.Open(recipientPrivateKey, kdf, aead, []byte("example"), ciphertext)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Decrypted message: %s\n", plaintext)
	// Output:
	// Decrypted message: |-()-|
}
```
