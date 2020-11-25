<template>
  <div class="grid">
    <div class="columns">
      <div class="info-col">
        <div class="card user-card">
          <img :src="user.avatar" alt="avatar" />
          <h3 class="primary">
            {{ user.display_name }}
            <b-badge v-if="user.role == 2" class="background-variant"
              >Admin</b-badge
            >
          </h3>
        </div>
      </div>
      <div class="main-col">
        <div class="user-info">
          <b-tabs card>
            <b-tab title="Posts" active>
              <router-link
                :to="'/post/' + post.id"
                v-for="post in posts"
                :key="post.id"
                class="card user-card"
                tag="div"
              >
                <h5>
                  <strong>{{ post.title }}</strong>
                </h5>
                {{ post.content }}
              </router-link></b-tab
            >
            <b-tab title="Replies"
              ><router-link
                :to="'/post/' + reply.post"
                class="card user-card"
                v-for="(reply, index) in replies"
                :key="index"
                tag="div"
              >
                <h5>
                  {{ reply.content }}
                </h5>
              </router-link></b-tab
            >
          </b-tabs>
        </div>
      </div>
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
.user-info .user-card {
  cursor: pointer;
}

.user-info .user-card:hover {
  opacity: 0.8;
}

.user-card img {
  display: block;
  margin-left: auto;
  margin-right: auto;
  margin-bottom: 15px;
  width: 150px;
  height: 150px;
  border-radius: 150px;
}

.user-card h3 {
  text-align: center;
}

.card-body {
  padding: 0;
}

.nav-tabs .nav-item.show .nav-link,
.nav-tabs .nav-link.active {
  background-color: #278ea5;
  border-color: #278ea5;
  color: white;
}

.nav-link {
  color: #278ea5;
}

.nav-tabs .nav-item.show .nav-link,
.nav-tabs .nav-link:not(.active):hover {
  color: #278ea5;
  border-color: #278ea5;
  opacity: 0.8;
}
</style>
