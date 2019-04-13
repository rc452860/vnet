package obfs

import (
	"bytes"
	"encoding/base64"
	"github.com/rc452860/vnet/common/ciphers"
	"github.com/rc452860/vnet/utils/bytesx"
	"github.com/rc452860/vnet/utils/randomx"
	"math"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/rc452860/vnet/common/cache"
	"github.com/rc452860/vnet/common/log"
	"github.com/rc452860/vnet/utils/binaryx"
)

const (
	MAX_INT  = (1 << 64) - 1
	MOV_MASK = (1 << (64 - 23)) - 1
)

type XorShift128Plus struct {
	V0 uint64
	V1 uint64
}

func NewXorShift128Plus() *XorShift128Plus {
	return &XorShift128Plus{
		V0: 0,
		V1: 0,
	}
}

func (xs1p *XorShift128Plus) Next() uint64 {
	x := xs1p.V0
	y := xs1p.V1
	xs1p.V0 = y
	x ^= ((x & MOV_MASK) << 23)
	x ^= (y ^ (x >> 17) ^ (y >> 26))
	xs1p.V1 = x
	return (x + y) & MAX_INT
}

func (xs1p *XorShift128Plus) InitFromBin(bin []byte) {
	if len(bin) < 16 {
		bin = conbineToBytes(bin, bytes.Repeat([]byte{byte(0x00)}, 16))
	}
	xs1p.V0 = binaryx.LEBytesToUint64(bin[0:8])
	xs1p.V1 = binaryx.LEBytesToUint64(bin[8:16])
}

func (xs1p *XorShift128Plus) InitFromBinLen(bin []byte, length int) {
	if len(bin) < 16 {
		bin = conbineToBytes(bin, bytes.Repeat([]byte{byte(0x00)}, 16))
	}
	xs1p.V0 = binaryx.LEBytesToUint64(conbineToBytes(binaryx.LEUInt16ToBytes(uint16(length)), bin[2:8]))
	xs1p.V1 = binaryx.LEBytesToUint64(bin[8:16])
	for i := 0; i < 4; i++ {
		xs1p.Next()
	}
}

/*--------------------------AuthBase----------------------------*/

type AuthBase struct {
	Plain
	Method             string
	NoCompatibleMethod string
	Overhead           int
	RawTrans           bool
}

func NewAuthBase(method string) *AuthBase {
	return &AuthBase{
		Plain:    NewPlain(method),
		Method:   method,
		Overhead: 4,
	}
}

func (authBase *AuthBase) GetOverhead(direction bool) int {
	return authBase.Overhead
}

func (authBase *AuthBase) NotMatchReturn(buf []byte) ([]byte, bool) {
	authBase.RawTrans = true
	authBase.Overhead = 0
	if authBase.GetMethod() == authBase.NoCompatibleMethod {
		return bytes.Repeat([]byte{byte('E')}, 2048), false
	}
	return buf, false
}

type ClientQueue struct {
	Front      int
	Back       int
	Alloc      *sync.Map
	Enable     bool
	LastUpdate time.Time
	Ref        int
}

func NewClientQueue(beginID int) *ClientQueue {
	return &ClientQueue{
		Front:      beginID - 64,
		Back:       beginID + 1,
		Alloc:      new(sync.Map),
		Enable:     true,
		LastUpdate: time.Now(),
		Ref:        0,
	}
}

func (c *ClientQueue) Update() {
	c.LastUpdate = time.Now()
}

func (c *ClientQueue) AddRef() {
	c.Ref += 1
}

func (c *ClientQueue) DelRef() {
	if c.Ref > 0 {
		c.Ref -= 1
	}
}

func (c *ClientQueue) IsActive() bool {
	return c.Ref > 0 && time.Now().Sub(c.LastUpdate).Seconds() < 60*10
}

func (c *ClientQueue) ReEnable(connectionID int) {
	c.Enable = true
	c.Front = connectionID - 64
	c.Back = connectionID + 1
	c.Alloc = new(sync.Map)
}

