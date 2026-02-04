package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/yacobolo/datastar-go-blueprint/internal/config"
	"github.com/yacobolo/datastar-go-blueprint/web/resources"

	"github.com/evanw/esbuild/pkg/api"
	"golang.org/x/sync/errgroup"
)

var (
	watch = false
)

func main() {
	flag.BoolVar(&watch, "watch", watch, "Enable watcher mode")
	flag.Parse()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	defer stop()

	if err := run(ctx); err != nil {
		slog.Error("failure", "error", err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	eg, egctx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		return build(egctx)
	})

	return eg.Wait()
}

func build(ctx context.Context) error {
	opts := api.BuildOptions{
		EntryPointsAdvanced: []api.EntryPoint{
			{
				InputPath:  resources.LibsDirectoryPath + "/index.ts",
				OutputPath: "libs/index",
			},
			{
				InputPath:  "web/ui/src/plugins/theme-switcher.ts",
				OutputPath: "theme-switcher",
			},
			{
				InputPath:  "web/ui/styles/main.css",
				OutputPath: "styles",
			},
		},
		Bundle:            true,
		Format:            api.FormatESModule,
		LogLevel:          api.LogLevelInfo,
		MinifyIdentifiers: true,
		MinifySyntax:      true,
		MinifyWhitespace:  true,
		NodePaths:         []string{"web/node_modules"},
		Outdir:            resources.StaticDirectoryPath,
		Sourcemap:         api.SourceMapLinked,
		Target:            api.ESNext,
		Write:             true,
	}

	if watch {
		slog.Info("watching...")

		opts.Plugins = append(opts.Plugins, api.Plugin{
			Name: "hotreload",
			Setup: func(build api.PluginBuild) {
				build.OnEnd(func(result *api.BuildResult) (api.OnEndResult, error) {
					slog.Info("build complete", "errors", len(result.Errors), "warnings", len(result.Warnings))
					if len(result.Errors) == 0 {
						http.Get(fmt.Sprintf("http://%s:%s/hotreload", config.Global.Host, config.Global.Port))
					}
					return api.OnEndResult{}, nil
				})
			},
		})

		buildCtx, err := api.Context(opts)
		if err != nil {
			return err
		}
		defer buildCtx.Dispose()

		if err := buildCtx.Watch(api.WatchOptions{}); err != nil {
			return err
		}

		<-ctx.Done()
		return nil
	}

	slog.Info("building...")

	result := api.Build(opts)

	if len(result.Errors) > 0 {
		errs := make([]error, len(result.Errors))
		for i, err := range result.Errors {
			errs[i] = errors.New(err.Text)
		}
		return errors.Join(errs...)
	}

	return nil
}
