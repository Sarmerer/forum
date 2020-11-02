<template>
  <div class="home ">
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
          <div v-if="posts.length > 0" class="card">
            <b-button @click="sortDisplayPosts()"
              >sort by date:
              {{ sorter.byDate ? "ascending" : "descending" }}</b-button
            >
          </div>
          <!-- Start of posts -->
          <router-link
            :to="'/post/' + post.id"
            v-for="(post, index) in posts"
            :key="index"
            class="card"
            tag="div"
            style="cursor: pointer;"
          >
            <b-row>
              <b-col cols="1">
                <Rating
                  v-if="!isMobile()"
                  v-on:update="
                    args => {
                      post.rating = args.new_rating;
                      post.your_reaction = args.new_your_reaction;
                    }
                  "
                  :postID="post.id"
                  :rating="post.rating"
                  :yourReaction="post.your_reaction"
                />
              </b-col>
              <b-col cols="11" class="post-content">
                <h5>
                  {{ post.title }}
                </h5>
                <p>{{ post.content }}</p>
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

              <!-- TO-DO: Make this look decent -->
              <!-- style is embedded here for responsiveness. MB fix later -->
              <small v-if="isMobile()"
                ><Rating style="flex-direction:row; margin:0;"
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
              v-for="(post, index) in posts.slice(0, 5)"
              :key="index"
              tag="div"
              ><p>{{ post.title }}</p>
              <small
                ><timeago
                  :datetime="post.created"
                  :auto-update="10"
                ></timeago></small
            ></router-link>
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

export default {
  name: "Home",
  computed: {
    ...mapGetters({
      user: "auth/user"
    })
  },
  components: {
    Rating,
    Error
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
        callback: Function
      }
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
        .then(response => {
          this.error.show = false;
          this.posts = response.data.data || [];
        })
        .catch(error => {
          this.error.show = true;
          this.error.status = error.response.status;
          this.error.message = error.response.statusText;
          this.error.callback = this.getPosts;
        });
    },
    async getCategories() {
      return await axios
        .get("categories")
        .then(response => {
          this.categories = response.data.data || [];
        })
        .catch(error => {
          console.log(error);
        });
    },
    sortDisplayPosts() {
      if (this.sorter.byDate) {
        this.posts
          .sort((a, b) => {
            return new Date(b.created) - new Date(a.created);
          })
          .reverse();
      } else {
        this.posts.sort((a, b) => {
          return new Date(b.created) - new Date(a.created);
        });
      }
      this.sorter.byDate = !this.sorter.byDate;
    }
  }
};
</script>
<style lang="scss" scoped>
/* The hero image */
.hero-image {
  background-image: url("../assets/img/home-hero.jpg");
  box-shadow: 0 5px 6px 2px rgba(10, 10, 10, 0.3);
  height: 200px;
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

@media (max-width: 500px) {
  .post-content {
    margin-left: -20px;
  }
}
</style>
