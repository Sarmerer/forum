import "@babel/polyfill";
import "mutationobserver-shim";
import Vue from "vue";
import "./plugins/bootstrap-vue";
import App from "./App.vue";
import router from "./router";
import store from "./store";
import axios from "axios";
import "./assets/styles.scss";

import VueSidebarMenu from "vue-sidebar-menu";
import "vue-sidebar-menu/dist/vue-sidebar-menu.css";

Vue.use(VueSidebarMenu);

Vue.config.productionTip = false;
axios.defaults.baseURL = "/api/";
axios.defaults.withCredinentials = true;
store.dispatch("auth/attempt", localStorage.getItem("ilgn")).then(() => {
  new Vue({
    router,
    store,
    render: (h) => h(App),
  }).$mount("#app");
});
