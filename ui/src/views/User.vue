<template>
  <div class="grid">
    <div class="user-card">
      <img :src="user.avatar" alt="avatar" />
      <h3 class="primary">
        {{ user.display_name }}
        <b-badge v-if="user.role == 2" class="background-variant"
          >Admin</b-badge
        >
      </h3>
      <p>Last online: {{ user.last_online | formatDate }}</p>
    </div>
    <div class="user-info">
      <b-tabs pills card>
        <b-tab title="Posts" active>
          <router-link
            :to="'/post/' + post.id"
            v-for="post in posts"
            :key="post.id"
            class="user-card"
            tag="div"
          >
            <h5>
              <strong>{{ post.title }}</strong>
            </h5>
            {{ post.content }}
          </router-link></b-tab
        >
        <b-tab title="Replies"
          ><div
            class="user-card"
            v-for="(reply, index) in replies"
            :key="index"
          >
            <h5>
              {{ reply.content }}
            </h5>
          </div></b-tab
        >
      </b-tabs>
    </div>
  </div>
</template>
<script>
import axios from "axios";
import moment from "moment";

export default {
  data() {
    return {
      user: {},
      posts: [],
      replies: [],
    };
  },
  mounted() {
    this.getUser();
    this.getUserPosts();
    this.getReplies();
  },
  filters: {
    formatDate: function (value) {
      if (value) {
        return moment(String(value)).format("MM/DD/YYYY hh:mm");
      }
    },
  },
  methods: {
    async getUser() {
      return await axios
        .get("user", { params: { id: this.$route.params.id } })
        .then((response) => {
          this.user = response.data.data;
        })
        .catch((error) => {
          console.log(error);
        });
    },
    async getUserPosts() {
      return await axios
        .post("post/find", {
          by: "author",
          author: Number.parseInt(this.$route.params.id),
        })
        .then((response) => {
          this.posts = response.data.data || [];
          console.log(response);
        })
        .catch((error) => {
          console.log(error);
        });
    },
    async getReplies() {
      return await axios
        .post("comments/find", {
          by: "user_id",
          id: Number.parseInt(this.$route.params.id),
        })
        .then((response) => {
          this.replies = response.data.data || [];
          console.log(response.data.data);
        })
        .catch((error) => {
          console.log(error);
        });
    },
  },
};
</script>
<style lang="scss">
.user-card {
  margin: 15px 0;
  padding: 15px;
  background-color: rgba(255, 255, 255, 0.05);
  box-shadow: 5px 5px 6px 2px rgba(10, 10, 10, 0.3);
}

.user-info .user-card {
  cursor: pointer;
}

.user-info .user-card:hover {
  opacity: 0.8;
}

.user-card img {
  width: 100px;
  height: 100px;
  border-radius: 200px;
}

.card-body {
  padding: 0;
}

.nav-pills .nav-link.active,
.nav-pills .show > .nav-link {
  background-color: #278ea5;
}

.nav-link {
  color: #278ea5;
}

.nav-link:hover {
  color: #278ea5;
  opacity: 0.8;
}
</style>
