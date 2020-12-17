<template>
  <div id="app">
    <Navbar v-if="isMobile()" id="nav" />
    <Sidenav v-else />
    <keep-alive include="Home">
      <router-view :class="`d-${search ? 'none' : 'block'}`" />
    </keep-alive>
    <Footer />
  </div>
</template>
<script>
import Sidenav from "@/components/Sidenav";
import Navbar from "@/components/Navbar";
import Footer from "@/components/Footer";

export default {
  name: "App",
  components: {
    Sidenav,
    Navbar,
    Footer,
  },
  data() {
    return {
      search: false,
    };
  },
};
</script>
