<template>
  <div class="grid">
    <div class="columns">
      <div v-if="isMobile()" class="info-col">
        <div class="card">
          <p>Comments: {{ post.comments_count }}</p>
          <p>1 participant</p>
          <p>Last comment from: {{ post.last_comment_from_name }}</p>
          <p>Last activity:</p>
        </div>
      </div>
      <div class="main-col">
        <div class="card">
          <b-row>
            <b-col cols="1">
              <Rating :callback="rate" :entity="post" type="comment" />
            </b-col>
            <b-col cols="11" style="margin-left: -35px">
              <h3 class="primary">{{ post.title }}</h3>
              <p style="color: white">{{ post.content }}</p>
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
          <p>
            {{ post.comments_count }} comment{{
              post.comments_count == 1 ? "" : "s"
            }}
          </p>
          <p>
            {{ post.total_participants }} participant{{
              post.total_participants == 1 ? "" : "s"
            }}
          </p>
          <p>Last comment from: {{ post.last_comment_from_name }}</p>
          <p>Last activity: {{ post.last_comment_date | formatDate }}</p>
        </div>
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
    postData: { type: Object, required: false },
  },
  components: {
    Rating,
    CommentsSection,
  },
  computed: {
    ...mapGetters({
      user: "auth/user",
    }),
  },
  data() {
    return {
      post: {},
      categories: [],
      postID: Number.parseInt(this.$route.params.id),
      modal: { show: false, deleting: false },
    };
  },
  created() {
    this.postData ? (this.post = this.postData) : this.getPost();
  },
  methods: {
    async getPost() {
      return await axios
        .post("post/find", {
          by: "id",
          id: this.postID,
        })
        .then((response) => {
          //TODO create error page for post
          const result = response.data.data;
          this.post = result;
        })
        .catch((error) => {
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
    },
    async rate(reaction, post) {
      if (this.requesting) return;
      let r = reaction == "up" ? 1 : -1;
      if (
        (reaction == "up" && post.your_reaction == 1) ||
        (reaction == "down" && post.your_reaction == -1)
      ) {
        r = 0;
      }
      await axios
        .post("post/rate", { id: this.postID, reaction: r })
        .then((response) => {
          post.your_reaction = response.data.data.your_reaction;
          post.rating = response.data.data.rating;
        })
        .catch((error) => {
          console.log(error);
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
