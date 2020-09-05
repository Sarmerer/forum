<template>
  <b-list-group-item href="#" disabled class="flex-column">
    <div class="d-flex w-100 justify-content-between">
      <h5 class="mb-1">{{ post.Title }}</h5>
      <small>{{ post.Created }}</small>
    </div>

    <p class="mb-1">{{ post.Content }}</p>

    <small>Author: {{ post.Author }}</small>
  </b-list-group-item>
</template>
<script>
export default {
  props: ["Post"],
  data() {
    return {
      post: this.Post,
    };
  },
};
</script>
