package obfs

import (
	"bytes"
	"math"
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

func (o *ObfsAuthChainData) SetMaxClient(maxClient int){
	o.MaxClient = maxClient
	o.MaxBuffer = int(math.Max(float64(maxClient),1024))
}
func (o *ObfsAuthChainData) Insert(userID ,clientID,connectionID int) bool{
	if o.UserID[userID] == nil{
		o.UserID[userID] = cache.NewLruCache(60 * time.Second)
	}
	localClientID := o.UserID[userID]
	var r ,_ = localClientID.Get(clientID).(*ClientQueue)
	if r != nil || !r.Enable{
		if localClientID.First() == nil || localClientID.Len() < o.MaxClient{
			if !localClientID.IsExist(clientID){
				// TODO check
				localClientID.Put(clientID,NewClientQueue(connectionID),60*time.Second)
			}else{
				localClientID.Get(clientID).(ClientQueue).ReEnable(connectionID)
			}
			return localClientID.Get(clientID).(ClientQueue).Insert(connectionID)
		}

		localClientIDFirst := localClientID.First()
		if !localClientID.Get(localClientIDFirst).(ClientQueue).IsActive(){
			localClientID.Delete(localClientIDFirst)
			if !localClientID.IsExist(clientID){
				// TODO check
				localClientID.Put(clientID,NewClientQueue(connectionID),60*time.Second)
			}else{
				localClientID.Get(clientID).(ClientQueue).ReEnable(connectionID)
			}
			return localClientID.Get(clientID).(ClientQueue).Insert(connectionID)
		}

		log.Warn("%s: no inactive client",o.Name)
		return false
	}else{
		return localClientID.Get(clientID).(ClientQueue).Insert(connectionID)
	}
}

func  (o *ObfsAuthChainData) Remove(userID,clientID int){
	localClientID := o.UserID[userID]
	if localClientID != nil{
		if localClientID.IsExist(clientID){
			localClientID.Get(clientID).(ClientQueue).DelRef()
		}
	}
}

/*----------------------------------AuthChainA----------------------------------*/

type AuthChainA struct {
	AuthBase
	RecvBuf        []byte
	UnintLen       int
	HasSentHeader  bool

	HasRecvHeader  bool
	ClientID       int
	ConnectionID   int
	MaxTimeDif     int
	Salt           []byte
	PackID         int
	RecvID         int
	UserID         int
	UserIDNum      int
	UserKey        []byte
	ClientOverhead int
	LastClientHash []byte
	LastServerHash []byte
	RandomClient   *XorShift128Plus
	RandomServer   *XorShift128Plus
}

func NewAuthChainA() *AuthChainA {
	return &AuthChainA{
		AuthBase: AuthBase{
			RawTrans:           false,
			Overhead:           4,
			NoCompatibleMethod: "auth_chain_a",
		},
		RecvBuf:        []byte{},
		UnintLen:       2800,
		HasRecvHeader:  false,
		HasSentHeader:  false,
		ClientID:       0,
		ConnectionID:   0,
		MaxTimeDif:     60 * 60 * 24,
		Salt:           []byte("auth_chain_a"),
		PackID:         1,
		RecvID:         1,
		UserIDNum:      0,
		ClientOverhead: 4,
		LastClientHash: []byte{},
		LastServerHash: []byte{},
		RandomClient:   NewXorShift128Plus(),
		RandomServer:   NewXorShift128Plus(),
	}
}

func (a *AuthChainA) InitData() []byte {
	panic("not implemented")
}

func (a *AuthChainA) GetMethod() string {
	panic("not implemented")
}

func (a *AuthChainA) SetMethod(method string) {
	panic("not implemented")
}

func (a *AuthChainA) GetOverhead(direction bool) int {
	panic("not implemented")
}

func (a *AuthChainA) GetServerInfo() serverInfo {
	panic("not implemented")
}

func (a *AuthChainA) SetServerInfo(s ServerInfo) {
	panic("not implemented")
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
