import { glob } from "glob";

export async function register() {
  const files = glob.sync(process.cwd() + "/functions/**/*.js");
  const registry = {};
  for (const file of files) {
    const name = file.split("/").splice(-2, 1).pop();
    await import(file).then((module) => {
      registry[name] = module.default;
    });
  }

  return registry;
}
