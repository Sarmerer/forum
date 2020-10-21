import "@babel/polyfill";
import "mutationobserver-shim";
import Vue from "vue";
import "./plugins/bootstrap-vue";
import App from "./App.vue";
import router from "./router";
import store from "./store";
import axios from "axios";
import "./assets/styles.scss"; // Global styles

import VueTimeago from "vue-timeago";
Vue.use(VueTimeago, {
  name: "Timeago", // Component name, `Timeago` by default
  locale: "en", // Default locale
});

Vue.config.productionTip = false;
axios.defaults.baseURL = "https://forum-api-sarmerer.herokuapp.com/api/";
axios.defaults.withCredinentials = true;
store.dispatch("auth/attempt").then(() => {
  new Vue({
    router,
    store,
    render: (h) => h(App),
  }).$mount("#app");
});
