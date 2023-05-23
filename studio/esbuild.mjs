import esbuild from "esbuild";
import {postcssModules, sassPlugin} from "esbuild-sass-plugin";
import { swcPlugin } from "esbuild-plugin-swc";
import polyfill from "esbuild-plugin-node-polyfills";
import tailwindcss from 'tailwindcss';
import autoprefixer from 'autoprefixer';
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
        "./src/styles/globals.css"
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
        // TODO breadchris use swc over tsc
        // swcPlugin(),
        polyfill,
        sassPlugin({
            filter: /\.css$/,
            type: "style",
            async transform(source, resolveDir, filePath) {
                const transformed = await postcss([
                    tailwindcss({
                        content: [
                            './src/**/*.{js,jsx,ts,tsx}',
                        ],
                        theme: {
                            extend: {},
                        },
                        plugins: [],
                    }),
                    autoprefixer,
                ]).process(source, {from: filePath});
                return transformed.css;
            },
        }),
    ],
    minify: false,
    sourcemap: "linked",
    define: {
        "global": "window",
    },
    logLevel: 'info'
};

if (prodBuild) {
    await esbuild.build(options);
} else {
    try {
        const context = await esbuild
            .context(options);

        await context.rebuild()
        context.serve({
            servedir: 'public',
        })
        await context.watch()
    } catch (e) {
        console.error('failed to build: ' + e)
    }
}
