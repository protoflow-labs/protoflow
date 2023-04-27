export default ({ value }) => {
  console.log("Called {{ .Node.Name }}");
  return {
    value: "{{ .Node.Name }} got " + value
  };
}
