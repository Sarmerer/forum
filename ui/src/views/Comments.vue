<template>
  <b-skeleton-wrapper :loading="loading">
    <template #loading>
      <CommentsSkeleton />
    </template>
    <div :class="isMobile() ? 'card-m' : 'card'">
      <!-- Comment form-start -->
      <b-form v-if="authenticated" @submit.prevent="addComment(form.comment)">
        <b-input-group class="mt-1">
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
                      prevRoute: $route.path,
                    },
                  })
                "
                variant="success"
                size="sm"
                :class="{ 'mr-2': !isMobile() }"
                >Sign Up
              </b-button>
            </b-container>
            <b-container
              :tag="isMobile() ? 'div' : 'span'"
              class="p-0 m-0"
              :class="{ 'mt-2': isMobile() }"
            >
              <b-button
                @click="
                  $router.push({
                    name: 'Auth',
                    params: {
                      prevRoute: $route.path,
                    },
                  })
                "
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
          <b-input-group v-if="comments.length > 1" class="mt-1">
            <b-dropdown
              :disabled="requesting"
              size="sm"
              variant="dark"
              :text="options.active.text"
            >
              <b-dropdown-item
                v-for="(option, index) in options.filters"
                :key="index"
                :disabled="options.active.text === option.text"
                @click="order(option)"
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

      <CommentGroup :comments="comments" />
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
    comments() {
      return [...this.fetchedComments, ...this.tempComments].filter(
        (c) => !c.deleted
      );
    },
  },
  components: {
    CommentsSkeleton,
    CommentGroup,
  },
  data() {
    return {
      loading: true,
      fetchedComments: [],
      requesting: false,
      tempComments: [],
      totalRows: 0,
      loadedRows: 0,
      limit: 10,
      offset: 0,
      ordered: false,
      form: { comment: "" },
      maxCommentLength: 200,
      minCommentLength: 1,
      options: {
        active: { text: "Most rated", type: "rating", direction: "desc" },
        filters: [
          { text: "Most rated", type: "rating", direction: "desc" },
          { text: "Least rated", type: "rating", direction: "asc" },
          { text: "Newest", type: "created", direction: "desc" },
          { text: "Oldest", type: "created", direction: "asc" },
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
    async order(option) {
      if (this.requesting) return;
      this.requesting = true;
      Object.assign(this.options.active, option);
      return await api
        .post("/comments", {
          post_id: this.postID,
          offset: 0,
          limit: this.limit,
          order_by: this.options.active.type,
          direction: this.options.active.direction,
        })
        .then((response) => {
          this.tempComments = [];
          this.fetchedComments = response?.data?.data?.comments || [];
          this.offset = 0;
          this.totalRows = response?.data?.data?.total_rows || 0;
          this.loadedRows = response?.data?.data?.loaded_rows || 0;
          this.requesting = false;
          this.ordered = true;
        });
    },
    async getComments() {
      if (this.requesting) return;
      this.requesting = true;
      return await api
        .post("/comments", {
          post_id: this.postID,
          offset: this.offset,
          limit: this.limit,
          order_by: this.options.active.type,
          direction: this.options.active.direction,
        })
        .then((response) => {
          this.ordered
            ? (this.fetchedComments = response?.data?.data?.comments || [])
            : this.fetchedComments.push(
                ...(response?.data?.data?.comments || [])
              );
          this.ordered = false;
          for (var i = this.tempComments.length - 1; i >= 0; i--) {
            for (var j = 0; j < this.fetchedComments.length; j++) {
              if (this.tempComments[i]?.id === this.fetchedComments[j]?.id) {
                this.tempComments.splice(i, 1);
              }
            }
          }
          this.offset += this.limit;
          this.totalRows = response?.data?.data?.total_rows || 0;
          this.loadedRows += response?.data?.data?.loaded_rows || 0;
          this.requesting = false;
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
      if (parent.requesting || !parent || !content) return;
      this.$set(parent, "requesting", true);
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
      this.$set(comment, "requesting", true);
      return await api
        .delete("comment/delete", {
          params: { id: comment.id },
        })
        .then(() => {
          this.$set(comment, "deleted", true);
          this.offset--;
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
        })
        .then(() => (comment.requesting = false));
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
