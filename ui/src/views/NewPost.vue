<template>
  <div class="grid">
    <div class="columns">
      <div class="main-col" id="new-post">
        <b-form @submit="onSubmit">
          <b-form-group label-for="title">
            <b-form-input
              id="title"
              v-model="form.title"
              autocomplete="off"
              required
              placeholder="Enter title"
            ></b-form-input>
          </b-form-group>
          <b-form-group id="input-group-2" label-for="input-2" fluid>
            <b-form-textarea
              id="textarea-auto-height"
              v-model="form.content"
              placeholder="Enter content"
              reqired
              :state="form.content.length >= 10 && form.content.length <= 2000"
              rows="5"
              max-rows="10"
            ></b-form-textarea>
          </b-form-group>
          <b-form-tags
            input-id="tags-basic"
            remove-on-delete
            v-model="form.categories"
            tag-variant="dark"
          ></b-form-tags>
          <b-button type="submit" class="mt-3">Submit</b-button>
        </b-form>
      </div>
      <div class="info-col">
        <div class="card user-card">
          <img :src="user.avatar" alt="avatar" />
          <h3 class="primary">
            {{ user.display_name }}
            <b-badge v-if="user.role == 2" class="background-variant"
              >Admin</b-badge
            >
          </h3>
        </div>
      </div>
    </div>
  </div>
</template>
<script>
import axios from "axios";
axios.defaults.withCredentials = true;
import { mapGetters } from "vuex";

export default {
  data() {
    return {
      form: {
        title: "",
        amount: 1,
        content: "",
        categories: [],
      },
    };
  },
  methods: {
    onSubmit(e) {
      e.preventDefault();
      axios
        .post("post/create", {
          title: this.form.title,
          content: this.form.content,
          categories: this.form.categories,
        })
        .then((response) => {
          this.resetForm();
          this.$router.push({
            name: "Post",
            params: { id: response.data.data.id, postData: response.data.data },
          });
        })
        .catch((error) => {
          alert(error.response.data.code + " " + error.response.data.message);
        });
    },
    resetForm() {
      this.form.title = "";
      this.form.content = "";
      this.form.categories = [];
    },
  },
  watch: {
    "$route.params": {
      handler(newID) {
        const { id } = newID;
        this.postID = Number.parseInt(id);
      },
      immediate: true,
    },
  },
  computed: {
    ...mapGetters({
      user: "auth/user",
    }),
  },
};
</script>
<style lang="scss">
form .btn {
  background-color: #278ea5;
  border: none;
  // display: block;
  // width: 100%;
}
form .btn:hover {
  background-color: #278ea5;
  opacity: 0.8;
}
.user-card img {
  display: block;
  margin-left: auto;
  margin-right: auto;
  margin-bottom: 15px;
  width: 150px;
  height: 150px;
  border-radius: 150px;
}

#new-post {
  margin: 20px;
}
</style>
