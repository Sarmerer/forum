<template>
  <div class="columns">
    <div class="main-col">
      <div class="card">
        <b-row>
          <b-col cols="1">
            <Rating
              v-on:update="
                args => {
                  post.rating = args.new_rating;
                  post.your_reaction = args.new_your_reaction;
                }
              "
              :postID="post.id"
              :rating="post.rating"
              :yourReaction="post.your_reaction"
          /></b-col>
          <b-col cols="11" style="margin-left: -35px;">
            <h3 class="primary">{{ post.title }}</h3>
            <p>{{ post.content }}</p>
            <div>
              <b-form-tag
                v-for="category in categories"
                disabled
                :key="category"
                :title="category"
                variant="dark"
                class="mr-1 mb-1"
              >
                {{ category }}
              </b-form-tag>
              <div class="controls">
                <b-overlay
                  :show="modal.deleting"
                  rounded
                  opacity="0.5"
                  spinner-small
                  variant="light"
                  spinner-variant="primary"
                  class="d-inline-block"
                >
                  <transition name="fade">
                    <b-button-group
                      v-if="
                        (user
                          ? post.author_id == user.id || user.role > 0
                          : false) && !modal.show
                      "
                      size="sm"
                    >
                      <b-button
                        size="sm"
                        lg="1"
                        class="controls-button"
                        variant="light"
                        title="Edit"
                      >
                        <img
                          src="@/assets/svg/post/edit.svg"
                          alt="edit"
                          srcset=""
                        />
                      </b-button>
                      <b-button
                        size="sm"
                        variant="outline-danger"
                        lg="2"
                        @click="modal.show = !modal.show"
                        class="controls-button"
                        title="Delete"
                      >
                        <img
                          src="@/assets/svg/post/delete.svg"
                          alt="delete"
                          srcset=""
                        />
                      </b-button>
                    </b-button-group>

                    <b-button-group
                      v-if="
                        (user
                          ? post.author_id == user.id || user.role > 0
                          : false) && modal.show
                      "
                      size="sm"
                    >
                      <b-button
                        size="sm"
                        variant="outline-success"
                        lg="2"
                        class="confirm"
                        @click="deletePost()"
                        title="Confirm"
                      >
                        <img
                          src="@/assets/svg/post/confirm.svg"
                          alt="delete"
                          srcset=""
                        />
                      </b-button>
                      <b-button
                        size="sm"
                        lg="1"
                        variant="outline-danger"
                        @click="modal.show = !modal.show"
                        class="confirm"
                        title="Dismiss"
                      >
                        <img
                          src="@/assets/svg/post/dismiss.svg"
                          alt="edit"
                          srcset=""
                        />
                      </b-button>
                    </b-button-group>
                  </transition>
                </b-overlay>
              </div></div
          ></b-col>
        </b-row>
      </div>
      <CommentsSection :postID="postID" />
    </div>
    <div v-if="!isMobile()" class="info-col">
      <div class="card">
        <p>1 reply</p>
        <p>1 participant</p>
        <p>Last reply from:</p>
        <p>Last activity:</p>
      </div>
    </div>
  </div>
</template>
<script>
import axios from "axios";
import { mapGetters } from "vuex";
import Rating from "@/components/Rating";
import CommentsSection from "@/components/CommentsSection";

export default {
  props: {
    postID: { type: Number, required: true }
  },
  components: {
    Rating,
    CommentsSection
  },
  computed: {
    ...mapGetters({
      user: "auth/user"
    })
  },
  data() {
    return {
      post: {},
      categories: [],
      modal: { show: false, deleting: false }
    };
  },
  created() {
    this.getPost();
  },
  methods: {
    async getPost() {
      return await axios
        .post("post/find", {
          by: "id",
          id: this.postID
        })
        .then(response => {
          //TODO create error page for post
          //? because response.data.data is an array of objects
          const result = response.data.data;
          this.post = result;
          if (result.replies) {
            this.comments = result.replies.sort(function(a, b) {
              return new Date(b.created) - new Date(a.created);
            });
          }
          if (result.categories) {
            result.categories.forEach(c => {
              this.categories.push(c.name);
            });
          }
        })
        .catch(error => {
          console.log(error);
          this.$router.push("/");
        });
    },
    async deletePost() {
      this.modal.deleting = true;
      return await axios
        .delete("post/delete", { params: { id: this.post.id } })
        .then(() => {
          this.$router.push("/");
        })
        .catch(() => {
          this.modal.show = false;
          // TODO show error notification
        })
        .finally(() => {
          this.modal.deleting = false;
          this.modal.show = false;
        });
    }
  }
};
</script>
<style lang="scss" scoped>
.controls {
  position: absolute;
  top: 5px;
  right: -15px;
}
.confirm {
  box-shadow: none;
  -moz-box-shadow: none;
  -webkit-box-shadow: none;
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

.fade-enter-active {
  transition: opacity 0.5s;
}
.fade-enter, .fade-leave-to /* .fade-leave-active below version 2.1.8 */ {
  opacity: 0;
}
</style>
