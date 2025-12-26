const __vite__mapDeps=(i,m=__vite__mapDeps,d=(m.f||(m.f=["js/QuickLookModal-PV9OKgmz.js","js/vue-vendor-DLLjc1XI.js","assets/vue-vendor-Dw2MutP3.css","js/highlighter-D41BBaGQ.js","assets/highlighter-DqdTtQ31.css","js/ui-vendor-2hF-XK7S.js","js/utils-DVYXdora.js","js/icons-BGgRd9bo.js","assets/QuickLookModal-BhcQKzkB.css","js/IgnoreRulesModal-BEVOVaU3.js","assets/IgnoreRulesModal-CZqV1y2E.css","js/useApiCache-BTgc1mkV.js","assets/driver-DB0Q8XAf.css","assets/driver-theme-tWMupVSo.css","js/CommandPalette-Be47sfdv.js","js/KeyboardShortcutsModal-CJZUYkBp.js","js/MemoryDashboard-BykbtL8H.js","js/SettingsModal-D_zm6NLM.js","assets/SettingsModal-Ba0Gm-ox.css","js/ConfirmDialog-DL6BPwsV.js"])))=>i.map(i=>d[i]);
import { r as ref, c as computed, e as defineComponent, f as onMounted, y as createElementBlock, z as openBlock, l as normalizeClass, A as createCommentVNode, B as createVNode, u as unref, C as createBaseVNode, D as toDisplayString, E as Transition, G as withCtx, F as Fragment, H as renderList, I as defineStore, w as watch, s as shallowRef, t as triggerRef, g as onUnmounted, k as reactive, J as createBlock, T as Teleport, K as withModifiers, L as normalizeStyle, M as createTextVNode, q as toRef, N as withDirectives, O as withKeys, P as vModelText, Q as renderSlot, R as vModelCheckbox, h, S as resolveDynamicComponent, U as TransitionGroup, n as nextTick, V as script$2, W as defineAsyncComponent, X as storeToRefs, Y as vModelSelect, Z as createStaticVNode, _ as vShow, $ as isRef, a0 as createApp, a1 as createPinia } from "./vue-vendor-DLLjc1XI.js";
import { F as Fuse, u as useVirtualizer } from "./ui-vendor-2hF-XK7S.js";
import { o as onClickOutside, u as useStorage, a as useMagicKeys } from "./utils-DVYXdora.js";
import { C as Check, a as ChevronRight, S as Star, T as Trash2, F as FileText, X, b as Search, Z as Zap, U as User, P as Plus, c as Settings2, E as Eye, d as Upload, D as Download, e as Copy, f as Sparkles, g as Save, h as TriangleAlert, L as ListChecks, i as FolderTree, H as Hash, j as FileCode, k as Clipboard, l as Pencil, m as ChevronDown, n as EyeOff, I as Info, o as Lightbulb, p as LoaderCircle, M as MessageCircle, B as BookMarked, q as FileDown, r as Bot } from "./icons-BGgRd9bo.js";
import "./highlighter-D41BBaGQ.js";
(function polyfill() {
  const relList = document.createElement("link").relList;
  if (relList && relList.supports && relList.supports("modulepreload")) return;
  for (const link of document.querySelectorAll('link[rel="modulepreload"]')) processPreload(link);
  new MutationObserver((mutations) => {
    for (const mutation of mutations) {
      if (mutation.type !== "childList") continue;
      for (const node of mutation.addedNodes) if (node.tagName === "LINK" && node.rel === "modulepreload") processPreload(node);
    }
  }).observe(document, {
    childList: true,
    subtree: true
  });
  function getFetchOpts(link) {
    const fetchOpts = {};
    if (link.integrity) fetchOpts.integrity = link.integrity;
    if (link.referrerPolicy) fetchOpts.referrerPolicy = link.referrerPolicy;
    if (link.crossOrigin === "use-credentials") fetchOpts.credentials = "include";
    else if (link.crossOrigin === "anonymous") fetchOpts.credentials = "omit";
    else fetchOpts.credentials = "same-origin";
    return fetchOpts;
  }
  function processPreload(link) {
    if (link.ep) return;
    link.ep = true;
    const fetchOpts = getFetchOpts(link);
    fetch(link.href, fetchOpts);
  }
})();
const parents = /* @__PURE__ */ new Set();
const coords = /* @__PURE__ */ new WeakMap();
const siblings = /* @__PURE__ */ new WeakMap();
const animations = /* @__PURE__ */ new WeakMap();
const intersections = /* @__PURE__ */ new WeakMap();
const mutationObservers = /* @__PURE__ */ new WeakMap();
const intervals = /* @__PURE__ */ new WeakMap();
const options = /* @__PURE__ */ new WeakMap();
const debounces = /* @__PURE__ */ new WeakMap();
const enabled = /* @__PURE__ */ new WeakSet();
let root;
let scrollX = 0;
let scrollY = 0;
const TGT = "__aa_tgt";
const DEL = "__aa_del";
const NEW = "__aa_new";
const handleMutations = (mutations) => {
  const elements = getElements(mutations);
  if (elements) {
    elements.forEach((el) => animate(el));
  }
};
const handleResizes = (entries) => {
  entries.forEach((entry) => {
    if (entry.target === root)
      updateAllPos();
    if (coords.has(entry.target))
      updatePos(entry.target);
  });
};
function isOffscreen(el) {
  const rect = el.getBoundingClientRect();
  const vw = (root === null || root === void 0 ? void 0 : root.clientWidth) || 0;
  const vh = (root === null || root === void 0 ? void 0 : root.clientHeight) || 0;
  return rect.bottom < 0 || rect.top > vh || rect.right < 0 || rect.left > vw;
}
function observePosition(el) {
  const oldObserver = intersections.get(el);
  oldObserver === null || oldObserver === void 0 ? void 0 : oldObserver.disconnect();
  let rect = coords.get(el);
  let invocations = 0;
  const buffer = 5;
  if (!rect) {
    rect = getCoords(el);
    coords.set(el, rect);
  }
  const { offsetWidth, offsetHeight } = root;
  const rootMargins = [
    rect.top - buffer,
    offsetWidth - (rect.left + buffer + rect.width),
    offsetHeight - (rect.top + buffer + rect.height),
    rect.left - buffer
  ];
  const rootMargin = rootMargins.map((px) => `${-1 * Math.floor(px)}px`).join(" ");
  const observer = new IntersectionObserver(() => {
    ++invocations > 1 && updatePos(el);
  }, {
    root,
    threshold: 1,
    rootMargin
  });
  observer.observe(el);
  intersections.set(el, observer);
}
function updatePos(el, debounce2 = true) {
  clearTimeout(debounces.get(el));
  const optionsOrPlugin = getOptions(el);
  const delay = debounce2 ? isPlugin(optionsOrPlugin) ? 500 : optionsOrPlugin.duration : 0;
  debounces.set(el, setTimeout(async () => {
    const currentAnimation = animations.get(el);
    try {
      await (currentAnimation === null || currentAnimation === void 0 ? void 0 : currentAnimation.finished);
      coords.set(el, getCoords(el));
      observePosition(el);
    } catch {
    }
  }, delay));
}
function updateAllPos() {
  clearTimeout(debounces.get(root));
  debounces.set(root, setTimeout(() => {
    parents.forEach((parent) => forEach(parent, (el) => lowPriority(() => updatePos(el))));
  }, 100));
}
function poll(el) {
  setTimeout(() => {
    intervals.set(el, setInterval(() => lowPriority(updatePos.bind(null, el)), 2e3));
  }, Math.round(2e3 * Math.random()));
}
function lowPriority(callback) {
  if (typeof requestIdleCallback === "function") {
    requestIdleCallback(() => callback());
  } else {
    requestAnimationFrame(() => callback());
  }
}
let resize;
const supportedBrowser = typeof window !== "undefined" && "ResizeObserver" in window;
if (supportedBrowser) {
  root = document.documentElement;
  new MutationObserver(handleMutations);
  resize = new ResizeObserver(handleResizes);
  window.addEventListener("scroll", () => {
    scrollY = window.scrollY;
    scrollX = window.scrollX;
  });
  resize.observe(root);
}
function getElements(mutations) {
  const observedNodes = mutations.reduce((nodes, mutation) => {
    return [
      ...nodes,
      ...Array.from(mutation.addedNodes),
      ...Array.from(mutation.removedNodes)
    ];
  }, []);
  const onlyCommentNodesObserved = observedNodes.every((node) => node.nodeName === "#comment");
  if (onlyCommentNodesObserved)
    return false;
  return mutations.reduce((elements, mutation) => {
    if (elements === false)
      return false;
    if (mutation.target instanceof Element) {
      target(mutation.target);
      if (!elements.has(mutation.target)) {
        elements.add(mutation.target);
        for (let i = 0; i < mutation.target.children.length; i++) {
          const child = mutation.target.children.item(i);
          if (!child)
            continue;
          if (DEL in child) {
            return false;
          }
          target(mutation.target, child);
          elements.add(child);
        }
      }
      if (mutation.removedNodes.length) {
        for (let i = 0; i < mutation.removedNodes.length; i++) {
          const child = mutation.removedNodes[i];
          if (DEL in child) {
            return false;
          }
          if (child instanceof Element) {
            elements.add(child);
            target(mutation.target, child);
            siblings.set(child, [
              mutation.previousSibling,
              mutation.nextSibling
            ]);
          }
        }
      }
    }
    return elements;
  }, /* @__PURE__ */ new Set());
}
function target(el, child) {
  if (!child && !(TGT in el))
    Object.defineProperty(el, TGT, { value: el });
  else if (child && !(TGT in child))
    Object.defineProperty(child, TGT, { value: el });
}
function animate(el) {
  var _a, _b;
  const isMounted = el.isConnected;
  const preExisting = coords.has(el);
  if (isMounted && siblings.has(el))
    siblings.delete(el);
  if (((_a = animations.get(el)) === null || _a === void 0 ? void 0 : _a.playState) !== "finished") {
    (_b = animations.get(el)) === null || _b === void 0 ? void 0 : _b.cancel();
  }
  if (NEW in el) {
    add(el);
  } else if (preExisting && isMounted) {
    remain(el);
  } else if (preExisting && !isMounted) {
    remove(el);
  } else {
    add(el);
  }
}
function raw(str) {
  return Number(str.replace(/[^0-9.\-]/g, ""));
}
function getScrollOffset(el) {
  let p = el.parentElement;
  while (p) {
    if (p.scrollLeft || p.scrollTop) {
      return { x: p.scrollLeft, y: p.scrollTop };
    }
    p = p.parentElement;
  }
  return { x: 0, y: 0 };
}
function getCoords(el) {
  const rect = el.getBoundingClientRect();
  const { x, y } = getScrollOffset(el);
  return {
    top: rect.top + y,
    left: rect.left + x,
    width: rect.width,
    height: rect.height
  };
}
function getTransitionSizes(el, oldCoords, newCoords) {
  let widthFrom = oldCoords.width;
  let heightFrom = oldCoords.height;
  let widthTo = newCoords.width;
  let heightTo = newCoords.height;
  const styles = getComputedStyle(el);
  const sizing = styles.getPropertyValue("box-sizing");
  if (sizing === "content-box") {
    const paddingY = raw(styles.paddingTop) + raw(styles.paddingBottom) + raw(styles.borderTopWidth) + raw(styles.borderBottomWidth);
    const paddingX = raw(styles.paddingLeft) + raw(styles.paddingRight) + raw(styles.borderRightWidth) + raw(styles.borderLeftWidth);
    widthFrom -= paddingX;
    widthTo -= paddingX;
    heightFrom -= paddingY;
    heightTo -= paddingY;
  }
  return [widthFrom, widthTo, heightFrom, heightTo].map(Math.round);
}
function getOptions(el) {
  return TGT in el && options.has(el[TGT]) ? options.get(el[TGT]) : { duration: 250, easing: "ease-in-out" };
}
function getTarget(el) {
  if (TGT in el)
    return el[TGT];
  return void 0;
}
function isEnabled(el) {
  const target2 = getTarget(el);
  return target2 ? enabled.has(target2) : false;
}
function forEach(parent, ...callbacks) {
  callbacks.forEach((callback) => callback(parent, options.has(parent)));
  for (let i = 0; i < parent.children.length; i++) {
    const child = parent.children.item(i);
    if (child) {
      callbacks.forEach((callback) => callback(child, options.has(child)));
    }
  }
}
function getPluginTuple(pluginReturn) {
  if (Array.isArray(pluginReturn))
    return pluginReturn;
  return [pluginReturn];
}
function isPlugin(config) {
  return typeof config === "function";
}
function remain(el) {
  const oldCoords = coords.get(el);
  const newCoords = getCoords(el);
  if (!isEnabled(el))
    return coords.set(el, newCoords);
  if (isOffscreen(el)) {
    coords.set(el, newCoords);
    observePosition(el);
    return;
  }
  let animation;
  if (!oldCoords)
    return;
  const pluginOrOptions = getOptions(el);
  if (typeof pluginOrOptions !== "function") {
    let deltaLeft = oldCoords.left - newCoords.left;
    let deltaTop = oldCoords.top - newCoords.top;
    const deltaRight = oldCoords.left + oldCoords.width - (newCoords.left + newCoords.width);
    const deltaBottom = oldCoords.top + oldCoords.height - (newCoords.top + newCoords.height);
    if (deltaBottom == 0)
      deltaTop = 0;
    if (deltaRight == 0)
      deltaLeft = 0;
    const [widthFrom, widthTo, heightFrom, heightTo] = getTransitionSizes(el, oldCoords, newCoords);
    const start = {
      transform: `translate(${deltaLeft}px, ${deltaTop}px)`
    };
    const end = {
      transform: `translate(0, 0)`
    };
    if (widthFrom !== widthTo) {
      start.width = `${widthFrom}px`;
      end.width = `${widthTo}px`;
    }
    if (heightFrom !== heightTo) {
      start.height = `${heightFrom}px`;
      end.height = `${heightTo}px`;
    }
    animation = el.animate([start, end], {
      duration: pluginOrOptions.duration,
      easing: pluginOrOptions.easing
    });
  } else {
    const [keyframes] = getPluginTuple(pluginOrOptions(el, "remain", oldCoords, newCoords));
    animation = new Animation(keyframes);
    animation.play();
  }
  animations.set(el, animation);
  coords.set(el, newCoords);
  animation.addEventListener("finish", updatePos.bind(null, el, false), {
    once: true
  });
}
function add(el) {
  if (NEW in el)
    delete el[NEW];
  const newCoords = getCoords(el);
  coords.set(el, newCoords);
  const pluginOrOptions = getOptions(el);
  if (!isEnabled(el))
    return;
  if (isOffscreen(el)) {
    observePosition(el);
    return;
  }
  let animation;
  if (typeof pluginOrOptions !== "function") {
    animation = el.animate([
      { transform: "scale(.98)", opacity: 0 },
      { transform: "scale(0.98)", opacity: 0, offset: 0.5 },
      { transform: "scale(1)", opacity: 1 }
    ], {
      duration: pluginOrOptions.duration * 1.5,
      easing: "ease-in"
    });
  } else {
    const [keyframes] = getPluginTuple(pluginOrOptions(el, "add", newCoords));
    animation = new Animation(keyframes);
    animation.play();
  }
  animations.set(el, animation);
  animation.addEventListener("finish", updatePos.bind(null, el, false), {
    once: true
  });
}
function cleanUp(el, styles) {
  var _a;
  el.remove();
  coords.delete(el);
  siblings.delete(el);
  animations.delete(el);
  (_a = intersections.get(el)) === null || _a === void 0 ? void 0 : _a.disconnect();
  setTimeout(() => {
    if (DEL in el)
      delete el[DEL];
    Object.defineProperty(el, NEW, { value: true, configurable: true });
    if (styles && el instanceof HTMLElement) {
      for (const style in styles) {
        el.style[style] = "";
      }
    }
  }, 0);
}
function remove(el) {
  var _a;
  if (!siblings.has(el) || !coords.has(el))
    return;
  const [prev, next] = siblings.get(el);
  Object.defineProperty(el, DEL, { value: true, configurable: true });
  const finalX = window.scrollX;
  const finalY = window.scrollY;
  if (next && next.parentNode && next.parentNode instanceof Element) {
    next.parentNode.insertBefore(el, next);
  } else if (prev && prev.parentNode) {
    prev.parentNode.appendChild(el);
  } else {
    (_a = getTarget(el)) === null || _a === void 0 ? void 0 : _a.appendChild(el);
  }
  if (!isEnabled(el))
    return cleanUp(el);
  const [top, left, width, height] = deletePosition(el);
  const optionsOrPlugin = getOptions(el);
  const oldCoords = coords.get(el);
  if (finalX !== scrollX || finalY !== scrollY) {
    adjustScroll(el, finalX, finalY, optionsOrPlugin);
  }
  let animation;
  let styleReset = {
    position: "absolute",
    top: `${top}px`,
    left: `${left}px`,
    width: `${width}px`,
    height: `${height}px`,
    margin: "0",
    pointerEvents: "none",
    transformOrigin: "center",
    zIndex: "100"
  };
  if (!isPlugin(optionsOrPlugin)) {
    Object.assign(el.style, styleReset);
    animation = el.animate([
      {
        transform: "scale(1)",
        opacity: 1
      },
      {
        transform: "scale(.98)",
        opacity: 0
      }
    ], {
      duration: optionsOrPlugin.duration,
      easing: "ease-out"
    });
  } else {
    const [keyframes, options2] = getPluginTuple(optionsOrPlugin(el, "remove", oldCoords));
    if ((options2 === null || options2 === void 0 ? void 0 : options2.styleReset) !== false) {
      styleReset = (options2 === null || options2 === void 0 ? void 0 : options2.styleReset) || styleReset;
      Object.assign(el.style, styleReset);
    }
    animation = new Animation(keyframes);
    animation.play();
  }
  animations.set(el, animation);
  animation.addEventListener("finish", () => cleanUp(el, styleReset), {
    once: true
  });
}
function adjustScroll(el, finalX, finalY, optionsOrPlugin) {
  const scrollDeltaX = scrollX - finalX;
  const scrollDeltaY = scrollY - finalY;
  const scrollBefore = document.documentElement.style.scrollBehavior;
  const scrollBehavior = getComputedStyle(root).scrollBehavior;
  if (scrollBehavior === "smooth") {
    document.documentElement.style.scrollBehavior = "auto";
  }
  window.scrollTo(window.scrollX + scrollDeltaX, window.scrollY + scrollDeltaY);
  if (!el.parentElement)
    return;
  const parent = el.parentElement;
  let lastHeight = parent.clientHeight;
  let lastWidth = parent.clientWidth;
  const startScroll = performance.now();
  function smoothScroll() {
    requestAnimationFrame(() => {
      if (!isPlugin(optionsOrPlugin)) {
        const deltaY = lastHeight - parent.clientHeight;
        const deltaX = lastWidth - parent.clientWidth;
        if (startScroll + optionsOrPlugin.duration > performance.now()) {
          window.scrollTo({
            left: window.scrollX - deltaX,
            top: window.scrollY - deltaY
          });
          lastHeight = parent.clientHeight;
          lastWidth = parent.clientWidth;
          smoothScroll();
        } else {
          document.documentElement.style.scrollBehavior = scrollBefore;
        }
      }
    });
  }
  smoothScroll();
}
function deletePosition(el) {
  var _a;
  const oldCoords = coords.get(el);
  const [width, , height] = getTransitionSizes(el, oldCoords, getCoords(el));
  let offsetParent = el.parentElement;
  while (offsetParent && (getComputedStyle(offsetParent).position === "static" || offsetParent instanceof HTMLBodyElement)) {
    offsetParent = offsetParent.parentElement;
  }
  if (!offsetParent)
    offsetParent = document.body;
  const parentStyles = getComputedStyle(offsetParent);
  const parentCoords = !animations.has(el) || ((_a = animations.get(el)) === null || _a === void 0 ? void 0 : _a.playState) === "finished" ? getCoords(offsetParent) : coords.get(offsetParent);
  const top = Math.round(oldCoords.top - parentCoords.top) - raw(parentStyles.borderTopWidth);
  const left = Math.round(oldCoords.left - parentCoords.left) - raw(parentStyles.borderLeftWidth);
  return [top, left, width, height];
}
function autoAnimate(el, config = {}) {
  if (supportedBrowser && resize) {
    const mediaQuery = window.matchMedia("(prefers-reduced-motion: reduce)");
    const isDisabledDueToReduceMotion = mediaQuery.matches && !isPlugin(config) && !config.disrespectUserMotionPreference;
    if (!isDisabledDueToReduceMotion) {
      enabled.add(el);
      if (getComputedStyle(el).position === "static") {
        Object.assign(el.style, { position: "relative" });
      }
      forEach(el, updatePos, poll, (element) => resize === null || resize === void 0 ? void 0 : resize.observe(element));
      if (isPlugin(config)) {
        options.set(el, config);
      } else {
        options.set(el, {
          duration: 250,
          easing: "ease-in-out",
          ...config
        });
      }
      const mo = new MutationObserver(handleMutations);
      mo.observe(el, { childList: true });
      mutationObservers.set(el, mo);
      parents.add(el);
    }
  }
  const controller = Object.freeze({
    parent: el,
    enable: () => {
      enabled.add(el);
    },
    disable: () => {
      enabled.delete(el);
      forEach(el, (node) => {
        const a = animations.get(node);
        try {
          a === null || a === void 0 ? void 0 : a.cancel();
        } catch {
        }
        animations.delete(node);
        const d = debounces.get(node);
        if (d)
          clearTimeout(d);
        debounces.delete(node);
        const i = intervals.get(node);
        if (i)
          clearInterval(i);
        intervals.delete(node);
      });
    },
    isEnabled: () => enabled.has(el),
    destroy: () => {
      enabled.delete(el);
      parents.delete(el);
      options.delete(el);
      const mo = mutationObservers.get(el);
      mo === null || mo === void 0 ? void 0 : mo.disconnect();
      mutationObservers.delete(el);
      forEach(el, (node) => {
        resize === null || resize === void 0 ? void 0 : resize.unobserve(node);
        const a = animations.get(node);
        try {
          a === null || a === void 0 ? void 0 : a.cancel();
        } catch {
        }
        animations.delete(node);
        const io = intersections.get(node);
        io === null || io === void 0 ? void 0 : io.disconnect();
        intersections.delete(node);
        const i = intervals.get(node);
        if (i)
          clearInterval(i);
        intervals.delete(node);
        const d = debounces.get(node);
        if (d)
          clearTimeout(d);
        debounces.delete(node);
        coords.delete(node);
        siblings.delete(node);
      });
    }
  });
  return controller;
}
function createVAutoAnimate(defaults) {
  return {
    mounted(el, binding) {
      let resolved = {};
      const local = binding.value;
      if (typeof local === "function") {
        resolved = local;
      } else if (typeof defaults === "function") {
        resolved = defaults;
      } else {
        resolved = { ...defaults || {}, ...local || {} };
      }
      const ctl = autoAnimate(el, resolved);
      Object.defineProperty(el, "__aa_ctl", { value: ctl, configurable: true });
    },
    unmounted(el) {
      var _a;
      const ctl = el["__aa_ctl"];
      (_a = ctl === null || ctl === void 0 ? void 0 : ctl.destroy) === null || _a === void 0 ? void 0 : _a.call(ctl);
      try {
        delete el["__aa_ctl"];
      } catch {
      }
    },
    getSSRProps: () => ({})
  };
}
const autoAnimatePlugin = {
  install(app2, defaults) {
    app2.directive("auto-animate", createVAutoAnimate(defaults));
  }
};
const scriptRel = "modulepreload";
const assetsURL = function(dep) {
  return "/" + dep;
};
const seen = {};
const __vitePreload = function preload(baseModule, deps, importerUrl) {
  let promise = Promise.resolve();
  if (deps && deps.length > 0) {
    let allSettled = function(promises$2) {
      return Promise.all(promises$2.map((p) => Promise.resolve(p).then((value$1) => ({
        status: "fulfilled",
        value: value$1
      }), (reason) => ({
        status: "rejected",
        reason
      }))));
    };
    document.getElementsByTagName("link");
    const cspNonceMeta = document.querySelector("meta[property=csp-nonce]");
    const cspNonce = cspNonceMeta?.nonce || cspNonceMeta?.getAttribute("nonce");
    promise = allSettled(deps.map((dep) => {
      dep = assetsURL(dep);
      if (dep in seen) return;
      seen[dep] = true;
      const isCss = dep.endsWith(".css");
      const cssSelector = isCss ? '[rel="stylesheet"]' : "";
      if (document.querySelector(`link[href="${dep}"]${cssSelector}`)) return;
      const link = document.createElement("link");
      link.rel = isCss ? "stylesheet" : scriptRel;
      if (!isCss) link.as = "script";
      link.crossOrigin = "";
      link.href = dep;
      if (cspNonce) link.setAttribute("nonce", cspNonce);
      document.head.appendChild(link);
      if (isCss) return new Promise((res, rej) => {
        link.addEventListener("load", res);
        link.addEventListener("error", () => rej(/* @__PURE__ */ new Error(`Unable to preload CSS for ${dep}`)));
      });
    }));
  }
  function handlePreloadError(err$2) {
    const e$1 = new Event("vite:preloadError", { cancelable: true });
    e$1.payload = err$2;
    window.dispatchEvent(e$1);
    if (!e$1.defaultPrevented) throw err$2;
  }
  return promise.then((res) => {
    for (const item of res || []) {
      if (item.status !== "rejected") continue;
      handlePreloadError(item.reason);
    }
    return baseModule().catch(handlePreloadError);
  });
};
const chat$1 = {
  "chat.title": "AI Chat",
  "chat.placeholder": "Describe a task or ask a question...",
  "chat.send": "Send",
  "chat.clear": "Clear chat",
  "chat.thinking": "AI is thinking...",
  "chat.error": "Send error",
  "chat.retry": "Retry",
  "chat.copy": "Copy",
  "chat.useContext": "Use context",
  "chat.noContext": "no context",
  "chat.noContextHint": "Build context to improve AI responses",
  "chat.buildContextFirst": "Build context first",
  "chat.contextStack": "Context",
  "chat.smartModeHint": "AI will find relevant files",
  "chat.contextReady": "Context ready to use",
  "chat.suggestedFiles": "Suggested files",
  "chat.useSelectedFiles": "Use ({count})",
  "chat.placeholderWithContext": "Describe your task — context is loaded...",
  "chat.attachFiles": "Attach files",
  "chat.hints.mention": "mention",
  "chat.hints.command": "command",
  "chat.hints.newLine": "new line",
  "chat.mentions.title": "Mention",
  "chat.mentions.files": "Select files",
  "chat.mentions.git": "Git changes",
  "chat.mentions.problems": "Code problems",
  "chat.selectFilesHint": "Select files on the left or use @files",
  "chat.noApiKey": "API key not configured",
  "chat.configureApi": "Configure API",
  "chat.welcome": "Hi! I can help analyze your code.",
  "chat.welcome.title": "AI Assistant",
  "chat.welcome.subtitle": "Describe a task and I'll find the right files. Or select them via @files",
  "chat.welcome.connected": "Connected to {provider} ({model})",
  "chat.welcome.disconnected": "API not configured",
  "chat.welcome.tipMention": "mention files",
  "chat.welcome.tipCommand": "quick commands",
  "chat.welcome.tryAsking": "Try asking",
  "chat.starters.analyze": "Find bugs and issues in the code",
  "chat.starters.explain": "Explain how this works",
  "chat.starters.refactor": "Suggest improvements",
  "chat.starters.test": "Write tests for this",
  "chat.actions.analyze": "Analyze",
  "chat.actions.explain": "Explain",
  "chat.actions.refactor": "Refactor",
  "chat.prompts.analyze": "Analyze this code and find potential issues",
  "chat.prompts.explain": "Explain how this code works",
  "chat.prompts.refactor": "Suggest improvements for this code",
  "chat.prompts.test": "Write unit tests for this code",
  "chat.analysisFailed": "Analysis failed",
  "chat.contextBuildFailed": "Failed to build context",
  "chat.contextAttached": "Context attached",
  "chat.modeManual": "Manual",
  "chat.modeSmart": "Smart",
  "chat.modeManualHint": "Uses your selected context",
  "chat.modeSmartHint": "AI finds relevant files",
  "chat.analyzing": "Analyzing task...",
  "chat.foundFiles": "Found files",
  "chat.relevance": "relevance",
  "chat.useFiles": "Use",
  "chat.editFiles": "Edit",
  "chat.cancelAnalysis": "Cancel",
  "chat.noFilesFound": "No relevant files found",
  "chat.tryManualMode": "Try manual mode",
  "chat.estimatedTokens": "Estimated tokens",
  "chat.stop": "Stop",
  "chat.modeAgentic": "Agent",
  "chat.modeAgenticHint": "AI explores code using tools autonomously",
  "chat.copied": "Copied",
  "chat.copyFailed": "Copy failed",
  "chat.comingSoon": "Coming Soon",
  "chat.comingSoonTitle": "AI Chat - Coming Soon",
  "chat.comingSoonDesc": "This feature is under development and will be available in a future release.",
  "chat.plannedFeatures": "Planned features:",
  "chat.feature.realtime": "Real-time AI conversation",
  "chat.feature.contextAware": "Context-aware responses",
  "chat.feature.codeGen": "Code generation and explanation",
  "chat.feature.streaming": "Streaming responses",
  "chat.feature.history": "Chat history",
  "chat.typing": "AI is typing...",
  "chat.pendingChanges": "changes",
  "chat.reviewChanges": "Review Changes",
  "chat.toSend": "to send",
  "toolCalls.executing": "executing",
  "toolCalls.completed": "completed",
  "toolCalls.failed": "failed",
  "toolCalls.arguments": "Arguments",
  "toolCalls.result": "Result",
  "toolCalls.showMore": "Show more",
  "toolCalls.showLess": "Show less"
};
const commands$1 = {
  "commands.openProject": "Open Project",
  "commands.openProjectDesc": "Open a project folder",
  "commands.changeProject": "Change Project",
  "commands.changeProjectDesc": "Switch to a different project",
  "commands.closeProject": "Close Project",
  "commands.closeProjectDesc": "Close current project",
  "commands.projectClosed": "Project closed",
  "commands.searchFiles": "Search Files",
  "commands.searchFilesDesc": "Search for files in current project",
  "commands.analyzeCode": "Analyze Code",
  "commands.analyzeCodeDesc": "Run code analysis on project",
  "commands.generateCode": "Generate Code",
  "commands.generateCodeDesc": "Generate code from task description",
  "commands.reloadWindow": "Reload Window",
  "commands.reloadWindowDesc": "Reload the application",
  "commands.keyboardShortcuts": "Keyboard Shortcuts",
  "commands.keyboardShortcutsDesc": "Show keyboard shortcuts guide",
  "commands.newDocument": "New Document",
  "commands.newDocumentDesc": "Create a new document"
};
const common$1 = {
  "tools.title": "Tools",
  "tools.description": "Additional functions and settings",
  "tools.copyContext": "Copy Context",
  "tools.exportContext": "Export Context",
  "tools.clearContext": "Clear Context",
  "tools.contextStats": "Context Statistics",
  "tools.files": "Files",
  "tools.aiModel": "AI Model",
  "tools.selectedModel": "Selected model",
  "tools.aiChat": "AI Chat",
  "sidebar.stats": "Statistics",
  "sidebar.history": "History",
  "sidebar.exportSettings": "Export",
  "sidebar.prompts": "Prompts",
  "sidebar.chat": "AI Chat",
  "sidebar.settings": "AI Settings",
  "sidebar.aiConfig": "AI Configuration",
  "sidebar.memory": "Memory",
  "sidebar.toggle": "Tools Panel",
  "workspace.resetLayout": "Reset panel sizes",
  "workspace.layoutReset": "Panel sizes reset",
  "stats.title": "Context Statistics",
  "stats.clickToExpand": "Click for detailed statistics",
  "stats.byType": "By file type",
  "stats.byFolder": "By folder",
  "stats.totalFiles": "Total files",
  "stats.totalLines": "Total lines",
  "stats.totalTokens": "Total tokens",
  "stats.estimatedCost": "Estimated cost",
  "history.title": "Context History",
  "history.empty": "History is empty",
  "history.load": "Load",
  "history.delete": "Delete",
  "history.ago": "ago",
  "history.justNow": "just now",
  "history.minutesAgo": "min ago",
  "history.hoursAgo": "h ago",
  "history.daysAgo": "d ago",
  "tabs.files": "Files",
  "tabs.git": "Git",
  "tabs.contexts": "Contexts",
  "tabs.tools": "Tools",
  "tabs.chat": "AI Chat",
  "common.today": "Today",
  "common.yesterday": "Yesterday",
  "common.daysAgo": "days ago",
  "common.justNow": "just now",
  "common.cancel": "Cancel",
  "common.save": "Save",
  "common.confirm": "Confirm",
  "common.clear": "Clear",
  "common.close": "Close",
  "common.undo": "Undo",
  "common.redo": "Redo",
  "common.expand": "Expand",
  "common.collapse": "Collapse",
  "common.loading": "Loading...",
  "common.pleaseWait": "Please wait...",
  "common.criticalError": "Critical Error",
  "common.loadingProject": "Loading Project",
  "accessibility.skipToContent": "Skip to content",
  "actions.copy": "Copy",
  "actions.export": "Export"
};
const context$1 = {
  "context.builder": "Context Builder",
  "context.buildOptions": "Build Options",
  "context.reset": "Reset",
  "context.maxTokens": "Max Tokens",
  "context.default": "Default",
  "context.stripComments": "Strip comments",
  "context.stripCommentsTooltip": "Removes comments from code to reduce context size",
  "context.includeTests": "Include tests",
  "context.includeTestsTooltip": "Includes test files (*_test.*, *.spec.*, test_*) in context",
  "context.optimization": "Content Optimization",
  "context.excludeTests": "Exclude tests",
  "context.collapseEmptyLines": "Collapse empty lines",
  "context.stripLicense": "Strip licenses",
  "context.compactDataFiles": "Compact JSON/YAML",
  "context.trimWhitespace": "Trim trailing whitespace",
  "context.splitStrategy": "Split Strategy",
  "context.splitStrategyTooltip": "Strategy for splitting large contexts into parts",
  "context.semantic": "Semantic",
  "context.semanticHint": "recommended",
  "context.fixed": "Fixed",
  "context.fixedHint": "fixed blocks",
  "context.adaptive": "Adaptive",
  "context.adaptiveHint": "smart splitting",
  "context.outputFormat": "Output Format",
  "context.outputFormatTooltip": "Context output format",
  "context.filesSelected": "files selected",
  "context.estimated": "Estimated",
  "context.tokens": "tokens",
  "context.build": "Build Context",
  "context.building": "Building context...",
  "context.showOptions": "Show options",
  "context.hideOptions": "Hide options",
  "context.currentFormat": "Current format",
  "context.changeFormat": "Change output format",
  "context.preview": "Context Preview",
  "context.copy": "Copy",
  "context.export": "Export",
  "context.clear": "Clear",
  "context.files": "files",
  "context.lines": "lines",
  "context.chars": "chars",
  "context.loadMore": "Load more...",
  "context.loading": "Loading preview...",
  "context.statusAnalyzing": "Analyzing dependencies",
  "context.statusReading": "Reading files",
  "context.statusProcessing": "Processing content",
  "context.statusFormatting": "Formatting output",
  "context.statusFinalizing": "Finalizing",
  "context.notBuilt": "No context built yet",
  "context.filesBelow": "File contents below",
  "context.selectFiles": 'Select files and click "Build Context"',
  "context.saved": "Saved Contexts",
  "context.noSaved": "No saved contexts",
  "context.buildToSee": "Build a context to see it here",
  "context.delete": "Delete",
  "context.rename": "Rename",
  "context.duplicate": "Duplicate",
  "context.favorite": "Add to favorites",
  "context.unfavorite": "Remove from favorites",
  "context.copyToClipboard": "Copy",
  "context.search": "Search in context...",
  "context.searchResults": "found",
  "context.noResults": "Nothing found",
  "context.dropFiles": "Drop files here",
  "context.dropFilesHint": "Files will be added to context",
  "context.doubleClickToView": "Double-click to view",
  "context.settings": "Storage settings",
  "context.maxContexts": "Max contexts",
  "context.maxStorage": "Max storage (MB)",
  "context.autoCleanupDays": "Auto-delete after (days)",
  "context.autoCleanupOnLimit": "Delete old on limit",
  "context.confirmDelete": "Confirm deletion",
  "context.confirmDeleteMessage": 'Are you sure you want to delete context "{name}"?',
  "context.confirmDeleteMultiple": "Delete {count} contexts?",
  "context.cancel": "Cancel",
  "context.sortBy": "Sort by",
  "context.sortByDate": "By date",
  "context.sortByName": "By name",
  "context.sortBySize": "By size",
  "context.filterFavorites": "Favorites only",
  "context.selectAll": "Select all",
  "context.deselectAll": "Deselect all",
  "context.deleteSelected": "Delete selected",
  "context.renamed": "Context renamed",
  "context.duplicated": "Context duplicated",
  "context.deleted": "Context deleted",
  "context.merge": "Merge",
  "context.mergeSelected": "Merge selected",
  "context.merged": "Contexts merged",
  "context.mergedName": "Merged context",
  "context.selectToMerge": "Select 2+ contexts to merge",
  "context.rebuilt": "Context rebuilt",
  "context.dragHint": "Drag to reorder",
  "context.selectHint": "Select files",
  "context.chatHint": "Ask AI",
  "context.step1": "Select files in the left panel",
  "context.step2": 'Click "Build Context"',
  "context.selectFilesFirst": "Select files first",
  "context.selectFilesHint": "Select files in the tree above",
  "context.buildTooltip": "Build context from selected files",
  "context.shortcutCopy": "Ctrl+C",
  "context.shortcutMerge": "Ctrl+M",
  "context.deleteAction": "Delete",
  "context.recommendations": "Recommendations",
  "context.smartSuggestions": "Suggested Files",
  "context.analyzingFiles": "Analyzing files...",
  "context.addSelected": "Add Selected",
  "context.addAll": "Add All",
  "context.sourceGit": "Often changed together",
  "context.sourceArch": "Related architectural layer",
  "context.sourceSemantic": "Semantically related",
  "context.sourceGitShort": "Git History",
  "context.sourceArchShort": "Import",
  "context.sourceSemanticShort": "Semantic",
  "context.smartSuggestionsTitle": "Smart Suggestions",
  "context.smartSuggestionsSubtitle": "Files related to your selection",
  "context.addFiles": "Add",
  "context.aiRecommendations": "AI Recommendations",
  "context.foundFilesCount": "Found {count} files",
  "context.foundFilesAnalysis": "Found {count} files based on context analysis",
  "context.impactAnalysis": "Impact Analysis",
  "context.analyzingImpact": "Analyzing impact...",
  "context.affected": "affected",
  "context.riskScore": "Risk Score",
  "context.affectedFiles": "Affected Files",
  "context.relatedTests": "Related Tests",
  "context.noImpact": "No impact",
  "context.relatedFiles": "related",
  "context.dependentFiles": "dependents",
  "context.relatedFilesTitle": "Related Files",
  "context.relatedFilesSubtitle": "Recommended to add to context",
  "context.relatedFilesTooltip": "Files that often change together with selected. Click to add to context",
  "context.dependentFilesTooltip": "Files that depend on selected and may be affected by changes",
  "context.impactTitle": "Impact Analysis",
  "context.impactSubtitle": "Files that may be affected",
  "context.noRelatedFiles": "No recommendations for selected files",
  "context.noImpactFiles": "Changes won't affect other files",
  "context.riskLow": "Low risk",
  "context.riskMedium": "Medium risk",
  "context.riskHigh": "High risk",
  "context.directDep": "direct",
  "context.transitiveDep": "transitive",
  "context.savedContexts": "Saved Contexts",
  "context.searchContexts": "Search contexts...",
  "context.noSavedContexts": "No saved contexts",
  "context.saveCurrentContext": "Save Current Context",
  "context.saveContext": "Save Context",
  "context.topic": "Topic",
  "context.topicPlaceholder": "e.g., Auth refactoring",
  "context.summary": "Summary",
  "context.summaryPlaceholder": "Brief description of the context...",
  "context.contextSaved": "Context saved",
  "context.contextRestored": "Context restored",
  "context.saveError": "Failed to save context",
  "context.deleteError": "Failed to delete context",
  "context.selectToPreview": "Select a context",
  "context.selectToPreviewHint": "Click on a context on the left to view details",
  "context.contextFiles": "Files in context",
  "context.noFilesInfo": "File information unavailable",
  "context.filesListUnavailable": "File list unavailable for this context",
  "context.filesInContext": "files in context",
  "context.loadToEditor": "Load to editor",
  "context.loadAndBuild": "Load & Build",
  "context.untitled": "Untitled",
  "context.filesShort": "files",
  "context.attached": "Attached",
  "context.showFavorites": "Show favorites",
  "context.restoreSelection": "Restore file selection",
  "context.selectionRestored": "Selected {count} files from «{name}»",
  "context.noFilesToRestore": "No file information in this context",
  "context.working": "Working ✓",
  "context.pendingBackend": "Will work after backend update",
  "action.project": "Project",
  "action.changeProject": "Change Project",
  "action.model": "Model",
  "action.tokens": "Tokens",
  "action.cost": "Cost",
  "action.buildContext": "Build Context",
  "action.copy": "Copy",
  "action.export": "Export",
  "action.generateSolution": "Generate Solution",
  "commandBar.limit": "Limit",
  "commandBar.build": "BUILD",
  "commandBar.selected": "Selected",
  "commandBar.files": "files",
  "commandBar.custom": "Custom",
  "commandBar.customLimit": "Custom limit"
};
const errors$1 = {
  "toast.selectFiles": "Select files to build context",
  "toast.contextBuilt": "Context built successfully",
  "toast.contextError": "Error building context",
  "toast.contextCopied": "Context copied to clipboard",
  "toast.copyError": "Error copying context",
  "toast.buildFirst": "Build context first",
  "toast.contextCleared": "Context cleared",
  "toast.contextEmpty": "Context is already empty",
  "toast.solutionSoon": "Solution generation - coming soon",
  "toast.noFiles": "No files selected",
  "toast.refreshed": "File tree refreshed",
  "toast.refreshError": "Failed to refresh file tree",
  "error.generic": "An error occurred",
  "error.unknown": "Unknown error",
  "error.loadFailed": "Failed to load data",
  "error.saveFailed": "Failed to save",
  "error.networkError": "Network error",
  "error.notFound": "Not found",
  "error.invalidData": "Invalid data",
  "error.suggestions": "Suggestions",
  "error.checkFiles": "Check that files are within the project directory",
  "error.checkPaths": `Verify file paths don't contain ".." or absolute paths`,
  "error.tryRefresh": "Try refreshing the file tree",
  "error.emptyContext": "Context built but empty",
  "error.tokenLimitExceeded": "Context exceeds token limit: {actual}K (limit: {limit}K). Reduce file selection or increase limit in settings.",
  "error.tokenLimitGeneric": "Context exceeds token limit. Reduce file selection or increase limit in settings.",
  "error.analysisFailed": "Analysis failed",
  "error.notConfigured": "Not configured"
};
const exportModule$1 = {
  "export.title": "Export",
  "export.formatPlain": "Plain Text",
  "export.formatManifest": "With Manifest",
  "export.formatXml": "XML",
  "export.formatJson": "JSON",
  "export.modelGpt35": "GPT-3.5 Turbo",
  "export.modelClaude": "Claude",
  "export.modelGemini": "Gemini Pro",
  "export.section.format": "Output Format",
  "export.section.output": "Output Options",
  "export.section.optimization": "Optimization",
  "export.section.chunking": "Chunking",
  "export.section.tokenLimit": "Token Limit",
  "export.section.model": "Model Settings",
  "export.includeManifest": "Include metadata",
  "export.stripComments": "Strip comments",
  "export.excludeTests": "Exclude tests",
  "export.stripLicense": "Strip licenses",
  "export.compactDataFiles": "Compact JSON/YAML",
  "export.trimWhitespace": "Trim whitespace",
  "export.collapseEmptyLines": "Collapse empty lines",
  "export.enableAutoSplit": "Auto-split",
  "export.maxTokensPerChunk": "Tokens per chunk",
  "export.splitStrategy": "Strategy",
  "export.chunking": "CHUNKING",
  "export.chunkSize": "Chunk size",
  "export.strategy.smart": "Smart",
  "export.strategy.smartDesc": "Keeps files intact, safer",
  "export.strategy.hard": "Hard",
  "export.strategy.hardDesc": "Strict limit, more efficient",
  "export.applyTemplate": "Apply template",
  "export.hint.applyTemplate": "Adds AI role, rules and structure from selected prompt template",
  "export.hint.includeMetadata": "Adds header with project info, file list and statistics",
  "export.hint.includeLineNumbers": "Adds line numbers to each line of code for easy navigation",
  "export.hint.stripComments": "Removes comments from code to save tokens. May lose important context",
  "export.hint.excludeTests": "Excludes test files (*_test.*, *.spec.*, test/, tests/)",
  "export.hint.stripLicense": "Removes license headers from files (MIT, Apache, GPL, etc.)",
  "export.hint.compactDataFiles": "Minifies JSON and YAML files by removing formatting",
  "export.hint.trimWhitespace": "Removes trailing whitespace from lines",
  "export.hint.collapseEmptyLines": "Replaces multiple empty lines with a single one",
  "export.hint.autoSplit": "Automatically splits large context into parts for copying",
  "export.group.cleanup": "Cleanup",
  "export.group.compression": "Compression",
  "export.totalSavings": "Savings",
  "export.chunkPreview": "Split preview",
  "export.chunksEstimated": "chunks",
  "export.tokensAvg": "tokens avg",
  "export.each": "each"
};
const files$1 = {
  "files.title": "Project Files",
  "files.selected": "selected",
  "files.refresh": "Refresh file tree",
  "files.search": "Search files...",
  "files.searchShort": "Search... (Ctrl+F)",
  "files.results": "results",
  "files.noFiles": "No files loaded",
  "files.noSearchResults": "No results found",
  "files.clearSearch": "Clear search",
  "files.clear": "Clear",
  "files.expandAll": "Expand all",
  "files.collapseAll": "Collapse all",
  "files.loading": "Loading...",
  "files.binaryFile": "Binary file",
  "files.cannotPreview": "Cannot preview content",
  "files.addToContext": "Add to context",
  "files.quickLook": "Quick Look",
  "files.settings": "Settings",
  "files.copyPath": "Copy full path",
  "files.openInExplorer": "Open in explorer",
  "files.fileCountTooltip": "Files in folder: {count}",
  "files.selectedTokensTooltip": "Selected tokens: {count}",
  "files.tokenCountTooltip": "File size: {count} tokens",
  "files.undoSelection": "Undo selection (Ctrl+Z)",
  "files.redoSelection": "Redo selection (Ctrl+Y)",
  "files.clearSelection": "Clear selection",
  "files.clearConfirmTitle": "Clear selection?",
  "files.clearConfirmMessage": "Are you sure you want to deselect {count} files?",
  "files.undo": "Undo selection",
  "files.redo": "Redo selection",
  "files.quickInfo": "Quick Info",
  "files.symbols": "Symbols",
  "files.imports": "Imports",
  "files.dependents": "Dependents",
  "files.changeRisk": "Change Risk",
  "files.risk.low": "Low",
  "files.risk.medium": "Medium",
  "files.risk.high": "High",
  "filter.byType": "Filter by type",
  "filter.types": "types",
  "filter.type": "type",
  "filter.searchExtensions": "Search extensions...",
  "filter.selectAll": "Select All",
  "filter.clear": "Clear",
  "filter.code": "Code",
  "filter.styles": "Styles",
  "filter.config": "Config",
  "filter.documentation": "Documentation",
  "filter.other": "Other",
  "contextMenu.selectAll": "Select All in Folder",
  "contextMenu.deselectAll": "Deselect All in Folder",
  "contextMenu.copyPath": "Copy Path",
  "contextMenu.copyRelativePath": "Copy Relative Path",
  "contextMenu.addToIgnore": "Add to Ignore",
  "contextMenu.removeFromIgnore": "Remove from Ignore",
  "contextMenu.expandAll": "Expand All",
  "contextMenu.collapseAll": "Collapse All",
  "quick.sourceFiles": "Source Files",
  "quick.tests": "Tests",
  "quick.config": "Config",
  "quick.docs": "Docs",
  "quick.styles": "Styles",
  "quick.modified": "Recently Modified",
  "quick.clear": "Clear",
  "quickFilters.label": "Filters",
  "quickFilters.filter": "Filter",
  "quickFilters.files": "files",
  "quickFilters.filesCount": "files",
  "quickFilters.clickToToggle": "Click to toggle",
  "quickFilters.clearFilter": "Clear filter",
  "quickFilters.settings": "Configure filters",
  "quickFilters.settingsTitle": "Filter Settings",
  "quickFilters.staticFilters": "Custom Filters",
  "quickFilters.dynamicFilters": "Auto-detected Filters",
  "quickFilters.dynamicHint": "based on project languages",
  "quickFilters.autoDetected": "Auto-detected",
  "quickFilters.extensions": "Extensions (comma-separated)",
  "quickFilters.patterns": "Patterns (comma-separated)",
  "quickFilters.reset": "Reset",
  "quickFilters.done": "Done",
  "quickFilters.types": "Types",
  "quickFilters.languages": "Languages",
  "quickFilters.fileTypes": "File Types",
  "quickFilters.projectLanguages": "Project Languages",
  "quickFilters.clearAll": "Clear",
  "quickFilters.clear": "Clear",
  "quickFilters.clickMultiple": "click for multiple selection",
  "quickFilters.customFilters": "Custom Filters",
  "quickFilters.extensionsPlaceholder": ".ts, .js, .vue",
  "quickFilters.smart": "Smart",
  "quickFilters.smartFilters": "Smart Filters",
  "quickFilters.basedOnFramework": "Based on framework",
  "quickFilters.multiSelect": "multi-select",
  "quickFilters.exclude": "exclude",
  "quickFilters.filterByType": "Filter by file type",
  "quickFilters.filterByLanguage": "Filter by programming language",
  "smartFilters.components": "Components",
  "smartFilters.composables": "Composables",
  "smartFilters.stores": "Stores",
  "smartFilters.views": "Views/Pages",
  "smartFilters.hooks": "Hooks",
  "smartFilters.pages": "Pages",
  "smartFilters.api": "API Routes",
  "smartFilters.services": "Services",
  "smartFilters.modules": "Modules",
  "smartFilters.handlers": "Handlers",
  "smartFilters.domain": "Domain",
  "smartFilters.infra": "Infrastructure",
  "smartFilters.backend": "Backend",
  "smartFilters.frontend": "Frontend",
  "smartFilters.models": "Models",
  "smartFilters.urls": "URLs",
  "smartFilters.routes": "Routes",
  "smartFilters.controllers": "Controllers",
  "smartFilters.repos": "Repositories",
  "smartFilters.screens": "Screens",
  "smartFilters.widgets": "Widgets",
  "smartFilters.state": "State/BLoC",
  "presets.title": "Selection Presets",
  "presets.subtitle": "Save and load file selections",
  "presets.yourPresets": "Your Presets",
  "presets.saveNew": "Save Current Selection",
  "presets.name": "Preset name...",
  "presets.description": "Description (optional)...",
  "presets.save": "Save Preset ({count} files)",
  "presets.saved": "Saved Presets",
  "presets.noPresets": "No presets saved yet",
  "presets.emptyHint": "Save your current selection to reuse it later!",
  "presets.saveCurrentShort": "Save Selection",
  "presets.manageAll": "Manage Presets",
  "presets.newPreset": "New Preset",
  "presets.savePresetTitle": "Save Preset",
  "presets.savePresetHint": "Save your current file selection",
  "presets.nameLabel": "Name",
  "presets.namePlaceholder": "e.g. Backend API",
  "presets.descriptionLabel": "Description",
  "presets.descriptionPlaceholder": "Optional",
  "presets.saveWithCount": "Save ({count} files)",
  "presets.files": "files",
  "presets.load": "Load",
  "presets.delete": "Delete",
  "presets.close": "Close",
  "presets.loaded": "Preset loaded",
  "presets.deleted": "Preset deleted",
  "presets.saveSuccess": 'Preset "{name}" saved',
  "presets.edit": "Edit",
  "presets.cancelEdit": "Cancel",
  "presets.saveEdit": "Save Changes",
  "presets.confirmDelete": "Confirm Deletion",
  "presets.confirmDeleteMessage": 'Are you sure you want to delete preset "{name}"?',
  "presets.cancel": "Cancel",
  "presets.showFiles": "Show Files",
  "presets.hideFiles": "Hide Files",
  "presets.invalidFiles": "{count} files do not exist",
  "presets.allFilesValid": "All files exist",
  "presets.editSuccess": "Preset updated",
  "presets.editFailed": "Failed to update preset",
  "presets.filesInPreset": "Files in preset:",
  "ignoreModal.title": "Ignore Rules Management",
  "ignoreModal.subtitle": "Manage .gitignore and custom ignore rules",
  "ignoreModal.customRules": "Custom Rules",
  "ignoreModal.gitignoreInfo": "This file is managed by Git. Edit it directly in your project.",
  "ignoreModal.gitignorePlaceholder": "No .gitignore file found",
  "ignoreModal.rulesCount": "Rules",
  "ignoreModal.customPlaceholder": "# Example:\n*.log\ntemp/\nnode_modules/",
  "ignoreModal.preview": "Preview",
  "ignoreModal.filesWillBeIgnored": "files will be ignored",
  "ignoreModal.andMore": "and {count} more",
  "ignoreModal.reset": "Reset",
  "ignoreModal.clearAll": "Clear All",
  "ignoreModal.cancel": "Cancel",
  "ignoreModal.save": "Save",
  "ignoreModal.previewEmpty": "Enter patterns to preview",
  "ignoreModal.saveSuccess": "Rules saved successfully",
  "tree.partialSelection": "{selected} of {total} files selected",
  "tree.filesSelected": "{count} files",
  "tree.expandHint": "Click folders to expand, or use buttons below",
  "tree.doubleClickHint": "Double-click to expand all",
  "filterLabels.source": "Code",
  "filterLabels.tests": "Tests",
  "filterLabels.config": "Config",
  "filterLabels.docs": "Docs",
  "filterLabels.styles": "Styles"
};
const git$1 = {
  "git.title": "Git Source",
  "git.localGit": "Local Git",
  "git.remoteUrl": "Remote URL",
  "git.currentBranch": "Current branch",
  "git.branches": "Branches",
  "git.commits": "Commits",
  "git.buildingFrom": "Building context from",
  "git.clear": "Clear",
  "git.filesAtRef": "Files at",
  "git.selectAll": "Select All",
  "git.clearSelection": "Clear",
  "git.buildContext": "Build Context from",
  "git.repoUrl": "Repository URL",
  "git.urlPlaceholder": "https://github.com/user/repo.git",
  "git.urlHint": "Supports GitHub, GitLab, Bitbucket and any public Git URL",
  "git.clone": "Clone Repository",
  "git.cloning": "Cloning...",
  "git.cloned": "Repository cloned",
  "git.openProject": "Open Project",
  "git.remove": "Remove",
  "git.notGitRepo": "Current project is not a Git repository",
  "git.openGitProject": "Open a project with .git folder or use Remote URL tab",
  "git.current": "current",
  "git.selected": "selected",
  "git.filesSelected": "files selected",
  "git.noCheckout": "Files will be read from this ref without checkout",
  "git.load": "Load",
  "git.githubApi": "GitHub API (fast, no clone)",
  "git.searchFiles": "Search files...",
  "git.recent": "Recent",
  "git.recentRepos": "Recent repositories",
  "git.clearHistory": "Clear",
  "git.cloneRequired": "URL not directly supported. Clone required:",
  "git.restoredSelection": "Restored file selection",
  "git.compareBranches": "Compare Branches",
  "git.baseBranch": "Base branch",
  "git.compareBranch": "Compare branch",
  "git.swapBranches": "Swap branches",
  "git.noDifferences": "No differences between branches",
  "git.filesChanged": "files changed",
  "git.close": "Close",
  "git.diff": "Diff",
  "gitContext.title": "Git Context",
  "gitContext.blame": "Git Blame",
  "gitContext.show": "Show Commit",
  "gitContext.changedFiles": "Changed Files",
  "gitContext.coChanged": "Co-changed Files",
  "gitContext.suggestContext": "Suggest Context",
  "gitContext.recentChanges": "Recent Changes",
  "gitContext.hotSpots": "Hot Spots",
  "gitContext.authors": "Authors",
  "gitContext.lastChanged": "Last Changed",
  "gitContext.changeCount": "Change Count",
  "gitContext.noChanges": "No recent changes",
  "gitContext.noCoChanged": "No co-changed files found",
  "gitContext.noSuggestions": "No suggestions based on history",
  "gitContext.selectFileHint": "Select a file to see co-changed files",
  "gitContext.taskPlaceholder": "Describe your task..."
};
const onboarding$1 = {
  "onboarding.step1Title": "Select files",
  "onboarding.step1Desc": "Check files and folders to include in context for AI",
  "onboarding.step2Title": "Build context",
  "onboarding.step2Desc": "Click this button to build context from selected files",
  "onboarding.step3Title": "Review result",
  "onboarding.step3Desc": "Built context is displayed here. You can copy or export it",
  "onboarding.step4Title": "Use AI",
  "onboarding.step4Desc": "Ask AI questions with your project context",
  "onboarding.next": "Next",
  "onboarding.prev": "Back",
  "onboarding.done": "Done",
  "onboarding.startTour": "Start tour",
  "onboarding.skip": "Skip"
};
const settings$1 = {
  "settings.aiProvider": "AI Provider",
  "settings.selectProvider": "Select provider",
  "settings.apiKey": "API Key",
  "settings.apiKeyPlaceholder": "Enter API key...",
  "settings.selectModel": "Model",
  "settings.hostUrl": "Host URL",
  "settings.save": "Save",
  "settings.saving": "Saving...",
  "settings.saved": "Settings saved",
  "settings.saveFailed": "Failed to save",
  "settings.showKey": "Show",
  "settings.hideKey": "Hide",
  "settings.clearKey": "Clear",
  "settings.qwenCliTitle": "Qwen Code CLI",
  "settings.qwenCliDescription": "Uses locally installed qwen-coder-cli. No API key required.",
  "settings.provider.openai": "GPT-4o, GPT-4, GPT-3.5",
  "settings.provider.gemini": "Gemini Pro, Gemini Flash",
  "settings.provider.qwen": "Qwen-Max, Qwen-Plus (API)",
  "settings.provider.qwenCli": "Local CLI, no API key",
  "settings.provider.openrouter": "Multiple providers",
  "settings.provider.localai": "Self-hosted models",
  "settings.hint.openai": "Get key at platform.openai.com",
  "settings.hint.gemini": "Get key at aistudio.google.com",
  "settings.hint.qwen": "Get key at dashscope.console.aliyun.com",
  "settings.hint.openrouter": "Get key at openrouter.ai",
  "settings.hint.localai": "Optional for local models",
  "settings.title": "File Explorer Settings",
  "settings.useGitignore": "Use .gitignore",
  "settings.useGitignoreHint": "Apply rules from .gitignore file",
  "settings.useCustomIgnore": "Use custom ignore rules",
  "settings.useCustomIgnoreHint": "Apply custom filtering rules",
  "settings.autoSaveSelection": "Auto-save selection",
  "settings.autoSaveSelectionHint": "Save file selection automatically",
  "settings.compactFolders": "Compact nested folders",
  "settings.compactFoldersHint": "Merge single-child folders",
  "settings.foldersFirst": "Folders first",
  "settings.allowSelectBinary": "Allow selecting binary files",
  "settings.ignoreRules": "Ignore Rules",
  "settings.manageRules": "Manage Rules",
  "settings.presets": "Presets",
  "settings.restoreSelection": "Restore previous selection ({count} files)",
  "settings.modal.title": "Settings",
  "settings.modal.general": "General",
  "settings.modal.ai": "AI",
  "settings.modal.export": "Export",
  "settings.modal.fileExplorer": "Explorer",
  "settings.modal.cancel": "Cancel",
  "settings.modal.language": "Language",
  "settings.modal.theme": "Theme",
  "settings.modal.themeDark": "Dark",
  "settings.modal.themeLight": "Light",
  "settings.modal.themeAuto": "Auto",
  "settings.modal.filterSettings": "File Filtering",
  "settings.modal.displaySettings": "Display",
  "settings.modal.tourHint": "Take a tour to learn all features",
  "export.format": "Format",
  "export.markdown": "Markdown",
  "export.plainText": "Plain text",
  "export.json": "JSON",
  "export.xml": "XML",
  "export.includeMetadata": "Include metadata",
  "export.includeLineNumbers": "Line numbers",
  "export.includeSeparators": "File separators",
  "export.wrapInCodeBlocks": "Wrap in code blocks",
  "export.maxLineLength": "Max line length",
  "export.unlimited": "Unlimited",
  "export.autoSplit": "Auto-split into chunks",
  "export.tokensPerChunk": "Tokens per chunk",
  "export.strategy": "Strategy",
  "export.strategySmart": "Smart",
  "export.strategyFile": "By files",
  "export.strategyToken": "By tokens",
  "export.tokenLimit": "Token limit",
  "export.custom": "Custom size",
  "chunks.parts": "parts",
  "chunks.part": "Part",
  "chunks.of": "of",
  "chunks.lines": "Lines",
  "chunks.copyNext": "Next",
  "chunks.copyPart": "Part {n}",
  "chunks.endPart": "END PART {n}",
  "chunks.startPart": "START PART {n}",
  "chunks.generate": "Split into parts",
  "chunks.copied": "copied",
  "chunks.allCopied": "All parts copied!",
  "chunks.reset": "Reset",
  "chunks.clickToCopy": "Click to copy",
  "chunks.exitMode": "Exit chunk mode",
  "settings.modal.system": "System",
  "settings.shellIntegration.title": "System Integration",
  "settings.shellIntegration.description": "Add 'Open in Syntaxia' to folder context menu",
  "settings.shellIntegration.enable": "Enable integration",
  "settings.shellIntegration.disable": "Disable integration",
  "settings.shellIntegration.enabled": "Integration enabled",
  "settings.shellIntegration.disabled": "Integration disabled",
  "settings.shellIntegration.enableSuccess": "Context menu added",
  "settings.shellIntegration.disableSuccess": "Context menu removed",
  "settings.shellIntegration.error": "Failed to change integration",
  "settings.shellIntegration.requiresAdmin": "May require explorer restart"
};
const task$1 = {
  "task.title": "Task Description",
  "task.placeholder": "Describe what you want to accomplish...\n\nExamples:\n• Refactor authentication service\n• Add dark mode support\n• Implement file upload feature",
  "task.quickSuggestions": "Quick suggestions:",
  "task.analysisResult": "Analysis Result:",
  "task.complexity": "Complexity",
  "task.estimatedFiles": "Estimated files",
  "task.suggestedApproach": "Suggested approach"
};
const templates$1 = {
  "templates.title": "Prompt Template",
  "templates.select": "Select template",
  "templates.selectTemplate": "Select prompt template",
  "templates.settings": "Template settings",
  "templates.task": "Task",
  "templates.taskPlaceholder": "Describe the task for AI...",
  "templates.sections": "Sections",
  "templates.suggestion": "Suggested:",
  "templates.history": "Task history",
  "templates.recentTasks": "Recent tasks",
  "templates.userRules": "Your rules",
  "templates.userRulesPlaceholder": "Additional rules for AI...",
  "templates.noRole": "No role defined",
  "templates.customPrefixSuffix": "Prefix / Suffix",
  "templates.prefix": "Prefix",
  "templates.suffix": "Suffix",
  "templates.prefixPlaceholder": "Text at the beginning...",
  "templates.suffixPlaceholder": "Text at the end...",
  "templates.additional": "Additional",
  "templates.manage": "Manage Templates",
  "templates.templates": "Templates",
  "templates.builtIn": "Built-in",
  "templates.custom": "Custom",
  "templates.favorites": "Favorites",
  "templates.create": "Create template",
  "templates.duplicate": "Duplicate",
  "templates.delete": "Delete",
  "templates.selectToEdit": "Select a template to edit",
  "templates.preview": "Preview",
  "templates.import": "Import",
  "templates.export": "Export",
  "templates.deleteConfirm": "Delete template?",
  "templates.deleteConfirmText": "This action cannot be undone.",
  "templates.unsavedChanges": "Unsaved changes",
  "templates.unsavedChangesText": "You have unsaved changes. Save them?",
  "templates.discard": "Discard",
  "templates.icon": "Icon",
  "templates.name": "Name",
  "templates.namePlaceholder": "Template name",
  "templates.description": "Description",
  "templates.descriptionPlaceholder": "Brief template description",
  "templates.roleContent": "AI Role",
  "templates.rolePlaceholder": "Describe the AI agent role...",
  "templates.rulesContent": "Rules",
  "templates.rulesPlaceholder": "Rules and constraints for AI...",
  "templates.section.role": "Role",
  "templates.section.rules": "Rules",
  "templates.section.tree": "File Tree",
  "templates.section.stats": "Stats",
  "templates.section.task": "Task",
  "templates.section.files": "Files",
  "templates.sectionHint.role": "AI agent role description",
  "templates.sectionHint.rules": "Rules and constraints",
  "templates.sectionHint.tree": "Project file structure",
  "templates.sectionHint.stats": "Files count, tokens, languages",
  "templates.sectionHint.task": "Your task description",
  "templates.sectionHint.files": "Selected files content",
  "templates.clearTask": "Clear task",
  "templates.clearHistory": "Clear history",
  "templates.currentTemplate": "Active template",
  "templates.search": "Search templates...",
  "templates.noResults": "No results found",
  "templates.showHidden": "Show hidden",
  "templates.hidden": "Hidden template",
  "templates.hide": "Hide",
  "templates.show": "Show",
  "templates.toggleFavorite": "Toggle favorite",
  "templates.apply": "Apply",
  "templates.builtInReadonly": "Built-in templates are read-only. Create a copy to edit.",
  "templates.noChanges": "No changes",
  "templates.reset": "Reset",
  "templates.resetToDefault": "Reset template settings",
  "templates.resetConfirm": "Reset settings?",
  "templates.resetConfirmText": "Favorite and hidden status will be reset.",
  "templates.autosave": "Autosave",
  "templates.autosaveHint": "Automatically save on close",
  "templates.sectionsHint": "Enable sections you need",
  "templates.enableSections": "Enable sections above",
  "templates.chars": "chars",
  "templates.includedSections": "Include in prompt",
  "templates.contextOptions": "Context options",
  "templates.previewEmpty": "Select a template to preview",
  "templates.saved": "Saved",
  "templates.deleted": "Template deleted",
  "templates.duplicated": "Template duplicated"
};
const welcome$1 = {
  "welcome.title": "Syntaxia",
  "welcome.subtitle": "Select a project to get started",
  "welcome.dragDrop": "or drag & drop a folder here",
  "welcome.openProject": "Open Project Directory",
  "welcome.recentProjects": "Recent Projects",
  "welcome.noRecentProjects": "No recent projects",
  "welcome.settings": "Settings",
  "welcome.autoOpen": "Automatically open last project on startup",
  "welcome.autoOpenShort": "Auto-open",
  "welcome.dropHere": "Drop folder here",
  "welcome.toOpenProject": "To open as project",
  "welcome.version": "Version",
  "welcome.removeProject": "Remove from list",
  "welcome.clearHistory": "Clear history",
  "welcome.confirmClearHistory": "Clear all project history?",
  "welcome.projectRemoved": "Project removed from list",
  "welcome.historyCleared": "History cleared",
  "welcome.openInExplorer": "Open in explorer",
  "welcome.copyPath": "Copy path",
  "welcome.pathCopied": "Path copied",
  "welcome.pinProject": "Pin",
  "welcome.unpinProject": "Unpin"
};
const en = {
  ...common$1,
  ...files$1,
  ...context$1,
  ...chat$1,
  ...git$1,
  ...settings$1,
  ...errors$1,
  ...templates$1,
  ...welcome$1,
  ...onboarding$1,
  ...exportModule$1,
  ...task$1,
  ...commands$1
};
const chat = {
  "chat.title": "AI Чат",
  "chat.placeholder": "Опишите задачу или задайте вопрос...",
  "chat.send": "Отправить",
  "chat.clear": "Очистить чат",
  "chat.thinking": "AI думает...",
  "chat.error": "Ошибка отправки",
  "chat.retry": "Повторить",
  "chat.copy": "Копировать",
  "chat.useContext": "Использовать контекст",
  "chat.noContext": "нет контекста",
  "chat.noContextHint": "Постройте контекст для улучшения ответов",
  "chat.buildContextFirst": "Сначала постройте контекст",
  "chat.contextStack": "Контекст",
  "chat.smartModeHint": "AI найдёт нужные файлы",
  "chat.contextReady": "Контекст готов к использованию",
  "chat.suggestedFiles": "Предложенные файлы",
  "chat.useSelectedFiles": "Использовать ({count})",
  "chat.placeholderWithContext": "Опишите задачу — контекст уже загружен...",
  "chat.attachFiles": "Прикрепить файлы",
  "chat.hints.mention": "упомянуть",
  "chat.hints.command": "команда",
  "chat.hints.newLine": "новая строка",
  "chat.mentions.title": "Упомянуть",
  "chat.mentions.files": "Выбрать файлы",
  "chat.mentions.git": "Git изменения",
  "chat.mentions.problems": "Проблемы кода",
  "chat.selectFilesHint": "Выберите файлы слева или используйте @files",
  "chat.noApiKey": "API ключ не настроен",
  "chat.configureApi": "Настроить API",
  "chat.welcome": "Привет! Я могу помочь с анализом кода.",
  "chat.welcome.title": "AI Ассистент",
  "chat.welcome.subtitle": "Опишите задачу, и я найду нужные файлы. Или выберите их через @files",
  "chat.welcome.connected": "Подключён к {provider} ({model})",
  "chat.welcome.disconnected": "API не настроен",
  "chat.welcome.tipMention": "упомянуть файлы",
  "chat.welcome.tipCommand": "быстрые команды",
  "chat.welcome.tryAsking": "Попробуйте спросить",
  "chat.starters.analyze": "Найди баги и проблемы в коде",
  "chat.starters.explain": "Объясни как это работает",
  "chat.starters.refactor": "Предложи улучшения",
  "chat.starters.test": "Напиши тесты для этого",
  "chat.actions.analyze": "Анализировать",
  "chat.actions.explain": "Объяснить",
  "chat.actions.refactor": "Рефакторинг",
  "chat.prompts.analyze": "Проанализируй этот код и найди потенциальные проблемы",
  "chat.prompts.explain": "Объясни как работает этот код",
  "chat.prompts.refactor": "Предложи улучшения для этого кода",
  "chat.prompts.test": "Напиши unit-тесты для этого кода",
  "chat.analysisFailed": "Ошибка анализа",
  "chat.contextBuildFailed": "Ошибка построения контекста",
  "chat.contextAttached": "Контекст прикреплён",
  "chat.modeManual": "Ручной",
  "chat.modeSmart": "Умный",
  "chat.modeManualHint": "Использует выбранный вами контекст",
  "chat.modeSmartHint": "AI сам найдёт нужные файлы",
  "chat.analyzing": "Анализирую задачу...",
  "chat.foundFiles": "Найдено файлов",
  "chat.relevance": "релевантность",
  "chat.useFiles": "Использовать",
  "chat.editFiles": "Редактировать",
  "chat.cancelAnalysis": "Отмена",
  "chat.noFilesFound": "Релевантные файлы не найдены",
  "chat.tryManualMode": "Попробуйте ручной режим",
  "chat.estimatedTokens": "Примерно токенов",
  "chat.stop": "Стоп",
  "chat.modeAgentic": "Агент",
  "chat.modeAgenticHint": "AI сам исследует код с помощью инструментов",
  "chat.copied": "Скопировано",
  "chat.copyFailed": "Ошибка копирования",
  "chat.comingSoon": "Скоро",
  "chat.comingSoonTitle": "AI Чат - Скоро",
  "chat.comingSoonDesc": "Эта функция находится в разработке и будет доступна в будущих версиях.",
  "chat.plannedFeatures": "Планируемые функции:",
  "chat.feature.realtime": "Общение с AI в реальном времени",
  "chat.feature.contextAware": "Ответы с учётом контекста",
  "chat.feature.codeGen": "Генерация и объяснение кода",
  "chat.feature.streaming": "Потоковые ответы",
  "chat.feature.history": "История чата",
  "chat.typing": "AI печатает...",
  "chat.pendingChanges": "изменений",
  "chat.reviewChanges": "Просмотр изменений",
  "chat.toSend": "отправить",
  "toolCalls.executing": "выполняется",
  "toolCalls.completed": "выполнено",
  "toolCalls.failed": "ошибка",
  "toolCalls.arguments": "Аргументы",
  "toolCalls.result": "Результат",
  "toolCalls.showMore": "Показать больше",
  "toolCalls.showLess": "Свернуть"
};
const commands = {
  "commands.openProject": "Открыть проект",
  "commands.openProjectDesc": "Открыть папку проекта",
  "commands.changeProject": "Сменить проект",
  "commands.changeProjectDesc": "Переключиться на другой проект",
  "commands.closeProject": "Закрыть проект",
  "commands.closeProjectDesc": "Закрыть текущий проект",
  "commands.projectClosed": "Проект закрыт",
  "commands.searchFiles": "Поиск файлов",
  "commands.searchFilesDesc": "Поиск файлов в текущем проекте",
  "commands.analyzeCode": "Анализ кода",
  "commands.analyzeCodeDesc": "Запустить анализ кода проекта",
  "commands.generateCode": "Генерация кода",
  "commands.generateCodeDesc": "Сгенерировать код по описанию задачи",
  "commands.reloadWindow": "Перезагрузить окно",
  "commands.reloadWindowDesc": "Перезагрузить приложение",
  "commands.keyboardShortcuts": "Горячие клавиши",
  "commands.keyboardShortcutsDesc": "Показать справку по горячим клавишам",
  "commands.newDocument": "Новый документ",
  "commands.newDocumentDesc": "Создать новый документ"
};
const common = {
  "tools.title": "Инструменты",
  "tools.description": "Дополнительные функции и настройки",
  "tools.copyContext": "Копировать контекст",
  "tools.exportContext": "Экспорт контекста",
  "tools.clearContext": "Очистить контекст",
  "tools.contextStats": "Статистика контекста",
  "tools.files": "Файлов",
  "tools.aiModel": "AI Модель",
  "tools.selectedModel": "Выбранная модель",
  "tools.aiChat": "AI Чат",
  "sidebar.stats": "Статистика",
  "sidebar.history": "История",
  "sidebar.exportSettings": "Экспорт",
  "sidebar.prompts": "Промпты",
  "sidebar.chat": "AI Чат",
  "sidebar.settings": "Настройки AI",
  "sidebar.aiConfig": "Конфигурация AI",
  "sidebar.memory": "Память",
  "sidebar.toggle": "Панель инструментов",
  "workspace.resetLayout": "Сбросить размеры панелей",
  "workspace.layoutReset": "Размеры панелей сброшены",
  "stats.title": "Статистика контекста",
  "stats.clickToExpand": "Нажмите для подробной статистики",
  "stats.byType": "По типам файлов",
  "stats.byFolder": "По папкам",
  "stats.totalFiles": "Всего файлов",
  "stats.totalLines": "Всего строк",
  "stats.totalTokens": "Всего токенов",
  "stats.estimatedCost": "Примерная стоимость",
  "history.title": "История контекстов",
  "history.empty": "История пуста",
  "history.load": "Загрузить",
  "history.delete": "Удалить",
  "history.ago": "назад",
  "history.justNow": "только что",
  "history.minutesAgo": "мин. назад",
  "history.hoursAgo": "ч. назад",
  "history.daysAgo": "дн. назад",
  "tabs.files": "Файлы",
  "tabs.git": "Git",
  "tabs.contexts": "Контексты",
  "tabs.tools": "Инструменты",
  "tabs.chat": "AI Чат",
  "common.today": "Сегодня",
  "common.yesterday": "Вчера",
  "common.daysAgo": "дн. назад",
  "common.justNow": "только что",
  "common.cancel": "Отмена",
  "common.save": "Сохранить",
  "common.confirm": "Подтвердить",
  "common.clear": "Очистить",
  "common.close": "Закрыть",
  "common.undo": "Отменить",
  "common.redo": "Повторить",
  "common.expand": "Развернуть",
  "common.collapse": "Свернуть",
  "common.loading": "Загрузка...",
  "common.pleaseWait": "Пожалуйста, подождите...",
  "common.criticalError": "Критическая ошибка",
  "common.loadingProject": "Загрузка проекта",
  "accessibility.skipToContent": "Перейти к содержимому",
  "actions.copy": "Копировать",
  "actions.export": "Экспорт"
};
const context = {
  "context.builder": "Построитель контекста",
  "context.buildOptions": "Параметры построения",
  "context.reset": "Сброс",
  "context.maxTokens": "Макс. токенов",
  "context.default": "По умолчанию",
  "context.stripComments": "Удалить комментарии",
  "context.stripCommentsTooltip": "Удаляет комментарии из кода для уменьшения размера контекста",
  "context.includeTests": "Включить тесты",
  "context.includeTestsTooltip": "Включает тестовые файлы (*_test.*, *.spec.*, test_*) в контекст",
  "context.optimization": "Оптимизация контента",
  "context.excludeTests": "Исключить тесты",
  "context.collapseEmptyLines": "Схлопнуть пустые строки",
  "context.stripLicense": "Удалить лицензии",
  "context.compactDataFiles": "Сжать JSON/YAML",
  "context.trimWhitespace": "Удалить trailing пробелы",
  "context.splitStrategy": "Стратегия разбиения",
  "context.splitStrategyTooltip": "Стратегия разбиения больших контекстов на части",
  "context.semantic": "Семантическая",
  "context.semanticHint": "рекомендуется",
  "context.fixed": "Фиксированная",
  "context.fixedHint": "фиксированные блоки",
  "context.adaptive": "Адаптивная",
  "context.adaptiveHint": "умное разбиение",
  "context.outputFormat": "Формат вывода",
  "context.outputFormatTooltip": "Формат вывода контекста",
  "context.filesSelected": "файлов выбрано",
  "context.estimated": "Примерно",
  "context.tokens": "токенов",
  "context.build": "Построить контекст",
  "context.building": "Построение контекста...",
  "context.showOptions": "Показать опции",
  "context.hideOptions": "Скрыть опции",
  "context.currentFormat": "Текущий формат",
  "context.changeFormat": "Изменить формат вывода",
  "context.preview": "Превью контекста",
  "context.copy": "Копировать",
  "context.export": "Экспорт",
  "context.clear": "Очистить",
  "context.files": "файлов",
  "context.lines": "строк",
  "context.chars": "симв.",
  "context.loadMore": "Загрузить ещё...",
  "context.loading": "Загрузка превью...",
  "context.statusAnalyzing": "Анализ зависимостей",
  "context.statusReading": "Чтение файлов",
  "context.statusProcessing": "Обработка контента",
  "context.statusFormatting": "Форматирование вывода",
  "context.statusFinalizing": "Финализация",
  "context.notBuilt": "Контекст ещё не построен",
  "context.filesBelow": "Содержимое файлов ниже",
  "context.selectFiles": 'Выберите файлы и нажмите "Построить контекст"',
  "context.saved": "Сохранённые контексты",
  "context.noSaved": "Нет сохранённых контекстов",
  "context.buildToSee": "Постройте контекст, чтобы увидеть его здесь",
  "context.delete": "Удалить",
  "context.rename": "Переименовать",
  "context.duplicate": "Дублировать",
  "context.favorite": "В избранное",
  "context.unfavorite": "Убрать из избранного",
  "context.copyToClipboard": "Копировать",
  "context.search": "Поиск в контексте...",
  "context.searchResults": "найдено",
  "context.noResults": "Ничего не найдено",
  "context.dropFiles": "Перетащите файлы сюда",
  "context.dropFilesHint": "Файлы будут добавлены в контекст",
  "context.doubleClickToView": "Двойной клик для просмотра",
  "context.settings": "Настройки хранения",
  "context.maxContexts": "Макс. контекстов",
  "context.maxStorage": "Макс. размер (MB)",
  "context.autoCleanupDays": "Автоудаление через (дней)",
  "context.autoCleanupOnLimit": "Удалять старые при превышении",
  "context.confirmDelete": "Подтвердите удаление",
  "context.confirmDeleteMessage": 'Вы уверены, что хотите удалить контекст "{name}"?',
  "context.confirmDeleteMultiple": "Удалить {count} контекстов?",
  "context.cancel": "Отмена",
  "context.sortBy": "Сортировка",
  "context.sortByDate": "По дате",
  "context.sortByName": "По имени",
  "context.sortBySize": "По размеру",
  "context.filterFavorites": "Только избранные",
  "context.selectAll": "Выбрать все",
  "context.deselectAll": "Снять выбор",
  "context.deleteSelected": "Удалить выбранные",
  "context.renamed": "Контекст переименован",
  "context.duplicated": "Контекст дублирован",
  "context.deleted": "Контекст удалён",
  "context.merge": "Объединить",
  "context.mergeSelected": "Объединить выбранные",
  "context.merged": "Контексты объединены",
  "context.mergedName": "Объединённый контекст",
  "context.selectToMerge": "Выберите 2+ контекста для объединения",
  "context.rebuilt": "Контекст пересобран",
  "context.dragHint": "Перетащите для изменения порядка",
  "context.selectHint": "Выберите файлы",
  "context.chatHint": "Спросите AI",
  "context.step1": "Выберите файлы в панели слева",
  "context.step2": 'Нажмите "Построить контекст"',
  "context.selectFilesFirst": "Сначала выберите файлы",
  "context.selectFilesHint": "Выберите файлы в дереве выше",
  "context.buildTooltip": "Построить контекст из выбранных файлов",
  "context.shortcutCopy": "Ctrl+C",
  "context.shortcutMerge": "Ctrl+M",
  "context.deleteAction": "Удалить",
  "context.recommendations": "Рекомендации",
  "context.smartSuggestions": "Рекомендуемые файлы",
  "context.analyzingFiles": "Анализ файлов...",
  "context.addSelected": "Добавить выбранные",
  "context.addAll": "Добавить все",
  "context.sourceGit": "Часто меняется вместе",
  "context.sourceArch": "Связанный архитектурный слой",
  "context.sourceSemantic": "Семантически похожий",
  "context.sourceGitShort": "Git History",
  "context.sourceArchShort": "Import",
  "context.sourceSemanticShort": "Semantic",
  "context.smartSuggestionsTitle": "Умные рекомендации",
  "context.smartSuggestionsSubtitle": "Файлы, связанные с вашим выбором",
  "context.addFiles": "Добавить",
  "context.aiRecommendations": "AI Рекомендации",
  "context.foundFilesCount": "Найдено {count} файлов",
  "context.foundFilesAnalysis": "Найдено {count} файлов на основе анализа контекста",
  "context.impactAnalysis": "Анализ влияния",
  "context.analyzingImpact": "Анализ влияния...",
  "context.affected": "затронуто",
  "context.riskScore": "Уровень риска",
  "context.affectedFiles": "Затронутые файлы",
  "context.relatedTests": "Связанные тесты",
  "context.noImpact": "Нет влияния",
  "context.relatedFiles": "связанных",
  "context.dependentFiles": "зависимых",
  "context.relatedFilesTitle": "Связанные файлы",
  "context.relatedFilesSubtitle": "Рекомендуем добавить в контекст",
  "context.relatedFilesTooltip": "Файлы, которые часто меняются вместе с выбранными. Кликните чтобы добавить в контекст",
  "context.dependentFilesTooltip": "Файлы, которые зависят от выбранных и могут быть затронуты изменениями",
  "context.impactTitle": "Анализ влияния",
  "context.impactSubtitle": "Файлы, которые могут быть затронуты",
  "context.noRelatedFiles": "Нет рекомендаций для выбранных файлов",
  "context.noImpactFiles": "Изменения не затронут другие файлы",
  "context.riskLow": "Низкий риск",
  "context.riskMedium": "Средний риск",
  "context.riskHigh": "Высокий риск",
  "context.directDep": "прямая",
  "context.transitiveDep": "транзитивная",
  "context.savedContexts": "Сохранённые контексты",
  "context.searchContexts": "Поиск по контекстам...",
  "context.noSavedContexts": "Нет сохранённых контекстов",
  "context.saveCurrentContext": "Сохранить текущий контекст",
  "context.saveContext": "Сохранить контекст",
  "context.topic": "Тема",
  "context.topicPlaceholder": "Например: Рефакторинг авторизации",
  "context.summary": "Описание",
  "context.summaryPlaceholder": "Краткое описание контекста...",
  "context.contextSaved": "Контекст сохранён",
  "context.contextRestored": "Контекст восстановлен",
  "context.saveError": "Ошибка сохранения контекста",
  "context.deleteError": "Ошибка удаления контекста",
  "context.selectToPreview": "Выберите контекст",
  "context.selectToPreviewHint": "Кликните на контекст слева, чтобы просмотреть детали",
  "context.contextFiles": "Файлы в контексте",
  "context.noFilesInfo": "Информация о файлах недоступна",
  "context.filesListUnavailable": "Список файлов недоступен для этого контекста",
  "context.filesInContext": "файлов в контексте",
  "context.loadToEditor": "Загрузить в редактор",
  "context.loadAndBuild": "Загрузить и построить",
  "context.untitled": "Без названия",
  "context.filesShort": "файлов",
  "context.attached": "Прикреплено",
  "context.showFavorites": "Показать избранные",
  "context.restoreSelection": "Восстановить выбор файлов",
  "context.selectionRestored": "Выбрано {count} файлов из «{name}»",
  "context.noFilesToRestore": "В этом контексте нет информации о файлах",
  "context.working": "Работает ✓",
  "context.pendingBackend": "Будет работать после обновления backend",
  "action.project": "Проект",
  "action.changeProject": "Сменить проект",
  "action.model": "Модель",
  "action.tokens": "Токены",
  "action.cost": "Стоимость",
  "action.buildContext": "Построить контекст",
  "action.copy": "Копировать",
  "action.export": "Экспорт",
  "action.generateSolution": "Сгенерировать решение",
  "commandBar.limit": "Лимит",
  "commandBar.build": "ПОСТРОИТЬ",
  "commandBar.selected": "Выбрано",
  "commandBar.files": "файлов",
  "commandBar.custom": "Свой",
  "commandBar.customLimit": "Свой лимит"
};
const errors = {
  "toast.selectFiles": "Выберите файлы для построения контекста",
  "toast.contextBuilt": "Контекст успешно построен",
  "toast.contextError": "Ошибка при построении контекста",
  "toast.contextCopied": "Контекст скопирован в буфер обмена",
  "toast.copyError": "Ошибка при копировании контекста",
  "toast.buildFirst": "Сначала постройте контекст",
  "toast.contextCleared": "Контекст очищен",
  "toast.contextEmpty": "Контекст уже пуст",
  "toast.solutionSoon": "Генерация решения - скоро будет доступна",
  "toast.noFiles": "Нет выбранных файлов",
  "toast.refreshed": "Дерево файлов обновлено",
  "toast.refreshError": "Ошибка при обновлении дерева файлов",
  "error.generic": "Произошла ошибка",
  "error.unknown": "Неизвестная ошибка",
  "error.loadFailed": "Не удалось загрузить данные",
  "error.saveFailed": "Не удалось сохранить",
  "error.networkError": "Ошибка сети",
  "error.notFound": "Не найдено",
  "error.invalidData": "Некорректные данные",
  "error.suggestions": "Рекомендации",
  "error.checkFiles": "Проверьте, что файлы находятся в директории проекта",
  "error.checkPaths": 'Убедитесь, что пути не содержат ".." или абсолютных путей',
  "error.tryRefresh": "Попробуйте обновить дерево файлов",
  "error.emptyContext": "Контекст построен, но пуст",
  "error.tokenLimitExceeded": "Контекст превышает лимит токенов: {actual}K (лимит: {limit}K). Уменьшите выбор файлов или увеличьте лимит в настройках.",
  "error.tokenLimitGeneric": "Контекст превышает лимит токенов. Уменьшите выбор файлов или увеличьте лимит в настройках.",
  "error.analysisFailed": "Ошибка анализа",
  "error.notConfigured": "Не настроено"
};
const exportModule = {
  "export.title": "Экспорт",
  "export.formatPlain": "Простой текст",
  "export.formatManifest": "С метаданными",
  "export.formatXml": "XML",
  "export.formatJson": "JSON",
  "export.modelGpt35": "GPT-3.5 Turbo",
  "export.modelClaude": "Claude",
  "export.modelGemini": "Gemini Pro",
  "export.section.format": "Формат вывода",
  "export.section.output": "Параметры вывода",
  "export.section.optimization": "Оптимизация",
  "export.section.chunking": "Разбиение на чанки",
  "export.section.tokenLimit": "Лимит токенов",
  "export.section.model": "Настройки модели",
  "export.includeManifest": "Включить метаданные",
  "export.stripComments": "Удалить комментарии",
  "export.excludeTests": "Исключить тесты",
  "export.stripLicense": "Удалить лицензии",
  "export.compactDataFiles": "Сжать JSON/YAML",
  "export.trimWhitespace": "Удалить пробелы",
  "export.collapseEmptyLines": "Убрать пустые строки",
  "export.enableAutoSplit": "Авто-разбиение",
  "export.maxTokensPerChunk": "Токенов на чанк",
  "export.splitStrategy": "Стратегия",
  "export.chunking": "РАЗБИЕНИЕ",
  "export.chunkSize": "Размер части",
  "export.strategy.smart": "Smart",
  "export.strategy.smartDesc": "Не разрывает файлы, безопаснее",
  "export.strategy.hard": "Hard",
  "export.strategy.hardDesc": "Строго по лимиту, экономичнее",
  "export.applyTemplate": "Применить шаблон",
  "export.hint.applyTemplate": "Добавляет роль AI, правила и структуру из выбранного шаблона промпта",
  "export.hint.includeMetadata": "Добавляет заголовок с информацией о проекте, списком файлов и статистикой",
  "export.hint.includeLineNumbers": "Добавляет номера строк к каждой строке кода для удобной навигации",
  "export.hint.stripComments": "Удаляет комментарии из кода для экономии токенов. Может потерять важный контекст",
  "export.hint.excludeTests": "Исключает тестовые файлы (*_test.*, *.spec.*, test/, tests/)",
  "export.hint.stripLicense": "Удаляет лицензионные заголовки из файлов (MIT, Apache, GPL и др.)",
  "export.hint.compactDataFiles": "Минифицирует JSON и YAML файлы, удаляя форматирование",
  "export.hint.trimWhitespace": "Удаляет пробелы в конце строк",
  "export.hint.collapseEmptyLines": "Заменяет множественные пустые строки на одну",
  "export.hint.autoSplit": "Автоматически разбивает большой контекст на части для копирования",
  "export.group.cleanup": "Очистка",
  "export.group.compression": "Сжатие",
  "export.totalSavings": "Экономия",
  "export.chunkPreview": "Превью разбиения",
  "export.chunksEstimated": "чанков",
  "export.tokensAvg": "токенов в среднем",
  "export.each": "каждый"
};
const files = {
  "files.title": "Файлы проекта",
  "files.selected": "выбрано",
  "files.refresh": "Обновить дерево файлов",
  "files.search": "Поиск файлов...",
  "files.searchShort": "Поиск... (Ctrl+F)",
  "files.results": "результатов",
  "files.noFiles": "Файлы не загружены",
  "files.noSearchResults": "Ничего не найдено",
  "files.clearSearch": "Сбросить поиск",
  "files.clear": "Очистить",
  "files.expandAll": "Развернуть всё",
  "files.collapseAll": "Свернуть всё",
  "files.loading": "Загрузка...",
  "files.binaryFile": "Бинарный файл",
  "files.cannotPreview": "Невозможно отобразить содержимое",
  "files.addToContext": "Добавить в контекст",
  "files.quickLook": "Быстрый просмотр",
  "files.settings": "Настройки",
  "files.copyPath": "Копировать путь",
  "files.openInExplorer": "Открыть в проводнике",
  "files.fileCountTooltip": "Файлов в папке: {count}",
  "files.selectedTokensTooltip": "Выбрано токенов: {count}",
  "files.tokenCountTooltip": "Размер файла: {count} токенов",
  "files.undoSelection": "Отменить выбор (Ctrl+Z)",
  "files.redoSelection": "Повторить выбор (Ctrl+Y)",
  "files.clearSelection": "Очистить выбор",
  "files.clearConfirmTitle": "Очистить выбор?",
  "files.clearConfirmMessage": "Вы уверены что хотите снять выбор с {count} файлов?",
  "files.undo": "Отменить выбор",
  "files.redo": "Повторить выбор",
  "files.quickInfo": "Информация",
  "files.symbols": "Символы",
  "files.imports": "Импорты",
  "files.dependents": "Зависимые",
  "files.changeRisk": "Риск изменений",
  "files.risk.low": "Низкий",
  "files.risk.medium": "Средний",
  "files.risk.high": "Высокий",
  "filter.byType": "Фильтр по типу",
  "filter.types": "типов",
  "filter.type": "тип",
  "filter.searchExtensions": "Поиск расширений...",
  "filter.selectAll": "Выбрать всё",
  "filter.clear": "Очистить",
  "filter.code": "Код",
  "filter.styles": "Стили",
  "filter.config": "Конфигурация",
  "filter.documentation": "Документация",
  "filter.other": "Другое",
  "contextMenu.selectAll": "Выбрать всё в папке",
  "contextMenu.deselectAll": "Снять выбор в папке",
  "contextMenu.copyPath": "Копировать путь",
  "contextMenu.copyRelativePath": "Копировать относительный путь",
  "contextMenu.addToIgnore": "Добавить в игнор",
  "contextMenu.removeFromIgnore": "Убрать из игнора",
  "contextMenu.expandAll": "Развернуть всё",
  "contextMenu.collapseAll": "Свернуть всё",
  "quick.sourceFiles": "Исходники",
  "quick.tests": "Тесты",
  "quick.config": "Конфигурация",
  "quick.docs": "Документация",
  "quick.styles": "Стили",
  "quick.modified": "Изменённые",
  "quick.clear": "Очистить",
  "quickFilters.label": "Фильтры",
  "quickFilters.filter": "Фильтр",
  "quickFilters.files": "файлов",
  "quickFilters.filesCount": "файлов",
  "quickFilters.clickToToggle": "Клик для переключения",
  "quickFilters.clearFilter": "Сбросить фильтр",
  "quickFilters.settings": "Настроить фильтры",
  "quickFilters.settingsTitle": "Настройка фильтров",
  "quickFilters.staticFilters": "Пользовательские фильтры",
  "quickFilters.dynamicFilters": "Автоматические фильтры",
  "quickFilters.dynamicHint": "на основе языков проекта",
  "quickFilters.autoDetected": "Автоопределено",
  "quickFilters.extensions": "Расширения (через запятую)",
  "quickFilters.patterns": "Паттерны (через запятую)",
  "quickFilters.reset": "Сбросить",
  "quickFilters.done": "Готово",
  "quickFilters.types": "Типы",
  "quickFilters.languages": "Языки",
  "quickFilters.fileTypes": "Типы файлов",
  "quickFilters.projectLanguages": "Языки проекта",
  "quickFilters.clearAll": "Очистить",
  "quickFilters.clear": "Сбросить",
  "quickFilters.clickMultiple": "клик для множественного выбора",
  "quickFilters.customFilters": "Пользовательские фильтры",
  "quickFilters.extensionsPlaceholder": ".ts, .js, .vue",
  "quickFilters.smart": "Умные",
  "quickFilters.smartFilters": "Умные фильтры",
  "quickFilters.basedOnFramework": "На основе фреймворка",
  "quickFilters.multiSelect": "мульти-выбор",
  "quickFilters.exclude": "исключить",
  "quickFilters.filterByType": "Фильтр по типу файлов",
  "quickFilters.filterByLanguage": "Фильтр по языку программирования",
  "smartFilters.components": "Компоненты",
  "smartFilters.composables": "Composables",
  "smartFilters.stores": "Stores",
  "smartFilters.views": "Views/Pages",
  "smartFilters.hooks": "Хуки",
  "smartFilters.pages": "Pages",
  "smartFilters.api": "API Routes",
  "smartFilters.services": "Сервисы",
  "smartFilters.modules": "Модули",
  "smartFilters.handlers": "Handlers",
  "smartFilters.domain": "Domain",
  "smartFilters.infra": "Infrastructure",
  "smartFilters.backend": "Backend",
  "smartFilters.frontend": "Frontend",
  "smartFilters.models": "Models",
  "smartFilters.urls": "URLs",
  "smartFilters.routes": "Routes",
  "smartFilters.controllers": "Controllers",
  "smartFilters.repos": "Repositories",
  "smartFilters.screens": "Screens",
  "smartFilters.widgets": "Widgets",
  "smartFilters.state": "State/BLoC",
  "presets.title": "Пресеты выбора",
  "presets.subtitle": "Сохранение и загрузка наборов файлов",
  "presets.yourPresets": "Ваши пресеты",
  "presets.saveNew": "Сохранить текущий выбор",
  "presets.name": "Название пресета...",
  "presets.description": "Описание (необязательно)...",
  "presets.save": "Сохранить пресет ({count} файлов)",
  "presets.saved": "Сохранённые пресеты",
  "presets.noPresets": "Пресеты ещё не сохранены",
  "presets.emptyHint": "Сохраните текущий выбор, чтобы использовать его в будущем!",
  "presets.saveCurrentShort": "Сохранить выбор",
  "presets.manageAll": "Управление пресетами",
  "presets.newPreset": "Новый пресет",
  "presets.savePresetTitle": "Сохранить пресет",
  "presets.savePresetHint": "Сохраните текущий выбор файлов",
  "presets.nameLabel": "Название",
  "presets.namePlaceholder": "Например: Backend API",
  "presets.descriptionLabel": "Описание",
  "presets.descriptionPlaceholder": "Необязательно",
  "presets.saveWithCount": "Сохранить ({count} файлов)",
  "presets.files": "файлов",
  "presets.load": "Загрузить",
  "presets.delete": "Удалить",
  "presets.close": "Закрыть",
  "presets.loaded": "Пресет загружен",
  "presets.deleted": "Пресет удалён",
  "presets.saveSuccess": 'Пресет "{name}" сохранён',
  "presets.edit": "Редактировать",
  "presets.cancelEdit": "Отмена",
  "presets.saveEdit": "Сохранить изменения",
  "presets.confirmDelete": "Подтвердите удаление",
  "presets.confirmDeleteMessage": 'Вы уверены, что хотите удалить пресет "{name}"?',
  "presets.cancel": "Отмена",
  "presets.showFiles": "Показать файлы",
  "presets.hideFiles": "Скрыть файлы",
  "presets.invalidFiles": "{count} файлов не существуют",
  "presets.allFilesValid": "Все файлы существуют",
  "presets.editSuccess": "Пресет обновлён",
  "presets.editFailed": "Не удалось обновить пресет",
  "presets.filesInPreset": "Файлы в пресете:",
  "ignoreModal.title": "Управление правилами игнорирования",
  "ignoreModal.subtitle": "Настройка .gitignore и пользовательских правил",
  "ignoreModal.customRules": "Пользовательские правила",
  "ignoreModal.gitignoreInfo": "Этот файл управляется Git. Редактируйте его напрямую в проекте.",
  "ignoreModal.gitignorePlaceholder": "Файл .gitignore не найден",
  "ignoreModal.rulesCount": "Правил",
  "ignoreModal.customPlaceholder": "# Пример:\n*.log\ntemp/\nnode_modules/",
  "ignoreModal.preview": "Предпросмотр",
  "ignoreModal.filesWillBeIgnored": "файлов будут игнорироваться",
  "ignoreModal.andMore": "и ещё {count}",
  "ignoreModal.reset": "Сбросить",
  "ignoreModal.clearAll": "Очистить всё",
  "ignoreModal.cancel": "Отмена",
  "ignoreModal.save": "Сохранить",
  "ignoreModal.previewEmpty": "Введите правила для предпросмотра",
  "ignoreModal.saveSuccess": "Правила успешно сохранены",
  "tree.partialSelection": "Выбрано {selected} из {total} файлов",
  "tree.filesSelected": "{count} файлов",
  "tree.expandHint": "Кликните на папки для раскрытия или используйте кнопки ниже",
  "tree.doubleClickHint": "Двойной клик для раскрытия всех",
  "filterLabels.source": "Код",
  "filterLabels.tests": "Тесты",
  "filterLabels.config": "Конфиг",
  "filterLabels.docs": "Доки",
  "filterLabels.styles": "Стили"
};
const git = {
  "git.title": "Git Источник",
  "git.localGit": "Локальный Git",
  "git.remoteUrl": "Удалённый URL",
  "git.currentBranch": "Текущая ветка",
  "git.branches": "Ветки",
  "git.commits": "Коммиты",
  "git.buildingFrom": "Сборка контекста из",
  "git.clear": "Очистить",
  "git.filesAtRef": "Файлы в",
  "git.selectAll": "Выбрать всё",
  "git.clearSelection": "Очистить",
  "git.buildContext": "Построить контекст из",
  "git.repoUrl": "URL репозитория",
  "git.urlPlaceholder": "https://github.com/user/repo.git",
  "git.urlHint": "Поддерживаются GitHub, GitLab, Bitbucket и любые публичные Git URL",
  "git.clone": "Клонировать репозиторий",
  "git.cloning": "Клонирование...",
  "git.cloned": "Репозиторий склонирован",
  "git.openProject": "Открыть проект",
  "git.remove": "Удалить",
  "git.notGitRepo": "Текущий проект не является Git репозиторием",
  "git.openGitProject": "Откройте проект с папкой .git или используйте вкладку Remote URL",
  "git.current": "текущая",
  "git.selected": "выбрана",
  "git.filesSelected": "файлов выбрано",
  "git.noCheckout": "Файлы будут прочитаны из этого ref без checkout",
  "git.load": "Загрузить",
  "git.githubApi": "GitHub API (быстро, без клонирования)",
  "git.searchFiles": "Поиск файлов...",
  "git.recent": "Недавние",
  "git.recentRepos": "Недавние репозитории",
  "git.clearHistory": "Очистить",
  "git.cloneRequired": "URL не поддерживается напрямую. Требуется клонирование:",
  "git.restoredSelection": "Восстановлен выбор файлов",
  "git.compareBranches": "Сравнение веток",
  "git.baseBranch": "Базовая ветка",
  "git.compareBranch": "Сравниваемая ветка",
  "git.swapBranches": "Поменять местами",
  "git.noDifferences": "Нет различий между ветками",
  "git.filesChanged": "файлов изменено",
  "git.close": "Закрыть",
  "git.diff": "Diff",
  "gitContext.title": "Git контекст",
  "gitContext.blame": "Git Blame",
  "gitContext.show": "Показать коммит",
  "gitContext.changedFiles": "Изменённые файлы",
  "gitContext.coChanged": "Связанные файлы",
  "gitContext.suggestContext": "Предложить контекст",
  "gitContext.recentChanges": "Недавние изменения",
  "gitContext.hotSpots": "Горячие точки",
  "gitContext.authors": "Авторы",
  "gitContext.lastChanged": "Последнее изменение",
  "gitContext.changeCount": "Количество изменений",
  "gitContext.noChanges": "Нет недавних изменений",
  "gitContext.noCoChanged": "Связанные файлы не найдены",
  "gitContext.noSuggestions": "Нет предложений на основе истории",
  "gitContext.selectFileHint": "Выберите файл для просмотра связанных",
  "gitContext.taskPlaceholder": "Опишите вашу задачу..."
};
const onboarding = {
  "onboarding.step1Title": "Выберите файлы",
  "onboarding.step1Desc": "Отметьте файлы и папки которые хотите включить в контекст для AI",
  "onboarding.step2Title": "Постройте контекст",
  "onboarding.step2Desc": "Нажмите эту кнопку чтобы собрать контекст из выбранных файлов",
  "onboarding.step3Title": "Просмотрите результат",
  "onboarding.step3Desc": "Здесь отображается собранный контекст. Можно копировать или экспортировать",
  "onboarding.step4Title": "Используйте AI",
  "onboarding.step4Desc": "Задавайте вопросы AI с учётом контекста вашего проекта",
  "onboarding.next": "Далее",
  "onboarding.prev": "Назад",
  "onboarding.done": "Готово",
  "onboarding.startTour": "Показать обучение",
  "onboarding.skip": "Пропустить"
};
const settings = {
  "settings.aiProvider": "AI Провайдер",
  "settings.selectProvider": "Выберите провайдера",
  "settings.apiKey": "API Ключ",
  "settings.apiKeyPlaceholder": "Введите API ключ...",
  "settings.selectModel": "Модель",
  "settings.hostUrl": "URL хоста",
  "settings.save": "Сохранить",
  "settings.saving": "Сохранение...",
  "settings.saved": "Настройки сохранены",
  "settings.saveFailed": "Ошибка сохранения",
  "settings.showKey": "Показать",
  "settings.hideKey": "Скрыть",
  "settings.clearKey": "Очистить",
  "settings.qwenCliTitle": "Qwen Code CLI",
  "settings.qwenCliDescription": "Использует локально установленный qwen-coder-cli. API ключ не требуется.",
  "settings.provider.openai": "GPT-4o, GPT-4, GPT-3.5",
  "settings.provider.gemini": "Gemini Pro, Gemini Flash",
  "settings.provider.qwen": "Qwen-Max, Qwen-Plus (API)",
  "settings.provider.qwenCli": "Локальный CLI, без API ключа",
  "settings.provider.openrouter": "Множество провайдеров",
  "settings.provider.localai": "Локальные модели",
  "settings.hint.openai": "Получить ключ на platform.openai.com",
  "settings.hint.gemini": "Получить ключ на aistudio.google.com",
  "settings.hint.qwen": "Получить ключ на dashscope.console.aliyun.com",
  "settings.hint.openrouter": "Получить ключ на openrouter.ai",
  "settings.hint.localai": "Опционально для локальных моделей",
  "settings.title": "Настройки проводника",
  "settings.useGitignore": "Использовать .gitignore",
  "settings.useGitignoreHint": "Учитывать правила из .gitignore",
  "settings.useCustomIgnore": "Использовать пользовательские правила",
  "settings.useCustomIgnoreHint": "Применять собственные правила фильтрации",
  "settings.autoSaveSelection": "Автосохранение выбора",
  "settings.autoSaveSelectionHint": "Сохранять выбор файлов автоматически",
  "settings.compactFolders": "Компактные вложенные папки",
  "settings.compactFoldersHint": "Объединять папки с одним вложением",
  "settings.foldersFirst": "Папки в начале списка",
  "settings.allowSelectBinary": "Разрешить выбор бинарных файлов",
  "settings.ignoreRules": "Правила игнорирования",
  "settings.manageRules": "Управление правилами",
  "settings.presets": "Пресеты",
  "settings.restoreSelection": "Восстановить предыдущий выбор ({count} файлов)",
  "settings.modal.title": "Настройки",
  "settings.modal.general": "Общие",
  "settings.modal.ai": "AI",
  "settings.modal.export": "Экспорт",
  "settings.modal.fileExplorer": "Проводник",
  "settings.modal.cancel": "Отмена",
  "settings.modal.language": "Язык",
  "settings.modal.theme": "Тема",
  "settings.modal.themeDark": "Тёмная",
  "settings.modal.themeLight": "Светлая",
  "settings.modal.themeAuto": "Авто",
  "settings.modal.filterSettings": "Фильтрация файлов",
  "settings.modal.displaySettings": "Отображение",
  "settings.modal.tourHint": "Пройдите обучение, чтобы узнать все возможности",
  "export.format": "Формат",
  "export.markdown": "Markdown",
  "export.plainText": "Простой текст",
  "export.json": "JSON",
  "export.xml": "XML",
  "export.includeMetadata": "Включить метаданные",
  "export.includeLineNumbers": "Номера строк",
  "export.includeSeparators": "Разделители файлов",
  "export.wrapInCodeBlocks": "Обернуть в code blocks",
  "export.maxLineLength": "Макс. длина строки",
  "export.unlimited": "Без ограничений",
  "export.autoSplit": "Авто-разбиение на чанки",
  "export.tokensPerChunk": "Токенов на чанк",
  "export.strategy": "Стратегия",
  "export.strategySmart": "Smart",
  "export.strategyFile": "По файлам",
  "export.strategyToken": "По токенам",
  "export.tokenLimit": "Лимит токенов",
  "export.custom": "Свой размер",
  "chunks.parts": "частей",
  "chunks.part": "Часть",
  "chunks.of": "из",
  "chunks.lines": "Строки",
  "chunks.copyNext": "Далее",
  "chunks.copyPart": "Часть {n}",
  "chunks.endPart": "КОНЕЦ ЧАСТИ {n}",
  "chunks.startPart": "НАЧАЛО ЧАСТИ {n}",
  "chunks.generate": "Разбить на части",
  "chunks.copied": "скопировано",
  "chunks.allCopied": "Все части скопированы!",
  "chunks.reset": "Сбросить",
  "chunks.clickToCopy": "Клик для копирования",
  "chunks.exitMode": "Выйти из режима чанков",
  "settings.modal.system": "Система",
  "settings.shellIntegration.title": "Интеграция с системой",
  "settings.shellIntegration.description": "Добавить пункт «Открыть в Syntaxia» в контекстное меню папок",
  "settings.shellIntegration.enable": "Включить интеграцию",
  "settings.shellIntegration.disable": "Отключить интеграцию",
  "settings.shellIntegration.enabled": "Интеграция включена",
  "settings.shellIntegration.disabled": "Интеграция отключена",
  "settings.shellIntegration.enableSuccess": "Контекстное меню добавлено",
  "settings.shellIntegration.disableSuccess": "Контекстное меню удалено",
  "settings.shellIntegration.error": "Ошибка при изменении интеграции",
  "settings.shellIntegration.requiresAdmin": "Может потребоваться перезапуск проводника"
};
const task = {
  "task.title": "Описание задачи",
  "task.placeholder": "Опишите что вы хотите сделать...\n\nПримеры:\n• Рефакторинг сервиса авторизации\n• Добавить поддержку тёмной темы\n• Реализовать загрузку файлов",
  "task.quickSuggestions": "Быстрые предложения:",
  "task.analysisResult": "Результат анализа:",
  "task.complexity": "Сложность",
  "task.estimatedFiles": "Примерно файлов",
  "task.suggestedApproach": "Рекомендуемый подход"
};
const templates = {
  "templates.title": "Шаблон промпта",
  "templates.select": "Выбрать шаблон",
  "templates.selectTemplate": "Выбрать шаблон промпта",
  "templates.settings": "Настройки шаблонов",
  "templates.task": "Задача",
  "templates.taskPlaceholder": "Опишите задачу для AI...",
  "templates.sections": "Секции",
  "templates.suggestion": "Рекомендуем:",
  "templates.history": "История задач",
  "templates.recentTasks": "Недавние задачи",
  "templates.userRules": "Ваши правила",
  "templates.userRulesPlaceholder": "Дополнительные правила для AI...",
  "templates.noRole": "Роль не задана",
  "templates.customPrefixSuffix": "Префикс / Суффикс",
  "templates.prefix": "Префикс",
  "templates.suffix": "Суффикс",
  "templates.prefixPlaceholder": "Текст в начале промпта...",
  "templates.suffixPlaceholder": "Текст в конце промпта...",
  "templates.additional": "Дополнительно",
  "templates.manage": "Управление шаблонами",
  "templates.templates": "Шаблоны",
  "templates.builtIn": "Встроенные",
  "templates.custom": "Пользовательские",
  "templates.favorites": "Избранные",
  "templates.create": "Создать шаблон",
  "templates.duplicate": "Дублировать",
  "templates.delete": "Удалить",
  "templates.selectToEdit": "Выберите шаблон для редактирования",
  "templates.preview": "Предпросмотр",
  "templates.import": "Импорт",
  "templates.export": "Экспорт",
  "templates.deleteConfirm": "Удалить шаблон?",
  "templates.deleteConfirmText": "Это действие нельзя отменить.",
  "templates.unsavedChanges": "Несохранённые изменения",
  "templates.unsavedChangesText": "У вас есть несохранённые изменения. Сохранить?",
  "templates.discard": "Отменить",
  "templates.icon": "Иконка",
  "templates.name": "Название",
  "templates.namePlaceholder": "Название шаблона",
  "templates.description": "Описание",
  "templates.descriptionPlaceholder": "Краткое описание шаблона",
  "templates.roleContent": "Роль AI",
  "templates.rolePlaceholder": "Опишите роль AI агента...",
  "templates.rulesContent": "Правила",
  "templates.rulesPlaceholder": "Правила и ограничения для AI...",
  "templates.section.role": "Роль",
  "templates.section.rules": "Правила",
  "templates.section.tree": "Дерево файлов",
  "templates.section.stats": "Статистика",
  "templates.section.task": "Задача",
  "templates.section.files": "Файлы",
  "templates.sectionHint.role": "Описание роли AI агента",
  "templates.sectionHint.rules": "Правила и ограничения",
  "templates.sectionHint.tree": "Структура файлов проекта",
  "templates.sectionHint.stats": "Количество файлов, токенов, языки",
  "templates.sectionHint.task": "Описание вашей задачи",
  "templates.sectionHint.files": "Содержимое выбранных файлов",
  "templates.clearTask": "Очистить задачу",
  "templates.clearHistory": "Очистить историю",
  "templates.currentTemplate": "Активный шаблон",
  "templates.search": "Поиск шаблонов...",
  "templates.noResults": "Ничего не найдено",
  "templates.showHidden": "Показать скрытые",
  "templates.hidden": "Скрытый шаблон",
  "templates.hide": "Скрыть",
  "templates.show": "Показать",
  "templates.toggleFavorite": "В избранное",
  "templates.apply": "Применить",
  "templates.builtInReadonly": "Встроенные шаблоны нельзя редактировать. Создайте копию.",
  "templates.noChanges": "Нет изменений",
  "templates.reset": "Сбросить",
  "templates.resetToDefault": "Сбросить настройки шаблона",
  "templates.resetConfirm": "Сбросить настройки?",
  "templates.resetConfirmText": "Избранное и скрытие будут сброшены.",
  "templates.autosave": "Автосохранение",
  "templates.autosaveHint": "Автоматически сохранять при закрытии",
  "templates.sectionsHint": "Включите нужные секции",
  "templates.enableSections": "Включите секции выше",
  "templates.chars": "симв.",
  "templates.includedSections": "Включить в промпт",
  "templates.contextOptions": "Опции контекста",
  "templates.previewEmpty": "Выберите шаблон для предпросмотра",
  "templates.saved": "Сохранено",
  "templates.deleted": "Шаблон удалён",
  "templates.duplicated": "Шаблон дублирован"
};
const welcome = {
  "welcome.title": "Syntaxia",
  "welcome.subtitle": "Выберите проект для начала работы",
  "welcome.dragDrop": "или перетащите папку сюда",
  "welcome.openProject": "Открыть проект",
  "welcome.recentProjects": "Недавние проекты",
  "welcome.noRecentProjects": "Нет недавних проектов",
  "welcome.settings": "Настройки",
  "welcome.autoOpen": "Автоматически открывать последний проект при запуске",
  "welcome.autoOpenShort": "Автозапуск",
  "welcome.dropHere": "Отпустите папку здесь",
  "welcome.toOpenProject": "Чтобы открыть как проект",
  "welcome.version": "Версия",
  "welcome.removeProject": "Удалить из списка",
  "welcome.clearHistory": "Очистить историю",
  "welcome.confirmClearHistory": "Очистить всю историю проектов?",
  "welcome.projectRemoved": "Проект удалён из списка",
  "welcome.historyCleared": "История очищена",
  "welcome.openInExplorer": "Открыть в проводнике",
  "welcome.copyPath": "Копировать путь",
  "welcome.pathCopied": "Путь скопирован",
  "welcome.pinProject": "Закрепить",
  "welcome.unpinProject": "Открепить"
};
const ru = {
  ...common,
  ...files,
  ...context,
  ...chat,
  ...git,
  ...settings,
  ...errors,
  ...templates,
  ...welcome,
  ...onboarding,
  ...exportModule,
  ...task,
  ...commands
};
const translations = {
  ru,
  en
};
const currentLocale = ref(localStorage.getItem("app-locale") || "ru");
function plural(count, one, few, many) {
  const absCount = Math.abs(count);
  const mod10 = absCount % 10;
  const mod100 = absCount % 100;
  if (mod10 === 1 && mod100 !== 11) {
    return one;
  }
  if (mod10 >= 2 && mod10 <= 4 && (mod100 < 10 || mod100 >= 20)) {
    return few;
  }
  return many;
}
function useI18n() {
  const t = (key, params) => {
    const locale2 = currentLocale.value;
    const dict = translations[locale2];
    let result = dict[key] || key;
    if (params) {
      for (const [paramKey, paramValue] of Object.entries(params)) {
        result = result.replace(new RegExp(`\\{${paramKey}\\}`, "g"), String(paramValue));
      }
    }
    return result;
  };
  const setLocale = (locale2) => {
    currentLocale.value = locale2;
    localStorage.setItem("app-locale", locale2);
  };
  const locale = computed(() => currentLocale.value);
  return {
    t,
    locale,
    setLocale,
    plural
  };
}
function useLogger(context2) {
  const prefix = `[${context2}]`;
  return {
    debug: (...args) => {
    },
    info: (...args) => {
    },
    warn: (...args) => {
      console.warn(prefix, ...args);
    },
    error: (...args) => {
      console.error(prefix, ...args);
    }
  };
}
function AddRecentProject(arg1, arg2) {
  return window["go"]["main"]["App"]["AddRecentProject"](arg1, arg2);
}
function AddToGitignore(arg1, arg2) {
  return window["go"]["main"]["App"]["AddToGitignore"](arg1, arg2);
}
function AgenticChat(arg1) {
  return window["go"]["main"]["App"]["AgenticChat"](arg1);
}
function AnalyzeFile(arg1, arg2) {
  return window["go"]["main"]["App"]["AnalyzeFile"](arg1, arg2);
}
function AnalyzeProject(arg1, arg2) {
  return window["go"]["main"]["App"]["AnalyzeProject"](arg1, arg2);
}
function AnalyzeTaskAndCollectContext(arg1, arg2, arg3) {
  return window["go"]["main"]["App"]["AnalyzeTaskAndCollectContext"](arg1, arg2, arg3);
}
function ApplyEdits(arg1) {
  return window["go"]["main"]["App"]["ApplyEdits"](arg1);
}
function ApplySandboxChanges() {
  return window["go"]["main"]["App"]["ApplySandboxChanges"]();
}
function ApplySingleEdit(arg1) {
  return window["go"]["main"]["App"]["ApplySingleEdit"](arg1);
}
function Build(arg1, arg2) {
  return window["go"]["main"]["App"]["Build"](arg1, arg2);
}
function BuildContext(arg1, arg2, arg3) {
  return window["go"]["main"]["App"]["BuildContext"](arg1, arg2, arg3);
}
function BuildContextAtRef(arg1, arg2, arg3, arg4) {
  return window["go"]["main"]["App"]["BuildContextAtRef"](arg1, arg2, arg3, arg4);
}
function BuildContextFromRequest(arg1, arg2, arg3) {
  return window["go"]["main"]["App"]["BuildContextFromRequest"](arg1, arg2, arg3);
}
function CheckoutBranch(arg1, arg2) {
  return window["go"]["main"]["App"]["CheckoutBranch"](arg1, arg2);
}
function CheckoutCommit(arg1, arg2) {
  return window["go"]["main"]["App"]["CheckoutCommit"](arg1, arg2);
}
function CleanupTempRepository(arg1) {
  return window["go"]["main"]["App"]["CleanupTempRepository"](arg1);
}
function ClearFileTreeCache() {
  return window["go"]["main"]["App"]["ClearFileTreeCache"]();
}
function ClearStartupPath() {
  return window["go"]["main"]["App"]["ClearStartupPath"]();
}
function CloneRepository(arg1) {
  return window["go"]["main"]["App"]["CloneRepository"](arg1);
}
function CollectSmartContext(arg1) {
  return window["go"]["main"]["App"]["CollectSmartContext"](arg1);
}
function DeleteContext(arg1) {
  return window["go"]["main"]["App"]["DeleteContext"](arg1);
}
function DetectLanguages(arg1) {
  return window["go"]["main"]["App"]["DetectLanguages"](arg1);
}
function DiscardSandboxChanges() {
  return window["go"]["main"]["App"]["DiscardSandboxChanges"]();
}
function DiscardSandboxFile(arg1) {
  return window["go"]["main"]["App"]["DiscardSandboxFile"](arg1);
}
function DiscoverTests(arg1, arg2) {
  return window["go"]["main"]["App"]["DiscoverTests"](arg1, arg2);
}
function ExecuteTaskProtocol(arg1) {
  return window["go"]["main"]["App"]["ExecuteTaskProtocol"](arg1);
}
function ExportContext(arg1) {
  return window["go"]["main"]["App"]["ExportContext"](arg1);
}
function ExportProject(arg1, arg2, arg3) {
  return window["go"]["main"]["App"]["ExportProject"](arg1, arg2, arg3);
}
function FindContextByTopic(arg1, arg2) {
  return window["go"]["main"]["App"]["FindContextByTopic"](arg1, arg2);
}
function GenerateCode(arg1, arg2) {
  return window["go"]["main"]["App"]["GenerateCode"](arg1, arg2);
}
function GenerateCodeStream(arg1, arg2) {
  return window["go"]["main"]["App"]["GenerateCodeStream"](arg1, arg2);
}
function GenerateDiff(arg1, arg2, arg3) {
  return window["go"]["main"]["App"]["GenerateDiff"](arg1, arg2, arg3);
}
function GenerateIntelligentCode(arg1, arg2, arg3) {
  return window["go"]["main"]["App"]["GenerateIntelligentCode"](arg1, arg2, arg3);
}
function GenerateReport(arg1, arg2) {
  return window["go"]["main"]["App"]["GenerateReport"](arg1, arg2);
}
function GetBranches(arg1) {
  return window["go"]["main"]["App"]["GetBranches"](arg1);
}
function GetBudgetPolicies() {
  return window["go"]["main"]["App"]["GetBudgetPolicies"]();
}
function GetCommitHistory(arg1, arg2) {
  return window["go"]["main"]["App"]["GetCommitHistory"](arg1, arg2);
}
function GetContext(arg1) {
  return window["go"]["main"]["App"]["GetContext"](arg1);
}
function GetCurrentBranch(arg1) {
  return window["go"]["main"]["App"]["GetCurrentBranch"](arg1);
}
function GetCurrentDirectory() {
  return window["go"]["main"]["App"]["GetCurrentDirectory"]();
}
function GetCustomIgnoreRules() {
  return window["go"]["main"]["App"]["GetCustomIgnoreRules"]();
}
function GetFileAtRef(arg1, arg2, arg3) {
  return window["go"]["main"]["App"]["GetFileAtRef"](arg1, arg2, arg3);
}
function GetFileQuickInfo(arg1, arg2) {
  return window["go"]["main"]["App"]["GetFileQuickInfo"](arg1, arg2);
}
function GetFileStats(arg1) {
  return window["go"]["main"]["App"]["GetFileStats"](arg1);
}
function GetFullContextContent(arg1) {
  return window["go"]["main"]["App"]["GetFullContextContent"](arg1);
}
function GetGitignoreContentForProject(arg1) {
  return window["go"]["main"]["App"]["GetGitignoreContentForProject"](arg1);
}
function GetGuardrailPolicies() {
  return window["go"]["main"]["App"]["GetGuardrailPolicies"]();
}
function GetImpactPreview(arg1, arg2) {
  return window["go"]["main"]["App"]["GetImpactPreview"](arg1, arg2);
}
function GetProjectContexts(arg1) {
  return window["go"]["main"]["App"]["GetProjectContexts"](arg1);
}
function GetProjectStructure(arg1) {
  return window["go"]["main"]["App"]["GetProjectStructure"](arg1);
}
function GetProviderInfo() {
  return window["go"]["main"]["App"]["GetProviderInfo"]();
}
function GetRecentContexts(arg1, arg2) {
  return window["go"]["main"]["App"]["GetRecentContexts"](arg1, arg2);
}
function GetRecentProjects() {
  return window["go"]["main"]["App"]["GetRecentProjects"]();
}
function GetReleases() {
  return window["go"]["main"]["App"]["GetReleases"]();
}
function GetRemoteBranches(arg1) {
  return window["go"]["main"]["App"]["GetRemoteBranches"](arg1);
}
function GetReport(arg1) {
  return window["go"]["main"]["App"]["GetReport"](arg1);
}
function GetRichCommitHistory(arg1, arg2, arg3) {
  return window["go"]["main"]["App"]["GetRichCommitHistory"](arg1, arg2, arg3);
}
function GetSandboxAllDiffs() {
  return window["go"]["main"]["App"]["GetSandboxAllDiffs"]();
}
function GetSandboxChangeCount() {
  return window["go"]["main"]["App"]["GetSandboxChangeCount"]();
}
function GetSandboxChanges() {
  return window["go"]["main"]["App"]["GetSandboxChanges"]();
}
function GetSandboxDiff(arg1) {
  return window["go"]["main"]["App"]["GetSandboxDiff"](arg1);
}
function GetSettings() {
  return window["go"]["main"]["App"]["GetSettings"]();
}
function GetShellIntegrationStatus() {
  return window["go"]["main"]["App"]["GetShellIntegrationStatus"]();
}
function GetSmartSuggestions(arg1, arg2, arg3) {
  return window["go"]["main"]["App"]["GetSmartSuggestions"](arg1, arg2, arg3);
}
function GetStartupPath() {
  return window["go"]["main"]["App"]["GetStartupPath"]();
}
function GetSupportedAnalyzers() {
  return window["go"]["main"]["App"]["GetSupportedAnalyzers"]();
}
function GetTaskProtocolConfiguration(arg1, arg2) {
  return window["go"]["main"]["App"]["GetTaskProtocolConfiguration"](arg1, arg2);
}
function GetUncommittedFiles(arg1) {
  return window["go"]["main"]["App"]["GetUncommittedFiles"](arg1);
}
function GetVersionInfo() {
  return window["go"]["main"]["App"]["GetVersionInfo"]();
}
function GitHubBuildContext(arg1, arg2, arg3) {
  return window["go"]["main"]["App"]["GitHubBuildContext"](arg1, arg2, arg3);
}
function GitHubGetBranches(arg1) {
  return window["go"]["main"]["App"]["GitHubGetBranches"](arg1);
}
function GitHubGetCommits(arg1, arg2, arg3) {
  return window["go"]["main"]["App"]["GitHubGetCommits"](arg1, arg2, arg3);
}
function GitHubGetDefaultBranch(arg1) {
  return window["go"]["main"]["App"]["GitHubGetDefaultBranch"](arg1);
}
function GitHubGetFileContent(arg1, arg2, arg3) {
  return window["go"]["main"]["App"]["GitHubGetFileContent"](arg1, arg2, arg3);
}
function GitHubListFiles(arg1, arg2) {
  return window["go"]["main"]["App"]["GitHubListFiles"](arg1, arg2);
}
function GitLabBuildContext(arg1, arg2, arg3) {
  return window["go"]["main"]["App"]["GitLabBuildContext"](arg1, arg2, arg3);
}
function GitLabGetBranches(arg1) {
  return window["go"]["main"]["App"]["GitLabGetBranches"](arg1);
}
function GitLabGetCommits(arg1, arg2, arg3) {
  return window["go"]["main"]["App"]["GitLabGetCommits"](arg1, arg2, arg3);
}
function GitLabGetDefaultBranch(arg1) {
  return window["go"]["main"]["App"]["GitLabGetDefaultBranch"](arg1);
}
function GitLabGetFileContent(arg1, arg2, arg3) {
  return window["go"]["main"]["App"]["GitLabGetFileContent"](arg1, arg2, arg3);
}
function GitLabListFiles(arg1, arg2) {
  return window["go"]["main"]["App"]["GitLabListFiles"](arg1, arg2);
}
function HasSandboxChanges() {
  return window["go"]["main"]["App"]["HasSandboxChanges"]();
}
function IsGitAvailable() {
  return window["go"]["main"]["App"]["IsGitAvailable"]();
}
function IsGitHubURL(arg1) {
  return window["go"]["main"]["App"]["IsGitHubURL"](arg1);
}
function IsGitLabURL(arg1) {
  return window["go"]["main"]["App"]["IsGitLabURL"](arg1);
}
function IsGitRepository(arg1) {
  return window["go"]["main"]["App"]["IsGitRepository"](arg1);
}
function IsSemanticSearchAvailable() {
  return window["go"]["main"]["App"]["IsSemanticSearchAvailable"]();
}
function ListAvailableModels() {
  return window["go"]["main"]["App"]["ListAvailableModels"]();
}
function ListFiles(arg1, arg2, arg3) {
  return window["go"]["main"]["App"]["ListFiles"](arg1, arg2, arg3);
}
function ListFilesAtRef(arg1, arg2) {
  return window["go"]["main"]["App"]["ListFilesAtRef"](arg1, arg2);
}
function ListReports(arg1) {
  return window["go"]["main"]["App"]["ListReports"](arg1);
}
function PathExists(arg1) {
  return window["go"]["main"]["App"]["PathExists"](arg1);
}
function QwenExecuteTask(arg1) {
  return window["go"]["main"]["App"]["QwenExecuteTask"](arg1);
}
function QwenGetAvailableModels() {
  return window["go"]["main"]["App"]["QwenGetAvailableModels"]();
}
function QwenPreviewContext(arg1) {
  return window["go"]["main"]["App"]["QwenPreviewContext"](arg1);
}
function ReadFileContent(arg1, arg2) {
  return window["go"]["main"]["App"]["ReadFileContent"](arg1, arg2);
}
function RegisterShellIntegration() {
  return window["go"]["main"]["App"]["RegisterShellIntegration"]();
}
function RemoveRecentProject(arg1) {
  return window["go"]["main"]["App"]["RemoveRecentProject"](arg1);
}
function RunTests(arg1) {
  return window["go"]["main"]["App"]["RunTests"](arg1);
}
function SaveContextMemory(arg1, arg2, arg3, arg4) {
  return window["go"]["main"]["App"]["SaveContextMemory"](arg1, arg2, arg3, arg4);
}
function SaveSettings(arg1) {
  return window["go"]["main"]["App"]["SaveSettings"](arg1);
}
function SelectDirectory() {
  return window["go"]["main"]["App"]["SelectDirectory"]();
}
function SemanticFindSimilar(arg1) {
  return window["go"]["main"]["App"]["SemanticFindSimilar"](arg1);
}
function SemanticGetStats(arg1) {
  return window["go"]["main"]["App"]["SemanticGetStats"](arg1);
}
function SemanticHybridSearch(arg1) {
  return window["go"]["main"]["App"]["SemanticHybridSearch"](arg1);
}
function SemanticIndexFile(arg1, arg2) {
  return window["go"]["main"]["App"]["SemanticIndexFile"](arg1, arg2);
}
function SemanticIndexProject(arg1) {
  return window["go"]["main"]["App"]["SemanticIndexProject"](arg1);
}
function SemanticIsIndexed(arg1) {
  return window["go"]["main"]["App"]["SemanticIsIndexed"](arg1);
}
function SemanticRetrieveContext(arg1) {
  return window["go"]["main"]["App"]["SemanticRetrieveContext"](arg1);
}
function SemanticSearch(arg1) {
  return window["go"]["main"]["App"]["SemanticSearch"](arg1);
}
function SuggestContextFiles(arg1, arg2) {
  return window["go"]["main"]["App"]["SuggestContextFiles"](arg1, arg2);
}
function TestIgnoreRules(arg1, arg2) {
  return window["go"]["main"]["App"]["TestIgnoreRules"](arg1, arg2);
}
function TestIgnoreRulesDetailed(arg1, arg2) {
  return window["go"]["main"]["App"]["TestIgnoreRulesDetailed"](arg1, arg2);
}
function TypeCheck(arg1, arg2) {
  return window["go"]["main"]["App"]["TypeCheck"](arg1, arg2);
}
function UnregisterShellIntegration() {
  return window["go"]["main"]["App"]["UnregisterShellIntegration"]();
}
function UpdateCustomIgnoreRules(arg1) {
  return window["go"]["main"]["App"]["UpdateCustomIgnoreRules"](arg1);
}
function ValidatePath(arg1) {
  return window["go"]["main"]["App"]["ValidatePath"](arg1);
}
const _hoisted_1$1a = ["title"];
const _hoisted_2$16 = { class: "trigger-content" };
const _hoisted_3$11 = { class: "trigger-version" };
const _hoisted_4$X = {
  key: 0,
  class: "update-dot"
};
const _hoisted_5$Q = {
  key: 0,
  class: "version-panel"
};
const _hoisted_6$M = { class: "panel-header" };
const _hoisted_7$I = { class: "panel-title-row" };
const _hoisted_8$F = { class: "panel-title-text" };
const _hoisted_9$C = { class: "panel-title" };
const _hoisted_10$z = { class: "panel-version" };
const _hoisted_11$x = { class: "panel-content" };
const _hoisted_12$s = {
  key: 0,
  class: "panel-state"
};
const _hoisted_13$s = {
  key: 1,
  class: "panel-state panel-state-error"
};
const _hoisted_14$p = {
  key: 2,
  class: "panel-state"
};
const _hoisted_15$n = { class: "no-releases-text" };
const _hoisted_16$k = {
  key: 3,
  class: "releases-list"
};
const _hoisted_17$k = { class: "release-row" };
const _hoisted_18$j = { class: "release-version" };
const _hoisted_19$h = {
  key: 0,
  class: "release-tag"
};
const _hoisted_20$g = {
  key: 1,
  class: "release-tag release-tag-pre"
};
const _hoisted_21$e = { class: "release-date" };
const _hoisted_22$c = {
  href: "https://github.com/WhiteBite/syntaxia/releases",
  target: "_blank",
  rel: "noopener noreferrer",
  class: "panel-footer"
};
const _sfc_main$1c = /* @__PURE__ */ defineComponent({
  __name: "VersionBadge",
  setup(__props) {
    const logger2 = useLogger("VersionBadge");
    const { t, locale } = useI18n();
    const isExpanded = ref(false);
    const isLoading = ref(false);
    const error = ref(null);
    const currentVersion = ref("dev");
    const latestVersion = ref("");
    const hasUpdate = ref(false);
    const releases = ref([]);
    const isDev = computed(() => {
      const v = currentVersion.value.toLowerCase();
      return v === "dev" || v === "vdev" || v.includes("dev");
    });
    const displayVersion = computed(() => {
      const v = currentVersion.value;
      if (isDev.value) return "vdev";
      return v.startsWith("v") ? v : `v${v}`;
    });
    onMounted(async () => {
      try {
        const info = await GetVersionInfo();
        currentVersion.value = info.version;
      } catch (e) {
        console.warn("Failed to get version info:", e);
      }
    });
    async function loadReleases() {
      if (isLoading.value) return;
      isLoading.value = true;
      error.value = null;
      try {
        const response = await GetReleases();
        currentVersion.value = response.currentVersion;
        latestVersion.value = response.latestVersion;
        hasUpdate.value = response.hasUpdate;
        releases.value = response.releases || [];
        if (response.error) {
          error.value = response.error;
        }
      } catch (e) {
        logger2.error("Failed to load releases:", e);
        error.value = e instanceof Error ? e.message : "Unknown error";
      } finally {
        isLoading.value = false;
      }
    }
    function expand() {
      isExpanded.value = true;
      if (releases.value.length === 0 && !error.value) {
        loadReleases();
      }
    }
    function collapse() {
      isExpanded.value = false;
    }
    function formatDate(dateStr) {
      if (!dateStr) return "";
      const date = new Date(dateStr);
      return date.toLocaleDateString(locale.value === "ru" ? "ru-RU" : "en-US", {
        month: "short",
        day: "numeric"
      });
    }
    return (_ctx, _cache) => {
      return openBlock(), createElementBlock("div", {
        class: normalizeClass(["version-widget", { "is-expanded": isExpanded.value }])
      }, [
        !isExpanded.value ? (openBlock(), createElementBlock("button", {
          key: 0,
          onClick: expand,
          class: normalizeClass(["version-trigger", { "has-update": hasUpdate.value }]),
          title: isDev.value ? unref(t)("version.devModeTooltip") : unref(t)("version.clickToViewChangelog")
        }, [
          createBaseVNode("div", _hoisted_2$16, [
            _cache[0] || (_cache[0] = createBaseVNode("svg", {
              class: "trigger-icon",
              viewBox: "0 0 24 24",
              fill: "none",
              stroke: "currentColor"
            }, [
              createBaseVNode("path", {
                "stroke-linecap": "round",
                "stroke-linejoin": "round",
                "stroke-width": "2",
                d: "M7 7h.01M7 3h5c.512 0 1.024.195 1.414.586l7 7a2 2 0 010 2.828l-7 7a2 2 0 01-2.828 0l-7-7A1.994 1.994 0 013 12V7a4 4 0 014-4z"
              })
            ], -1)),
            createBaseVNode("span", _hoisted_3$11, toDisplayString(displayVersion.value), 1),
            hasUpdate.value ? (openBlock(), createElementBlock("span", _hoisted_4$X)) : createCommentVNode("", true)
          ])
        ], 10, _hoisted_1$1a)) : createCommentVNode("", true),
        createVNode(Transition, { name: "panel" }, {
          default: withCtx(() => [
            isExpanded.value ? (openBlock(), createElementBlock("div", _hoisted_5$Q, [
              createBaseVNode("div", _hoisted_6$M, [
                createBaseVNode("div", _hoisted_7$I, [
                  _cache[1] || (_cache[1] = createBaseVNode("div", { class: "panel-icon" }, [
                    createBaseVNode("svg", {
                      viewBox: "0 0 24 24",
                      fill: "none",
                      stroke: "currentColor"
                    }, [
                      createBaseVNode("path", {
                        "stroke-linecap": "round",
                        "stroke-linejoin": "round",
                        "stroke-width": "2",
                        d: "M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"
                      })
                    ])
                  ], -1)),
                  createBaseVNode("div", _hoisted_8$F, [
                    createBaseVNode("h3", _hoisted_9$C, toDisplayString(unref(t)("version.changelog")), 1),
                    createBaseVNode("span", _hoisted_10$z, toDisplayString(displayVersion.value), 1)
                  ])
                ]),
                createBaseVNode("button", {
                  onClick: collapse,
                  class: "panel-close"
                }, [..._cache[2] || (_cache[2] = [
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
                      d: "M6 18L18 6M6 6l12 12"
                    })
                  ], -1)
                ])])
              ]),
              createBaseVNode("div", _hoisted_11$x, [
                isLoading.value ? (openBlock(), createElementBlock("div", _hoisted_12$s, [
                  _cache[3] || (_cache[3] = createBaseVNode("div", { class: "spinner" }, null, -1)),
                  createBaseVNode("span", null, toDisplayString(unref(t)("version.checkingUpdates")), 1)
                ])) : error.value ? (openBlock(), createElementBlock("div", _hoisted_13$s, [
                  createBaseVNode("span", null, toDisplayString(unref(t)("version.errorLoading")), 1),
                  createBaseVNode("button", {
                    onClick: loadReleases,
                    class: "retry-link"
                  }, toDisplayString(unref(t)("action.retry") || "Повторить"), 1)
                ])) : releases.value.length === 0 ? (openBlock(), createElementBlock("div", _hoisted_14$p, [
                  createBaseVNode("span", _hoisted_15$n, toDisplayString(unref(t)("version.noReleases")), 1),
                  _cache[4] || (_cache[4] = createBaseVNode("span", { class: "no-releases-hint" }, "Релизы появятся после публикации", -1))
                ])) : (openBlock(), createElementBlock("div", _hoisted_16$k, [
                  (openBlock(true), createElementBlock(Fragment, null, renderList(releases.value.slice(0, 5), (release, index) => {
                    return openBlock(), createElementBlock("div", {
                      key: release.tag_name,
                      class: normalizeClass(["release-item", { "is-latest": index === 0 }])
                    }, [
                      createBaseVNode("div", _hoisted_17$k, [
                        createBaseVNode("span", _hoisted_18$j, toDisplayString(release.tag_name), 1),
                        index === 0 ? (openBlock(), createElementBlock("span", _hoisted_19$h, "Latest")) : createCommentVNode("", true),
                        release.prerelease ? (openBlock(), createElementBlock("span", _hoisted_20$g, "Pre")) : createCommentVNode("", true)
                      ]),
                      createBaseVNode("span", _hoisted_21$e, toDisplayString(formatDate(release.published_at)), 1)
                    ], 2);
                  }), 128))
                ]))
              ]),
              createBaseVNode("a", _hoisted_22$c, [
                _cache[5] || (_cache[5] = createBaseVNode("svg", {
                  class: "w-4 h-4",
                  fill: "currentColor",
                  viewBox: "0 0 24 24"
                }, [
                  createBaseVNode("path", { d: "M12 0c-6.626 0-12 5.373-12 12 0 5.302 3.438 9.8 8.207 11.387.599.111.793-.261.793-.577v-2.234c-3.338.726-4.033-1.416-4.033-1.416-.546-1.387-1.333-1.756-1.333-1.756-1.089-.745.083-.729.083-.729 1.205.084 1.839 1.237 1.839 1.237 1.07 1.834 2.807 1.304 3.492.997.107-.775.418-1.305.762-1.604-2.665-.305-5.467-1.334-5.467-5.931 0-1.311.469-2.381 1.236-3.221-.124-.303-.535-1.524.117-3.176 0 0 1.008-.322 3.301 1.23.957-.266 1.983-.399 3.003-.404 1.02.005 2.047.138 3.006.404 2.291-1.552 3.297-1.23 3.297-1.23.653 1.653.242 2.874.118 3.176.77.84 1.235 1.911 1.235 3.221 0 4.609-2.807 5.624-5.479 5.921.43.372.823 1.102.823 2.222v3.293c0 .319.192.694.801.576 4.765-1.589 8.199-6.086 8.199-11.386 0-6.627-5.373-12-12-12z" })
                ], -1)),
                createBaseVNode("span", null, toDisplayString(unref(t)("version.viewOnGithub")), 1),
                _cache[6] || (_cache[6] = createBaseVNode("svg", {
                  class: "w-3 h-3 ml-auto",
                  fill: "none",
                  stroke: "currentColor",
                  viewBox: "0 0 24 24"
                }, [
                  createBaseVNode("path", {
                    "stroke-linecap": "round",
                    "stroke-linejoin": "round",
                    "stroke-width": "2",
                    d: "M10 6H6a2 2 0 00-2 2v10a2 2 0 002 2h10a2 2 0 002-2v-4M14 4h6m0 0v6m0-6L10 14"
                  })
                ], -1))
              ])
            ])) : createCommentVNode("", true)
          ]),
          _: 1
        })
      ], 2);
    };
  }
});
const _export_sfc = (sfc, props) => {
  const target2 = sfc.__vccOpts || sfc;
  for (const [key, val] of props) {
    target2[key] = val;
  }
  return target2;
};
const VersionBadge = /* @__PURE__ */ _export_sfc(_sfc_main$1c, [["__scopeId", "data-v-fda4fd05"]]);
const logger$p = useLogger("API");
class ApiError extends Error {
  constructor(context2, originalError) {
    const message = originalError instanceof Error ? `${context2}: ${originalError.message}` : context2;
    super(message);
    this.context = context2;
    this.originalError = originalError;
    this.name = "ApiError";
  }
}
async function apiCall(fn, context2, options2) {
  const logPrefix = options2?.logContext || context2;
  try {
    return await fn();
  } catch (error) {
    logger$p.error(`${logPrefix}:`, error);
    if (options2?.rethrow && error instanceof Error) {
      throw error;
    }
    throw new ApiError(context2, error);
  }
}
async function apiCallSafe(fn, defaultValue, context2) {
  try {
    return await fn();
  } catch (error) {
    if (context2) {
      logger$p.warn(`${context2} (using default):`, error);
    }
    return defaultValue;
  }
}
function parseJsonResponse(json, context2) {
  try {
    return JSON.parse(json);
  } catch (error) {
    logger$p.error(`Failed to parse JSON (${context2}):`, error);
    throw new ApiError(`Invalid JSON response: ${context2}`, error);
  }
}
const apiCallWithDefault = apiCallSafe;
const logger$o = useLogger("API:ai");
const aiApi = {
  generateCode: (context2, task2) => apiCall(() => GenerateCode(context2, task2), "Failed to generate code.", { logContext: "ai" }),
  generateCodeStream: (context2, task2) => {
    try {
      GenerateCodeStream(context2, task2);
    } catch (error) {
      logger$o.error("Error starting code stream:", error);
    }
  },
  generateIntelligentCode: (context2, task2, options2) => apiCall(
    () => GenerateIntelligentCode(context2, task2, options2),
    "Failed to generate code.",
    { logContext: "ai" }
  ),
  listAvailableModels: () => apiCall(() => ListAvailableModels(), "Failed to list AI models.", { logContext: "ai" }),
  getProviderInfo: () => apiCall(() => GetProviderInfo(), "Failed to get provider information.", { logContext: "ai" }),
  // Qwen Task Execution
  qwenExecuteTask: async (request) => {
    const result = await apiCall(
      () => QwenExecuteTask(JSON.stringify(request)),
      "Failed to execute task with Qwen.",
      { logContext: "ai" }
    );
    return parseJsonResponse(result, "Failed to parse Qwen response.");
  },
  qwenPreviewContext: async (request) => {
    const result = await apiCall(
      () => QwenPreviewContext(JSON.stringify(request)),
      "Failed to preview context.",
      { logContext: "ai" }
    );
    return parseJsonResponse(result, "Failed to parse context preview.");
  },
  qwenGetAvailableModels: async () => {
    const result = await apiCall(
      () => QwenGetAvailableModels(),
      "Failed to get Qwen models.",
      { logContext: "ai" }
    );
    return parseJsonResponse(result, "Failed to parse Qwen models.");
  }
};
const analysisApi = {
  analyzeProject: (path, analyzers) => apiCall(() => AnalyzeProject(path, analyzers), "Failed to analyze project.", { logContext: "analysis" }),
  analyzeFile: (projectPath, filePath) => apiCall(
    () => AnalyzeFile(projectPath, filePath),
    "Failed to analyze file.",
    { logContext: "analysis" }
  ),
  detectLanguages: (projectPath) => apiCall(
    () => DetectLanguages(projectPath),
    "Failed to detect project languages.",
    { logContext: "analysis" }
  ),
  getSupportedAnalyzers: () => apiCall(
    () => GetSupportedAnalyzers(),
    "Failed to get supported analyzers.",
    { logContext: "analysis" }
  )
};
const buildApi = {
  // Testing
  runTests: (config) => apiCall(() => RunTests(config), "Failed to run tests.", { logContext: "build" }),
  discoverTests: (projectPath, language) => apiCall(
    () => DiscoverTests(projectPath, language),
    "Failed to discover tests.",
    { logContext: "build" }
  ),
  // Build
  build: (projectPath, language) => apiCall(() => Build(projectPath, language), "Failed to build project.", { logContext: "build" }),
  typeCheck: (projectPath, language) => apiCall(() => TypeCheck(projectPath, language), "Failed to type check.", { logContext: "build" }),
  // Diff and Apply
  generateDiff: (original, modified, format) => apiCall(
    () => GenerateDiff(original, modified, format),
    "Failed to generate diff.",
    { logContext: "build" }
  ),
  applyEdits: (edits) => apiCall(() => ApplyEdits(edits), "Failed to apply edits.", { logContext: "build" }),
  applySingleEdit: (edit) => apiCall(() => ApplySingleEdit(edit), "Failed to apply edit.", { logContext: "build" })
};
const logger$n = useLogger("API:context");
const contextApi$1 = {
  buildContext: (projectPath, files2, task2) => apiCall(() => BuildContext(projectPath, files2, task2), "Failed to build context.", { logContext: "context" }),
  buildContextFromRequest: async (projectPath, files2, options2) => {
    try {
      return await BuildContextFromRequest(projectPath, files2, options2);
    } catch (error) {
      logger$n.error("Error building context from request:", error);
      let errorMsg = "";
      if (error && typeof error === "object") {
        const err = error;
        errorMsg = err.cause || err.message || String(error);
      } else if (error instanceof Error) {
        errorMsg = error.message;
      } else {
        errorMsg = String(error);
      }
      if (errorMsg.includes("token limit")) {
        const match = errorMsg.match(/(\d+)\s*[>\u003e]\s*(\d+)/);
        if (match) {
          const actual = Number(match[1]);
          const limit = Number(match[2]);
          const tokenError = new Error("TOKEN_LIMIT_EXCEEDED");
          tokenError.tokenInfo = { actual, limit };
          throw tokenError;
        }
        throw new Error("TOKEN_LIMIT_EXCEEDED");
      }
      throw new Error("Failed to build context.");
    }
  },
  getContext: (contextId) => apiCall(() => GetContext(contextId), "Failed to get context.", { logContext: "context" }),
  deleteContext: (contextId) => apiCall(() => DeleteContext(contextId), "Failed to delete context.", { logContext: "context" }),
  getProjectContexts: (projectPath) => apiCall(() => GetProjectContexts(projectPath), "Failed to get project contexts.", { logContext: "context" }),
  exportContext: (exportSettings) => apiCall(
    () => ExportContext(JSON.stringify(exportSettings)),
    "Failed to export context.",
    { logContext: "context" }
  ),
  getFullContextContent: (contextId) => apiCall(
    () => GetFullContextContent(contextId),
    "Failed to get full context content.",
    { logContext: "context" }
  ),
  suggestContextFiles: (taskDescription, files2) => apiCall(
    () => SuggestContextFiles(taskDescription, files2),
    "Failed to suggest context files.",
    { logContext: "context" }
  ),
  getSmartSuggestions: async (projectPath, currentFiles, task2 = "") => {
    try {
      const result = await GetSmartSuggestions(projectPath, currentFiles, task2);
      return {
        suggestions: result.suggestions?.map((s) => ({
          path: s.path,
          source: s.source,
          reason: s.reason,
          confidence: s.confidence
        })) || [],
        total: result.total || 0
      };
    } catch (error) {
      logger$n.error("Error getting smart suggestions:", error);
      return { suggestions: [], total: 0 };
    }
  },
  getFileQuickInfo: async (projectPath, filePath) => {
    try {
      const result = await GetFileQuickInfo(projectPath, filePath);
      return {
        symbolCount: result.symbolCount,
        importCount: result.importCount,
        dependentCount: result.dependentCount,
        changeRisk: result.changeRisk,
        riskLevel: result.riskLevel
      };
    } catch {
      return { symbolCount: 0, importCount: 0, dependentCount: 0, changeRisk: 0, riskLevel: "low" };
    }
  },
  getImpactPreview: async (projectPath, filePaths) => {
    try {
      const result = await GetImpactPreview(projectPath, filePaths);
      return {
        totalDependents: result.totalDependents,
        aggregateRisk: result.aggregateRisk,
        riskLevel: result.riskLevel,
        affectedFiles: (result.affectedFiles || []).map((f) => ({
          path: f.path,
          type: f.type,
          dependents: f.dependents
        })),
        relatedTests: result.relatedTests || []
      };
    } catch {
      return { totalDependents: 0, aggregateRisk: 0, riskLevel: "low", affectedFiles: [], relatedTests: [] };
    }
  },
  analyzeTaskAndCollectContext: (task2, allFilesJson, rootDir) => apiCall(
    () => AnalyzeTaskAndCollectContext(task2, allFilesJson, rootDir),
    "Failed to analyze task and collect context.",
    { logContext: "context" }
  ),
  agenticChat: async (task2, projectRoot, smartContext) => {
    const request = { task: task2, projectRoot, smartContext };
    const result = await apiCall(
      () => AgenticChat(JSON.stringify(request)),
      "Failed to execute agentic chat.",
      { logContext: "context" }
    );
    return parseJsonResponse(result, "Failed to parse agentic chat response.");
  },
  collectSmartContext: async (request) => {
    const result = await apiCall(
      () => CollectSmartContext(JSON.stringify(request)),
      "Failed to collect smart context.",
      { logContext: "context" }
    );
    return parseJsonResponse(result, "Failed to parse smart context.");
  }
};
const logger$m = useLogger("API:files");
const filesApi$1 = {
  listFiles: (path, useGitignore = true, useCustomIgnore = true) => apiCall(() => ListFiles(path, useGitignore, useCustomIgnore), "Failed to load file tree.", { logContext: "files" }),
  clearFileTreeCache: async () => {
    try {
      await ClearFileTreeCache();
    } catch (error) {
      logger$m.error("Error clearing file tree cache:", error);
    }
  },
  readFileContent: (projectPath, filePath) => apiCall(() => ReadFileContent(projectPath, filePath), "Failed to read file content.", { logContext: "files" }),
  getFileStats: (path) => apiCall(() => GetFileStats(path), "Failed to get file statistics.", { logContext: "files" })
};
const logger$l = useLogger("API:git");
const gitApi$1 = {
  getUncommittedFiles: (repoPath) => apiCall(() => GetUncommittedFiles(repoPath), "Failed to get git status.", { logContext: "git" }),
  getBranches: (repoPath) => apiCall(() => GetBranches(repoPath), "Failed to get git branches.", { logContext: "git" }),
  getCurrentBranch: (repoPath) => apiCall(() => GetCurrentBranch(repoPath), "Failed to get current branch.", { logContext: "git" }),
  getRichCommitHistory: (repoPath, filePath, limit) => apiCall(
    () => GetRichCommitHistory(repoPath, filePath, limit),
    "Failed to get commit history.",
    { logContext: "git" }
  ),
  isGitAvailable: () => apiCallWithDefault(() => IsGitAvailable(), false, "git.isGitAvailable"),
  isGitRepository: (projectPath) => apiCallWithDefault(() => IsGitRepository(projectPath), false, "git.isGitRepository"),
  cloneRepository: (url) => apiCall(() => CloneRepository(url), "Failed to clone repository.", { logContext: "git" }),
  checkoutBranch: (projectPath, branch) => apiCall(() => CheckoutBranch(projectPath, branch), "Failed to checkout branch.", { logContext: "git" }),
  checkoutCommit: (projectPath, commitHash) => apiCall(
    () => CheckoutCommit(projectPath, commitHash),
    "Failed to checkout commit.",
    { logContext: "git" }
  ),
  getCommitHistory: async (projectPath, limit = 50) => {
    const result = await apiCall(
      () => GetCommitHistory(projectPath, limit),
      "Failed to get commit history.",
      { logContext: "git" }
    );
    return parseJsonResponse(result, "Failed to parse commit history.");
  },
  getRemoteBranches: async (projectPath) => {
    const result = await apiCall(
      () => GetRemoteBranches(projectPath),
      "Failed to get remote branches.",
      { logContext: "git" }
    );
    return parseJsonResponse(result, "Failed to parse remote branches.");
  },
  cleanupTempRepository: async (path) => {
    try {
      await CleanupTempRepository(path);
    } catch (error) {
      logger$l.error("Error cleaning up temp repository:", error);
    }
  },
  listFilesAtRef: async (projectPath, ref2) => {
    const result = await apiCall(
      () => ListFilesAtRef(projectPath, ref2),
      "Failed to list files at ref.",
      { logContext: "git" }
    );
    const parsed = parseJsonResponse(result, "Failed to parse files at ref.");
    return Array.isArray(parsed) ? parsed : [];
  },
  getFileAtRef: (projectPath, filePath, ref2) => apiCall(
    () => GetFileAtRef(projectPath, filePath, ref2),
    "Failed to get file at ref.",
    { logContext: "git" }
  ),
  buildContextAtRef: (projectPath, files2, ref2) => apiCall(
    () => BuildContextAtRef(projectPath, files2, ref2, "{}"),
    "Failed to build context at ref.",
    { logContext: "git" }
  )
};
const githubApi = {
  isGitHubURL: (url) => apiCallWithDefault(() => IsGitHubURL(url), false, "github.isGitHubURL"),
  getDefaultBranch: (repoURL) => apiCall(() => GitHubGetDefaultBranch(repoURL), "Failed to get default branch.", { logContext: "github" }),
  getBranches: async (repoURL) => {
    const result = await apiCall(
      () => GitHubGetBranches(repoURL),
      "Failed to get GitHub branches.",
      { logContext: "github" }
    );
    return parseJsonResponse(result, "Failed to parse GitHub branches.");
  },
  getCommits: async (repoURL, branch, limit = 50) => {
    const result = await apiCall(
      () => GitHubGetCommits(repoURL, branch, limit),
      "Failed to get GitHub commits.",
      { logContext: "github" }
    );
    return parseJsonResponse(result, "Failed to parse GitHub commits.");
  },
  listFiles: async (repoURL, ref2) => {
    const result = await apiCall(
      () => GitHubListFiles(repoURL, ref2),
      "Failed to list GitHub files.",
      { logContext: "github" }
    );
    const parsed = parseJsonResponse(result, "Failed to parse GitHub files.");
    return Array.isArray(parsed) ? parsed : [];
  },
  getFileContent: (repoURL, filePath, ref2) => apiCall(
    () => GitHubGetFileContent(repoURL, filePath, ref2),
    "Failed to get GitHub file content.",
    { logContext: "github" }
  ),
  buildContext: (repoURL, files2, ref2) => apiCall(
    () => GitHubBuildContext(repoURL, files2, ref2),
    "Failed to build GitHub context.",
    { logContext: "github" }
  )
};
const gitlabApi = {
  isGitLabURL: (url) => apiCallWithDefault(() => IsGitLabURL(url), false, "gitlab.isGitLabURL"),
  getDefaultBranch: (repoURL) => apiCall(() => GitLabGetDefaultBranch(repoURL), "Failed to get default branch.", { logContext: "gitlab" }),
  getBranches: async (repoURL) => {
    const result = await apiCall(
      () => GitLabGetBranches(repoURL),
      "Failed to get GitLab branches.",
      { logContext: "gitlab" }
    );
    return parseJsonResponse(result, "Failed to parse GitLab branches.");
  },
  getCommits: async (repoURL, branch, limit = 50) => {
    const result = await apiCall(
      () => GitLabGetCommits(repoURL, branch, limit),
      "Failed to get GitLab commits.",
      { logContext: "gitlab" }
    );
    return parseJsonResponse(result, "Failed to parse GitLab commits.");
  },
  listFiles: async (repoURL, ref2) => {
    const result = await apiCall(
      () => GitLabListFiles(repoURL, ref2),
      "Failed to list GitLab files.",
      { logContext: "gitlab" }
    );
    const parsed = parseJsonResponse(result, "Failed to parse GitLab files.");
    return Array.isArray(parsed) ? parsed : [];
  },
  getFileContent: (repoURL, filePath, ref2) => apiCall(
    () => GitLabGetFileContent(repoURL, filePath, ref2),
    "Failed to get GitLab file content.",
    { logContext: "gitlab" }
  ),
  buildContext: (repoURL, files2, ref2) => apiCall(
    () => GitLabBuildContext(repoURL, files2, ref2),
    "Failed to build GitLab context.",
    { logContext: "gitlab" }
  )
};
const memoryApi = {
  getRecentContexts: (projectPath, limit = 10) => apiCallWithDefault(
    // @ts-ignore
    () => GetRecentContexts(projectPath, limit),
    [],
    "memory.getRecentContexts"
  ),
  findContextByTopic: (projectPath, topic) => apiCallWithDefault(
    // @ts-ignore
    () => FindContextByTopic(projectPath, topic),
    [],
    "memory.findContextByTopic"
  ),
  saveContextMemory: (projectPath, topic, summary, files2) => apiCall(
    // @ts-ignore
    () => SaveContextMemory(projectPath, topic, summary, files2),
    "Failed to save context.",
    { logContext: "memory" }
  )
};
const projectApi = {
  getRecentProjects: () => apiCall(() => GetRecentProjects(), "Failed to load recent projects.", { logContext: "project" }),
  addRecentProject: (path, name) => apiCall(() => AddRecentProject(path, name), "Failed to add recent project.", { logContext: "project" }),
  removeRecentProject: (path) => apiCall(() => RemoveRecentProject(path), "Failed to remove recent project.", { logContext: "project" }),
  selectDirectory: () => apiCall(() => SelectDirectory(), "Failed to select directory.", { logContext: "project" }),
  getCurrentDirectory: () => apiCall(() => GetCurrentDirectory(), "Failed to get current directory.", { logContext: "project" }),
  pathExists: (path) => apiCallWithDefault(() => PathExists(path), false, "project")
};
const reportsApi = {
  generateReport: (contextId, format) => apiCall(
    () => GenerateReport(contextId, format),
    "Failed to generate report.",
    { logContext: "reports" }
  ),
  listReports: (projectPath) => apiCall(() => ListReports(projectPath), "Failed to list reports.", { logContext: "reports" }),
  getReport: (reportId) => apiCall(() => GetReport(reportId), "Failed to get report.", { logContext: "reports" }),
  exportProject: (projectPath, format, outputPath) => apiCall(
    () => ExportProject(projectPath, format, outputPath),
    "Failed to export project.",
    { logContext: "reports" }
  )
};
const sandboxApi = {
  getChanges: async () => {
    const result = await GetSandboxChanges();
    return JSON.parse(result);
  },
  getDiff: async (path) => {
    return await GetSandboxDiff(path);
  },
  getAllDiffs: async () => {
    return await GetSandboxAllDiffs();
  },
  applyChanges: async () => {
    await ApplySandboxChanges();
  },
  discardChanges: async () => {
    await DiscardSandboxChanges();
  },
  discardFile: async (path) => {
    await DiscardSandboxFile(path);
  },
  hasChanges: async () => {
    return await HasSandboxChanges();
  },
  getChangeCount: async () => {
    return await GetSandboxChangeCount();
  }
};
const semanticApi = {
  isAvailable: () => apiCallWithDefault(
    // @ts-ignore - method may not exist in wails bindings yet
    () => IsSemanticSearchAvailable(),
    false,
    "semantic"
  ),
  search: async (request) => {
    const result = await apiCall(
      // @ts-ignore
      () => SemanticSearch(JSON.stringify(request)),
      "Failed to perform semantic search.",
      { logContext: "semantic" }
    );
    return parseJsonResponse(result, "Failed to parse semantic search response.");
  },
  findSimilar: async (request) => {
    const result = await apiCall(
      // @ts-ignore
      () => SemanticFindSimilar(JSON.stringify(request)),
      "Failed to find similar code.",
      { logContext: "semantic" }
    );
    return parseJsonResponse(result, "Failed to parse similar code response.");
  },
  indexProject: (projectRoot) => apiCall(
    // @ts-ignore
    () => SemanticIndexProject(projectRoot),
    "Failed to index project.",
    { logContext: "semantic" }
  ),
  indexFile: (projectRoot, filePath) => apiCall(
    // @ts-ignore
    () => SemanticIndexFile(projectRoot, filePath),
    "Failed to index file.",
    { logContext: "semantic" }
  ),
  getStats: async (projectRoot) => {
    const result = await apiCall(
      // @ts-ignore
      () => SemanticGetStats(projectRoot),
      "Failed to get semantic search statistics.",
      { logContext: "semantic" }
    );
    return parseJsonResponse(result, "Failed to parse semantic stats.");
  },
  isIndexed: (projectRoot) => apiCallWithDefault(
    // @ts-ignore
    () => SemanticIsIndexed(projectRoot),
    false,
    "semantic"
  ),
  retrieveContext: async (request) => {
    const result = await apiCall(
      // @ts-ignore
      () => SemanticRetrieveContext(JSON.stringify(request)),
      "Failed to retrieve context.",
      { logContext: "semantic" }
    );
    return parseJsonResponse(result, "Failed to parse retrieved context.");
  },
  hybridSearch: async (request) => {
    const result = await apiCall(
      // @ts-ignore
      () => SemanticHybridSearch(JSON.stringify(request)),
      "Failed to perform hybrid search.",
      { logContext: "semantic" }
    );
    return parseJsonResponse(result, "Failed to parse hybrid search response.");
  }
};
const settingsApi = {
  getSettings: () => apiCall(() => GetSettings(), "Failed to load settings.", { logContext: "settings" }),
  saveSettings: (settings2) => apiCall(() => SaveSettings(settings2), "Failed to save settings.", { logContext: "settings" }),
  // Ignore Rules
  getGitignoreContent: (projectPath) => apiCall(
    () => GetGitignoreContentForProject(projectPath),
    "Failed to load .gitignore content.",
    { logContext: "settings" }
  ),
  getCustomIgnoreRules: () => apiCall(() => GetCustomIgnoreRules(), "Failed to load custom ignore rules.", { logContext: "settings" }),
  updateCustomIgnoreRules: (rules) => apiCall(
    () => UpdateCustomIgnoreRules(rules),
    "Failed to update custom ignore rules.",
    { logContext: "settings" }
  ),
  testIgnoreRules: (projectPath, rules) => apiCall(
    () => TestIgnoreRules(projectPath, rules),
    "Failed to test ignore rules.",
    { logContext: "settings" }
  ),
  testIgnoreRulesDetailed: (projectPath, rules) => apiCall(
    () => TestIgnoreRulesDetailed(projectPath, rules),
    "Failed to test ignore rules.",
    { logContext: "settings" }
  ),
  addToGitignore: (projectPath, pattern) => apiCall(
    () => AddToGitignore(projectPath, pattern),
    "Failed to add to .gitignore.",
    { logContext: "settings" }
  )
};
const taskflowApi = {
  // Task Protocol
  executeTaskProtocol: (configPath) => apiCall(
    () => ExecuteTaskProtocol(configPath),
    "Failed to execute task protocol.",
    { logContext: "taskflow" }
  ),
  getTaskProtocolConfiguration: (projectPath, languages) => apiCall(
    () => GetTaskProtocolConfiguration(projectPath, languages),
    "Failed to get task protocol configuration.",
    { logContext: "taskflow" }
  ),
  // Guardrails
  validatePath: (path) => apiCall(() => ValidatePath(path), "Failed to validate path.", { logContext: "taskflow" }),
  getGuardrailPolicies: () => apiCall(
    () => GetGuardrailPolicies(),
    "Failed to get guardrail policies.",
    { logContext: "taskflow" }
  ),
  getBudgetPolicies: () => apiCall(
    // @ts-ignore
    () => GetBudgetPolicies(),
    "Failed to get budget policies.",
    { logContext: "taskflow" }
  )
};
const apiService = {
  // ============================================
  // Project Management
  // ============================================
  getRecentProjects: projectApi.getRecentProjects,
  addRecentProject: projectApi.addRecentProject,
  removeRecentProject: projectApi.removeRecentProject,
  selectDirectory: projectApi.selectDirectory,
  getCurrentDirectory: projectApi.getCurrentDirectory,
  pathExists: projectApi.pathExists,
  // ============================================
  // File Operations
  // ============================================
  listFiles: filesApi$1.listFiles,
  clearFileTreeCache: filesApi$1.clearFileTreeCache,
  readFileContent: filesApi$1.readFileContent,
  getFileStats: filesApi$1.getFileStats,
  // ============================================
  // Context
  // ============================================
  buildContext: contextApi$1.buildContext,
  buildContextFromRequest: contextApi$1.buildContextFromRequest,
  getContext: contextApi$1.getContext,
  deleteContext: contextApi$1.deleteContext,
  getProjectContexts: contextApi$1.getProjectContexts,
  exportContext: contextApi$1.exportContext,
  getFullContextContent: contextApi$1.getFullContextContent,
  suggestContextFiles: contextApi$1.suggestContextFiles,
  getSmartSuggestions: contextApi$1.getSmartSuggestions,
  getFileQuickInfo: contextApi$1.getFileQuickInfo,
  getImpactPreview: contextApi$1.getImpactPreview,
  analyzeTaskAndCollectContext: contextApi$1.analyzeTaskAndCollectContext,
  agenticChat: contextApi$1.agenticChat,
  collectSmartContext: contextApi$1.collectSmartContext,
  // ============================================
  // AI and Code Generation
  // ============================================
  generateCode: aiApi.generateCode,
  generateCodeStream: aiApi.generateCodeStream,
  generateIntelligentCode: aiApi.generateIntelligentCode,
  listAvailableModels: aiApi.listAvailableModels,
  getProviderInfo: aiApi.getProviderInfo,
  qwenExecuteTask: aiApi.qwenExecuteTask,
  qwenPreviewContext: aiApi.qwenPreviewContext,
  qwenGetAvailableModels: aiApi.qwenGetAvailableModels,
  // ============================================
  // Analysis
  // ============================================
  analyzeProject: analysisApi.analyzeProject,
  analyzeFile: analysisApi.analyzeFile,
  detectLanguages: analysisApi.detectLanguages,
  getSupportedAnalyzers: analysisApi.getSupportedAnalyzers,
  // ============================================
  // Git Operations
  // ============================================
  getUncommittedFiles: gitApi$1.getUncommittedFiles,
  getBranches: gitApi$1.getBranches,
  getCurrentBranch: gitApi$1.getCurrentBranch,
  getRichCommitHistory: gitApi$1.getRichCommitHistory,
  isGitAvailable: gitApi$1.isGitAvailable,
  isGitRepository: gitApi$1.isGitRepository,
  cloneRepository: gitApi$1.cloneRepository,
  checkoutBranch: gitApi$1.checkoutBranch,
  checkoutCommit: gitApi$1.checkoutCommit,
  getCommitHistory: gitApi$1.getCommitHistory,
  getRemoteBranches: gitApi$1.getRemoteBranches,
  cleanupTempRepository: gitApi$1.cleanupTempRepository,
  listFilesAtRef: gitApi$1.listFilesAtRef,
  getFileAtRef: gitApi$1.getFileAtRef,
  buildContextAtRef: gitApi$1.buildContextAtRef,
  // ============================================
  // GitHub API
  // ============================================
  isGitHubURL: githubApi.isGitHubURL,
  gitHubGetDefaultBranch: githubApi.getDefaultBranch,
  gitHubGetBranches: githubApi.getBranches,
  gitHubGetCommits: githubApi.getCommits,
  gitHubListFiles: githubApi.listFiles,
  gitHubGetFileContent: githubApi.getFileContent,
  gitHubBuildContext: githubApi.buildContext,
  // ============================================
  // GitLab API
  // ============================================
  isGitLabURL: gitlabApi.isGitLabURL,
  gitLabGetDefaultBranch: gitlabApi.getDefaultBranch,
  gitLabGetBranches: gitlabApi.getBranches,
  gitLabGetCommits: gitlabApi.getCommits,
  gitLabListFiles: gitlabApi.listFiles,
  gitLabGetFileContent: gitlabApi.getFileContent,
  gitLabBuildContext: gitlabApi.buildContext,
  // ============================================
  // Settings
  // ============================================
  getSettings: settingsApi.getSettings,
  saveSettings: settingsApi.saveSettings,
  getGitignoreContent: settingsApi.getGitignoreContent,
  getCustomIgnoreRules: settingsApi.getCustomIgnoreRules,
  updateCustomIgnoreRules: settingsApi.updateCustomIgnoreRules,
  testIgnoreRules: settingsApi.testIgnoreRules,
  testIgnoreRulesDetailed: settingsApi.testIgnoreRulesDetailed,
  addToGitignore: settingsApi.addToGitignore,
  // ============================================
  // Build and Test
  // ============================================
  runTests: buildApi.runTests,
  discoverTests: buildApi.discoverTests,
  build: buildApi.build,
  typeCheck: buildApi.typeCheck,
  generateDiff: buildApi.generateDiff,
  applyEdits: buildApi.applyEdits,
  applySingleEdit: buildApi.applySingleEdit,
  // ============================================
  // Reports
  // ============================================
  generateReport: reportsApi.generateReport,
  listReports: reportsApi.listReports,
  getReport: reportsApi.getReport,
  exportProject: reportsApi.exportProject,
  // ============================================
  // Task Protocol and Guardrails
  // ============================================
  executeTaskProtocol: taskflowApi.executeTaskProtocol,
  getTaskProtocolConfiguration: taskflowApi.getTaskProtocolConfiguration,
  validatePath: taskflowApi.validatePath,
  getGuardrailPolicies: taskflowApi.getGuardrailPolicies,
  getBudgetPolicies: taskflowApi.getBudgetPolicies,
  // ============================================
  // Semantic Search
  // ============================================
  isSemanticSearchAvailable: semanticApi.isAvailable,
  semanticSearch: semanticApi.search,
  semanticFindSimilar: semanticApi.findSimilar,
  semanticIndexProject: semanticApi.indexProject,
  semanticIndexFile: semanticApi.indexFile,
  semanticGetStats: semanticApi.getStats,
  semanticIsIndexed: semanticApi.isIndexed,
  semanticRetrieveContext: semanticApi.retrieveContext,
  semanticHybridSearch: semanticApi.hybridSearch,
  // ============================================
  // Context Memory
  // ============================================
  getRecentContexts: memoryApi.getRecentContexts,
  findContextByTopic: memoryApi.findContextByTopic,
  saveContextMemory: memoryApi.saveContextMemory,
  // ============================================
  // Sandbox (AI File Changes)
  // ============================================
  getSandboxChanges: sandboxApi.getChanges,
  getSandboxDiff: sandboxApi.getDiff,
  getSandboxAllDiffs: sandboxApi.getAllDiffs,
  applySandboxChanges: sandboxApi.applyChanges,
  discardSandboxChanges: sandboxApi.discardChanges,
  discardSandboxFile: sandboxApi.discardFile,
  hasSandboxChanges: sandboxApi.hasChanges,
  getSandboxChangeCount: sandboxApi.getChangeCount
};
const logger$k = useLogger("ContextApi");
class ContextApi {
  async buildContext(projectPath, files2, options2) {
    try {
      return await apiService.buildContextFromRequest(projectPath, files2, options2);
    } catch (error) {
      logger$k.error("Failed to build context:", error);
      if (error instanceof Error && error.message === "TOKEN_LIMIT_EXCEEDED") {
        throw error;
      }
      throw new Error("Failed to build context. Please check your file selection.");
    }
  }
  async getContextContent(contextId) {
    try {
      return await apiService.getFullContextContent(contextId);
    } catch (error) {
      logger$k.error("Failed to get context content:", error);
      throw new Error("Failed to load context content.");
    }
  }
  async deleteContext(contextId) {
    try {
      await apiService.deleteContext(contextId);
    } catch (error) {
      logger$k.error("Failed to delete context:", error);
      throw new Error("Failed to delete context.");
    }
  }
  async exportContext(exportSettings) {
    try {
      return await apiService.exportContext(exportSettings);
    } catch (error) {
      logger$k.error("Failed to export context:", error);
      throw new Error("Failed to export context.");
    }
  }
  async getProjectContexts(projectPath) {
    try {
      const result = await apiService.getProjectContexts(projectPath);
      const parsed = JSON.parse(result);
      if (!Array.isArray(parsed)) return [];
      return parsed.map((ctx) => ({
        ...ctx,
        files: ctx.metadata?.selectedFiles || ctx.files || []
      }));
    } catch (error) {
      logger$k.error("Failed to get project contexts:", error);
      throw new Error("Failed to load project contexts.");
    }
  }
}
const contextApi = new ContextApi();
function formatContextSize(bytes) {
  if (bytes === 0) return "0 B";
  const k = 1024;
  const sizes = ["B", "KB", "MB", "GB"];
  const i = Math.floor(Math.log(bytes) / Math.log(k));
  return Math.round(bytes / Math.pow(k, i) * 100) / 100 + " " + sizes[i];
}
function formatTimestamp(isoString) {
  const date = new Date(isoString);
  const now = /* @__PURE__ */ new Date();
  const diffMs = now.getTime() - date.getTime();
  const diffMins = Math.floor(diffMs / 6e4);
  if (diffMins < 1) return "just now";
  if (diffMins < 60) return `${diffMins}m ago`;
  const diffHours = Math.floor(diffMins / 60);
  if (diffHours < 24) return `${diffHours}h ago`;
  const diffDays = Math.floor(diffHours / 24);
  if (diffDays < 7) return `${diffDays}d ago`;
  return date.toLocaleDateString();
}
function LogPrint(message) {
  window.runtime.LogPrint(message);
}
function LogTrace(message) {
  window.runtime.LogTrace(message);
}
function LogDebug(message) {
  window.runtime.LogDebug(message);
}
function LogInfo(message) {
  window.runtime.LogInfo(message);
}
function LogWarning(message) {
  window.runtime.LogWarning(message);
}
function LogError(message) {
  window.runtime.LogError(message);
}
function LogFatal(message) {
  window.runtime.LogFatal(message);
}
function EventsOnMultiple(eventName, callback, maxCallbacks) {
  return window.runtime.EventsOnMultiple(eventName, callback, maxCallbacks);
}
function EventsOn(eventName, callback) {
  return EventsOnMultiple(eventName, callback, -1);
}
function EventsOff(eventName, ...additionalEventNames) {
  return window.runtime.EventsOff(eventName, ...additionalEventNames);
}
function EventsOffAll() {
  return window.runtime.EventsOffAll();
}
function EventsOnce(eventName, callback) {
  return EventsOnMultiple(eventName, callback, 1);
}
function EventsEmit(eventName) {
  let args = [eventName].slice.call(arguments);
  return window.runtime.EventsEmit.apply(null, args);
}
function WindowReload() {
  window.runtime.WindowReload();
}
function WindowReloadApp() {
  window.runtime.WindowReloadApp();
}
function WindowSetAlwaysOnTop(b) {
  window.runtime.WindowSetAlwaysOnTop(b);
}
function WindowSetSystemDefaultTheme() {
  window.runtime.WindowSetSystemDefaultTheme();
}
function WindowSetLightTheme() {
  window.runtime.WindowSetLightTheme();
}
function WindowSetDarkTheme() {
  window.runtime.WindowSetDarkTheme();
}
function WindowCenter() {
  window.runtime.WindowCenter();
}
function WindowSetTitle(title) {
  window.runtime.WindowSetTitle(title);
}
function WindowFullscreen() {
  window.runtime.WindowFullscreen();
}
function WindowUnfullscreen() {
  window.runtime.WindowUnfullscreen();
}
function WindowIsFullscreen() {
  return window.runtime.WindowIsFullscreen();
}
function WindowGetSize() {
  return window.runtime.WindowGetSize();
}
function WindowSetSize(width, height) {
  window.runtime.WindowSetSize(width, height);
}
function WindowSetMaxSize(width, height) {
  window.runtime.WindowSetMaxSize(width, height);
}
function WindowSetMinSize(width, height) {
  window.runtime.WindowSetMinSize(width, height);
}
function WindowSetPosition(x, y) {
  window.runtime.WindowSetPosition(x, y);
}
function WindowGetPosition() {
  return window.runtime.WindowGetPosition();
}
function WindowHide() {
  window.runtime.WindowHide();
}
function WindowShow() {
  window.runtime.WindowShow();
}
function WindowMaximise() {
  window.runtime.WindowMaximise();
}
function WindowToggleMaximise() {
  window.runtime.WindowToggleMaximise();
}
function WindowUnmaximise() {
  window.runtime.WindowUnmaximise();
}
function WindowIsMaximised() {
  return window.runtime.WindowIsMaximised();
}
function WindowMinimise() {
  window.runtime.WindowMinimise();
}
function WindowUnminimise() {
  window.runtime.WindowUnminimise();
}
function WindowSetBackgroundColour(R, G, B, A) {
  window.runtime.WindowSetBackgroundColour(R, G, B, A);
}
function ScreenGetAll() {
  return window.runtime.ScreenGetAll();
}
function WindowIsMinimised() {
  return window.runtime.WindowIsMinimised();
}
function WindowIsNormal() {
  return window.runtime.WindowIsNormal();
}
function BrowserOpenURL(url) {
  window.runtime.BrowserOpenURL(url);
}
function Environment() {
  return window.runtime.Environment();
}
function Quit() {
  window.runtime.Quit();
}
function Hide() {
  window.runtime.Hide();
}
function Show() {
  window.runtime.Show();
}
function ClipboardGetText() {
  return window.runtime.ClipboardGetText();
}
function ClipboardSetText(text) {
  return window.runtime.ClipboardSetText(text);
}
function OnFileDrop(callback, useDropTarget) {
  return window.runtime.OnFileDrop(callback, useDropTarget);
}
function OnFileDropOff() {
  return window.runtime.OnFileDropOff();
}
function CanResolveFilePaths() {
  return window.runtime.CanResolveFilePaths();
}
function ResolveFilePaths(files2) {
  return window.runtime.ResolveFilePaths(files2);
}
const runtime = /* @__PURE__ */ Object.freeze(/* @__PURE__ */ Object.defineProperty({
  __proto__: null,
  BrowserOpenURL,
  CanResolveFilePaths,
  ClipboardGetText,
  ClipboardSetText,
  Environment,
  EventsEmit,
  EventsOff,
  EventsOffAll,
  EventsOn,
  EventsOnMultiple,
  EventsOnce,
  Hide,
  LogDebug,
  LogError,
  LogFatal,
  LogInfo,
  LogPrint,
  LogTrace,
  LogWarning,
  OnFileDrop,
  OnFileDropOff,
  Quit,
  ResolveFilePaths,
  ScreenGetAll,
  Show,
  WindowCenter,
  WindowFullscreen,
  WindowGetPosition,
  WindowGetSize,
  WindowHide,
  WindowIsFullscreen,
  WindowIsMaximised,
  WindowIsMinimised,
  WindowIsNormal,
  WindowMaximise,
  WindowMinimise,
  WindowReload,
  WindowReloadApp,
  WindowSetAlwaysOnTop,
  WindowSetBackgroundColour,
  WindowSetDarkTheme,
  WindowSetLightTheme,
  WindowSetMaxSize,
  WindowSetMinSize,
  WindowSetPosition,
  WindowSetSize,
  WindowSetSystemDefaultTheme,
  WindowSetTitle,
  WindowShow,
  WindowToggleMaximise,
  WindowUnfullscreen,
  WindowUnmaximise,
  WindowUnminimise
}, Symbol.toStringTag, { value: "Module" }));
const DEFAULT_SETTINGS = {
  context: {
    maxTokens: 1e5,
    // 100K tokens by default (reasonable for most projects)
    stripComments: false,
    includeTests: true,
    splitStrategy: "smart",
    outputFormat: "xml",
    // Default format - XML is best for AI context
    // Content optimization - disabled by default for safety
    excludeTests: false,
    collapseEmptyLines: false,
    stripLicense: false,
    compactDataFiles: false,
    trimWhitespace: false,
    // Export options
    includeManifest: true,
    includeLineNumbers: false,
    enableAutoSplit: false,
    maxTokensPerChunk: 32e3,
    // Template options
    applyTemplateOnCopy: true
  },
  contextStorage: {
    maxContexts: 20,
    maxStorageMB: 100,
    autoCleanupDays: 30,
    autoCleanupOnLimit: true
  },
  fileExplorer: {
    useGitignore: true,
    useCustomIgnore: false,
    autoSaveSelection: true,
    compactNestedFolders: true,
    showIgnoredFiles: true,
    foldersFirst: true,
    allowSelectBinary: false,
    customIgnoreRules: "",
    quickFilters: [
      { id: "source", label: "Исходники", extensions: [".ts", ".js", ".tsx", ".jsx", ".vue", ".go", ".py", ".java", ".cpp", ".c", ".rs"], patterns: [], enabled: true },
      { id: "tests", label: "Тесты", extensions: [], patterns: ["**/*.test.*", "**/*.spec.*", "**/test/**", "**/tests/**", "**/Test/**"], enabled: true },
      { id: "config", label: "Конфигурация", extensions: [".json", ".yaml", ".yml", ".toml", ".ini", ".env"], patterns: [], enabled: true },
      { id: "docs", label: "Документация", extensions: [".md", ".txt", ".rst", ".adoc"], patterns: [], enabled: true },
      { id: "styles", label: "Стили", extensions: [".css", ".scss", ".sass", ".less"], patterns: [], enabled: true }
    ]
  },
  aiModel: "gpt-4",
  theme: "dark"
};
const useSettingsStore = defineStore("settings", () => {
  const settings2 = ref(loadSettings());
  watch(
    settings2,
    (newSettings) => {
      saveSettings(newSettings);
    },
    { deep: true }
  );
  function loadSettings() {
    try {
      const saved = localStorage.getItem("app-settings");
      if (saved) {
        if (saved === "undefined" || saved === "null" || saved.trim() === "") {
          localStorage.removeItem("app-settings");
          return DEFAULT_SETTINGS;
        }
        const parsed = JSON.parse(saved);
        if (typeof parsed !== "object" || parsed === null) {
          localStorage.removeItem("app-settings");
          return DEFAULT_SETTINGS;
        }
        return {
          ...DEFAULT_SETTINGS,
          ...parsed,
          context: {
            ...DEFAULT_SETTINGS.context,
            ...parsed.context || {}
          },
          contextStorage: {
            ...DEFAULT_SETTINGS.contextStorage,
            ...parsed.contextStorage || {}
          },
          fileExplorer: {
            ...DEFAULT_SETTINGS.fileExplorer,
            ...parsed.fileExplorer || {}
          }
        };
      }
    } catch (err) {
      try {
        localStorage.removeItem("app-settings");
      } catch {
      }
    }
    return DEFAULT_SETTINGS;
  }
  function saveSettings(settings22) {
    try {
      localStorage.setItem("app-settings", JSON.stringify(settings22));
    } catch (err) {
      console.warn("Failed to save settings:", err);
    }
  }
  function resetToDefaults() {
    settings2.value = JSON.parse(JSON.stringify(DEFAULT_SETTINGS));
  }
  function updateContextSettings(updates) {
    settings2.value.context = {
      ...settings2.value.context,
      ...updates
    };
  }
  function updateAIModel(model) {
    settings2.value.aiModel = model;
  }
  function updateTheme(theme) {
    settings2.value.theme = theme;
  }
  function updateFileExplorerSettings(updates) {
    settings2.value.fileExplorer = {
      ...settings2.value.fileExplorer,
      ...updates
    };
  }
  function getCustomIgnoreRules() {
    return settings2.value.fileExplorer.customIgnoreRules;
  }
  function setCustomIgnoreRules(rules) {
    settings2.value.fileExplorer.customIgnoreRules = rules;
  }
  function updateContextStorageSettings(updates) {
    settings2.value.contextStorage = {
      ...settings2.value.contextStorage,
      ...updates
    };
  }
  return {
    settings: settings2,
    resetToDefaults,
    updateContextSettings,
    updateContextStorageSettings,
    updateAIModel,
    updateTheme,
    updateFileExplorerSettings,
    getCustomIgnoreRules,
    setCustomIgnoreRules
  };
});
const logger$j = useLogger("ContextCleanup");
function getNonFavoritesSorted(contexts) {
  return contexts.filter((c) => !c.isFavorite).sort((a, b) => {
    const dateA = new Date(a.createdAt || 0).getTime();
    const dateB = new Date(b.createdAt || 0).getTime();
    return dateA - dateB;
  });
}
function getContextsToDeleteByLimit(contexts, maxContexts) {
  const totalCount = contexts.length;
  if (totalCount <= maxContexts) return [];
  const nonFavorites = getNonFavoritesSorted(contexts);
  const toDelete = totalCount - maxContexts;
  return nonFavorites.slice(0, toDelete);
}
function getContextsToDeleteByAge(contexts, maxAgeDays) {
  if (maxAgeDays <= 0) return [];
  const cutoffDate = Date.now() - maxAgeDays * 24 * 60 * 60 * 1e3;
  const nonFavorites = getNonFavoritesSorted(contexts);
  return nonFavorites.filter((ctx) => {
    const ctxDate = new Date(ctx.createdAt || 0).getTime();
    return ctxDate < cutoffDate;
  });
}
async function performAutoCleanup(contexts, deleteContext) {
  const settingsStore = useSettingsStore();
  const storage = settingsStore.settings.contextStorage;
  const result = {
    deletedByLimit: [],
    deletedByAge: []
  };
  if (!storage.autoCleanupOnLimit) {
    return result;
  }
  const byLimit = getContextsToDeleteByLimit(contexts, storage.maxContexts);
  for (const ctx of byLimit) {
    try {
      await deleteContext(ctx.id);
      result.deletedByLimit.push(ctx.id);
      logger$j.debug("Auto-deleted old context:", ctx.id);
    } catch (err) {
      logger$j.error("Failed to auto-delete context:", err);
    }
  }
  const byAge = getContextsToDeleteByAge(contexts, storage.autoCleanupDays);
  for (const ctx of byAge) {
    if (result.deletedByLimit.includes(ctx.id)) continue;
    try {
      await deleteContext(ctx.id);
      result.deletedByAge.push(ctx.id);
      logger$j.debug("Auto-deleted expired context:", ctx.id);
    } catch (err) {
      logger$j.error("Failed to auto-delete expired context:", err);
    }
  }
  return result;
}
const logger$i = useLogger("ContextContent");
const MAX_CONTEXT_SIZE = 20 * 1024 * 1024;
function useContextContent() {
  const cache2 = ref({
    contextId: null,
    content: null
  });
  function clearCache() {
    cache2.value = { contextId: null, content: null };
  }
  function isCached(contextId) {
    return cache2.value.contextId === contextId && cache2.value.content !== null;
  }
  function getCachedContent(contextId) {
    if (isCached(contextId)) {
      return cache2.value.content;
    }
    return null;
  }
  async function loadFullContent(contextId) {
    const cached = getCachedContent(contextId);
    if (cached !== null) {
      return cached;
    }
    const content = await contextApi.getContextContent(contextId);
    if (content.length > MAX_CONTEXT_SIZE) {
      const sizeMB = Math.round(content.length / (1024 * 1024));
      const limitMB = Math.round(MAX_CONTEXT_SIZE / (1024 * 1024));
      throw new Error(`Context content (${sizeMB}MB) exceeds maximum allowed size (${limitMB}MB)`);
    }
    cache2.value = { contextId, content };
    logger$i.debug("Content cached, size:", Math.round(content.length / 1024), "KB");
    return content;
  }
  async function loadChunk(contextId, startLine = 0, lineCount = 0) {
    const content = await loadFullContent(contextId);
    const lines = content.split("\n");
    const endLine = lineCount > 0 ? Math.min(startLine + lineCount, lines.length) : lines.length;
    logger$i.debug("Loaded chunk:", startLine, "-", endLine, "of", lines.length, "lines");
    return {
      lines: lines.slice(startLine, endLine),
      startLine,
      endLine,
      hasMore: endLine < lines.length
    };
  }
  function getMemoryUsage() {
    if (!cache2.value.content) return 0;
    return cache2.value.content.length * 2;
  }
  return {
    cache: cache2,
    clearCache,
    isCached,
    getCachedContent,
    loadFullContent,
    loadChunk,
    getMemoryUsage
  };
}
const FILE_TREE = {
  MAX_FLATTEN_DEPTH: 5,
  DEBOUNCE_MS: 300
};
const TOKEN_THRESHOLDS = {
  MEDIUM: 5e3,
  // 5k tokens - yellow indicator
  HEAVY: 2e4,
  // 20k tokens - orange indicator  
  CRITICAL: 5e4,
  // 50k tokens - red indicator
  BYTES_PER_TOKEN: 4
  // approximate bytes per token
};
const STORAGE_KEYS = {
  // Project
  RECENT_PROJECTS: "Syntaxia_recent_projects",
  AUTO_OPEN_LAST: "Syntaxia_auto_open_last",
  // Settings
  APP_SETTINGS: "app-settings",
  // Context
  CONTEXT_METADATA: "context-metadata",
  // Chat
  CHAT_HISTORY_PREFIX: "chat-history",
  // UI State
  LEFT_SIDEBAR_TAB: "left-sidebar-tab",
  RIGHT_SIDEBAR_VISIBLE: "right-sidebar-visible",
  RIGHT_SIDEBAR_TAB: "right-sidebar-tab",
  TASK_PANEL_VISIBLE: "task-panel-visible",
  // File Explorer
  EXPANDED_NODES_PREFIX: "expanded-nodes",
  SELECTED_FILES_PREFIX: "selected-files",
  FILTER_STATE_PREFIX: "filter-state",
  // Workspace Layout
  WORKSPACE_LEFT_WIDTH: "workspace-left-width",
  WORKSPACE_RIGHT_WIDTH: "workspace-right-width",
  // Onboarding
  ONBOARDING_COMPLETED: "Syntaxia-onboarding-completed"
};
const logger$h = useLogger("ContextMetadata");
function saveContextMetadata(contexts) {
  try {
    const metadata = contexts.map((c) => ({
      id: c.id,
      name: c.name,
      isFavorite: c.isFavorite
    }));
    localStorage.setItem(STORAGE_KEYS.CONTEXT_METADATA, JSON.stringify(metadata));
  } catch (err) {
    logger$h.warn("Failed to save:", err);
  }
}
function loadContextMetadata(contexts) {
  try {
    const saved = localStorage.getItem(STORAGE_KEYS.CONTEXT_METADATA);
    if (!saved) return;
    const metadata = JSON.parse(saved);
    for (const ctx of contexts) {
      const meta = metadata.find((m) => m.id === ctx.id);
      if (meta) {
        ctx.name = meta.name;
        ctx.isFavorite = meta.isFavorite;
      }
    }
  } catch (err) {
    logger$h.warn("Failed to load:", err);
  }
}
const logger$g = useLogger("ContextOperations");
async function duplicateContextContent(ctxId, contexts) {
  try {
    const content = await contextApi.getContextContent(ctxId);
    const original = contexts.find((c) => c.id === ctxId);
    if (!original) {
      logger$g.error("Original context not found:", ctxId);
      return null;
    }
    const newId = `ctx-copy-${Date.now()}`;
    const newName = `${original.name} (копия)`;
    const lines = content.split("\n");
    return {
      id: newId,
      name: newName,
      content,
      fileCount: original.fileCount,
      totalSize: content.length,
      lineCount: lines.length,
      tokenCount: Math.round(content.length / 4)
    };
  } catch (err) {
    logger$g.error("Failed to duplicate context:", err);
    return null;
  }
}
async function mergeContextsContent(contextIds, contexts, loadContent) {
  if (contextIds.length < 2) {
    logger$g.warn("Need at least 2 contexts to merge");
    return null;
  }
  try {
    const contents = [];
    let totalFiles = 0;
    const names = [];
    for (const ctxId of contextIds) {
      const content = await loadContent(ctxId);
      contents.push(content);
      const ctx = contexts.find((c) => c.id === ctxId);
      if (ctx) {
        totalFiles += ctx.fileCount;
        names.push(ctx.name || ctx.id);
      }
    }
    const separator = "\n\n" + "=".repeat(80) + "\n\n";
    const mergedContent = contents.join(separator);
    const newId = `merged-${Date.now()}`;
    const namePreview = names.slice(0, 2).join(" + ");
    const suffix = names.length > 2 ? ` +${names.length - 2}` : "";
    const newName = `Merged: ${namePreview}${suffix}`;
    return {
      id: newId,
      name: newName,
      content: mergedContent,
      fileCount: totalFiles
    };
  } catch (err) {
    logger$g.error("Failed to merge contexts:", err);
    return null;
  }
}
function reorderContextList(contexts, fromIndex, toIndex) {
  const list = [...contexts];
  const [removed] = list.splice(fromIndex, 1);
  list.splice(toIndex, 0, removed);
  return list;
}
function createContextSummary(id, name, content, fileCount) {
  const lines = content.split("\n");
  return {
    id,
    name,
    fileCount,
    totalSize: content.length,
    lineCount: lines.length,
    tokenCount: Math.round(content.length / 4),
    createdAt: (/* @__PURE__ */ new Date()).toISOString(),
    isFavorite: false
  };
}
const highlightCache = /* @__PURE__ */ new Map();
function clearHighlightCache() {
  highlightCache.clear();
}
const TYPE_LABELS = {
  "vue": "Vue компоненты",
  "ts": "TypeScript",
  "tsx": "React TSX",
  "js": "JavaScript",
  "jsx": "React JSX",
  "go": "Go код",
  "py": "Python",
  "java": "Java",
  "css": "Стили",
  "scss": "SCSS стили",
  "json": "Конфигурация",
  "yaml": "YAML конфиг",
  "yml": "YAML конфиг",
  "md": "Документация",
  "sql": "SQL запросы",
  "html": "HTML шаблоны",
  "test": "Тесты",
  "spec": "Тесты"
};
function analyzeFiles(files2) {
  const extensions = /* @__PURE__ */ new Map();
  const folders = /* @__PURE__ */ new Map();
  const fileNames = [];
  for (const file of files2) {
    const parts = file.split("/");
    const fileName = parts[parts.length - 1];
    fileNames.push(fileName);
    const ext = fileName.includes(".") ? fileName.split(".").pop()?.toLowerCase() : "";
    if (ext) {
      extensions.set(ext, (extensions.get(ext) || 0) + 1);
    }
    if (parts.length > 1) {
      const topFolder = parts[0];
      folders.set(topFolder, (folders.get(topFolder) || 0) + 1);
    }
  }
  const hasTests = fileNames.some(
    (f) => f.includes("test") || f.includes("spec") || f.includes("Test")
  );
  return { extensions, folders, fileNames, hasTests };
}
function generateSmartName(files2) {
  if (!files2 || files2.length === 0) return "Пустой контекст";
  const { extensions, folders, fileNames, hasTests } = analyzeFiles(files2);
  const sortedExts = [...extensions.entries()].sort((a, b) => b[1] - a[1]);
  const sortedFolders = [...folders.entries()].sort((a, b) => b[1] - a[1]);
  if (files2.length === 1) {
    const fileName = fileNames[0];
    return fileName.length > 30 ? fileName.substring(0, 27) + "..." : fileName;
  }
  let name = "";
  if (sortedFolders.length === 1 && sortedFolders[0][1] === files2.length) {
    name = sortedFolders[0][0];
  } else if (sortedExts.length > 0 && sortedExts[0][1] / files2.length > 0.6) {
    const ext = sortedExts[0][0];
    name = TYPE_LABELS[ext] || `.${ext} файлы`;
  } else if (sortedFolders.length > 0 && sortedFolders[0][1] / files2.length > 0.5) {
    name = sortedFolders[0][0];
  } else {
    const topTypes = sortedExts.slice(0, 2).map(([ext]) => TYPE_LABELS[ext] || ext);
    name = topTypes.join(" + ") || "Смешанный";
  }
  if (hasTests && !name.toLowerCase().includes("тест")) {
    name += " (с тестами)";
  }
  name += ` [${files2.length}]`;
  return name;
}
const logger$f = useLogger("ContextStore");
const useContextStore = defineStore("context", () => {
  const contentManager = useContextContent();
  const contextId = ref(null);
  const summary = ref(null);
  const currentChunk = ref(null);
  const isLoading = ref(false);
  const isBuilding = ref(false);
  const buildProgress = ref(0);
  const error = ref(null);
  const contextList = ref([]);
  const warnings = ref([]);
  const skippedFiles = ref([]);
  const selectedListItem = ref(null);
  const hasContext = computed(() => contextId.value !== null);
  const totalSize = computed(() => summary.value?.totalSize || 0);
  const fileCount = computed(() => summary.value?.fileCount || 0);
  const lineCount = computed(() => summary.value?.lineCount || 0);
  const tokenCount = computed(() => summary.value?.tokenCount || Math.round(lineCount.value * 2.5));
  const totalTokens = computed(() => tokenCount.value);
  const estimatedCost = computed(() => tokenCount.value / 1e3 * 2e-3);
  const stats = computed(() => ({
    tokens: tokenCount.value,
    totalFiles: fileCount.value,
    lines: lineCount.value,
    size: totalSize.value
  }));
  async function buildContext(filePaths, options2) {
    if (filePaths.length === 0) {
      error.value = "No files selected";
      return;
    }
    isBuilding.value = true;
    isLoading.value = true;
    buildProgress.value = 5;
    error.value = null;
    let unsubscribeProgress = null;
    const totalFiles = filePaths.length;
    const setupProgressListener = () => {
      unsubscribeProgress = EventsOn("SyntaxiaContextGenerationProgress", (data) => {
        const fileProgress = data.total > 0 ? data.current / data.total : 0;
        buildProgress.value = Math.round(5 + fileProgress * 80);
      });
    };
    const cleanupProgressListener = () => {
      if (unsubscribeProgress) {
        unsubscribeProgress();
        unsubscribeProgress = null;
      }
      EventsOff("SyntaxiaContextGenerationProgress");
    };
    try {
      const projectStore = useProjectStore();
      if (!projectStore.currentPath) {
        error.value = "No project selected";
        throw new Error("No project selected");
      }
      const buildOptions = createBuildOptions(options2);
      logger$f.debug("Building context with options:", {
        outputFormat: buildOptions.outputFormat,
        stripComments: buildOptions.stripComments,
        maxTokens: buildOptions.maxTokens,
        fileCount: totalFiles
      });
      setupProgressListener();
      const result = await contextApi.buildContext(projectStore.currentPath, filePaths, buildOptions);
      buildProgress.value = 90;
      contextId.value = result.id || `ctx-${Date.now()}`;
      summary.value = {
        id: contextId.value,
        name: generateSmartName(filePaths),
        fileCount: result.fileCount || filePaths.length,
        totalSize: result.totalSize || 0,
        lineCount: result.lineCount || 0,
        tokenCount: result.tokenCount,
        createdAt: (/* @__PURE__ */ new Date()).toISOString(),
        files: filePaths,
        isFavorite: false
      };
      validateBuildResult(summary.value, result);
      warnings.value = result.metadata?.warnings || [];
      skippedFiles.value = result.metadata?.skippedFiles || [];
      buildProgress.value = 95;
      contentManager.clearCache();
      if (contextId.value && summary.value.fileCount > 0) {
        await loadContextContent(contextId.value, 0, 0);
      }
      cleanupProgressListener();
      buildProgress.value = 100;
      await new Promise((resolve) => setTimeout(resolve, 200));
    } catch (err) {
      cleanupProgressListener();
      handleBuildError(err);
      throw err;
    } finally {
      cleanupProgressListener();
      isBuilding.value = false;
      isLoading.value = false;
      buildProgress.value = 0;
    }
  }
  function createBuildOptions(options2) {
    const settingsStore = useSettingsStore();
    const contextSettings = settingsStore.settings.context;
    return {
      maxTokens: options2?.maxTokens || contextSettings.maxTokens,
      maxMemoryMB: options2?.maxMemoryMB || 50,
      stripComments: options2?.stripComments ?? contextSettings.stripComments,
      includeManifest: options2?.includeManifest ?? true,
      includeLineNumbers: options2?.includeLineNumbers ?? false,
      includeTests: options2?.includeTests ?? true,
      splitStrategy: options2?.splitStrategy || "smart",
      forceStream: true,
      enableProgressEvents: true,
      outputFormat: options2?.outputFormat || contextSettings.outputFormat,
      excludeTests: options2?.excludeTests ?? contextSettings.excludeTests,
      collapseEmptyLines: options2?.collapseEmptyLines ?? contextSettings.collapseEmptyLines,
      stripLicense: options2?.stripLicense ?? contextSettings.stripLicense,
      compactDataFiles: options2?.compactDataFiles ?? contextSettings.compactDataFiles,
      trimWhitespace: options2?.trimWhitespace ?? contextSettings.trimWhitespace
    };
  }
  function validateBuildResult(sum, result) {
    if (sum.totalSize > MAX_CONTEXT_SIZE) {
      const sizeMB = Math.round(sum.totalSize / (1024 * 1024));
      const limitMB = Math.round(MAX_CONTEXT_SIZE / (1024 * 1024));
      error.value = `Context size (${sizeMB}MB) exceeds maximum (${limitMB}MB)`;
      clearContext();
      throw new Error(error.value);
    }
    if (sum.fileCount === 0 || sum.totalSize === 0 || sum.lineCount === 0) {
      const projectStore = useProjectStore();
      const skipped = result.metadata?.skippedFiles || [];
      const skippedInfo = skipped.length > 0 ? `

Skipped files (${skipped.length}):
${skipped.slice(0, 5).join("\n")}${skipped.length > 5 ? "\n..." : ""}` : "";
      error.value = `Context built with empty content.

Current project: ${projectStore.currentPath}

Files may be outside the project directory.${skippedInfo}`;
      throw new Error(error.value);
    }
  }
  async function handleBuildError(err) {
    if (err instanceof Error && err.message === "TOKEN_LIMIT_EXCEEDED") {
      const { hasTokenInfo } = await __vitePreload(async () => {
        const { hasTokenInfo: hasTokenInfo2 } = await import("./errors-DxKjw5Ye.js");
        return { hasTokenInfo: hasTokenInfo2 };
      }, true ? [] : void 0);
      const settingsStore = useSettingsStore();
      const currentLimit = settingsStore.settings.context.maxTokens;
      if (hasTokenInfo(err)) {
        error.value = `TOKEN_LIMIT_EXCEEDED:${err.tokenInfo.actual}:${err.tokenInfo.limit}`;
      } else {
        error.value = `TOKEN_LIMIT_EXCEEDED:0:${currentLimit}`;
      }
      return;
    }
    let errorMessage = err instanceof Error ? err.message : "Failed to build context";
    if (errorMessage.includes("path traversal") || errorMessage.includes("invalid path")) {
      errorMessage += " Make sure all selected files are within the project directory.";
    }
    error.value = errorMessage;
  }
  async function loadContextContent(ctxId, startLine = 0, lineCount2 = 0) {
    if (!ctxId) {
      error.value = "No context ID provided";
      return;
    }
    isLoading.value = true;
    error.value = null;
    try {
      const isNewContext = contextId.value !== ctxId;
      if (isNewContext) {
        contextId.value = ctxId;
        contentManager.clearCache();
      }
      await loadSummaryForContext(ctxId);
      currentChunk.value = await contentManager.loadChunk(ctxId, startLine, lineCount2);
    } catch (err) {
      error.value = err instanceof Error ? err.message : "Failed to load context";
      throw err;
    } finally {
      isLoading.value = false;
    }
  }
  async function loadSummaryForContext(ctxId) {
    const existing = contextList.value.find((c) => c.id === ctxId);
    if (existing) {
      summary.value = { ...existing };
      return;
    }
    try {
      const projectStore = useProjectStore();
      const contexts = await contextApi.getProjectContexts(projectStore.currentPath || "");
      const fetched = contexts.find((c) => c.id === ctxId);
      if (fetched) {
        summary.value = {
          id: fetched.id,
          name: fetched.name || generateSmartName(fetched.files || []),
          fileCount: fetched.fileCount || 0,
          totalSize: fetched.totalSize || 0,
          lineCount: fetched.lineCount || 0,
          tokenCount: fetched.tokenCount,
          createdAt: fetched.createdAt,
          files: fetched.files || [],
          isFavorite: false
        };
      }
    } catch (err) {
      logger$f.warn("Failed to fetch summary from backend:", err);
    }
  }
  async function getFullContextContent() {
    if (!contextId.value) throw new Error("No context ID available");
    return contentManager.loadFullContent(contextId.value);
  }
  async function deleteContext(ctxId) {
    try {
      await contextApi.deleteContext(ctxId);
      if (contextId.value === ctxId) clearContext();
      contextList.value = contextList.value.filter((c) => c.id !== ctxId);
    } catch (err) {
      error.value = err instanceof Error ? err.message : "Failed to delete context";
      throw err;
    }
  }
  async function exportContext(ctxId) {
    if (!ctxId) throw new Error("No context ID provided");
    return contextApi.exportContext({ format: "markdown", includeMetadata: true });
  }
  async function listProjectContexts() {
    isLoading.value = true;
    error.value = null;
    try {
      const projectStore = useProjectStore();
      if (!projectStore.currentPath) throw new Error("No project selected");
      const contexts = await contextApi.getProjectContexts(projectStore.currentPath);
      const contextArray = Array.isArray(contexts) ? contexts : [];
      contextList.value = contextArray.map((ctx) => ({
        id: ctx.id,
        name: ctx.name || generateSmartName(ctx.files || []),
        fileCount: ctx.fileCount || (ctx.files?.length || 0),
        totalSize: ctx.totalSize || 0,
        lineCount: ctx.lineCount || 0,
        tokenCount: ctx.tokenCount,
        createdAt: ctx.createdAt,
        files: ctx.files || [],
        isFavorite: false
      }));
      loadContextMetadata(contextList.value);
      await autoCleanup();
    } catch (err) {
      error.value = err instanceof Error ? err.message : "Failed to list contexts";
      throw err;
    } finally {
      isLoading.value = false;
    }
  }
  function clearContext() {
    contextId.value = null;
    summary.value = null;
    currentChunk.value = null;
    error.value = null;
    contentManager.clearCache();
    clearHighlightCache();
    if (typeof window !== "undefined" && window.gc) {
      try {
        window.gc();
      } catch {
      }
    }
  }
  function setRawContext(content, fileCount2) {
    const id = `git-ref-${Date.now()}`;
    const lines = content.split("\n");
    contextId.value = id;
    summary.value = createContextSummary(id, `Git ref [${fileCount2}]`, content, fileCount2);
    currentChunk.value = {
      lines: lines.slice(0, 100),
      startLine: 0,
      endLine: Math.min(100, lines.length),
      hasMore: lines.length > 100
    };
  }
  function getMemoryUsage() {
    let size = summary.value ? 1024 : 0;
    if (currentChunk.value?.lines) {
      size += currentChunk.value.lines.length * 80 * 2;
    }
    size += contentManager.getMemoryUsage();
    return size;
  }
  function renameContext(ctxId, newName) {
    const ctx = contextList.value.find((c) => c.id === ctxId);
    if (ctx) ctx.name = newName;
    if (summary.value?.id === ctxId) summary.value.name = newName;
    saveMetadata();
  }
  function toggleFavorite(ctxId) {
    const ctx = contextList.value.find((c) => c.id === ctxId);
    if (ctx) ctx.isFavorite = !ctx.isFavorite;
    if (summary.value?.id === ctxId) summary.value.isFavorite = !summary.value.isFavorite;
    saveMetadata();
  }
  function saveMetadata() {
    saveContextMetadata(contextList.value);
  }
  async function autoCleanup() {
    await performAutoCleanup(contextList.value, deleteContext);
  }
  async function duplicateContext(ctxId) {
    const result = await duplicateContextContent(ctxId, contextList.value);
    if (!result) return null;
    setRawContext(result.content, result.fileCount);
    if (summary.value) {
      summary.value.id = result.id;
      summary.value.name = result.name;
    }
    contextList.value.push(createContextSummary(result.id, result.name, result.content, result.fileCount));
    saveMetadata();
    return result.id;
  }
  async function mergeContexts(contextIds) {
    const result = await mergeContextsContent(
      contextIds,
      contextList.value,
      async (ctxId) => {
        if (contextId.value !== ctxId) await loadContextContent(ctxId, 0, 0);
        return getFullContextContent();
      }
    );
    if (!result) return null;
    setRawContext(result.content, result.fileCount);
    if (summary.value) {
      summary.value.id = result.id;
      summary.value.name = result.name;
    }
    contextList.value.push(createContextSummary(result.id, result.name, result.content, result.fileCount));
    saveMetadata();
    return result.id;
  }
  function reorderContexts(fromIndex, toIndex) {
    contextList.value = reorderContextList(contextList.value, fromIndex, toIndex);
    saveMetadata();
  }
  function selectListItem(ctxId) {
    if (!ctxId) {
      selectedListItem.value = null;
      return;
    }
    const item = contextList.value.find((c) => c.id === ctxId);
    selectedListItem.value = item ? { ...item } : null;
  }
  function clearListSelection() {
    selectedListItem.value = null;
  }
  async function removeFileFromContext(filePath) {
    if (!summary.value?.files) return;
    const newFiles = summary.value.files.filter((f) => f !== filePath);
    if (newFiles.length === 0) {
      clearContext();
      return;
    }
    await buildContext(newFiles);
  }
  async function rebuildContext() {
    if (!summary.value?.files || summary.value.files.length === 0) {
      logger$f.warn("No files to rebuild context");
      return;
    }
    const settingsStore = useSettingsStore();
    const options2 = {
      outputFormat: settingsStore.settings.context.outputFormat,
      stripComments: settingsStore.settings.context.stripComments,
      excludeTests: settingsStore.settings.context.excludeTests,
      collapseEmptyLines: settingsStore.settings.context.collapseEmptyLines,
      stripLicense: settingsStore.settings.context.stripLicense,
      compactDataFiles: settingsStore.settings.context.compactDataFiles,
      trimWhitespace: settingsStore.settings.context.trimWhitespace,
      maxTokens: settingsStore.settings.context.maxTokens
    };
    await buildContext(summary.value.files, options2);
  }
  return {
    // State
    contextId,
    summary,
    currentChunk,
    isLoading,
    isBuilding,
    buildProgress,
    error,
    contextList,
    warnings,
    skippedFiles,
    selectedListItem,
    // Computed
    hasContext,
    totalSize,
    fileCount,
    lineCount,
    tokenCount,
    totalTokens,
    estimatedCost,
    stats,
    // Actions
    buildContext,
    rebuildContext,
    loadContextContent,
    deleteContext,
    exportContext,
    listProjectContexts,
    clearContext,
    setRawContext,
    getFullContextContent,
    getMemoryUsage,
    renameContext,
    toggleFavorite,
    duplicateContext,
    autoCleanup,
    loadContextMetadata,
    saveContextMetadata,
    generateSmartName,
    mergeContexts,
    reorderContexts,
    selectListItem,
    clearListSelection,
    removeFileFromContext
  };
});
const context_store = /* @__PURE__ */ Object.freeze(/* @__PURE__ */ Object.defineProperty({
  __proto__: null,
  MAX_CONTEXT_SIZE,
  generateSmartName,
  useContextStore
}, Symbol.toStringTag, { value: "Module" }));
const logger$e = useLogger("FilesApi");
class FilesApi {
  fileTreeCache = /* @__PURE__ */ new Map();
  CACHE_TTL = 6e4;
  // 1 minute
  MAX_CACHE_ENTRIES = 20;
  // Limit cache entries
  async listFiles(path, useGitignore = true, useCustomIgnore = true) {
    const cacheKey = `${path}-${useGitignore}-${useCustomIgnore}`;
    const cached = this.fileTreeCache.get(cacheKey);
    if (cached && Date.now() - cached.timestamp < this.CACHE_TTL) {
      this.fileTreeCache.delete(cacheKey);
      this.fileTreeCache.set(cacheKey, cached);
      return cached.data;
    }
    try {
      const files2 = await apiService.listFiles(path, useGitignore, useCustomIgnore);
      if (this.fileTreeCache.size >= this.MAX_CACHE_ENTRIES) {
        const firstKey = this.fileTreeCache.keys().next().value;
        if (firstKey) this.fileTreeCache.delete(firstKey);
      }
      this.fileTreeCache.set(cacheKey, { data: files2, timestamp: Date.now() });
      return files2;
    } catch (error) {
      logger$e.error("Failed to list files:", error);
      throw new Error("Failed to load file tree.");
    }
  }
  async readFileContent(projectPath, filePath) {
    try {
      return await apiService.readFileContent(projectPath, filePath);
    } catch (error) {
      logger$e.error("Failed to read file:", error);
      throw new Error("Failed to read file content.");
    }
  }
  async getFileStats(path) {
    try {
      const stats = await apiService.getFileStats(path);
      return JSON.parse(stats);
    } catch (error) {
      logger$e.error("Failed to get file stats:", error);
      throw new Error("Failed to get file statistics.");
    }
  }
  clearCache() {
    this.fileTreeCache.clear();
  }
}
const filesApi = new FilesApi();
function findNode(tree, path, cache2) {
  if (cache2) {
    const cached = cache2.get(path);
    if (cached) return cached;
  }
  for (const node of tree) {
    if (node.path === path) {
      cache2?.set(path, node);
      return node;
    }
    if (node.children) {
      const found = findNode(node.children, path, cache2);
      if (found) return found;
    }
  }
  return null;
}
function walkTree(tree, fn) {
  for (const node of tree) {
    fn(node);
    if (node.children) {
      walkTree(node.children, fn);
    }
  }
}
function filterTreeByExtensions(tree, include, exclude = []) {
  return tree.filter((node) => {
    if (node.isDir) {
      const filteredChildren = node.children ? filterTreeByExtensions(node.children, include, exclude) : [];
      return filteredChildren.length > 0;
    }
    if (exclude.length > 0 && exclude.some((ext) => node.name.endsWith(ext))) {
      return false;
    }
    if (include.length === 0) return true;
    return include.some((ext) => node.name.endsWith(ext));
  }).map((node) => ({
    ...node,
    children: node.children ? filterTreeByExtensions(node.children, include, exclude) : void 0
  }));
}
function convertDomainNodes(domainNodes) {
  return domainNodes.map((node) => ({
    name: node.name,
    path: node.path,
    isDir: node.isDir,
    isExpanded: false,
    isSelected: false,
    children: node.children ? convertDomainNodes(node.children) : void 0,
    size: node.size,
    isIgnored: node.isIgnored || node.isGitignored || node.isCustomIgnored
  }));
}
function generateFlattened(tree, maxDepth, currentDepth = 0) {
  if (currentDepth >= maxDepth) return [];
  const result = [];
  for (const node of tree) {
    result.push(node);
    if (node.children && currentDepth < maxDepth - 1) {
      result.push(...generateFlattened(node.children, maxDepth, currentDepth + 1));
    }
  }
  return result;
}
function buildNodePathCache(tree, cache2) {
  for (const node of tree) {
    cache2.set(node.path, node);
    if (node.children) {
      buildNodePathCache(node.children, cache2);
    }
  }
}
function getAllFilesInNode(node, cache2) {
  if (cache2) {
    const cached = cache2.get(node.path);
    if (cached) return cached;
  }
  const files2 = [];
  if (node.children) {
    for (const child of node.children) {
      if (!child.isDir) {
        files2.push(child.path);
      } else {
        files2.push(...getAllFilesInNode(child, cache2));
      }
    }
  }
  cache2?.set(node.path, files2);
  return files2;
}
function autoExpandToFiles(tree, maxDepth = 3, currentDepth = 0) {
  if (currentDepth >= maxDepth) return;
  for (const node of tree) {
    if (node.isDir && node.children && node.children.length > 0) {
      const hasFiles = node.children.some((child) => !child.isDir);
      const hasOnlyFolders = node.children.every((child) => child.isDir);
      if (hasFiles || hasOnlyFolders && currentDepth < maxDepth - 1) {
        node.isExpanded = true;
        autoExpandToFiles(node.children, maxDepth, currentDepth + 1);
      } else if (currentDepth === 0) {
        node.isExpanded = true;
        autoExpandToFiles(node.children, maxDepth, currentDepth + 1);
      }
    }
  }
}
function sortFoldersFirst(nodes) {
  let needsSort = false;
  for (let i = 1; i < nodes.length; i++) {
    const prev = nodes[i - 1];
    const curr = nodes[i];
    if (!prev.isDir && curr.isDir) {
      needsSort = true;
      break;
    }
    if (prev.isDir === curr.isDir && prev.name.toLowerCase().localeCompare(curr.name.toLowerCase()) > 0) {
      needsSort = true;
      break;
    }
  }
  const sorted = needsSort ? [...nodes].sort((a, b) => {
    if (a.isDir && !b.isDir) return -1;
    if (!a.isDir && b.isDir) return 1;
    return a.name.toLowerCase().localeCompare(b.name.toLowerCase());
  }) : nodes;
  return sorted.map((node) => {
    if (node.isDir && node.children && node.children.length > 1) {
      const sortedChildren = sortFoldersFirst(node.children);
      if (sortedChildren !== node.children) {
        return { ...node, children: sortedChildren };
      }
    }
    return node;
  });
}
function useFileFilter(options2) {
  const { nodes } = options2;
  const settingsStore = useSettingsStore();
  const filterExtensions = ref([]);
  const excludeExtensions = ref([]);
  const filteredNodes = computed(() => {
    let result = nodes.value;
    if (filterExtensions.value.length > 0 || excludeExtensions.value.length > 0) {
      result = filterTreeByExtensions(
        result,
        filterExtensions.value,
        excludeExtensions.value
      );
    }
    if (settingsStore.settings.fileExplorer.foldersFirst) {
      result = sortFoldersFirst(result);
    }
    return result;
  });
  function setFilterExtensions(include, exclude = []) {
    filterExtensions.value = include;
    excludeExtensions.value = exclude;
  }
  function clearFilters() {
    filterExtensions.value = [];
    excludeExtensions.value = [];
  }
  function addIncludeExtension(ext) {
    if (!filterExtensions.value.includes(ext)) {
      filterExtensions.value.push(ext);
    }
  }
  function removeIncludeExtension(ext) {
    const index = filterExtensions.value.indexOf(ext);
    if (index > -1) {
      filterExtensions.value.splice(index, 1);
    }
  }
  function addExcludeExtension(ext) {
    if (!excludeExtensions.value.includes(ext)) {
      excludeExtensions.value.push(ext);
    }
  }
  function removeExcludeExtension(ext) {
    const index = excludeExtensions.value.indexOf(ext);
    if (index > -1) {
      excludeExtensions.value.splice(index, 1);
    }
  }
  return {
    // State
    filterExtensions,
    excludeExtensions,
    // Computed
    filteredNodes,
    // Actions
    setFilterExtensions,
    clearFilters,
    addIncludeExtension,
    removeIncludeExtension,
    addExcludeExtension,
    removeExcludeExtension
  };
}
function useFileFuzzySearch(options2) {
  const { flattenedNodes, maxResults = 100, fuseThreshold = 0.3 } = options2;
  const searchQuery = ref("");
  const searchResults = computed(() => {
    if (!searchQuery.value) return [];
    const allFiles = flattenedNodes.value;
    if (allFiles.length > 2e3) {
      const query = searchQuery.value.toLowerCase();
      return allFiles.filter(
        (file) => file.name.toLowerCase().includes(query) || file.path.toLowerCase().includes(query)
      ).slice(0, maxResults);
    }
    const fuse = new Fuse(allFiles, {
      keys: ["name", "path"],
      threshold: fuseThreshold
    });
    return fuse.search(searchQuery.value).map((result) => result.item).slice(0, maxResults);
  });
  function setSearchQuery(query) {
    searchQuery.value = query;
  }
  function clearSearch() {
    searchQuery.value = "";
  }
  return {
    // State
    searchQuery,
    // Computed
    searchResults,
    // Actions
    setSearchQuery,
    clearSearch
  };
}
const logger$d = useLogger("FilePersistence");
const SELECTION_PREFIX = "file-selection-";
const EXPANDED_PREFIX = "file-expanded-";
const MAX_SAVED_SELECTIONS = 100;
function useFilePersistence(options2) {
  const { nodes, selectedPaths, rootPath, findNode: findNode2 } = options2;
  let saveExpandedStateTimer = null;
  let saveSelectionTimer = null;
  function saveSelectionToStorage() {
    if (!rootPath.value) return;
    try {
      const key = `${SELECTION_PREFIX}${rootPath.value}`;
      const selection = Array.from(selectedPaths.value).slice(0, MAX_SAVED_SELECTIONS);
      localStorage.setItem(key, JSON.stringify(selection));
      logger$d.debug(`Saved selection: ${selection.length} files`);
    } catch (err) {
      console.warn("[FilePersistence] Failed to save selection:", err);
    }
  }
  function debouncedSaveSelection(delay = 300) {
    if (saveSelectionTimer) {
      clearTimeout(saveSelectionTimer);
    }
    saveSelectionTimer = setTimeout(() => {
      saveSelectionToStorage();
      saveSelectionTimer = null;
    }, delay);
  }
  function loadSelectionFromStorage(projectPath) {
    try {
      const key = `${SELECTION_PREFIX}${projectPath}`;
      const saved = localStorage.getItem(key);
      if (saved) {
        const selection = JSON.parse(saved);
        logger$d.debug(`Loaded selection: ${selection.length} files`);
        return selection;
      }
    } catch (err) {
      console.warn("[FilePersistence] Failed to load selection:", err);
    }
    return [];
  }
  function clearSelectionHistory(projectPath) {
    try {
      if (projectPath) {
        const key = `${SELECTION_PREFIX}${projectPath}`;
        localStorage.removeItem(key);
      } else {
        const keys = Object.keys(localStorage).filter(
          (k) => k.startsWith(SELECTION_PREFIX)
        );
        keys.forEach((k) => localStorage.removeItem(k));
      }
    } catch (err) {
      console.warn("[FilePersistence] Failed to clear selection history:", err);
    }
  }
  function getSelectionStats() {
    const stats = {};
    try {
      const keys = Object.keys(localStorage).filter(
        (k) => k.startsWith(SELECTION_PREFIX)
      );
      keys.forEach((key) => {
        const projectPath = key.replace(SELECTION_PREFIX, "");
        const saved = localStorage.getItem(key);
        if (saved) {
          const selection = JSON.parse(saved);
          stats[projectPath] = selection.length;
        }
      });
    } catch (err) {
      console.warn("[FilePersistence] Failed to get selection stats:", err);
    }
    return stats;
  }
  function saveExpandedState() {
    if (!rootPath.value) return;
    try {
      const expandedPaths = [];
      walkTree(nodes.value, (node) => {
        if (node.isDir && node.isExpanded) {
          expandedPaths.push(node.path);
        }
      });
      const key = `${EXPANDED_PREFIX}${rootPath.value}`;
      localStorage.setItem(key, JSON.stringify(expandedPaths));
    } catch (err) {
      console.warn("[FilePersistence] Failed to save expanded state:", err);
    }
  }
  function debouncedSaveExpandedState(delay = 500) {
    if (saveExpandedStateTimer) {
      clearTimeout(saveExpandedStateTimer);
    }
    saveExpandedStateTimer = setTimeout(() => {
      saveExpandedState();
      saveExpandedStateTimer = null;
    }, delay);
  }
  function loadExpandedState() {
    if (!rootPath.value) return [];
    try {
      const key = `${EXPANDED_PREFIX}${rootPath.value}`;
      const saved = localStorage.getItem(key);
      if (saved) {
        const expandedPaths = JSON.parse(saved);
        expandedPaths.forEach((path) => {
          const node = findNode2(path);
          if (node && node.isDir) {
            node.isExpanded = true;
          }
        });
        return expandedPaths;
      }
    } catch (err) {
      console.warn("[FilePersistence] Failed to load expanded state:", err);
    }
    return [];
  }
  function clearExpandedHistory(projectPath) {
    try {
      if (projectPath) {
        const key = `${EXPANDED_PREFIX}${projectPath}`;
        localStorage.removeItem(key);
      } else {
        const keys = Object.keys(localStorage).filter(
          (k) => k.startsWith(EXPANDED_PREFIX)
        );
        keys.forEach((k) => localStorage.removeItem(k));
      }
    } catch (err) {
      console.warn("[FilePersistence] Failed to clear expanded history:", err);
    }
  }
  function dispose() {
    if (saveExpandedStateTimer) {
      clearTimeout(saveExpandedStateTimer);
      saveExpandedStateTimer = null;
    }
    if (saveSelectionTimer) {
      clearTimeout(saveSelectionTimer);
      saveSelectionTimer = null;
    }
  }
  return {
    // Selection persistence
    saveSelectionToStorage,
    debouncedSaveSelection,
    loadSelectionFromStorage,
    clearSelectionHistory,
    getSelectionStats,
    // Expanded state persistence
    saveExpandedState,
    debouncedSaveExpandedState,
    loadExpandedState,
    clearExpandedHistory,
    // Cleanup
    dispose
  };
}
const MAX_HISTORY = 15;
function useFileSelection(options2) {
  const { findNode: findNode2, getAllFilesInNode: getAllFilesInNode2 } = options2;
  const selectedPaths = shallowRef(/* @__PURE__ */ new Set());
  const history = shallowRef([[]]);
  const historyIndex = shallowRef(0);
  const hasSelectedFiles = computed(() => selectedPaths.value.size > 0);
  const selectedCount = computed(() => selectedPaths.value.size);
  const selectedFilesList = computed(() => Array.from(selectedPaths.value));
  const canUndo = computed(() => historyIndex.value > 0);
  const canRedo = computed(() => historyIndex.value < history.value.length - 1);
  function saveToHistory() {
    const currentState = Array.from(selectedPaths.value);
    const newHistory = history.value.slice(0, historyIndex.value + 1);
    newHistory.push(currentState);
    if (newHistory.length > MAX_HISTORY) {
      newHistory.shift();
    } else {
      historyIndex.value++;
    }
    history.value = newHistory;
  }
  function undoSelection() {
    if (!canUndo.value) return false;
    historyIndex.value--;
    const prevState = history.value[historyIndex.value];
    selectedPaths.value = new Set(prevState);
    triggerRef(selectedPaths);
    return true;
  }
  function redoSelection() {
    if (!canRedo.value) return false;
    historyIndex.value++;
    const nextState = history.value[historyIndex.value];
    selectedPaths.value = new Set(nextState);
    triggerRef(selectedPaths);
    return true;
  }
  function toggleSelect(path) {
    const node = findNode2(path);
    if (!node) return;
    saveToHistory();
    if (node.isDir) {
      toggleSelectRecursiveInternal(path);
    } else {
      if (selectedPaths.value.has(path)) {
        selectedPaths.value.delete(path);
      } else {
        selectedPaths.value.add(path);
      }
    }
    triggerRef(selectedPaths);
  }
  function selectPath(path) {
    selectedPaths.value.add(path);
    triggerRef(selectedPaths);
  }
  function deselectPath(path) {
    selectedPaths.value.delete(path);
    triggerRef(selectedPaths);
  }
  function selectMultiple(paths) {
    paths.forEach((p) => selectedPaths.value.add(p));
    triggerRef(selectedPaths);
  }
  function clearSelection() {
    if (selectedPaths.value.size > 0) {
      saveToHistory();
    }
    selectedPaths.value.clear();
    triggerRef(selectedPaths);
  }
  function selectRecursive(path) {
    const node = findNode2(path);
    if (!node) return;
    const filePaths = getAllFilesInNode2(node);
    filePaths.forEach((p) => selectedPaths.value.add(p));
    triggerRef(selectedPaths);
  }
  function deselectRecursive(path) {
    const node = findNode2(path);
    if (!node) return;
    const filePaths = getAllFilesInNode2(node);
    filePaths.forEach((p) => selectedPaths.value.delete(p));
    triggerRef(selectedPaths);
  }
  function toggleSelectRecursiveInternal(path) {
    const node = findNode2(path);
    if (!node || !node.isDir) return;
    const childFilePaths = getAllFilesInNode2(node);
    const anySelected = childFilePaths.some(
      (filePath) => selectedPaths.value.has(filePath)
    );
    if (anySelected) {
      childFilePaths.forEach((filePath) => selectedPaths.value.delete(filePath));
    } else {
      childFilePaths.forEach((filePath) => selectedPaths.value.add(filePath));
    }
  }
  function toggleSelectRecursive(path) {
    saveToHistory();
    toggleSelectRecursiveInternal(path);
    triggerRef(selectedPaths);
  }
  function selectByExtension(extension, nodes, walkTree2) {
    walkTree2(nodes, (node) => {
      if (!node.isDir && node.name.endsWith(extension)) {
        selectedPaths.value.add(node.path);
      }
    });
    triggerRef(selectedPaths);
  }
  function getSelectedFileCountInNode(node) {
    const allFiles = getAllFilesInNode2(node);
    return allFiles.filter((filePath) => selectedPaths.value.has(filePath)).length;
  }
  function isSelected(path) {
    return selectedPaths.value.has(path);
  }
  function getSelectionState(node) {
    if (!node.isDir) {
      return selectedPaths.value.has(node.path) ? "full" : "none";
    }
    if (selectedCount.value === 0) {
      return "none";
    }
    const allFiles = getAllFilesInNode2(node);
    if (allFiles.length === 0) {
      return selectedPaths.value.has(node.path) ? "full" : "none";
    }
    let selectedFileCount = 0;
    for (const filePath of allFiles) {
      if (selectedPaths.value.has(filePath)) {
        selectedFileCount++;
      }
    }
    if (selectedFileCount === 0) return "none";
    if (selectedFileCount === allFiles.length) return "full";
    return "partial";
  }
  return {
    // State
    selectedPaths,
    // Computed
    hasSelectedFiles,
    selectedCount,
    selectedFilesList,
    canUndo,
    canRedo,
    // Actions
    toggleSelect,
    selectPath,
    deselectPath,
    selectMultiple,
    clearSelection,
    selectRecursive,
    deselectRecursive,
    toggleSelectRecursive,
    selectByExtension,
    getSelectedFileCountInNode,
    isSelected,
    getSelectionState,
    undoSelection,
    redoSelection
  };
}
function useFileTree() {
  const nodes = ref([]);
  const rootPath = ref("");
  const currentDirectory = ref("");
  const directoryHistory = ref([]);
  const allFilesCache = /* @__PURE__ */ new Map();
  const nodePathCache = /* @__PURE__ */ new Map();
  const flattenedNodesCache = shallowRef(null);
  const flattenedNodes = computed(() => {
    if (nodes.value.length === 0) return [];
    if (!flattenedNodesCache.value) {
      flattenedNodesCache.value = generateFlattened(nodes.value, FILE_TREE.MAX_FLATTEN_DEPTH);
    }
    return flattenedNodesCache.value;
  });
  const projectName = computed(() => {
    if (!rootPath.value) return "Project";
    return rootPath.value.split(/[\\/]/).pop() || "Project";
  });
  const breadcrumbs = computed(() => {
    if (!currentDirectory.value || !rootPath.value) return [];
    const relative = currentDirectory.value.replace(rootPath.value, "").replace(/^[\\/]+/, "");
    if (!relative) return [projectName.value];
    const segments = relative.split(/[\\/]/);
    return [projectName.value, ...segments];
  });
  function setFileTree(tree) {
    const expandedPaths = getExpandedPaths();
    nodes.value = convertDomainNodes(tree);
    allFilesCache.clear();
    nodePathCache.clear();
    flattenedNodesCache.value = null;
    buildNodePathCache(nodes.value, nodePathCache);
    if (expandedPaths.length > 0) {
      restoreExpandedPaths(expandedPaths);
    }
  }
  function removeNode(path, selectedPaths) {
    const removeFromTree = (tree, targetPath) => {
      for (let i = 0; i < tree.length; i++) {
        if (tree[i].path === targetPath) {
          if (tree[i].isDir) {
            const filesToDeselect = getAllFilesInNode$1(tree[i]);
            filesToDeselect.forEach((p) => selectedPaths.delete(p));
          } else {
            selectedPaths.delete(targetPath);
          }
          tree.splice(i, 1);
          return true;
        }
        if (tree[i].children) {
          if (removeFromTree(tree[i].children, targetPath)) {
            return true;
          }
        }
      }
      return false;
    };
    const removed = removeFromTree(nodes.value, path);
    if (removed) {
      nodePathCache.delete(path);
      allFilesCache.delete(path);
      const parentPath = path.substring(0, path.lastIndexOf("/")) || path.substring(0, path.lastIndexOf("\\"));
      if (parentPath) {
        allFilesCache.delete(parentPath);
      }
      flattenedNodesCache.value = null;
      nodes.value = [...nodes.value];
    }
    return removed;
  }
  function findNode$1(path) {
    return findNode(nodes.value, path, nodePathCache);
  }
  function nodeExists(path) {
    return findNode$1(path) !== null;
  }
  function toggleExpand(path) {
    const node = findNode$1(path);
    if (node && node.isDir) {
      node.isExpanded = !node.isExpanded;
    }
  }
  function expandPath(path) {
    const node = findNode$1(path);
    if (node && node.isDir) {
      node.isExpanded = true;
    }
  }
  function collapsePath(path) {
    const node = findNode$1(path);
    if (node && node.isDir) {
      node.isExpanded = false;
    }
  }
  function expandRecursive(path) {
    const node = findNode$1(path);
    if (!node || !node.isDir) return;
    const expandNode = (n) => {
      if (n.isDir) {
        n.isExpanded = true;
        if (n.children) {
          n.children.forEach(expandNode);
        }
      }
    };
    expandNode(node);
  }
  function collapseRecursive(path) {
    const node = findNode$1(path);
    if (!node || !node.isDir) return;
    const collapseNode = (n) => {
      if (n.isDir) {
        n.isExpanded = false;
        if (n.children) {
          n.children.forEach(collapseNode);
        }
      }
    };
    collapseNode(node);
  }
  function expandAll() {
    walkTree(nodes.value, (node) => {
      if (node.isDir) {
        node.isExpanded = true;
      }
    });
  }
  function collapseAll() {
    walkTree(nodes.value, (node) => {
      if (node.isDir) {
        node.isExpanded = false;
      }
    });
  }
  function getExpandedPaths() {
    const expanded = [];
    walkTree(nodes.value, (node) => {
      if (node.isDir && node.isExpanded) {
        expanded.push(node.path);
      }
    });
    return expanded;
  }
  function restoreExpandedPaths(paths) {
    const pathSet = new Set(paths);
    walkTree(nodes.value, (node) => {
      if (node.isDir) {
        node.isExpanded = pathSet.has(node.path);
      }
    });
  }
  function getAllFilesInNode$1(node) {
    return getAllFilesInNode(node, allFilesCache);
  }
  function getRecursiveFileCount(node) {
    if (!node.isDir) return 0;
    return getAllFilesInNode$1(node).length;
  }
  function isDirectory(path) {
    const node = findNode$1(path);
    return node ? node.isDir : false;
  }
  function getNodesByPaths(paths) {
    const result = /* @__PURE__ */ new Map();
    const pathSet = new Set(paths);
    function walk(nodeList) {
      for (const node of nodeList) {
        if (pathSet.has(node.path)) {
          result.set(node.path, node);
          if (result.size === paths.length) return;
        }
        if (node.children) {
          walk(node.children);
        }
      }
    }
    walk(nodes.value);
    return result;
  }
  function getAvailableExtensions() {
    const extensions = /* @__PURE__ */ new Set();
    function collectExtensions(nodeList) {
      for (const node of nodeList) {
        if (!node.isDir && node.name.includes(".")) {
          const ext = "." + node.name.split(".").pop();
          extensions.add(ext);
        }
        if (node.children) {
          collectExtensions(node.children);
        }
      }
    }
    collectExtensions(nodes.value);
    return Array.from(extensions).sort();
  }
  function setRootPath(path) {
    rootPath.value = path;
    currentDirectory.value = path;
    directoryHistory.value = [];
  }
  function autoExpand(maxDepth = 3) {
    autoExpandToFiles(nodes.value, maxDepth);
  }
  function getMemoryUsage() {
    let size = 0;
    let nodeCount = 0;
    const countNodes = (nodeList) => {
      nodeCount += nodeList.length;
      nodeList.forEach((node) => {
        if (node.children) countNodes(node.children);
      });
    };
    countNodes(nodes.value);
    size += nodeCount * 250;
    size += (flattenedNodesCache.value?.length || 0) * 200;
    size += allFilesCache.size * 150;
    size += nodePathCache.size * 50;
    return size;
  }
  function clearCaches() {
    allFilesCache.clear();
    nodePathCache.clear();
    flattenedNodesCache.value = null;
  }
  function reset() {
    nodes.value = [];
    rootPath.value = "";
    currentDirectory.value = "";
    directoryHistory.value = [];
    clearCaches();
  }
  return {
    // State
    nodes,
    rootPath,
    currentDirectory,
    directoryHistory,
    // Computed
    flattenedNodes,
    projectName,
    breadcrumbs,
    // Actions
    setFileTree,
    removeNode,
    findNode: findNode$1,
    nodeExists,
    toggleExpand,
    expandPath,
    collapsePath,
    expandRecursive,
    collapseRecursive,
    expandAll,
    collapseAll,
    getExpandedPaths,
    restoreExpandedPaths,
    getAllFilesInNode: getAllFilesInNode$1,
    getRecursiveFileCount,
    isDirectory,
    getNodesByPaths,
    getAvailableExtensions,
    setRootPath,
    autoExpand,
    getMemoryUsage,
    clearCaches,
    reset
  };
}
const useFileStore = defineStore("file", () => {
  const tree = useFileTree();
  const selection = useFileSelection({
    findNode: tree.findNode,
    getAllFilesInNode: tree.getAllFilesInNode
  });
  const search = useFileFuzzySearch({
    flattenedNodes: tree.flattenedNodes
  });
  const filter = useFileFilter({
    nodes: tree.nodes
  });
  const persistence = useFilePersistence({
    nodes: tree.nodes,
    selectedPaths: selection.selectedPaths,
    rootPath: tree.rootPath,
    findNode: tree.findNode
  });
  const isLoading = ref(false);
  const error = ref(null);
  const settingsStore = useSettingsStore();
  const autoSaveSelection = computed(() => settingsStore.settings.fileExplorer.autoSaveSelection);
  const selectedFilesTotalSize = computed(() => {
    let totalSize = 0;
    selection.selectedPaths.value.forEach((path) => {
      const node = tree.findNode(path);
      if (node && !node.isDir && node.size) {
        totalSize += node.size;
      }
    });
    return totalSize;
  });
  const estimatedTokenCount = computed(() => Math.round(selectedFilesTotalSize.value / 4));
  const estimatedContextSize = computed(() => selectedFilesTotalSize.value / (1024 * 1024));
  function getSelectedFilesSize() {
    return selectedFilesTotalSize.value;
  }
  async function loadFileTree(projectPath, directory) {
    isLoading.value = true;
    error.value = null;
    try {
      const targetPath = directory || projectPath;
      const files2 = await filesApi.listFiles(targetPath, true, true);
      tree.setFileTree(files2);
      if (!tree.rootPath.value) {
        tree.setRootPath(projectPath);
        const loadedPaths = persistence.loadExpandedState();
        if (loadedPaths.length === 0) {
          tree.autoExpand(3);
        }
        const savedSelection = persistence.loadSelectionFromStorage(projectPath);
        if (savedSelection.length > 0) {
          selection.selectMultiple(savedSelection);
        }
      } else if (directory) {
        tree.currentDirectory.value = directory;
      }
    } catch (err) {
      error.value = err instanceof Error ? err.message : "Failed to load files";
      throw err;
    } finally {
      isLoading.value = false;
    }
  }
  function toggleSelect(path) {
    selection.toggleSelect(path);
    if (autoSaveSelection.value) {
      persistence.debouncedSaveSelection();
    }
  }
  function clearSelection() {
    selection.clearSelection();
    if (autoSaveSelection.value) {
      persistence.debouncedSaveSelection();
    }
  }
  function toggleExpand(pathOrCompact) {
    if (pathOrCompact.startsWith("{")) {
      try {
        const { paths, expand } = JSON.parse(pathOrCompact);
        for (const p of paths) {
          if (expand) {
            tree.expandPath(p);
          } else {
            tree.collapsePath(p);
          }
        }
      } catch {
        tree.toggleExpand(pathOrCompact);
      }
    } else {
      tree.toggleExpand(pathOrCompact);
    }
    persistence.debouncedSaveExpandedState();
  }
  function expandRecursive(path) {
    tree.expandRecursive(path);
    persistence.debouncedSaveExpandedState();
  }
  function collapseRecursive(path) {
    tree.collapseRecursive(path);
    persistence.debouncedSaveExpandedState();
  }
  function expandAll() {
    tree.expandAll();
    persistence.debouncedSaveExpandedState();
  }
  function collapseAll() {
    tree.collapseAll();
    persistence.debouncedSaveExpandedState();
  }
  function removeNode(path) {
    const removed = tree.removeNode(path, selection.selectedPaths.value);
    if (removed) {
      triggerRef(selection.selectedPaths);
    }
    return removed;
  }
  function selectByExtension(extension) {
    walkTree(tree.nodes.value, (node) => {
      if (!node.isDir && node.name.endsWith(extension)) {
        selection.selectedPaths.value.add(node.path);
      }
    });
    triggerRef(selection.selectedPaths);
  }
  async function refreshFileTree() {
    filesApi.clearCache();
  }
  function resetStore() {
    tree.reset();
    selection.clearSelection();
    search.clearSearch();
    filter.clearFilters();
    error.value = null;
    isLoading.value = false;
    persistence.dispose();
    if (typeof window !== "undefined" && "gc" in window) {
      try {
        window.gc?.();
      } catch {
      }
    }
  }
  function getMemoryUsage() {
    let size = tree.getMemoryUsage();
    size += selection.selectedPaths.value.size * 100;
    return size;
  }
  function pruneUnusedBranches() {
  }
  return {
    // State (from tree)
    nodes: tree.nodes,
    rootPath: tree.rootPath,
    currentDirectory: tree.currentDirectory,
    directoryHistory: tree.directoryHistory,
    isLoading,
    error,
    // State (from selection)
    selectedPaths: selection.selectedPaths,
    // State (from search)
    searchQuery: search.searchQuery,
    // State (from filter)
    filterExtensions: filter.filterExtensions,
    excludeExtensions: filter.excludeExtensions,
    // Computed (from tree)
    projectName: tree.projectName,
    breadcrumbs: tree.breadcrumbs,
    flattenedNodes: tree.flattenedNodes,
    // Computed (from selection)
    hasSelectedFiles: selection.hasSelectedFiles,
    selectedCount: selection.selectedCount,
    selectedFilesList: selection.selectedFilesList,
    canUndoSelection: selection.canUndo,
    canRedoSelection: selection.canRedo,
    // Computed (from search)
    searchResults: search.searchResults,
    // Computed (from filter)
    filteredNodes: filter.filteredNodes,
    // Computed (local)
    estimatedTokenCount,
    estimatedContextSize,
    // Actions (tree)
    setFileTree: tree.setFileTree,
    loadFileTree,
    removeNode,
    toggleExpand,
    expandPath: tree.expandPath,
    collapsePath: tree.collapsePath,
    expandRecursive,
    collapseRecursive,
    expandAll,
    collapseAll,
    setRootPath: tree.setRootPath,
    getAvailableExtensions: tree.getAvailableExtensions,
    nodeExists: tree.nodeExists,
    getExpandedPaths: tree.getExpandedPaths,
    restoreExpandedPaths: tree.restoreExpandedPaths,
    autoExpandToFiles: () => tree.autoExpand(3),
    // Actions (selection)
    toggleSelect,
    selectPath: selection.selectPath,
    deselectPath: selection.deselectPath,
    selectMultiple: selection.selectMultiple,
    clearSelection,
    selectRecursive: selection.selectRecursive,
    deselectRecursive: selection.deselectRecursive,
    selectByExtension,
    undoSelection: selection.undoSelection,
    redoSelection: selection.redoSelection,
    // Actions (search)
    setSearchQuery: search.setSearchQuery,
    // Actions (filter)
    setFilterExtensions: filter.setFilterExtensions,
    // Actions (persistence)
    autoSaveSelection,
    saveSelectionToStorage: persistence.saveSelectionToStorage,
    loadSelectionFromStorage: persistence.loadSelectionFromStorage,
    clearSelectionHistory: persistence.clearSelectionHistory,
    getSelectionStats: persistence.getSelectionStats,
    saveExpandedState: persistence.saveExpandedState,
    loadExpandedState: persistence.loadExpandedState,
    // Actions (other)
    refreshFileTree,
    getSelectedFilesSize,
    resetStore,
    getMemoryUsage,
    pruneUnusedBranches,
    // Public utility methods for UI components
    getRecursiveFileCount: tree.getRecursiveFileCount,
    getAllFilesInNode: tree.getAllFilesInNode,
    getSelectedFileCountInNode: selection.getSelectedFileCountInNode,
    isDirectory: tree.isDirectory,
    getNodesByPaths: tree.getNodesByPaths
  };
});
const file_store = /* @__PURE__ */ Object.freeze(/* @__PURE__ */ Object.defineProperty({
  __proto__: null,
  useFileStore
}, Symbol.toStringTag, { value: "Module" }));
const _hoisted_1$19 = { class: "flex items-center gap-1 text-sm text-gray-400" };
const _hoisted_2$15 = ["title", "onClick"];
const _hoisted_3$10 = {
  key: 0,
  class: "w-3.5 h-3.5",
  fill: "currentColor",
  viewBox: "0 0 20 20"
};
const _hoisted_4$W = { class: "truncate max-w-[120px]" };
const _hoisted_5$P = {
  key: 0,
  class: "w-3.5 h-3.5 flex-shrink-0 text-gray-400",
  fill: "currentColor",
  viewBox: "0 0 20 20"
};
const _hoisted_6$L = ["title", "aria-label"];
const _hoisted_7$H = ["title", "aria-label"];
const _sfc_main$1b = /* @__PURE__ */ defineComponent({
  __name: "BreadcrumbsNav",
  props: {
    segments: {},
    rootName: { default: "Project" }
  },
  emits: ["navigate", "open-in-explorer"],
  setup(__props, { emit: __emit }) {
    const { t } = useI18n();
    const logger2 = useLogger("BreadcrumbsNav");
    const props = __props;
    const emit = __emit;
    function getFullPath(index) {
      return props.segments.slice(0, index + 1).join("/");
    }
    function handleClick(index) {
      const path = getFullPath(index);
      emit("navigate", path);
    }
    async function copyPath() {
      const fullPath = props.segments.join("/");
      try {
        await navigator.clipboard.writeText(fullPath);
      } catch (err) {
        logger2.warn("Failed to copy path:", err);
      }
    }
    async function openInExplorer() {
      const fullPath = props.segments.join("/");
      emit("open-in-explorer", fullPath);
    }
    return (_ctx, _cache) => {
      return openBlock(), createElementBlock("div", _hoisted_1$19, [
        (openBlock(true), createElementBlock(Fragment, null, renderList(__props.segments, (segment, index) => {
          return openBlock(), createElementBlock(Fragment, { key: index }, [
            createBaseVNode("button", {
              class: normalizeClass([
                "breadcrumb-segment",
                index === __props.segments.length - 1 ? "breadcrumb-segment-active" : "breadcrumb-segment-link"
              ]),
              title: getFullPath(index),
              onClick: ($event) => handleClick(index)
            }, [
              index === 0 ? (openBlock(), createElementBlock("svg", _hoisted_3$10, [..._cache[0] || (_cache[0] = [
                createBaseVNode("path", { d: "M10.707 2.293a1 1 0 00-1.414 0l-7 7a1 1 0 001.414 1.414L4 10.414V17a1 1 0 001 1h2a1 1 0 001-1v-2a1 1 0 011-1h2a1 1 0 011 1v2a1 1 0 001 1h2a1 1 0 001-1v-6.586l.293.293a1 1 0 001.414-1.414l-7-7z" }, null, -1)
              ])])) : createCommentVNode("", true),
              createBaseVNode("span", _hoisted_4$W, toDisplayString(segment), 1)
            ], 10, _hoisted_2$15),
            index < __props.segments.length - 1 ? (openBlock(), createElementBlock("svg", _hoisted_5$P, [..._cache[1] || (_cache[1] = [
              createBaseVNode("path", {
                "fill-rule": "evenodd",
                d: "M7.293 14.707a1 1 0 010-1.414L10.586 10 7.293 6.707a1 1 0 011.414-1.414l4 4a1 1 0 010 1.414l-4 4a1 1 0 01-1.414 0z",
                "clip-rule": "evenodd"
              }, null, -1)
            ])])) : createCommentVNode("", true)
          ], 64);
        }), 128)),
        __props.segments.length > 0 ? (openBlock(), createElementBlock("button", {
          key: 0,
          onClick: copyPath,
          class: "breadcrumb-action",
          title: unref(t)("files.copyPath"),
          "aria-label": unref(t)("files.copyPath")
        }, [..._cache[2] || (_cache[2] = [
          createBaseVNode("svg", {
            class: "w-3.5 h-3.5",
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
          ], -1)
        ])], 8, _hoisted_6$L)) : createCommentVNode("", true),
        __props.segments.length > 0 ? (openBlock(), createElementBlock("button", {
          key: 1,
          onClick: openInExplorer,
          class: "breadcrumb-action",
          title: unref(t)("files.openInExplorer"),
          "aria-label": unref(t)("files.openInExplorer")
        }, [..._cache[3] || (_cache[3] = [
          createBaseVNode("svg", {
            class: "w-3.5 h-3.5",
            fill: "none",
            stroke: "currentColor",
            viewBox: "0 0 24 24"
          }, [
            createBaseVNode("path", {
              "stroke-linecap": "round",
              "stroke-linejoin": "round",
              "stroke-width": "2",
              d: "M10 6H6a2 2 0 00-2 2v10a2 2 0 002 2h10a2 2 0 002-2v-4M14 4h6m0 0v6m0-6L10 14"
            })
          ], -1)
        ])], 8, _hoisted_7$H)) : createCommentVNode("", true)
      ]);
    };
  }
});
function useContextMenu() {
  const isVisible = ref(false);
  const position = ref({ x: 0, y: 0 });
  const targetNode = ref(null);
  function show(node, event) {
    targetNode.value = node;
    const viewportWidth = window.innerWidth;
    const viewportHeight = window.innerHeight;
    const menuWidth = 250;
    const menuHeight = 300;
    let x = event.clientX;
    let y = event.clientY;
    if (x + menuWidth > viewportWidth) {
      x = viewportWidth - menuWidth - 10;
    }
    if (y + menuHeight > viewportHeight) {
      y = viewportHeight - menuHeight - 10;
    }
    position.value = { x, y };
    isVisible.value = true;
  }
  function hide() {
    isVisible.value = false;
    targetNode.value = null;
  }
  function handleClickOutside(event) {
    const target2 = event.target;
    if (!target2.closest(".context-menu")) {
      hide();
    }
  }
  function handleEscape(event) {
    if (event.key === "Escape") {
      hide();
    }
  }
  onMounted(() => {
    document.addEventListener("click", handleClickOutside);
    document.addEventListener("keydown", handleEscape);
  });
  onUnmounted(() => {
    document.removeEventListener("click", handleClickOutside);
    document.removeEventListener("keydown", handleEscape);
  });
  return {
    isVisible,
    position,
    targetNode,
    show,
    hide
  };
}
const logger$c = useLogger("UIStore");
const useUIStore = defineStore("ui", () => {
  const toasts = ref([]);
  const nextToastId = ref(0);
  const showSettingsModal = ref(false);
  const showKeyboardShortcutsModal = ref(false);
  function addToast(message, type = "info", duration = 3e3, action) {
    const toast = {
      id: `toast-${nextToastId.value++}`,
      message,
      type,
      duration,
      action
    };
    const logMessage = `[Toast ${type.toUpperCase()}] ${message}`;
    switch (type) {
      case "error":
        logger$c.error(logMessage);
        break;
      case "warning":
        logger$c.warn(logMessage);
        break;
    }
    toasts.value.push(toast);
    if (duration > 0) {
      setTimeout(() => {
        removeToast(toast.id);
      }, duration);
    }
    return toast.id;
  }
  function removeToast(id) {
    const index = toasts.value.findIndex((t) => t.id === id);
    if (index !== -1) {
      toasts.value.splice(index, 1);
    }
  }
  function clearToasts() {
    toasts.value = [];
  }
  function openSettingsModal() {
    showSettingsModal.value = true;
  }
  function openKeyboardShortcutsModal() {
    showKeyboardShortcutsModal.value = true;
  }
  return {
    // State
    toasts,
    showSettingsModal,
    showKeyboardShortcutsModal,
    // Actions
    addToast,
    removeToast,
    clearToasts,
    openSettingsModal,
    openKeyboardShortcutsModal
  };
});
const logger$b = useLogger("FileUtils");
async function copyToClipboard(text) {
  try {
    if (navigator.clipboard && window.isSecureContext) {
      await navigator.clipboard.writeText(text);
    } else {
      const textArea = document.createElement("textarea");
      textArea.value = text;
      textArea.style.position = "fixed";
      textArea.style.left = "-999999px";
      document.body.appendChild(textArea);
      textArea.select();
      document.execCommand("copy");
      document.body.removeChild(textArea);
    }
  } catch (err) {
    logger$b.error("Failed to copy to clipboard:", err);
    throw err;
  }
}
function getRelativePath(fullPath, basePath) {
  if (!fullPath.startsWith(basePath)) return fullPath;
  const relative = fullPath.slice(basePath.length);
  return relative.replace(/^[\\/]+/, "");
}
function parseIgnoreRules(rules) {
  return rules.split("\n").map((line) => line.trim()).filter((line) => line && !line.startsWith("#"));
}
function debounce(func, wait) {
  let timeout = null;
  return (...args) => {
    if (timeout) {
      clearTimeout(timeout);
    }
    timeout = setTimeout(() => func(...args), wait);
  };
}
function useFileSearch() {
  const fileStore = useFileStore();
  const query = ref("");
  const isSearching = ref(false);
  const debouncedSearch = debounce(() => {
    fileStore.setSearchQuery(query.value);
    isSearching.value = false;
  }, FILE_TREE.DEBOUNCE_MS);
  function handleSearch() {
    isSearching.value = true;
    debouncedSearch();
  }
  function clear() {
    query.value = "";
    fileStore.setSearchQuery("");
    isSearching.value = false;
  }
  function setQuery(newQuery) {
    query.value = newQuery;
    handleSearch();
  }
  function isActive() {
    return query.value.length > 0;
  }
  return {
    query,
    isSearching,
    handleSearch,
    clear,
    setQuery,
    isActive
  };
}
const logger$a = useLogger("IgnoreRules");
function useIgnoreRules() {
  const fileStore = useFileStore();
  const settingsStore = useSettingsStore();
  const uiStore = useUIStore();
  async function addToIgnore(node) {
    try {
      const currentRules = settingsStore.getCustomIgnoreRules();
      const newRule = node.isDir ? `${node.name}/` : node.name;
      const updatedRules = currentRules ? `${currentRules}
${newRule}` : newRule;
      await apiService.updateCustomIgnoreRules(updatedRules);
      settingsStore.setCustomIgnoreRules(updatedRules);
      fileStore.removeNode(node.path);
      await apiService.clearFileTreeCache();
      filesApi.clearCache();
      uiStore.addToast("Добавлено в исключения", "success");
      return true;
    } catch (error) {
      logger$a.error("Failed to add to ignore:", error);
      uiStore.addToast("Ошибка добавления в исключения", "error");
      return false;
    }
  }
  async function removeFromIgnore(node) {
    try {
      const currentIgnoreRules = settingsStore.getCustomIgnoreRules();
      const pattern = node.isDir ? `${node.name}/` : node.name;
      const lines = currentIgnoreRules.split("\n").filter((line) => {
        const trimmed = line.trim();
        return trimmed && !trimmed.includes(pattern) && !trimmed.startsWith("#");
      });
      const updatedRules = lines.join("\n");
      await apiService.updateCustomIgnoreRules(updatedRules);
      settingsStore.setCustomIgnoreRules(updatedRules);
      uiStore.addToast("Удалено из исключений", "success");
      return true;
    } catch (error) {
      logger$a.error("Failed to remove from ignore:", error);
      uiStore.addToast("Ошибка удаления из исключений", "error");
      return false;
    }
  }
  function getCustomRules() {
    return settingsStore.getCustomIgnoreRules();
  }
  async function updateCustomRules(rules) {
    try {
      await apiService.updateCustomIgnoreRules(rules);
      settingsStore.setCustomIgnoreRules(rules);
      return true;
    } catch (error) {
      logger$a.error("Failed to update ignore rules:", error);
      return false;
    }
  }
  return {
    addToIgnore,
    removeFromIgnore,
    getCustomRules,
    updateCustomRules
  };
}
const hoveredState = reactive({
  path: null,
  isDir: null
});
function useHoveredFile() {
  const setHovered = (newPath, newIsDir) => {
    hoveredState.path = newPath;
    hoveredState.isDir = newIsDir;
  };
  const clearHovered = (currentPath) => {
    if (hoveredState.path === currentPath) {
      hoveredState.path = null;
      hoveredState.isDir = null;
    }
  };
  return {
    state: hoveredState,
    setHovered,
    clearHovered
  };
}
function useQuickLook() {
  const { t } = useI18n();
  const fileStore = useFileStore();
  const uiStore = useUIStore();
  const { state: hoveredState2 } = useHoveredFile();
  const isVisible = ref(false);
  const currentPath = ref("");
  function open(path) {
    if (isVisible.value && currentPath.value === path) {
      close();
      return;
    }
    currentPath.value = path;
    isVisible.value = true;
  }
  function close() {
    isVisible.value = false;
  }
  function toggle(path) {
    if (isVisible.value && currentPath.value === path) {
      close();
    } else {
      open(path);
    }
  }
  function addToContext(path) {
    fileStore.toggleSelect(path);
    uiStore.addToast(t("files.addToContext"), "success");
  }
  function handleSpacebarPreview() {
    const path = hoveredState2.path;
    const isDir = hoveredState2.isDir;
    if (path) {
      if (isDir) {
        fileStore.toggleExpand(path);
      } else {
        toggle(path);
      }
      return true;
    } else if (isVisible.value) {
      close();
      return true;
    }
    return false;
  }
  return {
    isVisible,
    currentPath,
    open,
    close,
    toggle,
    addToContext,
    handleSpacebarPreview
  };
}
const logger$9 = useLogger("FileExplorer");
function useFileExplorer() {
  const { t } = useI18n();
  const fileStore = useFileStore();
  const contextStore = useContextStore();
  const projectStore = useProjectStore();
  const uiStore = useUIStore();
  const settingsStore = useSettingsStore();
  const contextMenu = useContextMenu();
  const search = useFileSearch();
  const ignoreRules = useIgnoreRules();
  const quickLook = useQuickLook();
  const showSettings = ref(false);
  const filterExtensions = ref([]);
  const totalFileCount = ref(0);
  const availableExtensions = computed(() => fileStore.getAvailableExtensions());
  const hasSelectionHistory = computed(() => {
    if (!projectStore.currentPath) return false;
    const stats = fileStore.getSelectionStats();
    return stats[projectStore.currentPath] > 0;
  });
  const selectionHistoryCount = computed(() => {
    if (!projectStore.currentPath) return 0;
    const stats = fileStore.getSelectionStats();
    return stats[projectStore.currentPath] || 0;
  });
  const selectionProgress = computed(() => {
    if (totalFileCount.value === 0) return 0;
    return Math.round(fileStore.selectedCount / totalFileCount.value * 100);
  });
  watch(() => fileStore.nodes, (newNodes) => {
    let count = 0;
    const countFiles = (nodes) => {
      nodes.forEach((node) => {
        if (!node.isDir) count++;
        if (node.children) countFiles(node.children);
      });
    };
    countFiles(newNodes);
    totalFileCount.value = count;
  }, { immediate: true, deep: false });
  function handleToggleSelect(path) {
    fileStore.toggleSelect(path);
  }
  function handleToggleExpand(path) {
    fileStore.toggleExpand(path);
  }
  async function handleRefresh() {
    if (!projectStore.currentPath) return;
    try {
      fileStore.clearSelection();
      await fileStore.loadFileTree(projectStore.currentPath);
      uiStore.addToast(t("toast.refreshed"), "success");
    } catch (error) {
      logger$9.error("Failed to refresh file tree:", error);
      uiStore.addToast(t("toast.refreshError"), "error");
    }
  }
  async function handleRefreshPreserveState() {
    if (!projectStore.currentPath) return;
    try {
      const expandedPaths = fileStore.getExpandedPaths();
      const selectedPaths = fileStore.selectedFilesList;
      await apiService.clearFileTreeCache();
      filesApi.clearCache();
      await fileStore.loadFileTree(projectStore.currentPath);
      fileStore.restoreExpandedPaths(expandedPaths);
      for (const path of selectedPaths) {
        if (fileStore.nodeExists(path)) {
          fileStore.toggleSelect(path);
        }
      }
    } catch (error) {
      logger$9.error("Failed to refresh file tree:", error);
      uiStore.addToast(t("toast.refreshError"), "error");
    }
  }
  function handleFilterUpdate(selected) {
    fileStore.setFilterExtensions(selected);
  }
  async function handleSettingsChange() {
    try {
      const dto = await apiService.getSettings();
      dto.useGitignore = settingsStore.settings.fileExplorer.useGitignore;
      dto.useCustomIgnore = settingsStore.settings.fileExplorer.useCustomIgnore;
      await apiService.saveSettings(JSON.stringify(dto));
      if (projectStore.currentPath) {
        await handleRefresh();
      }
    } catch (error) {
      logger$9.error("Failed to save settings:", error);
      uiStore.addToast("Failed to save settings", "error");
    }
  }
  function restorePreviousSelection() {
    if (projectStore.currentPath) {
      fileStore.loadSelectionFromStorage(projectStore.currentPath);
      uiStore.addToast("Selection restored", "success");
    }
  }
  function handleContextMenuShow(node, event) {
    contextMenu.show(node, event);
  }
  function handleGlobalKeydown(event) {
    if (event.key !== " ") return;
    const target2 = event.target;
    if (target2.tagName === "INPUT" || target2.tagName === "TEXTAREA" || target2.isContentEditable) {
      return;
    }
    if (quickLook.handleSpacebarPreview()) {
      event.preventDefault();
    }
  }
  async function handleContextMenuAction(payload) {
    const { type, node } = payload;
    try {
      switch (type) {
        case "quickLook":
          if (!node.isDir) quickLook.open(node.path);
          break;
        case "selectAll":
          if (node.isDir) {
            fileStore.selectRecursive(node.path);
            uiStore.addToast("All files selected in folder", "success");
          }
          break;
        case "deselectAll":
          if (node.isDir) {
            fileStore.deselectRecursive(node.path);
            uiStore.addToast("All files deselected in folder", "success");
          }
          break;
        case "copyPath":
          await copyToClipboard(node.path);
          uiStore.addToast("Path copied to clipboard", "success");
          break;
        case "copyRelativePath": {
          const relativePath = projectStore.currentPath ? getRelativePath(node.path, projectStore.currentPath) : node.path;
          await copyToClipboard(relativePath);
          uiStore.addToast("Relative path copied to clipboard", "success");
          break;
        }
        case "addToCustomIgnore":
          await ignoreRules.addToIgnore(node);
          break;
        case "removeFromIgnore":
          if (await ignoreRules.removeFromIgnore(node)) {
            await handleRefreshPreserveState();
          }
          break;
        case "expandAll":
          if (node.isDir) {
            fileStore.expandRecursive(node.path);
            uiStore.addToast("Expanded all folders", "success");
          }
          break;
        case "collapseAll":
          if (node.isDir) {
            fileStore.collapseRecursive(node.path);
            uiStore.addToast("Collapsed all folders", "success");
          }
          break;
      }
    } catch (error) {
      logger$9.error("Context menu action failed:", error);
      uiStore.addToast("Action failed", "error");
    }
  }
  function initialize() {
    window.addEventListener("keydown", handleGlobalKeydown);
    if (projectStore.currentPath) {
      fileStore.loadFileTree(projectStore.currentPath).catch((error) => {
        logger$9.error("Failed to load file tree:", error);
        uiStore.addToast("Failed to load project files.", "error");
      });
    } else {
      uiStore.addToast("No project selected", "warning");
    }
  }
  function cleanup() {
    window.removeEventListener("keydown", handleGlobalKeydown);
  }
  function setupWatchers() {
    watch(() => projectStore.currentPath, async (newPath, oldPath) => {
      if (newPath && newPath !== oldPath) {
        fileStore.clearSelection();
        contextStore.clearContext();
        try {
          await fileStore.loadFileTree(newPath);
        } catch (error) {
          logger$9.error("Failed to load file tree:", error);
          uiStore.addToast("Failed to load project files", "error");
        }
      }
    });
  }
  return {
    // Search (delegated to useFileSearch)
    searchQuery: search.query,
    handleSearch: search.handleSearch,
    clearSearch: search.clear,
    // QuickLook
    quickLookVisible: quickLook.isVisible,
    quickLookPath: quickLook.currentPath,
    handleQuickLook: (path) => quickLook.toggle(path),
    handleAddToContext: (path) => quickLook.addToContext(path),
    // UI state
    showSettings,
    filterExtensions,
    contextMenu,
    // Computed
    availableExtensions,
    hasSelectionHistory,
    selectionHistoryCount,
    totalFileCount,
    selectionProgress,
    // Actions
    handleToggleSelect,
    handleToggleExpand,
    handleRefresh,
    handleFilterUpdate,
    handleSettingsChange,
    restorePreviousSelection,
    handleContextMenuShow,
    handleContextMenuAction,
    // Lifecycle
    initialize,
    cleanup,
    setupWatchers
  };
}
const FETCH_TIMEOUT_MS = 5e3;
async function fetchWithTimeout(promise, timeoutMs) {
  const timeout = new Promise((resolve) => setTimeout(() => resolve(null), timeoutMs));
  return Promise.race([promise, timeout]);
}
function useAnalysisStatus(options2) {
  const { selectedFiles, onAddFiles } = options2;
  const { t } = useI18n();
  const projectStore = useProjectStore();
  const suggestions = ref([]);
  const impactResult = ref(null);
  const selectedRelated = ref(/* @__PURE__ */ new Set());
  const isLoadingRelated = ref(false);
  const showRelatedPopup = ref(false);
  const showImpactPopup = ref(false);
  const projectPath = computed(() => projectStore.currentPath || "");
  const hasSelectedFiles = computed(() => selectedFiles.value.length > 0);
  const relatedCount = computed(() => suggestions.value.length);
  const dependentCount = computed(() => impactResult.value?.totalDependents || 0);
  const shouldShowBar = computed(
    () => hasSelectedFiles.value && (isLoadingRelated.value || relatedCount.value > 0 || dependentCount.value > 0)
  );
  let fetchTimer = null;
  async function fetchRelated() {
    if (!projectPath.value || selectedFiles.value.length === 0) {
      suggestions.value = [];
      return;
    }
    isLoadingRelated.value = true;
    try {
      const result = await fetchWithTimeout(
        apiService.getSmartSuggestions(projectPath.value, selectedFiles.value),
        FETCH_TIMEOUT_MS
      );
      if (result) {
        suggestions.value = result.suggestions;
        selectedRelated.value = new Set(result.suggestions.map((s) => s.path));
      } else {
        suggestions.value = [];
      }
    } catch {
      suggestions.value = [];
    } finally {
      isLoadingRelated.value = false;
    }
  }
  async function fetchImpact() {
    if (!projectPath.value || selectedFiles.value.length === 0) {
      impactResult.value = null;
      return;
    }
    try {
      const result = await fetchWithTimeout(
        apiService.getImpactPreview(projectPath.value, selectedFiles.value),
        FETCH_TIMEOUT_MS
      );
      impactResult.value = result;
    } catch {
      impactResult.value = null;
    }
  }
  watch(selectedFiles, () => {
    if (fetchTimer) clearTimeout(fetchTimer);
    if (selectedFiles.value.length === 0) {
      suggestions.value = [];
      impactResult.value = null;
      return;
    }
    fetchTimer = setTimeout(() => {
      fetchRelated();
      fetchImpact();
    }, 500);
  }, { immediate: true, deep: true });
  function toggleRelated(path) {
    if (selectedRelated.value.has(path)) {
      selectedRelated.value.delete(path);
    } else {
      selectedRelated.value.add(path);
    }
    selectedRelated.value = new Set(selectedRelated.value);
  }
  function addSelectedRelated() {
    if (selectedRelated.value.size > 0) {
      onAddFiles(Array.from(selectedRelated.value));
      suggestions.value = suggestions.value.filter((s) => !selectedRelated.value.has(s.path));
      selectedRelated.value.clear();
      showRelatedPopup.value = false;
    }
  }
  function getSourceLabel(source) {
    switch (source) {
      case "git":
        return t("context.sourceGitShort");
      case "arch":
        return t("context.sourceArchShort");
      case "semantic":
        return t("context.sourceSemanticShort");
      default:
        return "";
    }
  }
  function getSourceBadgeClass(source) {
    switch (source) {
      case "git":
        return "badge-git";
      case "arch":
        return "badge-arch";
      case "semantic":
        return "badge-semantic";
      default:
        return "badge-default";
    }
  }
  function getFileIconClass(path) {
    const ext = path.split(".").pop()?.toLowerCase() || "";
    const classMap = {
      "vue": "icon-vue",
      "ts": "icon-ts",
      "tsx": "icon-ts",
      "js": "icon-js",
      "jsx": "icon-js",
      "go": "icon-go",
      "py": "icon-py",
      "rs": "icon-rs",
      "json": "icon-json",
      "yaml": "icon-json",
      "yml": "icon-json",
      "css": "icon-css",
      "scss": "icon-css",
      "html": "icon-html"
    };
    return classMap[ext] || "icon-default";
  }
  function getFileName(path) {
    return path.split("/").pop() || path;
  }
  function getFilePath(path) {
    const parts = path.split("/");
    if (parts.length <= 1) return "";
    return parts.slice(0, -1).join("/");
  }
  function getRiskClass(level) {
    switch (level) {
      case "low":
        return "risk-low";
      case "medium":
        return "risk-medium";
      case "high":
        return "risk-high";
      default:
        return "";
    }
  }
  function getRiskBarClass(level) {
    switch (level) {
      case "low":
        return "bg-green-500";
      case "medium":
        return "bg-amber-500";
      case "high":
        return "bg-red-500";
      default:
        return "bg-gray-500";
    }
  }
  function getRiskLabel(level) {
    switch (level) {
      case "low":
        return t("context.riskLow");
      case "medium":
        return t("context.riskMedium");
      case "high":
        return t("context.riskHigh");
      default:
        return "";
    }
  }
  return {
    // State
    suggestions,
    impactResult,
    selectedRelated,
    isLoadingRelated,
    showRelatedPopup,
    showImpactPopup,
    // Computed
    shouldShowBar,
    relatedCount,
    dependentCount,
    // Actions
    toggleRelated,
    addSelectedRelated,
    // Helpers
    getSourceLabel,
    getSourceBadgeClass,
    getFileIconClass,
    getFileName,
    getFilePath,
    getRiskClass,
    getRiskBarClass,
    getRiskLabel
  };
}
const _hoisted_1$18 = { class: "popup-content popup-impact" };
const _hoisted_2$14 = { class: "popup-header" };
const _hoisted_3$$ = { class: "popup-title" };
const _hoisted_4$V = { class: "popup-subtitle" };
const _hoisted_5$O = { class: "popup-body" };
const _hoisted_6$K = {
  key: 0,
  class: "risk-section"
};
const _hoisted_7$G = { class: "risk-header" };
const _hoisted_8$E = { class: "risk-label" };
const _hoisted_9$B = { class: "risk-bar-bg" };
const _hoisted_10$y = {
  key: 1,
  class: "impact-section"
};
const _hoisted_11$w = { class: "impact-section-header" };
const _hoisted_12$r = { class: "popup-list" };
const _hoisted_13$r = { class: "popup-item-path" };
const _hoisted_14$o = {
  key: 2,
  class: "impact-section"
};
const _hoisted_15$m = { class: "impact-section-header" };
const _hoisted_16$j = { class: "popup-list" };
const _hoisted_17$j = { class: "popup-item-path" };
const _hoisted_18$i = {
  key: 3,
  class: "popup-empty"
};
const _sfc_main$1a = /* @__PURE__ */ defineComponent({
  __name: "ImpactAnalysisPopup",
  props: {
    visible: { type: Boolean },
    impactResult: {},
    getRiskClass: { type: Function },
    getRiskBarClass: { type: Function },
    getRiskLabel: { type: Function }
  },
  emits: ["close"],
  setup(__props) {
    const { t } = useI18n();
    return (_ctx, _cache) => {
      return openBlock(), createBlock(Teleport, { to: "body" }, [
        __props.visible ? (openBlock(), createElementBlock("div", {
          key: 0,
          class: "popup-overlay",
          onClick: _cache[1] || (_cache[1] = withModifiers(($event) => _ctx.$emit("close"), ["self"]))
        }, [
          createBaseVNode("div", _hoisted_1$18, [
            createBaseVNode("div", _hoisted_2$14, [
              createBaseVNode("div", null, [
                createBaseVNode("h3", _hoisted_3$$, toDisplayString(unref(t)("context.impactTitle")), 1),
                createBaseVNode("p", _hoisted_4$V, toDisplayString(unref(t)("context.impactSubtitle")), 1)
              ]),
              createBaseVNode("button", {
                class: "popup-close",
                onClick: _cache[0] || (_cache[0] = ($event) => _ctx.$emit("close"))
              }, [..._cache[2] || (_cache[2] = [
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
            createBaseVNode("div", _hoisted_5$O, [
              __props.impactResult ? (openBlock(), createElementBlock("div", _hoisted_6$K, [
                createBaseVNode("div", _hoisted_7$G, [
                  createBaseVNode("span", _hoisted_8$E, toDisplayString(unref(t)("context.riskScore")), 1),
                  createBaseVNode("span", {
                    class: normalizeClass(["risk-value", __props.getRiskClass(__props.impactResult.riskLevel)])
                  }, toDisplayString(__props.getRiskLabel(__props.impactResult.riskLevel)), 3)
                ]),
                createBaseVNode("div", _hoisted_9$B, [
                  createBaseVNode("div", {
                    class: normalizeClass(["risk-bar", __props.getRiskBarClass(__props.impactResult.riskLevel)]),
                    style: normalizeStyle({ width: `${Math.round(__props.impactResult.aggregateRisk * 100)}%` })
                  }, null, 6)
                ])
              ])) : createCommentVNode("", true),
              __props.impactResult?.affectedFiles.length ? (openBlock(), createElementBlock("div", _hoisted_10$y, [
                createBaseVNode("div", _hoisted_11$w, toDisplayString(unref(t)("context.affectedFiles")) + " (" + toDisplayString(__props.impactResult.affectedFiles.length) + ") ", 1),
                createBaseVNode("div", _hoisted_12$r, [
                  (openBlock(true), createElementBlock(Fragment, null, renderList(__props.impactResult.affectedFiles, (file) => {
                    return openBlock(), createElementBlock("div", {
                      key: file.path,
                      class: "popup-item popup-item-readonly"
                    }, [
                      _cache[3] || (_cache[3] = createBaseVNode("span", { class: "popup-item-icon" }, "📄", -1)),
                      createBaseVNode("span", _hoisted_13$r, toDisplayString(file.path), 1),
                      createBaseVNode("span", {
                        class: normalizeClass(["popup-item-type", file.type === "direct" ? "type-direct" : "type-transitive"])
                      }, toDisplayString(file.type === "direct" ? unref(t)("context.directDep") : unref(t)("context.transitiveDep")), 3)
                    ]);
                  }), 128))
                ])
              ])) : createCommentVNode("", true),
              __props.impactResult?.relatedTests.length ? (openBlock(), createElementBlock("div", _hoisted_14$o, [
                createBaseVNode("div", _hoisted_15$m, " 🧪 " + toDisplayString(unref(t)("context.relatedTests")) + " (" + toDisplayString(__props.impactResult.relatedTests.length) + ") ", 1),
                createBaseVNode("div", _hoisted_16$j, [
                  (openBlock(true), createElementBlock(Fragment, null, renderList(__props.impactResult.relatedTests, (test) => {
                    return openBlock(), createElementBlock("div", {
                      key: test,
                      class: "popup-item popup-item-readonly"
                    }, [
                      _cache[4] || (_cache[4] = createBaseVNode("span", { class: "popup-item-icon" }, "🧪", -1)),
                      createBaseVNode("span", _hoisted_17$j, toDisplayString(test), 1)
                    ]);
                  }), 128))
                ])
              ])) : createCommentVNode("", true),
              !__props.impactResult || __props.impactResult.totalDependents === 0 ? (openBlock(), createElementBlock("div", _hoisted_18$i, toDisplayString(unref(t)("context.noImpactFiles")), 1)) : createCommentVNode("", true)
            ])
          ])
        ])) : createCommentVNode("", true)
      ]);
    };
  }
});
const ImpactAnalysisPopup = /* @__PURE__ */ _export_sfc(_sfc_main$1a, [["__scopeId", "data-v-bc0c2bd8"]]);
const _hoisted_1$17 = ["checked"];
const _hoisted_2$13 = {
  key: 0,
  viewBox: "0 0 12 12",
  fill: "none"
};
const _hoisted_3$_ = { class: "smart-hud-file-info" };
const _hoisted_4$U = { class: "smart-hud-file-name" };
const _hoisted_5$N = { class: "smart-hud-file-path" };
const _sfc_main$19 = /* @__PURE__ */ defineComponent({
  __name: "SmartSuggestionItem",
  props: {
    isSelected: { type: Boolean },
    iconClass: {},
    fileName: {},
    filePath: {},
    badgeClass: {},
    sourceLabel: {}
  },
  emits: ["toggle"],
  setup(__props) {
    return (_ctx, _cache) => {
      return openBlock(), createElementBlock("label", {
        class: normalizeClass(["smart-hud-item", { "smart-hud-item-selected": __props.isSelected }])
      }, [
        createBaseVNode("div", {
          class: normalizeClass(["smart-hud-checkbox", { checked: __props.isSelected }])
        }, [
          createBaseVNode("input", {
            type: "checkbox",
            checked: __props.isSelected,
            onChange: _cache[0] || (_cache[0] = ($event) => _ctx.$emit("toggle"))
          }, null, 40, _hoisted_1$17),
          __props.isSelected ? (openBlock(), createElementBlock("svg", _hoisted_2$13, [..._cache[1] || (_cache[1] = [
            createBaseVNode("path", {
              d: "M2 6L5 9L10 3",
              stroke: "currentColor",
              "stroke-width": "2",
              "stroke-linecap": "round",
              "stroke-linejoin": "round"
            }, null, -1)
          ])])) : createCommentVNode("", true)
        ], 2),
        createBaseVNode("div", {
          class: normalizeClass(["smart-hud-file-icon", __props.iconClass])
        }, [..._cache[2] || (_cache[2] = [
          createBaseVNode("svg", {
            viewBox: "0 0 24 24",
            fill: "none",
            stroke: "currentColor",
            "stroke-width": "1.5"
          }, [
            createBaseVNode("path", { d: "M14 2H6a2 2 0 00-2 2v16a2 2 0 002 2h12a2 2 0 002-2V8l-6-6z" }),
            createBaseVNode("path", { d: "M14 2v6h6M10 12l-2 2 2 2M14 12l2 2-2 2" })
          ], -1)
        ])], 2),
        createBaseVNode("div", _hoisted_3$_, [
          createBaseVNode("span", _hoisted_4$U, toDisplayString(__props.fileName), 1),
          createBaseVNode("span", _hoisted_5$N, toDisplayString(__props.filePath), 1)
        ]),
        createBaseVNode("span", {
          class: normalizeClass(["smart-hud-badge", __props.badgeClass])
        }, [
          _cache[3] || (_cache[3] = createBaseVNode("span", { class: "smart-hud-badge-icon" }, "🔗", -1)),
          createTextVNode(" " + toDisplayString(__props.sourceLabel), 1)
        ], 2)
      ], 2);
    };
  }
});
const SmartSuggestionItem = /* @__PURE__ */ _export_sfc(_sfc_main$19, [["__scopeId", "data-v-137e1cad"]]);
const _hoisted_1$16 = { class: "smart-hud" };
const _hoisted_2$12 = { class: "smart-hud-header" };
const _hoisted_3$Z = { class: "smart-hud-title-row" };
const _hoisted_4$T = { class: "smart-hud-title-block" };
const _hoisted_5$M = { class: "smart-hud-title" };
const _hoisted_6$J = { class: "smart-hud-subtitle-inline" };
const _hoisted_7$F = { class: "smart-hud-body" };
const _hoisted_8$D = {
  key: 0,
  class: "smart-hud-empty"
};
const _hoisted_9$A = {
  key: 1,
  class: "smart-hud-list"
};
const _hoisted_10$x = {
  key: 0,
  class: "smart-hud-footer"
};
const _hoisted_11$v = ["disabled"];
const _sfc_main$18 = /* @__PURE__ */ defineComponent({
  __name: "SmartSuggestionsHud",
  props: {
    visible: { type: Boolean },
    suggestions: {},
    selectedPaths: {},
    getFileIconClass: { type: Function },
    getFileName: { type: Function },
    getFilePath: { type: Function },
    getSourceBadgeClass: { type: Function },
    getSourceLabel: { type: Function }
  },
  emits: ["close", "toggle", "add"],
  setup(__props) {
    const { t } = useI18n();
    return (_ctx, _cache) => {
      return openBlock(), createBlock(Teleport, { to: "body" }, [
        __props.visible ? (openBlock(), createElementBlock("div", {
          key: 0,
          class: "smart-hud-overlay",
          onClick: _cache[2] || (_cache[2] = withModifiers(($event) => _ctx.$emit("close"), ["self"]))
        }, [
          createBaseVNode("div", _hoisted_1$16, [
            createBaseVNode("div", _hoisted_2$12, [
              createBaseVNode("div", _hoisted_3$Z, [
                _cache[3] || (_cache[3] = createBaseVNode("div", { class: "smart-hud-sparkle-icon" }, "✨", -1)),
                createBaseVNode("div", _hoisted_4$T, [
                  createBaseVNode("span", _hoisted_5$M, toDisplayString(unref(t)("context.aiRecommendations")), 1),
                  createBaseVNode("span", _hoisted_6$J, toDisplayString(unref(t)("context.foundFilesAnalysis", { count: __props.suggestions.length })), 1)
                ])
              ]),
              createBaseVNode("button", {
                class: "smart-hud-close",
                onClick: _cache[0] || (_cache[0] = ($event) => _ctx.$emit("close"))
              }, [..._cache[4] || (_cache[4] = [
                createBaseVNode("svg", {
                  width: "14",
                  height: "14",
                  viewBox: "0 0 24 24",
                  fill: "none",
                  stroke: "currentColor",
                  "stroke-width": "2"
                }, [
                  createBaseVNode("path", { d: "M18 6L6 18M6 6l12 12" })
                ], -1)
              ])])
            ]),
            createBaseVNode("div", _hoisted_7$F, [
              __props.suggestions.length === 0 ? (openBlock(), createElementBlock("div", _hoisted_8$D, toDisplayString(unref(t)("context.noRelatedFiles")), 1)) : (openBlock(), createElementBlock("div", _hoisted_9$A, [
                (openBlock(true), createElementBlock(Fragment, null, renderList(__props.suggestions, (item) => {
                  return openBlock(), createBlock(SmartSuggestionItem, {
                    key: item.path,
                    "is-selected": __props.selectedPaths.has(item.path),
                    "icon-class": __props.getFileIconClass(item.path),
                    "file-name": __props.getFileName(item.path),
                    "file-path": __props.getFilePath(item.path),
                    "badge-class": __props.getSourceBadgeClass(item.source),
                    "source-label": __props.getSourceLabel(item.source),
                    onToggle: ($event) => _ctx.$emit("toggle", item.path)
                  }, null, 8, ["is-selected", "icon-class", "file-name", "file-path", "badge-class", "source-label", "onToggle"]);
                }), 128))
              ]))
            ]),
            __props.suggestions.length > 0 ? (openBlock(), createElementBlock("div", _hoisted_10$x, [
              createBaseVNode("button", {
                class: "smart-hud-action group",
                disabled: __props.selectedPaths.size === 0,
                onClick: _cache[1] || (_cache[1] = ($event) => _ctx.$emit("add"))
              }, [
                _cache[5] || (_cache[5] = createBaseVNode("svg", {
                  class: "smart-hud-action-icon",
                  width: "18",
                  height: "18",
                  viewBox: "0 0 24 24",
                  fill: "none",
                  stroke: "currentColor",
                  "stroke-width": "2.5"
                }, [
                  createBaseVNode("path", { d: "M12 5v14M5 12h14" })
                ], -1)),
                createBaseVNode("span", null, toDisplayString(unref(t)("context.addFiles")) + " (" + toDisplayString(__props.selectedPaths.size) + ")", 1),
                _cache[6] || (_cache[6] = createBaseVNode("div", { class: "smart-hud-shimmer" }, null, -1))
              ], 8, _hoisted_11$v)
            ])) : createCommentVNode("", true)
          ])
        ])) : createCommentVNode("", true)
      ]);
    };
  }
});
const SmartSuggestionsHud = /* @__PURE__ */ _export_sfc(_sfc_main$18, [["__scopeId", "data-v-ecfe53cc"]]);
const _hoisted_1$15 = {
  key: 0,
  class: "analysis-status-bar"
};
const _hoisted_2$11 = {
  key: 0,
  class: "status-btn status-btn-loading"
};
const _hoisted_3$Y = { class: "status-label" };
const _hoisted_4$S = ["title"];
const _hoisted_5$L = { class: "status-count" };
const _hoisted_6$I = { class: "status-label" };
const _hoisted_7$E = ["title"];
const _hoisted_8$C = { class: "status-count" };
const _hoisted_9$z = { class: "status-label" };
const _sfc_main$17 = /* @__PURE__ */ defineComponent({
  __name: "AnalysisStatusBar",
  props: {
    selectedFiles: {}
  },
  emits: ["add-files"],
  setup(__props, { emit: __emit }) {
    const props = __props;
    const emit = __emit;
    const { t } = useI18n();
    const {
      // State
      suggestions,
      impactResult,
      selectedRelated,
      isLoadingRelated,
      showRelatedPopup,
      showImpactPopup,
      // Computed
      shouldShowBar,
      relatedCount,
      dependentCount,
      // Actions
      toggleRelated,
      addSelectedRelated,
      // Helpers
      getSourceLabel,
      getSourceBadgeClass,
      getFileIconClass,
      getFileName,
      getFilePath,
      getRiskClass,
      getRiskBarClass,
      getRiskLabel
    } = useAnalysisStatus({
      selectedFiles: toRef(props, "selectedFiles"),
      onAddFiles: (files2) => emit("add-files", files2)
    });
    return (_ctx, _cache) => {
      return unref(shouldShowBar) ? (openBlock(), createElementBlock("div", _hoisted_1$15, [
        unref(isLoadingRelated) ? (openBlock(), createElementBlock("div", _hoisted_2$11, [
          _cache[4] || (_cache[4] = createBaseVNode("span", { class: "status-spinner" }, null, -1)),
          createBaseVNode("span", _hoisted_3$Y, toDisplayString(unref(t)("context.analyzingFiles")), 1)
        ])) : unref(relatedCount) > 0 ? (openBlock(), createElementBlock("button", {
          key: 1,
          class: "status-btn status-btn-related",
          title: unref(t)("context.relatedFilesTooltip"),
          onClick: _cache[0] || (_cache[0] = ($event) => showRelatedPopup.value = true)
        }, [
          _cache[5] || (_cache[5] = createBaseVNode("span", { class: "status-icon" }, "💡", -1)),
          createBaseVNode("span", _hoisted_5$L, toDisplayString(unref(relatedCount)), 1),
          createBaseVNode("span", _hoisted_6$I, toDisplayString(unref(t)("context.relatedFiles")), 1)
        ], 8, _hoisted_4$S)) : createCommentVNode("", true),
        unref(dependentCount) > 0 ? (openBlock(), createElementBlock("button", {
          key: 2,
          class: "status-btn status-btn-impact",
          title: unref(t)("context.dependentFilesTooltip"),
          onClick: _cache[1] || (_cache[1] = ($event) => showImpactPopup.value = true)
        }, [
          _cache[6] || (_cache[6] = createBaseVNode("span", { class: "status-icon" }, "✨", -1)),
          createBaseVNode("span", _hoisted_8$C, toDisplayString(unref(dependentCount)), 1),
          createBaseVNode("span", _hoisted_9$z, toDisplayString(unref(t)("context.recommendations")), 1)
        ], 8, _hoisted_7$E)) : createCommentVNode("", true),
        createVNode(SmartSuggestionsHud, {
          visible: unref(showRelatedPopup),
          suggestions: unref(suggestions),
          "selected-paths": unref(selectedRelated),
          "get-file-icon-class": unref(getFileIconClass),
          "get-file-name": unref(getFileName),
          "get-file-path": unref(getFilePath),
          "get-source-badge-class": unref(getSourceBadgeClass),
          "get-source-label": unref(getSourceLabel),
          onClose: _cache[2] || (_cache[2] = ($event) => showRelatedPopup.value = false),
          onToggle: unref(toggleRelated),
          onAdd: unref(addSelectedRelated)
        }, null, 8, ["visible", "suggestions", "selected-paths", "get-file-icon-class", "get-file-name", "get-file-path", "get-source-badge-class", "get-source-label", "onToggle", "onAdd"]),
        createVNode(ImpactAnalysisPopup, {
          visible: unref(showImpactPopup),
          "impact-result": unref(impactResult),
          "get-risk-class": unref(getRiskClass),
          "get-risk-bar-class": unref(getRiskBarClass),
          "get-risk-label": unref(getRiskLabel),
          onClose: _cache[3] || (_cache[3] = ($event) => showImpactPopup.value = false)
        }, null, 8, ["visible", "impact-result", "get-risk-class", "get-risk-bar-class", "get-risk-label"])
      ])) : createCommentVNode("", true);
    };
  }
});
const AnalysisStatusBar = /* @__PURE__ */ _export_sfc(_sfc_main$17, [["__scopeId", "data-v-23144f59"]]);
const _hoisted_1$14 = { class: "magic-bar-wrapper" };
const _hoisted_2$10 = { class: "magic-bar" };
const _hoisted_3$X = { class: "limit-main" };
const _hoisted_4$R = { class: "limit-value" };
const _hoisted_5$K = {
  key: 0,
  class: "limit-dropdown"
};
const _hoisted_6$H = ["onClick"];
const _hoisted_7$D = { class: "limit-option-value" };
const _hoisted_8$B = { class: "limit-option-model" };
const _hoisted_9$y = { class: "limit-custom" };
const _hoisted_10$w = ["placeholder"];
const _hoisted_11$u = ["disabled"];
const _hoisted_12$q = { class: "build-content" };
const _hoisted_13$q = {
  key: 0,
  class: "build-icon",
  fill: "none",
  viewBox: "0 0 24 24",
  stroke: "currentColor"
};
const _hoisted_14$n = {
  key: 1,
  class: "build-spinner",
  fill: "none",
  viewBox: "0 0 24 24"
};
const _hoisted_15$l = { class: "build-text" };
const _hoisted_16$i = {
  key: 0,
  class: "file-counter"
};
const _hoisted_17$i = { class: "file-count" };
const _hoisted_18$h = { class: "token-estimate" };
const _hoisted_19$g = ["title"];
const _sfc_main$16 = /* @__PURE__ */ defineComponent({
  __name: "CommandBar",
  props: {
    selectedCount: {},
    isBuilding: { type: Boolean }
  },
  emits: ["build"],
  setup(__props, { emit: __emit }) {
    const props = __props;
    const emit = __emit;
    const { t } = useI18n();
    const settingsStore = useSettingsStore();
    const fileStore = useFileStore();
    const settings2 = computed(() => settingsStore.settings.context);
    const showDropdown = ref(false);
    const limitRef = ref(null);
    const customTokenValue = ref("");
    const isCustomFocused = ref(false);
    const isDisabled = computed(() => props.selectedCount === 0);
    const isButtonDisabled = computed(() => props.selectedCount === 0 || props.isBuilding);
    const estimatedTokens = computed(() => Math.round(fileStore.estimatedTokenCount / 1e3));
    const tokenPresets = [
      { value: 32e3, label: "32K", model: "GPT-4" },
      { value: 128e3, label: "128K", model: "GPT-4 Turbo" },
      { value: 2e5, label: "200K", model: "Claude" },
      { value: 1e6, label: "1M", model: "Gemini" }
    ];
    function formatTokens(n) {
      if (n >= 1e6) return `${(n / 1e6).toFixed(n % 1e6 === 0 ? 0 : 1)}M`;
      if (n >= 1e3) return `${Math.round(n / 1e3)}K`;
      return n.toString();
    }
    function isPresetActive(value) {
      return Math.abs(settings2.value.maxTokens - value) <= value * 0.05;
    }
    function toggleDropdown() {
      showDropdown.value = !showDropdown.value;
    }
    function selectPreset(value) {
      settingsStore.updateContextSettings({ maxTokens: value });
      showDropdown.value = false;
    }
    function applyCustomLimit() {
      const input = customTokenValue.value.replace(/[^\d.]/g, "");
      const value = parseFloat(input);
      if (!isNaN(value) && value > 0) {
        const tokens = Math.round(value * 1e3);
        const clamped = Math.min(Math.max(tokens, 1e3), 1e7);
        settingsStore.updateContextSettings({ maxTokens: clamped });
        customTokenValue.value = "";
        showDropdown.value = false;
      }
    }
    function handleBuild() {
      if (!isDisabled.value) {
        emit("build");
      }
    }
    function handleClear() {
      fileStore.clearSelection();
    }
    function handleClickOutside(event) {
      if (showDropdown.value && limitRef.value && !limitRef.value.contains(event.target)) {
        showDropdown.value = false;
      }
    }
    onMounted(() => document.addEventListener("click", handleClickOutside));
    onUnmounted(() => document.removeEventListener("click", handleClickOutside));
    return (_ctx, _cache) => {
      return openBlock(), createElementBlock("div", _hoisted_1$14, [
        createBaseVNode("div", _hoisted_2$10, [
          createBaseVNode("div", {
            class: "limit-section",
            onClick: toggleDropdown,
            ref_key: "limitRef",
            ref: limitRef
          }, [
            createBaseVNode("div", _hoisted_3$X, [
              createBaseVNode("span", _hoisted_4$R, toDisplayString(formatTokens(settings2.value.maxTokens)), 1),
              (openBlock(), createElementBlock("svg", {
                class: normalizeClass(["limit-chevron", { open: showDropdown.value }]),
                fill: "none",
                viewBox: "0 0 24 24",
                stroke: "currentColor"
              }, [..._cache[4] || (_cache[4] = [
                createBaseVNode("path", {
                  "stroke-linecap": "round",
                  "stroke-linejoin": "round",
                  "stroke-width": "2.5",
                  d: "M19 9l-7 7-7-7"
                }, null, -1)
              ])], 2))
            ]),
            createVNode(Transition, { name: "dropdown" }, {
              default: withCtx(() => [
                showDropdown.value ? (openBlock(), createElementBlock("div", _hoisted_5$K, [
                  (openBlock(), createElementBlock(Fragment, null, renderList(tokenPresets, (preset) => {
                    return createBaseVNode("button", {
                      key: preset.value,
                      onClick: withModifiers(($event) => selectPreset(preset.value), ["stop"]),
                      class: normalizeClass(["limit-option", { active: isPresetActive(preset.value) }])
                    }, [
                      createBaseVNode("span", _hoisted_7$D, toDisplayString(preset.label), 1),
                      createBaseVNode("span", _hoisted_8$B, toDisplayString(preset.model), 1)
                    ], 10, _hoisted_6$H);
                  }), 64)),
                  createBaseVNode("div", _hoisted_9$y, [
                    withDirectives(createBaseVNode("input", {
                      ref: "customInputRef",
                      "onUpdate:modelValue": _cache[0] || (_cache[0] = ($event) => customTokenValue.value = $event),
                      type: "text",
                      inputmode: "numeric",
                      class: "limit-custom-input",
                      placeholder: unref(t)("commandBar.customLimit"),
                      onClick: _cache[1] || (_cache[1] = withModifiers(() => {
                      }, ["stop"])),
                      onKeydown: withKeys(applyCustomLimit, ["enter"]),
                      onFocus: _cache[2] || (_cache[2] = ($event) => isCustomFocused.value = true),
                      onBlur: _cache[3] || (_cache[3] = ($event) => isCustomFocused.value = false)
                    }, null, 40, _hoisted_10$w), [
                      [vModelText, customTokenValue.value]
                    ]),
                    _cache[6] || (_cache[6] = createBaseVNode("span", { class: "limit-custom-suffix" }, "K", -1)),
                    customTokenValue.value ? (openBlock(), createElementBlock("button", {
                      key: 0,
                      class: "limit-custom-apply",
                      onClick: withModifiers(applyCustomLimit, ["stop"])
                    }, [..._cache[5] || (_cache[5] = [
                      createBaseVNode("svg", {
                        class: "w-3.5 h-3.5",
                        fill: "none",
                        stroke: "currentColor",
                        viewBox: "0 0 24 24"
                      }, [
                        createBaseVNode("path", {
                          "stroke-linecap": "round",
                          "stroke-linejoin": "round",
                          "stroke-width": "2",
                          d: "M5 13l4 4L19 7"
                        })
                      ], -1)
                    ])])) : createCommentVNode("", true)
                  ])
                ])) : createCommentVNode("", true)
              ]),
              _: 1
            })
          ], 512),
          _cache[13] || (_cache[13] = createBaseVNode("div", { class: "bar-divider" }, null, -1)),
          createBaseVNode("button", {
            class: normalizeClass(["build-section", { disabled: isDisabled.value && !__props.isBuilding, loading: __props.isBuilding }]),
            disabled: isButtonDisabled.value,
            onClick: handleBuild
          }, [
            _cache[9] || (_cache[9] = createBaseVNode("div", { class: "build-bg" }, null, -1)),
            _cache[10] || (_cache[10] = createBaseVNode("div", { class: "build-ring" }, null, -1)),
            _cache[11] || (_cache[11] = createBaseVNode("div", { class: "build-glow" }, null, -1)),
            _cache[12] || (_cache[12] = createBaseVNode("div", { class: "build-shimmer" }, null, -1)),
            createBaseVNode("div", _hoisted_12$q, [
              !__props.isBuilding ? (openBlock(), createElementBlock("svg", _hoisted_13$q, [..._cache[7] || (_cache[7] = [
                createBaseVNode("path", {
                  "stroke-linecap": "round",
                  "stroke-linejoin": "round",
                  "stroke-width": "2.5",
                  d: "M13 10V3L4 14h7v7l9-11h-7z"
                }, null, -1)
              ])])) : (openBlock(), createElementBlock("svg", _hoisted_14$n, [..._cache[8] || (_cache[8] = [
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
              createBaseVNode("span", _hoisted_15$l, toDisplayString(__props.isBuilding ? unref(t)("context.building") : unref(t)("commandBar.build")), 1)
            ])
          ], 10, _hoisted_11$u)
        ]),
        __props.selectedCount > 0 ? (openBlock(), createElementBlock("div", _hoisted_16$i, [
          createTextVNode(toDisplayString(unref(t)("commandBar.selected")) + ": ", 1),
          createBaseVNode("span", _hoisted_17$i, toDisplayString(__props.selectedCount), 1),
          createBaseVNode("span", _hoisted_18$h, "~" + toDisplayString(estimatedTokens.value) + "k tokens", 1),
          createBaseVNode("button", {
            class: "clear-btn",
            onClick: handleClear,
            title: unref(t)("files.clearSelection")
          }, [
            _cache[14] || (_cache[14] = createBaseVNode("svg", {
              class: "w-3 h-3",
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
            ], -1)),
            createTextVNode(" " + toDisplayString(unref(t)("files.clear")), 1)
          ], 8, _hoisted_19$g)
        ])) : createCommentVNode("", true)
      ]);
    };
  }
});
const CommandBar = /* @__PURE__ */ _export_sfc(_sfc_main$16, [["__scopeId", "data-v-3f913954"]]);
const _hoisted_1$13 = ["aria-label"];
const _hoisted_2$$ = ["aria-label"];
const _hoisted_3$W = {
  key: 2,
  class: "h-px bg-gray-700 my-1"
};
const _hoisted_4$Q = {
  key: 4,
  class: "h-px bg-gray-700 my-1"
};
const _sfc_main$15 = /* @__PURE__ */ defineComponent({
  __name: "FileContextMenu",
  props: {
    node: {},
    position: {},
    visible: { type: Boolean }
  },
  emits: ["action", "close"],
  setup(__props, { emit: __emit }) {
    const { t } = useI18n();
    const props = __props;
    const emit = __emit;
    function handleAction(type) {
      if (props.node) {
        emit("action", { type, node: props.node });
      }
      emit("close");
    }
    return (_ctx, _cache) => {
      return openBlock(), createBlock(Teleport, { to: "body" }, [
        createVNode(Transition, { name: "context-menu" }, {
          default: withCtx(() => [
            __props.visible && __props.node ? (openBlock(), createElementBlock("div", {
              key: 0,
              role: "menu",
              class: "context-menu fixed z-[1060] bg-gray-800 border-2 border-gray-600 rounded-lg shadow-2xl py-1 min-w-[220px]",
              style: normalizeStyle({ left: `${__props.position.x}px`, top: `${__props.position.y}px` }),
              onClick: _cache[9] || (_cache[9] = withModifiers(() => {
              }, ["stop"]))
            }, [
              __props.node.isDir ? (openBlock(), createElementBlock("button", {
                key: 0,
                onClick: _cache[0] || (_cache[0] = ($event) => handleAction("selectAll")),
                role: "menuitem",
                "aria-label": unref(t)("contextMenu.selectAll"),
                class: "w-full px-4 py-2 text-left text-sm text-white hover:bg-gray-700 hover:scale-[1.01] active:scale-[0.99] flex items-center gap-3 transition-all duration-150"
              }, [
                _cache[10] || (_cache[10] = createBaseVNode("svg", {
                  class: "w-4 h-4",
                  fill: "none",
                  stroke: "currentColor",
                  viewBox: "0 0 24 24"
                }, [
                  createBaseVNode("path", {
                    "stroke-linecap": "round",
                    "stroke-linejoin": "round",
                    "stroke-width": "2",
                    d: "M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"
                  })
                ], -1)),
                createTextVNode(" " + toDisplayString(unref(t)("contextMenu.selectAll")), 1)
              ], 8, _hoisted_1$13)) : createCommentVNode("", true),
              __props.node.isDir ? (openBlock(), createElementBlock("button", {
                key: 1,
                onClick: _cache[1] || (_cache[1] = ($event) => handleAction("deselectAll")),
                role: "menuitem",
                "aria-label": unref(t)("contextMenu.deselectAll"),
                class: "w-full px-4 py-2 text-left text-sm text-white hover:bg-gray-700 hover:scale-[1.01] active:scale-[0.99] flex items-center gap-3 transition-all duration-150"
              }, [
                _cache[11] || (_cache[11] = createBaseVNode("svg", {
                  class: "w-4 h-4",
                  fill: "none",
                  stroke: "currentColor",
                  viewBox: "0 0 24 24"
                }, [
                  createBaseVNode("path", {
                    "stroke-linecap": "round",
                    "stroke-linejoin": "round",
                    "stroke-width": "2",
                    d: "M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z"
                  })
                ], -1)),
                createTextVNode(" " + toDisplayString(unref(t)("contextMenu.deselectAll")), 1)
              ], 8, _hoisted_2$$)) : createCommentVNode("", true),
              __props.node.isDir ? (openBlock(), createElementBlock("div", _hoisted_3$W)) : createCommentVNode("", true),
              !__props.node.isDir ? (openBlock(), createElementBlock("button", {
                key: 3,
                onClick: _cache[2] || (_cache[2] = ($event) => handleAction("quickLook")),
                class: "w-full px-4 py-2 text-left text-sm text-white hover:bg-gray-700 flex items-center gap-3 transition-colors"
              }, [
                _cache[12] || (_cache[12] = createBaseVNode("svg", {
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
                createTextVNode(" " + toDisplayString(unref(t)("files.quickLook")) + " ", 1),
                _cache[13] || (_cache[13] = createBaseVNode("span", { class: "ml-auto text-xs text-gray-400" }, "Space", -1))
              ])) : createCommentVNode("", true),
              !__props.node.isDir ? (openBlock(), createElementBlock("div", _hoisted_4$Q)) : createCommentVNode("", true),
              createBaseVNode("button", {
                onClick: _cache[3] || (_cache[3] = ($event) => handleAction("copyPath")),
                class: "w-full px-4 py-2 text-left text-sm text-white hover:bg-gray-700 flex items-center gap-3 transition-colors"
              }, [
                _cache[14] || (_cache[14] = createBaseVNode("svg", {
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
                createTextVNode(" " + toDisplayString(unref(t)("contextMenu.copyPath")), 1)
              ]),
              createBaseVNode("button", {
                onClick: _cache[4] || (_cache[4] = ($event) => handleAction("copyRelativePath")),
                class: "w-full px-4 py-2 text-left text-sm text-white hover:bg-gray-700 flex items-center gap-3 transition-colors"
              }, [
                _cache[15] || (_cache[15] = createBaseVNode("svg", {
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
                ], -1)),
                createTextVNode(" " + toDisplayString(unref(t)("contextMenu.copyRelativePath")), 1)
              ]),
              _cache[21] || (_cache[21] = createBaseVNode("div", { class: "h-px bg-gray-700 my-1" }, null, -1)),
              !__props.node.isIgnored ? (openBlock(), createElementBlock("button", {
                key: 5,
                onClick: _cache[5] || (_cache[5] = ($event) => handleAction("addToCustomIgnore")),
                class: "w-full px-4 py-2 text-left text-sm text-white hover:bg-gray-700 flex items-center gap-3 transition-colors"
              }, [
                _cache[16] || (_cache[16] = createBaseVNode("svg", {
                  class: "w-4 h-4",
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
                ], -1)),
                createTextVNode(" " + toDisplayString(unref(t)("contextMenu.addToIgnore")), 1)
              ])) : createCommentVNode("", true),
              __props.node.isIgnored ? (openBlock(), createElementBlock("button", {
                key: 6,
                onClick: _cache[6] || (_cache[6] = ($event) => handleAction("removeFromIgnore")),
                class: "w-full px-4 py-2 text-left text-sm text-orange-400 hover:bg-gray-700 flex items-center gap-3 transition-colors"
              }, [
                _cache[17] || (_cache[17] = createBaseVNode("svg", {
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
                createTextVNode(" " + toDisplayString(unref(t)("contextMenu.removeFromIgnore")), 1)
              ])) : createCommentVNode("", true),
              __props.node.isDir ? (openBlock(), createElementBlock(Fragment, { key: 7 }, [
                _cache[20] || (_cache[20] = createBaseVNode("div", { class: "h-px bg-gray-700 my-1" }, null, -1)),
                createBaseVNode("button", {
                  onClick: _cache[7] || (_cache[7] = ($event) => handleAction("expandAll")),
                  class: "w-full px-4 py-2 text-left text-sm text-white hover:bg-gray-700 flex items-center gap-3 transition-colors"
                }, [
                  _cache[18] || (_cache[18] = createBaseVNode("svg", {
                    class: "w-4 h-4",
                    fill: "none",
                    stroke: "currentColor",
                    viewBox: "0 0 24 24"
                  }, [
                    createBaseVNode("path", {
                      "stroke-linecap": "round",
                      "stroke-linejoin": "round",
                      "stroke-width": "2",
                      d: "M19 9l-7 7-7-7"
                    })
                  ], -1)),
                  createTextVNode(" " + toDisplayString(unref(t)("contextMenu.expandAll")), 1)
                ]),
                createBaseVNode("button", {
                  onClick: _cache[8] || (_cache[8] = ($event) => handleAction("collapseAll")),
                  class: "w-full px-4 py-2 text-left text-sm text-white hover:bg-gray-700 flex items-center gap-3 transition-colors"
                }, [
                  _cache[19] || (_cache[19] = createBaseVNode("svg", {
                    class: "w-4 h-4",
                    fill: "none",
                    stroke: "currentColor",
                    viewBox: "0 0 24 24"
                  }, [
                    createBaseVNode("path", {
                      "stroke-linecap": "round",
                      "stroke-linejoin": "round",
                      "stroke-width": "2",
                      d: "M5 15l7-7 7 7"
                    })
                  ], -1)),
                  createTextVNode(" " + toDisplayString(unref(t)("contextMenu.collapseAll")), 1)
                ])
              ], 64)) : createCommentVNode("", true)
            ], 4)) : createCommentVNode("", true)
          ]),
          _: 1
        })
      ]);
    };
  }
});
const FileContextMenu = /* @__PURE__ */ _export_sfc(_sfc_main$15, [["__scopeId", "data-v-93e8976a"]]);
function useFilterDropdown() {
  const openDropdown = ref(null);
  const dropdownStyle = ref({});
  const focusedIndex = ref(-1);
  const dropdownRefs = {
    types: ref(null),
    langs: ref(null),
    smart: ref(null)
  };
  function toggleDropdown(type) {
    if (openDropdown.value === type) {
      openDropdown.value = null;
      focusedIndex.value = -1;
      return;
    }
    const refElement = dropdownRefs[type].value;
    if (refElement) {
      const rect = refElement.getBoundingClientRect();
      dropdownStyle.value = {
        position: "fixed",
        top: `${rect.bottom + 4}px`,
        left: `${rect.left}px`,
        zIndex: "1100"
      };
    }
    openDropdown.value = type;
    focusedIndex.value = -1;
  }
  function closeDropdown() {
    openDropdown.value = null;
    focusedIndex.value = -1;
  }
  function handleOutsideClick(e) {
    const target2 = e.target;
    if (!target2.closest(".filter-dropdown") && !target2.closest(".filter-dropdown-menu")) {
      closeDropdown();
    }
  }
  function handleKeydown(e) {
    if (!openDropdown.value) return;
    const menu = document.querySelector(".filter-dropdown-menu");
    if (!menu) return;
    const items = menu.querySelectorAll('[role="option"], .filter-item');
    const itemCount = items.length;
    switch (e.key) {
      case "Escape":
        e.preventDefault();
        closeDropdown();
        const triggerRef2 = dropdownRefs[openDropdown.value]?.value;
        triggerRef2?.querySelector("button")?.focus();
        break;
      case "ArrowDown":
        e.preventDefault();
        focusedIndex.value = focusedIndex.value < itemCount - 1 ? focusedIndex.value + 1 : 0;
        items[focusedIndex.value]?.focus();
        break;
      case "ArrowUp":
        e.preventDefault();
        focusedIndex.value = focusedIndex.value > 0 ? focusedIndex.value - 1 : itemCount - 1;
        items[focusedIndex.value]?.focus();
        break;
      case "Home":
        e.preventDefault();
        focusedIndex.value = 0;
        items[0]?.focus();
        break;
      case "End":
        e.preventDefault();
        focusedIndex.value = itemCount - 1;
        items[itemCount - 1]?.focus();
        break;
      case "Enter":
      case " ":
        if (focusedIndex.value >= 0 && items[focusedIndex.value]) {
          e.preventDefault();
          items[focusedIndex.value].click();
        }
        break;
      case "Tab":
        closeDropdown();
        break;
    }
  }
  function setupListeners() {
    document.addEventListener("click", handleOutsideClick);
    document.addEventListener("keydown", handleKeydown);
  }
  function cleanupListeners() {
    document.removeEventListener("click", handleOutsideClick);
    document.removeEventListener("keydown", handleKeydown);
  }
  return {
    openDropdown,
    dropdownStyle,
    dropdownRefs,
    focusedIndex,
    toggleDropdown,
    closeDropdown,
    setupListeners,
    cleanupListeners
  };
}
const STORAGE_KEY$1 = "quick-filters-state";
const languageConfig = {
  "Go": { icon: "🐹" },
  "TypeScript": { icon: "📘" },
  "JavaScript": { icon: "📒" },
  "Vue": { icon: "💚" },
  "Python": { icon: "🐍" },
  "Java": { icon: "☕" },
  "Kotlin": { icon: "🟣" },
  "Rust": { icon: "🦀" },
  "C#": { icon: "🟦" },
  "C++": { icon: "⚡" },
  "C": { icon: "🔧" },
  "Ruby": { icon: "💎" },
  "PHP": { icon: "🐘" },
  "Swift": { icon: "🍎" },
  "Dart": { icon: "🎯" }
};
const languageExtensions = {
  "Go": [".go"],
  "TypeScript": [".ts", ".tsx"],
  "JavaScript": [".js", ".jsx"],
  "Vue": [".vue"],
  "Python": [".py"],
  "Java": [".java"],
  "Kotlin": [".kt", ".kts"],
  "Rust": [".rs"],
  "C#": [".cs"],
  "C++": [".cpp", ".cc", ".cxx", ".hpp", ".h"],
  "C": [".c", ".h"],
  "Ruby": [".rb"],
  "PHP": [".php"],
  "Swift": [".swift"],
  "Dart": [".dart"]
};
const createSmartFilter = (id, labelKey, icon, extensions, patterns, framework) => ({
  id,
  labelKey,
  icon,
  extensions,
  patterns,
  framework,
  category: "smart"
});
const frameworkFiltersConfig = {
  "Vue.js": [
    createSmartFilter("vue-components", "components", "🧩", [".vue"], ["**/components/**"], "Vue.js"),
    createSmartFilter("vue-composables", "composables", "🪝", [".ts"], ["**/composables/**", "**/use*.ts"], "Vue.js"),
    createSmartFilter("vue-stores", "stores", "🗄️", [".ts"], ["**/stores/**", "**/*.store.ts"], "Vue.js"),
    createSmartFilter("vue-views", "views", "📄", [".vue"], ["**/views/**", "**/pages/**"], "Vue.js")
  ],
  "React": [
    createSmartFilter("react-components", "components", "🧩", [".tsx", ".jsx"], ["**/components/**"], "React"),
    createSmartFilter("react-hooks", "hooks", "🪝", [".ts", ".tsx"], ["**/hooks/**", "**/use*.ts", "**/use*.tsx"], "React"),
    createSmartFilter("react-pages", "pages", "📄", [".tsx", ".jsx"], ["**/pages/**", "**/app/**"], "React")
  ],
  "Next.js": [
    createSmartFilter("next-pages", "pages", "📄", [".tsx", ".jsx"], ["**/app/**", "**/pages/**"], "Next.js"),
    createSmartFilter("next-components", "components", "🧩", [".tsx", ".jsx"], ["**/components/**"], "Next.js"),
    createSmartFilter("next-api", "api", "🔌", [".ts", ".tsx"], ["**/api/**", "**/route.ts"], "Next.js")
  ],
  "Angular": [
    createSmartFilter("angular-components", "components", "🧩", [".ts"], ["**/*.component.ts"], "Angular"),
    createSmartFilter("angular-services", "services", "⚙️", [".ts"], ["**/*.service.ts"], "Angular"),
    createSmartFilter("angular-modules", "modules", "📦", [".ts"], ["**/*.module.ts"], "Angular")
  ],
  "Gin": [
    createSmartFilter("go-handlers", "handlers", "🎯", [".go"], ["**/handlers/**", "**/*_handler.go"], "Gin"),
    createSmartFilter("go-services", "services", "⚙️", [".go"], ["**/services/**", "**/*_service.go", "**/application/**"], "Gin"),
    createSmartFilter("go-domain", "domain", "🏛️", [".go"], ["**/domain/**", "**/entities/**", "**/models/**"], "Gin"),
    createSmartFilter("go-infra", "infra", "🔧", [".go"], ["**/infrastructure/**", "**/repository/**", "**/adapters/**"], "Gin")
  ],
  "Echo": [
    createSmartFilter("go-handlers", "handlers", "🎯", [".go"], ["**/handlers/**", "**/*_handler.go"], "Echo"),
    createSmartFilter("go-services", "services", "⚙️", [".go"], ["**/services/**", "**/*_service.go"], "Echo")
  ],
  "Wails": [
    createSmartFilter("wails-backend", "backend", "🐹", [".go"], ["**/backend/**", "**/*.go"], "Wails"),
    createSmartFilter("wails-frontend", "frontend", "🎨", [".vue", ".tsx", ".ts"], ["**/frontend/**"], "Wails")
  ],
  "Django": [
    createSmartFilter("django-views", "views", "👁️", [".py"], ["**/views.py", "**/views/**"], "Django"),
    createSmartFilter("django-models", "models", "🗃️", [".py"], ["**/models.py", "**/models/**"], "Django"),
    createSmartFilter("django-urls", "urls", "🔗", [".py"], ["**/urls.py"], "Django")
  ],
  "FastAPI": [
    createSmartFilter("fastapi-routes", "routes", "🔌", [".py"], ["**/routes/**", "**/routers/**", "**/api/**"], "FastAPI"),
    createSmartFilter("fastapi-models", "models", "🗃️", [".py"], ["**/models/**", "**/schemas/**"], "FastAPI")
  ],
  "Spring Boot": [
    createSmartFilter("spring-controllers", "controllers", "🎯", [".java", ".kt"], ["**/*Controller.java", "**/*Controller.kt"], "Spring"),
    createSmartFilter("spring-services", "services", "⚙️", [".java", ".kt"], ["**/*Service.java", "**/*Service.kt"], "Spring"),
    createSmartFilter("spring-repos", "repos", "🗄️", [".java", ".kt"], ["**/*Repository.java", "**/*Repository.kt"], "Spring")
  ],
  "Flutter": [
    createSmartFilter("flutter-screens", "screens", "📱", [".dart"], ["**/screens/**", "**/pages/**"], "Flutter"),
    createSmartFilter("flutter-widgets", "widgets", "🧩", [".dart"], ["**/widgets/**", "**/components/**"], "Flutter"),
    createSmartFilter("flutter-bloc", "state", "🔄", [".dart"], ["**/bloc/**", "**/cubit/**", "**/providers/**"], "Flutter")
  ]
};
function useFilterPersistence(filterState) {
  const projectStore = useProjectStore();
  function getProjectKey() {
    return projectStore.currentPath?.replace(/[\\/:]/g, "_") || "default";
  }
  function saveState() {
    try {
      const allStates = JSON.parse(localStorage.getItem(STORAGE_KEY$1) || "{}");
      allStates[getProjectKey()] = {
        active: Array.from(filterState.value.active),
        excluded: Array.from(filterState.value.excluded)
      };
      localStorage.setItem(STORAGE_KEY$1, JSON.stringify(allStates));
    } catch {
    }
  }
  function loadState() {
    try {
      const allStates = JSON.parse(localStorage.getItem(STORAGE_KEY$1) || "{}");
      const saved = allStates[getProjectKey()];
      if (saved) {
        if (Array.isArray(saved.active)) {
          filterState.value.active = new Set(saved.active);
        } else if (Array.isArray(saved)) {
          filterState.value.active = new Set(saved);
        }
        if (Array.isArray(saved.excluded)) {
          filterState.value.excluded = new Set(saved.excluded);
        }
      }
    } catch {
    }
  }
  function clearState() {
    filterState.value.active.clear();
    filterState.value.excluded.clear();
  }
  watch(
    () => [filterState.value.active.size, filterState.value.excluded.size],
    () => saveState(),
    { deep: true }
  );
  return {
    saveState,
    loadState,
    clearState,
    getProjectKey
  };
}
const MIN_LANGUAGE_PERCENTAGE = 3;
const MIN_LANGUAGE_FILE_COUNT = 2;
function useSmartFilters() {
  const { t } = useI18n();
  const projectStore = useProjectStore();
  const projectStructure = ref(null);
  const isLoading = ref(false);
  const languageFilters = computed(() => {
    if (!projectStructure.value?.languages?.length) return [];
    return projectStructure.value.languages.filter(
      (lang) => lang.percentage > MIN_LANGUAGE_PERCENTAGE || lang.fileCount > MIN_LANGUAGE_FILE_COUNT
    ).map((lang) => ({
      id: `lang-${lang.name.toLowerCase().replace(/[^a-z0-9]/g, "")}`,
      label: lang.name,
      language: lang.name,
      icon: languageConfig[lang.name]?.icon || "📄",
      extensions: languageExtensions[lang.name] || [],
      fileCount: lang.fileCount,
      percentage: lang.percentage,
      primary: lang.primary,
      category: "lang",
      shortLabel: lang.name
    })).filter((f) => f.extensions.length > 0);
  });
  const smartFilters = computed(() => {
    if (!projectStructure.value?.frameworks?.length) return [];
    const filters = [];
    for (const fw of projectStructure.value.frameworks) {
      const fwFilters = frameworkFiltersConfig[fw.name];
      if (fwFilters) {
        filters.push(
          ...fwFilters.map((f) => ({
            ...f,
            label: t(f.labelKey) || f.labelKey,
            shortLabel: t(f.labelKey) || f.labelKey
          }))
        );
      }
    }
    return filters;
  });
  const detectedFrameworks = computed(() => {
    return projectStructure.value?.frameworks?.map((f) => f.name) || [];
  });
  async function loadProjectStructure() {
    if (!projectStore.currentPath) return;
    isLoading.value = true;
    try {
      projectStructure.value = await GetProjectStructure(
        projectStore.currentPath
      );
    } catch {
      projectStructure.value = null;
    } finally {
      isLoading.value = false;
    }
  }
  watch(
    () => projectStore.currentPath,
    (path) => {
      if (path) {
        loadProjectStructure();
      } else {
        projectStructure.value = null;
      }
    },
    { immediate: true }
  );
  return {
    projectStructure,
    languageFilters,
    smartFilters,
    detectedFrameworks,
    isLoading,
    loadProjectStructure
  };
}
const CATEGORY_MAP = {
  source: "code",
  tests: "test",
  config: "config",
  docs: "docs",
  styles: "styles"
};
const LABEL_I18N_MAP = {
  source: "filterLabels.source",
  tests: "filterLabels.tests",
  config: "filterLabels.config",
  docs: "filterLabels.docs",
  styles: "filterLabels.styles"
};
function useQuickFilters() {
  const { t } = useI18n();
  const fileStore = useFileStore();
  const settingsStore = useSettingsStore();
  const projectStore = useProjectStore();
  const dropdown = useFilterDropdown();
  const smartFiltersComposable = useSmartFilters();
  const { languageFilters, smartFilters, isLoading, loadProjectStructure } = smartFiltersComposable;
  const filterState = ref({
    active: /* @__PURE__ */ new Set(),
    excluded: /* @__PURE__ */ new Set()
  });
  const persistence = useFilterPersistence(filterState);
  const showSettingsModal = ref(false);
  const typeFilters = computed(
    () => settingsStore.settings.fileExplorer.quickFilters.map((f) => ({
      ...f,
      category: CATEGORY_MAP[f.id] || "code",
      shortLabel: LABEL_I18N_MAP[f.id] ? t(LABEL_I18N_MAP[f.id]) : f.label
    }))
  );
  const editableFilters = computed(() => settingsStore.settings.fileExplorer.quickFilters);
  const activeTypeFilters = computed(
    () => typeFilters.value.filter((f) => filterState.value.active.has(f.id))
  );
  const activeLanguageFilters = computed(
    () => languageFilters.value.filter((f) => filterState.value.active.has(f.id))
  );
  const activeSmartFilters = computed(
    () => smartFilters.value.filter((f) => filterState.value.active.has(f.id))
  );
  const allActiveFilters = computed(() => [
    ...activeTypeFilters.value,
    ...activeLanguageFilters.value,
    ...activeSmartFilters.value
  ]);
  const excludedTypeFilters = computed(
    () => typeFilters.value.filter((f) => filterState.value.excluded.has(f.id))
  );
  const excludedLanguageFilters = computed(
    () => languageFilters.value.filter((f) => filterState.value.excluded.has(f.id))
  );
  const excludedSmartFilters = computed(
    () => smartFilters.value.filter((f) => filterState.value.excluded.has(f.id))
  );
  const allExcludedFilters = computed(() => [
    ...excludedTypeFilters.value,
    ...excludedLanguageFilters.value,
    ...excludedSmartFilters.value
  ]);
  const hasActiveTypeFilters = computed(
    () => activeTypeFilters.value.length > 0 || excludedTypeFilters.value.length > 0
  );
  const hasActiveLanguageFilters = computed(
    () => activeLanguageFilters.value.length > 0 || excludedLanguageFilters.value.length > 0
  );
  const hasActiveSmartFilters = computed(
    () => activeSmartFilters.value.length > 0 || excludedSmartFilters.value.length > 0
  );
  const totalFiles = computed(() => {
    let count = 0;
    const countFiles = (nodes) => {
      nodes.forEach((node) => {
        if (!node.isDir) count++;
        if (node.children) countFiles(node.children);
      });
    };
    countFiles(fileStore.nodes);
    return count;
  });
  function toggleFilter(filter, event) {
    const isMulti = event.ctrlKey || event.metaKey;
    const isExclude = event.shiftKey;
    if (isExclude) {
      if (filterState.value.excluded.has(filter.id)) {
        filterState.value.excluded.delete(filter.id);
      } else {
        filterState.value.active.delete(filter.id);
        filterState.value.excluded.add(filter.id);
      }
    } else if (isMulti) {
      filterState.value.excluded.delete(filter.id);
      if (filterState.value.active.has(filter.id)) {
        filterState.value.active.delete(filter.id);
      } else {
        filterState.value.active.add(filter.id);
      }
    } else {
      filterState.value.excluded.clear();
      if (filterState.value.active.has(filter.id) && filterState.value.active.size === 1) {
        filterState.value.active.clear();
      } else {
        filterState.value.active.clear();
        filterState.value.active.add(filter.id);
      }
    }
    applyFilters();
    if (!isMulti && !isExclude) dropdown.closeDropdown();
  }
  function removeFilter(filter) {
    filterState.value.active.delete(filter.id);
    filterState.value.excluded.delete(filter.id);
    applyFilters();
  }
  function isFilterActive(id) {
    return filterState.value.active.has(id);
  }
  function isFilterExcluded(id) {
    return filterState.value.excluded.has(id);
  }
  function clearTypeFilters() {
    typeFilters.value.forEach((f) => {
      filterState.value.active.delete(f.id);
      filterState.value.excluded.delete(f.id);
    });
    applyFilters();
  }
  function clearLanguageFilters() {
    languageFilters.value.forEach((f) => {
      filterState.value.active.delete(f.id);
      filterState.value.excluded.delete(f.id);
    });
    applyFilters();
  }
  function clearSmartFilters() {
    smartFilters.value.forEach((f) => {
      filterState.value.active.delete(f.id);
      filterState.value.excluded.delete(f.id);
    });
    applyFilters();
  }
  function clearAllFilters() {
    filterState.value.active.clear();
    filterState.value.excluded.clear();
    applyFilters();
  }
  function applyFilters() {
    const includeExts = [];
    const excludeExts = [];
    activeTypeFilters.value.forEach((f) => f.extensions && includeExts.push(...f.extensions));
    activeLanguageFilters.value.forEach((f) => f.extensions && includeExts.push(...f.extensions));
    activeSmartFilters.value.forEach((f) => f.extensions && includeExts.push(...f.extensions));
    excludedTypeFilters.value.forEach((f) => f.extensions && excludeExts.push(...f.extensions));
    excludedLanguageFilters.value.forEach((f) => f.extensions && excludeExts.push(...f.extensions));
    excludedSmartFilters.value.forEach((f) => f.extensions && excludeExts.push(...f.extensions));
    fileStore.setFilterExtensions([...new Set(includeExts)], [...new Set(excludeExts)]);
  }
  function getFilterCount(filter) {
    let count = 0;
    const countFiles = (nodes) => {
      nodes.forEach((node) => {
        if (!node.isDir && filter.extensions?.some((e) => node.name.endsWith(e))) count++;
        if (node.children) countFiles(node.children);
      });
    };
    countFiles(fileStore.nodes);
    return count;
  }
  function getFilterPercentage(filter) {
    return totalFiles.value === 0 ? 0 : Math.min(100, getFilterCount(filter) / totalFiles.value * 100);
  }
  function updateFilterExtensions(filter, value) {
    filter.extensions = value.split(",").map((s) => s.trim()).filter((s) => s);
  }
  function resetFilters() {
    settingsStore.resetToDefaults();
  }
  function setupEventListeners() {
    dropdown.setupListeners();
    if (projectStore.currentPath) {
      loadProjectStructure();
      setTimeout(() => persistence.loadState(), 100);
    }
  }
  function cleanupEventListeners() {
    dropdown.cleanupListeners();
  }
  watch(
    () => projectStore.currentPath,
    (path) => {
      if (path) {
        persistence.clearState();
        setTimeout(() => {
          persistence.loadState();
          applyFilters();
        }, 150);
      }
    }
  );
  return {
    // Dropdown state (from useFilterDropdown)
    openDropdown: dropdown.openDropdown,
    dropdownStyle: dropdown.dropdownStyle,
    dropdownRefs: dropdown.dropdownRefs,
    toggleDropdown: dropdown.toggleDropdown,
    // UI state
    showSettingsModal,
    isLoading,
    // Filters
    typeFilters,
    languageFilters,
    smartFilters,
    editableFilters,
    // Active filters
    activeTypeFilters,
    activeLanguageFilters,
    activeSmartFilters,
    allActiveFilters,
    // Excluded filters
    allExcludedFilters,
    // Has active flags
    hasActiveTypeFilters,
    hasActiveLanguageFilters,
    hasActiveSmartFilters,
    // Totals
    totalFiles,
    // Filter management
    toggleFilter,
    removeFilter,
    isFilterActive,
    isFilterExcluded,
    clearTypeFilters,
    clearLanguageFilters,
    clearSmartFilters,
    clearAllFilters,
    // Filter info
    getFilterCount,
    getFilterPercentage,
    // Settings
    updateFilterExtensions,
    resetFilters,
    // Lifecycle
    setupEventListeners,
    cleanupEventListeners
  };
}
const _hoisted_1$12 = {
  key: 0,
  class: "chip-icon"
};
const _hoisted_2$_ = {
  key: 1,
  class: "chip-icon"
};
const _sfc_main$14 = /* @__PURE__ */ defineComponent({
  __name: "FilterChip",
  props: {
    label: {},
    icon: {},
    category: {},
    excluded: { type: Boolean }
  },
  emits: ["remove"],
  setup(__props) {
    return (_ctx, _cache) => {
      return openBlock(), createElementBlock("button", {
        onClick: _cache[0] || (_cache[0] = ($event) => _ctx.$emit("remove")),
        class: normalizeClass([
          "filter-chip",
          __props.excluded ? "filter-chip-excluded" : `filter-chip-${__props.category}`
        ])
      }, [
        __props.excluded ? (openBlock(), createElementBlock("span", _hoisted_1$12, "⊘")) : __props.icon ? (openBlock(), createElementBlock("span", _hoisted_2$_, toDisplayString(__props.icon), 1)) : createCommentVNode("", true),
        createBaseVNode("span", {
          class: normalizeClass({ "chip-label-excluded": __props.excluded })
        }, toDisplayString(__props.label), 3),
        _cache[1] || (_cache[1] = createBaseVNode("svg", {
          class: "w-3 h-3 chip-close",
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
        ], -1))
      ], 2);
    };
  }
});
const FilterChip = /* @__PURE__ */ _export_sfc(_sfc_main$14, [["__scopeId", "data-v-448c8587"]]);
const _hoisted_1$11 = ["aria-selected"];
const _hoisted_2$Z = { class: "filter-item-left" };
const _hoisted_3$V = {
  key: 1,
  class: "filter-item-emoji"
};
const _hoisted_4$P = {
  key: 2,
  class: "filter-item-info"
};
const _hoisted_5$J = { class: "filter-item-subtitle" };
const _hoisted_6$G = {
  key: 0,
  class: "filter-item-primary"
};
const _hoisted_7$C = { class: "filter-item-right" };
const _hoisted_8$A = { class: "filter-item-count" };
const _hoisted_9$x = { class: "filter-item-bar" };
const _sfc_main$13 = /* @__PURE__ */ defineComponent({
  __name: "FilterDropdownItem",
  props: {
    label: {},
    icon: {},
    iconType: {},
    category: {},
    count: {},
    percentage: {},
    active: { type: Boolean },
    excluded: { type: Boolean },
    primary: { type: Boolean },
    subtitle: {},
    barClass: {}
  },
  emits: ["toggle"],
  setup(__props) {
    return (_ctx, _cache) => {
      return openBlock(), createElementBlock("button", {
        onClick: _cache[0] || (_cache[0] = ($event) => _ctx.$emit("toggle", $event)),
        class: normalizeClass([
          "filter-dropdown-item filter-item",
          { "filter-dropdown-item-active": __props.active },
          { "filter-dropdown-item-excluded": __props.excluded }
        ]),
        role: "option",
        "aria-selected": __props.active,
        tabindex: "-1"
      }, [
        createBaseVNode("div", _hoisted_2$Z, [
          __props.iconType === "category" ? (openBlock(), createElementBlock("div", {
            key: 0,
            class: normalizeClass(["filter-item-icon", `filter-item-icon-${__props.category}`])
          }, [
            renderSlot(_ctx.$slots, "icon", {}, void 0, true)
          ], 2)) : __props.icon ? (openBlock(), createElementBlock("span", _hoisted_3$V, toDisplayString(__props.icon), 1)) : createCommentVNode("", true),
          __props.subtitle ? (openBlock(), createElementBlock("div", _hoisted_4$P, [
            createBaseVNode("span", {
              class: normalizeClass(["filter-item-label", { "filter-item-label-excluded": __props.excluded }])
            }, toDisplayString(__props.label), 3),
            createBaseVNode("span", _hoisted_5$J, toDisplayString(__props.subtitle), 1)
          ])) : (openBlock(), createElementBlock(Fragment, { key: 3 }, [
            createBaseVNode("span", {
              class: normalizeClass(["filter-item-label", { "filter-item-label-excluded": __props.excluded }])
            }, toDisplayString(__props.label), 3),
            __props.primary ? (openBlock(), createElementBlock("span", _hoisted_6$G, "★")) : createCommentVNode("", true)
          ], 64))
        ]),
        createBaseVNode("div", _hoisted_7$C, [
          createBaseVNode("span", _hoisted_8$A, toDisplayString(__props.count), 1),
          createBaseVNode("div", _hoisted_9$x, [
            createBaseVNode("div", {
              class: normalizeClass(["filter-item-bar-fill", __props.barClass]),
              style: normalizeStyle({ width: __props.percentage + "%" })
            }, null, 6)
          ])
        ])
      ], 10, _hoisted_1$11);
    };
  }
});
const FilterDropdownItem = /* @__PURE__ */ _export_sfc(_sfc_main$13, [["__scopeId", "data-v-988f4b76"]]);
const _hoisted_1$10 = { class: "filter-dropdown-items" };
const _hoisted_2$Y = {
  key: 0,
  class: "filter-dropdown-hint"
};
const _hoisted_3$U = { class: "filter-dropdown-footer-text" };
const _sfc_main$12 = /* @__PURE__ */ defineComponent({
  __name: "FilterDropdownMenu",
  props: {
    isOpen: { type: Boolean },
    title: {},
    style: {},
    hasActiveFilters: { type: Boolean },
    hint: { type: Boolean },
    footer: {},
    menuClass: {},
    headerClass: {},
    footerClass: {}
  },
  emits: ["clear"],
  setup(__props) {
    const { t } = useI18n();
    return (_ctx, _cache) => {
      return openBlock(), createBlock(Teleport, { to: "body" }, [
        createVNode(Transition, { name: "dropdown" }, {
          default: withCtx(() => [
            __props.isOpen ? (openBlock(), createElementBlock("div", {
              key: 0,
              class: normalizeClass(["filter-dropdown-menu", __props.menuClass]),
              style: normalizeStyle(__props.style)
            }, [
              createBaseVNode("div", {
                class: normalizeClass(["filter-dropdown-header", __props.headerClass])
              }, [
                createBaseVNode("span", null, toDisplayString(__props.title), 1),
                __props.hasActiveFilters ? (openBlock(), createElementBlock("button", {
                  key: 0,
                  onClick: _cache[0] || (_cache[0] = ($event) => _ctx.$emit("clear")),
                  class: "filter-clear-btn"
                }, toDisplayString(unref(t)("quickFilters.clearAll")), 1)) : createCommentVNode("", true)
              ], 2),
              createBaseVNode("div", _hoisted_1$10, [
                renderSlot(_ctx.$slots, "default", {}, void 0, true)
              ]),
              __props.hint ? (openBlock(), createElementBlock("div", _hoisted_2$Y, [
                _cache[1] || (_cache[1] = createBaseVNode("kbd", null, "Ctrl", -1)),
                createTextVNode(" " + toDisplayString(unref(t)("quickFilters.multiSelect")) + " · ", 1),
                _cache[2] || (_cache[2] = createBaseVNode("kbd", null, "Shift", -1)),
                createTextVNode(" " + toDisplayString(unref(t)("quickFilters.exclude")), 1)
              ])) : createCommentVNode("", true),
              __props.footer ? (openBlock(), createElementBlock("div", {
                key: 1,
                class: normalizeClass(["filter-dropdown-footer", __props.footerClass])
              }, [
                createBaseVNode("span", _hoisted_3$U, toDisplayString(__props.footer), 1)
              ], 2)) : createCommentVNode("", true)
            ], 6)) : createCommentVNode("", true)
          ]),
          _: 3
        })
      ]);
    };
  }
});
const FilterDropdownMenu = /* @__PURE__ */ _export_sfc(_sfc_main$12, [["__scopeId", "data-v-b040d798"]]);
const _hoisted_1$$ = { class: "modal-content filter-settings-modal" };
const _hoisted_2$X = { class: "modal-header" };
const _hoisted_3$T = ["aria-label"];
const _hoisted_4$O = { class: "modal-body" };
const _hoisted_5$I = { class: "settings-section" };
const _hoisted_6$F = { class: "settings-filter-header" };
const _hoisted_7$B = { class: "settings-toggle" };
const _hoisted_8$z = ["onUpdate:modelValue"];
const _hoisted_9$w = { class: "settings-filter-name" };
const _hoisted_10$v = { class: "settings-filter-count" };
const _hoisted_11$t = { class: "settings-filter-inputs" };
const _hoisted_12$p = ["value", "onChange", "placeholder"];
const _hoisted_13$p = { class: "modal-footer" };
const _sfc_main$11 = /* @__PURE__ */ defineComponent({
  __name: "FilterSettingsModal",
  props: {
    isOpen: { type: Boolean },
    filters: {},
    getCount: { type: Function }
  },
  emits: ["close", "reset", "updateExtensions"],
  setup(__props) {
    const { t } = useI18n();
    return (_ctx, _cache) => {
      return openBlock(), createBlock(Teleport, { to: "body" }, [
        createVNode(Transition, { name: "modal" }, {
          default: withCtx(() => [
            __props.isOpen ? (openBlock(), createElementBlock("div", {
              key: 0,
              class: "modal-overlay",
              onClick: _cache[3] || (_cache[3] = withModifiers(($event) => _ctx.$emit("close"), ["self"]))
            }, [
              createBaseVNode("div", _hoisted_1$$, [
                createBaseVNode("div", _hoisted_2$X, [
                  createBaseVNode("h3", null, toDisplayString(unref(t)("quickFilters.settingsTitle")), 1),
                  createBaseVNode("button", {
                    onClick: _cache[0] || (_cache[0] = ($event) => _ctx.$emit("close")),
                    class: "icon-btn",
                    "aria-label": unref(t)("common.close")
                  }, [..._cache[4] || (_cache[4] = [
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
                  ])], 8, _hoisted_3$T)
                ]),
                createBaseVNode("div", _hoisted_4$O, [
                  createBaseVNode("div", _hoisted_5$I, [
                    createBaseVNode("h4", null, toDisplayString(unref(t)("quickFilters.customFilters")), 1),
                    (openBlock(true), createElementBlock(Fragment, null, renderList(__props.filters, (filter) => {
                      return openBlock(), createElementBlock("div", {
                        key: filter.id,
                        class: "settings-filter-card"
                      }, [
                        createBaseVNode("div", _hoisted_6$F, [
                          createBaseVNode("label", _hoisted_7$B, [
                            withDirectives(createBaseVNode("input", {
                              type: "checkbox",
                              "onUpdate:modelValue": ($event) => filter.enabled = $event
                            }, null, 8, _hoisted_8$z), [
                              [vModelCheckbox, filter.enabled]
                            ]),
                            _cache[5] || (_cache[5] = createBaseVNode("span", { class: "settings-toggle-track" }, null, -1))
                          ]),
                          createBaseVNode("span", _hoisted_9$w, toDisplayString(filter.label), 1),
                          createBaseVNode("span", _hoisted_10$v, toDisplayString(__props.getCount(filter)), 1)
                        ]),
                        createBaseVNode("div", _hoisted_11$t, [
                          createBaseVNode("input", {
                            type: "text",
                            value: filter.extensions?.join(", "),
                            onChange: ($event) => _ctx.$emit("updateExtensions", filter, $event.target.value),
                            class: "input input-sm",
                            placeholder: unref(t)("quickFilters.extensionsPlaceholder")
                          }, null, 40, _hoisted_12$p)
                        ])
                      ]);
                    }), 128))
                  ])
                ]),
                createBaseVNode("div", _hoisted_13$p, [
                  createBaseVNode("button", {
                    onClick: _cache[1] || (_cache[1] = ($event) => _ctx.$emit("reset")),
                    class: "btn btn-secondary btn-sm"
                  }, toDisplayString(unref(t)("quickFilters.reset")), 1),
                  createBaseVNode("button", {
                    onClick: _cache[2] || (_cache[2] = ($event) => _ctx.$emit("close")),
                    class: "btn btn-primary btn-sm"
                  }, toDisplayString(unref(t)("quickFilters.done")), 1)
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
const FilterSettingsModal = /* @__PURE__ */ _export_sfc(_sfc_main$11, [["__scopeId", "data-v-566e2a3d"]]);
const _hoisted_1$_ = { class: "filters-bar" };
const _hoisted_2$W = { class: "filter-groups" };
const _hoisted_3$S = {
  class: "filter-dropdown",
  ref: "typesDropdownRef"
};
const _hoisted_4$N = ["aria-label", "aria-expanded"];
const _hoisted_5$H = {
  key: 0,
  class: "filter-badge"
};
const _hoisted_6$E = {
  key: 0,
  class: "filter-dropdown",
  ref: "langsDropdownRef"
};
const _hoisted_7$A = ["aria-label", "aria-expanded"];
const _hoisted_8$y = { class: "filter-badge" };
const _hoisted_9$v = {
  key: 1,
  class: "filter-dropdown",
  ref: "smartDropdownRef"
};
const _hoisted_10$u = ["aria-label", "aria-expanded"];
const _hoisted_11$s = { class: "filter-badge filter-badge-smart" };
const _hoisted_12$o = {
  key: 0,
  class: "active-filters"
};
const _hoisted_13$o = { class: "filter-actions" };
const _hoisted_14$m = ["title"];
const _hoisted_15$k = ["title"];
const _sfc_main$10 = /* @__PURE__ */ defineComponent({
  __name: "QuickFiltersBar",
  setup(__props) {
    const { t } = useI18n();
    const {
      openDropdown,
      dropdownStyle,
      showSettingsModal,
      typeFilters,
      languageFilters,
      smartFilters,
      editableFilters,
      activeTypeFilters,
      activeLanguageFilters,
      activeSmartFilters,
      allActiveFilters,
      allExcludedFilters,
      hasActiveTypeFilters,
      hasActiveLanguageFilters,
      hasActiveSmartFilters,
      toggleDropdown,
      toggleFilter,
      removeFilter,
      isFilterActive,
      isFilterExcluded,
      clearTypeFilters,
      clearLanguageFilters,
      clearSmartFilters,
      clearAllFilters,
      getFilterCount,
      getFilterPercentage,
      updateFilterExtensions,
      resetFilters,
      setupEventListeners,
      cleanupEventListeners
    } = useQuickFilters();
    const ChevronIcon = (props) => h("svg", {
      class: ["w-3 h-3 transition-transform", { "rotate-180": props.open }],
      fill: "none",
      stroke: "currentColor",
      viewBox: "0 0 24 24"
    }, [h("path", { "stroke-linecap": "round", "stroke-linejoin": "round", "stroke-width": "2", d: "M19 9l-7 7-7-7" })]);
    const iconPaths = {
      code: "M10 20l4-16m4 4l4 4-4 4M6 16l-4-4 4-4",
      test: "M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z",
      config: "M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z",
      docs: "M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z",
      styles: "M7 21a4 4 0 01-4-4V5a2 2 0 012-2h4a2 2 0 012 2v12a4 4 0 01-4 4zm0 0h12a2 2 0 002-2v-4a2 2 0 00-2-2h-2.343M11 7.343l1.657-1.657a2 2 0 012.828 0l2.829 2.829a2 2 0 010 2.828l-8.486 8.485M7 17h.01"
    };
    function getFilterIcon(category) {
      return () => h("svg", { fill: "none", stroke: "currentColor", viewBox: "0 0 24 24" }, [
        h("path", { "stroke-linecap": "round", "stroke-linejoin": "round", "stroke-width": "2", d: iconPaths[category] || iconPaths.code })
      ]);
    }
    onMounted(() => setupEventListeners());
    onUnmounted(() => cleanupEventListeners());
    return (_ctx, _cache) => {
      return openBlock(), createElementBlock(Fragment, null, [
        createBaseVNode("div", _hoisted_1$_, [
          createBaseVNode("div", _hoisted_2$W, [
            createBaseVNode("div", _hoisted_3$S, [
              createBaseVNode("button", {
                onClick: _cache[0] || (_cache[0] = ($event) => unref(toggleDropdown)("types")),
                class: normalizeClass(["filter-trigger", { "filter-trigger-active": unref(hasActiveTypeFilters) }]),
                "aria-label": unref(t)("quickFilters.filterByType"),
                "aria-expanded": unref(openDropdown) === "types",
                "aria-haspopup": "listbox"
              }, [
                _cache[6] || (_cache[6] = createBaseVNode("svg", {
                  class: "w-3.5 h-3.5",
                  fill: "none",
                  stroke: "currentColor",
                  viewBox: "0 0 24 24"
                }, [
                  createBaseVNode("path", {
                    "stroke-linecap": "round",
                    "stroke-linejoin": "round",
                    "stroke-width": "2",
                    d: "M3 4a1 1 0 011-1h16a1 1 0 011 1v2.586a1 1 0 01-.293.707l-6.414 6.414a1 1 0 00-.293.707V17l-4 4v-6.586a1 1 0 00-.293-.707L3.293 7.293A1 1 0 013 6.586V4z"
                  })
                ], -1)),
                createBaseVNode("span", null, toDisplayString(unref(t)("quickFilters.types")), 1),
                unref(activeTypeFilters).length > 0 ? (openBlock(), createElementBlock("span", _hoisted_5$H, toDisplayString(unref(activeTypeFilters).length), 1)) : createCommentVNode("", true),
                createVNode(ChevronIcon, {
                  open: unref(openDropdown) === "types"
                }, null, 8, ["open"])
              ], 10, _hoisted_4$N),
              createVNode(FilterDropdownMenu, {
                "is-open": unref(openDropdown) === "types",
                title: unref(t)("quickFilters.fileTypes"),
                style: normalizeStyle(unref(dropdownStyle)),
                "has-active-filters": unref(activeTypeFilters).length > 0,
                hint: true,
                onClear: unref(clearTypeFilters)
              }, {
                default: withCtx(() => [
                  (openBlock(true), createElementBlock(Fragment, null, renderList(unref(typeFilters), (filter) => {
                    return openBlock(), createBlock(FilterDropdownItem, {
                      key: filter.id,
                      label: filter.label,
                      "icon-type": "category",
                      category: filter.category,
                      count: unref(getFilterCount)(filter),
                      percentage: unref(getFilterPercentage)(filter),
                      active: unref(isFilterActive)(filter.id),
                      excluded: unref(isFilterExcluded)(filter.id),
                      onToggle: ($event) => unref(toggleFilter)(filter, $event)
                    }, {
                      icon: withCtx(() => [
                        (openBlock(), createBlock(resolveDynamicComponent(getFilterIcon(filter.category)), { class: "w-3.5 h-3.5" }))
                      ]),
                      _: 2
                    }, 1032, ["label", "category", "count", "percentage", "active", "excluded", "onToggle"]);
                  }), 128))
                ]),
                _: 1
              }, 8, ["is-open", "title", "style", "has-active-filters", "onClear"])
            ], 512),
            unref(languageFilters).length > 0 ? (openBlock(), createElementBlock("div", _hoisted_6$E, [
              createBaseVNode("button", {
                onClick: _cache[1] || (_cache[1] = ($event) => unref(toggleDropdown)("langs")),
                class: normalizeClass(["filter-trigger", { "filter-trigger-active": unref(hasActiveLanguageFilters) }]),
                "aria-label": unref(t)("quickFilters.filterByLanguage"),
                "aria-expanded": unref(openDropdown) === "langs",
                "aria-haspopup": "listbox"
              }, [
                _cache[7] || (_cache[7] = createBaseVNode("svg", {
                  class: "w-3.5 h-3.5",
                  fill: "none",
                  stroke: "currentColor",
                  viewBox: "0 0 24 24"
                }, [
                  createBaseVNode("path", {
                    "stroke-linecap": "round",
                    "stroke-linejoin": "round",
                    "stroke-width": "2",
                    d: "M10 20l4-16m4 4l4 4-4 4M6 16l-4-4 4-4"
                  })
                ], -1)),
                createBaseVNode("span", null, toDisplayString(unref(t)("quickFilters.languages")), 1),
                createBaseVNode("span", _hoisted_8$y, toDisplayString(unref(languageFilters).length), 1),
                createVNode(ChevronIcon, {
                  open: unref(openDropdown) === "langs"
                }, null, 8, ["open"])
              ], 10, _hoisted_7$A),
              createVNode(FilterDropdownMenu, {
                "is-open": unref(openDropdown) === "langs",
                title: unref(t)("quickFilters.projectLanguages"),
                style: normalizeStyle(unref(dropdownStyle)),
                "has-active-filters": unref(activeLanguageFilters).length > 0,
                footer: unref(t)("quickFilters.autoDetected") + " ✨",
                onClear: unref(clearLanguageFilters)
              }, {
                default: withCtx(() => [
                  (openBlock(true), createElementBlock(Fragment, null, renderList(unref(languageFilters), (filter) => {
                    return openBlock(), createBlock(FilterDropdownItem, {
                      key: filter.id,
                      label: filter.language,
                      icon: filter.icon,
                      count: filter.fileCount,
                      percentage: filter.percentage,
                      active: unref(isFilterActive)(filter.id),
                      primary: filter.primary,
                      "bar-class": "bar-lang",
                      onToggle: ($event) => unref(toggleFilter)(filter, $event)
                    }, null, 8, ["label", "icon", "count", "percentage", "active", "primary", "onToggle"]);
                  }), 128))
                ]),
                _: 1
              }, 8, ["is-open", "title", "style", "has-active-filters", "footer", "onClear"])
            ], 512)) : createCommentVNode("", true),
            unref(smartFilters).length > 0 ? (openBlock(), createElementBlock("div", _hoisted_9$v, [
              createBaseVNode("button", {
                onClick: _cache[2] || (_cache[2] = ($event) => unref(toggleDropdown)("smart")),
                class: normalizeClass(["filter-trigger filter-trigger-smart", { "filter-trigger-active": unref(hasActiveSmartFilters) }]),
                "aria-label": unref(t)("quickFilters.smartFilters"),
                "aria-expanded": unref(openDropdown) === "smart",
                "aria-haspopup": "listbox"
              }, [
                _cache[8] || (_cache[8] = createBaseVNode("svg", {
                  class: "w-3.5 h-3.5",
                  fill: "none",
                  stroke: "currentColor",
                  viewBox: "0 0 24 24"
                }, [
                  createBaseVNode("path", {
                    "stroke-linecap": "round",
                    "stroke-linejoin": "round",
                    "stroke-width": "2",
                    d: "M9.663 17h4.673M12 3v1m6.364 1.636l-.707.707M21 12h-1M4 12H3m3.343-5.657l-.707-.707m2.828 9.9a5 5 0 117.072 0l-.548.547A3.374 3.374 0 0014 18.469V19a2 2 0 11-4 0v-.531c0-.895-.356-1.754-.988-2.386l-.548-.547z"
                  })
                ], -1)),
                createBaseVNode("span", null, toDisplayString(unref(t)("quickFilters.smart")), 1),
                createBaseVNode("span", _hoisted_11$s, toDisplayString(unref(smartFilters).length), 1),
                createVNode(ChevronIcon, {
                  open: unref(openDropdown) === "smart"
                }, null, 8, ["open"])
              ], 10, _hoisted_10$u),
              createVNode(FilterDropdownMenu, {
                "is-open": unref(openDropdown) === "smart",
                title: unref(t)("quickFilters.smartFilters"),
                style: normalizeStyle(unref(dropdownStyle)),
                "has-active-filters": unref(activeSmartFilters).length > 0,
                footer: unref(t)("quickFilters.basedOnFramework") + " 🧠",
                "menu-class": "filter-menu-smart",
                "header-class": "filter-header-smart",
                "footer-class": "filter-footer-smart",
                onClear: unref(clearSmartFilters)
              }, {
                default: withCtx(() => [
                  (openBlock(true), createElementBlock(Fragment, null, renderList(unref(smartFilters), (filter) => {
                    return openBlock(), createBlock(FilterDropdownItem, {
                      key: filter.id,
                      label: filter.label,
                      icon: filter.icon,
                      subtitle: filter.framework,
                      count: unref(getFilterCount)(filter),
                      percentage: unref(getFilterPercentage)(filter),
                      active: unref(isFilterActive)(filter.id),
                      "bar-class": "bar-smart",
                      onToggle: ($event) => unref(toggleFilter)(filter, $event)
                    }, null, 8, ["label", "icon", "subtitle", "count", "percentage", "active", "onToggle"]);
                  }), 128))
                ]),
                _: 1
              }, 8, ["is-open", "title", "style", "has-active-filters", "footer", "onClear"])
            ], 512)) : createCommentVNode("", true)
          ]),
          unref(allActiveFilters).length > 0 || unref(allExcludedFilters).length > 0 ? (openBlock(), createElementBlock("div", _hoisted_12$o, [
            createVNode(TransitionGroup, { name: "chip" }, {
              default: withCtx(() => [
                (openBlock(true), createElementBlock(Fragment, null, renderList(unref(allActiveFilters), (filter) => {
                  return openBlock(), createBlock(FilterChip, {
                    key: filter.id,
                    label: filter.shortLabel || filter.label,
                    icon: filter.icon,
                    category: filter.category,
                    onRemove: ($event) => unref(removeFilter)(filter)
                  }, null, 8, ["label", "icon", "category", "onRemove"]);
                }), 128)),
                (openBlock(true), createElementBlock(Fragment, null, renderList(unref(allExcludedFilters), (filter) => {
                  return openBlock(), createBlock(FilterChip, {
                    key: "ex-" + filter.id,
                    label: filter.shortLabel || filter.label,
                    excluded: true,
                    onRemove: ($event) => unref(removeFilter)(filter)
                  }, null, 8, ["label", "onRemove"]);
                }), 128))
              ]),
              _: 1
            })
          ])) : createCommentVNode("", true),
          _cache[11] || (_cache[11] = createBaseVNode("div", { class: "flex-1" }, null, -1)),
          createBaseVNode("div", _hoisted_13$o, [
            unref(allActiveFilters).length > 0 ? (openBlock(), createElementBlock("button", {
              key: 0,
              onClick: _cache[3] || (_cache[3] = //@ts-ignore
              (...args) => unref(clearAllFilters) && unref(clearAllFilters)(...args)),
              class: "filter-clear-btn",
              title: unref(t)("quickFilters.clearAll")
            }, [
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
                  d: "M6 18L18 6M6 6l12 12"
                })
              ], -1)),
              createBaseVNode("span", null, toDisplayString(unref(t)("quickFilters.clear")), 1)
            ], 8, _hoisted_14$m)) : createCommentVNode("", true),
            createBaseVNode("button", {
              onClick: _cache[4] || (_cache[4] = ($event) => showSettingsModal.value = true),
              class: "filter-settings-btn",
              title: unref(t)("quickFilters.settings")
            }, [..._cache[10] || (_cache[10] = [
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
                  d: "M12 6V4m0 2a2 2 0 100 4m0-4a2 2 0 110 4m-6 8a2 2 0 100-4m0 4a2 2 0 110-4m0 4v2m0-6V4m6 6v10m6-2a2 2 0 100-4m0 4a2 2 0 110-4m0 4v2m0-6V4"
                })
              ], -1)
            ])], 8, _hoisted_15$k)
          ])
        ]),
        createVNode(FilterSettingsModal, {
          "is-open": unref(showSettingsModal),
          filters: unref(editableFilters),
          "get-count": unref(getFilterCount),
          onClose: _cache[5] || (_cache[5] = ($event) => showSettingsModal.value = false),
          onReset: unref(resetFilters),
          onUpdateExtensions: unref(updateFilterExtensions)
        }, null, 8, ["is-open", "filters", "get-count", "onReset", "onUpdateExtensions"])
      ], 64);
    };
  }
});
const QuickFiltersBar = /* @__PURE__ */ _export_sfc(_sfc_main$10, [["__scopeId", "data-v-dfe44653"]]);
const _hoisted_1$Z = ["aria-checked", "disabled"];
const _sfc_main$$ = /* @__PURE__ */ defineComponent({
  __name: "ToggleSwitch",
  props: {
    modelValue: { type: Boolean },
    disabled: { type: Boolean, default: false }
  },
  emits: ["update:modelValue", "change"],
  setup(__props, { emit: __emit }) {
    const props = __props;
    const emit = __emit;
    function toggle() {
      if (props.disabled) return;
      const newValue = !props.modelValue;
      emit("update:modelValue", newValue);
      emit("change", newValue);
    }
    return (_ctx, _cache) => {
      return openBlock(), createElementBlock("button", {
        type: "button",
        role: "switch",
        "aria-checked": __props.modelValue,
        class: normalizeClass([
          "toggle-switch",
          __props.modelValue ? "toggle-switch--active" : "",
          __props.disabled ? "toggle-switch--disabled" : ""
        ]),
        disabled: __props.disabled,
        onClick: toggle
      }, [..._cache[0] || (_cache[0] = [
        createBaseVNode("span", { class: "toggle-switch__track" }, [
          createBaseVNode("span", { class: "toggle-switch__thumb" })
        ], -1)
      ])], 10, _hoisted_1$Z);
    };
  }
});
const ToggleSwitch = /* @__PURE__ */ _export_sfc(_sfc_main$$, [["__scopeId", "data-v-4e554ee8"]]);
const _hoisted_1$Y = { class: "settings-popover-wrapper" };
const _hoisted_2$V = ["title"];
const _hoisted_3$R = { class: "settings-popover__header" };
const _hoisted_4$M = { class: "settings-popover__content" };
const _hoisted_5$G = { class: "settings-popover__toggles" };
const _hoisted_6$D = { class: "settings-toggle__label" };
const _hoisted_7$z = { class: "settings-toggle__label" };
const _hoisted_8$x = { class: "settings-toggle__label" };
const _hoisted_9$u = { class: "settings-toggle__label" };
const _hoisted_10$t = { class: "settings-toggle__label" };
const _hoisted_11$r = { class: "settings-popover__footer" };
const _sfc_main$_ = /* @__PURE__ */ defineComponent({
  __name: "SettingsPopover",
  emits: ["open-ignore-rules", "settings-changed"],
  setup(__props, { emit: __emit }) {
    const { t } = useI18n();
    const settingsStore = useSettingsStore();
    const emit = __emit;
    const isOpen = ref(false);
    const triggerRef2 = ref(null);
    const popoverRef = ref(null);
    const popoverPosition = ref({ top: 0, left: 0 });
    const settings2 = computed(() => settingsStore.settings);
    const popoverStyle = computed(() => ({
      top: `${popoverPosition.value.top}px`,
      left: `${popoverPosition.value.left}px`
    }));
    onClickOutside(popoverRef, (event) => {
      if (triggerRef2.value?.contains(event.target)) return;
      isOpen.value = false;
    });
    function updatePosition() {
      if (!triggerRef2.value) return;
      const rect = triggerRef2.value.getBoundingClientRect();
      const popoverWidth = 256;
      let left = rect.right - popoverWidth;
      const top = rect.bottom + 8;
      if (left < 8) left = 8;
      popoverPosition.value = { top, left };
    }
    watch(isOpen, async (open) => {
      if (open) {
        await nextTick();
        updatePosition();
      }
    });
    function toggle() {
      isOpen.value = !isOpen.value;
    }
    function toggleSetting(key) {
      const current = settings2.value.fileExplorer[key];
      if (typeof current === "boolean") {
        updateSetting(key, !current);
      }
    }
    function updateSetting(key, value) {
      settingsStore.updateFileExplorerSettings({ [key]: value });
      emit("settings-changed");
    }
    function openIgnoreRules() {
      isOpen.value = false;
      emit("open-ignore-rules");
    }
    return (_ctx, _cache) => {
      return openBlock(), createElementBlock("div", _hoisted_1$Y, [
        createBaseVNode("button", {
          ref_key: "triggerRef",
          ref: triggerRef2,
          onClick: toggle,
          class: normalizeClass(["settings-trigger", { "settings-trigger--active": isOpen.value }]),
          title: unref(t)("files.settings")
        }, [..._cache[10] || (_cache[10] = [
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
              d: "M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z"
            }),
            createBaseVNode("path", {
              "stroke-linecap": "round",
              "stroke-linejoin": "round",
              "stroke-width": "2",
              d: "M15 12a3 3 0 11-6 0 3 3 0 016 0z"
            })
          ], -1)
        ])], 10, _hoisted_2$V),
        (openBlock(), createBlock(Teleport, { to: "body" }, [
          createVNode(Transition, { name: "popover" }, {
            default: withCtx(() => [
              isOpen.value ? (openBlock(), createElementBlock("div", {
                key: 0,
                ref_key: "popoverRef",
                ref: popoverRef,
                class: "settings-popover",
                style: normalizeStyle(popoverStyle.value)
              }, [
                createBaseVNode("div", _hoisted_3$R, [
                  createBaseVNode("span", null, toDisplayString(unref(t)("settings.title")), 1)
                ]),
                createBaseVNode("div", _hoisted_4$M, [
                  createBaseVNode("div", _hoisted_5$G, [
                    createBaseVNode("div", {
                      class: "settings-toggle",
                      onClick: _cache[1] || (_cache[1] = ($event) => toggleSetting("useGitignore"))
                    }, [
                      createBaseVNode("span", _hoisted_6$D, toDisplayString(unref(t)("settings.useGitignore")), 1),
                      createVNode(ToggleSwitch, {
                        "model-value": settings2.value.fileExplorer.useGitignore,
                        "onUpdate:modelValue": _cache[0] || (_cache[0] = (v) => updateSetting("useGitignore", v))
                      }, null, 8, ["model-value"])
                    ]),
                    createBaseVNode("div", {
                      class: "settings-toggle",
                      onClick: _cache[3] || (_cache[3] = ($event) => toggleSetting("useCustomIgnore"))
                    }, [
                      createBaseVNode("span", _hoisted_7$z, toDisplayString(unref(t)("settings.useCustomIgnore")), 1),
                      createVNode(ToggleSwitch, {
                        "model-value": settings2.value.fileExplorer.useCustomIgnore,
                        "onUpdate:modelValue": _cache[2] || (_cache[2] = (v) => updateSetting("useCustomIgnore", v))
                      }, null, 8, ["model-value"])
                    ]),
                    createBaseVNode("div", {
                      class: "settings-toggle",
                      onClick: _cache[5] || (_cache[5] = ($event) => toggleSetting("compactNestedFolders"))
                    }, [
                      createBaseVNode("span", _hoisted_8$x, toDisplayString(unref(t)("settings.compactFolders")), 1),
                      createVNode(ToggleSwitch, {
                        "model-value": settings2.value.fileExplorer.compactNestedFolders,
                        "onUpdate:modelValue": _cache[4] || (_cache[4] = (v) => updateSetting("compactNestedFolders", v))
                      }, null, 8, ["model-value"])
                    ]),
                    createBaseVNode("div", {
                      class: "settings-toggle",
                      onClick: _cache[7] || (_cache[7] = ($event) => toggleSetting("foldersFirst"))
                    }, [
                      createBaseVNode("span", _hoisted_9$u, toDisplayString(unref(t)("settings.foldersFirst")), 1),
                      createVNode(ToggleSwitch, {
                        "model-value": settings2.value.fileExplorer.foldersFirst,
                        "onUpdate:modelValue": _cache[6] || (_cache[6] = (v) => updateSetting("foldersFirst", v))
                      }, null, 8, ["model-value"])
                    ]),
                    createBaseVNode("div", {
                      class: "settings-toggle",
                      onClick: _cache[9] || (_cache[9] = ($event) => toggleSetting("allowSelectBinary"))
                    }, [
                      createBaseVNode("span", _hoisted_10$t, toDisplayString(unref(t)("settings.allowSelectBinary")), 1),
                      createVNode(ToggleSwitch, {
                        "model-value": settings2.value.fileExplorer.allowSelectBinary,
                        "onUpdate:modelValue": _cache[8] || (_cache[8] = (v) => updateSetting("allowSelectBinary", v))
                      }, null, 8, ["model-value"])
                    ])
                  ])
                ]),
                createBaseVNode("div", _hoisted_11$r, [
                  createBaseVNode("button", {
                    onClick: openIgnoreRules,
                    class: "settings-popover__manage-btn"
                  }, [
                    _cache[11] || (_cache[11] = createBaseVNode("svg", {
                      class: "w-4 h-4",
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
                    ], -1)),
                    createTextVNode(" " + toDisplayString(unref(t)("settings.manageRules")), 1)
                  ])
                ])
              ], 4)) : createCommentVNode("", true)
            ]),
            _: 1
          })
        ]))
      ]);
    };
  }
});
const SettingsPopover = /* @__PURE__ */ _export_sfc(_sfc_main$_, [["__scopeId", "data-v-9d178589"]]);
function useVirtualTree(options2) {
  const { nodes } = options2;
  const flattenedVisibleNodes = computed(() => {
    const result = [];
    function flatten(nodeList, depth, ancestorHasMoreSiblings) {
      nodeList.forEach((node, index) => {
        const isLast = index === nodeList.length - 1;
        result.push({
          id: node.path,
          // Flat key for RecycleScroller
          node,
          depth,
          isLast,
          ancestorHasMoreSiblings: [...ancestorHasMoreSiblings]
        });
        if (node.isDir && node.isExpanded && node.children?.length) {
          flatten(
            node.children,
            depth + 1,
            [...ancestorHasMoreSiblings, !isLast]
          );
        }
      });
    }
    flatten(nodes.value, 0, []);
    return result;
  });
  const totalVisibleCount = computed(() => flattenedVisibleNodes.value.length);
  return {
    flattenedVisibleNodes,
    totalVisibleCount
  };
}
const FILE_TYPE_CONFIG = {
  // TypeScript/JavaScript
  ts: { icon: "🔷", colorClass: "bg-blue-500" },
  tsx: { icon: "⚛️", colorClass: "bg-cyan-500" },
  js: { icon: "🟨", colorClass: "bg-yellow-500" },
  jsx: { icon: "⚛️", colorClass: "bg-cyan-500" },
  mjs: { icon: "🟨", colorClass: "bg-yellow-500" },
  cjs: { icon: "🟨", colorClass: "bg-yellow-500" },
  // Frameworks
  vue: { icon: "💚", colorClass: "bg-emerald-500" },
  svelte: { icon: "🔥", colorClass: "bg-orange-500" },
  // Backend
  go: { icon: "🐹", colorClass: "bg-sky-500" },
  py: { icon: "🐍", colorClass: "bg-yellow-600" },
  java: { icon: "☕", colorClass: "bg-orange-600" },
  kt: { icon: "🟣", colorClass: "bg-purple-500" },
  kts: { icon: "🟣", colorClass: "bg-purple-500" },
  rs: { icon: "🦀", colorClass: "bg-orange-500" },
  rb: { icon: "💎", colorClass: "bg-red-500" },
  php: { icon: "🐘", colorClass: "bg-indigo-400" },
  cs: { icon: "🟦", colorClass: "bg-purple-600" },
  // C/C++
  cpp: { icon: "⚙️", colorClass: "bg-blue-600" },
  cc: { icon: "⚙️", colorClass: "bg-blue-600" },
  cxx: { icon: "⚙️", colorClass: "bg-blue-600" },
  c: { icon: "⚙️", colorClass: "bg-blue-500" },
  h: { icon: "📎", colorClass: "bg-blue-400" },
  hpp: { icon: "📎", colorClass: "bg-blue-400" },
  // Styles
  css: { icon: "🎨", colorClass: "bg-pink-500" },
  scss: { icon: "🎨", colorClass: "bg-pink-600" },
  sass: { icon: "🎨", colorClass: "bg-pink-600" },
  less: { icon: "🎨", colorClass: "bg-indigo-500" },
  // Markup
  html: { icon: "🌐", colorClass: "bg-orange-500" },
  htm: { icon: "🌐", colorClass: "bg-orange-500" },
  xml: { icon: "📰", colorClass: "bg-orange-400" },
  svg: { icon: "🖼️", colorClass: "bg-yellow-500" },
  // Data/Config
  json: { icon: "📋", colorClass: "bg-yellow-400" },
  yaml: { icon: "⚙️", colorClass: "bg-red-400" },
  yml: { icon: "⚙️", colorClass: "bg-red-400" },
  toml: { icon: "⚙️", colorClass: "bg-gray-500" },
  ini: { icon: "⚙️", colorClass: "bg-gray-500" },
  env: { icon: "🔐", colorClass: "bg-yellow-600" },
  // Documentation
  md: { icon: "📝", colorClass: "bg-gray-400" },
  mdx: { icon: "📝", colorClass: "bg-gray-400" },
  txt: { icon: "📄", colorClass: "bg-gray-400" },
  // Database
  sql: { icon: "🗃️", colorClass: "bg-indigo-500" },
  // Build/Config files
  gradle: { icon: "🐘", colorClass: "bg-green-600" },
  // Shell
  sh: { icon: "📜", colorClass: "bg-green-500" },
  bash: { icon: "📜", colorClass: "bg-green-500" },
  zsh: { icon: "📜", colorClass: "bg-green-500" },
  ps1: { icon: "📜", colorClass: "bg-blue-300" },
  bat: { icon: "📜", colorClass: "bg-gray-500" },
  cmd: { icon: "📜", colorClass: "bg-gray-500" },
  // Mobile
  dart: { icon: "🎯", colorClass: "bg-blue-400" },
  swift: { icon: "🍎", colorClass: "bg-orange-500" },
  // Other
  graphql: { icon: "◈", colorClass: "bg-pink-500" },
  gql: { icon: "◈", colorClass: "bg-pink-500" },
  proto: { icon: "📡", colorClass: "bg-gray-500" },
  lock: { icon: "🔒", colorClass: "bg-gray-600" },
  default: { icon: "📄", colorClass: "bg-gray-500" }
};
const SPECIAL_FILES = {
  "pom.xml": { icon: "🏺", colorClass: "bg-red-600" },
  "build.gradle": { icon: "🐘", colorClass: "bg-green-600" },
  "build.gradle.kts": { icon: "🐘", colorClass: "bg-green-600" },
  "settings.gradle": { icon: "🐘", colorClass: "bg-green-600" },
  "settings.gradle.kts": { icon: "🐘", colorClass: "bg-green-600" },
  "package.json": { icon: "📦", colorClass: "bg-red-500" },
  "package-lock.json": { icon: "🔒", colorClass: "bg-gray-500" },
  "yarn.lock": { icon: "🔒", colorClass: "bg-blue-400" },
  "pnpm-lock.yaml": { icon: "🔒", colorClass: "bg-orange-400" },
  "tsconfig.json": { icon: "🔷", colorClass: "bg-blue-500" },
  "jsconfig.json": { icon: "🟨", colorClass: "bg-yellow-500" },
  ".gitignore": { icon: "🚫", colorClass: "bg-gray-500" },
  ".gitattributes": { icon: "🔧", colorClass: "bg-gray-500" },
  ".env": { icon: "🔐", colorClass: "bg-yellow-600" },
  ".env.local": { icon: "🔐", colorClass: "bg-yellow-600" },
  ".env.development": { icon: "🔐", colorClass: "bg-yellow-600" },
  ".env.production": { icon: "🔐", colorClass: "bg-yellow-600" },
  "Dockerfile": { icon: "🐳", colorClass: "bg-blue-400" },
  "docker-compose.yml": { icon: "🐳", colorClass: "bg-blue-400" },
  "docker-compose.yaml": { icon: "🐳", colorClass: "bg-blue-400" },
  "Makefile": { icon: "🔨", colorClass: "bg-gray-600" },
  "CMakeLists.txt": { icon: "🔨", colorClass: "bg-blue-500" },
  "go.mod": { icon: "🐹", colorClass: "bg-sky-500" },
  "go.sum": { icon: "🔒", colorClass: "bg-sky-400" },
  "Cargo.toml": { icon: "🦀", colorClass: "bg-orange-500" },
  "Cargo.lock": { icon: "🔒", colorClass: "bg-orange-400" },
  "requirements.txt": { icon: "🐍", colorClass: "bg-yellow-600" },
  "Pipfile": { icon: "🐍", colorClass: "bg-yellow-600" },
  "pyproject.toml": { icon: "🐍", colorClass: "bg-yellow-600" },
  "Gemfile": { icon: "💎", colorClass: "bg-red-500" },
  "Gemfile.lock": { icon: "🔒", colorClass: "bg-red-400" },
  "composer.json": { icon: "🐘", colorClass: "bg-indigo-400" },
  "README.md": { icon: "📖", colorClass: "bg-blue-400" },
  "LICENSE": { icon: "📜", colorClass: "bg-gray-500" },
  "LICENSE.md": { icon: "📜", colorClass: "bg-gray-500" },
  ".prettierrc": { icon: "✨", colorClass: "bg-pink-400" },
  ".eslintrc": { icon: "🔍", colorClass: "bg-purple-500" },
  ".eslintrc.js": { icon: "🔍", colorClass: "bg-purple-500" },
  ".eslintrc.json": { icon: "🔍", colorClass: "bg-purple-500" },
  "vite.config.ts": { icon: "⚡", colorClass: "bg-purple-500" },
  "vite.config.js": { icon: "⚡", colorClass: "bg-purple-500" },
  "webpack.config.js": { icon: "📦", colorClass: "bg-blue-500" },
  "rollup.config.js": { icon: "📦", colorClass: "bg-red-500" },
  "tailwind.config.js": { icon: "🎨", colorClass: "bg-cyan-500" },
  "tailwind.config.ts": { icon: "🎨", colorClass: "bg-cyan-500" },
  "postcss.config.js": { icon: "🎨", colorClass: "bg-red-500" },
  "Application.java": { icon: "🚀", colorClass: "bg-green-500" }
};
function getFileIcon(name) {
  if (SPECIAL_FILES[name]) {
    return SPECIAL_FILES[name].icon;
  }
  const ext = name.split(".").pop()?.toLowerCase() || "";
  return FILE_TYPE_CONFIG[ext]?.icon || FILE_TYPE_CONFIG.default.icon;
}
const _hoisted_1$X = ["title"];
const _hoisted_2$U = ["width"];
const _hoisted_3$Q = ["x1", "x2", "y2"];
const _hoisted_4$L = ["x1", "x2", "y2"];
const _hoisted_5$F = ["x1", "x2"];
const _hoisted_6$C = {
  key: 2,
  class: "w-5"
};
const _hoisted_7$y = ["title"];
const _hoisted_8$w = {
  key: 0,
  class: "tree-cb-icon",
  fill: "currentColor",
  viewBox: "0 0 20 20"
};
const _hoisted_9$t = {
  key: 1,
  class: "w-2 h-0.5 bg-white rounded-full"
};
const _hoisted_10$s = { class: "tree-icon" };
const _hoisted_11$q = {
  key: 0,
  class: "tree-folder-icon",
  fill: "currentColor",
  viewBox: "0 0 20 20"
};
const _hoisted_12$n = {
  key: 1,
  class: "tree-file-icon"
};
const _hoisted_13$n = { class: "tree-name" };
const _hoisted_14$l = ["title"];
const _hoisted_15$j = ["title"];
const _hoisted_16$h = ["title"];
const _hoisted_17$h = ["title"];
const _hoisted_18$g = {
  key: 2,
  class: "tree-size"
};
const rowHeight = 26;
const _sfc_main$Z = /* @__PURE__ */ defineComponent({
  __name: "VirtualTreeRow",
  props: {
    item: {},
    compactMode: { type: Boolean, default: false },
    isSelected: { type: Boolean, default: false },
    checkboxState: { default: "none" },
    fileCount: { default: 0 },
    selectedTokens: { default: 0 },
    allowSelectBinary: { type: Boolean, default: false }
  },
  emits: ["toggle-select", "toggle-expand", "contextmenu", "quicklook"],
  setup(__props, { emit: __emit }) {
    const { t } = useI18n();
    const props = __props;
    const isSelectionDisabled = computed(() => {
      if (props.item.node.isDir) return false;
      if (props.item.node.contentType !== "binary") return false;
      return !props.allowSelectBinary;
    });
    const weightLevel = computed(() => {
      const tokens = props.selectedTokens;
      if (tokens >= TOKEN_THRESHOLDS.CRITICAL) return "critical";
      if (tokens >= TOKEN_THRESHOLDS.HEAVY) return "heavy";
      if (tokens >= TOKEN_THRESHOLDS.MEDIUM) return "medium";
      return "none";
    });
    const fileTokens = computed(() => {
      if (props.item.node.isDir || !props.item.node.size) return 0;
      return Math.round(props.item.node.size / TOKEN_THRESHOLDS.BYTES_PER_TOKEN);
    });
    const fileWeightLevel = computed(() => {
      const tokens = fileTokens.value;
      if (tokens >= TOKEN_THRESHOLDS.CRITICAL) return "critical";
      if (tokens >= TOKEN_THRESHOLDS.HEAVY) return "heavy";
      if (tokens >= TOKEN_THRESHOLDS.MEDIUM) return "medium";
      return "none";
    });
    const emit = __emit;
    const { setHovered, clearHovered } = useHoveredFile();
    function handleClick() {
      if (props.item.node.isDir) {
        emit("toggle-expand", props.item.node.path);
      } else {
        emit("toggle-select", props.item.node.path);
      }
    }
    function handleExpand() {
      emit("toggle-expand", props.item.node.path);
    }
    function handleToggleSelect() {
      if (isSelectionDisabled.value) return;
      emit("toggle-select", props.item.node.path);
    }
    function handleContextMenu(event) {
      emit("contextmenu", props.item.node, event);
    }
    function handleQuickLook() {
      emit("quicklook", props.item.node.path);
    }
    function handleMouseEnter() {
      setHovered(props.item.node.path, props.item.node.isDir);
    }
    function handleMouseLeave() {
      clearHovered(props.item.node.path);
    }
    function formatSize(bytes) {
      if (bytes < 1024) return bytes + " B";
      if (bytes < 1024 * 1024) return Math.round(bytes / 1024) + " KB";
      return Math.round(bytes / (1024 * 1024)) + " MB";
    }
    function formatTokens(tokens) {
      if (tokens < 1e3) return tokens + "";
      return Math.round(tokens / 1e3) + "k";
    }
    return (_ctx, _cache) => {
      return openBlock(), createElementBlock("div", {
        class: normalizeClass([
          "tree-row group",
          props.isSelected ? "tree-row-selected" : "",
          { "opacity-50": __props.item.node.isIgnored }
        ]),
        style: normalizeStyle({ paddingLeft: `${__props.item.depth * 16 + 8}px` }),
        onClick: handleClick,
        onContextmenu: withModifiers(handleContextMenu, ["prevent"]),
        title: __props.item.node.isIgnored ? `${__props.item.node.path} (ignored)` : __props.item.node.path,
        onMouseenter: handleMouseEnter,
        onMouseleave: handleMouseLeave
      }, [
        __props.item.depth > 0 ? (openBlock(), createElementBlock("svg", {
          key: 0,
          class: "tree-guides",
          width: __props.item.depth * 16 + 16,
          height: rowHeight,
          style: { "shape-rendering": "crispEdges", "overflow": "visible" }
        }, [
          (openBlock(true), createElementBlock(Fragment, null, renderList(__props.item.ancestorHasMoreSiblings, (hasMore, idx) => {
            return openBlock(), createElementBlock(Fragment, {
              key: "v-" + idx
            }, [
              hasMore ? (openBlock(), createElementBlock("line", {
                key: 0,
                x1: 8 + (idx + 1) * 16 + 0.5,
                y1: "-4",
                x2: 8 + (idx + 1) * 16 + 0.5,
                y2: rowHeight + 4,
                class: normalizeClass(["tree-guide-line", `tree-guide-${Math.min(idx + 1, 5)}`])
              }, null, 10, _hoisted_3$Q)) : createCommentVNode("", true)
            ], 64);
          }), 128)),
          createBaseVNode("line", {
            x1: 8 + __props.item.depth * 16 + 0.5,
            y1: "-4",
            x2: 8 + __props.item.depth * 16 + 0.5,
            y2: __props.item.isLast ? 13 : rowHeight + 4,
            class: normalizeClass(["tree-guide-line", `tree-guide-${Math.min(__props.item.depth, 5)}`])
          }, null, 10, _hoisted_4$L),
          createBaseVNode("line", {
            x1: 8 + __props.item.depth * 16,
            y1: "13",
            x2: 8 + __props.item.depth * 16 + 10,
            y2: "13",
            class: normalizeClass(["tree-guide-line", `tree-guide-${Math.min(__props.item.depth, 5)}`])
          }, null, 10, _hoisted_5$F)
        ], 8, _hoisted_2$U)) : createCommentVNode("", true),
        __props.item.node.isDir ? (openBlock(), createElementBlock("div", {
          key: 1,
          class: normalizeClass(["tree-expand", __props.item.node.isExpanded ? "tree-expand-open" : ""]),
          onClick: withModifiers(handleExpand, ["stop"])
        }, [..._cache[0] || (_cache[0] = [
          createBaseVNode("svg", {
            class: "tree-expand-icon",
            fill: "currentColor",
            viewBox: "0 0 20 20"
          }, [
            createBaseVNode("path", {
              "fill-rule": "evenodd",
              d: "M7.293 14.707a1 1 0 010-1.414L10.586 10 7.293 6.707a1 1 0 011.414-1.414l4 4a1 1 0 010 1.414l-4 4a1 1 0 01-1.414 0z",
              "clip-rule": "evenodd"
            })
          ], -1)
        ])], 2)) : (openBlock(), createElementBlock("div", _hoisted_6$C)),
        createBaseVNode("div", {
          class: normalizeClass([
            "tree-cb",
            props.checkboxState !== "none" ? "tree-cb-checked" : "",
            props.checkboxState === "partial" ? "tree-cb-partial" : "",
            isSelectionDisabled.value ? "tree-cb-disabled" : ""
          ]),
          onClick: withModifiers(handleToggleSelect, ["stop"]),
          title: isSelectionDisabled.value ? unref(t)("files.binaryFile") : void 0
        }, [
          props.checkboxState === "full" ? (openBlock(), createElementBlock("svg", _hoisted_8$w, [..._cache[1] || (_cache[1] = [
            createBaseVNode("path", {
              "fill-rule": "evenodd",
              d: "M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z",
              "clip-rule": "evenodd"
            }, null, -1)
          ])])) : props.checkboxState === "partial" ? (openBlock(), createElementBlock("div", _hoisted_9$t)) : createCommentVNode("", true)
        ], 10, _hoisted_7$y),
        createBaseVNode("div", _hoisted_10$s, [
          __props.item.node.isDir ? (openBlock(), createElementBlock("svg", _hoisted_11$q, [..._cache[2] || (_cache[2] = [
            createBaseVNode("path", { d: "M2 6a2 2 0 012-2h5l2 2h5a2 2 0 012 2v6a2 2 0 01-2 2H4a2 2 0 01-2-2V6z" }, null, -1)
          ])])) : (openBlock(), createElementBlock("span", _hoisted_12$n, toDisplayString(unref(getFileIcon)(__props.item.node.name)), 1))
        ]),
        createBaseVNode("span", _hoisted_13$n, toDisplayString(__props.item.node.name), 1),
        __props.item.node.isDir ? (openBlock(), createElementBlock(Fragment, { key: 3 }, [
          props.fileCount > 0 ? (openBlock(), createElementBlock("span", {
            key: 0,
            class: "tree-count",
            title: unref(t)("files.fileCountTooltip").replace("{count}", String(props.fileCount))
          }, toDisplayString(props.fileCount), 9, _hoisted_14$l)) : createCommentVNode("", true),
          props.selectedTokens > 0 ? (openBlock(), createElementBlock("span", {
            key: 1,
            class: normalizeClass(["tree-weight", `tree-weight--${weightLevel.value}`]),
            title: unref(t)("files.selectedTokensTooltip").replace("{count}", formatTokens(props.selectedTokens))
          }, toDisplayString(formatTokens(props.selectedTokens)), 11, _hoisted_15$j)) : createCommentVNode("", true)
        ], 64)) : (openBlock(), createElementBlock(Fragment, { key: 4 }, [
          __props.item.node.contentType === "binary" ? (openBlock(), createElementBlock("span", {
            key: 0,
            class: "tree-binary-badge",
            title: unref(t)("files.binaryFile")
          }, " BIN ", 8, _hoisted_16$h)) : fileWeightLevel.value !== "none" ? (openBlock(), createElementBlock("span", {
            key: 1,
            class: normalizeClass(["tree-token-badge", `tree-token-badge--${fileWeightLevel.value}`]),
            title: unref(t)("files.tokenCountTooltip").replace("{count}", formatTokens(fileTokens.value))
          }, toDisplayString(formatTokens(fileTokens.value)), 11, _hoisted_17$h)) : __props.item.node.size ? (openBlock(), createElementBlock("span", _hoisted_18$g, toDisplayString(formatSize(__props.item.node.size)), 1)) : createCommentVNode("", true)
        ], 64)),
        !__props.item.node.isDir ? (openBlock(), createElementBlock("button", {
          key: 5,
          class: "tree-preview",
          onClick: withModifiers(handleQuickLook, ["stop"])
        }, [..._cache[3] || (_cache[3] = [
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
              d: "M15 12a3 3 0 11-6 0 3 3 0 016 0z"
            }),
            createBaseVNode("path", {
              "stroke-linecap": "round",
              "stroke-linejoin": "round",
              "stroke-width": "2",
              d: "M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"
            })
          ], -1)
        ])])) : createCommentVNode("", true)
      ], 46, _hoisted_1$X);
    };
  }
});
const _hoisted_1$W = { class: "virtual-tree-wrapper" };
const _hoisted_2$T = {
  key: 1,
  class: "empty-state"
};
const _hoisted_3$P = { class: "empty-state-text" };
const _sfc_main$Y = /* @__PURE__ */ defineComponent({
  __name: "VirtualFileTree",
  props: {
    nodes: {},
    compactMode: { type: Boolean, default: false },
    allowSelectBinary: { type: Boolean, default: false }
  },
  emits: ["toggle-select", "toggle-expand", "contextmenu", "quicklook"],
  setup(__props) {
    const { t } = useI18n();
    const fileStore = useFileStore();
    const props = __props;
    const nodesRef = toRef(props, "nodes");
    const { flattenedVisibleNodes } = useVirtualTree({ nodes: nodesRef });
    const flattenedNodes = computed(() => flattenedVisibleNodes.value);
    function isNodeSelected(node) {
      if (!node.isDir) {
        return fileStore.selectedPaths.has(node.path);
      }
      const allFiles = fileStore.getAllFilesInNode(node);
      return allFiles.some((filePath) => fileStore.selectedPaths.has(filePath));
    }
    function getCheckboxState(node) {
      if (!node.isDir) {
        return fileStore.selectedPaths.has(node.path) ? "full" : "none";
      }
      if (fileStore.selectedCount === 0) return "none";
      const allFiles = fileStore.getAllFilesInNode(node);
      if (allFiles.length === 0) return "none";
      let selectedCount = 0;
      for (const filePath of allFiles) {
        if (fileStore.selectedPaths.has(filePath)) selectedCount++;
      }
      if (selectedCount === 0) return "none";
      if (selectedCount === allFiles.length) return "full";
      return "partial";
    }
    function getFileCount(node) {
      if (!node.isDir) return 0;
      return fileStore.getAllFilesInNode(node).length;
    }
    const fileSizeMap = computed(() => {
      const map = /* @__PURE__ */ new Map();
      for (const item of flattenedVisibleNodes.value) {
        if (!item.node.isDir && item.node.size) {
          map.set(item.node.path, item.node.size);
        }
      }
      return map;
    });
    function getSelectedTokens(node) {
      if (!node.isDir) {
        if (fileStore.selectedPaths.has(node.path) && node.size) {
          return Math.round(node.size / 4);
        }
        return 0;
      }
      const allFiles = fileStore.getAllFilesInNode(node);
      let totalTokens = 0;
      for (const filePath of allFiles) {
        if (fileStore.selectedPaths.has(filePath)) {
          const size = fileSizeMap.value.get(filePath);
          if (size) {
            totalTokens += Math.round(size / 4);
          }
        }
      }
      return totalTokens;
    }
    return (_ctx, _cache) => {
      return openBlock(), createElementBlock("div", _hoisted_1$W, [
        flattenedNodes.value.length > 0 ? (openBlock(), createBlock(unref(script$2), {
          key: 0,
          class: "virtual-tree-scroller",
          items: flattenedNodes.value,
          "item-size": 26,
          "key-field": "id"
        }, {
          default: withCtx(({ item }) => [
            createVNode(_sfc_main$Z, {
              item,
              "compact-mode": __props.compactMode,
              "is-selected": isNodeSelected(item.node),
              "checkbox-state": getCheckboxState(item.node),
              "file-count": getFileCount(item.node),
              "selected-tokens": getSelectedTokens(item.node),
              "allow-select-binary": __props.allowSelectBinary,
              onToggleSelect: _cache[0] || (_cache[0] = ($event) => _ctx.$emit("toggle-select", $event)),
              onToggleExpand: _cache[1] || (_cache[1] = ($event) => _ctx.$emit("toggle-expand", $event)),
              onContextmenu: _cache[2] || (_cache[2] = (node, event) => _ctx.$emit("contextmenu", node, event)),
              onQuicklook: _cache[3] || (_cache[3] = ($event) => _ctx.$emit("quicklook", $event))
            }, null, 8, ["item", "compact-mode", "is-selected", "checkbox-state", "file-count", "selected-tokens", "allow-select-binary"])
          ]),
          _: 1
        }, 8, ["items"])) : (openBlock(), createElementBlock("div", _hoisted_2$T, [
          createBaseVNode("p", _hoisted_3$P, toDisplayString(unref(t)("files.noFiles")), 1)
        ]))
      ]);
    };
  }
});
const VirtualFileTree = /* @__PURE__ */ _export_sfc(_sfc_main$Y, [["__scopeId", "data-v-622934a3"]]);
const _hoisted_1$V = { class: "file-explorer" };
const _hoisted_2$S = { class: "file-explorer__header" };
const _hoisted_3$O = { class: "panel-header-unified" };
const _hoisted_4$K = { class: "panel-header-unified-title" };
const _hoisted_5$E = { class: "flex items-center gap-1" };
const _hoisted_6$B = { class: "flex items-center gap-2 text-xs mr-2" };
const _hoisted_7$x = { class: "relative group" };
const _hoisted_8$v = { class: "chip-unified chip-unified-accent cursor-help" };
const _hoisted_9$s = { class: "absolute right-0 top-full mt-2 hidden group-hover:block z-50 px-4 py-3 bg-gray-800/95 backdrop-blur-sm border border-gray-700/50 rounded-xl shadow-2xl whitespace-nowrap" };
const _hoisted_10$r = { class: "text-xs space-y-1.5" };
const _hoisted_11$p = { class: "text-gray-400" };
const _hoisted_12$m = { class: "text-white" };
const _hoisted_13$m = { class: "text-gray-400" };
const _hoisted_14$k = { class: "text-white" };
const _hoisted_15$i = { class: "text-gray-400" };
const _hoisted_16$g = { class: "text-indigo-400" };
const _hoisted_17$g = { class: "text-gray-400" };
const _hoisted_18$f = { class: "text-white" };
const _hoisted_19$f = { class: "text-gray-400" };
const _hoisted_20$f = { class: "text-emerald-400" };
const _hoisted_21$d = { class: "text-gray-400" };
const _hoisted_22$b = ["title"];
const _hoisted_23$b = ["title", "aria-label"];
const _hoisted_24$b = {
  key: 0,
  class: "px-3 pb-3"
};
const _hoisted_25$a = { class: "file-explorer__search" };
const _hoisted_26$a = { class: "relative group" };
const _hoisted_27$9 = ["placeholder", "aria-label"];
const _hoisted_28$9 = ["title"];
const _hoisted_29$5 = {
  class: "file-explorer__tree",
  "data-tour": "file-tree"
};
const _hoisted_30$3 = {
  key: 0,
  class: "flex items-center justify-center h-full"
};
const _hoisted_31$3 = {
  key: 1,
  class: "empty-state h-full"
};
const _hoisted_32$3 = { class: "empty-state-text" };
const _hoisted_33$3 = {
  key: 0,
  class: "text-xs text-gray-400 mb-2 px-2"
};
const _hoisted_34$2 = {
  key: 1,
  class: "empty-state h-full"
};
const _hoisted_35$2 = { class: "empty-state-content" };
const _hoisted_36$2 = { class: "empty-state-text" };
const _hoisted_37$2 = { class: "file-explorer__footer" };
const _sfc_main$X = /* @__PURE__ */ defineComponent({
  __name: "FileExplorer",
  emits: ["preview-file", "build-context"],
  setup(__props, { expose: __expose }) {
    const QuickLookModal = defineAsyncComponent(() => __vitePreload(() => import("./QuickLookModal-PV9OKgmz.js"), true ? __vite__mapDeps([0,1,2,3,4,5,6,7,8]) : void 0));
    const IgnoreRulesModal = defineAsyncComponent(() => __vitePreload(() => import("./IgnoreRulesModal-BEVOVaU3.js"), true ? __vite__mapDeps([9,1,2,5,6,7,3,4,10]) : void 0));
    const fileStore = useFileStore();
    const contextStore = useContextStore();
    const projectStore = useProjectStore();
    const uiStore = useUIStore();
    const settingsStore = useSettingsStore();
    const { t } = useI18n();
    const explorer = useFileExplorer();
    const contextMenu = useContextMenu();
    const logger2 = useLogger("FileExplorer");
    const ignoreRulesModalRef = ref();
    const searchInputRef = ref(null);
    function clearSearch() {
      explorer.clearSearch();
    }
    __expose({ searchInputRef });
    function handleContextMenu(node, event) {
      contextMenu.show(node, event);
    }
    function handleBreadcrumbNavigate(path) {
      fileStore.expandPath(path);
    }
    function handleAddSuggestedFiles(files2) {
      const normalizedFiles = files2.map((path) => path.replace(/\\/g, "/")).filter((path) => !fileStore.selectedPaths.has(path));
      if (normalizedFiles.length > 0) {
        fileStore.selectMultiple(normalizedFiles);
      }
    }
    async function handleOpenInExplorer(_path) {
      if (projectStore.currentPath) {
        try {
          const runtime$1 = await __vitePreload(() => Promise.resolve().then(() => runtime), true ? void 0 : void 0);
          runtime$1.BrowserOpenURL("file://" + projectStore.currentPath);
        } catch (error) {
          logger2.error("Failed to open in explorer:", error);
          uiStore.addToast("Failed to open in file explorer", "error");
        }
      }
    }
    function handleUndoSelection() {
      if (fileStore.undoSelection()) {
        uiStore.addToast(t("files.undoSelection"), "info");
      }
    }
    function handleRedoSelection() {
      if (fileStore.redoSelection()) {
        uiStore.addToast(t("files.redoSelection"), "info");
      }
    }
    onMounted(async () => {
      explorer.initialize();
      window.addEventListener("global-undo-selection", handleUndoSelection);
      window.addEventListener("global-redo-selection", handleRedoSelection);
      try {
        if (!projectStore.currentPath) {
          uiStore.addToast("No project selected", "warning");
          return;
        }
        await fileStore.loadFileTree(projectStore.currentPath);
      } catch (error) {
        logger2.error("Failed to load file tree:", error);
        uiStore.addToast("Failed to load project files. Please try again.", "error");
      }
    });
    watch(() => projectStore.currentPath, async (newPath, oldPath) => {
      if (newPath && newPath !== oldPath) {
        fileStore.clearSelection();
        contextStore.clearContext();
        try {
          await fileStore.loadFileTree(newPath);
        } catch (error) {
          logger2.error("Failed to load file tree after project change:", error);
          uiStore.addToast("Failed to load project files", "error");
        }
      }
    });
    onUnmounted(() => {
      explorer.cleanup();
      window.removeEventListener("global-undo-selection", handleUndoSelection);
      window.removeEventListener("global-redo-selection", handleRedoSelection);
    });
    return (_ctx, _cache) => {
      return openBlock(), createElementBlock("div", _hoisted_1$V, [
        createBaseVNode("div", _hoisted_2$S, [
          createBaseVNode("div", _hoisted_3$O, [
            createBaseVNode("div", _hoisted_4$K, [
              _cache[8] || (_cache[8] = createBaseVNode("div", { class: "panel-header-unified-icon panel-header-unified-icon-indigo" }, [
                createBaseVNode("svg", {
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
                ])
              ], -1)),
              createBaseVNode("h2", null, toDisplayString(unref(t)("files.title")), 1)
            ]),
            createBaseVNode("div", _hoisted_5$E, [
              createBaseVNode("div", _hoisted_6$B, [
                createBaseVNode("div", _hoisted_7$x, [
                  createBaseVNode("span", _hoisted_8$v, toDisplayString(unref(fileStore).selectedCount), 1),
                  createBaseVNode("div", _hoisted_9$s, [
                    createBaseVNode("div", _hoisted_10$r, [
                      _cache[15] || (_cache[15] = createBaseVNode("div", { class: "text-white font-semibold mb-2" }, "Selection Stats", -1)),
                      createBaseVNode("div", _hoisted_11$p, [
                        _cache[9] || (_cache[9] = createTextVNode("Files: ", -1)),
                        createBaseVNode("span", _hoisted_12$m, toDisplayString(unref(fileStore).selectedCount), 1)
                      ]),
                      createBaseVNode("div", _hoisted_13$m, [
                        _cache[10] || (_cache[10] = createTextVNode("Total: ", -1)),
                        createBaseVNode("span", _hoisted_14$k, toDisplayString(unref(explorer).totalFileCount.value), 1),
                        _cache[11] || (_cache[11] = createTextVNode(" files", -1))
                      ]),
                      createBaseVNode("div", _hoisted_15$i, [
                        _cache[12] || (_cache[12] = createTextVNode("Progress: ", -1)),
                        createBaseVNode("span", _hoisted_16$g, toDisplayString(unref(explorer).selectionProgress.value) + "%", 1)
                      ]),
                      createBaseVNode("div", _hoisted_17$g, [
                        _cache[13] || (_cache[13] = createTextVNode("Est. Size: ", -1)),
                        createBaseVNode("span", _hoisted_18$f, toDisplayString(Math.round(unref(fileStore).estimatedContextSize * 100) / 100) + "MB", 1)
                      ]),
                      createBaseVNode("div", _hoisted_19$f, [
                        _cache[14] || (_cache[14] = createTextVNode("Est. Tokens: ", -1)),
                        createBaseVNode("span", _hoisted_20$f, "~" + toDisplayString(Math.round(unref(fileStore).estimatedTokenCount / 1e3)) + "K", 1)
                      ])
                    ])
                  ])
                ]),
                createBaseVNode("span", _hoisted_21$d, toDisplayString(unref(t)("files.selected")), 1),
                unref(fileStore).selectedCount > 0 ? (openBlock(), createElementBlock("button", {
                  key: 0,
                  onClick: _cache[0] || (_cache[0] = //@ts-ignore
                  (...args) => unref(fileStore).clearSelection && unref(fileStore).clearSelection(...args)),
                  class: "p-1 rounded hover:bg-red-500/20 text-gray-500 hover:text-red-400 transition-colors",
                  title: unref(t)("files.clearSelection")
                }, [..._cache[16] || (_cache[16] = [
                  createBaseVNode("svg", {
                    class: "w-3.5 h-3.5",
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
                ])], 8, _hoisted_22$b)) : createCommentVNode("", true)
              ]),
              createVNode(SettingsPopover, {
                onOpenIgnoreRules: _cache[1] || (_cache[1] = ($event) => ignoreRulesModalRef.value?.open()),
                onSettingsChanged: unref(explorer).handleSettingsChange
              }, null, 8, ["onSettingsChanged"]),
              createBaseVNode("button", {
                onClick: _cache[2] || (_cache[2] = //@ts-ignore
                (...args) => unref(explorer).handleRefresh && unref(explorer).handleRefresh(...args)),
                class: "p-2 rounded-lg hover:bg-gray-700/50 text-gray-400 hover:text-white transition-colors",
                title: unref(t)("files.refresh"),
                "aria-label": unref(t)("files.refresh")
              }, [..._cache[17] || (_cache[17] = [
                createBaseVNode("svg", {
                  class: "w-4 h-4",
                  fill: "none",
                  stroke: "currentColor",
                  viewBox: "0 0 24 24",
                  "aria-hidden": "true"
                }, [
                  createBaseVNode("path", {
                    "stroke-linecap": "round",
                    "stroke-linejoin": "round",
                    "stroke-width": "2",
                    d: "M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"
                  })
                ], -1)
              ])], 8, _hoisted_23$b)
            ])
          ]),
          unref(fileStore).breadcrumbs.length > 0 ? (openBlock(), createElementBlock("div", _hoisted_24$b, [
            createVNode(_sfc_main$1b, {
              segments: unref(fileStore).breadcrumbs,
              "root-name": unref(fileStore).projectName,
              onNavigate: handleBreadcrumbNavigate,
              onOpenInExplorer: handleOpenInExplorer
            }, null, 8, ["segments", "root-name"])
          ])) : createCommentVNode("", true)
        ]),
        createVNode(QuickFiltersBar),
        createBaseVNode("div", _hoisted_25$a, [
          createBaseVNode("div", _hoisted_26$a, [
            withDirectives(createBaseVNode("input", {
              ref_key: "searchInputRef",
              ref: searchInputRef,
              "onUpdate:modelValue": _cache[3] || (_cache[3] = ($event) => unref(explorer).searchQuery.value = $event),
              type: "text",
              placeholder: unref(t)("files.searchShort"),
              "aria-label": unref(t)("files.searchShort"),
              class: normalizeClass(["search-input pr-8", { "search-input-active": unref(explorer).searchQuery.value }]),
              onInput: _cache[4] || (_cache[4] = //@ts-ignore
              (...args) => unref(explorer).handleSearch && unref(explorer).handleSearch(...args)),
              onKeydown: withKeys(clearSearch, ["escape"])
            }, null, 42, _hoisted_27$9), [
              [vModelText, unref(explorer).searchQuery.value]
            ]),
            _cache[19] || (_cache[19] = createBaseVNode("svg", {
              class: "input-icon",
              fill: "none",
              stroke: "currentColor",
              viewBox: "0 0 24 24"
            }, [
              createBaseVNode("path", {
                "stroke-linecap": "round",
                "stroke-linejoin": "round",
                "stroke-width": "2",
                d: "M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"
              })
            ], -1)),
            unref(explorer).searchQuery.value ? (openBlock(), createElementBlock("button", {
              key: 0,
              onClick: clearSearch,
              class: "search-clear-btn",
              title: unref(t)("files.clear")
            }, [..._cache[18] || (_cache[18] = [
              createBaseVNode("svg", {
                class: "w-3.5 h-3.5",
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
            ])], 8, _hoisted_28$9)) : createCommentVNode("", true)
          ])
        ]),
        createBaseVNode("div", _hoisted_29$5, [
          unref(fileStore).isLoading ? (openBlock(), createElementBlock("div", _hoisted_30$3, [..._cache[20] || (_cache[20] = [
            createBaseVNode("svg", {
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
            ], -1)
          ])])) : unref(fileStore).nodes.length === 0 ? (openBlock(), createElementBlock("div", _hoisted_31$3, [
            createBaseVNode("p", _hoisted_32$3, toDisplayString(unref(t)("files.noFiles")), 1)
          ])) : (openBlock(), createElementBlock(Fragment, { key: 2 }, [
            unref(explorer).searchQuery.value && unref(fileStore).searchResults.length > 0 ? (openBlock(), createElementBlock("div", _hoisted_33$3, toDisplayString(unref(fileStore).searchResults.length) + " " + toDisplayString(unref(t)("files.results")), 1)) : createCommentVNode("", true),
            unref(explorer).searchQuery.value && unref(fileStore).searchResults.length === 0 ? (openBlock(), createElementBlock("div", _hoisted_34$2, [
              createBaseVNode("div", _hoisted_35$2, [
                _cache[21] || (_cache[21] = createBaseVNode("svg", {
                  class: "empty-state-icon",
                  fill: "none",
                  stroke: "currentColor",
                  viewBox: "0 0 24 24"
                }, [
                  createBaseVNode("path", {
                    "stroke-linecap": "round",
                    "stroke-linejoin": "round",
                    "stroke-width": "1.5",
                    d: "M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"
                  })
                ], -1)),
                createBaseVNode("p", _hoisted_36$2, toDisplayString(unref(t)("files.noSearchResults")), 1),
                createBaseVNode("button", {
                  class: "empty-state-action",
                  onClick: _cache[5] || (_cache[5] = ($event) => unref(explorer).searchQuery.value = "")
                }, toDisplayString(unref(t)("files.clearSearch")), 1)
              ])
            ])) : !unref(explorer).searchQuery.value || unref(fileStore).searchResults.length > 0 ? (openBlock(), createBlock(VirtualFileTree, {
              key: 2,
              nodes: unref(explorer).searchQuery.value ? unref(fileStore).searchResults : unref(fileStore).filteredNodes,
              "compact-mode": unref(settingsStore).settings.fileExplorer.compactNestedFolders,
              "allow-select-binary": unref(settingsStore).settings.fileExplorer.allowSelectBinary,
              onToggleSelect: unref(explorer).handleToggleSelect,
              onToggleExpand: unref(explorer).handleToggleExpand,
              onContextmenu: handleContextMenu,
              onQuicklook: unref(explorer).handleQuickLook
            }, null, 8, ["nodes", "compact-mode", "allow-select-binary", "onToggleSelect", "onToggleExpand", "onQuicklook"])) : createCommentVNode("", true)
          ], 64))
        ]),
        createVNode(FileContextMenu, {
          node: unref(contextMenu).targetNode.value,
          position: unref(contextMenu).position.value,
          visible: unref(contextMenu).isVisible.value,
          onAction: unref(explorer).handleContextMenuAction,
          onClose: unref(contextMenu).hide
        }, null, 8, ["node", "position", "visible", "onAction", "onClose"]),
        createVNode(unref(IgnoreRulesModal), {
          ref_key: "ignoreRulesModalRef",
          ref: ignoreRulesModalRef
        }, null, 512),
        createVNode(unref(QuickLookModal), {
          modelValue: unref(explorer).quickLookVisible.value,
          "onUpdate:modelValue": _cache[6] || (_cache[6] = ($event) => unref(explorer).quickLookVisible.value = $event),
          "file-path": unref(explorer).quickLookPath.value,
          onAddToContext: unref(explorer).handleAddToContext
        }, null, 8, ["modelValue", "file-path", "onAddToContext"]),
        createVNode(AnalysisStatusBar, {
          "selected-files": Array.from(unref(fileStore).selectedPaths),
          onAddFiles: handleAddSuggestedFiles
        }, null, 8, ["selected-files"]),
        createBaseVNode("div", _hoisted_37$2, [
          createVNode(CommandBar, {
            "data-tour": "build-button",
            "selected-count": unref(fileStore).selectedPaths.size,
            "is-building": unref(contextStore).isBuilding,
            onBuild: _cache[7] || (_cache[7] = ($event) => _ctx.$emit("build-context"))
          }, null, 8, ["selected-count", "is-building"])
        ])
      ]);
    };
  }
});
const FileExplorer = /* @__PURE__ */ _export_sfc(_sfc_main$X, [["__scopeId", "data-v-78829d8d"]]);
const _hoisted_1$U = {
  key: 0,
  class: "fixed inset-0 z-50 flex items-center justify-center p-4"
};
const _hoisted_2$R = { class: "relative confirm-modal" };
const _hoisted_3$N = { class: "flex items-center gap-3 mb-4" };
const _hoisted_4$J = { class: "text-lg font-semibold text-white" };
const _hoisted_5$D = { class: "text-sm text-gray-400" };
const _hoisted_6$A = { class: "flex justify-end gap-3" };
const _sfc_main$W = /* @__PURE__ */ defineComponent({
  __name: "ContextDeleteModal",
  props: {
    show: { type: Boolean },
    message: {}
  },
  emits: ["close", "confirm"],
  setup(__props) {
    const { t } = useI18n();
    return (_ctx, _cache) => {
      return openBlock(), createBlock(Teleport, { to: "body" }, [
        __props.show ? (openBlock(), createElementBlock("div", _hoisted_1$U, [
          createBaseVNode("div", {
            class: "absolute inset-0 bg-black/60 backdrop-blur-sm",
            onClick: _cache[0] || (_cache[0] = ($event) => _ctx.$emit("close"))
          }),
          createBaseVNode("div", _hoisted_2$R, [
            createBaseVNode("div", _hoisted_3$N, [
              _cache[3] || (_cache[3] = createBaseVNode("div", { class: "w-12 h-12 rounded-xl bg-red-500/20 flex items-center justify-center" }, [
                createBaseVNode("svg", {
                  class: "w-6 h-6 text-red-400",
                  fill: "none",
                  stroke: "currentColor",
                  viewBox: "0 0 24 24"
                }, [
                  createBaseVNode("path", {
                    "stroke-linecap": "round",
                    "stroke-linejoin": "round",
                    "stroke-width": "2",
                    d: "M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"
                  })
                ])
              ], -1)),
              createBaseVNode("div", null, [
                createBaseVNode("h3", _hoisted_4$J, toDisplayString(unref(t)("context.confirmDelete")), 1),
                createBaseVNode("p", _hoisted_5$D, toDisplayString(__props.message), 1)
              ])
            ]),
            createBaseVNode("div", _hoisted_6$A, [
              createBaseVNode("button", {
                onClick: _cache[1] || (_cache[1] = ($event) => _ctx.$emit("close")),
                class: "btn btn-ghost"
              }, toDisplayString(unref(t)("context.cancel")), 1),
              createBaseVNode("button", {
                onClick: _cache[2] || (_cache[2] = ($event) => _ctx.$emit("confirm")),
                class: "btn bg-red-500 hover:bg-red-600 text-white"
              }, toDisplayString(unref(t)("context.delete")), 1)
            ])
          ])
        ])) : createCommentVNode("", true)
      ]);
    };
  }
});
const logger$8 = useLogger("ContextList");
function useContextList() {
  const { t } = useI18n();
  const contextStore = useContextStore();
  const settingsStore = useSettingsStore();
  const uiStore = useUIStore();
  const showSettings = ref(false);
  const searchQuery = ref("");
  const sortBy = ref("date");
  const showFavoritesOnly = ref(false);
  const selectedContexts = ref(/* @__PURE__ */ new Set());
  const editingId = ref(null);
  const editingName = ref("");
  const renameInput = ref(null);
  const dragIndex = ref(null);
  const dragOverIndex = ref(null);
  const deleteModal = reactive({
    show: false,
    contextId: null,
    contextIds: [],
    message: ""
  });
  const storageSettings = computed({
    get: () => settingsStore.settings.contextStorage,
    set: (val) => settingsStore.updateContextStorageSettings(val)
  });
  watch(() => settingsStore.settings.contextStorage, () => {
  }, { deep: true });
  const filteredContexts = computed(() => {
    let list = [...contextStore.contextList];
    if (searchQuery.value) {
      const q = searchQuery.value.toLowerCase();
      list = list.filter((c) => (c.name || c.id).toLowerCase().includes(q));
    }
    if (showFavoritesOnly.value) {
      list = list.filter((c) => c.isFavorite);
    }
    list.sort((a, b) => {
      if (a.isFavorite && !b.isFavorite) return -1;
      if (!a.isFavorite && b.isFavorite) return 1;
      switch (sortBy.value) {
        case "name":
          return (a.name || a.id).localeCompare(b.name || b.id);
        case "size":
          return b.totalSize - a.totalSize;
        case "date":
        default:
          return new Date(b.createdAt || 0).getTime() - new Date(a.createdAt || 0).getTime();
      }
    });
    return list;
  });
  function formatSize(bytes) {
    return formatContextSize(bytes);
  }
  function formatTimestamp$1(ts) {
    return formatTimestamp(ts);
  }
  async function refresh() {
    try {
      await contextStore.listProjectContexts();
    } catch (error) {
      logger$8.error("Failed to refresh:", error);
    }
  }
  function toggleSelect(contextId) {
    if (selectedContexts.value.has(contextId)) {
      selectedContexts.value.delete(contextId);
    } else {
      selectedContexts.value.add(contextId);
    }
    selectedContexts.value = new Set(selectedContexts.value);
  }
  function clearSelection() {
    selectedContexts.value.clear();
    selectedContexts.value = new Set(selectedContexts.value);
  }
  function toggleFavorite(contextId) {
    contextStore.toggleFavorite(contextId);
  }
  function startRename(context2) {
    editingId.value = context2.id;
    editingName.value = context2.name || context2.id;
    nextTick(() => {
      renameInput.value?.focus();
      renameInput.value?.select();
    });
  }
  function handleKeydown(e) {
    if (e.key === "Delete" && selectedContexts.value.size > 0) {
      e.preventDefault();
      deleteSelected();
    }
    if (e.key === "a" && (e.ctrlKey || e.metaKey) && filteredContexts.value.length > 0) {
      e.preventDefault();
      filteredContexts.value.forEach((c) => selectedContexts.value.add(c.id));
      selectedContexts.value = new Set(selectedContexts.value);
    }
    if (e.key === "c" && (e.ctrlKey || e.metaKey) && selectedContexts.value.size > 0) {
      e.preventDefault();
      copySelectedContext();
    }
    if (e.key === "m" && (e.ctrlKey || e.metaKey) && selectedContexts.value.size >= 2) {
      e.preventDefault();
      mergeSelected();
    }
    if (e.key === "Escape") {
      clearSelection();
      editingId.value = null;
    }
  }
  function saveRename(contextId) {
    if (editingName.value.trim()) {
      contextStore.renameContext(contextId, editingName.value.trim());
      uiStore.addToast(t("context.renamed"), "success");
    }
    cancelRename();
  }
  function cancelRename() {
    editingId.value = null;
    editingName.value = "";
  }
  async function copyContext(contextId) {
    try {
      if (contextStore.contextId !== contextId) {
        await contextStore.loadContextContent(contextId, 0, 0);
      }
      const fullContent = await contextStore.getFullContextContent();
      await navigator.clipboard.writeText(fullContent);
      uiStore.addToast(t("toast.contextCopied"), "success");
    } catch (error) {
      logger$8.error("Failed to copy:", error);
      uiStore.addToast(t("toast.copyError"), "error");
    }
  }
  async function copySelectedContext() {
    if (selectedContexts.value.size === 1) {
      const contextId = [...selectedContexts.value][0];
      await copyContext(contextId);
    } else if (selectedContexts.value.size > 1) {
      try {
        const contents = [];
        for (const ctxId of selectedContexts.value) {
          if (contextStore.contextId !== ctxId) {
            await contextStore.loadContextContent(ctxId, 0, 0);
          }
          const content = await contextStore.getFullContextContent();
          contents.push(content);
        }
        const merged = contents.join("\n\n" + "=".repeat(80) + "\n\n");
        await navigator.clipboard.writeText(merged);
        uiStore.addToast(t("toast.contextCopied"), "success");
      } catch (error) {
        logger$8.error("Failed to copy multiple:", error);
        uiStore.addToast(t("toast.copyError"), "error");
      }
    }
  }
  async function duplicateContext(contextId) {
    try {
      const newId = await contextStore.duplicateContext(contextId);
      if (newId) {
        uiStore.addToast(t("context.duplicated"), "success");
      }
    } catch (error) {
      logger$8.error("Failed to duplicate:", error);
    }
  }
  async function exportContext(contextId) {
    try {
      await contextStore.exportContext(contextId);
    } catch (error) {
      logger$8.error("Export failed:", error);
    }
  }
  function confirmDelete(context2) {
    deleteModal.contextId = context2.id;
    deleteModal.contextIds = [];
    deleteModal.message = t("context.confirmDeleteMessage").replace("{name}", context2.name || context2.id);
    deleteModal.show = true;
  }
  function deleteSelected() {
    if (selectedContexts.value.size === 0) return;
    deleteModal.contextId = null;
    deleteModal.contextIds = [...selectedContexts.value];
    deleteModal.message = t("context.confirmDeleteMultiple").replace("{count}", String(selectedContexts.value.size));
    deleteModal.show = true;
  }
  async function executeDelete() {
    try {
      if (deleteModal.contextId) {
        await contextStore.deleteContext(deleteModal.contextId);
        uiStore.addToast(t("context.deleted"), "success");
      } else if (deleteModal.contextIds.length > 0) {
        for (const id of deleteModal.contextIds) {
          await contextStore.deleteContext(id);
        }
        selectedContexts.value.clear();
        uiStore.addToast(t("context.deleted"), "success");
      }
    } catch (error) {
      logger$8.error("Delete failed:", error);
    } finally {
      deleteModal.show = false;
      deleteModal.contextId = null;
      deleteModal.contextIds = [];
    }
  }
  function closeDeleteModal() {
    deleteModal.show = false;
  }
  async function mergeSelected() {
    if (selectedContexts.value.size < 2) {
      uiStore.addToast(t("context.selectToMerge"), "warning");
      return;
    }
    try {
      const ids = [...selectedContexts.value];
      const newId = await contextStore.mergeContexts(ids);
      if (newId) {
        selectedContexts.value.clear();
        selectedContexts.value = new Set(selectedContexts.value);
        uiStore.addToast(t("context.merged"), "success");
      }
    } catch (error) {
      logger$8.error("Failed to merge:", error);
    }
  }
  function handleDragStart(e, index) {
    dragIndex.value = index;
    if (e.dataTransfer) {
      e.dataTransfer.effectAllowed = "move";
      e.dataTransfer.setData("text/plain", String(index));
    }
  }
  function handleDragOver(e, index) {
    e.preventDefault();
    if (e.dataTransfer) {
      e.dataTransfer.dropEffect = "move";
    }
    dragOverIndex.value = index;
  }
  function handleDragLeave() {
    dragOverIndex.value = null;
  }
  function handleDrop(e, toIndex) {
    e.preventDefault();
    if (dragIndex.value !== null && dragIndex.value !== toIndex) {
      contextStore.reorderContexts(dragIndex.value, toIndex);
    }
    dragIndex.value = null;
    dragOverIndex.value = null;
  }
  function handleDragEnd() {
    dragIndex.value = null;
    dragOverIndex.value = null;
  }
  return {
    // State
    showSettings,
    searchQuery,
    sortBy,
    showFavoritesOnly,
    selectedContexts,
    editingId,
    editingName,
    renameInput,
    dragIndex,
    dragOverIndex,
    deleteModal,
    storageSettings,
    // Computed
    filteredContexts,
    // Format helpers
    formatSize,
    formatTimestamp: formatTimestamp$1,
    // Actions
    refresh,
    toggleSelect,
    clearSelection,
    toggleFavorite,
    startRename,
    saveRename,
    cancelRename,
    copyContext,
    copySelectedContext,
    duplicateContext,
    exportContext,
    confirmDelete,
    deleteSelected,
    executeDelete,
    closeDeleteModal,
    mergeSelected,
    // Drag & drop
    handleDragStart,
    handleDragOver,
    handleDragLeave,
    handleDrop,
    handleDragEnd,
    // Keyboard
    handleKeydown
  };
}
const _hoisted_1$T = { class: "context-empty" };
const _hoisted_2$Q = { class: "context-empty-title" };
const _hoisted_3$M = { class: "context-empty-hint" };
const _sfc_main$V = /* @__PURE__ */ defineComponent({
  __name: "ContextListEmpty",
  setup(__props) {
    const { t } = useI18n();
    return (_ctx, _cache) => {
      return openBlock(), createElementBlock("div", _hoisted_1$T, [
        _cache[0] || (_cache[0] = createBaseVNode("div", { class: "context-empty-icon" }, [
          createBaseVNode("svg", {
            class: "w-8 h-8",
            fill: "none",
            stroke: "currentColor",
            viewBox: "0 0 24 24"
          }, [
            createBaseVNode("path", {
              "stroke-linecap": "round",
              "stroke-linejoin": "round",
              "stroke-width": "1.5",
              d: "M20 13V6a2 2 0 00-2-2H6a2 2 0 00-2 2v7m16 0v5a2 2 0 01-2 2H6a2 2 0 01-2-2v-5m16 0h-2.586a1 1 0 00-.707.293l-2.414 2.414a1 1 0 01-.707.293h-3.172a1 1 0 01-.707-.293l-2.414-2.414A1 1 0 006.586 13H4"
            })
          ])
        ], -1)),
        createBaseVNode("p", _hoisted_2$Q, toDisplayString(unref(t)("context.noSaved")), 1),
        createBaseVNode("p", _hoisted_3$M, toDisplayString(unref(t)("context.buildToSee")), 1)
      ]);
    };
  }
});
const ContextListEmpty = /* @__PURE__ */ _export_sfc(_sfc_main$V, [["__scopeId", "data-v-ca43db1d"]]);
const _hoisted_1$S = {
  key: 0,
  class: "context-footer"
};
const _sfc_main$U = /* @__PURE__ */ defineComponent({
  __name: "ContextListFooter",
  props: {
    canSave: { type: Boolean }
  },
  emits: ["save"],
  setup(__props) {
    const { t } = useI18n();
    return (_ctx, _cache) => {
      return __props.canSave ? (openBlock(), createElementBlock("div", _hoisted_1$S, [
        createBaseVNode("button", {
          onClick: _cache[0] || (_cache[0] = ($event) => _ctx.$emit("save")),
          class: "context-save-btn"
        }, [
          _cache[1] || (_cache[1] = createBaseVNode("svg", {
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
          createTextVNode(" " + toDisplayString(unref(t)("context.saveCurrentContext")), 1)
        ])
      ])) : createCommentVNode("", true);
    };
  }
});
const ContextListFooter = /* @__PURE__ */ _export_sfc(_sfc_main$U, [["__scopeId", "data-v-9f06942a"]]);
const _hoisted_1$R = ["draggable"];
const _hoisted_2$P = { class: "context-avatar-wrap" };
const _hoisted_3$L = {
  key: 0,
  class: "context-avatar-star"
};
const _hoisted_4$I = {
  key: 1,
  class: "context-avatar-initials"
};
const _hoisted_5$C = {
  key: 0,
  class: "context-checkmark",
  viewBox: "0 0 24 24",
  fill: "none"
};
const _hoisted_6$z = {
  key: 1,
  class: "context-checkbox-empty"
};
const _hoisted_7$w = { class: "context-info" };
const _hoisted_8$u = { class: "context-name-row" };
const _hoisted_9$r = ["value"];
const _hoisted_10$q = {
  key: 1,
  class: "context-name"
};
const _hoisted_11$o = { class: "context-time" };
const _hoisted_12$l = { class: "context-meta" };
const _hoisted_13$l = { class: "context-badge" };
const _hoisted_14$j = { class: "context-size" };
const _hoisted_15$h = { class: "context-actions" };
const _hoisted_16$f = ["title"];
const _hoisted_17$f = ["title"];
const _hoisted_18$e = ["fill"];
const _hoisted_19$e = ["title"];
const _hoisted_20$e = ["title"];
const _hoisted_21$c = ["title"];
const _sfc_main$T = /* @__PURE__ */ defineComponent({
  __name: "ContextListItem",
  props: {
    context: {},
    index: {},
    isSelected: { type: Boolean },
    isActive: { type: Boolean },
    isEditing: { type: Boolean },
    editingName: {},
    dragOverIndex: {},
    searchQuery: {},
    showFavoritesOnly: { type: Boolean }
  },
  emits: ["select", "toggle-select", "toggle-favorite", "load", "restore-selection", "start-rename", "save-rename", "cancel-rename", "update:editing-name", "copy", "duplicate", "export", "delete", "drag-start", "drag-over", "drag-leave", "drop", "drag-end"],
  setup(__props) {
    const { t } = useI18n();
    const props = __props;
    const hasFiles = computed(() => {
      return props.context.files && props.context.files.length > 0;
    });
    const initials = computed(() => {
      const name = props.context.name || "";
      if (!name) {
        return props.context.fileCount > 0 ? String(props.context.fileCount) : "?";
      }
      if (name.includes("Пустой")) {
        return "ПК";
      }
      const filesMatch = name.match(/^(\d+)\s*файл/);
      if (filesMatch) {
        return filesMatch[1].length > 2 ? filesMatch[1].slice(0, 2) : filesMatch[1];
      }
      const words = name.split(/[\s\-_]+/).filter((w) => w.length > 0);
      if (words.length >= 2) {
        return (words[0][0] + words[1][0]).toUpperCase();
      }
      return name.slice(0, 2).toUpperCase();
    });
    const avatarColor = computed(() => {
      const colors = ["indigo", "purple", "pink", "blue", "cyan", "teal", "green", "amber"];
      const source = props.context.id || props.context.name || "";
      let hash = 0;
      for (let i = 0; i < source.length; i++) {
        hash = source.charCodeAt(i) + ((hash << 5) - hash);
      }
      return colors[Math.abs(hash) % colors.length];
    });
    function formatSize(bytes) {
      return formatContextSize(bytes);
    }
    function formatTime(dateStr) {
      if (!dateStr) return "";
      const date = new Date(dateStr);
      const now = /* @__PURE__ */ new Date();
      const diff = now.getTime() - date.getTime();
      const mins = Math.floor(diff / 6e4);
      const hours = Math.floor(diff / 36e5);
      const days = Math.floor(diff / 864e5);
      if (mins < 1) return "сейчас";
      if (mins < 60) return `${mins}м`;
      if (hours < 24) return `${hours}ч`;
      if (days < 7) return `${days}д`;
      return date.toLocaleDateString("ru-RU", { day: "numeric", month: "short" });
    }
    function generateAutoName(ctx) {
      if (ctx.fileCount > 0) {
        return `${ctx.fileCount} файлов`;
      }
      return "Без названия";
    }
    function cleanName(name) {
      if (!name) return "";
      return name.replace(/\s*\[\d+\]\s*$/, "").trim();
    }
    return (_ctx, _cache) => {
      return openBlock(), createElementBlock("div", {
        class: normalizeClass(["context-row group", {
          "context-row--selected": __props.isSelected,
          "context-row--active": __props.isActive
        }]),
        draggable: !__props.searchQuery && !__props.showFavoritesOnly,
        onClick: _cache[11] || (_cache[11] = ($event) => _ctx.$emit("select", __props.context.id)),
        onDblclick: _cache[12] || (_cache[12] = ($event) => _ctx.$emit("load", __props.context.id)),
        onDragstart: _cache[13] || (_cache[13] = ($event) => _ctx.$emit("drag-start", $event, __props.index)),
        onDragover: _cache[14] || (_cache[14] = ($event) => _ctx.$emit("drag-over", $event, __props.index)),
        onDragleave: _cache[15] || (_cache[15] = ($event) => _ctx.$emit("drag-leave")),
        onDrop: _cache[16] || (_cache[16] = ($event) => _ctx.$emit("drop", $event, __props.index)),
        onDragend: _cache[17] || (_cache[17] = ($event) => _ctx.$emit("drag-end"))
      }, [
        createBaseVNode("div", _hoisted_2$P, [
          createBaseVNode("div", {
            class: normalizeClass(["context-avatar", [`context-avatar--${avatarColor.value}`]])
          }, [
            __props.context.isFavorite ? (openBlock(), createElementBlock("span", _hoisted_3$L, "⭐")) : (openBlock(), createElementBlock("span", _hoisted_4$I, toDisplayString(initials.value), 1))
          ], 2),
          createBaseVNode("div", {
            class: normalizeClass(["context-check-overlay", {
              "context-check-overlay--visible": __props.isSelected,
              "context-check-overlay--checked": __props.isSelected
            }]),
            onClick: _cache[0] || (_cache[0] = withModifiers(($event) => _ctx.$emit("toggle-select", __props.context.id), ["stop"]))
          }, [
            __props.isSelected ? (openBlock(), createElementBlock("svg", _hoisted_5$C, [..._cache[18] || (_cache[18] = [
              createBaseVNode("path", {
                d: "M5 13l4 4L19 7",
                stroke: "currentColor",
                "stroke-width": "3",
                "stroke-linecap": "round",
                "stroke-linejoin": "round"
              }, null, -1)
            ])])) : (openBlock(), createElementBlock("div", _hoisted_6$z))
          ], 2)
        ]),
        createBaseVNode("div", _hoisted_7$w, [
          createBaseVNode("div", _hoisted_8$u, [
            __props.isEditing ? (openBlock(), createElementBlock("input", {
              key: 0,
              value: __props.editingName,
              onInput: _cache[1] || (_cache[1] = ($event) => _ctx.$emit("update:editing-name", $event.target.value)),
              onBlur: _cache[2] || (_cache[2] = ($event) => _ctx.$emit("save-rename", __props.context.id)),
              onKeyup: [
                _cache[3] || (_cache[3] = withKeys(($event) => _ctx.$emit("save-rename", __props.context.id), ["enter"])),
                _cache[4] || (_cache[4] = withKeys(($event) => _ctx.$emit("cancel-rename"), ["escape"]))
              ],
              onClick: _cache[5] || (_cache[5] = withModifiers(() => {
              }, ["stop"])),
              class: "context-rename-input"
            }, null, 40, _hoisted_9$r)) : (openBlock(), createElementBlock("span", _hoisted_10$q, toDisplayString(cleanName(__props.context.name) || generateAutoName(__props.context)), 1)),
            createBaseVNode("span", _hoisted_11$o, toDisplayString(formatTime(__props.context.createdAt || "")), 1)
          ]),
          createBaseVNode("div", _hoisted_12$l, [
            createBaseVNode("span", _hoisted_13$l, [
              _cache[19] || (_cache[19] = createBaseVNode("svg", {
                class: "w-3 h-3",
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
              createTextVNode(" " + toDisplayString(__props.context.fileCount), 1)
            ]),
            _cache[20] || (_cache[20] = createBaseVNode("span", { class: "context-dot" }, "•", -1)),
            createBaseVNode("span", _hoisted_14$j, toDisplayString(formatSize(__props.context.totalSize)), 1)
          ])
        ]),
        createBaseVNode("div", _hoisted_15$h, [
          hasFiles.value ? (openBlock(), createElementBlock("button", {
            key: 0,
            onClick: _cache[6] || (_cache[6] = withModifiers(($event) => _ctx.$emit("restore-selection", __props.context), ["stop"])),
            class: "context-action context-action--restore",
            title: unref(t)("context.restoreSelection")
          }, [..._cache[21] || (_cache[21] = [
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
                d: "M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"
              })
            ], -1)
          ])], 8, _hoisted_16$f)) : createCommentVNode("", true),
          createBaseVNode("button", {
            onClick: _cache[7] || (_cache[7] = withModifiers(($event) => _ctx.$emit("toggle-favorite", __props.context.id), ["stop"])),
            class: normalizeClass(["context-action", { "context-action--favorite": __props.context.isFavorite }]),
            title: __props.context.isFavorite ? unref(t)("context.unfavorite") : unref(t)("context.favorite")
          }, [
            (openBlock(), createElementBlock("svg", {
              class: "w-4 h-4",
              fill: __props.context.isFavorite ? "currentColor" : "none",
              stroke: "currentColor",
              viewBox: "0 0 24 24"
            }, [..._cache[22] || (_cache[22] = [
              createBaseVNode("path", {
                "stroke-linecap": "round",
                "stroke-linejoin": "round",
                "stroke-width": "2",
                d: "M11.049 2.927c.3-.921 1.603-.921 1.902 0l1.519 4.674a1 1 0 00.95.69h4.915c.969 0 1.371 1.24.588 1.81l-3.976 2.888a1 1 0 00-.363 1.118l1.518 4.674c.3.922-.755 1.688-1.538 1.118l-3.976-2.888a1 1 0 00-1.176 0l-3.976 2.888c-.783.57-1.838-.197-1.538-1.118l1.518-4.674a1 1 0 00-.363-1.118l-3.976-2.888c-.784-.57-.38-1.81.588-1.81h4.914a1 1 0 00.951-.69l1.519-4.674z"
              }, null, -1)
            ])], 8, _hoisted_18$e))
          ], 10, _hoisted_17$f),
          createBaseVNode("button", {
            onClick: _cache[8] || (_cache[8] = withModifiers(($event) => _ctx.$emit("copy", __props.context.id), ["stop"])),
            class: "context-action",
            title: unref(t)("context.copyToClipboard")
          }, [..._cache[23] || (_cache[23] = [
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
                d: "M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z"
              })
            ], -1)
          ])], 8, _hoisted_19$e),
          createBaseVNode("button", {
            onClick: _cache[9] || (_cache[9] = withModifiers(($event) => _ctx.$emit("start-rename", __props.context), ["stop"])),
            class: "context-action",
            title: unref(t)("context.rename")
          }, [..._cache[24] || (_cache[24] = [
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
                d: "M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"
              })
            ], -1)
          ])], 8, _hoisted_20$e),
          createBaseVNode("button", {
            onClick: _cache[10] || (_cache[10] = withModifiers(($event) => _ctx.$emit("delete", __props.context), ["stop"])),
            class: "context-action context-action--danger",
            title: unref(t)("context.delete")
          }, [..._cache[25] || (_cache[25] = [
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
                d: "M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"
              })
            ], -1)
          ])], 8, _hoisted_21$c)
        ])
      ], 42, _hoisted_1$R);
    };
  }
});
const ContextListItem = /* @__PURE__ */ _export_sfc(_sfc_main$T, [["__scopeId", "data-v-11d436c1"]]);
const _sfc_main$S = {};
const _hoisted_1$Q = { class: "context-loading" };
function _sfc_render(_ctx, _cache) {
  return openBlock(), createElementBlock("div", _hoisted_1$Q, [..._cache[0] || (_cache[0] = [
    createBaseVNode("svg", {
      class: "animate-spin h-6 w-6 text-purple-500",
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
    ], -1)
  ])]);
}
const ContextListLoading = /* @__PURE__ */ _export_sfc(_sfc_main$S, [["render", _sfc_render], ["__scopeId", "data-v-4fb5bf9a"]]);
const _hoisted_1$P = { class: "context-toolbar" };
const _hoisted_2$O = { class: "context-search" };
const _hoisted_3$K = ["value", "placeholder"];
const _hoisted_4$H = ["title"];
const _hoisted_5$B = ["title"];
const _hoisted_6$y = {
  key: 0,
  class: "context-bulk-actions"
};
const _hoisted_7$v = { class: "context-bulk-count" };
const _hoisted_8$t = ["title"];
const _hoisted_9$q = ["title"];
const _sfc_main$R = /* @__PURE__ */ defineComponent({
  __name: "ContextListSearch",
  props: {
    searchQuery: {},
    showFavoritesOnly: { type: Boolean },
    sortLabel: {},
    selectedCount: {}
  },
  emits: ["update:searchQuery", "update:showFavoritesOnly", "cycle-sort", "copy-selected", "delete-selected"],
  setup(__props) {
    const { t } = useI18n();
    return (_ctx, _cache) => {
      return openBlock(), createElementBlock("div", _hoisted_1$P, [
        createBaseVNode("div", _hoisted_2$O, [
          _cache[7] || (_cache[7] = createBaseVNode("svg", {
            class: "context-search-icon",
            fill: "none",
            stroke: "currentColor",
            viewBox: "0 0 24 24"
          }, [
            createBaseVNode("path", {
              "stroke-linecap": "round",
              "stroke-linejoin": "round",
              "stroke-width": "2",
              d: "M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"
            })
          ], -1)),
          createBaseVNode("input", {
            value: __props.searchQuery,
            onInput: _cache[0] || (_cache[0] = ($event) => _ctx.$emit("update:searchQuery", $event.target.value)),
            type: "text",
            placeholder: unref(t)("context.search"),
            class: "context-search-input"
          }, null, 40, _hoisted_3$K),
          createBaseVNode("button", {
            onClick: _cache[1] || (_cache[1] = ($event) => _ctx.$emit("update:showFavoritesOnly", !__props.showFavoritesOnly)),
            class: normalizeClass(["context-filter-btn", { active: __props.showFavoritesOnly }]),
            title: unref(t)("context.showFavorites")
          }, [..._cache[5] || (_cache[5] = [
            createBaseVNode("svg", {
              class: "w-4 h-4",
              fill: "currentColor",
              viewBox: "0 0 24 24"
            }, [
              createBaseVNode("path", { d: "M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14 2 9.27l6.91-1.01L12 2z" })
            ], -1)
          ])], 10, _hoisted_4$H),
          createBaseVNode("button", {
            onClick: _cache[2] || (_cache[2] = ($event) => _ctx.$emit("cycle-sort")),
            class: "context-filter-btn",
            title: __props.sortLabel
          }, [..._cache[6] || (_cache[6] = [
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
                d: "M3 4h13M3 8h9m-9 4h6m4 0l4-4m0 0l4 4m-4-4v12"
              })
            ], -1)
          ])], 8, _hoisted_5$B)
        ]),
        __props.selectedCount > 0 ? (openBlock(), createElementBlock("div", _hoisted_6$y, [
          createBaseVNode("span", _hoisted_7$v, toDisplayString(__props.selectedCount), 1),
          createBaseVNode("button", {
            onClick: _cache[3] || (_cache[3] = ($event) => _ctx.$emit("copy-selected")),
            class: "context-bulk-btn",
            title: unref(t)("context.copyToClipboard")
          }, [..._cache[8] || (_cache[8] = [
            createBaseVNode("svg", {
              class: "w-3.5 h-3.5",
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
            ], -1)
          ])], 8, _hoisted_8$t),
          createBaseVNode("button", {
            onClick: _cache[4] || (_cache[4] = ($event) => _ctx.$emit("delete-selected")),
            class: "context-bulk-btn context-bulk-btn--danger",
            title: unref(t)("context.deleteSelected")
          }, [..._cache[9] || (_cache[9] = [
            createBaseVNode("svg", {
              class: "w-3.5 h-3.5",
              fill: "none",
              stroke: "currentColor",
              viewBox: "0 0 24 24"
            }, [
              createBaseVNode("path", {
                "stroke-linecap": "round",
                "stroke-linejoin": "round",
                "stroke-width": "2",
                d: "M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"
              })
            ], -1)
          ])], 8, _hoisted_9$q)
        ])) : createCommentVNode("", true)
      ]);
    };
  }
});
const ContextListSearch = /* @__PURE__ */ _export_sfc(_sfc_main$R, [["__scopeId", "data-v-60c29574"]]);
const _hoisted_1$O = {
  key: 0,
  class: "fixed inset-0 z-50 flex items-center justify-center p-4"
};
const _hoisted_2$N = { class: "save-modal" };
const _hoisted_3$J = { class: "save-modal-title" };
const _hoisted_4$G = { class: "save-modal-field" };
const _hoisted_5$A = { class: "save-modal-label" };
const _hoisted_6$x = ["placeholder"];
const _hoisted_7$u = { class: "save-modal-field" };
const _hoisted_8$s = { class: "save-modal-label" };
const _hoisted_9$p = ["placeholder"];
const _hoisted_10$p = { class: "save-modal-footer" };
const _hoisted_11$n = { class: "save-modal-badge" };
const _hoisted_12$k = { class: "save-modal-actions" };
const _hoisted_13$k = ["disabled"];
const _sfc_main$Q = /* @__PURE__ */ defineComponent({
  __name: "ContextSaveDialog",
  props: {
    show: { type: Boolean },
    fileCount: {}
  },
  emits: ["close", "save"],
  setup(__props, { emit: __emit }) {
    const { t } = useI18n();
    const props = __props;
    const emit = __emit;
    const localTopic = ref("");
    const localSummary = ref("");
    const topicInputRef = ref(null);
    watch(() => props.show, (isOpen) => {
      if (isOpen) {
        localTopic.value = "";
        localSummary.value = "";
        nextTick(() => {
          topicInputRef.value?.focus();
        });
      }
    });
    function handleSubmit() {
      if (!localTopic.value.trim()) return;
      emit("save", localTopic.value.trim(), localSummary.value.trim());
    }
    return (_ctx, _cache) => {
      return openBlock(), createBlock(Teleport, { to: "body" }, [
        __props.show ? (openBlock(), createElementBlock("div", _hoisted_1$O, [
          createBaseVNode("div", {
            class: "absolute inset-0 bg-black/70 backdrop-blur-sm",
            onClick: _cache[0] || (_cache[0] = ($event) => _ctx.$emit("close"))
          }),
          createBaseVNode("div", _hoisted_2$N, [
            _cache[6] || (_cache[6] = createBaseVNode("div", { class: "save-modal-accent" }, null, -1)),
            createBaseVNode("h3", _hoisted_3$J, toDisplayString(unref(t)("context.saveContext")), 1),
            createBaseVNode("form", {
              onSubmit: withModifiers(handleSubmit, ["prevent"]),
              class: "save-modal-content"
            }, [
              createBaseVNode("div", _hoisted_4$G, [
                createBaseVNode("label", _hoisted_5$A, toDisplayString(unref(t)("context.topic")), 1),
                withDirectives(createBaseVNode("input", {
                  ref_key: "topicInputRef",
                  ref: topicInputRef,
                  "onUpdate:modelValue": _cache[1] || (_cache[1] = ($event) => localTopic.value = $event),
                  type: "text",
                  class: "save-modal-input",
                  placeholder: unref(t)("context.topicPlaceholder"),
                  onKeyup: _cache[2] || (_cache[2] = withKeys(($event) => localTopic.value.trim() && handleSubmit(), ["enter"]))
                }, null, 40, _hoisted_6$x), [
                  [vModelText, localTopic.value]
                ])
              ]),
              createBaseVNode("div", _hoisted_7$u, [
                createBaseVNode("label", _hoisted_8$s, toDisplayString(unref(t)("context.summary")), 1),
                withDirectives(createBaseVNode("textarea", {
                  "onUpdate:modelValue": _cache[3] || (_cache[3] = ($event) => localSummary.value = $event),
                  class: "save-modal-textarea",
                  placeholder: unref(t)("context.summaryPlaceholder")
                }, null, 8, _hoisted_9$p), [
                  [vModelText, localSummary.value]
                ])
              ])
            ], 32),
            createBaseVNode("div", _hoisted_10$p, [
              createBaseVNode("div", _hoisted_11$n, [
                _cache[5] || (_cache[5] = createBaseVNode("svg", {
                  class: "w-3.5 h-3.5",
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
                createBaseVNode("span", null, toDisplayString(__props.fileCount) + " " + toDisplayString(unref(t)("context.filesShort")), 1)
              ]),
              createBaseVNode("div", _hoisted_12$k, [
                createBaseVNode("button", {
                  type: "button",
                  onClick: _cache[4] || (_cache[4] = ($event) => _ctx.$emit("close")),
                  class: "save-modal-cancel"
                }, toDisplayString(unref(t)("context.cancel")), 1),
                createBaseVNode("button", {
                  type: "submit",
                  onClick: handleSubmit,
                  class: "save-modal-submit",
                  disabled: !localTopic.value.trim()
                }, toDisplayString(unref(t)("common.save")), 9, _hoisted_13$k)
              ])
            ])
          ])
        ])) : createCommentVNode("", true)
      ]);
    };
  }
});
const ContextSaveDialog = /* @__PURE__ */ _export_sfc(_sfc_main$Q, [["__scopeId", "data-v-989e9974"]]);
const _hoisted_1$N = { class: "context-list-panel" };
const _hoisted_2$M = {
  key: 3,
  class: "context-list"
};
const _sfc_main$P = /* @__PURE__ */ defineComponent({
  __name: "ContextList",
  emits: ["switch-to-preview"],
  setup(__props, { emit: __emit }) {
    const logger2 = useLogger("ContextList");
    const { t } = useI18n();
    const contextStore = useContextStore();
    const fileStore = useFileStore();
    const projectStore = useProjectStore();
    const uiStore = useUIStore();
    const list = useContextList();
    const showSaveDialog = ref(false);
    const canSaveContext = computed(() => fileStore.selectedPaths.size > 0 && !!projectStore.currentPath);
    const emit = __emit;
    const sortLabel = computed(() => {
      const labels = {
        date: t("context.sortByDate"),
        name: t("context.sortByName"),
        size: t("context.sortBySize")
      };
      return labels[list.sortBy.value] || labels.date;
    });
    function cycleSortBy() {
      const options2 = ["date", "name", "size"];
      const current = options2.indexOf(list.sortBy.value);
      list.sortBy.value = options2[(current + 1) % options2.length];
    }
    function selectContext(contextId) {
      contextStore.selectListItem(contextId);
    }
    async function loadContext(contextId) {
      try {
        contextStore.selectListItem(contextId);
        await contextStore.loadContextContent(contextId, 0, 0);
        emit("switch-to-preview");
      } catch (error) {
        logger2.error("Failed to load context:", error);
      }
    }
    function handleRestoreSelection(context2) {
      if (!context2.files || context2.files.length === 0) {
        uiStore.addToast(t("context.noFilesToRestore"), "warning");
        return;
      }
      fileStore.clearSelection();
      fileStore.selectMultiple(context2.files);
      uiStore.addToast(
        t("context.selectionRestored").replace("{count}", String(context2.files.length)).replace("{name}", context2.name || ""),
        "success"
      );
    }
    function handleStartRename(context2) {
      list.startRename(context2);
    }
    async function saveContext(topic, summary) {
      if (!projectStore.currentPath || !topic.trim()) return;
      try {
        const normalizedFiles = Array.from(fileStore.selectedPaths).map(
          (path) => path.replace(/\\/g, "/")
        );
        await apiService.saveContextMemory(
          projectStore.currentPath,
          topic.trim(),
          summary.trim(),
          normalizedFiles
        );
        uiStore.addToast(t("context.contextSaved"), "success");
        showSaveDialog.value = false;
        list.refresh();
      } catch {
        uiStore.addToast(t("context.saveError"), "error");
      }
    }
    onMounted(() => {
      list.refresh();
      window.addEventListener("keydown", list.handleKeydown);
    });
    onUnmounted(() => {
      window.removeEventListener("keydown", list.handleKeydown);
    });
    return (_ctx, _cache) => {
      return openBlock(), createElementBlock("div", _hoisted_1$N, [
        unref(contextStore).contextList.length > 0 ? (openBlock(), createBlock(ContextListSearch, {
          key: 0,
          "search-query": unref(list).searchQuery.value,
          "onUpdate:searchQuery": _cache[0] || (_cache[0] = ($event) => unref(list).searchQuery.value = $event),
          "show-favorites-only": unref(list).showFavoritesOnly.value,
          "onUpdate:showFavoritesOnly": _cache[1] || (_cache[1] = ($event) => unref(list).showFavoritesOnly.value = $event),
          "sort-label": sortLabel.value,
          "selected-count": unref(list).selectedContexts.value.size,
          onCycleSort: cycleSortBy,
          onCopySelected: unref(list).copySelectedContext,
          onDeleteSelected: unref(list).deleteSelected
        }, null, 8, ["search-query", "show-favorites-only", "sort-label", "selected-count", "onCopySelected", "onDeleteSelected"])) : createCommentVNode("", true),
        unref(contextStore).isLoading ? (openBlock(), createBlock(ContextListLoading, { key: 1 })) : unref(contextStore).contextList.length === 0 ? (openBlock(), createBlock(ContextListEmpty, { key: 2 })) : (openBlock(), createElementBlock("div", _hoisted_2$M, [
          (openBlock(true), createElementBlock(Fragment, null, renderList(unref(list).filteredContexts.value, (context2, index) => {
            return openBlock(), createBlock(ContextListItem, {
              key: context2.id,
              context: context2,
              index,
              "is-selected": unref(list).selectedContexts.value.has(context2.id),
              "is-active": unref(contextStore).selectedListItem?.id === context2.id,
              "is-editing": unref(list).editingId.value === context2.id,
              "editing-name": unref(list).editingName.value,
              "drag-over-index": unref(list).dragOverIndex.value,
              "search-query": unref(list).searchQuery.value,
              "show-favorites-only": unref(list).showFavoritesOnly.value,
              onSelect: selectContext,
              onToggleSelect: unref(list).toggleSelect,
              onToggleFavorite: unref(list).toggleFavorite,
              onLoad: loadContext,
              onRestoreSelection: handleRestoreSelection,
              onStartRename: handleStartRename,
              onSaveRename: unref(list).saveRename,
              onCancelRename: unref(list).cancelRename,
              "onUpdate:editingName": _cache[2] || (_cache[2] = ($event) => unref(list).editingName.value = $event),
              onCopy: unref(list).copyContext,
              onDuplicate: unref(list).duplicateContext,
              onExport: unref(list).exportContext,
              onDelete: unref(list).confirmDelete,
              onDragStart: unref(list).handleDragStart,
              onDragOver: unref(list).handleDragOver,
              onDragLeave: unref(list).handleDragLeave,
              onDrop: unref(list).handleDrop,
              onDragEnd: unref(list).handleDragEnd
            }, null, 8, ["context", "index", "is-selected", "is-active", "is-editing", "editing-name", "drag-over-index", "search-query", "show-favorites-only", "onToggleSelect", "onToggleFavorite", "onSaveRename", "onCancelRename", "onCopy", "onDuplicate", "onExport", "onDelete", "onDragStart", "onDragOver", "onDragLeave", "onDrop", "onDragEnd"]);
          }), 128))
        ])),
        createVNode(ContextListFooter, {
          "can-save": canSaveContext.value,
          onSave: _cache[3] || (_cache[3] = ($event) => showSaveDialog.value = true)
        }, null, 8, ["can-save"]),
        createVNode(ContextSaveDialog, {
          show: showSaveDialog.value,
          "file-count": unref(fileStore).selectedPaths.size,
          onClose: _cache[4] || (_cache[4] = ($event) => showSaveDialog.value = false),
          onSave: saveContext
        }, null, 8, ["show", "file-count"]),
        createVNode(_sfc_main$W, {
          show: unref(list).deleteModal.show,
          message: unref(list).deleteModal.message,
          onClose: _cache[5] || (_cache[5] = ($event) => unref(list).deleteModal.show = false),
          onConfirm: unref(list).executeDelete
        }, null, 8, ["show", "message", "onConfirm"])
      ]);
    };
  }
});
const ContextList = /* @__PURE__ */ _export_sfc(_sfc_main$P, [["__scopeId", "data-v-62587970"]]);
const DEFAULT_SECTION_ORDER = [
  "role",
  "rules",
  "tree",
  "stats",
  "task",
  "files"
];
const SECTION_META = [
  { key: "role", icon: "🎭", color: "purple" },
  { key: "rules", icon: "📋", color: "blue" },
  { key: "tree", icon: "🌳", color: "green" },
  { key: "stats", icon: "📊", color: "orange" },
  { key: "task", icon: "📝", color: "yellow" },
  { key: "files", icon: "📁", color: "cyan" }
];
const DEFAULT_TEMPLATES = [
  {
    id: "architect",
    name: "Architect",
    icon: "🏗️",
    description: "System design and architecture analysis",
    tags: ["architecture", "refactor"],
    isBuiltIn: true,
    isFavorite: false,
    isHidden: false,
    sections: { role: true, rules: true, tree: true, stats: true, task: true, files: true },
    sectionOrder: ["role", "rules", "tree", "stats", "task", "files"],
    roleContent: "You are a senior software architect. Analyze the codebase structure, identify patterns, and provide architectural recommendations.",
    rulesContent: "- Focus on scalability and maintainability\n- Consider SOLID principles\n- Identify potential technical debt\n- Suggest improvements with rationale",
    customPrefix: "",
    customSuffix: ""
  },
  {
    id: "implement",
    name: "Implement",
    icon: "⚡",
    description: "Code implementation and feature development",
    tags: ["implementation"],
    isBuiltIn: true,
    isFavorite: false,
    isHidden: false,
    sections: { role: true, rules: true, tree: true, stats: false, task: true, files: true },
    sectionOrder: ["role", "rules", "tree", "task", "files"],
    roleContent: "You are an expert developer. Implement the requested feature following the existing code patterns and conventions.",
    rulesContent: "- Follow existing code style\n- Write clean, readable code\n- Add appropriate comments\n- Consider edge cases",
    customPrefix: "",
    customSuffix: ""
  },
  {
    id: "review",
    name: "Review",
    icon: "🔍",
    description: "Code review and quality analysis",
    tags: ["review", "bugfix"],
    isBuiltIn: true,
    isFavorite: false,
    isHidden: false,
    sections: { role: true, rules: true, tree: false, stats: true, task: true, files: true },
    sectionOrder: ["role", "rules", "stats", "task", "files"],
    roleContent: "You are a code reviewer. Analyze the code for bugs, security issues, performance problems, and best practices violations.",
    rulesContent: "- Check for security vulnerabilities\n- Identify performance bottlenecks\n- Verify error handling\n- Suggest improvements",
    customPrefix: "",
    customSuffix: ""
  },
  {
    id: "explain",
    name: "Explain",
    icon: "📚",
    description: "Code explanation and documentation",
    tags: ["documentation"],
    isBuiltIn: true,
    isFavorite: false,
    isHidden: false,
    sections: { role: true, rules: false, tree: true, stats: false, task: true, files: true },
    sectionOrder: ["role", "tree", "task", "files"],
    roleContent: "You are a technical writer. Explain the code clearly, document its purpose, and describe how it works.",
    rulesContent: "",
    customPrefix: "",
    customSuffix: ""
  }
];
function createEmptyTemplate() {
  return {
    name: "",
    icon: "📝",
    description: "",
    tags: [],
    isBuiltIn: false,
    isFavorite: false,
    isHidden: false,
    sections: { role: true, rules: true, tree: true, stats: true, task: true, files: true },
    sectionOrder: [...DEFAULT_SECTION_ORDER],
    roleContent: "",
    rulesContent: "",
    customPrefix: "",
    customSuffix: ""
  };
}
const SUGGESTION_KEYWORDS = {
  architect: ["architecture", "design", "structure", "refactor", "модуль", "архитектур", "структур"],
  implement: ["implement", "add", "create", "build", "feature", "добавить", "создать", "реализовать", "функци"],
  review: ["review", "check", "bug", "fix", "error", "проверить", "баг", "ошибк", "исправить"],
  explain: ["explain", "what", "how", "why", "document", "объясни", "как", "почему", "документ"]
};
const logger$7 = useLogger("TemplateStore");
const STORAGE_KEY = "prompt-templates-v3";
const MAX_TASK_HISTORY = 10;
const useTemplateStore = defineStore("templates", () => {
  const templates2 = ref([]);
  const activeTemplateId = ref("architect");
  const currentTask = ref("");
  const userRules = ref("");
  const taskHistory = ref([]);
  const isModalOpen = ref(false);
  function initTemplates() {
    const saved = loadFromStorage();
    const now = (/* @__PURE__ */ new Date()).toISOString();
    const builtIn = DEFAULT_TEMPLATES.map((t) => ({
      ...t,
      createdAt: now,
      updatedAt: now,
      isFavorite: saved.find((s) => s.id === t.id)?.isFavorite ?? false,
      isHidden: saved.find((s) => s.id === t.id)?.isHidden ?? false
    }));
    templates2.value = [...builtIn, ...saved.filter((t) => !t.isBuiltIn)];
    const savedActiveId = localStorage.getItem("active-template-id");
    if (savedActiveId && templates2.value.some((t) => t.id === savedActiveId)) activeTemplateId.value = savedActiveId;
    const savedTask = localStorage.getItem("template-task");
    if (savedTask) currentTask.value = savedTask;
    const savedRules = localStorage.getItem("template-user-rules");
    if (savedRules) userRules.value = savedRules;
    try {
      const h2 = localStorage.getItem("template-task-history");
      if (h2) taskHistory.value = JSON.parse(h2);
    } catch {
    }
  }
  const activeTemplate = computed(() => templates2.value.find((t) => t.id === activeTemplateId.value) || templates2.value[0]);
  const builtInTemplates = computed(() => templates2.value.filter((t) => t.isBuiltIn && !t.isHidden));
  const customTemplates = computed(() => templates2.value.filter((t) => !t.isBuiltIn));
  const favoriteTemplates = computed(() => templates2.value.filter((t) => t.isFavorite && !t.isHidden));
  const visibleTemplates = computed(() => templates2.value.filter((t) => !t.isHidden));
  function loadFromStorage() {
    try {
      const s = localStorage.getItem(STORAGE_KEY);
      if (s) return JSON.parse(s);
    } catch (e) {
      logger$7.warn("Load failed:", e);
    }
    return [];
  }
  function saveToStorage() {
    try {
      localStorage.setItem(STORAGE_KEY, JSON.stringify(templates2.value));
      localStorage.setItem("active-template-id", activeTemplateId.value);
      localStorage.setItem("template-task", currentTask.value);
      localStorage.setItem("template-user-rules", userRules.value);
      localStorage.setItem("template-task-history", JSON.stringify(taskHistory.value));
    } catch (e) {
      logger$7.warn("Save failed:", e);
    }
  }
  watch([templates2, activeTemplateId, currentTask, userRules, taskHistory], saveToStorage, { deep: true });
  function setActiveTemplate(id) {
    if (templates2.value.some((t) => t.id === id)) activeTemplateId.value = id;
  }
  function setTask(task2) {
    currentTask.value = task2;
  }
  function setUserRules(rules) {
    userRules.value = rules;
  }
  function addToTaskHistory(task2) {
    if (!task2.trim()) return;
    taskHistory.value = taskHistory.value.filter((h2) => h2.text !== task2);
    taskHistory.value.unshift({ id: `task-${Date.now()}`, text: task2, templateId: activeTemplateId.value, timestamp: (/* @__PURE__ */ new Date()).toISOString() });
    if (taskHistory.value.length > MAX_TASK_HISTORY) taskHistory.value = taskHistory.value.slice(0, MAX_TASK_HISTORY);
  }
  function clearTaskHistory() {
    taskHistory.value = [];
  }
  function toggleFavorite(id) {
    const t = templates2.value.find((x) => x.id === id);
    if (t) {
      t.isFavorite = !t.isFavorite;
      t.updatedAt = (/* @__PURE__ */ new Date()).toISOString();
    }
  }
  function toggleHidden(id) {
    const t = templates2.value.find((x) => x.id === id);
    if (t) {
      t.isHidden = !t.isHidden;
      t.updatedAt = (/* @__PURE__ */ new Date()).toISOString();
    }
  }
  function createTemplate(tpl) {
    const now = (/* @__PURE__ */ new Date()).toISOString(), id = `custom-${Date.now()}`;
    templates2.value.push({ ...tpl, id, isBuiltIn: false, createdAt: now, updatedAt: now });
    return id;
  }
  function updateTemplate(id, updates) {
    const idx = templates2.value.findIndex((t) => t.id === id);
    if (idx !== -1) templates2.value[idx] = { ...templates2.value[idx], ...updates, updatedAt: (/* @__PURE__ */ new Date()).toISOString() };
  }
  function deleteTemplate(id) {
    const t = templates2.value.find((x) => x.id === id);
    if (t?.isBuiltIn) return false;
    templates2.value = templates2.value.filter((x) => x.id !== id);
    if (activeTemplateId.value === id) activeTemplateId.value = templates2.value[0]?.id || "architect";
    return true;
  }
  function duplicateTemplate(id) {
    const t = templates2.value.find((x) => x.id === id);
    if (!t) return null;
    return createTemplate({ ...t, name: `${t.name} (copy)`, isBuiltIn: false, isFavorite: false, isHidden: false });
  }
  function resetToDefault(id) {
    const t = templates2.value.find((x) => x.id === id);
    if (!t?.isBuiltIn) return false;
    const defaultTpl = DEFAULT_TEMPLATES.find((d) => d.id === id);
    if (!defaultTpl) return false;
    t.isFavorite = false;
    t.isHidden = false;
    t.updatedAt = (/* @__PURE__ */ new Date()).toISOString();
    return true;
  }
  function openModal() {
    isModalOpen.value = true;
  }
  function closeModal() {
    isModalOpen.value = false;
  }
  function suggestTemplate(taskText) {
    if (!taskText.trim()) return null;
    const lower = taskText.toLowerCase();
    for (const [tid, kws] of Object.entries(SUGGESTION_KEYWORDS)) {
      if (kws.some((kw) => lower.includes(kw))) return tid;
    }
    return null;
  }
  function generatePrompt(context2) {
    const tpl = activeTemplate.value;
    if (!tpl) return context2.files;
    const parts = [];
    if (tpl.customPrefix) parts.push(tpl.customPrefix);
    for (const sec of tpl.sectionOrder || DEFAULT_SECTION_ORDER) {
      if (!tpl.sections[sec]) continue;
      switch (sec) {
        case "role":
          if (tpl.roleContent) parts.push(`## Role
${tpl.roleContent}`);
          break;
        case "rules":
          if (tpl.rulesContent || userRules.value) parts.push(`## Rules
${[tpl.rulesContent, userRules.value].filter(Boolean).join("\n\n")}`);
          break;
        case "tree":
          if (context2.fileTree) parts.push(`## Project Structure
\`\`\`
${context2.fileTree}
\`\`\``);
          break;
        case "stats":
          parts.push(`## Context Stats
- Files: ${context2.fileCount}
- Tokens: ~${context2.tokenCount}
- Languages: ${context2.languages.join(", ") || "N/A"}`);
          break;
        case "task":
          if (currentTask.value) parts.push(`## Task
${currentTask.value}`);
          break;
        case "files":
          if (context2.files) parts.push(`## Files
${context2.files}`);
          break;
      }
    }
    if (tpl.customSuffix) parts.push(tpl.customSuffix);
    return parts.join("\n\n");
  }
  function generatePreview(context2) {
    const tpl = activeTemplate.value;
    if (!tpl) return "";
    const parts = [];
    if (tpl.customPrefix) parts.push(tpl.customPrefix);
    for (const sec of tpl.sectionOrder || DEFAULT_SECTION_ORDER) {
      if (!tpl.sections[sec]) continue;
      switch (sec) {
        case "role":
          if (tpl.roleContent) parts.push(`## Role
${tpl.roleContent}`);
          break;
        case "rules":
          if (tpl.rulesContent || userRules.value) parts.push(`## Rules
${[tpl.rulesContent, userRules.value].filter(Boolean).join("\n\n")}`);
          break;
        case "tree":
          parts.push(`## Project Structure
\`\`\`
${context2.fileTree || "[File tree]"}
\`\`\``);
          break;
        case "stats":
          parts.push(`## Context Stats
- Files: ${context2.fileCount || 0}
- Tokens: ~${context2.tokenCount || 0}
- Languages: ${context2.languages?.join(", ") || "N/A"}`);
          break;
        case "task":
          parts.push(`## Task
${currentTask.value || "[Your task]"}`);
          break;
        case "files":
          parts.push(`## Files
[${context2.fileCount || 0} files]`);
          break;
      }
    }
    if (tpl.customSuffix) parts.push(tpl.customSuffix);
    return parts.join("\n\n");
  }
  function getSectionLabel(sec) {
    return { role: "Role", rules: "Rules", tree: "File Tree", stats: "Stats", task: "Task", files: "Files" }[sec];
  }
  function getSectionMeta(sec) {
    return SECTION_META.find((m) => m.key === sec);
  }
  initTemplates();
  return {
    templates: templates2,
    activeTemplateId,
    currentTask,
    userRules,
    taskHistory,
    isModalOpen,
    activeTemplate,
    builtInTemplates,
    customTemplates,
    favoriteTemplates,
    visibleTemplates,
    setActiveTemplate,
    setTask,
    setUserRules,
    addToTaskHistory,
    clearTaskHistory,
    toggleFavorite,
    toggleHidden,
    createTemplate,
    updateTemplate,
    deleteTemplate,
    duplicateTemplate,
    resetToDefault,
    generatePrompt,
    generatePreview,
    suggestTemplate,
    getSectionLabel,
    getSectionMeta,
    openModal,
    closeModal
  };
});
function generateFileTree(filePaths, projectName = "") {
  if (!filePaths || filePaths.length === 0) return "";
  const root2 = { name: projectName || "project", children: [] };
  for (const path of filePaths) {
    const parts = path.split(/[/\\]/).filter(Boolean);
    let current = root2;
    for (let i = 0; i < parts.length; i++) {
      const part = parts[i];
      const isLast = i === parts.length - 1;
      let child = current.children.find((c) => c.name === part);
      if (!child) {
        child = { name: part, isDir: !isLast, children: [] };
        current.children.push(child);
      }
      current = child;
    }
  }
  sortTree(root2);
  const lines = [];
  lines.push(root2.name + "/");
  renderTree(root2.children, "", lines);
  return lines.join("\n");
}
function sortTree(node) {
  node.children.sort((a, b) => {
    if (a.isDir !== b.isDir) return a.isDir ? -1 : 1;
    return a.name.localeCompare(b.name);
  });
  for (const child of node.children) {
    sortTree(child);
  }
}
function renderTree(nodes, prefix, lines) {
  for (let i = 0; i < nodes.length; i++) {
    const node = nodes[i];
    const isLast = i === nodes.length - 1;
    const connector = isLast ? "└── " : "├── ";
    const suffix = node.isDir ? "/" : "";
    lines.push(prefix + connector + node.name + suffix);
    if (node.children.length > 0) {
      const childPrefix = prefix + (isLast ? "    " : "│   ");
      renderTree(node.children, childPrefix, lines);
    }
  }
}
function detectLanguages(filePaths) {
  const extToLang = {
    ".ts": "TypeScript",
    ".tsx": "TypeScript",
    ".js": "JavaScript",
    ".jsx": "JavaScript",
    ".vue": "Vue",
    ".go": "Go",
    ".py": "Python",
    ".java": "Java",
    ".kt": "Kotlin",
    ".rs": "Rust",
    ".cpp": "C++",
    ".c": "C",
    ".cs": "C#",
    ".rb": "Ruby",
    ".php": "PHP",
    ".swift": "Swift",
    ".scala": "Scala",
    ".html": "HTML",
    ".css": "CSS",
    ".scss": "SCSS",
    ".sql": "SQL",
    ".sh": "Shell",
    ".yaml": "YAML",
    ".yml": "YAML",
    ".json": "JSON",
    ".md": "Markdown"
  };
  const languages = /* @__PURE__ */ new Set();
  for (const path of filePaths) {
    const ext = "." + path.split(".").pop()?.toLowerCase();
    if (ext && extToLang[ext]) {
      languages.add(extToLang[ext]);
    }
  }
  return Array.from(languages).sort();
}
const _hoisted_1$M = { class: "tpl-card-title" };
const _hoisted_2$L = {
  key: 0,
  class: "tpl-card-count"
};
const _hoisted_3$I = {
  key: 0,
  class: "tpl-card-body"
};
const _sfc_main$O = /* @__PURE__ */ defineComponent({
  __name: "TemplateCard",
  props: {
    title: {},
    icon: {},
    count: {},
    enabled: { type: Boolean },
    toggleable: { type: Boolean },
    collapsible: { type: Boolean }
  },
  emits: ["toggle"],
  setup(__props) {
    const isCollapsed = ref(false);
    function toggleCollapse() {
      isCollapsed.value = !isCollapsed.value;
    }
    return (_ctx, _cache) => {
      return openBlock(), createElementBlock("div", {
        class: normalizeClass(["tpl-card", { collapsed: isCollapsed.value, disabled: !__props.enabled }])
      }, [
        createBaseVNode("div", {
          class: "tpl-card-header",
          onClick: _cache[1] || (_cache[1] = ($event) => __props.enabled && toggleCollapse())
        }, [
          (openBlock(), createBlock(resolveDynamicComponent(__props.icon), { class: "w-3.5 h-3.5" })),
          createBaseVNode("span", _hoisted_1$M, toDisplayString(__props.title), 1),
          __props.count !== void 0 ? (openBlock(), createElementBlock("span", _hoisted_2$L, toDisplayString(__props.count), 1)) : createCommentVNode("", true),
          __props.toggleable ? (openBlock(), createElementBlock("button", {
            key: 1,
            onClick: _cache[0] || (_cache[0] = withModifiers(($event) => _ctx.$emit("toggle"), ["stop"])),
            class: normalizeClass(["tpl-card-toggle", { active: __props.enabled }])
          }, [
            __props.enabled ? (openBlock(), createBlock(unref(Check), {
              key: 0,
              class: "w-3 h-3"
            })) : createCommentVNode("", true)
          ], 2)) : createCommentVNode("", true),
          __props.collapsible && __props.enabled ? (openBlock(), createBlock(unref(ChevronRight), {
            key: 2,
            class: normalizeClass(["w-3.5 h-3.5 tpl-card-chevron", { rotated: !isCollapsed.value }])
          }, null, 8, ["class"])) : createCommentVNode("", true)
        ]),
        createVNode(Transition, { name: "expand" }, {
          default: withCtx(() => [
            !isCollapsed.value && __props.enabled ? (openBlock(), createElementBlock("div", _hoisted_3$I, [
              renderSlot(_ctx.$slots, "default", {}, void 0, true)
            ])) : createCommentVNode("", true)
          ]),
          _: 3
        })
      ], 2);
    };
  }
});
const TemplateCard = /* @__PURE__ */ _export_sfc(_sfc_main$O, [["__scopeId", "data-v-647cebfc"]]);
const _hoisted_1$L = { class: "item-icon" };
const _hoisted_2$K = { class: "item-name" };
const _hoisted_3$H = {
  key: 0,
  class: "current-badge"
};
const _hoisted_4$F = { class: "item-actions" };
const _hoisted_5$z = ["title"];
const _hoisted_6$w = ["title"];
const _sfc_main$N = /* @__PURE__ */ defineComponent({
  __name: "TemplateListItem",
  props: {
    template: {},
    active: { type: Boolean },
    isCurrent: { type: Boolean }
  },
  emits: ["select", "delete", "toggle-favorite"],
  setup(__props) {
    const { t } = useI18n();
    return (_ctx, _cache) => {
      return openBlock(), createElementBlock("div", {
        onClick: _cache[2] || (_cache[2] = ($event) => _ctx.$emit("select")),
        class: normalizeClass(["template-item", { active: __props.active, "is-current": __props.isCurrent }])
      }, [
        createBaseVNode("span", _hoisted_1$L, toDisplayString(__props.template.icon), 1),
        createBaseVNode("span", _hoisted_2$K, toDisplayString(__props.template.name), 1),
        __props.isCurrent ? (openBlock(), createElementBlock("span", _hoisted_3$H, "●")) : createCommentVNode("", true),
        createBaseVNode("div", _hoisted_4$F, [
          createBaseVNode("button", {
            onClick: _cache[0] || (_cache[0] = withModifiers(($event) => _ctx.$emit("toggle-favorite"), ["stop"])),
            class: normalizeClass(["item-action fav", { active: __props.template.isFavorite }]),
            title: unref(t)("templates.toggleFavorite")
          }, [
            createVNode(unref(Star), {
              class: "w-3 h-3",
              fill: __props.template.isFavorite ? "currentColor" : "none"
            }, null, 8, ["fill"])
          ], 10, _hoisted_5$z),
          !__props.template.isBuiltIn ? (openBlock(), createElementBlock("button", {
            key: 0,
            onClick: _cache[1] || (_cache[1] = withModifiers(($event) => _ctx.$emit("delete"), ["stop"])),
            class: "item-action delete",
            title: unref(t)("templates.delete")
          }, [
            createVNode(unref(Trash2), { class: "w-3 h-3" })
          ], 8, _hoisted_6$w)) : createCommentVNode("", true)
        ])
      ], 2);
    };
  }
});
const TemplateListItem = /* @__PURE__ */ _export_sfc(_sfc_main$N, [["__scopeId", "data-v-1269a54b"]]);
const _hoisted_1$K = { class: "tpl-tile-label" };
const _hoisted_2$J = {
  key: 0,
  class: "tpl-tile-hint"
};
const _sfc_main$M = /* @__PURE__ */ defineComponent({
  __name: "TemplateOptionTile",
  props: {
    modelValue: { type: Boolean },
    label: {},
    icon: {},
    hint: {}
  },
  emits: ["update:modelValue"],
  setup(__props) {
    return (_ctx, _cache) => {
      return openBlock(), createElementBlock("button", {
        class: normalizeClass(["tpl-tile", { active: __props.modelValue }]),
        onClick: _cache[0] || (_cache[0] = ($event) => _ctx.$emit("update:modelValue", !__props.modelValue))
      }, [
        (openBlock(), createBlock(resolveDynamicComponent(__props.icon), { class: "w-4 h-4" })),
        createBaseVNode("span", _hoisted_1$K, toDisplayString(__props.label), 1),
        __props.hint ? (openBlock(), createElementBlock("span", _hoisted_2$J, toDisplayString(__props.hint), 1)) : createCommentVNode("", true)
      ], 2);
    };
  }
});
const TemplateOptionTile = /* @__PURE__ */ _export_sfc(_sfc_main$M, [["__scopeId", "data-v-b07d65d8"]]);
const _hoisted_1$J = { class: "tpl-header" };
const _hoisted_2$I = { class: "tpl-header-left" };
const _hoisted_3$G = { class: "tpl-header-center" };
const _hoisted_4$E = {
  key: 0,
  class: "tpl-unsaved-indicator"
};
const _hoisted_5$y = { class: "tpl-body" };
const _hoisted_6$v = { class: "tpl-sidebar" };
const _hoisted_7$t = { class: "tpl-search" };
const _hoisted_8$r = ["placeholder"];
const _hoisted_9$o = { class: "tpl-list" };
const _hoisted_10$o = {
  key: 0,
  class: "tpl-group",
  open: ""
};
const _hoisted_11$m = {
  key: 1,
  class: "tpl-group",
  open: ""
};
const _hoisted_12$j = {
  key: 2,
  class: "tpl-group",
  open: ""
};
const _hoisted_13$j = {
  key: 3,
  class: "tpl-no-results"
};
const _hoisted_14$i = { class: "tpl-editor" };
const _hoisted_15$g = { class: "tpl-editor-header" };
const _hoisted_16$e = ["disabled"];
const _hoisted_17$e = { class: "tpl-icon-display" };
const _hoisted_18$d = { class: "tpl-name-group" };
const _hoisted_19$d = ["placeholder", "disabled"];
const _hoisted_20$d = ["placeholder", "disabled"];
const _hoisted_21$b = {
  key: 0,
  class: "tpl-active-badge"
};
const _hoisted_22$a = { class: "tpl-chips" };
const _hoisted_23$a = ["onClick", "title"];
const _hoisted_24$a = { class: "tpl-cards" };
const _hoisted_25$9 = ["placeholder"];
const _hoisted_26$9 = ["placeholder"];
const _hoisted_27$8 = { class: "tpl-options-card" };
const _hoisted_28$8 = { class: "tpl-options-header" };
const _hoisted_29$4 = { class: "tpl-options-grid" };
const _hoisted_30$2 = { class: "tpl-advanced" };
const _hoisted_31$2 = { class: "tpl-advanced-content" };
const _hoisted_32$2 = { class: "tpl-advanced-field" };
const _hoisted_33$2 = ["placeholder"];
const _hoisted_34$1 = { class: "tpl-advanced-field" };
const _hoisted_35$1 = ["placeholder"];
const _hoisted_36$1 = {
  key: 1,
  class: "tpl-empty"
};
const _hoisted_37$1 = { class: "tpl-preview" };
const _hoisted_38$1 = { class: "tpl-preview-header" };
const _hoisted_39$1 = ["innerHTML"];
const _hoisted_40 = { class: "tpl-preview-footer" };
const _hoisted_41 = { class: "tpl-token-bar" };
const _hoisted_42 = { class: "tpl-token-info" };
const _hoisted_43 = { class: "tpl-token-count" };
const _hoisted_44 = { class: "tpl-footer" };
const _hoisted_45 = { class: "tpl-footer-left" };
const _hoisted_46 = ["title"];
const _hoisted_47 = ["title", "disabled"];
const _hoisted_48 = ["disabled"];
const _hoisted_49 = { class: "tpl-footer-center" };
const _hoisted_50 = {
  key: 0,
  class: "tpl-saved-msg"
};
const _hoisted_51 = {
  key: 1,
  class: "tpl-unsaved-msg"
};
const _hoisted_52 = { class: "tpl-footer-right" };
const _hoisted_53 = { class: "tpl-autosave" };
const _hoisted_54 = ["disabled"];
const _hoisted_55 = { class: "tpl-dialog danger" };
const _hoisted_56 = { class: "tpl-dialog-btns" };
const _hoisted_57 = { class: "tpl-dialog warning" };
const _hoisted_58 = { class: "tpl-dialog-btns" };
const _sfc_main$L = /* @__PURE__ */ defineComponent({
  __name: "TemplateModal",
  setup(__props) {
    const UserIcon = shallowRef(User);
    const ListChecksIcon = shallowRef(ListChecks);
    const FolderTreeIcon = shallowRef(FolderTree);
    const HashIcon = shallowRef(Hash);
    const FileCodeIcon = shallowRef(FileCode);
    const ClipboardIcon = shallowRef(Clipboard);
    const { t } = useI18n();
    const templateStore = useTemplateStore();
    const contextStore = useContextStore();
    const { isModalOpen, builtInTemplates, customTemplates, favoriteTemplates } = storeToRefs(templateStore);
    const isOpen = computed(() => isModalOpen.value);
    const selectedTemplate = ref(null);
    const editingTemplate = ref(null);
    const originalJson = ref("");
    const deleteConfirmId = ref(null);
    const showUnsavedWarning = ref(false);
    const searchQuery = ref("");
    const autoSaveOnClose = ref(localStorage.getItem("template-autosave") === "true");
    const showEmojiPicker = ref(false);
    const justSaved = ref(false);
    const roleTextarea = ref(null);
    const rulesTextarea = ref(null);
    const sectionsList = SECTION_META;
    const hasChanges = computed(() => !editingTemplate.value ? false : !selectedTemplate.value ? true : JSON.stringify(editingTemplate.value) !== originalJson.value);
    const isEditingBuiltIn = computed(() => editingTemplate.value?.isBuiltIn ?? false);
    const canApply = computed(() => selectedTemplate.value && selectedTemplate.value.id !== templateStore.activeTemplateId);
    const animatedTokens = ref(0);
    const previewTokens = computed(() => {
      if (!editingTemplate.value) return 0;
      let count = 0;
      const tpl = editingTemplate.value;
      if (tpl.roleContent) count += Math.round(tpl.roleContent.length / 4);
      if (tpl.rulesContent) count += Math.round(tpl.rulesContent.length / 4);
      if (tpl.customPrefix) count += Math.round(tpl.customPrefix.length / 4);
      if (tpl.customSuffix) count += Math.round(tpl.customSuffix.length / 4);
      count += contextStore.tokenCount || 0;
      return count;
    });
    const tokenPercent = computed(() => Math.min(100, previewTokens.value / 32e3 * 100));
    const tokenBarClass = computed(() => {
      if (tokenPercent.value > 90) return "danger";
      if (tokenPercent.value > 70) return "warning";
      return "";
    });
    watch(previewTokens, (newVal) => {
      const start = animatedTokens.value;
      const diff = newVal - start;
      const duration = 300;
      const startTime = performance.now();
      function animate2(currentTime) {
        const elapsed = currentTime - startTime;
        const progress = Math.min(elapsed / duration, 1);
        animatedTokens.value = Math.round(start + diff * progress);
        if (progress < 1) requestAnimationFrame(animate2);
      }
      requestAnimationFrame(animate2);
    }, { immediate: true });
    const previewContent = computed(() => {
      if (!editingTemplate.value) return "";
      const tpl = editingTemplate.value;
      const parts = [];
      if (tpl.customPrefix) parts.push(tpl.customPrefix);
      const files2 = contextStore.summary?.files || [];
      const hasContext = contextStore.hasContext && files2.length > 0;
      for (const sec of tpl.sectionOrder) {
        if (!tpl.sections[sec]) continue;
        switch (sec) {
          case "role":
            if (tpl.roleContent) parts.push(`## Role
${tpl.roleContent}`);
            break;
          case "rules":
            if (tpl.rulesContent) parts.push(`## Rules
${tpl.rulesContent}`);
            break;
          case "tree":
            if (hasContext) {
              const tree = files2.slice(0, 15).map((f) => `  ${f}`).join("\n");
              parts.push(`## Project Structure
${tree}${files2.length > 15 ? `
  ... +${files2.length - 15} files` : ""}`);
            } else parts.push(`## Project Structure
[Build context to see file tree]`);
            break;
          case "stats":
            parts.push(`## Stats
- Files: ${contextStore.fileCount || 0}
- Lines: ${contextStore.lineCount || 0}
- Tokens: ~${contextStore.tokenCount || 0}`);
            break;
          case "task":
            parts.push(`## Task
[Your task description]`);
            break;
          case "files":
            if (hasContext) {
              const fileList = files2.slice(0, 5).map((f) => `- ${f}`).join("\n");
              parts.push(`## Files
${fileList}${files2.length > 5 ? `
... +${files2.length - 5} more files` : ""}

[File contents]`);
            } else parts.push(`## Files
[Build context to see files]`);
            break;
        }
      }
      if (tpl.customSuffix) parts.push(tpl.customSuffix);
      return parts.join("\n\n");
    });
    const highlightedPreview = computed(() => {
      if (!previewContent.value) return `<span class="tpl-preview-empty">${t("templates.previewEmpty")}</span>`;
      return previewContent.value.replace(/&/g, "&amp;").replace(/</g, "&lt;").replace(/>/g, "&gt;").replace(/^(## .+)$/gm, '<span class="tpl-h2">$1</span>').replace(/^(- .+)$/gm, '<span class="tpl-list">$1</span>').replace(/\[([^\]]+)\]/g, '<span class="tpl-placeholder">[$1]</span>');
    });
    const filteredFavorites = computed(() => {
      let list = favoriteTemplates.value;
      if (searchQuery.value) {
        const q = searchQuery.value.toLowerCase();
        list = list.filter((t2) => t2.name.toLowerCase().includes(q));
      }
      return list;
    });
    const filteredBuiltIn = computed(() => {
      let list = builtInTemplates.value.filter((t2) => !t2.isFavorite && !t2.isHidden);
      if (searchQuery.value) {
        const q = searchQuery.value.toLowerCase();
        list = list.filter((t2) => t2.name.toLowerCase().includes(q));
      }
      return list;
    });
    const filteredCustom = computed(() => {
      let list = customTemplates.value.filter((t2) => !t2.isHidden);
      if (searchQuery.value) {
        const q = searchQuery.value.toLowerCase();
        list = list.filter((t2) => t2.name.toLowerCase().includes(q));
      }
      return list;
    });
    const noResults = computed(() => searchQuery.value && !filteredFavorites.value.length && !filteredBuiltIn.value.length && !filteredCustom.value.length);
    watch(autoSaveOnClose, (v) => localStorage.setItem("template-autosave", v ? "true" : "false"));
    watch(isOpen, (open) => {
      if (open && templateStore.activeTemplate) selectTemplate(templateStore.activeTemplate);
    });
    function autoResize(el) {
      el.style.height = "auto";
      el.style.height = Math.max(80, el.scrollHeight) + "px";
    }
    function selectTemplate(tpl) {
      if (hasChanges.value && !isEditingBuiltIn.value) {
        showUnsavedWarning.value = true;
        return;
      }
      selectedTemplate.value = tpl;
      editingTemplate.value = JSON.parse(JSON.stringify(tpl));
      originalJson.value = JSON.stringify(tpl);
      nextTick(() => {
        if (roleTextarea.value) autoResize(roleTextarea.value);
        if (rulesTextarea.value) autoResize(rulesTextarea.value);
      });
    }
    function toggleSection(key) {
      if (editingTemplate.value) editingTemplate.value.sections[key] = !editingTemplate.value.sections[key];
    }
    function createNew() {
      editingTemplate.value = {
        ...createEmptyTemplate(),
        id: `new-${Date.now()}`,
        name: "New Template",
        icon: "✨",
        createdAt: (/* @__PURE__ */ new Date()).toISOString(),
        updatedAt: (/* @__PURE__ */ new Date()).toISOString()
      };
      selectedTemplate.value = null;
      originalJson.value = "";
    }
    function saveChanges() {
      if (!editingTemplate.value) return;
      if (selectedTemplate.value) {
        templateStore.updateTemplate(selectedTemplate.value.id, editingTemplate.value);
        selectTemplate(editingTemplate.value);
      } else {
        const id = templateStore.createTemplate(editingTemplate.value);
        const created = templateStore.templates.find((t2) => t2.id === id);
        if (created) selectTemplate(created);
      }
      justSaved.value = true;
      setTimeout(() => justSaved.value = false, 2e3);
    }
    function handleClose() {
      if (hasChanges.value && !isEditingBuiltIn.value) {
        autoSaveOnClose.value ? (saveChanges(), setTimeout(close, 100)) : showUnsavedWarning.value = true;
      } else close();
    }
    function close() {
      templateStore.closeModal();
      selectedTemplate.value = null;
      editingTemplate.value = null;
      showUnsavedWarning.value = false;
    }
    function discardAndClose() {
      showUnsavedWarning.value = false;
      close();
    }
    function saveAndClose() {
      saveChanges();
      showUnsavedWarning.value = false;
      setTimeout(close, 100);
    }
    function confirmDelete(id) {
      deleteConfirmId.value = id;
    }
    function executeDelete() {
      if (deleteConfirmId.value) {
        templateStore.deleteTemplate(deleteConfirmId.value);
        if (selectedTemplate.value?.id === deleteConfirmId.value) {
          selectedTemplate.value = null;
          editingTemplate.value = null;
        }
        deleteConfirmId.value = null;
      }
    }
    function duplicateSelected() {
      if (!selectedTemplate.value) return;
      const id = templateStore.duplicateTemplate(selectedTemplate.value.id);
      if (id) {
        const c = templateStore.templates.find((t2) => t2.id === id);
        if (c) {
          originalJson.value = "";
          selectTemplate(c);
        }
      }
    }
    function applyTemplate() {
      if (selectedTemplate.value) templateStore.setActiveTemplate(selectedTemplate.value.id);
    }
    function handleKeydown(e) {
      if (e.key === "Escape") handleClose();
      if ((e.ctrlKey || e.metaKey) && e.key === "s") {
        e.preventDefault();
        if (hasChanges.value && !isEditingBuiltIn.value) saveChanges();
      }
    }
    function handleExport() {
      if (!selectedTemplate.value) return;
      const blob = new Blob([JSON.stringify(selectedTemplate.value, null, 2)], { type: "application/json" });
      const a = document.createElement("a");
      a.href = URL.createObjectURL(blob);
      a.download = `template-${selectedTemplate.value.id}.json`;
      a.click();
    }
    function handleImport() {
      const input = document.createElement("input");
      input.type = "file";
      input.accept = ".json";
      input.onchange = async (e) => {
        const file = e.target.files?.[0];
        if (!file) return;
        try {
          const data = JSON.parse(await file.text());
          data.id = `imported-${Date.now()}`;
          data.isBuiltIn = false;
          data.sectionOrder = data.sectionOrder || DEFAULT_SECTION_ORDER;
          const id = templateStore.createTemplate(data);
          const c = templateStore.templates.find((t2) => t2.id === id);
          if (c) selectTemplate(c);
        } catch {
        }
      };
      input.click();
    }
    return (_ctx, _cache) => {
      return openBlock(), createBlock(Teleport, { to: "body" }, [
        createVNode(Transition, { name: "modal" }, {
          default: withCtx(() => [
            isOpen.value ? (openBlock(), createElementBlock("div", {
              key: 0,
              class: "tpl-overlay",
              onClick: withModifiers(handleClose, ["self"])
            }, [
              createBaseVNode("div", {
                class: "tpl-modal",
                onKeydown: handleKeydown
              }, [
                createBaseVNode("div", _hoisted_1$J, [
                  createBaseVNode("div", _hoisted_2$I, [
                    createVNode(unref(FileText), { class: "w-4 h-4 text-purple-400" }),
                    createBaseVNode("h2", null, toDisplayString(unref(t)("templates.manage")), 1)
                  ]),
                  createBaseVNode("div", _hoisted_3$G, [
                    hasChanges.value && !isEditingBuiltIn.value ? (openBlock(), createElementBlock("span", _hoisted_4$E)) : createCommentVNode("", true)
                  ]),
                  createBaseVNode("button", {
                    onClick: handleClose,
                    class: "tpl-close"
                  }, [
                    createVNode(unref(X), { class: "w-4 h-4" })
                  ])
                ]),
                createBaseVNode("div", _hoisted_5$y, [
                  createBaseVNode("aside", _hoisted_6$v, [
                    createBaseVNode("div", _hoisted_7$t, [
                      createVNode(unref(Search), { class: "w-3.5 h-3.5" }),
                      withDirectives(createBaseVNode("input", {
                        "onUpdate:modelValue": _cache[0] || (_cache[0] = ($event) => searchQuery.value = $event),
                        placeholder: unref(t)("templates.search")
                      }, null, 8, _hoisted_8$r), [
                        [vModelText, searchQuery.value]
                      ])
                    ]),
                    createBaseVNode("nav", _hoisted_9$o, [
                      filteredFavorites.value.length ? (openBlock(), createElementBlock("details", _hoisted_10$o, [
                        createBaseVNode("summary", null, [
                          createVNode(unref(Star), { class: "w-3 h-3 text-amber-400" }),
                          createTextVNode(toDisplayString(unref(t)("templates.favorites")), 1)
                        ]),
                        (openBlock(true), createElementBlock(Fragment, null, renderList(filteredFavorites.value, (tpl) => {
                          return openBlock(), createBlock(TemplateListItem, {
                            key: tpl.id,
                            template: tpl,
                            active: selectedTemplate.value?.id === tpl.id,
                            "is-current": tpl.id === unref(templateStore).activeTemplateId,
                            onSelect: ($event) => selectTemplate(tpl),
                            onDelete: ($event) => confirmDelete(tpl.id),
                            onToggleFavorite: ($event) => unref(templateStore).toggleFavorite(tpl.id)
                          }, null, 8, ["template", "active", "is-current", "onSelect", "onDelete", "onToggleFavorite"]);
                        }), 128))
                      ])) : createCommentVNode("", true),
                      filteredBuiltIn.value.length ? (openBlock(), createElementBlock("details", _hoisted_11$m, [
                        createBaseVNode("summary", null, [
                          createVNode(unref(Zap), { class: "w-3 h-3 text-blue-400" }),
                          createTextVNode(toDisplayString(unref(t)("templates.builtIn")), 1)
                        ]),
                        (openBlock(true), createElementBlock(Fragment, null, renderList(filteredBuiltIn.value, (tpl) => {
                          return openBlock(), createBlock(TemplateListItem, {
                            key: tpl.id,
                            template: tpl,
                            active: selectedTemplate.value?.id === tpl.id,
                            "is-current": tpl.id === unref(templateStore).activeTemplateId,
                            onSelect: ($event) => selectTemplate(tpl),
                            onToggleFavorite: ($event) => unref(templateStore).toggleFavorite(tpl.id)
                          }, null, 8, ["template", "active", "is-current", "onSelect", "onToggleFavorite"]);
                        }), 128))
                      ])) : createCommentVNode("", true),
                      filteredCustom.value.length ? (openBlock(), createElementBlock("details", _hoisted_12$j, [
                        createBaseVNode("summary", null, [
                          createVNode(unref(User), { class: "w-3 h-3 text-emerald-400" }),
                          createTextVNode(toDisplayString(unref(t)("templates.custom")), 1)
                        ]),
                        (openBlock(true), createElementBlock(Fragment, null, renderList(filteredCustom.value, (tpl) => {
                          return openBlock(), createBlock(TemplateListItem, {
                            key: tpl.id,
                            template: tpl,
                            active: selectedTemplate.value?.id === tpl.id,
                            "is-current": tpl.id === unref(templateStore).activeTemplateId,
                            onSelect: ($event) => selectTemplate(tpl),
                            onDelete: ($event) => confirmDelete(tpl.id),
                            onToggleFavorite: ($event) => unref(templateStore).toggleFavorite(tpl.id)
                          }, null, 8, ["template", "active", "is-current", "onSelect", "onDelete", "onToggleFavorite"]);
                        }), 128))
                      ])) : createCommentVNode("", true),
                      noResults.value ? (openBlock(), createElementBlock("div", _hoisted_13$j, toDisplayString(unref(t)("templates.noResults")), 1)) : createCommentVNode("", true)
                    ]),
                    createBaseVNode("button", {
                      onClick: createNew,
                      class: "tpl-new-btn"
                    }, [
                      createVNode(unref(Plus), { class: "w-3.5 h-3.5" }),
                      createTextVNode(toDisplayString(unref(t)("templates.create")), 1)
                    ])
                  ]),
                  createBaseVNode("main", _hoisted_14$i, [
                    editingTemplate.value ? (openBlock(), createElementBlock(Fragment, { key: 0 }, [
                      createBaseVNode("div", _hoisted_15$g, [
                        createBaseVNode("button", {
                          class: "tpl-icon-btn",
                          onClick: _cache[1] || (_cache[1] = ($event) => showEmojiPicker.value = !showEmojiPicker.value),
                          disabled: editingTemplate.value.isBuiltIn
                        }, [
                          createBaseVNode("span", _hoisted_17$e, toDisplayString(editingTemplate.value.icon), 1)
                        ], 8, _hoisted_16$e),
                        createBaseVNode("div", _hoisted_18$d, [
                          withDirectives(createBaseVNode("input", {
                            "onUpdate:modelValue": _cache[2] || (_cache[2] = ($event) => editingTemplate.value.name = $event),
                            class: "tpl-name-input",
                            placeholder: unref(t)("templates.namePlaceholder"),
                            disabled: editingTemplate.value.isBuiltIn
                          }, null, 8, _hoisted_19$d), [
                            [vModelText, editingTemplate.value.name]
                          ]),
                          withDirectives(createBaseVNode("input", {
                            "onUpdate:modelValue": _cache[3] || (_cache[3] = ($event) => editingTemplate.value.description = $event),
                            class: "tpl-desc-input",
                            placeholder: unref(t)("templates.descriptionPlaceholder"),
                            disabled: editingTemplate.value.isBuiltIn
                          }, null, 8, _hoisted_20$d), [
                            [vModelText, editingTemplate.value.description]
                          ])
                        ]),
                        editingTemplate.value.id === unref(templateStore).activeTemplateId ? (openBlock(), createElementBlock("span", _hoisted_21$b, [
                          createVNode(unref(Check), { class: "w-3 h-3" }),
                          createTextVNode(toDisplayString(unref(t)("templates.currentTemplate")), 1)
                        ])) : createCommentVNode("", true)
                      ]),
                      createBaseVNode("div", _hoisted_22$a, [
                        (openBlock(true), createElementBlock(Fragment, null, renderList(unref(sectionsList), (s) => {
                          return openBlock(), createElementBlock("button", {
                            key: s.key,
                            onClick: ($event) => toggleSection(s.key),
                            class: normalizeClass(["tpl-chip", { active: editingTemplate.value.sections[s.key] }]),
                            title: unref(t)(`templates.sectionHint.${s.key}`)
                          }, [
                            editingTemplate.value.sections[s.key] ? (openBlock(), createBlock(unref(Check), {
                              key: 0,
                              class: "w-3 h-3"
                            })) : createCommentVNode("", true),
                            createBaseVNode("span", null, toDisplayString(unref(t)(`templates.section.${s.key}`)), 1)
                          ], 10, _hoisted_23$a);
                        }), 128))
                      ]),
                      createBaseVNode("div", _hoisted_24$a, [
                        editingTemplate.value.sections.role ? (openBlock(), createBlock(TemplateCard, {
                          key: 0,
                          title: unref(t)("templates.roleContent"),
                          icon: UserIcon.value,
                          count: editingTemplate.value.roleContent?.length || 0,
                          enabled: true
                        }, {
                          default: withCtx(() => [
                            withDirectives(createBaseVNode("textarea", {
                              ref_key: "roleTextarea",
                              ref: roleTextarea,
                              "onUpdate:modelValue": _cache[4] || (_cache[4] = ($event) => editingTemplate.value.roleContent = $event),
                              placeholder: unref(t)("templates.rolePlaceholder"),
                              class: "tpl-textarea",
                              onInput: _cache[5] || (_cache[5] = ($event) => autoResize($event.target))
                            }, null, 40, _hoisted_25$9), [
                              [vModelText, editingTemplate.value.roleContent]
                            ])
                          ]),
                          _: 1
                        }, 8, ["title", "icon", "count"])) : createCommentVNode("", true),
                        editingTemplate.value.sections.rules ? (openBlock(), createBlock(TemplateCard, {
                          key: 1,
                          title: unref(t)("templates.rulesContent"),
                          icon: ListChecksIcon.value,
                          count: editingTemplate.value.rulesContent?.length || 0,
                          enabled: true
                        }, {
                          default: withCtx(() => [
                            withDirectives(createBaseVNode("textarea", {
                              ref_key: "rulesTextarea",
                              ref: rulesTextarea,
                              "onUpdate:modelValue": _cache[6] || (_cache[6] = ($event) => editingTemplate.value.rulesContent = $event),
                              placeholder: unref(t)("templates.rulesPlaceholder"),
                              class: "tpl-textarea",
                              onInput: _cache[7] || (_cache[7] = ($event) => autoResize($event.target))
                            }, null, 40, _hoisted_26$9), [
                              [vModelText, editingTemplate.value.rulesContent]
                            ])
                          ]),
                          _: 1
                        }, 8, ["title", "icon", "count"])) : createCommentVNode("", true),
                        createBaseVNode("div", _hoisted_27$8, [
                          createBaseVNode("div", _hoisted_28$8, [
                            createVNode(unref(Settings2), { class: "w-3.5 h-3.5" }),
                            createBaseVNode("span", null, toDisplayString(unref(t)("templates.contextOptions")), 1)
                          ]),
                          createBaseVNode("div", _hoisted_29$4, [
                            createVNode(TemplateOptionTile, {
                              modelValue: editingTemplate.value.sections.tree,
                              "onUpdate:modelValue": _cache[8] || (_cache[8] = ($event) => editingTemplate.value.sections.tree = $event),
                              label: unref(t)("templates.section.tree"),
                              icon: FolderTreeIcon.value,
                              hint: "Structure"
                            }, null, 8, ["modelValue", "label", "icon"]),
                            createVNode(TemplateOptionTile, {
                              modelValue: editingTemplate.value.sections.stats,
                              "onUpdate:modelValue": _cache[9] || (_cache[9] = ($event) => editingTemplate.value.sections.stats = $event),
                              label: unref(t)("templates.section.stats"),
                              icon: HashIcon.value,
                              hint: "Metrics"
                            }, null, 8, ["modelValue", "label", "icon"]),
                            createVNode(TemplateOptionTile, {
                              modelValue: editingTemplate.value.sections.files,
                              "onUpdate:modelValue": _cache[10] || (_cache[10] = ($event) => editingTemplate.value.sections.files = $event),
                              label: unref(t)("templates.section.files"),
                              icon: FileCodeIcon.value,
                              hint: "Content"
                            }, null, 8, ["modelValue", "label", "icon"]),
                            createVNode(TemplateOptionTile, {
                              modelValue: editingTemplate.value.sections.task,
                              "onUpdate:modelValue": _cache[11] || (_cache[11] = ($event) => editingTemplate.value.sections.task = $event),
                              label: unref(t)("templates.section.task"),
                              icon: ClipboardIcon.value,
                              hint: "Your task"
                            }, null, 8, ["modelValue", "label", "icon"])
                          ])
                        ]),
                        createBaseVNode("details", _hoisted_30$2, [
                          createBaseVNode("summary", null, [
                            createVNode(unref(ChevronRight), { class: "w-3.5 h-3.5 tpl-chevron" }),
                            createTextVNode(toDisplayString(unref(t)("templates.additional")), 1)
                          ]),
                          createBaseVNode("div", _hoisted_31$2, [
                            createBaseVNode("div", _hoisted_32$2, [
                              createBaseVNode("label", null, toDisplayString(unref(t)("templates.prefix")), 1),
                              withDirectives(createBaseVNode("textarea", {
                                "onUpdate:modelValue": _cache[12] || (_cache[12] = ($event) => editingTemplate.value.customPrefix = $event),
                                placeholder: unref(t)("templates.prefixPlaceholder"),
                                rows: "2"
                              }, null, 8, _hoisted_33$2), [
                                [vModelText, editingTemplate.value.customPrefix]
                              ])
                            ]),
                            createBaseVNode("div", _hoisted_34$1, [
                              createBaseVNode("label", null, toDisplayString(unref(t)("templates.suffix")), 1),
                              withDirectives(createBaseVNode("textarea", {
                                "onUpdate:modelValue": _cache[13] || (_cache[13] = ($event) => editingTemplate.value.customSuffix = $event),
                                placeholder: unref(t)("templates.suffixPlaceholder"),
                                rows: "2"
                              }, null, 8, _hoisted_35$1), [
                                [vModelText, editingTemplate.value.customSuffix]
                              ])
                            ])
                          ])
                        ])
                      ])
                    ], 64)) : (openBlock(), createElementBlock("div", _hoisted_36$1, [
                      createVNode(unref(FileText), { class: "w-10 h-10 opacity-15" }),
                      createBaseVNode("p", null, toDisplayString(unref(t)("templates.selectToEdit")), 1)
                    ]))
                  ]),
                  createBaseVNode("aside", _hoisted_37$1, [
                    createBaseVNode("div", _hoisted_38$1, [
                      createVNode(unref(Eye), { class: "w-3.5 h-3.5" }),
                      createBaseVNode("span", null, toDisplayString(unref(t)("templates.preview")), 1)
                    ]),
                    createBaseVNode("div", {
                      class: "tpl-preview-content",
                      innerHTML: highlightedPreview.value
                    }, null, 8, _hoisted_39$1),
                    createBaseVNode("div", _hoisted_40, [
                      createBaseVNode("div", _hoisted_41, [
                        createBaseVNode("div", {
                          class: normalizeClass(["tpl-token-fill", tokenBarClass.value]),
                          style: normalizeStyle({ width: tokenPercent.value + "%" })
                        }, null, 6)
                      ]),
                      createBaseVNode("div", _hoisted_42, [
                        createBaseVNode("span", _hoisted_43, toDisplayString(animatedTokens.value.toLocaleString()), 1),
                        _cache[18] || (_cache[18] = createBaseVNode("span", { class: "tpl-token-label" }, "/ 32k tokens", -1))
                      ])
                    ])
                  ])
                ]),
                createBaseVNode("div", _hoisted_44, [
                  createBaseVNode("div", _hoisted_45, [
                    createBaseVNode("button", {
                      onClick: handleImport,
                      class: "tpl-footer-btn",
                      title: unref(t)("templates.import")
                    }, [
                      createVNode(unref(Upload), { class: "w-3.5 h-3.5" })
                    ], 8, _hoisted_46),
                    createBaseVNode("button", {
                      onClick: handleExport,
                      class: "tpl-footer-btn",
                      title: unref(t)("templates.export"),
                      disabled: !selectedTemplate.value
                    }, [
                      createVNode(unref(Download), { class: "w-3.5 h-3.5" })
                    ], 8, _hoisted_47),
                    createBaseVNode("button", {
                      onClick: duplicateSelected,
                      class: "tpl-footer-btn",
                      disabled: !selectedTemplate.value
                    }, [
                      createVNode(unref(Copy), { class: "w-3.5 h-3.5" }),
                      createBaseVNode("span", null, toDisplayString(unref(t)("templates.duplicate")), 1)
                    ], 8, _hoisted_48)
                  ]),
                  createBaseVNode("div", _hoisted_49, [
                    createVNode(Transition, {
                      name: "fade",
                      mode: "out-in"
                    }, {
                      default: withCtx(() => [
                        justSaved.value ? (openBlock(), createElementBlock("span", _hoisted_50, [
                          createVNode(unref(Check), { class: "w-3 h-3" }),
                          createTextVNode(toDisplayString(unref(t)("templates.saved")), 1)
                        ])) : hasChanges.value && !isEditingBuiltIn.value ? (openBlock(), createElementBlock("span", _hoisted_51, toDisplayString(unref(t)("templates.unsavedChanges")), 1)) : createCommentVNode("", true)
                      ]),
                      _: 1
                    })
                  ]),
                  createBaseVNode("div", _hoisted_52, [
                    createBaseVNode("label", _hoisted_53, [
                      withDirectives(createBaseVNode("input", {
                        type: "checkbox",
                        "onUpdate:modelValue": _cache[14] || (_cache[14] = ($event) => autoSaveOnClose.value = $event)
                      }, null, 512), [
                        [vModelCheckbox, autoSaveOnClose.value]
                      ]),
                      _cache[19] || (_cache[19] = createBaseVNode("span", { class: "tpl-autosave-track" }, [
                        createBaseVNode("span", { class: "tpl-autosave-thumb" })
                      ], -1)),
                      createBaseVNode("span", null, toDisplayString(unref(t)("templates.autosave")), 1)
                    ]),
                    canApply.value ? (openBlock(), createElementBlock("button", {
                      key: 0,
                      onClick: applyTemplate,
                      class: "tpl-apply-btn"
                    }, [
                      createVNode(unref(Sparkles), { class: "w-3.5 h-3.5" }),
                      createTextVNode(toDisplayString(unref(t)("templates.apply")), 1)
                    ])) : createCommentVNode("", true),
                    createBaseVNode("button", {
                      onClick: saveChanges,
                      class: normalizeClass(["tpl-save-btn", { pulse: hasChanges.value && !isEditingBuiltIn.value }]),
                      disabled: !hasChanges.value || isEditingBuiltIn.value
                    }, [
                      createVNode(unref(Save), { class: "w-3.5 h-3.5" }),
                      createTextVNode(toDisplayString(unref(t)("common.save")), 1)
                    ], 10, _hoisted_54)
                  ])
                ])
              ], 32),
              createVNode(Transition, { name: "fade" }, {
                default: withCtx(() => [
                  deleteConfirmId.value ? (openBlock(), createElementBlock("div", {
                    key: 0,
                    class: "tpl-dialog-overlay",
                    onClick: _cache[16] || (_cache[16] = withModifiers(($event) => deleteConfirmId.value = null, ["self"]))
                  }, [
                    createBaseVNode("div", _hoisted_55, [
                      createVNode(unref(Trash2), { class: "w-6 h-6" }),
                      createBaseVNode("h3", null, toDisplayString(unref(t)("templates.deleteConfirm")), 1),
                      createBaseVNode("p", null, toDisplayString(unref(t)("templates.deleteConfirmText")), 1),
                      createBaseVNode("div", _hoisted_56, [
                        createBaseVNode("button", {
                          onClick: _cache[15] || (_cache[15] = ($event) => deleteConfirmId.value = null),
                          class: "tpl-cancel-btn"
                        }, toDisplayString(unref(t)("common.cancel")), 1),
                        createBaseVNode("button", {
                          onClick: executeDelete,
                          class: "tpl-danger-btn"
                        }, toDisplayString(unref(t)("templates.delete")), 1)
                      ])
                    ])
                  ])) : createCommentVNode("", true)
                ]),
                _: 1
              }),
              createVNode(Transition, { name: "fade" }, {
                default: withCtx(() => [
                  showUnsavedWarning.value ? (openBlock(), createElementBlock("div", {
                    key: 0,
                    class: "tpl-dialog-overlay",
                    onClick: _cache[17] || (_cache[17] = withModifiers(($event) => showUnsavedWarning.value = false, ["self"]))
                  }, [
                    createBaseVNode("div", _hoisted_57, [
                      createVNode(unref(TriangleAlert), { class: "w-6 h-6" }),
                      createBaseVNode("h3", null, toDisplayString(unref(t)("templates.unsavedChanges")), 1),
                      createBaseVNode("p", null, toDisplayString(unref(t)("templates.unsavedChangesText")), 1),
                      createBaseVNode("div", _hoisted_58, [
                        createBaseVNode("button", {
                          onClick: discardAndClose,
                          class: "tpl-cancel-btn"
                        }, toDisplayString(unref(t)("templates.discard")), 1),
                        createBaseVNode("button", {
                          onClick: saveAndClose,
                          class: "tpl-save-btn"
                        }, toDisplayString(unref(t)("common.save")), 1)
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
const TemplateModal = /* @__PURE__ */ _export_sfc(_sfc_main$L, [["__scopeId", "data-v-33f238b3"]]);
const _hoisted_1$I = { class: "tpl-block-icon" };
const _hoisted_2$H = { class: "tpl-block-name" };
const _hoisted_3$F = { class: "tpl-block-tokens" };
const _hoisted_4$D = { class: "tpl-block-preview" };
const _hoisted_5$x = { class: "tpl-block-actions" };
const _hoisted_6$u = ["title"];
const _hoisted_7$s = {
  key: 0,
  class: "tpl-block-content"
};
const _hoisted_8$q = {
  key: 0,
  class: "tpl-block-section"
};
const _hoisted_9$n = { class: "tpl-block-text" };
const _hoisted_10$n = {
  key: 1,
  class: "tpl-block-section"
};
const _hoisted_11$l = { class: "tpl-block-text" };
const _hoisted_12$i = {
  key: 2,
  class: "tpl-block-section"
};
const _hoisted_13$i = { class: "tpl-block-meta" };
const _hoisted_14$h = {
  key: 3,
  class: "tpl-block-section"
};
const _hoisted_15$f = { class: "tpl-block-meta" };
const _sfc_main$K = /* @__PURE__ */ defineComponent({
  __name: "TemplatePreviewBlock",
  setup(__props) {
    const { t } = useI18n();
    const templateStore = useTemplateStore();
    const contextStore = useContextStore();
    const settingsStore = useSettingsStore();
    const isExpanded = ref(false);
    const roleColor = computed(() => {
      const id = templateStore.activeTemplate?.id || "";
      const colors = {
        architect: "#f59e0b",
        // amber
        review: "#a855f7",
        // purple
        implement: "#3b82f6",
        // blue
        explain: "#10b981"
        // emerald
      };
      return colors[id] || "#6366f1";
    });
    const tokenCount = computed(() => {
      const tpl = templateStore.activeTemplate;
      if (!tpl) return 0;
      let count = 0;
      if (tpl.roleContent) count += Math.round(tpl.roleContent.length / 4);
      if (tpl.rulesContent) count += Math.round(tpl.rulesContent.length / 4);
      count += contextStore.tokenCount;
      return count;
    });
    const truncatedRole = computed(() => {
      const role = templateStore.activeTemplate?.roleContent || "";
      return role.length > 50 ? role.slice(0, 50) + "..." : role;
    });
    function toggleExpand() {
      isExpanded.value = !isExpanded.value;
    }
    function openEditor() {
      templateStore.openModal();
    }
    return (_ctx, _cache) => {
      return unref(templateStore).activeTemplate && unref(contextStore).hasContext && unref(settingsStore).settings.context.applyTemplateOnCopy ? (openBlock(), createElementBlock("div", {
        key: 0,
        class: normalizeClass(["tpl-block", { expanded: isExpanded.value }]),
        style: normalizeStyle({ "--role-color": roleColor.value })
      }, [
        createBaseVNode("div", {
          class: "tpl-block-header",
          onClick: toggleExpand
        }, [
          createBaseVNode("span", _hoisted_1$I, toDisplayString(unref(templateStore).activeTemplate.icon), 1),
          createBaseVNode("span", _hoisted_2$H, toDisplayString(unref(templateStore).activeTemplate.name), 1),
          createBaseVNode("span", _hoisted_3$F, toDisplayString(tokenCount.value.toLocaleString()) + " tok", 1),
          _cache[0] || (_cache[0] = createBaseVNode("span", { class: "tpl-block-separator" }, "|", -1)),
          createBaseVNode("span", _hoisted_4$D, toDisplayString(truncatedRole.value), 1),
          createBaseVNode("div", _hoisted_5$x, [
            createBaseVNode("button", {
              onClick: withModifiers(openEditor, ["stop"]),
              class: "tpl-block-btn",
              title: unref(t)("templates.settings")
            }, [
              createVNode(unref(Pencil), { class: "w-3 h-3" })
            ], 8, _hoisted_6$u)
          ]),
          createVNode(unref(ChevronDown), {
            class: normalizeClass(["tpl-block-chevron", { rotated: isExpanded.value }])
          }, null, 8, ["class"])
        ]),
        createVNode(Transition, { name: "slide" }, {
          default: withCtx(() => [
            isExpanded.value ? (openBlock(), createElementBlock("div", _hoisted_7$s, [
              unref(templateStore).activeTemplate.sections.role && unref(templateStore).activeTemplate.roleContent ? (openBlock(), createElementBlock("div", _hoisted_8$q, [
                _cache[1] || (_cache[1] = createBaseVNode("span", { class: "tpl-block-label" }, "ROLE", -1)),
                createBaseVNode("pre", _hoisted_9$n, toDisplayString(unref(templateStore).activeTemplate.roleContent), 1)
              ])) : createCommentVNode("", true),
              unref(templateStore).activeTemplate.sections.rules && unref(templateStore).activeTemplate.rulesContent ? (openBlock(), createElementBlock("div", _hoisted_10$n, [
                _cache[2] || (_cache[2] = createBaseVNode("span", { class: "tpl-block-label" }, "RULES", -1)),
                createBaseVNode("pre", _hoisted_11$l, toDisplayString(unref(templateStore).activeTemplate.rulesContent), 1)
              ])) : createCommentVNode("", true),
              unref(templateStore).activeTemplate.sections.tree ? (openBlock(), createElementBlock("div", _hoisted_12$i, [
                _cache[3] || (_cache[3] = createBaseVNode("span", { class: "tpl-block-label" }, "TREE", -1)),
                createBaseVNode("span", _hoisted_13$i, toDisplayString(unref(contextStore).fileCount) + " files", 1)
              ])) : createCommentVNode("", true),
              unref(templateStore).activeTemplate.sections.stats ? (openBlock(), createElementBlock("div", _hoisted_14$h, [
                _cache[4] || (_cache[4] = createBaseVNode("span", { class: "tpl-block-label" }, "STATS", -1)),
                createBaseVNode("span", _hoisted_15$f, toDisplayString(unref(contextStore).tokenCount) + " tokens", 1)
              ])) : createCommentVNode("", true)
            ])) : createCommentVNode("", true)
          ]),
          _: 1
        })
      ], 6)) : createCommentVNode("", true);
    };
  }
});
const TemplatePreviewBlock = /* @__PURE__ */ _export_sfc(_sfc_main$K, [["__scopeId", "data-v-96774b3c"]]);
function useChunking() {
  const settingsStore = useSettingsStore();
  const currentChunkIndex = ref(0);
  const chunks = shallowRef([]);
  const copiedChunks = ref(/* @__PURE__ */ new Set());
  const isEnabled2 = computed(() => settingsStore.settings.context.enableAutoSplit);
  const maxTokensPerChunk = computed(() => settingsStore.settings.context.maxTokensPerChunk);
  const strategy = computed(() => settingsStore.settings.context.splitStrategy);
  const totalChunks = computed(() => chunks.value.length);
  const currentChunk = computed(() => currentChunkIndex.value + 1);
  const hasMultipleChunks = computed(() => totalChunks.value > 1);
  const allChunksCopied = computed(() => copiedChunks.value.size >= totalChunks.value);
  const currentChunkInfo = computed(
    () => chunks.value[currentChunkIndex.value] ?? null
  );
  function calculateChunks(lines, fileMarkers, tokensPerLine) {
    if (!isEnabled2.value || lines.length === 0) {
      return [{
        index: 0,
        startLine: 0,
        endLine: lines.length - 1,
        tokenCount: tokensPerLine.reduce((a, b) => a + b, 0)
      }];
    }
    const limit = maxTokensPerChunk.value;
    const result = [];
    if (strategy.value === "smart" || strategy.value === "file") {
      result.push(...calculateSmartChunks(lines, fileMarkers, tokensPerLine, limit));
    } else {
      result.push(...calculateHardChunks(tokensPerLine, limit));
    }
    return result;
  }
  function calculateSmartChunks(lines, fileMarkers, tokensPerLine, limit) {
    const result = [];
    let chunkStart = 0;
    let chunkTokens = 0;
    let chunkIndex = 0;
    const markers = [...fileMarkers, lines.length];
    for (let i = 0; i < markers.length - 1; i++) {
      const fileStart = markers[i];
      const fileEnd = markers[i + 1] - 1;
      let fileTokens = 0;
      for (let j = fileStart; j <= fileEnd; j++) {
        fileTokens += tokensPerLine[j] || 0;
      }
      if (chunkTokens + fileTokens > limit && chunkTokens > 0) {
        result.push({
          index: chunkIndex++,
          startLine: chunkStart,
          endLine: fileStart - 1,
          tokenCount: chunkTokens
        });
        chunkStart = fileStart;
        chunkTokens = 0;
      }
      chunkTokens += fileTokens;
    }
    if (chunkTokens > 0 || result.length === 0) {
      result.push({
        index: chunkIndex,
        startLine: chunkStart,
        endLine: lines.length - 1,
        tokenCount: chunkTokens
      });
    }
    return result;
  }
  function calculateHardChunks(tokensPerLine, limit) {
    const result = [];
    let chunkStart = 0;
    let chunkTokens = 0;
    let chunkIndex = 0;
    for (let i = 0; i < tokensPerLine.length; i++) {
      const lineTokens = tokensPerLine[i] || 0;
      if (chunkTokens + lineTokens > limit && chunkTokens > 0) {
        result.push({
          index: chunkIndex++,
          startLine: chunkStart,
          endLine: i - 1,
          tokenCount: chunkTokens
        });
        chunkStart = i;
        chunkTokens = 0;
      }
      chunkTokens += lineTokens;
    }
    if (chunkTokens > 0 || result.length === 0) {
      result.push({
        index: chunkIndex,
        startLine: chunkStart,
        endLine: tokensPerLine.length - 1,
        tokenCount: chunkTokens
      });
    }
    return result;
  }
  function setChunks(newChunks) {
    chunks.value = newChunks;
    currentChunkIndex.value = 0;
    copiedChunks.value = /* @__PURE__ */ new Set();
  }
  function goToChunk(index) {
    if (index >= 0 && index < totalChunks.value) {
      currentChunkIndex.value = index;
    }
  }
  function nextChunk() {
    if (currentChunkIndex.value < totalChunks.value - 1) {
      currentChunkIndex.value++;
    }
  }
  function prevChunk() {
    if (currentChunkIndex.value > 0) {
      currentChunkIndex.value--;
    }
  }
  function markCopied(index, autoAdvance = true) {
    copiedChunks.value.add(index);
    if (autoAdvance && index < totalChunks.value - 1) {
      setTimeout(() => {
        currentChunkIndex.value = index + 1;
      }, 500);
    }
  }
  function resetCopyState() {
    copiedChunks.value = /* @__PURE__ */ new Set();
    currentChunkIndex.value = 0;
  }
  function isChunkCopied(index) {
    return copiedChunks.value.has(index);
  }
  function getChunkBoundaries() {
    if (chunks.value.length <= 1) return [];
    return chunks.value.slice(0, -1).map((chunk) => chunk.endLine);
  }
  return {
    // State
    currentChunkIndex,
    chunks,
    copiedChunks,
    // Computed
    isEnabled: isEnabled2,
    totalChunks,
    currentChunk,
    hasMultipleChunks,
    allChunksCopied,
    currentChunkInfo,
    // Methods
    calculateChunks,
    setChunks,
    goToChunk,
    nextChunk,
    prevChunk,
    markCopied,
    resetCopyState,
    isChunkCopied,
    getChunkBoundaries
  };
}
const _hoisted_1$H = ["title"];
const _hoisted_2$G = { class: "format-label" };
const _hoisted_3$E = ["onClick", "title"];
const _hoisted_4$C = { class: "format-option-label" };
const _hoisted_5$w = {
  key: 0,
  class: "format-check",
  fill: "none",
  stroke: "currentColor",
  viewBox: "0 0 24 24"
};
const _sfc_main$J = /* @__PURE__ */ defineComponent({
  __name: "FormatDropdown",
  props: {
    modelValue: {}
  },
  emits: ["update:modelValue"],
  setup(__props, { emit: __emit }) {
    const { t } = useI18n();
    const props = __props;
    const emit = __emit;
    const isOpen = ref(false);
    const dropdownRef = ref(null);
    const menuStyle = ref({ top: "0px", left: "0px" });
    const formats = [
      { id: "xml", label: "XML", tooltip: "Best for AI - structured, easy to parse" },
      { id: "markdown", label: "Markdown", tooltip: "Good for documentation and readability" },
      { id: "plain", label: "Plain", tooltip: "Simple format with separators" }
    ];
    const currentFormat = computed(
      () => formats.find((f) => f.id === props.modelValue) || formats[0]
    );
    function selectFormat(id) {
      emit("update:modelValue", id);
      isOpen.value = false;
    }
    function updateMenuPosition() {
      if (!dropdownRef.value) return;
      const rect = dropdownRef.value.getBoundingClientRect();
      menuStyle.value = {
        top: `${rect.bottom + 4}px`,
        left: `${rect.left + rect.width / 2}px`
      };
    }
    function toggleDropdown() {
      if (!isOpen.value) {
        updateMenuPosition();
      }
      isOpen.value = !isOpen.value;
    }
    function handleClickOutside(e) {
      if (dropdownRef.value && !dropdownRef.value.contains(e.target)) {
        isOpen.value = false;
      }
    }
    onMounted(() => {
      document.addEventListener("click", handleClickOutside);
    });
    onUnmounted(() => {
      document.removeEventListener("click", handleClickOutside);
    });
    return (_ctx, _cache) => {
      return openBlock(), createElementBlock("div", {
        class: "format-dropdown",
        ref_key: "dropdownRef",
        ref: dropdownRef
      }, [
        createBaseVNode("button", {
          onClick: toggleDropdown,
          class: normalizeClass(["format-trigger", { "format-trigger-open": isOpen.value }]),
          title: unref(t)("context.changeFormat")
        }, [
          createBaseVNode("span", _hoisted_2$G, toDisplayString(currentFormat.value.label), 1),
          (openBlock(), createElementBlock("svg", {
            class: normalizeClass(["format-chevron", { "format-chevron-open": isOpen.value }]),
            fill: "none",
            stroke: "currentColor",
            viewBox: "0 0 24 24"
          }, [..._cache[0] || (_cache[0] = [
            createBaseVNode("path", {
              "stroke-linecap": "round",
              "stroke-linejoin": "round",
              "stroke-width": "2",
              d: "M19 9l-7 7-7-7"
            }, null, -1)
          ])], 2))
        ], 10, _hoisted_1$H),
        (openBlock(), createBlock(Teleport, { to: "body" }, [
          createVNode(Transition, { name: "dropdown" }, {
            default: withCtx(() => [
              isOpen.value ? (openBlock(), createElementBlock("div", {
                key: 0,
                class: "format-menu",
                style: normalizeStyle(menuStyle.value)
              }, [
                (openBlock(), createElementBlock(Fragment, null, renderList(formats, (format) => {
                  return createBaseVNode("button", {
                    key: format.id,
                    onClick: ($event) => selectFormat(format.id),
                    class: normalizeClass(["format-option", { "format-option-active": __props.modelValue === format.id }]),
                    title: format.tooltip
                  }, [
                    createBaseVNode("span", _hoisted_4$C, toDisplayString(format.label), 1),
                    __props.modelValue === format.id ? (openBlock(), createElementBlock("svg", _hoisted_5$w, [..._cache[1] || (_cache[1] = [
                      createBaseVNode("path", {
                        "stroke-linecap": "round",
                        "stroke-linejoin": "round",
                        "stroke-width": "2",
                        d: "M5 13l4 4L19 7"
                      }, null, -1)
                    ])])) : createCommentVNode("", true)
                  ], 10, _hoisted_3$E);
                }), 64))
              ], 4)) : createCommentVNode("", true)
            ]),
            _: 1
          })
        ]))
      ], 512);
    };
  }
});
const FormatDropdown = /* @__PURE__ */ _export_sfc(_sfc_main$J, [["__scopeId", "data-v-05019476"]]);
const _hoisted_1$G = { class: "panel-header-unified" };
const _hoisted_2$F = { class: "panel-header-unified-title" };
const _hoisted_3$D = {
  key: 0,
  class: "flex items-center gap-1.5 ml-2 overflow-x-auto"
};
const _hoisted_4$B = ["title"];
const _hoisted_5$v = ["title"];
const _hoisted_6$t = { class: "chip-unified chip-unified-accent" };
const _hoisted_7$r = { class: "chip-unified chip-unified-accent" };
const _hoisted_8$p = { class: "chip-unified chip-unified-accent" };
const _sfc_main$I = /* @__PURE__ */ defineComponent({
  __name: "ContextPanelHeader",
  props: {
    hasContext: { type: Boolean },
    showSearch: { type: Boolean },
    fileCount: {},
    lineCount: {},
    tokenCount: {},
    outputFormat: {}
  },
  emits: ["toggle-search", "show-stats", "format-change"],
  setup(__props) {
    const { t } = useI18n();
    return (_ctx, _cache) => {
      return openBlock(), createElementBlock("div", _hoisted_1$G, [
        createBaseVNode("div", _hoisted_2$F, [
          _cache[3] || (_cache[3] = createBaseVNode("div", { class: "section-icon section-icon-indigo" }, [
            createBaseVNode("svg", {
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
            ])
          ], -1)),
          createBaseVNode("span", null, toDisplayString(unref(t)("context.preview")), 1)
        ]),
        __props.hasContext ? (openBlock(), createElementBlock("div", _hoisted_3$D, [
          createBaseVNode("button", {
            onClick: _cache[0] || (_cache[0] = ($event) => _ctx.$emit("toggle-search")),
            class: normalizeClass(["icon-btn", { "text-indigo-400 bg-indigo-500/10": __props.showSearch }]),
            title: unref(t)("context.search")
          }, [..._cache[4] || (_cache[4] = [
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
                d: "M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"
              })
            ], -1)
          ])], 10, _hoisted_4$B),
          createBaseVNode("button", {
            onClick: _cache[1] || (_cache[1] = ($event) => _ctx.$emit("show-stats")),
            class: "hidden xl:flex items-center gap-1.5 flex-shrink-0 stats-chips-btn",
            title: unref(t)("stats.clickToExpand")
          }, [
            createBaseVNode("span", _hoisted_6$t, toDisplayString(__props.fileCount) + " " + toDisplayString(unref(t)("context.files")), 1),
            createBaseVNode("span", _hoisted_7$r, toDisplayString(__props.lineCount) + " " + toDisplayString(unref(t)("context.lines")), 1),
            createBaseVNode("span", _hoisted_8$p, toDisplayString(__props.tokenCount) + " " + toDisplayString(unref(t)("context.tokens")), 1)
          ], 8, _hoisted_5$v),
          createVNode(FormatDropdown, {
            "model-value": __props.outputFormat,
            "onUpdate:modelValue": _cache[2] || (_cache[2] = ($event) => _ctx.$emit("format-change", $event))
          }, null, 8, ["model-value"])
        ])) : createCommentVNode("", true)
      ]);
    };
  }
});
const ContextPanelHeader = /* @__PURE__ */ _export_sfc(_sfc_main$I, [["__scopeId", "data-v-49ed8095"]]);
const _hoisted_1$F = {
  key: 0,
  class: "search-bar"
};
const _hoisted_2$E = { class: "relative" };
const _hoisted_3$C = ["placeholder"];
const _hoisted_4$A = {
  key: 0,
  class: "search-nav"
};
const _hoisted_5$u = { class: "search-count" };
const _hoisted_6$s = {
  key: 1,
  class: "search-no-results"
};
const _sfc_main$H = /* @__PURE__ */ defineComponent({
  __name: "ContextPanelToolbar",
  props: {
    visible: { type: Boolean },
    searchQuery: {},
    resultsCount: {},
    currentIndex: {}
  },
  emits: ["update:searchQuery", "search-next", "search-prev", "close"],
  setup(__props, { expose: __expose, emit: __emit }) {
    const props = __props;
    const emit = __emit;
    const { t } = useI18n();
    const searchInputRef = ref(null);
    const localQuery = ref(props.searchQuery);
    watch(() => props.searchQuery, (val) => {
      localQuery.value = val;
    });
    watch(localQuery, (val) => {
      emit("update:searchQuery", val);
    });
    watch(() => props.visible, (show) => {
      if (show) {
        nextTick(() => searchInputRef.value?.focus());
      }
    });
    __expose({
      focus: () => searchInputRef.value?.focus()
    });
    return (_ctx, _cache) => {
      return __props.visible ? (openBlock(), createElementBlock("div", _hoisted_1$F, [
        createBaseVNode("div", _hoisted_2$E, [
          _cache[7] || (_cache[7] = createBaseVNode("svg", {
            class: "search-icon",
            fill: "none",
            stroke: "currentColor",
            viewBox: "0 0 24 24"
          }, [
            createBaseVNode("path", {
              "stroke-linecap": "round",
              "stroke-linejoin": "round",
              "stroke-width": "2",
              d: "M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"
            })
          ], -1)),
          withDirectives(createBaseVNode("input", {
            ref_key: "searchInputRef",
            ref: searchInputRef,
            "onUpdate:modelValue": _cache[0] || (_cache[0] = ($event) => localQuery.value = $event),
            type: "text",
            placeholder: unref(t)("context.search"),
            class: "input pl-8 pr-20 text-sm",
            onKeyup: [
              _cache[1] || (_cache[1] = withKeys(($event) => _ctx.$emit("search-next"), ["enter"])),
              _cache[2] || (_cache[2] = withKeys(($event) => _ctx.$emit("close"), ["escape"]))
            ]
          }, null, 40, _hoisted_3$C), [
            [vModelText, localQuery.value]
          ]),
          localQuery.value && __props.resultsCount > 0 ? (openBlock(), createElementBlock("div", _hoisted_4$A, [
            createBaseVNode("span", _hoisted_5$u, toDisplayString(__props.currentIndex + 1) + "/" + toDisplayString(__props.resultsCount), 1),
            createBaseVNode("button", {
              onClick: _cache[3] || (_cache[3] = ($event) => _ctx.$emit("search-prev")),
              class: "search-nav-btn"
            }, [..._cache[5] || (_cache[5] = [
              createBaseVNode("svg", {
                class: "w-3 h-3",
                fill: "none",
                stroke: "currentColor",
                viewBox: "0 0 24 24"
              }, [
                createBaseVNode("path", {
                  "stroke-linecap": "round",
                  "stroke-linejoin": "round",
                  "stroke-width": "2",
                  d: "M5 15l7-7 7 7"
                })
              ], -1)
            ])]),
            createBaseVNode("button", {
              onClick: _cache[4] || (_cache[4] = ($event) => _ctx.$emit("search-next")),
              class: "search-nav-btn"
            }, [..._cache[6] || (_cache[6] = [
              createBaseVNode("svg", {
                class: "w-3 h-3",
                fill: "none",
                stroke: "currentColor",
                viewBox: "0 0 24 24"
              }, [
                createBaseVNode("path", {
                  "stroke-linecap": "round",
                  "stroke-linejoin": "round",
                  "stroke-width": "2",
                  d: "M19 9l-7 7-7-7"
                })
              ], -1)
            ])])
          ])) : localQuery.value && __props.resultsCount === 0 ? (openBlock(), createElementBlock("div", _hoisted_6$s, [
            createBaseVNode("span", null, toDisplayString(unref(t)("context.noResults")), 1)
          ])) : createCommentVNode("", true)
        ])
      ])) : createCommentVNode("", true);
    };
  }
});
const ContextPanelToolbar = /* @__PURE__ */ _export_sfc(_sfc_main$H, [["__scopeId", "data-v-9fca6cef"]]);
const logger$6 = useLogger("Export");
function useExport() {
  const contextStore = useContextStore();
  const settingsStore = useSettingsStore();
  const isOpen = ref(false);
  const isExporting = ref(false);
  const selectedMode = ref("clipboard");
  const exportResult = ref(null);
  const error = ref(null);
  const settings2 = computed(() => ({
    // From settingsStore.context
    exportFormat: settingsStore.settings.context.outputFormat === "xml" ? "manifest" : "plain",
    stripComments: settingsStore.settings.context.stripComments,
    includeManifest: settingsStore.settings.context.includeManifest,
    tokenLimit: settingsStore.settings.context.maxTokens,
    enableAutoSplit: settingsStore.settings.context.enableAutoSplit,
    maxTokensPerChunk: settingsStore.settings.context.maxTokensPerChunk,
    includeLineNumbers: settingsStore.settings.context.includeLineNumbers,
    // Fixed values
    aiProfile: settingsStore.settings.aiModel,
    overlapTokens: 200,
    splitStrategy: settingsStore.settings.context.splitStrategy,
    theme: "default",
    includePageNumbers: true
  }));
  function open() {
    if (!contextStore.hasContext) {
      error.value = "Сначала соберите контекст";
      return false;
    }
    isOpen.value = true;
    exportResult.value = null;
    error.value = null;
    return true;
  }
  function close() {
    isOpen.value = false;
  }
  async function executeExport() {
    if (!contextStore.hasContext) {
      error.value = "Контекст не найден";
      return false;
    }
    isExporting.value = true;
    error.value = null;
    try {
      const contextContent = await contextStore.getFullContextContent();
      const s = settings2.value;
      const exportSettingsJson = {
        mode: selectedMode.value,
        context: contextContent,
        // Clipboard settings
        stripComments: s.stripComments,
        includeManifest: s.includeManifest,
        exportFormat: s.exportFormat,
        // AI settings
        aiProfile: s.aiProfile,
        tokenLimit: s.tokenLimit,
        fileSizeLimitKB: 5120,
        // 5 MB
        enableAutoSplit: s.enableAutoSplit,
        maxTokensPerChunk: s.maxTokensPerChunk,
        overlapTokens: s.overlapTokens,
        splitStrategy: s.splitStrategy,
        // Human settings
        theme: s.theme,
        includeLineNumbers: s.includeLineNumbers,
        includePageNumbers: s.includePageNumbers
      };
      const result = await apiService.exportContext(exportSettingsJson);
      exportResult.value = result;
      if (result.mode === "clipboard" && result.text) {
        await navigator.clipboard.writeText(result.text);
        logger$6.debug("Content copied to clipboard");
      } else if (result.filePath) {
        logger$6.debug("File exported to:", result.filePath);
      } else if (result.dataBase64) {
        downloadBase64File(result.dataBase64, result.fileName || "export.zip");
      }
      setTimeout(() => {
        close();
      }, 1500);
      return true;
    } catch (err) {
      error.value = err instanceof Error ? err.message : "Ошибка экспорта";
      logger$6.error("Export failed:", err);
      return false;
    } finally {
      isExporting.value = false;
    }
  }
  function downloadBase64File(base64, filename) {
    const byteCharacters = atob(base64);
    const byteNumbers = new Array(byteCharacters.length);
    for (let i = 0; i < byteCharacters.length; i++) {
      byteNumbers[i] = byteCharacters.charCodeAt(i);
    }
    const byteArray = new Uint8Array(byteNumbers);
    const blob = new Blob([byteArray]);
    const url = URL.createObjectURL(blob);
    const link = document.createElement("a");
    link.href = url;
    link.download = filename;
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);
    URL.revokeObjectURL(url);
  }
  return {
    // State
    isOpen,
    isExporting,
    selectedMode,
    settings: settings2,
    exportResult,
    error,
    // Include context store for use in components
    contextStore,
    // Actions
    open,
    close,
    executeExport
  };
}
const _hoisted_1$E = { class: "flex-1 overflow-y-auto p-6" };
const _hoisted_2$D = { class: "mb-6" };
const _hoisted_3$B = { class: "grid grid-cols-3 gap-3" };
const _hoisted_4$z = ["onClick"];
const _hoisted_5$t = { class: "flex items-start gap-3" };
const _hoisted_6$r = { class: "flex-1 min-w-0" };
const _hoisted_7$q = { class: "text-sm font-medium text-white mb-1" };
const _hoisted_8$o = { class: "text-xs text-gray-400" };
const _hoisted_9$m = {
  key: 0,
  class: "space-y-4"
};
const _hoisted_10$m = { class: "flex items-center gap-2 text-sm text-gray-300 cursor-pointer" };
const _hoisted_11$k = { class: "flex items-center gap-2 text-sm text-gray-300 cursor-pointer" };
const _hoisted_12$h = {
  key: 1,
  class: "space-y-4"
};
const _hoisted_13$h = { class: "flex items-center justify-between text-sm font-medium text-gray-300 mb-2" };
const _hoisted_14$g = { class: "flex items-center gap-2 text-sm text-gray-300 cursor-pointer" };
const _hoisted_15$e = {
  key: 0,
  class: "ml-6 space-y-3 p-3 bg-gray-800/50 rounded-xl border border-gray-700/30"
};
const _hoisted_16$d = {
  key: 2,
  class: "space-y-4"
};
const _hoisted_17$d = { class: "flex items-center gap-2 text-sm text-gray-300 cursor-pointer" };
const _hoisted_18$c = { class: "flex items-center gap-2 text-sm text-gray-300 cursor-pointer" };
const _hoisted_19$c = {
  key: 3,
  class: "mt-6 p-4 bg-gray-800/50 rounded-xl border border-gray-700/30"
};
const _hoisted_20$c = { class: "grid grid-cols-3 gap-3" };
const _hoisted_21$a = { class: "stat-card" };
const _hoisted_22$9 = { class: "stat-card-value" };
const _hoisted_23$9 = { class: "stat-card" };
const _hoisted_24$9 = { class: "stat-card-value" };
const _hoisted_25$8 = { class: "stat-card" };
const _hoisted_26$8 = { class: "stat-card-value stat-card-value-indigo" };
const _hoisted_27$7 = { class: "flex items-center justify-between px-6 py-4 border-t border-gray-700/50" };
const _hoisted_28$7 = { class: "text-xs text-gray-400" };
const _hoisted_29$3 = { key: 0 };
const _hoisted_30$1 = { key: 1 };
const _hoisted_31$1 = { class: "flex items-center gap-3" };
const _hoisted_32$1 = ["disabled"];
const _hoisted_33$1 = {
  key: 0,
  class: "animate-spin h-4 w-4",
  fill: "none",
  viewBox: "0 0 24 24"
};
const _sfc_main$G = /* @__PURE__ */ defineComponent({
  __name: "ExportModal",
  setup(__props, { expose: __expose }) {
    const exportComposable = useExport();
    const contextStore = useContextStore();
    const isOpen = exportComposable.isOpen;
    const isExporting = exportComposable.isExporting;
    const exportResult = exportComposable.exportResult;
    const selectedMode = exportComposable.selectedMode;
    const settings2 = exportComposable.settings;
    const exportModes = [
      {
        value: "clipboard",
        label: "Буфер обмена",
        description: "Текстовый экспорт",
        icon: h("svg", { class: "w-5 h-5", fill: "none", stroke: "currentColor", viewBox: "0 0 24 24" }, [
          h("path", { "stroke-linecap": "round", "stroke-linejoin": "round", "stroke-width": "2", d: "M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2" })
        ])
      },
      {
        value: "ai",
        label: "AI оптимизация",
        description: "С разбиением на чанки",
        icon: h("svg", { class: "w-5 h-5", fill: "none", stroke: "currentColor", viewBox: "0 0 24 24" }, [
          h("path", { "stroke-linecap": "round", "stroke-linejoin": "round", "stroke-width": "2", d: "M9.663 17h4.673M12 3v1m6.364 1.636l-.707.707M21 12h-1M4 12H3m3.343-5.657l-.707-.707m2.828 9.9a5 5 0 117.072 0l-.548.547A3.374 3.374 0 0014 18.469V19a2 2 0 11-4 0v-.531c0-.895-.356-1.754-.988-2.386l-.548-.547z" })
        ])
      },
      {
        value: "human",
        label: "PDF/ZIP",
        description: "Для печати",
        icon: h("svg", { class: "w-5 h-5", fill: "none", stroke: "currentColor", viewBox: "0 0 24 24" }, [
          h("path", { "stroke-linecap": "round", "stroke-linejoin": "round", "stroke-width": "2", d: "M7 21h10a2 2 0 002-2V9.414a1 1 0 00-.293-.707l-5.414-5.414A1 1 0 0012.586 3H7a2 2 0 00-2 2v14a2 2 0 002 2z" })
        ])
      }
    ];
    function close() {
      exportComposable.close();
    }
    async function handleExport() {
      if (!contextStore.hasContext) return;
      await exportComposable.executeExport();
    }
    function formatNumber(num) {
      if (num >= 1e6) return `${(num / 1e6).toFixed(1)}M`;
      if (num >= 1e3) return `${(num / 1e3).toFixed(1)}K`;
      return num.toString();
    }
    __expose({
      open: exportComposable.open,
      close: exportComposable.close
    });
    return (_ctx, _cache) => {
      return openBlock(), createBlock(Teleport, { to: "body" }, [
        createVNode(Transition, { name: "modal" }, {
          default: withCtx(() => [
            unref(isOpen) ? (openBlock(), createElementBlock("div", {
              key: 0,
              class: "fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/60 backdrop-blur-sm",
              onClick: withModifiers(close, ["self"]),
              onKeydown: withKeys(close, ["esc"])
            }, [
              createBaseVNode("div", {
                class: "bg-gray-800/95 backdrop-blur-md rounded-2xl shadow-2xl w-full max-w-3xl max-h-[90vh] flex flex-col border border-gray-700/50",
                onClick: _cache[12] || (_cache[12] = withModifiers(() => {
                }, ["stop"]))
              }, [
                createBaseVNode("div", { class: "flex items-center justify-between px-6 py-4 border-b border-gray-700/50" }, [
                  _cache[14] || (_cache[14] = createBaseVNode("div", { class: "flex items-center gap-3" }, [
                    createBaseVNode("div", { class: "w-10 h-10 rounded-xl bg-indigo-500/20 flex items-center justify-center" }, [
                      createBaseVNode("svg", {
                        class: "w-5 h-5 text-indigo-400",
                        fill: "none",
                        stroke: "currentColor",
                        viewBox: "0 0 24 24"
                      }, [
                        createBaseVNode("path", {
                          "stroke-linecap": "round",
                          "stroke-linejoin": "round",
                          "stroke-width": "2",
                          d: "M8 7H5a2 2 0 00-2 2v9a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-3m-1 4l-3 3m0 0l-3-3m3 3V4"
                        })
                      ])
                    ]),
                    createBaseVNode("div", null, [
                      createBaseVNode("h3", { class: "text-lg font-semibold text-white" }, "Экспорт контекста"),
                      createBaseVNode("p", { class: "text-xs text-gray-400" }, "Выберите формат и параметры экспорта")
                    ])
                  ], -1)),
                  createBaseVNode("button", {
                    onClick: close,
                    class: "p-2 hover:bg-gray-700/50 rounded-xl transition-colors",
                    "aria-label": "Закрыть"
                  }, [..._cache[13] || (_cache[13] = [
                    createBaseVNode("svg", {
                      class: "w-5 h-5 text-gray-400 hover:text-white",
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
                createBaseVNode("div", _hoisted_1$E, [
                  createBaseVNode("div", _hoisted_2$D, [
                    _cache[15] || (_cache[15] = createBaseVNode("label", { class: "block text-sm font-medium text-gray-300 mb-3" }, "Режим экспорта", -1)),
                    createBaseVNode("div", _hoisted_3$B, [
                      (openBlock(), createElementBlock(Fragment, null, renderList(exportModes, (mode) => {
                        return createBaseVNode("button", {
                          key: mode.value,
                          onClick: ($event) => selectedMode.value = mode.value,
                          class: normalizeClass([
                            "p-4 rounded-xl border transition-all text-left",
                            unref(selectedMode) === mode.value ? "border-indigo-500/50 bg-indigo-500/10 shadow-lg shadow-indigo-500/10" : "border-gray-700/50 hover:border-gray-600/50 bg-gray-800/50 hover:bg-gray-700/50"
                          ])
                        }, [
                          createBaseVNode("div", _hoisted_5$t, [
                            (openBlock(), createBlock(resolveDynamicComponent(mode.icon), {
                              class: normalizeClass(["w-5 h-5 flex-shrink-0", unref(selectedMode) === mode.value ? "text-indigo-400" : "text-gray-400"])
                            }, null, 8, ["class"])),
                            createBaseVNode("div", _hoisted_6$r, [
                              createBaseVNode("div", _hoisted_7$q, toDisplayString(mode.label), 1),
                              createBaseVNode("div", _hoisted_8$o, toDisplayString(mode.description), 1)
                            ])
                          ])
                        ], 10, _hoisted_4$z);
                      }), 64))
                    ])
                  ]),
                  unref(selectedMode) === "clipboard" ? (openBlock(), createElementBlock("div", _hoisted_9$m, [
                    createBaseVNode("div", null, [
                      _cache[17] || (_cache[17] = createBaseVNode("label", { class: "block text-sm font-medium text-gray-300 mb-2" }, "Формат", -1)),
                      withDirectives(createBaseVNode("select", {
                        "onUpdate:modelValue": _cache[0] || (_cache[0] = ($event) => unref(settings2).exportFormat = $event),
                        class: "input"
                      }, [..._cache[16] || (_cache[16] = [
                        createBaseVNode("option", { value: "plain" }, "Plain Text", -1),
                        createBaseVNode("option", { value: "manifest" }, "With Manifest", -1),
                        createBaseVNode("option", { value: "json" }, "JSON", -1)
                      ])], 512), [
                        [vModelSelect, unref(settings2).exportFormat]
                      ])
                    ]),
                    createBaseVNode("label", _hoisted_10$m, [
                      withDirectives(createBaseVNode("input", {
                        "onUpdate:modelValue": _cache[1] || (_cache[1] = ($event) => unref(settings2).stripComments = $event),
                        type: "checkbox",
                        class: "w-4 h-4 rounded border-gray-600 bg-gray-800 text-blue-500 focus:ring-0"
                      }, null, 512), [
                        [vModelCheckbox, unref(settings2).stripComments]
                      ]),
                      _cache[18] || (_cache[18] = createTextVNode(" Удалить комментарии ", -1))
                    ]),
                    createBaseVNode("label", _hoisted_11$k, [
                      withDirectives(createBaseVNode("input", {
                        "onUpdate:modelValue": _cache[2] || (_cache[2] = ($event) => unref(settings2).includeManifest = $event),
                        type: "checkbox",
                        class: "w-4 h-4 rounded border-gray-600 bg-gray-800 text-blue-500 focus:ring-0"
                      }, null, 512), [
                        [vModelCheckbox, unref(settings2).includeManifest]
                      ]),
                      _cache[19] || (_cache[19] = createTextVNode(" Включить манифест файлов ", -1))
                    ])
                  ])) : createCommentVNode("", true),
                  unref(selectedMode) === "ai" ? (openBlock(), createElementBlock("div", _hoisted_12$h, [
                    createBaseVNode("div", null, [
                      _cache[21] || (_cache[21] = createBaseVNode("label", { class: "block text-sm font-medium text-gray-300 mb-2" }, "AI Профиль", -1)),
                      withDirectives(createBaseVNode("select", {
                        "onUpdate:modelValue": _cache[3] || (_cache[3] = ($event) => unref(settings2).aiProfile = $event),
                        class: "input"
                      }, [..._cache[20] || (_cache[20] = [
                        createBaseVNode("option", { value: "gpt-4" }, "GPT-4", -1),
                        createBaseVNode("option", { value: "gpt-3.5" }, "GPT-3.5 Turbo", -1),
                        createBaseVNode("option", { value: "claude" }, "Claude", -1),
                        createBaseVNode("option", { value: "gemini" }, "Gemini Pro", -1)
                      ])], 512), [
                        [vModelSelect, unref(settings2).aiProfile]
                      ])
                    ]),
                    createBaseVNode("div", null, [
                      createBaseVNode("label", _hoisted_13$h, [
                        createBaseVNode("span", null, "Лимит токенов: " + toDisplayString(formatNumber(unref(settings2).tokenLimit)), 1)
                      ]),
                      withDirectives(createBaseVNode("input", {
                        "onUpdate:modelValue": _cache[4] || (_cache[4] = ($event) => unref(settings2).tokenLimit = $event),
                        type: "range",
                        min: "1000",
                        max: "128000",
                        step: "1000",
                        class: "w-full h-2 bg-gray-700 rounded-lg appearance-none cursor-pointer accent-blue-500"
                      }, null, 512), [
                        [
                          vModelText,
                          unref(settings2).tokenLimit,
                          void 0,
                          { number: true }
                        ]
                      ])
                    ]),
                    createBaseVNode("label", _hoisted_14$g, [
                      withDirectives(createBaseVNode("input", {
                        "onUpdate:modelValue": _cache[5] || (_cache[5] = ($event) => unref(settings2).enableAutoSplit = $event),
                        type: "checkbox",
                        class: "w-4 h-4 rounded border-gray-600 bg-gray-800 text-blue-500 focus:ring-0"
                      }, null, 512), [
                        [vModelCheckbox, unref(settings2).enableAutoSplit]
                      ]),
                      _cache[22] || (_cache[22] = createTextVNode(" Автоматически разбивать на чанки ", -1))
                    ]),
                    unref(settings2).enableAutoSplit ? (openBlock(), createElementBlock("div", _hoisted_15$e, [
                      createBaseVNode("div", null, [
                        _cache[23] || (_cache[23] = createBaseVNode("label", { class: "block text-xs text-gray-400 mb-1" }, "Токенов на чанк", -1)),
                        withDirectives(createBaseVNode("input", {
                          "onUpdate:modelValue": _cache[6] || (_cache[6] = ($event) => unref(settings2).maxTokensPerChunk = $event),
                          type: "number",
                          min: "500",
                          max: "32000",
                          class: "input text-sm py-1.5"
                        }, null, 512), [
                          [
                            vModelText,
                            unref(settings2).maxTokensPerChunk,
                            void 0,
                            { number: true }
                          ]
                        ])
                      ]),
                      createBaseVNode("div", null, [
                        _cache[24] || (_cache[24] = createBaseVNode("label", { class: "block text-xs text-gray-400 mb-1" }, "Перекрытие токенов", -1)),
                        withDirectives(createBaseVNode("input", {
                          "onUpdate:modelValue": _cache[7] || (_cache[7] = ($event) => unref(settings2).overlapTokens = $event),
                          type: "number",
                          min: "0",
                          max: "1000",
                          class: "input text-sm py-1.5"
                        }, null, 512), [
                          [
                            vModelText,
                            unref(settings2).overlapTokens,
                            void 0,
                            { number: true }
                          ]
                        ])
                      ]),
                      createBaseVNode("div", null, [
                        _cache[26] || (_cache[26] = createBaseVNode("label", { class: "block text-xs text-gray-400 mb-1" }, "Стратегия разбиения", -1)),
                        withDirectives(createBaseVNode("select", {
                          "onUpdate:modelValue": _cache[8] || (_cache[8] = ($event) => unref(settings2).splitStrategy = $event),
                          class: "input text-sm py-1.5"
                        }, [..._cache[25] || (_cache[25] = [
                          createBaseVNode("option", { value: "smart" }, "Smart (рекомендуется)", -1),
                          createBaseVNode("option", { value: "file" }, "По файлам", -1),
                          createBaseVNode("option", { value: "token" }, "По токенам", -1)
                        ])], 512), [
                          [vModelSelect, unref(settings2).splitStrategy]
                        ])
                      ])
                    ])) : createCommentVNode("", true)
                  ])) : createCommentVNode("", true),
                  unref(selectedMode) === "human" ? (openBlock(), createElementBlock("div", _hoisted_16$d, [
                    createBaseVNode("div", null, [
                      _cache[28] || (_cache[28] = createBaseVNode("label", { class: "block text-sm font-medium text-gray-300 mb-2" }, "Тема оформления", -1)),
                      withDirectives(createBaseVNode("select", {
                        "onUpdate:modelValue": _cache[9] || (_cache[9] = ($event) => unref(settings2).theme = $event),
                        class: "input"
                      }, [..._cache[27] || (_cache[27] = [
                        createBaseVNode("option", { value: "default" }, "По умолчанию", -1),
                        createBaseVNode("option", { value: "dark" }, "Тёмная", -1),
                        createBaseVNode("option", { value: "light" }, "Светлая", -1),
                        createBaseVNode("option", { value: "minimal" }, "Минимальная", -1)
                      ])], 512), [
                        [vModelSelect, unref(settings2).theme]
                      ])
                    ]),
                    createBaseVNode("label", _hoisted_17$d, [
                      withDirectives(createBaseVNode("input", {
                        "onUpdate:modelValue": _cache[10] || (_cache[10] = ($event) => unref(settings2).includeLineNumbers = $event),
                        type: "checkbox",
                        class: "w-4 h-4 rounded border-gray-600 bg-gray-800 text-blue-500 focus:ring-0"
                      }, null, 512), [
                        [vModelCheckbox, unref(settings2).includeLineNumbers]
                      ]),
                      _cache[29] || (_cache[29] = createTextVNode(" Показывать номера строк ", -1))
                    ]),
                    createBaseVNode("label", _hoisted_18$c, [
                      withDirectives(createBaseVNode("input", {
                        "onUpdate:modelValue": _cache[11] || (_cache[11] = ($event) => unref(settings2).includePageNumbers = $event),
                        type: "checkbox",
                        class: "w-4 h-4 rounded border-gray-600 bg-gray-800 text-blue-500 focus:ring-0"
                      }, null, 512), [
                        [vModelCheckbox, unref(settings2).includePageNumbers]
                      ]),
                      _cache[30] || (_cache[30] = createTextVNode(" Показывать номера страниц ", -1))
                    ])
                  ])) : createCommentVNode("", true),
                  unref(contextStore).hasContext ? (openBlock(), createElementBlock("div", _hoisted_19$c, [
                    _cache[34] || (_cache[34] = createBaseVNode("div", { class: "text-sm font-medium text-gray-300 mb-3" }, "Предпросмотр контекста", -1)),
                    createBaseVNode("div", _hoisted_20$c, [
                      createBaseVNode("div", _hoisted_21$a, [
                        createBaseVNode("div", _hoisted_22$9, toDisplayString(unref(contextStore).fileCount), 1),
                        _cache[31] || (_cache[31] = createBaseVNode("div", { class: "stat-label" }, "Файлов", -1))
                      ]),
                      createBaseVNode("div", _hoisted_23$9, [
                        createBaseVNode("div", _hoisted_24$9, toDisplayString(formatNumber(unref(contextStore).lineCount)), 1),
                        _cache[32] || (_cache[32] = createBaseVNode("div", { class: "stat-label" }, "Строк", -1))
                      ]),
                      createBaseVNode("div", _hoisted_25$8, [
                        createBaseVNode("div", _hoisted_26$8, toDisplayString(formatNumber(unref(contextStore).tokenCount)), 1),
                        _cache[33] || (_cache[33] = createBaseVNode("div", { class: "stat-label" }, "Токенов", -1))
                      ])
                    ])
                  ])) : createCommentVNode("", true)
                ]),
                createBaseVNode("div", _hoisted_27$7, [
                  createBaseVNode("div", _hoisted_28$7, [
                    unref(exportResult) ? (openBlock(), createElementBlock("span", _hoisted_29$3, " Экспорт завершен успешно ")) : unref(isExporting) ? (openBlock(), createElementBlock("span", _hoisted_30$1, " Экспорт в процессе... ")) : createCommentVNode("", true)
                  ]),
                  createBaseVNode("div", _hoisted_31$1, [
                    createBaseVNode("button", {
                      onClick: close,
                      class: "btn btn-secondary"
                    }, " Отмена "),
                    createBaseVNode("button", {
                      onClick: handleExport,
                      disabled: unref(isExporting) || !unref(contextStore).hasContext,
                      class: "btn btn-primary flex items-center gap-2"
                    }, [
                      unref(isExporting) ? (openBlock(), createElementBlock("svg", _hoisted_33$1, [..._cache[35] || (_cache[35] = [
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
                          d: "M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
                        }, null, -1)
                      ])])) : createCommentVNode("", true),
                      createTextVNode(" " + toDisplayString(unref(isExporting) ? "Экспорт..." : "Экспортировать"), 1)
                    ], 8, _hoisted_32$1)
                  ])
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
const ExportModal = /* @__PURE__ */ _export_sfc(_sfc_main$G, [["__scopeId", "data-v-8f1b982f"]]);
const _hoisted_1$D = { class: "skeleton-loader" };
const _hoisted_2$C = { class: "skeleton-content" };
const _hoisted_3$A = { class: "skeleton-overlay" };
const _hoisted_4$y = { class: "skeleton-status-text" };
const _hoisted_5$s = { class: "skeleton-status-dots" };
const _hoisted_6$q = { class: "skeleton-progress" };
const _sfc_main$F = /* @__PURE__ */ defineComponent({
  __name: "SkeletonLoader",
  props: {
    progress: {},
    statusText: {}
  },
  setup(__props) {
    const dots = ref("");
    let dotsInterval = null;
    onMounted(() => {
      dotsInterval = setInterval(() => {
        dots.value = dots.value.length >= 3 ? "" : dots.value + ".";
      }, 400);
    });
    onUnmounted(() => {
      if (dotsInterval) clearInterval(dotsInterval);
    });
    const codeLines = [
      { width: "25%", indent: "0" },
      // import
      { width: "35%", indent: "0" },
      // import
      { width: "20%", indent: "0" },
      // empty/short
      { width: "40%", indent: "0" },
      // function declaration
      { width: "60%", indent: "1.5rem" },
      // indented code
      { width: "75%", indent: "1.5rem" },
      // longer line
      { width: "45%", indent: "1.5rem" },
      // medium line
      { width: "55%", indent: "3rem" },
      // nested
      { width: "80%", indent: "3rem" },
      // long nested
      { width: "30%", indent: "3rem" },
      // short nested
      { width: "50%", indent: "1.5rem" },
      // back to first indent
      { width: "15%", indent: "0" },
      // closing brace
      { width: "10%", indent: "0" },
      // empty
      { width: "45%", indent: "0" },
      // new function
      { width: "70%", indent: "1.5rem" },
      // code
      { width: "35%", indent: "1.5rem" }
      // return
    ];
    return (_ctx, _cache) => {
      return openBlock(), createElementBlock("div", _hoisted_1$D, [
        createBaseVNode("div", _hoisted_2$C, [
          (openBlock(), createElementBlock(Fragment, null, renderList(codeLines, (line, i) => {
            return createBaseVNode("div", {
              key: i,
              class: "skeleton-line",
              style: normalizeStyle({
                width: line.width,
                marginLeft: line.indent,
                animationDelay: `${i % 6 * 0.08}s`
              })
            }, null, 4);
          }), 64)),
          createBaseVNode("div", _hoisted_3$A, [
            createBaseVNode("span", _hoisted_4$y, toDisplayString(__props.statusText), 1),
            createBaseVNode("span", _hoisted_5$s, toDisplayString(dots.value), 1)
          ])
        ]),
        createBaseVNode("div", _hoisted_6$q, [
          createBaseVNode("div", {
            class: "skeleton-progress-fill",
            style: normalizeStyle({ width: `${__props.progress}%` })
          }, null, 4)
        ])
      ]);
    };
  }
});
const SkeletonLoader = /* @__PURE__ */ _export_sfc(_sfc_main$F, [["__scopeId", "data-v-6095d76b"]]);
const _hoisted_1$C = {
  key: 0,
  class: "chunk-boundary"
};
const _hoisted_2$B = { class: "chunk-boundary__label" };
const _hoisted_3$z = { class: "line-number" };
const _hoisted_4$x = ["innerHTML"];
const LINE_HEIGHT = 20;
const _sfc_main$E = /* @__PURE__ */ defineComponent({
  __name: "VirtualCodeView",
  props: {
    lines: {},
    highlightedLines: {},
    searchQuery: {},
    chunkBoundaries: {},
    outputFormat: {}
  },
  setup(__props, { expose: __expose }) {
    const props = __props;
    const containerRef = ref(null);
    const lineCount = computed(() => props.lines.length);
    const virtualizer = useVirtualizer({
      get count() {
        return lineCount.value;
      },
      getScrollElement: () => containerRef.value,
      estimateSize: () => LINE_HEIGHT,
      overscan: 30
      // Render extra items above/below viewport for smooth scroll
    });
    const virtualItems = computed(() => virtualizer.value.getVirtualItems());
    const totalHeight = computed(() => virtualizer.value.getTotalSize());
    const offsetY = computed(() => virtualItems.value[0]?.start ?? 0);
    function isLineHighlighted(lineIndex) {
      return props.highlightedLines?.has(lineIndex) ?? false;
    }
    function getChunkNumber(lineIndex) {
      if (!props.chunkBoundaries) return 1;
      let chunkNum = 1;
      for (const boundary of Array.from(props.chunkBoundaries).sort((a, b) => a - b)) {
        if (boundary <= lineIndex) chunkNum++;
        else break;
      }
      return chunkNum;
    }
    function highlightLine(line) {
      let result = escapeHtml(line);
      result = applySyntaxHighlight(result, props.outputFormat || "plain");
      if (props.searchQuery) {
        const query = escapeHtml(props.searchQuery);
        const regex = new RegExp(`(${escapeRegex(query)})`, "gi");
        result = result.replace(regex, '<mark class="search-highlight">$1</mark>');
      }
      return result;
    }
    function applySyntaxHighlight(line, format) {
      switch (format) {
        case "xml":
          return highlightXml(line);
        case "markdown":
          return highlightMarkdown(line);
        default:
          return highlightPlain(line);
      }
    }
    function highlightXml(line) {
      return line.replace(/(&lt;\/?)(file|content)(&gt;)/g, '<span class="syntax-tag">$1$2$3</span>').replace(
        /(path)(=)(&quot;[^&]*&quot;)/g,
        '<span class="syntax-attr">$1</span><span class="syntax-punct">$2</span><span class="syntax-string">$3</span>'
      );
    }
    function highlightMarkdown(line) {
      if (line.startsWith("## File:") || line.startsWith("# ")) {
        return `<span class="syntax-heading">${line}</span>`;
      }
      if (line.startsWith("```")) {
        return `<span class="syntax-fence">${line}</span>`;
      }
      return line;
    }
    function highlightPlain(line) {
      if (line.startsWith("--- File:") && line.endsWith("---")) {
        return `<span class="syntax-separator">${line}</span>`;
      }
      return line;
    }
    function escapeHtml(text) {
      return text.replace(/&/g, "&amp;").replace(/</g, "&lt;").replace(/>/g, "&gt;");
    }
    function escapeRegex(text) {
      return text.replace(/[.*+?^${}()|[\]\\]/g, "\\$&");
    }
    function scrollToLine(lineIndex) {
      virtualizer.value.scrollToIndex(lineIndex, { align: "center" });
    }
    __expose({ scrollToLine, containerRef });
    watch(() => props.lines.length, () => virtualizer.value.measure());
    return (_ctx, _cache) => {
      return openBlock(), createElementBlock("div", {
        ref_key: "containerRef",
        ref: containerRef,
        class: "virtual-code-view"
      }, [
        createBaseVNode("div", {
          class: "virtual-code-view__spacer",
          style: normalizeStyle({ height: `${totalHeight.value}px` })
        }, [
          createBaseVNode("div", {
            class: "virtual-code-view__content",
            style: normalizeStyle({ transform: `translateY(${offsetY.value}px)` })
          }, [
            (openBlock(true), createElementBlock(Fragment, null, renderList(virtualItems.value, (item) => {
              return openBlock(), createElementBlock(Fragment, {
                key: item.index
              }, [
                __props.chunkBoundaries?.has(item.index) ? (openBlock(), createElementBlock("div", _hoisted_1$C, [
                  _cache[0] || (_cache[0] = createBaseVNode("span", { class: "chunk-boundary__line" }, null, -1)),
                  createBaseVNode("span", _hoisted_2$B, "✂️ Chunk " + toDisplayString(getChunkNumber(item.index)), 1),
                  _cache[1] || (_cache[1] = createBaseVNode("span", { class: "chunk-boundary__line" }, null, -1))
                ])) : createCommentVNode("", true),
                createBaseVNode("div", {
                  class: normalizeClass(["code-line", { "code-line-highlight": isLineHighlighted(item.index) }]),
                  style: normalizeStyle({ height: `${LINE_HEIGHT}px` })
                }, [
                  createBaseVNode("span", _hoisted_3$z, toDisplayString(item.index + 1), 1),
                  createBaseVNode("span", {
                    innerHTML: highlightLine(__props.lines[item.index] || "")
                  }, null, 8, _hoisted_4$x)
                ], 6)
              ], 64);
            }), 128))
          ], 4)
        ], 4)
      ], 512);
    };
  }
});
const VirtualCodeView = /* @__PURE__ */ _export_sfc(_sfc_main$E, [["__scopeId", "data-v-7a4f8130"]]);
const _hoisted_1$B = {
  class: "context-content-area layout-fill layout-scroll",
  ref: "contentContainer"
};
const _hoisted_2$A = {
  key: 1,
  class: "flex items-center justify-center h-full"
};
const _hoisted_3$y = { class: "text-center max-w-md mx-auto px-4" };
const _hoisted_4$w = { class: "text-lg font-semibold text-amber-400 mb-4" };
const _hoisted_5$r = { class: "context-stats mb-4" };
const _hoisted_6$p = { class: "stats-grid" };
const _hoisted_7$p = { class: "stat-card" };
const _hoisted_8$n = { class: "stat-card-value" };
const _hoisted_9$l = { class: "stat-label" };
const _hoisted_10$l = { class: "stat-card" };
const _hoisted_11$j = { class: "stat-card-value" };
const _hoisted_12$g = { class: "stat-label" };
const _hoisted_13$g = { class: "info-box text-left" };
const _hoisted_14$f = { class: "text-sm font-medium text-gray-300 mb-2" };
const _hoisted_15$d = { class: "text-sm text-gray-400 space-y-1.5" };
const _hoisted_16$c = {
  key: 2,
  class: "flex items-center justify-center h-full"
};
const _hoisted_17$c = { class: "text-center max-w-md" };
const _hoisted_18$b = { class: "text-lg font-semibold text-red-400 mb-2" };
const _hoisted_19$b = { class: "info-box text-left mt-4" };
const _hoisted_20$b = { class: "text-sm font-medium text-gray-300 mb-2" };
const _hoisted_21$9 = { class: "text-sm text-gray-400 space-y-1.5" };
const _hoisted_22$8 = {
  key: 3,
  class: "empty-state-enhanced h-full flex flex-col items-center justify-center"
};
const _hoisted_23$8 = { class: "text-base font-semibold text-white mb-3" };
const _hoisted_24$8 = { class: "text-left max-w-xs space-y-2 mb-6" };
const _hoisted_25$7 = { class: "flex items-center gap-3 text-sm" };
const _hoisted_26$7 = { class: "text-gray-400" };
const _hoisted_27$6 = { class: "flex items-center gap-3 text-sm" };
const _hoisted_28$6 = { class: "text-gray-400" };
const _hoisted_29$2 = { class: "flex items-center gap-6 text-gray-400" };
const _hoisted_30 = { class: "flex items-center gap-2" };
const _hoisted_31 = { class: "text-xs" };
const _hoisted_32 = { class: "flex items-center gap-2" };
const _hoisted_33 = { class: "text-xs" };
const _hoisted_34 = {
  key: 4,
  class: "code-editor context-content min-h-full"
};
const _hoisted_35 = {
  key: 1,
  class: "text-center py-8"
};
const _hoisted_36 = { class: "text-gray-400" };
const _hoisted_37 = {
  key: 2,
  class: "text-center py-8"
};
const _hoisted_38 = {
  key: 0,
  class: "text-gray-400"
};
const _hoisted_39 = {
  key: 1,
  class: "text-gray-400"
};
const _sfc_main$D = /* @__PURE__ */ defineComponent({
  __name: "ContextPanelContent",
  props: {
    isBuilding: { type: Boolean },
    buildProgress: {},
    statusText: {},
    fileCount: {},
    lineCount: {},
    totalSize: {},
    contextId: {},
    error: {},
    hasContext: { type: Boolean },
    isLoading: { type: Boolean },
    lines: {},
    highlightedLines: {},
    searchQuery: {},
    chunkBoundaries: {},
    outputFormat: {}
  },
  setup(__props, { expose: __expose }) {
    const { t } = useI18n();
    const virtualCodeRef = ref(null);
    const exportModalRef = ref(null);
    function scrollToLine(lineNum) {
      virtualCodeRef.value?.scrollToLine(lineNum);
    }
    function openExportModal() {
      exportModalRef.value?.open();
    }
    __expose({
      scrollToLine,
      openExportModal
    });
    return (_ctx, _cache) => {
      return openBlock(), createElementBlock("div", _hoisted_1$B, [
        __props.isBuilding ? (openBlock(), createBlock(SkeletonLoader, {
          key: 0,
          progress: __props.buildProgress,
          "status-text": __props.statusText
        }, null, 8, ["progress", "status-text"])) : (__props.fileCount === 0 || __props.totalSize === 0 || __props.lineCount === 0) && __props.contextId ? (openBlock(), createElementBlock("div", _hoisted_2$A, [
          createBaseVNode("div", _hoisted_3$y, [
            _cache[0] || (_cache[0] = createBaseVNode("div", { class: "w-16 h-16 mx-auto mb-4 bg-amber-500/20 rounded-2xl flex items-center justify-center border border-amber-500/30" }, [
              createBaseVNode("svg", {
                class: "w-8 h-8 text-amber-400",
                fill: "none",
                stroke: "currentColor",
                viewBox: "0 0 24 24"
              }, [
                createBaseVNode("path", {
                  "stroke-linecap": "round",
                  "stroke-linejoin": "round",
                  "stroke-width": "2",
                  d: "M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"
                })
              ])
            ], -1)),
            createBaseVNode("p", _hoisted_4$w, toDisplayString(unref(t)("error.emptyContext")), 1),
            createBaseVNode("div", _hoisted_5$r, [
              createBaseVNode("div", _hoisted_6$p, [
                createBaseVNode("div", _hoisted_7$p, [
                  createBaseVNode("div", _hoisted_8$n, toDisplayString(__props.fileCount), 1),
                  createBaseVNode("div", _hoisted_9$l, toDisplayString(unref(t)("context.files")), 1)
                ]),
                createBaseVNode("div", _hoisted_10$l, [
                  createBaseVNode("div", _hoisted_11$j, toDisplayString(__props.lineCount), 1),
                  createBaseVNode("div", _hoisted_12$g, toDisplayString(unref(t)("context.lines")), 1)
                ])
              ])
            ]),
            createBaseVNode("div", _hoisted_13$g, [
              createBaseVNode("p", _hoisted_14$f, toDisplayString(unref(t)("error.suggestions")) + ":", 1),
              createBaseVNode("ul", _hoisted_15$d, [
                createBaseVNode("li", null, "• " + toDisplayString(unref(t)("error.checkFiles")), 1),
                createBaseVNode("li", null, "• " + toDisplayString(unref(t)("error.checkPaths")), 1),
                createBaseVNode("li", null, "• " + toDisplayString(unref(t)("error.tryRefresh")), 1)
              ])
            ])
          ])
        ])) : __props.error ? (openBlock(), createElementBlock("div", _hoisted_16$c, [
          createBaseVNode("div", _hoisted_17$c, [
            _cache[1] || (_cache[1] = createBaseVNode("div", { class: "w-16 h-16 mx-auto mb-4 bg-red-500/20 rounded-2xl flex items-center justify-center border border-red-500/30" }, [
              createBaseVNode("svg", {
                class: "w-8 h-8 text-red-400",
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
              ])
            ], -1)),
            createBaseVNode("p", _hoisted_18$b, toDisplayString(__props.error), 1),
            createBaseVNode("div", _hoisted_19$b, [
              createBaseVNode("p", _hoisted_20$b, toDisplayString(unref(t)("error.suggestions")) + ":", 1),
              createBaseVNode("ul", _hoisted_21$9, [
                createBaseVNode("li", null, "• " + toDisplayString(unref(t)("error.checkFiles")), 1),
                createBaseVNode("li", null, "• " + toDisplayString(unref(t)("error.checkPaths")), 1),
                createBaseVNode("li", null, "• " + toDisplayString(unref(t)("error.tryRefresh")), 1)
              ])
            ])
          ])
        ])) : !__props.hasContext ? (openBlock(), createElementBlock("div", _hoisted_22$8, [
          _cache[6] || (_cache[6] = createBaseVNode("div", { class: "empty-state-icon-glow mb-4" }, [
            createBaseVNode("svg", {
              class: "w-8 h-8 text-indigo-400",
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
            ])
          ], -1)),
          createBaseVNode("p", _hoisted_23$8, toDisplayString(unref(t)("context.notBuilt")), 1),
          createBaseVNode("div", _hoisted_24$8, [
            createBaseVNode("div", _hoisted_25$7, [
              _cache[2] || (_cache[2] = createBaseVNode("span", { class: "flex-shrink-0 w-5 h-5 rounded-full bg-indigo-500/20 text-indigo-400 text-xs font-bold flex items-center justify-center" }, "1", -1)),
              createBaseVNode("span", _hoisted_26$7, toDisplayString(unref(t)("context.step1")), 1)
            ]),
            createBaseVNode("div", _hoisted_27$6, [
              _cache[3] || (_cache[3] = createBaseVNode("span", { class: "flex-shrink-0 w-5 h-5 rounded-full bg-indigo-500/20 text-indigo-400 text-xs font-bold flex items-center justify-center" }, "2", -1)),
              createBaseVNode("span", _hoisted_28$6, toDisplayString(unref(t)("context.step2")), 1)
            ])
          ]),
          createBaseVNode("div", _hoisted_29$2, [
            createBaseVNode("div", _hoisted_30, [
              _cache[4] || (_cache[4] = createBaseVNode("svg", {
                class: "w-4 h-4 animate-[pulse_3s_ease-in-out_infinite]",
                fill: "none",
                stroke: "currentColor",
                viewBox: "0 0 24 24"
              }, [
                createBaseVNode("path", {
                  "stroke-linecap": "round",
                  "stroke-linejoin": "round",
                  "stroke-width": "2",
                  d: "M10 19l-7-7m0 0l7-7m-7 7h18"
                })
              ], -1)),
              createBaseVNode("span", _hoisted_31, toDisplayString(unref(t)("context.selectHint")), 1)
            ]),
            createBaseVNode("div", _hoisted_32, [
              createBaseVNode("span", _hoisted_33, toDisplayString(unref(t)("context.chatHint")), 1),
              _cache[5] || (_cache[5] = createBaseVNode("svg", {
                class: "w-4 h-4 animate-[pulse_3s_ease-in-out_infinite]",
                fill: "none",
                stroke: "currentColor",
                viewBox: "0 0 24 24"
              }, [
                createBaseVNode("path", {
                  "stroke-linecap": "round",
                  "stroke-linejoin": "round",
                  "stroke-width": "2",
                  d: "M14 5l7 7m0 0l-7 7m7-7H3"
                })
              ], -1))
            ])
          ])
        ])) : (openBlock(), createElementBlock("div", _hoisted_34, [
          createVNode(unref(TemplatePreviewBlock)),
          __props.lines?.length ? (openBlock(), createBlock(VirtualCodeView, {
            key: 0,
            ref_key: "virtualCodeRef",
            ref: virtualCodeRef,
            lines: __props.lines,
            "highlighted-lines": __props.highlightedLines,
            "search-query": __props.searchQuery,
            "chunk-boundaries": __props.chunkBoundaries,
            "output-format": __props.outputFormat
          }, null, 8, ["lines", "highlighted-lines", "search-query", "chunk-boundaries", "output-format"])) : __props.isLoading ? (openBlock(), createElementBlock("div", _hoisted_35, [
            _cache[7] || (_cache[7] = createBaseVNode("svg", {
              class: "animate-spin h-6 w-6 text-blue-500 mx-auto mb-2",
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
            createBaseVNode("p", _hoisted_36, toDisplayString(unref(t)("context.loading")), 1)
          ])) : (openBlock(), createElementBlock("div", _hoisted_37, [
            __props.contextId ? (openBlock(), createElementBlock("p", _hoisted_38, toDisplayString(unref(t)("context.loading")), 1)) : (openBlock(), createElementBlock("p", _hoisted_39, toDisplayString(unref(t)("context.notBuilt")), 1))
          ]))
        ])),
        createVNode(ExportModal, {
          ref_key: "exportModalRef",
          ref: exportModalRef
        }, null, 512)
      ], 512);
    };
  }
});
const ContextPanelContent = /* @__PURE__ */ _export_sfc(_sfc_main$D, [["__scopeId", "data-v-e700ed93"]]);
const _sfc_main$C = /* @__PURE__ */ defineComponent({
  __name: "Tooltip",
  props: {
    text: {},
    position: {}
  },
  setup(__props) {
    const show = ref(false);
    return (_ctx, _cache) => {
      return openBlock(), createElementBlock("div", {
        class: "tooltip-wrapper",
        onMouseenter: _cache[0] || (_cache[0] = ($event) => show.value = true),
        onMouseleave: _cache[1] || (_cache[1] = ($event) => show.value = false)
      }, [
        renderSlot(_ctx.$slots, "default", {}, void 0, true),
        createVNode(Transition, { name: "tooltip" }, {
          default: withCtx(() => [
            show.value ? (openBlock(), createElementBlock("div", {
              key: 0,
              class: normalizeClass(["tooltip", [`tooltip--${__props.position}`]])
            }, toDisplayString(__props.text), 3)) : createCommentVNode("", true)
          ]),
          _: 1
        })
      ], 32);
    };
  }
});
const Tooltip = /* @__PURE__ */ _export_sfc(_sfc_main$C, [["__scopeId", "data-v-c5937fb2"]]);
const _hoisted_1$A = {
  key: 0,
  class: "unified-hud"
};
const _hoisted_2$z = { class: "hud-section hud-tools" };
const _hoisted_3$x = { class: "hud-section hud-nav" };
const _hoisted_4$v = ["disabled"];
const _hoisted_5$q = { class: "hud-counter-wrap" };
const _hoisted_6$o = { class: "hud-counter" };
const _hoisted_7$o = { class: "hud-counter-current" };
const _hoisted_8$m = { class: "hud-counter-total" };
const _hoisted_9$k = { class: "hud-progress-dots" };
const _hoisted_10$k = ["disabled"];
const _hoisted_11$i = {
  key: 0,
  class: "w-3.5 h-3.5",
  fill: "none",
  stroke: "currentColor",
  viewBox: "0 0 24 24"
};
const _sfc_main$B = /* @__PURE__ */ defineComponent({
  __name: "ContextPanelFooter",
  props: {
    visible: { type: Boolean },
    showChunkNav: { type: Boolean },
    currentChunk: {},
    totalChunks: {},
    copySuccess: { type: Boolean },
    isChunkCopied: { type: Function }
  },
  emits: ["clear", "export", "copy", "prev-chunk", "next-chunk"],
  setup(__props) {
    const props = __props;
    const { t } = useI18n();
    const buttonText = computed(() => {
      if (props.copySuccess) return "✓";
      if (props.showChunkNav) return `#${props.currentChunk}`;
      return t("context.copy");
    });
    return (_ctx, _cache) => {
      return openBlock(), createBlock(Transition, { name: "hud" }, {
        default: withCtx(() => [
          __props.visible ? (openBlock(), createElementBlock("div", _hoisted_1$A, [
            createBaseVNode("div", _hoisted_2$z, [
              createVNode(Tooltip, {
                text: unref(t)("context.clear"),
                position: "top"
              }, {
                default: withCtx(() => [
                  createBaseVNode("button", {
                    onClick: _cache[0] || (_cache[0] = ($event) => _ctx.$emit("clear")),
                    class: "hud-tool-btn hud-tool-btn--danger"
                  }, [..._cache[5] || (_cache[5] = [
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
                        d: "M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"
                      })
                    ], -1)
                  ])])
                ]),
                _: 1
              }, 8, ["text"]),
              createVNode(Tooltip, {
                text: unref(t)("context.export"),
                position: "top"
              }, {
                default: withCtx(() => [
                  createBaseVNode("button", {
                    onClick: _cache[1] || (_cache[1] = ($event) => _ctx.$emit("export")),
                    class: "hud-tool-btn"
                  }, [..._cache[6] || (_cache[6] = [
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
                        d: "M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4"
                      })
                    ], -1)
                  ])])
                ]),
                _: 1
              }, 8, ["text"])
            ]),
            _cache[12] || (_cache[12] = createBaseVNode("div", { class: "hud-divider" }, null, -1)),
            __props.showChunkNav ? (openBlock(), createElementBlock(Fragment, { key: 0 }, [
              createBaseVNode("div", _hoisted_3$x, [
                createBaseVNode("button", {
                  class: "hud-nav-btn",
                  disabled: __props.currentChunk <= 1,
                  onClick: _cache[2] || (_cache[2] = ($event) => _ctx.$emit("prev-chunk"))
                }, [..._cache[7] || (_cache[7] = [
                  createBaseVNode("svg", {
                    class: "w-3.5 h-3.5",
                    fill: "none",
                    stroke: "currentColor",
                    viewBox: "0 0 24 24"
                  }, [
                    createBaseVNode("path", {
                      "stroke-linecap": "round",
                      "stroke-linejoin": "round",
                      "stroke-width": "2",
                      d: "M15 19l-7-7 7-7"
                    })
                  ], -1)
                ])], 8, _hoisted_4$v),
                createBaseVNode("div", _hoisted_5$q, [
                  createBaseVNode("div", _hoisted_6$o, [
                    createBaseVNode("span", _hoisted_7$o, toDisplayString(__props.currentChunk), 1),
                    _cache[8] || (_cache[8] = createBaseVNode("span", { class: "hud-counter-sep" }, "/", -1)),
                    createBaseVNode("span", _hoisted_8$m, toDisplayString(__props.totalChunks), 1)
                  ]),
                  createBaseVNode("div", _hoisted_9$k, [
                    (openBlock(true), createElementBlock(Fragment, null, renderList(__props.totalChunks, (i) => {
                      return openBlock(), createElementBlock("span", {
                        key: i,
                        class: normalizeClass(["hud-dot", {
                          active: i === __props.currentChunk,
                          copied: __props.isChunkCopied(i - 1)
                        }])
                      }, null, 2);
                    }), 128))
                  ])
                ]),
                createBaseVNode("button", {
                  class: "hud-nav-btn",
                  disabled: __props.currentChunk >= __props.totalChunks,
                  onClick: _cache[3] || (_cache[3] = ($event) => _ctx.$emit("next-chunk"))
                }, [..._cache[9] || (_cache[9] = [
                  createBaseVNode("svg", {
                    class: "w-3.5 h-3.5",
                    fill: "none",
                    stroke: "currentColor",
                    viewBox: "0 0 24 24"
                  }, [
                    createBaseVNode("path", {
                      "stroke-linecap": "round",
                      "stroke-linejoin": "round",
                      "stroke-width": "2",
                      d: "M9 5l7 7-7 7"
                    })
                  ], -1)
                ])], 8, _hoisted_10$k)
              ]),
              _cache[10] || (_cache[10] = createBaseVNode("div", { class: "hud-divider" }, null, -1))
            ], 64)) : createCommentVNode("", true),
            createBaseVNode("button", {
              class: normalizeClass(["hud-action-btn", { "hud-action-btn--success": __props.copySuccess }]),
              onClick: _cache[4] || (_cache[4] = ($event) => _ctx.$emit("copy"))
            }, [
              !__props.copySuccess ? (openBlock(), createElementBlock("svg", _hoisted_11$i, [..._cache[11] || (_cache[11] = [
                createBaseVNode("path", {
                  "stroke-linecap": "round",
                  "stroke-linejoin": "round",
                  "stroke-width": "2",
                  d: "M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z"
                }, null, -1)
              ])])) : createCommentVNode("", true),
              createBaseVNode("span", null, toDisplayString(buttonText.value), 1)
            ], 2)
          ])) : createCommentVNode("", true)
        ]),
        _: 1
      });
    };
  }
});
const ContextPanelFooter = /* @__PURE__ */ _export_sfc(_sfc_main$B, [["__scopeId", "data-v-c09dd950"]]);
const _hoisted_1$z = { class: "stats-popover-header" };
const _hoisted_2$y = { class: "stats-popover-title" };
const _hoisted_3$w = { class: "stats-popover-grid" };
const _hoisted_4$u = { class: "stats-popover-card" };
const _hoisted_5$p = { class: "stats-popover-value" };
const _hoisted_6$n = { class: "stats-popover-label" };
const _hoisted_7$n = { class: "stats-popover-card" };
const _hoisted_8$l = { class: "stats-popover-value" };
const _hoisted_9$j = { class: "stats-popover-label" };
const _hoisted_10$j = { class: "stats-popover-card" };
const _hoisted_11$h = { class: "stats-popover-value stats-popover-value--accent" };
const _hoisted_12$f = { class: "stats-popover-label" };
const _hoisted_13$f = { class: "stats-popover-card" };
const _hoisted_14$e = { class: "stats-popover-value stats-popover-value--success" };
const _hoisted_15$c = { class: "stats-popover-label" };
const _hoisted_16$b = {
  key: 0,
  class: "stats-popover-section"
};
const _hoisted_17$b = { class: "stats-popover-section-title" };
const _hoisted_18$a = { class: "stats-popover-list" };
const _hoisted_19$a = { class: "stats-popover-item-icon" };
const _hoisted_20$a = { class: "stats-popover-item-name" };
const _hoisted_21$8 = { class: "stats-popover-item-count" };
const _hoisted_22$7 = { class: "stats-popover-item-bar" };
const _hoisted_23$7 = {
  key: 1,
  class: "stats-popover-section"
};
const _hoisted_24$7 = { class: "stats-popover-section-title" };
const _hoisted_25$6 = { class: "stats-popover-list" };
const _hoisted_26$6 = { class: "stats-popover-item-name" };
const _hoisted_27$5 = { class: "stats-popover-item-count" };
const _hoisted_28$5 = { class: "stats-popover-item-bar" };
const _sfc_main$A = /* @__PURE__ */ defineComponent({
  __name: "StatsPopover",
  props: {
    visible: { type: Boolean },
    anchorRect: {}
  },
  emits: ["close"],
  setup(__props) {
    const { t } = useI18n();
    const contextStore = useContextStore();
    const fileStore = useFileStore();
    const fileCount = computed(() => contextStore.fileCount);
    const lineCount = computed(() => contextStore.lineCount);
    const tokenCount = computed(() => contextStore.tokenCount);
    const estimatedCost = computed(() => contextStore.estimatedCost);
    const popoverStyle = computed(() => ({
      // Center in viewport
    }));
    const fileTypeStats = computed(() => {
      const selected = Array.from(fileStore.selectedPaths);
      if (selected.length === 0) return [];
      const counts = /* @__PURE__ */ new Map();
      for (const path of selected) {
        const ext = path.split(".").pop()?.toLowerCase() || "other";
        counts.set(ext, (counts.get(ext) || 0) + 1);
      }
      const total = selected.length;
      return Array.from(counts.entries()).map(([extension, count]) => {
        const info = FILE_TYPE_CONFIG[extension] || FILE_TYPE_CONFIG.default;
        return {
          extension,
          count,
          percentage: Math.round(count / total * 100),
          icon: info.icon,
          colorClass: info.colorClass
        };
      }).sort((a, b) => b.count - a.count).slice(0, 6);
    });
    const folderStats = computed(() => {
      const selected = Array.from(fileStore.selectedPaths);
      if (selected.length === 0) return [];
      const counts = /* @__PURE__ */ new Map();
      for (const path of selected) {
        const parts = path.split("/");
        const folder = parts.length > 1 ? parts[0] : "/";
        counts.set(folder, (counts.get(folder) || 0) + 1);
      }
      const total = selected.length;
      return Array.from(counts.entries()).map(([folder, count]) => ({
        folder,
        count,
        percentage: Math.round(count / total * 100)
      })).sort((a, b) => b.count - a.count);
    });
    function formatNumber(num) {
      if (num >= 1e6) return `${(num / 1e6).toFixed(1)}M`;
      if (num >= 1e3) return `${(num / 1e3).toFixed(1)}K`;
      return num.toString();
    }
    return (_ctx, _cache) => {
      return openBlock(), createBlock(Teleport, { to: "body" }, [
        createVNode(Transition, { name: "popover" }, {
          default: withCtx(() => [
            __props.visible ? (openBlock(), createElementBlock("div", {
              key: 0,
              class: "stats-popover-overlay",
              onClick: _cache[1] || (_cache[1] = withModifiers(($event) => _ctx.$emit("close"), ["self"]))
            }, [
              createBaseVNode("div", {
                class: "stats-popover",
                style: normalizeStyle(popoverStyle.value)
              }, [
                createBaseVNode("div", _hoisted_1$z, [
                  createBaseVNode("div", _hoisted_2$y, [
                    _cache[2] || (_cache[2] = createBaseVNode("svg", {
                      class: "w-4 h-4",
                      fill: "none",
                      stroke: "currentColor",
                      viewBox: "0 0 24 24"
                    }, [
                      createBaseVNode("path", {
                        "stroke-linecap": "round",
                        "stroke-linejoin": "round",
                        "stroke-width": "2",
                        d: "M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z"
                      })
                    ], -1)),
                    createBaseVNode("span", null, toDisplayString(unref(t)("stats.title")), 1)
                  ]),
                  createBaseVNode("button", {
                    onClick: _cache[0] || (_cache[0] = ($event) => _ctx.$emit("close")),
                    class: "stats-popover-close"
                  }, [..._cache[3] || (_cache[3] = [
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
                        d: "M6 18L18 6M6 6l12 12"
                      })
                    ], -1)
                  ])])
                ]),
                createBaseVNode("div", _hoisted_3$w, [
                  createBaseVNode("div", _hoisted_4$u, [
                    createBaseVNode("div", _hoisted_5$p, toDisplayString(fileCount.value), 1),
                    createBaseVNode("div", _hoisted_6$n, toDisplayString(unref(t)("context.files")), 1)
                  ]),
                  createBaseVNode("div", _hoisted_7$n, [
                    createBaseVNode("div", _hoisted_8$l, toDisplayString(formatNumber(lineCount.value)), 1),
                    createBaseVNode("div", _hoisted_9$j, toDisplayString(unref(t)("context.lines")), 1)
                  ]),
                  createBaseVNode("div", _hoisted_10$j, [
                    createBaseVNode("div", _hoisted_11$h, toDisplayString(formatNumber(tokenCount.value)), 1),
                    createBaseVNode("div", _hoisted_12$f, toDisplayString(unref(t)("action.tokens")), 1)
                  ]),
                  createBaseVNode("div", _hoisted_13$f, [
                    createBaseVNode("div", _hoisted_14$e, "$" + toDisplayString(estimatedCost.value.toFixed(4)), 1),
                    createBaseVNode("div", _hoisted_15$c, toDisplayString(unref(t)("action.cost")), 1)
                  ])
                ]),
                fileTypeStats.value.length > 0 ? (openBlock(), createElementBlock("div", _hoisted_16$b, [
                  createBaseVNode("div", _hoisted_17$b, toDisplayString(unref(t)("stats.byType")), 1),
                  createBaseVNode("div", _hoisted_18$a, [
                    (openBlock(true), createElementBlock(Fragment, null, renderList(fileTypeStats.value, (stat) => {
                      return openBlock(), createElementBlock("div", {
                        key: stat.extension,
                        class: "stats-popover-item"
                      }, [
                        createBaseVNode("span", _hoisted_19$a, toDisplayString(stat.icon), 1),
                        createBaseVNode("span", _hoisted_20$a, "." + toDisplayString(stat.extension), 1),
                        createBaseVNode("span", _hoisted_21$8, toDisplayString(stat.count), 1),
                        createBaseVNode("div", _hoisted_22$7, [
                          createBaseVNode("div", {
                            class: normalizeClass(["stats-popover-item-fill", stat.colorClass]),
                            style: normalizeStyle({ width: `${stat.percentage}%` })
                          }, null, 6)
                        ])
                      ]);
                    }), 128))
                  ])
                ])) : createCommentVNode("", true),
                folderStats.value.length > 0 ? (openBlock(), createElementBlock("div", _hoisted_23$7, [
                  createBaseVNode("div", _hoisted_24$7, toDisplayString(unref(t)("stats.byFolder")), 1),
                  createBaseVNode("div", _hoisted_25$6, [
                    (openBlock(true), createElementBlock(Fragment, null, renderList(folderStats.value.slice(0, 5), (stat) => {
                      return openBlock(), createElementBlock("div", {
                        key: stat.folder,
                        class: "stats-popover-item"
                      }, [
                        _cache[4] || (_cache[4] = createBaseVNode("svg", {
                          class: "w-4 h-4 text-blue-400 flex-shrink-0",
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
                        createBaseVNode("span", _hoisted_26$6, toDisplayString(stat.folder || "/"), 1),
                        createBaseVNode("span", _hoisted_27$5, toDisplayString(stat.count), 1),
                        createBaseVNode("div", _hoisted_28$5, [
                          createBaseVNode("div", {
                            class: "stats-popover-item-fill stats-popover-item-fill--folder",
                            style: normalizeStyle({ width: `${stat.percentage}%` })
                          }, null, 4)
                        ])
                      ]);
                    }), 128))
                  ])
                ])) : createCommentVNode("", true)
              ], 4)
            ])) : createCommentVNode("", true)
          ]),
          _: 1
        })
      ]);
    };
  }
});
const StatsPopover = /* @__PURE__ */ _export_sfc(_sfc_main$A, [["__scopeId", "data-v-42afbac0"]]);
const _hoisted_1$y = {
  key: 0,
  class: "drop-zone-overlay"
};
const _hoisted_2$x = { class: "drop-zone-content" };
const _hoisted_3$v = { class: "text-lg font-semibold text-white" };
const _hoisted_4$t = { class: "text-sm text-gray-400" };
const _sfc_main$z = /* @__PURE__ */ defineComponent({
  __name: "ContextPanel",
  setup(__props) {
    const logger2 = useLogger("ContextPanel");
    const contextStore = useContextStore();
    const settingsStore = useSettingsStore();
    const templateStore = useTemplateStore();
    const projectStore = useProjectStore();
    const uiStore = useUIStore();
    const { t } = useI18n();
    const contentRef = ref(null);
    const copySuccess = ref(false);
    const buildStatusText = computed(() => {
      const progress = contextStore.buildProgress;
      if (progress < 20) return t("context.statusAnalyzing");
      if (progress < 50) return t("context.statusReading");
      if (progress < 80) return t("context.statusProcessing");
      if (progress < 95) return t("context.statusFormatting");
      return t("context.statusFinalizing");
    });
    const chunking = useChunking();
    const showChunkHUD = computed(() => {
      if (!contextStore.hasContext) return false;
      if (!settingsStore.settings.context.enableAutoSplit) return false;
      const totalTokens = contextStore.tokenCount;
      const maxPerChunk = settingsStore.settings.context.maxTokensPerChunk;
      return totalTokens > maxPerChunk;
    });
    watch(
      () => [
        contextStore.tokenCount,
        contextStore.lineCount,
        contextStore.fileCount,
        settingsStore.settings.context.enableAutoSplit,
        settingsStore.settings.context.maxTokensPerChunk,
        settingsStore.settings.context.splitStrategy
      ],
      () => {
        if (!settingsStore.settings.context.enableAutoSplit) {
          chunking.setChunks([]);
          return;
        }
        const totalTokens = contextStore.tokenCount;
        const totalLines = contextStore.lineCount;
        const maxPerChunk = settingsStore.settings.context.maxTokensPerChunk;
        if (totalTokens <= 0 || totalLines <= 0) {
          chunking.setChunks([]);
          return;
        }
        const numChunks = Math.ceil(totalTokens / maxPerChunk);
        if (numChunks <= 1) {
          chunking.setChunks([{
            index: 0,
            startLine: 0,
            endLine: totalLines - 1,
            tokenCount: totalTokens
          }]);
          return;
        }
        const tokensPerChunk = Math.ceil(totalTokens / numChunks);
        const linesPerChunk = Math.ceil(totalLines / numChunks);
        const chunks = [];
        for (let i = 0; i < numChunks; i++) {
          const startLine = i * linesPerChunk;
          const endLine = Math.min((i + 1) * linesPerChunk - 1, totalLines - 1);
          const chunkTokens = i === numChunks - 1 ? totalTokens - tokensPerChunk * (numChunks - 1) : tokensPerChunk;
          chunks.push({ index: i, startLine, endLine, tokenCount: chunkTokens });
        }
        chunking.setChunks(chunks);
      },
      { immediate: true }
    );
    const showSearch = ref(false);
    const showStatsPopover = ref(false);
    const searchQuery = ref("");
    const searchResults = ref([]);
    const currentSearchIndex = ref(0);
    const highlightedLinesSet = computed(() => {
      if (searchResults.value.length === 0) return /* @__PURE__ */ new Set();
      const currentLine = searchResults.value[currentSearchIndex.value];
      return currentLine !== void 0 ? /* @__PURE__ */ new Set([currentLine]) : /* @__PURE__ */ new Set();
    });
    const chunkBoundaries = computed(() => {
      if (!showChunkHUD.value || chunking.chunks.value.length <= 1) return /* @__PURE__ */ new Set();
      return new Set(chunking.chunks.value.slice(1).map((c) => c.startLine));
    });
    const isDragging = ref(false);
    let dragCounter = 0;
    async function handleFormatChange(format) {
      settingsStore.updateContextSettings({ outputFormat: format });
      if (contextStore.contextId) {
        await contextStore.rebuildContext();
      }
    }
    watch(searchQuery, (query) => {
      if (!query || !contextStore.currentChunk?.lines) {
        searchResults.value = [];
        currentSearchIndex.value = 0;
        return;
      }
      const results = [];
      const lowerQuery = query.toLowerCase();
      contextStore.currentChunk.lines.forEach((line, index) => {
        if (line.toLowerCase().includes(lowerQuery)) {
          results.push(contextStore.currentChunk.startLine + index);
        }
      });
      searchResults.value = results;
      currentSearchIndex.value = 0;
      if (results.length > 0) {
        scrollToLine(results[0]);
      }
    });
    watch(showSearch, (show) => {
      if (!show) {
        searchQuery.value = "";
      }
    });
    function goToPrevChunk() {
      chunking.prevChunk();
      scrollToCurrentChunk();
    }
    function goToNextChunk() {
      chunking.nextChunk();
      scrollToCurrentChunk();
    }
    function scrollToCurrentChunk() {
      nextTick(() => {
        const chunkInfo = chunking.currentChunkInfo.value;
        if (!chunkInfo) return;
        contentRef.value?.scrollToLine(chunkInfo.startLine);
      });
    }
    function scrollToLine(lineNum) {
      nextTick(() => {
        contentRef.value?.scrollToLine(lineNum);
      });
    }
    function searchNext() {
      if (searchResults.value.length === 0) return;
      currentSearchIndex.value = (currentSearchIndex.value + 1) % searchResults.value.length;
      scrollToLine(searchResults.value[currentSearchIndex.value]);
    }
    function searchPrev() {
      if (searchResults.value.length === 0) return;
      currentSearchIndex.value = (currentSearchIndex.value - 1 + searchResults.value.length) % searchResults.value.length;
      scrollToLine(searchResults.value[currentSearchIndex.value]);
    }
    async function handleExport() {
      try {
        contentRef.value?.openExportModal();
      } catch (error) {
        logger2.error("Failed to open export modal:", error);
      }
    }
    async function handleCopyText() {
      if (!contextStore.contextId) return;
      try {
        const filesContent = await contextStore.getFullContextContent();
        let content;
        if (settingsStore.settings.context.applyTemplateOnCopy && templateStore.activeTemplate) {
          const files2 = contextStore.summary?.files || [];
          const templateContext = {
            fileTree: generateFileTree(files2, projectStore.projectName),
            files: filesContent,
            task: templateStore.currentTask,
            userRules: templateStore.userRules,
            fileCount: contextStore.fileCount,
            tokenCount: contextStore.tokenCount,
            languages: detectLanguages(files2),
            projectName: projectStore.projectName
          };
          content = templateStore.generatePrompt(templateContext);
        } else {
          content = filesContent;
        }
        await navigator.clipboard.writeText(content);
        showCopySuccess();
        uiStore.addToast(t("toast.contextCopied"), "success");
      } catch (error) {
        logger2.error("Failed to copy context:", error);
        uiStore.addToast(t("toast.copyError"), "error");
      }
    }
    async function handleCopyCurrentChunk() {
      if (!contextStore.contextId || !chunking.currentChunkInfo.value) return;
      try {
        const chunkInfo = chunking.currentChunkInfo.value;
        const lines = contextStore.currentChunk?.lines ?? [];
        const chunkContent = lines.slice(chunkInfo.startLine, chunkInfo.endLine + 1).join("\n");
        await navigator.clipboard.writeText(chunkContent);
        showCopySuccess();
        chunking.markCopied(chunkInfo.index, true);
        uiStore.addToast(t("chunks.copied"), "success");
      } catch (error) {
        logger2.error("Failed to copy chunk:", error);
        uiStore.addToast(t("toast.copyError"), "error");
      }
    }
    function showCopySuccess() {
      copySuccess.value = true;
      setTimeout(() => {
        copySuccess.value = false;
      }, 1500);
    }
    function handleDragOver(e) {
      dragCounter++;
      isDragging.value = true;
      if (e.dataTransfer?.types.includes("Files") || e.dataTransfer?.types.includes("text/plain")) {
        e.dataTransfer.dropEffect = "copy";
      }
    }
    function handleDragLeave() {
      dragCounter--;
      if (dragCounter <= 0) {
        dragCounter = 0;
        isDragging.value = false;
      }
    }
    async function handleDrop(e) {
      isDragging.value = false;
      dragCounter = 0;
      if (!e.dataTransfer) return;
      const textData = e.dataTransfer.getData("text/plain");
      if (textData) {
        const paths = textData.split("\n").filter((p) => p.trim());
        if (paths.length > 0) {
          logger2.debug("Dropped file paths:", paths.length);
          window.dispatchEvent(new CustomEvent("add-files-to-context", { detail: { paths } }));
        }
      }
    }
    return (_ctx, _cache) => {
      return openBlock(), createElementBlock("div", {
        class: "context-panel-root layout-fill layout-column layout-clip",
        "data-tour": "context-preview",
        onDragover: withModifiers(handleDragOver, ["prevent"]),
        onDragleave: handleDragLeave,
        onDrop: withModifiers(handleDrop, ["prevent"])
      }, [
        isDragging.value ? (openBlock(), createElementBlock("div", _hoisted_1$y, [
          createBaseVNode("div", _hoisted_2$x, [
            _cache[6] || (_cache[6] = createBaseVNode("svg", {
              class: "w-12 h-12 text-indigo-400 mb-3",
              fill: "none",
              stroke: "currentColor",
              viewBox: "0 0 24 24"
            }, [
              createBaseVNode("path", {
                "stroke-linecap": "round",
                "stroke-linejoin": "round",
                "stroke-width": "2",
                d: "M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M15 13l-3-3m0 0l-3 3m3-3v12"
              })
            ], -1)),
            createBaseVNode("p", _hoisted_3$v, toDisplayString(unref(t)("context.dropFiles")), 1),
            createBaseVNode("p", _hoisted_4$t, toDisplayString(unref(t)("context.dropFilesHint")), 1)
          ])
        ])) : createCommentVNode("", true),
        createVNode(ContextPanelHeader, {
          "has-context": unref(contextStore).hasContext,
          "show-search": showSearch.value,
          "file-count": unref(contextStore).fileCount,
          "line-count": unref(contextStore).lineCount,
          "token-count": unref(contextStore).tokenCount,
          "output-format": unref(settingsStore).settings.context.outputFormat,
          onToggleSearch: _cache[0] || (_cache[0] = ($event) => showSearch.value = !showSearch.value),
          onShowStats: _cache[1] || (_cache[1] = ($event) => showStatsPopover.value = true),
          onFormatChange: handleFormatChange
        }, null, 8, ["has-context", "show-search", "file-count", "line-count", "token-count", "output-format"]),
        createVNode(ContextPanelToolbar, {
          visible: showSearch.value,
          "search-query": searchQuery.value,
          "results-count": searchResults.value.length,
          "current-index": currentSearchIndex.value,
          "onUpdate:searchQuery": _cache[2] || (_cache[2] = ($event) => searchQuery.value = $event),
          onSearchNext: searchNext,
          onSearchPrev: searchPrev,
          onClose: _cache[3] || (_cache[3] = ($event) => showSearch.value = false)
        }, null, 8, ["visible", "search-query", "results-count", "current-index"]),
        createVNode(ContextPanelContent, {
          ref_key: "contentRef",
          ref: contentRef,
          "is-building": unref(contextStore).isBuilding,
          "build-progress": unref(contextStore).buildProgress,
          "status-text": buildStatusText.value,
          "file-count": unref(contextStore).fileCount,
          "line-count": unref(contextStore).lineCount,
          "total-size": unref(contextStore).totalSize,
          "context-id": unref(contextStore).contextId,
          error: unref(contextStore).error,
          "has-context": unref(contextStore).hasContext,
          "is-loading": unref(contextStore).isLoading,
          lines: unref(contextStore).currentChunk?.lines,
          "highlighted-lines": highlightedLinesSet.value,
          "search-query": searchQuery.value,
          "chunk-boundaries": showChunkHUD.value ? chunkBoundaries.value : void 0,
          "output-format": unref(settingsStore).settings.context.outputFormat
        }, null, 8, ["is-building", "build-progress", "status-text", "file-count", "line-count", "total-size", "context-id", "error", "has-context", "is-loading", "lines", "highlighted-lines", "search-query", "chunk-boundaries", "output-format"]),
        createVNode(ContextPanelFooter, {
          visible: unref(contextStore).hasContext,
          "show-chunk-nav": showChunkHUD.value,
          "current-chunk": unref(chunking).currentChunk.value,
          "total-chunks": unref(chunking).totalChunks.value,
          "copy-success": copySuccess.value,
          "is-chunk-copied": unref(chunking).isChunkCopied,
          onClear: unref(contextStore).clearContext,
          onExport: handleExport,
          onCopy: _cache[4] || (_cache[4] = ($event) => showChunkHUD.value ? handleCopyCurrentChunk() : handleCopyText()),
          onPrevChunk: goToPrevChunk,
          onNextChunk: goToNextChunk
        }, null, 8, ["visible", "show-chunk-nav", "current-chunk", "total-chunks", "copy-success", "is-chunk-copied", "onClear"]),
        createVNode(StatsPopover, {
          visible: showStatsPopover.value,
          onClose: _cache[5] || (_cache[5] = ($event) => showStatsPopover.value = false)
        }, null, 8, ["visible"])
      ], 32);
    };
  }
});
const ContextPanel = /* @__PURE__ */ _export_sfc(_sfc_main$z, [["__scopeId", "data-v-c3e1e2c9"]]);
const logger$5 = useLogger("ProjectStore");
const RECENT_PROJECTS_KEY = "Syntaxia_recent_projects";
const MAX_RECENT = 10;
const useProjectStore = defineStore("project", () => {
  const currentPath = ref(null);
  const currentName = ref(null);
  const recentProjects = ref([]);
  const isLoading = ref(false);
  const error = ref(null);
  const autoOpenLast = ref(false);
  const hasProject = computed(() => currentPath.value !== null);
  const projectName = computed(() => currentName.value || "");
  const projectPath = computed(() => currentPath.value || "");
  const hasRecentProjects = computed(() => recentProjects.value.length > 0);
  const lastProjectPath = computed(() => recentProjects.value[0]?.path || null);
  async function openProjectByPath(path) {
    isLoading.value = true;
    error.value = null;
    try {
      const exists = await apiService.pathExists(path);
      if (!exists) {
        error.value = `Path does not exist: ${path}`;
        recentProjects.value = recentProjects.value.filter((p) => p.path !== path);
        saveRecentProjects();
        return false;
      }
      currentPath.value = path;
      currentName.value = path.split(/[\\/]/).pop() || path;
      const fileStore = useFileStore();
      const contextStore = useContextStore();
      fileStore.resetStore();
      contextStore.clearContext();
      addToRecent(path);
      try {
        const name = path.split(/[\\/]/).pop() || path;
        await apiService.addRecentProject(path, name);
      } catch (backendError) {
        logger$5.error("Failed to save project to backend:", backendError);
      }
      saveRecentProjects();
      return true;
    } catch (err) {
      error.value = err instanceof Error ? err.message : "Failed to load project";
      return false;
    } finally {
      isLoading.value = false;
    }
  }
  function addToRecent(path) {
    const name = path.split(/[\\/]/).pop() || path;
    const existing = recentProjects.value.findIndex((p) => p.path === path);
    if (existing !== -1) {
      recentProjects.value.splice(existing, 1);
    }
    recentProjects.value.unshift({
      path,
      name,
      lastOpened: Date.now()
    });
    if (recentProjects.value.length > MAX_RECENT) {
      recentProjects.value = recentProjects.value.slice(0, MAX_RECENT);
    }
  }
  function loadRecentProjects() {
    try {
      const stored = localStorage.getItem(RECENT_PROJECTS_KEY);
      if (stored) {
        recentProjects.value = JSON.parse(stored);
      }
      const auto = localStorage.getItem("Syntaxia_auto_open_last");
      if (auto !== null) {
        autoOpenLast.value = auto === "true";
      }
    } catch (err) {
      console.warn("Failed to load recent projects:", err);
    }
  }
  async function fetchRecentProjects() {
    try {
      isLoading.value = true;
      error.value = null;
      const projectsJson = await apiService.getRecentProjects();
      if (projectsJson) {
        const projects = JSON.parse(projectsJson);
        if (Array.isArray(projects)) {
          recentProjects.value = projects;
          saveRecentProjects();
        }
      }
    } catch (err) {
      logger$5.error("Failed to fetch recent projects from backend:", err);
      error.value = err instanceof Error ? err.message : "Failed to load recent projects";
    } finally {
      isLoading.value = false;
    }
  }
  function saveRecentProjects() {
    try {
      localStorage.setItem(RECENT_PROJECTS_KEY, JSON.stringify(recentProjects.value));
    } catch (err) {
      console.warn("Failed to save recent projects:", err);
    }
  }
  function setAutoOpenLast(value) {
    autoOpenLast.value = value;
    try {
      localStorage.setItem("Syntaxia_auto_open_last", String(value));
    } catch (err) {
      console.warn("Failed to save auto-open setting:", err);
    }
  }
  async function maybeAutoOpenLastProject() {
    if (currentPath.value) return true;
    if (!autoOpenLast.value) return false;
    if (!lastProjectPath.value) return false;
    return await openProjectByPath(lastProjectPath.value);
  }
  async function removeFromRecent(path) {
    recentProjects.value = recentProjects.value.filter((p) => p.path !== path);
    try {
      await apiService.removeRecentProject(path);
    } catch (backendError) {
      logger$5.error("Failed to remove project from backend:", backendError);
    }
    saveRecentProjects();
  }
  function clearRecent() {
    recentProjects.value = [];
    saveRecentProjects();
  }
  function clearProject() {
    currentPath.value = null;
    currentName.value = null;
    error.value = null;
    const fileStore = useFileStore();
    const contextStore = useContextStore();
    fileStore.resetStore();
    contextStore.clearContext();
  }
  loadRecentProjects();
  return {
    // State
    currentPath,
    currentName,
    recentProjects,
    isLoading,
    error,
    autoOpenLast,
    // Computed
    hasProject,
    projectName,
    projectPath,
    hasRecentProjects,
    lastProjectPath,
    // Actions
    openProjectByPath,
    fetchRecentProjects,
    // Новый action для загрузки из бэкенда
    removeFromRecent,
    clearRecent,
    clearProject,
    setAutoOpenLast,
    maybeAutoOpenLastProject
  };
});
const _hoisted_1$x = { class: "absolute top-4 right-4 z-10 flex items-center gap-3" };
const _hoisted_2$w = ["title"];
const _hoisted_3$u = { class: "toggle-wrapper-sm" };
const _hoisted_4$s = ["checked"];
const _hoisted_5$o = ["title"];
const _hoisted_6$m = { class: "font-medium" };
const _hoisted_7$m = {
  key: 0,
  class: "drop-overlay"
};
const _hoisted_8$k = { class: "drop-content" };
const _hoisted_9$i = { class: "text-2xl font-bold text-white mt-4" };
const _hoisted_10$i = { class: "text-purple-200 mt-2" };
const _hoisted_11$g = { class: "content-wrapper" };
const _hoisted_12$e = { class: "header-section" };
const _hoisted_13$e = { class: "app-title" };
const _hoisted_14$d = { class: "app-subtitle" };
const _hoisted_15$b = { class: "app-hint" };
const _hoisted_16$a = { class: "cta-section" };
const _hoisted_17$a = { class: "recent-section" };
const _hoisted_18$9 = { class: "section-header" };
const _hoisted_19$9 = { class: "section-title" };
const _hoisted_20$9 = ["title"];
const _hoisted_21$7 = { class: "projects-list" };
const _hoisted_22$6 = ["onClick", "onContextmenu"];
const _hoisted_23$6 = ["title"];
const _hoisted_24$6 = { class: "project-name" };
const _hoisted_25$5 = { class: "project-path" };
const _hoisted_26$5 = ["onClick", "title"];
const _hoisted_27$4 = {
  key: 1,
  class: "empty-projects"
};
const _hoisted_28$4 = { class: "absolute bottom-4 left-4 z-10" };
const _sfc_main$y = /* @__PURE__ */ defineComponent({
  __name: "ProjectSelector",
  emits: ["opened"],
  setup(__props, { emit: __emit }) {
    const logger2 = useLogger("ProjectSelector");
    const emit = __emit;
    const projectStore = useProjectStore();
    const uiStore = useUIStore();
    const { t, locale, setLocale } = useI18n();
    const isDragging = ref(false);
    const contextMenu = ref({
      visible: false,
      x: 0,
      y: 0,
      project: null
    });
    const recentProjects = computed(() => projectStore.recentProjects);
    function shortenPath(path, maxLength = 50) {
      if (path.length <= maxLength) return path;
      const separator = path.includes("\\") ? "\\" : "/";
      const parts = path.split(separator);
      if (parts.length <= 3) return path;
      const first = parts[0];
      const last = parts.slice(-2).join(separator);
      const shortened = `${first}${separator}...${separator}${last}`;
      return shortened.length < path.length ? shortened : path;
    }
    onMounted(async () => {
      try {
        await projectStore.fetchRecentProjects();
      } catch (error) {
        logger2.error("Failed to load recent projects:", error);
      }
    });
    function toggleLanguage() {
      const newLocale = locale.value === "ru" ? "en" : "ru";
      setLocale(newLocale);
      uiStore.addToast(
        newLocale === "ru" ? "Язык изменён на русский" : "Language changed to English",
        "success"
      );
    }
    function onToggleAutoOpen(e) {
      const target2 = e.target;
      projectStore.setAutoOpenLast(target2.checked);
    }
    async function handleDrop(e) {
      isDragging.value = false;
      const items = e.dataTransfer?.items;
      if (!items) return;
      for (let i = 0; i < items.length; i++) {
        const item = items[i];
        if (item.kind === "file") {
          const entry = item.webkitGetAsEntry?.();
          if (entry?.isDirectory) {
            const path = entry.fullPath;
            const success = await projectStore.openProjectByPath(path);
            if (success) {
              emit("opened", path);
            }
            return;
          }
        }
      }
      uiStore.addToast("Please drop a folder, not a file", "warning");
    }
    async function selectProject() {
      try {
        const dirPath = await apiService.selectDirectory();
        if (!dirPath || dirPath === "") return;
        const success = await projectStore.openProjectByPath(dirPath);
        if (success && projectStore.currentPath) {
          emit("opened", projectStore.currentPath);
        }
      } catch (error) {
        logger2.error("Failed to select project:", error);
        const errorMessage = error instanceof Error ? error.message : "Unknown error";
        uiStore.addToast(`Failed to select directory: ${errorMessage}`, "error");
      }
    }
    async function openRecentProject(path) {
      try {
        await projectStore.openProjectByPath(path);
        emit("opened", path);
      } catch (error) {
        logger2.error("Failed to open recent project:", error);
        uiStore.addToast("Failed to open project", "error");
      }
    }
    async function removeProject(path) {
      try {
        await projectStore.removeFromRecent(path);
        uiStore.addToast(t("welcome.projectRemoved"), "success");
      } catch (error) {
        logger2.error("Failed to remove project:", error);
      }
    }
    async function clearAllHistory() {
      if (confirm(t("welcome.confirmClearHistory"))) {
        projectStore.clearRecent();
        uiStore.addToast(t("welcome.historyCleared"), "success");
      }
    }
    function showContextMenu(event, project) {
      contextMenu.value = { visible: true, x: event.clientX, y: event.clientY, project };
    }
    function hideContextMenu() {
      contextMenu.value.visible = false;
      contextMenu.value.project = null;
    }
    async function copyProjectPath() {
      if (contextMenu.value.project) {
        try {
          await navigator.clipboard.writeText(contextMenu.value.project.path);
          uiStore.addToast(t("welcome.pathCopied"), "success");
        } catch (error) {
          logger2.error("Failed to copy path:", error);
        }
      }
      hideContextMenu();
    }
    async function removeProjectFromMenu() {
      if (contextMenu.value.project) {
        await removeProject(contextMenu.value.project.path);
      }
      hideContextMenu();
    }
    if (typeof window !== "undefined") {
      window.addEventListener("click", hideContextMenu);
    }
    return (_ctx, _cache) => {
      return openBlock(), createElementBlock("div", {
        class: normalizeClass(["project-selector", { "drag-over": isDragging.value }]),
        onDrop: withModifiers(handleDrop, ["prevent"]),
        onDragover: _cache[1] || (_cache[1] = withModifiers(($event) => isDragging.value = true, ["prevent"])),
        onDragleave: _cache[2] || (_cache[2] = withModifiers(($event) => isDragging.value = false, ["prevent"]))
      }, [
        _cache[17] || (_cache[17] = createStaticVNode('<div class="bg-decoration" data-v-5260acfc><div class="bg-glow bg-glow-1" data-v-5260acfc></div><div class="bg-glow bg-glow-2" data-v-5260acfc></div><div class="bg-glow bg-glow-3" data-v-5260acfc></div><div class="bg-grid" data-v-5260acfc></div></div>', 1)),
        createBaseVNode("div", _hoisted_1$x, [
          createBaseVNode("label", {
            class: "auto-open-toggle",
            title: unref(t)("welcome.autoOpen")
          }, [
            createBaseVNode("div", _hoisted_3$u, [
              createBaseVNode("input", {
                type: "checkbox",
                class: "sr-only peer",
                checked: unref(projectStore).autoOpenLast,
                onChange: onToggleAutoOpen
              }, null, 40, _hoisted_4$s),
              _cache[3] || (_cache[3] = createBaseVNode("div", { class: "toggle-track-sm peer-checked:bg-gradient-to-r peer-checked:from-purple-600 peer-checked:to-pink-600" }, null, -1)),
              _cache[4] || (_cache[4] = createBaseVNode("div", { class: "toggle-thumb-sm peer-checked:translate-x-3" }, null, -1))
            ]),
            createBaseVNode("span", null, toDisplayString(unref(t)("welcome.autoOpenShort")), 1)
          ], 8, _hoisted_2$w),
          createBaseVNode("button", {
            onClick: toggleLanguage,
            class: "lang-switcher",
            title: unref(locale) === "ru" ? "Switch to English" : "Переключить на русский"
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
                d: "M3 5h12M9 3v2m1.048 9.5A18.022 18.022 0 016.412 9m6.088 9h7M11 21l5-10 5 10M12.751 5C11.783 10.77 8.07 15.61 3 18.129"
              })
            ], -1)),
            createBaseVNode("span", _hoisted_6$m, toDisplayString(unref(locale) === "ru" ? "RU" : "EN"), 1)
          ], 8, _hoisted_5$o)
        ]),
        createVNode(Transition, { name: "fade" }, {
          default: withCtx(() => [
            isDragging.value ? (openBlock(), createElementBlock("div", _hoisted_7$m, [
              createBaseVNode("div", _hoisted_8$k, [
                _cache[6] || (_cache[6] = createBaseVNode("div", { class: "drop-icon" }, [
                  createBaseVNode("svg", {
                    class: "w-12 h-12",
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
                  ])
                ], -1)),
                createBaseVNode("p", _hoisted_9$i, toDisplayString(unref(t)("welcome.dropHere")), 1),
                createBaseVNode("p", _hoisted_10$i, toDisplayString(unref(t)("welcome.toOpenProject")), 1)
              ])
            ])) : createCommentVNode("", true)
          ]),
          _: 1
        }),
        createBaseVNode("div", _hoisted_11$g, [
          createBaseVNode("div", _hoisted_12$e, [
            _cache[7] || (_cache[7] = createStaticVNode('<div class="logo-container" data-v-5260acfc><div class="logo-glow" data-v-5260acfc></div><div class="logo" data-v-5260acfc><svg class="w-8 h-8 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24" data-v-5260acfc><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 20l4-16m4 4l4 4-4 4M6 16l-4-4 4-4" data-v-5260acfc></path></svg></div></div>', 1)),
            createBaseVNode("h1", _hoisted_13$e, toDisplayString(unref(t)("welcome.title")), 1),
            createBaseVNode("p", _hoisted_14$d, toDisplayString(unref(t)("welcome.subtitle")), 1),
            createBaseVNode("p", _hoisted_15$b, toDisplayString(unref(t)("welcome.dragDrop")), 1)
          ]),
          createBaseVNode("div", _hoisted_16$a, [
            createBaseVNode("button", {
              onClick: selectProject,
              class: "cta-button"
            }, [
              _cache[8] || (_cache[8] = createBaseVNode("svg", {
                class: "w-5 h-5",
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
              createTextVNode(" " + toDisplayString(unref(t)("welcome.openProject")), 1)
            ])
          ]),
          createBaseVNode("div", _hoisted_17$a, [
            createBaseVNode("div", _hoisted_18$9, [
              createBaseVNode("h2", _hoisted_19$9, [
                _cache[9] || (_cache[9] = createBaseVNode("svg", {
                  class: "w-5 h-5 text-gray-400",
                  fill: "none",
                  stroke: "currentColor",
                  viewBox: "0 0 24 24"
                }, [
                  createBaseVNode("path", {
                    "stroke-linecap": "round",
                    "stroke-linejoin": "round",
                    "stroke-width": "2",
                    d: "M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"
                  })
                ], -1)),
                createTextVNode(" " + toDisplayString(unref(t)("welcome.recentProjects")), 1)
              ]),
              recentProjects.value.length > 0 ? (openBlock(), createElementBlock("button", {
                key: 0,
                onClick: clearAllHistory,
                class: "clear-btn",
                title: unref(t)("welcome.clearHistory")
              }, [
                _cache[10] || (_cache[10] = createBaseVNode("svg", {
                  class: "w-4 h-4",
                  fill: "none",
                  stroke: "currentColor",
                  viewBox: "0 0 24 24"
                }, [
                  createBaseVNode("path", {
                    "stroke-linecap": "round",
                    "stroke-linejoin": "round",
                    "stroke-width": "2",
                    d: "M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"
                  })
                ], -1)),
                createTextVNode(" " + toDisplayString(unref(t)("welcome.clearHistory")), 1)
              ], 8, _hoisted_20$9)) : createCommentVNode("", true)
            ]),
            createBaseVNode("div", _hoisted_21$7, [
              recentProjects.value.length > 0 ? (openBlock(true), createElementBlock(Fragment, { key: 0 }, renderList(recentProjects.value, (project, index) => {
                return openBlock(), createElementBlock("div", {
                  key: project.path,
                  class: "project-card",
                  style: normalizeStyle({ animationDelay: `${index * 60}ms` }),
                  onClick: ($event) => openRecentProject(project.path),
                  onContextmenu: withModifiers(($event) => showContextMenu($event, project), ["prevent"])
                }, [
                  _cache[12] || (_cache[12] = createBaseVNode("div", { class: "project-icon" }, [
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
                        d: "M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-6l-2-2H5a2 2 0 00-2 2z"
                      })
                    ])
                  ], -1)),
                  createBaseVNode("div", {
                    class: "project-info",
                    title: project.path
                  }, [
                    createBaseVNode("div", _hoisted_24$6, toDisplayString(project.name), 1),
                    createBaseVNode("div", _hoisted_25$5, toDisplayString(shortenPath(project.path)), 1)
                  ], 8, _hoisted_23$6),
                  createBaseVNode("button", {
                    onClick: withModifiers(($event) => removeProject(project.path), ["stop"]),
                    class: "project-remove",
                    title: unref(t)("welcome.removeProject")
                  }, [..._cache[11] || (_cache[11] = [
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
                        d: "M6 18L18 6M6 6l12 12"
                      })
                    ], -1)
                  ])], 8, _hoisted_26$5),
                  _cache[13] || (_cache[13] = createBaseVNode("svg", {
                    class: "project-arrow",
                    fill: "none",
                    stroke: "currentColor",
                    viewBox: "0 0 24 24"
                  }, [
                    createBaseVNode("path", {
                      "stroke-linecap": "round",
                      "stroke-linejoin": "round",
                      "stroke-width": "2",
                      d: "M9 5l7 7-7 7"
                    })
                  ], -1))
                ], 44, _hoisted_22$6);
              }), 128)) : (openBlock(), createElementBlock("div", _hoisted_27$4, [
                _cache[14] || (_cache[14] = createBaseVNode("div", { class: "empty-icon" }, [
                  createBaseVNode("svg", {
                    class: "w-6 h-6",
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
                  ])
                ], -1)),
                createBaseVNode("p", null, toDisplayString(unref(t)("welcome.noRecentProjects")), 1)
              ]))
            ])
          ]),
          (openBlock(), createBlock(Teleport, { to: "body" }, [
            createVNode(Transition, { name: "fade" }, {
              default: withCtx(() => [
                contextMenu.value.visible ? (openBlock(), createElementBlock("div", {
                  key: 0,
                  class: "context-menu",
                  style: normalizeStyle({ left: contextMenu.value.x + "px", top: contextMenu.value.y + "px" }),
                  onClick: _cache[0] || (_cache[0] = withModifiers(() => {
                  }, ["stop"]))
                }, [
                  createBaseVNode("button", {
                    onClick: copyProjectPath,
                    class: "context-menu-item"
                  }, [
                    _cache[15] || (_cache[15] = createBaseVNode("svg", {
                      class: "w-4 h-4",
                      fill: "none",
                      stroke: "currentColor",
                      viewBox: "0 0 24 24"
                    }, [
                      createBaseVNode("path", {
                        "stroke-linecap": "round",
                        "stroke-linejoin": "round",
                        "stroke-width": "2",
                        d: "M8 5H6a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2v-1M8 5a2 2 0 002 2h2a2 2 0 002-2M8 5a2 2 0 012-2h2a2 2 0 012 2m0 0h2a2 2 0 012 2v3m2 4H10m0 0l3-3m-3 3l3 3"
                      })
                    ], -1)),
                    createTextVNode(" " + toDisplayString(unref(t)("welcome.copyPath")), 1)
                  ]),
                  createBaseVNode("button", {
                    onClick: removeProjectFromMenu,
                    class: "context-menu-item context-menu-item-danger"
                  }, [
                    _cache[16] || (_cache[16] = createBaseVNode("svg", {
                      class: "w-4 h-4",
                      fill: "none",
                      stroke: "currentColor",
                      viewBox: "0 0 24 24"
                    }, [
                      createBaseVNode("path", {
                        "stroke-linecap": "round",
                        "stroke-linejoin": "round",
                        "stroke-width": "2",
                        d: "M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"
                      })
                    ], -1)),
                    createTextVNode(" " + toDisplayString(unref(t)("welcome.removeProject")), 1)
                  ])
                ], 4)) : createCommentVNode("", true)
              ]),
              _: 1
            })
          ]))
        ]),
        createBaseVNode("div", _hoisted_28$4, [
          createVNode(VersionBadge)
        ])
      ], 34);
    };
  }
});
const ProjectSelector = /* @__PURE__ */ _export_sfc(_sfc_main$y, [["__scopeId", "data-v-5260acfc"]]);
const DEFAULT_MAX_PANEL_WIDTH_PERCENT = 0.4;
function getEffectiveMaxWidth(maxWidth, maxWidthPercent) {
  const viewportMax = Math.floor(window.innerWidth * maxWidthPercent);
  return Math.min(maxWidth, viewportMax);
}
function useResizablePanel(options2) {
  const {
    minWidth = 200,
    maxWidth = 800,
    defaultWidth = 300,
    storageKey,
    invertDirection = false,
    maxWidthPercent = DEFAULT_MAX_PANEL_WIDTH_PERCENT
  } = options2;
  const panelRef = ref();
  const width = useStorage(storageKey, defaultWidth, localStorage, {
    mergeDefaults: true,
    serializer: {
      read: (v) => {
        const parsed = Number(v);
        if (isNaN(parsed)) return defaultWidth;
        const effectiveMax = getEffectiveMaxWidth(maxWidth, maxWidthPercent);
        return Math.max(minWidth, Math.min(effectiveMax, parsed));
      },
      write: (v) => String(v)
    }
  });
  const isResizing = ref(false);
  const validateWidth = () => {
    const effectiveMax = getEffectiveMaxWidth(maxWidth, maxWidthPercent);
    if (width.value > effectiveMax) {
      width.value = effectiveMax;
    }
  };
  onMounted(() => {
    validateWidth();
    window.addEventListener("resize", validateWidth);
  });
  let startX = 0;
  let startWidth = 0;
  const onMouseDown = (e) => {
    if (!panelRef.value) return;
    isResizing.value = true;
    startX = e.clientX;
    startWidth = width.value;
    document.body.style.userSelect = "none";
    document.body.style.cursor = "col-resize";
    document.addEventListener("mousemove", onMouseMove);
    document.addEventListener("mouseup", onMouseUp);
  };
  const onMouseMove = (e) => {
    if (!isResizing.value) return;
    const delta = e.clientX - startX;
    const adjustedDelta = invertDirection ? -delta : delta;
    const effectiveMax = getEffectiveMaxWidth(maxWidth, maxWidthPercent);
    const newWidth = Math.max(minWidth, Math.min(effectiveMax, startWidth + adjustedDelta));
    width.value = newWidth;
  };
  const onMouseUp = () => {
    if (!isResizing.value) return;
    isResizing.value = false;
    document.body.style.userSelect = "";
    document.body.style.cursor = "";
    document.removeEventListener("mousemove", onMouseMove);
    document.removeEventListener("mouseup", onMouseUp);
  };
  onUnmounted(() => {
    document.removeEventListener("mousemove", onMouseMove);
    document.removeEventListener("mouseup", onMouseUp);
    window.removeEventListener("resize", validateWidth);
  });
  const resetToDefault = () => {
    width.value = defaultWidth;
  };
  return {
    panelRef,
    width,
    isResizing,
    onMouseDown,
    resetToDefault,
    defaultWidth
  };
}
const _hoisted_1$w = { class: "action-bar" };
const _hoisted_2$v = { class: "action-bar-left" };
const _hoisted_3$t = ["title"];
const _hoisted_4$r = { class: "project-name" };
const _hoisted_5$n = { class: "template-quick-switch" };
const _hoisted_6$l = ["onClick", "title"];
const _hoisted_7$l = { class: "action-bar-right" };
const _hoisted_8$j = ["title"];
const _hoisted_9$h = ["title"];
const _hoisted_10$h = ["title"];
const _sfc_main$x = /* @__PURE__ */ defineComponent({
  __name: "ActionBar",
  emits: ["open-export", "reset-layout"],
  setup(__props, { emit: __emit }) {
    const projectStore = useProjectStore();
    const templateStore = useTemplateStore();
    const uiStore = useUIStore();
    const { t, locale, setLocale } = useI18n();
    const quickTemplates = computed(() => {
      const favs = templateStore.favoriteTemplates;
      if (favs.length > 0) return favs.slice(0, 4);
      return templateStore.visibleTemplates.slice(0, 4);
    });
    function toggleLanguage() {
      const newLocale = locale.value === "ru" ? "en" : "ru";
      setLocale(newLocale);
      uiStore.addToast(newLocale === "ru" ? "Язык изменён на русский" : "Language changed to English", "success");
    }
    function changeProject() {
      projectStore.clearProject();
      uiStore.addToast(t("action.changeProject"), "info");
    }
    function openSettings() {
      uiStore.openSettingsModal();
    }
    const emit = __emit;
    async function handleResetLayout() {
      emit("reset-layout");
      try {
        if (window.go?.main?.App?.ResetWindowState) {
          await window.go.main.App.ResetWindowState();
        }
      } catch (error) {
        console.warn("Failed to reset window state:", error);
      }
    }
    return (_ctx, _cache) => {
      return openBlock(), createElementBlock("div", _hoisted_1$w, [
        createBaseVNode("div", _hoisted_2$v, [
          createBaseVNode("button", {
            onClick: changeProject,
            class: "project-btn",
            title: unref(t)("hotkey.changeProject")
          }, [
            _cache[0] || (_cache[0] = createBaseVNode("div", { class: "project-icon" }, [
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
                  d: "M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-6l-2-2H5a2 2 0 00-2 2z"
                })
              ])
            ], -1)),
            createBaseVNode("span", _hoisted_4$r, toDisplayString(unref(projectStore).projectName), 1),
            _cache[1] || (_cache[1] = createBaseVNode("svg", {
              class: "w-3 h-3 text-gray-400 group-hover:text-gray-300",
              fill: "none",
              stroke: "currentColor",
              viewBox: "0 0 24 24"
            }, [
              createBaseVNode("path", {
                "stroke-linecap": "round",
                "stroke-linejoin": "round",
                "stroke-width": "2",
                d: "M19 9l-7 7-7-7"
              })
            ], -1))
          ], 8, _hoisted_3$t),
          createBaseVNode("div", _hoisted_5$n, [
            (openBlock(true), createElementBlock(Fragment, null, renderList(quickTemplates.value, (tpl) => {
              return openBlock(), createElementBlock("button", {
                key: tpl.id,
                onClick: ($event) => unref(templateStore).setActiveTemplate(tpl.id),
                class: normalizeClass(["quick-tpl-btn", { active: tpl.id === unref(templateStore).activeTemplateId }]),
                title: tpl.name
              }, toDisplayString(tpl.icon), 11, _hoisted_6$l);
            }), 128))
          ])
        ]),
        _cache[4] || (_cache[4] = createBaseVNode("div", { class: "action-bar-separator" }, null, -1)),
        createBaseVNode("div", _hoisted_7$l, [
          createBaseVNode("button", {
            onClick: openSettings,
            class: "toolbar-btn",
            title: unref(t)("settings.modal.title") + " (Ctrl+,)"
          }, [..._cache[2] || (_cache[2] = [
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
                d: "M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z"
              }),
              createBaseVNode("path", {
                "stroke-linecap": "round",
                "stroke-linejoin": "round",
                "stroke-width": "2",
                d: "M15 12a3 3 0 11-6 0 3 3 0 016 0z"
              })
            ], -1)
          ])], 8, _hoisted_8$j),
          createBaseVNode("button", {
            onClick: handleResetLayout,
            class: "toolbar-btn",
            title: unref(t)("workspace.resetLayout")
          }, [..._cache[3] || (_cache[3] = [
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
                d: "M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"
              })
            ], -1)
          ])], 8, _hoisted_9$h),
          createBaseVNode("button", {
            onClick: toggleLanguage,
            class: "toolbar-btn",
            title: unref(locale) === "ru" ? "Switch to English" : "Переключить на русский"
          }, toDisplayString(unref(locale).toUpperCase()), 9, _hoisted_10$h)
        ])
      ]);
    };
  }
});
const ActionBar = /* @__PURE__ */ _export_sfc(_sfc_main$x, [["__scopeId", "data-v-65fdffd2"]]);
const logger$4 = useLogger("TaskStore");
const useTaskStore = defineStore("task", () => {
  const taskDescription = ref("");
  const taskType = ref("feature");
  const isAnalyzing = ref(false);
  const analysisResult = ref(null);
  const suggestions = ref([
    "Add user authentication",
    "Implement data caching",
    "Refactor API calls",
    "Add unit tests"
  ]);
  const error = ref(null);
  let saveTimeout = null;
  watch(taskDescription, () => {
    if (saveTimeout) clearTimeout(saveTimeout);
    saveTimeout = setTimeout(() => {
      saveTaskDraft();
    }, 500);
  });
  async function analyzeTask() {
    if (taskDescription.value.trim().length < 10) {
      error.value = "Task description is too short";
      return;
    }
    isAnalyzing.value = true;
    error.value = null;
    try {
      const projectPath = await apiService.getCurrentDirectory();
      const fileStore = useFileStore();
      if (fileStore.nodes.length === 0) {
        await fileStore.loadFileTree(projectPath);
      }
      const domainFiles = convertToDomainFiles(fileStore.nodes);
      const suggestedFilePaths = await apiService.suggestContextFiles(
        taskDescription.value,
        domainFiles
      );
      let complexity = "low";
      if (suggestedFilePaths.length > 10) {
        complexity = "high";
      } else if (suggestedFilePaths.length > 5) {
        complexity = "medium";
      }
      const estimatedTime = suggestedFilePaths.length * 5 + (complexity === "high" ? 30 : complexity === "medium" ? 15 : 5);
      analysisResult.value = {
        suggestedFiles: suggestedFilePaths,
        complexity,
        estimatedTime,
        recommendations: generateRecommendations(complexity, suggestedFilePaths.length)
      };
      if (suggestedFilePaths.length > 0) {
        fileStore.clearSelection();
        fileStore.selectMultiple(suggestedFilePaths);
      }
    } catch (err) {
      error.value = err instanceof Error ? err.message : "Failed to analyze task";
      throw err;
    } finally {
      isAnalyzing.value = false;
    }
  }
  function convertToDomainFiles(nodes) {
    return nodes.map((node) => ({
      name: node.name,
      path: node.path,
      relPath: node.path,
      isDir: node.isDir,
      children: node.children ? convertToDomainFiles(node.children) : void 0,
      isGitignored: node.isIgnored || false,
      isCustomIgnored: false,
      isIgnored: node.isIgnored || false,
      size: node.size || 0
    }));
  }
  function generateRecommendations(complexity, fileCount) {
    const recommendations = [];
    if (complexity === "high") {
      recommendations.push("Consider breaking down into smaller tasks");
      recommendations.push("This is a complex task that may require multiple iterations");
    }
    if (fileCount > 15) {
      recommendations.push("Large number of files involved - ensure proper testing");
    }
    if (fileCount === 0) {
      recommendations.push("No relevant files found - try refining your task description");
    } else {
      recommendations.push(`${fileCount} relevant files identified`);
    }
    recommendations.push("Review existing similar implementations");
    recommendations.push("Add tests for new functionality");
    return recommendations;
  }
  function saveTaskDraft() {
    try {
      const draft = {
        description: taskDescription.value,
        type: taskType.value,
        timestamp: (/* @__PURE__ */ new Date()).toISOString()
      };
      localStorage.setItem("task-draft", JSON.stringify(draft));
    } catch (err) {
      logger$4.warn("Failed to save task draft:", err);
    }
  }
  function loadTaskDraft() {
    try {
      const draftJson = localStorage.getItem("task-draft");
      if (draftJson) {
        const draft = JSON.parse(draftJson);
        taskDescription.value = draft.description;
        taskType.value = draft.type || "feature";
      }
    } catch (err) {
      logger$4.warn("Failed to load task draft:", err);
    }
  }
  function clearTask() {
    taskDescription.value = "";
    taskType.value = "feature";
    analysisResult.value = null;
    error.value = null;
    localStorage.removeItem("task-draft");
  }
  function applySuggestion(suggestion) {
    taskDescription.value = suggestion;
  }
  return {
    // State
    taskDescription,
    taskType,
    isAnalyzing,
    analysisResult,
    suggestions,
    error,
    // Actions
    analyzeTask,
    saveTaskDraft,
    loadTaskDraft,
    clearTask,
    applySuggestion
  };
});
const _hoisted_1$v = { class: "h-full flex flex-col bg-gray-800/50 backdrop-blur-sm p-4" };
const _hoisted_2$u = { class: "flex items-center justify-between mb-3" };
const _hoisted_3$s = { class: "flex items-center gap-2" };
const _hoisted_4$q = { class: "text-lg font-semibold text-white" };
const _hoisted_5$m = ["placeholder"];
const _hoisted_6$k = {
  key: 0,
  class: "mt-3 space-y-2"
};
const _hoisted_7$k = { class: "text-xs text-gray-400 font-medium" };
const _hoisted_8$i = { class: "flex flex-wrap gap-2" };
const _hoisted_9$g = ["onClick"];
const _hoisted_10$g = {
  key: 1,
  class: "mt-3 p-3 bg-gray-700/50 rounded-lg text-sm"
};
const _hoisted_11$f = { class: "text-gray-300 font-medium mb-2" };
const _hoisted_12$d = { class: "space-y-1 text-xs text-gray-400" };
const _hoisted_13$d = { class: "text-white" };
const _hoisted_14$c = { key: 0 };
const _hoisted_15$a = { class: "text-white" };
const _hoisted_16$9 = {
  key: 1,
  class: "mt-2"
};
const _hoisted_17$9 = { class: "list-disc list-inside space-y-0.5" };
const _hoisted_18$8 = { class: "mt-4 flex gap-2" };
const _hoisted_19$8 = ["disabled"];
const _hoisted_20$8 = {
  key: 0,
  class: "h-4 w-4",
  fill: "none",
  stroke: "currentColor",
  viewBox: "0 0 24 24"
};
const _hoisted_21$6 = {
  key: 1,
  class: "animate-spin h-4 w-4",
  fill: "none",
  viewBox: "0 0 24 24"
};
const _sfc_main$w = /* @__PURE__ */ defineComponent({
  __name: "TaskPanel",
  setup(__props) {
    const { t } = useI18n();
    const taskStore = useTaskStore();
    const uiStore = useUIStore();
    const textareaRef = ref();
    const canAnalyze = computed(() => {
      return taskStore.taskDescription.trim().length > 10 && !taskStore.isAnalyzing;
    });
    function autoResize() {
      if (!textareaRef.value) return;
      textareaRef.value.style.height = "auto";
      textareaRef.value.style.height = textareaRef.value.scrollHeight + "px";
    }
    async function handleAnalyze() {
      if (!canAnalyze.value) return;
      try {
        await taskStore.analyzeTask();
        uiStore.addToast("Task analyzed successfully", "success");
      } catch (error) {
        uiStore.addToast("Failed to analyze task", "error");
      }
    }
    onMounted(() => {
      taskStore.loadTaskDraft();
      autoResize();
    });
    return (_ctx, _cache) => {
      return openBlock(), createElementBlock("div", _hoisted_1$v, [
        createBaseVNode("div", _hoisted_2$u, [
          createBaseVNode("div", _hoisted_3$s, [
            _cache[1] || (_cache[1] = createBaseVNode("div", { class: "w-1 h-5 bg-blue-500 rounded-full" }, null, -1)),
            createBaseVNode("h2", _hoisted_4$q, toDisplayString(unref(t)("task.title")), 1)
          ]),
          createBaseVNode("span", {
            class: normalizeClass(["text-xs text-gray-400", { "text-yellow-400": unref(taskStore).taskDescription.length > 4500, "text-red-400": unref(taskStore).taskDescription.length >= 5e3 }])
          }, toDisplayString(unref(taskStore).taskDescription.length) + " / 5000 ", 3)
        ]),
        withDirectives(createBaseVNode("textarea", {
          ref_key: "textareaRef",
          ref: textareaRef,
          "onUpdate:modelValue": _cache[0] || (_cache[0] = ($event) => unref(taskStore).taskDescription = $event),
          placeholder: unref(t)("task.placeholder"),
          class: "task-textarea flex-1 w-full px-4 py-3 bg-gray-800/90 border border-gray-700/50 rounded-lg text-white placeholder-gray-500 focus:outline-none focus:border-blue-500/50 focus:bg-gray-800 resize-none transition-all duration-200",
          maxlength: "5000",
          onInput: autoResize
        }, null, 40, _hoisted_5$m), [
          [vModelText, unref(taskStore).taskDescription]
        ]),
        unref(taskStore).suggestions.length > 0 && !unref(taskStore).taskDescription ? (openBlock(), createElementBlock("div", _hoisted_6$k, [
          createBaseVNode("p", _hoisted_7$k, toDisplayString(unref(t)("task.quickSuggestions")), 1),
          createBaseVNode("div", _hoisted_8$i, [
            (openBlock(true), createElementBlock(Fragment, null, renderList(unref(taskStore).suggestions, (suggestion) => {
              return openBlock(), createElementBlock("button", {
                key: suggestion,
                onClick: ($event) => unref(taskStore).applySuggestion(suggestion),
                class: "chip chip-default"
              }, toDisplayString(suggestion), 9, _hoisted_9$g);
            }), 128))
          ])
        ])) : createCommentVNode("", true),
        unref(taskStore).analysisResult ? (openBlock(), createElementBlock("div", _hoisted_10$g, [
          createBaseVNode("p", _hoisted_11$f, toDisplayString(unref(t)("task.analysisResult")), 1),
          createBaseVNode("div", _hoisted_12$d, [
            createBaseVNode("p", null, [
              createTextVNode(toDisplayString(unref(t)("task.complexity")) + ": ", 1),
              createBaseVNode("span", _hoisted_13$d, toDisplayString(unref(taskStore).analysisResult.complexity), 1)
            ]),
            unref(taskStore).analysisResult.estimatedTime ? (openBlock(), createElementBlock("p", _hoisted_14$c, [
              _cache[2] || (_cache[2] = createTextVNode(" Estimated time: ", -1)),
              createBaseVNode("span", _hoisted_15$a, toDisplayString(unref(taskStore).analysisResult.estimatedTime) + "min", 1)
            ])) : createCommentVNode("", true),
            unref(taskStore).analysisResult.recommendations.length > 0 ? (openBlock(), createElementBlock("div", _hoisted_16$9, [
              _cache[3] || (_cache[3] = createBaseVNode("p", { class: "font-medium text-gray-300 mb-1" }, "Recommendations:", -1)),
              createBaseVNode("ul", _hoisted_17$9, [
                (openBlock(true), createElementBlock(Fragment, null, renderList(unref(taskStore).analysisResult.recommendations, (rec) => {
                  return openBlock(), createElementBlock("li", { key: rec }, toDisplayString(rec), 1);
                }), 128))
              ])
            ])) : createCommentVNode("", true)
          ])
        ])) : createCommentVNode("", true),
        createBaseVNode("div", _hoisted_18$8, [
          createBaseVNode("button", {
            onClick: handleAnalyze,
            disabled: !canAnalyze.value,
            class: "btn btn-primary flex-1 py-2"
          }, [
            !unref(taskStore).isAnalyzing ? (openBlock(), createElementBlock("svg", _hoisted_20$8, [..._cache[4] || (_cache[4] = [
              createBaseVNode("path", {
                "stroke-linecap": "round",
                "stroke-linejoin": "round",
                "stroke-width": "2",
                d: "M9.663 17h4.673M12 3v1m6.364 1.636l-.707.707M21 12h-1M4 12H3m3.343-5.657l-.707-.707m2.828 9.9a5 5 0 117.072 0l-.548.547A3.374 3.374 0 0014 18.469V19a2 2 0 11-4 0v-.531c0-.895-.356-1.754-.988-2.386l-.548-.547z"
              }, null, -1)
            ])])) : (openBlock(), createElementBlock("svg", _hoisted_21$6, [..._cache[5] || (_cache[5] = [
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
                d: "M4 12a8 8 0 018-8V0C5.373 0 5.373 0 12h4z"
              }, null, -1)
            ])])),
            createBaseVNode("span", null, toDisplayString(unref(taskStore).isAnalyzing ? "Analyzing..." : "Analyze Task"), 1)
          ], 8, _hoisted_19$8)
        ])
      ]);
    };
  }
});
const TaskPanel = /* @__PURE__ */ _export_sfc(_sfc_main$w, [["__scopeId", "data-v-17134cbf"]]);
const _hoisted_1$u = { class: "center-container layout-full layout-column layout-clip" };
const _hoisted_2$t = {
  key: 0,
  class: "task-panel"
};
const _sfc_main$v = /* @__PURE__ */ defineComponent({
  __name: "CenterWorkspace",
  setup(__props) {
    const showTaskPanel = ref(false);
    watch(showTaskPanel, (visible) => {
      try {
        localStorage.setItem("task-panel-visible", visible.toString());
      } catch (err) {
        console.warn("Failed to save task panel visibility:", err);
      }
    });
    try {
      const savedVisibility = localStorage.getItem("task-panel-visible");
      if (savedVisibility) {
        showTaskPanel.value = savedVisibility === "true";
      } else {
        showTaskPanel.value = false;
      }
    } catch (err) {
      console.warn("Failed to load task panel visibility:", err);
    }
    return (_ctx, _cache) => {
      return openBlock(), createElementBlock("div", _hoisted_1$u, [
        showTaskPanel.value ? (openBlock(), createElementBlock("div", _hoisted_2$t, [
          createVNode(unref(TaskPanel))
        ])) : createCommentVNode("", true),
        createBaseVNode("div", {
          class: normalizeClass(["layout-fill layout-column layout-clip", showTaskPanel.value ? "context-panel-with-task" : "context-panel-full"]),
          "data-tour": "context-preview"
        }, [
          createVNode(unref(ContextPanel))
        ], 2)
      ]);
    };
  }
});
const CenterWorkspace = /* @__PURE__ */ _export_sfc(_sfc_main$v, [["__scopeId", "data-v-e110c943"]]);
const _hoisted_1$t = { class: "relative w-full max-w-5xl max-h-[90vh] bg-gray-900 rounded-xl border border-gray-700 shadow-2xl flex flex-col overflow-hidden" };
const _hoisted_2$s = { class: "flex items-center justify-between p-4 border-b border-gray-700" };
const _hoisted_3$r = { class: "flex items-center gap-3" };
const _hoisted_4$p = { class: "text-sm font-semibold text-white" };
const _hoisted_5$l = { class: "p-4 border-b border-gray-700 flex items-center gap-4" };
const _hoisted_6$j = { class: "flex-1" };
const _hoisted_7$j = { class: "block text-xs text-gray-400 mb-1" };
const _hoisted_8$h = ["value"];
const _hoisted_9$f = { class: "flex-1" };
const _hoisted_10$f = { class: "block text-xs text-gray-400 mb-1" };
const _hoisted_11$e = ["value"];
const _hoisted_12$c = ["title"];
const _hoisted_13$c = { class: "flex-1 overflow-auto" };
const _hoisted_14$b = {
  key: 0,
  class: "flex items-center justify-center h-64"
};
const _hoisted_15$9 = {
  key: 1,
  class: "flex items-center justify-center h-64 text-red-400"
};
const _hoisted_16$8 = { class: "text-center" };
const _hoisted_17$8 = {
  key: 2,
  class: "flex items-center justify-center h-64 text-gray-400"
};
const _hoisted_18$7 = { class: "text-center" };
const _hoisted_19$7 = {
  key: 3,
  class: "divide-y divide-gray-700"
};
const _hoisted_20$7 = { class: "flex items-center gap-3" };
const _hoisted_21$5 = { class: "text-sm text-white flex-1 truncate font-mono" };
const _hoisted_22$5 = {
  key: 0,
  class: "text-xs text-emerald-400"
};
const _hoisted_23$5 = {
  key: 1,
  class: "text-xs text-red-400"
};
const _hoisted_24$5 = { class: "px-4 py-3 border-t border-gray-700 flex items-center justify-between" };
const _hoisted_25$4 = { class: "text-xs text-gray-400" };
const _hoisted_26$4 = { key: 0 };
const _hoisted_27$3 = {
  key: 0,
  class: "text-emerald-400 ml-2"
};
const _hoisted_28$3 = {
  key: 1,
  class: "text-red-400 ml-1"
};
const _sfc_main$u = /* @__PURE__ */ defineComponent({
  __name: "BranchDiffModal",
  props: {
    isOpen: { type: Boolean },
    branches: {},
    projectPath: {},
    currentBranch: {}
  },
  emits: ["close"],
  setup(__props, { emit: __emit }) {
    const logger2 = useLogger("BranchDiffModal");
    const props = __props;
    const emit = __emit;
    const { t } = useI18n();
    const baseBranch = ref("");
    const compareBranch = ref("");
    const diffFiles = ref([]);
    const isLoading = ref(false);
    const error = ref("");
    const totalAdditions = computed(() => diffFiles.value.reduce((sum, f) => sum + (f.additions || 0), 0));
    const totalDeletions = computed(() => diffFiles.value.reduce((sum, f) => sum + (f.deletions || 0), 0));
    watch(() => props.isOpen, (isOpen) => {
      if (isOpen && props.branches.length >= 2) {
        baseBranch.value = props.currentBranch || props.branches[0];
        compareBranch.value = props.branches.find((b) => b !== baseBranch.value) || props.branches[1];
        loadDiff();
      }
    });
    function close() {
      emit("close");
    }
    function swapBranches() {
      const temp = baseBranch.value;
      baseBranch.value = compareBranch.value;
      compareBranch.value = temp;
      loadDiff();
    }
    async function loadDiff() {
      if (!baseBranch.value || !compareBranch.value || baseBranch.value === compareBranch.value) {
        diffFiles.value = [];
        return;
      }
      isLoading.value = true;
      error.value = "";
      diffFiles.value = [];
      try {
        const [baseFiles, compareFiles] = await Promise.all([
          apiService.listFilesAtRef(props.projectPath, baseBranch.value),
          apiService.listFilesAtRef(props.projectPath, compareBranch.value)
        ]);
        const baseSet = new Set(baseFiles);
        const compareSet = new Set(compareFiles);
        const allFiles = /* @__PURE__ */ new Set([...baseFiles, ...compareFiles]);
        const diff = [];
        for (const file of allFiles) {
          const inBase = baseSet.has(file);
          const inCompare = compareSet.has(file);
          if (!inBase && inCompare) {
            diff.push({ path: file, status: "added" });
          } else if (inBase && !inCompare) {
            diff.push({ path: file, status: "deleted" });
          } else {
            diff.push({ path: file, status: "modified" });
          }
        }
        diffFiles.value = diff.filter((f) => f.status !== "modified" || Math.random() > 0.7).sort((a, b) => {
          const order = { added: 0, modified: 1, deleted: 2, renamed: 3 };
          return order[a.status] - order[b.status];
        });
      } catch (err) {
        logger2.error("Failed to load diff:", err);
        error.value = t("error.loadFailed");
      } finally {
        isLoading.value = false;
      }
    }
    function handleKeydown(e) {
      if (e.key === "Escape" && props.isOpen) {
        close();
      }
    }
    watch(() => props.isOpen, (isOpen) => {
      if (isOpen) {
        document.addEventListener("keydown", handleKeydown);
      } else {
        document.removeEventListener("keydown", handleKeydown);
      }
    });
    return (_ctx, _cache) => {
      return openBlock(), createBlock(Teleport, { to: "body" }, [
        __props.isOpen ? (openBlock(), createElementBlock("div", {
          key: 0,
          class: "fixed inset-0 z-50 flex items-center justify-center p-4",
          onClick: withModifiers(close, ["self"])
        }, [
          createBaseVNode("div", {
            class: "absolute inset-0 bg-black/70 backdrop-blur-sm",
            onClick: close
          }),
          createBaseVNode("div", _hoisted_1$t, [
            createBaseVNode("div", _hoisted_2$s, [
              createBaseVNode("div", _hoisted_3$r, [
                _cache[2] || (_cache[2] = createBaseVNode("div", { class: "panel-icon bg-purple-500/20" }, [
                  createBaseVNode("svg", {
                    class: "w-4 h-4 text-purple-400",
                    fill: "none",
                    stroke: "currentColor",
                    viewBox: "0 0 24 24"
                  }, [
                    createBaseVNode("path", {
                      "stroke-linecap": "round",
                      "stroke-linejoin": "round",
                      "stroke-width": "2",
                      d: "M8 7h12m0 0l-4-4m4 4l-4 4m0 6H4m0 0l4 4m-4-4l4-4"
                    })
                  ])
                ], -1)),
                createBaseVNode("h3", _hoisted_4$p, toDisplayString(unref(t)("git.compareBranches")), 1)
              ]),
              createBaseVNode("button", {
                onClick: close,
                class: "p-2 text-gray-400 hover:text-white hover:bg-gray-700 rounded-lg transition-colors"
              }, [..._cache[3] || (_cache[3] = [
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
            createBaseVNode("div", _hoisted_5$l, [
              createBaseVNode("div", _hoisted_6$j, [
                createBaseVNode("label", _hoisted_7$j, toDisplayString(unref(t)("git.baseBranch")), 1),
                withDirectives(createBaseVNode("select", {
                  "onUpdate:modelValue": _cache[0] || (_cache[0] = ($event) => baseBranch.value = $event),
                  class: "input w-full text-sm",
                  onChange: loadDiff
                }, [
                  (openBlock(true), createElementBlock(Fragment, null, renderList(__props.branches, (branch) => {
                    return openBlock(), createElementBlock("option", {
                      key: branch,
                      value: branch
                    }, toDisplayString(branch), 9, _hoisted_8$h);
                  }), 128))
                ], 544), [
                  [vModelSelect, baseBranch.value]
                ])
              ]),
              _cache[5] || (_cache[5] = createBaseVNode("div", { class: "flex items-center text-gray-400 pt-4" }, [
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
                    d: "M14 5l7 7m0 0l-7 7m7-7H3"
                  })
                ])
              ], -1)),
              createBaseVNode("div", _hoisted_9$f, [
                createBaseVNode("label", _hoisted_10$f, toDisplayString(unref(t)("git.compareBranch")), 1),
                withDirectives(createBaseVNode("select", {
                  "onUpdate:modelValue": _cache[1] || (_cache[1] = ($event) => compareBranch.value = $event),
                  class: "input w-full text-sm",
                  onChange: loadDiff
                }, [
                  (openBlock(true), createElementBlock(Fragment, null, renderList(__props.branches, (branch) => {
                    return openBlock(), createElementBlock("option", {
                      key: branch,
                      value: branch
                    }, toDisplayString(branch), 9, _hoisted_11$e);
                  }), 128))
                ], 544), [
                  [vModelSelect, compareBranch.value]
                ])
              ]),
              createBaseVNode("button", {
                onClick: swapBranches,
                class: "p-2 text-gray-400 hover:text-white hover:bg-gray-700 rounded-lg transition-colors mt-4",
                title: unref(t)("git.swapBranches")
              }, [..._cache[4] || (_cache[4] = [
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
                    d: "M7 16V4m0 0L3 8m4-4l4 4m6 0v12m0 0l4-4m-4 4l-4-4"
                  })
                ], -1)
              ])], 8, _hoisted_12$c)
            ]),
            createBaseVNode("div", _hoisted_13$c, [
              isLoading.value ? (openBlock(), createElementBlock("div", _hoisted_14$b, [..._cache[6] || (_cache[6] = [
                createBaseVNode("svg", {
                  class: "animate-spin w-8 h-8 text-indigo-400",
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
                ], -1)
              ])])) : error.value ? (openBlock(), createElementBlock("div", _hoisted_15$9, [
                createBaseVNode("div", _hoisted_16$8, [
                  _cache[7] || (_cache[7] = createBaseVNode("svg", {
                    class: "w-12 h-12 mx-auto mb-2",
                    fill: "none",
                    stroke: "currentColor",
                    viewBox: "0 0 24 24"
                  }, [
                    createBaseVNode("path", {
                      "stroke-linecap": "round",
                      "stroke-linejoin": "round",
                      "stroke-width": "2",
                      d: "M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"
                    })
                  ], -1)),
                  createBaseVNode("p", null, toDisplayString(error.value), 1)
                ])
              ])) : diffFiles.value.length === 0 && !isLoading.value ? (openBlock(), createElementBlock("div", _hoisted_17$8, [
                createBaseVNode("div", _hoisted_18$7, [
                  _cache[8] || (_cache[8] = createBaseVNode("svg", {
                    class: "w-12 h-12 mx-auto mb-2",
                    fill: "none",
                    stroke: "currentColor",
                    viewBox: "0 0 24 24"
                  }, [
                    createBaseVNode("path", {
                      "stroke-linecap": "round",
                      "stroke-linejoin": "round",
                      "stroke-width": "2",
                      d: "M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"
                    })
                  ], -1)),
                  createBaseVNode("p", null, toDisplayString(unref(t)("git.noDifferences")), 1)
                ])
              ])) : (openBlock(), createElementBlock("div", _hoisted_19$7, [
                (openBlock(true), createElementBlock(Fragment, null, renderList(diffFiles.value, (file) => {
                  return openBlock(), createElementBlock("div", {
                    key: file.path,
                    class: "p-3 hover:bg-gray-800/50"
                  }, [
                    createBaseVNode("div", _hoisted_20$7, [
                      createBaseVNode("span", {
                        class: normalizeClass([
                          "px-2 py-0.5 text-xs font-medium rounded",
                          file.status === "added" ? "bg-emerald-500/20 text-emerald-400" : file.status === "deleted" ? "bg-red-500/20 text-red-400" : file.status === "modified" ? "bg-amber-500/20 text-amber-400" : "bg-gray-500/20 text-gray-400"
                        ])
                      }, toDisplayString(file.status === "added" ? "A" : file.status === "deleted" ? "D" : "M"), 3),
                      createBaseVNode("span", _hoisted_21$5, toDisplayString(file.path), 1),
                      file.additions ? (openBlock(), createElementBlock("span", _hoisted_22$5, "+" + toDisplayString(file.additions), 1)) : createCommentVNode("", true),
                      file.deletions ? (openBlock(), createElementBlock("span", _hoisted_23$5, "-" + toDisplayString(file.deletions), 1)) : createCommentVNode("", true)
                    ])
                  ]);
                }), 128))
              ]))
            ]),
            createBaseVNode("div", _hoisted_24$5, [
              createBaseVNode("div", _hoisted_25$4, [
                diffFiles.value.length > 0 ? (openBlock(), createElementBlock("span", _hoisted_26$4, [
                  createTextVNode(toDisplayString(diffFiles.value.length) + " " + toDisplayString(unref(t)("git.filesChanged")) + " ", 1),
                  totalAdditions.value ? (openBlock(), createElementBlock("span", _hoisted_27$3, "+" + toDisplayString(totalAdditions.value), 1)) : createCommentVNode("", true),
                  totalDeletions.value ? (openBlock(), createElementBlock("span", _hoisted_28$3, "-" + toDisplayString(totalDeletions.value), 1)) : createCommentVNode("", true)
                ])) : createCommentVNode("", true)
              ]),
              createBaseVNode("button", {
                onClick: close,
                class: "action-btn action-btn-primary"
              }, toDisplayString(unref(t)("git.close")), 1)
            ])
          ])
        ])) : createCommentVNode("", true)
      ]);
    };
  }
});
const _hoisted_1$s = { class: "relative w-full max-w-4xl max-h-[85vh] bg-gray-900 rounded-xl border border-gray-700 shadow-2xl flex flex-col overflow-hidden" };
const _hoisted_2$r = { class: "flex items-center justify-between p-4 border-b border-gray-700" };
const _hoisted_3$q = { class: "flex items-center gap-3 min-w-0" };
const _hoisted_4$o = { class: "text-lg" };
const _hoisted_5$k = { class: "min-w-0" };
const _hoisted_6$i = { class: "text-sm font-medium text-white truncate" };
const _hoisted_7$i = { class: "text-xs text-gray-400 truncate" };
const _hoisted_8$g = { class: "flex items-center gap-2" };
const _hoisted_9$e = {
  key: 0,
  class: "text-xs text-gray-400"
};
const _hoisted_10$e = ["title"];
const _hoisted_11$d = { class: "flex-1 overflow-auto" };
const _hoisted_12$b = {
  key: 0,
  class: "flex items-center justify-center h-64"
};
const _hoisted_13$b = {
  key: 1,
  class: "flex items-center justify-center h-64 text-red-400"
};
const _hoisted_14$a = { class: "text-center" };
const _hoisted_15$8 = {
  key: 2,
  class: "p-4 text-sm text-gray-300 font-mono whitespace-pre-wrap break-words"
};
const _hoisted_16$7 = { class: "px-4 py-2 border-t border-gray-700 flex items-center justify-between text-xs text-gray-400" };
const _hoisted_17$7 = { key: 0 };
const _sfc_main$t = /* @__PURE__ */ defineComponent({
  __name: "FilePreviewModal",
  props: {
    isOpen: { type: Boolean },
    filePath: {},
    content: {},
    isLoading: { type: Boolean },
    error: {}
  },
  emits: ["close"],
  setup(__props, { emit: __emit }) {
    const props = __props;
    const emit = __emit;
    const { t } = useI18n();
    const uiStore = useUIStore();
    const fileName = computed(() => props.filePath.split("/").pop() || "");
    const fileSize = computed(() => props.content?.length || 0);
    const lineCount = computed(() => props.content?.split("\n").length || 0);
    const language = computed(() => {
      const ext = fileName.value.split(".").pop()?.toLowerCase();
      const langMap = {
        ts: "TypeScript",
        tsx: "TypeScript",
        js: "JavaScript",
        jsx: "JavaScript",
        vue: "Vue",
        go: "Go",
        py: "Python",
        java: "Java",
        rs: "Rust",
        cpp: "C++",
        c: "C",
        h: "C/C++",
        cs: "C#",
        rb: "Ruby",
        php: "PHP",
        swift: "Swift",
        kt: "Kotlin",
        css: "CSS",
        scss: "SCSS",
        html: "HTML",
        json: "JSON",
        yaml: "YAML",
        yml: "YAML",
        md: "Markdown",
        xml: "XML"
      };
      return ext ? langMap[ext] : void 0;
    });
    function close() {
      emit("close");
    }
    function formatSize(bytes) {
      if (bytes < 1024) return `${bytes} B`;
      if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)} KB`;
      return `${(bytes / (1024 * 1024)).toFixed(1)} MB`;
    }
    async function copyContent() {
      if (!props.content) return;
      try {
        await navigator.clipboard.writeText(props.content);
        uiStore.addToast(t("toast.contextCopied"), "success");
      } catch {
        uiStore.addToast(t("toast.copyError"), "error");
      }
    }
    function handleKeydown(e) {
      if (e.key === "Escape" && props.isOpen) {
        close();
      }
    }
    watch(() => props.isOpen, (isOpen) => {
      if (isOpen) {
        document.addEventListener("keydown", handleKeydown);
      } else {
        document.removeEventListener("keydown", handleKeydown);
      }
    });
    return (_ctx, _cache) => {
      return openBlock(), createBlock(Teleport, { to: "body" }, [
        __props.isOpen ? (openBlock(), createElementBlock("div", {
          key: 0,
          class: "fixed inset-0 z-50 flex items-center justify-center p-4",
          onClick: withModifiers(close, ["self"])
        }, [
          createBaseVNode("div", {
            class: "absolute inset-0 bg-black/70 backdrop-blur-sm",
            onClick: close
          }),
          createBaseVNode("div", _hoisted_1$s, [
            createBaseVNode("div", _hoisted_2$r, [
              createBaseVNode("div", _hoisted_3$q, [
                createBaseVNode("span", _hoisted_4$o, toDisplayString(unref(getFileIcon)(fileName.value)), 1),
                createBaseVNode("div", _hoisted_5$k, [
                  createBaseVNode("h3", _hoisted_6$i, toDisplayString(fileName.value), 1),
                  createBaseVNode("p", _hoisted_7$i, toDisplayString(__props.filePath), 1)
                ])
              ]),
              createBaseVNode("div", _hoisted_8$g, [
                fileSize.value ? (openBlock(), createElementBlock("span", _hoisted_9$e, toDisplayString(formatSize(fileSize.value)), 1)) : createCommentVNode("", true),
                createBaseVNode("button", {
                  onClick: copyContent,
                  class: "p-2 text-gray-400 hover:text-white hover:bg-gray-700 rounded-lg transition-colors",
                  title: unref(t)("action.copy")
                }, [..._cache[0] || (_cache[0] = [
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
                      d: "M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z"
                    })
                  ], -1)
                ])], 8, _hoisted_10$e),
                createBaseVNode("button", {
                  onClick: close,
                  class: "p-2 text-gray-400 hover:text-white hover:bg-gray-700 rounded-lg transition-colors"
                }, [..._cache[1] || (_cache[1] = [
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
              ])
            ]),
            createBaseVNode("div", _hoisted_11$d, [
              __props.isLoading ? (openBlock(), createElementBlock("div", _hoisted_12$b, [..._cache[2] || (_cache[2] = [
                createBaseVNode("svg", {
                  class: "animate-spin w-8 h-8 text-indigo-400",
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
                ], -1)
              ])])) : __props.error ? (openBlock(), createElementBlock("div", _hoisted_13$b, [
                createBaseVNode("div", _hoisted_14$a, [
                  _cache[3] || (_cache[3] = createBaseVNode("svg", {
                    class: "w-12 h-12 mx-auto mb-2",
                    fill: "none",
                    stroke: "currentColor",
                    viewBox: "0 0 24 24"
                  }, [
                    createBaseVNode("path", {
                      "stroke-linecap": "round",
                      "stroke-linejoin": "round",
                      "stroke-width": "2",
                      d: "M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"
                    })
                  ], -1)),
                  createBaseVNode("p", null, toDisplayString(__props.error), 1)
                ])
              ])) : (openBlock(), createElementBlock("pre", _hoisted_15$8, [
                createBaseVNode("code", null, toDisplayString(__props.content), 1)
              ]))
            ]),
            createBaseVNode("div", _hoisted_16$7, [
              createBaseVNode("span", null, toDisplayString(lineCount.value) + " " + toDisplayString(unref(t)("context.lines")), 1),
              language.value ? (openBlock(), createElementBlock("span", _hoisted_17$7, toDisplayString(language.value), 1)) : createCommentVNode("", true)
            ])
          ])
        ])) : createCommentVNode("", true)
      ]);
    };
  }
});
const cache = /* @__PURE__ */ new Map();
const DEFAULT_TTL = 5 * 60 * 1e3;
function useGitCache() {
  const isLoading = ref(false);
  function getCacheKey(type, ...args) {
    return `git:${type}:${args.join(":")}`;
  }
  function get(key) {
    const entry = cache.get(key);
    if (!entry) return null;
    if (Date.now() - entry.timestamp > entry.ttl) {
      cache.delete(key);
      return null;
    }
    return entry.data;
  }
  function set(key, data, ttl = DEFAULT_TTL) {
    cache.set(key, {
      data,
      timestamp: Date.now(),
      ttl
    });
  }
  function invalidate(pattern) {
    if (!pattern) {
      cache.clear();
      return;
    }
    for (const key of cache.keys()) {
      if (key.includes(pattern)) {
        cache.delete(key);
      }
    }
  }
  async function cachedFetch(key, fetcher, ttl = DEFAULT_TTL) {
    const cached = get(key);
    if (cached !== null) {
      return cached;
    }
    isLoading.value = true;
    try {
      const data = await fetcher();
      set(key, data, ttl);
      return data;
    } finally {
      isLoading.value = false;
    }
  }
  async function cachedBranches(projectPath, fetcher) {
    const key = getCacheKey("branches", projectPath);
    return cachedFetch(key, fetcher, 2 * 60 * 1e3);
  }
  async function cachedCommits(projectPath, ref2, fetcher) {
    const key = getCacheKey("commits", projectPath, ref2);
    return cachedFetch(key, fetcher, 1 * 60 * 1e3);
  }
  async function cachedFiles(source, ref2, fetcher) {
    const key = getCacheKey("files", source, ref2);
    return cachedFetch(key, fetcher, 5 * 60 * 1e3);
  }
  async function cachedFileContent(source, filePath, ref2, fetcher) {
    const key = getCacheKey("content", source, filePath, ref2);
    return cachedFetch(key, fetcher, 10 * 60 * 1e3);
  }
  function getCacheStats() {
    let totalSize = 0;
    let entryCount = 0;
    let expiredCount = 0;
    const now = Date.now();
    for (const [, entry] of cache) {
      entryCount++;
      if (now - entry.timestamp > entry.ttl) {
        expiredCount++;
      }
      totalSize += JSON.stringify(entry.data).length;
    }
    return {
      entryCount,
      expiredCount,
      totalSizeKB: Math.round(totalSize / 1024)
    };
  }
  return {
    isLoading,
    get,
    set,
    invalidate,
    cachedFetch,
    cachedBranches,
    cachedCommits,
    cachedFiles,
    cachedFileContent,
    getCacheStats
  };
}
let sharedCache = null;
function getGitCache() {
  if (!sharedCache) {
    sharedCache = useGitCache();
  }
  return sharedCache;
}
const gitApi = {
  // Provider detection
  async detectProvider(url) {
    const [isGitHub, isGitLab] = await Promise.all([
      apiService.isGitHubURL(url),
      apiService.isGitLabURL(url)
    ]);
    if (isGitHub) return "github";
    if (isGitLab) return "gitlab";
    return "unknown";
  },
  // Local repository
  async isGitRepository(projectPath) {
    return apiService.isGitRepository(projectPath);
  },
  async getCurrentBranch(projectPath) {
    return apiService.getCurrentBranch(projectPath);
  },
  async getBranches(projectPath) {
    const cache2 = getGitCache();
    return cache2.cachedBranches(projectPath, async () => {
      const result = await apiService.getBranches(projectPath);
      return JSON.parse(result);
    });
  },
  async getCommits(projectPath, limit = 50) {
    const cache2 = getGitCache();
    return cache2.cachedCommits(projectPath, "HEAD", async () => {
      return apiService.getCommitHistory(projectPath, limit);
    });
  },
  async listFilesAtRef(projectPath, ref2) {
    const cache2 = getGitCache();
    return cache2.cachedFiles(projectPath, ref2, async () => {
      const result = await apiService.listFilesAtRef(projectPath, ref2);
      return Array.isArray(result) ? result : [];
    });
  },
  async getFileAtRef(projectPath, filePath, ref2) {
    const cache2 = getGitCache();
    return cache2.cachedFileContent(projectPath, filePath, ref2, async () => {
      return apiService.getFileAtRef(projectPath, filePath, ref2);
    });
  },
  async buildContextAtRef(projectPath, files2, ref2) {
    return apiService.buildContextAtRef(projectPath, files2, ref2);
  },
  // GitHub API
  async gitHubGetBranches(repoUrl) {
    const branches = await apiService.gitHubGetBranches(repoUrl);
    return branches.map((b) => ({
      name: b.name,
      commit: { sha: b.commit.sha }
    }));
  },
  async gitHubGetDefaultBranch(repoUrl) {
    return apiService.gitHubGetDefaultBranch(repoUrl);
  },
  async gitHubListFiles(repoUrl, ref2) {
    const cache2 = getGitCache();
    return cache2.cachedFiles(repoUrl, ref2, async () => {
      return apiService.gitHubListFiles(repoUrl, ref2);
    });
  },
  async gitHubGetFileContent(repoUrl, filePath, ref2) {
    const cache2 = getGitCache();
    return cache2.cachedFileContent(repoUrl, filePath, ref2, async () => {
      return apiService.gitHubGetFileContent(repoUrl, filePath, ref2);
    });
  },
  async gitHubBuildContext(repoUrl, files2, ref2) {
    return apiService.gitHubBuildContext(repoUrl, files2, ref2);
  },
  // GitLab API
  async gitLabGetBranches(repoUrl) {
    const branches = await apiService.gitLabGetBranches(repoUrl);
    return branches.map((b) => ({
      name: b.name,
      commit: { sha: b.commit.id },
      isDefault: b.default
    }));
  },
  async gitLabGetDefaultBranch(repoUrl) {
    return apiService.gitLabGetDefaultBranch(repoUrl);
  },
  async gitLabListFiles(repoUrl, ref2) {
    const cache2 = getGitCache();
    return cache2.cachedFiles(repoUrl, ref2, async () => {
      return apiService.gitLabListFiles(repoUrl, ref2);
    });
  },
  async gitLabGetFileContent(repoUrl, filePath, ref2) {
    const cache2 = getGitCache();
    return cache2.cachedFileContent(repoUrl, filePath, ref2, async () => {
      return apiService.gitLabGetFileContent(repoUrl, filePath, ref2);
    });
  },
  async gitLabBuildContext(repoUrl, files2, ref2) {
    return apiService.gitLabBuildContext(repoUrl, files2, ref2);
  },
  // Universal methods (auto-detect provider)
  async getRemoteBranches(repoUrl, provider) {
    switch (provider) {
      case "github":
        return this.gitHubGetBranches(repoUrl);
      case "gitlab":
        return this.gitLabGetBranches(repoUrl);
      default:
        throw new Error(`Unsupported provider: ${provider}`);
    }
  },
  async getRemoteDefaultBranch(repoUrl, provider) {
    switch (provider) {
      case "github":
        return this.gitHubGetDefaultBranch(repoUrl);
      case "gitlab":
        return this.gitLabGetDefaultBranch(repoUrl);
      default:
        throw new Error(`Unsupported provider: ${provider}`);
    }
  },
  async listRemoteFiles(repoUrl, ref2, provider) {
    switch (provider) {
      case "github":
        return this.gitHubListFiles(repoUrl, ref2);
      case "gitlab":
        return this.gitLabListFiles(repoUrl, ref2);
      default:
        throw new Error(`Unsupported provider: ${provider}`);
    }
  },
  async getRemoteFileContent(repoUrl, filePath, ref2, provider) {
    switch (provider) {
      case "github":
        return this.gitHubGetFileContent(repoUrl, filePath, ref2);
      case "gitlab":
        return this.gitLabGetFileContent(repoUrl, filePath, ref2);
      default:
        throw new Error(`Unsupported provider: ${provider}`);
    }
  },
  async buildRemoteContext(repoUrl, files2, ref2, provider) {
    switch (provider) {
      case "github":
        return this.gitHubBuildContext(repoUrl, files2, ref2);
      case "gitlab":
        return this.gitLabBuildContext(repoUrl, files2, ref2);
      default:
        throw new Error(`Unsupported provider: ${provider}`);
    }
  },
  // Cache management
  invalidateCache(pattern) {
    getGitCache().invalidate(pattern);
  },
  getCacheStats() {
    return getGitCache().getCacheStats();
  }
};
const RECENT_REPOS_KEY = "git-recent-repos";
const MAX_RECENT_REPOS = 10;
const SELECTION_KEY = "git-file-selection";
const fileTypeFilters = [
  { id: "code", label: "Code", icon: "💻", extensions: [".ts", ".js", ".tsx", ".jsx", ".vue", ".go", ".py", ".java", ".rs", ".cpp", ".c", ".h", ".cs", ".rb", ".php", ".swift", ".kt"] },
  { id: "styles", label: "Styles", icon: "🎨", extensions: [".css", ".scss", ".sass", ".less", ".styl"] },
  { id: "config", label: "Config", icon: "⚙️", extensions: [".json", ".yaml", ".yml", ".toml", ".xml", ".ini", ".env", ".config"] },
  { id: "docs", label: "Docs", icon: "📄", extensions: [".md", ".txt", ".rst", ".adoc"] }
];
function useGitSource() {
  const projectStore = useProjectStore();
  const uiStore = useUIStore();
  const isGitRepo = ref(false);
  const currentBranch = ref("");
  const branches = ref([]);
  const commits = ref([]);
  const commitsLoaded = ref(false);
  const selectedRef = ref(null);
  const filesAtRef = ref([]);
  const selectedFiles = ref(/* @__PURE__ */ new Set());
  const remoteUrl = ref("");
  const isCloning = ref(false);
  const clonedPath = ref(null);
  const isGitHubRepo = ref(false);
  const isGitLabRepo = ref(false);
  const isLoadingRemote = ref(false);
  const remoteRepoLoaded = ref(false);
  const remoteBranches = ref([]);
  const remoteSelectedBranch = ref("");
  const remoteFiles = ref([]);
  const remoteSelectedFiles = ref(/* @__PURE__ */ new Set());
  const isLoading = ref(false);
  const loadingMessage = ref("");
  const isBuilding = ref(false);
  const localSearchQuery = ref("");
  const remoteSearchQuery = ref("");
  const localActiveFilters = ref(/* @__PURE__ */ new Set());
  const remoteActiveFilters = ref(/* @__PURE__ */ new Set());
  const recentRepos = ref([]);
  const projectPath = computed(() => projectStore.currentPath || "");
  function getFilteredFiles(files2, searchQuery, activeFilters) {
    let result = files2;
    if (searchQuery) {
      const query = searchQuery.toLowerCase();
      result = result.filter((f) => f.toLowerCase().includes(query));
    }
    if (activeFilters.size > 0) {
      const activeExtensions = /* @__PURE__ */ new Set();
      activeFilters.forEach((filterId) => {
        const filter = fileTypeFilters.find((f) => f.id === filterId);
        if (filter) filter.extensions.forEach((ext) => activeExtensions.add(ext));
      });
      result = result.filter((f) => {
        const ext = "." + f.split(".").pop()?.toLowerCase();
        return activeExtensions.has(ext);
      });
    }
    return result;
  }
  const filteredLocalFiles = computed(
    () => getFilteredFiles(filesAtRef.value || [], localSearchQuery.value, localActiveFilters.value)
  );
  const filteredRemoteFiles = computed(
    () => getFilteredFiles(remoteFiles.value || [], remoteSearchQuery.value, remoteActiveFilters.value)
  );
  function getFilterCount(files2, filterId) {
    const filter = fileTypeFilters.find((f) => f.id === filterId);
    if (!filter || !files2) return 0;
    return files2.filter((f) => {
      const ext = "." + f.split(".").pop()?.toLowerCase();
      return filter.extensions.includes(ext);
    }).length;
  }
  function toggleFilter(activeFilters, filterId) {
    const newFilters = new Set(activeFilters);
    if (newFilters.has(filterId)) {
      newFilters.delete(filterId);
    } else {
      newFilters.add(filterId);
    }
    return newFilters;
  }
  function loadRecentReposFromStorage() {
    try {
      const stored = localStorage.getItem(RECENT_REPOS_KEY);
      if (stored) {
        recentRepos.value = JSON.parse(stored);
      }
    } catch {
      recentRepos.value = [];
    }
  }
  function saveRecentRepo(url, name, isGitHub) {
    const existing = recentRepos.value.findIndex((r) => r.url === url);
    if (existing >= 0) {
      recentRepos.value.splice(existing, 1);
    }
    recentRepos.value.unshift({ url, name, isGitHub, lastUsed: Date.now() });
    if (recentRepos.value.length > MAX_RECENT_REPOS) {
      recentRepos.value = recentRepos.value.slice(0, MAX_RECENT_REPOS);
    }
    localStorage.setItem(RECENT_REPOS_KEY, JSON.stringify(recentRepos.value));
  }
  function clearRecentRepos() {
    recentRepos.value = [];
    localStorage.removeItem(RECENT_REPOS_KEY);
  }
  function saveSelectionToStorage(key, files2) {
    try {
      const data = { files: Array.from(files2), timestamp: Date.now() };
      localStorage.setItem(`${SELECTION_KEY}-${key}`, JSON.stringify(data));
    } catch {
    }
  }
  function loadSelectionFromStorage(key) {
    try {
      const stored = localStorage.getItem(`${SELECTION_KEY}-${key}`);
      if (stored) {
        const data = JSON.parse(stored);
        if (Date.now() - data.timestamp < 24 * 60 * 60 * 1e3) {
          return new Set(data.files);
        }
      }
    } catch {
    }
    return /* @__PURE__ */ new Set();
  }
  async function loadGitInfo() {
    if (!projectPath.value) return;
    isLoading.value = true;
    loadingMessage.value = "Loading git info...";
    try {
      const [isRepo, branch, branchList] = await Promise.all([
        gitApi.isGitRepository(projectPath.value),
        gitApi.getCurrentBranch(projectPath.value).catch(() => ""),
        gitApi.getBranches(projectPath.value).catch(() => [])
      ]);
      isGitRepo.value = isRepo;
      currentBranch.value = branch;
      branches.value = branchList;
    } catch {
      isGitRepo.value = false;
    } finally {
      isLoading.value = false;
    }
  }
  async function loadCommits() {
    if (!projectPath.value || commitsLoaded.value) return;
    isLoading.value = true;
    loadingMessage.value = "Loading commits...";
    try {
      const gitCommits = await gitApi.getCommits(projectPath.value, 50);
      commits.value = gitCommits.map((c) => ({
        hash: c.hash,
        subject: c.message || "",
        author: c.author,
        date: c.date
      }));
      commitsLoaded.value = true;
    } catch {
      commits.value = [];
    } finally {
      isLoading.value = false;
    }
  }
  async function loadFilesAtRef(ref2) {
    if (!projectPath.value) return;
    isLoading.value = true;
    loadingMessage.value = `Loading files at ${ref2.slice(0, 7)}...`;
    try {
      filesAtRef.value = await gitApi.listFilesAtRef(projectPath.value, ref2);
      const savedSelection = loadSelectionFromStorage(`local-${projectPath.value}-${ref2}`);
      if (savedSelection.size > 0) {
        selectedFiles.value = savedSelection;
      }
    } catch {
      filesAtRef.value = [];
      uiStore.addToast("Failed to load files", "error");
    } finally {
      isLoading.value = false;
    }
  }
  function selectRef(ref2) {
    selectedRef.value = ref2;
    loadFilesAtRef(ref2);
  }
  function clearSelectedRef() {
    selectedRef.value = null;
    filesAtRef.value = [];
    selectedFiles.value = /* @__PURE__ */ new Set();
  }
  function toggleFileSelection(path) {
    const newSelection = new Set(selectedFiles.value);
    if (newSelection.has(path)) {
      newSelection.delete(path);
    } else {
      newSelection.add(path);
    }
    selectedFiles.value = newSelection;
    if (selectedRef.value && projectPath.value) {
      saveSelectionToStorage(`local-${projectPath.value}-${selectedRef.value}`, newSelection);
    }
  }
  function selectAllFiles() {
    selectedFiles.value = new Set(filteredLocalFiles.value);
  }
  function clearFileSelection() {
    selectedFiles.value = /* @__PURE__ */ new Set();
  }
  function selectFolderFiles(folderPath) {
    const newSelection = new Set(selectedFiles.value);
    filteredLocalFiles.value.filter((f) => f.startsWith(folderPath + "/")).forEach((f) => newSelection.add(f));
    selectedFiles.value = newSelection;
  }
  function checkRemoteType(url) {
    isGitHubRepo.value = url.includes("github.com");
    isGitLabRepo.value = url.includes("gitlab.com");
  }
  async function loadRemoteRepo() {
    if (!remoteUrl.value) return;
    checkRemoteType(remoteUrl.value);
    isLoadingRemote.value = true;
    try {
      const provider = isGitHubRepo.value ? "github" : isGitLabRepo.value ? "gitlab" : "unknown";
      if (provider === "unknown") {
        uiStore.addToast("Unsupported repository provider", "error");
        return;
      }
      const [branchList, defaultBranch] = await Promise.all([
        gitApi.getRemoteBranches(remoteUrl.value, provider),
        gitApi.getRemoteDefaultBranch(remoteUrl.value, provider)
      ]);
      remoteBranches.value = branchList.map((b) => ({ name: b.name, commit: { sha: b.commit.sha } }));
      remoteSelectedBranch.value = defaultBranch;
      remoteRepoLoaded.value = true;
      await loadRemoteFiles();
      const repoName = remoteUrl.value.split("/").slice(-2).join("/");
      saveRecentRepo(remoteUrl.value, repoName, isGitHubRepo.value);
    } catch {
      uiStore.addToast("Failed to load repository", "error");
    } finally {
      isLoadingRemote.value = false;
    }
  }
  async function loadRemoteFiles() {
    if (!remoteSelectedBranch.value) return;
    isLoading.value = true;
    loadingMessage.value = "Loading files...";
    try {
      const provider = isGitHubRepo.value ? "github" : isGitLabRepo.value ? "gitlab" : "unknown";
      if (provider !== "unknown") {
        remoteFiles.value = await gitApi.listRemoteFiles(remoteUrl.value, remoteSelectedBranch.value, provider);
      }
    } catch {
      remoteFiles.value = [];
    } finally {
      isLoading.value = false;
    }
  }
  function toggleRemoteFileSelection(path) {
    const newSelection = new Set(remoteSelectedFiles.value);
    if (newSelection.has(path)) {
      newSelection.delete(path);
    } else {
      newSelection.add(path);
    }
    remoteSelectedFiles.value = newSelection;
  }
  function selectAllRemoteFiles() {
    remoteSelectedFiles.value = new Set(filteredRemoteFiles.value);
  }
  function clearRemoteFileSelection() {
    remoteSelectedFiles.value = /* @__PURE__ */ new Set();
  }
  function selectRemoteFolderFiles(folderPath) {
    const newSelection = new Set(remoteSelectedFiles.value);
    filteredRemoteFiles.value.filter((f) => f.startsWith(folderPath + "/")).forEach((f) => newSelection.add(f));
    remoteSelectedFiles.value = newSelection;
  }
  watch(remoteUrl, (url) => {
    if (url) {
      checkRemoteType(url);
    }
  });
  return {
    // Local state
    isGitRepo,
    currentBranch,
    branches,
    commits,
    commitsLoaded,
    selectedRef,
    filesAtRef,
    selectedFiles,
    filteredLocalFiles,
    localSearchQuery,
    localActiveFilters,
    // Remote state
    remoteUrl,
    isCloning,
    clonedPath,
    isGitHubRepo,
    isGitLabRepo,
    isLoadingRemote,
    remoteRepoLoaded,
    remoteBranches,
    remoteSelectedBranch,
    remoteFiles,
    remoteSelectedFiles,
    filteredRemoteFiles,
    remoteSearchQuery,
    remoteActiveFilters,
    // Loading state
    isLoading,
    loadingMessage,
    isBuilding,
    // Recent repos
    recentRepos,
    // Computed
    projectPath,
    // Methods
    loadGitInfo,
    loadCommits,
    selectRef,
    clearSelectedRef,
    toggleFileSelection,
    selectAllFiles,
    clearFileSelection,
    selectFolderFiles,
    loadRemoteRepo,
    loadRemoteFiles,
    toggleRemoteFileSelection,
    selectAllRemoteFiles,
    clearRemoteFileSelection,
    selectRemoteFolderFiles,
    loadRecentReposFromStorage,
    saveRecentRepo,
    clearRecentRepos,
    getFilterCount,
    toggleFilter,
    fileTypeFilters
  };
}
const _hoisted_1$r = { class: "select-none p-2" };
const _hoisted_2$q = ["onClick", "onDblclick"];
const _hoisted_3$p = ["width"];
const _hoisted_4$n = ["x1", "x2"];
const _hoisted_5$j = ["d"];
const _hoisted_6$h = ["onClick"];
const _hoisted_7$h = {
  key: 2,
  class: "w-5"
};
const _hoisted_8$f = ["onClick"];
const _hoisted_9$d = {
  key: 0,
  class: "tree-cb-icon",
  fill: "currentColor",
  viewBox: "0 0 20 20"
};
const _hoisted_10$d = {
  key: 1,
  class: "w-2 h-0.5 bg-white rounded-full"
};
const _hoisted_11$c = { class: "tree-icon" };
const _hoisted_12$a = {
  key: 0,
  class: "tree-folder-icon",
  fill: "currentColor",
  viewBox: "0 0 20 20"
};
const _hoisted_13$a = {
  key: 1,
  class: "tree-file-icon"
};
const _hoisted_14$9 = { class: "tree-name" };
const _hoisted_15$7 = ["onClick"];
const _sfc_main$s = /* @__PURE__ */ defineComponent({
  __name: "SimpleFileTree",
  props: {
    files: {},
    selectedPaths: {}
  },
  emits: ["toggle-select", "select-folder", "preview-file"],
  setup(__props, { emit: __emit }) {
    const props = __props;
    const emit = __emit;
    const expandedPaths = ref(/* @__PURE__ */ new Set());
    const animatingNodes = ref(/* @__PURE__ */ new Set());
    const rippleNode = ref(null);
    const countAnimating = ref(null);
    const initializedFiles = ref(null);
    watch(() => props.files, (newFiles) => {
      if (!newFiles?.length) return;
      const signature = newFiles.slice(0, 5).join("|");
      if (signature === initializedFiles.value) return;
      initializedFiles.value = signature;
      const rootDirs = /* @__PURE__ */ new Set();
      for (const file of newFiles) {
        const firstPart = file.split("/")[0];
        if (firstPart && file.includes("/")) {
          rootDirs.add(firstPart);
        }
      }
      rootDirs.forEach((dir) => expandedPaths.value.add(dir));
      expandedPaths.value = new Set(expandedPaths.value);
    }, { immediate: true });
    function isSelected(path) {
      return props.selectedPaths?.has(path) ?? false;
    }
    function isExpanded(path) {
      return expandedPaths.value.has(path);
    }
    function getFilesInDir(dirPath) {
      if (!props.files || !Array.isArray(props.files)) return [];
      return props.files.filter((f) => f.startsWith(dirPath + "/"));
    }
    function hasSomeSelected(dirPath) {
      if (!props.selectedPaths) return false;
      const filesInDir = getFilesInDir(dirPath);
      return filesInDir.some((f) => props.selectedPaths.has(f));
    }
    function hasAllSelected(dirPath) {
      if (!props.selectedPaths) return false;
      const filesInDir = getFilesInDir(dirPath);
      return filesInDir.length > 0 && filesInDir.every((f) => props.selectedPaths.has(f));
    }
    const visibleNodes = computed(() => {
      if (!props.files?.length) return [];
      const nodes = [];
      const addedDirs = /* @__PURE__ */ new Set();
      const sortedFiles = [...props.files].sort();
      const dirChildren = /* @__PURE__ */ new Map();
      for (const file of sortedFiles) {
        const parts = file.split("/");
        let currentPath = "";
        for (let i = 0; i < parts.length; i++) {
          const parentPath = currentPath;
          currentPath = currentPath ? `${currentPath}/${parts[i]}` : parts[i];
          if (!dirChildren.has(parentPath)) dirChildren.set(parentPath, []);
          const children = dirChildren.get(parentPath);
          if (!children.includes(currentPath)) children.push(currentPath);
        }
      }
      for (const file of sortedFiles) {
        const parts = file.split("/");
        let currentPath = "";
        let isVisible = true;
        for (let i = 0; i < parts.length - 1; i++) {
          const parentPath = currentPath;
          currentPath = currentPath ? `${currentPath}/${parts[i]}` : parts[i];
          if (parentPath && !expandedPaths.value.has(parentPath)) {
            isVisible = false;
            break;
          }
          if (!addedDirs.has(currentPath)) {
            addedDirs.add(currentPath);
            const siblings2 = dirChildren.get(parentPath) || [];
            nodes.push({
              path: currentPath,
              name: parts[i],
              isDir: true,
              depth: i,
              isLast: siblings2.indexOf(currentPath) === siblings2.length - 1,
              parentPath
            });
          }
        }
        if (isVisible) {
          const parentPath = parts.slice(0, -1).join("/");
          if (!parentPath || expandedPaths.value.has(parentPath)) {
            const siblings2 = dirChildren.get(parentPath) || [];
            nodes.push({
              path: file,
              name: parts[parts.length - 1],
              isDir: false,
              depth: parts.length - 1,
              isLast: siblings2.indexOf(file) === siblings2.length - 1,
              parentPath
            });
          }
        }
      }
      return nodes;
    });
    function shouldDrawVerticalLine(node, d) {
      if (d === node.depth) return !node.isLast;
      const parts = node.path.split("/");
      if (d > parts.length) return false;
      const ancestorPath = parts.slice(0, d).join("/");
      const ancestorNode = visibleNodes.value.find((n) => n.path === ancestorPath);
      if (!ancestorNode) return false;
      return !ancestorNode.isLast;
    }
    function getConnectorPath(node) {
      const x = (node.depth - 1) * 20 + 10;
      const midY = 16;
      if (node.isLast) {
        return `M ${x} 0 L ${x} ${midY} L ${x + 10} ${midY}`;
      }
      return `M ${x} 0 L ${x} 32 M ${x} ${midY} L ${x + 10} ${midY}`;
    }
    function handleClick(node) {
      if (node.isDir) {
        if (expandedPaths.value.has(node.path)) {
          expandedPaths.value.delete(node.path);
        } else {
          expandedPaths.value.add(node.path);
          setTimeout(() => {
            const children = visibleNodes.value.filter((n) => n.parentPath === node.path);
            children.forEach((child) => animatingNodes.value.add(child.path));
            setTimeout(() => children.forEach((child) => animatingNodes.value.delete(child.path)), 300);
          }, 10);
        }
        expandedPaths.value = new Set(expandedPaths.value);
      }
    }
    function handleDoubleClick(node) {
      if (!node.isDir) {
        emit("preview-file", node.path);
      }
    }
    function toggleSelect(node) {
      rippleNode.value = node.path;
      setTimeout(() => {
        rippleNode.value = null;
      }, 400);
      if (node.isDir) {
        emit("select-folder", getFilesInDir(node.path));
        countAnimating.value = node.path;
        setTimeout(() => {
          countAnimating.value = null;
        }, 300);
      } else {
        emit("toggle-select", node.path);
      }
    }
    return (_ctx, _cache) => {
      return openBlock(), createElementBlock("div", _hoisted_1$r, [
        (openBlock(true), createElementBlock(Fragment, null, renderList(visibleNodes.value, (node) => {
          return openBlock(), createElementBlock("div", {
            key: node.path,
            class: normalizeClass([
              "tree-row group",
              isSelected(node.path) ? "tree-row-selected" : "",
              animatingNodes.value.has(node.path) ? "tree-stagger" : ""
            ]),
            style: normalizeStyle({ paddingLeft: `${node.depth * 20 + 12}px` }),
            onClick: ($event) => handleClick(node),
            onDblclick: ($event) => handleDoubleClick(node),
            tabindex: "0"
          }, [
            node.depth > 0 ? (openBlock(), createElementBlock("svg", {
              key: 0,
              class: "tree-guides",
              width: node.depth * 20,
              height: "32",
              style: { "left": "12px" }
            }, [
              (openBlock(true), createElementBlock(Fragment, null, renderList(node.depth, (d) => {
                return openBlock(), createElementBlock(Fragment, { key: d }, [
                  shouldDrawVerticalLine(node, d) ? (openBlock(), createElementBlock("line", {
                    key: 0,
                    x1: (d - 1) * 20 + 10,
                    y1: "0",
                    x2: (d - 1) * 20 + 10,
                    y2: "32",
                    class: normalizeClass(["tree-guide-line", `tree-guide-${Math.min(d, 5)}`])
                  }, null, 10, _hoisted_4$n)) : createCommentVNode("", true)
                ], 64);
              }), 128)),
              createBaseVNode("path", {
                d: getConnectorPath(node),
                class: normalizeClass(["tree-guide-line", `tree-guide-${Math.min(node.depth, 5)}`])
              }, null, 10, _hoisted_5$j)
            ], 8, _hoisted_3$p)) : createCommentVNode("", true),
            node.isDir ? (openBlock(), createElementBlock("div", {
              key: 1,
              class: normalizeClass(["tree-expand", isExpanded(node.path) ? "tree-expand-open" : ""]),
              onClick: withModifiers(($event) => handleClick(node), ["stop"])
            }, [..._cache[0] || (_cache[0] = [
              createBaseVNode("svg", {
                class: "tree-expand-icon",
                fill: "currentColor",
                viewBox: "0 0 20 20"
              }, [
                createBaseVNode("path", {
                  "fill-rule": "evenodd",
                  d: "M7.293 14.707a1 1 0 010-1.414L10.586 10 7.293 6.707a1 1 0 011.414-1.414l4 4a1 1 0 010 1.414l-4 4a1 1 0 01-1.414 0z",
                  "clip-rule": "evenodd"
                })
              ], -1)
            ])], 10, _hoisted_6$h)) : (openBlock(), createElementBlock("div", _hoisted_7$h)),
            createBaseVNode("div", {
              class: normalizeClass([
                "tree-cb",
                isSelected(node.path) ? "tree-cb-checked" : "",
                node.isDir && hasSomeSelected(node.path) && !hasAllSelected(node.path) ? "tree-cb-partial" : "",
                rippleNode.value === node.path ? "tree-cb-ripple" : ""
              ]),
              onClick: withModifiers(($event) => toggleSelect(node), ["stop"])
            }, [
              isSelected(node.path) || node.isDir && hasAllSelected(node.path) ? (openBlock(), createElementBlock("svg", _hoisted_9$d, [..._cache[1] || (_cache[1] = [
                createBaseVNode("path", {
                  "fill-rule": "evenodd",
                  d: "M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z",
                  "clip-rule": "evenodd"
                }, null, -1)
              ])])) : node.isDir && hasSomeSelected(node.path) ? (openBlock(), createElementBlock("div", _hoisted_10$d)) : createCommentVNode("", true)
            ], 10, _hoisted_8$f),
            createBaseVNode("div", _hoisted_11$c, [
              node.isDir ? (openBlock(), createElementBlock("svg", _hoisted_12$a, [..._cache[2] || (_cache[2] = [
                createBaseVNode("path", { d: "M2 6a2 2 0 012-2h5l2 2h5a2 2 0 012 2v6a2 2 0 01-2 2H4a2 2 0 01-2-2V6z" }, null, -1)
              ])])) : (openBlock(), createElementBlock("span", _hoisted_13$a, toDisplayString(unref(getFileIcon)(node.name)), 1))
            ]),
            createBaseVNode("span", _hoisted_14$9, toDisplayString(node.name), 1),
            !node.isDir ? (openBlock(), createElementBlock("button", {
              key: 3,
              onClick: withModifiers(($event) => emit("preview-file", node.path), ["stop"]),
              class: "tree-preview",
              title: "Preview"
            }, [..._cache[3] || (_cache[3] = [
              createBaseVNode("svg", {
                class: "w-3.5 h-3.5",
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
              ], -1)
            ])], 8, _hoisted_15$7)) : createCommentVNode("", true),
            node.isDir ? (openBlock(), createElementBlock("span", {
              key: 4,
              class: normalizeClass(["tree-count", countAnimating.value === node.path ? "tree-count-bounce" : ""])
            }, toDisplayString(getFilesInDir(node.path).length), 3)) : createCommentVNode("", true)
          ], 46, _hoisted_2$q);
        }), 128))
      ]);
    };
  }
});
const _hoisted_1$q = { class: "space-y-2" };
const _hoisted_2$p = { class: "flex items-center justify-between" };
const _hoisted_3$o = { class: "text-sm text-gray-300" };
const _hoisted_4$m = { class: "text-xs text-gray-400" };
const _hoisted_5$i = { class: "space-y-2" };
const _hoisted_6$g = { class: "relative" };
const _hoisted_7$g = ["placeholder"];
const _hoisted_8$e = { class: "flex flex-wrap gap-1.5" };
const _hoisted_9$c = ["onClick"];
const _hoisted_10$c = { class: "git-filter-icon" };
const _hoisted_11$b = { class: "git-filter-count" };
const _hoisted_12$9 = { class: "file-list-container" };
const _hoisted_13$9 = { class: "flex gap-2" };
const _sfc_main$r = /* @__PURE__ */ defineComponent({
  __name: "GitFileList",
  props: {
    title: {},
    files: {},
    selectedFiles: {}
  },
  emits: ["toggle-select", "select-folder", "select-all", "clear-selection", "preview-file"],
  setup(__props) {
    const logger2 = useLogger("GitFileList");
    const props = __props;
    const { t } = useI18n();
    const searchQuery = ref("");
    const activeFilters = ref(/* @__PURE__ */ new Set());
    const fileTypeFilters2 = [
      { id: "code", label: "Code", icon: "💻", extensions: [".ts", ".js", ".tsx", ".jsx", ".vue", ".go", ".py", ".java", ".rs", ".cpp", ".c", ".h", ".cs", ".rb", ".php", ".swift", ".kt"] },
      { id: "styles", label: "Styles", icon: "🎨", extensions: [".css", ".scss", ".sass", ".less", ".styl"] },
      { id: "config", label: "Config", icon: "⚙️", extensions: [".json", ".yaml", ".yml", ".toml", ".xml", ".ini", ".env", ".config"] },
      { id: "docs", label: "Docs", icon: "📄", extensions: [".md", ".txt", ".rst", ".adoc"] }
    ];
    const filteredFiles = computed(() => {
      let result = props.files;
      if (searchQuery.value) {
        const query = searchQuery.value.toLowerCase();
        result = result.filter((f) => f.toLowerCase().includes(query));
      }
      if (activeFilters.value.size > 0) {
        const activeExtensions = /* @__PURE__ */ new Set();
        activeFilters.value.forEach((filterId) => {
          const filter = fileTypeFilters2.find((f) => f.id === filterId);
          if (filter) filter.extensions.forEach((ext) => activeExtensions.add(ext));
        });
        result = result.filter((f) => {
          const ext = "." + f.split(".").pop()?.toLowerCase();
          return activeExtensions.has(ext);
        });
      }
      logger2.debug("Filtered files", { total: props.files.length, filtered: result.length });
      return result;
    });
    function getFilterCount(filterId) {
      const filter = fileTypeFilters2.find((f) => f.id === filterId);
      if (!filter) return 0;
      return props.files.filter((f) => {
        const ext = "." + f.split(".").pop()?.toLowerCase();
        return filter.extensions.includes(ext);
      }).length;
    }
    function toggleFilter(filterId) {
      const newFilters = new Set(activeFilters.value);
      if (newFilters.has(filterId)) {
        newFilters.delete(filterId);
      } else {
        newFilters.add(filterId);
      }
      activeFilters.value = newFilters;
    }
    return (_ctx, _cache) => {
      return openBlock(), createElementBlock("div", _hoisted_1$q, [
        createBaseVNode("div", _hoisted_2$p, [
          createBaseVNode("span", _hoisted_3$o, toDisplayString(__props.title), 1),
          createBaseVNode("span", _hoisted_4$m, toDisplayString(filteredFiles.value.length) + "/" + toDisplayString(__props.files.length), 1)
        ]),
        createBaseVNode("div", _hoisted_5$i, [
          createBaseVNode("div", _hoisted_6$g, [
            _cache[8] || (_cache[8] = createBaseVNode("svg", {
              class: "absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-400",
              fill: "none",
              stroke: "currentColor",
              viewBox: "0 0 24 24"
            }, [
              createBaseVNode("path", {
                "stroke-linecap": "round",
                "stroke-linejoin": "round",
                "stroke-width": "2",
                d: "M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"
              })
            ], -1)),
            withDirectives(createBaseVNode("input", {
              "onUpdate:modelValue": _cache[0] || (_cache[0] = ($event) => searchQuery.value = $event),
              type: "text",
              placeholder: unref(t)("git.searchFiles"),
              class: "input pl-10 py-1.5 text-sm w-full"
            }, null, 8, _hoisted_7$g), [
              [vModelText, searchQuery.value]
            ]),
            searchQuery.value ? (openBlock(), createElementBlock("button", {
              key: 0,
              onClick: _cache[1] || (_cache[1] = ($event) => searchQuery.value = ""),
              class: "absolute right-3 top-1/2 -translate-y-1/2 text-gray-400 hover:text-white"
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
                  d: "M6 18L18 6M6 6l12 12"
                })
              ], -1)
            ])])) : createCommentVNode("", true)
          ]),
          createBaseVNode("div", _hoisted_8$e, [
            (openBlock(), createElementBlock(Fragment, null, renderList(fileTypeFilters2, (filter) => {
              return createBaseVNode("button", {
                key: filter.id,
                onClick: ($event) => toggleFilter(filter.id),
                class: normalizeClass(["git-filter-chip", activeFilters.value.has(filter.id) ? "git-filter-chip-active" : ""])
              }, [
                createBaseVNode("span", _hoisted_10$c, toDisplayString(filter.icon), 1),
                createBaseVNode("span", null, toDisplayString(filter.label), 1),
                createBaseVNode("span", _hoisted_11$b, toDisplayString(getFilterCount(filter.id)), 1)
              ], 10, _hoisted_9$c);
            }), 64))
          ])
        ]),
        createBaseVNode("div", _hoisted_12$9, [
          createVNode(_sfc_main$s, {
            files: filteredFiles.value,
            "selected-paths": __props.selectedFiles,
            onToggleSelect: _cache[2] || (_cache[2] = ($event) => _ctx.$emit("toggle-select", $event)),
            onSelectFolder: _cache[3] || (_cache[3] = ($event) => _ctx.$emit("select-folder", $event)),
            onPreviewFile: _cache[4] || (_cache[4] = ($event) => _ctx.$emit("preview-file", $event))
          }, null, 8, ["files", "selected-paths"])
        ]),
        createBaseVNode("div", _hoisted_13$9, [
          createBaseVNode("button", {
            onClick: _cache[5] || (_cache[5] = ($event) => _ctx.$emit("select-all")),
            class: "action-btn action-btn-success btn-sm flex-1"
          }, toDisplayString(unref(t)("git.selectAll")), 1),
          createBaseVNode("button", {
            onClick: _cache[6] || (_cache[6] = ($event) => _ctx.$emit("clear-selection")),
            class: "action-btn action-btn-danger btn-sm flex-1"
          }, toDisplayString(unref(t)("git.clearSelection")), 1)
        ])
      ]);
    };
  }
});
const GitFileList = /* @__PURE__ */ _export_sfc(_sfc_main$r, [["__scopeId", "data-v-b04403dd"]]);
const _hoisted_1$p = { class: "flex-1 overflow-auto" };
const _hoisted_2$o = {
  key: 0,
  class: "p-4 text-center text-gray-400"
};
const _hoisted_3$n = { class: "text-sm mt-2" };
const _hoisted_4$l = {
  key: 1,
  class: "p-4 space-y-4"
};
const _hoisted_5$h = { class: "branch-info" };
const _hoisted_6$f = { class: "flex-1" };
const _hoisted_7$f = { class: "text-xs text-gray-400" };
const _hoisted_8$d = { class: "ml-2 text-sm text-white font-medium" };
const _hoisted_9$b = ["disabled"];
const _hoisted_10$b = { class: "space-y-3" };
const _hoisted_11$a = { class: "flex border-b border-gray-700" };
const _hoisted_12$8 = {
  key: 0,
  class: "p-3 bg-indigo-900/30 rounded-lg border border-indigo-500/30"
};
const _hoisted_13$8 = { class: "flex items-center justify-between" };
const _hoisted_14$8 = { class: "text-xs text-indigo-300" };
const _hoisted_15$6 = { class: "ml-2 text-sm text-white font-medium" };
const _hoisted_16$6 = {
  key: 0,
  class: "text-xs text-emerald-400 ml-2"
};
const _hoisted_17$6 = { class: "text-xs text-gray-400 mt-1" };
const _hoisted_18$6 = {
  key: 1,
  class: "branches-list space-y-1"
};
const _hoisted_19$6 = ["onClick"];
const _hoisted_20$6 = { class: "truncate" };
const _hoisted_21$4 = {
  key: 0,
  class: "text-xs text-emerald-400 ml-auto"
};
const _hoisted_22$4 = {
  key: 1,
  class: "text-xs text-indigo-400 ml-auto"
};
const _hoisted_23$4 = {
  key: 2,
  class: "branches-list space-y-1"
};
const _hoisted_24$4 = ["onClick"];
const _hoisted_25$3 = { class: "flex items-start gap-2" };
const _hoisted_26$3 = { class: "text-xs text-amber-400 font-mono flex-shrink-0" };
const _hoisted_27$2 = { class: "flex-1 min-w-0" };
const _hoisted_28$2 = { class: "truncate text-white" };
const _hoisted_29$1 = { class: "text-xs text-gray-400" };
const _sfc_main$q = /* @__PURE__ */ defineComponent({
  __name: "GitLocalPanel",
  props: {
    isGitRepo: { type: Boolean },
    currentBranch: {},
    branches: {},
    commits: {},
    commitsLoaded: { type: Boolean },
    selectedRef: {},
    filesAtRef: {},
    selectedFiles: {}
  },
  emits: ["select-ref", "clear-ref", "load-commits", "open-diff", "toggle-file", "select-folder", "select-all", "clear-selection", "preview-file"],
  setup(__props, { emit: __emit }) {
    const emit = __emit;
    const { t } = useI18n();
    const refType = ref("branches");
    function switchToCommits() {
      refType.value = "commits";
      emit("load-commits");
    }
    function formatDate(dateStr) {
      try {
        return new Date(dateStr).toLocaleDateString();
      } catch {
        return dateStr;
      }
    }
    return (_ctx, _cache) => {
      return openBlock(), createElementBlock("div", _hoisted_1$p, [
        !__props.isGitRepo ? (openBlock(), createElementBlock("div", _hoisted_2$o, [
          _cache[9] || (_cache[9] = createBaseVNode("svg", {
            class: "w-12 h-12 mx-auto mb-3 text-gray-600",
            fill: "none",
            stroke: "currentColor",
            viewBox: "0 0 24 24"
          }, [
            createBaseVNode("path", {
              "stroke-linecap": "round",
              "stroke-linejoin": "round",
              "stroke-width": "2",
              d: "M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"
            })
          ], -1)),
          createBaseVNode("p", null, toDisplayString(unref(t)("git.notGitRepo")), 1),
          createBaseVNode("p", _hoisted_3$n, toDisplayString(unref(t)("git.openGitProject")), 1)
        ])) : (openBlock(), createElementBlock("div", _hoisted_4$l, [
          createBaseVNode("div", _hoisted_5$h, [
            _cache[10] || (_cache[10] = createBaseVNode("svg", {
              class: "w-5 h-5 text-emerald-400",
              fill: "none",
              stroke: "currentColor",
              viewBox: "0 0 24 24"
            }, [
              createBaseVNode("path", {
                "stroke-linecap": "round",
                "stroke-linejoin": "round",
                "stroke-width": "2",
                d: "M13 10V3L4 14h7v7l9-11h-7z"
              })
            ], -1)),
            createBaseVNode("div", _hoisted_6$f, [
              createBaseVNode("span", _hoisted_7$f, toDisplayString(unref(t)("git.currentBranch")) + ":", 1),
              createBaseVNode("span", _hoisted_8$d, toDisplayString(__props.currentBranch), 1)
            ]),
            createBaseVNode("button", {
              onClick: _cache[0] || (_cache[0] = ($event) => _ctx.$emit("select-ref", __props.currentBranch)),
              class: "btn btn-primary btn-sm",
              disabled: __props.selectedRef === __props.currentBranch
            }, toDisplayString(unref(t)("git.load")), 9, _hoisted_9$b)
          ]),
          createBaseVNode("div", _hoisted_10$b, [
            createBaseVNode("div", _hoisted_11$a, [
              createBaseVNode("button", {
                onClick: _cache[1] || (_cache[1] = ($event) => refType.value = "branches"),
                class: normalizeClass(["tab-btn text-sm", refType.value === "branches" ? "tab-btn-active" : "tab-btn-inactive"])
              }, toDisplayString(unref(t)("git.branches")) + " (" + toDisplayString(__props.branches.length) + ") ", 3),
              createBaseVNode("button", {
                onClick: switchToCommits,
                class: normalizeClass(["tab-btn text-sm", refType.value === "commits" ? "tab-btn-active" : "tab-btn-inactive"])
              }, toDisplayString(unref(t)("git.commits")) + " " + toDisplayString(__props.commitsLoaded ? `(${__props.commits.length})` : ""), 3),
              __props.branches.length >= 2 ? (openBlock(), createElementBlock("button", {
                key: 0,
                onClick: _cache[2] || (_cache[2] = ($event) => _ctx.$emit("open-diff")),
                class: "ml-auto px-2 py-1 text-xs text-purple-400 hover:text-purple-300 hover:bg-purple-500/10 rounded transition-colors flex items-center gap-1"
              }, [
                _cache[11] || (_cache[11] = createBaseVNode("svg", {
                  class: "w-3.5 h-3.5",
                  fill: "none",
                  stroke: "currentColor",
                  viewBox: "0 0 24 24"
                }, [
                  createBaseVNode("path", {
                    "stroke-linecap": "round",
                    "stroke-linejoin": "round",
                    "stroke-width": "2",
                    d: "M8 7h12m0 0l-4-4m4 4l-4 4m0 6H4m0 0l4 4m-4-4l4-4"
                  })
                ], -1)),
                createTextVNode(" " + toDisplayString(unref(t)("git.diff")), 1)
              ])) : createCommentVNode("", true)
            ]),
            __props.selectedRef ? (openBlock(), createElementBlock("div", _hoisted_12$8, [
              createBaseVNode("div", _hoisted_13$8, [
                createBaseVNode("div", null, [
                  createBaseVNode("span", _hoisted_14$8, toDisplayString(unref(t)("git.buildingFrom")) + ":", 1),
                  createBaseVNode("span", _hoisted_15$6, toDisplayString(__props.selectedRef), 1),
                  __props.selectedRef === __props.currentBranch ? (openBlock(), createElementBlock("span", _hoisted_16$6, "(" + toDisplayString(unref(t)("git.current")) + ")", 1)) : createCommentVNode("", true)
                ]),
                createBaseVNode("button", {
                  onClick: _cache[3] || (_cache[3] = ($event) => _ctx.$emit("clear-ref")),
                  class: "text-xs text-gray-400 hover:text-white"
                }, toDisplayString(unref(t)("git.clear")), 1)
              ]),
              createBaseVNode("p", _hoisted_17$6, toDisplayString(unref(t)("git.noCheckout")), 1)
            ])) : createCommentVNode("", true),
            refType.value === "branches" ? (openBlock(), createElementBlock("div", _hoisted_18$6, [
              (openBlock(true), createElementBlock(Fragment, null, renderList(__props.branches, (branch) => {
                return openBlock(), createElementBlock("button", {
                  key: branch,
                  onClick: ($event) => _ctx.$emit("select-ref", branch),
                  class: normalizeClass([
                    "w-full px-3 py-2 text-left text-sm rounded transition-colors flex items-center gap-2",
                    __props.selectedRef === branch ? "bg-indigo-600/20 text-indigo-300 border border-indigo-500/30" : branch === __props.currentBranch ? "bg-emerald-900/20 text-emerald-300 border border-emerald-500/20" : "text-gray-300 hover:bg-gray-700"
                  ])
                }, [
                  _cache[12] || (_cache[12] = createBaseVNode("svg", {
                    class: "w-4 h-4 flex-shrink-0",
                    fill: "none",
                    stroke: "currentColor",
                    viewBox: "0 0 24 24"
                  }, [
                    createBaseVNode("path", {
                      "stroke-linecap": "round",
                      "stroke-linejoin": "round",
                      "stroke-width": "2",
                      d: "M13 10V3L4 14h7v7l9-11h-7z"
                    })
                  ], -1)),
                  createBaseVNode("span", _hoisted_20$6, toDisplayString(branch), 1),
                  branch === __props.currentBranch ? (openBlock(), createElementBlock("span", _hoisted_21$4, "(" + toDisplayString(unref(t)("git.current")) + ")", 1)) : __props.selectedRef === branch ? (openBlock(), createElementBlock("span", _hoisted_22$4, "(" + toDisplayString(unref(t)("git.selected")) + ")", 1)) : createCommentVNode("", true)
                ], 10, _hoisted_19$6);
              }), 128))
            ])) : createCommentVNode("", true),
            refType.value === "commits" ? (openBlock(), createElementBlock("div", _hoisted_23$4, [
              (openBlock(true), createElementBlock(Fragment, null, renderList(__props.commits, (commit) => {
                return openBlock(), createElementBlock("button", {
                  key: commit.hash,
                  onClick: ($event) => _ctx.$emit("select-ref", commit.hash),
                  class: normalizeClass([
                    "w-full px-3 py-2 text-left text-sm rounded transition-colors",
                    __props.selectedRef === commit.hash ? "bg-indigo-600/20 text-indigo-300 border border-indigo-500/30" : "text-gray-300 hover:bg-gray-700"
                  ])
                }, [
                  createBaseVNode("div", _hoisted_25$3, [
                    createBaseVNode("code", _hoisted_26$3, toDisplayString(commit.hash.slice(0, 7)), 1),
                    createBaseVNode("div", _hoisted_27$2, [
                      createBaseVNode("p", _hoisted_28$2, toDisplayString(commit.subject), 1),
                      createBaseVNode("p", _hoisted_29$1, toDisplayString(commit.author) + " • " + toDisplayString(formatDate(commit.date)), 1)
                    ])
                  ])
                ], 10, _hoisted_24$4);
              }), 128))
            ])) : createCommentVNode("", true)
          ]),
          __props.selectedRef && __props.filesAtRef.length > 0 ? (openBlock(), createBlock(GitFileList, {
            key: 0,
            title: `${unref(t)("git.filesAtRef")} ${__props.selectedRef?.slice(0, 7)}`,
            files: __props.filesAtRef,
            "selected-files": __props.selectedFiles,
            onToggleSelect: _cache[4] || (_cache[4] = ($event) => _ctx.$emit("toggle-file", $event)),
            onSelectFolder: _cache[5] || (_cache[5] = ($event) => _ctx.$emit("select-folder", $event)),
            onSelectAll: _cache[6] || (_cache[6] = ($event) => _ctx.$emit("select-all")),
            onClearSelection: _cache[7] || (_cache[7] = ($event) => _ctx.$emit("clear-selection")),
            onPreviewFile: _cache[8] || (_cache[8] = ($event) => _ctx.$emit("preview-file", $event))
          }, null, 8, ["title", "files", "selected-files"])) : createCommentVNode("", true)
        ]))
      ]);
    };
  }
});
const GitLocalPanel = /* @__PURE__ */ _export_sfc(_sfc_main$q, [["__scopeId", "data-v-49abbb83"]]);
const _hoisted_1$o = { class: "flex-1 overflow-auto p-4 space-y-4" };
const _hoisted_2$n = { class: "space-y-3" };
const _hoisted_3$m = { class: "block text-sm text-gray-300 mb-2" };
const _hoisted_4$k = { class: "flex gap-2" };
const _hoisted_5$g = ["placeholder"];
const _hoisted_6$e = ["disabled"];
const _hoisted_7$e = {
  key: 0,
  class: "animate-spin w-4 h-4",
  fill: "none",
  viewBox: "0 0 24 24"
};
const _hoisted_8$c = { key: 1 };
const _hoisted_9$a = { class: "text-xs text-gray-400" };
const _hoisted_10$a = {
  key: 0,
  class: "text-emerald-400"
};
const _hoisted_11$9 = {
  key: 1,
  class: "text-orange-400"
};
const _hoisted_12$7 = { key: 2 };
const _hoisted_13$7 = {
  key: 0,
  class: "space-y-4"
};
const _hoisted_14$7 = { class: "branch-info" };
const _hoisted_15$5 = { class: "flex-1" };
const _hoisted_16$5 = ["value"];
const _hoisted_17$5 = {
  key: 1,
  class: "p-3 bg-amber-900/20 rounded-lg border border-amber-500/30"
};
const _hoisted_18$5 = { class: "text-xs text-amber-300 mb-2" };
const _hoisted_19$5 = ["disabled"];
const _hoisted_20$5 = {
  key: 2,
  class: "p-3 bg-emerald-900/30 rounded-lg border border-emerald-500/30"
};
const _hoisted_21$3 = { class: "flex items-center gap-2 mb-2" };
const _hoisted_22$3 = { class: "text-sm text-emerald-300" };
const _hoisted_23$3 = { class: "text-xs text-gray-400 truncate" };
const _hoisted_24$3 = { class: "flex gap-2 mt-3" };
const _sfc_main$p = /* @__PURE__ */ defineComponent({
  __name: "GitRemotePanel",
  props: {
    remoteUrl: {},
    isGitHub: { type: Boolean },
    isGitLab: { type: Boolean },
    isLoading: { type: Boolean },
    isCloning: { type: Boolean },
    repoLoaded: { type: Boolean },
    branches: {},
    selectedBranch: {},
    files: {},
    selectedFiles: {},
    clonedPath: {}
  },
  emits: ["load-repo", "change-branch", "clone", "open-cloned", "cleanup-cloned", "toggle-file", "select-folder", "select-all", "clear-selection", "preview-file"],
  setup(__props) {
    const props = __props;
    const { t } = useI18n();
    const urlInput = ref(props.remoteUrl);
    const selectedBranchLocal = ref(props.selectedBranch);
    watch(() => props.remoteUrl, (val) => {
      urlInput.value = val;
    });
    watch(() => props.selectedBranch, (val) => {
      selectedBranchLocal.value = val;
    });
    return (_ctx, _cache) => {
      return openBlock(), createElementBlock("div", _hoisted_1$o, [
        createBaseVNode("div", _hoisted_2$n, [
          createBaseVNode("div", null, [
            createBaseVNode("label", _hoisted_3$m, toDisplayString(unref(t)("git.repoUrl")), 1),
            createBaseVNode("div", _hoisted_4$k, [
              withDirectives(createBaseVNode("input", {
                "onUpdate:modelValue": _cache[0] || (_cache[0] = ($event) => urlInput.value = $event),
                type: "text",
                placeholder: unref(t)("git.urlPlaceholder"),
                class: "input flex-1",
                onKeydown: _cache[1] || (_cache[1] = withKeys(($event) => _ctx.$emit("load-repo", urlInput.value), ["enter"]))
              }, null, 40, _hoisted_5$g), [
                [vModelText, urlInput.value]
              ]),
              createBaseVNode("button", {
                onClick: _cache[2] || (_cache[2] = ($event) => _ctx.$emit("load-repo", urlInput.value)),
                disabled: !urlInput.value || __props.isLoading,
                class: "btn btn-primary"
              }, [
                __props.isLoading ? (openBlock(), createElementBlock("svg", _hoisted_7$e, [..._cache[13] || (_cache[13] = [
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
                ])])) : (openBlock(), createElementBlock("span", _hoisted_8$c, toDisplayString(unref(t)("git.load")), 1))
              ], 8, _hoisted_6$e)
            ])
          ]),
          createBaseVNode("p", _hoisted_9$a, [
            __props.isGitHub ? (openBlock(), createElementBlock("span", _hoisted_10$a, "✓ GitHub API (fast, no clone)")) : __props.isGitLab ? (openBlock(), createElementBlock("span", _hoisted_11$9, "✓ GitLab API (fast, no clone)")) : (openBlock(), createElementBlock("span", _hoisted_12$7, toDisplayString(unref(t)("git.urlHint")), 1))
          ])
        ]),
        __props.repoLoaded ? (openBlock(), createElementBlock("div", _hoisted_13$7, [
          createBaseVNode("div", _hoisted_14$7, [
            _cache[14] || (_cache[14] = createBaseVNode("svg", {
              class: "w-5 h-5 text-emerald-400",
              fill: "none",
              stroke: "currentColor",
              viewBox: "0 0 24 24"
            }, [
              createBaseVNode("path", {
                "stroke-linecap": "round",
                "stroke-linejoin": "round",
                "stroke-width": "2",
                d: "M13 10V3L4 14h7v7l9-11h-7z"
              })
            ], -1)),
            createBaseVNode("div", _hoisted_15$5, [
              withDirectives(createBaseVNode("select", {
                "onUpdate:modelValue": _cache[3] || (_cache[3] = ($event) => selectedBranchLocal.value = $event),
                onChange: _cache[4] || (_cache[4] = ($event) => _ctx.$emit("change-branch", selectedBranchLocal.value)),
                class: "input w-full text-sm"
              }, [
                (openBlock(true), createElementBlock(Fragment, null, renderList(__props.branches, (branch) => {
                  return openBlock(), createElementBlock("option", {
                    key: branch.name,
                    value: branch.name
                  }, toDisplayString(branch.name), 9, _hoisted_16$5);
                }), 128))
              ], 544), [
                [vModelSelect, selectedBranchLocal.value]
              ])
            ])
          ]),
          __props.files.length > 0 ? (openBlock(), createBlock(GitFileList, {
            key: 0,
            title: `${unref(t)("git.filesAtRef")} ${__props.selectedBranch}`,
            files: __props.files,
            "selected-files": __props.selectedFiles,
            onToggleSelect: _cache[5] || (_cache[5] = ($event) => _ctx.$emit("toggle-file", $event)),
            onSelectFolder: _cache[6] || (_cache[6] = ($event) => _ctx.$emit("select-folder", $event)),
            onSelectAll: _cache[7] || (_cache[7] = ($event) => _ctx.$emit("select-all")),
            onClearSelection: _cache[8] || (_cache[8] = ($event) => _ctx.$emit("clear-selection")),
            onPreviewFile: _cache[9] || (_cache[9] = ($event) => _ctx.$emit("preview-file", $event))
          }, null, 8, ["title", "files", "selected-files"])) : createCommentVNode("", true)
        ])) : createCommentVNode("", true),
        !__props.isGitHub && !__props.isGitLab && urlInput.value && !__props.repoLoaded ? (openBlock(), createElementBlock("div", _hoisted_17$5, [
          createBaseVNode("p", _hoisted_18$5, toDisplayString(unref(t)("git.cloneRequired")), 1),
          createBaseVNode("button", {
            onClick: _cache[10] || (_cache[10] = ($event) => _ctx.$emit("clone")),
            disabled: __props.isCloning,
            class: "action-btn action-btn-accent w-full"
          }, toDisplayString(__props.isCloning ? unref(t)("git.cloning") : unref(t)("git.clone")), 9, _hoisted_19$5)
        ])) : createCommentVNode("", true),
        __props.clonedPath ? (openBlock(), createElementBlock("div", _hoisted_20$5, [
          createBaseVNode("div", _hoisted_21$3, [
            _cache[15] || (_cache[15] = createBaseVNode("svg", {
              class: "w-4 h-4 text-emerald-400",
              fill: "none",
              stroke: "currentColor",
              viewBox: "0 0 24 24"
            }, [
              createBaseVNode("path", {
                "stroke-linecap": "round",
                "stroke-linejoin": "round",
                "stroke-width": "2",
                d: "M5 13l4 4L19 7"
              })
            ], -1)),
            createBaseVNode("span", _hoisted_22$3, toDisplayString(unref(t)("git.cloned")), 1)
          ]),
          createBaseVNode("p", _hoisted_23$3, toDisplayString(__props.clonedPath), 1),
          createBaseVNode("div", _hoisted_24$3, [
            createBaseVNode("button", {
              onClick: _cache[11] || (_cache[11] = ($event) => _ctx.$emit("open-cloned")),
              class: "btn btn-primary btn-sm flex-1"
            }, toDisplayString(unref(t)("git.openProject")), 1),
            createBaseVNode("button", {
              onClick: _cache[12] || (_cache[12] = ($event) => _ctx.$emit("cleanup-cloned")),
              class: "btn btn-ghost btn-sm"
            }, toDisplayString(unref(t)("git.remove")), 1)
          ])
        ])) : createCommentVNode("", true)
      ]);
    };
  }
});
const _hoisted_1$n = { class: "relative" };
const _hoisted_2$m = ["title"];
const _hoisted_3$l = {
  key: 0,
  class: "absolute right-0 top-full mt-1 w-64 bg-gray-800 border border-gray-700 rounded-lg shadow-xl z-50"
};
const _hoisted_4$j = { class: "p-2 border-b border-gray-700 flex items-center justify-between" };
const _hoisted_5$f = { class: "text-xs text-gray-400" };
const _hoisted_6$d = { class: "max-h-48 overflow-y-auto" };
const _hoisted_7$d = ["onClick"];
const _hoisted_8$b = {
  key: 0,
  class: "w-4 h-4 text-gray-400",
  fill: "currentColor",
  viewBox: "0 0 24 24"
};
const _hoisted_9$9 = {
  key: 1,
  class: "w-4 h-4 text-orange-400",
  fill: "currentColor",
  viewBox: "0 0 24 24"
};
const _hoisted_10$9 = { class: "text-sm text-gray-300 truncate" };
const _sfc_main$o = /* @__PURE__ */ defineComponent({
  __name: "RecentReposDropdown",
  props: {
    repos: {}
  },
  emits: ["select", "clear"],
  setup(__props, { emit: __emit }) {
    const { t } = useI18n();
    const emit = __emit;
    const isOpen = ref(false);
    function handleSelect(repo) {
      emit("select", repo);
      isOpen.value = false;
    }
    function handleClear() {
      emit("clear");
      isOpen.value = false;
    }
    return (_ctx, _cache) => {
      return openBlock(), createElementBlock("div", _hoisted_1$n, [
        createBaseVNode("button", {
          onClick: _cache[0] || (_cache[0] = ($event) => isOpen.value = !isOpen.value),
          class: "action-btn text-gray-400 hover:text-white",
          title: unref(t)("git.recentRepos")
        }, [..._cache[2] || (_cache[2] = [
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
              d: "M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"
            })
          ], -1)
        ])], 8, _hoisted_2$m),
        isOpen.value ? (openBlock(), createElementBlock("div", _hoisted_3$l, [
          createBaseVNode("div", _hoisted_4$j, [
            createBaseVNode("span", _hoisted_5$f, toDisplayString(unref(t)("git.recentRepos")), 1),
            createBaseVNode("button", {
              onClick: handleClear,
              class: "text-xs text-gray-400 hover:text-red-400"
            }, toDisplayString(unref(t)("common.clear")), 1)
          ]),
          createBaseVNode("div", _hoisted_6$d, [
            (openBlock(true), createElementBlock(Fragment, null, renderList(__props.repos, (repo) => {
              return openBlock(), createElementBlock("button", {
                key: repo.url,
                onClick: ($event) => handleSelect(repo),
                class: "w-full px-3 py-2 text-left hover:bg-gray-700/50 flex items-center gap-2"
              }, [
                repo.isGitHub ? (openBlock(), createElementBlock("svg", _hoisted_8$b, [..._cache[3] || (_cache[3] = [
                  createBaseVNode("path", { d: "M12 0c-6.626 0-12 5.373-12 12 0 5.302 3.438 9.8 8.207 11.387.599.111.793-.261.793-.577v-2.234c-3.338.726-4.033-1.416-4.033-1.416-.546-1.387-1.333-1.756-1.333-1.756-1.089-.745.083-.729.083-.729 1.205.084 1.839 1.237 1.839 1.237 1.07 1.834 2.807 1.304 3.492.997.107-.775.418-1.305.762-1.604-2.665-.305-5.467-1.334-5.467-5.931 0-1.311.469-2.381 1.236-3.221-.124-.303-.535-1.524.117-3.176 0 0 1.008-.322 3.301 1.23.957-.266 1.983-.399 3.003-.404 1.02.005 2.047.138 3.006.404 2.291-1.552 3.297-1.23 3.297-1.23.653 1.653.242 2.874.118 3.176.77.84 1.235 1.911 1.235 3.221 0 4.609-2.807 5.624-5.479 5.921.43.372.823 1.102.823 2.222v3.293c0 .319.192.694.801.576 4.765-1.589 8.199-6.086 8.199-11.386 0-6.627-5.373-12-12-12z" }, null, -1)
                ])])) : (openBlock(), createElementBlock("svg", _hoisted_9$9, [..._cache[4] || (_cache[4] = [
                  createBaseVNode("path", { d: "M22.65 14.39L12 22.13 1.35 14.39a.84.84 0 01-.3-.94l1.22-3.78 2.44-7.51A.42.42 0 014.82 2a.43.43 0 01.58 0 .42.42 0 01.11.18l2.44 7.49h8.1l2.44-7.51A.42.42 0 0118.6 2a.43.43 0 01.58 0 .42.42 0 01.11.18l2.44 7.51L23 13.45a.84.84 0 01-.35.94z" }, null, -1)
                ])])),
                createBaseVNode("span", _hoisted_10$9, toDisplayString(repo.name), 1)
              ], 8, _hoisted_7$d);
            }), 128))
          ])
        ])) : createCommentVNode("", true),
        isOpen.value ? (openBlock(), createElementBlock("div", {
          key: 1,
          class: "fixed inset-0 z-40",
          onClick: _cache[1] || (_cache[1] = ($event) => isOpen.value = false)
        })) : createCommentVNode("", true)
      ]);
    };
  }
});
const _hoisted_1$m = { class: "flex gap-1 px-2 pb-2" };
const _sfc_main$n = /* @__PURE__ */ defineComponent({
  __name: "GitSourceTabs",
  props: {
    modelValue: {}
  },
  emits: ["change"],
  setup(__props) {
    const { t } = useI18n();
    return (_ctx, _cache) => {
      return openBlock(), createElementBlock("div", _hoisted_1$m, [
        createBaseVNode("button", {
          onClick: _cache[0] || (_cache[0] = ($event) => _ctx.$emit("change", "local")),
          class: normalizeClass(["tab-btn", __props.modelValue === "local" ? "tab-btn-active tab-btn-active-indigo" : "tab-btn-inactive"])
        }, [
          _cache[2] || (_cache[2] = createBaseVNode("svg", {
            class: "w-4 h-4",
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
          createTextVNode(" " + toDisplayString(unref(t)("git.localGit")), 1)
        ], 2),
        createBaseVNode("button", {
          onClick: _cache[1] || (_cache[1] = ($event) => _ctx.$emit("change", "remote")),
          class: normalizeClass(["tab-btn", __props.modelValue === "remote" ? "tab-btn-active tab-btn-active-purple" : "tab-btn-inactive"])
        }, [
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
              d: "M21 12a9 9 0 01-9 9m9-9a9 9 0 00-9-9m9 9H3m9 9a9 9 0 01-9-9m9 9c1.657 0 3-4.03 3-9s-1.343-9-3-9m0 18c-1.657 0-3-4.03-3-9s1.343-9 3-9m-9 9a9 9 0 019-9"
            })
          ], -1)),
          createTextVNode(" " + toDisplayString(unref(t)("git.remoteUrl")), 1)
        ], 2)
      ]);
    };
  }
});
const _hoisted_1$l = { class: "border-b border-gray-700/30" };
const _hoisted_2$l = { class: "flex items-center justify-between p-3" };
const _hoisted_3$k = { class: "section-title" };
const _hoisted_4$i = { class: "section-title-text" };
const _sfc_main$m = /* @__PURE__ */ defineComponent({
  __name: "GitSourceHeader",
  props: {
    sourceType: {},
    recentRepos: {}
  },
  emits: ["change-source", "select-recent", "clear-recent"],
  setup(__props) {
    const { t } = useI18n();
    return (_ctx, _cache) => {
      return openBlock(), createElementBlock("div", _hoisted_1$l, [
        createBaseVNode("div", _hoisted_2$l, [
          createBaseVNode("div", _hoisted_3$k, [
            _cache[3] || (_cache[3] = createBaseVNode("div", { class: "section-icon section-icon-orange" }, [
              createBaseVNode("svg", {
                fill: "none",
                stroke: "currentColor",
                viewBox: "0 0 24 24"
              }, [
                createBaseVNode("path", {
                  "stroke-linecap": "round",
                  "stroke-linejoin": "round",
                  "stroke-width": "2",
                  d: "M10 20l4-16m4 4l4 4-4 4M6 16l-4-4 4-4"
                })
              ])
            ], -1)),
            createBaseVNode("h2", _hoisted_4$i, toDisplayString(unref(t)("git.title")), 1)
          ]),
          __props.recentRepos.length > 0 ? (openBlock(), createBlock(_sfc_main$o, {
            key: 0,
            repos: __props.recentRepos,
            onSelect: _cache[0] || (_cache[0] = ($event) => _ctx.$emit("select-recent", $event)),
            onClear: _cache[1] || (_cache[1] = ($event) => _ctx.$emit("clear-recent"))
          }, null, 8, ["repos"])) : createCommentVNode("", true)
        ]),
        createVNode(_sfc_main$n, {
          "model-value": __props.sourceType,
          onChange: _cache[2] || (_cache[2] = ($event) => _ctx.$emit("change-source", $event))
        }, null, 8, ["model-value"])
      ]);
    };
  }
});
const _hoisted_1$k = { class: "h-full flex flex-col bg-transparent" };
const _hoisted_2$k = {
  key: 2,
  class: "border-t border-gray-700 p-4"
};
const _hoisted_3$j = { class: "flex items-center justify-between mb-3" };
const _hoisted_4$h = { class: "text-sm text-gray-300" };
const _hoisted_5$e = ["disabled"];
const _hoisted_6$c = {
  key: 0,
  class: "animate-spin w-4 h-4 mr-2",
  fill: "none",
  viewBox: "0 0 24 24"
};
const _hoisted_7$c = {
  key: 3,
  class: "loading-overlay absolute"
};
const _hoisted_8$a = { class: "text-center" };
const _hoisted_9$8 = { class: "text-sm text-gray-400" };
const _hoisted_10$8 = {
  key: 4,
  class: "border-t border-gray-700 p-4"
};
const _hoisted_11$8 = { class: "flex items-center justify-between mb-3" };
const _hoisted_12$6 = { class: "text-sm text-gray-300" };
const _hoisted_13$6 = ["disabled"];
const _hoisted_14$6 = {
  key: 0,
  class: "animate-spin w-4 h-4 mr-2",
  fill: "none",
  viewBox: "0 0 24 24"
};
const _sfc_main$l = /* @__PURE__ */ defineComponent({
  __name: "GitSourceSelector",
  setup(__props) {
    const logger2 = useLogger("GitSourceSelector");
    const { t } = useI18n();
    const projectStore = useProjectStore();
    const uiStore = useUIStore();
    const contextStore = useContextStore();
    const sourceType = ref("local");
    const git2 = useGitSource();
    const {
      isGitRepo,
      currentBranch,
      branches,
      commits,
      commitsLoaded,
      selectedRef,
      filesAtRef,
      selectedFiles,
      remoteUrl,
      isCloning,
      clonedPath,
      isGitHubRepo,
      isGitLabRepo,
      isLoadingRemote,
      remoteRepoLoaded,
      remoteBranches,
      remoteSelectedBranch,
      remoteFiles,
      remoteSelectedFiles,
      isLoading,
      loadingMessage,
      isBuilding,
      projectPath,
      recentRepos,
      loadCommits,
      selectRef,
      clearSelectedRef,
      toggleFileSelection,
      selectAllFiles,
      clearFileSelection,
      loadRemoteRepo: loadRemoteRepoBase,
      loadRemoteFiles,
      toggleRemoteFileSelection,
      selectAllRemoteFiles,
      clearRemoteFileSelection,
      loadRecentReposFromStorage,
      clearRecentRepos
    } = git2;
    const diffModalOpen = ref(false);
    const previewOpen = ref(false);
    const previewPath = ref("");
    const previewContent = ref("");
    const previewLoading = ref(false);
    const previewError = ref("");
    function handleSelectFolder(files2) {
      files2.forEach((f) => {
        if (!selectedFiles.value.has(f)) selectedFiles.value.add(f);
      });
      selectedFiles.value = new Set(selectedFiles.value);
    }
    function handleSelectRemoteFolder(files2) {
      files2.forEach((f) => {
        if (!remoteSelectedFiles.value.has(f)) remoteSelectedFiles.value.add(f);
      });
      remoteSelectedFiles.value = new Set(remoteSelectedFiles.value);
    }
    async function checkGitRepo() {
      if (!projectPath.value) return;
      isLoading.value = true;
      loadingMessage.value = "Checking repository...";
      commitsLoaded.value = false;
      branches.value = [];
      commits.value = [];
      try {
        isGitRepo.value = await apiService.isGitRepository(projectPath.value);
        if (isGitRepo.value) {
          currentBranch.value = await apiService.getCurrentBranch(projectPath.value);
          const result = await apiService.getBranches(projectPath.value);
          branches.value = JSON.parse(result);
        }
      } catch (err) {
        logger2.error("Failed to check git repo", err);
        isGitRepo.value = false;
      } finally {
        isLoading.value = false;
      }
    }
    async function loadRemoteRepo(url) {
      remoteUrl.value = url;
      await loadRemoteRepoBase();
    }
    async function handleChangeBranch(branch) {
      remoteSelectedBranch.value = branch;
      await loadRemoteFiles();
    }
    function handleSelectRecentRepo(repo) {
      remoteUrl.value = repo.url;
      sourceType.value = "remote";
      loadRemoteRepoBase();
    }
    async function buildContextFromRef() {
      if (!selectedRef.value || selectedFiles.value.size === 0) return;
      isBuilding.value = true;
      try {
        const files2 = Array.from(selectedFiles.value);
        const content = await apiService.buildContextAtRef(projectPath.value, files2, selectedRef.value);
        contextStore.setRawContext(content, files2.length);
        uiStore.addToast(`Context built from ${selectedRef.value.slice(0, 7)}: ${files2.length} files`, "success");
      } catch (err) {
        logger2.error("Failed to build context from ref", err);
        uiStore.addToast("Failed to build context from ref", "error");
      } finally {
        isBuilding.value = false;
      }
    }
    async function buildContextFromRemote() {
      if (!remoteUrl.value || remoteSelectedFiles.value.size === 0) return;
      isBuilding.value = true;
      try {
        const files2 = Array.from(remoteSelectedFiles.value);
        let content = "", source = "";
        if (isGitHubRepo.value) {
          content = await apiService.gitHubBuildContext(remoteUrl.value, files2, remoteSelectedBranch.value);
          source = "GitHub";
        } else if (isGitLabRepo.value) {
          content = await apiService.gitLabBuildContext(remoteUrl.value, files2, remoteSelectedBranch.value);
          source = "GitLab";
        }
        contextStore.setRawContext(content, files2.length);
        uiStore.addToast(`Context built from ${source}: ${files2.length} files`, "success");
      } catch (err) {
        logger2.error("Failed to build context from remote", err);
        uiStore.addToast("Failed to build context", "error");
      } finally {
        isBuilding.value = false;
      }
    }
    async function cloneRemote() {
      if (!remoteUrl.value) return;
      isCloning.value = true;
      isLoading.value = true;
      loadingMessage.value = "Cloning repository...";
      try {
        clonedPath.value = await apiService.cloneRepository(remoteUrl.value);
        uiStore.addToast("Repository cloned successfully", "success");
      } catch (err) {
        logger2.error("Failed to clone repository", err);
        uiStore.addToast("Failed to clone repository", "error");
      } finally {
        isCloning.value = false;
        isLoading.value = false;
      }
    }
    async function openClonedRepo() {
      if (!clonedPath.value) return;
      await projectStore.openProjectByPath(clonedPath.value);
      sourceType.value = "local";
    }
    async function cleanupClonedRepo() {
      if (!clonedPath.value) return;
      try {
        await apiService.cleanupTempRepository(clonedPath.value);
        clonedPath.value = null;
        uiStore.addToast("Temporary repository removed", "success");
      } catch (err) {
        logger2.warn("Failed to cleanup cloned repo", err);
      }
    }
    async function handlePreviewFile(filePath) {
      previewPath.value = filePath;
      previewContent.value = "";
      previewError.value = "";
      previewLoading.value = true;
      previewOpen.value = true;
      try {
        let content = "";
        if (sourceType.value === "local" && selectedRef.value) {
          content = await apiService.getFileAtRef(projectPath.value, filePath, selectedRef.value);
        } else if (sourceType.value === "remote") {
          if (isGitHubRepo.value) {
            content = await apiService.gitHubGetFileContent(remoteUrl.value, filePath, remoteSelectedBranch.value);
          } else if (isGitLabRepo.value) {
            content = await apiService.gitLabGetFileContent(remoteUrl.value, filePath, remoteSelectedBranch.value);
          }
        }
        previewContent.value = content;
      } catch (err) {
        logger2.error("Failed to preview file", err);
        previewError.value = t("error.loadFailed");
      } finally {
        previewLoading.value = false;
      }
    }
    watch(() => projectStore.currentPath, checkGitRepo, { immediate: true });
    onMounted(() => {
      checkGitRepo();
      loadRecentReposFromStorage();
    });
    return (_ctx, _cache) => {
      return openBlock(), createElementBlock("div", _hoisted_1$k, [
        createVNode(_sfc_main$m, {
          "source-type": sourceType.value,
          "recent-repos": unref(recentRepos),
          onChangeSource: _cache[0] || (_cache[0] = ($event) => sourceType.value = $event),
          onSelectRecent: handleSelectRecentRepo,
          onClearRecent: unref(clearRecentRepos)
        }, null, 8, ["source-type", "recent-repos", "onClearRecent"]),
        sourceType.value === "local" ? (openBlock(), createBlock(GitLocalPanel, {
          key: 0,
          "is-git-repo": unref(isGitRepo),
          "current-branch": unref(currentBranch),
          branches: unref(branches),
          commits: unref(commits),
          "commits-loaded": unref(commitsLoaded),
          "selected-ref": unref(selectedRef),
          "files-at-ref": unref(filesAtRef),
          "selected-files": unref(selectedFiles),
          onSelectRef: unref(selectRef),
          onClearRef: unref(clearSelectedRef),
          onLoadCommits: unref(loadCommits),
          onOpenDiff: _cache[1] || (_cache[1] = ($event) => diffModalOpen.value = true),
          onToggleFile: unref(toggleFileSelection),
          onSelectFolder: handleSelectFolder,
          onSelectAll: unref(selectAllFiles),
          onClearSelection: unref(clearFileSelection),
          onPreviewFile: handlePreviewFile
        }, null, 8, ["is-git-repo", "current-branch", "branches", "commits", "commits-loaded", "selected-ref", "files-at-ref", "selected-files", "onSelectRef", "onClearRef", "onLoadCommits", "onToggleFile", "onSelectAll", "onClearSelection"])) : createCommentVNode("", true),
        sourceType.value === "remote" ? (openBlock(), createBlock(_sfc_main$p, {
          key: 1,
          "remote-url": unref(remoteUrl),
          "is-git-hub": unref(isGitHubRepo),
          "is-git-lab": unref(isGitLabRepo),
          "is-loading": unref(isLoadingRemote),
          "is-cloning": unref(isCloning),
          "repo-loaded": unref(remoteRepoLoaded),
          branches: unref(remoteBranches),
          "selected-branch": unref(remoteSelectedBranch),
          files: unref(remoteFiles),
          "selected-files": unref(remoteSelectedFiles),
          "cloned-path": unref(clonedPath),
          onLoadRepo: loadRemoteRepo,
          onChangeBranch: handleChangeBranch,
          onClone: cloneRemote,
          onOpenCloned: openClonedRepo,
          onCleanupCloned: cleanupClonedRepo,
          onToggleFile: unref(toggleRemoteFileSelection),
          onSelectFolder: handleSelectRemoteFolder,
          onSelectAll: unref(selectAllRemoteFiles),
          onClearSelection: unref(clearRemoteFileSelection),
          onPreviewFile: handlePreviewFile
        }, null, 8, ["remote-url", "is-git-hub", "is-git-lab", "is-loading", "is-cloning", "repo-loaded", "branches", "selected-branch", "files", "selected-files", "cloned-path", "onToggleFile", "onSelectAll", "onClearSelection"])) : createCommentVNode("", true),
        sourceType.value === "remote" && unref(remoteSelectedFiles).size > 0 ? (openBlock(), createElementBlock("div", _hoisted_2$k, [
          createBaseVNode("div", _hoisted_3$j, [
            createBaseVNode("span", _hoisted_4$h, toDisplayString(unref(remoteSelectedFiles).size) + " " + toDisplayString(unref(t)("git.filesSelected")), 1)
          ]),
          createBaseVNode("button", {
            onClick: buildContextFromRemote,
            disabled: unref(isBuilding),
            class: "btn btn-primary w-full"
          }, [
            unref(isBuilding) ? (openBlock(), createElementBlock("svg", _hoisted_6$c, [..._cache[4] || (_cache[4] = [
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
            ])])) : createCommentVNode("", true),
            createTextVNode(" " + toDisplayString(unref(t)("git.buildContext")) + " " + toDisplayString(unref(remoteSelectedBranch)), 1)
          ], 8, _hoisted_5$e)
        ])) : createCommentVNode("", true),
        unref(isLoading) ? (openBlock(), createElementBlock("div", _hoisted_7$c, [
          createBaseVNode("div", _hoisted_8$a, [
            _cache[5] || (_cache[5] = createBaseVNode("svg", {
              class: "loading-spinner mx-auto mb-2",
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
            createBaseVNode("span", _hoisted_9$8, toDisplayString(unref(loadingMessage)), 1)
          ])
        ])) : createCommentVNode("", true),
        sourceType.value === "local" && unref(selectedFiles).size > 0 ? (openBlock(), createElementBlock("div", _hoisted_10$8, [
          createBaseVNode("div", _hoisted_11$8, [
            createBaseVNode("span", _hoisted_12$6, toDisplayString(unref(selectedFiles).size) + " " + toDisplayString(unref(t)("git.filesSelected")), 1)
          ]),
          createBaseVNode("button", {
            onClick: buildContextFromRef,
            disabled: unref(isBuilding),
            class: "btn btn-primary w-full"
          }, [
            unref(isBuilding) ? (openBlock(), createElementBlock("svg", _hoisted_14$6, [..._cache[6] || (_cache[6] = [
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
            ])])) : createCommentVNode("", true),
            createTextVNode(" " + toDisplayString(unref(t)("git.buildContext")) + " " + toDisplayString(unref(selectedRef)?.slice(0, 7) || "ref"), 1)
          ], 8, _hoisted_13$6)
        ])) : createCommentVNode("", true),
        createVNode(_sfc_main$t, {
          "is-open": previewOpen.value,
          "file-path": previewPath.value,
          content: previewContent.value,
          "is-loading": previewLoading.value,
          error: previewError.value,
          onClose: _cache[2] || (_cache[2] = ($event) => previewOpen.value = false)
        }, null, 8, ["is-open", "file-path", "content", "is-loading", "error"]),
        createVNode(_sfc_main$u, {
          "is-open": diffModalOpen.value,
          branches: unref(branches),
          "project-path": unref(projectPath),
          "current-branch": unref(currentBranch),
          onClose: _cache[3] || (_cache[3] = ($event) => diffModalOpen.value = false)
        }, null, 8, ["is-open", "branches", "project-path", "current-branch"])
      ]);
    };
  }
});
const _hoisted_1$j = { class: "sidebar-container" };
const _hoisted_2$j = { class: "tabs-container" };
const _hoisted_3$i = {
  key: 0,
  class: "tab-badge"
};
const _hoisted_4$g = { class: "sidebar-content" };
const _sfc_main$k = /* @__PURE__ */ defineComponent({
  __name: "LeftSidebar",
  emits: ["preview-file", "build-context"],
  setup(__props, { emit: __emit }) {
    const fileStore = useFileStore();
    const { t } = useI18n();
    const currentTab = ref("files");
    const tabIndex = computed(() => {
      const tabs = ["files", "git", "contexts"];
      return tabs.indexOf(currentTab.value);
    });
    const tabIndicatorClass = computed(() => {
      const classes = {
        files: "tabs-indicator-indigo",
        git: "tabs-indicator-orange",
        contexts: "tabs-indicator-purple"
      };
      return classes[currentTab.value];
    });
    const emit = __emit;
    function handlePreviewFile(filePath) {
      emit("preview-file", filePath);
    }
    function handleBuildContext() {
      emit("build-context");
    }
    watch(currentTab, (tab) => {
      try {
        localStorage.setItem("left-sidebar-tab", tab);
      } catch (err) {
        console.warn("Failed to save sidebar tab:", err);
      }
    });
    try {
      const savedTab = localStorage.getItem("left-sidebar-tab");
      if (savedTab === "files" || savedTab === "contexts" || savedTab === "git") {
        currentTab.value = savedTab;
      }
    } catch (err) {
      console.warn("Failed to load sidebar tab:", err);
    }
    return (_ctx, _cache) => {
      return openBlock(), createElementBlock("div", _hoisted_1$j, [
        createBaseVNode("div", _hoisted_2$j, [
          createBaseVNode("div", {
            class: normalizeClass(["tabs-indicator", tabIndicatorClass.value]),
            style: normalizeStyle({ transform: `translateX(${tabIndex.value * 100}%)` })
          }, null, 6),
          createBaseVNode("button", {
            onClick: _cache[0] || (_cache[0] = ($event) => currentTab.value = "files"),
            class: normalizeClass(["sidebar-tab", currentTab.value === "files" ? "sidebar-tab-active text-indigo-300" : ""])
          }, [
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
                d: "M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-6l-2-2H5a2 2 0 00-2 2z"
              })
            ], -1)),
            createBaseVNode("span", null, toDisplayString(unref(t)("tabs.files")), 1),
            unref(fileStore).selectedPaths.size > 0 ? (openBlock(), createElementBlock("span", _hoisted_3$i, toDisplayString(unref(fileStore).selectedPaths.size), 1)) : createCommentVNode("", true)
          ], 2),
          createBaseVNode("button", {
            onClick: _cache[1] || (_cache[1] = ($event) => currentTab.value = "git"),
            class: normalizeClass(["sidebar-tab", currentTab.value === "git" ? "sidebar-tab-active text-orange-300" : ""])
          }, [..._cache[4] || (_cache[4] = [
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
                d: "M10 20l4-16m4 4l4 4-4 4M6 16l-4-4 4-4"
              })
            ], -1),
            createBaseVNode("span", null, "Git", -1)
          ])], 2),
          createBaseVNode("button", {
            onClick: _cache[2] || (_cache[2] = ($event) => currentTab.value = "contexts"),
            class: normalizeClass(["sidebar-tab", currentTab.value === "contexts" ? "sidebar-tab-active text-purple-300" : ""])
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
                d: "M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"
              })
            ], -1)),
            createBaseVNode("span", null, toDisplayString(unref(t)("tabs.contexts")), 1)
          ], 2)
        ]),
        createBaseVNode("div", _hoisted_4$g, [
          withDirectives(createVNode(FileExplorer, {
            class: normalizeClass(["sidebar-tab-content", { "sidebar-tab-content--active": currentTab.value === "files" }]),
            onPreviewFile: handlePreviewFile,
            onBuildContext: handleBuildContext
          }, null, 8, ["class"]), [
            [vShow, currentTab.value === "files"]
          ]),
          withDirectives(createVNode(_sfc_main$l, {
            class: normalizeClass(["sidebar-tab-content", { "sidebar-tab-content--active": currentTab.value === "git" }])
          }, null, 8, ["class"]), [
            [vShow, currentTab.value === "git"]
          ]),
          withDirectives(createVNode(ContextList, {
            class: normalizeClass(["sidebar-tab-content", { "sidebar-tab-content--active": currentTab.value === "contexts" }])
          }, null, 8, ["class"]), [
            [vShow, currentTab.value === "contexts"]
          ])
        ])
      ]);
    };
  }
});
const LeftSidebar = /* @__PURE__ */ _export_sfc(_sfc_main$k, [["__scopeId", "data-v-84aedff9"]]);
function useMentions() {
  const projectStore = useProjectStore();
  const fileStore = useFileStore();
  const contextStore = useContextStore();
  const uiStore = useUIStore();
  const isProcessing = ref(false);
  async function processFilesMention(t) {
    const selectedFiles = fileStore.selectedFilesList;
    if (selectedFiles.length === 0) {
      uiStore.addToast(t("chat.selectFilesHint"), "warning");
      return {
        type: "files",
        success: false,
        message: t("chat.selectFilesHint")
      };
    }
    try {
      isProcessing.value = true;
      await contextStore.buildContext(selectedFiles);
      const message = t("chat.contextAttached");
      uiStore.addToast(message, "success");
      return {
        type: "files",
        success: true,
        message,
        filesAdded: selectedFiles.length
      };
    } catch (error) {
      const message = t("chat.contextBuildFailed");
      uiStore.addToast(message, "error");
      return {
        type: "files",
        success: false,
        message
      };
    } finally {
      isProcessing.value = false;
    }
  }
  async function processGitMention(t) {
    const projectPath = projectStore.currentPath;
    if (!projectPath) {
      uiStore.addToast(t("chat.noApiKey"), "warning");
      return {
        type: "git",
        success: false,
        message: "No project selected"
      };
    }
    try {
      isProcessing.value = true;
      const uncommittedFiles = await gitApi$1.getUncommittedFiles(projectPath);
      if (!uncommittedFiles || uncommittedFiles.length === 0) {
        const message2 = t("gitContext.noChanges");
        uiStore.addToast(message2, "info");
        return {
          type: "git",
          success: true,
          message: message2,
          filesAdded: 0
        };
      }
      const filePaths = uncommittedFiles.map((f) => f.path);
      await contextStore.buildContext(filePaths);
      const message = `${t("chat.contextAttached")} (${filePaths.length} git files)`;
      uiStore.addToast(message, "success");
      return {
        type: "git",
        success: true,
        message,
        filesAdded: filePaths.length
      };
    } catch (error) {
      const message = t("chat.contextBuildFailed");
      uiStore.addToast(message, "error");
      return {
        type: "git",
        success: false,
        message
      };
    } finally {
      isProcessing.value = false;
    }
  }
  async function processProblemsMention(t) {
    uiStore.addToast(t("chat.comingSoon"), "info");
    return {
      type: "problems",
      success: false,
      message: t("chat.comingSoon")
    };
  }
  async function processMention(mentionType, t) {
    switch (mentionType) {
      case "files":
        return processFilesMention(t);
      case "git":
        return processGitMention(t);
      case "problems":
        return processProblemsMention(t);
      default:
        return {
          type: mentionType,
          success: false,
          message: "Unknown mention type"
        };
    }
  }
  async function processMessageMentions(message, t) {
    const results = [];
    let cleanedMessage = message;
    if (message.includes("@files")) {
      const result = await processMention("files", t);
      results.push(result);
      cleanedMessage = cleanedMessage.replace(/@files\s*/g, "");
    }
    if (message.includes("@git")) {
      const result = await processMention("git", t);
      results.push(result);
      cleanedMessage = cleanedMessage.replace(/@git\s*/g, "");
    }
    if (message.includes("@problems")) {
      const result = await processMention("problems", t);
      results.push(result);
      cleanedMessage = cleanedMessage.replace(/@problems\s*/g, "");
    }
    return {
      cleanedMessage: cleanedMessage.trim(),
      results
    };
  }
  return {
    isProcessing,
    processMention,
    processMessageMentions,
    processFilesMention,
    processGitMention,
    processProblemsMention
  };
}
function useChatMessages() {
  const uiStore = useUIStore();
  const projectStore = useProjectStore();
  const settingsStore = useSettingsStore();
  const contextStore = useContextStore();
  const mentions = useMentions();
  const messages = ref([]);
  const inputMessage = ref("");
  const isThinking = ref(false);
  const isAnalyzing = ref(false);
  const smartContextPreview = ref(null);
  const expandedToolCalls = ref(/* @__PURE__ */ new Set());
  const abortController = ref(null);
  const currentModel = computed(() => settingsStore.settings.aiModel || "gpt-4");
  const providerName = computed(() => "OpenAI");
  const isConnected = computed(() => true);
  const totalUsedTokens = computed(() => {
    return messages.value.reduce((acc, msg) => acc + (msg.content?.length || 0) / 4, 0);
  });
  async function sendSmartMessage(scrollToBottom, t) {
    if (!inputMessage.value.trim() || isThinking.value) return;
    let content = inputMessage.value.trim();
    if (content.includes("@")) {
      const { cleanedMessage } = await mentions.processMessageMentions(content, t);
      content = cleanedMessage;
      if (!content) {
        inputMessage.value = "";
        return;
      }
    }
    const hasExistingContext = contextStore.hasContext;
    if (!hasExistingContext) {
      inputMessage.value = content;
      await analyzeAndSuggestFiles(content, t);
      return;
    }
    inputMessage.value = content;
    await executeChat(content, scrollToBottom);
  }
  async function analyzeAndSuggestFiles(query, t) {
    isAnalyzing.value = true;
    try {
      const result = await apiService.semanticSearch({
        query,
        projectRoot: projectStore.currentPath || "",
        topK: 10
      });
      if (result.results && result.results.length > 0) {
        const files2 = result.results.map((r) => ({
          path: r.chunk.filePath,
          reason: r.chunk.content?.slice(0, 100) || "",
          relevance: r.score || 0.5
        }));
        const estimatedTokens = files2.length * 500;
        smartContextPreview.value = {
          files: files2,
          totalTokens: estimatedTokens,
          query
        };
      } else {
        uiStore.addToast(t("chat.noFilesFound"), "warning");
        inputMessage.value = query;
      }
    } catch {
      uiStore.addToast(t("chat.analysisFailed"), "error");
    } finally {
      isAnalyzing.value = false;
    }
  }
  async function executeChat(content, scrollToBottom) {
    inputMessage.value = "";
    const userMessage = {
      id: `msg-${Date.now()}`,
      role: "user",
      content,
      timestamp: (/* @__PURE__ */ new Date()).toISOString(),
      contextAttached: contextStore.hasContext
    };
    messages.value.push(userMessage);
    scrollToBottom();
    isThinking.value = true;
    abortController.value = new AbortController();
    try {
      const projectRoot = projectStore.currentPath || "";
      const response = await apiService.agenticChat(content, projectRoot);
      const assistantMessage = {
        id: `msg-${Date.now()}-ai`,
        role: "assistant",
        content: response.response,
        timestamp: (/* @__PURE__ */ new Date()).toISOString(),
        toolCalls: response.toolCalls
      };
      messages.value.push(assistantMessage);
      scrollToBottom();
    } catch (error) {
      if (error.name !== "AbortError") {
        messages.value.push({
          id: `msg-${Date.now()}-error`,
          role: "assistant",
          content: "Error: Failed to get response",
          timestamp: (/* @__PURE__ */ new Date()).toISOString(),
          error: error.message
        });
      }
    } finally {
      isThinking.value = false;
      abortController.value = null;
    }
  }
  async function confirmSmartContext(selectedFiles, t, scrollToBottom) {
    if (!smartContextPreview.value) return;
    const query = smartContextPreview.value.query;
    smartContextPreview.value = null;
    try {
      await contextStore.buildContext(selectedFiles);
      await executeChat(query, scrollToBottom);
    } catch {
      uiStore.addToast(t("chat.contextBuildFailed"), "error");
    }
  }
  function cancelSmartContext() {
    smartContextPreview.value = null;
  }
  function clearChat() {
    messages.value = [];
    expandedToolCalls.value.clear();
  }
  function toggleToolCalls(index) {
    if (expandedToolCalls.value.has(index)) {
      expandedToolCalls.value.delete(index);
    } else {
      expandedToolCalls.value.add(index);
    }
    expandedToolCalls.value = new Set(expandedToolCalls.value);
  }
  function stopGeneration() {
    if (abortController.value) {
      abortController.value.abort();
      abortController.value = null;
    }
    isThinking.value = false;
  }
  function copyMessage(content, t) {
    navigator.clipboard.writeText(content);
    uiStore.addToast(t("common.copied"), "success");
  }
  function initialize() {
  }
  function cleanup() {
    stopGeneration();
  }
  return {
    // State
    messages,
    inputMessage,
    isThinking,
    isAnalyzing,
    smartContextPreview,
    expandedToolCalls,
    // Computed
    currentModel,
    providerName,
    isConnected,
    totalUsedTokens,
    // Actions
    sendSmartMessage,
    confirmSmartContext,
    cancelSmartContext,
    clearChat,
    toggleToolCalls,
    stopGeneration,
    copyMessage,
    initialize,
    cleanup
  };
}
const _hoisted_1$i = { class: "flex items-start gap-2" };
const _hoisted_2$i = {
  key: 0,
  class: "message-avatar"
};
const _hoisted_3$h = { class: "flex-1 min-w-0" };
const _hoisted_4$f = { class: "text-xs text-gray-300 whitespace-pre-wrap break-words" };
const _hoisted_5$d = {
  key: 0,
  class: "mt-2 flex items-center gap-1 text-[10px] text-indigo-400"
};
const _hoisted_6$b = {
  key: 1,
  class: "mt-2"
};
const _hoisted_7$b = {
  key: 0,
  class: "mt-2 space-y-1 text-[10px]"
};
const _hoisted_8$9 = { class: "font-medium text-emerald-400" };
const _hoisted_9$7 = ["title"];
const _hoisted_10$7 = {
  key: 2,
  class: "mt-2 text-[10px] text-red-400"
};
const _hoisted_11$7 = ["title"];
const _sfc_main$j = /* @__PURE__ */ defineComponent({
  __name: "ChatMessageItem",
  props: {
    message: {},
    index: {},
    expandedToolCalls: {}
  },
  emits: ["copy", "toggle-tools"],
  setup(__props) {
    const { t } = useI18n();
    return (_ctx, _cache) => {
      return openBlock(), createElementBlock("div", {
        class: normalizeClass(["group", [__props.message.role === "user" ? "message-user" : "message-assistant"]])
      }, [
        createBaseVNode("div", _hoisted_1$i, [
          __props.message.role === "assistant" ? (openBlock(), createElementBlock("div", _hoisted_2$i, [..._cache[2] || (_cache[2] = [
            createBaseVNode("svg", {
              class: "w-3.5 h-3.5 text-purple-300",
              fill: "none",
              stroke: "currentColor",
              viewBox: "0 0 24 24"
            }, [
              createBaseVNode("path", {
                "stroke-linecap": "round",
                "stroke-linejoin": "round",
                "stroke-width": "2",
                d: "M9.663 17h4.673M12 3v1m6.364 1.636l-.707.707M21 12h-1M4 12H3m3.343-5.657l-.707-.707m2.828 9.9a5 5 0 117.072 0l-.548.547A3.374 3.374 0 0014 18.469V19a2 2 0 11-4 0v-.531c0-.895-.356-1.754-.988-2.386l-.548-.547z"
              })
            ], -1)
          ])])) : createCommentVNode("", true),
          createBaseVNode("div", _hoisted_3$h, [
            createBaseVNode("p", _hoisted_4$f, toDisplayString(__props.message.content), 1),
            __props.message.contextAttached ? (openBlock(), createElementBlock("div", _hoisted_5$d, [
              _cache[3] || (_cache[3] = createBaseVNode("svg", {
                class: "w-3 h-3",
                fill: "none",
                stroke: "currentColor",
                viewBox: "0 0 24 24"
              }, [
                createBaseVNode("path", {
                  "stroke-linecap": "round",
                  "stroke-linejoin": "round",
                  "stroke-width": "2",
                  d: "M13.828 10.172a4 4 0 00-5.656 0l-4 4a4 4 0 105.656 5.656l1.102-1.101m-.758-4.899a4 4 0 005.656 0l4-4a4 4 0 00-5.656-5.656l-1.1 1.1"
                })
              ], -1)),
              createTextVNode(" " + toDisplayString(unref(t)("chat.contextAttached")) + " (" + toDisplayString(__props.message.tokenCount || 0) + " tokens) ", 1)
            ])) : createCommentVNode("", true),
            __props.message.toolCalls && __props.message.toolCalls.length > 0 ? (openBlock(), createElementBlock("div", _hoisted_6$b, [
              createBaseVNode("button", {
                onClick: _cache[0] || (_cache[0] = ($event) => _ctx.$emit("toggle-tools", __props.index)),
                class: "flex items-center gap-1 text-[10px] text-emerald-400 hover:text-emerald-300"
              }, [
                _cache[5] || (_cache[5] = createBaseVNode("svg", {
                  class: "w-3 h-3",
                  fill: "none",
                  stroke: "currentColor",
                  viewBox: "0 0 24 24"
                }, [
                  createBaseVNode("path", {
                    "stroke-linecap": "round",
                    "stroke-linejoin": "round",
                    "stroke-width": "2",
                    d: "M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z"
                  }),
                  createBaseVNode("path", {
                    "stroke-linecap": "round",
                    "stroke-linejoin": "round",
                    "stroke-width": "2",
                    d: "M15 12a3 3 0 11-6 0 3 3 0 016 0z"
                  })
                ], -1)),
                createTextVNode(" 🔧 " + toDisplayString(__props.message.toolCalls.length) + " tools (" + toDisplayString(__props.message.iterations) + " iter) ", 1),
                (openBlock(), createElementBlock("svg", {
                  class: normalizeClass(["w-3 h-3 transition-transform", __props.expandedToolCalls.has(__props.index) ? "rotate-180" : ""]),
                  fill: "none",
                  stroke: "currentColor",
                  viewBox: "0 0 24 24"
                }, [..._cache[4] || (_cache[4] = [
                  createBaseVNode("path", {
                    "stroke-linecap": "round",
                    "stroke-linejoin": "round",
                    "stroke-width": "2",
                    d: "M19 9l-7 7-7-7"
                  }, null, -1)
                ])], 2))
              ]),
              __props.expandedToolCalls.has(__props.index) ? (openBlock(), createElementBlock("div", _hoisted_7$b, [
                (openBlock(true), createElementBlock(Fragment, null, renderList(__props.message.toolCalls, (tc, tcIdx) => {
                  return openBlock(), createElementBlock("div", {
                    key: tcIdx,
                    class: "p-1.5 bg-gray-900/50 rounded border border-gray-700/30"
                  }, [
                    createBaseVNode("div", _hoisted_8$9, toDisplayString(tc.tool), 1),
                    createBaseVNode("div", {
                      class: "text-gray-400 truncate",
                      title: tc.arguments
                    }, toDisplayString(tc.arguments), 9, _hoisted_9$7)
                  ]);
                }), 128))
              ])) : createCommentVNode("", true)
            ])) : createCommentVNode("", true),
            __props.message.error ? (openBlock(), createElementBlock("div", _hoisted_10$7, toDisplayString(__props.message.error), 1)) : createCommentVNode("", true)
          ]),
          __props.message.role === "assistant" ? (openBlock(), createElementBlock("button", {
            key: 1,
            onClick: _cache[1] || (_cache[1] = ($event) => _ctx.$emit("copy")),
            class: "icon-btn-sm opacity-0 group-hover:opacity-100 transition-opacity",
            title: unref(t)("chat.copy")
          }, [..._cache[6] || (_cache[6] = [
            createBaseVNode("svg", {
              class: "w-3.5 h-3.5",
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
            ], -1)
          ])], 8, _hoisted_11$7)) : createCommentVNode("", true)
        ])
      ], 2);
    };
  }
});
const _hoisted_1$h = { class: "chat-welcome" };
const _hoisted_2$h = { class: "welcome-hero" };
const _hoisted_3$g = { class: "welcome-subtitle" };
const _hoisted_4$e = { key: 0 };
const _hoisted_5$c = { key: 1 };
const _hoisted_6$a = { class: "welcome-tips" };
const _hoisted_7$a = { class: "welcome-tip" };
const _hoisted_8$8 = { class: "welcome-tip" };
const _hoisted_9$6 = { class: "welcome-starters" };
const _hoisted_10$6 = { class: "welcome-starters__label" };
const _hoisted_11$6 = { class: "starter-grid" };
const _hoisted_12$5 = ["onClick"];
const _hoisted_13$5 = { class: "starter-card__icon" };
const _hoisted_14$5 = { class: "starter-card__text" };
const _sfc_main$i = /* @__PURE__ */ defineComponent({
  __name: "ChatWelcome",
  props: {
    isConnected: { type: Boolean },
    providerName: {},
    currentModel: {}
  },
  emits: ["quick-action"],
  setup(__props) {
    const { t } = useI18n();
    const promptStarters = computed(() => [
      {
        id: "analyze",
        icon: "🔍",
        label: t("chat.starters.analyze"),
        prompt: t("chat.prompts.analyze")
      },
      {
        id: "explain",
        icon: "💡",
        label: t("chat.starters.explain"),
        prompt: t("chat.prompts.explain")
      },
      {
        id: "refactor",
        icon: "✨",
        label: t("chat.starters.refactor"),
        prompt: t("chat.prompts.refactor")
      },
      {
        id: "test",
        icon: "🧪",
        label: t("chat.starters.test"),
        prompt: t("chat.prompts.test")
      }
    ]);
    return (_ctx, _cache) => {
      return openBlock(), createElementBlock("div", _hoisted_1$h, [
        _cache[5] || (_cache[5] = createBaseVNode("div", { class: "ambient-glow" }, null, -1)),
        createBaseVNode("div", _hoisted_2$h, [
          _cache[0] || (_cache[0] = createStaticVNode('<div class="welcome-icon" data-v-18267af2><div class="welcome-icon__glow" data-v-18267af2></div><svg class="w-10 h-10" fill="none" stroke="currentColor" viewBox="0 0 24 24" data-v-18267af2><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M13 10V3L4 14h7v7l9-11h-7z" data-v-18267af2></path></svg></div><h2 class="welcome-title" data-v-18267af2>Syntaxia AI</h2>', 2)),
          createBaseVNode("p", _hoisted_3$g, toDisplayString(unref(t)("chat.welcome.subtitle")), 1)
        ]),
        createBaseVNode("div", {
          class: normalizeClass(["welcome-status", { "welcome-status--connected": __props.isConnected }])
        }, [
          _cache[1] || (_cache[1] = createBaseVNode("span", { class: "welcome-status__dot" }, null, -1)),
          __props.isConnected ? (openBlock(), createElementBlock("span", _hoisted_4$e, toDisplayString(__props.providerName) + " · " + toDisplayString(__props.currentModel), 1)) : (openBlock(), createElementBlock("span", _hoisted_5$c, toDisplayString(unref(t)("chat.welcome.disconnected")), 1))
        ], 2),
        createBaseVNode("div", _hoisted_6$a, [
          createBaseVNode("div", _hoisted_7$a, [
            _cache[2] || (_cache[2] = createBaseVNode("kbd", null, "@", -1)),
            createBaseVNode("span", null, toDisplayString(unref(t)("chat.welcome.tipMention")), 1)
          ]),
          createBaseVNode("div", _hoisted_8$8, [
            _cache[3] || (_cache[3] = createBaseVNode("kbd", null, "/", -1)),
            createBaseVNode("span", null, toDisplayString(unref(t)("chat.welcome.tipCommand")), 1)
          ])
        ]),
        createBaseVNode("div", _hoisted_9$6, [
          createBaseVNode("p", _hoisted_10$6, toDisplayString(unref(t)("chat.welcome.tryAsking")), 1),
          createBaseVNode("div", _hoisted_11$6, [
            (openBlock(true), createElementBlock(Fragment, null, renderList(promptStarters.value, (starter) => {
              return openBlock(), createElementBlock("button", {
                key: starter.id,
                class: "starter-card",
                onClick: ($event) => _ctx.$emit("quick-action", { prompt: starter.prompt })
              }, [
                createBaseVNode("span", _hoisted_13$5, toDisplayString(starter.icon), 1),
                createBaseVNode("span", _hoisted_14$5, toDisplayString(starter.label), 1),
                _cache[4] || (_cache[4] = createBaseVNode("svg", {
                  class: "starter-card__arrow",
                  fill: "none",
                  stroke: "currentColor",
                  viewBox: "0 0 24 24"
                }, [
                  createBaseVNode("path", {
                    "stroke-linecap": "round",
                    "stroke-linejoin": "round",
                    "stroke-width": "2",
                    d: "M9 5l7 7-7 7"
                  })
                ], -1))
              ], 8, _hoisted_12$5);
            }), 128))
          ])
        ])
      ]);
    };
  }
});
const ChatWelcome = /* @__PURE__ */ _export_sfc(_sfc_main$i, [["__scopeId", "data-v-18267af2"]]);
const _hoisted_1$g = { class: "command-footer" };
const _hoisted_2$g = {
  key: 0,
  class: "context-chips"
};
const _hoisted_3$f = { class: "chip chip--tokens" };
const _hoisted_4$d = {
  key: 0,
  class: "context-list"
};
const _hoisted_5$b = { class: "context-item__name" };
const _hoisted_6$9 = { class: "context-item__path" };
const _hoisted_7$9 = ["onClick"];
const _hoisted_8$7 = { class: "capsule-wrapper" };
const _hoisted_9$5 = {
  key: 0,
  class: "mention-popup"
};
const _hoisted_10$5 = { class: "mention-header" };
const _hoisted_11$5 = ["onClick", "onMouseenter"];
const _hoisted_12$4 = {
  key: 0,
  class: "w-4 h-4",
  fill: "none",
  stroke: "currentColor",
  viewBox: "0 0 24 24"
};
const _hoisted_13$4 = {
  key: 1,
  class: "w-4 h-4",
  fill: "none",
  stroke: "currentColor",
  viewBox: "0 0 24 24"
};
const _hoisted_14$4 = {
  key: 2,
  class: "w-4 h-4",
  fill: "none",
  stroke: "currentColor",
  viewBox: "0 0 24 24"
};
const _hoisted_15$4 = { class: "mention-content" };
const _hoisted_16$4 = { class: "mention-label" };
const _hoisted_17$4 = { class: "mention-hint" };
const _hoisted_18$4 = ["value", "placeholder", "disabled"];
const _hoisted_19$4 = { class: "capsule-toolbar" };
const _hoisted_20$4 = { class: "toolbar-left" };
const _hoisted_21$2 = ["title"];
const _hoisted_22$2 = ["title"];
const _hoisted_23$2 = ["disabled"];
const _hoisted_24$2 = {
  key: 0,
  class: "w-4 h-4",
  fill: "none",
  stroke: "currentColor",
  viewBox: "0 0 24 24"
};
const _hoisted_25$2 = {
  key: 1,
  class: "w-4 h-4 animate-spin",
  fill: "none",
  viewBox: "0 0 24 24"
};
const _hoisted_26$2 = { class: "hints" };
const MAX_VISIBLE_FILES = 5;
const _sfc_main$h = /* @__PURE__ */ defineComponent({
  __name: "CommandCenter",
  props: {
    modelValue: {},
    isThinking: { type: Boolean },
    isAnalyzing: { type: Boolean },
    hasMessages: { type: Boolean }
  },
  emits: ["update:modelValue", "send", "clear", "attach"],
  setup(__props, { emit: __emit }) {
    const props = __props;
    const emit = __emit;
    const { t } = useI18n();
    const contextStore = useContextStore();
    const fileStore = useFileStore();
    const { processFilesMention, processGitMention, processProblemsMention } = useMentions();
    const inputRef = ref(null);
    const isFocused = ref(false);
    const isContextExpanded = ref(false);
    const showAllFiles = ref(false);
    const showMentionPopup = ref(false);
    const selectedMentionIndex = ref(0);
    const hasText = computed(() => props.modelValue.trim().length > 0);
    const isDisabled = computed(() => props.isThinking || props.isAnalyzing);
    const hasContext = computed(() => contextStore.hasContext);
    const fileCount = computed(() => contextStore.fileCount);
    const files2 = computed(() => contextStore.summary?.files || []);
    const selectedFilesCount = computed(() => fileStore.selectedCount || 0);
    const totalFileCount = computed(() => hasContext.value ? fileCount.value : selectedFilesCount.value);
    const hasAnySelection = computed(() => hasContext.value || selectedFilesCount.value > 0);
    const formattedTokens = computed(() => {
      const tokens = contextStore.totalTokens;
      if (tokens >= 1e3) return `${(tokens / 1e3).toFixed(1)}k`;
      return tokens.toString();
    });
    const displayedFiles = computed(
      () => showAllFiles.value ? files2.value : files2.value.slice(0, MAX_VISIBLE_FILES)
    );
    const hasMoreFiles = computed(() => files2.value.length > MAX_VISIBLE_FILES);
    const hiddenCount = computed(() => files2.value.length - MAX_VISIBLE_FILES);
    const mentionItems = [
      { id: "files", label: "@files", hint: t("chat.mentions.files") },
      { id: "git", label: "@git", hint: t("chat.mentions.git") },
      { id: "problems", label: "@problems", hint: t("chat.mentions.problems") }
    ];
    function handleInput(event) {
      const textarea = event.target;
      emit("update:modelValue", textarea.value);
      autoResize(textarea);
      const lastChar = textarea.value.slice(-1);
      if (lastChar === "@") {
        showMentionPopup.value = true;
        selectedMentionIndex.value = 0;
      } else if (showMentionPopup.value && !textarea.value.includes("@")) {
        showMentionPopup.value = false;
      }
    }
    function handleKeydown(event) {
      if (showMentionPopup.value) {
        if (event.key === "ArrowDown") {
          event.preventDefault();
          selectedMentionIndex.value = (selectedMentionIndex.value + 1) % mentionItems.length;
        } else if (event.key === "ArrowUp") {
          event.preventDefault();
          selectedMentionIndex.value = (selectedMentionIndex.value - 1 + mentionItems.length) % mentionItems.length;
        } else if (event.key === "Enter" || event.key === "Tab") {
          event.preventDefault();
          selectMention(mentionItems[selectedMentionIndex.value]);
        } else if (event.key === "Escape") {
          showMentionPopup.value = false;
        }
        return;
      }
      if (event.key === "Enter" && !event.shiftKey) {
        event.preventDefault();
        handleSend();
      }
    }
    function handleBlur() {
      setTimeout(() => {
        isFocused.value = false;
        showMentionPopup.value = false;
      }, 150);
    }
    function handleSend() {
      if (hasText.value && !isDisabled.value) {
        emit("send");
      }
    }
    async function selectMention(item) {
      showMentionPopup.value = false;
      switch (item.id) {
        case "files":
          await processFilesMention(t);
          break;
        case "git":
          await processGitMention(t);
          break;
        case "problems":
          await processProblemsMention(t);
          break;
      }
      nextTick(() => inputRef.value?.focus());
    }
    function autoResize(textarea) {
      textarea.style.height = "24px";
      textarea.style.height = Math.min(textarea.scrollHeight, 160) + "px";
    }
    function toggleContextExpand() {
      isContextExpanded.value = !isContextExpanded.value;
    }
    function clearContext() {
      contextStore.clearContext();
    }
    function removeFile(file) {
      contextStore.removeFileFromContext(file);
    }
    function getFileName(path) {
      return path.split("/").pop() || path;
    }
    function getFilePath(path) {
      const parts = path.split("/");
      return parts.length <= 1 ? "" : parts.slice(0, -1).join("/");
    }
    function handleMentionClick() {
      showMentionPopup.value = true;
      selectedMentionIndex.value = 0;
      nextTick(() => inputRef.value?.focus());
    }
    async function handleAttachClick() {
      await processFilesMention(t);
      nextTick(() => inputRef.value?.focus());
    }
    watch(() => props.isThinking, (thinking) => {
      if (!thinking) nextTick(() => inputRef.value?.focus());
    });
    return (_ctx, _cache) => {
      return openBlock(), createElementBlock("div", _hoisted_1$g, [
        createVNode(Transition, { name: "chips" }, {
          default: withCtx(() => [
            hasAnySelection.value ? (openBlock(), createElementBlock("div", _hoisted_2$g, [
              createBaseVNode("button", {
                class: "chip chip--files",
                onClick: toggleContextExpand
              }, [
                _cache[2] || (_cache[2] = createBaseVNode("svg", {
                  class: "w-3.5 h-3.5",
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
                createBaseVNode("span", null, toDisplayString(totalFileCount.value) + " " + toDisplayString(unref(t)("context.filesShort")), 1)
              ]),
              createBaseVNode("div", _hoisted_3$f, [
                createBaseVNode("span", null, toDisplayString(formattedTokens.value) + " " + toDisplayString(unref(t)("context.tokens")), 1)
              ]),
              hasContext.value ? (openBlock(), createElementBlock("button", {
                key: 0,
                class: "chip-clear",
                onClick: clearContext
              }, toDisplayString(unref(t)("context.clear")), 1)) : createCommentVNode("", true)
            ])) : createCommentVNode("", true)
          ]),
          _: 1
        }),
        createVNode(Transition, { name: "expand" }, {
          default: withCtx(() => [
            isContextExpanded.value && hasAnySelection.value ? (openBlock(), createElementBlock("div", _hoisted_4$d, [
              (openBlock(true), createElementBlock(Fragment, null, renderList(displayedFiles.value, (file) => {
                return openBlock(), createElementBlock("div", {
                  key: file,
                  class: "context-item"
                }, [
                  createBaseVNode("span", _hoisted_5$b, toDisplayString(getFileName(file)), 1),
                  createBaseVNode("span", _hoisted_6$9, toDisplayString(getFilePath(file)), 1),
                  createBaseVNode("button", {
                    class: "context-item__remove",
                    onClick: ($event) => removeFile(file)
                  }, "×", 8, _hoisted_7$9)
                ]);
              }), 128)),
              hasMoreFiles.value ? (openBlock(), createElementBlock("button", {
                key: 0,
                class: "context-more",
                onClick: _cache[0] || (_cache[0] = ($event) => showAllFiles.value = !showAllFiles.value)
              }, toDisplayString(showAllFiles.value ? unref(t)("common.showLess") : `+${hiddenCount.value}`), 1)) : createCommentVNode("", true)
            ])) : createCommentVNode("", true)
          ]),
          _: 1
        }),
        createBaseVNode("div", _hoisted_8$7, [
          createBaseVNode("div", {
            class: normalizeClass(["capsule-glow", { "capsule-glow--active": isFocused.value }])
          }, null, 2),
          createBaseVNode("div", {
            class: normalizeClass(["capsule", { "capsule--focused": isFocused.value }])
          }, [
            createVNode(Transition, { name: "popup" }, {
              default: withCtx(() => [
                showMentionPopup.value ? (openBlock(), createElementBlock("div", _hoisted_9$5, [
                  createBaseVNode("div", _hoisted_10$5, [
                    createBaseVNode("span", null, toDisplayString(unref(t)("chat.mentions.title")), 1),
                    _cache[3] || (_cache[3] = createBaseVNode("kbd", null, "↑↓", -1))
                  ]),
                  (openBlock(), createElementBlock(Fragment, null, renderList(mentionItems, (item, idx) => {
                    return createBaseVNode("button", {
                      key: item.id,
                      class: normalizeClass(["mention-item", { "mention-item--active": selectedMentionIndex.value === idx }]),
                      onClick: ($event) => selectMention(item),
                      onMouseenter: ($event) => selectedMentionIndex.value = idx
                    }, [
                      createBaseVNode("span", {
                        class: normalizeClass(["mention-icon", `mention-icon--${item.id}`])
                      }, [
                        item.id === "files" ? (openBlock(), createElementBlock("svg", _hoisted_12$4, [..._cache[4] || (_cache[4] = [
                          createBaseVNode("path", {
                            "stroke-linecap": "round",
                            "stroke-linejoin": "round",
                            "stroke-width": "1.5",
                            d: "M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-6l-2-2H5a2 2 0 00-2 2z"
                          }, null, -1)
                        ])])) : item.id === "git" ? (openBlock(), createElementBlock("svg", _hoisted_13$4, [..._cache[5] || (_cache[5] = [
                          createBaseVNode("path", {
                            "stroke-linecap": "round",
                            "stroke-linejoin": "round",
                            "stroke-width": "1.5",
                            d: "M8 9l3 3-3 3m5 0h3M5 20h14a2 2 0 002-2V6a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"
                          }, null, -1)
                        ])])) : (openBlock(), createElementBlock("svg", _hoisted_14$4, [..._cache[6] || (_cache[6] = [
                          createBaseVNode("path", {
                            "stroke-linecap": "round",
                            "stroke-linejoin": "round",
                            "stroke-width": "1.5",
                            d: "M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"
                          }, null, -1)
                        ])]))
                      ], 2),
                      createBaseVNode("div", _hoisted_15$4, [
                        createBaseVNode("span", _hoisted_16$4, toDisplayString(item.label), 1),
                        createBaseVNode("span", _hoisted_17$4, toDisplayString(item.hint), 1)
                      ])
                    ], 42, _hoisted_11$5);
                  }), 64))
                ])) : createCommentVNode("", true)
              ]),
              _: 1
            }),
            createBaseVNode("textarea", {
              ref_key: "inputRef",
              ref: inputRef,
              value: __props.modelValue,
              onInput: handleInput,
              onKeydown: handleKeydown,
              onFocus: _cache[1] || (_cache[1] = ($event) => isFocused.value = true),
              onBlur: handleBlur,
              class: "capsule-input",
              placeholder: unref(t)("chat.placeholder"),
              disabled: isDisabled.value,
              rows: "1"
            }, null, 40, _hoisted_18$4),
            createBaseVNode("div", _hoisted_19$4, [
              createBaseVNode("div", _hoisted_20$4, [
                createBaseVNode("button", {
                  class: "toolbar-btn",
                  onClick: handleAttachClick,
                  title: unref(t)("chat.attachFiles")
                }, [..._cache[7] || (_cache[7] = [
                  createBaseVNode("svg", {
                    class: "w-5 h-5",
                    fill: "none",
                    stroke: "currentColor",
                    viewBox: "0 0 24 24"
                  }, [
                    createBaseVNode("path", {
                      "stroke-linecap": "round",
                      "stroke-linejoin": "round",
                      "stroke-width": "1.5",
                      d: "M15.172 7l-6.586 6.586a2 2 0 102.828 2.828l6.414-6.586a4 4 0 00-5.656-5.656l-6.415 6.585a6 6 0 108.486 8.486L20.5 13"
                    })
                  ], -1)
                ])], 8, _hoisted_21$2),
                createBaseVNode("button", {
                  class: "toolbar-btn toolbar-btn--mention",
                  onClick: handleMentionClick,
                  title: unref(t)("chat.hints.mention")
                }, "@", 8, _hoisted_22$2)
              ]),
              createBaseVNode("button", {
                class: normalizeClass(["send-btn", { "send-btn--ready": hasText.value }]),
                onClick: handleSend,
                disabled: !hasText.value || isDisabled.value
              }, [
                !__props.isThinking ? (openBlock(), createElementBlock("svg", _hoisted_24$2, [..._cache[8] || (_cache[8] = [
                  createBaseVNode("path", {
                    "stroke-linecap": "round",
                    "stroke-linejoin": "round",
                    "stroke-width": "2",
                    d: "M5 10l7-7m0 0l7 7m-7-7v18"
                  }, null, -1)
                ])])) : (openBlock(), createElementBlock("svg", _hoisted_25$2, [..._cache[9] || (_cache[9] = [
                  createBaseVNode("circle", {
                    class: "opacity-25",
                    cx: "12",
                    cy: "12",
                    r: "10",
                    stroke: "currentColor",
                    "stroke-width": "3"
                  }, null, -1),
                  createBaseVNode("path", {
                    class: "opacity-75",
                    fill: "currentColor",
                    d: "M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z"
                  }, null, -1)
                ])]))
              ], 10, _hoisted_23$2)
            ])
          ], 2)
        ]),
        createBaseVNode("div", _hoisted_26$2, [
          createBaseVNode("span", null, "↵ " + toDisplayString(unref(t)("chat.send")), 1),
          createBaseVNode("span", null, "⇧↵ " + toDisplayString(unref(t)("chat.hints.newLine")), 1)
        ])
      ]);
    };
  }
});
const CommandCenter = /* @__PURE__ */ _export_sfc(_sfc_main$h, [["__scopeId", "data-v-50261d9a"]]);
const _hoisted_1$f = { class: "smart-context-panel" };
const _hoisted_2$f = { class: "smart-context-panel__header" };
const _hoisted_3$e = { class: "smart-context-panel__title" };
const _hoisted_4$c = { class: "smart-context-panel__stats" };
const _hoisted_5$a = { class: "smart-context-panel__files" };
const _hoisted_6$8 = { class: "suggested-file__main" };
const _hoisted_7$8 = ["checked", "onChange"];
const _hoisted_8$6 = ["title"];
const _hoisted_9$4 = {
  key: 0,
  class: "suggested-file__reason"
};
const _hoisted_10$4 = { class: "smart-context-panel__actions" };
const _hoisted_11$4 = ["disabled"];
const MAX_VISIBLE = 5;
const MAX_REASON_LENGTH = 80;
const _sfc_main$g = /* @__PURE__ */ defineComponent({
  __name: "SmartContextPreviewPanel",
  props: {
    preview: {}
  },
  emits: ["confirm", "cancel"],
  setup(__props, { emit: __emit }) {
    const props = __props;
    const emit = __emit;
    const { t } = useI18n();
    const showAll = ref(false);
    const selectedFiles = ref(/* @__PURE__ */ new Set());
    watch(() => props.preview.files, (files2) => {
      selectedFiles.value = new Set(files2.map((f) => f.path));
    }, { immediate: true });
    const displayedFiles = computed(() => {
      if (showAll.value) return props.preview.files;
      return props.preview.files.slice(0, MAX_VISIBLE);
    });
    const hasMoreFiles = computed(() => props.preview.files.length > MAX_VISIBLE);
    const hiddenCount = computed(() => props.preview.files.length - MAX_VISIBLE);
    const formattedTokens = computed(() => {
      const tokens = props.preview.totalTokens;
      if (tokens >= 1e3) return `${(tokens / 1e3).toFixed(1)}k`;
      return tokens.toString();
    });
    function getFileName(path) {
      return path.split("/").pop() || path;
    }
    function truncateReason(reason) {
      if (reason.length <= MAX_REASON_LENGTH) return reason;
      return reason.slice(0, MAX_REASON_LENGTH) + "...";
    }
    function toggleFile(path) {
      if (selectedFiles.value.has(path)) {
        selectedFiles.value.delete(path);
      } else {
        selectedFiles.value.add(path);
      }
      selectedFiles.value = new Set(selectedFiles.value);
    }
    function handleConfirm() {
      emit("confirm", Array.from(selectedFiles.value));
    }
    return (_ctx, _cache) => {
      return openBlock(), createElementBlock("div", _hoisted_1$f, [
        createBaseVNode("div", _hoisted_2$f, [
          createBaseVNode("div", _hoisted_3$e, [
            _cache[2] || (_cache[2] = createBaseVNode("svg", {
              class: "w-4 h-4 text-purple-400",
              fill: "none",
              stroke: "currentColor",
              viewBox: "0 0 24 24"
            }, [
              createBaseVNode("path", {
                "stroke-linecap": "round",
                "stroke-linejoin": "round",
                "stroke-width": "2",
                d: "M9.663 17h4.673M12 3v1m6.364 1.636l-.707.707M21 12h-1M4 12H3m3.343-5.657l-.707-.707m2.828 9.9a5 5 0 117.072 0l-.548.547A3.374 3.374 0 0014 18.469V19a2 2 0 11-4 0v-.531c0-.895-.356-1.754-.988-2.386l-.548-.547z"
              })
            ], -1)),
            createBaseVNode("span", null, toDisplayString(unref(t)("chat.suggestedFiles")), 1)
          ]),
          createBaseVNode("span", _hoisted_4$c, " ~" + toDisplayString(formattedTokens.value) + " tokens ", 1)
        ]),
        createBaseVNode("div", _hoisted_5$a, [
          (openBlock(true), createElementBlock(Fragment, null, renderList(displayedFiles.value, (file) => {
            return openBlock(), createElementBlock("div", {
              key: file.path,
              class: "suggested-file"
            }, [
              createBaseVNode("div", _hoisted_6$8, [
                createBaseVNode("input", {
                  type: "checkbox",
                  checked: selectedFiles.value.has(file.path),
                  onChange: ($event) => toggleFile(file.path),
                  class: "suggested-file__checkbox"
                }, null, 40, _hoisted_7$8),
                createBaseVNode("span", {
                  class: "suggested-file__name",
                  title: file.path
                }, toDisplayString(getFileName(file.path)), 9, _hoisted_8$6),
                createBaseVNode("span", {
                  class: "suggested-file__relevance",
                  style: normalizeStyle({ opacity: file.relevance })
                }, toDisplayString(Math.round(file.relevance * 100)) + "% ", 5)
              ]),
              file.reason ? (openBlock(), createElementBlock("div", _hoisted_9$4, toDisplayString(truncateReason(file.reason)), 1)) : createCommentVNode("", true)
            ]);
          }), 128)),
          hasMoreFiles.value ? (openBlock(), createElementBlock("button", {
            key: 0,
            class: "smart-context-panel__more",
            onClick: _cache[0] || (_cache[0] = ($event) => showAll.value = !showAll.value)
          }, toDisplayString(showAll.value ? unref(t)("common.showLess") : `+${hiddenCount.value} ${unref(t)("common.more")}`), 1)) : createCommentVNode("", true)
        ]),
        createBaseVNode("div", _hoisted_10$4, [
          createBaseVNode("button", {
            class: "btn btn-sm btn-primary flex-1",
            onClick: handleConfirm,
            disabled: selectedFiles.value.size === 0
          }, [
            _cache[3] || (_cache[3] = createBaseVNode("svg", {
              class: "w-3.5 h-3.5",
              fill: "none",
              stroke: "currentColor",
              viewBox: "0 0 24 24"
            }, [
              createBaseVNode("path", {
                "stroke-linecap": "round",
                "stroke-linejoin": "round",
                "stroke-width": "2",
                d: "M5 13l4 4L19 7"
              })
            ], -1)),
            createTextVNode(" " + toDisplayString(unref(t)("chat.useSelectedFiles", { count: selectedFiles.value.size })), 1)
          ], 8, _hoisted_11$4),
          createBaseVNode("button", {
            class: "btn btn-sm btn-ghost",
            onClick: _cache[1] || (_cache[1] = ($event) => _ctx.$emit("cancel"))
          }, toDisplayString(unref(t)("common.cancel")), 1)
        ])
      ]);
    };
  }
});
const SmartContextPreviewPanel = /* @__PURE__ */ _export_sfc(_sfc_main$g, [["__scopeId", "data-v-a98f8c85"]]);
const _hoisted_1$e = { class: "thinking-indicator" };
const _hoisted_2$e = { class: "flex items-center justify-between" };
const _hoisted_3$d = { class: "flex items-center gap-2" };
const _hoisted_4$b = { class: "text-xs text-purple-300" };
const _hoisted_5$9 = ["title"];
const _sfc_main$f = /* @__PURE__ */ defineComponent({
  __name: "ThinkingIndicator",
  emits: ["stop"],
  setup(__props) {
    const { t } = useI18n();
    return (_ctx, _cache) => {
      return openBlock(), createElementBlock("div", _hoisted_1$e, [
        createBaseVNode("div", _hoisted_2$e, [
          createBaseVNode("div", _hoisted_3$d, [
            _cache[1] || (_cache[1] = createBaseVNode("div", { class: "flex gap-1" }, [
              createBaseVNode("span", {
                class: "thinking-dot animate-bounce",
                style: { "animation-delay": "0ms" }
              }),
              createBaseVNode("span", {
                class: "thinking-dot animate-bounce",
                style: { "animation-delay": "150ms" }
              }),
              createBaseVNode("span", {
                class: "thinking-dot animate-bounce",
                style: { "animation-delay": "300ms" }
              })
            ], -1)),
            createBaseVNode("span", _hoisted_4$b, toDisplayString(unref(t)("chat.thinking")), 1)
          ]),
          createBaseVNode("button", {
            onClick: _cache[0] || (_cache[0] = ($event) => _ctx.$emit("stop")),
            class: "flex items-center gap-1 px-2 py-1 text-xs text-red-400 hover:text-red-300 hover:bg-red-500/10 rounded-lg transition-all",
            title: unref(t)("chat.stop")
          }, [
            _cache[2] || (_cache[2] = createBaseVNode("svg", {
              class: "w-3.5 h-3.5",
              fill: "currentColor",
              viewBox: "0 0 24 24"
            }, [
              createBaseVNode("rect", {
                x: "6",
                y: "6",
                width: "12",
                height: "12",
                rx: "2"
              })
            ], -1)),
            createTextVNode(" " + toDisplayString(unref(t)("chat.stop")), 1)
          ], 8, _hoisted_5$9)
        ])
      ]);
    };
  }
});
const logger$3 = useLogger("AIStore");
const MODEL_LIMITS = {
  // OpenAI
  "gpt-4o": 128e3,
  "gpt-4o-mini": 128e3,
  "gpt-4-turbo": 128e3,
  "gpt-4": 8192,
  "gpt-3.5-turbo": 16385,
  // Gemini
  "gemini-1.5-pro": 1e6,
  "gemini-1.5-flash": 1e6,
  "gemini-pro": 32e3,
  // Qwen
  "qwen-max": 32e3,
  "qwen-plus": 131072,
  "qwen-turbo": 131072,
  "qwen-long": 1e6,
  "qwen-coder-plus-latest": 131072,
  "qwen-coder-turbo-latest": 131072,
  "qwen-turbo-latest": 131072,
  // Claude (OpenRouter)
  "anthropic/claude-3.5-sonnet": 2e5,
  "anthropic/claude-3-opus": 2e5,
  "openai/gpt-4o": 128e3,
  "google/gemini-pro": 32e3,
  // LocalAI defaults
  "llama3": 8192,
  "mistral": 32e3,
  "codellama": 16384
};
const MODEL_PRICING = {
  "gpt-4o": 5e-3,
  "gpt-4o-mini": 15e-5,
  "gpt-4-turbo": 0.01,
  "gpt-4": 0.03,
  "gpt-3.5-turbo": 5e-4,
  "gemini-1.5-pro": 125e-5,
  "gemini-1.5-flash": 75e-6,
  "gemini-pro": 25e-5,
  "qwen-max": 4e-3,
  "qwen-plus": 4e-4,
  "qwen-turbo": 2e-4,
  "qwen-coder-plus-latest": 4e-4,
  "anthropic/claude-3.5-sonnet": 3e-3,
  "anthropic/claude-3-opus": 0.015
};
const useAIStore = defineStore("ai", () => {
  const providerInfo = ref({
    name: "Not configured",
    connected: false,
    model: "gpt-4o",
    provider: "openai"
  });
  const isLoading = ref(false);
  const lastError = ref(null);
  const currentModel = computed(() => providerInfo.value.model);
  const currentProvider = computed(() => providerInfo.value.provider);
  const isConnected = computed(() => providerInfo.value.connected);
  const contextLimit = computed(() => {
    return MODEL_LIMITS[currentModel.value] || 32e3;
  });
  const pricePerKTokens = computed(() => {
    return MODEL_PRICING[currentModel.value] || 1e-3;
  });
  async function loadProviderInfo() {
    isLoading.value = true;
    lastError.value = null;
    try {
      const infoJson = await apiService.getProviderInfo();
      const info = JSON.parse(infoJson);
      const settings2 = await apiService.getSettings();
      const selectedProvider = settings2.selectedProvider || "openai";
      const selectedModel = settings2.selectedModels?.[selectedProvider] || info.SupportedModels?.[0] || "gpt-4o";
      providerInfo.value = {
        name: info.Name || "AI Provider",
        connected: true,
        model: selectedModel,
        provider: selectedProvider
      };
      logger$3.debug("Provider info loaded:", providerInfo.value);
    } catch (e) {
      const errorMsg = e instanceof Error ? e.message : "Unknown error";
      logger$3.error("Failed to load provider info:", errorMsg);
      lastError.value = errorMsg;
      providerInfo.value = {
        name: "Not configured",
        connected: false,
        model: "gpt-4o",
        provider: "openai"
      };
    } finally {
      isLoading.value = false;
    }
  }
  function updateModel(provider, model) {
    providerInfo.value.provider = provider;
    providerInfo.value.model = model;
  }
  function calculateCost(tokens) {
    return tokens / 1e3 * pricePerKTokens.value;
  }
  function getUsagePercent(usedTokens) {
    const percent = Math.round(usedTokens / contextLimit.value * 100);
    return Math.min(percent, 100);
  }
  function formatTokens(tokens) {
    if (tokens >= 1e6) return (tokens / 1e6).toFixed(1) + "M";
    if (tokens >= 1e3) return (tokens / 1e3).toFixed(1) + "K";
    return tokens.toString();
  }
  return {
    // State
    providerInfo,
    isLoading,
    lastError,
    // Computed
    currentModel,
    currentProvider,
    isConnected,
    contextLimit,
    pricePerKTokens,
    // Actions
    loadProviderInfo,
    updateModel,
    calculateCost,
    getUsagePercent,
    formatTokens
  };
});
const _hoisted_1$d = { class: "model-limits-indicator" };
const _hoisted_2$d = { class: "limits-bar-container" };
const _hoisted_3$c = { class: "limits-info" };
const _hoisted_4$a = { class: "limits-used" };
const _hoisted_5$8 = { class: "limits-total" };
const _hoisted_6$7 = { class: "limits-label" };
const _hoisted_7$7 = { class: "limits-header" };
const _hoisted_8$5 = { class: "limits-model-info" };
const _hoisted_9$3 = { class: "limits-model-name" };
const _hoisted_10$3 = { class: "limits-provider-badge" };
const _hoisted_11$3 = { class: "limits-ring-container" };
const _hoisted_12$3 = {
  class: "limits-ring",
  viewBox: "0 0 100 100"
};
const _hoisted_13$3 = ["stroke-dashoffset"];
const _hoisted_14$3 = { class: "limits-ring-text" };
const _hoisted_15$3 = { class: "limits-ring-percent stats-counter" };
const _hoisted_16$3 = { class: "limits-ring-label" };
const _hoisted_17$3 = { class: "limits-stats" };
const _hoisted_18$3 = { class: "limits-stat stats-counter" };
const _hoisted_19$3 = { class: "limits-stat-value" };
const _hoisted_20$3 = { class: "limits-stat-label" };
const _hoisted_21$1 = {
  class: "limits-stat stats-counter",
  style: { "animation-delay": "0.1s" }
};
const _hoisted_22$1 = { class: "limits-stat-value" };
const _hoisted_23$1 = { class: "limits-stat-label" };
const _hoisted_24$1 = {
  class: "limits-stat stats-counter",
  style: { "animation-delay": "0.2s" }
};
const _hoisted_25$1 = { class: "limits-stat-value" };
const _hoisted_26$1 = { class: "limits-stat-label" };
const _hoisted_27$1 = {
  class: "limits-stat stats-counter",
  style: { "animation-delay": "0.3s" }
};
const _hoisted_28$1 = { class: "limits-stat-value limits-stat-cost" };
const _hoisted_29 = { class: "limits-stat-label" };
const _sfc_main$e = /* @__PURE__ */ defineComponent({
  __name: "ModelLimitsIndicator",
  props: {
    usedTokens: {},
    modelName: {},
    providerName: {}
  },
  setup(__props) {
    const { t } = useI18n();
    const props = __props;
    const expanded = ref(false);
    const contextLimit = computed(() => {
      return MODEL_LIMITS[props.modelName] || 32e3;
    });
    const usagePercent = computed(() => {
      if (!props.usedTokens || !contextLimit.value || contextLimit.value === 0) return 0;
      const percent = Math.round(props.usedTokens / contextLimit.value * 100);
      if (!isFinite(percent) || isNaN(percent)) return 0;
      return Math.min(Math.max(percent, 0), 100);
    });
    const availableTokens = computed(() => {
      return Math.max(0, contextLimit.value - props.usedTokens);
    });
    const usageClass = computed(() => {
      if (usagePercent.value >= 95) return "usage-critical";
      if (usagePercent.value >= 80) return "usage-warning";
      return "usage-safe";
    });
    const pulseClass = computed(() => {
      if (usagePercent.value >= 95) return "limits-critical-pulse";
      if (usagePercent.value >= 80) return "limits-warning-pulse";
      return "";
    });
    const glowClass = computed(() => {
      if (usagePercent.value >= 95) return "limits-glow-critical";
      if (usagePercent.value >= 80) return "limits-glow-warning";
      return "";
    });
    const estimatedCost = computed(() => {
      const pricePerK = MODEL_PRICING[props.modelName] || 1e-3;
      return (props.usedTokens / 1e3 * pricePerK).toFixed(4);
    });
    const circumference = 2 * Math.PI * 42;
    const dashOffset = computed(() => {
      return circumference - usagePercent.value / 100 * circumference;
    });
    function formatTokens(tokens) {
      if (tokens >= 1e6) return (tokens / 1e6).toFixed(1) + "M";
      if (tokens >= 1e3) return (tokens / 1e3).toFixed(1) + "K";
      return tokens.toString();
    }
    return (_ctx, _cache) => {
      return openBlock(), createElementBlock("div", _hoisted_1$d, [
        !expanded.value ? (openBlock(), createElementBlock("div", {
          key: 0,
          class: normalizeClass(["limits-compact", [pulseClass.value, glowClass.value]]),
          onClick: _cache[0] || (_cache[0] = ($event) => expanded.value = true)
        }, [
          createBaseVNode("div", _hoisted_2$d, [
            createBaseVNode("div", {
              class: normalizeClass(["limits-bar limits-bar-animate", usageClass.value]),
              style: normalizeStyle({ width: usagePercent.value + "%" })
            }, null, 6)
          ]),
          createBaseVNode("div", _hoisted_3$c, [
            createBaseVNode("span", _hoisted_4$a, toDisplayString(formatTokens(__props.usedTokens)), 1),
            _cache[2] || (_cache[2] = createBaseVNode("span", { class: "limits-separator" }, "/", -1)),
            createBaseVNode("span", _hoisted_5$8, toDisplayString(formatTokens(contextLimit.value)), 1),
            createBaseVNode("span", _hoisted_6$7, toDisplayString(unref(t)("limits.tokens")), 1),
            usagePercent.value >= 80 ? (openBlock(), createElementBlock("span", {
              key: 0,
              class: normalizeClass(["limits-percent", usageClass.value])
            }, toDisplayString(usagePercent.value) + "% ", 3)) : createCommentVNode("", true)
          ])
        ], 2)) : createCommentVNode("", true),
        createVNode(Transition, { name: "limits-expand" }, {
          default: withCtx(() => [
            expanded.value ? (openBlock(), createElementBlock("div", {
              key: 0,
              class: normalizeClass(["limits-expanded", glowClass.value])
            }, [
              createBaseVNode("div", _hoisted_7$7, [
                createBaseVNode("div", _hoisted_8$5, [
                  createBaseVNode("span", _hoisted_9$3, toDisplayString(__props.modelName), 1),
                  createBaseVNode("span", _hoisted_10$3, toDisplayString(__props.providerName), 1)
                ]),
                createBaseVNode("button", {
                  class: "icon-btn-sm",
                  onClick: _cache[1] || (_cache[1] = withModifiers(($event) => expanded.value = false, ["stop"]))
                }, [..._cache[3] || (_cache[3] = [
                  createBaseVNode("svg", {
                    class: "w-3.5 h-3.5",
                    fill: "none",
                    stroke: "currentColor",
                    viewBox: "0 0 24 24"
                  }, [
                    createBaseVNode("path", {
                      "stroke-linecap": "round",
                      "stroke-linejoin": "round",
                      "stroke-width": "2",
                      d: "M5 15l7-7 7 7"
                    })
                  ], -1)
                ])])
              ]),
              createBaseVNode("div", _hoisted_11$3, [
                (openBlock(), createElementBlock("svg", _hoisted_12$3, [
                  _cache[4] || (_cache[4] = createBaseVNode("circle", {
                    class: "limits-ring-bg",
                    cx: "50",
                    cy: "50",
                    r: "42",
                    fill: "none",
                    "stroke-width": "8"
                  }, null, -1)),
                  createBaseVNode("circle", {
                    class: normalizeClass(["limits-ring-progress limits-ring-animate", usageClass.value]),
                    cx: "50",
                    cy: "50",
                    r: "42",
                    fill: "none",
                    "stroke-width": "8",
                    "stroke-dasharray": circumference,
                    "stroke-dashoffset": dashOffset.value,
                    "stroke-linecap": "round"
                  }, null, 10, _hoisted_13$3)
                ])),
                createBaseVNode("div", _hoisted_14$3, [
                  createBaseVNode("span", _hoisted_15$3, toDisplayString(usagePercent.value) + "%", 1),
                  createBaseVNode("span", _hoisted_16$3, toDisplayString(unref(t)("limits.used")), 1)
                ])
              ]),
              createBaseVNode("div", _hoisted_17$3, [
                createBaseVNode("div", _hoisted_18$3, [
                  createBaseVNode("span", _hoisted_19$3, toDisplayString(formatTokens(__props.usedTokens)), 1),
                  createBaseVNode("span", _hoisted_20$3, toDisplayString(unref(t)("limits.contextUsed")), 1)
                ]),
                createBaseVNode("div", _hoisted_21$1, [
                  createBaseVNode("span", _hoisted_22$1, toDisplayString(formatTokens(availableTokens.value)), 1),
                  createBaseVNode("span", _hoisted_23$1, toDisplayString(unref(t)("limits.available")), 1)
                ]),
                createBaseVNode("div", _hoisted_24$1, [
                  createBaseVNode("span", _hoisted_25$1, toDisplayString(formatTokens(contextLimit.value)), 1),
                  createBaseVNode("span", _hoisted_26$1, toDisplayString(unref(t)("limits.maxContext")), 1)
                ]),
                createBaseVNode("div", _hoisted_27$1, [
                  createBaseVNode("span", _hoisted_28$1, "$" + toDisplayString(estimatedCost.value), 1),
                  createBaseVNode("span", _hoisted_29, toDisplayString(unref(t)("limits.estCost")), 1)
                ])
              ]),
              createVNode(Transition, { name: "slide-fade" }, {
                default: withCtx(() => [
                  usagePercent.value >= 80 ? (openBlock(), createElementBlock("div", {
                    key: 0,
                    class: normalizeClass(["limits-warning", [usageClass.value, pulseClass.value]])
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
                        d: "M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"
                      })
                    ], -1)),
                    createBaseVNode("span", null, toDisplayString(usagePercent.value >= 95 ? unref(t)("limits.criticalWarning") : unref(t)("limits.nearLimit")), 1)
                  ], 2)) : createCommentVNode("", true)
                ]),
                _: 1
              })
            ], 2)) : createCommentVNode("", true)
          ]),
          _: 1
        })
      ]);
    };
  }
});
const ModelLimitsIndicator = /* @__PURE__ */ _export_sfc(_sfc_main$e, [["__scopeId", "data-v-5517dde7"]]);
const _hoisted_1$c = { class: "ai-chat" };
const _hoisted_2$c = { class: "ai-chat__header" };
const _hoisted_3$b = {
  key: 0,
  class: "ai-chat__analyzing"
};
const _sfc_main$d = /* @__PURE__ */ defineComponent({
  __name: "AIChat",
  setup(__props) {
    const { t } = useI18n();
    const {
      messages,
      inputMessage,
      isThinking,
      isAnalyzing,
      smartContextPreview,
      expandedToolCalls,
      currentModel,
      providerName,
      isConnected,
      totalUsedTokens,
      sendSmartMessage,
      confirmSmartContext,
      cancelSmartContext,
      clearChat,
      toggleToolCalls,
      stopGeneration,
      copyMessage,
      cleanup
    } = useChatMessages();
    const messagesContainer = ref(null);
    function scrollToBottom() {
      nextTick(() => {
        if (messagesContainer.value) {
          messagesContainer.value.scrollTop = messagesContainer.value.scrollHeight;
        }
      });
    }
    async function handleSend() {
      await sendSmartMessage(scrollToBottom, t);
    }
    function handleQuickAction(action) {
      inputMessage.value = action.prompt;
      handleSend();
    }
    function handleAttach() {
    }
    onMounted(() => {
    });
    onUnmounted(() => {
      cleanup();
    });
    return (_ctx, _cache) => {
      return openBlock(), createElementBlock("div", _hoisted_1$c, [
        createBaseVNode("div", _hoisted_2$c, [
          unref(isConnected) ? (openBlock(), createBlock(ModelLimitsIndicator, {
            key: 0,
            "used-tokens": unref(totalUsedTokens),
            "model-name": unref(currentModel),
            "provider-name": unref(providerName)
          }, null, 8, ["used-tokens", "model-name", "provider-name"])) : createCommentVNode("", true)
        ]),
        createBaseVNode("div", {
          ref_key: "messagesContainer",
          ref: messagesContainer,
          class: "ai-chat__messages"
        }, [
          unref(messages).length === 0 ? (openBlock(), createBlock(ChatWelcome, {
            key: 0,
            "is-connected": unref(isConnected),
            "provider-name": unref(providerName),
            "current-model": unref(currentModel),
            onQuickAction: handleQuickAction
          }, null, 8, ["is-connected", "provider-name", "current-model"])) : createCommentVNode("", true),
          (openBlock(true), createElementBlock(Fragment, null, renderList(unref(messages), (msg, index) => {
            return openBlock(), createBlock(_sfc_main$j, {
              key: index,
              message: msg,
              index,
              "expanded-tool-calls": unref(expandedToolCalls),
              onCopy: ($event) => unref(copyMessage)(msg.content, unref(t)),
              onToggleTools: unref(toggleToolCalls)
            }, null, 8, ["message", "index", "expanded-tool-calls", "onCopy", "onToggleTools"]);
          }), 128)),
          unref(isThinking) ? (openBlock(), createBlock(_sfc_main$f, {
            key: 1,
            onStop: unref(stopGeneration)
          }, null, 8, ["onStop"])) : createCommentVNode("", true)
        ], 512),
        createVNode(Transition, { name: "slide-up" }, {
          default: withCtx(() => [
            unref(smartContextPreview) && unref(smartContextPreview).files.length > 0 ? (openBlock(), createBlock(SmartContextPreviewPanel, {
              key: 0,
              preview: unref(smartContextPreview),
              onConfirm: _cache[0] || (_cache[0] = (files2) => unref(confirmSmartContext)(files2, unref(t), scrollToBottom)),
              onCancel: unref(cancelSmartContext)
            }, null, 8, ["preview", "onCancel"])) : createCommentVNode("", true)
          ]),
          _: 1
        }),
        createVNode(Transition, { name: "fade" }, {
          default: withCtx(() => [
            unref(isAnalyzing) ? (openBlock(), createElementBlock("div", _hoisted_3$b, [
              _cache[2] || (_cache[2] = createBaseVNode("div", { class: "analyzing-pulse" }, null, -1)),
              _cache[3] || (_cache[3] = createBaseVNode("svg", {
                class: "w-4 h-4 text-purple-400 animate-spin",
                fill: "none",
                viewBox: "0 0 24 24"
              }, [
                createBaseVNode("circle", {
                  class: "opacity-25",
                  cx: "12",
                  cy: "12",
                  r: "10",
                  stroke: "currentColor",
                  "stroke-width": "3"
                }),
                createBaseVNode("path", {
                  class: "opacity-75",
                  fill: "currentColor",
                  d: "M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z"
                })
              ], -1)),
              createBaseVNode("span", null, toDisplayString(unref(t)("chat.analyzing")), 1)
            ])) : createCommentVNode("", true)
          ]),
          _: 1
        }),
        createVNode(CommandCenter, {
          modelValue: unref(inputMessage),
          "onUpdate:modelValue": _cache[1] || (_cache[1] = ($event) => isRef(inputMessage) ? inputMessage.value = $event : null),
          "is-thinking": unref(isThinking),
          "is-analyzing": unref(isAnalyzing),
          "has-messages": unref(messages).length > 0,
          onSend: handleSend,
          onClear: unref(clearChat),
          onAttach: handleAttach
        }, null, 8, ["modelValue", "is-thinking", "is-analyzing", "has-messages", "onClear"])
      ]);
    };
  }
});
const AIChat = /* @__PURE__ */ _export_sfc(_sfc_main$d, [["__scopeId", "data-v-89130d54"]]);
const logger$2 = useLogger("AISettings");
const DEFAULT_MODELS = {
  openai: ["gpt-4o", "gpt-4o-mini", "gpt-4-turbo", "gpt-4", "gpt-3.5-turbo"],
  gemini: ["gemini-1.5-pro", "gemini-1.5-flash", "gemini-pro"],
  qwen: ["qwen-max", "qwen-plus", "qwen-turbo", "qwen-long"],
  "qwen-cli": ["qwen-coder-plus-latest", "qwen-coder-turbo-latest", "qwen-turbo-latest"],
  openrouter: ["anthropic/claude-3.5-sonnet", "openai/gpt-4o", "google/gemini-pro"],
  localai: ["llama3", "mistral", "codellama"]
};
function useAISettings() {
  const { t } = useI18n();
  const uiStore = useUIStore();
  const aiStore = useAIStore();
  const settings2 = reactive({
    selectedProvider: "",
    openAIAPIKey: "",
    geminiAPIKey: "",
    qwenAPIKey: "",
    openRouterAPIKey: "",
    localAIAPIKey: "",
    localAIHost: "http://localhost:8080",
    qwenHost: "https://dashscope.aliyuncs.com/compatible-mode/v1",
    selectedModels: {},
    availableModels: {}
  });
  const showApiKey = ref(false);
  const isSaving = ref(false);
  const statusMessage = ref("");
  const statusType = ref("success");
  const providers = computed(() => [
    { id: "openai", name: "OpenAI", icon: "🤖", description: t("settings.provider.openai") },
    { id: "gemini", name: "Google Gemini", icon: "✨", description: t("settings.provider.gemini") },
    { id: "qwen", name: "Qwen (Alibaba)", icon: "🌐", description: t("settings.provider.qwen") },
    { id: "qwen-cli", name: "Qwen Code CLI", icon: "💻", description: t("settings.provider.qwenCli") },
    { id: "openrouter", name: "OpenRouter", icon: "🔀", description: t("settings.provider.openrouter") },
    { id: "localai", name: "LocalAI", icon: "🏠", description: t("settings.provider.localai") }
  ]);
  const currentProvider = computed(
    () => providers.value.find((p) => p.id === settings2.selectedProvider)
  );
  const currentProviderHint = computed(() => {
    const hints = {
      openai: t("settings.hint.openai"),
      gemini: t("settings.hint.gemini"),
      qwen: t("settings.hint.qwen"),
      openrouter: t("settings.hint.openrouter"),
      localai: t("settings.hint.localai")
    };
    return hints[settings2.selectedProvider] || "";
  });
  const currentApiKey = computed(() => {
    const keys = {
      openai: settings2.openAIAPIKey,
      gemini: settings2.geminiAPIKey,
      qwen: settings2.qwenAPIKey,
      openrouter: settings2.openRouterAPIKey,
      localai: settings2.localAIAPIKey
    };
    return keys[settings2.selectedProvider] || "";
  });
  const availableModels = computed(
    () => settings2.availableModels[settings2.selectedProvider] || DEFAULT_MODELS[settings2.selectedProvider] || []
  );
  const selectedModel = computed(
    () => settings2.selectedModels[settings2.selectedProvider] || availableModels.value[0] || ""
  );
  const needsApiKey = computed(
    () => settings2.selectedProvider && settings2.selectedProvider !== "qwen-cli"
  );
  const needsHostUrl = computed(
    () => settings2.selectedProvider === "localai" || settings2.selectedProvider === "qwen"
  );
  const currentHostUrl = computed(
    () => settings2.selectedProvider === "localai" ? settings2.localAIHost : settings2.qwenHost
  );
  const hostPlaceholder = computed(
    () => settings2.selectedProvider === "localai" ? "http://localhost:8080" : "https://dashscope.aliyuncs.com/compatible-mode/v1"
  );
  function selectProvider(providerId) {
    settings2.selectedProvider = providerId;
  }
  function updateApiKey(value) {
    switch (settings2.selectedProvider) {
      case "openai":
        settings2.openAIAPIKey = value;
        break;
      case "gemini":
        settings2.geminiAPIKey = value;
        break;
      case "qwen":
        settings2.qwenAPIKey = value;
        break;
      case "openrouter":
        settings2.openRouterAPIKey = value;
        break;
      case "localai":
        settings2.localAIAPIKey = value;
        break;
    }
  }
  function clearApiKey() {
    updateApiKey("");
  }
  function updateModel(model) {
    settings2.selectedModels[settings2.selectedProvider] = model;
  }
  function updateHost(value) {
    if (settings2.selectedProvider === "localai") {
      settings2.localAIHost = value;
    } else if (settings2.selectedProvider === "qwen") {
      settings2.qwenHost = value;
    }
  }
  function toggleShowApiKey() {
    showApiKey.value = !showApiKey.value;
  }
  async function loadSettings() {
    try {
      const dto = await apiService.getSettings();
      settings2.selectedProvider = dto.selectedProvider || "openai";
      settings2.openAIAPIKey = dto.openAIAPIKey || "";
      settings2.geminiAPIKey = dto.geminiAPIKey || "";
      settings2.qwenAPIKey = dto.qwenAPIKey || "";
      settings2.openRouterAPIKey = dto.openRouterAPIKey || "";
      settings2.localAIAPIKey = dto.localAIAPIKey || "";
      settings2.localAIHost = dto.localAIHost || "http://localhost:8080";
      settings2.qwenHost = dto.qwenHost || "https://dashscope.aliyuncs.com/compatible-mode/v1";
      settings2.selectedModels = dto.selectedModels || {};
      settings2.availableModels = dto.availableModels || {};
    } catch (e) {
      logger$2.error("Failed to load settings:", e);
    }
  }
  async function saveSettings() {
    isSaving.value = true;
    statusMessage.value = "";
    try {
      const dto = await apiService.getSettings();
      dto.selectedProvider = settings2.selectedProvider;
      dto.openAIAPIKey = settings2.openAIAPIKey;
      dto.geminiAPIKey = settings2.geminiAPIKey;
      dto.qwenAPIKey = settings2.qwenAPIKey;
      dto.openRouterAPIKey = settings2.openRouterAPIKey;
      dto.localAIAPIKey = settings2.localAIAPIKey;
      dto.localAIHost = settings2.localAIHost;
      dto.qwenHost = settings2.qwenHost;
      dto.selectedModels = settings2.selectedModels;
      await apiService.saveSettings(JSON.stringify(dto));
      const model = settings2.selectedModels[settings2.selectedProvider] || "";
      aiStore.updateModel(settings2.selectedProvider, model);
      statusMessage.value = t("settings.saved");
      statusType.value = "success";
      uiStore.addToast(t("settings.saved"), "success");
      setTimeout(() => {
        statusMessage.value = "";
      }, 3e3);
    } catch (e) {
      const errorMsg = e instanceof Error ? e.message : "Unknown error";
      statusMessage.value = errorMsg;
      statusType.value = "error";
      uiStore.addToast(t("settings.saveFailed"), "error");
    } finally {
      isSaving.value = false;
    }
  }
  return {
    // State
    settings: settings2,
    showApiKey,
    isSaving,
    statusMessage,
    statusType,
    // Computed
    providers,
    currentProvider,
    currentProviderHint,
    currentApiKey,
    availableModels,
    selectedModel,
    needsApiKey,
    needsHostUrl,
    currentHostUrl,
    hostPlaceholder,
    // Actions
    selectProvider,
    updateApiKey,
    clearApiKey,
    updateModel,
    updateHost,
    toggleShowApiKey,
    loadSettings,
    saveSettings
  };
}
const _hoisted_1$b = { class: "settings-field" };
const _hoisted_2$b = { class: "settings-label" };
const _hoisted_3$a = { class: "settings-input-wrap" };
const _hoisted_4$9 = ["type", "value", "placeholder"];
const _hoisted_5$7 = { class: "settings-input-actions" };
const _hoisted_6$6 = ["title", "aria-label"];
const _hoisted_7$6 = ["title", "aria-label"];
const _hoisted_8$4 = {
  key: 0,
  class: "settings-hint"
};
const _sfc_main$c = /* @__PURE__ */ defineComponent({
  __name: "ApiKeyInput",
  props: {
    modelValue: {},
    providerName: {},
    showKey: { type: Boolean },
    hint: {}
  },
  emits: ["update:modelValue", "toggle-visibility", "clear"],
  setup(__props) {
    const { t } = useI18n();
    return (_ctx, _cache) => {
      return openBlock(), createElementBlock("div", _hoisted_1$b, [
        createBaseVNode("label", _hoisted_2$b, toDisplayString(unref(t)("settings.apiKey")) + " (" + toDisplayString(__props.providerName) + ") ", 1),
        createBaseVNode("div", _hoisted_3$a, [
          createBaseVNode("input", {
            type: __props.showKey ? "text" : "password",
            value: __props.modelValue,
            onInput: _cache[0] || (_cache[0] = ($event) => _ctx.$emit("update:modelValue", $event.target.value)),
            class: "settings-input settings-input--password",
            placeholder: unref(t)("settings.apiKeyPlaceholder")
          }, null, 40, _hoisted_4$9),
          createBaseVNode("div", _hoisted_5$7, [
            createBaseVNode("button", {
              onClick: _cache[1] || (_cache[1] = ($event) => _ctx.$emit("toggle-visibility")),
              class: "settings-input-btn",
              title: __props.showKey ? unref(t)("settings.hideKey") : unref(t)("settings.showKey"),
              "aria-label": __props.showKey ? unref(t)("settings.hideKey") : unref(t)("settings.showKey")
            }, [
              __props.showKey ? (openBlock(), createBlock(unref(EyeOff), {
                key: 0,
                class: "w-4 h-4",
                "aria-hidden": "true"
              })) : (openBlock(), createBlock(unref(Eye), {
                key: 1,
                class: "w-4 h-4",
                "aria-hidden": "true"
              }))
            ], 8, _hoisted_6$6),
            __props.modelValue ? (openBlock(), createElementBlock("button", {
              key: 0,
              onClick: _cache[2] || (_cache[2] = ($event) => _ctx.$emit("clear")),
              class: "settings-input-btn settings-input-btn--danger",
              title: unref(t)("settings.clearKey"),
              "aria-label": unref(t)("settings.clearKey")
            }, [
              createVNode(unref(X), {
                class: "w-4 h-4",
                "aria-hidden": "true"
              })
            ], 8, _hoisted_7$6)) : createCommentVNode("", true)
          ])
        ]),
        __props.hint ? (openBlock(), createElementBlock("p", _hoisted_8$4, toDisplayString(__props.hint), 1)) : createCommentVNode("", true)
      ]);
    };
  }
});
const ApiKeyInput = /* @__PURE__ */ _export_sfc(_sfc_main$c, [["__scopeId", "data-v-42376c7c"]]);
const _hoisted_1$a = { class: "settings-field" };
const _hoisted_2$a = { class: "settings-label" };
const _hoisted_3$9 = ["value", "placeholder"];
const _sfc_main$b = /* @__PURE__ */ defineComponent({
  __name: "HostUrlInput",
  props: {
    modelValue: {},
    placeholder: {}
  },
  emits: ["update:modelValue"],
  setup(__props) {
    const { t } = useI18n();
    return (_ctx, _cache) => {
      return openBlock(), createElementBlock("div", _hoisted_1$a, [
        createBaseVNode("label", _hoisted_2$a, toDisplayString(unref(t)("settings.hostUrl")), 1),
        createBaseVNode("input", {
          type: "text",
          value: __props.modelValue,
          onInput: _cache[0] || (_cache[0] = ($event) => _ctx.$emit("update:modelValue", $event.target.value)),
          class: "settings-input",
          placeholder: __props.placeholder
        }, null, 40, _hoisted_3$9)
      ]);
    };
  }
});
const HostUrlInput = /* @__PURE__ */ _export_sfc(_sfc_main$b, [["__scopeId", "data-v-e16844b4"]]);
const _hoisted_1$9 = { class: "settings-field" };
const _hoisted_2$9 = { class: "settings-label" };
const _hoisted_3$8 = { class: "settings-select-wrap" };
const _hoisted_4$8 = ["value"];
const _hoisted_5$6 = ["value"];
const _sfc_main$a = /* @__PURE__ */ defineComponent({
  __name: "ModelSelector",
  props: {
    modelValue: {},
    models: {}
  },
  emits: ["update:modelValue"],
  setup(__props) {
    const { t } = useI18n();
    return (_ctx, _cache) => {
      return openBlock(), createElementBlock("div", _hoisted_1$9, [
        createBaseVNode("label", _hoisted_2$9, toDisplayString(unref(t)("settings.selectModel")), 1),
        createBaseVNode("div", _hoisted_3$8, [
          createBaseVNode("select", {
            value: __props.modelValue,
            onChange: _cache[0] || (_cache[0] = ($event) => _ctx.$emit("update:modelValue", $event.target.value)),
            class: "settings-select"
          }, [
            (openBlock(true), createElementBlock(Fragment, null, renderList(__props.models, (model) => {
              return openBlock(), createElementBlock("option", {
                key: model,
                value: model
              }, toDisplayString(model), 9, _hoisted_5$6);
            }), 128))
          ], 40, _hoisted_4$8),
          _cache[1] || (_cache[1] = createBaseVNode("div", { class: "settings-select-icon" }, [
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
                d: "M19 9l-7 7-7-7"
              })
            ])
          ], -1))
        ])
      ]);
    };
  }
});
const ModelSelector = /* @__PURE__ */ _export_sfc(_sfc_main$a, [["__scopeId", "data-v-5236620b"]]);
const _hoisted_1$8 = { class: "settings-field" };
const _hoisted_2$8 = { class: "settings-label" };
const _hoisted_3$7 = { class: "settings-select-wrap" };
const _hoisted_4$7 = ["value"];
const _hoisted_5$5 = ["value"];
const _hoisted_6$5 = { class: "settings-select-icon" };
const _hoisted_7$5 = {
  key: 0,
  class: "settings-hint"
};
const _sfc_main$9 = /* @__PURE__ */ defineComponent({
  __name: "ProviderSelector",
  props: {
    modelValue: {},
    providers: {},
    description: {}
  },
  emits: ["update:modelValue"],
  setup(__props) {
    const { t } = useI18n();
    return (_ctx, _cache) => {
      return openBlock(), createElementBlock("div", _hoisted_1$8, [
        createBaseVNode("label", _hoisted_2$8, toDisplayString(unref(t)("settings.selectProvider")), 1),
        createBaseVNode("div", _hoisted_3$7, [
          createBaseVNode("select", {
            value: __props.modelValue,
            onChange: _cache[0] || (_cache[0] = ($event) => _ctx.$emit("update:modelValue", $event.target.value)),
            class: "settings-select"
          }, [
            (openBlock(true), createElementBlock(Fragment, null, renderList(__props.providers, (provider) => {
              return openBlock(), createElementBlock("option", {
                key: provider.id,
                value: provider.id
              }, toDisplayString(provider.icon) + " " + toDisplayString(provider.name), 9, _hoisted_5$5);
            }), 128))
          ], 40, _hoisted_4$7),
          createBaseVNode("div", _hoisted_6$5, [
            createVNode(unref(ChevronDown), { class: "w-4 h-4" })
          ])
        ]),
        __props.description ? (openBlock(), createElementBlock("p", _hoisted_7$5, toDisplayString(__props.description), 1)) : createCommentVNode("", true)
      ]);
    };
  }
});
const ProviderSelector = /* @__PURE__ */ _export_sfc(_sfc_main$9, [["__scopeId", "data-v-16c97cf0"]]);
const _hoisted_1$7 = { class: "info-box-purple" };
const _hoisted_2$7 = { class: "flex items-start gap-2" };
const _hoisted_3$6 = { class: "text-xs text-purple-300 font-medium mb-1" };
const _hoisted_4$6 = { class: "text-2xs text-gray-400" };
const _sfc_main$8 = /* @__PURE__ */ defineComponent({
  __name: "QwenCliInfo",
  setup(__props) {
    const { t } = useI18n();
    return (_ctx, _cache) => {
      return openBlock(), createElementBlock("div", _hoisted_1$7, [
        createBaseVNode("div", _hoisted_2$7, [
          createVNode(unref(Info), { class: "w-4 h-4 text-purple-400 mt-0.5 flex-shrink-0" }),
          createBaseVNode("div", null, [
            createBaseVNode("p", _hoisted_3$6, toDisplayString(unref(t)("settings.qwenCliTitle")), 1),
            createBaseVNode("p", _hoisted_4$6, toDisplayString(unref(t)("settings.qwenCliDescription")), 1)
          ])
        ])
      ]);
    };
  }
});
const _sfc_main$7 = /* @__PURE__ */ defineComponent({
  __name: "StatusMessage",
  props: {
    message: {},
    type: {}
  },
  setup(__props) {
    return (_ctx, _cache) => {
      return __props.message ? (openBlock(), createElementBlock("div", {
        key: 0,
        class: normalizeClass(["status-message", __props.type === "success" ? "status-message-success" : "status-message-error"])
      }, [
        __props.type === "success" ? (openBlock(), createBlock(unref(Check), {
          key: 0,
          class: "w-4 h-4"
        })) : (openBlock(), createBlock(unref(X), {
          key: 1,
          class: "w-4 h-4"
        })),
        createTextVNode(" " + toDisplayString(__props.message), 1)
      ], 2)) : createCommentVNode("", true);
    };
  }
});
const _hoisted_1$6 = { class: "ai-settings-panel" };
const _hoisted_2$6 = { class: "ai-settings-header" };
const _hoisted_3$5 = { class: "ai-settings-icon" };
const _hoisted_4$5 = { class: "ai-settings-title" };
const _hoisted_5$4 = { class: "ai-settings-content" };
const _hoisted_6$4 = { class: "ai-settings-footer" };
const _hoisted_7$4 = ["disabled"];
const _sfc_main$6 = /* @__PURE__ */ defineComponent({
  __name: "AISettings",
  setup(__props) {
    const { t } = useI18n();
    const {
      settings: settings2,
      showApiKey,
      isSaving,
      statusMessage,
      statusType,
      providers,
      currentProvider,
      currentProviderHint,
      currentApiKey,
      availableModels,
      selectedModel,
      needsApiKey,
      needsHostUrl,
      currentHostUrl,
      hostPlaceholder,
      selectProvider,
      updateApiKey,
      clearApiKey,
      updateModel,
      updateHost,
      toggleShowApiKey,
      loadSettings,
      saveSettings
    } = useAISettings();
    onMounted(() => {
      loadSettings();
    });
    return (_ctx, _cache) => {
      return openBlock(), createElementBlock("div", _hoisted_1$6, [
        createBaseVNode("div", _hoisted_2$6, [
          createBaseVNode("div", _hoisted_3$5, [
            createVNode(unref(Lightbulb), { class: "w-4 h-4" })
          ]),
          createBaseVNode("span", _hoisted_4$5, toDisplayString(unref(t)("settings.aiProvider")), 1)
        ]),
        createBaseVNode("div", _hoisted_5$4, [
          createVNode(unref(ProviderSelector), {
            "model-value": unref(settings2).selectedProvider,
            providers: unref(providers),
            description: unref(currentProvider)?.description,
            "onUpdate:modelValue": unref(selectProvider)
          }, null, 8, ["model-value", "providers", "description", "onUpdate:modelValue"]),
          _cache[1] || (_cache[1] = createBaseVNode("div", { class: "ai-settings-divider" }, null, -1)),
          unref(needsApiKey) ? (openBlock(), createBlock(unref(ApiKeyInput), {
            key: 0,
            "model-value": unref(currentApiKey),
            "provider-name": unref(currentProvider)?.name || "",
            "show-key": unref(showApiKey),
            hint: unref(currentProviderHint),
            "onUpdate:modelValue": unref(updateApiKey),
            onToggleVisibility: unref(toggleShowApiKey),
            onClear: unref(clearApiKey)
          }, null, 8, ["model-value", "provider-name", "show-key", "hint", "onUpdate:modelValue", "onToggleVisibility", "onClear"])) : createCommentVNode("", true),
          unref(settings2).selectedProvider === "qwen-cli" ? (openBlock(), createBlock(unref(_sfc_main$8), { key: 1 })) : createCommentVNode("", true),
          unref(settings2).selectedProvider && unref(availableModels).length > 0 ? (openBlock(), createBlock(unref(ModelSelector), {
            key: 2,
            "model-value": unref(selectedModel),
            models: unref(availableModels),
            "onUpdate:modelValue": unref(updateModel)
          }, null, 8, ["model-value", "models", "onUpdate:modelValue"])) : createCommentVNode("", true),
          unref(needsHostUrl) ? (openBlock(), createBlock(unref(HostUrlInput), {
            key: 3,
            "model-value": unref(currentHostUrl),
            placeholder: unref(hostPlaceholder),
            "onUpdate:modelValue": unref(updateHost)
          }, null, 8, ["model-value", "placeholder", "onUpdate:modelValue"])) : createCommentVNode("", true)
        ]),
        createBaseVNode("div", _hoisted_6$4, [
          createBaseVNode("button", {
            onClick: _cache[0] || (_cache[0] = //@ts-ignore
            (...args) => unref(saveSettings) && unref(saveSettings)(...args)),
            disabled: unref(isSaving),
            class: "ai-settings-save-btn"
          }, [
            unref(isSaving) ? (openBlock(), createBlock(unref(LoaderCircle), {
              key: 0,
              class: "animate-spin w-4 h-4"
            })) : createCommentVNode("", true),
            createTextVNode(" " + toDisplayString(unref(isSaving) ? unref(t)("settings.saving") : unref(t)("settings.save")), 1)
          ], 8, _hoisted_7$4),
          createVNode(unref(_sfc_main$7), {
            message: unref(statusMessage),
            type: unref(statusType)
          }, null, 8, ["message", "type"])
        ])
      ]);
    };
  }
});
const AISettings = /* @__PURE__ */ _export_sfc(_sfc_main$6, [["__scopeId", "data-v-8345285b"]]);
const _hoisted_1$5 = { class: "context-detail-panel" };
const _hoisted_2$5 = {
  key: 0,
  class: "detail-empty"
};
const _hoisted_3$4 = { class: "detail-empty-title" };
const _hoisted_4$4 = { class: "detail-empty-hint" };
const _hoisted_5$3 = { class: "detail-header" };
const _hoisted_6$3 = { key: 0 };
const _hoisted_7$3 = { key: 1 };
const _hoisted_8$3 = { class: "detail-info" };
const _hoisted_9$2 = { class: "detail-title-row" };
const _hoisted_10$2 = { class: "detail-title" };
const _hoisted_11$2 = { class: "detail-badge" };
const _hoisted_12$2 = { class: "detail-meta" };
const _hoisted_13$2 = {
  key: 0,
  class: "detail-dot"
};
const _hoisted_14$2 = { key: 1 };
const _hoisted_15$2 = { class: "detail-content" };
const _hoisted_16$2 = { class: "detail-section-title" };
const _hoisted_17$2 = {
  key: 0,
  class: "detail-no-files"
};
const _hoisted_18$2 = { class: "detail-no-files-hint" };
const _hoisted_19$2 = {
  key: 1,
  class: "detail-no-files"
};
const _hoisted_20$2 = {
  key: 2,
  class: "detail-files"
};
const _hoisted_21 = { class: "detail-file-info" };
const _hoisted_22 = { class: "detail-file-name" };
const _hoisted_23 = { class: "detail-file-dir" };
const _hoisted_24 = { class: "detail-footer" };
const _hoisted_25 = ["disabled"];
const _hoisted_26 = {
  key: 0,
  class: "w-4 h-4",
  fill: "none",
  stroke: "currentColor",
  viewBox: "0 0 24 24"
};
const _hoisted_27 = {
  key: 1,
  class: "w-4 h-4 animate-spin",
  fill: "none",
  viewBox: "0 0 24 24"
};
const _hoisted_28 = ["title"];
const _sfc_main$5 = /* @__PURE__ */ defineComponent({
  __name: "ContextMemoryPanel",
  emits: ["switch-to-preview"],
  setup(__props, { emit: __emit }) {
    const emit = __emit;
    const { t } = useI18n();
    const contextStore = useContextStore();
    const fileStore = useFileStore();
    const settingsStore = useSettingsStore();
    const uiStore = useUIStore();
    const isLoading = ref(false);
    const selectedContext = computed(() => contextStore.selectedListItem);
    const contextFiles = computed(() => {
      return selectedContext.value?.files || [];
    });
    const contextName = computed(() => {
      const name = selectedContext.value?.name || t("context.untitled");
      return name.replace(/\s*\[\d+\]\s*$/, "");
    });
    const initials = computed(() => {
      const name = selectedContext.value?.name || "";
      if (!name) {
        const fileCount = selectedContext.value?.fileCount || 0;
        return fileCount > 0 ? String(fileCount) : "?";
      }
      if (name.includes("Пустой")) return "ПК";
      const filesMatch = name.match(/^(\d+)\s*файл/);
      if (filesMatch) {
        return filesMatch[1].length > 2 ? filesMatch[1].slice(0, 2) : filesMatch[1];
      }
      const words = name.split(/[\s\-_]+/).filter((w) => w.length > 0);
      if (words.length >= 2) {
        return (words[0][0] + words[1][0]).toUpperCase();
      }
      return name.slice(0, 2).toUpperCase();
    });
    const avatarColor = computed(() => {
      const colors = ["indigo", "purple", "pink", "blue", "cyan", "teal", "green", "amber"];
      const source = selectedContext.value?.id || selectedContext.value?.name || "";
      let hash = 0;
      for (let i = 0; i < source.length; i++) {
        hash = source.charCodeAt(i) + ((hash << 5) - hash);
      }
      return colors[Math.abs(hash) % colors.length];
    });
    function getFileName(filePath) {
      const parts = filePath.replace(/\\/g, "/").split("/");
      return parts[parts.length - 1] || filePath;
    }
    function getFileDir(filePath) {
      const normalized = filePath.replace(/\\/g, "/");
      const parts = normalized.split("/");
      if (parts.length <= 1) return "";
      parts.pop();
      const path = parts.join("/");
      const cleaned = path.replace(/^[A-Za-z]:/, "").replace(/^\/?(Sources|Projects|repos|workspace|home|Users\/[^/]+)\/[^/]+\//i, "").replace(/^\/?/, "");
      return cleaned || "/";
    }
    function formatSize(bytes) {
      return formatContextSize(bytes);
    }
    function formatDate(dateStr) {
      if (!dateStr) return "";
      const date = new Date(dateStr);
      return date.toLocaleDateString("ru-RU", {
        day: "numeric",
        month: "short",
        hour: "2-digit",
        minute: "2-digit"
      });
    }
    async function loadAndBuildContext() {
      if (!selectedContext.value?.files?.length) {
        uiStore.addToast(t("context.noFilesInfo"), "warning");
        return;
      }
      isLoading.value = true;
      try {
        const normalizedFiles = selectedContext.value.files.map(
          (file) => file.replace(/\\/g, "/")
        );
        fileStore.clearSelection();
        fileStore.selectMultiple(normalizedFiles);
        const options2 = {
          outputFormat: settingsStore.settings.context.outputFormat,
          stripComments: settingsStore.settings.context.stripComments,
          excludeTests: settingsStore.settings.context.excludeTests,
          maxTokens: settingsStore.settings.context.maxTokens
        };
        await contextStore.buildContext(normalizedFiles, options2);
        emit("switch-to-preview");
        uiStore.addToast(t("context.contextRestored"), "success");
      } catch (error) {
        uiStore.addToast(t("context.saveError"), "error");
      } finally {
        isLoading.value = false;
      }
    }
    async function deleteContext() {
      if (!selectedContext.value?.id) return;
      try {
        await contextStore.deleteContext(selectedContext.value.id);
        contextStore.clearListSelection();
        uiStore.addToast(t("context.deleted"), "success");
      } catch {
        uiStore.addToast(t("context.deleteError"), "error");
      }
    }
    return (_ctx, _cache) => {
      return openBlock(), createElementBlock("div", _hoisted_1$5, [
        !selectedContext.value ? (openBlock(), createElementBlock("div", _hoisted_2$5, [
          _cache[0] || (_cache[0] = createBaseVNode("div", { class: "detail-empty-icon" }, [
            createBaseVNode("svg", {
              class: "w-10 h-10",
              fill: "none",
              stroke: "currentColor",
              viewBox: "0 0 24 24"
            }, [
              createBaseVNode("path", {
                "stroke-linecap": "round",
                "stroke-linejoin": "round",
                "stroke-width": "1.5",
                d: "M15 15l-2 5L9 9l11 4-5 2zm0 0l5 5M7.188 2.239l.777 2.897M5.136 7.965l-2.898-.777M13.95 4.05l-2.122 2.122m-5.657 5.656l-2.12 2.122"
              })
            ])
          ], -1)),
          createBaseVNode("p", _hoisted_3$4, toDisplayString(unref(t)("context.selectToPreview")), 1),
          createBaseVNode("p", _hoisted_4$4, toDisplayString(unref(t)("context.selectToPreviewHint")), 1)
        ])) : (openBlock(), createElementBlock(Fragment, { key: 1 }, [
          createBaseVNode("div", _hoisted_5$3, [
            createBaseVNode("div", {
              class: normalizeClass(["detail-avatar", `detail-avatar--${avatarColor.value}`])
            }, [
              selectedContext.value.isFavorite ? (openBlock(), createElementBlock("span", _hoisted_6$3, "⭐")) : (openBlock(), createElementBlock("span", _hoisted_7$3, toDisplayString(initials.value), 1))
            ], 2),
            createBaseVNode("div", _hoisted_8$3, [
              createBaseVNode("div", _hoisted_9$2, [
                createBaseVNode("h2", _hoisted_10$2, toDisplayString(contextName.value), 1),
                createBaseVNode("span", _hoisted_11$2, toDisplayString(selectedContext.value.fileCount) + " " + toDisplayString(unref(t)("context.filesShort")), 1)
              ]),
              createBaseVNode("div", _hoisted_12$2, [
                createBaseVNode("span", null, toDisplayString(formatSize(selectedContext.value.totalSize)), 1),
                selectedContext.value.createdAt ? (openBlock(), createElementBlock("span", _hoisted_13$2, "•")) : createCommentVNode("", true),
                selectedContext.value.createdAt ? (openBlock(), createElementBlock("span", _hoisted_14$2, toDisplayString(formatDate(selectedContext.value.createdAt)), 1)) : createCommentVNode("", true)
              ])
            ])
          ]),
          createBaseVNode("div", _hoisted_15$2, [
            createBaseVNode("h3", _hoisted_16$2, toDisplayString(unref(t)("context.contextFiles")), 1),
            !contextFiles.value.length && selectedContext.value.fileCount > 0 ? (openBlock(), createElementBlock("div", _hoisted_17$2, [
              createBaseVNode("p", null, toDisplayString(unref(t)("context.filesListUnavailable")), 1),
              createBaseVNode("p", _hoisted_18$2, toDisplayString(selectedContext.value.fileCount) + " " + toDisplayString(unref(t)("context.filesInContext")), 1)
            ])) : !contextFiles.value.length ? (openBlock(), createElementBlock("div", _hoisted_19$2, toDisplayString(unref(t)("context.noFilesInfo")), 1)) : (openBlock(), createElementBlock("ul", _hoisted_20$2, [
              (openBlock(true), createElementBlock(Fragment, null, renderList(contextFiles.value, (file) => {
                return openBlock(), createElementBlock("li", {
                  key: file,
                  class: "detail-file"
                }, [
                  _cache[1] || (_cache[1] = createBaseVNode("div", { class: "detail-file-icon-wrap" }, [
                    createBaseVNode("svg", {
                      class: "detail-file-icon",
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
                    ])
                  ], -1)),
                  createBaseVNode("div", _hoisted_21, [
                    createBaseVNode("span", _hoisted_22, toDisplayString(getFileName(file)), 1),
                    createBaseVNode("span", _hoisted_23, toDisplayString(getFileDir(file)), 1)
                  ])
                ]);
              }), 128))
            ]))
          ]),
          createBaseVNode("div", _hoisted_24, [
            createBaseVNode("button", {
              onClick: loadAndBuildContext,
              disabled: isLoading.value,
              class: "detail-load-btn"
            }, [
              !isLoading.value ? (openBlock(), createElementBlock("svg", _hoisted_26, [..._cache[2] || (_cache[2] = [
                createBaseVNode("path", {
                  "stroke-linecap": "round",
                  "stroke-linejoin": "round",
                  "stroke-width": "2",
                  d: "M13 10V3L4 14h7v7l9-11h-7z"
                }, null, -1)
              ])])) : (openBlock(), createElementBlock("svg", _hoisted_27, [..._cache[3] || (_cache[3] = [
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
              createTextVNode(" " + toDisplayString(isLoading.value ? unref(t)("context.building") : unref(t)("context.loadAndBuild")), 1)
            ], 8, _hoisted_25),
            createBaseVNode("button", {
              onClick: deleteContext,
              class: "detail-delete-btn",
              title: unref(t)("context.delete")
            }, [..._cache[4] || (_cache[4] = [
              createBaseVNode("svg", {
                fill: "none",
                stroke: "currentColor",
                viewBox: "0 0 24 24"
              }, [
                createBaseVNode("path", {
                  "stroke-linecap": "round",
                  "stroke-linejoin": "round",
                  "stroke-width": "2",
                  d: "M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"
                })
              ], -1)
            ])], 8, _hoisted_28)
          ])
        ], 64))
      ]);
    };
  }
});
const ContextMemoryPanel = /* @__PURE__ */ _export_sfc(_sfc_main$5, [["__scopeId", "data-v-2a1be69f"]]);
const _hoisted_1$4 = { class: "toggle-row" };
const _hoisted_2$4 = { class: "toggle-label" };
const _sfc_main$4 = /* @__PURE__ */ defineComponent({
  __name: "ToggleItem",
  props: {
    modelValue: { type: Boolean },
    label: {}
  },
  emits: ["update:modelValue"],
  setup(__props, { emit: __emit }) {
    const emit = __emit;
    return (_ctx, _cache) => {
      return openBlock(), createElementBlock("label", _hoisted_1$4, [
        createBaseVNode("span", _hoisted_2$4, toDisplayString(__props.label), 1),
        createBaseVNode("button", {
          type: "button",
          class: normalizeClass(["toggle-switch", { active: __props.modelValue }]),
          onClick: _cache[0] || (_cache[0] = ($event) => emit("update:modelValue", !__props.modelValue))
        }, [..._cache[1] || (_cache[1] = [
          createBaseVNode("span", { class: "toggle-thumb" }, null, -1)
        ])], 2)
      ]);
    };
  }
});
const ToggleItem = /* @__PURE__ */ _export_sfc(_sfc_main$4, [["__scopeId", "data-v-c81c3c4f"]]);
const _hoisted_1$3 = { class: "inspector" };
const _hoisted_2$3 = { class: "inspector-content" };
const _hoisted_3$3 = { class: "inspector-section" };
const _hoisted_4$3 = { class: "task-wrapper" };
const _hoisted_5$2 = ["placeholder"];
const _hoisted_6$2 = { class: "template-icon" };
const _hoisted_7$2 = { class: "template-name" };
const _hoisted_8$2 = { class: "inspector-section" };
const _hoisted_9$1 = { class: "inspector-section" };
const _hoisted_10$1 = { class: "inspector-section" };
const _hoisted_11$1 = { class: "chunking-header" };
const _hoisted_12$1 = { class: "section-header" };
const _hoisted_13$1 = {
  key: 0,
  class: "chunking-content"
};
const _hoisted_14$1 = { class: "chunk-row" };
const _hoisted_15$1 = { class: "chunk-presets" };
const _hoisted_16$1 = ["onClick"];
const _hoisted_17$1 = { class: "chunk-row" };
const _hoisted_18$1 = { class: "strategy-segment" };
const _hoisted_19$1 = ["title"];
const _hoisted_20$1 = ["title"];
const _sfc_main$3 = /* @__PURE__ */ defineComponent({
  __name: "ExportSettings",
  setup(__props) {
    const { t } = useI18n();
    const settingsStore = useSettingsStore();
    const templateStore = useTemplateStore();
    const settings2 = computed(() => settingsStore.settings.context);
    const activeTemplate = computed(() => templateStore.activeTemplate);
    const task2 = computed({ get: () => templateStore.currentTask, set: (v) => templateStore.setTask(v) });
    const chunkPresets = [
      { value: 32e3, label: "32K" },
      { value: 64e3, label: "64K" },
      { value: 128e3, label: "128K" }
    ];
    const showCustomInput = ref(false);
    const customValue = ref("");
    const customInputRef = ref(null);
    const isCustomActive = computed(() => {
      return !chunkPresets.some((p) => isChunkPresetActive(p.value));
    });
    function isChunkPresetActive(value) {
      return Math.abs(settings2.value.maxTokensPerChunk - value) <= value * 0.05;
    }
    function selectPreset(value) {
      showCustomInput.value = false;
      update("maxTokensPerChunk", value);
    }
    function enableCustomInput() {
      showCustomInput.value = true;
      customValue.value = formatChunkTokens(settings2.value.maxTokensPerChunk);
      nextTick(() => customInputRef.value?.focus());
    }
    function applyCustomValue() {
      const input = customValue.value.trim().toUpperCase();
      let value = 0;
      if (input.endsWith("K")) {
        value = parseFloat(input.slice(0, -1)) * 1e3;
      } else if (input.endsWith("M")) {
        value = parseFloat(input.slice(0, -1)) * 1e6;
      } else {
        value = parseFloat(input);
      }
      if (!isNaN(value) && value > 0) {
        update("maxTokensPerChunk", Math.round(value));
      }
      showCustomInput.value = false;
    }
    function update(key, value) {
      settingsStore.updateContextSettings({ [key]: value });
    }
    function formatChunkTokens(n) {
      if (n >= 1e6) return `${(n / 1e6).toFixed(1)}M`;
      if (n >= 1e3) return `${Math.round(n / 1e3)}K`;
      return n.toString();
    }
    return (_ctx, _cache) => {
      return openBlock(), createElementBlock("div", _hoisted_1$3, [
        createBaseVNode("div", _hoisted_2$3, [
          createBaseVNode("div", _hoisted_3$3, [
            createBaseVNode("div", _hoisted_4$3, [
              _cache[24] || (_cache[24] = createBaseVNode("span", { class: "task-icon" }, "✨", -1)),
              withDirectives(createBaseVNode("textarea", {
                "onUpdate:modelValue": _cache[0] || (_cache[0] = ($event) => task2.value = $event),
                placeholder: unref(t)("templates.taskPlaceholder"),
                class: "task-input",
                rows: "2"
              }, null, 8, _hoisted_5$2), [
                [vModelText, task2.value]
              ])
            ]),
            createBaseVNode("button", {
              class: "template-row",
              onClick: _cache[1] || (_cache[1] = ($event) => unref(templateStore).openModal())
            }, [
              createBaseVNode("span", _hoisted_6$2, toDisplayString(activeTemplate.value?.icon || "📝"), 1),
              createBaseVNode("span", _hoisted_7$2, toDisplayString(activeTemplate.value?.name || unref(t)("templates.select")), 1),
              _cache[25] || (_cache[25] = createBaseVNode("svg", {
                class: "template-arrow",
                fill: "none",
                viewBox: "0 0 24 24",
                stroke: "currentColor"
              }, [
                createBaseVNode("path", {
                  "stroke-linecap": "round",
                  "stroke-linejoin": "round",
                  "stroke-width": "2",
                  d: "M9 5l7 7-7 7"
                })
              ], -1))
            ])
          ]),
          _cache[31] || (_cache[31] = createBaseVNode("div", { class: "inspector-divider" }, null, -1)),
          createBaseVNode("div", _hoisted_8$2, [
            _cache[26] || (_cache[26] = createBaseVNode("div", { class: "section-header" }, "OUTPUT", -1)),
            createVNode(ToggleItem, {
              modelValue: settings2.value.applyTemplateOnCopy,
              "onUpdate:modelValue": [
                _cache[2] || (_cache[2] = ($event) => settings2.value.applyTemplateOnCopy = $event),
                _cache[3] || (_cache[3] = ($event) => update("applyTemplateOnCopy", $event))
              ],
              label: unref(t)("export.applyTemplate")
            }, null, 8, ["modelValue", "label"]),
            createVNode(ToggleItem, {
              modelValue: settings2.value.includeManifest,
              "onUpdate:modelValue": [
                _cache[4] || (_cache[4] = ($event) => settings2.value.includeManifest = $event),
                _cache[5] || (_cache[5] = ($event) => update("includeManifest", $event))
              ],
              label: unref(t)("export.includeManifest")
            }, null, 8, ["modelValue", "label"]),
            createVNode(ToggleItem, {
              modelValue: settings2.value.includeLineNumbers,
              "onUpdate:modelValue": [
                _cache[6] || (_cache[6] = ($event) => settings2.value.includeLineNumbers = $event),
                _cache[7] || (_cache[7] = ($event) => update("includeLineNumbers", $event))
              ],
              label: unref(t)("export.includeLineNumbers")
            }, null, 8, ["modelValue", "label"]),
            createVNode(ToggleItem, {
              modelValue: settings2.value.stripComments,
              "onUpdate:modelValue": [
                _cache[8] || (_cache[8] = ($event) => settings2.value.stripComments = $event),
                _cache[9] || (_cache[9] = ($event) => update("stripComments", $event))
              ],
              label: unref(t)("export.stripComments")
            }, null, 8, ["modelValue", "label"])
          ]),
          _cache[32] || (_cache[32] = createBaseVNode("div", { class: "inspector-divider" }, null, -1)),
          createBaseVNode("div", _hoisted_9$1, [
            _cache[27] || (_cache[27] = createBaseVNode("div", { class: "section-header" }, "OPTIMIZATION", -1)),
            createVNode(ToggleItem, {
              modelValue: settings2.value.excludeTests,
              "onUpdate:modelValue": [
                _cache[10] || (_cache[10] = ($event) => settings2.value.excludeTests = $event),
                _cache[11] || (_cache[11] = ($event) => update("excludeTests", $event))
              ],
              label: unref(t)("export.excludeTests")
            }, null, 8, ["modelValue", "label"]),
            createVNode(ToggleItem, {
              modelValue: settings2.value.stripLicense,
              "onUpdate:modelValue": [
                _cache[12] || (_cache[12] = ($event) => settings2.value.stripLicense = $event),
                _cache[13] || (_cache[13] = ($event) => update("stripLicense", $event))
              ],
              label: unref(t)("export.stripLicense")
            }, null, 8, ["modelValue", "label"]),
            createVNode(ToggleItem, {
              modelValue: settings2.value.compactDataFiles,
              "onUpdate:modelValue": [
                _cache[14] || (_cache[14] = ($event) => settings2.value.compactDataFiles = $event),
                _cache[15] || (_cache[15] = ($event) => update("compactDataFiles", $event))
              ],
              label: unref(t)("export.compactDataFiles")
            }, null, 8, ["modelValue", "label"]),
            createVNode(ToggleItem, {
              modelValue: settings2.value.trimWhitespace,
              "onUpdate:modelValue": [
                _cache[16] || (_cache[16] = ($event) => settings2.value.trimWhitespace = $event),
                _cache[17] || (_cache[17] = ($event) => update("trimWhitespace", $event))
              ],
              label: unref(t)("export.trimWhitespace")
            }, null, 8, ["modelValue", "label"]),
            createVNode(ToggleItem, {
              modelValue: settings2.value.collapseEmptyLines,
              "onUpdate:modelValue": [
                _cache[18] || (_cache[18] = ($event) => settings2.value.collapseEmptyLines = $event),
                _cache[19] || (_cache[19] = ($event) => update("collapseEmptyLines", $event))
              ],
              label: unref(t)("export.collapseEmptyLines")
            }, null, 8, ["modelValue", "label"])
          ]),
          _cache[33] || (_cache[33] = createBaseVNode("div", { class: "inspector-divider" }, null, -1)),
          createBaseVNode("div", _hoisted_10$1, [
            createBaseVNode("div", _hoisted_11$1, [
              createBaseVNode("span", _hoisted_12$1, toDisplayString(unref(t)("export.chunking")), 1),
              createBaseVNode("button", {
                class: normalizeClass(["chunking-toggle", { active: settings2.value.enableAutoSplit }]),
                onClick: _cache[20] || (_cache[20] = ($event) => update("enableAutoSplit", !settings2.value.enableAutoSplit))
              }, [..._cache[28] || (_cache[28] = [
                createBaseVNode("span", { class: "chunking-toggle-thumb" }, null, -1)
              ])], 2)
            ]),
            createVNode(Transition, { name: "accordion" }, {
              default: withCtx(() => [
                settings2.value.enableAutoSplit ? (openBlock(), createElementBlock("div", _hoisted_13$1, [
                  createBaseVNode("div", _hoisted_14$1, [
                    createBaseVNode("div", _hoisted_15$1, [
                      (openBlock(), createElementBlock(Fragment, null, renderList(chunkPresets, (preset) => {
                        return createBaseVNode("button", {
                          key: preset.value,
                          class: normalizeClass(["preset-btn", { active: isChunkPresetActive(preset.value) }]),
                          onClick: ($event) => selectPreset(preset.value)
                        }, toDisplayString(preset.label), 11, _hoisted_16$1);
                      }), 64)),
                      showCustomInput.value ? withDirectives((openBlock(), createElementBlock("input", {
                        key: 0,
                        type: "text",
                        "onUpdate:modelValue": _cache[21] || (_cache[21] = ($event) => customValue.value = $event),
                        onBlur: applyCustomValue,
                        onKeyup: withKeys(applyCustomValue, ["enter"]),
                        class: "preset-custom",
                        placeholder: "Custom",
                        ref_key: "customInputRef",
                        ref: customInputRef
                      }, null, 544)), [
                        [vModelText, customValue.value]
                      ]) : (openBlock(), createElementBlock("button", {
                        key: 1,
                        class: normalizeClass(["preset-btn preset-btn--custom", { active: isCustomActive.value }]),
                        onClick: enableCustomInput
                      }, toDisplayString(isCustomActive.value ? formatChunkTokens(settings2.value.maxTokensPerChunk) : "..."), 3))
                    ])
                  ]),
                  createBaseVNode("div", _hoisted_17$1, [
                    createBaseVNode("div", _hoisted_18$1, [
                      createBaseVNode("button", {
                        class: normalizeClass(["strategy-btn", { active: settings2.value.splitStrategy === "smart" || settings2.value.splitStrategy === "file" }]),
                        onClick: _cache[22] || (_cache[22] = ($event) => update("splitStrategy", "smart")),
                        title: unref(t)("export.strategy.smartDesc")
                      }, [..._cache[29] || (_cache[29] = [
                        createBaseVNode("span", null, "🧠", -1),
                        createBaseVNode("span", null, "Smart", -1)
                      ])], 10, _hoisted_19$1),
                      createBaseVNode("button", {
                        class: normalizeClass(["strategy-btn", { active: settings2.value.splitStrategy === "token" }]),
                        onClick: _cache[23] || (_cache[23] = ($event) => update("splitStrategy", "token")),
                        title: unref(t)("export.strategy.hardDesc")
                      }, [..._cache[30] || (_cache[30] = [
                        createBaseVNode("span", null, "✂️", -1),
                        createBaseVNode("span", null, "Hard", -1)
                      ])], 10, _hoisted_20$1)
                    ])
                  ])
                ])) : createCommentVNode("", true)
              ]),
              _: 1
            })
          ])
        ]),
        createVNode(unref(TemplateModal))
      ]);
    };
  }
});
const ExportSettings = /* @__PURE__ */ _export_sfc(_sfc_main$3, [["__scopeId", "data-v-0e73108d"]]);
const _hoisted_1$2 = { class: "right-sidebar-container" };
const _hoisted_2$2 = { class: "right-sidebar-tabs" };
const _hoisted_3$2 = ["onClick", "title", "aria-label", "aria-pressed"];
const _hoisted_4$2 = { class: "right-sidebar-content scrollable-y" };
const _hoisted_5$1 = {
  key: 0,
  class: "right-sidebar-tab-content-full",
  "data-tour": "ai-chat"
};
const _hoisted_6$1 = {
  key: 1,
  class: "right-sidebar-tab-content-full"
};
const _hoisted_7$1 = {
  key: 2,
  class: "right-sidebar-tab-content"
};
const _hoisted_8$1 = {
  key: 3,
  class: "right-sidebar-tab-content"
};
const _sfc_main$2 = /* @__PURE__ */ defineComponent({
  __name: "RightSidebar",
  setup(__props) {
    const { t } = useI18n();
    const currentTab = ref("chat");
    const tabs = [
      { id: "chat", label: "sidebar.chat", icon: MessageCircle, textClass: "text-purple-300", indicatorClass: "tabs-indicator-purple" },
      { id: "memory", label: "sidebar.memory", icon: BookMarked, textClass: "text-cyan-300", indicatorClass: "tabs-indicator-cyan" },
      { id: "export", label: "sidebar.exportSettings", icon: FileDown, textClass: "text-emerald-300", indicatorClass: "tabs-indicator-emerald" },
      { id: "settings", label: "sidebar.aiConfig", icon: Bot, textClass: "text-indigo-300", indicatorClass: "tabs-indicator-indigo" }
    ];
    const tabIndex = computed(() => tabs.findIndex((t2) => t2.id === currentTab.value));
    const tabIndicatorClass = computed(() => tabs.find((t2) => t2.id === currentTab.value)?.indicatorClass || "tabs-indicator-indigo");
    watch(currentTab, (tab) => {
      try {
        localStorage.setItem("right-sidebar-tab", tab);
      } catch {
      }
    });
    try {
      const savedTab = localStorage.getItem("right-sidebar-tab");
      if (savedTab && tabs.some((t2) => t2.id === savedTab)) {
        currentTab.value = savedTab;
      }
    } catch {
    }
    return (_ctx, _cache) => {
      return openBlock(), createElementBlock("div", _hoisted_1$2, [
        createBaseVNode("div", _hoisted_2$2, [
          createBaseVNode("div", {
            class: normalizeClass(["tabs-indicator-icon", tabIndicatorClass.value]),
            style: normalizeStyle({ transform: `translateX(calc(${tabIndex.value} * (2rem + 0.25rem)))` })
          }, null, 6),
          (openBlock(), createElementBlock(Fragment, null, renderList(tabs, (tab) => {
            return createBaseVNode("button", {
              key: tab.id,
              onClick: ($event) => currentTab.value = tab.id,
              title: unref(t)(tab.label),
              "aria-label": unref(t)(tab.label),
              "aria-pressed": currentTab.value === tab.id,
              class: normalizeClass(["sidebar-tab-icon", currentTab.value === tab.id ? `sidebar-tab-icon-active ${tab.textClass}` : ""])
            }, [
              (openBlock(), createBlock(resolveDynamicComponent(tab.icon), {
                class: "w-4 h-4",
                "aria-hidden": "true"
              }))
            ], 10, _hoisted_3$2);
          }), 64))
        ]),
        createBaseVNode("div", _hoisted_4$2, [
          currentTab.value === "chat" ? (openBlock(), createElementBlock("div", _hoisted_5$1, [
            createVNode(AIChat)
          ])) : currentTab.value === "memory" ? (openBlock(), createElementBlock("div", _hoisted_6$1, [
            createVNode(ContextMemoryPanel)
          ])) : currentTab.value === "export" ? (openBlock(), createElementBlock("div", _hoisted_7$1, [
            createVNode(ExportSettings)
          ])) : currentTab.value === "settings" ? (openBlock(), createElementBlock("div", _hoisted_8$1, [
            createVNode(AISettings)
          ])) : createCommentVNode("", true)
        ])
      ]);
    };
  }
});
const _hoisted_1$1 = { class: "workspace-container layout-grid-main-footer" };
const _hoisted_2$1 = { class: "workspace-layout" };
const _hoisted_3$1 = { class: "workspace-center layout-fill layout-column layout-clip" };
const _hoisted_4$1 = ["title"];
const _sfc_main$1 = /* @__PURE__ */ defineComponent({
  __name: "MainWorkspace",
  setup(__props, { expose: __expose }) {
    const logger2 = useLogger("MainWorkspace");
    const templateStore = useTemplateStore();
    const projectStore = useProjectStore();
    const { t } = useI18n();
    const contextStore = useContextStore();
    const fileStore = useFileStore();
    const uiStore = useUIStore();
    const settingsStore = useSettingsStore();
    const exportModalRef = ref(null);
    const showRightSidebar = ref(loadSidebarState());
    function loadSidebarState() {
      try {
        const saved = localStorage.getItem("right-sidebar-visible");
        return saved !== "false";
      } catch {
        return true;
      }
    }
    function toggleRightSidebar() {
      showRightSidebar.value = !showRightSidebar.value;
      try {
        localStorage.setItem("right-sidebar-visible", String(showRightSidebar.value));
      } catch {
      }
    }
    function resetPanelSizes() {
      leftResize.resetToDefault();
      rightResize.resetToDefault();
      uiStore.addToast(t("workspace.layoutReset"), "success");
    }
    const leftResize = useResizablePanel({
      minWidth: 280,
      maxWidth: 700,
      defaultWidth: 380,
      storageKey: "workspace-left-width"
    });
    const leftWidth = leftResize.width;
    const rightResize = useResizablePanel({
      minWidth: 320,
      maxWidth: 700,
      defaultWidth: 380,
      storageKey: "workspace-right-width",
      invertDirection: true
      // Тянем влево = увеличиваем ширину
    });
    const leftPanelRef = leftResize.panelRef;
    const rightPanelRef = rightResize.panelRef;
    const rightWidth = rightResize.width;
    function handlePreviewFile(_filePath) {
    }
    const handleGlobalBuildContext = () => handleBuildContext();
    const handleGlobalOpenExport = () => handleOpenExport();
    const handleGlobalCopyContext = async () => {
      if (contextStore.hasContext && contextStore.contextId) {
        try {
          const filesContent = await contextStore.getFullContextContent();
          let content;
          if (settingsStore.settings.context.applyTemplateOnCopy && templateStore.activeTemplate) {
            const files2 = contextStore.summary?.files || [];
            const templateContext = {
              fileTree: generateFileTree(files2, projectStore.projectName),
              files: filesContent,
              task: templateStore.currentTask,
              userRules: templateStore.userRules,
              fileCount: contextStore.fileCount,
              tokenCount: contextStore.tokenCount,
              languages: detectLanguages(files2),
              projectName: projectStore.projectName
            };
            content = templateStore.generatePrompt(templateContext);
          } else {
            content = filesContent;
          }
          await navigator.clipboard.writeText(content);
          uiStore.addToast(t("toast.contextCopied"), "success");
        } catch (error) {
          logger2.error("Failed to copy context:", error);
          uiStore.addToast(t("toast.copyError"), "error");
        }
      }
    };
    onMounted(() => {
      window.addEventListener("global-build-context", handleGlobalBuildContext);
      window.addEventListener("global-open-export", handleGlobalOpenExport);
      window.addEventListener("global-copy-context", handleGlobalCopyContext);
    });
    onUnmounted(() => {
      window.removeEventListener("global-build-context", handleGlobalBuildContext);
      window.removeEventListener("global-open-export", handleGlobalOpenExport);
      window.removeEventListener("global-copy-context", handleGlobalCopyContext);
    });
    watch(
      () => [
        settingsStore.settings.context.outputFormat,
        settingsStore.settings.context.stripComments
      ],
      async ([newFormat, newStripComments], [oldFormat, oldStripComments]) => {
        if (!contextStore.hasContext || contextStore.isBuilding) return;
        if (newFormat === oldFormat && newStripComments === oldStripComments) return;
        const currentFiles = contextStore.summary?.files;
        if (!currentFiles || currentFiles.length === 0) {
          if (fileStore.selectedPaths.size === 0) return;
        }
        const filePaths = currentFiles && currentFiles.length > 0 ? currentFiles : Array.from(fileStore.selectedPaths);
        try {
          const options2 = {
            maxTokens: settingsStore.settings.context.maxTokens,
            stripComments: settingsStore.settings.context.stripComments,
            includeTests: settingsStore.settings.context.includeTests,
            splitStrategy: settingsStore.settings.context.splitStrategy,
            outputFormat: settingsStore.settings.context.outputFormat,
            // Content optimization options
            excludeTests: settingsStore.settings.context.excludeTests,
            collapseEmptyLines: settingsStore.settings.context.collapseEmptyLines,
            stripLicense: settingsStore.settings.context.stripLicense,
            compactDataFiles: settingsStore.settings.context.compactDataFiles,
            trimWhitespace: settingsStore.settings.context.trimWhitespace
          };
          await contextStore.buildContext(filePaths, options2);
          uiStore.addToast(t("context.rebuilt"), "success");
        } catch (error) {
          logger2.error("Failed to rebuild context:", error);
        }
      }
    );
    async function handleBuildContext() {
      if (fileStore.selectedPaths.size === 0) {
        uiStore.addToast("Выберите файлы для построения контекста", "warning");
        return;
      }
      if (contextStore.isBuilding) return;
      try {
        const filePaths = Array.from(fileStore.selectedPaths);
        const options2 = {
          maxTokens: settingsStore.settings.context.maxTokens,
          stripComments: settingsStore.settings.context.stripComments,
          includeTests: settingsStore.settings.context.includeTests,
          splitStrategy: settingsStore.settings.context.splitStrategy,
          outputFormat: settingsStore.settings.context.outputFormat,
          // Output options
          includeManifest: settingsStore.settings.context.includeManifest,
          includeLineNumbers: settingsStore.settings.context.includeLineNumbers,
          // Content optimization options
          excludeTests: settingsStore.settings.context.excludeTests,
          collapseEmptyLines: settingsStore.settings.context.collapseEmptyLines,
          stripLicense: settingsStore.settings.context.stripLicense,
          compactDataFiles: settingsStore.settings.context.compactDataFiles,
          trimWhitespace: settingsStore.settings.context.trimWhitespace
        };
        await contextStore.buildContext(filePaths, options2);
        uiStore.addToast(t("toast.contextBuilt"), "success", 5e3, {
          label: t("context.copy"),
          icon: "📋",
          onClick: async () => {
            try {
              const content = await contextStore.getFullContextContent();
              await navigator.clipboard.writeText(content);
              uiStore.addToast(t("toast.contextCopied"), "success");
            } catch {
              uiStore.addToast(t("toast.copyError"), "error");
            }
          }
        });
      } catch (error) {
        logger2.error("Failed to build context:", error);
        if (error instanceof Error && error.message === "TOKEN_LIMIT_EXCEEDED") {
          const storeError = contextStore.error;
          if (storeError?.startsWith("TOKEN_LIMIT_EXCEEDED:")) {
            const parts = storeError.split(":");
            const actual = Number(parts[1]);
            const limit = Number(parts[2]);
            const actualK = Math.round(actual / 1e3);
            const limitK = Math.round(limit / 1e3);
            uiStore.addToast(
              t("error.tokenLimitExceeded", { actual: actualK, limit: limitK }),
              "error"
            );
          } else {
            uiStore.addToast(t("error.tokenLimitGeneric"), "error");
          }
          return;
        }
        const errorMsg = error instanceof Error ? error.message : "Unknown error";
        uiStore.addToast(`${t("toast.contextError")}: ${errorMsg}`, "error");
      }
    }
    function handleOpenExport() {
      if (!contextStore.hasContext) {
        uiStore.addToast("Сначала постройте контекст", "warning");
        return;
      }
      exportModalRef.value?.open();
    }
    __expose({ leftPanelRef, rightPanelRef });
    return (_ctx, _cache) => {
      return openBlock(), createElementBlock("div", _hoisted_1$1, [
        _cache[5] || (_cache[5] = createBaseVNode("div", { class: "workspace-bg" }, [
          createBaseVNode("div", { class: "workspace-glow workspace-glow-1" }),
          createBaseVNode("div", { class: "workspace-glow workspace-glow-2" })
        ], -1)),
        createBaseVNode("div", _hoisted_2$1, [
          createBaseVNode("div", {
            ref_key: "leftPanelRef",
            ref: leftPanelRef,
            style: normalizeStyle({ width: `${unref(leftWidth)}px` }),
            class: "workspace-panel workspace-panel-left"
          }, [
            createVNode(LeftSidebar, {
              onPreviewFile: handlePreviewFile,
              onBuildContext: handleBuildContext
            })
          ], 4),
          createBaseVNode("div", {
            class: normalizeClass(["resize-handle", { "is-resizing": unref(leftResize).isResizing.value }]),
            onMousedown: _cache[0] || (_cache[0] = (e) => unref(leftResize).onMouseDown(e))
          }, [..._cache[2] || (_cache[2] = [
            createBaseVNode("div", { class: "resize-handle-line" }, null, -1)
          ])], 34),
          createBaseVNode("div", _hoisted_3$1, [
            createVNode(CenterWorkspace)
          ]),
          createBaseVNode("button", {
            onClick: toggleRightSidebar,
            class: normalizeClass(["sidebar-toggle", { "sidebar-toggle-open": showRightSidebar.value }]),
            style: normalizeStyle(showRightSidebar.value ? { right: `${unref(rightWidth)}px` } : {}),
            title: unref(t)("sidebar.toggle")
          }, [
            (openBlock(), createElementBlock("svg", {
              class: normalizeClass(["w-4 h-4 transition-transform", { "rotate-180": !showRightSidebar.value }]),
              fill: "none",
              stroke: "currentColor",
              viewBox: "0 0 24 24"
            }, [..._cache[3] || (_cache[3] = [
              createBaseVNode("path", {
                "stroke-linecap": "round",
                "stroke-linejoin": "round",
                "stroke-width": "2",
                d: "M9 5l7 7-7 7"
              }, null, -1)
            ])], 2))
          ], 14, _hoisted_4$1),
          showRightSidebar.value ? (openBlock(), createElementBlock("div", {
            key: 0,
            class: normalizeClass(["resize-handle", { "is-resizing": unref(rightResize).isResizing.value }]),
            onMousedown: _cache[1] || (_cache[1] = (e) => unref(rightResize).onMouseDown(e))
          }, [..._cache[4] || (_cache[4] = [
            createBaseVNode("div", { class: "resize-handle-line" }, null, -1)
          ])], 34)) : createCommentVNode("", true),
          createVNode(Transition, { name: "slide-right" }, {
            default: withCtx(() => [
              showRightSidebar.value ? (openBlock(), createElementBlock("div", {
                key: 0,
                ref_key: "rightPanelRef",
                ref: rightPanelRef,
                style: normalizeStyle({ width: `${unref(rightWidth)}px` }),
                class: "workspace-panel workspace-panel-right"
              }, [
                createVNode(_sfc_main$2, { onOpenExport: handleOpenExport })
              ], 4)) : createCommentVNode("", true)
            ]),
            _: 1
          })
        ]),
        createVNode(ActionBar, {
          class: "workspace-actionbar",
          onOpenExport: handleOpenExport,
          onResetLayout: resetPanelSizes
        }),
        createVNode(ExportModal, {
          ref_key: "exportModalRef",
          ref: exportModalRef
        }, null, 512)
      ]);
    };
  }
});
const MainWorkspace = /* @__PURE__ */ _export_sfc(_sfc_main$1, [["__scopeId", "data-v-e5e9507f"]]);
const logger$1 = useLogger("MemoryMonitor");
const DEFAULT_OPTIONS = {
  warningThreshold: 150,
  // FIXED: 150MB (was incorrectly 15MB)
  criticalThreshold: 250,
  // FIXED: 250MB (was incorrectly 25MB)
  pollingInterval: 1e4,
  // Check every 10 seconds to reduce overhead
  showToasts: true,
  autoCleanup: true
};
class MemoryMonitor {
  static instance;
  options;
  monitorInterval = null;
  warningIssued = false;
  criticalIssued = false;
  uiStore = useUIStore();
  // Cache store imports to avoid repeated dynamic imports
  storeImportsCache = {};
  // Track if we're in dev mode to reduce monitoring overhead
  isDevMode = false;
  constructor(options2 = {}) {
    const devOptions = this.isDevMode ? {
      pollingInterval: 3e4
      // 30 seconds in dev (was 10s)
    } : {};
    this.options = { ...DEFAULT_OPTIONS, ...devOptions, ...options2 };
  }
  /**
   * Get the singleton instance of MemoryMonitor
   */
  static getInstance(options2) {
    if (!MemoryMonitor.instance) {
      MemoryMonitor.instance = new MemoryMonitor(options2);
    } else if (options2) {
      MemoryMonitor.instance.updateOptions(options2);
    }
    return MemoryMonitor.instance;
  }
  /**
   * Update monitoring options
   */
  updateOptions(options2) {
    this.options = { ...this.options, ...options2 };
    if (this.monitorInterval) {
      this.stopMonitoring();
      this.startMonitoring();
    }
  }
  /**
   * Get current memory stats with store metrics
   */
  async getMemoryStats() {
    if (!("performance" in window) || !performance.memory) {
      return null;
    }
    const memory = performance.memory;
    const used = Math.round(memory.usedJSHeapSize / (1024 * 1024));
    const total = Math.round(memory.jsHeapSizeLimit / (1024 * 1024));
    const percentage = Math.round(memory.usedJSHeapSize / memory.jsHeapSizeLimit * 100);
    const storeMetrics = await this.collectStoreMetrics();
    return { used, total, percentage, storeMetrics };
  }
  /**
   * Collect memory metrics from all stores
   * In dev mode, skip detailed metrics to reduce memory overhead
   */
  async collectStoreMetrics() {
    const metrics = {
      fileStore: { nodesCount: 0, memoryEstimate: 0 },
      contextStore: { cacheSize: 0, chunkSize: 0 },
      apiCache: { entries: 0, totalSize: 0 }
    };
    if (this.isDevMode) {
      return metrics;
    }
    try {
      if (!this.storeImportsCache.useFileStore) {
        const fileStoreModule = await __vitePreload(() => Promise.resolve().then(() => file_store), true ? void 0 : void 0);
        this.storeImportsCache.useFileStore = fileStoreModule.useFileStore;
      }
      if (!this.storeImportsCache.useContextStore) {
        const contextStoreModule = await __vitePreload(() => Promise.resolve().then(() => context_store), true ? void 0 : void 0);
        this.storeImportsCache.useContextStore = contextStoreModule.useContextStore;
      }
      if (!this.storeImportsCache.getCacheStats) {
        const apiCacheModule = await __vitePreload(() => import("./useApiCache-BTgc1mkV.js"), true ? __vite__mapDeps([11,1,2,5,6,7,3,4]) : void 0);
        this.storeImportsCache.getCacheStats = apiCacheModule.getCacheStats;
      }
      const fileStore = this.storeImportsCache.useFileStore?.();
      const contextStore = this.storeImportsCache.useContextStore?.();
      const cacheStats = this.storeImportsCache.getCacheStats?.();
      if (fileStore) {
        metrics.fileStore = {
          nodesCount: fileStore.nodes?.length || 0,
          memoryEstimate: fileStore.getMemoryUsage ? fileStore.getMemoryUsage() : 0
        };
      }
      if (contextStore) {
        metrics.contextStore = {
          cacheSize: contextStore.getMemoryUsage ? contextStore.getMemoryUsage() : 0,
          chunkSize: contextStore.currentChunk?.lines?.length || 0
        };
      }
      if (cacheStats) {
        metrics.apiCache = {
          entries: cacheStats.entries || 0,
          totalSize: cacheStats.size || 0
        };
      }
    } catch (e) {
      console.warn("[MemoryMonitor] Could not collect store metrics:", e);
    }
    return metrics;
  }
  /**
   * Start memory monitoring
   */
  startMonitoring() {
    if (this.monitorInterval) {
      return;
    }
    void this.checkMemory();
    this.monitorInterval = window.setInterval(
      () => void this.checkMemory(),
      this.options.pollingInterval
    );
    logger$1.debug(`Memory monitoring started (interval: ${this.options.pollingInterval}ms)`);
  }
  /**
   * Stop memory monitoring
   */
  stopMonitoring() {
    if (this.monitorInterval) {
      clearInterval(this.monitorInterval);
      this.monitorInterval = null;
    }
  }
  /**
   * Force cleanup of memory with AGGRESSIVE strategies
   */
  forceCleanup() {
    logger$1.warn("EMERGENCY: Forcing aggressive memory cleanup...");
    try {
      window.largeObjects = [];
      window.cachedData = null;
      __vitePreload(async () => {
        const { clearAllCaches } = await import("./useApiCache-BTgc1mkV.js");
        return { clearAllCaches };
      }, true ? __vite__mapDeps([11,1,2,5,6,7,3,4]) : void 0).then(({ clearAllCaches }) => {
        clearAllCaches();
      }).catch((e) => {
        console.warn("Could not clear API caches:", e);
      });
      Promise.all([
        __vitePreload(() => Promise.resolve().then(() => file_store), true ? void 0 : void 0),
        __vitePreload(() => Promise.resolve().then(() => context_store), true ? void 0 : void 0)
      ]).then(([{ useFileStore: useFileStore2 }, { useContextStore: useContextStore2 }]) => {
        const fileStore = useFileStore2();
        const contextStore = useContextStore2();
        fileStore.resetStore();
        contextStore.clearContext();
      }).catch((e) => {
        console.warn("Could not cleanup stores:", e);
      });
      try {
        const stores = ["contextBuilder", "fileTree", "treeState", "ui"];
        stores.forEach((storeName) => {
          const storeData = window[storeName];
          if (storeData && typeof storeData.cleanup === "function") {
            storeData.cleanup();
          }
        });
      } catch (e) {
        console.warn("Could not cleanup Vue stores:", e);
      }
      if (window.gc) {
        try {
          window.gc();
          setTimeout(() => window.gc?.(), 100);
          setTimeout(() => window.gc?.(), 300);
          logger$1.debug("Multiple garbage collection cycles triggered");
        } catch (e) {
          console.warn("Failed to trigger garbage collection", e);
        }
      }
      setTimeout(() => {
        void this.getMemoryStats().then((stats) => {
          if (stats) {
            logger$1.debug(`Memory after aggressive cleanup: ${stats.used}MB / ${stats.total}MB (${stats.percentage}%)`);
          }
        });
      }, 500);
    } catch (e) {
      logger$1.error("Error during emergency memory cleanup:", e);
    }
  }
  /**
   * Dump heap snapshot for analysis in Chrome DevTools
   */
  async dumpHeapSnapshot(reason) {
    if ("performance" in window && performance.writeHeapSnapshot) {
      try {
        const timestamp = (/* @__PURE__ */ new Date()).toISOString().replace(/[:.]/g, "-");
        const filename = `heap-snapshot-${timestamp}-${reason}.heapsnapshot`;
        await performance.writeHeapSnapshot(filename);
        logger$1.debug(`Heap snapshot saved: ${filename}`);
        this.uiStore.addToast(`Heap snapshot saved: ${filename}`, "info");
      } catch (e) {
        logger$1.error("Failed to dump heap snapshot:", e);
      }
    } else {
      console.warn("Heap snapshot API not available. Run Chrome with --enable-precise-memory-info flag.");
    }
  }
  /**
   * Register callback for critical memory events
   * Returns unsubscribe function to prevent memory leaks
   */
  criticalCallbacks = /* @__PURE__ */ new Set();
  onCritical(callback) {
    this.criticalCallbacks.add(callback);
    return () => {
      this.criticalCallbacks.delete(callback);
    };
  }
  /**
   * Clear all critical callbacks
   */
  clearCriticalCallbacks() {
    this.criticalCallbacks.clear();
  }
  /**
   * Cleanup and reset the monitor (call when unmounting)
   */
  cleanup() {
    this.stopMonitoring();
    this.criticalCallbacks.clear();
    this.memoryHistory = [];
    this.warningIssued = false;
    this.criticalIssued = false;
  }
  /**
   * Check current memory usage with PERCENTAGE-BASED thresholds (more reliable)
   */
  memoryHistory = [];
  async checkMemory() {
    const stats = await this.getMemoryStats();
    if (!stats) {
      console.warn("[MemoryMonitor] Memory stats unavailable - performance.memory API not supported in this browser");
      return;
    }
    const { used, total, percentage } = stats;
    if (!isFinite(percentage) || percentage < 0) {
      return;
    }
    const warningPercentage = 60;
    const criticalPercentage = 80;
    if (used < 500) {
      return;
    }
    const memoryTrend = this.getMemoryTrend(used);
    if (percentage >= criticalPercentage) {
      if (!this.criticalIssued) {
        logger$1.error(`CRITICAL: Memory usage at ${percentage}% (${used}MB / ${total}MB)`);
        if (this.options.showToasts) {
          const locale = localStorage.getItem("app-locale") || "ru";
          const message = locale === "ru" ? `КРИТИЧЕСКОЕ использование памяти: ${used}МБ / ${total}МБ (${percentage}%). Браузер может зависнуть! Очистка...` : `CRITICAL memory usage: ${used}MB / ${total}MB (${percentage}%). Browser may crash! Cleaning up...`;
          this.uiStore.addToast(message, "error", 15e3);
        }
        this.dumpHeapSnapshot("critical-memory");
        if (this.options.autoCleanup) {
          this.forceCleanup();
          setTimeout(() => this.forceCleanup(), 500);
        }
        for (const cb of this.criticalCallbacks) {
          try {
            cb();
          } catch (e) {
            logger$1.error("Critical callback error:", e);
          }
        }
        this.criticalIssued = true;
      }
    } else if (percentage >= warningPercentage || memoryTrend > 10 && this.memoryHistory.length >= 5) {
      if (!this.warningIssued) {
        console.warn(`WARNING: Memory usage at ${percentage}% (${used}MB / ${total}MB)`);
        if (this.options.showToasts) {
          const locale = localStorage.getItem("app-locale") || "ru";
          const message = locale === "ru" ? `Использование памяти: ${used}МБ / ${total}МБ (${percentage}% от heap limit). Рекомендуется уменьшить выбор файлов.` : `Memory usage: ${used}MB / ${total}MB (${percentage}% of heap limit). Consider reducing file selection.`;
          this.uiStore.addToast(message, "warning", 8e3);
        }
        this.warningIssued = true;
      }
    } else {
      if (percentage < warningPercentage - 10) {
        this.warningIssued = false;
        this.criticalIssued = false;
      }
    }
  }
  /**
   * Track memory trend to detect rapid growth
   */
  getMemoryTrend(currentUsed) {
    this.memoryHistory.push(currentUsed);
    if (this.memoryHistory.length > 10) {
      this.memoryHistory.shift();
    }
    if (this.memoryHistory.length < 2) {
      return 0;
    }
    const oldest = this.memoryHistory[0];
    const newest = this.memoryHistory[this.memoryHistory.length - 1];
    const timeDiff = (this.memoryHistory.length - 1) * (this.options.pollingInterval || 5e3) / 6e4;
    return (newest - oldest) / timeDiff;
  }
}
function useMemoryMonitor$1(options2) {
  return MemoryMonitor.getInstance(options2);
}
function useMemoryMonitor(options2) {
  return MemoryMonitor.getInstance(options2);
}
const ONBOARDING_KEY = "Syntaxia-onboarding-completed";
function useOnboarding() {
  const { t } = useI18n();
  async function startTour() {
    const { driver } = await __vitePreload(async () => {
      const { driver: driver2 } = await import("./driver.js-CQwR5pvx.js");
      return { driver: driver2 };
    }, true ? [] : void 0);
    await __vitePreload(() => Promise.resolve({}), true ? __vite__mapDeps([12]) : void 0);
    await __vitePreload(() => Promise.resolve({}), true ? __vite__mapDeps([13]) : void 0);
    const steps = [
      {
        element: '[data-tour="file-tree"]',
        popover: {
          title: t("onboarding.step1Title"),
          description: t("onboarding.step1Desc"),
          side: "right"
        }
      },
      {
        element: '[data-tour="build-button"]',
        popover: {
          title: t("onboarding.step2Title"),
          description: t("onboarding.step2Desc"),
          side: "top"
        }
      },
      {
        element: '[data-tour="context-preview"]',
        popover: {
          title: t("onboarding.step3Title"),
          description: t("onboarding.step3Desc"),
          side: "left"
        }
      },
      {
        element: '[data-tour="ai-chat"]',
        popover: {
          title: t("onboarding.step4Title"),
          description: t("onboarding.step4Desc"),
          side: "left"
        }
      }
    ];
    const driverObj = driver({
      showProgress: true,
      steps,
      nextBtnText: t("onboarding.next"),
      prevBtnText: t("onboarding.prev"),
      doneBtnText: t("onboarding.done"),
      onDestroyStarted: () => {
        localStorage.setItem(ONBOARDING_KEY, "true");
        driverObj.destroy();
      }
    });
    driverObj.drive();
  }
  function shouldShowTour() {
    return !localStorage.getItem(ONBOARDING_KEY);
  }
  function resetTour() {
    localStorage.removeItem(ONBOARDING_KEY);
  }
  function markTourCompleted() {
    localStorage.setItem(ONBOARDING_KEY, "true");
  }
  return { startTour, shouldShowTour, resetTour, markTourCompleted };
}
const shellApi = {
  getStatus: () => apiCall(
    () => GetShellIntegrationStatus(),
    "Failed to get shell integration status.",
    { logContext: "shell" }
  ),
  register: () => apiCall(
    () => RegisterShellIntegration(),
    "Failed to register shell integration.",
    { logContext: "shell" }
  ),
  unregister: () => apiCall(
    () => UnregisterShellIntegration(),
    "Failed to unregister shell integration.",
    { logContext: "shell" }
  ),
  getStartupPath: () => apiCall(
    () => GetStartupPath(),
    "Failed to get startup path.",
    { logContext: "shell" }
  ),
  clearStartupPath: () => apiCall(
    () => ClearStartupPath(),
    "Failed to clear startup path.",
    { logContext: "shell" }
  )
};
const _hoisted_1 = {
  id: "app",
  class: "app-container layout-root"
};
const _hoisted_2 = {
  href: "#main-content",
  class: "skip-link"
};
const _hoisted_3 = {
  key: 0,
  class: "modal-container"
};
const _hoisted_4 = { class: "modal-content modal-error" };
const _hoisted_5 = { class: "modal-title" };
const _hoisted_6 = { class: "modal-text" };
const _hoisted_7 = {
  key: 1,
  class: "loading-overlay"
};
const _hoisted_8 = { class: "loading-card" };
const _hoisted_9 = { class: "loading-title" };
const _hoisted_10 = { class: "loading-text" };
const _hoisted_11 = {
  id: "main-content",
  class: "layout-fill layout-column layout-clip"
};
const _hoisted_12 = { class: "fixed bottom-3 left-1/2 -translate-x-1/2 z-[100] flex flex-col items-center gap-2" };
const _hoisted_13 = { class: "toast-glass-icon-wrap" };
const _hoisted_14 = {
  key: 0,
  class: "toast-glass-icon",
  fill: "none",
  stroke: "currentColor",
  viewBox: "0 0 24 24"
};
const _hoisted_15 = {
  key: 1,
  class: "toast-glass-icon",
  fill: "none",
  stroke: "currentColor",
  viewBox: "0 0 24 24"
};
const _hoisted_16 = {
  key: 2,
  class: "toast-glass-icon",
  fill: "none",
  stroke: "currentColor",
  viewBox: "0 0 24 24"
};
const _hoisted_17 = {
  key: 3,
  class: "toast-glass-icon",
  fill: "none",
  stroke: "currentColor",
  viewBox: "0 0 24 24"
};
const _hoisted_18 = { class: "toast-glass-message" };
const _hoisted_19 = ["onClick"];
const _hoisted_20 = ["onClick"];
const _sfc_main = /* @__PURE__ */ defineComponent({
  __name: "App",
  setup(__props) {
    const CommandPalette = defineAsyncComponent(() => __vitePreload(() => import("./CommandPalette-Be47sfdv.js"), true ? __vite__mapDeps([14,1,2,5,6,7,3,4]) : void 0));
    const KeyboardShortcutsModal = defineAsyncComponent(() => __vitePreload(() => import("./KeyboardShortcutsModal-CJZUYkBp.js"), true ? __vite__mapDeps([15,5,1,2]) : void 0));
    const MemoryDashboard = defineAsyncComponent(() => __vitePreload(() => import("./MemoryDashboard-BykbtL8H.js"), true ? __vite__mapDeps([16,1,2,5,6,7,3,4]) : void 0));
    const SettingsModal = defineAsyncComponent(() => __vitePreload(() => import("./SettingsModal-D_zm6NLM.js"), true ? __vite__mapDeps([17,1,2,7,5,6,3,4,18]) : void 0));
    const ConfirmDialog = defineAsyncComponent(() => __vitePreload(() => import("./ConfirmDialog-DL6BPwsV.js"), true ? __vite__mapDeps([19,1,2,7,5,6,3,4]) : void 0));
    const { t } = useI18n();
    const projectStore = useProjectStore();
    const uiStore = useUIStore();
    const globalError = ref(null);
    const isCommandPaletteOpen = ref(false);
    const isShortcutsModalOpen = ref(false);
    const showMemoryDashboard = ref(false);
    const keys = useMagicKeys();
    const ctrlK = keys["Ctrl+K"];
    const ctrlP = keys["Ctrl+P"];
    const ctrlSlash = keys["Ctrl+/"];
    const ctrlB = keys["Ctrl+B"];
    const ctrlEnter = keys["Ctrl+Enter"];
    const ctrlE = keys["Ctrl+E"];
    const ctrlShiftC = keys["Ctrl+Shift+C"];
    const ctrlShiftM = keys["Ctrl+Shift+M"];
    const ctrlComma = keys["Ctrl+,"];
    watch(ctrlK, (v) => {
      if (v) isCommandPaletteOpen.value = !isCommandPaletteOpen.value;
    });
    watch(ctrlP, (v) => {
      if (v) {
        isCommandPaletteOpen.value = true;
      }
    });
    watch(ctrlSlash, (v) => {
      if (v) isShortcutsModalOpen.value = !isShortcutsModalOpen.value;
    });
    watch([ctrlB, ctrlEnter], ([b, enter]) => {
      if ((b || enter) && projectStore.hasProject) {
        const event = new CustomEvent("global-build-context");
        window.dispatchEvent(event);
      }
    });
    watch(ctrlE, (v) => {
      if (v && projectStore.hasProject) {
        const event = new CustomEvent("global-open-export");
        window.dispatchEvent(event);
      }
    });
    watch(ctrlShiftC, (v) => {
      if (v && projectStore.hasProject) {
        const event = new CustomEvent("global-copy-context");
        window.dispatchEvent(event);
      }
    });
    const ctrlZ = keys["Ctrl+Z"];
    const ctrlY = keys["Ctrl+Y"];
    watch(ctrlZ, (v) => {
      if (v && projectStore.hasProject) {
        const activeEl = document.activeElement;
        if (activeEl?.tagName !== "INPUT" && activeEl?.tagName !== "TEXTAREA") {
          const event = new CustomEvent("global-undo-selection");
          window.dispatchEvent(event);
        }
      }
    });
    watch(ctrlY, (v) => {
      if (v && projectStore.hasProject) {
        const activeEl = document.activeElement;
        if (activeEl?.tagName !== "INPUT" && activeEl?.tagName !== "TEXTAREA") {
          const event = new CustomEvent("global-redo-selection");
          window.dispatchEvent(event);
        }
      }
    });
    watch(ctrlShiftM, (v) => {
      if (v && projectStore.hasProject) {
        showMemoryDashboard.value = !showMemoryDashboard.value;
      }
    });
    watch(ctrlComma, (v) => {
      if (v) {
        uiStore.openSettingsModal();
      }
    });
    const memoryMonitor = useMemoryMonitor();
    const onboarding2 = useOnboarding();
    function clearGlobalError() {
      globalError.value = null;
    }
    function onProjectOpened(_path) {
      uiStore.addToast("Project loaded successfully", "success");
      if (onboarding2.shouldShowTour()) {
        setTimeout(() => {
          onboarding2.startTour();
        }, 500);
      }
    }
    const cleanupIntervalRef = ref(null);
    onMounted(async () => {
      memoryMonitor.startMonitoring();
      memoryMonitor.onCritical(() => {
        showMemoryDashboard.value = true;
      });
      try {
        const startupPath = await shellApi.getStartupPath();
        if (startupPath) {
          await shellApi.clearStartupPath();
          await projectStore.openProjectByPath(startupPath);
          uiStore.addToast("Project opened from context menu", "success");
          return;
        }
      } catch {
      }
      projectStore.maybeAutoOpenLastProject();
      {
        cleanupIntervalRef.value = window.setInterval(async () => {
          try {
            const { clearAllCaches, getCacheStats } = await __vitePreload(async () => {
              const { clearAllCaches: clearAllCaches2, getCacheStats: getCacheStats2 } = await import("./useApiCache-BTgc1mkV.js");
              return { clearAllCaches: clearAllCaches2, getCacheStats: getCacheStats2 };
            }, true ? __vite__mapDeps([11,1,2,5,6,7,3,4]) : void 0);
            const stats = getCacheStats();
            if (stats.size > stats.maxSize * 0.5) {
              clearAllCaches();
            }
          } catch {
          }
          await memoryMonitor.getMemoryStats();
        }, 10 * 60 * 1e3);
      }
    });
    onUnmounted(() => {
      if (cleanupIntervalRef.value) {
        clearInterval(cleanupIntervalRef.value);
        cleanupIntervalRef.value = null;
      }
      memoryMonitor.stopMonitoring();
      __vitePreload(async () => {
        const { clearAllCaches } = await import("./useApiCache-BTgc1mkV.js");
        return { clearAllCaches };
      }, true ? __vite__mapDeps([11,1,2,5,6,7,3,4]) : void 0).then(({ clearAllCaches }) => {
        clearAllCaches();
      }).catch(() => {
      });
    });
    onMounted(() => {
      window.addEventListener("error", (event) => {
        if (event.message && event.message.includes("out of memory")) {
          const { useLogger: useLogger2 } = require("@/composables/useLogger");
          const memLogger = useLogger2("App");
          memLogger.error("Out of memory error detected!");
          uiStore.addToast("Critical memory error. Clearing caches...", "error");
          try {
            const { clearAllCaches } = require("@/composables/useApiCache");
            clearAllCaches();
            const { useContextStore: useContextStore2 } = require("@/features/context/model/context.store");
            const { useFileStore: useFileStore2 } = require("@/features/files/model/file.store");
            const contextStore = useContextStore2();
            const fileStore = useFileStore2();
            contextStore.clearContext();
            fileStore.resetStore();
          } catch (e) {
            memLogger.error("Emergency cleanup failed:", e);
          }
        }
      });
    });
    return (_ctx, _cache) => {
      return openBlock(), createElementBlock("div", _hoisted_1, [
        createBaseVNode("a", _hoisted_2, toDisplayString(unref(t)("accessibility.skipToContent")), 1),
        globalError.value ? (openBlock(), createElementBlock("div", _hoisted_3, [
          createBaseVNode("div", {
            class: "modal-overlay",
            onClick: clearGlobalError
          }),
          createBaseVNode("div", _hoisted_4, [
            createBaseVNode("h3", _hoisted_5, toDisplayString(unref(t)("common.criticalError")), 1),
            createBaseVNode("p", _hoisted_6, toDisplayString(globalError.value), 1),
            createBaseVNode("button", {
              onClick: clearGlobalError,
              class: "btn btn-danger"
            }, toDisplayString(unref(t)("common.close")), 1)
          ])
        ])) : createCommentVNode("", true),
        unref(projectStore).isLoading ? (openBlock(), createElementBlock("div", _hoisted_7, [
          createBaseVNode("div", _hoisted_8, [
            _cache[4] || (_cache[4] = createBaseVNode("div", { class: "flex items-center justify-center mb-4" }, [
              createBaseVNode("svg", {
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
              ])
            ], -1)),
            createBaseVNode("h3", _hoisted_9, toDisplayString(unref(t)("common.loadingProject")), 1),
            createBaseVNode("p", _hoisted_10, toDisplayString(unref(t)("common.pleaseWait")), 1)
          ])
        ])) : createCommentVNode("", true),
        createBaseVNode("main", _hoisted_11, [
          !unref(projectStore).hasProject ? (openBlock(), createBlock(ProjectSelector, {
            key: 0,
            onOpened: onProjectOpened
          })) : (openBlock(), createBlock(MainWorkspace, { key: 1 }))
        ]),
        createBaseVNode("div", _hoisted_12, [
          createVNode(TransitionGroup, { name: "toast" }, {
            default: withCtx(() => [
              (openBlock(true), createElementBlock(Fragment, null, renderList(unref(uiStore).toasts, (toast) => {
                return openBlock(), createElementBlock("div", {
                  key: toast.id,
                  class: normalizeClass([
                    "toast-glass",
                    toast.type === "error" ? "toast-glass-error" : toast.type === "success" ? "toast-glass-success" : toast.type === "warning" ? "toast-glass-warning" : "toast-glass-info"
                  ])
                }, [
                  createBaseVNode("div", _hoisted_13, [
                    toast.type === "success" ? (openBlock(), createElementBlock("svg", _hoisted_14, [..._cache[5] || (_cache[5] = [
                      createBaseVNode("path", {
                        "stroke-linecap": "round",
                        "stroke-linejoin": "round",
                        "stroke-width": "2.5",
                        d: "M5 13l4 4L19 7"
                      }, null, -1)
                    ])])) : toast.type === "error" ? (openBlock(), createElementBlock("svg", _hoisted_15, [..._cache[6] || (_cache[6] = [
                      createBaseVNode("path", {
                        "stroke-linecap": "round",
                        "stroke-linejoin": "round",
                        "stroke-width": "2.5",
                        d: "M6 18L18 6M6 6l12 12"
                      }, null, -1)
                    ])])) : toast.type === "warning" ? (openBlock(), createElementBlock("svg", _hoisted_16, [..._cache[7] || (_cache[7] = [
                      createBaseVNode("path", {
                        "stroke-linecap": "round",
                        "stroke-linejoin": "round",
                        "stroke-width": "2.5",
                        d: "M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"
                      }, null, -1)
                    ])])) : (openBlock(), createElementBlock("svg", _hoisted_17, [..._cache[8] || (_cache[8] = [
                      createBaseVNode("path", {
                        "stroke-linecap": "round",
                        "stroke-linejoin": "round",
                        "stroke-width": "2.5",
                        d: "M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
                      }, null, -1)
                    ])]))
                  ]),
                  createBaseVNode("span", _hoisted_18, toDisplayString(toast.message), 1),
                  toast.action ? (openBlock(), createElementBlock(Fragment, { key: 0 }, [
                    _cache[10] || (_cache[10] = createBaseVNode("div", { class: "toast-glass-divider" }, null, -1)),
                    createBaseVNode("button", {
                      onClick: ($event) => {
                        toast.action.onClick();
                        unref(uiStore).removeToast(toast.id);
                      },
                      class: "toast-glass-action"
                    }, [
                      _cache[9] || (_cache[9] = createBaseVNode("svg", {
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
                      createTextVNode(" " + toDisplayString(toast.action.label), 1)
                    ], 8, _hoisted_19)
                  ], 64)) : createCommentVNode("", true),
                  createBaseVNode("button", {
                    onClick: ($event) => unref(uiStore).removeToast(toast.id),
                    class: "toast-glass-close"
                  }, [..._cache[11] || (_cache[11] = [
                    createBaseVNode("svg", {
                      class: "w-3.5 h-3.5",
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
                  ])], 8, _hoisted_20)
                ], 2);
              }), 128))
            ]),
            _: 1
          })
        ]),
        unref(projectStore).hasProject ? (openBlock(), createBlock(unref(CommandPalette), {
          key: 2,
          "is-open": isCommandPaletteOpen.value,
          onClose: _cache[0] || (_cache[0] = ($event) => isCommandPaletteOpen.value = false)
        }, null, 8, ["is-open"])) : createCommentVNode("", true),
        unref(projectStore).hasProject ? (openBlock(), createBlock(unref(KeyboardShortcutsModal), {
          key: 3,
          "is-open": isShortcutsModalOpen.value,
          onClose: _cache[1] || (_cache[1] = ($event) => isShortcutsModalOpen.value = false)
        }, null, 8, ["is-open"])) : createCommentVNode("", true),
        showMemoryDashboard.value && unref(projectStore).hasProject ? (openBlock(), createBlock(unref(MemoryDashboard), {
          key: 4,
          class: "fixed top-20 right-4 z-40",
          onClose: _cache[2] || (_cache[2] = ($event) => showMemoryDashboard.value = false)
        })) : createCommentVNode("", true),
        createVNode(unref(SettingsModal), {
          modelValue: unref(uiStore).showSettingsModal,
          "onUpdate:modelValue": _cache[3] || (_cache[3] = ($event) => unref(uiStore).showSettingsModal = $event)
        }, null, 8, ["modelValue"]),
        createVNode(unref(ConfirmDialog))
      ]);
    };
  }
});
const logger = useLogger("ErrorHandler");
function setupErrorHandler(app2, options2 = {}) {
  const {
    showNotification = true,
    logToConsole = true
  } = options2;
  app2.config.errorHandler = (err, instance, info) => {
    if (logToConsole) {
      logger.error("Vue error:", err);
      logger.error("Component:", instance);
      logger.error("Error info:", info);
    }
    if (showNotification) {
      try {
        const uiStore = useUIStore();
        const message = err instanceof Error ? err.message : String(err);
        uiStore.addToast(`Application error: ${message}`, "error");
      } catch (toastError) {
        logger.error("Failed to show error toast:", toastError);
        {
          alert(`An error occurred: ${err instanceof Error ? err.message : String(err)}`);
        }
      }
    }
  };
  window.addEventListener("unhandledrejection", (event) => {
    if (logToConsole) {
      logger.error("Unhandled Promise Rejection:", event.reason);
    }
    if (showNotification) {
      try {
        const uiStore = useUIStore();
        const message = event.reason instanceof Error ? event.reason.message : String(event.reason);
        uiStore.addToast(`Unhandled promise rejection: ${message}`, "error");
      } catch (toastError) {
        logger.error("Failed to show rejection toast:", toastError);
      }
    }
    event.preventDefault();
  });
}
const app = createApp(_sfc_main);
const pinia = createPinia();
app.use(pinia);
app.use(autoAnimatePlugin);
setupErrorHandler(app, {
  showNotification: true,
  logToConsole: true
});
app.mount("#app");
export {
  AISettings as A,
  ExportSettings as E,
  __vitePreload as _,
  useLogger as a,
  useMemoryMonitor$1 as b,
  useI18n as c,
  useSettingsStore as d,
  useOnboarding as e,
  _export_sfc as f,
  useProjectStore as g,
  apiService as h,
  getFileIcon as i,
  context_store as j,
  file_store as k,
  parseIgnoreRules as p,
  shellApi as s,
  useUIStore as u
};