func (c *ClientQueue) Insert(connectionID int) bool {
	if !c.Enable {
		log.Warn("obfs auth: not enable")
		return false
	}
	if !c.IsActive() {
		c.ReEnable(connectionID)
	}
	c.Update()
	if connectionID < c.Front {
		log.Warn("obfs auth: deprecated ID, someone replay attack")
		return false
	}
	if connectionID > c.Front+0x4000 {
		log.Warn("obfs auth: wrong ID")
		return false
	}
	if _, ok := c.Alloc.Load(connectionID); ok {
		log.Warn("obfs auth: deprecated ID, someone replay attack")
		return false
	}
	if c.Back <= connectionID {
		c.Back = connectionID + 1
	}
	c.Alloc.Store(connectionID, 1)
	for {
		if _, ok := c.Alloc.Load(c.Back); !ok || c.Front+0x1000 >= c.Back {
			break
		}
		if _, ok := c.Alloc.Load(c.Front); ok {
			c.Alloc.Delete(c.Front)
		}
		c.Front += 1
	}
	c.AddRef()
	return true
}

type ObfsAuthChainData struct {
	Name         string
	UserID       map[int]*cache.LRU
	LastClientID []byte
	ConnectionID int
	MaxClient    int
	MaxBuffer    int
}

func NewObfsAuthChainData(name string) *ObfsAuthChainData {
	result := &ObfsAuthChainData{
		Name:         name,
		UserID:       make(map[int]*cache.LRU),
		LastClientID: []byte{},
		ConnectionID: 0,
	}
	result.SetMaxClient(64)
	return result
}

func (o *ObfsAuthChainData) Update(userID, clientID, connectionID int) {
	if o.UserID[userID] == nil {
		o.UserID[userID] = cache.NewLruCache(60 * time.Second)
	}
	localClientID := o.UserID[userID]
	var r *ClientQueue = nil
	if localClientID != nil {
		r, _ = localClientID.Get(clientID).(*ClientQueue)
	}
	if r != nil {
		r.Update()
	}
}

func (o *ObfsAuthChainData) SetMaxClient(maxClient int) {
	o.MaxClient = maxClient
	o.MaxBuffer = int(math.Max(float64(maxClient), 1024))
}
func (o *ObfsAuthChainData) Insert(userID, clientID, connectionID int) bool {
	if o.UserID[userID] == nil {
		o.UserID[userID] = cache.NewLruCache(60 * time.Second)
	}
	localClientID := o.UserID[userID]
	var r, _ = localClientID.Get(clientID).(*ClientQueue)
	if r != nil || !r.Enable {
		if localClientID.First() == nil || localClientID.Len() < o.MaxClient {
			if !localClientID.IsExist(clientID) {
				// TODO check
				localClientID.Put(clientID, NewClientQueue(connectionID), 60*time.Second)
			} else {
				localClientID.Get(clientID).(*ClientQueue).ReEnable(connectionID)
			}
			return localClientID.Get(clientID).(*ClientQueue).Insert(connectionID)
		}

		localClientIDFirst := localClientID.First()
		if !localClientID.Get(localClientIDFirst).(*ClientQueue).IsActive() {
			localClientID.Delete(localClientIDFirst)
			if !localClientID.IsExist(clientID) {
				// TODO check
				localClientID.Put(clientID, NewClientQueue(connectionID), 60*time.Second)
			} else {
				localClientID.Get(clientID).(*ClientQueue).ReEnable(connectionID)
			}
			return localClientID.Get(clientID).(*ClientQueue).Insert(connectionID)
		}

		log.Warn("%s: no inactive client", o.Name)
		return false
	} else {
		return localClientID.Get(clientID).(*ClientQueue).Insert(connectionID)
	}
}

func (o *ObfsAuthChainData) Remove(userID, clientID int) {
	localClientID := o.UserID[userID]
	if localClientID != nil {
		if localClientID.IsExist(clientID) {
			localClientID.Get(clientID).(*ClientQueue).DelRef()
		}
	}
}

/*----------------------------------AuthChainA----------------------------------*/
type AuthChainA struct {
	AuthBase
	RecvBuf       []byte
	UnintLen      int
	HasSentHeader bool

	HasRecvHeader     bool
	ClientID          int
	ConnectionID      int
	MaxTimeDif        int
	Salt              []byte
	PackID            int
	RecvID            int
	UserID            int
	UserIDNum         int
	UserKey           []byte
	ClientOverhead    int
	LastClientHash    []byte
	LastServerHash    []byte
	RandomClient      *XorShift128Plus
	RandomServer      *XorShift128Plus
	ObfsAuthChainData *ObfsAuthChainData
	Encryptor         *ciphers.Encryptor
}

