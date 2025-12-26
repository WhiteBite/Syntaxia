"use strict";
var __create = Object.create;
var __defProp = Object.defineProperty;
var __getOwnPropDesc = Object.getOwnPropertyDescriptor;
var __getOwnPropNames = Object.getOwnPropertyNames;
var __getProtoOf = Object.getPrototypeOf;
var __hasOwnProp = Object.prototype.hasOwnProperty;
var __copyProps = (to, from, except, desc) => {
  if (from && typeof from === "object" || typeof from === "function") {
    for (let key of __getOwnPropNames(from))
      if (!__hasOwnProp.call(to, key) && key !== except)
        __defProp(to, key, { get: () => from[key], enumerable: !(desc = __getOwnPropDesc(from, key)) || desc.enumerable });
  }
  return to;
};
var __toESM = (mod, isNodeMode, target) => (target = mod != null ? __create(__getProtoOf(mod)) : {}, __copyProps(
  // If the importer is in node compatibility mode or this is not an ESM
  // file that has been converted to a CommonJS file using a Babel-
  // compatible transform (i.e. "__esModule" has not been set), then set
  // "default" to the CommonJS "module.exports" for node compatibility.
  isNodeMode || !mod || !mod.__esModule ? __defProp(target, "default", { value: mod, enumerable: true }) : target,
  mod
));

// src/cli.ts
var import_node_process = __toESM(require("process"));

// src/cli-lib.ts
var import_node_fs3 = __toESM(require("fs"));
var import_node_path2 = __toESM(require("path"));
var import_fast_glob = __toESM(require("fast-glob"));

// package.json
var package_default = {
  name: "@unovue/detypes",
  version: "0.8.5",
  packageManager: "pnpm@9.5.0",
  description: "Removes TypeScript type annotations but keeps the formatting",
  author: "Zernonia <zernonia@gmail.com>",
  license: "MIT",
  homepage: "https://github.com/unovue/detypes#readme",
  repository: {
    type: "git",
    url: "git+https://github.com/unovue/detypes.git"
  },
  bugs: {
    url: "https://github.com/unovue/detypes/issues"
  },
  keywords: [
    "typescript",
    "formatting",
    "vue",
    "sfc"
  ],
  exports: {
    ".": {
      types: "./dist/index.d.ts",
      import: "./dist/index.mjs",
      require: "./dist/index.js"
    }
  },
  main: "dist/index.js",
  bin: "detype.js",
  files: [
    "dist/**/*",
    "index.d.ts"
  ],
  engines: {
    node: ">=18"
  },
  scripts: {
    prepack: "rimraf dist && pnpm build",
    build: "tsup",
    dev: "tsup --watch",
    test: "vitest",
    typecheck: "tsc -p tsconfig.json --noEmit",
    lint: "eslint .",
    "lint:fix": "eslint . --fix",
    release: "bumpp package.json"
  },
  dependencies: {
    "@babel/core": "^7.24.5",
    "@babel/preset-typescript": "^7.24.1",
    "@vue/compiler-dom": "^3.4.27",
    "@vue/compiler-sfc": "^3.4.27",
    "@vuedx/template-ast-types": "0.7.1",
    "fast-glob": "^3.3.2",
    prettier: "^3.2.5",
    typescript: "^5.4.5"
  },
  devDependencies: {
    "@antfu/eslint-config": "^2.17.0",
    "@babel/traverse": "^7.24.5",
    "@types/babel__core": "^7.20.5",
    "@types/babel__traverse": "^7.20.5",
    "@types/node": "18.15.11",
    bumpp: "^9.4.1",
    eslint: "^9.2.0",
    "npm-run-all": "^4.1.5",
    rimraf: "^5.0.7",
    tsup: "^8.0.2",
    vitest: "^1.6.0"
  }
};

// src/transformFile.ts
var import_node_fs2 = __toESM(require("fs"));
var import_prettier2 = require("prettier");

