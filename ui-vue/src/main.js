import "@babel/polyfill";
import "mutationobserver-shim";
import Vue from "vue";
import "./plugins/bootstrap-vue";
import App from "./App.vue";
import router from "./router";
import store from "./store";
import axios from "axios";
import "./assets/styles.scss";
import vuetify from "./plugins/vuetify";
Vue.config.productionTip = false;
axios.defaults.baseURL = "/api/";
axios.defaults.withCredinentials = true;
store.dispatch("auth/attempt", localStorage.getItem("ilgn")).then(() => {
  new Vue({
    router,
    store,
    vuetify,
    render: h => h(App)
  }).$mount("#app");
});
