import axios from "axios";
import store from "@/store";

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
    }
    return Promise.reject(error.response);
  }
);

export default api;
