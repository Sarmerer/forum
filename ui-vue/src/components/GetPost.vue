<template>
  <div class="columns">
    <div class="info-col">
      <div class="card">
        <h3 class="primary">AUTHOR</h3>
        <p>Author info</p>
      </div>
    </div>
    <div class="post-col">
      <div class="card">
        <h3 class="primary">{{ post.Title }}</h3>
        <pre style="color: white">{{ post.Content }}</pre>
        <div>
          <b-form-tags
            v-model="categories"
            tag-variant="primary"
            tag-pills
            size="sm"
            disabled
            separator=" "
          ></b-form-tags>
        </div>
      </div>
    </div>
  </div>
</template>
<script>
import axios from "axios";
export default {
  data() {
    return {
      post: {},
      categories: [],
    };
  },
  mounted() {
    this.getPost();
  },
  methods: {
    async getPost() {
      return await axios
        .post("post/find", {
          by: "id",
          id: Number.parseInt(this.$route.params.id),
        })
        .then((response) => {
          response.data.data.forEach((res) => {
            this.post = res.post;
            res.categories.forEach((c) => {
              this.categories.push(c.name);
            });
          });
        })
        .catch(() => {
          this.$router.push("/");
        });
    },
  },
};
</script>
