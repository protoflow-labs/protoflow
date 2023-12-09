import * as vscode from 'vscode';

export function activate(context: vscode.ExtensionContext) {
    let disposable = vscode.commands.registerCommand('extension.chatgpt', async () => {
        const panel = vscode.window.createWebviewPanel(
            'chatGPT',
            'ChatGPT',
            vscode.ViewColumn.One,
            {}
        );

        panel.webview.html = getWebviewContent();

        // Handle messages from the webview
        panel.webview.onDidReceiveMessage(
            async (message) => {
                switch (message.command) {
                    case 'prompt':
                        const response = await getChatGptResponse(message.text);
                        panel.webview.postMessage({ command: 'response', text: response });
                        return;
                }
            },
            undefined,
            context.subscriptions
        );
    });

    context.subscriptions.push(disposable);
}

function getWebviewContent() {
    return `<!DOCTYPE html>
            <html lang="en">
            <head>
                <meta charset="UTF-8">
                <meta name="viewport" content="width=device-width, initial-scale=1.0">
                <title>ChatGPT</title>
            </head>
            <body>
                <input id="input" type="text" size="50">
                <button id="sendButton">Send</button>
                <div id="response"></div>
                
                <script>
                    const vscode = acquireVsCodeApi();
                    document.getElementById('sendButton').addEventListener('click', () => {
                        const input = document.getElementById('input').value;
                        vscode.postMessage({
                            command: 'prompt',
                            text: input
                        });
                    });

                    window.addEventListener('message', event => {
                        const message = event.data;
                        if (message.command === 'response') {
                            document.getElementById('response').innerText = message.text;
                        }
                    });
                </script>
            </body>
            </html>`;
}

async function getChatGptResponse(prompt: string): Promise<string> {
    return '';
}

export function deactivate() {}
