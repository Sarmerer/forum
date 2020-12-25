<template>
  <div>
    <div
      v-for="(comment, index) in comments"
      :key="comment.id"
      class="ml-3 mt-2 position-relative"
    >
      <div
        :class="
          isMobile() || comment.deleted ? 'collapse-line-m' : 'collapse-line'
        "
        v-if="comment.children && !comment.collapsed"
        @click="$set(comment, 'collapsed', true)"
      ></div>
      <div v-if="comment.deleted && comment.children">
        <span class="text-white-50">deleted</span>
      </div>
      <div v-if="!comment.deleted">
        <ControlModal
          v-on:delete-event="deleteComment(comment.id, index)"
          v-on:edit-event="
            (editor.editing = comment.id), (editor.content = comment.content)
          "
          :modalID="'modal-menu' + index"
        />
        <div v-if="editor.editing !== comment.id" class="mb-2">
          <b-row :class="isMobile() ? 'm-0' : ''">
            <b-col v-if="!isMobile()" cols="start" class="ml-n2 mr-1">
              <Rating :entity="comment" size="md" endpoint="comment" />
            </b-col>
            <b-col>
              <b-row align-h="between">
                <b-col cols="start">
                  <router-link
                    :to="`/user/${comment.author.id}`"
                    class="text-white-50 text-underscore"
                  >
                    <small>
                      <user-popover
                        noAvatar
                        :userData="comment.author"
                        :popoverID="'c' + index"
                      >
                      </user-popover>
                    </small>
                  </router-link>
                  <small v-if="!isMobile()" class="text-white-50">
                    {{ comment.rating }}
                    {{
                      comment.rating % 10 === 1 || comment.rating % 10 === -1
                        ? "point"
                        : "points"
                    }}
                  </small>
                  <small class="text-white-50">
                    <time-ago :datetime="comment.created" tooltip="right">
                    </time-ago>
                    <small
                      v-b-tooltip.hover="formatDate(comment.edited)"
                      v-if="comment.edited != comment.created"
                      class="text-muted"
                      >(edited)
                    </small>
                  </small>
                </b-col>
                <b-col v-if="isMobile()" cols="end">
                  <ControlButtons
                    :hasPermission="hasPermission(comment)"
                    :disabled="requesting"
                    compact
                    :modalID="'modal-menu' + index"
                  />
                </b-col>
              </b-row>
              <b-row>
                <pre class="text-break m-0">
                {{ comment.content }}
              </pre
                >
              </b-row>
              <b-row align-h="between">
                <b-col v-if="isMobile()" cols="start">
                  <Rating :entity="comment" size="sm" endpoint="comment" />
                </b-col>
                <b-col class="pl-1">
                  <!-- TODO add redirect to auth page if not authorized -->
                  <small v-if="!comment.replying">
                    <a
                      class="secondary"
                      @click="
                        $set(comment, 'reply', ''),
                          $set(comment, 'replying', true)
                      "
                    >
                      reply
                    </a>
                  </small>
                </b-col>
              </b-row>
            </b-col>
            <b-col v-if="!isMobile()" cols="end" class="mr-2">
              <ControlButtons
                :hasPermission="hasPermission(comment)"
                v-on:delete-event="deleteComment(comment.id, index)"
                v-on:edit-event="
                  (editor.editing = comment.id),
                    (editor.content = comment.content)
                "
                :disabled="requesting"
                :compact="isMobile()"
                :modalID="'modal-menu' + index"
              />
            </b-col>
          </b-row>
        </div>
        <b-row v-if="hasPermission(comment) && editor.editing === comment.id">
          <b-col class="ml-n2">
            <b-input-group>
              <b-form-textarea
                class="textarea"
                autofocus
                v-model="editor.content"
                @keydown.enter.exact.prevent
                @keyup.enter.exact="
                  updateComment(comment.id, index, editor.content)
                "
                keydown.enter.shift.exact="newline"
                rows="1"
                no-resize
                :disabled="requesting"
                max-rows="5"
              ></b-form-textarea>
              <template #append>
                <b-button-group
                  size="sm"
                  vertical
                  v-if="editor.editing === comment.id"
                >
                  <b-button
                    :disabled="
                      editor.content == comment.content ||
                        !properEditorLength ||
                        requesting
                    "
                    variant="outline-success"
                    @click="updateComment(comment.id, index, editor.content)"
                  >
                    Save
                  </b-button>
                  <b-button
                    class="m-0"
                    variant="outline-danger"
                    @click="editor.editing = -1"
                    :disabled="requesting"
                    >Cancel</b-button
                  >
                </b-button-group>
              </template>
            </b-input-group>
            <small v-if="properEditorLength"
              >{{ editorLength }}/{{ maxCommentLength }}</small
            >
            <small v-else style="color: red"
              >{{ editorLength }}/{{ maxCommentLength }}</small
            >
          </b-col>
        </b-row>
        <div v-if="comment.replying">
          <b-input-group>
            <b-form-textarea
              class="textarea"
              autofocus
              v-model="comment.reply"
              @keydown.enter.exact.prevent
              @keyup.enter.exact="
                reply(comment.post_id, comment.id, index, comment.reply)
              "
              keydown.enter.shift.exact="newline"
              rows="1"
              no-resize
              :disabled="requesting"
              max-rows="5"
            ></b-form-textarea>
            <template #append>
              <b-button-group size="sm" vertical>
                <b-button
                  :disabled="!properReplyLength(comment.reply)"
                  variant="outline-success"
                  @click="
                    reply(comment.post_id, comment.id, index, comment.reply)
                  "
                >
                  Save
                </b-button>
                <b-button
                  class="m-0"
                  variant="outline-danger"
                  @click="comment.replying = false"
                  :disabled="requesting"
                  >Cancel</b-button
                >
              </b-button-group>
            </template>
          </b-input-group>
          <small v-if="properReplyLength(comment.reply)"
            >{{ comment.reply.length }}/{{ maxCommentLength }}</small
          >
          <small v-else style="color: red"
            >{{ comment.reply.length }}/{{ maxCommentLength }}</small
          >
        </div>
      </div>
      <CommentGroup
        v-if="comment.children"
        v-show="!comment.collapsed"
        :editor="editor"
        :comments="comment.children"
      />
      <div
        v-if="comment.collapsed"
        @click="$set(comment, 'collapsed', false)"
        class="ml-2"
      >
        <small>
          <b-icon-plus-circle class="secondary"></b-icon-plus-circle>
          <span class="text-white-50">
            {{ comment.children_length }} more
            {{ comment.children_length % 10 === 1 ? "comment" : "comments" }}
          </span>
        </small>
      </div>
    </div>
  </div>
