<template>
  <div class="error-container">
    <h3 class="mb-4">{{ status }} {{ message }}</h3>
    <!-- TODO add overlay to prevent request spam -->
    <b-button @click="callback()" variant="outline-light">retry</b-button>
  </div>
</template>
<script>
export default {
  props: ["message", "status", "callback"],
};
</script>
<style lang="scss">
.error-container {
  margin: 20px;
  padding: 10px;
  text-align: center;
}
</style>
