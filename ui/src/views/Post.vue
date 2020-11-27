<template>
  <div class="grid">
    <div class="columns">
      <div v-if="isMobile()" class="info-col">
        <PostStats :stats="postStats" />
      </div>
      <div class="main-col">
        <div class="card">
          <b-row v-if="!editor.editing">
            <b-col cols="start">
              <Rating :callback="rate" :entity="post" type="post" />
            </b-col>
            <b-col class="ml-2">
              <b-row>
                <b-col>
                  <h3 class="primary">{{ post.title }}</h3>
                </b-col>
                <b-col cols="end">
                  <ControlButtons
                    :hasPermission="hasPermission"
                    :deleteCallback="{ callback: deletePost }"
                    :editCallback="{ callback: edit }"
                    :disabled="requesting"
                  />
                </b-col>
              </b-row>
              <p color="white">{{ post.content }}</p>
              <div>
                <b-form-tag
                  v-for="category in post.categories"
                  disabled
                  :key="category.id"
                  :title="category.name"
                  variant="dark"
                  class="mr-1 mb-1"
                  >{{ category.name }}
                </b-form-tag>
              </div>
            </b-col>
            <b-col cols="end" align-v="start" class="mr-3" display="flex">
            </b-col>
          </b-row>
          <b-form v-if="editor.editing">
            <b-form-row>
              <b-col>
                <b-form-group label-for="title">
                  <b-form-input
                    id="title"
                    class="mb-3"
                    v-model="editor.title"
                    autocomplete="off"
                    placeholder="Enter title"
                  ></b-form-input>
                  <b-form-textarea
                    id="textarea-auto-height"
                    v-model="editor.content"
                    rows="1"
                    max-rows="10"
                  >
                  </b-form-textarea>
                </b-form-group>
                <b-form-tags
                  input-id="tags-basic"
                  remove-on-delete
                  v-model="editor.categories"
                  tag-variant="dark"
                ></b-form-tags>
              </b-col>
              <b-col cols="end">
                <b-button-group size="sm" vertical>
                  <b-button variant="outline-success" @click="updatePost()">
                    Save
                  </b-button>
                  <b-button
                    variant="outline-danger"
                    @click="editor.editing = false"
                  >
                    Cancel
                  </b-button>
                </b-button-group>
              </b-col>
            </b-form-row>
          </b-form>
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
axios.defaults.withCredentials = true;

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
      return this.post?.author?.id == this.user?.id || this.user?.role > 0;
    },
  },
  data() {
    return {
      post: {},
      postStats: {},
      categories: [],
      editor: {
        title: "",
        content: "",
        categories: [],
        editing: false,
      },
      postID: Number.parseInt(this.$route.params.id),
      requesting: false,
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
      this.requesting = true;
      return await axios
        .delete("post/delete", { params: { id: this.post.id } })
        .then(() => {
          this.$router.push("/");
        })
        .catch(console.log)
        .finally(() => {
          this.requesting = false;
        });
    },
    async updatePost() {
      this.requesting = true;
      return await axios
        .put("post/update", {
          id: this.post.id,
          title: this.editor.title,
          content: this.editor.content,
          categories: this.editor.categories,
        })
        .then((response) => {
          this.post = response.data.data;
        })
        .catch(console.log)
        .then(() => {
          this.requesting = false;
          this.editor.editing = false;
        });
    },
    edit() {
      this.editor.title = this.post.title;
      this.editor.content = this.post.content;
      this.editor.categories = this.post.categories
        ? this.post.categories.map((c) => c.name)
        : [];
      this.editor.editing = true;
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
