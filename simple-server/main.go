package main

import (
	"flag"
	"os"

	"github.com/sirupsen/logrus"
)

type options struct {
	etcdURL            string
	etcdServicePath    string
	etcdTimeoutSeconds int
}

func main() {

	logLevel := flag.String("log-level", "debug", "debug, info, warning, error")
	etcdURL0 := flag.String("etcd-url", "", "ETCD URLs. ex: http://etcd0:2379")
	etcdServicePath0 := flag.String("etcd-service-path", "", "ETCD base path where nodes for this service were registered. ex.: /registry/myservice. required")
	etcdTimeoutSeconds0 := flag.Int("etcd-timeout-seconds", 30, "ETCD timeout seconds. defaults to 30")
	flag.Parse()

	switch *logLevel {
	case "debug":
		logrus.SetLevel(logrus.DebugLevel)
		break
	case "warning":
		logrus.SetLevel(logrus.WarnLevel)
		break
	case "error":
		logrus.SetLevel(logrus.ErrorLevel)
		break
	default:
		logrus.SetLevel(logrus.InfoLevel)
	}

	opt := options{
		etcdURL:            *etcdURL0,
		etcdServicePath:    *etcdServicePath0,
		etcdTimeoutSeconds: *etcdTimeoutSeconds0,
	}

	if opt.etcdServicePath == "" {
		logrus.Errorf("'--etcd-url' is required")
		os.Exit(1)
	}

	if opt.etcdServicePath == "" {
		logrus.Errorf("'--etcd-service-path' is required")
		os.Exit(1)
	}

	logrus.Infof("Starting ETCD hashring simple server...")

	h := NewHTTPServer(opt)
	err := h.Start()
	if err != nil {
		logrus.Errorf("Error starting server. err=%s", err)
		os.Exit(1)
	}

}
