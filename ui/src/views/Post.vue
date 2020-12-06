<template>
  <b-skeleton-wrapper :loading="loading">
    <template #loading>
      <PostSkeleton />
    </template>
    <div class="grid">
      <ControlModal
        v-on:edit-event="edit()"
        v-on:delete-event="deletePost()"
        modalID="modal-menu"
      />
      <div class="columns">
        <div v-if="isMobile()" class="info-col">
          <UserCard v-if="post.author" link :userData="post.author" />
        </div>
        <div class="main-col">
          <div :class="`${isMobile() ? 'card-m' : 'card'}`">
            <b-row v-if="!editor.editing">
              <b-col cols="start">
                <Rating :callback="rate" :entity="post" class="ml-n4" />
              </b-col>
              <b-col class="ml-2">
                <b-row>
                  <b-col>
                    <h3 class="primary text-break">{{ post.title }}</h3>
                  </b-col>
                  <b-col cols="end" class="mr-2">
                    <ControlButtons
                      :hasPermission="hasPermission"
                      v-on:delete-event="deletePost()"
                      v-on:edit-event="edit()"
                      :disabled="requesting"
                      :compact="isMobile()"
                      modalID="modal-menu"
                    />
                  </b-col>
                </b-row>
                <pre color="white" class="text-break">{{ post.content }}</pre>
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
            </b-row>
            <b-row v-if="isMobile() && !editor.editing" class="ml-2">
              <Rating :callback="rate" :entity="post" compact />
            </b-row>
            <b-form v-if="editor.editing" @submit.prevent="updatePost()">
              <b-form-row>
                <b-col>
                  <PostForm
                    :form="editor"
                    v-on:valid-form="editor.valid = $event"
                  />
                </b-col>
              </b-form-row>
              <b-form-row class="mt-2">
                <b-col align="end">
                  <b-button-group size="sm" v-if="!editor.confirmCancel">
                    <b-button
                      variant="outline-danger"
                      @click="editor.confirmCancel = true"
                    >
                      Cancel
                    </b-button>
                    <b-button
                      variant="outline-success"
                      :disabled="!editor.valid"
                      type="submit"
                      class="px-3"
                    >
                      Save
                    </b-button>
                  </b-button-group>
                </b-col>
                <b-col align="end" v-if="editor.confirmCancel">
                  <p class="m-0">Cancel editor?</p>
                  <b-button-group size="sm">
                    <b-button
                      variant="outline-danger"
                      @click="editor.confirmCancel = false"
                    >
                      <b-icon-x></b-icon-x> No
                    </b-button>
                    <b-button
                      variant="outline-success"
                      @click="editor.editing = false"
                    >
                      <b-icon-check2></b-icon-check2> Yes
                    </b-button>
                  </b-button-group>
                </b-col>
              </b-form-row>
            </b-form>
          </div>
          <CommentsSection
            v-if="post.author && comments"
            :comments="comments"
            :postID="postID"
            :postAuthorID="post.author.id"
          />
        </div>
        <div v-if="!isMobile()" class="info-col">
          <UserCard v-if="post.author" link :userData="post.author" />
        </div>
      </div>
    </div>
  </b-skeleton-wrapper>
</template>
<script>
import PostSkeleton from "@/components/skeletons/PostSkeleton";
import CommentsSection from "@/components/CommentsSection";
import ControlButtons from "@/components/ControlButtons";
import ControlModal from "@/components/ControlModal";
import PostForm from "@/components/forms/PostForm";
import UserCard from "@/components/UserCard";
import Rating from "@/components/Rating";
import { mapGetters } from "vuex";
import api from "@/router/api";

export default {
  props: {
    postData: { type: Object },
  },
  components: {
    CommentsSection,
    ControlButtons,
    PostSkeleton,
    ControlModal,
    PostForm,
    UserCard,
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
      categories: [],
      comments: [],
      editor: {
        title: "",
        content: "",
        categories: [],
        editing: false,
        valid: false,
        confirmCancel: false,
      },
      postID: Number.parseInt(this.$route.params.id),
      requesting: false,
      loading: true,
    };
  },
  created() {
    if (this.postData) {
      document.title = this.postData.title;
      this.post = this.postData;
      Promise.all([this.getComments()]).then(() => {
        setTimeout(() => {
          this.loading = false;
        }, 500);
      });
    } else {
      let p = this.getPost();
      let p1 = this.getComments();
      Promise.all([p, p1]).then(() => {
        setTimeout(() => {
          this.loading = false;
        }, 500);
      });
    }
  },
  methods: {
    async getPost() {
      return await api
        .post("post/find", {
          by: "id",
          id: this.postID,
        })
        .then((response) => {
          let result = response?.data?.data;
          document.title = result.title;
          this.post = result;
        })
        .catch(() => {
          this.$router.push("/");
        });
    },
    async getComments() {
      return await api
        .get("/comments", {
          params: { id: this.postID },
        })
        .then((response) => {
          this.comments = response.data.data || [];
        })
        .catch((error) => {
          console.log(error.response.data);
        });
    },
    async deletePost() {
      this.requesting = true;
      return await api
        .delete("post/delete", {
          params: { id: this.post.id },
        })
        .then(() => {
          this.$router.push("/");
        })
        .finally(() => {
          this.requesting = false;
        });
    },
    async updatePost() {
      this.requesting = true;
      return await api
        .put("post/update", {
          id: this.post.id,
          title: this.editor.title,
          content: this.editor.content,
          categories: this.editor.categories,
        })
        .then((response) => {
          this.post = response.data.data;
        })
        .catch((error) => {
          if (error.status === 403)
            this.$bvToast.toast("You need to be logged in, to update posts!", {
              title: "Oops!",
              variant: "danger",
              solid: true,
            });
        })
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
      await api
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
