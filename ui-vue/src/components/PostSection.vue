<template>
  <div>
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
        >
          {{ category }}
        </b-form-tag>
        <div class="controls">
          <b-button-group
            v-if="user ? post.author_id == user.id || user.role > 0 : false"
            size="sm"
          >
            <b-button size="sm" lg="1" class="controls-button" variant="light">
              <img src="@/assets/svg/post/edit.svg" alt="edit" srcset="" />
            </b-button>
            <b-button
              size="sm"
              variant="outline-danger"
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
  </div>
</template>
<script>
import axios from "axios";
import { mapGetters } from "vuex";

export default {
  props: {
    postID: { type: Number },
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
      modal: { show: false, deleting: false },
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
          id: this.postID,
        })
        .then((response) => {
          //? because response.data.data is an array of objects
          const result = response.data.data[0];
          this.post = result.post;
          if (result.replies) {
            this.comments = result.replies.sort(function(a, b) {
              return new Date(b.created) - new Date(a.created);
            });
          }
          result.categories.forEach((c) => {
            //TODO create error page for post
            this.categories.push(c.name);
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
  },
};
</script>
<style lang="scss" scoped>
.controls {
  position: absolute;
  top: 5px;
  right: 10px;
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
</style>
