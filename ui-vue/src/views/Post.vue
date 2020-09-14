<template>
  <div class="wrapper">
    <div v-if="$route.params.id === 'new'">
      <b-form @submit="onSubmit">
        <b-form-group id="input-group-2" label="Title:" label-for="input-2">
          <b-form-input
            id="input-2"
            v-model="form.title"
            autocomplete="off"
            required
            placeholder="Enter title"
          ></b-form-input>
        </b-form-group>
        <b-form-group id="input-group-2" label="Content:" label-for="input-2" fluid>
          <b-form-textarea
            id="textarea-auto-height"
            v-model="form.content"
            placeholder="Enter content"
            :state="form.content.length >= 10 && form.content.length <= 2000"
            rows="3"
            max-rows="10"
          ></b-form-textarea>
        </b-form-group>

        <label for="tags-basic">Type a new tag and press enter</label>
        <b-form-tags input-id="tags-basic" v-model="form.categories"></b-form-tags>
        <br />
        <b-button type="submit" variant="primary">Submit</b-button>
      </b-form>
      <b-card class="mt-3" header="Form Data Result">
        <pre class="m-0" style="color: white">{{ form }}</pre>
      </b-card>
    </div>
    <div v-else><h1>get post</h1></div>
  </div>
</template>

<script>
import axios from "axios";
export default {
  data() {
    return {
      form: {
        title: "",
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
          Title: this.form.title,
          Content: this.form.content,
          Categories: this.form.categories,
        })
        .then((response) => {
          console.log(response.data);
          this.resetForm();
          this.$router.push("/post/" + response.data.data);
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
};
</script>
<style lang="scss">
.wrapper {
  text-align: center;
  width: 800px;
  margin: 0 25%;
}
</style>
