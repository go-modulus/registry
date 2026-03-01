package internal

import (
	"encoding/json"
	"io/fs"
	"os"

	"github.com/go-modulus/modulus/module"
)

type Registry struct {
	Name        string             `json:"name"`
	Description string             `json:"description"`
	Version     string             `json:"version"`
	Modules     []module.Manifesto `json:"modules"`
}

func (m *Registry) ReadFromJSON(data []byte) error {
	return json.Unmarshal(data, &m)
}

func (m *Registry) WriteToJSON() ([]byte, error) {
	return json.MarshalIndent(m, "", "  ")
}

func (m *Registry) AddModule(moduleManifesto module.Manifesto) {
	m.Modules = append(m.Modules, moduleManifesto)
}

func (m *Registry) UpdateModule(module module.Manifesto) {
	for i, mod := range m.Modules {
		if mod.Package == module.Package {
			m.Modules[i] = module
			return
		}
	}
	m.AddModule(module)
}

func NewFromFs(manifestFs fs.FS, filename string) (*Registry, error) {
	data, err := fs.ReadFile(manifestFs, filename)
	if err != nil {
		return nil, err
	}
	m := &Registry{}
	err = m.ReadFromJSON(data)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func LoadLocalRegistry(projPath string) (Registry, error) {
	res := Registry{
		Modules:     make([]module.Manifesto, 0),
		Version:     "1.0.0",
		Name:        "Modulus framework modules manifest",
		Description: "List of installed modules for the Modulus framework",
	}
	if fileExists(projPath + "/modules.json") {
		projFs := os.DirFS(projPath)
		manifest, err := NewFromFs(projFs, "modules.json")
		if err != nil {
			return res, err
		}
		return *manifest, nil
	}
	return res, nil
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func (m *Registry) SaveAsLocalFile(projPath string) error {
	data, err := m.WriteToJSON()
	if err != nil {
		return err
	}
	return os.WriteFile(projPath+"/modules.json", data, 0644)
}
