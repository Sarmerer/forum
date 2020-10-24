<template>
  <div id="app">
    <!-- <div id="nav">
      <b-navbar
        toggleable="lg"
        type="light"
        style="background: rgba(255, 255, 255, 0.05)"
      >
        <b-navbar-brand href="#">forum</b-navbar-brand>
        <b-navbar-toggle target="nav-collapse"></b-navbar-toggle>
        <b-collapse id="nav-collapse" is-nav>
          <b-navbar-nav>
            <b-nav-item><router-link to="/">Home</router-link></b-nav-item>
          </b-navbar-nav>
          <b-navbar-nav>
            <b-nav-item
              ><router-link to="/post/new">New post</router-link></b-nav-item
            >
          </b-navbar-nav>
          <b-navbar-nav v-if="!authenticated" class="ml-auto">
            <Login class="ml-auto"></Login>
          </b-navbar-nav>
          <b-navbar-nav v-else class="ml-auto">
            <Me />
          </b-navbar-nav>
        </b-collapse>
      </b-navbar>
    </div> -->
    <div class="sidenav">
      <router-link to="/"
        ><span class="primary"><b-icon icon="chat-left-dots"></b-icon></span
      ></router-link>
      <router-link to="/post/new"><b-icon icon="pen"></b-icon></router-link>
      <div v-if="!authenticated">
        <Login />
      </div>
      <div v-if="authenticated">
        <Me />
      </div>
      <!-- <a href="#about"><b-icon icon="person"></b-icon></a>
      <a href="#about"><b-icon icon="gear"></b-icon></a>
      <router-link to="/post/new"
        ><b-icon icon="door-open"></b-icon
      ></router-link> -->
    </div>
    <router-view class="main" />
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
      authenticated: "auth/authenticated"
    })
  },
  components: {
    Login,
    Me
  }
};
</script>
<style lang="scss">
.sidenav {
  height: 100%;
  width: 65px;
  position: fixed;
  z-index: 1;
  top: 0;
  left: 0;
  background-color: #111;
  overflow-x: hidden;
  padding-top: 30px;
}

.sidenav a {
  padding: 10px 8px 10px 16px;
  text-decoration: none;
  font-size: 25px;
  color: #818181;
  display: block;
}

.sidenav a:hover {
  color: #f1f1f1;
}

.main {
  margin-left: 65px;
  margin-top: 30px;
}
</style>
