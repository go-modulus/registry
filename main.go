package main

import (
	"cmp"
	"fmt"
	"slices"

	"github.com/fatih/color"
	"github.com/go-modulus/auth"
	"github.com/go-modulus/auth/providers/email"
	"github.com/go-modulus/auth/providers/google"
	"github.com/go-modulus/chihttp"
	"github.com/go-modulus/graphql"
	"github.com/go-modulus/modulus/captcha"
	"github.com/go-modulus/modulus/cli"
	"github.com/go-modulus/modulus/http"
	"github.com/go-modulus/modulus/logger"
	"github.com/go-modulus/modulus/module"
	"github.com/go-modulus/modulus/translation"
	"github.com/go-modulus/pgx"
	"github.com/go-modulus/pgx/migrator"
	"github.com/go-modulus/registry/internal"
	"github.com/go-modulus/temporal"

	"os"
)

func main() {
	modules := []module.Manifesto{
		cli.NewManifesto(),
		pgx.NewManifesto(),
		migrator.NewManifesto(),
		logger.NewManifesto(),
		http.NewManifesto(),
		chihttp.NewManifesto(),
		graphql.NewManifesto(),
		auth.NewManifesto(),
		email.NewManifesto(),
		temporal.NewManifesto(),
		captcha.NewManifesto(),
		google.NewManifesto(),
		translation.NewManifesto(),
	}

	registry := internal.Registry{
		Modules:     make([]module.Manifesto, 0),
		Version:     "1.0.0",
		Name:        "Modulus framework modules registry",
		Description: "List of available modules for the Modulus framework",
	}

	slices.SortFunc(modules, func(a, b module.Manifesto) int { return cmp.Compare(a.Name, b.Name) })

	for _, currentModule := range modules {
		fmt.Println("Updating module", color.BlueString(currentModule.Name))
		registry.UpdateModule(currentModule)
	}

	err := registry.SaveAsLocalFile("./")
	if err != nil {
		fmt.Println("Cannot save the registry file modules.json:", color.RedString(err.Error()))
		os.Exit(1)
	}
}
