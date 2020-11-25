<template>
  <div class="grid">
    <div class="columns">
      <div v-if="isMobile()" class="info-col">
        <PostStats :stats="postStats" />
      </div>
      <div class="main-col">
        <div class="card">
          <b-row>
            <b-col cols="1">
              <Rating :callback="rate" :entity="post" type="comment" />
            </b-col>
            <b-col cols="11" style="margin-left: -35px">
              <h3 class="primary">{{ post.title }}</h3>
              <p color="white">{{ post.content }}</p>
              <div>
                <b-form-tag
                  v-for="category in post.categories"
                  disabled
                  :key="category.id"
                  :title="category.name"
                  variant="dark"
                  class="mr-1 mb-1"
                >
                  {{ category.name }}
                </b-form-tag>
                <div class="controls">
                  <ControlButtons
                    :hasPermission="hasPermission"
                    :deleteCallback="{ callback: deletePost }"
                    :editCallback="{}"
                    :disabled="deleting"
                  />
                </div></div
            ></b-col>
          </b-row>
        </div>
        <CommentsSection :postID="postID" />
      </div>
      <div v-if="!isMobile()" class="info-col">
        <PostStats :stats="postStats" />
      </div>
    </div>
  </div>
</template>
<script>
import CommentsSection from "@/components/CommentsSection";
import ControlButtons from "@/components/ControlButtons";
import PostStats from "@/components/PostStats";
import Rating from "@/components/Rating";
import { mapGetters } from "vuex";
import axios from "axios";

export default {
  props: {
    postData: { type: Object, required: false },
  },
  components: {
    CommentsSection,
    ControlButtons,
    PostStats,
    Rating,
  },
  computed: {
    ...mapGetters({
      user: "auth/user",
    }),
    hasPermission() {
      if (!this.user) return false;
      return this.post.author_id == this.user.id || this.user.role > 0;
    },
  },
  data() {
    return {
      post: {},
      postStats: {},
      categories: [],
      postID: Number.parseInt(this.$route.params.id),
      deleting: false,
    };
  },
  created() {
    this.postData ? (this.post = this.postData) : this.getPost();
  },
  methods: {
    async getPost() {
      return await axios
        .post("post/find", {
          by: "id",
          id: this.postID,
        })
        .then((response) => {
          //TODO create error page for post response.data.data;
          let result = response.data.data;
          this.post = result;
          this.postStats = {
            commentsCount: result.comments_count,
            participantsCount: result.participants_count,
            lastCommentFromID: result.last_comment_from_id,
            lastCommentFromName: result.last_comment_from_name,
            lastCommentDate: result.last_comment_date,
          };
        })
        .catch((error) => {
          console.log(error);
          this.$router.push("/");
        });
    },
    async deletePost() {
      this.deleting = true;
      return await axios
        .delete("post/delete", { params: { id: this.post.id } })
        .then(() => {
          this.$router.push("/");
        })
        .catch(() => {
          // TODO show error notification
        })
        .finally(() => {
          this.deleting = false;
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
        .post("post/rate", { id: this.postID, reaction: r })
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
.controls {
  position: absolute;
  top: 5px;
  right: 10px;
}
</style>
