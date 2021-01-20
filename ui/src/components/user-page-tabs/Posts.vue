<template>
  <div>
    <div
      v-for="post in filteredPosts"
      :key="post.id"
      :class="isMobile() ? 'card-m' : 'card'"
    >
      <ControlModal
        v-if="isMobile() && hasPermission(post.author)"
        v-on:delete-event="deletePost(post)"
        v-on:edit-event="edit(post)"
        :modalID="`modal-menu${post.id}`"
      />
      <div v-if="!post.editing">
        <b-row>
          <b-col>
            <router-link
              :to="'/post/' + post.id"
              tag="h4"
              class="primary text-break"
              >{{ post.title }}
            </router-link>
          </b-col>
          <b-col cols="end">
            <ControlButtons
              :class="isMobile() ? 'mr-4' : 'mr-2'"
              :hasPermission="hasPermission(post.author)"
              v-on:delete-event="deletePost(post)"
              v-on:edit-event="edit(post)"
              :disabled="false"
              :compact="isMobile()"
              :modalID="`modal-menu${post.id}`"
            />
          </b-col>
        </b-row>
        <router-link :to="'/post/' + post.id" tag="pre">
          {{ post.content }}
        </router-link>
        <b-row no-gutters>
          <b-form-tag
            v-for="category in post.categories"
            disabled
            :key="category.id"
            :title="category.name"
            variant="dark"
            class="mr-1 mb-1"
            >{{ category.name }}
          </b-form-tag>
        </b-row>
        <b-row no-gutters>
          <small>
            <span v-b-tooltip.hover title="Rating">
              <b-icon
                :icon="reactionIcon(post.your_reaction)"
                :color="reactionColor(post.your_reaction)"
              >
              </b-icon
              >{{ post.rating }}
            </span>
            <span v-b-tooltip.hover title="Comments">
              <b-icon-chat></b-icon-chat> {{ post.comments_count }}
            </span>
            <span v-b-tooltip.hover title="Participants">
              <b-icon-people></b-icon-people>
              {{ post.participants_count }}
            </span>
            <time-ago :datetime="post.created" tooltip="right"> </time-ago>
          </small>
        </b-row>
      </div>
      <b-form v-if="post.editing">
        <b-form-row>
          <b-col>
            <PostForm
              :form="post.editor"
              v-on:valid-form="$set(post, 'valid', $event)"
            />
          </b-col>
        </b-form-row>
        <b-form-row class="mt-2">
          <b-col align="end">
            <b-button-group size="sm" v-if="!post.confirmCancel">
              <b-button
                variant="outline-danger"
                @click="$set(post, 'confirmCancel', true)"
              >
                Cancel
              </b-button>
              <b-button
                variant="outline-success"
                :disabled="!post.valid"
                @click.prevent="updatePost(post)"
                class="px-3"
              >
                Save
              </b-button>
            </b-button-group>
          </b-col>
          <b-col align="end" v-if="post.confirmCancel">
            <p class="m-0">Cancel editor?</p>
            <b-button-group size="sm">
              <b-button
                variant="outline-danger"
                @click="post.confirmCancel = false"
              >
                <b-icon-x></b-icon-x> No
              </b-button>
              <b-button
                variant="outline-success"
                @click="
                  (post.editing = false),
                    (post.confirmCancel = false),
                    (post.editor = null)
                "
              >
                <b-icon-check2></b-icon-check2> Yes
              </b-button>
            </b-button-group>
          </b-col>
        </b-form-row>
      </b-form>
    </div>
  </div>
</template>
<script>
import ControlButtons from "@/components/ControlButtons";
import ControlModal from "@/components/ControlModal";
import PostForm from "@/components/forms/PostForm";
import TimeAgo from "@/components/TimeAgo";
import { mapGetters } from "vuex";
import api from "@/api/api";

export default {
  name: "PostsTab",
  data() {
    return {
      posts: [],
    };
  },
  components: {
    ControlButtons,
    ControlModal,
    PostForm,
    TimeAgo,
  },
  computed: {
    filteredPosts: function() {
      return this.posts.filter((p) => !p.deleted);
    },
    postsLength: function() {
      return this.filteredPosts.length;
    },
    ...mapGetters({
      user: "auth/user",
      authenticated: "auth/authenticated",
    }),
  },
  mounted() {
    this.getPosts();
  },
  methods: {
    hasPermission(author) {
      if (!author?.id || !this.user?.id) return false;
      return this.user.id === author.id;
    },
    reactionColor(yourReaction) {
      return yourReaction === 1 ? "green" : yourReaction === -1 ? "red" : "";
    },
    reactionIcon(yourReaction) {
      return yourReaction === -1 ? "arrow-down" : "arrow-up";
    },
    async getPosts() {
      return await api
        .post("post/find", {
          by: "author",
          author: Number.parseInt(this.$route.params.userID),
        })
        .then((response) => {
          this.posts = response.data.data || [];
        })
        .catch((error) => {
          console.log(error);
        })
        .then(() => (this.madeRequest = true));
    },
    async deletePost(post) {
      return await api
        .delete("post/delete", {
          params: { id: post.id },
        })
        .then(() => {
          this.$set(post, "deleted", true);
        });
    },
    async updatePost(post) {
      return await api
        .put("post/update", {
          id: post.id,
          title: post.editor.title,
          content: post.editor.content,
          categories: post.editor.categories,
        })
        .then((response) => {
          if (response?.data?.data) {
            if (!response?.data?.data) return;
            post.title = response.data.data.title;
            post.content = response.data.data.content;
            post.categories = response.data.data.categories;
            post.updatred = response.data.data.updated;
            document.title = post.title;
            post.editing = false;
            post.confirmCancel = false;
            post.editor = null;
          }
        })
        .catch((error) => {
          if (error.status === 403)
            this.$bvToast.toast("You need to be logged in, to update posts!", {
              title: "Oops!",
              variant: "danger",
              solid: true,
            });
        });
    },
    edit(post) {
      this.$set(post, "editing", true);
      this.$set(post, "editor", {
        title: post.title,
        content: post.content,
        categories: post.categories ? post.categories.map((c) => c.name) : [],
      });
    },
  },
};
</script>
