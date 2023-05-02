package editor

import (
	"github.com/evanw/esbuild/pkg/api"
	"strconv"
	"strings"
)

// TODO breadchris this is a WIP, theoretically we could build the editor with esbuild from go
func Build(prodBuild bool) api.BuildResult {
	options := api.BuildOptions{
		EntryPoints: []string{
			"./src/index.tsx",
			"./src/styles/globals.css",
		},
		Outdir: "public/build/",
		Bundle: true,
		Loader: map[string]api.Loader{
			".ts":    api.LoaderTS,
			".tsx":   api.LoaderTSX,
			".woff2": api.LoaderFile,
			".woff":  api.LoaderFile,
		},
		// TODO breadchris how do you use js plugins from inside go?
		Plugins: []func(*api.Plugin){
			postcssModules(),
			sassPlugin(api.SassPluginOptions{
				Filter: func(filename string) bool {
					return strings.HasSuffix(filename, ".css")
				},
				Async: func(source []byte, resolveDir string, _ string) (api.SassPluginResult, error) {
					result, err := postcss([]postcss.Plugin{
						tailwindcss.Tailwind("./tailwind.config.js"),
						autoprefixer.Autoprefixer(),
					}).Process(source)
					if err != nil {
						return esbuild.SassPluginResult{}, err
					}
					return esbuild.SassPluginResult{
						Contents: result.CSS,
					}, nil
				},
			}),
			swcPlugin(swc.Config{}),
		},
		MinifySyntax: prodBuild,
		Sourcemap:    api.SourceMapLinked,
		Define: map[string]string{
			"global":               "window",
			"process":              "{}",
			"process.env":          "{}",
			"process.env.NODE_ENV": strings.Title(strings.ToLower(strconv.FormatBool(prodBuild))),
		},
		LogLevel: api.LogLevelInfo,
	}

	var buildRes api.BuildResult
	if prodBuild {
		buildRes = api.Build(options)
	} else {
		serveOptions := api.ServeOptions{
			ServeDir: "public",
		}
		c, err := api.NewWatchContext(options)
		if err != nil {
			return err
		}
		c.WatchFiles(func(_ string) {
			_ = c.Restart()
		})
		err = c.Serve(serveOptions)
		if err != nil {
			return err
		}
	}
	return buildRes
}
