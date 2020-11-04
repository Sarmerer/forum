<template>
  <div>
    <div v-if="$route.params.id === 'new'" class="grid">
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
            <!-- <b-form-group label-for="amount">
          <b-form-input
            id="amount"
            v-model="form.amount"
            autocomplete="off"
            required
            placeholder="Enter title"
          ></b-form-input>
        </b-form-group> -->
            <b-form-group id="input-group-2" label-for="input-2" fluid>
              <b-form-textarea
                id="textarea-auto-height"
                v-model="form.content"
                placeholder="Enter content"
                reqired
                :state="
                  form.content.length >= 10 && form.content.length <= 2000
                "
                rows="5"
                max-rows="10"
              ></b-form-textarea>
            </b-form-group>

            <!-- <label for="tags-basic">Type a new tag and press enter</label> -->
            <b-form-tags
              input-id="tags-basic"
              remove-on-delete
              v-model="form.categories"
              tag-variant="dark"
            ></b-form-tags>
            <b-button type="submit" class="mt-3">Submit</b-button>
          </b-form>
          <!-- <b-card class="mt-3" header="Form Data Result">
        <pre class="m-0" style="color: white">{{ form }}</pre>
      </b-card> -->
        </div>
        <div class="info-col">
          <div class="card" id="user-info">
            <img :src="user.avatar" alt="avatar" />
            <h3 class="primary">
              {{ user.display_name }}
            </h3>
            <h5>
              <b-badge v-if="user.role == 2" class="background-variant"
                >Admin</b-badge
              >
            </h5>
          </div>
        </div>
      </div>
    </div>
    <div v-else class="grid">
      <PostSection :postID="postID" />
    </div>
  </div>
</template>

<script>
import axios from "axios";
import PostSection from "@/components/PostSection";
import { mapGetters } from "vuex";
export default {
  components: {
    PostSection
  },
  computed: {
    ...mapGetters({
      user: "auth/user"
    })
  },
  data() {
    return {
      postID: Number.parseInt(this.$route.params.id),
      form: {
        title: "",
        amount: 1,
        content: "",
        categories: []
      }
    };
  },
  methods: {
    onSubmit(e) {
      e.preventDefault();
      for (let i = 0; i < this.form.amount; i++) {
        axios
          .post("post/create", {
            title: this.form.title,
            content: this.form.content,
            categories: this.form.categories
          })
          .then(response => {
            console.log(response.data);
            this.resetForm();
            this.$router.push("/post/" + response.data.data);
          })
          .catch(error => {
            alert(error.response.data.code + " " + error.response.data.message);
          });
      }
    },
    resetForm() {
      this.form.title = "";
      this.form.content = "";
      this.form.categories = [];
    }
  },
  watch: {
    "$route.params": {
      handler(newID) {
        const { id } = newID;
        this.postID = Number.parseInt(id);
      },
      immediate: true
    }
  }
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

#new-post {
  margin: 20px;
}

#user-info {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  text-align: center;
}

#user-info img {
  width: 100px;
  height: 100px;
  border-radius: 200px;
}
</style>
