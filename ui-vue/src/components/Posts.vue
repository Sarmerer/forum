<template>
  <div>
    //FIXME post gets weird actual ID
    <div class="card" v-for="post in posts" :key="post.id" :id="post.id">
      <h3 class="primary">
        <router-link to="/post/2">{{ post.post.Title }}</router-link>
      </h3>
      <hr />
      <p>{{ post.post.Content }}</p>
      <p>ID in the list: {{ post.id }}</p>
      <p>ID in the database: {{ post.post.ID }}</p>
      <button @click="deletePost(post.id, post.post.ID)">Delete</button>
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
  methods: {
    deletePost(IDInTheList, actualID) {
      axios
        .delete("post/delete", { params: { ID: actualID } })
        .then((response) => {
          console.log(response.data);
          this.$delete(this.posts, IDInTheList);
        })
        .catch((error) => {
          alert(error.response.data.code + " " + error.response.data.message);
        });
    },
  },
};
</script>
