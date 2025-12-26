import { w as watch, f as onMounted, n as nextTick, m as readonly, j as getCurrentInstance, q as toRef$1, r as ref, v as customRef, k as reactive, c as computed, x as toValue, s as shallowRef, u as unref } from "./vue-vendor-DLLjc1XI.js";
const isClient = typeof window !== "undefined" && typeof document !== "undefined";
typeof WorkerGlobalScope !== "undefined" && globalThis instanceof WorkerGlobalScope;
const toString = Object.prototype.toString;
const isObject = (val) => toString.call(val) === "[object Object]";
const noop = () => {
};
const isIOS = /* @__PURE__ */ getIsIOS();
function getIsIOS() {
  var _window, _window2, _window3;
  return isClient && !!((_window = window) === null || _window === void 0 || (_window = _window.navigator) === null || _window === void 0 ? void 0 : _window.userAgent) && (/iP(?:ad|hone|od)/.test(window.navigator.userAgent) || ((_window2 = window) === null || _window2 === void 0 || (_window2 = _window2.navigator) === null || _window2 === void 0 ? void 0 : _window2.maxTouchPoints) > 2 && /iPad|Macintosh/.test((_window3 = window) === null || _window3 === void 0 ? void 0 : _window3.navigator.userAgent));
}
function toRef(...args) {
  if (args.length !== 1) return toRef$1(...args);
  const r = args[0];
  return typeof r === "function" ? readonly(customRef(() => ({
    get: r,
    set: noop
  }))) : ref(r);
}
function createFilterWrapper(filter, fn) {
  function wrapper(...args) {
    return new Promise((resolve, reject) => {
      Promise.resolve(filter(() => fn.apply(this, args), {
        fn,
        thisArg: this,
        args
      })).then(resolve).catch(reject);
    });
  }
  return wrapper;
}
const bypassFilter = (invoke$1) => {
  return invoke$1();
};
function pausableFilter(extendFilter = bypassFilter, options = {}) {
  const { initialState = "active" } = options;
  const isActive = toRef(initialState === "active");
  function pause() {
    isActive.value = false;
  }
  function resume() {
    isActive.value = true;
  }
  const eventFilter = (...args) => {
    if (isActive.value) extendFilter(...args);
  };
  return {
    isActive: readonly(isActive),
    pause,
    resume,
    eventFilter
  };
}
function toArray(value) {
  return Array.isArray(value) ? value : [value];
}
function getLifeCycleTarget(target) {
  return getCurrentInstance();
}
function watchWithFilter(source, cb, options = {}) {
  const { eventFilter = bypassFilter, ...watchOptions } = options;
  return watch(source, createFilterWrapper(eventFilter, cb), watchOptions);
}
function watchPausable(source, cb, options = {}) {
  const { eventFilter: filter, initialState = "active", ...watchOptions } = options;
  const { eventFilter, pause, resume, isActive } = pausableFilter(filter, { initialState });
  return {
    stop: watchWithFilter(source, cb, {
      ...watchOptions,
      eventFilter
    }),
    pause,
    resume,
    isActive
  };
}
const pausableWatch = watchPausable;
function tryOnMounted(fn, sync = true, target) {
  if (getLifeCycleTarget()) onMounted(fn, target);
  else if (sync) fn();
  else nextTick(fn);
}
function watchImmediate(source, cb, options) {
  return watch(source, cb, {
    ...options,
    immediate: true
  });
}
const defaultWindow = isClient ? window : void 0;
function unrefElement(elRef) {
  var _$el;
  const plain = toValue(elRef);
  return (_$el = plain === null || plain === void 0 ? void 0 : plain.$el) !== null && _$el !== void 0 ? _$el : plain;
}
function useEventListener(...args) {
  const register = (el, event, listener, options) => {
    el.addEventListener(event, listener, options);
    return () => el.removeEventListener(event, listener, options);
  };
  const firstParamTargets = computed(() => {
    const test = toArray(toValue(args[0])).filter((e) => e != null);
    return test.every((e) => typeof e !== "string") ? test : void 0;
  });
  return watchImmediate(() => {
    var _firstParamTargets$va, _firstParamTargets$va2;
    return [
      (_firstParamTargets$va = (_firstParamTargets$va2 = firstParamTargets.value) === null || _firstParamTargets$va2 === void 0 ? void 0 : _firstParamTargets$va2.map((e) => unrefElement(e))) !== null && _firstParamTargets$va !== void 0 ? _firstParamTargets$va : [defaultWindow].filter((e) => e != null),
      toArray(toValue(firstParamTargets.value ? args[1] : args[0])),
      toArray(unref(firstParamTargets.value ? args[2] : args[1])),
      toValue(firstParamTargets.value ? args[3] : args[2])
    ];
  }, ([raw_targets, raw_events, raw_listeners, raw_options], _, onCleanup) => {
    if (!(raw_targets === null || raw_targets === void 0 ? void 0 : raw_targets.length) || !(raw_events === null || raw_events === void 0 ? void 0 : raw_events.length) || !(raw_listeners === null || raw_listeners === void 0 ? void 0 : raw_listeners.length)) return;
    const optionsClone = isObject(raw_options) ? { ...raw_options } : raw_options;
    const cleanups = raw_targets.flatMap((el) => raw_events.flatMap((event) => raw_listeners.map((listener) => register(el, event, listener, optionsClone))));
    onCleanup(() => {
      cleanups.forEach((fn) => fn());
    });
  }, { flush: "post" });
}
let _iOSWorkaround = false;
function onClickOutside(target, handler, options = {}) {
  const { window: window$1 = defaultWindow, ignore = [], capture = true, detectIframe = false, controls = false } = options;
  if (!window$1) return controls ? {
    stop: noop,
    cancel: noop,
    trigger: noop
  } : noop;
  if (isIOS && !_iOSWorkaround) {
    _iOSWorkaround = true;
    const listenerOptions = { passive: true };
    Array.from(window$1.document.body.children).forEach((el) => el.addEventListener("click", noop, listenerOptions));
    window$1.document.documentElement.addEventListener("click", noop, listenerOptions);
  }
  let shouldListen = true;
  const shouldIgnore = (event) => {
    return toValue(ignore).some((target$1) => {
      if (typeof target$1 === "string") return Array.from(window$1.document.querySelectorAll(target$1)).some((el) => el === event.target || event.composedPath().includes(el));
      else {
        const el = unrefElement(target$1);
        return el && (event.target === el || event.composedPath().includes(el));
      }
    });
  };
  function hasMultipleRoots(target$1) {
    const vm = toValue(target$1);
    return vm && vm.$.subTree.shapeFlag === 16;
  }
  function checkMultipleRoots(target$1, event) {
    const vm = toValue(target$1);
    const children = vm.$.subTree && vm.$.subTree.children;
    if (children == null || !Array.isArray(children)) return false;
    return children.some((child) => child.el === event.target || event.composedPath().includes(child.el));
  }
  const listener = (event) => {
    const el = unrefElement(target);
    if (event.target == null) return;
    if (!(el instanceof Element) && hasMultipleRoots(target) && checkMultipleRoots(target, event)) return;
    if (!el || el === event.target || event.composedPath().includes(el)) return;
    if ("detail" in event && event.detail === 0) shouldListen = !shouldIgnore(event);
    if (!shouldListen) {
      shouldListen = true;
      return;
    }
    handler(event);
  };
  let isProcessingClick = false;
  const cleanup = [
    useEventListener(window$1, "click", (event) => {
      if (!isProcessingClick) {
        isProcessingClick = true;
        setTimeout(() => {
          isProcessingClick = false;
        }, 0);
        listener(event);
      }
    }, {
      passive: true,
      capture
    }),
    useEventListener(window$1, "pointerdown", (e) => {
      const el = unrefElement(target);
      shouldListen = !shouldIgnore(e) && !!(el && !e.composedPath().includes(el));
    }, { passive: true }),
    detectIframe && useEventListener(window$1, "blur", (event) => {
      setTimeout(() => {
        var _window$document$acti;
        const el = unrefElement(target);
        if (((_window$document$acti = window$1.document.activeElement) === null || _window$document$acti === void 0 ? void 0 : _window$document$acti.tagName) === "IFRAME" && !(el === null || el === void 0 ? void 0 : el.contains(window$1.document.activeElement))) handler(event);
      }, 0);
    }, { passive: true })
  ].filter(Boolean);
  const stop = () => cleanup.forEach((fn) => fn());
  if (controls) return {
    stop,
    cancel: () => {
      shouldListen = false;
    },
    trigger: (event) => {
      shouldListen = true;
      listener(event);
      shouldListen = false;
    }
  };
  return stop;
}
const _global = typeof globalThis !== "undefined" ? globalThis : typeof window !== "undefined" ? window : typeof global !== "undefined" ? global : typeof self !== "undefined" ? self : {};
const globalKey = "__vueuse_ssr_handlers__";
const handlers = /* @__PURE__ */ getHandlers();
function getHandlers() {
  if (!(globalKey in _global)) _global[globalKey] = _global[globalKey] || {};
  return _global[globalKey];
}
function getSSRHandler(key, fallback) {
  return handlers[key] || fallback;
}
function guessSerializerType(rawInit) {
  return rawInit == null ? "any" : rawInit instanceof Set ? "set" : rawInit instanceof Map ? "map" : rawInit instanceof Date ? "date" : typeof rawInit === "boolean" ? "boolean" : typeof rawInit === "string" ? "string" : typeof rawInit === "object" ? "object" : !Number.isNaN(rawInit) ? "number" : "any";
}
const StorageSerializers = {
  boolean: {
    read: (v) => v === "true",
    write: (v) => String(v)
  },
  object: {
    read: (v) => JSON.parse(v),
    write: (v) => JSON.stringify(v)
  },
  number: {
    read: (v) => Number.parseFloat(v),
    write: (v) => String(v)
  },
  any: {
    read: (v) => v,
    write: (v) => String(v)
  },
  string: {
    read: (v) => v,
    write: (v) => String(v)
  },
  map: {
    read: (v) => new Map(JSON.parse(v)),
    write: (v) => JSON.stringify(Array.from(v.entries()))
  },
  set: {
    read: (v) => new Set(JSON.parse(v)),
    write: (v) => JSON.stringify(Array.from(v))
  },
  date: {
    read: (v) => new Date(v),
    write: (v) => v.toISOString()
  }
};
const customStorageEventName = "vueuse-storage";
function useStorage(key, defaults$1, storage, options = {}) {
  var _options$serializer;
  const { flush = "pre", deep = true, listenToStorageChanges = true, writeDefaults = true, mergeDefaults = false, shallow, window: window$1 = defaultWindow, eventFilter, onError = (e) => {
    console.error(e);
  }, initOnMounted } = options;
  const data = (shallow ? shallowRef : ref)(typeof defaults$1 === "function" ? defaults$1() : defaults$1);
  const keyComputed = computed(() => toValue(key));
  if (!storage) try {
    storage = getSSRHandler("getDefaultStorage", () => defaultWindow === null || defaultWindow === void 0 ? void 0 : defaultWindow.localStorage)();
  } catch (e) {
    onError(e);
  }
  if (!storage) return data;
  const rawInit = toValue(defaults$1);
  const type = guessSerializerType(rawInit);
  const serializer = (_options$serializer = options.serializer) !== null && _options$serializer !== void 0 ? _options$serializer : StorageSerializers[type];
  const { pause: pauseWatch, resume: resumeWatch } = pausableWatch(data, (newValue) => write(newValue), {
    flush,
    deep,
    eventFilter
  });
  watch(keyComputed, () => update(), { flush });
  let firstMounted = false;
  const onStorageEvent = (ev) => {
    if (initOnMounted && !firstMounted) return;
    update(ev);
  };
  const onStorageCustomEvent = (ev) => {
    if (initOnMounted && !firstMounted) return;
    updateFromCustomEvent(ev);
  };
  if (window$1 && listenToStorageChanges) if (storage instanceof Storage) useEventListener(window$1, "storage", onStorageEvent, { passive: true });
  else useEventListener(window$1, customStorageEventName, onStorageCustomEvent);
  if (initOnMounted) tryOnMounted(() => {
    firstMounted = true;
    update();
  });
  else update();
  function dispatchWriteEvent(oldValue, newValue) {
    if (window$1) {
      const payload = {
        key: keyComputed.value,
        oldValue,
        newValue,
        storageArea: storage
      };
      window$1.dispatchEvent(storage instanceof Storage ? new StorageEvent("storage", payload) : new CustomEvent(customStorageEventName, { detail: payload }));
    }
  }
  function write(v) {
    try {
      const oldValue = storage.getItem(keyComputed.value);
      if (v == null) {
        dispatchWriteEvent(oldValue, null);
        storage.removeItem(keyComputed.value);
      } else {
        const serialized = serializer.write(v);
        if (oldValue !== serialized) {
          storage.setItem(keyComputed.value, serialized);
          dispatchWriteEvent(oldValue, serialized);
        }
      }
    } catch (e) {
      onError(e);
    }
  }
  function read(event) {
    const rawValue = event ? event.newValue : storage.getItem(keyComputed.value);
    if (rawValue == null) {
      if (writeDefaults && rawInit != null) storage.setItem(keyComputed.value, serializer.write(rawInit));
      return rawInit;
    } else if (!event && mergeDefaults) {
      const value = serializer.read(rawValue);
      if (typeof mergeDefaults === "function") return mergeDefaults(value, rawInit);
      else if (type === "object" && !Array.isArray(value)) return {
        ...rawInit,
        ...value
      };
      return value;
    } else if (typeof rawValue !== "string") return rawValue;
    else return serializer.read(rawValue);
  }
  function update(event) {
    if (event && event.storageArea !== storage) return;
    if (event && event.key == null) {
      data.value = rawInit;
      return;
    }
    if (event && event.key !== keyComputed.value) return;
    pauseWatch();
    try {
      const serializedData = serializer.write(data.value);
      if (event === void 0 || (event === null || event === void 0 ? void 0 : event.newValue) !== serializedData) data.value = read(event);
    } catch (e) {
      onError(e);
    } finally {
      if (event) nextTick(resumeWatch);
      else resumeWatch();
    }
  }
  function updateFromCustomEvent(event) {
    update(event.detail);
  }
  return data;
}
const DefaultMagicKeysAliasMap = {
  ctrl: "control",
  command: "meta",
  cmd: "meta",
  option: "alt",
  up: "arrowup",
  down: "arrowdown",
  left: "arrowleft",
  right: "arrowright"
};
function useMagicKeys(options = {}) {
  const { reactive: useReactive = false, target = defaultWindow, aliasMap = DefaultMagicKeysAliasMap, passive = true, onEventFired = noop } = options;
  const current = reactive(/* @__PURE__ */ new Set());
  const obj = {
    toJSON() {
      return {};
    },
    current
  };
  const refs = useReactive ? reactive(obj) : obj;
  const metaDeps = /* @__PURE__ */ new Set();
  const depsMap = /* @__PURE__ */ new Map([
    ["Meta", metaDeps],
    ["Shift", /* @__PURE__ */ new Set()],
    ["Alt", /* @__PURE__ */ new Set()]
  ]);
  const usedKeys = /* @__PURE__ */ new Set();
  function setRefs(key, value) {
    if (key in refs) if (useReactive) refs[key] = value;
    else refs[key].value = value;
  }
  function reset() {
    current.clear();
    for (const key of usedKeys) setRefs(key, false);
  }
  function updateDeps(value, e, keys$1) {
    if (!value || typeof e.getModifierState !== "function") return;
    for (const [modifier, depsSet] of depsMap) if (e.getModifierState(modifier)) {
      keys$1.forEach((key) => depsSet.add(key));
      break;
    }
  }
  function clearDeps(value, key) {
    if (value) return;
    const depsMapKey = `${key[0].toUpperCase()}${key.slice(1)}`;
    const deps = depsMap.get(depsMapKey);
    if (!["shift", "alt"].includes(key) || !deps) return;
    const depsArray = Array.from(deps);
    const depsIndex = depsArray.indexOf(key);
    depsArray.forEach((key$1, index) => {
      if (index >= depsIndex) {
        current.delete(key$1);
        setRefs(key$1, false);
      }
    });
    deps.clear();
  }
  function updateRefs(e, value) {
    var _e$key, _e$code;
    const key = (_e$key = e.key) === null || _e$key === void 0 ? void 0 : _e$key.toLowerCase();
    const values = [(_e$code = e.code) === null || _e$code === void 0 ? void 0 : _e$code.toLowerCase(), key].filter(Boolean);
    if (key === "") return;
    if (key) if (value) current.add(key);
    else current.delete(key);
    for (const key$1 of values) {
      usedKeys.add(key$1);
      setRefs(key$1, value);
    }
    updateDeps(value, e, [...current, ...values]);
    clearDeps(value, key);
    if (key === "meta" && !value) {
      metaDeps.forEach((key$1) => {
        current.delete(key$1);
        setRefs(key$1, false);
      });
      metaDeps.clear();
    }
  }
  useEventListener(target, "keydown", (e) => {
    updateRefs(e, true);
    return onEventFired(e);
  }, { passive });
  useEventListener(target, "keyup", (e) => {
    updateRefs(e, false);
    return onEventFired(e);
  }, { passive });
  useEventListener("blur", reset, { passive });
  useEventListener("focus", reset, { passive });
  const proxy = new Proxy(refs, { get(target$1, prop, rec) {
    if (typeof prop !== "string") return Reflect.get(target$1, prop, rec);
    prop = prop.toLowerCase();
    if (prop in aliasMap) prop = aliasMap[prop];
    if (!(prop in refs)) if (/[+_-]/.test(prop)) {
      const keys$1 = prop.split(/[+_-]/g).map((i) => i.trim());
      refs[prop] = computed(() => keys$1.map((key) => toValue(proxy[key])).every(Boolean));
    } else refs[prop] = shallowRef(false);
    const r = Reflect.get(target$1, prop, rec);
    return useReactive ? toValue(r) : r;
  } });
  return proxy;
}
export {
  useMagicKeys as a,
  onClickOutside as o,
  useStorage as u
};
