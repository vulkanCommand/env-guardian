"use strict";

const childProcess = require("child_process");
const vscode = require("vscode");

const outputChannelName = "Env Guardian";
const installCommand = "curl -fsSL https://raw.githubusercontent.com/vulkanCommand/env-guardian/main/scripts/install.sh | sh";
const githubIssuesUrl = "https://github.com/vulkanCommand/env-guardian/issues";
const ansiPattern = /\x1b\[[0-9;?]*[ -/]*[@-~]/g;

function activate(context) {
	const outputChannel = vscode.window.createOutputChannel(outputChannelName);
	const statusBarItem = vscode.window.createStatusBarItem(vscode.StatusBarAlignment.Left, 100);

	statusBarItem.text = "$(shield) Env Guardian";
	statusBarItem.tooltip = "Run Env Guardian validation";
	statusBarItem.command = "envGuardian.validate";
	statusBarItem.show();

	context.subscriptions.push(
		outputChannel,
		statusBarItem,
		registerCommand("envGuardian.validate", () => runValidate(outputChannel)),
		registerCommand("envGuardian.validateAll", () => runValidateAll(outputChannel)),
		registerCommand("envGuardian.ci", () => runCI(outputChannel)),
		registerCommand("envGuardian.security", () => runSecurity(outputChannel)),
		registerCommand("envGuardian.logScan", () => runLogScan(outputChannel)),
		registerCommand("envGuardian.version", () => runVersion(outputChannel)),
	);
}

function deactivate() {}

function registerCommand(command, callback) {
	return vscode.commands.registerCommand(command, callback);
}

function runValidate(outputChannel) {
	const config = getConfig();
	const args = ["validate", "--file", config.envFile, "--example", config.exampleFile];

	if (config.useJson) {
		args.push("--json");
	}

	runEnvGuard("Validate", args, outputChannel);
}

function runValidateAll(outputChannel) {
	const config = getConfig();
	const args = ["validate", "--all"];

	if (config.useJson) {
		args.push("--json");
	}

	runEnvGuard("Validate All Environments", args, outputChannel);
}

function runCI(outputChannel) {
	const config = getConfig();
	const args = ["ci", "--file", config.envFile, "--example", config.exampleFile];

	if (config.useJson) {
		args.push("--json");
	}

	runEnvGuard("CI Check", args, outputChannel);
}

function runSecurity(outputChannel) {
	const config = getConfig();
	const args = ["security", "--dir", config.rootDirectory, "--file", config.envFile];

	if (config.useJson) {
		args.push("--json");
	}

	runEnvGuard("Security Scan", args, outputChannel);
}

function runLogScan(outputChannel) {
	const config = getConfig();
	const args = ["log-scan", "--dir", config.rootDirectory];

	if (config.useJson) {
		args.push("--json");
	}

	runEnvGuard("Log Exposure Scan", args, outputChannel);
}

function runVersion(outputChannel) {
	runEnvGuard("Version", ["version"], outputChannel);
}

function runEnvGuard(label, args, outputChannel) {
	const workspaceFolder = getWorkspaceFolder();
	if (!workspaceFolder) {
		vscode.window.showErrorMessage("Env Guardian requires an open workspace folder.");
		return;
	}

	const config = getConfig();
	const commandLine = [config.executablePath].concat(args).join(" ");
	let failedToStart = false;

	outputChannel.clear();
	outputChannel.appendLine("Env Guardian");
	outputChannel.appendLine("------------");
	outputChannel.appendLine(`$ ${commandLine}`);
	outputChannel.appendLine("");
	outputChannel.show(true);

	const child = childProcess.spawn(config.executablePath, args, {
		cwd: workspaceFolder.uri.fsPath,
		env: Object.assign({}, process.env, {
			NO_COLOR: "1",
		}),
		windowsHide: true,
	});

	child.stdout.on("data", (chunk) => {
		outputChannel.append(stripAnsi(chunk.toString()));
	});

	child.stderr.on("data", (chunk) => {
		outputChannel.append(stripAnsi(chunk.toString()));
	});

	child.on("error", (error) => {
		failedToStart = true;
		showMissingCLIMessage(config.executablePath, error, outputChannel);
	});

	child.on("close", (code) => {
		if (failedToStart) {
			return;
		}

		if (code === 0) {
			vscode.window.showInformationMessage(`Env Guardian ${label} passed.`);
			return;
		}

		vscode.window.showWarningMessage(`Env Guardian ${label} completed with issues (exit code ${code}).`);
	});
}

function stripAnsi(value) {
	return value.replace(ansiPattern, "");
}

function showMissingCLIMessage(executablePath, error, outputChannel) {
	outputChannel.appendLine("");
	outputChannel.appendLine("Env Guardian CLI could not be started.");
	outputChannel.appendLine(`Configured executable: ${executablePath}`);
	outputChannel.appendLine(`Reason: ${error.message}`);
	outputChannel.appendLine("");
	outputChannel.appendLine("Install Env Guardian:");
	outputChannel.appendLine(`  ${installCommand}`);
	outputChannel.appendLine("");
	outputChannel.appendLine("If Env Guardian is already installed, set envGuardian.executablePath to the full envguard binary path.");

	vscode.window.showErrorMessage(
		"Env Guardian CLI is not installed or VS Code cannot find it.",
		"Copy Install Command",
		"Open Settings",
		"Open GitHub Issues",
	).then((selection) => {
		if (selection === "Copy Install Command") {
			vscode.env.clipboard.writeText(installCommand);
			vscode.window.showInformationMessage("Env Guardian install command copied.");
			return;
		}

		if (selection === "Open Settings") {
			vscode.commands.executeCommand("workbench.action.openSettings", "envGuardian.executablePath");
			return;
		}

		if (selection === "Open GitHub Issues") {
			vscode.env.openExternal(vscode.Uri.parse(githubIssuesUrl));
		}
	});
}

function getWorkspaceFolder() {
	const activeEditor = vscode.window.activeTextEditor;
	if (activeEditor) {
		const activeFolder = vscode.workspace.getWorkspaceFolder(activeEditor.document.uri);
		if (activeFolder) {
			return activeFolder;
		}
	}

	if (vscode.workspace.workspaceFolders && vscode.workspace.workspaceFolders.length > 0) {
		return vscode.workspace.workspaceFolders[0];
	}

	return undefined;
}

function getConfig() {
	const config = vscode.workspace.getConfiguration("envGuardian");

	return {
		executablePath: config.get("executablePath", "envguard"),
		envFile: config.get("envFile", ".env"),
		exampleFile: config.get("exampleFile", ".env.example"),
		rootDirectory: config.get("rootDirectory", "."),
		useJson: config.get("useJson", false),
	};
}

module.exports = {
	activate,
	deactivate,
};
