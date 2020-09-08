<template>
  <div>
    <div class="card" v-for="post in posts" :key="post.post.ID" :id="post.post.ID">
      <h3 class="primary">{{ post.post.Title }}</h3>
      <hr />
      <p>{{ post.post.Content }}</p>
    </div>
  </div>
</template>
<script>
import axios from "axios";
export default {
  data: function() {
    return {
      posts: null,
    };
  },
  created() {
    axios.get("posts").then((response) => (this.posts = response.data.data));
  },
  methods: {},
};
</script>
