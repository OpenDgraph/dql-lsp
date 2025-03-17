const path = require('path');
const vscode = require('vscode');
const { LanguageClient, TransportKind } = require('vscode-languageclient/node');

let client;

function activate(context) {

    const isDev = process.env.DEBUG === "true";

    const serverCommand = isDev ? 'go' : path.join(__dirname, '..', 'dql-lsp');
    const serverArgs = isDev ? ['run', './cmd/main.go'] : [];
    

    const serverOptions = {
        run: {
            command: serverCommand,
            args: serverArgs,
            transport: TransportKind.stdio,
            options: {
                cwd: path.join(__dirname, '..'), 
                env: {
                    ...process.env,
                    DEBUG: "true"
                }
            }
        },
        debug: {
            command: serverCommand,
            args: serverArgs,
            transport: TransportKind.stdio,
            options: {
                cwd: path.join(__dirname, '..'),
                env: {
                    ...process.env,
                    DEBUG: "true"
                }
            }
        }
    };

    const clientOptions = {
        documentSelector: [
            { scheme: 'file', language: 'dql' },
            { scheme: 'file', language: 'schema' },
            { scheme: 'file', pattern: '**/*.dql' },
            { scheme: 'file', pattern: '**/*.schema' }
        ]
    };
    client = new LanguageClient('dqlLsp', 'DQL Language Server', serverOptions, clientOptions);
    client.start();
}

function deactivate() {
    if (!client) {
        return undefined;
    }
    return client.stop();
}

module.exports = { activate, deactivate };
