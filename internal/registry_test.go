package internal

import (
	"os"
	"path/filepath"
	"testing"
	"testing/fstest"

	"github.com/go-modulus/modulus/module"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func makeManifesto(pkg string, name string, isLocal bool) module.Manifesto {
	return module.Manifesto{
		Name:          name,
		Package:       pkg,
		Description:   "desc",
		Version:       "1.0.0",
		IsLocalModule: isLocal,
	}
}

func TestReadFromJSON(t *testing.T) {
	data := []byte(`{
		"name": "test",
		"description": "test desc",
		"version": "1.0.0",
		"modules": [
			{"name": "mod1", "package": "github.com/example/mod1", "version": "1.0.0"}
		]
	}`)

	r := &Registry{}
	err := r.ReadFromJSON(data)
	require.NoError(t, err)
	assert.Equal(t, "test", r.Name)
	assert.Equal(t, "test desc", r.Description)
	assert.Equal(t, "1.0.0", r.Version)
	require.Len(t, r.Modules, 1)
	assert.Equal(t, "mod1", r.Modules[0].Name)
	assert.Equal(t, "github.com/example/mod1", r.Modules[0].Package)
}

func TestReadFromJSON_Invalid(t *testing.T) {
	r := &Registry{}
	err := r.ReadFromJSON([]byte(`not json`))
	assert.Error(t, err)
}

func TestWriteToJSON(t *testing.T) {
	r := &Registry{
		Name:        "test",
		Description: "desc",
		Version:     "1.0.0",
		Modules:     []module.Manifesto{makeManifesto("github.com/example/mod1", "mod1", false)},
	}

	data, err := r.WriteToJSON()
	require.NoError(t, err)

	r2 := &Registry{}
	require.NoError(t, r2.ReadFromJSON(data))
	assert.Equal(t, r.Name, r2.Name)
	assert.Equal(t, r.Version, r2.Version)
	require.Len(t, r2.Modules, 1)
	assert.Equal(t, "github.com/example/mod1", r2.Modules[0].Package)
}

func TestAddModule(t *testing.T) {
	r := &Registry{}
	m1 := makeManifesto("github.com/example/mod1", "mod1", false)
	m2 := makeManifesto("github.com/example/mod2", "mod2", false)

	r.AddModule(m1)
	r.AddModule(m2)

	require.Len(t, r.Modules, 2)
	assert.Equal(t, "github.com/example/mod1", r.Modules[0].Package)
	assert.Equal(t, "github.com/example/mod2", r.Modules[1].Package)
}

func TestUpdateModule_ExistingModule(t *testing.T) {
	r := &Registry{}
	r.AddModule(makeManifesto("github.com/example/mod1", "mod1", false))

	updated := makeManifesto("github.com/example/mod1", "mod1-updated", false)
	r.UpdateModule(updated)

	require.Len(t, r.Modules, 1)
	assert.Equal(t, "mod1-updated", r.Modules[0].Name)
}

func TestUpdateModule_NewModule(t *testing.T) {
	r := &Registry{}
	r.AddModule(makeManifesto("github.com/example/mod1", "mod1", false))

	r.UpdateModule(makeManifesto("github.com/example/mod2", "mod2", false))

	require.Len(t, r.Modules, 2)
}

func TestNewFromFs(t *testing.T) {
	jsonData := `{
		"name": "test",
		"description": "desc",
		"version": "2.0.0",
		"modules": []
	}`
	memFs := fstest.MapFS{
		"modules.json": &fstest.MapFile{Data: []byte(jsonData)},
	}

	reg, err := NewFromFs(memFs, "modules.json")
	require.NoError(t, err)
	assert.Equal(t, "test", reg.Name)
	assert.Equal(t, "2.0.0", reg.Version)
}

func TestNewFromFs_FileNotFound(t *testing.T) {
	memFs := fstest.MapFS{}
	_, err := NewFromFs(memFs, "missing.json")
	assert.Error(t, err)
}

func TestNewFromFs_InvalidJSON(t *testing.T) {
	memFs := fstest.MapFS{
		"modules.json": &fstest.MapFile{Data: []byte(`bad json`)},
	}
	_, err := NewFromFs(memFs, "modules.json")
	assert.Error(t, err)
}

func TestSaveAsLocalFile(t *testing.T) {
	dir := t.TempDir()
	r := &Registry{
		Name:    "saved",
		Version: "1.0.0",
		Modules: []module.Manifesto{makeManifesto("github.com/example/mod1", "mod1", false)},
	}

	err := r.SaveAsLocalFile(dir)
	require.NoError(t, err)

	data, err := os.ReadFile(filepath.Join(dir, "modules.json"))
	require.NoError(t, err)

	r2 := &Registry{}
	require.NoError(t, r2.ReadFromJSON(data))
	assert.Equal(t, "saved", r2.Name)
	require.Len(t, r2.Modules, 1)
}

func TestLoadLocalRegistry_WithModulesJSON(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, os.MkdirAll(filepath.Join(dir, "cmd", "app"), 0755))
	require.NoError(t, os.WriteFile(filepath.Join(dir, "cmd", "app", "main.go"), []byte("package main"), 0644))

	jsonData := `{
		"name": "local-registry",
		"description": "local desc",
		"version": "1.2.3",
		"modules": [{"name": "mod1", "package": "github.com/example/mod1", "version": "1.0.0"}]
	}`
	require.NoError(t, os.WriteFile(filepath.Join(dir, "modules.json"), []byte(jsonData), 0644))

	reg, err := LoadLocalRegistry(dir)
	require.NoError(t, err)
	assert.Equal(t, "local-registry", reg.Name)
	assert.Equal(t, "1.2.3", reg.Version)
	require.Len(t, reg.Modules, 1)
}

func TestLoadLocalRegistry_WithoutModulesJSON(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, os.MkdirAll(filepath.Join(dir, "cmd", "app"), 0755))
	require.NoError(t, os.WriteFile(filepath.Join(dir, "cmd", "app", "main.go"), []byte("package main"), 0644))

	reg, err := LoadLocalRegistry(dir)
	require.NoError(t, err)
	assert.Equal(t, "Modulus framework modules manifest", reg.Name)
	assert.Equal(t, "1.0.0", reg.Version)
	assert.Empty(t, reg.Modules)
}
