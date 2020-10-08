<template>
  <div class="error-container">
    <h3 class="mb-4">{{ errorData.status }} {{ errorData.message }}</h3>
    <!-- TODO add overlay to prevent request spam -->
    <b-button @click="errorData.callback()" variant="outline-light">retry</b-button>
  </div>
</template>
<script>
export default {
  props: {
    errorData: { type: Object, required: true },
  },
};
</script>
<style lang="scss">
.error-container {
  margin: 20px;
  padding: 10px;
  text-align: center;
}
</style>
