import esbuild from "esbuild";

const prodBuild = process.env.BUILD === 'true'
const target = process.env.TARGET || 'site'
const buildDir = prodBuild ? 'dist' : 'build'

const buildVSCode = target === 'vscode' || prodBuild
const buildSite = target === 'site' || prodBuild

async function doBuild(options, serve) {
    if (prodBuild) {
        await esbuild.build(options);
    } else {
        try {
            const context = await esbuild
                .context(options);

            await context.rebuild()
            if (serve) {
                console.log('serving', `${buildDir}/site`)
                context.serve({
                    port: 8001,
                    servedir: `${buildDir}/site`,
                    fallback: `${buildDir}/site/index.html`,
                    onRequest: args => {
                        console.log(args.method, args.path)
                    }
                })
            }
            await context.watch()
        } catch (e) {
            console.error('failed to build: ' + e)
        }
    }
}

const baseOptions = {
    bundle: true,
    loader: {
        ".ts": "tsx",
        ".tsx": "tsx",
        ".woff2": "file",
        ".woff": "file",
        ".html": "copy",
        ".json": "copy",
        ".ico": "copy",
    },
    plugins: [],
    minify: false,
    sourcemap: "linked",
    define: {
        "global": "window",
    },
    logLevel: 'info'
};

if (buildSite) {
    await doBuild({
        ...baseOptions,
        entryPoints: [
            "./src/index.tsx",
            "./src/styles/globals.css",
            "./src/favicon.ico",
            "./src/index.html",
        ],
        outdir: `${buildDir}/site/`,
    }, true);
}

if (buildVSCode) {
    await doBuild({
        ...baseOptions,
        entryPoints: [
            "./src/extension/extension.ts",
        ],
        outdir: `${buildDir}/extension/`,
    }, false);
}
