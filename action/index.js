const core = require('@actions/core');
const github = require('@actions/github');
const child_process = require('child_process');
const fs = require('fs');

try {
  const cmd = core.getInput('cmd');
  const functionType = core.getInput('functionType');
  const validateMapping = core.getInput('validateMapping');
  // const runInContainer = core.getInput('runInContainer');

  // Install conformance client binary.
  const installErr = run('go install github.com/GoogleCloudPlatform/functions-framework-conformance/client');
  if (installErr) {
    throw installErr;
  }

  // Run the client with the specified parameters.
  const runErr = run('go run github.com/GoogleCloudPlatform/functions-framework-conformance/client --cmd=\'' +
      cmd + '\' --type=' + functionType +
      ' --validate-mapping=' + validateMapping);
  if (runErr) {
    throw runErr;
  }

} catch (error) {
  core.setFailed(error.message);
}

/**
 * Run a specified command.
 * @param {string} cmd - command to run
 */
function run(cmd) {
  child_process.exec(cmd, (error, stdout, stderr) => {
    if (stderr) {
      console.log(`stderr: ${stderr}`);
    }
    if (error) {
      console.log(`error: ${error}`);
      if (fs.existsSync('serverlog_stdout.txt')) {
        fs.readFileSync('serverlog_stdout.txt', 'utf8', (err, data) => {
          if (err) {
            throw err; // print and move on
          }
          console.log(`server stdout: ${data}`);
        });
      }
      if (fs.existsSync('serverlog_stderr.txt')) {
        fs.readFileSync('serverlog_stderr.txt', 'utf8', (err, data) => {
          if (err) {
            throw err;
          }
          console.log(`server stderr: ${data}`);
        });
      }
      if (fs.existsSync('function_output.json')) {
        fs.readFileSync('function_output.json', 'utf8', (err, data) => {
          if (err) {
            throw err;
          }
          console.log(`function output: ${data}`);
        });
      }
      return error;
    }
    console.log(`stdout: ${stdout}`);
  });
}
