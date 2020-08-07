"use strict";
var __createBinding = (this && this.__createBinding) || (Object.create ? (function(o, m, k, k2) {
    if (k2 === undefined) k2 = k;
    Object.defineProperty(o, k2, { enumerable: true, get: function() { return m[k]; } });
}) : (function(o, m, k, k2) {
    if (k2 === undefined) k2 = k;
    o[k2] = m[k];
}));
var __setModuleDefault = (this && this.__setModuleDefault) || (Object.create ? (function(o, v) {
    Object.defineProperty(o, "default", { enumerable: true, value: v });
}) : function(o, v) {
    o["default"] = v;
});
var __importStar = (this && this.__importStar) || function (mod) {
    if (mod && mod.__esModule) return mod;
    var result = {};
    if (mod != null) for (var k in mod) if (k !== "default" && Object.hasOwnProperty.call(mod, k)) __createBinding(result, mod, k);
    __setModuleDefault(result, mod);
    return result;
};
var __awaiter = (this && this.__awaiter) || function (thisArg, _arguments, P, generator) {
    function adopt(value) { return value instanceof P ? value : new P(function (resolve) { resolve(value); }); }
    return new (P || (P = Promise))(function (resolve, reject) {
        function fulfilled(value) { try { step(generator.next(value)); } catch (e) { reject(e); } }
        function rejected(value) { try { step(generator["throw"](value)); } catch (e) { reject(e); } }
        function step(result) { result.done ? resolve(result.value) : adopt(result.value).then(fulfilled, rejected); }
        step((generator = generator.apply(thisArg, _arguments || [])).next());
    });
};
Object.defineProperty(exports, "__esModule", { value: true });
const core = __importStar(require("@actions/core"));
const childProcess = __importStar(require("child_process"));
const fs = __importStar(require("fs"));
/**
 * writeFileToConsole contents of file to console.
 * @param {string} path - filepath to write to the console
 */
function writeFileToConsole(path) {
    try {
        const data = fs.readFileSync(path, 'utf8');
        console.log(`${path}: ${data}`);
    }
    catch (e) {
        console.log(`$unable to read {path}, skipping: ${e}`);
    }
}
/**
 * Run a specified command.
 * @param {string} cmd - command to run
 */
function runCmd(cmd) {
    try {
        childProcess.execSync(cmd);
    }
    catch (error) {
        writeFileToConsole('serverlog_stdout.txt');
        writeFileToConsole('serverlog_stderr.txt');
        writeFileToConsole('function_output.json');
        throw error;
    }
}
function run() {
    return __awaiter(this, void 0, void 0, function* () {
        try {
            const functionType = core.getInput('functionType');
            const validateMapping = core.getInput('validateMapping');
            const source = core.getInput('source');
            const target = core.getInput('target');
            const runtime = core.getInput('runtime');
            const tag = core.getInput('tag');
            // Install conformance client binary.
            runCmd('go install github.com/GoogleCloudPlatform/functions-framework-conformance/client');
            // Run the client with the specified parameters.
            runCmd([
                `go run github.com/GoogleCloudPlatform/functions-framework-conformance/client`,
                `-type=${functionType}`,
                `-validate-mapping=${validateMapping}`,
                `-builder-source=${source}`,
                `-builder-target=${target}`,
                `-builder-runtime=${runtime}`,
                `-builder-tag=${tag}`,
            ].join(' '));
        }
        catch (error) {
            core.setFailed(error.message);
        }
    });
}
run();
