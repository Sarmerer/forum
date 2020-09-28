<template>
  <div id="app">
    <sidebar-menu disableHover collapsed :menu="menu" />
    <div id="nav">
      <b-navbar toggleable="lg" type="light" style="background: #278ea5">
        <b-navbar-brand href="#">forum</b-navbar-brand>
        <b-navbar-toggle target="nav-collapse"></b-navbar-toggle>
        <b-collapse id="nav-collapse" is-nav>
          <b-navbar-nav>
            <b-nav-item><router-link to="/">Home</router-link></b-nav-item>
          </b-navbar-nav>
          <b-navbar-nav>
            <b-nav-item><router-link to="/post/new">New post</router-link></b-nav-item>
          </b-navbar-nav>
          <b-navbar-nav v-if="!authenticated" class="ml-auto">
            <Login class="ml-auto"></Login>
          </b-navbar-nav>
          <b-navbar-nav v-else class="ml-auto">
            <Me />
          </b-navbar-nav>
        </b-collapse>
      </b-navbar>
    </div>
    <router-view />
  </div>
</template>
<script>
import Login from "../src/components/Login";
import Me from "../src/components/Me";
import { mapGetters } from "vuex";

export default {
  name: "App",
  computed: {
    ...mapGetters({
      authenticated: "auth/authenticated",
    }),
  },
  components: {
    Login,
    Me,
  },
  data() {
    return {
      menu: [
        {
          href: "/",
          title: "Home",
          icon: "fa fa-home",
        },
        {
          href: "/",
          title: "Profile",
          icon: "fa fa-user",
        },
        {
          href: "/",
          title: "New post",
          icon: "fa fa-pencil-alt",
        },
        {
          href: "/",
          title: "Sign out",
          icon: {
            element: "span",
            class: "fa fa-sign-out-alt",
            attributes: {
              style: {
                color: "red",
              },
            },
            // text: ''
          },
        },
      ],
    };
  },
};
</script>
