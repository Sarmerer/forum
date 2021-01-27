import Vue from "vue";
import VueRouter from "vue-router";
import Home from "@/views/Home.vue";
import store from "@/store";

Vue.use(VueRouter);

const ROLES = {};
ROLES.ADMIN = 2;

const routes = [
  {
    path: "/",
    name: "Home",
    component: Home,
  },
  {
    path: "*",
    name: "Not Found",
    component: () =>
      import(/* webpackChunkName: "not-found" */ "@/components/NotFound.vue"),
  },
  {
    path: "/servers-down",
    name: "Oh no...",
    component: () =>
      import(/* webpackChunkName: "servers-down" */ "@/views/ServersDown.vue"),
  },
  {
    path: "/post/:id",
    name: "Post",
    props: true,
    meta: {
      titleSetInComponent: true,
    },
    // route level code-splitting
    // this generates a separate chunk (about.[hash].js) for this route
    // which is lazy-loaded when the route is visited.
    component: () => import(/* webpackChunkName: "post" */ "@/views/Post.vue"),
  },
  {
    path: "/new-post",
    name: "New Post",
    props: true,
    component: () =>
      import(/* webpackChunkName: "new-post" */ "@/views/NewPost.vue"),
    beforeEnter(to, from, next) {
      if (!store.getters["auth/authenticated"]) {
        next({ name: "Auth", params: { prevRoute: "/new-post" } });
      } else {
        next();
      }
    },
  },
  {
    path: "/user/:userID",
    name: "User",
    meta: {
      titleSetInComponent: true,
    },
    component: () => import(/* webpackChunkName: "user" */ "@/views/User.vue"),
  },
  {
    path: "/auth",
    name: "Auth",
    props: true,
    component: () => import(/* webpackChunkName: "auth" */ "@/views/Auth.vue"),
    beforeEnter(to, from, next) {
      store.getters["auth/authenticated"] ? next("/") : next();
    },
  },
  {
    path: "/auth/verify",
    name: "AuthMerge",
    props: true,
    component: () =>
      import(/* webpackChunkName: "authMerge" */ "@/views/AuthVerify.vue"),
    beforeEnter(to, from, next) {
      store.getters["auth/authenticated"] ? next("/") : next();
    },
  },
  {
    path: "/auth/merge",
    name: "AuthMerge",
    props: true,
    component: () =>
      import(/* webpackChunkName: "authMerge" */ "@/views/AuthMerge.vue"),
    beforeEnter(to, from, next) {
      store.getters["auth/authenticated"] ? next("/") : next();
    },
  },
  {
    path: "/dashboard/:role",
    name: "Admin Dashboard",
    component: () =>
      import(/* webpackChunkName: "dashboard" */ "@/views/Dashboard.vue"),
    beforeEnter: (to, from, next) => {
      let user = store.getters["auth/user"];
      !store.getters["auth/authenticated"] || !user || user.role < ROLES.ADMIN
        ? next("/")
        : next();
    },
  },
];

const router = new VueRouter({
  mode: "history",
  base: process.env.BASE_URL,
  routes,
});

router.beforeEach((to, _from, next) => {
  if (!to.meta?.titleSetInComponent) document.title = to.name;
  next();
});

export default router;
