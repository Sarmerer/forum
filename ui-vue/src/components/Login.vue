<template>
  <b-form inline @submit="onSubmit">
    <b-input-group class="mb-2 mr-sm-2 mb-sm-0">
      <b-input style="width: 100px;height: 25px" v-model="form.login" type="text" placeholder="Login"></b-input>
    </b-input-group>

    <b-input-group class="mb-2 mr-sm-2 mb-sm-0">
      <b-input style="width: 100px;height: 25px" v-model="form.password" type="password" placeholder="Password"></b-input>
    </b-input-group>
    <b-form-checkbox class="mb-2 mr-sm-2 mb-sm-0">Remember me</b-form-checkbox>
    <b-button type="submit" variant="primary">Submit</b-button>
  </b-form>
</template>
<script>
import axios from "axios";

export default {
  data() {
    return {
      form: {
        login: "",
        password: "",
      },
      response: "",
    };
  },
  methods: {
    onSubmit(e) {
      e.preventDefault();
      axios.defaults.withCredentials = true;
      const body = { login: this.form.login, password: this.form.password };
      axios(`/api/auth/signin`, {
        method: "POST",
        data: body,
        withCredentials: true,
      })
        .then(function(response) {
          console.log(response);
        })
        .catch(function(error) {
          console.log(error);
        });
    },
  },
};
</script>
