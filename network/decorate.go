package network

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/pkg/errors"
	"github.com/rc452860/vnet/common/ciphers"
	"github.com/rc452860/vnet/common/obfs"
	"github.com/rc452860/vnet/utils/addr"
	"github.com/rc452860/vnet/utils/binaryx"
	"github.com/sirupsen/logrus"
	"net"
)

func NewShadowsocksRDecorate(request *Request, obfsMethod, cryptMethod, key, protocolMethod, obfsParam, protocolParam, host string, port int, isLocal bool, users map[string]string) (ssrd *ShadowsocksRDecorate, err error) {
	// init essential parameters
	ssrd = &ShadowsocksRDecorate{
		Request:       request,
		ObfsParam:     obfsParam,
		ProtocolParam: protocolParam,
		Host:          host,
		Port:          port,
		ISLocal:       isLocal,
		Users:         users,
		recvBuf:       new(bytes.Buffer),
	}

	// init obfs protocol encrypto component
	ssrd.obfs = obfs.GetObfs(obfsMethod)
	ssrd.protocol = obfs.GetObfs(protocolMethod)
	ssrd.encryptor, err = ciphers.NewEncryptor(cryptMethod, key)

	ssrd.Overhead = ssrd.obfs.GetOverhead(isLocal) + ssrd.protocol.GetOverhead(isLocal)
	// set serverinfo
	ssrd.obfs.SetServerInfo(ssrd.getServerInfo(true))
	ssrd.protocol.SetServerInfo(ssrd.getServerInfo(false))
	return ssrd, err
}

type ShadowsocksRDecorate struct {
	*Request
	obfs          obfs.Plain
	protocol      obfs.Plain
	encryptor     *ciphers.Encryptor
	Host          string
	Port          int
	ObfsParam     string
	ProtocolParam string
	Users         map[string]string
	Overhead      int
	ISLocal       bool
	recvBuf       *bytes.Buffer
}

