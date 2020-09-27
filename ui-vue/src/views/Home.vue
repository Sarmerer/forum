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
          <!-- Start of posts -->
          <div class="card" v-for="(post, index) in posts" :key="index">
            <b-link>{{ post.post.AuthorName }}</b-link>
            <h2 class="primary">
              <router-link :to="'/post/' + post.post.ID" style="text-decoration: none;">
                {{ post.post.Title }}
              </router-link>
            </h2>
            <hr />
            <p style="color: white">{{ post.post.Content }}</p>
            <sub
              ><img src="@/assets/svg/post/comments.svg" alt="comments" srcset="" />
              {{ post.replies }} replies</sub
            >
          </div>
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
              <div v-if="!categories"><p>None</p></div>
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
export default {
  name: "Home",
  computed: {
    ...mapGetters({
      user: "auth/user",
    }),
  },
  data() {
    return {
      posts: [],
      categories: [],
      deleting: false,
    };
  },
  mounted() {
    this.getPosts();
    this.getCategories();
  },
  methods: {
    async getPosts() {
      return await axios.get("posts").then((response) => (this.posts = response.data.data));
    },
    async getCategories() {
      return await axios
        .get("categories")
        .then((response) => (this.categories = response.data.data));
    },
  },
};
</script>
<style lang="scss">
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
