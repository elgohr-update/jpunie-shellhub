package sshd

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"sync"
	"time"

	sshserver "github.com/gliderlabs/ssh"
	"github.com/shellhub-io/shellhub/agent/pkg/osauth"
	"github.com/shellhub-io/shellhub/pkg/api/client"
	"github.com/shellhub-io/shellhub/pkg/models"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
)

type sshConn struct {
	net.Conn
	closeCallback func(string)
	ctx           sshserver.Context
}

func (c *sshConn) Close() error {
	if id, ok := c.ctx.Value(sshserver.ContextKeySessionID).(string); ok {
		c.closeCallback(id)
	}

	return c.Conn.Close()
}

type Server struct {
	sshd               *sshserver.Server
	api                client.Client
	authData           *models.DeviceAuthResponse
	cmds               map[string]*exec.Cmd
	Sessions           map[string]net.Conn
	deviceName         string
	mu                 sync.Mutex
	keepAliveInterval  int
	singleUserPassword string
}

func NewServer(api client.Client, authData *models.DeviceAuthResponse, privateKey string, keepAliveInterval int, singleUserPassword string) *Server {
	s := &Server{
		api:               api,
		authData:          authData,
		cmds:              make(map[string]*exec.Cmd),
		Sessions:          make(map[string]net.Conn),
		keepAliveInterval: keepAliveInterval,
	}

	s.sshd = &sshserver.Server{
		PasswordHandler:  s.passwordHandler,
		PublicKeyHandler: s.publicKeyHandler,
		Handler:          s.sessionHandler,
		RequestHandlers:  sshserver.DefaultRequestHandlers,
		ChannelHandlers:  sshserver.DefaultChannelHandlers,
		ConnCallback: func(ctx sshserver.Context, conn net.Conn) net.Conn {
			closeCallback := func(id string) {
				s.mu.Lock()
				defer s.mu.Unlock()

				if v, ok := s.cmds[id]; ok {
					v.Process.Kill() // nolint:errcheck
					delete(s.cmds, id)
				}
			}

			return &sshConn{conn, closeCallback, ctx}
		},
	}

	err := s.sshd.SetOption(sshserver.HostKeyFile(privateKey))
	if err != nil {
		logrus.Warn(err)
	}

	return s
}

func (s *Server) ListenAndServe() error {
	return s.sshd.ListenAndServe()
}

func (s *Server) HandleConn(conn net.Conn) {
	s.sshd.HandleConn(conn)
}

func (s *Server) SetDeviceName(name string) {
	s.deviceName = name
}

