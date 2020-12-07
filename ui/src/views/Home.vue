<template>
  <div>
    <div class="hero-image" v-if="!isMobile()">
      <div class="hero-text">
        <h1>
          WELCOME<br />
          TO <span class="primary">FORUM</span>
        </h1>
      </div>
    </div>
    <b-skeleton-wrapper :loading="showSkeleton">
      <template #loading>
        <HomeSkeleton />
      </template>
      <div class="grid">
        <div class="columns">
          <div class="main-col">
            <div>
              <Error v-if="error.show" :errorData="error" />
            </div>
            <div>
              <div
                v-if="sorter.filtered"
                :class="isMobile() ? 'card-m' : 'card'"
              >
                <b-row>
                  <b-col>
                    Posts with categories:
                    <b-form-tag
                      v-for="(category, index) in sorter.prevCategories"
                      disabled
                      :key="index"
                      :title="category"
                      variant="dark"
                      class="mr-1 mb-1"
                    >
                      {{ category }}
                    </b-form-tag>
                  </b-col>
                  <b-col cols="end">
                    <b-button
                      class="mr-3"
                      @click="resetCategories()"
                      variant="outline-light"
                      size="sm"
                    >
                      <b-icon-arrow-clockwise></b-icon-arrow-clockwise>
                    </b-button>
                  </b-col>
                </b-row>
              </div>
              <b-row class="mx-3 mt-3" v-if="!sorter.filtered && !isMobile()">
                <b-col>
                  <b-pagination
                    v-if="pagination.totalPages > pagination.perPage"
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
                <b-col v-if="pagination.totalPages > 1" cols="end">
                  <PostFilters
                    :orderCallback="order"
                    :sortCallback="sort"
                    :sorter="sorter"
                  />
                </b-col>
              </b-row>
              <b-container class="mt-2" v-if="!sorter.filtered && isMobile()">
                <b-row
                  v-if="pagination.totalPages > pagination.perPage"
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
                <b-row v-if="pagination.totalPages > 1" align-h="center">
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
                v-for="(post, index) in posts"
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
                      <router-link
                        :to="'/user/' + post.author.id"
                        class="secondary"
                      >
                        <b-img :src="post.author.avatar" width="15px"></b-img>
                        {{ post.author.display_name }}
                      </router-link>
                      <time-ago :datetime="post.created"></time-ago>
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
                      <router-link
                        :to="'/user/' + post.author.id"
                        class="secondary"
                      >
                        <user-popover
                          :userData="post.author"
                          :popoverID="'p' + index"
                        >
                        </user-popover>
                      </router-link>
                      <time-ago :datetime="post.created" tooltip="bottom">
                      </time-ago>
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
              <b-img-lazy fluid src="@/assets/img/empty.png" class="w-100">
              </b-img-lazy>
              <p>It's so empty here...</p>
            </b-container>
          </div>

          <div class="info-col">
            <div
              :class="
                `text-break ${
                  isMobile()
                    ? sorter.filtered
                      ? 'card-m d-none'
                      : 'card-m d-block'
                    : 'card'
                }`
              "
            >
              <h3 class="primary">RECENT</h3>
              <span v-if="!recent.length"
                >None...
                <router-link to="/new-post" class="secondary">yet</router-link>
              </span>
              <span v-else>
                <router-link
                  :to="'/post/' + post.id"
                  v-for="(post, index) in recent"
                  :key="index"
                  tag="div"
                  style="cursor: pointer"
                >
                  <small>
                    <router-link
                      :to="'/user/' + post.author.id"
                      class="secondary"
                    >
                      <user-popover
                        :userData="post.author"
                        :popoverID="'r' + index"
                        popoverDirection="right"
                        noAvatar
                      >
                      </user-popover>
                    </router-link>
                    <time-ago :datetime="post.created" tooltip="right">
                    </time-ago>
                  </small>
                  <p>
                    {{ post.title }}
                  </p>
                </router-link>
              </span>
            </div>
            <!-- Start of categories -->
            <div :class="`text-break ${isMobile() ? 'card-m' : 'card'}`">
              <b-row>
                <b-col>
                  <h3 class="primary">
                    TAGS
                  </h3>
                </b-col>
                <b-col cols="end" class="mr-3">
                  <b-button
                    size="sm"
                    @click="sortByCategories()"
                    v-if="sorter.categories.length || sorter.filtered"
                    :disabled="!sorter.categories"
                    variant="outline-info"
                  >
                    <b-icon-filter> </b-icon-filter>
                  </b-button>
                </b-col>
              </b-row>
              <span v-if="!categories.length"
                >None...
                <router-link to="/new-post" class="secondary">yet</router-link>
              </span>
              <b-container v-else class="ml-0 pl-0">
                <div class="categories">
                  <b-form-checkbox-group
                    v-for="c in categories"
                    :key="c.ID"
                    class="category-name"
                    size="sm"
                    v-model="sorter.categories"
                    buttons
                    ><b-form-checkbox :value="c.name">{{
                      c.name
                    }}</b-form-checkbox>
                    <b-form-checkbox disabled class="category-count">
                      {{ c.use_count }}
                    </b-form-checkbox>
                  </b-form-checkbox-group>
                </div>
              </b-container>
              <!-- End of categories -->
            </div>
          </div>
        </div>
      </div>
    </b-skeleton-wrapper>
  </div>
