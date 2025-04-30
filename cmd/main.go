package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"

	"github.com/cloudfresco/sc-ubl/internal/common"
	"github.com/cloudfresco/sc-ubl/internal/config"
	_ "github.com/go-sql-driver/mysql" // mysql
	"go.uber.org/zap"
)

func getKeys(log *zap.Logger, caCertPath string, certPath string, keyPath string) *tls.Config {
	caCert, err := os.ReadFile(caCertPath)
	if err != nil {
		log.Error("Error",
			zap.Error(err))
	}

	caCertpool := x509.NewCertPool()
	caCertpool.AppendCertsFromPEM(caCert)

	// LoadX509KeyPair reads files, so we give it the paths
	serverCert, err := tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		log.Error("Error",
			zap.Error(err))
	}

	tlsConfig := tls.Config{
		ClientCAs:    caCertpool,
		Certificates: []tls.Certificate{serverCert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
	}

	return &tlsConfig
}

func main() {
	v, err := config.GetViper()
	if err != nil {
		os.Exit(1)
	}

	logOpt, err := config.GetLogConfig(v)
	if err != nil {
		os.Exit(1)
	}

	log := config.SetUpLogging(logOpt.Path)

	serverOpt, err := config.GetServerConfig(log, v)
	if err != nil {
		log.Error("Error",
			zap.Int("msgnum", 103),
			zap.Error(err))
		os.Exit(1)
	}

	jwtOpt, err := config.GetJWTConfig(log, v, false, "SC_UBL_JWT_KEY", "SC_UBL_JWT_DURATION")
	if err != nil {
		log.Error("Error", zap.Error(err))
		os.Exit(1)
	}

	common.SetJWTOpt(jwtOpt)

	// Create a ServeMux for routing
	mux := http.NewServeMux()

	// Chain middleware
	chain := common.ChainMiddlewares(
		common.HandleCacheControl,
		common.CorsMiddleware,
		common.ValidateToken(serverOpt.Auth0Audience, serverOpt.Auth0Domain),
	)

	var backendServiceURL string

	if serverOpt.ServerTLS == "true" {
		backendServiceURL = "https://" + serverOpt.BackendServerAddr
	} else {
		backendServiceURL = "http://" + serverOpt.BackendServerAddr
	}

	// Create a proxy handler for the backend service
	proxyHandler := common.NewProxyHandler(backendServiceURL)

	mux.Handle("/v2.3/users", chain(proxyHandler))
	mux.Handle("/v2.3/users/", chain(proxyHandler))
	mux.Handle("/v2.3/parties", chain(proxyHandler))
	mux.Handle("/v2.3/parties/", chain(proxyHandler))
	mux.Handle("/v2.3/invoices", chain(proxyHandler))
	mux.Handle("/v2.3/invoices/", chain(proxyHandler))
	mux.Handle("/v2.3/credit-notes", chain(proxyHandler))
	mux.Handle("/v2.3/credit-notes/", chain(proxyHandler))
	mux.Handle("/v2.3/debit-notes", chain(proxyHandler))
	mux.Handle("/v2.3/debit-notes/", chain(proxyHandler))
	mux.Handle("/v2.3/consignments", chain(proxyHandler))
	mux.Handle("/v2.3/consignments/", chain(proxyHandler))
	mux.Handle("/v2.3/receipt-advices", chain(proxyHandler))
	mux.Handle("/v2.3/receipt-advices/", chain(proxyHandler))
	mux.Handle("/v2.3/despatches", chain(proxyHandler))
	mux.Handle("/v2.3/despatches/", chain(proxyHandler))
	mux.Handle("/v2.3/shipments", chain(proxyHandler))
	mux.Handle("/v2.3/shipments/", chain(proxyHandler))
	mux.Handle("/v2.3/purchase-orders", chain(proxyHandler))
	mux.Handle("/v2.3/purchase-orders/", chain(proxyHandler))
	mux.Handle("/v2.3/tax-schemes", chain(proxyHandler))
	mux.Handle("/v2.3/tax-schemes/", chain(proxyHandler))

	if serverOpt.ServerTLS == "true" {
		var caCertPath, certPath, keyPath string
		var tlsConfig *tls.Config
		pwd, _ := os.Getwd()
		caCertPath = pwd + filepath.FromSlash(serverOpt.CaCertPath)
		certPath = pwd + filepath.FromSlash(serverOpt.CertPath)
		keyPath = pwd + filepath.FromSlash(serverOpt.KeyPath)

		tlsConfig = getKeys(log, caCertPath, certPath, keyPath)

		srv := &http.Server{
			Addr:      serverOpt.ApigServerAddr,
			Handler:   mux,
			TLSConfig: tlsConfig,
		}

		idleConnsClosed := make(chan struct{})
		go func() {
			sigint := make(chan os.Signal, 1)
			signal.Notify(sigint, os.Interrupt)
			<-sigint

			// We received an interrupt signal, shut down.
			if err := srv.Shutdown(context.Background()); err != nil {
				// Error from closing listeners, or context timeout:
				log.Error("Error",
					zap.Int("msgnum", 104),
					zap.Error(errors.New("HTTP server Shutdown")))
			}
			close(idleConnsClosed)
		}()

		if err := srv.ListenAndServeTLS(certPath, keyPath); err != http.ErrServerClosed {
			// Error starting or closing listener:
			log.Error("Error",
				zap.Int("msgnum", 105),
				zap.Error(err))
		}
		log.Error("Error",
			zap.Int("msgnum", 106),
			zap.Error(err))

		<-idleConnsClosed
	} else {

		srv := &http.Server{
			Addr:    serverOpt.ApigServerAddr,
			Handler: mux,
		}

		idleConnsClosed := make(chan struct{})
		go func() {
			sigint := make(chan os.Signal, 1)
			signal.Notify(sigint, os.Interrupt)
			<-sigint

			// We received an interrupt signal, shut down.
			if err := srv.Shutdown(context.Background()); err != nil {
				// Error from closing listeners, or context timeout:
				log.Error("Error",
					zap.Int("msgnum", 107),
					zap.Error(err))
			}
			close(idleConnsClosed)
		}()

		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			// Error starting or closing listener:
			log.Error("Error",
				zap.Int("msgnum", 108),
				zap.Error(errors.New("HTTP server ListenAndServe")))
		}

		log.Error("Error",
			zap.Int("msgnum", 109),
			zap.Error(errors.New("server shutting down")))

		<-idleConnsClosed

	}
}
