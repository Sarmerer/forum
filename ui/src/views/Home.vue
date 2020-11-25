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

    <div class="grid">
      <div class="columns">
        <div class="main-col">
          <div>
            <Error v-if="error.show" :errorData="error" />
          </div>
          <b-row
            align-h="between"
            align-v="end"
            class="mx-3 mt-3"
            v-if="posts.length > 0"
          >
            <b-pagination
              v-model="pagination.currentPage"
              :total-rows="pagination.totalPages"
              :per-page="pagination.perPage"
              aria-controls="my-table"
              @change="handlePageChange"
              first-number
              last-number
            ></b-pagination>
            <div>
              <b-button
                variant="dark"
                @click="sort(), toast('b-toaster-bottom-center', true)"
                :disabled="sorter.throttled"
                class="mx-2"
                v-b-tooltip.hover
                :title="sorter.asc ? 'Ascending' : 'Descending'"
              >
                <b-icon
                  :icon="sorter.asc ? 'sort-up' : 'sort-down-alt'"
                ></b-icon>
              </b-button>
              <b-button-group
                ><b-button
                  :disabled="sorter.throttled"
                  @click="order('rating')"
                  v-b-tooltip.hover
                  title="Most liked"
                  :variant="sorter.orderBy == 'rating' ? 'info' : 'dark'"
                  ><b-icon-heart></b-icon-heart></b-button
                ><b-button
                  :disabled="sorter.throttled"
                  @click="order('created')"
                  v-b-tooltip.hover
                  title="Most recent"
                  :variant="sorter.orderBy == 'created' ? 'info' : 'dark'"
                  ><b-icon-clock></b-icon-clock>
                </b-button>
                <b-button
                  :disabled="sorter.throttled"
                  @click="order('comments_count')"
                  v-b-tooltip.hover
                  title="Most commented"
                  :variant="
                    sorter.orderBy == 'comments_count' ? 'info' : 'dark'
                  "
                  ><b-icon-chat-left></b-icon-chat-left
                ></b-button>
                <b-button
                  :disabled="sorter.throttled"
                  @click="order('total_participants')"
                  v-b-tooltip.hover
                  title="Most participants"
                  :variant="
                    sorter.orderBy == 'total_participants' ? 'info' : 'dark'
                  "
                  ><b-icon-people></b-icon-people></b-button
              ></b-button-group>
            </div>
          </b-row>
          <!-- Start of posts -->
          <router-link
            :to="'/post/' + post.id"
            v-for="post in posts"
            :key="post.id"
            class="card"
            tag="div"
            style="cursor: pointer"
          >
            <b-row>
              <b-col cols="1">
                <Rating :callback="rate" :entity="post" type="comment" />
              </b-col>
              <b-col cols="11" class="post-content" style="flex-wrap: wrap">
                <h5>
                  {{ post.title }}
                </h5>
                <p v-if="!isMobile()">{{ post.content }}</p>
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
            <div class="post-footer">
              <small>
                <img
                  src="@/assets/svg/post/comments.svg"
                  alt="comments"
                  srcset=""
                />
                {{ post.comments_count }}
                <!-- {{ post.comments_count == 1 ? "comment" : "comments" }} -->
              </small>

              <!-- TODO: Make this look decent -->
              <!-- style is embedded here for responsiveness. MB fix later -->
              <small v-if="isMobile()"
                ><Rating
                  style="flex-direction: row; margin: 0"
                  type="post"
                  :callbakc="rate"
                  :entity="post"
              /></small>
              <small
                >by
                <router-link :to="'/user/' + post.author_id" class="secondary">
                  {{ post.author_name }}
                </router-link>
                <timeago :datetime="post.created" :auto-update="60"></timeago
              ></small>
            </div>
          </router-link>
          <!-- End of posts -->
        </div>

        <div class="info-col">
          <div class="card">
            <h3 class="primary">RECENT</h3>
            <router-link
              :to="'/post/' + post.id"
              v-for="(post, index) in recent"
              :key="index"
              tag="div"
              style="cursor: pointer"
              ><p>
                {{ post.title }}<br /><small class="text-muted"
                  ><timeago :datetime="post.created" :auto-update="10"></timeago
                ></small>
              </p>
            </router-link>
          </div>
          <div class="card">
            <h3 class="primary">CATEGORIES</h3>
            <!-- Start of categories -->
            <span v-if="categories.length == 0">None</span>
            <b-container v-else>
              <div class="categories">
                <b-button-group
                  v-for="c in categories"
                  :key="c.ID"
                  :id="c.ID"
                  class="category-name"
                  size="sm"
                  ><b-button disabled>{{ c.name }}</b-button>
                  <b-button disabled class="category-count">
                    {{ c.use_count }}
                  </b-button></b-button-group
                >
              </div>
            </b-container>
            <!-- End of categories -->
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import axios from "axios";
import { mapGetters } from "vuex";
import Error from "@/components/Error";
import Rating from "@/components/Rating";
export default {
  name: "Home",
  computed: {
    ...mapGetters({
      user: "auth/user",
    }),
  },
  components: {
    Rating,
    Error,
  },
  data() {
    return {
      posts: [],
      recent: [],
      categories: [],
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
  },
  methods: {
    async handlePageChange(value) {
      await this.getPosts(value);
    },
    async getPosts(currentPage) {
      return await axios
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
      return await axios
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
    toast(toaster, append = true) {
      var message = this.sorter.asc ? "ascending" : "descending";
      this.$bvToast.toast(`Posts sorted in ${message} order.`, {
        // title: `YAY`,
        toaster: toaster,
        solid: true,
        // variant: "secondary",
        appendToast: append,
        noCloseButton: true,
      });
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
      await axios
        .post("post/rate", { id: post.id, reaction: r })
        .then((response) => {
          post.your_reaction = response.data.data.your_reaction;
          post.rating = response.data.data.rating;
        })
        .catch((error) => {
          console.log(error);
        });
    },
  },
};
</script>
<style lang="scss" scoped>
/* The hero image */
.hero-image {
  background-image: url("../assets/img/home-hero.jpg");
  box-shadow: 0 5px 6px 2px rgba(10, 10, 10, 0.3);
  height: 200px;
  width: 100%;
  background-position: center;
  background-repeat: no-repeat;
  background-size: cover;
  position: relative;
  margin-bottom: 25px;
}

.hero-text {
  text-align: center;
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  font-size: large;
}

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

ul#filters {
  padding: 0;
  margin: 0;
}

ul#filters li {
  display: inline;
  padding: 7px;
  font-size: 22px;
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

/* TOGGLE SWITCH START */

input[type="checkbox"] {
  height: 0;
  width: 0;
  visibility: hidden;
}

label {
  cursor: pointer;
  /* bottom: 5px; */
  margin: 0;
  text-indent: -9999px;
  width: 35px;
  height: 20px;
  background: grey;
  border-radius: 50px;
  position: relative;
}

label:after {
  content: "";
  position: absolute;
  top: -5px;
  left: 0;
  width: 15px;
  height: 15px;
  background: #278ea5;
  border-radius: 30px;
  transition: 0.3s;
}

input:checked + label:after {
  left: calc(100%);
  transform: translateX(-100%);
}

.btn-dark {
  // color: #fff;
  background-color: rgba(255, 255, 255, 0.05);
  border: none;
}

.btn-dark:hover {
  background-color: rgba(255, 255, 255, 0.03);
}

/* TOGGLE SWITCH END */

@media (max-width: 500px) {
  .post-content {
    margin-left: -20px;
  }
}
</style>
