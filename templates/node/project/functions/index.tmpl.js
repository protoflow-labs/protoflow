{{- range .Methods}}
import {{.Name}} from './{{.Name}}/index.js';
{{- end}}

export default {
  {{- range .Methods}}
  "{{lowercaseFirstLetter .Name}}": {{.Name}},
  {{- end}}
}
