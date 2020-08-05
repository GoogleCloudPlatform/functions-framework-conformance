const core = require('@actions/core');
const github = require('@actions/github');
const child_process = require("child_process");
const  fs = require('fs')

try {
  // `who-to-greet` input defined in action metadata file
  const cmd = core.getInput('cmd');
  const functionType = core.getInput('functionType');
  const validateMapping = core.getInput('validateMapping');
  // const runInContainer = core.getInput('runInContainer');

  // Install conformance client binary.
  run("go install github.com/GoogleCloudPlatform/functions-framework-conformance/client");

  // Run the client with the specified parameters.
  run("go run github.com/GoogleCloudPlatform/functions-framework-conformance/client --cmd='" + cmd + "' --type=" + functionType + " --validate-mapping=" + validateMapping);

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
          fs.readFile('serverlog_stdout.txt', 'utf8', function (err,data) {
            if (err) {
              throw err;
            }
            console.log("server stdout: " + data);
          });

          fs.readFile('serverlog_stderr.txt', 'utf8', function (err,data) {
            if (err) {
              throw err;
            }
            console.log("server stderr: "+ data);
          });

          fs.readFile('function_output.txt', 'utf8', function (err,data) {
            if (err) {
              throw err;
            }
            console.log("function output: " + data);
          });

          throw error;
      }
      console.log(`stdout: ${stdout}`);
  });
}
