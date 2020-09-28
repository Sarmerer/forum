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
            <!-- Emulate built in modal footer ok and cancel button actions -->
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
          <h3 class="primary">{{ post.Title }}</h3>
          <p style="color: white">{{ post.Content }}</p>
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
              <button class="contols-button">
                <img src="@/assets/svg/post/edit.svg" alt="edit" srcset="" />
              </button>
              <button @click="modal.show = !modal.show" class="controls-button">
                <img src="@/assets/svg/post/delete.svg" alt="delete" srcset="" />
              </button>
            </div>
          </div>
        </div>
      </div>
      <div class="post-col">
        <div class="card">
          <b-form @submit="onSubmit" inline>
            <b-input
              id="inline-form-input-name"
              class="mb-2 mr-sm-2 mb-sm-0"
              placeholder="Comment"
              v-model="form.comment"
              style="width: 85%"
            ></b-input>

            <b-button type="submit" variant="dark">submit</b-button>
          </b-form>
          <div v-for="comment in comments" :key="comment.ID">
            <h5 style="display: inline">{{ comment.Content }}</h5>
            <sup style="display: inline"> {{ comment.Created }}</sup>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
<script>
import axios from "axios";
export default {
  data() {
    return {
      modal: { show: false, deleting: false },
      post: {},
      comments: {},
      categories: [],
      form: {
        comment: "",
      },
    };
  },
  mounted() {
    this.getPost().then(() => {
      console.log(this.comments);
    });
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
            this.comments = res.replies;
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
        .delete("post/delete1", { params: { ID: this.post.ID } })
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
    async onSubmit(e) {
      console.log(this.comments);
      e.preventDefault();
      return await axios
        .post("comment/add", { pid: this.post.ID, content: this.form.comment })
        .then((response) => {
          console.log(response);
          this.form.comment = "";
        })
        .catch((error) => {
          console.log(error);
        });
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
  margin-left: 20px;
}
</style>