func NewAuthChainA(method string) *AuthChainA {
	return &AuthChainA{
		AuthBase: AuthBase{
			Method:             method,
			RawTrans:           false,
			Overhead:           4,
			NoCompatibleMethod: "auth_chain_a",
		},
		RecvBuf:           []byte{},
		UnintLen:          2800,
		HasRecvHeader:     false,
		HasSentHeader:     false,
		ClientID:          0,
		ConnectionID:      0,
		MaxTimeDif:        60 * 60 * 24,
		Salt:              []byte("auth_chain_a"),
		PackID:            1,
		RecvID:            1,
		UserIDNum:         0,
		ClientOverhead:    4,
		LastClientHash:    []byte{},
		LastServerHash:    []byte{},
		RandomClient:      NewXorShift128Plus(),
		RandomServer:      NewXorShift128Plus(),
		ObfsAuthChainData: NewObfsAuthChainData(method),
	}
}

func (a *AuthChainA) InitData() []byte {
	panic("not implemented")
}

func (a *AuthChainA) GetOverhead(direction bool) int {
	return a.Overhead
}

func (a *AuthChainA) SetServerInfo(s ServerInfo) {
	a.SetServerInfo(s)
	var maxClient int
	maxClient, err := strconv.Atoi(strings.Split(s.GetProtocolParam(), "#")[0])
	if err != nil {
		maxClient = 64
	}
	a.ObfsAuthChainData.SetMaxClient(maxClient)
}

func (a *AuthChainA) ClientPreEncrypt(buf []byte) ([]byte, error) {
	panic("not implemented")
}

func (a *AuthChainA) ClientEncode(buf []byte) ([]byte, error) {
	panic("not implemented")
}

func (a *AuthChainA) ClientDecode(buf []byte) ([]byte, bool, error) {
	panic("not implemented")
}

func (a *AuthChainA) ClientPostDecrypt(buf []byte) ([]byte, error) {
	panic("not implemented")
}

func (a *AuthChainA) ServerPreEncrypt(buf []byte) ([]byte, error) {
	panic("not implemented")
}

func (a *AuthChainA) ServerEncode(buf []byte) ([]byte, error) {
	panic("not implemented")
}

func (a *AuthChainA) ServerDecode(buf []byte) ([]byte, bool, bool, error) {
	panic("not implemented")
}

func (a *AuthChainA) ServerPostDecrypt(buf []byte) ([]byte, bool, error) {
	panic("not implemented")
}

func (a *AuthChainA) ClientUDPPreEncrypt(buf []byte) ([]byte, error) {
	panic("not implemented")
}

func (a *AuthChainA) ClientUDPPostDecrypt(buf []byte) ([]byte, error) {
	panic("not implemented")
}

func (a *AuthChainA) ServerUDPPreEncrypt(buf []byte) ([]byte, error) {
	panic("not implemented")
}

func (a *AuthChainA) ServerUDPPostDecrypt(buf []byte) ([]byte, string, error) {
	panic("not implemented")
}

func (a *AuthChainA) Dispose() {
	panic("not implemented")
}

func (a *AuthChainA) GetHeadSize(buf []byte, defaultValue int) int {
	panic("not implemented")
}

func (a *AuthChainA) trapezoidRandomFloat(d float64) float64 {
	if d == 0 {
		return randomx.Float64()
	}
	s := randomx.Float64()
	tmp := 1 - d
	return (math.Sqrt(tmp*tmp+4*d*s) - tmp) / (2 * d)
}

func (a *AuthChainA) trapezoidRandomInt(maxVal, d float64) int {
	v := a.trapezoidRandomFloat(d)
	return int(v * maxVal)
}

