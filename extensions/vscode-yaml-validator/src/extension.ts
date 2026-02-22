import { workspace } from 'vscode';
import {
  LanguageClient,
  LanguageClientOptions,
  ServerOptions,
  TransportKind,
} from 'vscode-languageclient/node';

let client: LanguageClient | undefined;

export function activate(): void {
  const config = workspace.getConfiguration('yamlValidator');
  const serverPath = config.get<string>('serverPath', 'yaml-validator');

  const serverOptions: ServerOptions = {
    command: serverPath,
    args: ['lsp'],
    transport: TransportKind.stdio,
  };

  const clientOptions: LanguageClientOptions = {
    documentSelector: [
      { scheme: 'file', language: 'yaml' },
      { scheme: 'file', pattern: '**/*.yml' },
    ],
  };

  client = new LanguageClient(
    'yaml-validator',
    'YAML Validator',
    serverOptions,
    clientOptions
  );
  client.start();
}

export function deactivate(): Promise<void> | undefined {
  return client?.stop();
}
