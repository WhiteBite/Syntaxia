// src/transform.ts
import { existsSync, readFileSync, statSync } from "node:fs";
import { dirname, isAbsolute, resolve } from "node:path";
import {
  transformAsync
} from "@babel/core";
import { format } from "prettier";
import {
  MagicString,
  compileScript,
  parse,
  registerTS
} from "@vue/compiler-sfc";
import {
  isComponentNode as isVueComponentNode,
  isSimpleExpressionNode as isVueSimpleExpressionNode,
  traverse as traverseVueAst
} from "@vuedx/template-ast-types";
import babelTs from "@babel/preset-typescript";
function getDefinePropsObject(content) {
  const matched = /\sprops:\s*\{/m.exec(content);
  if (matched) {
    const startContentIndex = matched.index + matched[0].length - 1;
    let leftBracketCount = 1;
    let endContentIndex = startContentIndex + 1;
    while (leftBracketCount) {
      if (content.charAt(endContentIndex) === "{")
        leftBracketCount++;
      else if (content.charAt(endContentIndex) === "}")
        leftBracketCount--;
      endContentIndex++;
    }
    return content.substring(startContentIndex, endContentIndex);
  }
  return "";
}
async function transform(code, fileName, options = {}) {
  const { prettierOptions, ...removeTypeOptions } = options;
  const originalFileName = fileName;
  code = code.replaceAll("\r\n", "\n");
  if (fileName.endsWith(".vue"))
    code = await transformVue(code, fileName, options);
  else
    code = await removeTypes(code, fileName, removeTypeOptions);
  return await format(code, {
    ...prettierOptions,
    filepath: originalFileName
  });
}
async function transformVue(code, fileName, options = {}) {
  const parsedVue = parse(code, { filename: fileName });
  if (parsedVue.descriptor.script?.lang !== "ts" && parsedVue.descriptor.scriptSetup?.lang !== "ts") {
    return code;
  }
  const { script: script1, scriptSetup: script2 } = parsedVue.descriptor;
  let propsContent = "";
  let emitsContent = "";
  const isContainsDefinePropsType = script2?.content.match(/defineProps\s*</m);
  const isContainsDefineEmitType = script2?.content.match(/defineEmits\s*</m);
  if (isContainsDefinePropsType || isContainsDefineEmitType) {
    const typescript = await import("typescript");
    registerTS(() => typescript.default);
    const resolveFile = (file) => {
      return resolve(
        isAbsolute(file) || file.startsWith("node_modules") ? "" : dirname(fileName),
        file
      );
    };
    const { content } = compileScript(parsedVue.descriptor, {
      id: fileName,
      fs: {
        fileExists(file) {
          const resolvedFile = resolveFile(file);
          if (!existsSync(resolvedFile))
            return false;
          return !!statSync(resolveFile(file), {
            throwIfNoEntry: false
          })?.isFile();
        },
        readFile: (file) => {
          return readFileSync(resolveFile(file)).toString();
        }
      }
    });
    if (isContainsDefinePropsType)
      propsContent = getDefinePropsObject(content);
    if (isContainsDefineEmitType)
      emitsContent = content.match(/\semits:\s(\[.*\]?)/m)?.[1] || "";
  }
  const removeVueSfcScriptOptions = {
    ...options,
    customizeBabelConfig(config) {
      config.plugins ||= [];
      config.plugins?.push({
        name: "detype-remove-with-defaults",
        visitor: {
          CallExpression(path) {
            const callee = path.get("callee");
            if (callee.isIdentifier()) {
              if (callee.node.name === "defineProps") {
                const parentPath = path.parentPath;
                if (parentPath.isCallExpression()) {
                  const callee2 = parentPath.get("callee");
                  if (callee2.isIdentifier() && callee2.node.name === "withDefaults") {
                    parentPath.replaceWith(path.node);
                    parentPath.stop();
                  }
                }
              }
            }
          }
        }
      });
      options.customizeBabelConfig?.(config);
    }
  };
  const ms = new MagicString(code);
  const locs = [];
  const expressionCodeList = [];
  const template = parsedVue.descriptor.template;
  let expressionCode = "";
  if (template?.ast) {
    traverseVueAst(template.ast, {
      enter(node) {
        if (isVueSimpleExpressionNode(node) && !node.isStatic) {
          const ForOfOrInRE = /\s+(of|in)\s+/;
          if (node.content.match(ForOfOrInRE)) {
            const parts = node.content.split(ForOfOrInRE);
            if (parts.length === 3) {
              const content = parts[parts.length - 1];
              expressionCodeList.push(content);
              locs.push([node.loc.start.offset + node.content.length - content.length, node.loc.end.offset]);
            } else {
              expressionCodeList.push(node.content);
              locs.push([node.loc.start.offset, node.loc.end.offset]);
            }
          } else {
            expressionCodeList.push(node.loc.source);
            locs.push([node.loc.start.offset, node.loc.end.offset]);
          }
        } else if (isVueComponentNode(node)) {
          const content = node.tag;
          expressionCodeList.push(node.tag);
          let start = node.loc.start.offset + 1;
          locs.push([start, start += content.length]);
        }
      }
    });
    const delimiter = `['---detypes-delimiter---'];`;
    expressionCode = await removeTypes(expressionCodeList.map((c) => `[${c}]`).join(`;${delimiter}`), `${fileName}.ts`, options);
    const lines = expressionCode.split(delimiter);
    for (let i = 0; i < locs.length; i++) {
      const loc = locs[i];
      const line = lines[i].trim();
      ms.update(loc[0], loc[1], line.substring(1, line.length - 2));
    }
  }
  for (const script of [script1, script2].filter(Boolean)) {
    code = await removeTypesFromVueSfcScript(
      expressionCode,
      fileName,
      script,
      removeVueSfcScriptOptions
    );
    if (script?.attrs.generic)
      ms.replace(` generic="${script.attrs.generic}"`, "");
    ms.update(script.loc.start.offset, script.loc.end.offset, code);
  }
  ms.replaceAll(/\blang\s*=\s*["']ts["']/g, "");
  let result = ms.toString();
  if (propsContent)
    result = result.replace("defineProps(", (str) => `${str}${propsContent}`);
  if (emitsContent)
    result = result.replace("defineEmits(", (str) => `${str}${emitsContent}`);
  return result;
}
async function removeTypes(code, fileName, options) {
  code = code.replace(
    /\n\n+/g,
    (match) => `
/* @detype: empty-line=${match.length} */
`
  );
  code = processMagicComments(code);
  const removeComments = {
    enter(p) {
      if (!p.node.leadingComments)
        return;
      for (let i = p.node.leadingComments.length - 1; i >= 0; i--) {
        const comment = p.node.leadingComments[i];
        if (code.slice(comment.end).match(/^\s*\n\s*\n/) || comment.value.includes("@detype: empty-line")) {
          break;
        }
        comment.value = "@detype: remove-me";
      }
    }
  };
  const babelConfig = {
    filename: fileName,
    retainLines: true,
    plugins: [
      // Plugin to remove leading comments attached to TypeScript-only constructs
      {
        name: "detype-comment-remover",
        visitor: {
          TSTypeAliasDeclaration: removeComments,
          TSInterfaceDeclaration: removeComments,
          TSDeclareFunction: removeComments,
          TSDeclareMethod: removeComments,
          TSImportType: removeComments
        }
      }
    ].filter(Boolean),
    presets: [babelTs],
    generatorOpts: {
      shouldPrintComment: (comment) => comment !== "@detype: remove-me" && (!options.removeTsComments || !comment.match(/^\s*(@ts-ignore|@ts-expect-error)/))
    }
  };
  if (options.customizeBabelConfig)
    options.customizeBabelConfig(babelConfig);
  const babelOutput = await transformAsync(code, babelConfig);
  if (!babelOutput || babelOutput.code === void 0 || babelOutput.code === null)
    throw new Error("Babel error");
  return babelOutput.code.replaceAll(/\n\n*/g, "\n").replace(/\/\* @detype: empty-line=([0-9]+) \*\//g, (_match, p1) => `
`.repeat(p1 - 2));
}
async function removeTypesFromVueSfcScript(expressionCode, fileName, script, options) {
  if (script === null || script.lang !== "ts")
    return script.content;
  script.content += `/* @detype: remove-after-this */${expressionCode}`;
  let scriptCode = await removeTypes(script.content, `${fileName}.ts`, options);
  const removeAfterIndex = scriptCode.indexOf(
    "/* @detype: remove-after-this */"
  );
  if (removeAfterIndex >= 0)
    scriptCode = scriptCode.slice(0, removeAfterIndex);
  return scriptCode;
}
function processMagicComments(input) {
  const REPLACE_COMMENT = "// @detype: replace\n";
  const WITH_COMMENT = "// @detype: with\n";
  const END_COMMENT = "// @detype: end\n";
  let start = input.indexOf(REPLACE_COMMENT);
  while (start >= 0) {
    const middle = input.indexOf(WITH_COMMENT, start);
    if (middle < 0)
      return input;
    const middleEnd = middle + WITH_COMMENT.length;
    const end = input.indexOf(END_COMMENT, middleEnd);
    if (end < 0)
      return input;
    const endEnd = end + END_COMMENT.length;
    const before = input.slice(0, start);
    const newText = input.slice(middleEnd, end).replaceAll(/^\s*\/\//gm, "");
    const after = input.slice(endEnd);
    input = before + newText + after;
    start = input.indexOf(REPLACE_COMMENT, before.length + newText.length);
  }
  return input;
}
async function removeMagicComments(code, fileName, prettierOptions) {
  const REPLACE_COMMENT = "// @detype: replace\n";
  const WITH_COMMENT = "// @detype: with\n";
  const END_COMMENT = "// @detype: end\n";
  let start = code.indexOf(REPLACE_COMMENT);
  let startEnd = start + REPLACE_COMMENT.length;
  while (start >= 0) {
    const middle = code.indexOf(WITH_COMMENT, start);
    if (middle < 0)
      return code;
    const middleEnd = middle + WITH_COMMENT.length;
    const end = code.indexOf(END_COMMENT, middleEnd);
    if (end < 0)
      return code;
    const endEnd = end + END_COMMENT.length;
    const before = code.slice(0, start);
    const keptText = code.slice(startEnd, middle);
    const after = code.slice(endEnd);
    code = before + keptText + after;
    start = code.indexOf(REPLACE_COMMENT, before.length + keptText.length);
    startEnd = start + REPLACE_COMMENT.length;
  }
  code = await format(code, {
    ...prettierOptions,
    filepath: fileName
  });
  return code;
}

// src/transformFile.ts
import fs from "node:fs";
import { resolveConfig } from "prettier";
var { readFile, writeFile } = fs.promises;
async function transformFile(inputFileName, outputFileName, options = {}) {
  const code = await readFile(inputFileName, "utf-8");
  const prettierOptions = await resolveConfig(inputFileName);
  const output = await transform(code, inputFileName, {
    prettierOptions,
    ...options
  });
  await writeFile(outputFileName, output, "utf-8");
}
async function removeMagicCommentsFromFile(inputFileName, outputFileName) {
  const code = await readFile(inputFileName, "utf-8");
  const prettierConfig = await resolveConfig(inputFileName);
  const output = await removeMagicComments(code, inputFileName, prettierConfig);
  await writeFile(outputFileName, output, "utf-8");
}

export {
  transform,
  removeMagicComments,
  transformFile,
  removeMagicCommentsFromFile
};
