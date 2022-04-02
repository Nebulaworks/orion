package server

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/charmbracelet/wish"
	"github.com/charmbracelet/wish/bubbletea"
	"github.com/charmbracelet/wish/logging"
	"github.com/charmbracelet/wish/scp"
	"github.com/gliderlabs/ssh"
	"github.com/nebulaworks/orion/apps/term-apply/pkg/applicant"
	"github.com/nebulaworks/orion/apps/term-apply/pkg/auth"
	"github.com/nebulaworks/orion/apps/term-apply/pkg/transfer"
	"github.com/nebulaworks/orion/apps/term-apply/pkg/ui"
)

type Server struct {
	ws   *ssh.Server
	host string
	port int
}

func NewServer(host, uploadDir, dataFile string, port int) (*Server, error) {
	am, err := applicant.NewApplicantManager(dataFile, uploadDir)
	if err != nil {
		return nil, err
	}
	tm := ui.NewTeaManager(am)

	const SECONDS_FIVE_MINUTES = 300
	ws, err := wish.NewServer(
		ssh.PublicKeyAuth(auth.PkHandler),
		wish.WithAddress(fmt.Sprintf("%s:%d", host, port)),
		wish.WithHostKeyPath(".ssh/term_info_ed25519"),
		wish.WithMaxTimeout(time.Second*time.Duration(SECONDS_FIVE_MINUTES)),
		wish.WithMiddleware(
			scp.Middleware(
				transfer.NewNilCopyHandler(),
				transfer.NewCopyFromClientHandler(uploadDir)),
			bubbletea.Middleware(tm.TeaHandler),
			logging.Middleware(),
		),
	)
	if err != nil {
		return &Server{}, err
	}
	return &Server{
		ws:   ws,
		host: host,
		port: port,
	}, nil
}

func (s *Server) Start() {
	log.Printf("Starting SSH server on %s:%d", s.host, s.port)
	go func() {
		if err := s.ws.ListenAndServe(); err != nil {
			log.Printf("Server failed: %v", err)
			log.Fatal("exiting...")
		}
	}()
}

func (s *Server) Stop(ctx context.Context) {
	log.Println("Stopping SSH server")
	if err := s.ws.Shutdown(ctx); err != nil {
		log.Printf("Server failed to stop %v", err)
		log.Fatal("exiting...")
	}
}