// src/transform.ts
var import_node_fs = require("fs");
var import_node_path = require("path");
var import_core = require("@babel/core");
var import_prettier = require("prettier");
var import_compiler_sfc = require("@vue/compiler-sfc");
var import_template_ast_types = require("@vuedx/template-ast-types");
var import_preset_typescript = __toESM(require("@babel/preset-typescript"));
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
  return await (0, import_prettier.format)(code, {
    ...prettierOptions,
    filepath: originalFileName
  });
}
async function transformVue(code, fileName, options = {}) {
  const parsedVue = (0, import_compiler_sfc.parse)(code, { filename: fileName });
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
    (0, import_compiler_sfc.registerTS)(() => typescript.default);
    const resolveFile = (file) => {
      return (0, import_node_path.resolve)(
        (0, import_node_path.isAbsolute)(file) || file.startsWith("node_modules") ? "" : (0, import_node_path.dirname)(fileName),
        file
      );
    };
    const { content } = (0, import_compiler_sfc.compileScript)(parsedVue.descriptor, {
      id: fileName,
      fs: {
        fileExists(file) {
          const resolvedFile = resolveFile(file);
          if (!(0, import_node_fs.existsSync)(resolvedFile))
            return false;
          return !!(0, import_node_fs.statSync)(resolveFile(file), {
            throwIfNoEntry: false
          })?.isFile();
        },
        readFile: (file) => {
          return (0, import_node_fs.readFileSync)(resolveFile(file)).toString();
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
          CallExpression(path2) {
            const callee = path2.get("callee");
            if (callee.isIdentifier()) {
              if (callee.node.name === "defineProps") {
                const parentPath = path2.parentPath;
                if (parentPath.isCallExpression()) {
                  const callee2 = parentPath.get("callee");
                  if (callee2.isIdentifier() && callee2.node.name === "withDefaults") {
                    parentPath.replaceWith(path2.node);
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
  const ms = new import_compiler_sfc.MagicString(code);
  const locs = [];
  const expressionCodeList = [];
  const template = parsedVue.descriptor.template;
  let expressionCode = "";
  if (template?.ast) {
    (0, import_template_ast_types.traverse)(template.ast, {
      enter(node) {
        if ((0, import_template_ast_types.isSimpleExpressionNode)(node) && !node.isStatic) {
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
        } else if ((0, import_template_ast_types.isComponentNode)(node)) {
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
    presets: [import_preset_typescript.default],
    generatorOpts: {
      shouldPrintComment: (comment) => comment !== "@detype: remove-me" && (!options.removeTsComments || !comment.match(/^\s*(@ts-ignore|@ts-expect-error)/))
    }
  };
  if (options.customizeBabelConfig)
    options.customizeBabelConfig(babelConfig);
  const babelOutput = await (0, import_core.transformAsync)(code, babelConfig);
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
  code = await (0, import_prettier.format)(code, {
    ...prettierOptions,
    filepath: fileName
  });
  return code;
}

// src/transformFile.ts
var { readFile, writeFile } = import_node_fs2.default.promises;
async function transformFile(inputFileName, outputFileName, options = {}) {
  const code = await readFile(inputFileName, "utf-8");
  const prettierOptions = await (0, import_prettier2.resolveConfig)(inputFileName);
  const output = await transform(code, inputFileName, {
    prettierOptions,
    ...options
  });
  await writeFile(outputFileName, output, "utf-8");
}
async function removeMagicCommentsFromFile(inputFileName, outputFileName) {
  const code = await readFile(inputFileName, "utf-8");
  const prettierConfig = await (0, import_prettier2.resolveConfig)(inputFileName);
  const output = await removeMagicComments(code, inputFileName, prettierConfig);
  await writeFile(outputFileName, output, "utf-8");
}

// src/cli-lib.ts
var { stat, mkdir } = import_node_fs3.default.promises;
var USAGE = `Usage:

detypes [-m | --remove-magic-comments] <INPUT> [OUTPUT]

INPUT   Input file or directory

OUTPUT  Output file or directory
(optional if it can be inferred and it won't overwrite the source file)

-t, --remove-ts-comments
Remove @ts-ignore and @ts-expect-error comments

-m, --remove-magic-comments
Remove magic comments only, don't perform ts > js transform

detypes [-v | --version]

Print version and exit

detypes [-h | --help]

Print this help and exit`;
var VERSION = package_default.version;
async function cli(...args2) {
  let dashDash = false;
  const params = [];
  const flags = [];
  for (const arg of args2) {
    if (arg === "--")
      dashDash = true;
    else if (dashDash || !arg.startsWith("-"))
      params.push(arg);
    else
      flags.push(arg);
  }
  if (params.length > 2) {
    console.error("Too many arguments");
    return false;
  }
  let [input, output] = params;
  if (params.length === 0 || flags.some((flag) => flag === "-h" || flag === "--help")) {
    printUsage();
    return params.length > 0;
  }
  if (flags.some((flag) => flag === "-v" || flag === "--version")) {
    console.log(VERSION);
    return true;
  }
  const removeMagic = flags.some(
    (flag) => flag === "-m" || flag === "--remove-magic-comments"
  );
  const removeTsComments = flags.some(
    (flag) => flag === "-t" || flag === "--remove-ts-comments"
  );
  if (removeMagic && removeTsComments) {
    console.warn(
      "--remove-ts-comments has no effect when --remove-magic-comments is used"
    );
  }
  if (!removeMagic)
    [input, output] = args2;
  const inputStat = await stat(input);
  if (inputStat.isDirectory()) {
    if (!output) {
      console.error("No output directory given");
      printUsage();
      return false;
    }
    const files = (await (0, import_fast_glob.default)(unixify(`${input}/**/*.{ts,tsx,vue}`))).filter(
      (file) => !file.endsWith(".d.ts")
    );
    const dirs = [...new Set(files.map((file) => import_node_path2.default.dirname(file)))].sort();
    await mkdir(import_node_path2.default.normalize(output), { recursive: true });
    for (const dir of dirs) {
      const outDir = import_node_path2.default.join(output, import_node_path2.default.relative(input, dir));
      if (outDir === output)
        continue;
      await mkdir(import_node_path2.default.normalize(outDir), { recursive: true });
    }
    for (const file of files) {
      const inputDir = import_node_path2.default.dirname(import_node_path2.default.relative(input, file));
      const outputName = inferName(file, import_node_path2.default.join(output, inputDir));
      removeMagic ? await removeMagicCommentsFromFile(
        import_node_path2.default.normalize(file),
        import_node_path2.default.normalize(outputName)
      ) : await transformFile(
        import_node_path2.default.normalize(file),
        import_node_path2.default.normalize(outputName),
        { removeTsComments }
      );
    }
    return true;
  }
  if (output) {
    const outputStat = await stat(output).catch((error) => {
      if (error && error.code === "ENOENT")
        return null;
      throw error;
    });
    if (outputStat && outputStat.isDirectory())
      output = inferName(input, output);
  } else {
    if (removeMagic) {
      console.error(
        "Output file name is required when removing magic comments"
      );
      return false;
    }
    if (input.endsWith(".vue")) {
      console.error("Output file name is required for .vue files");
      return false;
    }
    output = inferName(input);
  }
  const outputDir = import_node_path2.default.dirname(output);
  if (outputDir)
    await mkdir(import_node_path2.default.normalize(outputDir), { recursive: true });
  removeMagic ? await removeMagicCommentsFromFile(
    import_node_path2.default.normalize(input),
    import_node_path2.default.normalize(output)
  ) : await transformFile(import_node_path2.default.normalize(input), import_node_path2.default.normalize(output), {
    removeTsComments
  });
  return true;
  function inferName(input2, outputDir2) {
    let output2;
    const { dir, name, ext } = import_node_path2.default.parse(input2);
    if (removeMagic)
      output2 = import_node_path2.default.join(outputDir2 ?? dir, `${name}${ext}`);
    else if (ext === ".ts")
      output2 = import_node_path2.default.join(outputDir2 ?? dir, `${name}.js`);
    else if (ext === ".tsx")
      output2 = import_node_path2.default.join(outputDir2 ?? dir, `${name}.jsx`);
    else if (ext === ".vue")
      output2 = import_node_path2.default.join(outputDir2 ?? dir, `${name}.vue`);
    else
      throw new Error(`Unknwon file extension ${input2}`);
    return output2;
  }
}
function printUsage() {
  console.error(USAGE);
}
function unixify(name) {
  return name.replaceAll(import_node_path2.default.sep, "/");
}

// src/cli.ts
var args = import_node_process.default.argv.slice(2);
cli(...args).then((success) => import_node_process.default.exit(success ? 0 : 1));
