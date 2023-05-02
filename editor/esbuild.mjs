import esbuild from "esbuild";
import {postcssModules, sassPlugin} from "esbuild-sass-plugin";
import tailwindcss from 'tailwindcss';
import autoprefixer from 'autoprefixer';
import path from "path";
import postcss from "postcss";

const prodBuild = process.env.BUILD === 'true'

const watch = !prodBuild ? {
  onRebuild: () => {
    console.log("rebuilt!");
  },
} : undefined;

const minify = prodBuild;

const nodeEnv = prodBuild ? "'production'" : "'development'";

const options = {
      entryPoints: [
          "./src/index.tsx",
          "./src/styles/globals.scss"
      ],
      outdir: "public/build/",
      bundle: true,
      loader: {
        ".ts": "tsx",
        ".tsx": "tsx",
        ".woff2": "file",
        ".woff": "file",
      },
      plugins: [
        sassPlugin({
        filter: /\.scss$/,
        nonce: "window.__esbuild_nonce__",
        type: "style",
        async transform(source, resolveDir, filePath) {
          // Transform TailwindCSS
          const transformed = await postcss(
            autoprefixer,
            tailwindcss()
          ).process(source, { from: filePath });

          // Run CSS modules transformer
          const result = postcssModules({});
          return result(transformed.css, resolveDir, filePath);
        },
      }),
      ],
      minify: false,
      sourcemap: "linked",
      define: {
        "global": "window",
        "process": "{}",
        "process.env": "{}",
        "process.env.NODE_ENV": nodeEnv,
      },
      logLevel: 'info'
    };

if (prodBuild) {
  await esbuild.build(options);
} else {
  const context = await esbuild
    .context(options);

  const result = await context.rebuild()
  context.serve({
      servedir: 'public',
  })
  await context.watch()
  // maybe think of live reload? https://esbuild.github.io/api/#live-reload
  // process.stdin.on('data', async () => {
  //   try {
  //     // Cancel the already-running build
  //     await context.cancel()

  //     // Then start a new build
  //     console.log('build:', await context.rebuild())
  //   } catch (err) {
  //     console.error(err)
  //   }
  // })
}
