import Vue from "vue";
import App from "./App.vue";
// import "./registerServiceWorker";
import router from "./router";
import store from "./store";

// Bug with UIKit. Imports dont work with TypeScripts
// ISSUE: https://github.com/uikit/uikit/issues/3606
// FIX/Hack: https://stackoverflow.com/questions/41292559/could-not-find-a-declaration-file-for-module-module-name-path-to-module-nam
// const UIKit = require("uikit");
// const Icons = require("uikit/dist/js/uikit-icons");
// import * as UIkit from "uikit";
// import * as Icons from "uikit/dist/js/uikit-icons";

// UIKit.use(Icons);

Vue.config.productionTip = false;

new Vue({
  router,
  store,
  render: h => h(App)
}).$mount("#app");
