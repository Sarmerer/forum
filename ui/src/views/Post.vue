<template>
  <b-skeleton-wrapper :loading="loading">
    <template #loading>
      <PostSkeleton />
    </template>
    <div>
      <div
        v-if="distanceToWindowTop > 200 && !isMobile()"
        @click="scrollToTop"
        class="back-to-top h-100 text-center"
      >
        <b-icon-chevron-up class="mt-4"></b-icon-chevron-up>
      </div>
      <div class="grid">
        <NotFound v-if="notFound" />
        <div v-else>
          <ControlModal
            v-on:edit-event="$set(post, 'editing', true)"
            v-on:delete-event="deletePost()"
            modalID="modal-menu"
          />
          <div class="columns">
            <div v-if="isMobile()" class="info-col">
              <UserCard v-if="post.author" link :userData="post.author" />
            </div>
            <div class="main-col">
              <div :class="`${isMobile() ? 'card-m' : 'card'}`">
                <b-row v-if="!post.editing">
                  <b-col cols="start">
                    <Rating
                      :entity="post"
                      class="ml-n4"
                      size="lg"
                      endpoint="post"
                    />
                  </b-col>
                  <b-col>
                    <b-row>
                      <b-col>
                        <h3 class="primary text-break">{{ post.title }}</h3>
                      </b-col>
                      <b-col cols="end">
                        <ControlButtons
                          :class="isMobile() ? 'mr-4' : 'mr-2'"
                          :hasPermission="hasPermission"
                          v-on:delete-event="deletePost()"
                          v-on:edit-event="$set(post, 'editing', true)"
                          :disabled="requesting"
                          :compact="isMobile()"
                          modalID="modal-menu"
                        />
                      </b-col>
                    </b-row>
                    <pre v-if="!post.is_image" class="mb-1">
                      {{ post.content }}
                    </pre>
                    <b-img-lazy
                      v-if="post.is_image"
                      :src="post.content"
                      class="mb-2"
                      center
                      rounded
                      fluid-grow
                    >
                    </b-img-lazy>
                    <b-form-tag
                      v-for="category in post.categories"
                      disabled
                      :key="category.id"
                      :title="category.name"
                      variant="dark"
                      class="mr-1 mb-1"
                      >{{ category.name }}
                    </b-form-tag>
                    <b-row>
                      <b-col>
                        <small>
                          <span v-b-tooltip.hover title="Comments">
                            <b-icon-chat></b-icon-chat>
                            {{ post.comments_count }}
                          </span>
                          <span v-b-tooltip.hover title="Participants">
                            <b-icon-people></b-icon-people>
                            {{ post.participants_count }}
                          </span>
                        </small>
                      </b-col>
                      <b-col v-if="isMobile()" cols="end" class="mr-3">
                        <small>
                          <Rating size="sm" :entity="post" endpoint="post" />
                        </small>
                      </b-col>
                    </b-row>
                  </b-col>
                </b-row>
                <div v-if="post.editing">
                  <h4 align="center">Edit post</h4>
                  <PostForm :formData="post" edit v-on:post-update="updatePost">
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
                                (post.editing = false),
                                  (post.confirmCancel = false)
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
              </div>
              <Comments v-if="!noComments" :postID="postID" />
            </div>

            <UserCard
              v-if="!isMobile() && post.author"
              class="info-col"
              link
              :userData="post.author"
            />
          </div>
        </div>
      </div>
    </div>
  </b-skeleton-wrapper>
</template>
<script>
import PostSkeleton from "@/components/skeletons/PostSkeleton";
import ControlButtons from "@/components/ControlButtons";
import ControlModal from "@/components/ControlModal";
import PostForm from "@/components/forms/PostForm";
import NotFound from "@/components/NotFound";
import UserCard from "@/components/UserCard";
import Rating from "@/components/Rating";
import Comments from "@/views/Comments";
import { mapGetters } from "vuex";
import api from "@/api/api";

export default {
  props: {
    postData: Object,
    noComments: Boolean,
  },
  components: {
    ControlButtons,
    PostSkeleton,
    ControlModal,
    Comments,
    NotFound,
    PostForm,
    UserCard,
    Rating,
  },
  computed: {
    ...mapGetters({
      user: "auth/user",
    }),
    hasPermission() {
      return this.post?.author?.id == this.user?.id || this.user?.role > 0;
    },
  },
  data() {
    return {
      postID: Number.parseInt(this.$route.params.id),
      post: {},
      categories: [],
      requesting: false,
      loading: true,
      notFound: false,
      distanceToWindowTop: 0,
    };
  },
  created() {
    if (this.postData) {
      document.title = this.postData.title;
      this.post = this.postData;
      this.loading = false;
    } else {
      this.getPost().then(() => {
        setTimeout(() => {
          this.loading = false;
        }, 500);
      });
    }
  },
  mounted() {
    window.addEventListener("scroll", this.onScroll);
  },
  beforeDestroy() {
    window.removeEventListener("scroll", this.onScroll);
  },

  methods: {
    onScroll(e) {
      this.distanceToWindowTop = e.target.documentElement.scrollTop;
    },
    scrollToTop() {
      window.scrollTo(0, 0);
    },
    async getPost() {
      return await api
        .post("post/find", {
          by: "id",
          id: this.postID,
        })
        .then((response) => {
          let result = response?.data?.data;
          document.title = result.title;
          this.post = result;
        })
        .catch((error) => {
          if (error.status === 404) this.notFound = true;
        });
    },
    async deletePost() {
      this.requesting = true;
      return await api
        .delete("post/delete", {
          params: { id: this.post.id },
        })
        .then(() => {
          this.$router.push("/");
        })
        .finally(() => {
          this.requesting = false;
        });
    },
    updatePost(post) {
      this.post = post;
      document.title = this.post.title;
      this.requesting = false;
    },
  },
};
</script>
<style lang="scss">
.back-to-top {
  position: fixed;
  width: 100px;
  top: 0;
}

.back-to-top:hover {
  background-color: #2a2a2a;
}
</style>
