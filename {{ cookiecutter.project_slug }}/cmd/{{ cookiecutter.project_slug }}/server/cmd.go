package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"{{ cookiecutter.project_slug }}/internal/config"
	"{{ cookiecutter.project_slug }}/pkg/logger"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var CmdRun = &cobra.Command{
	Use:   "run",
	Short: "Run the {{ cookiecutter.project_slug }} server",
	Run: func(*cobra.Command, []string) {
		////////////////SETTING LOGS && Sentry ///////////////////
		logger.Setup()
		//////////////////////////////////////////////////////////
		app := InitApp()
		log.Info().Msg("starting http server")
		go app.Listen(fmt.Sprintf(":%d", config.Get().HTTPPort)) //nolint:errcheck

		ctx, cancel := context.WithCancel(context.Background())
		wg := &sync.WaitGroup{}
		// Start the app
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-ctx.Done()
			if err := app.Shutdown(); err != nil {
				log.Fatal().Err(err).Msg("failed to shutdown")
			}
		}()

		prometheusAddr := fmt.Sprintf(":%d", config.Get().PrometheusPort)
		mux := http.NewServeMux()
		mux.Handle("/metrics", promhttp.Handler())
		server := &http.Server{Addr: prometheusAddr, Handler: mux}

		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Fatal().Err(err).Msg("failed to start prometheus server")
			}
		}()

		wg.Add(1)
		go func() {
			defer wg.Done()
			<-ctx.Done()
			if err := server.Shutdown(ctx); err != nil {
				log.Fatal().Err(err).Msg("failed to shutdown prometheus server")
			}
		}()

		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

		log.Info().Msg("initialize complete")
		<-sigs // Blocks here until interrupted
		// Handle shutdown
		fmt.Println("Shutdown signal received")
		cancel()
		wg.Wait()
		fmt.Println("All workers done, shutting down!")
	},
}
