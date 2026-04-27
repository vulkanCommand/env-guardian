# VS Code Marketplace

The VS Code extension lives in `vscode-extension/`.

## Package Locally

```bash
cd vscode-extension
npm run package
```

This creates a `.vsix` file that can be installed manually in VS Code.

## Publish to Marketplace

Marketplace publishing requires:

- a Visual Studio Marketplace publisher
- `@vscode/vsce`
- a Personal Access Token

Publish:

```bash
cd vscode-extension
npm run publish
```

## GitHub Actions

The repository includes `.github/workflows/vscode-extension.yml`.

On version tags, it packages the extension and uploads the `.vsix` as a workflow artifact.

To enable automatic publishing later, add a `VSCE_PAT` repository secret and extend the workflow with:

```bash
npx @vscode/vsce publish -p "$VSCE_PAT"
```

## Support

- Email: `gdkalyan2109@gmail.com`
- Issues: `https://github.com/vulkanCommand/env-guardian/issues`

## References

- `https://code.visualstudio.com/api/working-with-extensions/publishing-extension`
