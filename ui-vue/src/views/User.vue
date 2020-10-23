<template>
  <div class="info-col">
    <div class="card">
      <img
        style="width: 100px; height:100px; border-radius: 200px"
        :src="user.avatar"
        alt="avatar"
        srcset=""
      />
      <h3 class="primary">
        {{ user.display_name }}
        <b-badge v-if="user.role == 2" variant="primary">Admin</b-badge>
      </h3>
      <p>Last online: {{ user.last_online | formatDate }}</p>
    </div>
    <div class="user-info">
      <b-tabs pills card>
        <b-tab title="Posts" active>
          <div class="post" v-for="(post, index) in posts" :key="index">
            <h5>
              <strong>{{ post.title }}</strong>
            </h5>
            {{ post.content }}
          </div></b-tab
        >
        <b-tab title="Replies"
          ><div class="post" v-for="(reply, index) in replies" :key="index">
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
      replies: []
    };
  },
  mounted() {
    this.getUser();
    this.getUserPosts();
    this.getReplies();
  },
  filters: {
    formatDate: function(value) {
      if (value) {
        return moment(String(value)).format("MM/DD/YYYY hh:mm");
      }
    }
  },
  methods: {
    async getUser() {
      return await axios
        .get("user", { params: { id: this.$route.params.id } })
        .then(response => {
          this.user = response.data.data;
        })
        .catch(error => {
          console.log(error);
        });
    },
    async getUserPosts() {
      return await axios
        .post("post/find", {
          by: "author",
          author: Number.parseInt(this.$route.params.id)
        })
        .then(response => {
          this.posts = response.data.data || [];
          console.log(response);
        })
        .catch(error => {
          console.log(error);
        });
    },
    async getReplies() {
      return await axios
        .get("comments", { params: { id: this.$route.params.id } })
        .then(response => {
          this.replies = response.data.data || [];
          console.log(response.data.data);
        })
        .catch(error => {
          console.log(error);
        });
    }
  }
};
</script>
<style lang="scss">
.wrapper {
  text-align: center;
  width: 800px;
  margin: 0 25%;
}

.columns {
  display: flex;
  flex-wrap: wrap;
  margin: 0 17%;
}

.columns > * {
  flex-basis: calc(calc(750px - 100%) * 999);
}

.card {
  margin: 20px 200px;
  padding: 20px;
  background-color: rgba(255, 255, 255, 0.05);
  box-shadow: 5px 5px 6px 2px rgba(10, 10, 10, 0.3);
}

.card-body {
  padding: 0;
}
.user-info {
  margin: 20px 200px;
}
.post {
  margin: 15px 0;
  padding: 15px;
  background-color: rgba(255, 255, 255, 0.05);
  box-shadow: 5px 5px 6px 2px rgba(10, 10, 10, 0.3);
}

.post-col {
  flex-grow: 1;
}

.info-col {
  flex-grow: 0.4;
}

hr {
  opacity: 0.3;
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