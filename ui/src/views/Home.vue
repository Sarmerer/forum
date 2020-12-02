<template>
  <div>
    <div class="hero-image">
      <div class="hero-text">
        <h1>
          WELCOME<br />
          TO <span class="primary">FORUM</span>
        </h1>
      </div>
    </div>

    <b-skeleton-wrapper :loading="loading">
      <template #loading>
        <HomeSkeleton v-bind:postsLength="posts.length" />
      </template>
      <div class="grid">
        <div class="columns">
          <div class="main-col">
            <div>
              <Error v-if="error.show" :errorData="error" />
            </div>
            <div>
              <b-row class="mx-3 mt-3" v-if="posts.length > 0 && !isMobile()">
                <b-col
                  align=" start"
                  v-if="posts.total_rows > pagination.perPage"
                >
                  <b-pagination
                    v-model="pagination.currentPage"
                    :total-rows="pagination.totalPages"
                    :per-page="pagination.perPage"
                    aria-controls="my-table"
                    @change="handlePageChange"
                    first-number
                    last-number
                  >
                  </b-pagination>
                </b-col>
                <b-col v-if="posts.length > 1" align="end">
                  <PostFilters
                    :orderCallback="order"
                    :sortCallback="sort"
                    :sorter="sorter"
                  />
                </b-col>
              </b-row>
              <b-container v-if="posts.length > 0 && isMobile()">
                <b-row
                  v-if="posts.total_rows > pagination.perPage"
                  align-h="center"
                  class="mb-2"
                >
                  <b-pagination
                    v-model="pagination.currentPage"
                    :total-rows="pagination.totalPages"
                    :per-page="pagination.perPage"
                    aria-controls="my-table"
                    @change="handlePageChange"
                    first-number
                    last-number
                  >
                  </b-pagination>
                </b-row>
                <b-row v-if="posts.length > 1" align-h="center">
                  <PostFilters
                    :orderCallback="order"
                    :sortCallback="sort"
                    :sorter="sorter"
                  />
                </b-row>
              </b-container>
              <!-- Start of posts -->
              <router-link
                :to="'/post/' + post.id"
                v-for="post in posts"
                :key="post.id"
                :class="`text-break ${isMobile() ? 'card-m' : 'card'}`"
                tag="div"
                style="cursor: pointer"
              >
                <b-row>
                  <b-col v-if="!isMobile()" cols="start">
                    <Rating :callback="rate" :entity="post" class="ml-n4" />
                  </b-col>
                  <b-col class="ml-2">
                    <small v-if="isMobile()">
                      <b-img :src="post.author.avatar" width="15px"></b-img>
                      <router-link
                        :to="'/user/' + post.author.id"
                        class="secondary"
                      >
                        {{ post.author.display_name }}
                      </router-link>
                      <time-ago tooltip :datetime="post.created"></time-ago>
                    </small>
                    <h5>
                      {{ post.title }}
                    </h5>
                    <pre>{{ post.content }}</pre>
                    <b-form-tag
                      v-for="(category, index) in post.categories"
                      disabled
                      :key="index"
                      :title="category.name"
                      variant="dark"
                      class="mr-1 mb-1"
                    >
                      {{ category.name }}
                    </b-form-tag>
                  </b-col>
                </b-row>
                <b-row class="ml-1">
                  <b-col>
                    <small>
                      <span v-b-tooltip.hover title="Comments">
                        <b-icon-chat></b-icon-chat> {{ post.comments_count }}
                      </span>
                      <span v-b-tooltip.hover title="Participants">
                        <b-icon-people></b-icon-people>
                        {{ post.participants_count }}
                      </span>
                    </small>
                  </b-col>
                  <b-col cols="end" class="mr-4">
                    <small v-if="!isMobile()"
                      >by
                      <b-img :src="post.author.avatar" width="15px"></b-img>
                      <router-link
                        :to="'/user/' + post.author.id"
                        class="secondary"
                      >
                        {{ post.author.display_name }}
                      </router-link>
                      <time-ago
                        tooltip
                        :datetime="post.created"
                        long
                      ></time-ago>
                    </small>
                  </b-col>
                  <b-col v-if="isMobile()" cols="end" class="mr-4">
                    <small>
                      <Rating
                        class="mr-2"
                        compact
                        :callback="rate"
                        :entity="post"
                      />
                    </small>
                  </b-col>
                </b-row>
              </router-link>
              <!-- End of posts -->
            </div>
            <b-container v-if="posts.length === 0" align="center">
              <b-img-lazy fluid src="@/assets/img/empty.png"> </b-img-lazy>
              <p>It's so empty here...</p>
            </b-container>
          </div>

          <div class="info-col">
            <div :class="`text-break ${isMobile() ? 'card-m' : 'card'}`">
              <h3 class="primary">RECENT</h3>
              <span v-if="recent.length == 0"
                >None...
                <router-link to="/new-post" class="secondary">yet</router-link>
              </span>
              <span v-else>
                <router-link
                  :to="'/post/' + post.id"
                  v-for="(post, index) in recent"
                  :key="index"
                  tag="div"
                  class="ml-2"
                  style="cursor: pointer"
                >
                  <small>
                    <router-link
                      :to="'/user/' + post.author.id"
                      class="secondary"
                    >
                      {{ post.author.display_name }}
                    </router-link>
                    <time-ago
                      tooltip
                      :datetime="post.created"
                      :long="!isMobile()"
                    >
                    </time-ago>
                  </small>
                  <p>
                    {{ post.title }}
                  </p>
                </router-link>
              </span>
            </div>
            <div :class="`text-break ${isMobile() ? 'card-m' : 'card'}`">
              <h3 class="primary">
                CATEGORIES<b-button id="popover-filter-button">
                  <b-icon-three-dots-vertical></b-icon-three-dots-vertical>
                </b-button>
              </h3>
              <!-- Start of categories -->
              <span v-if="categories.length == 0"
                >None...
                <router-link to="/new-post" class="secondary">yet</router-link>
              </span>
              <b-container v-else>
                <div class="categories">
                  <b-form-checkbox-group
                    v-for="c in categories"
                    :key="c.ID"
                    :id="c.ID"
                    class="category-name"
                    size="sm"
                    v-model="selectedCategories"
                    buttons
                    ><b-form-checkbox :value="c.name">{{
                      c.name
                    }}</b-form-checkbox>
                    <b-form-checkbox disabled class="category-count">
                      {{ c.use_count }}
                    </b-form-checkbox></b-form-checkbox-group
                  >
                </div>
              </b-container>
              <!-- End of categories -->
            </div>
            <b-popover
              target="popover-filter-button"
              triggers="focus"
              variant="dark"
            >
              <b-button
                v-on:click="sortByCategories()"
                class="mb-1"
                style="width: 135px"
              >
                <b-icon-filter> </b-icon-filter> filter</b-button
              >
              <br />
              <b-button v-on:click="resetCategories()" style="width: 135px">
                <b-icon-arrow-clockwise></b-icon-arrow-clockwise>
                reset</b-button
              >
            </b-popover>
          </div>
        </div>
      </div>
    </b-skeleton-wrapper>
  </div>