func (a *AuthChainA) rndDataLen(bufSize int, lastHash []byte, random *XorShift128Plus) int {
	if bufSize > 1440 {
		return 0
	}
	random.InitFromBinLen(lastHash, bufSize)
	if bufSize > 1300 {
		return int(random.Next()) % 31
	}
	if bufSize > 900 {
		return int(random.Next()) % 127
	}
	if bufSize > 400 {
		return int(random.Next()) % 521
	}
	return int(random.Next()) % 1024
}

func (a *AuthChainA) udpRndDataLen(lastHash []byte, random XorShift128Plus) int {
	random.InitFromBin(lastHash)
	return int(random.Next()) % 127
}

func (a *AuthChainA) rndStartPos(randLen int, random *XorShift128Plus) int {
	if randLen > 0 {
		return int(randomx.Int64() % 8589934609 % int64(randLen))
	}
	return 0
}

func (a *AuthChainA) rndData(bufSize int, buf []byte, lastHashe []byte, random *XorShift128Plus) []byte {
	randLen := a.rndDataLen(bufSize, lastHashe, random)
	rndDataBuf := randomx.RandomBytes(randLen)
	if bufSize == 0 {
		return rndDataBuf
	} else {
		if randLen > 0 {
			startPos := a.rndStartPos(randLen, random)
			return conbineToBytes(rndDataBuf[:startPos], buf, rndDataBuf[startPos:])
		} else {
			return buf
		}
	}
}

func (a *AuthChainA) packClientData(buf []byte) ([]byte, error) {
	buf, err := a.Encryptor.Decrypt(buf)
	if err != nil {
		return nil, err
	}
	data := a.rndData(len(buf), buf, a.LastClientHash, a.RandomClient)
	macKey := bytesx.ContactSlice(a.UserKey, binaryx.BEUInt32ToBytes(uint32((a.PackID))))
	length := len(buf) ^ int(binaryx.LEBytesToUint16(a.LastServerHash[14:]))
	data = bytesx.ContactSlice(binaryx.LEUInt16ToBytes(uint16(length)), data)
	a.LastClientHash = hmacmd5(macKey, data)
	data = bytesx.ContactSlice(data, a.LastClientHash[:2])
	a.PackID = a.PackID + 1&0xFFFFFFFF
	return data, nil
}

func (a *AuthChainA) packAuthData(authData, buf []byte) (result []byte, err error) {
	data := authData
	data = bytesx.ContactSlice(data, binaryx.LEUInt16ToBytes(uint16(a.GetServerInfo().GetOverhead())), binaryx.LEUInt16ToBytes(0))
	mac_key := bytesx.ContactSlice(a.GetServerInfo().GetIv(), a.GetServerInfo().GetKey())

	check_head := randomx.RandomBytes(4)
	a.LastClientHash = hmacmd5(mac_key, check_head)
	check_head = bytesx.ContactSlice(check_head, a.LastClientHash[:8])

	param := a.GetServerInfo().GetProtocolParam()
	if strings.Contains(param, ":") {
		items := strings.Split(param, ":")
		var uidBytes []byte
		if len(items) > 1 {
			a.UserKey = []byte(items[1])
			uidInt, err := strconv.Atoi(items[0])
			if err != nil {
				return nil, err
			}
			uidBytes = binaryx.LEUInt16ToBytes(uint16((uidInt)))
		} else {
			uidBytes = randomx.RandomBytes(4)
		}
		if a.UserKey == nil {
			a.UserKey = a.GetServerInfo().GetKey()
		}
		encryptor, err := ciphers.NewEncryptorWithIv("aes-128-cbc",
			string(bytesx.ContactSlice([]byte(base64.StdEncoding.EncodeToString(a.UserKey)), a.Salt)),
			bytes.Repeat([]byte{0x00}, 16))
		if err != nil {
			return nil, err
		}
		uid := binaryx.LEBytesToUint16(uidBytes) ^ binaryx.LEBytesToUint16(a.LastClientHash[8:12])
		uidBytes := binaryx.LEUInt16ToBytes(uid)
		dataCipherText,err := encryptor.Encrypt(data)
		if err != nil{
			return nil,err
		}
		data = bytesx.ContactSlice(uidBytes,dataCipherText)
		a.LastServerHash = hmacmd5(a.UserKey,data)
		data = bytesx.ContactSlice(check_head,data,a.LastServerHash[:4])
		a.Encryptor = encryptor.

	}
}