//func (ssrd *ShadowsocksRDecorate) Read(buf []byte)(n int,err error){
//	// ServerDecode return buffer_to_recv, is_need_decrypt, is_need_to_encode_and_send_back
//	if ssrd.recvBuf.Len() > 0 {
//		n,err =  ssrd.recvBuf.Read(buf)
//		logrus.WithFields(logrus.Fields{
//			"readData":string(buf[:n]),
//		}).Debug()
//		return n,err
//	}
//	bufTmp := make([]byte, 2048)
//	n, err = ssrd.TCPConn.Read(bufTmp)
//	if err != nil {
//		return 0, err
//	}
//	data := bufTmp[:n]
//	data,err = ssrd.encryptor.Decrypt(data)
//	if err !=nil{
//		return 0,nil
//	}
//	ssrd.recvBuf.Write(data)
//	n,err = ssrd.recvBuf.Read(buf)
//	logrus.WithFields(logrus.Fields{
//		"readData":string(buf[:n]),
//	}).Debug()
//	return n,err
//}
func (ssrd *ShadowsocksRDecorate) Read(buf []byte) (n int, err error) {

	// ServerDecode return buffer_to_recv, is_need_decrypt, is_need_to_encode_and_send_back
	if ssrd.recvBuf.Len() > 0 {
		return ssrd.recvBuf.Read(buf)
	}

	bufTmp := make([]byte, 2048)
	n, err = ssrd.TCPConn.Read(bufTmp)
	if err != nil {
		return 0, err
	}
	data := bufTmp[:n]
	unobfsData, needDecrypt, needSendBack, err := ssrd.obfs.ServerDecode(data)
	logrus.WithFields(logrus.Fields{
		"requestId":    ssrd.RequestID,
		"data":         hex.EncodeToString(data),
		"unobfsData":   hex.EncodeToString(unobfsData),
		"needDecrypt":  needDecrypt,
		"needSendBack": needSendBack,
	}).Debug("shadowsocksr obfs ServerDecode")
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"err": err,
		}).Debugf("ShadowsocksRDecorate obfs decrypt error.")
		return 0, errors.New(fmt.Sprintf("[%s] shadowsocksr obfs decrypt error.", ssrd.RequestID))
	}

	if needSendBack {
		result, err := ssrd.obfs.ServerEncode([]byte{})
		if err != nil {
			return 0, err
		}
		_, err = ssrd.TCPConn.Write(result)
		if err != nil {
			return 0, errors.Wrap(err, fmt.Sprintf("[%s] ShadowsocksRDecorate obfs sendback error.", ssrd.RequestID))
		}
		return ssrd.Read(buf)
	}

	if needDecrypt {
		//if ssrd.protocol.GetServerInfo().GetRecvIv() == nil || len(ssrd.protocol.GetServerInfo().GetRecvIv()) == 0 {
		//	ivLen := len(ssrd.protocol.GetServerInfo().GetIv())
		//	ssrd.protocol.GetServerInfo().SetRecvIv(unobfsData[:ivLen])
		//}
		cleartext, err := ssrd.encryptor.Decrypt(unobfsData)
		if ssrd.protocol.GetServerInfo().GetRecvIv() == nil || len(ssrd.protocol.GetServerInfo().GetRecvIv()) == 0 {
			ssrd.protocol.GetServerInfo().SetRecvIv(ssrd.encryptor.IVIn)
		}
		logrus.WithFields(logrus.Fields{
			"cleartextHexEncode": hex.EncodeToString(cleartext),
		}).Debug("ShadowsocksRDecorate encryptor Decrypt")
		if err != nil {
			return 0, errors.Wrap(err, fmt.Sprintf("[%s] ShadowsocksRDecorate obfs decrypt error.", ssrd.RequestID))
		}
		data = cleartext
	} else {
		data = unobfsData
	}

	data, sendback, err := ssrd.protocol.ServerPostDecrypt(data)
	logrus.WithFields(logrus.Fields{
		"serverPostDecryptHex": hex.EncodeToString(data),
		"sendback":             sendback,
	}).Debug("ShadowsocksRDecorate protocol ServerPostDecrypt")
	if err != nil {
		return 0, errors.Wrap(err, fmt.Sprintf("[%s] ShadowsocksRDecorate protocol post decrypt error.", ssrd.RequestID))
	}
	if sendback {
		backdata, err := ssrd.protocol.ServerPreEncrypt([]byte{})
		logrus.WithFields(logrus.Fields{
			"backdata": hex.EncodeToString(backdata),
			//"LastServerHash":hex.EncodeToString(ssrd.protocol.(*obfs.AuthChainA).LastServerHash),
		}).Debug("shadowoscksr Read ServerPreEncrypt")
		if err != nil {
			return 0, errors.Wrap(err, fmt.Sprintf("[%s] ShadowsocksRDecorate protocol pre encode error.", ssrd.RequestID))
		}
		backdata, err = ssrd.encryptor.Encrypt(backdata)
		if err != nil {
			return 0, errors.Wrap(err, fmt.Sprintf("[%s] ShadowsocksRDecorate encrypter encrypt error.", ssrd.RequestID))
		}
		logrus.WithFields(logrus.Fields{
			"ReadEncryptData": hex.EncodeToString(backdata),
		}).Debug("shadowoscksr Read Encrypt")
		backdata, err = ssrd.obfs.ServerEncode(backdata)
		if err != nil {
			return 0, errors.Wrap(err, fmt.Sprintf("[%s] ShadowsocksRDecorate obfs server encode error.", ssrd.RequestID))
		}
		logrus.WithFields(logrus.Fields{
			"ReadServerEncodeData": hex.EncodeToString(backdata),
		}).Debug("shadowoscksr Read ServerEncode")
		_, err = ssrd.TCPConn.Write(backdata)
		if err != nil {
			return 0, errors.Wrap(err, fmt.Sprintf("[%s] ShadowsocksRDecorate obfs sendback error.", ssrd.RequestID))
		}
	}
	ssrd.recvBuf.Write(data)
	n, err = ssrd.recvBuf.Read(buf)
	return n, err
}

