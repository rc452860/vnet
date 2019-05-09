package network

import (
	"encoding/hex"
	"fmt"
	"github.com/pkg/errors"
	"github.com/rc452860/vnet/common/ciphers"
	"github.com/rc452860/vnet/common/obfs"
	"github.com/rc452860/vnet/utils/addr"
	"github.com/sirupsen/logrus"
	"net"
)

func NewShadowsocksRDecorate(request *Request, obfsMethod, cryptMethod, key, protocolMethod, obfsParam, protocolParam, host string, port int, isLocal bool) (ssrd *ShadowsocksRDecorate, err error) {
	// init essential parameters
	ssrd = &ShadowsocksRDecorate{
		Request:       request,
		ObfsParam:     obfsParam,
		ProtocolParam: protocolParam,
		Host:          host,
		Port:          port,
		ISLocal:       isLocal,
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
}

func (ssrd *ShadowsocksRDecorate) Read(buf []byte) (n int, err error) {

	// ServerDecode return buffer_to_recv, is_need_decrypt, is_need_to_encode_and_send_back
	n, err = ssrd.TCPConn.Read(buf)
	if err != nil {
		return 0, err
	}
	data := buf[:n]
	unobfsData, needDecrypt, needSendBack, err := ssrd.obfs.ServerDecode(data)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"requestId":    ssrd.RequestID,
			"data":         hex.EncodeToString(data),
			"needDecrypt":  needDecrypt,
			"needSendBack": needSendBack,
			"err":          err,
		}).Debugf("ShadowsocksRDecorate obfs decrypt error.")
		return 0, errors.New(fmt.Sprintf("[%s] shadowsocksr obfs decrypt error.", ssrd.RequestID))
	}

	if needSendBack {
		result, err := ssrd.obfs.ServerEncode([]byte{})
		if err != nil {
			return 0, err
		}
		if len(result) > len(buf) {
			logrus.WithFields(logrus.Fields{
				"requestId":     ssrd.RequestID,
				"resultLen":     len(result),
				"unobfsDataLen": len(unobfsData),
				"unobfsData":    hex.EncodeToString(unobfsData),
				"result":        hex.EncodeToString(result),
			}).Debugf("")
			return 0, errors.New(fmt.Sprintf("[%s] ShadowsocksRDecorate obfs buf is too short.", ssrd.RequestID))
		}
		_, err = ssrd.TCPConn.Write(result)
		if err != nil {
			return 0, errors.Wrap(err, fmt.Sprintf("[%s] ShadowsocksRDecorate obfs sendback error.", ssrd.RequestID))
		}
	}

	if needDecrypt {
		if ssrd.protocol.GetServerInfo().GetRecvIv() == nil || len(ssrd.protocol.GetServerInfo().GetRecvIv()) == 0 {
			ivLen := len(ssrd.protocol.GetServerInfo().GetRecvIv())
			ssrd.protocol.GetServerInfo().SetRecvIv(unobfsData[:ivLen])
		}
		cleartext, err := ssrd.encryptor.Decrypt(unobfsData)
		if err != nil {
			return 0, errors.Wrap(err, fmt.Sprintf("[%s] ShadowsocksRDecorate obfs decrypt error.", ssrd.RequestID))
		}
		copySize := copy(data, cleartext)
		data = data[:copySize]
	} else {
		data = unobfsData
	}

	data, sendback, err := ssrd.protocol.ServerPostDecrypt(data)
	if err != nil {
		return 0, errors.Wrap(err, fmt.Sprintf("[%s] ShadowsocksRDecorate protocol post decrypt error.", ssrd.RequestID))
	}
	if sendback {
		backdata, err := ssrd.protocol.ServerPreEncrypt([]byte{})
		if err != nil {
			return 0, errors.Wrap(err, fmt.Sprintf("[%s] ShadowsocksRDecorate protocol pre encode error.", ssrd.RequestID))
		}
		backdata, err = ssrd.encryptor.Encrypt(backdata)
		if err != nil {
			return 0, errors.Wrap(err, fmt.Sprintf("[%s] ShadowsocksRDecorate encrypter encrypt error.", ssrd.RequestID))
		}
		backdata, err = ssrd.obfs.ServerEncode(backdata)
		if err != nil {
			return 0, errors.Wrap(err, fmt.Sprintf("[%s] ShadowsocksRDecorate obfs server encode error.", ssrd.RequestID))
		}
		_, err = ssrd.TCPConn.Write(backdata)
		if err != nil {
			return 0, errors.Wrap(err, fmt.Sprintf("[%s] ShadowsocksRDecorate obfs sendback error.", ssrd.RequestID))
		}
	}
	copy(buf, data)
	return len(data), nil
}

func (ssrd *ShadowsocksRDecorate) Write(buf []byte) (n int, err error) {
	data, err := ssrd.protocol.ServerPreEncrypt(buf)
	if err != nil {
		return 0, errors.Wrap(err, fmt.Sprintf("[%s] ShadowsocksRDecorate protocol server encode error.", ssrd.RequestID))
	}
	data, err = ssrd.encryptor.Encrypt(data)
	if err != nil {
		return 0, errors.Wrap(err, fmt.Sprintf("[%s] ShadowsocksRDecorate encryptor encrypt error.", ssrd.RequestID))
	}
	data, err = ssrd.obfs.ServerEncode(data)
	if err != nil {
		return 0, errors.Wrap(err, fmt.Sprintf("[%s] ShadowsocksRDecorate obfs server encode error.", ssrd.RequestID))
	}
	return ssrd.TCPConn.Write(data)
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
	return serverInfo
}

func (ssrd *ShadowsocksRDecorate) UpdateUser(uid []byte) {
	// TODO: update user callback
}
