package httputil_test

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/json"
	"encoding/pem"
	"math/big"
	"net"
	"net/http"
	"net/netip"
	"net/url"
	"testing"
	"time"

	"github.com/AdguardTeam/golibs/logutil/slogutil"
	"github.com/AdguardTeam/golibs/netutil"
	"github.com/AdguardTeam/golibs/netutil/httputil"
	"github.com/AdguardTeam/golibs/netutil/urlutil"
	"github.com/AdguardTeam/golibs/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	// testServerName is the common server name for test TLS configuration.
	testServerName = "example.org"
)

// newTLSConfig is a helper that creates a TLS config with default values and
// returns a PEM-encoded certificate and the config itself.
func newTLSConfig(tb testing.TB) (conf *tls.Config, certPem []byte) {
	tb.Helper()

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(tb, err)

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	require.NoError(tb, err)

	notBefore := time.Now()
	notAfter := notBefore.Add(5 * 365 * time.Hour * 24)

	keyUsage := x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign
	template := x509.Certificate{
		SerialNumber:          serialNumber,
		Subject:               pkix.Name{Organization: []string{"AdGuard Tests"}},
		NotBefore:             notBefore,
		NotAfter:              notAfter,
		KeyUsage:              keyUsage,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		IsCA:                  true,
		DNSNames:              []string{testServerName},
	}

	derBytes, err := x509.CreateCertificate(
		rand.Reader,
		&template,
		&template,
		&privateKey.PublicKey,
		privateKey,
	)
	require.NoError(tb, err)

	certPem = pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: derBytes,
	})
	keyPem := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})

	cert, err := tls.X509KeyPair(certPem, keyPem)
	require.NoError(tb, err)

	return &tls.Config{Certificates: []tls.Certificate{cert}, ServerName: testServerName}, certPem
}

// runTestServer helper starts the specified server and waits for it to be ready
// to serve requests, then returns its listening address.
func runTestServer(tb testing.TB, srv *httputil.Server) (addr net.Addr) {
	ctx := testutil.ContextWithTimeout(tb, testTimeout)
	go func() {
		err := srv.Start(ctx)
		if err != nil {
			require.ErrorIs(tb, http.ErrServerClosed, err)
		}
	}()

	testutil.CleanupAndRequireSuccess(tb, func() (err error) {
		return srv.Shutdown(ctx)
	})

	for addr == nil {
		addr = srv.LocalAddr()
	}

	return addr
}

// logRecord represents server log record.  It is used for unmarshalling
// [httputil.Server] logs in tests.
type logRecord struct {
	Server *url.URL `json:"server"`
}

// serverURLFromLog is a helper that parses server log record and returns
// parsed url.
func serverURLFromLog(tb testing.TB, log []byte) (url *url.URL) {
	tb.Helper()

	var record logRecord
	err := json.Unmarshal(log, &record)
	require.NoError(tb, err)

	return record.Server
}

func TestServer(t *testing.T) {
	logOutput := &bytes.Buffer{}
	logger := slogutil.New(&slogutil.Config{
		Output: logOutput,
		Format: slogutil.FormatJSON,
	})

	tlsConfig, caPem := newTLSConfig(t)
	roots := x509.NewCertPool()
	require.True(t, roots.AppendCertsFromPEM(caPem))

	tlsTransport := &http.Transport{
		TLSClientConfig: &tls.Config{
			RootCAs:    roots,
			ServerName: testServerName,
		},
	}

	require.True(t, t.Run("non_zero_port", func(t *testing.T) {
		srv := httputil.NewServer(&httputil.ServerConfig{
			BaseLogger:     logger,
			InitialAddress: netip.AddrPortFrom(netutil.IPv4Localhost(), 1234),
			Server: &http.Server{
				Handler: httputil.HealthCheckHandler,
			},
		})

		addr := runTestServer(t, srv)
		testURL := &url.URL{
			Scheme: urlutil.SchemeHTTP,
			Host:   addr.String(),
		}

		ctx := testutil.ContextWithTimeout(t, testTimeout)
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, testURL.String(), nil)
		require.NoError(t, err)

		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)

		lines := bytes.Split(logOutput.Bytes(), []byte("\n"))
		require.Len(t, lines, 2)

		url := serverURLFromLog(t, lines[0])
		assert.Equal(t, urlutil.SchemeHTTP, url.Scheme)
		assert.Equal(t, "1234", url.Port())
	}))

	require.True(t, t.Run("tls_zero_port", func(t *testing.T) {
		logOutput.Reset()
		client := http.Client{
			Transport: tlsTransport,
		}

		srv := httputil.NewServer(&httputil.ServerConfig{
			BaseLogger:     logger,
			InitialAddress: netip.AddrPortFrom(netutil.IPv4Localhost(), 0),
			Server: &http.Server{
				Handler:   httputil.HealthCheckHandler,
				TLSConfig: tlsConfig,
			},
		})

		addr := runTestServer(t, srv)
		testURL := &url.URL{
			Scheme: urlutil.SchemeHTTPS,
			Host:   addr.String(),
		}

		ctx := testutil.ContextWithTimeout(t, testTimeout)
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, testURL.String(), nil)
		require.NoError(t, err)

		resp, err := client.Do(req)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)

		lines := bytes.Split(logOutput.Bytes(), []byte("\n"))
		require.Len(t, lines, 2)

		url := serverURLFromLog(t, lines[0])
		assert.Equal(t, urlutil.SchemeHTTPS, url.Scheme)
		assert.NotEqual(t, "0", url.Port())
	}))
}
