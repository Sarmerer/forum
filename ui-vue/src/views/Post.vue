<template>
  <div>
    <div v-if="$route.params.id === 'new'" class="wrapper">
      <b-form @submit="onSubmit">
        <b-form-group label="Title:" label-for="title">
          <b-form-input
            id="title"
            v-model="form.title"
            autocomplete="off"
            required
            placeholder="Enter title"
          ></b-form-input>
        </b-form-group>
        <b-form-group label="Copies:" label-for="amount">
          <b-form-input
            id="amount"
            v-model="form.amount"
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
            reqired
            :state="form.content.length >= 10 && form.content.length <= 2000"
            rows="3"
            max-rows="10"
          ></b-form-textarea>
        </b-form-group>

        <label for="tags-basic">Type a new tag and press enter</label>
        <b-form-tags input-id="tags-basic" remove-on-delete v-model="form.categories"></b-form-tags>
        <br />
        <b-button type="submit" variant="primary">Submit</b-button>
      </b-form>
      <b-card class="mt-3" header="Form Data Result">
        <pre class="m-0" style="color: white">{{ form }}</pre>
      </b-card>
    </div>
    <div v-else>
      <Post />
    </div>
  </div>
</template>

<script>
import axios from "axios";
import Post from "@/components/GetPost";
export default {
  components: {
    Post,
  },
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
      for (let i = 0; i < this.form.amount; i++) {
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
      }
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

.columns {
  display: flex;
  flex-wrap: wrap;
  margin: 0 17%;
}

.columns > * {
  flex-basis: calc(calc(750px - 100%) * 999);
}

.card {
  margin: 20px;
  padding: 10px;
  background-color: rgba(255, 255, 255, 0.05);
  box-shadow: 5px 5px 6px 2px rgba(10, 10, 10, 0.3);
}

.post-col {
  flex-grow: 1;
}

.info-col {
  flex-grow: 0.4;
}

hr {
  opacity: 0.3;
}
</style>
