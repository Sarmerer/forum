//FIXME error page is scrollable
<template>
  <div class="grid error">
    <h2>Oops! Unknown page.</h2>
    <h5>The page you are looking for doesn't seem to exist...</h5>
    <h5>Instead, kindly enjoy this exceptionally beautiful masterpiece.</h5>
    <img src="../assets/img/masterpiece.png" alt="masterpiece" />
    <b-button variant="outline-light" @click="redirect">Take me home</b-button>
  </div>
</template>
<script>
export default {
  methods: {
    redirect() {
      this.$router.push("/");
    },
  },
};
</script>
<style scoped>
.error {
  text-align: center;
  margin-top: 30px;
}
img {
  width: 100%;
  height: auto;
}
</style>
