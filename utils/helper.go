package utils

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"math/big"
	"net"
	"strings"
)

// NewTLSConfigFromFile 通过证书文件生成证书信息
func NewTLSConfigFromFile(tlsCertFile, tlsKeyFile string, insecureSkipVerifyForClient ...bool) (*tls.Config, error) {
	cert, err := tls.LoadX509KeyPair(tlsCertFile, tlsKeyFile)
	if err != nil {
		return nil, err
	}
	return newTLSConfig(cert, insecureSkipVerifyForClient...), nil
}

func newTLSConfig(cert tls.Certificate, insecureSkipVerifyForClient ...bool) *tls.Config {
	var insecureSkipVerify bool
	if len(insecureSkipVerifyForClient) > 0 {
		insecureSkipVerify = insecureSkipVerifyForClient[0]
	}
	return &tls.Config{
		InsecureSkipVerify:       insecureSkipVerify,
		Certificates:             []tls.Certificate{cert},
		NextProtos:               []string{"http/1.1", "h2"},
		PreferServerCipherSuites: true,
		CurvePreferences: []tls.CurveID{
			tls.CurveP256,
			tls.X25519,
		},
		MinVersion: tls.VersionTLS12,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_RSA_WITH_AES_128_GCM_SHA256,
		},
	}
}

// FakeAddr 是一个虚地址对象，实现了net.Add
type FakeAddr struct {
	network string
	addr    string
	host    string
	port    string
	udpAddr *net.UDPAddr
}

var _ net.Addr = (*FakeAddr)(nil)

// NewFakeAddr 创建一个虚地址对象
func NewFakeAddr(network, host, port string) *FakeAddr {
	if network == "" {
		network = "tcp"
	}
	if host == "" {
		host = "0.0.0.0"
	}
	if port == "" {
		port = "0"
	}
	addr := net.JoinHostPort(host, port)
	return &FakeAddr{
		network: network,
		addr:    addr,
		host:    host,
		port:    port,
	}
}

// NewFakeAddr2 创建另一个不同参数的虚地址对象
func NewFakeAddr2(network, addr string) (*FakeAddr, error) {
	if addr == "" {
		return NewFakeAddr(network, "", ""), nil
	}
	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		return nil, err
	}
	return NewFakeAddr(network, host, port), nil
}

func (f *FakeAddr) Network() string {
	return f.network
}

func (f *FakeAddr) String() string {
	return f.addr
}

func (f *FakeAddr) Host() string {
	return f.host
}

func (f *FakeAddr) Port() string {
	return f.port
}

func (f *FakeAddr) UdpAddr() *net.UDPAddr {
	return f.udpAddr
}

func (f *FakeAddr) SetUdpAddr(udpAddr *net.UDPAddr) {
	f.udpAddr = udpAddr
}

func GenerateTLSConfigForClient() *tls.Config {
	return &tls.Config{InsecureSkipVerify: true}
}

// GenerateTLSConfigForServer 生成tls证书
func GenerateTLSConfigForServer() *tls.Config {
	key, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		panic(err)
	}
	template := x509.Certificate{SerialNumber: big.NewInt(1)}
	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &key.PublicKey, key)
	if err != nil {
		panic(err)
	}
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)})
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})

	cert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		panic(err)
	}
	return newTLSConfig(cert)
}

// IsSameAddr 判断两个地址是否相同
func IsSameAddr(addOne, addTwo net.Addr) bool {
	if addOne.Network() != addTwo.Network() {
		return false
	}
	addOneStr := addOne.String()
	addTwoStr := addTwo.String()

	if addOneStr == addTwoStr {
		return true
	}
	//去掉地址上的ipv6前缀
	const ipv6prefix = "[::]"
	addOneStr = strings.TrimPrefix(addOneStr, ipv6prefix)
	addTwoStr = strings.TrimPrefix(addTwoStr, ipv6prefix)

	//去掉地址上的ipv4前缀
	const ipv4prefix = "0.0.0.0"
	addOneStr = strings.TrimPrefix(addOneStr, ipv4prefix)
	addTwoStr = strings.TrimPrefix(addTwoStr, ipv4prefix)

	//判断去掉前缀后的地址是否相等
	return addOneStr == addTwoStr
}