</template>
<script>
import ControlButtons from "@/components/ControlButtons";
import ControlModal from "@/components/ControlModal";
import UserPopover from "@/components/UserPopover";
import TimeAgo from "@/components/TimeAgo";
import Rating from "@/components/Rating";
import moment from "moment-shortformat";
import { mapGetters } from "vuex";
import api from "@/router/api";

export default {
  name: "CommentGroup",
  props: {
    comments: { type: Array, required: true },
    editor: Object,
  },
  components: {
    ControlButtons,
    ControlModal,
    UserPopover,
    TimeAgo,
    Rating,
  },
  data() {
    return {
      maxCommentLength: 200,
      minCommentLength: 5,
      requesting: false,
    };
  },
  computed: {
    ...mapGetters({
      user: "auth/user",
      authenticated: "auth/authenticated",
    }),
    editorLength() {
      return this.editor.content.replace(/(\r\n|\n|\r|\s)/g, "").length;
    },
    properEditorLength() {
      let el = this.editorLength;
      return el >= this.minCommentLength && el <= this.maxCommentLength;
    },
  },
  methods: {
    hasPermission(comment) {
      return comment?.author?.id == this.user?.id || this.user?.role > 0;
    },
    formatDate(date) {
      return moment.unix(date).format("ddd, d MMM y, hh:mm");
    },
    replyLength(reply) {
      return reply.replace(/(\r\n|\n|\r|\s)/g, "").length;
    },
    properReplyLength(reply) {
      let rl = this.replyLength(reply);
      return rl >= this.minCommentLength && rl <= this.maxCommentLength;
    },
    async reply(postID, parentID, parentIndex, content) {
      if (!parentID || !content || !this.properReplyLength(content)) return;
      return await api
        .post("comment/add", {
          id: postID,
          parent: parentID,
          content: content,
        })
        .then((response) => {
          if (response?.data?.data) {
            if (this.comments[parentIndex].children) {
              this.comments[parentIndex].children.push(response.data.data);
              this.comments[parentIndex].children_length++;
            } else {
              this.$set(this.comments[parentIndex], "children", [
                response.data.data,
              ]);
              this.$set(this.comments[parentIndex], "children_length", 1);
            }
          }
          this.comments[parentIndex].replying = false;
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
        });
    },
    async updateComment(commentID, commentIndex, newContent) {
      if (this.requesting || !this.properEditorLength) return;
      this.requesting = true;
      return await api
        .put("comment/update", {
          id: commentID,
          content: newContent,
        })
        .then((response) => {
          if (response?.data?.data) {
            this.comments[commentIndex].content = response.data.data.content;
            this.comments[commentIndex].edited = response.data.data.edited;
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
          this.editor.editing = -1;
          this.requesting = false;
        });
    },
    async deleteComment(commentID, commentIndex) {
      if (this.requesting) return;
      this.requesting = true;
      return await api
        .delete("comment/delete", {
          params: { id: commentID },
        })
        .then(() => {
          if (this.comments[commentIndex].children) {
            this.$set(this.comments[commentIndex], "deleted", true);
          } else {
            this.comments.splice(commentIndex, 1);
          }
          this.requesting = false;
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
.collapse-line {
  height: calc(100% - 50px);
  top: 50px;
  left: -16px;
  position: absolute;
  width: 15px;
  background-image: linear-gradient(
    rgba(255, 255, 255, 0.1),
    rgba(255, 255, 255, 0.1)
  );
  background-size: 4px 100%;
  background-repeat: no-repeat;
  background-position: center center;
}

.collapse-line-m {
  height: calc(100% - 5px);
  margin-left: -16px;
  top: 5px;
  width: 15px;
  position: absolute;
  background-image: linear-gradient(
    rgba(255, 255, 255, 0.1),
    rgba(255, 255, 255, 0.1)
  );
  background-size: 4px 100%;
  background-repeat: no-repeat;
  background-position: center center;
}

.collapse-line:hover,
.collapse-line-m:hover {
  background-image: linear-gradient(grey, grey);
}
</style>
