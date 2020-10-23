<template>
  <div class="home ">
    <!-- <div class="hero-image">
      <div class="hero-text">
        <h1>
          WELCOME<br />
          TO <span class="primary">FORUM</span>
        </h1>
      </div>
    </div> -->

    <div class="main-content">
      <div class="columns">
        <div class="post-col">
          <div>
            <Error v-if="error.show" :errorData="error" />
          </div>
          <div v-if="posts.length > 0" class="card">
            <b-button @click="sortPosts()"
              >sort by date:
              {{ sorter.byDate ? "ascending" : "descending" }}</b-button
            >
          </div>
          <!-- Start of posts -->
          <router-link
            :to="'/post/' + post.id"
            v-for="(post, index) in posts"
            :key="index"
            class="card post"
            tag="div"
            style="cursor: pointer"
          >
            <Rating
              v-on:update="
                (args) => {
                  post.rating = args.new_rating;
                  post.your_reaction = args.new_your_reaction;
                }
              "
              :postID="post.id"
              :rating="post.rating"
              :yourReaction="post.your_reaction"
            />
            <div style="max-width: 95%">
              <small
                >by
                <router-link
                  :to="'/user/' + post.author_id"
                  style="text-decoration: none;"
                >
                  {{ post.author_name }}
                </router-link>
                <timeago :datetime="post.created" :auto-update="60"></timeago
              ></small>
              <h2 class="primary">
                {{ post.title }}
              </h2>
              <hr />
              <p style="color: white">{{ post.content }}</p>
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
              <br />
              <sub v-if="post.comments_count > 0"
                ><img
                  src="@/assets/svg/post/comments.svg"
                  alt="comments"
                  srcset=""
                />
                {{ post.comments_count }}
                {{ post.comments_count == 1 ? "comment" : "comments" }}</sub
              >
            </div>
          </router-link>
          <!-- End of posts -->
        </div>
        <div class="info-col">
          <div class="card">
            <h3 class="primary">RECENT</h3>
            <hr />
            <p>col 2</p>
          </div>
          <div class="card">
            <h3 class="primary">CATEGORIES</h3>
            <!-- Start of categories -->
            <b-overlay variant="transparent" :show="deleting" rounded="sm">
              <span v-if="categories.length == 0">None</span>
              <b-container v-else>
                <b-row>
                  <b-col>Name</b-col>
                  <b-col>Posts</b-col>
                </b-row>
                <b-row v-for="c in categories" :key="c.ID" :id="c.ID">
                  <b-col>{{ c.name }}</b-col>
                  <b-col>{{ c.use_count }}</b-col>
                </b-row>
              </b-container>
              <!-- End of categories -->
            </b-overlay>
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

axios.defaults.withCredentials = true;

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
      categories: [],
      deleting: false,
      sorter: { byDate: false },
      error: {
        show: false,
        status: Number,
        message: String,
        callback: Function,
      },
    };
  },
  created() {
    this.getPosts();
    this.getCategories();
  },
  methods: {
    async getPosts() {
      return await axios
        .get("posts")
        .then((response) => {
          this.error.show = false;
          this.posts = response.data.data || [];
        })
        .catch((error) => {
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
    sortPosts() {
      if (this.sorter.byDate) {
        this.posts
          .sort((a, b) => {
            return new Date(b.created) - new Date(a.created);
          })
          .reverse();
      } else {
        this.posts.sort((a, b) => {
          return new Date(b.post.created) - new Date(a.post.created);
        });
      }
      this.sorter.byDate = !this.sorter.byDate;
    },
  },
};
</script>
<style lang="scss" scoped>
/* The hero image */
.hero-image {
  background-image: url("../assets/img/home-hero.jpg");
  box-shadow: 0 5px 6px 2px rgba(10, 10, 10, 0.3);
  height: 25%;
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

.columns {
  display: flex;
  flex-wrap: wrap;
  margin: 0 17%;
}

.columns > * {
  flex-basis: calc(calc(750px - 100%) * 999);
}

.card {
  margin: 20px;
  padding: 10px;
  background-color: rgba(255, 255, 255, 0.05);
  box-shadow: 5px 5px 6px 2px rgba(10, 10, 10, 0.3);
}
.card.post {
  display: flex;
  flex-wrap: wrap;
  flex-direction: row;
}
.card.post.rating-column {
  flex-grow: 1;
}
a.card {
  color: inherit;
  text-decoration: none;
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
</style>