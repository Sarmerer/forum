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
      import(/* webpackChunkName: "Not Found" */ "@/views/NotFound.vue"),
  },
  {
    path: "/post/:id",
    name: "Post",
    // route level code-splitting
    // this generates a separate chunk (about.[hash].js) for this route
    // which is lazy-loaded when the route is visited.
    component: () => import(/* webpackChunkName: "post" */ "@/views/Post.vue"),
  },
  {
    path: "/new-post",
    name: "New Post",
    component: () =>
      import(/* webpackChunkName: "user" */ "@/views/NewPost.vue"),
  },
  {
    path: "/user/:id",
    name: "User",
    component: () => import(/* webpackChunkName: "user" */ "@/views/User.vue"),
  },
  {
    path: "/dashboard/:role",
    name: "Admin Dashnoard",
    component: () =>
      import(/* webpackChunkName: "user" */ "@/views/Dashboard.vue"),
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

export default router;
