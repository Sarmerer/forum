<template>
  <b-skeleton-wrapper :loading="loading">
    <template #loading>
      <UserSkeleton />
    </template>
    <div class="grid">
      <div class="columns">
        <div class="info-col" v-if="user && isMobile()">
          <UserCard :userData="user" />
        </div>
        <div class="main-col">
          <div class="user-info">
            <b-tabs card v-if="user.posts || user.comments">
              <b-tab
                v-if="user.posts > 0"
                title="Posts"
                :active="user.posts > 0"
                @click="getPosts()"
              >
                <router-link
                  :to="'/post/' + post.id"
                  v-for="post in posts"
                  :key="post.id"
                  :class="`user-card text-break ${
                    isMobile() ? 'card-m' : 'card'
                  }`"
                  tag="div"
                >
                  <h5>
                    <strong>{{ post.title }}</strong>
                  </h5>
                  <pre>{{ post.content }}</pre>
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
                    <time-ago :datetime="post.created" tooltip="right">
                    </time-ago>
                  </small>
                </router-link>
              </b-tab>
              <b-tab
                v-if="user.comments > 0"
                title="Comments"
                :active="!user.posts"
                @click="getComments()"
              >
                <router-link
                  :to="'/post/' + comment.post"
                  :class="`user-card text-break ${
                    isMobile() ? 'card-m' : 'card'
                  }`"
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
                    <time-ago :datetime="comment.created" tooltip="right">
                    </time-ago>
                  </small>
                </router-link>
              </b-tab>
            </b-tabs>
            <b-container
              v-if="!posts.length && !comments.length && madeRequest"
              align="center"
            >
              <b-img-lazy fluid src="@/assets/img/empty.png"> </b-img-lazy>
              <p>It's so empty here...</p>
            </b-container>
          </div>
        </div>
        <div class="info-col" v-if="user && !isMobile()">
          <UserCard :userData="user" />
        </div>
      </div>
    </div>
  </b-skeleton-wrapper>
</template>
<script>
import api from "@/router/api";
import TimeAgo from "@/components/TimeAgo";
import UserCard from "@/components/UserCard";
import UserSkeleton from "@/components/skeletons/UserSkeleton";

export default {
  watch: {
    "$route.params.id": function () {
      this.user = {};
      this.posts = [];
      this.comments = [];
      this.activeTab = "";
      this.madeRequest = false;
      this.getUser().then(() => {
        this.user.posts > 0 ? this.getPosts() : this.getComments();
      });
    },
  },
  data() {
    return {
      loading: true,
      user: {},
      posts: [],
      comments: [],
      activeTab: "",
      madeRequest: false,
    };
  },
  created() {
    let p = this.getUser();
    let p1 = p.then(() => {
      this.user.posts > 0
        ? this.getPosts()
        : this.user.comments > 0
        ? this.getComments()
        : (this.madeRequest = true);
    });
    Promise.all([p, p1]).then(() =>
      setTimeout(() => {
        this.loading = false;
      }, 500)
    );
  },
  components: {
    TimeAgo,
    UserCard,
    UserSkeleton,
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
        .post("user/find", {
          by: "id",
          id: Number.parseInt(this.$route.params.id),
        })
        .then((response) => {
          this.user = response.data.data;
          document.title = this.user.display_name;
        })
        .catch((error) => {
          if (error.status === 404) this.$router.push("/not-found");
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
        })
        .then(() => (this.madeRequest = true));
    },
    async getComments() {
      if (this.activeTab === "comments") return;
      this.activeTab = "comments";
      return await api
        .post("comments/find", {
          by: "author",
          author: Number.parseInt(this.$route.params.id),
        })
        .then((response) => {
          this.comments = response.data.data || [];
        })
        .catch((error) => {
          console.log(error);
        })
        .then(() => (this.madeRequest = true));
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
