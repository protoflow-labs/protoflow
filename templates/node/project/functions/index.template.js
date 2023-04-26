{{- range index . "FunctionNodes"}}
import {{.Name}} from './{{.Name}}/index.js';
{{- end}}

export default {
  {{- range index . "FunctionNodes"}}
  {{lowercaseFirstLetter .Name}}: {{.Name}},
  {{- end}}
}
