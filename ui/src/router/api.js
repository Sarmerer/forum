import axios from "axios";
import store from "@/store";
import router from "@/router";
let errorsCounter = 0;

const api = axios.create({
  baseURL: process.env.VUE_APP_API_URL,
  withCredentials: true,
});

api.interceptors.response.use(
  function(res) {
    return res;
  },
  function(error) {
    if (error.response.status === 403) {
      store.commit("auth/setUser", null);
    } else if (error.response.status === 500) {
      errorsCounter++;
      if (errorsCounter > 2) {
        errorsCounter = 0;
        router.push("/servers-down");
      }
    }
    return Promise.reject(error.response);
  }
);

export default api;
