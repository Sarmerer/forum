<template>
  <b-skeleton-wrapper :loading="loading">
    <template #loading>
      <UserSkeleton />
    </template>
    <NotFound v-if="notFound" />
    <div v-else>
      <div class="grid">
        <div class="columns">
          <div class="info-col" v-if="user && isMobile()">
            <UserCard :userData="user" />
          </div>
          <div class="main-col">
            <div class="user-info" v-if="user.posts || user.comments">
              <div :class="isMobile() ? 'card-m' : 'card'">
                <b-row>
                  <b-col>
                    <b-button
                      v-for="tab in tabs"
                      :key="tab.title"
                      :disabled="!user[tab.prop]"
                      @click="tab.callback(), (activeTab = tab.prop)"
                      :size="isMobile() ? 'sm' : 'md'"
                      :variant="
                        `outline-${activeTab === tab.prop ? 'info' : 'light'}`
                      "
                      class="mr-2"
                    >
                      {{ tab.title }}
                      <b-badge
                        v-if="user[tab.prop]"
                        :variant="activeTab === tab.prop ? 'info' : 'light'"
                        >{{ user[tab.prop] }}</b-badge
                      >
                    </b-button>
                  </b-col>
                </b-row>
              </div>
              <div v-if="activeTab === 'posts'">
                <router-link
                  :to="'/post/' + post.id"
                  v-for="post in posts"
                  :key="post.id"
                  :class="
                    `user-card text-break ${isMobile() ? 'card-m' : 'card'}`
                  "
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
              </div>
              <div v-if="activeTab === 'comments'">
                <router-link
                  :to="'/post/' + comment.post_id"
                  :class="
                    `user-card text-break ${isMobile() ? 'card-m' : 'card'}`
                  "
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
              </div>
            </div>
            <b-container
              v-if="!posts.length && !comments.length && madeRequest"
              align="center"
            >
              <b-img-lazy fluid src="@/assets/img/empty.png"> </b-img-lazy>
              <p>It's so empty here...</p>
            </b-container>
          </div>
          <div class="info-col" v-if="user && !isMobile()">
            <UserCard :userData="user" />
          </div>
        </div>
      </div>
    </div>
  </b-skeleton-wrapper>
</template>
<script>
import UserSkeleton from "@/components/skeletons/UserSkeleton";
import NotFound from "@/components/NotFound";
import UserCard from "@/components/UserCard";
import TimeAgo from "@/components/TimeAgo";
import api from "@/router/api";

export default {
  watch: {
    "$route.params.id": function() {
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
      notFound: false,
      user: {},
      posts: [],
      comments: [],
      activeTab: "",
      madeRequest: false,
      tabs: [
        { title: "Posts", prop: "posts", callback: this.getPosts },
        { title: "Comments", prop: "comments", callback: this.getComments },
        { title: "Saved", prop: "saved", callback: this.getComments },
      ],
    };
  },
  created() {
    let p = this.getUser();
    let p1 = p.then(() => {
      this.user.posts > 0
        ? (this.getPosts(), (this.activeTab = "posts"))
        : this.user.comments > 0
        ? (this.getComments(), (this.activeTab = "comments"))
        : (this.madeRequest = true);
    });
    Promise.all([p, p1]).then(() =>
      setTimeout(() => {
        this.loading = false;
      }, 500)
    );
  },
  components: {
    UserSkeleton,
    UserCard,
    NotFound,
    TimeAgo,
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
          document.title = this.user.alias;
        })
        .catch((error) => {
          if (error.status === 404) this.notFound = true;
        });
    },
    async getPosts() {
      if (this.activeTab === "posts") return;
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
<style lang="scss" scoped>
.user-info .user-card {
  cursor: pointer;
}

.user-info .user-card:hover {
  border: white;
}

.user-card h3 {
  text-align: center;
}
</style>
