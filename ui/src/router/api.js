import axios from "axios";

const api = axios.create({
  baseURL: process.env.VUE_APP_API_URL,
  withCredentials: true,
});

api.interceptors.response.use(
  (res) => {
    return res;
  },
  (error) => {
    console.log(error);
    return error;
  }
);

export default api;
