<template>
  <div>
    <div
      v-for="post in filteredPosts"
      :key="post.id"
      :class="{
        'card-m': isMobile(),
        'card card-hover': !isMobile(),
      }"
    >
      <ControlModal
        v-if="isMobile() && hasPermission(post.author)"
        v-on:delete-event="deletePost(post)"
        v-on:edit-event="$set(post, 'editing', true)"
        :modalID="`modal-menu${post.id}`"
      />
      <b-overlay
        :show="post.requesting"
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
        <div v-if="!post.editing">
          <b-row>
            <b-col>
              <router-link
                :to="`/post/${post.id}`"
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
                v-on:edit-event="$set(post, 'editing', true)"
                :disabled="post.requesting || false"
                :compact="isMobile()"
                :modalID="`modal-menu${post.id}`"
              />
            </b-col>
          </b-row>
          <router-link :to="`/post/${post.id}`">
            <pre v-if="!post.is_image" class="mb-1">{{ post.content }}</pre>
            <b-img-lazy
              v-if="post.is_image"
              :src="post.content"
              class="mb-2"
              center
              rounded
              fluid-grow
            >
            </b-img-lazy>
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
        <div v-if="post.editing">
          <h4 align="center">Edit post</h4>
          <PostForm
            edit
            :formData="post"
            v-on:post-update="updatePost(post, $event)"
          >
            <template slot="buttons" slot-scope="props">
              <b-row class="mt-2">
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
                      :disabled="!props.validForm"
                      type="submit"
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
                      variant="outline-success"
                      @click="post.confirmCancel = false"
                    >
                      <b-icon-x></b-icon-x> No
                    </b-button>
                    <b-button
                      variant="outline-danger"
                      @click="
                        (post.editing = false), (post.confirmCancel = false)
                      "
                    >
                      <b-icon-check2></b-icon-check2> Yes
                    </b-button>
                  </b-button-group>
                </b-col>
              </b-row>
            </template>
          </PostForm>
        </div>
      </b-overlay>
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
      return this.user?.id === author?.id || this.user?.role > 0;
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
      if (post.requesting) return;
      this.$set(post, "requesting", true);
      return await api
        .delete("post/delete", {
          params: { id: post.id },
        })
        .then(() => {
          this.$set(post, "deleted", true);
          this.$set(post, "requesting", false);
        });
    },
    updatePost(post, newData) {
      Object.assign(post, newData);
      this.$set(post, "editing", false);
    },
  },
};
</script>