</template>
<script>
import HomeSkeleton from "@/components/skeletons/HomeSkeleton";
import UserPopover from "@/components/UserPopover";
import PostFilters from "@/components/PostFilters";
import TimeAgo from "@/components/TimeAgo";
import Rating from "@/components/Rating";
import Error from "@/components/Error";
import { mapGetters } from "vuex";
import api from "@/router/api";

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
    UserPopover,
    TimeAgo,
    Rating,
    Error,
  },
  data() {
    return {
      showSkeleton: true,
      posts: [],
      recent: [],
      categories: [],
      sorter: {
        orderBy: "rating",
        asc: true,
        throttled: false,
        categories: [],
        prevCategories: [],
        filtered: false,
      },
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
    Promise.all([this.getPosts(0), this.getCategories()]).then(() => {
      setTimeout(() => {
        this.showSkeleton = false;
      }, 500);
    });
  },
  methods: {
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
          this.posts = response.data?.data?.hot || [];
          this.recent = response.data?.data?.recent || [];
          this.pagination.totalPages = response?.data?.data?.total_rows || 0;
        })
        .catch((error) => {
          this.posts = [];
          this.recent = [];
          this.categories = [];
          this.error.show = true;
          this.error.status = error.response?.status;
          this.error.message = error.response?.statusText;
          this.error.callback = this.getPosts;
        });
    },
    async getCategories() {
      return await api.get("categories").then((response) => {
        this.categories = response.data.data || [];
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
        }, 500)
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
      if (this.sorter.categories === this.sorter.prevCategories) return;
      await api
        .post("post/find", {
          by: "categories",
          categories: this.sorter.categories,
        })
        .then((response) => {
          this.posts = response.data.data || [];
          this.sorter.prevCategories = this.sorter.categories;
          this.sorter.filtered = true;
        })
        .catch((error) => {
          console.log(error);
        });
    },
    async resetCategories() {
      this.sorter.categories = [];
      this.sorter.prevCategories = [];
      this.sorter.filtered = false;
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
  font-size: 0.65rem !important;
  // padding: 0.1rem 0.4rem !important;
}
.categories .btn-secondary {
  background-color: #343a40;
  border-color: #343a40;
  color: rgba(255, 255, 255, 0.87);
}

.categories .btn-secondary.disabled.category-count {
  background-color: #278ea5;
  opacity: 1;
}

.categories .btn-secondary.disabled {
  background-color: #343a40;
  border: none;
}

.categories .btn-secondary.active {
  background-color: #262a2e;
  border-color: #278ea5;
  color: #278ea5;
}

.btn-dark {
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

.popover {
  background-color: #282828;
  opacity: 0.87;
  border: none;
}
</style>
