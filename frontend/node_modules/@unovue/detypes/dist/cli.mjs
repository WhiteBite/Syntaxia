import {
  removeMagicCommentsFromFile,
  transformFile
} from "./chunk-FDB77I6Q.mjs";

// src/cli.ts
import process from "node:process";

// src/cli-lib.ts
import fs from "node:fs";
import path from "node:path";
import glob from "fast-glob";

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

// src/cli-lib.ts
var { stat, mkdir } = fs.promises;
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
    const files = (await glob(unixify(`${input}/**/*.{ts,tsx,vue}`))).filter(
      (file) => !file.endsWith(".d.ts")
    );
    const dirs = [...new Set(files.map((file) => path.dirname(file)))].sort();
    await mkdir(path.normalize(output), { recursive: true });
    for (const dir of dirs) {
      const outDir = path.join(output, path.relative(input, dir));
      if (outDir === output)
        continue;
      await mkdir(path.normalize(outDir), { recursive: true });
    }
    for (const file of files) {
      const inputDir = path.dirname(path.relative(input, file));
      const outputName = inferName(file, path.join(output, inputDir));
      removeMagic ? await removeMagicCommentsFromFile(
        path.normalize(file),
        path.normalize(outputName)
      ) : await transformFile(
        path.normalize(file),
        path.normalize(outputName),
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
  const outputDir = path.dirname(output);
  if (outputDir)
    await mkdir(path.normalize(outputDir), { recursive: true });
  removeMagic ? await removeMagicCommentsFromFile(
    path.normalize(input),
    path.normalize(output)
  ) : await transformFile(path.normalize(input), path.normalize(output), {
    removeTsComments
  });
  return true;
  function inferName(input2, outputDir2) {
    let output2;
    const { dir, name, ext } = path.parse(input2);
    if (removeMagic)
      output2 = path.join(outputDir2 ?? dir, `${name}${ext}`);
    else if (ext === ".ts")
      output2 = path.join(outputDir2 ?? dir, `${name}.js`);
    else if (ext === ".tsx")
      output2 = path.join(outputDir2 ?? dir, `${name}.jsx`);
    else if (ext === ".vue")
      output2 = path.join(outputDir2 ?? dir, `${name}.vue`);
    else
      throw new Error(`Unknwon file extension ${input2}`);
    return output2;
  }
}
function printUsage() {
  console.error(USAGE);
}
function unixify(name) {
  return name.replaceAll(path.sep, "/");
}

// src/cli.ts
var args = process.argv.slice(2);
cli(...args).then((success) => process.exit(success ? 0 : 1));
