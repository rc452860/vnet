package ciphers

import (
	"bytes"
	"crypto/cipher"
	"crypto/rand"
	"github.com/rc452860/vnet/utils/bytesx"
	"io"

	"github.com/pkg/errors"
	"github.com/rc452860/vnet/common/ciphers/aead"
	"github.com/rc452860/vnet/common/ciphers/stream"
)

const AEAD_MAX_SEGMENT_LENGTH = 0x3FFF

type ICipher interface {
	GetKey() []byte
	GetIv() []byte
	Encrypto([]byte) ([]byte, error)
	Decrypto([]byte) ([]byte, error)
}

const (
	OP_ENCRYPT = 1
	OP_DECRYPT = 2
)

type SSCipher struct {
	OP int
	cipher.Stream
	cipher.AEAD
	IV             []byte
	AEADWriteNonce []byte
	AEADReadNonce  []byte
	AEADReadBuffer *bytes.Buffer
	Method         string
}

func NewCipher(op int, method string, key, iv []byte) (*SSCipher, error) {
	if cp := stream.GetStreamCipher(method); cp != nil {
		cpInstance, err := cp.NewStream(key, iv)
		if err != nil {
			return nil, err
		}
		return &SSCipher{
			Stream: cpInstance,
			Method: method,
			IV:     iv,
			OP:     op,
		}, nil
	}
	if cp := aead.GetAEADCipher(method); cp != nil {
		cpInstance, err := cp.NewAEAD(key, iv)
		if err != nil {
			return nil, err
		}
		return &SSCipher{
			AEAD:           cpInstance,
			Method:         method,
			IV:             iv,
			OP:             op,
			AEADReadBuffer: new(bytes.Buffer),
			AEADWriteNonce: make([]byte, cpInstance.NonceSize()),
			AEADReadNonce:  make([]byte, cpInstance.NonceSize()),
		}, nil
	}
	return nil, errors.WithStack(errors.New("unreachable code"))
}

func (s *SSCipher) Encrypt(src []byte) (result []byte, err error) {
	if s.OP != OP_ENCRYPT {
		return nil, errors.WithStack(errors.New("operation not support."))
	}
	if s.Stream != nil {
		dst := make([]byte, len(src))
		s.Stream.XORKeyStream(dst, src)
		return dst, nil
	}
	if s.AEAD != nil {
		dst := new(bytes.Buffer)
		r := bytes.NewBuffer(src)
		rn := 0
		overHead := s.AEAD.Overhead()
		for {
			buf := make([]byte, 2+overHead+AEAD_MAX_SEGMENT_LENGTH+overHead)
			dataBuf := buf[2+overHead : 2+overHead+AEAD_MAX_SEGMENT_LENGTH]
			rn, err = r.Read(dataBuf)
			if err != nil {
				if err == io.EOF {
					err = nil
				}
				break
			}
			if rn > 0 {
				buf = buf[:2+overHead+rn+overHead]
				dataBuf = dataBuf[:rn]
				buf[0], buf[1] = byte(rn>>8), byte(rn&0xffff)
				s.AEAD.Seal(buf[:0], s.AEADWriteNonce, buf[:2], nil)
				increment(s.AEADWriteNonce)

				s.AEAD.Seal(dataBuf[:0], s.AEADWriteNonce, dataBuf, nil)
				increment(s.AEADWriteNonce)

				_, ew := dst.Write(buf)
				if ew != nil {
					err = ew
					break
				}
			} else {
				break
			}
		}
		return dst.Bytes(), err
	}
	return nil, errors.WithStack(errors.New("unreachable code"))
}

