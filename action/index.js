const core = require('@actions/core');
const github = require('@actions/github');
const child_process = require('child_process');
const fs = require('fs');

/**
 * Dump contents of file to console.
 * @param {string} f - file to dump
 */
function dump(f) {
    if (!fs.existsSync(f)) {
      return;
    }
      fs.readFileSync(f, 'utf8', (err, data) => {
      if (err) {
        console.log(`error reading ${f}: ${err}`);
      } else {
        console.log(`${f}: ${data}`);
      }
    });
}

/**
 * Run a specified command.
 * @param {string} cmd - command to run
 */
function run(cmd) {
  child_process.exec(cmd, (error, stdout, stderr) => {
    console.log(`stdout: ${stdout}`);
    if (stderr) {
      console.log(`stderr: ${stderr}`);
    }

    dump('serverlog_stdout.txt');
    dump('serverlog_stdout.txt');
    dump('function_output.json');

    if (error) {
      throw error;
    }
  });
}

try {
  const cmd = core.getInput('cmd');
  const functionType = core.getInput('functionType');
  const validateMapping = core.getInput('validateMapping');
  // const runInContainer = core.getInput('runInContainer');

  // Install conformance client binary.
  run('go install github.com/GoogleCloudPlatform/functions-framework-conformance/client');

  // Run the client with the specified parameters.
  run('go run github.com/GoogleCloudPlatform/functions-framework-conformance/client --cmd=\'' +
      cmd + '\' --type=' + functionType +
      ' --validate-mapping=' + validateMapping);

} catch (error) {
  core.setFailed(error.message);
}