func (s *Server) sessionHandler(session sshserver.Session) {
	sspty, winCh, isPty := session.Pty()

	log := logrus.WithFields(logrus.Fields{
		"user": session.User(),
		"pty":  isPty,
	})

	log.Info("New session request")

	go StartKeepAliveLoop(time.Second*time.Duration(s.keepAliveInterval), session)

	if isPty { //nolint:nestif
		scmd := newShellCmd(s, session.User(), sspty.Term)

		pts, err := startPty(scmd, session, winCh)
		if err != nil {
			logrus.Warn(err)
		}

		u := osauth.LookupUser(session.User())

		err = os.Chown(pts.Name(), int(u.UID), -1)
		if err != nil {
			logrus.Warn(err)
		}

		remoteAddr := session.RemoteAddr()

		logrus.WithFields(logrus.Fields{
			"user":       session.User(),
			"pty":        pts.Name(),
			"remoteaddr": remoteAddr,
			"localaddr":  session.LocalAddr(),
		}).Info("Session started")

		ut := utmpStartSession(
			pts.Name(),
			session.User(),
			remoteAddr.String(),
		)

		s.mu.Lock()
		s.cmds[session.Context().Value(sshserver.ContextKeySessionID).(string)] = scmd
		s.mu.Unlock()

		if err := scmd.Wait(); err != nil {
			logrus.Warn(err)
		}

		logrus.WithFields(logrus.Fields{
			"user":       session.User(),
			"pty":        pts.Name(),
			"remoteaddr": remoteAddr,
			"localaddr":  session.LocalAddr(),
		}).Info("Session ended")

		utmpEndSession(ut)
	} else {
		u := osauth.LookupUser(session.User())
		cmd := newCmd(u, "", "", s.deviceName, session.Command()...)

		stdout, _ := cmd.StdoutPipe()
		stdin, _ := cmd.StdinPipe()

		logrus.WithFields(logrus.Fields{
			"user":        session.User(),
			"remoteaddr":  session.RemoteAddr(),
			"localaddr":   session.LocalAddr(),
			"Raw command": session.RawCommand(),
		}).Info("Command started")

		err := cmd.Start()
		if err != nil {
			logrus.Warn(err)
		}

		go func() {
			if _, err := io.Copy(stdin, session); err != nil {
				fmt.Println(err) //nolint:forbidigo
			}
		}()

		go func() {
			if _, err := io.Copy(session, stdout); err != nil {
				fmt.Println(err) //nolint:forbidigo
			}
		}()

		err = cmd.Wait()
		if err != nil {
			logrus.Warn(err)
		}

		logrus.WithFields(logrus.Fields{
			"user":        session.User(),
			"remoteaddr":  session.RemoteAddr(),
			"localaddr":   session.LocalAddr(),
			"Raw command": session.RawCommand(),
		}).Info("Command ended")
	}
}

func (s *Server) passwordHandler(ctx sshserver.Context, pass string) bool {
	log := logrus.WithFields(logrus.Fields{
		"user": ctx.User(),
	})
	var ok bool

	if s.singleUserPassword == "" {
		ok = osauth.AuthUser(ctx.User(), pass)
	} else {
		ok = osauth.VerifyPasswordHash(s.singleUserPassword, pass)
	}

	if ok {
		log.Info("Accepted password")
	} else {
		log.Info("Failed password")
	}

	return ok
}

func (s *Server) publicKeyHandler(ctx sshserver.Context, key sshserver.PublicKey) bool {
	if osauth.LookupUser(ctx.User()) == nil {
		return false
	}

	type Signature struct {
		Username  string
		Namespace string
	}

	sig := &Signature{
		Username:  ctx.User(),
		Namespace: s.deviceName,
	}

	sigBytes, err := json.Marshal(sig)
	if err != nil {
		return false
	}

	sigHash := sha256.Sum256(sigBytes)

	res, err := s.api.AuthPublicKey(&models.PublicKeyAuthRequest{
		Fingerprint: ssh.FingerprintLegacyMD5(key),
		Data:        string(sigBytes),
	}, s.authData.Token)
	if err != nil {
		return false
	}

	digest, err := base64.StdEncoding.DecodeString(res.Signature)
	if err != nil {
		return false
	}

	cryptoKey, ok := key.(ssh.CryptoPublicKey)
	if !ok {
		return false
	}

	pubCrypto := cryptoKey.CryptoPublicKey()

	pubKey, ok := pubCrypto.(*rsa.PublicKey)
	if !ok {
		return false
	}

	if err = rsa.VerifyPKCS1v15(pubKey, crypto.SHA256, sigHash[:], digest); err != nil {
		return false
	}

	return true
}

func (s *Server) CloseSession(id string) {
	if session, ok := s.Sessions[id]; ok {
		session.Close()
		delete(s.Sessions, id)
	}
}

func newShellCmd(s *Server, username, term string) *exec.Cmd {
	shell := os.Getenv("SHELL")

	u := osauth.LookupUser(username)

	if shell == "" {
		shell = u.Shell
	}

	if term == "" {
		term = "xterm"
	}

	cmd := newCmd(u, shell, term, s.deviceName, shell, "--login")

	return cmd
}
