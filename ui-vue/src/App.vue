<template>
  <div id="app">
    <div v-if="isMobile()" id="nav">
      <Navbar />
    </div>
    <div v-else>
      <Sidenav />
    </div>
    <router-view />
  </div>
</template>
<script>
import Sidenav from "../src/components/Sidenav";
import Navbar from "../src/components/Navbar";

export default {
  name: "App",
  components: {
    Sidenav,
    Navbar
  }
};
</script>
