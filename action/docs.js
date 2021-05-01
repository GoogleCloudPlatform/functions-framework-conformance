// Copyright 2020 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Script to generate the README.md docs from action.yml. Run via: 
//   `npm run docs`
const yaml = require('js-yaml');
const fs   = require('fs');

const START_HINT = "<!--BEGIN GENERATED DOCS-->";
const END_HINT = "<!--END GENERATED DOCS-->"

const stringifyDefaultValue = (value) => {
    const strVal = typeof value === "string" ? `'${value}'` : `${value}`;
    return "`" + strVal + "`";
};

const markdownDoc = (input, description, defaultValue) => (`
### \`${input}\`

${description}. Default value: ${stringifyDefaultValue(defaultValue)}.
`);

try {
  const doc = yaml.load(fs.readFileSync('./action.yml', 'utf8'));
  const markdown = Object.keys(doc.inputs)
    .map((input) => markdownDoc(input, doc.inputs[input].description, doc.inputs[input].default))
    .join("\n");

  let readme = fs.readFileSync('README.md', 'utf-8').split(/\r?\n/);
  let newReadme = [];
  let skipLine = false;
  for (const line of readme) {
    if (line.indexOf(START_HINT) > -1) {
        skipLine = true;
        newReadme.push(line);
        newReadme.push(markdown);
    }
    if (line.indexOf(END_HINT) > -1) {
        skipLine = false;
    }
    if (!skipLine) {
        newReadme.push(line);
    }
  }
  fs.writeFileSync('README.md', newReadme.join("\n"));
} catch (e) {
  console.log(e);
}