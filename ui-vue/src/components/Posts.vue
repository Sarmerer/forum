<template>
  <div>
    <ul>
      <li v-for="(Title, Content) in posts" :key="Title">{{ Title }} | {{ Content }}</li>
    </ul>
    <!-- <pre>{{ posts }}</pre> -->
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
