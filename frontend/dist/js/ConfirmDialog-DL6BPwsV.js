import { r as ref, e as defineComponent, w as watch, n as nextTick, c as computed, J as createBlock, z as openBlock, B as createVNode, G as withCtx, y as createElementBlock, A as createCommentVNode, u as unref, O as withKeys, K as withModifiers, C as createBaseVNode, l as normalizeClass, D as toDisplayString, E as Transition, T as Teleport } from "./vue-vendor-DLLjc1XI.js";
import { c as useI18n } from "./index-CCibvNLt.js";
import { h as TriangleAlert, I as Info } from "./icons-BGgRd9bo.js";
import "./ui-vendor-2hF-XK7S.js";
import "./utils-DVYXdora.js";
import "./highlighter-D41BBaGQ.js";
const isOpen = ref(false);
const options = ref(null);
let resolvePromise = null;
function useConfirm() {
  function confirm(opts) {
    options.value = opts;
    isOpen.value = true;
    return new Promise((resolve) => {
      resolvePromise = resolve;
    });
  }
  function handleConfirm() {
    isOpen.value = false;
    resolvePromise?.(true);
    resolvePromise = null;
  }
  function handleCancel() {
    isOpen.value = false;
    resolvePromise?.(false);
    resolvePromise = null;
  }
  return {
    isOpen,
    options,
    confirm,
    handleConfirm,
    handleCancel
  };
}
const _hoisted_1 = {
  class: "modal-content max-w-sm",
  role: "alertdialog",
  "aria-modal": "true",
  "aria-labelledby": "confirm-title"
};
const _hoisted_2 = { class: "flex justify-center mb-4" };
const _hoisted_3 = {
  id: "confirm-title",
  class: "text-lg font-semibold text-white text-center mb-2"
};
const _hoisted_4 = { class: "text-sm text-gray-400 text-center mb-6" };
const _hoisted_5 = { class: "flex gap-3" };
const _sfc_main = /* @__PURE__ */ defineComponent({
  __name: "ConfirmDialog",
  setup(__props) {
    const { t } = useI18n();
    const { isOpen: isOpen2, options: options2, handleConfirm, handleCancel } = useConfirm();
    const cancelBtnRef = ref(null);
    watch(isOpen2, (open) => {
      if (open) {
        nextTick(() => cancelBtnRef.value?.focus());
      }
    });
    const variant = computed(() => options2.value?.variant || "info");
    const iconBgClass = computed(() => {
      switch (variant.value) {
        case "danger":
          return "bg-red-500/20";
        case "warning":
          return "bg-yellow-500/20";
        default:
          return "bg-blue-500/20";
      }
    });
    const iconClass = computed(() => {
      switch (variant.value) {
        case "danger":
          return "text-red-400";
        case "warning":
          return "text-yellow-400";
        default:
          return "text-blue-400";
      }
    });
    const confirmBtnClass = computed(() => {
      switch (variant.value) {
        case "danger":
          return "btn-unified-danger";
        case "warning":
          return "btn-unified-warning";
        default:
          return "btn-unified-primary";
      }
    });
    return (_ctx, _cache) => {
      return openBlock(), createBlock(Teleport, { to: "body" }, [
        createVNode(Transition, { name: "modal" }, {
          default: withCtx(() => [
            unref(isOpen2) ? (openBlock(), createElementBlock("div", {
              key: 0,
              class: "modal-container",
              onClick: _cache[2] || (_cache[2] = withModifiers(
                //@ts-ignore
                (...args) => unref(handleCancel) && unref(handleCancel)(...args),
                ["self"]
              )),
              onKeydown: [
                _cache[3] || (_cache[3] = withKeys(
                  //@ts-ignore
                  (...args) => unref(handleCancel) && unref(handleCancel)(...args),
                  ["escape"]
                )),
                _cache[4] || (_cache[4] = withKeys(
                  //@ts-ignore
                  (...args) => unref(handleConfirm) && unref(handleConfirm)(...args),
                  ["enter"]
                ))
              ]
            }, [
              createBaseVNode("div", _hoisted_1, [
                createBaseVNode("div", _hoisted_2, [
                  createBaseVNode("div", {
                    class: normalizeClass(["w-12 h-12 rounded-full flex items-center justify-center", iconBgClass.value])
                  }, [
                    variant.value === "warning" || variant.value === "danger" ? (openBlock(), createBlock(unref(TriangleAlert), {
                      key: 0,
                      class: normalizeClass(["w-6 h-6", iconClass.value])
                    }, null, 8, ["class"])) : (openBlock(), createBlock(unref(Info), {
                      key: 1,
                      class: normalizeClass(["w-6 h-6", iconClass.value])
                    }, null, 8, ["class"]))
                  ], 2)
                ]),
                createBaseVNode("h3", _hoisted_3, toDisplayString(unref(options2)?.title), 1),
                createBaseVNode("p", _hoisted_4, toDisplayString(unref(options2)?.message), 1),
                createBaseVNode("div", _hoisted_5, [
                  createBaseVNode("button", {
                    ref_key: "cancelBtnRef",
                    ref: cancelBtnRef,
                    onClick: _cache[0] || (_cache[0] = //@ts-ignore
                    (...args) => unref(handleCancel) && unref(handleCancel)(...args)),
                    class: "flex-1 btn-unified btn-unified-secondary"
                  }, toDisplayString(unref(options2)?.cancelText || unref(t)("common.cancel")), 513),
                  createBaseVNode("button", {
                    onClick: _cache[1] || (_cache[1] = //@ts-ignore
                    (...args) => unref(handleConfirm) && unref(handleConfirm)(...args)),
                    class: normalizeClass(["flex-1 btn-unified", confirmBtnClass.value])
                  }, toDisplayString(unref(options2)?.confirmText || unref(t)("common.confirm")), 3)
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
export {
  _sfc_main as default
};
