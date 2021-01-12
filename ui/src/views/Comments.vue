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
      <b-row>
        <b-col>
          <b-input-group>
            <b-dropdown size="sm" variant="dark" :text="options.active">
              <b-dropdown-item
                v-for="(option, index) in options.filters"
                :key="index"
                :disabled="options.active === option.text"
                @click="options.active = option.text"
                >{{ option.text }}</b-dropdown-item
              >
            </b-dropdown>
            <b-button
              size="sm"
              variant="dark"
              class="ml-1"
              :disabled="true"
              @click="this.getComments()"
            >
              <b-icon icon="arrow-counterclockwise"></b-icon>
            </b-button>
          </b-input-group>
        </b-col>
      </b-row>
      <!-- Not authenticated-end -->
      <!-- Comments-start -->

      <CommentGroup :comments="[...comments, ...tempComments]" />
      <!-- Comments-end -->
      <b-row no-gutters align-h="center">
        <b-button
          v-if="totalRows > loadedRows"
          size="sm"
          variant="dark"
          @click="getComments()"
          >load {{ totalRows - loadedRows }} more
          {{ (totalRows - loadedRows) % 10 !== 1 ? "comments" : "comment" }}
        </b-button>
      </b-row>
    </div>
  </b-skeleton-wrapper>
</template>
<script>
import CommentsSkeleton from "@/components/skeletons/CommentsSkeleton";
import CommentGroup from "@/components/CommentGroup";
import { mapGetters } from "vuex";
import api from "@/api/api";
import eventBus from "@/event-bus";

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
      tempComments: [],
      totalRows: 0,
      loadedRows: 0,
      limit: 10,
      offset: 0,
      form: { comment: "" },
      maxCommentLength: 200,
      minCommentLength: 1,
      options: {
        active: "Most rated",
        filters: [
          { text: "Most rated", type: "rating" },
          { text: "Least rated", type: "rating", descending: true },
          { text: "Newest", type: "created" },
          { text: "Oldest", type: "created", descending: true },
        ],
      },
      editor: { editing: -1, content: "" },
    };
  },
  created() {
    Promise.all([this.getComments()]).then(() => {
      setTimeout(() => {
        this.loading = false;
      }, 500);
    });
    eventBus.$on("update-event", (event) => {
      this.updateComment(...event);
    });
    eventBus.$on("reply-event", (event) => {
      this.reply(...event);
    });
    eventBus.$on("delete-event", (comment) => {
      this.deleteComment(comment);
    });
  },
  methods: {
    hasPermission(comment) {
      return comment?.author?.id == this.user?.id || this.user?.role > 0;
    },
    async getComments() {
      return await api
        .post("/comments", {
          post_id: this.postID,
          offset: this.offset,
          limit: this.limit,
        })
        .then((response) => {
          this.comments.push(...(response?.data?.data?.comments || []));
          for (var i = this.tempComments.length - 1; i >= 0; i--) {
            for (var j = 0; j < this.comments.length; j++) {
              if (this.tempComments[i]?.id === this.comments[j]?.id) {
                this.tempComments.splice(i, 1);
              }
            }
          }
          this.totalRows = response?.data?.data?.total_rows || 0;
          this.loadedRows += response?.data?.data?.loaded_rows || 0;
          this.offset += this.limit;
        });
    },
    async addComment(content) {
      if (this.requesting || !this.properCommentLength) return;
      this.requesting = true;
      return await api
        .post("comment/add", {
          post_id: this.postID,
          content: content,
        })
        .then((response) => {
          this.form.comment = "";
          if (response?.data?.data) this.tempComments.push(response.data.data);
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
    async reply(parent, content) {
      if (!parent || !content) return;
      return await api
        .post("comment/add", {
          post_id: parent.post_id,
          parent: {
            id: parent.id,
            depth: parent.depth,
            lineage: parent.lineage,
          },
          content: content,
        })
        .then((response) => {
          if (response?.data?.data) {
            if (parent.children) {
              parent.children.push(response.data.data);
              parent.childrenLength++;
            } else {
              this.$set(parent, "children", [response.data.data]);
              this.$set(parent, "children_length", 1);
            }
          }
          parent.replying = false;
        })
        .catch((error) => {
          if (error.status === 403)
            this.$bvToast.toast(
              "You need to be logged in, to reply to comments!",
              {
                title: "Oops!",
                variant: "danger",
                solid: true,
              }
            );
        })
        .then(() => {
          parent.replying = false;
          parent.requesting = false;
        });
    },
    async updateComment(comment, newContent) {
      if (comment.requesting) return;
      this.$set(comment, "requesting", true);
      comment.editing = false;
      return await api
        .put("comment/update", {
          id: comment.id,
          content: newContent,
        })
        .then((response) => {
          if (response?.data?.data) {
            comment.content = response.data.data.content;
            comment.edited = response.data.data.edited;
          }
        })
        .catch((error) => {
          if (error.status === 403)
            this.$bvToast.toast(
              "You need to be logged in, to update comments!",
              {
                title: "Oops!",
                variant: "danger",
                solid: true,
              }
            );
        })
        .then(() => {
          comment.requesting = false;
        });
    },
    async deleteComment(comment) {
      if (comment.requesting) return;
      comment.requesting = true;
      return await api
        .delete("comment/delete", {
          params: { id: comment.id },
        })
        .then(() => {
          this.$set(comment, "deleted", true);
          comment.requesting = false;
        })
        .catch((error) => {
          if (error.status === 403)
            this.$bvToast.toast(
              "You need to be logged in, to delete comments!",
              {
                title: "Oops!",
                variant: "danger",
                solid: true,
              }
            );
        });
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
