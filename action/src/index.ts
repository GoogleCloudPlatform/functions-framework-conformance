import * as core from '@actions/core';
import * as github from '@actions/github';
import * as childProcess from 'child_process';
import * as fs from 'fs';

/**
 * writeFileToConsole contents of file to console.
 * @param {string} path - filepath to write to the console
 */
function writeFileToConsole(path: string) {
  try {
    const data = fs.readFileSync(path, 'utf8');
    console.log(`${path}: ${data}`);
  } catch (e) {
    console.log(`$unable to read {path}, skipping: ${e}`);
  }
}

/**
 * Run a specified command.
 * @param {string} cmd - command to run
 */
function runCmd(cmd: string) {
  try {
    childProcess.execSync(cmd);
  } catch (error) {
    writeFileToConsole('serverlog_stdout.txt');
    writeFileToConsole('serverlog_stderr.txt');
    writeFileToConsole('function_output.json');
    core.setFailed(error.message);
  }
}

async function run() {
  const outputFile = core.getInput('outputFile');
  const functionType = core.getInput('functionType');
  const validateMapping = core.getInput('validateMapping');
  const source = core.getInput('source');
  const target = core.getInput('target');
  const runtime = core.getInput('runtime');
  const tag = core.getInput('tag');
  const useBuildpacks = core.getInput('useBuildpacks');
  const cmd = core.getInput('cmd');
  const startDelay = core.getInput('startDelay');

  // Install conformance client binary.
  runCmd(
      'go get github.com/GoogleCloudPlatform/functions-framework-conformance/client');

  // Run the client with the specified parameters.
  runCmd([
    `go run github.com/GoogleCloudPlatform/functions-framework-conformance/client`,
    `-output-file=${outputFile}`,
    `-type=${functionType}`,
    `-validate-mapping=${validateMapping}`,
    `-builder-source=${source}`,
    `-builder-target=${target}`,
    `-builder-runtime=${runtime}`,
    `-builder-tag=${tag}`,
    `-buildpacks=${useBuildpacks}`,
    `-cmd=${cmd}`,
    `-start-delay=${startDelay}`,
  ].join(' '));
}

run();
