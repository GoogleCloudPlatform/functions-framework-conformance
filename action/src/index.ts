import * as core from '@actions/core';
import * as childProcess from 'child_process';
import * as fs from 'fs';
import * as process from 'process';

/**
 * Run a specified command.
 * @param {string} cmd - command to run
 */
function runCmd(cmd: string) {
  try {
    console.log(`RUNNING: "${cmd}"`)
    childProcess.execSync(cmd);
  } catch (error) {
    core.setFailed(error.message);
  }
}

async function run() {
  const version = core.getInput('version');
  const outputFile = core.getInput('outputFile');
  const functionType = core.getInput('functionType');
  const declarativeType = core.getInput('declarativeType');
  const validateMapping = core.getInput('validateMapping');
  const validateConcurrency = core.getInput('validateConcurrency');
  const source = core.getInput('source');
  const target = core.getInput('target');
  const runtime = core.getInput('runtime');
  const runtimeVersion = core.getInput('runtimeVersion');
  const tag = core.getInput('tag');
  const useBuildpacks = core.getInput('useBuildpacks');
  const cmd = core.getInput('cmd');
  const startDelay = core.getInput('startDelay');
  const workingDirectory = core.getInput('workingDirectory');
  const runtimeEnvs = core.getInput('runtimeEnvs');

  let cwd = process.cwd();

  // Build conformance client binary from source.
  let repo = 'functions-framework-conformance';
  if (!fs.existsSync(repo)) {
    runCmd(`git clone https://github.com/GoogleCloudPlatform/functions-framework-conformance.git`);
  }
  process.chdir('functions-framework-conformance/client');
  if (version) {
    runCmd(`git fetch origin refs/tags/${version} && git checkout ${version}`);
  } else {
    // Checkout latest release tag.
    runCmd('git fetch --tags && git checkout $(git describe --tags $(git rev-list --tags --max-count=1))')
  }
  runCmd(`go build -o ~/client`);

  process.chdir(cwd);
  if (workingDirectory) {
    process.chdir(workingDirectory);
  }
  // Run the client with the specified parameters.
  runCmd([
    `~/client`,
    `-output-file=${outputFile}`,
    `-type=${functionType}`,
    `-declarative-type=${declarativeType}`,
    `-validate-mapping=${validateMapping}`,
    `-validate-concurrency=${validateConcurrency}`,
    `-builder-source=${source}`,
    `-builder-target=${target}`,
    `-builder-runtime=${runtime}`,
    `-builder-runtime-version=${runtimeVersion}`,
    `-builder-tag=${tag}`,
    `-buildpacks=${useBuildpacks}`,
    `-cmd=${cmd}`,
    `-start-delay=${startDelay}`,
    `-envs=${runtimeEnvs}`,
  ].filter((x) => !!x).join(' '));
}

run();
