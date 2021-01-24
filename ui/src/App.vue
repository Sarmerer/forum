<template>
  <div id="app">
    <Navbar v-if="isMobile()" id="nav" />
    <Sidenav v-else />
    <keep-alive include="Home,UserPage">
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
  created() {
    if (!this.isMobile()) document.body.style.marginLeft = "65px";
  },
};
</script>
