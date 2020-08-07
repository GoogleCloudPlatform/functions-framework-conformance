const core = require('@actions/core');
const github = require('@actions/github');
const child_process = require('child_process');
const fs = require('fs');

/**
 * writeFileToConsole contents of file to console.
 * @param {string} path - filepath to write to the console
 */
function writeFileToConsole(path) {
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
function run(cmd) {
  try {
    child_process.execSync(cmd);
  } catch (error) {
    writeFileToConsole('serverlog_stdout.txt');
    writeFileToConsole('serverlog_stderr.txt');
    writeFileToConsole('function_output.json');
    throw error;
  }
}

try {
  const functionType = core.getInput('functionType');
  const validateMapping = core.getInput('validateMapping');
  const source = core.getInput('source');
  const target = core.getInput('target');
  const runtime = core.getInput('runtime');
  const tag = core.getInput('tag');

  // Install conformance client binary.
  run('go install github.com/GoogleCloudPlatform/functions-framework-conformance/client');

  // Run the client with the specified parameters.
  run([
    `go run github.com/GoogleCloudPlatform/functions-framework-conformance/client`,
    `-type=${functionType}`,
    `-validate-mapping=${validateMapping}`,
    `-builder-source=${source}`,
    `-builder-target=${target}`,
    `-builder-runtime=${runtime}`,
    `-builder-tag=${tag}`,
  ].join(' '));

} catch (error) {
  core.setFailed(error.message);
}
