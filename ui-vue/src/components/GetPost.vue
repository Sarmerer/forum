// TODO improve and split all this shit
<template>
  <div>
    <b-overlay :show="modal.deleting">
      <b-modal
        id="modal-1"
        v-model="modal.show"
        title="Delete this post?"
        header-bg-variant="dark"
        body-bg-variant="dark"
        footer-bg-variant="dark"
        body-class="position-static"
      >
        <p class="my-4">This action can <span style="color: red">NOT</span> be undone</p>
        <template v-slot:modal-footer="{ hide }">
          <b-button
            :disabled="modal.deleting"
            size="sm"
            variant="outline-secondary"
            @click="hide('forget')"
          >
            Cancel
          </b-button>
          <b-overlay
            :show="modal.deleting"
            rounded="sm"
            spinner-small
            spinner-variant="success"
            class="d-inline-block"
          >
            <b-button :disabled="modal.deleting" size="sm" variant="success" @click="deletePost()">
              Yes!
            </b-button>
          </b-overlay>
        </template>
      </b-modal>
    </b-overlay>
    <div class="columns">
      <div class="info-col">
        <div class="card">
          <h3 class="primary">AUTHOR</h3>
          <p>Author info</p>
        </div>
      </div>
      <div class="post-col">
        <div class="card">
          <h3 class="primary">{{ post.title }}</h3>
          <p style="color: white">{{ post.content }}</p>
          <div>
            <b-form-tag
              v-for="category in categories"
              disabled
              :key="category"
              :title="category"
              variant="dark"
              class="mr-1"
              >{{ category }}</b-form-tag
            >
            <div class="controls">
              <b-button-group v-if="user && (post.author_id == user.id || user.role > 0)" size="sm">
                <b-button size="sm" lg="1" class="controls-button" variant="light">
                  <img src="@/assets/svg/post/edit.svg" alt="edit" srcset="" />
                </b-button>
                <b-button
                  size="sm"
                  variant="danger"
                  lg="2"
                  @click="modal.show = !modal.show"
                  class="controls-button"
                >
                  <img src="@/assets/svg/post/delete.svg" alt="delete" srcset="" />
                </b-button>
              </b-button-group>
            </div>
          </div>
        </div>
        <div class="card">
          <b-form @submit="leaveComment">
            <b-input-group class="mt-1">
              <b-textarea
                type="text"
                placeholder="Comment this post"
                v-model="form.comment"
                rows="1"
                max-rows="3"
                size="sm"
                no-resize
              ></b-textarea>
              <b-input-group-append>
                <b-button
                  variant="outline-light"
                  type="submit"
                  :disabled="form.comment.length < 5 || form.comment.length > 200"
                  >Say</b-button
                >
              </b-input-group-append>
            </b-input-group>
            <div v-if="form.comment.length > 0">
              <small v-if="form.comment.length >= 5 && form.comment.length <= 200"
                >{{ form.comment.length }}/200</small
              >
              <small v-else style="color: red">{{ form.comment.length }}/200</small>
            </div>
          </b-form>
          <div v-for="(comment, index) in comments" :key="index">
            <div style="margin: 0.3rem; position: relative">
              <p
                v-if="index != currComment.editing"
                style="margin-bottom: 0px; display: block; margin-right: 4rem"
              >
                {{ comment.author_name }} says: {{ comment.content }}
              </p>
              <b-form-textarea
                v-else
                v-model="currComment.editingContent"
                rows="1"
                no-resize
                max-rows="10"
                style="margin-bottom: 0px; display: block; width: 85%"
              ></b-form-textarea>
              <small v-b-tooltip.hover :title="comment.created"
                ><timeago :datetime="comment.created" :auto-update="60"></timeago
              ></small>
              <b-button-group
                v-if="
                  user &&
                    (post.author_id == user.id || user.role > 0) &&
                    index != currComment.editing
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
                  :disabled="currComment.deleting"
                  @click="
                    () => {
                      currComment.editing = index;
                      currComment.editingContent = comment.content;
                    }
                  "
                >
                  <img src="@/assets/svg/post/edit.svg" alt="edit" srcset="" />
                </b-button>
                <b-button
                  variant="danger"
                  :disabled="currComment.deleting"
                  class="controls-button"
                  @click="deleteComment(comment.id, index)"
                  ><img src="@/assets/svg/post/delete.svg" alt="delete" srcset=""
                /></b-button>
              </b-button-group>
              <b-button-group
                size="sm"
                vertical
                v-if="index == currComment.editing"
                style="position: absolute; right: 0px; top: 2px"
              >
                <b-button
                  :disabled="currComment.editingContent == comment.content"
                  variant="success"
                  @click="updateComment(comment.id)"
                >
                  Save
                </b-button>
                <b-button variant="outline-danger" @click="currComment.editing = -1"
                  >Cancel</b-button
                >
              </b-button-group>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
