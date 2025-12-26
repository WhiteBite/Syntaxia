import { Y as Ye, h as he, G as Ge, V as Ve, S as Se } from "./ui-vendor-2hF-XK7S.js";
import { e as defineComponent, J as createBlock, z as openBlock, G as withCtx, B as createVNode, u as unref, C as createBaseVNode, M as createTextVNode, y as createElementBlock, F as Fragment, H as renderList, D as toDisplayString } from "./vue-vendor-DLLjc1XI.js";
const _hoisted_1 = { class: "fixed inset-0 overflow-y-auto" };
const _hoisted_2 = { class: "flex min-h-full items-center justify-center p-4" };
const _hoisted_3 = { class: "px-6 py-4 border-b border-gray-200 dark:border-gray-700" };
const _hoisted_4 = { class: "px-6 py-4 max-h-[70vh] overflow-y-auto" };
const _hoisted_5 = { class: "grid grid-cols-1 md:grid-cols-2 gap-6" };
const _hoisted_6 = { class: "text-sm font-semibold text-gray-400 dark:text-gray-400 uppercase tracking-wider mb-3" };
const _hoisted_7 = { class: "space-y-2" };
const _hoisted_8 = { class: "text-sm text-gray-700 dark:text-gray-300" };
const _hoisted_9 = { class: "flex gap-1" };
const _hoisted_10 = { class: "px-6 py-4 bg-gray-50 dark:bg-gray-900 border-t border-gray-200 dark:border-gray-700 flex justify-end" };
const _sfc_main = /* @__PURE__ */ defineComponent({
  __name: "KeyboardShortcutsModal",
  props: {
    isOpen: { type: Boolean }
  },
  emits: ["close"],
  setup(__props) {
    const categories = [
      {
        name: "General",
        shortcuts: [
          { key: "Ctrl+K", description: "Open command palette" },
          { key: "Ctrl+/", description: "Show keyboard shortcuts" },
          { key: "Escape", description: "Close modal" }
        ]
      },
      {
        name: "File",
        shortcuts: [
          { key: "Ctrl+N", description: "New file" },
          { key: "Ctrl+O", description: "Open file" },
          { key: "Ctrl+S", description: "Save file" },
          { key: "Ctrl+P", description: "Quick open file" }
        ]
      },
      {
        name: "Tasks",
        shortcuts: [
          { key: "Ctrl+Enter", description: "Analyze task" },
          { key: "Ctrl+B", description: "Build context" },
          { key: "Ctrl+Shift+A", description: "New task" },
          { key: "Ctrl+Shift+D", description: "Duplicate task" }
        ]
      },
      {
        name: "Navigation",
        shortcuts: [
          { key: "Ctrl+1", description: "Focus file explorer" },
          { key: "Ctrl+2", description: "Focus task panel" },
          { key: "Ctrl+3", description: "Focus output" },
          { key: "Alt+←", description: "Go back" },
          { key: "Alt+→", description: "Go forward" }
        ]
      },
      {
        name: "View",
        shortcuts: [
          { key: "Ctrl++", description: "Zoom in" },
          { key: "Ctrl+-", description: "Zoom out" },
          { key: "Ctrl+0", description: "Reset zoom" },
          { key: "F11", description: "Toggle fullscreen" }
        ]
      },
      {
        name: "Window",
        shortcuts: [
          { key: "Ctrl+R", description: "Reload window" },
          { key: "Ctrl+,", description: "Open settings" }
        ]
      }
    ];
    return (_ctx, _cache) => {
      return openBlock(), createBlock(unref(Se), {
        show: __props.isOpen,
        as: "template"
      }, {
        default: withCtx(() => [
          createVNode(unref(Ye), {
            onClose: _cache[1] || (_cache[1] = ($event) => _ctx.$emit("close")),
            class: "relative z-50"
          }, {
            default: withCtx(() => [
              createVNode(unref(he), {
                as: "template",
                enter: "ease-out duration-200",
                "enter-from": "opacity-0",
                "enter-to": "opacity-100",
                leave: "ease-in duration-150",
                "leave-from": "opacity-100",
                "leave-to": "opacity-0"
              }, {
                default: withCtx(() => [..._cache[2] || (_cache[2] = [
                  createBaseVNode("div", { class: "fixed inset-0 bg-black/50 backdrop-blur-sm" }, null, -1)
                ])]),
                _: 1
              }),
              createBaseVNode("div", _hoisted_1, [
                createBaseVNode("div", _hoisted_2, [
                  createVNode(unref(he), {
                    as: "template",
                    enter: "ease-out duration-200",
                    "enter-from": "opacity-0 scale-95",
                    "enter-to": "opacity-100 scale-100",
                    leave: "ease-in duration-150",
                    "leave-from": "opacity-100 scale-100",
                    "leave-to": "opacity-0 scale-95"
                  }, {
                    default: withCtx(() => [
                      createVNode(unref(Ge), { class: "w-full max-w-4xl transform overflow-hidden rounded-xl bg-white dark:bg-gray-800 shadow-2xl transition-all" }, {
                        default: withCtx(() => [
                          createBaseVNode("div", _hoisted_3, [
                            createVNode(unref(Ve), { class: "text-xl font-semibold text-gray-900 dark:text-gray-100" }, {
                              default: withCtx(() => [..._cache[3] || (_cache[3] = [
                                createTextVNode(" Keyboard Shortcuts ", -1)
                              ])]),
                              _: 1
                            })
                          ]),
                          createBaseVNode("div", _hoisted_4, [
                            createBaseVNode("div", _hoisted_5, [
                              (openBlock(), createElementBlock(Fragment, null, renderList(categories, (category) => {
                                return createBaseVNode("div", {
                                  key: category.name
                                }, [
                                  createBaseVNode("h3", _hoisted_6, toDisplayString(category.name), 1),
                                  createBaseVNode("div", _hoisted_7, [
                                    (openBlock(true), createElementBlock(Fragment, null, renderList(category.shortcuts, (shortcut) => {
                                      return openBlock(), createElementBlock("div", {
                                        key: shortcut.key,
                                        class: "flex items-center justify-between py-2 px-3 rounded-lg hover:bg-gray-50 dark:hover:bg-gray-700/50 transition-colors"
                                      }, [
                                        createBaseVNode("span", _hoisted_8, toDisplayString(shortcut.description), 1),
                                        createBaseVNode("div", _hoisted_9, [
                                          (openBlock(true), createElementBlock(Fragment, null, renderList(shortcut.key.split("+"), (key, index) => {
                                            return openBlock(), createElementBlock("kbd", {
                                              key: index,
                                              class: "px-2 py-1 text-xs font-semibold text-gray-600 dark:text-gray-300 bg-gray-100 dark:bg-gray-700 border border-gray-300 dark:border-gray-600 rounded shadow-sm"
                                            }, toDisplayString(key), 1);
                                          }), 128))
                                        ])
                                      ]);
                                    }), 128))
                                  ])
                                ]);
                              }), 64))
                            ])
                          ]),
                          createBaseVNode("div", _hoisted_10, [
                            createBaseVNode("button", {
                              onClick: _cache[0] || (_cache[0] = ($event) => _ctx.$emit("close")),
                              class: "px-4 py-2 text-sm font-medium text-white bg-blue-600 hover:bg-blue-700 rounded-lg transition-colors"
                            }, " Close ")
                          ])
                        ]),
                        _: 1
                      })
                    ]),
                    _: 1
                  })
                ])
              ])
            ]),
            _: 1
          })
        ]),
        _: 1
      }, 8, ["show"]);
    };
  }
});
export {
  _sfc_main as default
};
