package templates

import (
	"embed"
	"fmt"
	"github.com/protoflow-labs/protoflow/pkg/util"
	"html/template"
	"io/fs"
	"os"
	"path"
	"strings"
)

//go:embed *
var Templates embed.FS

func TemplateFile(templateFile string, dest string, data any) error {
	_, filename := path.Split(templateFile)
	tmpl, err := template.New(filename).Funcs(templateHelpers()).ParseFS(Templates, templateFile)
	if err != nil {
		return err
	}

	os.MkdirAll(path.Dir(dest), 0700)

	file, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer file.Close()

	err = tmpl.Execute(file, data)
	if err != nil {
		return err
	}

	return nil
}

func TemplateDir(templateDir string, destDir string, data any) error {
	return fs.WalkDir(Templates, templateDir, func(templateFile string, entry fs.DirEntry, _ error) error {
		if !entry.IsDir() {
			relativePath := strings.ReplaceAll(templateFile, templateDir, "")
			destFile := fmt.Sprintf("%s/%s", destDir, strings.ReplaceAll(relativePath, "template.", ""))

			err := TemplateFile(templateFile, destFile, data)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func templateHelpers() template.FuncMap {
	return template.FuncMap{
		"lowercaseFirstLetter": func(s string) string {
			if len(s) == 0 {
				return s
			}
			first := strings.ToLower(string(s[0]))
			return first + s[1:]
		},
		"titlecase": func(s string) string {
			return util.ToTitleCase(s)
		},
	}
}
