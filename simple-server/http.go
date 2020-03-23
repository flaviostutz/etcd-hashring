package main

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
	"github.com/sirupsen/logrus"

	cors "github.com/itsjamie/gin-cors"
)

type HTTPServer struct {
	server *http.Server
	router *gin.Engine
}

func NewHTTPServer(opt options) *HTTPServer {
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

	etcdHashring := ETCDHashring()

	h := &HTTPServer{
		server: &http.Server{
					Addr:    ":3000",
					Handler: router,
				}, 
				router: router,
				etcdHashring: etcdHashring,
		}

	router.GET("/element/:elementId/node", h.getNode())

	return h
}

func (h *HTTPServer) getNode() func(*gin.Context) {
	return func(c *gin.Context) {
		elementId := c.Param("elementId")

		

		c.Header("Cache-Control", "no-cache")
		c.Render(
			http.StatusOK,
			render.Data{
				ContentType: "application/json",
				Data:        layerBytes,
			}
		)
	}
}

//Start the main HTTP Server entry
func (s *HTTPServer) Start() error {
	logrus.Infof("Starting HTTP Server on port 3000")
	return s.server.ListenAndServe()
}

