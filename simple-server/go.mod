module github.com/flaviostutz/etcd-hashring-simple-server

go 1.14

require (
	github.com/flaviostutz/etcd-hashring v0.0.0-20200323221257-30e9224a842e
	github.com/gin-gonic/gin v1.6.1
	github.com/itsjamie/gin-cors v0.0.0-20160420130702-97b4a9da7933
	github.com/sirupsen/logrus v1.5.0
)

replace (
	github.com/coreos/go-systemd => github.com/coreos/go-systemd/v22 v22.0.0
	google.golang.org/grpc => google.golang.org/grpc v1.26.0
	github.com/flaviostutz/etcd-hashring => ../
)