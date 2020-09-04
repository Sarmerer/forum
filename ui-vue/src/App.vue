<template>
  <div id="app">
    <div id="nav">
      <b-nav>
        <b-nav-item> <router-link to="/">Home</router-link></b-nav-item>
        <Login class="ml-auto"></Login>
      </b-nav>
    </div>
    <router-view />
  </div>
</template>
<script>
import Login from "../src/components/Login";

export default {
  name: "App",
  components: {
    Login,
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

#nav {
  padding: 30px;
}

#nav a {
  font-weight: bold;
  color: #2c3e50;
}

#nav a.router-link-exact-active {
  color: #42b983;
}
</style>
