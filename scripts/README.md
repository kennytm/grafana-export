# User script for exporting Grafana dashboard snapshot as JSON file

This package contains the user script to insert the "Export" dialog into Grafana.

## Development

This script is built using [Yarn](https://yarnpkg.com/getting-started/install) and TypeScript.

| Command        | Explanation                                                            |
|----------------|------------------------------------------------------------------------|
| `yarn`         | Initiate the development environment                                   |
| `yarn build`   | Compile the TypeScript sources, produces `dist/index.js`               |
| `yarn release` | Generate minimized user-script for each language, rewrite `../index.html` |
