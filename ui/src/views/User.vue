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
                      :disabled="tabDisabled(tab)"
                      @click="activeTab = tab.name"
                      :size="isMobile() ? 'sm' : 'md'"
                      :variant="
                        `outline-${
                          activeTab === tab.name && !tab.disabled
                            ? 'info'
                            : tabDisabled(tab)
                            ? 'dark'
                            : 'light'
                        }`
                      "
                      class="mr-2"
                    >
                      {{ tab.title }}
                      <b-badge
                        v-if="user[tab.prop]"
                        :variant="activeTab === tab.name ? 'info' : 'light'"
                        >{{ user[tab.prop] }}
                      </b-badge>
                    </b-button>
                  </b-col>
                </b-row>
              </div>
              <component :is="activeTab"></component>
            </div>
            <b-container v-if="!user.posts && !user.comments" align="center">
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
import ControlButtons from "@/components/ControlButtons";
import NotFound from "@/components/NotFound";
import UserCard from "@/components/UserCard";
import TimeAgo from "@/components/TimeAgo";
import api from "@/api/api";

export default {
  name: "UserPage",
  watch: {
    "$route.params.userID": function(id) {
      if (!id || this.prevUserID == id) return;
      this.user = {};
      this.getUser().then(() => {
        this.activeTab = this.user.posts > 0 ? "PostsTab" : "CommentsTab";
      });
    },
  },
  data() {
    return {
      loading: true,
      notFound: false,
      user: {},
      prevUserID: 0,
      activeTab: "",
      tabs: [
        {
          title: "Posts",
          name: "PostsTab",
          prop: "posts",
          component: this.PostTab,
        },
        {
          title: "Comments",
          name: "CommentsTab",
          prop: "comments",
          component: this.PostTab,
        },
        {
          title: "Saved",
          prop: "saved",
          component: this.PostTab,
          disabled: true,
        },
      ],
    };
  },
  computed: {
    userID: function() {
      return Number.parseInt(this.$route.params.userID);
    },
  },
  activated() {
    this.getUser().then(() => {
      if (!this.activeTab)
        this.activeTab =
          this.user.posts > 0
            ? "PostsTab"
            : this.user.comments > 0
            ? "CommentsTab"
            : "";
      setTimeout(() => {
        this.loading = false;
        this.madeRequest = true;
      }, 500);
    });
  },
  components: {
    ControlButtons,
    UserSkeleton,
    UserCard,
    NotFound,
    PostsTab: () => import("@/components/user-page-tabs/Posts"),
    CommentsTab: () => import("@/components/user-page-tabs/Comments"),
    TimeAgo,
  },
  methods: {
    tabDisabled(tab) {
      if (!tab) return true;
      return (
        this.activeTab === tab.name || tab.disabled || !this.user[tab.prop]
      );
    },
    async getUser() {
      return await api
        .post("user/find", {
          by: "id",
          id: this.userID,
        })
        .then((response) => {
          this.user = response.data.data;
          this.prevUserID = this.user.id;
          document.title = this.user.alias;
        })
        .catch((error) => {
          if (error.status === 404) this.notFound = true;
        });
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
