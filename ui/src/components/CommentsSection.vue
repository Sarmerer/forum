<template>
  <div class="card">
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
      <span>Want to leave a comment?</span>
      <b-button
        v-b-modal.signup-modal
        class="float - right"
        variant="success"
        size="sm"
        >Sign Up
      </b-button>
      <b-button
        v-b-modal.auth-modal
        class="float-right mr-1"
        variant="outline-info"
        size="sm"
        >Sign In
      </b-button>
    </div>
    <!-- Not authenticated-end -->
    <!-- Comment form-end -->

    <!-- Comments-start -->
    <div v-for="(comment, index) in comments" :key="index">
      <div style="margin: 0.3rem; position: relative">
        <Rating :callback="rate" :entity="comment" type="comment" />
        <div v-if="index != editor.editing">
          <b-img :src="comment.author.avatar" width="15px"></b-img>
          <p style="margin-bottom: 0px; display: block; margin-right: 4rem">
            <router-link :to="`/user/${comment.author.id}`">
              {{ comment.author.display_name }}</router-link
            >
            says: {{ comment.content }}
          </p>
          <small v-b-tooltip.hover :title="comment.created">
            <timeago
              v-if="comment.created"
              :datetime="comment.created"
              :auto-update="60"
            >
            </timeago>
            <small v-if="comment.edited" class="text-muted"> edited</small>
          </small>
          <ControlButtons
            :hasPermission="hasPermission(comment)"
            :deleteCallback="{
              callback: deleteComment,
              args: [comment.id, index],
            }"
            :editCallback="{ callback: edit, args: [index, comment.content] }"
            :disabled="deleting"
          />
        </div>
        <div v-if="hasPermission(comment) && index == editor.editing">
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
          <small v-if="properEditorLength"
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
      </div>
    </div>
    <!-- Comments-end -->
  </div>
</template>
<script>
import ControlButtons from "@/components/ControlButtons";
import Rating from "@/components/Rating";
import { mapGetters } from "vuex";
import axios from "axios";
axios.defaults.withCredentials = true;

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
    properCommentLength() {
      let cl = this.commentLength;
      return cl >= this.minCommentLength && cl <= this.maxCommentLength;
    },
    editorLength() {
      return this.editor.editingContent.replace(/(\r\n|\n|\r)/g, "").length;
    },
    properEditorLength() {
      let el = this.editorLength;
      return el >= this.minCommentLength && el <= this.maxCommentLength;
    },
  },
  components: {
    ControlButtons,
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
      return await axios
        .get("/comments", { params: { id: this.postID } })
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
      return await axios
        .delete("comment/delete", { params: { id: actualID } })
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
      return await axios
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
      await axios
        .post("comment/rate", { id: comment.id, reaction: r })
        .then((response) => {
          comment.rating = response.data.data.rating;
          comment.your_reaction = response.data.data.your_reaction;
        })
        .catch((error) => {
          console.log(error);
          //TODO show alert saying that you need to be logged in to rate
        })
        .finally(() => {
          self.requesting = false;
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
</style>
