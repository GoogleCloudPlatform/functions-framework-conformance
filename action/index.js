import {} from '@actions/core';
import {} from '@actions/github';
import {} from 'child_process';
import {} from 'fs';

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
  const cmd = core.getInput('cmd');
  const functionType = core.getInput('functionType');
  const validateMapping = core.getInput('validateMapping');
  // const runInContainer = core.getInput('runInContainer');

  // Install conformance client binary.
  run('go install github.com/GoogleCloudPlatform/functions-framework-conformance/client');

  // Run the client with the specified parameters.
  run([
    `go run github.com/GoogleCloudPlatform/functions-framework-conformance/client`,
    `--cmd=${cmd}`,
    `--type=${functionType}`,
    `--validate-mappings=${validateMapping}`,
  ].join(' '))

} catch (error) {
  core.setFailed(error.message);
}