func (s *SSCipher) Decrypt(ciphertext []byte) (result []byte, err error) {
	if s.OP != OP_DECRYPT {
		return nil, errors.WithStack(errors.New("operation not support."))
	}
	if s.Stream != nil {
		buf := make([]byte, len(ciphertext))
		s.Stream.XORKeyStream(buf, ciphertext)
		return buf, nil
	}

	if s.AEAD != nil {
		overHead := s.AEAD.Overhead()
		s.AEADReadBuffer.Write(ciphertext)
		dst := new(bytes.Buffer)
		sizeTmp := make([]byte, 2+overHead)
		payloadTmp := make([]byte, AEAD_MAX_SEGMENT_LENGTH+overHead)
		for {
			if s.AEADReadBuffer.Len() == 0 {
				break
			}
			if s.AEADReadBuffer.Len() < 2+overHead {
				return nil, errors.WithStack(errors.New("buf is too short"))
			}
			_, err := s.AEADReadBuffer.Read(sizeTmp)
			if err != nil {
				return nil, err
			}
			_, err = s.AEAD.Open(sizeTmp[:0], s.AEADReadNonce, sizeTmp, nil)
			if err != nil {
				return nil, err
			}
			increment(s.AEADReadNonce)
			size := (int(sizeTmp[0])<<8 + int(sizeTmp[1])) & AEAD_MAX_SEGMENT_LENGTH
			if s.AEADReadBuffer.Len() < size+overHead {
				return nil, errors.WithStack(errors.New("buf is too short"))
			}
			_, err = io.ReadFull(s.AEADReadBuffer, payloadTmp[:size+overHead])
			if err != nil {
				return nil, errors.WithStack(err)
			}
			result, err := s.AEAD.Open(payloadTmp[:0], s.AEADReadNonce, payloadTmp[:size+overHead], nil)
			increment(s.AEADReadNonce)
			if err != nil {
				return nil, err
			}
			dst.Write(result)
		}
		return dst.Bytes(), nil
	}
	return nil, errors.WithStack(errors.New("unknow ciphers"))
}

type Encryptor struct {
	Key          []byte
	Method       string
	IVOut        []byte
	IVIn         []byte
	IVSent       bool
	IVLen        int
	IVBuf        *bytes.Buffer
	EncodeCipher *SSCipher
	DecodeCipher *SSCipher
}

func NewEncryptor(method, key string) (result *Encryptor, err error) {
	result = new(Encryptor)
	result.IVSent = false
	result.Method = method
	// if method is stream then
	if cp := stream.GetStreamCipher(method); cp != nil {
		result.IVOut = make([]byte, cp.IVLen())
		if _, err := io.ReadFull(rand.Reader, result.IVOut); err != nil {
			return nil, err
		}
		result.Key = evpBytesToKey(key, cp.KeyLen())
		result.IVLen = cp.IVLen()
	}

	if cp := aead.GetAEADCipher(method); cp != nil {
		result.IVOut = make([]byte, cp.SaltSize())
		if _, err := io.ReadFull(rand.Reader, result.IVOut); err != nil {
			return nil, err
		}
		result.Key = evpBytesToKey(key, cp.KeySize())
		result.IVLen = cp.SaltSize()
	}

	result.EncodeCipher, err = NewCipher(OP_ENCRYPT, method, result.Key, result.IVOut)
	result.IVBuf = new(bytes.Buffer)
	return result, err
}

func (e *Encryptor) Encrypt(src []byte) (result []byte, err error) {
	result, err = e.EncodeCipher.Encrypt(src)
	if err != nil {
		return result, err
	}
	if !e.IVSent {
		return bytesx.ContactSlice(e.IVOut, result), nil
	} else {
		return result, nil
	}
}

func (e *Encryptor) Decrypt(ciphertext []byte) (result []byte, err error) {
	if len(ciphertext) == 0 {
		return ciphertext, nil
	}
	if e.DecodeCipher != nil {
		return e.DecodeCipher.Decrypt(ciphertext)
	}

	if e.IVBuf.Len()<e.IVLen{
		e.IVBuf.Write(ciphertext)
	}

	if e.IVBuf.Len() > e.IVLen {
		buf := e.IVBuf.Bytes()
		decipherIV := buf[:e.IVLen]
		e.DecodeCipher, err = NewCipher(OP_DECRYPT, e.Method, e.Key, decipherIV)
		if err != nil {
			return nil, err
		}
		remainBuf := buf[e.IVLen:]
		return e.DecodeCipher.Decrypt(remainBuf)
	} else {
		return []byte{}, nil
	}
}
