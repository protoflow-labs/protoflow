import esbuild from "esbuild";
import {spawn, spawnSync} from "child_process";

const prodBuild = process.env.BUILD === 'true'
const target = process.env.TARGET || 'site'
const buildDir = prodBuild ? 'dist' : 'build'

const buildVSCode = target === 'vscode' || prodBuild
const buildSite = target === 'site' || prodBuild

const runTailwindBuild = (watch, outfile) => {
    console.log("Building Tailwind CSS...");
    try {
        const command = 'npx';
        const args = [
            'tailwindcss',
            'build',
            '-i', 'src/styles/tailwind.css',
            '-o', outfile
        ];

        if (watch) {
            args.push('--watch')
            spawn(command, args, {
                stdio: 'inherit'
            })
        } else {
            spawnSync(command, args, {
                stdio: 'inherit'
            });
        }
        console.log("Tailwind CSS build successful!");
    } catch (error) {
        console.error("Error building Tailwind CSS:", error.message);
    }
};

async function doBuild(options, serve) {
    // TODO breadchris support tailwind for extension
    if (buildSite) {
        runTailwindBuild(!prodBuild, `${options.outdir}/tailwind.css`);
    }
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

// if (buildVSCode) {
//     await doBuild({
//         ...baseOptions,
//         entryPoints: [
//             "./src/extension/extension.ts",
//         ],
//         outdir: `${buildDir}/extension/`,
//     }, false);
// }
