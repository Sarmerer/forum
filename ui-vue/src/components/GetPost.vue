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
          <b-form-input placeholder="Comment this post" variant="dark"></b-form-input>
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
      categories: [],
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
