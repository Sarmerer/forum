import api from "@/router/api";

export default {
  namespaced: true,
  state: {
    user: null,
  },
  getters: {
    authenticated(state) {
      return state.user != null;
    },
    user(state) {
      return state.user;
    },
  },
  mutations: {
    setUser(state, user) {
      state.user = user;
    },
  },
  actions: {
    async signUp({ dispatch }, credentials) {
      await api
        .post(`auth/signup`, credentials)
        .then(() => {
          dispatch("attempt");
        })
        .catch((error) => console.log(error));
    },
    async signIn({ dispatch }, credentials) {
      await api
        .post(`auth/signin`, credentials)
        .then(() => {
          dispatch("attempt");
        })
        .catch((error) => console.log(error));
    },
    async signOut({ commit }) {
      //FIXME fix case when you log in from a new place, and then try to log out from the first session
      return api.post("auth/signout").finally(() => {
        commit("setUser", null);
      });
    },
    async attempt({ commit }) {
      return await api
        .get("auth/me")
        .then((response) => commit("setUser", response.data.data))
        .catch(() => commit("setUser", null));
    },
  },
};
