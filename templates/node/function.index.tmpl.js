export default ({ value }) => {
  console.log("Called {{ .Name }}");
  return {
    value: "{{ .Name }} got " + value
  };
}
