import { createApp } from "vue";
import { createPinia } from "pinia";
import { autoAnimatePlugin } from "@formkit/auto-animate/vue";

import App from "./App.vue";
import "./assets/main.css";
import "./assets/design-tokens.css";
import { setupErrorHandler } from "./plugins/errorHandler";
import 'highlight.js/styles/github-dark.css';

const app = createApp(App);
const pinia = createPinia();

app.use(pinia);
app.use(autoAnimatePlugin);

// Setup global error handler
setupErrorHandler(app, {
  showNotification: true,
  logToConsole: true
});

app.mount("#app");