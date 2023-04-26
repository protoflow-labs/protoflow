package templates

import (
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"os"
	"path"
	"strings"
)

//go:embed *
var Templates embed.FS

func TemplateGlob(templatePattern string, dest string, data any) error {
	matches, err := fs.Glob(Templates, templatePattern)
	if err != nil {
		return err
	}

	for _, match := range matches {
		_, fileName := path.Split(match)
		tmpl, err := template.ParseFS(Templates, match)
		if err != nil {
			return err
		}

		destFile := fmt.Sprintf("%s/%s", dest, strings.ReplaceAll(fileName, "template.", ""))

		file, err := os.Create(destFile)
		if err != nil {
			return err
		}

		err = tmpl.Execute(file, data)
		if err != nil {
			return err
		}
	}

	return nil
}
