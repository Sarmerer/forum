<template>
  <div>
    <div
      v-for="comment in comments"
      :key="comment.id"
      class="mt-2 position-relative"
      :class="{ 'ml-2': isMobile(), 'ml-3': !isMobile() }"
    >
      <div
        :class="
          isMobile() || comment.deleted ? 'collapse-line-m' : 'collapse-line'
        "
        v-if="comment.children && !comment.collapsed"
        @click="$set(comment, 'collapsed', true)"
      ></div>
      <div v-if="comment.deleted" v-show="!comment.collapsed">
        <span class="text-white-50">deleted</span>
      </div>
      <div v-if="!comment.deleted" v-show="!comment.collapsed">
        <ControlModal
          v-if="isMobile() && hasPermission(comment.author)"
          v-on:delete-event="deleteComment(comment)"
          v-on:edit-event="
            $set(comment, 'editing', true),
              $set(comment, 'editor', comment.content)
          "
          :modalID="'modal-menu' + comment.id"
        />
        <div v-if="!comment.editing" class="mb-2">
          <b-row :class="{ 'm-0': isMobile() }">
            <b-col v-if="!isMobile()" cols="start" class="ml-n2 mr-1">
              <Rating :entity="comment" size="md" endpoint="comment" />
            </b-col>
            <b-col>
              <b-row>
                <b-col cols="start">
                  <small>
                    <user-popover
                      noAvatar
                      popoverDirection="right"
                      :userData="comment.author"
                      :popoverID="'c' + comment.id"
                    >
                    </user-popover>
                  </small>

                  <small
                    v-if="!isMobile() && comment.rating != 0"
                    class="text-white-50"
                  >
                    {{ comment.rating }}
                    {{
                      comment.rating % 10 === 1 || comment.rating % 10 === -1
                        ? "point"
                        : "points"
                    }}
                  </small>
                  <small class="text-white-50">
                    •
                    <time-ago :datetime="comment.created" tooltip="right">
                    </time-ago>
                    <span
                      v-b-tooltip.hover="formatDate(comment.edited)"
                      v-if="comment.edited != comment.created"
                      class="text-white-50"
                    >
                      • edited
                    </span>
                  </small>
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
                <b-col class="pl-0">
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
            <b-col cols="end">
              <ControlButtons
                :hasPermission="hasPermission(comment.author)"
                v-on:delete-event="deleteComment(comment)"
                v-on:edit-event="
                  $set(comment, 'editing', true),
                    $set(comment, 'editor', comment.content)
                "
                :disabled="requesting"
                :compact="isMobile()"
                :modalID="'modal-menu' + comment.id"
              />
            </b-col>
          </b-row>
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
                :disabled="requesting"
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
                    :disabled="requesting"
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

        <div v-if="comment.replying">
          <b-input-group>
            <b-form-textarea
              class="textarea"
              autofocus
              v-model="comment.reply"
              @keydown.enter.exact.prevent
              @keyup.enter.exact="reply(comment, comment.reply)"
              keydown.enter.shift.exact="newline"
              @keyup.escape.exact="comment.replying = false"
              rows="1"
              no-resize
              :disabled="requesting"
              max-rows="5"
            ></b-form-textarea>
            <template #append>
              <b-button-group size="sm" vertical>
                <b-button
                  :disabled="!properReplyLength(comment.reply)"
                  variant="outline-light"
                  @click="reply(comment, comment.reply)"
                >
                  Say
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
        :comments="comment.children"
      />
      <div
        v-if="comment.collapsed"
        @click="$set(comment, 'collapsed', false)"
        class="ml-n2"
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

import eventBus from "@/event-bus";

export default {
  name: "CommentGroup",
  props: {
    comments: { type: Array, required: true },
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
      minCommentLength: 1,
      requesting: false,
    };
  },
  computed: {
    ...mapGetters({
      user: "auth/user",
      authenticated: "auth/authenticated",
    }),
  },
  methods: {
    editorLength(editor) {
      return editor.replace(/(\r\n|\n|\r|\s)/g, "").length;
    },
    properEditorLength(editor) {
      let el = this.editorLength(editor);
      return el >= this.minCommentLength && el <= this.maxCommentLength;
    },
    hasPermission(author) {
      return author?.id === this.user?.id || this.user?.role > 0;
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
    reply(parent, content) {
      if (!this.properReplyLength(content)) return;
      eventBus.$emit("reply-event", [parent, content]);
    },
    updateComment(comment) {
      if (!this.properEditorLength(comment.editor)) return;
      eventBus.$emit("update-event", [comment, comment.editor]);
    },
    deleteComment(comment) {
      eventBus.$emit("delete-event", comment);
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
  background-size: 2px 100%;
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
  background-size: 2px 100%;
  background-repeat: no-repeat;
  background-position: center center;
}

.collapse-line:hover,
.collapse-line-m:hover {
  background-size: 3px 100%;
  background-image: linear-gradient(grey, grey);
}
</style>
