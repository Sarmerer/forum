<template>
  <div id="app">
    <div v-if="isMobile()" id="nav">
      <Navbar />
    </div>
    <div v-else>
      <Sidenav />
    </div>
    <router-view />
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
};
</script>
