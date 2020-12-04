<template>
  <div :class="isMobile() ? 'card-m' : 'card'">
    <!-- Comment form-start -->
    <b-form v-if="authenticated" @submit.prevent="addComment">
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
          @keyup.enter.exact="addComment()"
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
                    prevRoute: $router.currentRoute.path,
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
              @click="
                $router.push({
                  name: 'Auth',
                  params: { prevRoute: $router.currentRoute.path },
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
    <!-- Not authenticated-end -->
    <!-- Comment form-end -->

    <!-- Comments-start -->
    <div
      v-for="(comment, index) in comments"
      :key="index"
      class="ml-3 mr-3 mt-2"
    >
      <ControlModal
        v-on:edit-event="edit(index, comment.content)"
        v-on:delete-event="deleteComment(comment.id, index)"
        :modalID="'modal-menu' + index"
      />
      <div>
        <div v-if="index != editor.editing">
          <b-row>
            <b-col cols="start">
              <router-link :to="`/user/${comment.author.id}`" class="secondary">
                <user-avatar
                  :userData="comment.author"
                  :popoverID="'c' + index"
                >
                </user-avatar>
              </router-link>
              <span
                v-if="comment.author.id === postAuthorID"
                v-b-tooltip.hover.right
                title="Post author"
              >
                <b-icon-pencil width="13px"></b-icon-pencil>
              </span>
            </b-col>
            <b-col>
              <small>
                <time-ago
                  :datetime="comment.created"
                  :long="!isMobile()"
                  tooltip="right"
                ></time-ago>
                <small
                  v-b-tooltip.hover
                  :title="comment.edited"
                  v-if="comment.edited != comment.created"
                  class="text-muted"
                >
                  edited</small
                >
              </small>
            </b-col>
            <b-col cols="end" class="mr-n2">
              <ControlButtons
                :hasPermission="hasPermission(comment)"
                v-on:delete-event="deleteComment(comment.id, index)"
                v-on:edit-event="edit(index, comment.content)"
                :disabled="deleting"
                :compact="isMobile()"
                :modalID="'modal-menu' + index"
              />
            </b-col>
          </b-row>
          <b-row class="mt-n2">
            <span>
              {{ comment.content }}
            </span>
          </b-row>
          <b-row>
            <Rating :callback="rate" :entity="comment" compact />
          </b-row>
        </div>
        <div v-if="hasPermission(comment) && index == editor.editing">
          <b-row>
            <b-col class="ml-n3">
              <b-form-textarea
                class="textarea"
                ref="editComment"
                v-model="editor.editingContent"
                @keydown.enter.exact.prevent
                @keyup.enter.exact="updateComment(comment.id)"
                keydown.enter.shift.exact="newline"
                rows="1"
                no-resize
                :disabled="editor.requesting"
                max-rows="10"
              ></b-form-textarea>
              <small v-if="properEditorLength"
                >{{ editorLength }}/{{ maxCommentLength }}</small
              >
              <small v-else style="color: red"
                >{{ editorLength }}/{{ maxCommentLength }}</small
              >
            </b-col>
            <b-col cols="end">
              <b-button-group size="sm" vertical v-if="index == editor.editing">
                <b-button
                  :disabled="
                    editor.editingContent == comment.content ||
                      !properEditorLength ||
                      editor.requesting
                  "
                  variant="outline-success"
                  @click="updateComment(comment.id, comment.content)"
                >
                  Save
                </b-button>
                <b-button
                  variant="outline-danger"
                  @click="editor.editing = -1"
                  :disabled="editor.requesting"
                  >Cancel</b-button
                >
              </b-button-group>
            </b-col>
          </b-row>
        </div>
      </div>
      <hr v-if="index != comments.length - 1" class="comment-divider" />
    </div>
    <!-- Comments-end -->
  </div>
</template>
<script>
import ControlButtons from "@/components/ControlButtons";
import UserAvatar from "@/components/UserAvatar";
import ControlModal from "./ControlModal";
import Rating from "@/components/Rating";
import TimeAgo from "vue2-timeago";
import { mapGetters } from "vuex";
import api from "@/router/api";

export default {
  props: {
    postID: { type: Number, required: true },
    postAuthorID: Number,
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
    editorLength() {
      return this.editor.editingContent.replace(/(\r\n|\n|\r|\s)/g, "").length;
    },
    properEditorLength() {
      let el = this.editorLength;
      return el >= this.minCommentLength && el <= this.maxCommentLength;
    },
  },
  components: {
    ControlButtons,
    ControlModal,
    UserAvatar,
    TimeAgo,
    Rating,
  },
  data() {
    return {
      maxCommentLength: 200,
      minCommentLength: 5,
      requesting: false,
      editor: { editing: -1, editingContent: "", requesting: false },
      deleting: false,
      comments: [],
      form: {
        comment: "",
      },
    };
  },
  created() {
    this.getComments();
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
        })
        .catch((error) => {
          console.log(error.response.data);
        });
    },
    async deleteComment(actualID, IDInList) {
      this.editor.editing = -1;
      this.deleting = true;
      return await api
        .delete("comment/delete", {
          params: { id: actualID },
        })
        .then(() => {
          this.comments.splice(IDInList, 1);
        })
        .catch((error) => {
          console.log(error.response.data);
          //TODO show alert
        })
        .then(() => (this.deleting = false));
    },
    async addComment() {
      if (!this.properCommentLength) return;
      this.editor.editing = -1;
      return await api
        .post("comment/add", { id: this.postID, content: this.form.comment })
        .then((response) => {
          if (response.data.data)
            this.comments = [response.data.data, ...this.comments];
          this.form.comment = "";
        })
        .catch((error) => {
          console.log(error.response.data);
        });
    },
    async updateComment(actualID, oldCommentContent) {
      if (!this.properEditorLength) return;
      if (this.editor.editingContent == oldCommentContent) {
        this.comments[this.editor.editing].content = this.editor.editingContent;
        this.comments[this.editor.editing].edited = true;
        this.editor.editing = -1;
        return;
      }
      this.editor.requesting = true;
      return await api
        .put("comment/update", {
          id: actualID,
          content: this.editor.editingContent,
        })
        .then((response) => {
          if (response.data.data)
            this.comments[this.editor.editing] = response.data.data;
          this.editor.editing = -1;
        })
        .catch((error) => {
          console.log(error);
        })
        .then(() => {
          this.editor.requesting = false;
        });
    },
    edit(index, content) {
      this.editor.editing = index;
      this.editor.editingContent = content;
      let interval = setInterval(() => {
        try {
          this.$refs.editComment[0].focus();
          clearInterval(interval);
        } catch (error) {
          console.log(error);
        }
      }, 10);
    },
    async rate(reaction, comment) {
      if (this.requesting) return;
      this.requesting = true;
      let r = reaction == "up" ? 1 : -1;
      if (
        (reaction == "up" && comment.your_reaction == 1) ||
        (reaction == "down" && comment.your_reaction == -1)
      ) {
        r = 0;
      }
      await api
        .post("comment/rate", { id: comment.id, reaction: r })
        .then((response) => {
          comment.rating = response.data.data.rating;
          comment.your_reaction = response.data.data.your_reaction;
        })
        .catch((error) => {
          console.log(error);
          //TODO show alert saying that you need to be logged in to rate
        })
        .then(() => {
          this.requesting = false;
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
.controls-button {
  background-color: transparent;
  border-color: transparent;
  outline: none !important;
  outline-width: 0 !important;
  box-shadow: none;
  -moz-box-shadow: none;
  -webkit-box-shadow: none;
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

.authorize-button {
  color: #121212;
  background-color: #21e6c1;
  border: none;
}

.authorize-button:hover {
  background-color: #21e6c1;
  opacity: 0.8;
}

.comment-divider {
  margin: 5px -10px 5px;
  border: 0;
  border-top: 1px solid #121212;
}
</style>
