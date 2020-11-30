<template>
  <div class="grid">
    <div class="columns">
      <div class="info-col" v-if="user && isMobile()">
        <UserCard :userData="user" />
      </div>
      <div class="main-col">
        <div class="user-info">
          <b-tabs card v-if="posts.length > 0 || comments.length > 0">
            <b-tab v-if="posts.length" title="Posts" active @click="getPosts()">
              <router-link
                :to="'/post/' + post.id"
                v-for="post in posts"
                :key="post.id"
                :class="`user-card ${isMobile() ? 'card-m' : 'card'}`"
                tag="div"
              >
                <h5>
                  <strong>{{ post.title }}</strong>
                </h5>

                {{ post.content }}
                <small>
                  <span v-b-tooltip.hover title="Rating">
                    <b-icon
                      :icon="reactionIcon(post.your_reaction)"
                      :color="reactionColor(post.your_reaction)"
                    >
                    </b-icon
                    >{{ post.rating }}
                  </span>
                  <span v-b-tooltip.hover title="Comments">
                    <b-icon-chat></b-icon-chat> {{ post.comments_count }}
                  </span>
                  <span v-b-tooltip.hover title="Participants">
                    <b-icon-people></b-icon-people>
                    {{ post.participants_count }}
                  </span>
                  <time-ago
                    v-b-tooltip.hover
                    :title="post.created"
                    :datetime="post.created"
                    :long="!isMobile()"
                  >
                  </time-ago>
                </small>
              </router-link>
            </b-tab>
            <b-tab
              v-if="comments.length"
              title="Comments"
              @click="getComments()"
            >
              <router-link
                :to="'/post/' + comment.post"
                :class="`user-card ${isMobile() ? 'card-m' : 'card'}`"
                v-for="comment in comments"
                :key="comment.id"
                tag="div"
              >
                <h5>
                  {{ comment.content }}
                </h5>
                <small>
                  <span v-b-tooltip.hover title="Rating">
                    <b-icon
                      :icon="reactionIcon(comment.your_reaction)"
                      :color="reactionColor(comment.your_reaction)"
                    >
                    </b-icon
                    >{{ comment.rating }}
                  </span>
                  <time-ago
                    v-b-tooltip.hover
                    :title="comment.created"
                    :datetime="comment.created"
                    :long="!isMobile()"
                  >
                  </time-ago>
                </small>
              </router-link>
            </b-tab>
          </b-tabs>
          <b-container v-else align="center">
            <b-img fluid src="@/assets/img/masterpiece.png"> </b-img>
            <p>It's so empty here...</p>
          </b-container>
        </div>
      </div>
      <div class="info-col" v-if="user && !isMobile()">
        <UserCard :userData="user" />
      </div>
    </div>
  </div>
</template>
<script>
import api from "@/router/api";
import TimeAgo from "vue2-timeago";
import UserCard from "@/components/UserCard";

export default {
  watch: {
    "$route.params.id": function() {
      this.$router.go();
    },
  },
  data() {
    return {
      user: {},
      posts: [],
      comments: [],
      activeTab: "",
    };
  },
  mounted() {
    this.getUser();
    this.getPosts();
  },
  components: {
    TimeAgo,
    UserCard,
  },
  methods: {
    reactionColor(yourReaction) {
      return yourReaction === 1 ? "green" : yourReaction === -1 ? "red" : "";
    },
    reactionIcon(yourReaction) {
      return yourReaction === -1 ? "arrow-down" : "arrow-up";
    },
    async getUser() {
      return await api
        .get("user", {
          params: { id: this.$route.params.id },
        })
        .then((response) => {
          this.user = response.data.data;
          document.title = this.user.display_name;
        })
        .catch((error) => {
          console.log(error);
        });
    },
    async getPosts() {
      if (this.activeTab === "posts") return;
      this.activeTab = "posts";
      return await api
        .post("post/find", {
          by: "author",
          author: Number.parseInt(this.$route.params.id),
        })
        .then((response) => {
          this.posts = response.data.data || [];
        })
        .catch((error) => {
          console.log(error);
        });
    },
    async getComments() {
      if (this.activeTab === "comments") return;
      this.activeTab = "comments";
      return await api
        .post("comments/find", {
          by: "user_id",
          id: Number.parseInt(this.$route.params.id),
        })
        .then((response) => {
          this.comments = response.data.data || [];
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
