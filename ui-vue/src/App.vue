<template>
  <div id="app">
    <div v-if="isMobile()" id="nav">
      <b-navbar type="dark" fixed="bottom">
        <b-navbar-nav
          ><b-nav-item>
            <router-link to="/"
              ><span class="primary"
                ><b-icon
                  icon="chat-left-dots"
                ></b-icon></span></router-link></b-nav-item
        ></b-navbar-nav>
        <b-navbar-nav
          ><b-nav-item>
            <router-link to="/post/new"
              ><b-icon icon="pen"></b-icon></router-link></b-nav-item
        ></b-navbar-nav>
        <b-navbar-nav v-if="!authenticated"
          ><b-nav-item>
            <div>
              <Login /></div></b-nav-item
        ></b-navbar-nav>
        <b-navbar-nav v-if="authenticated"
          ><b-nav-item>
            <div>
              <Me /></div></b-nav-item
        ></b-navbar-nav>
      </b-navbar>
    </div>
    <div v-if="!isMobile()" class="sidenav">
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
    </div>
    <!-- why the hack is this not working??? -->
    <router-view class="{'main': !isMobile}" />
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
  },
  methods: {
    isMobile() {
      if (
        /Android|webOS|iPhone|iPad|iPod|BlackBerry|IEMobile|Opera Mini/i.test(
          navigator.userAgent
        )
      ) {
        return true;
      } else {
        return false;
      }
    }
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

.navbar {
  background-color: #111;
  justify-content: space-around;
}
// .navbar > .navbar-nav > .nav-link {
//   color: #f1f1f1;
// }

.main {
  margin-left: 65px;
}
</style>
