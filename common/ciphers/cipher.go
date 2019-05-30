package ciphers

import (
	"bytes"
	"crypto/cipher"
	"crypto/rand"
	"github.com/rc452860/vnet/common/ciphers/block"
	"github.com/rc452860/vnet/utils/bytesx"
	"github.com/sirupsen/logrus"
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
	OP_ENCRYPT = 0
	OP_DECRYPT = 1
)

type SSCipher struct {
	OP int
	cipher.Stream
	cipher.AEAD
	cipher.BlockMode
	IV             []byte
	AEADWriteNonce []byte
	AEADReadNonce  []byte
	AEADReadBuffer *bytes.Buffer
	Method         string
}

func NewCipher(op int, method string, key, iv []byte) (*SSCipher, error) {
	if cp := stream.GetStreamCipher(method); cp != nil {
		cpInstance, err := cp.NewStream(key, iv, op)
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
		cpInstance, err := cp.NewAEAD(key, iv, op)
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
	if cp := block.GetBlockCipher(method); cp != nil {
		cpInstance, err := cp.NewBlock(key, iv, op)
		if err != nil {
			return nil, err
		}
		return &SSCipher{
			BlockMode: cpInstance,
			Method:    method,
			IV:        iv,
			OP:        op,
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
	if s.BlockMode != nil {
		dst := make([]byte, len(src))
		s.BlockMode.CryptBlocks(dst, src)
		return dst, nil
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
	if s.BlockMode != nil {
		dst := make([]byte, len(ciphertext))
		s.BlockMode.CryptBlocks(dst, ciphertext)
		return dst, nil
	}
	return nil, errors.WithStack(errors.New("unknow ciphers"))
}

type Encryptor struct {
	Key          []byte
	KeyStr       string
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
	result.KeyStr = key
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

	if cp := block.GetBlockCipher(method); cp != nil {
		result.IVOut = make([]byte, cp.IVLen())
		if _, err := io.ReadFull(rand.Reader, result.IVOut); err != nil {
			return nil, err
		}
		result.Key = evpBytesToKey(key, cp.KeyLen())
		result.IVLen = cp.IVLen()
	}

	result.EncodeCipher, err = NewCipher(OP_ENCRYPT, method, result.Key, result.IVOut)
	result.IVBuf = new(bytes.Buffer)
	return result, err
}

func NewEncryptorWithIv(method, key string, iv []byte) (result *Encryptor, err error) {
	result = new(Encryptor)
	result.IVSent = false
	result.Method = method
	// if method is stream then
	if cp := stream.GetStreamCipher(method); cp != nil {
		result.Key = evpBytesToKey(key, cp.KeyLen())
		result.IVLen = cp.IVLen()
	}

	if cp := aead.GetAEADCipher(method); cp != nil {
		result.Key = evpBytesToKey(key, cp.KeySize())
		result.IVLen = cp.SaltSize()
	}

	if cp := block.GetBlockCipher(method); cp != nil {
		result.Key = evpBytesToKey(key, cp.KeyLen())
		result.IVLen = cp.IVLen()
	}

	result.IVOut = iv[:result.IVLen]
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
		e.IVSent = true
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

	if e.IVBuf.Len() <= e.IVLen {
		e.IVBuf.Write(ciphertext)
	}

	if e.IVBuf.Len() > e.IVLen {
		buf := e.IVBuf.Bytes()
		decipherIV := buf[:e.IVLen]
		e.IVIn = decipherIV
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

func (e *Encryptor) EncryptAll(src,iv []byte) (result []byte, err error) {
	encrypter,err := NewEncryptorWithIv(e.Method,e.KeyStr,iv)
	if err != nil{
		return nil,err
	}
	return encrypter.Encrypt(src)
}

func (e *Encryptor) DecryptAll(ciphertext []byte) (result,iv []byte, err error) {
	encrypter,err := NewEncryptor(e.Method,e.KeyStr)
	if err != nil{
		return nil,nil,err
	}
	data,err := encrypter.Decrypt(ciphertext)
	if err != nil{
		return nil,nil,err
	}
	return data,encrypter.IVIn,nil
}

func (e *Encryptor) NewIV()(iv []byte,err error){
	iv = make([]byte, e.IVLen)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil,err
	}
	return iv,err
}

func (e *Encryptor) MustNewIV() []byte{
	iv,err := e.NewIV()
	if err != nil{
		logrus.WithFields(logrus.Fields{
			"error":err,
		}).Error("cipher Encryptor MustNewIV error")
		return nil
	}
	return iv
}