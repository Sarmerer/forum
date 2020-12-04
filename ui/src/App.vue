<template>
  <div id="app">
    <div class="hero-image" v-if="!isMobile()">
      <div class="hero-text">
        <h1 v-if="!isMobile()">
          WELCOME<br />
          TO <span class="primary">FORUM</span>
        </h1>
        <h3 v-if="isMobile()">
          WELCOME<br />
          TO <span class="primary">FORUM</span>
        </h3>
      </div>
    </div>

    <div class="card">
      <SearchBar v-on:toggle-search="search = !search" />
    </div>
    <Navbar v-if="isMobile()" id="nav" />
    <Sidenav v-else />
    <router-view :class="`d-${search ? 'none' : 'block'}`" />
    <Footer />
  </div>
</template>
<script>
import Sidenav from "@/components/Sidenav";
import Navbar from "@/components/Navbar";
import Footer from "@/components/Footer";
import SearchBar from "@/components/SearchBar";

export default {
  name: "App",
  components: {
    SearchBar,
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