<script>
import axios from "axios";
import { mapGetters } from "vuex";
export default {
  computed: {
    ...mapGetters({
      user: "auth/user",
    }),
    commentValidation() {
      return this.form.comment.length > 2 ? true : false;
    },
  },
  data() {
    return {
      modal: { show: false, deleting: false },
      currComment: { deleting: false, editing: -1, editingContent: "" },
      post: {},
      comments: [],
      categories: [],
      form: {
        comment: "",
      },
    };
  },
  mounted() {
    this.getPost();
  },
  methods: {
    async getPost() {
      return await axios
        .post("post/find", {
          by: "id",
          id: Number.parseInt(this.$route.params.id),
        })
        .then((response) => {
          response.data.data.forEach((res) => {
            this.post = res.post;
            if (res.replies) {
              this.comments = res.replies.sort(function(a, b) {
                return new Date(b.created) - new Date(a.created);
              });
            }
            res.categories.forEach((c) => {
              this.categories.push(c.name);
            });
          });
        })
        .catch(() => {
          this.$router.push("/");
        });
    },
    async deletePost() {
      this.modal.deleting = true;
      return await axios
        .delete("post/delete", { params: { ID: this.post.id } })
        .then(() => {
          this.modal.show = false;
          this.modal.deleting = false;
          this.$router.push("/");
        })
        .catch(() => {
          this.modal.show = false;
          this.modal.deleting = false;
          // TODO show error notification
        });
    },
    async deleteComment(actualID, IDInList) {
      this.currComment.editing = -1;
      this.currComment.deleting = true;
      return await axios
        .delete("comment/delete", { params: { ID: actualID } })
        .then(() => {
          this.comments.splice(IDInList, 1);
          this.currComment.deleting = false;
        })
        .catch((error) => {
          console.log(error);
          this.currComment.deleting = false;
          //TODO show alert
        });
    },
    async leaveComment(e) {
      e.preventDefault();
      return await axios
        .post("comment/add", { pid: this.post.id, content: this.form.comment })
        .then((response) => {
          console.log(response);
          this.appendComment();
          this.form.comment = "";
        })
        .catch((error) => {
          console.log(error);
        });
    },
    async updateComment(actualID) {
      return await axios
        .put("comment/update", {
          id: actualID,
          content: this.currComment.editingContent,
        })
        .then(() => {
          console.log(this.comments, this.currComment.editing);
          this.comments[this.currComment.editing].content = this.currComment.editingContent;
          this.currComment.editing = -1;
        })
        .catch((error) => {
          console.log(error);
        });
    },
    appendComment() {
      let comment = {
        author_id: this.user.id,
        author_name: this.user.display_name,
        content: this.form.comment,
        created: Date.now(),
        id: this.comments.length !== 0 ? this.comments[0].id + 1 : 1,
        post: this.post.id,
      };
      this.comments = [comment, ...this.comments];
    },
  },
};
</script>
<style lang="scss">
.controls {
  position: absolute;
  top: 5px;
  right: 10px;
}
.controls-button {
  background-color: transparent;
  border-color: transparent;
}
.controls-button:hover {
  background: transparent;
}
</style>
