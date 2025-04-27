package hasher

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

var (
	ErrInvalidHashFormat    = errors.New("invalid hash format")
	ErrIncompatibleVersion  = errors.New("incompatible argon2 version")
	ErrHashingFailed        = errors.New("failed to hash password")
	ErrGetParamsFailed      = errors.New("failed to parse parameters")
	ErrInvalidParams        = errors.New("invalid argon2 parameters")
	ErrDecodingFailed       = errors.New("failed to decode base64 component")
	ErrGenerateSaltFailed   = errors.New("failed to generate salt")
	ErrHashGenerationFailed = errors.New("failed to generate hash key")
)

type Argon2idParams struct {
	Memory      uint32
	Iterations  uint32
	Parallelism uint8
	SaltLength  uint32
	KeyLength   uint32
}

var DefaultParams = Argon2idParams{
	Memory:      64 * 1024,
	Iterations:  3,
	Parallelism: 4,
	SaltLength:  16,
	KeyLength:   32,
}

type Argon2idHasher struct {
	Params Argon2idParams
}

func NewArgon2IDHasher() *Argon2idHasher {
	return &Argon2idHasher{Params: DefaultParams}
}

func (h *Argon2idHasher) Hash(password string) (string, error) {
	salt := make([]byte, h.Params.SaltLength)
	_, err := rand.Read(salt)
	if err != nil {
		return "", fmt.Errorf("%w: %w", ErrGenerateSaltFailed, err)
	}

	hash := argon2.IDKey([]byte(password), salt, h.Params.Iterations, h.Params.Memory, h.Params.Parallelism, h.Params.KeyLength)
	if hash == nil {
		return "", ErrHashGenerationFailed
	}

	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	encodedHash := fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version,
		h.Params.Memory,
		h.Params.Iterations,
		h.Params.Parallelism,
		b64Salt,
		b64Hash,
	)

	return encodedHash, nil
}

func (h *Argon2idHasher) Check(password, encodedHash string) bool {
	params, salt, hash, err := h.decodeHash(encodedHash)
	if err != nil {
		return false
	}

	otherHash := argon2.IDKey([]byte(password), salt, params.Iterations, params.Memory, params.Parallelism, params.KeyLength)
	if otherHash == nil {
		return false
	}

	if uint32(len(hash)) != params.KeyLength || uint32(len(otherHash)) != params.KeyLength {
		return false
	}

	if subtle.ConstantTimeCompare(hash, otherHash) == 1 {
		return true
	}

	return false
}

func (h *Argon2idHasher) decodeHash(encodedHash string) (params *Argon2idParams, salt, hash []byte, err error) {
	vals := strings.Split(encodedHash, "$")
	if len(vals) != 6 {
		return nil, nil, nil, fmt.Errorf("%w: expected 6 parts, got %d", ErrInvalidHashFormat, len(vals))
	}

	if vals[1] != "argon2id" {
		return nil, nil, nil, fmt.Errorf("%w: unsupported algorithm %s", ErrInvalidHashFormat, vals[1])
	}

	var version int
	_, err = fmt.Sscanf(vals[2], "v=%d", &version)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("%w: failed parsing version: %w", ErrInvalidHashFormat, err)
	}
	if version != argon2.Version {
		return nil, nil, nil, fmt.Errorf("%w: expected version %d, got %d", ErrIncompatibleVersion, argon2.Version, version)
	}

	params = &Argon2idParams{}
	_, err = fmt.Sscanf(vals[3], "m=%d,t=%d,p=%d", &params.Memory, &params.Iterations, &params.Parallelism)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("%w: %w", ErrGetParamsFailed, err)
	}

	if params.Memory == 0 || params.Iterations == 0 || params.Parallelism == 0 {
		return nil, nil, nil, ErrInvalidParams
	}

	salt, err = base64.RawStdEncoding.DecodeString(vals[4])
	if err != nil {
		return nil, nil, nil, fmt.Errorf("%w: failed decoding salt: %w", ErrDecodingFailed, err)
	}
	params.SaltLength = uint32(len(salt))

	hash, err = base64.RawStdEncoding.DecodeString(vals[5])
	if err != nil {
		return nil, nil, nil, fmt.Errorf("%w: failed decoding hash: %w", ErrDecodingFailed, err)
	}
	params.KeyLength = uint32(len(hash))

	return params, salt, hash, nil
}
