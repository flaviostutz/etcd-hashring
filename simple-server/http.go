package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	hashring "github.com/flaviostutz/etcd-hashring"
	cors "github.com/itsjamie/gin-cors"
)

type HTTPServer struct {
	server       *http.Server
	router       *gin.Engine
	etcdHashring *hashring.ETCDHashring
}

func NewHTTPServer(opt options) (*HTTPServer, error) {
	router := gin.Default()

	router.Use(cors.Middleware(cors.Config{
		Origins:         "*",
		Methods:         "GET",
		RequestHeaders:  "Origin, Content-Type",
		ExposedHeaders:  "",
		MaxAge:          1 * time.Hour,
		Credentials:     false,
		ValidateHeaders: false,
	}))

	etcdHashring, err := hashring.NewETCDHashring(opt.etcdURLs, opt.etcdServicePath, opt.etcdTimeout)
	if err != nil {
		return nil, err
	}

	logrus.Debugf("Running ETCD Hashring...")
	go etcdHashring.Run()

	h := &HTTPServer{
		server: &http.Server{
			Addr:    ":3000",
			Handler: router,
		},
		router:       router,
		etcdHashring: etcdHashring,
	}

	logrus.Debugf("Registering /config/:elementID route")
	router.GET("/config/:elementID", h.getNode())

	return h, nil
}

func (h *HTTPServer) getNode() func(*gin.Context) {
	return func(c *gin.Context) {
		elementID := c.Param("elementID")

		logrus.Debugf("Getting node for element %s", elementID)
		node, err := h.etcdHashring.GetNode(elementID)
		if err != nil {
			logrus.Debugf("Error getting node for element %s. err=%s", elementID, err)
			c.Header("Cache-Control", "no-cache")
			c.JSON(http.StatusNotFound, gin.H{"message": fmt.Sprintf("Could not select a node for element %s", elementID)})
			return
		}

		c.Header("Cache-Control", "no-cache")
		c.JSON(http.StatusOK, gin.H{"date": time.Now(), "nodeName": node.Name, "nodeInfo": node.Info})
	}
}

//Start the main HTTP Server entry
func (s *HTTPServer) Start() error {
	logrus.Infof("Starting HTTP Server on port 3000")
	return s.server.ListenAndServe()
}
