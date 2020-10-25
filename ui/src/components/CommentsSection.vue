<template>
  <div class="card">
    <!-- Comment form -->
    <b-form v-if="authenticated" @submit="leaveComment">
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
          @keyup.enter.exact="leaveComment()"
          keydown.enter.shift.exact="newline"
          rows="1"
          max-rows="3"
          size="sm"
          no-resize
        ></b-textarea>
        <b-input-group-append>
          <transition name="slide-fade">
            <b-button
              v-if="commentLength >= minCommentLength"
              variant="outline-light"
              type="submit"
              :disabled="
                !(
                  commentLength > minCommentLength &&
                  commentLength < maxCommentLength
                )
              "
              >Say</b-button
            >
          </transition>
        </b-input-group-append>
      </b-input-group>
      <div v-if="commentLength > 0">
        <small
          v-if="
            commentLength >= minCommentLength &&
              commentLength <= maxCommentLength
          "
          >{{ commentLength }}/{{ maxCommentLength }}</small
        >
        <small v-else style="color: red"
          >{{ commentLength }}/{{ maxCommentLength }}</small
        >
      </div>
    </b-form>
    <!-- End of comment form -->
    <div v-if="!authenticated" class="border border-dark rounded p-2">
      <span>Want to leave a comment?</span>
      <b-button class="float-right" size="sm" variant="success"
        >Sign Up</b-button
      >
      <b-button class="float-right mr-1" size="sm" variant="outline-primary"
        >Sign In</b-button
      >
    </div>
    <!-- Post comments -->
    <div v-for="(comment, index) in comments" :key="index">
      <div style="margin: 0.3rem; position: relative">
        <div v-if="index != editor.editing">
          <p style="margin-bottom: 0px; display: block; margin-right: 4rem">
            {{ comment.author_name }} says: {{ comment.content }}
          </p>
          <small v-b-tooltip.hover :title="comment.created"
            ><timeago :datetime="comment.created" :auto-update="60"></timeago>
            <small v-if="comment.edited"> edited</small>
          </small>
          <b-button-group
            v-if="
              (user ? comment.author_id == user.id || user.role > 0 : false) &&
                index != editor.editing
            "
            size="sm"
            class="controls-button"
            style="position: absolute; right: 0px; top: 10px"
          >
            <b-button
              size="sm"
              lg="1"
              variant="light"
              class="controls-button"
              :disabled="deletingComment"
              @click="edit(index, comment.content)"
            >
              <img src="@/assets/svg/post/edit.svg" alt="edit" srcset="" />
            </b-button>
            <b-button
              variant="danger"
              :disabled="deletingComment"
              class="controls-button"
              @click="deleteComment(comment.id, index)"
              ><img src="@/assets/svg/post/delete.svg" alt="delete" srcset=""
            /></b-button>
          </b-button-group>
        </div>
        <div
          v-if="
            (user ? comment.author_id == user.id || user.role > 0 : false) &&
              index == editor.editing
          "
        >
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
            style="margin-bottom: 0px; display: block; width: 85%"
          ></b-form-textarea>
          <small
            v-if="
              editorLength > minCommentLength && editorLength < maxCommentLength
            "
            >{{ editorLength }}/{{ maxCommentLength }}</small
          >
          <small v-else style="color: red"
            >{{ editorLength }}/{{ maxCommentLength }}</small
          >
        </div>
        <b-button-group
          size="sm"
          vertical
          v-if="index == editor.editing"
          style="position: absolute; right: 0px; top: 2px"
        >
          <b-button
            :disabled="
              editor.editingContent == comment.content ||
                editorLength < minCommentLength ||
                editorLength > maxCommentLength ||
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
      </div>
    </div>
    <!-- End of post comments -->
  </div>
</template>
<script>
import axios from "axios";

axios.defaults.withCredentials = true;

import { mapGetters } from "vuex";

export default {
  props: {
    postID: { type: Number, required: true },
  },
  computed: {
    ...mapGetters({
      user: "auth/user",
      authenticated: "auth/authenticated",
    }),
    commentLength() {
      return this.form.comment.replace(/(\r\n|\n|\r)/g, "").length;
    },
    editorLength() {
      return this.editor.editingContent.replace(/(\r\n|\n|\r)/g, "").length;
    },
  },
  data() {
    return {
      maxCommentLength: 200,
      minCommentLength: 5,
      editor: { editing: -1, editingContent: "", requesting: false },
      deletingComment: false,
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
    async getComments() {
      return await axios
        .get("/comments", { params: { id: this.postID } })
        .then((response) => {
          this.comments = response.data.data || [];
        })
        .catch((error) => {
          console.log(error);
        });
    },
    async deleteComment(actualID, IDInList) {
      this.editor.editing = -1;
      this.deletingComment = true;
      return await axios
        .delete("comment/delete", { params: { id: actualID } })
        .then(() => {
          this.comments.splice(IDInList, 1);
        })
        .catch((error) => {
          console.log(error);
          //TODO show alert
        })
        .finally(() => (this.deletingComment = false));
    },
    async leaveComment(e) {
      if (e) e.preventDefault();
      if (
        this.commentLength < this.minCommentLength ||
        this.commentLength > this.maxCommentLength
      )
        return;
      this.editor.editing = -1;
      return await axios
        .post("comment/add", { pid: this.postID, content: this.form.comment })
        .then(() => {
          this.appendComment();
          this.form.comment = "";
        })
        .catch((error) => {
          console.log(error);
        });
    },
    async updateComment(actualID, oldCommentContent) {
      if (
        this.editorLength < this.minCommentLength ||
        this.editorLength > this.maxCommentLength
      )
        return;
      if (this.editor.editingContent == oldCommentContent) {
        this.comments[this.editor.editing].content = this.editor.editingContent;
        this.comments[this.editor.editing].edited = true;
        this.editor.editing = -1;
        return;
      }
      this.editor.requesting = true;
      return await axios
        .put("comment/update", {
          id: actualID,
          content: this.editor.editingContent,
        })
        .then(() => {
          console.log(this.comments);
          this.comments[
            this.editor.editing
          ].content = this.editor.editingContent;
          this.comments[this.editor.editing].edited = true;
          this.editor.editing = -1;
        })
        .catch((error) => {
          console.log(error);
        })
        .finally(() => {
          this.editor.requesting = false;
        });
    },
    appendComment() {
      let comment = {
        id: this.comments.length != 0 ? this.comments[0].id + 1 : 1,
        author_id: this.user.id,
        author_name: this.user.display_name,
        created: Date.now(),
        content: this.form.comment,
        post: this.postID,
        edited: false,
      };
      this.comments = [comment, ...this.comments];
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
</style>
