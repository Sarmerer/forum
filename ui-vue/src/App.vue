<template>
  <div id="app">
    <div id="nav">
      <b-nav>
        <b-nav-item> <router-link to="/">Home</router-link> |</b-nav-item>
        <b-nav-item><router-link to="/login">Login</router-link></b-nav-item>
      </b-nav>
    </div>
    <router-view />
  </div>
</template>

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