</template>

<script>
import api from "@/router/api";
import { mapGetters } from "vuex";
import Error from "@/components/Error";
import Rating from "@/components/Rating";
import TimeAgo from "vue2-timeago";
import PostFilters from "@/components/PostFilters";
import HomeSkeleton from "@/components/HomeSkeleton";

export default {
  name: "Home",
  computed: {
    ...mapGetters({
      user: "auth/user",
    }),
  },
  components: {
    HomeSkeleton,
    PostFilters,
    TimeAgo,
    Rating,
    Error,
  },
  data() {
    return {
      loading: false,
      loadingTime: 0,
      maxLoadingTime: 1,
      posts: [],
      recent: [],
      categories: [],
      selectedCategories: [],
      sorter: { orderBy: "rating", asc: true, throttled: false },
      pagination: { currentPage: 1, totalPages: 1, perPage: 7 },
      error: {
        show: false,
        status: Number,
        message: String,
        callback: Function,
      },
      checked: false,
    };
  },
  created() {
    this.getPosts(0);
    this.getCategories();
    this.$_loadingTimeInterval = null;
  },
  watch: {
    loading(newVal, oldValue) {
      if (newVal !== oldValue) {
        this.clearLoadingTimeInterval();

        if (newVal) {
          this.$_loadingTimeInterval = setInterval(() => {
            this.loadingTime++;
          }, 1000);
        }
      }
    },
    loadingTime(newVal, oldValue) {
      if (newVal !== oldValue) {
        if (newVal === this.maxLoadingTime) {
          this.loading = false;
        }
      }
    },
  },

  mounted() {
    this.startLoading();
  },

  methods: {
    clearLoadingTimeInterval() {
      clearInterval(this.$_loadingTimeInterval);
      this.$_loadingTimeInterval = null;
    },
    startLoading() {
      this.loading = true;
      this.loadingTime = 0;
    },
    async handlePageChange(value) {
      await this.getPosts(value);
    },
    async getPosts(currentPage) {
      return await api
        .post("posts", {
          per_page: this.pagination.perPage,
          current_page: currentPage,
          order_by: this.sorter.orderBy,
          ascending: this.sorter.asc,
        })
        .then((response) => {
          this.error.show = false;
          this.posts = response.data.data.hot || [];
          this.recent = response.data.data.recent || [];
          this.pagination.totalPages = response.data.data.total_rows || 5;
        })
        .catch((error) => {
          this.posts = [];
          this.recent = [];
          this.categories = [];
          this.error.show = true;
          this.error.status = error.response.status;
          this.error.message = error.response.statusText;
          this.error.callback = this.getPosts;
        });
    },
    async getCategories() {
      return await api
        .get("categories")
        .then((response) => {
          this.categories = response.data.data || [];
        })
        .catch((error) => {
          console.log(error);
        });
    },
    sort() {
      this.sorter.asc = !this.sorter.asc;
      this.throttle();
    },
    order(by) {
      if (this.sorter.orderBy === by) return;
      this.sorter.orderBy = this.sorter.orderBy = by;
      this.throttle();
    },
    throttle() {
      if (this.sorter.throttled) return;
      this.sorter.throttled = true;
      this.getPosts().then(
        setTimeout(() => {
          this.sorter.throttled = false;
        }, 1000)
      );
    },
    async rate(reaction, post) {
      if (this.requesting) return;
      let r = reaction == "up" ? 1 : -1;
      if (
        (reaction == "up" && post.your_reaction == 1) ||
        (reaction == "down" && post.your_reaction == -1)
      ) {
        r = 0;
      }
      await api
        .post("post/rate", { id: post.id, reaction: r })
        .then((response) => {
          post.your_reaction = response.data.data.your_reaction;
          post.rating = response.data.data.rating;
        })
        .catch((error) => {
          console.log(error);
        });
    },
    async sortByCategories() {
      await api
        .post("post/find", {
          by: "categories",
          categories: this.selectedCategories,
        })
        .then((response) => {
          this.posts = response.data.data || [];
        })
        .catch((error) => {
          console.log(error);
        });
    },
    async resetCategories() {
      this.selectedCategories = [];
      this.getPosts(0);
    },
  },
};
</script>
<style lang="scss" scoped>
.post-content {
  margin-left: -40px;
}

.post-footer {
  line-height: 36px;
  overflow: hidden;
  padding: 2px 18px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  flex-wrap: wrap;
}

.categories {
  display: flex;
  flex-wrap: wrap;
}

.category-name {
  margin: 2px;
}
.categories .btn-secondary.disabled.category-count {
  background-color: #278ea5;
}

.categories .btn-secondary.disabled {
  background-color: #343a40;
  border: none;
}

.btn-dark {
  // color: #fff;
  background-color: rgba(255, 255, 255, 0.05);
  border: none;
}

.btn-dark:hover {
  background-color: rgba(255, 255, 255, 0.03);
}

#popover-filter-button {
  background-color: rgba(10, 10, 10, 0);
  border: none;
  color: #278ea5;
  padding: 0 2px 4px 2px;
  font-size: 23px;
}

@media (max-width: 500px) {
  .post-content {
    margin-left: -20px;
  }
}
</style>
