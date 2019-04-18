package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"
	
	"github.com/kardianos/service"

	_ "net/http/pprof"

	"github.com/docker/distribution/registry"
	_ "github.com/docker/distribution/registry/auth/htpasswd"
	_ "github.com/docker/distribution/registry/auth/silly"
	_ "github.com/docker/distribution/registry/auth/token"
	_ "github.com/docker/distribution/registry/proxy"
	_ "github.com/docker/distribution/registry/storage/driver/azure"
	_ "github.com/docker/distribution/registry/storage/driver/filesystem"
	_ "github.com/docker/distribution/registry/storage/driver/gcs"
	_ "github.com/docker/distribution/registry/storage/driver/inmemory"
	_ "github.com/docker/distribution/registry/storage/driver/middleware/cloudfront"
	_ "github.com/docker/distribution/registry/storage/driver/middleware/redirect"
	_ "github.com/docker/distribution/registry/storage/driver/oss"
	_ "github.com/docker/distribution/registry/storage/driver/s3-aws"
	_ "github.com/docker/distribution/registry/storage/driver/swift"
)

var logger service.Logger

// Program structures.
//  Define Start and Stop methods.
type program struct {
	exit chan struct{}
}

func (p *program) Start(s service.Service) error {
	if service.Interactive() {
		logger.Info("Running in terminal.")
	} else {
		logger.Info("Running under service manager.")
	}
	p.exit = make(chan struct{})

	// Start should not block. Do the actual work async.
	go p.run()
	return nil
}
func (p *program) run() error {
	configPath, err := getConfigPath()
	if err != nil {
		log.Fatal(err)
		return nil
	}
	registry.RootCmd.SetArgs([]string{"serve", configPath})
	registry.RootCmd.Execute()
	return nil
}
func (p *program) Stop(s service.Service) error {
	// Any work in Stop should be quick, usually a few seconds at most.
	logger.Info("Docker Registry is Stopping!")
	close(p.exit)
	return nil
}

func getConfigPath() (string, error) {
	fullexecpath, err := os.Executable()
	if err != nil {
		return "", err
	}

	dir, _ := filepath.Split(fullexecpath)

	return filepath.Join(dir, "config.yml"), nil
}

// Service setup.
//   Define service config.
//   Create the service.
//   Setup the logger.
//   Handle service controls (optional).
//   Run the service.
//   ./registry -service install
//   ./registry -service uninstall
func main() {
	svcFlag := flag.String("service", "", "Control the system service.")
	flag.Parse()

	svcConfig := &service.Config{
		Name:        "DockerRegistry",
		DisplayName: "Docker Registry",
		Description: "Docker Registry",
	}

	prg := &program{}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		log.Fatal(err)
	}
	errs := make(chan error, 5)
	logger, err = s.Logger(errs)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for {
			err := <-errs
			if err != nil {
				log.Print(err)
			}
		}
	}()

	if len(*svcFlag) != 0 {
		err := service.Control(s, *svcFlag)
		if err != nil {
			log.Printf("Valid actions: %q\n", service.ControlAction)
			log.Fatal(err)
		}
		return
	}
	err = s.Run()
	if err != nil {
		logger.Error(err)
	}
}
