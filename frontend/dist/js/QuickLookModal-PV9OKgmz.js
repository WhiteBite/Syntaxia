import { e as defineComponent, c as computed, r as ref, w as watch, J as createBlock, z as openBlock, B as createVNode, G as withCtx, y as createElementBlock, A as createCommentVNode, O as withKeys, K as withModifiers, C as createBaseVNode, D as toDisplayString, u as unref, M as createTextVNode, E as Transition, T as Teleport } from "./vue-vendor-DLLjc1XI.js";
import { c as useI18n, g as useProjectStore, h as apiService, i as getFileIcon, f as _export_sfc } from "./index-CCibvNLt.js";
import { H as HighlightJS, t as typescript, j as javascript, x as xml, g as go, p as python, r as rust, a as java, c as cpp, b as c, d as csharp, e as ruby, f as php, s as swift, k as kotlin, h as scala, i as bash, y as yaml, l as json, m as css, n as scss, o as sql, q as markdown, u as plaintext } from "./highlighter-D41BBaGQ.js";
import "./ui-vendor-2hF-XK7S.js";
import "./utils-DVYXdora.js";
import "./icons-BGgRd9bo.js";
HighlightJS.registerLanguage("typescript", typescript);
HighlightJS.registerLanguage("javascript", javascript);
HighlightJS.registerLanguage("xml", xml);
HighlightJS.registerLanguage("go", go);
HighlightJS.registerLanguage("python", python);
HighlightJS.registerLanguage("rust", rust);
HighlightJS.registerLanguage("java", java);
HighlightJS.registerLanguage("cpp", cpp);
HighlightJS.registerLanguage("c", c);
HighlightJS.registerLanguage("csharp", csharp);
HighlightJS.registerLanguage("ruby", ruby);
HighlightJS.registerLanguage("php", php);
HighlightJS.registerLanguage("swift", swift);
HighlightJS.registerLanguage("kotlin", kotlin);
HighlightJS.registerLanguage("scala", scala);
HighlightJS.registerLanguage("bash", bash);
HighlightJS.registerLanguage("yaml", yaml);
HighlightJS.registerLanguage("json", json);
HighlightJS.registerLanguage("css", css);
HighlightJS.registerLanguage("scss", scss);
HighlightJS.registerLanguage("sql", sql);
HighlightJS.registerLanguage("markdown", markdown);
HighlightJS.registerLanguage("plaintext", plaintext);
const LANG_MAP = {
  ts: "typescript",
  tsx: "typescript",
  js: "javascript",
  jsx: "javascript",
  vue: "xml",
  go: "go",
  py: "python",
  rs: "rust",
  java: "java",
  cpp: "cpp",
  c: "c",
  cs: "csharp",
  rb: "ruby",
  php: "php",
  swift: "swift",
  kt: "kotlin",
  scala: "scala",
  sh: "bash",
  yaml: "yaml",
  yml: "yaml",
  json: "json",
  xml: "xml",
  html: "xml",
  css: "css",
  scss: "scss",
  sql: "sql",
  md: "markdown"
};
function highlight(code, extension) {
  const lang = LANG_MAP[extension] || "plaintext";
  const result = HighlightJS.highlight(code, { language: lang, ignoreIllegals: true });
  return result.value;
}
const _hoisted_1 = { class: "quicklook-modal" };
const _hoisted_2 = { class: "quicklook-header" };
const _hoisted_3 = { class: "flex items-center gap-2 min-w-0" };
const _hoisted_4 = { class: "text-lg" };
const _hoisted_5 = { class: "quicklook-filename" };
const _hoisted_6 = { class: "quicklook-path" };
const _hoisted_7 = { class: "flex items-center gap-2" };
const _hoisted_8 = {
  key: 0,
  class: "chip-unified chip-unified-accent"
};
const _hoisted_9 = ["aria-label"];
const _hoisted_10 = { class: "quicklook-content scrollable-y" };
const _hoisted_11 = {
  key: 0,
  class: "quicklook-loading"
};
const _hoisted_12 = {
  key: 1,
  class: "quicklook-error"
};
const _hoisted_13 = {
  key: 2,
  class: "quicklook-binary"
};
const _hoisted_14 = { class: "text-lg font-medium text-gray-300" };
const _hoisted_15 = { class: "text-sm text-gray-400" };
const _hoisted_16 = {
  key: 3,
  class: "quicklook-code"
};
const _hoisted_17 = ["innerHTML"];
const _hoisted_18 = { class: "quicklook-footer" };
const _hoisted_19 = ["disabled"];
const _hoisted_20 = ["disabled"];
const _sfc_main = /* @__PURE__ */ defineComponent({
  __name: "QuickLookModal",
  props: {
    modelValue: { type: Boolean },
    filePath: {}
  },
  emits: ["update:modelValue", "add-to-context"],
  setup(__props, { emit: __emit }) {
    const props = __props;
    const emit = __emit;
    const { t } = useI18n();
    const projectStore = useProjectStore();
    const isOpen = computed({
      get: () => props.modelValue,
      set: (value) => emit("update:modelValue", value)
    });
    const content = ref("");
    const isLoading = ref(false);
    const error = ref(null);
    const fileSize = ref(0);
    const fileName = computed(() => {
      const parts = props.filePath.split("/");
      return parts[parts.length - 1];
    });
    const fileExtension = computed(() => {
      const parts = fileName.value.split(".");
      return parts.length > 1 ? parts[parts.length - 1].toLowerCase() : "";
    });
    const isBinary = computed(() => {
      const binaryExtensions = ["png", "jpg", "jpeg", "gif", "webp", "ico", "svg", "pdf", "zip", "tar", "gz", "exe", "dll", "so", "dylib", "woff", "woff2", "ttf", "eot", "mp3", "mp4", "wav", "avi", "mov"];
      return binaryExtensions.includes(fileExtension.value);
    });
    const highlightedContent = computed(() => {
      if (!content.value || isBinary.value) return "";
      try {
        return highlight(content.value, fileExtension.value);
      } catch {
        return escapeHtml(content.value);
      }
    });
    function escapeHtml(text) {
      return text.replace(/&/g, "&amp;").replace(/</g, "&lt;").replace(/>/g, "&gt;");
    }
    watch(() => [props.modelValue, props.filePath], async ([open, path]) => {
      if (open && path) {
        await loadContent();
      }
    }, { immediate: true });
    async function loadContent() {
      if (!props.filePath || !projectStore.currentPath) return;
      isLoading.value = true;
      error.value = null;
      content.value = "";
      try {
        if (isBinary.value) {
          isLoading.value = false;
          return;
        }
        const result = await apiService.readFileContent(projectStore.currentPath, props.filePath);
        content.value = result;
        fileSize.value = result.length;
      } catch (err) {
        error.value = err instanceof Error ? err.message : t("error.loadFailed");
      } finally {
        isLoading.value = false;
      }
    }
    function close() {
      isOpen.value = false;
    }
    async function copyContent() {
      if (content.value) {
        await navigator.clipboard.writeText(content.value);
      }
    }
    function addToContext() {
      emit("add-to-context", props.filePath);
      close();
    }
    function formatSize(bytes) {
      if (bytes < 1024) return `${bytes} B`;
      if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)} KB`;
      return `${(bytes / (1024 * 1024)).toFixed(1)} MB`;
    }
    return (_ctx, _cache) => {
      return openBlock(), createBlock(Teleport, { to: "body" }, [
        createVNode(Transition, { name: "modal" }, {
          default: withCtx(() => [
            isOpen.value ? (openBlock(), createElementBlock("div", {
              key: 0,
              class: "quicklook-overlay",
              onClick: withModifiers(close, ["self"]),
              onKeydown: withKeys(close, ["escape"])
            }, [
              createBaseVNode("div", _hoisted_1, [
                createBaseVNode("div", _hoisted_2, [
                  createBaseVNode("div", _hoisted_3, [
                    createBaseVNode("span", _hoisted_4, toDisplayString(unref(getFileIcon)(fileName.value)), 1),
                    createBaseVNode("span", _hoisted_5, toDisplayString(fileName.value), 1),
                    createBaseVNode("span", _hoisted_6, toDisplayString(__props.filePath), 1)
                  ]),
                  createBaseVNode("div", _hoisted_7, [
                    fileSize.value ? (openBlock(), createElementBlock("span", _hoisted_8, toDisplayString(formatSize(fileSize.value)), 1)) : createCommentVNode("", true),
                    createBaseVNode("button", {
                      onClick: close,
                      class: "icon-btn icon-btn-danger",
                      "aria-label": unref(t)("common.close")
                    }, [..._cache[0] || (_cache[0] = [
                      createBaseVNode("svg", {
                        class: "w-5 h-5",
                        fill: "none",
                        stroke: "currentColor",
                        viewBox: "0 0 24 24",
                        "aria-hidden": "true"
                      }, [
                        createBaseVNode("path", {
                          "stroke-linecap": "round",
                          "stroke-linejoin": "round",
                          "stroke-width": "2",
                          d: "M6 18L18 6M6 6l12 12"
                        })
                      ], -1)
                    ])], 8, _hoisted_9)
                  ])
                ]),
                createBaseVNode("div", _hoisted_10, [
                  isLoading.value ? (openBlock(), createElementBlock("div", _hoisted_11, [
                    _cache[1] || (_cache[1] = createBaseVNode("svg", {
                      class: "loading-spinner",
                      xmlns: "http://www.w3.org/2000/svg",
                      fill: "none",
                      viewBox: "0 0 24 24"
                    }, [
                      createBaseVNode("circle", {
                        class: "opacity-25",
                        cx: "12",
                        cy: "12",
                        r: "10",
                        stroke: "currentColor",
                        "stroke-width": "4"
                      }),
                      createBaseVNode("path", {
                        class: "opacity-75",
                        fill: "currentColor",
                        d: "M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z"
                      })
                    ], -1)),
                    createBaseVNode("p", null, toDisplayString(unref(t)("files.loading")), 1)
                  ])) : error.value ? (openBlock(), createElementBlock("div", _hoisted_12, [
                    _cache[2] || (_cache[2] = createBaseVNode("svg", {
                      class: "w-12 h-12 text-red-400 mb-2",
                      fill: "none",
                      stroke: "currentColor",
                      viewBox: "0 0 24 24"
                    }, [
                      createBaseVNode("path", {
                        "stroke-linecap": "round",
                        "stroke-linejoin": "round",
                        "stroke-width": "2",
                        d: "M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
                      })
                    ], -1)),
                    createBaseVNode("p", null, toDisplayString(error.value), 1)
                  ])) : isBinary.value ? (openBlock(), createElementBlock("div", _hoisted_13, [
                    _cache[3] || (_cache[3] = createBaseVNode("svg", {
                      class: "w-16 h-16 text-gray-400 mb-4",
                      fill: "none",
                      stroke: "currentColor",
                      viewBox: "0 0 24 24"
                    }, [
                      createBaseVNode("path", {
                        "stroke-linecap": "round",
                        "stroke-linejoin": "round",
                        "stroke-width": "2",
                        d: "M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"
                      })
                    ], -1)),
                    createBaseVNode("p", _hoisted_14, toDisplayString(unref(t)("files.binaryFile")), 1),
                    createBaseVNode("p", _hoisted_15, toDisplayString(unref(t)("files.cannotPreview")), 1)
                  ])) : (openBlock(), createElementBlock("pre", _hoisted_16, [
                    createBaseVNode("code", { innerHTML: highlightedContent.value }, null, 8, _hoisted_17)
                  ]))
                ]),
                createBaseVNode("div", _hoisted_18, [
                  createBaseVNode("button", {
                    onClick: copyContent,
                    class: "btn-unified btn-unified-secondary",
                    disabled: !content.value
                  }, [
                    _cache[4] || (_cache[4] = createBaseVNode("svg", {
                      class: "w-4 h-4",
                      fill: "none",
                      stroke: "currentColor",
                      viewBox: "0 0 24 24"
                    }, [
                      createBaseVNode("path", {
                        "stroke-linecap": "round",
                        "stroke-linejoin": "round",
                        "stroke-width": "2",
                        d: "M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z"
                      })
                    ], -1)),
                    createTextVNode(" " + toDisplayString(unref(t)("context.copy")), 1)
                  ], 8, _hoisted_19),
                  createBaseVNode("button", {
                    onClick: addToContext,
                    class: "btn-unified btn-unified-primary",
                    disabled: !content.value
                  }, [
                    _cache[5] || (_cache[5] = createBaseVNode("svg", {
                      class: "w-4 h-4",
                      fill: "none",
                      stroke: "currentColor",
                      viewBox: "0 0 24 24"
                    }, [
                      createBaseVNode("path", {
                        "stroke-linecap": "round",
                        "stroke-linejoin": "round",
                        "stroke-width": "2",
                        d: "M12 4v16m8-8H4"
                      })
                    ], -1)),
                    createTextVNode(" " + toDisplayString(unref(t)("files.addToContext")), 1)
                  ], 8, _hoisted_20)
                ])
              ])
            ], 32)) : createCommentVNode("", true)
          ]),
          _: 1
        })
      ]);
    };
  }
});
const QuickLookModal = /* @__PURE__ */ _export_sfc(_sfc_main, [["__scopeId", "data-v-9bbf12d1"]]);
export {
  QuickLookModal as default
};
