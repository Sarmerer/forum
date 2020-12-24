<template>
  <b-skeleton-wrapper :loading="loading">
    <template #loading>
      <CommentsSkeleton />
    </template>
    <div :class="isMobile() ? 'card-m' : 'card'">
      <!-- Comment form-start -->
      <b-form v-if="authenticated" @submit.prevent="addComment(form.comment)">
        <b-input-group class="mt-1 mb-1">
          <b-textarea
            type="text"
            class="textarea"
            :placeholder="
              comments.length == 0
                ? 'Be the first to comment this post'
                : 'What you think?'
            "
            v-model="form.comment"
            @keydown.enter.exact.prevent
            @keyup.enter.exact="addComment(form.comment)"
            keydown.enter.shift.exact="newline"
            rows="1"
            max-rows="3"
            size="sm"
            no-resize
          ></b-textarea>
          <b-input-group-append>
            <transition name="slide-fade">
              <b-button
                v-if="properCommentLength"
                variant="outline-light"
                type="submit"
                :disabled="!properCommentLength"
                >Say</b-button
              >
            </transition>
          </b-input-group-append>
        </b-input-group>
        <div v-if="commentLength > 0">
          <small :style="properCommentLength ? 'color: white' : 'color: red'">
            {{ commentLength }}/{{ maxCommentLength }}
          </small>
        </div>
      </b-form>
      <!-- Comment form-end -->
      <!-- Not authenticated-start -->
      <div v-if="!authenticated" class="border border-dark rounded p-2">
        <b-row>
          <b-col>
            <span>Want to leave a comment?</span>
          </b-col>
          <b-col cols="end" class="mr-3">
            <b-container :tag="isMobile() ? 'div' : 'span'" class="p-0 m-0">
              <b-button
                @click="
                  $router.push({
                    name: 'Auth',
                    params: {
                      signUpPage: true,
                    },
                  })
                "
                variant="success"
                size="sm"
                :class="!isMobile() ? 'mr-2' : ''"
                >Sign Up
              </b-button>
            </b-container>
            <b-container
              :tag="isMobile() ? 'div' : 'span'"
              :class="`p-0 m-0 ${isMobile() ? 'mt-2' : ''}`"
            >
              <b-button
                @click="$router.push('/auth')"
                variant="outline-info"
                size="sm"
                >Sign In
              </b-button>
            </b-container>
          </b-col>
        </b-row>
      </div>
      <!-- Not authenticated-end -->
      <!-- Comments-start -->
      <CommentGroup :comments="comments" :editor="editor" />
      <!-- Comments-end -->
    </div>
  </b-skeleton-wrapper>
</template>
<script>
import CommentsSkeleton from "@/components/skeletons/CommentsSkeleton";
import CommentGroup from "@/components/CommentGroup";
import { mapGetters } from "vuex";
import api from "@/router/api";

export default {
  name: "Comments",
  props: {
    postID: { type: Number, required: true },
  },
  computed: {
    ...mapGetters({
      user: "auth/user",
      authenticated: "auth/authenticated",
    }),

    commentLength() {
      return this.form.comment.replace(/(\r\n|\n|\r|\s)/g, "").length;
    },
    properCommentLength() {
      let cl = this.commentLength;
      return cl >= this.minCommentLength && cl <= this.maxCommentLength;
    },
  },
  components: {
    CommentsSkeleton,
    CommentGroup,
  },
  data() {
    return {
      loading: true,
      comments: [],
      form: { comment: "" },
      maxCommentLength: 200,
      minCommentLength: 5,
      editor: { editing: -1, content: "" },
    };
  },
  created() {
    Promise.all([this.getComments()]).then(() => {
      setTimeout(() => {
        this.loading = false;
      }, 500);
    });
  },
  methods: {
    hasPermission(comment) {
      return comment?.author?.id == this.user?.id || this.user?.role > 0;
    },
    async getComments() {
      return await api
        .get("/comments", {
          params: { id: this.postID },
        })
        .then((response) => {
          this.comments = response.data.data || [];
        });
    },
    async addComment(content, parent = 0) {
      if (this.requesting || !this.properCommentLength) return;
      this.requesting = true;
      return await api
        .post("comment/add", {
          id: this.postID,
          content: content,
          parent: parent,
        })
        .then((response) => {
          this.form.comment = "";
          if (response?.data?.data) this.comments.push(response.data.data);
        })
        .catch((error) => {
          if (error.status === 403)
            this.$bvToast.toast(
              "You need to be logged in, to leave comments!",
              {
                title: "Oops!",
                variant: "danger",
                solid: true,
              }
            );
        })
        .then(() => (this.requesting = false));
    },
  },
};
</script>
<style lang="scss" scoped>
.authorize-button {
  color: #121212;
  background-color: #21e6c1;
  border: none;
}

.authorize-button:hover {
  background-color: #21e6c1;
  opacity: 0.8;
}

.slide-fade-enter-active {
  transition: all 0.5s ease;
}

.slide-fade-leave-active {
  transition: all 0.3s cubic-bezier(1, 0.5, 0.8, 1);
}

.slide-fade-enter,
.slide-fade-leave-to {
  transform: translateX(10px);
  opacity: 0;
}
</style>
