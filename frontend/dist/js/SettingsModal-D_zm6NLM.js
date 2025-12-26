import { e as defineComponent, r as ref, f as onMounted, y as createElementBlock, z as openBlock, C as createBaseVNode, B as createVNode, u as unref, A as createCommentVNode, D as toDisplayString, J as createBlock, M as createTextVNode, l as normalizeClass, c as computed, w as watch, g as onUnmounted, G as withCtx, F as Fragment, H as renderList, S as resolveDynamicComponent, E as Transition, N as withDirectives, Y as vModelSelect, R as vModelCheckbox, T as Teleport } from "./vue-vendor-DLLjc1XI.js";
import { c as useI18n, u as useUIStore, s as shellApi, d as useSettingsStore, e as useOnboarding, A as AISettings, E as ExportSettings, f as _export_sfc } from "./index-CCibvNLt.js";
import { s as Monitor, p as LoaderCircle, P as Plus, T as Trash2, t as Settings, o as Lightbulb, F as FileText, i as FolderTree, X, G as Globe, u as CircleQuestionMark, f as Sparkles, v as Funnel } from "./icons-BGgRd9bo.js";
import "./ui-vendor-2hF-XK7S.js";
import "./utils-DVYXdora.js";
import "./highlighter-D41BBaGQ.js";
const _hoisted_1$1 = { class: "space-y-4" };
const _hoisted_2$1 = { class: "flex items-start gap-3 p-4 rounded-lg bg-gray-800/50 border border-gray-700/30" };
const _hoisted_3$1 = { class: "flex-1 min-w-0" };
const _hoisted_4$1 = { class: "text-sm font-medium text-white mb-1" };
const _hoisted_5$1 = { class: "text-xs text-gray-400 mb-3" };
const _hoisted_6$1 = { class: "flex items-center gap-3" };
const _hoisted_7$1 = ["disabled"];
const _hoisted_8$1 = ["disabled"];
const _hoisted_9$1 = {
  key: 0,
  class: "text-xs text-gray-500 mt-2"
};
const _sfc_main$1 = /* @__PURE__ */ defineComponent({
  __name: "ShellIntegrationSettings",
  setup(__props) {
    const { t } = useI18n();
    const uiStore = useUIStore();
    const isRegistered = ref(false);
    const isLoading = ref(false);
    const currentOS = ref("");
    async function loadStatus() {
      try {
        const status = await shellApi.getStatus();
        isRegistered.value = status.isRegistered;
        currentOS.value = status.currentOS;
      } catch {
        uiStore.addToast(t("settings.shellIntegration.error"), "error");
      }
    }
    async function handleRegister() {
      isLoading.value = true;
      try {
        await shellApi.register();
        isRegistered.value = true;
        uiStore.addToast(t("settings.shellIntegration.enableSuccess"), "success");
      } catch {
        uiStore.addToast(t("settings.shellIntegration.error"), "error");
      } finally {
        isLoading.value = false;
      }
    }
    async function handleUnregister() {
      isLoading.value = true;
      try {
        await shellApi.unregister();
        isRegistered.value = false;
        uiStore.addToast(t("settings.shellIntegration.disableSuccess"), "success");
      } catch {
        uiStore.addToast(t("settings.shellIntegration.error"), "error");
      } finally {
        isLoading.value = false;
      }
    }
    onMounted(loadStatus);
    return (_ctx, _cache) => {
      return openBlock(), createElementBlock("div", _hoisted_1$1, [
        createBaseVNode("div", _hoisted_2$1, [
          createVNode(unref(Monitor), { class: "w-5 h-5 text-indigo-400 mt-0.5 flex-shrink-0" }),
          createBaseVNode("div", _hoisted_3$1, [
            createBaseVNode("h3", _hoisted_4$1, toDisplayString(unref(t)("settings.shellIntegration.title")), 1),
            createBaseVNode("p", _hoisted_5$1, toDisplayString(unref(t)("settings.shellIntegration.description")), 1),
            createBaseVNode("div", _hoisted_6$1, [
              !isRegistered.value ? (openBlock(), createElementBlock("button", {
                key: 0,
                onClick: handleRegister,
                disabled: isLoading.value,
                class: "btn-unified btn-unified-primary text-sm"
              }, [
                isLoading.value ? (openBlock(), createBlock(unref(LoaderCircle), {
                  key: 0,
                  class: "w-4 h-4 animate-spin"
                })) : (openBlock(), createBlock(unref(Plus), {
                  key: 1,
                  class: "w-4 h-4"
                })),
                createTextVNode(" " + toDisplayString(unref(t)("settings.shellIntegration.enable")), 1)
              ], 8, _hoisted_7$1)) : (openBlock(), createElementBlock("button", {
                key: 1,
                onClick: handleUnregister,
                disabled: isLoading.value,
                class: "btn-unified btn-unified-secondary text-sm"
              }, [
                isLoading.value ? (openBlock(), createBlock(unref(LoaderCircle), {
                  key: 0,
                  class: "w-4 h-4 animate-spin"
                })) : (openBlock(), createBlock(unref(Trash2), {
                  key: 1,
                  class: "w-4 h-4"
                })),
                createTextVNode(" " + toDisplayString(unref(t)("settings.shellIntegration.disable")), 1)
              ], 8, _hoisted_8$1)),
              createBaseVNode("span", {
                class: normalizeClass(["text-xs px-2 py-1 rounded", isRegistered.value ? "bg-green-500/20 text-green-400" : "bg-gray-600/30 text-gray-400"])
              }, toDisplayString(isRegistered.value ? unref(t)("settings.shellIntegration.enabled") : unref(t)("settings.shellIntegration.disabled")), 3)
            ]),
            currentOS.value === "windows" ? (openBlock(), createElementBlock("p", _hoisted_9$1, toDisplayString(unref(t)("settings.shellIntegration.requiresAdmin")), 1)) : createCommentVNode("", true)
          ])
        ])
      ]);
    };
  }
});
const _hoisted_1 = {
  key: 0,
  class: "fixed inset-0 z-50 flex items-center justify-center p-4"
};
const _hoisted_2 = { class: "settings-modal relative" };
const _hoisted_3 = { class: "settings-header" };
const _hoisted_4 = { class: "flex items-center gap-2" };
const _hoisted_5 = { class: "settings-header-icon" };
const _hoisted_6 = { class: "text-base font-semibold text-white" };
const _hoisted_7 = { class: "settings-tabs" };
const _hoisted_8 = { class: "settings-tabs-container" };
const _hoisted_9 = ["onClick"];
const _hoisted_10 = { class: "settings-content" };
const _hoisted_11 = {
  key: "general",
  class: "settings-section"
};
const _hoisted_12 = { class: "settings-group" };
const _hoisted_13 = { class: "settings-group-header" };
const _hoisted_14 = { class: "settings-group" };
const _hoisted_15 = { class: "settings-group-header" };
const _hoisted_16 = { class: "settings-hint" };
const _hoisted_17 = {
  key: "ai",
  class: "settings-section"
};
const _hoisted_18 = {
  key: "export",
  class: "settings-section"
};
const _hoisted_19 = {
  key: "fileExplorer",
  class: "settings-section"
};
const _hoisted_20 = { class: "settings-group" };
const _hoisted_21 = { class: "settings-group-header" };
const _hoisted_22 = { class: "settings-toggle-list" };
const _hoisted_23 = { class: "settings-toggle-item" };
const _hoisted_24 = { class: "settings-toggle-info" };
const _hoisted_25 = { class: "settings-toggle-label" };
const _hoisted_26 = { class: "settings-toggle-hint" };
const _hoisted_27 = { class: "settings-toggle" };
const _hoisted_28 = { class: "settings-toggle-item" };
const _hoisted_29 = { class: "settings-toggle-info" };
const _hoisted_30 = { class: "settings-toggle-label" };
const _hoisted_31 = { class: "settings-toggle-hint" };
const _hoisted_32 = { class: "settings-toggle" };
const _hoisted_33 = { class: "settings-group" };
const _hoisted_34 = { class: "settings-group-header" };
const _hoisted_35 = { class: "settings-toggle-list" };
const _hoisted_36 = { class: "settings-toggle-item" };
const _hoisted_37 = { class: "settings-toggle-info" };
const _hoisted_38 = { class: "settings-toggle-label" };
const _hoisted_39 = { class: "settings-toggle-hint" };
const _hoisted_40 = { class: "settings-toggle" };
const _hoisted_41 = { class: "settings-toggle-item" };
const _hoisted_42 = { class: "settings-toggle-info" };
const _hoisted_43 = { class: "settings-toggle-label" };
const _hoisted_44 = { class: "settings-toggle-hint" };
const _hoisted_45 = { class: "settings-toggle" };
const _hoisted_46 = {
  key: "system",
  class: "settings-section"
};
const _hoisted_47 = { class: "settings-footer" };
const _sfc_main = /* @__PURE__ */ defineComponent({
  __name: "SettingsModal",
  props: {
    modelValue: { type: Boolean }
  },
  emits: ["update:modelValue"],
  setup(__props, { emit: __emit }) {
    const props = __props;
    const emit = __emit;
    const { t, setLocale, locale } = useI18n();
    const settingsStore = useSettingsStore();
    const { startTour, resetTour } = useOnboarding();
    const tabs = computed(() => [
      { id: "general", label: t("settings.modal.general"), icon: Settings },
      { id: "ai", label: t("settings.modal.ai"), icon: Lightbulb },
      { id: "export", label: t("settings.modal.export"), icon: FileText },
      { id: "fileExplorer", label: t("settings.modal.fileExplorer"), icon: FolderTree },
      { id: "system", label: t("settings.modal.system"), icon: Monitor }
    ]);
    const activeTab = ref("general");
    const selectedLanguage = ref(locale.value);
    watch(
      () => props.modelValue,
      (isOpen) => {
        if (isOpen) {
          selectedLanguage.value = locale.value;
          activeTab.value = "general";
        }
      }
    );
    watch(selectedLanguage, (newLang) => {
      if (newLang !== locale.value) {
        setLocale(newLang);
      }
    });
    function close() {
      emit("update:modelValue", false);
    }
    function handleStartTour() {
      close();
      resetTour();
      setTimeout(() => {
        startTour();
      }, 300);
    }
    function handleKeydown(e) {
      if (e.key === "Escape" && props.modelValue) {
        close();
      }
    }
    onMounted(() => {
      window.addEventListener("keydown", handleKeydown);
    });
    onUnmounted(() => {
      window.removeEventListener("keydown", handleKeydown);
    });
    return (_ctx, _cache) => {
      return openBlock(), createBlock(Teleport, { to: "body" }, [
        createVNode(Transition, { name: "modal" }, {
          default: withCtx(() => [
            __props.modelValue ? (openBlock(), createElementBlock("div", _hoisted_1, [
              createBaseVNode("div", {
                class: "absolute inset-0 bg-black/70 backdrop-blur-md",
                onClick: close
              }),
              createBaseVNode("div", _hoisted_2, [
                createBaseVNode("div", _hoisted_3, [
                  createBaseVNode("div", _hoisted_4, [
                    createBaseVNode("div", _hoisted_5, [
                      createVNode(unref(Settings), { class: "w-4 h-4" })
                    ]),
                    createBaseVNode("h2", _hoisted_6, toDisplayString(unref(t)("settings.modal.title")), 1)
                  ]),
                  createBaseVNode("button", {
                    onClick: close,
                    class: "settings-close-btn"
                  }, [
                    createVNode(unref(X), { class: "w-4 h-4" })
                  ])
                ]),
                createBaseVNode("div", _hoisted_7, [
                  createBaseVNode("div", _hoisted_8, [
                    (openBlock(true), createElementBlock(Fragment, null, renderList(tabs.value, (tab) => {
                      return openBlock(), createElementBlock("button", {
                        key: tab.id,
                        onClick: ($event) => activeTab.value = tab.id,
                        class: normalizeClass(["settings-tab", { "settings-tab-active": activeTab.value === tab.id }])
                      }, [
                        (openBlock(), createBlock(resolveDynamicComponent(tab.icon), { class: "w-4 h-4" })),
                        createBaseVNode("span", null, toDisplayString(tab.label), 1)
                      ], 10, _hoisted_9);
                    }), 128))
                  ])
                ]),
                createBaseVNode("div", _hoisted_10, [
                  createVNode(Transition, {
                    name: "tab-fade",
                    mode: "out-in"
                  }, {
                    default: withCtx(() => [
                      activeTab.value === "general" ? (openBlock(), createElementBlock("div", _hoisted_11, [
                        createBaseVNode("div", _hoisted_12, [
                          createBaseVNode("div", _hoisted_13, [
                            createVNode(unref(Globe), { class: "w-4 h-4 text-indigo-400" }),
                            createBaseVNode("span", null, toDisplayString(unref(t)("settings.modal.language")), 1)
                          ]),
                          withDirectives(createBaseVNode("select", {
                            "onUpdate:modelValue": _cache[0] || (_cache[0] = ($event) => selectedLanguage.value = $event),
                            class: "settings-select"
                          }, [..._cache[5] || (_cache[5] = [
                            createBaseVNode("option", { value: "ru" }, "Русский", -1),
                            createBaseVNode("option", { value: "en" }, "English", -1)
                          ])], 512), [
                            [vModelSelect, selectedLanguage.value]
                          ])
                        ]),
                        _cache[6] || (_cache[6] = createBaseVNode("div", { class: "settings-divider" }, null, -1)),
                        createBaseVNode("div", _hoisted_14, [
                          createBaseVNode("div", _hoisted_15, [
                            createVNode(unref(CircleQuestionMark), { class: "w-4 h-4 text-purple-400" }),
                            createBaseVNode("span", null, toDisplayString(unref(t)("onboarding.startTour")), 1)
                          ]),
                          createBaseVNode("p", _hoisted_16, toDisplayString(unref(t)("settings.modal.tourHint")), 1),
                          createBaseVNode("button", {
                            onClick: handleStartTour,
                            class: "settings-action-btn"
                          }, [
                            createVNode(unref(Sparkles), { class: "w-4 h-4" }),
                            createTextVNode(" " + toDisplayString(unref(t)("onboarding.startTour")), 1)
                          ])
                        ])
                      ])) : activeTab.value === "ai" ? (openBlock(), createElementBlock("div", _hoisted_17, [
                        createVNode(AISettings)
                      ])) : activeTab.value === "export" ? (openBlock(), createElementBlock("div", _hoisted_18, [
                        createVNode(ExportSettings)
                      ])) : activeTab.value === "fileExplorer" ? (openBlock(), createElementBlock("div", _hoisted_19, [
                        createBaseVNode("div", _hoisted_20, [
                          createBaseVNode("div", _hoisted_21, [
                            createVNode(unref(Funnel), { class: "w-4 h-4 text-emerald-400" }),
                            createBaseVNode("span", null, toDisplayString(unref(t)("settings.modal.filterSettings")), 1)
                          ]),
                          createBaseVNode("div", _hoisted_22, [
                            createBaseVNode("label", _hoisted_23, [
                              createBaseVNode("div", _hoisted_24, [
                                createBaseVNode("span", _hoisted_25, toDisplayString(unref(t)("settings.useGitignore")), 1),
                                createBaseVNode("span", _hoisted_26, toDisplayString(unref(t)("settings.useGitignoreHint")), 1)
                              ]),
                              createBaseVNode("div", _hoisted_27, [
                                withDirectives(createBaseVNode("input", {
                                  type: "checkbox",
                                  "onUpdate:modelValue": _cache[1] || (_cache[1] = ($event) => unref(settingsStore).settings.fileExplorer.useGitignore = $event),
                                  class: "sr-only peer"
                                }, null, 512), [
                                  [vModelCheckbox, unref(settingsStore).settings.fileExplorer.useGitignore]
                                ]),
                                _cache[7] || (_cache[7] = createBaseVNode("div", { class: "settings-toggle-track peer-checked:bg-indigo-500" }, null, -1)),
                                _cache[8] || (_cache[8] = createBaseVNode("div", { class: "settings-toggle-thumb peer-checked:translate-x-5" }, null, -1))
                              ])
                            ]),
                            createBaseVNode("label", _hoisted_28, [
                              createBaseVNode("div", _hoisted_29, [
                                createBaseVNode("span", _hoisted_30, toDisplayString(unref(t)("settings.useCustomIgnore")), 1),
                                createBaseVNode("span", _hoisted_31, toDisplayString(unref(t)("settings.useCustomIgnoreHint")), 1)
                              ]),
                              createBaseVNode("div", _hoisted_32, [
                                withDirectives(createBaseVNode("input", {
                                  type: "checkbox",
                                  "onUpdate:modelValue": _cache[2] || (_cache[2] = ($event) => unref(settingsStore).settings.fileExplorer.useCustomIgnore = $event),
                                  class: "sr-only peer"
                                }, null, 512), [
                                  [vModelCheckbox, unref(settingsStore).settings.fileExplorer.useCustomIgnore]
                                ]),
                                _cache[9] || (_cache[9] = createBaseVNode("div", { class: "settings-toggle-track peer-checked:bg-indigo-500" }, null, -1)),
                                _cache[10] || (_cache[10] = createBaseVNode("div", { class: "settings-toggle-thumb peer-checked:translate-x-5" }, null, -1))
                              ])
                            ])
                          ])
                        ]),
                        _cache[15] || (_cache[15] = createBaseVNode("div", { class: "settings-divider" }, null, -1)),
                        createBaseVNode("div", _hoisted_33, [
                          createBaseVNode("div", _hoisted_34, [
                            createVNode(unref(FolderTree), { class: "w-4 h-4 text-orange-400" }),
                            createBaseVNode("span", null, toDisplayString(unref(t)("settings.modal.displaySettings")), 1)
                          ]),
                          createBaseVNode("div", _hoisted_35, [
                            createBaseVNode("label", _hoisted_36, [
                              createBaseVNode("div", _hoisted_37, [
                                createBaseVNode("span", _hoisted_38, toDisplayString(unref(t)("settings.autoSaveSelection")), 1),
                                createBaseVNode("span", _hoisted_39, toDisplayString(unref(t)("settings.autoSaveSelectionHint")), 1)
                              ]),
                              createBaseVNode("div", _hoisted_40, [
                                withDirectives(createBaseVNode("input", {
                                  type: "checkbox",
                                  "onUpdate:modelValue": _cache[3] || (_cache[3] = ($event) => unref(settingsStore).settings.fileExplorer.autoSaveSelection = $event),
                                  class: "sr-only peer"
                                }, null, 512), [
                                  [vModelCheckbox, unref(settingsStore).settings.fileExplorer.autoSaveSelection]
                                ]),
                                _cache[11] || (_cache[11] = createBaseVNode("div", { class: "settings-toggle-track peer-checked:bg-indigo-500" }, null, -1)),
                                _cache[12] || (_cache[12] = createBaseVNode("div", { class: "settings-toggle-thumb peer-checked:translate-x-5" }, null, -1))
                              ])
                            ]),
                            createBaseVNode("label", _hoisted_41, [
                              createBaseVNode("div", _hoisted_42, [
                                createBaseVNode("span", _hoisted_43, toDisplayString(unref(t)("settings.compactFolders")), 1),
                                createBaseVNode("span", _hoisted_44, toDisplayString(unref(t)("settings.compactFoldersHint")), 1)
                              ]),
                              createBaseVNode("div", _hoisted_45, [
                                withDirectives(createBaseVNode("input", {
                                  type: "checkbox",
                                  "onUpdate:modelValue": _cache[4] || (_cache[4] = ($event) => unref(settingsStore).settings.fileExplorer.compactNestedFolders = $event),
                                  class: "sr-only peer"
                                }, null, 512), [
                                  [vModelCheckbox, unref(settingsStore).settings.fileExplorer.compactNestedFolders]
                                ]),
                                _cache[13] || (_cache[13] = createBaseVNode("div", { class: "settings-toggle-track peer-checked:bg-indigo-500" }, null, -1)),
                                _cache[14] || (_cache[14] = createBaseVNode("div", { class: "settings-toggle-thumb peer-checked:translate-x-5" }, null, -1))
                              ])
                            ])
                          ])
                        ])
                      ])) : activeTab.value === "system" ? (openBlock(), createElementBlock("div", _hoisted_46, [
                        createVNode(_sfc_main$1)
                      ])) : createCommentVNode("", true)
                    ]),
                    _: 1
                  })
                ]),
                createBaseVNode("div", _hoisted_47, [
                  createBaseVNode("button", {
                    onClick: close,
                    class: "settings-close-action"
                  }, toDisplayString(unref(t)("settings.modal.cancel")), 1)
                ])
              ])
            ])) : createCommentVNode("", true)
          ]),
          _: 1
        })
      ]);
    };
  }
});
const SettingsModal = /* @__PURE__ */ _export_sfc(_sfc_main, [["__scopeId", "data-v-87160c16"]]);
export {
  SettingsModal as default
};
