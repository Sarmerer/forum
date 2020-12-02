import axios from "axios";

const api = axios.create({
  baseURL: process.env.VUE_APP_API_URL,
  withCredentials: true,
});

api.interceptors.response.use(
  function(res) {
    return res;
  },
  function(error) {
    return Promise.reject(error.response);
  }
);

export default api;
