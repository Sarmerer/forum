<template>
  <div>
    <b-list-group style="align-items: center; justify-content: center">
      <Card style="width: 800px;" v-for="post in posts" :key="post.post.ID" :id="post.post.ID" v-bind:Post="post.post" />
    </b-list-group>
  </div>
</template>
<script>
import axios from "axios";
import Card from "@/components/Card";
export default {
  components: {
    Card,
  },
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
