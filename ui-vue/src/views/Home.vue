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
              >sort by date: {{ sorter.byDate ? "ascending" : "descending" }}</b-button
            >
          </div>
          <!-- Start of posts -->
          <router-link
            :to="'/post/' + post.post.id"
            v-for="(post, index) in posts"
            :key="index"
            class="card post"
            tag="div"
            style="cursor: pointer"
          >
            <div class="rating-column mr-2" style="text-align:center;">
              <svg
                style="display:block"
                @click="rate($event, 'up', index)"
                :disabled="rating.requesting"
                width="1.5em"
                height="1.5em"
                viewBox="0 0 16 16"
                class="bi bi-chevron-up"
                :fill="post.post.your_reaction == 1 ? 'green' : 'white'"
                xmlns="http://www.w3.org/2000/svg"
              >
                <path
                  fill-rule="evenodd"
                  d="M7.646 4.646a.5.5 0 0 1 .708 0l6 6a.5.5 0 0 1-.708.708L8 5.707l-5.646 5.647a.5.5 0 0 1-.708-.708l6-6z"
                />
              </svg>
              <span>{{ post.post.rating }}</span>
              <span
                ><svg
                  style="display:block"
                  @click="rate($event, 'down', index)"
                  :disabled="rating.requesting"
                  width="1.5em"
                  height="1.5em"
                  viewBox="0 0 16 16"
                  class="bi bi-chevron-down"
                  :fill="post.post.your_reaction == -1 ? 'red' : 'white'"
                  xmlns="http://www.w3.org/2000/svg"
                >
                  <path
                    fill-rule="evenodd"
                    d="M1.646 4.646a.5.5 0 0 1 .708 0L8 10.293l5.646-5.647a.5.5 0 0 1 .708.708l-6 6a.5.5 0 0 1-.708 0l-6-6a.5.5 0 0 1 0-.708z"
                  /></svg
              ></span>
            </div>
            <div style="max-width: 95%">
              <small
                >by
                <router-link :to="'/user/' + post.post.author_id" style="text-decoration: none;">
                  {{ post.post.author_name }}
                </router-link>
                <timeago :datetime="post.post.created" :auto-update="60"></timeago
              ></small>
              <h2 class="primary">
                {{ post.post.title }}
              </h2>
              <hr />
              <p style="color: white">{{ post.post.content }}</p>
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
              <sub
                ><img src="@/assets/svg/post/comments.svg" alt="comments" srcset="" />
                {{ post.replies }} replies</sub
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

export default {
  name: "Home",
  computed: {
    ...mapGetters({
      user: "auth/user",
    }),
  },
  components: {
    Error,
  },
  data() {
    return {
      posts: [],
      categories: [],
      deleting: false,
      rating: { requesting: false },
      sorter: { byDate: false },
      error: { show: false, status: Number, message: String, callback: Function },
    };
  },
  created() {
    this.getPosts();
    this.getCategories();
  },
  methods: {
    async rate(e, reaction, index) {
      e.preventDefault();
      if (this.rating.requesting) return;
      this.rating.requesting = true;
      let self = this;
      let r = reaction == "up" ? 1 : -1;
      if (
        (reaction == "up" && this.posts[index].post.your_reaction == 1) ||
        (reaction == "down" && this.posts[index].post.your_reaction == -1)
      ) {
        r = 0;
      }
      await axios
        .post("post/rate", { pid: self.posts[index].post.id, reaction: r })
        .then((response) => {
          self.posts[index].post.your_reaction = response.data.data.your_reaction;
          self.posts[index].post.rating = response.data.data.rating;
        })
        .catch((error) => {
          console.log(error);
          //TODO show alert saying that you need to be logged in to rate
        })
        .finally(() => {
          self.rating.requesting = false;
        });
    },
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