//func (ssrd *ShadowsocksRDecorate) Write(buf []byte) (n int,err error){
//	logrus.WithFields(logrus.Fields{
//		"WriteData":string(buf),
//	}).Debug()
//	result,err := ssrd.encryptor.Encrypt(buf)
//	if err !=nil{
//		return 0,err
//	}
//	_,err = ssrd.TCPConn.Write(result)
//	if err !=nil{
//		return 0,err
//	}
//	return len(buf),nil
//}
func (ssrd *ShadowsocksRDecorate) Write(buf []byte) (n int, err error) {
	data, err := ssrd.protocol.ServerPreEncrypt(buf)

	if err != nil {
		return 0, errors.Wrap(err, fmt.Sprintf("[%s] ShadowsocksRDecorate protocol server encode error.", ssrd.RequestID))
	}
	logrus.WithFields(logrus.Fields{
		"ServerPreEncryptWriteData": hex.EncodeToString(data),
		//"LastServerHash":hex.EncodeToString(ssrd.protocol.(*obfs.AuthChainA).LastServerHash),
	}).Debug("shadowoscksr Write ServerPreEncrypt")
	data, err = ssrd.encryptor.Encrypt(data)
	if err != nil {
		return 0, errors.Wrap(err, fmt.Sprintf("[%s] ShadowsocksRDecorate encryptor encrypt error.", ssrd.RequestID))
	}
	logrus.WithFields(logrus.Fields{
		"EncryptWriteData": hex.EncodeToString(data),
	}).Debug("shadowoscksr Write Encrypt")
	data, err = ssrd.obfs.ServerEncode(data)
	if err != nil {
		return 0, errors.Wrap(err, fmt.Sprintf("[%s] ShadowsocksRDecorate obfs server encode error.", ssrd.RequestID))
	}
	logrus.WithFields(logrus.Fields{
		"ServerEncodeWriteData": hex.EncodeToString(data),
	}).Debug("shadowoscksr Write ServerEncode")
	n, err = ssrd.TCPConn.Write(data)
	if err != nil {
		return 0, err
	}
	return len(buf), nil
}

func (ssrd *ShadowsocksRDecorate) getServerInfo(isObfs bool) obfs.ServerInfo {
	serverInfo := obfs.NewServerInfo()
	serverInfo.SetHost(ssrd.Host)
	serverInfo.SetPort(ssrd.Port)
	serverInfo.SetClient(net.ParseIP(addr.GetIPFromAddr(ssrd.TCPConn.RemoteAddr())))
	serverInfo.SetPort(addr.GetPortFromAddr(ssrd.TCPConn.RemoteAddr()))
	if isObfs {
		serverInfo.SetObfsParam(ssrd.ObfsParam)
		serverInfo.SetProtocolParam("")
	} else {
		serverInfo.SetObfsParam("")
		serverInfo.SetProtocolParam(ssrd.ProtocolParam)
	}
	serverInfo.SetIv(ssrd.encryptor.IVOut)
	serverInfo.SetRecvIv([]byte{})
	serverInfo.SetKeyStr(ssrd.encryptor.KeyStr)
	serverInfo.SetKey(ssrd.encryptor.Key)
	serverInfo.SetHeadLen(obfs.DEFAULT_HEAD_LEN)
	// TODO: need calculate,for now, I don't know how to implement it on windows
	serverInfo.SetTCPMss(obfs.TCP_MSS)
	serverInfo.SetBufferSize(obfs.BUF_SIZE - ssrd.Overhead)
	serverInfo.SetOverhead(ssrd.Overhead)
	serverInfo.SetUpdateUserFunc(ssrd.UpdateUser)
	serverInfo.SetUsers(ssrd.Users)
	return serverInfo
}

func (ssrd *ShadowsocksRDecorate) UpdateUser(uid []byte) {
	// TODO: update user callback
	uidInt := binaryx.LEBytesToUInt32(uid)
	logrus.Infof("ShadowsocksRDecorate update uid: %v", uidInt)
}

/*below deprecated*/
func (ssrd *ShadowsocksRDecorate) AddUser(uid uint32, passwd string) {
	uidPack := string(binaryx.LEUint32ToBytes(uid))
	ssrd.Users[uidPack] = passwd
	ssrd.protocol.GetServerInfo().SetUsers(ssrd.Users)
}

func (ssrd *ShadowsocksRDecorate) DelUser(uid uint32) {
	uidPack := string(binaryx.LEUint32ToBytes(uid))
	delete(ssrd.Users, uidPack)
	ssrd.protocol.GetServerInfo().SetUsers(ssrd.Users)
}
func (ssrd *ShadowsocksRDecorate) Reload(users map[string]string) {
	ssrd.Users = users
	ssrd.protocol.GetServerInfo().SetUsers(ssrd.Users)
}
