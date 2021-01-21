<template>
  <div>
    <div
      :class="
        `user-card text-break ${isMobile() ? 'card-m' : 'card card-hover'}`
      "
      v-for="comment in filteredComments"
      :key="comment.id"
    >
      <ControlModal
        v-if="isMobile() && hasPermission(comment.author)"
        v-on:delete-event="deleteComment(comment)"
        v-on:edit-event="
          $set(comment, 'editing', true),
            $set(comment, 'editor', comment.content)
        "
        :modalID="'modal-menu' + comment.id"
      />
      <b-overlay
        :show="comment.requesting"
        rounded
        opacity="0.6"
        spinner-small
        variant="dark"
        spinner-variant="light"
      >
        <template #overlay>
          <div class="text-center">
            <b-icon icon="stopwatch" font-scale="2" animation="cylon"></b-icon>
            <p id="cancel-label">Please wait...</p>
          </div>
        </template>
        <div v-if="!comment.editing">
          <b-row no-gutters>
            <b-col>
              <router-link :to="'/post/' + comment.post_id" tag="h5">
                {{ comment.content }}
              </router-link>
            </b-col>
            <b-col cols="end">
              <ControlButtons
                :hasPermission="hasPermission(comment.author)"
                v-on:delete-event="deleteComment(comment)"
                v-on:edit-event="
                  $set(comment, 'editing', true),
                    $set(comment, 'editor', comment.content)
                "
                :disabled="false"
                :compact="isMobile()"
                :modalID="'modal-menu' + comment.id"
              />
            </b-col>
          </b-row>
          <small>
            <span v-b-tooltip.hover title="Rating">
              <b-icon
                :icon="reactionIcon(comment.your_reaction)"
                :color="reactionColor(comment.your_reaction)"
              >
              </b-icon
              >{{ comment.rating }}
            </span>
            <time-ago :datetime="comment.created" tooltip="right"> </time-ago>
          </small>
        </div>

        <b-row v-if="hasPermission(comment.author) && comment.editing">
          <b-col>
            <b-input-group>
              <b-form-textarea
                class="textarea"
                autofocus
                v-model="comment.editor"
                @keydown.enter.exact.prevent
                @keyup.enter.exact="updateComment(comment)"
                keydown.enter.shift.exact="newline"
                rows="1"
                no-resize
                :disabled="comment.requesting"
                max-rows="5"
              ></b-form-textarea>
              <template #append>
                <b-button-group size="sm" vertical v-if="comment.editing">
                  <b-button
                    :disabled="
                      comment.editor == comment.content ||
                        !properEditorLength(comment.editor) ||
                        comment.requesting
                    "
                    variant="outline-success"
                    @click="updateComment(comment)"
                  >
                    Save
                  </b-button>
                  <b-button
                    class="m-0"
                    variant="outline-danger"
                    @click="comment.editing = false"
                    :disabled="comment.requesting"
                    >Cancel</b-button
                  >
                </b-button-group>
              </template>
            </b-input-group>
            <small v-if="properEditorLength(comment.editor)">
              {{ editorLength(comment.editor) }}/{{ maxCommentLength }}
            </small>
            <small v-else style="color: red">
              {{ editorLength(comment.editor) }}/{{ maxCommentLength }}
            </small>
          </b-col>
        </b-row>
      </b-overlay>
    </div>
  </div>
</template>
<script>
import ControlButtons from "@/components/ControlButtons";
import ControlModal from "@/components/ControlModal";
import TimeAgo from "@/components/TimeAgo";
import { mapGetters } from "vuex";
import api from "@/api/api";

export default {
  name: "CommentsTab",
  mounted() {
    this.getComments();
  },
  data() {
    return {
      comments: [],
      maxCommentLength: 200,
      minCommentLength: 1,
    };
  },
  components: { ControlModal, ControlButtons, TimeAgo },
  computed: {
    filteredComments: function() {
      return this.comments.filter((c) => !c.deleted);
    },
    ...mapGetters({
      user: "auth/user",
      authenticated: "auth/authenticated",
    }),
  },
  methods: {
    hasPermission(author) {
      return this.user?.id === author?.id || this.user?.role > 0;
    },
    editorLength(editor) {
      return editor.replace(/(\r\n|\n|\r|\s)/g, "").length;
    },
    properEditorLength(editor) {
      let el = this.editorLength(editor);
      return el >= this.minCommentLength && el <= this.maxCommentLength;
    },
    reactionColor(yourReaction) {
      return yourReaction === 1 ? "green" : yourReaction === -1 ? "red" : "";
    },
    reactionIcon(yourReaction) {
      return yourReaction === -1 ? "arrow-down" : "arrow-up";
    },
    async getComments() {
      if (this.activeTab === "Comments") return;
      return await api
        .post("comments/find", {
          by: "author",
          author: Number.parseInt(this.$route.params.userID),
        })
        .then((response) => {
          this.comments = response.data.data || [];
        })
        .catch((error) => {
          console.log(error);
        })
        .then(() => (this.madeRequest = true));
    },
    async updateComment(comment) {
      if (comment.requesting) return;
      this.$set(comment, "requesting", true);
      comment.editing = false;
      return await api
        .put("comment/update", {
          id: comment.id,
          content: comment.editor,
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
