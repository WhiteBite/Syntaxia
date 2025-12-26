import { e as defineComponent, r as ref, c as computed, y as createElementBlock, z as openBlock, C as createBaseVNode, A as createCommentVNode, D as toDisplayString, u as unref, l as normalizeClass, F as Fragment, H as renderList, w as watch, J as createBlock, B as createVNode, G as withCtx, K as withModifiers, E as Transition, M as createTextVNode, N as withDirectives, P as vModelText, T as Teleport } from "./vue-vendor-DLLjc1XI.js";
import { c as useI18n, f as _export_sfc, d as useSettingsStore, g as useProjectStore, u as useUIStore, p as parseIgnoreRules, h as apiService } from "./index-CCibvNLt.js";
import "./ui-vendor-2hF-XK7S.js";
import "./utils-DVYXdora.js";
import "./icons-BGgRd9bo.js";
import "./highlighter-D41BBaGQ.js";
const _hoisted_1$1 = { class: "preview-panel" };
const _hoisted_2$1 = { class: "preview-panel__header" };
const _hoisted_3$1 = { class: "preview-panel__title" };
const _hoisted_4$1 = {
  key: 0,
  class: "preview-panel__count"
};
const _hoisted_5$1 = {
  key: 0,
  class: "preview-panel__tabs"
};
const _hoisted_6$1 = {
  key: 0,
  class: "preview-panel__loading"
};
const _hoisted_7$1 = {
  key: 1,
  class: "preview-panel__empty"
};
const _hoisted_8$1 = {
  key: 2,
  class: "preview-panel__content"
};
const _hoisted_9$1 = {
  key: 0,
  class: "preview-list"
};
const _hoisted_10$1 = { class: "preview-dir-item__name" };
const _hoisted_11$1 = { class: "preview-dir-item__count" };
const _hoisted_12$1 = {
  key: 1,
  class: "preview-list"
};
const _hoisted_13$1 = { class: "preview-rule-item__pattern" };
const _hoisted_14$1 = { class: "preview-rule-item__count" };
const _hoisted_15$1 = {
  key: 2,
  class: "preview-list"
};
const _hoisted_16$1 = { class: "preview-file-item__path" };
const _hoisted_17$1 = {
  key: 0,
  class: "preview-more"
};
const _sfc_main$1 = /* @__PURE__ */ defineComponent({
  __name: "IgnorePreviewPanel",
  props: {
    result: {},
    loading: { type: Boolean }
  },
  setup(__props) {
    const { t } = useI18n();
    const props = __props;
    const viewMode = ref("dirs");
    const topRules = computed(() => {
      if (!props.result?.byRule) return {};
      const entries = Object.entries(props.result.byRule);
      entries.sort((a, b) => b[1] - a[1]);
      return Object.fromEntries(entries.slice(0, 10));
    });
    function formatNumber(n) {
      return n.toLocaleString("ru-RU");
    }
    return (_ctx, _cache) => {
      return openBlock(), createElementBlock("div", _hoisted_1$1, [
        createBaseVNode("div", _hoisted_2$1, [
          createBaseVNode("div", _hoisted_3$1, [
            _cache[3] || (_cache[3] = createBaseVNode("svg", {
              class: "w-4 h-4",
              fill: "none",
              stroke: "currentColor",
              viewBox: "0 0 24 24"
            }, [
              createBaseVNode("path", {
                "stroke-linecap": "round",
                "stroke-linejoin": "round",
                "stroke-width": "2",
                d: "M15 12a3 3 0 11-6 0 3 3 0 016 0z"
              }),
              createBaseVNode("path", {
                "stroke-linecap": "round",
                "stroke-linejoin": "round",
                "stroke-width": "2",
                d: "M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"
              })
            ], -1)),
            createBaseVNode("span", null, toDisplayString(unref(t)("ignoreModal.preview")), 1),
            __props.result?.totalFiles ? (openBlock(), createElementBlock("span", _hoisted_4$1, toDisplayString(formatNumber(__props.result.totalFiles)), 1)) : createCommentVNode("", true)
          ]),
          __props.result?.totalFiles ? (openBlock(), createElementBlock("div", _hoisted_5$1, [
            createBaseVNode("button", {
              class: normalizeClass(["preview-tab", viewMode.value === "dirs" ? "preview-tab--active" : ""]),
              onClick: _cache[0] || (_cache[0] = ($event) => viewMode.value = "dirs")
            }, "Папки", 2),
            createBaseVNode("button", {
              class: normalizeClass(["preview-tab", viewMode.value === "rules" ? "preview-tab--active" : ""]),
              onClick: _cache[1] || (_cache[1] = ($event) => viewMode.value = "rules")
            }, "Правила", 2),
            createBaseVNode("button", {
              class: normalizeClass(["preview-tab", viewMode.value === "files" ? "preview-tab--active" : ""]),
              onClick: _cache[2] || (_cache[2] = ($event) => viewMode.value = "files")
            }, "Файлы", 2)
          ])) : createCommentVNode("", true)
        ]),
        __props.loading ? (openBlock(), createElementBlock("div", _hoisted_6$1, [..._cache[4] || (_cache[4] = [
          createBaseVNode("svg", {
            class: "w-5 h-5 animate-spin",
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
          ], -1),
          createBaseVNode("span", null, "Анализ...", -1)
        ])])) : !__props.result?.totalFiles ? (openBlock(), createElementBlock("div", _hoisted_7$1, [
          _cache[5] || (_cache[5] = createBaseVNode("svg", {
            class: "w-8 h-8 text-gray-600",
            fill: "none",
            stroke: "currentColor",
            viewBox: "0 0 24 24"
          }, [
            createBaseVNode("path", {
              "stroke-linecap": "round",
              "stroke-linejoin": "round",
              "stroke-width": "1.5",
              d: "M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"
            })
          ], -1)),
          createBaseVNode("p", null, toDisplayString(unref(t)("ignoreModal.previewEmpty")), 1)
        ])) : (openBlock(), createElementBlock("div", _hoisted_8$1, [
          viewMode.value === "dirs" ? (openBlock(), createElementBlock("div", _hoisted_9$1, [
            (openBlock(true), createElementBlock(Fragment, null, renderList(__props.result.topDirs, (dir) => {
              return openBlock(), createElementBlock("div", {
                key: dir.dir,
                class: "preview-dir-item"
              }, [
                _cache[6] || (_cache[6] = createBaseVNode("svg", {
                  class: "w-4 h-4 text-amber-400/70",
                  fill: "none",
                  stroke: "currentColor",
                  viewBox: "0 0 24 24"
                }, [
                  createBaseVNode("path", {
                    "stroke-linecap": "round",
                    "stroke-linejoin": "round",
                    "stroke-width": "2",
                    d: "M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-6l-2-2H5a2 2 0 00-2 2z"
                  })
                ], -1)),
                createBaseVNode("span", _hoisted_10$1, toDisplayString(dir.dir), 1),
                createBaseVNode("span", _hoisted_11$1, toDisplayString(formatNumber(dir.count)), 1)
              ]);
            }), 128))
          ])) : viewMode.value === "rules" ? (openBlock(), createElementBlock("div", _hoisted_12$1, [
            (openBlock(true), createElementBlock(Fragment, null, renderList(topRules.value, (count, rule) => {
              return openBlock(), createElementBlock("div", {
                key: rule,
                class: "preview-rule-item"
              }, [
                createBaseVNode("code", _hoisted_13$1, toDisplayString(rule), 1),
                createBaseVNode("span", _hoisted_14$1, toDisplayString(formatNumber(count)), 1)
              ]);
            }), 128))
          ])) : (openBlock(), createElementBlock("div", _hoisted_15$1, [
            (openBlock(true), createElementBlock(Fragment, null, renderList(__props.result.sampleFiles, (file) => {
              return openBlock(), createElementBlock("div", {
                key: file,
                class: "preview-file-item"
              }, [
                _cache[7] || (_cache[7] = createBaseVNode("svg", {
                  class: "w-3.5 h-3.5 text-red-400/70 flex-shrink-0",
                  fill: "none",
                  stroke: "currentColor",
                  viewBox: "0 0 24 24"
                }, [
                  createBaseVNode("path", {
                    "stroke-linecap": "round",
                    "stroke-linejoin": "round",
                    "stroke-width": "2",
                    d: "M18.364 18.364A9 9 0 005.636 5.636m12.728 12.728A9 9 0 015.636 5.636m12.728 12.728L5.636 5.636"
                  })
                ], -1)),
                createBaseVNode("span", _hoisted_16$1, toDisplayString(file), 1)
              ]);
            }), 128)),
            __props.result.totalFiles > __props.result.sampleFiles.length ? (openBlock(), createElementBlock("div", _hoisted_17$1, " и ещё " + toDisplayString(formatNumber(__props.result.totalFiles - __props.result.sampleFiles.length)) + " файлов... ", 1)) : createCommentVNode("", true)
          ]))
        ]))
      ]);
    };
  }
});
const IgnorePreviewPanel = /* @__PURE__ */ _export_sfc(_sfc_main$1, [["__scopeId", "data-v-adc61bae"]]);
const _hoisted_1 = { class: "ignore-modal__header" };
const _hoisted_2 = { class: "flex items-center gap-3" };
const _hoisted_3 = { class: "text-lg font-semibold text-white" };
const _hoisted_4 = { class: "text-xs text-gray-400" };
const _hoisted_5 = { class: "ignore-modal__tabs" };
const _hoisted_6 = { class: "ignore-modal__content" };
const _hoisted_7 = {
  key: 0,
  class: "ignore-custom-layout"
};
const _hoisted_8 = { class: "ignore-editor-section" };
const _hoisted_9 = { class: "ignore-editor-header" };
const _hoisted_10 = { class: "ignore-editor-label" };
const _hoisted_11 = { class: "ignore-editor" };
const _hoisted_12 = ["placeholder"];
const _hoisted_13 = { class: "ignore-editor__hint" };
const _hoisted_14 = {
  key: 1,
  class: "ignore-custom-layout"
};
const _hoisted_15 = { class: "ignore-editor-section" };
const _hoisted_16 = { class: "ignore-editor-header" };
const _hoisted_17 = { class: "ignore-editor-label" };
const _hoisted_18 = { class: "ignore-editor" };
const _hoisted_19 = ["placeholder"];
const _hoisted_20 = { class: "ignore-modal__footer" };
const _hoisted_21 = { class: "ignore-footer__left" };
const _hoisted_22 = { class: "ignore-footer__right" };
const _hoisted_23 = ["disabled"];
const _hoisted_24 = {
  key: 0,
  class: "w-4 h-4",
  fill: "none",
  stroke: "currentColor",
  viewBox: "0 0 24 24"
};
const _hoisted_25 = {
  key: 1,
  class: "w-4 h-4 animate-spin",
  fill: "none",
  viewBox: "0 0 24 24"
};
const _sfc_main = /* @__PURE__ */ defineComponent({
  __name: "IgnoreRulesModal",
  setup(__props, { expose: __expose }) {
    const { t } = useI18n();
    const settingsStore = useSettingsStore();
    const projectStore = useProjectStore();
    const uiStore = useUIStore();
    const isOpen = ref(false);
    const currentTab = ref("custom");
    const gitignoreContent = ref("");
    const customRules = ref("");
    const isSaving = ref(false);
    const gitignorePreview = ref(null);
    const customPreview = ref(null);
    const gitignoreLoading = ref(false);
    const customLoading = ref(false);
    const gitignoreLinesRef = ref(null);
    const gitignoreTextareaRef = ref(null);
    const customLinesRef = ref(null);
    const customTextareaRef = ref(null);
    function syncGitignoreScroll() {
      if (gitignoreLinesRef.value && gitignoreTextareaRef.value) {
        gitignoreLinesRef.value.scrollTop = gitignoreTextareaRef.value.scrollTop;
      }
    }
    function syncCustomScroll() {
      if (customLinesRef.value && customTextareaRef.value) {
        customLinesRef.value.scrollTop = customTextareaRef.value.scrollTop;
      }
    }
    const rulesCount = computed(() => parseIgnoreRules(customRules.value).length);
    const gitignoreRulesCount = computed(() => parseIgnoreRules(gitignoreContent.value).length);
    const gitignoreLineCount = computed(() => Math.max(gitignoreContent.value.split("\n").length, 10));
    const customLineCount = computed(() => Math.max(customRules.value.split("\n").length, 10));
    function isCommentLine(lineNum) {
      const lines = customRules.value.split("\n");
      return lines[lineNum - 1]?.trim().startsWith("#") || false;
    }
    function isGitignoreCommentLine(lineNum) {
      const lines = gitignoreContent.value.split("\n");
      return lines[lineNum - 1]?.trim().startsWith("#") || false;
    }
    let debounceTimer = null;
    watch(customRules, (newRules) => {
      if (debounceTimer) clearTimeout(debounceTimer);
      debounceTimer = setTimeout(() => {
        if (newRules.trim() && projectStore.currentPath) {
          loadCustomPreview();
        } else {
          customPreview.value = null;
        }
      }, 500);
    });
    async function loadCustomPreview() {
      if (!projectStore.currentPath) return;
      customLoading.value = true;
      try {
        customPreview.value = await apiService.testIgnoreRulesDetailed(projectStore.currentPath, customRules.value);
      } catch {
        customPreview.value = null;
      } finally {
        customLoading.value = false;
      }
    }
    async function loadGitignorePreview() {
      if (!projectStore.currentPath || !gitignoreContent.value.trim()) return;
      gitignoreLoading.value = true;
      try {
        gitignorePreview.value = await apiService.testIgnoreRulesDetailed(projectStore.currentPath, gitignoreContent.value);
      } catch {
        gitignorePreview.value = null;
      } finally {
        gitignoreLoading.value = false;
      }
    }
    async function open() {
      isOpen.value = true;
      customRules.value = settingsStore.getCustomIgnoreRules();
      if (projectStore.currentPath) {
        try {
          gitignoreContent.value = await apiService.getGitignoreContent(projectStore.currentPath);
          loadGitignorePreview();
        } catch {
          gitignoreContent.value = "# No .gitignore file found";
        }
      }
      if (customRules.value.trim()) loadCustomPreview();
    }
    function close() {
      isOpen.value = false;
      gitignorePreview.value = null;
      customPreview.value = null;
    }
    async function save() {
      isSaving.value = true;
      try {
        await apiService.updateCustomIgnoreRules(customRules.value);
        settingsStore.setCustomIgnoreRules(customRules.value);
        uiStore.addToast(t("ignoreModal.saveSuccess"), "success");
        close();
      } catch {
        uiStore.addToast("Failed to save ignore rules", "error");
      } finally {
        isSaving.value = false;
      }
    }
    function resetToDefaults() {
      customRules.value = "";
      customPreview.value = null;
    }
    function clearAll() {
      customRules.value = "";
      customPreview.value = null;
    }
    __expose({ open, close });
    return (_ctx, _cache) => {
      return openBlock(), createBlock(Teleport, { to: "body" }, [
        createVNode(Transition, { name: "modal-backdrop" }, {
          default: withCtx(() => [
            isOpen.value ? (openBlock(), createElementBlock("div", {
              key: 0,
              class: "fixed inset-0 bg-black/60 backdrop-blur-sm z-40 flex items-center justify-center p-4",
              onClick: withModifiers(close, ["self"])
            }, [
              createVNode(Transition, { name: "modal" }, {
                default: withCtx(() => [
                  isOpen.value ? (openBlock(), createElementBlock("div", {
                    key: 0,
                    class: "ignore-modal",
                    onClick: _cache[4] || (_cache[4] = withModifiers(() => {
                    }, ["stop"]))
                  }, [
                    createBaseVNode("div", _hoisted_1, [
                      createBaseVNode("div", _hoisted_2, [
                        _cache[5] || (_cache[5] = createBaseVNode("div", { class: "ignore-modal__icon" }, [
                          createBaseVNode("svg", {
                            class: "w-5 h-5",
                            fill: "none",
                            stroke: "currentColor",
                            viewBox: "0 0 24 24"
                          }, [
                            createBaseVNode("path", {
                              "stroke-linecap": "round",
                              "stroke-linejoin": "round",
                              "stroke-width": "2",
                              d: "M13.875 18.825A10.05 10.05 0 0112 19c-4.478 0-8.268-2.943-9.543-7a9.97 9.97 0 011.563-3.029m5.858.908a3 3 0 114.243 4.243M9.878 9.878l4.242 4.242M9.88 9.88l-3.29-3.29m7.532 7.532l3.29 3.29M3 3l3.59 3.59m0 0A9.953 9.953 0 0112 5c4.478 0 8.268 2.943 9.543 7a10.025 10.025 0 01-4.132 5.411m0 0L21 21"
                            })
                          ])
                        ], -1)),
                        createBaseVNode("div", null, [
                          createBaseVNode("h3", _hoisted_3, toDisplayString(unref(t)("ignoreModal.title")), 1),
                          createBaseVNode("p", _hoisted_4, toDisplayString(unref(t)("ignoreModal.subtitle")), 1)
                        ])
                      ]),
                      createBaseVNode("button", {
                        onClick: close,
                        class: "ignore-modal__close"
                      }, [..._cache[6] || (_cache[6] = [
                        createBaseVNode("svg", {
                          class: "w-5 h-5",
                          fill: "none",
                          stroke: "currentColor",
                          viewBox: "0 0 24 24"
                        }, [
                          createBaseVNode("path", {
                            "stroke-linecap": "round",
                            "stroke-linejoin": "round",
                            "stroke-width": "2",
                            d: "M6 18L18 6M6 6l12 12"
                          })
                        ], -1)
                      ])])
                    ]),
                    createBaseVNode("div", _hoisted_5, [
                      createBaseVNode("button", {
                        onClick: _cache[0] || (_cache[0] = ($event) => currentTab.value = "gitignore"),
                        class: normalizeClass(["ignore-tab", currentTab.value === "gitignore" ? "ignore-tab--active" : ""])
                      }, [..._cache[7] || (_cache[7] = [
                        createBaseVNode("svg", {
                          class: "w-4 h-4",
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
                        ], -1),
                        createTextVNode(" .gitignore ", -1)
                      ])], 2),
                      createBaseVNode("button", {
                        onClick: _cache[1] || (_cache[1] = ($event) => currentTab.value = "custom"),
                        class: normalizeClass(["ignore-tab", currentTab.value === "custom" ? "ignore-tab--active" : ""])
                      }, [
                        _cache[8] || (_cache[8] = createBaseVNode("svg", {
                          class: "w-4 h-4",
                          fill: "none",
                          stroke: "currentColor",
                          viewBox: "0 0 24 24"
                        }, [
                          createBaseVNode("path", {
                            "stroke-linecap": "round",
                            "stroke-linejoin": "round",
                            "stroke-width": "2",
                            d: "M12 6V4m0 2a2 2 0 100 4m0-4a2 2 0 110 4m-6 8a2 2 0 100-4m0 4a2 2 0 110-4m0 4v2m0-6V4m6 6v10m6-2a2 2 0 100-4m0 4a2 2 0 110-4m0 4v2m0-6V4"
                          })
                        ], -1)),
                        createTextVNode(" " + toDisplayString(unref(t)("ignoreModal.customRules")), 1)
                      ], 2)
                    ]),
                    createBaseVNode("div", _hoisted_6, [
                      currentTab.value === "gitignore" ? (openBlock(), createElementBlock("div", _hoisted_7, [
                        createBaseVNode("div", _hoisted_8, [
                          createBaseVNode("div", _hoisted_9, [
                            createBaseVNode("span", _hoisted_10, [
                              createTextVNode(toDisplayString(unref(t)("ignoreModal.rulesCount")) + ": ", 1),
                              createBaseVNode("strong", null, toDisplayString(gitignoreRulesCount.value), 1)
                            ])
                          ]),
                          createBaseVNode("div", _hoisted_11, [
                            createBaseVNode("div", {
                              ref_key: "gitignoreLinesRef",
                              ref: gitignoreLinesRef,
                              class: "ignore-editor__lines"
                            }, [
                              (openBlock(true), createElementBlock(Fragment, null, renderList(gitignoreLineCount.value, (n) => {
                                return openBlock(), createElementBlock("span", {
                                  key: n,
                                  class: normalizeClass({ "ignore-editor__line--comment": isGitignoreCommentLine(n) })
                                }, toDisplayString(n), 3);
                              }), 128))
                            ], 512),
                            withDirectives(createBaseVNode("textarea", {
                              ref_key: "gitignoreTextareaRef",
                              ref: gitignoreTextareaRef,
                              "onUpdate:modelValue": _cache[2] || (_cache[2] = ($event) => gitignoreContent.value = $event),
                              readonly: "",
                              class: "ignore-editor__textarea",
                              placeholder: unref(t)("ignoreModal.gitignorePlaceholder"),
                              onScroll: syncGitignoreScroll
                            }, null, 40, _hoisted_12), [
                              [vModelText, gitignoreContent.value]
                            ])
                          ]),
                          createBaseVNode("p", _hoisted_13, [
                            _cache[9] || (_cache[9] = createBaseVNode("svg", {
                              class: "w-3.5 h-3.5",
                              fill: "none",
                              stroke: "currentColor",
                              viewBox: "0 0 24 24"
                            }, [
                              createBaseVNode("path", {
                                "stroke-linecap": "round",
                                "stroke-linejoin": "round",
                                "stroke-width": "2",
                                d: "M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
                              })
                            ], -1)),
                            createTextVNode(" " + toDisplayString(unref(t)("ignoreModal.gitignoreInfo")), 1)
                          ])
                        ]),
                        createVNode(IgnorePreviewPanel, {
                          result: gitignorePreview.value,
                          loading: gitignoreLoading.value
                        }, null, 8, ["result", "loading"])
                      ])) : createCommentVNode("", true),
                      currentTab.value === "custom" ? (openBlock(), createElementBlock("div", _hoisted_14, [
                        createBaseVNode("div", _hoisted_15, [
                          createBaseVNode("div", _hoisted_16, [
                            createBaseVNode("span", _hoisted_17, [
                              createTextVNode(toDisplayString(unref(t)("ignoreModal.rulesCount")) + ": ", 1),
                              createBaseVNode("strong", null, toDisplayString(rulesCount.value), 1)
                            ])
                          ]),
                          createBaseVNode("div", _hoisted_18, [
                            createBaseVNode("div", {
                              ref_key: "customLinesRef",
                              ref: customLinesRef,
                              class: "ignore-editor__lines"
                            }, [
                              (openBlock(true), createElementBlock(Fragment, null, renderList(customLineCount.value, (n) => {
                                return openBlock(), createElementBlock("span", {
                                  key: n,
                                  class: normalizeClass({ "ignore-editor__line--comment": isCommentLine(n) })
                                }, toDisplayString(n), 3);
                              }), 128))
                            ], 512),
                            withDirectives(createBaseVNode("textarea", {
                              ref_key: "customTextareaRef",
                              ref: customTextareaRef,
                              "onUpdate:modelValue": _cache[3] || (_cache[3] = ($event) => customRules.value = $event),
                              class: "ignore-editor__textarea ignore-editor__textarea--editable",
                              placeholder: unref(t)("ignoreModal.customPlaceholder"),
                              spellcheck: "false",
                              onScroll: syncCustomScroll
                            }, null, 40, _hoisted_19), [
                              [vModelText, customRules.value]
                            ])
                          ])
                        ]),
                        createVNode(IgnorePreviewPanel, {
                          result: customPreview.value,
                          loading: customLoading.value
                        }, null, 8, ["result", "loading"])
                      ])) : createCommentVNode("", true)
                    ]),
                    createBaseVNode("div", _hoisted_20, [
                      createBaseVNode("div", _hoisted_21, [
                        currentTab.value === "custom" ? (openBlock(), createElementBlock("button", {
                          key: 0,
                          onClick: resetToDefaults,
                          class: "ignore-footer__danger-btn"
                        }, toDisplayString(unref(t)("ignoreModal.reset")), 1)) : createCommentVNode("", true),
                        currentTab.value === "custom" ? (openBlock(), createElementBlock("button", {
                          key: 1,
                          onClick: clearAll,
                          class: "ignore-footer__danger-btn"
                        }, toDisplayString(unref(t)("ignoreModal.clearAll")), 1)) : createCommentVNode("", true)
                      ]),
                      createBaseVNode("div", _hoisted_22, [
                        createBaseVNode("button", {
                          onClick: close,
                          class: "btn btn-secondary"
                        }, toDisplayString(unref(t)("ignoreModal.cancel")), 1),
                        currentTab.value === "custom" ? (openBlock(), createElementBlock("button", {
                          key: 0,
                          onClick: save,
                          disabled: isSaving.value,
                          class: "ignore-footer__save-btn"
                        }, [
                          !isSaving.value ? (openBlock(), createElementBlock("svg", _hoisted_24, [..._cache[10] || (_cache[10] = [
                            createBaseVNode("path", {
                              "stroke-linecap": "round",
                              "stroke-linejoin": "round",
                              "stroke-width": "2",
                              d: "M5 13l4 4L19 7"
                            }, null, -1)
                          ])])) : (openBlock(), createElementBlock("svg", _hoisted_25, [..._cache[11] || (_cache[11] = [
                            createBaseVNode("circle", {
                              class: "opacity-25",
                              cx: "12",
                              cy: "12",
                              r: "10",
                              stroke: "currentColor",
                              "stroke-width": "4"
                            }, null, -1),
                            createBaseVNode("path", {
                              class: "opacity-75",
                              fill: "currentColor",
                              d: "M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z"
                            }, null, -1)
                          ])])),
                          createTextVNode(" " + toDisplayString(unref(t)("ignoreModal.save")), 1)
                        ], 8, _hoisted_23)) : createCommentVNode("", true)
                      ])
                    ])
                  ])) : createCommentVNode("", true)
                ]),
                _: 1
              })
            ])) : createCommentVNode("", true)
          ]),
          _: 1
        })
      ]);
    };
  }
});
const IgnoreRulesModal = /* @__PURE__ */ _export_sfc(_sfc_main, [["__scopeId", "data-v-da8ea94e"]]);
export {
  IgnoreRulesModal as default
};
