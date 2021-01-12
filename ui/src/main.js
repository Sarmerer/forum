import "@babel/polyfill";
import "mutationobserver-shim";
import Vue from "vue";
import "./plugins/bootstrap-vue";
import App from "./App.vue";
import router from "./router";
import store from "./store";
import "./assets/styles.scss"; // Global styles
import { BootstrapVue, BootstrapVueIcons } from "bootstrap-vue";
import api from "@/api/api";

Vue.use(BootstrapVue);
Vue.use(BootstrapVueIcons);

Vue.config.productionTip = false;

store.dispatch("auth/attempt").then(() => {
  new Vue({
    router,
    api,
    store,
    render: (h) => h(App),
  }).$mount("#app");
});

Vue.mixin({
  methods: {
    isMobile: function() {
      if (
        /Android|webOS|iPhone|iPad|iPod|BlackBerry|IEMobile|Opera Mini/i.test(
          navigator.userAgent
        )
      ) {
        return true;
      } else {
        //Margin-left is equal to the width of the .sidenav element
        document.body.style.marginLeft = "65px";
        return false;
      }
    },
  },
});
