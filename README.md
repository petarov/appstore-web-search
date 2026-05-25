# App Store Web Search

A [web app](https://vexelon.net/asws) that queries the Apple App Store in your browser.

<img src="demo/shot1.png" width="300">

# Development

```
npm install
npm run dev
```

Opens a Vite dev server at `http://localhost:5173/asws/`.

# Build

```
npm run build
```

Runs EditorConfig check, compiles and obfuscates the app into `dist/`, then verifies the bundle is ES6-compatible. Output is a static `dist/index.html` ready to deploy to any web host.

```
npm run preview
```

Serves the built `dist/` locally to sanity-check the artifact.

# Test

```
npm run test
```

Alias of `npm run build` — runs both the EditorConfig and ES6 checks.

# License

[MIT](LICENSE)
