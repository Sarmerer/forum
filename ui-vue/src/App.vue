<template>
  <div id="app">
    <div id="nav">
      <b-navbar toggleable="lg">
        <b-navbar-brand href="#">forum</b-navbar-brand>
        <b-navbar-toggle target="nav-collapse"></b-navbar-toggle>
        <b-collapse id="nav-collapse" is-nav>
          <b-navbar-nav>
            <b-nav-item> <router-link to="/">Home</router-link></b-nav-item>
          </b-navbar-nav>
          <b-navbar-nav v-if="!authenticated" class="ml-auto">
            <Login class="ml-auto"></Login>
          </b-navbar-nav>
          <b-navbar-nav v-if="authenticated" class="ml-auto">
            <Me />
            <b-nav-item @click.prevent="signOut" href="#">Logout</b-nav-item>
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
import { mapActions } from "vuex";
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
  methods: {
    ...mapActions({
      signOut: "auth/signOut",
    }),
    submit() {
      this.signOut();
    },
  },
};
</script>
<style>
@import url("https://fonts.googleapis.com/css?family=Roboto");

#app {
  font-family: "Roboto", sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
}

#nav a {
  font-weight: bold;
  color: #2c3e50;
}

#nav a.router-link-exact-active {
  color: #42b983;
}
</style>
