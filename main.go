package main

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/go-modulus/modulus/auth"
	"github.com/go-modulus/modulus/auth/providers/email"
	"github.com/go-modulus/modulus/auth/providers/google"
	"github.com/go-modulus/modulus/captcha"
	"github.com/go-modulus/modulus/cli"
	"github.com/go-modulus/modulus/db/embedded"
	"github.com/go-modulus/modulus/db/migrator"
	"github.com/go-modulus/modulus/db/pgx"
	"github.com/go-modulus/modulus/graphql"
	"github.com/go-modulus/modulus/http"
	httpMiddleware "github.com/go-modulus/modulus/http/middleware"
	"github.com/go-modulus/modulus/logger"
	"github.com/go-modulus/modulus/module"
	"github.com/go-modulus/modulus/temporal"
	"github.com/go-modulus/modulus/translation"
	"github.com/go-modulus/registry/internal"

	"os"
)

func main() {
	modules := []module.Manifesto{
		module.NewManifesto(
			cli.NewModule(),
			"github.com/go-modulus/modulus/cli",
			"Adds ability to create cli applications in the Modulus framework.",
			"1.0.0",
		),
		pgx.NewManifesto(),
		module.NewManifesto(
			logger.NewModule(),
			"github.com/go-modulus/modulus/logger",
			"Adds a slog logger with a zap backend to the Modulus framework.",
			"1.0.0",
		),
		module.NewManifesto(
			migrator.NewModule(),
			"github.com/go-modulus/modulus/db/migrator",
			"Several CLI commands to use DBMate (https://github.com/amacneil/dbmate) migration tool inside your application.",
			"1.0.0",
		),
		module.NewManifesto(
			http.NewModule(),
			"github.com/go-modulus/modulus/http",
			"HTTP module based on the Chi router.",
			"1.0.0",
		),
		module.NewManifesto(
			httpMiddleware.NewModule(),
			"github.com/go-modulus/modulus/http/middleware",
			"Various useful middlewares",
			"1.0.0",
		),
		graphql.NewManifesto(),
		auth.NewManifesto(),
		email.NewManifesto(),
		embedded.NewManifesto(),
		temporal.NewManifesto(),
		captcha.NewManifesto(),
		google.NewManifesto(),
		translation.NewManifesto(),
	}

	registry, err := internal.LoadLocalRegistry("./")
	if err != nil {
		fmt.Println("Cannot load the registry file modules.json:", color.RedString(err.Error()))
		os.Exit(1)
	}

	for _, currentModule := range modules {
		fmt.Println("Updating module", color.BlueString(currentModule.Name))
		registry.UpdateModule(currentModule)
	}

	err = registry.SaveAsLocalFile("./")
	if err != nil {
		fmt.Println("Cannot save the registry file modules.json:", color.RedString(err.Error()))
		os.Exit(1)
	}
}
