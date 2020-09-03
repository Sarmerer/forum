<template>
  <div>
    <b-form @submit="onSubmit">
      <b-form-group id="input-group-1" label="Login:" label-for="input-1">
        <b-form-input id="input-1" v-model="form.login" type="text" required placeholder="Enter login"></b-form-input>
      </b-form-group>

      <b-form-group id="input-group-2" label="Password:" label-for="input-2">
        <b-form-input id="input-2" v-model="form.password" type="password" required placeholder="Enter password"></b-form-input>
      </b-form-group>

      <b-button type="submit" variant="primary">Submit</b-button>
    </b-form>
    <b-card class="mt-3" header="Form Data Result">
      <pre class="m-0">{{ response }}</pre>
    </b-card>
  </div>
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
    //FIXME browser ignores Set-Cookie header
    onSubmit(evt) {
      evt.preventDefault();
      let obj = this
      const body = { login: this.form.login, password: this.form.password };
      axios
        .post(`https://localhost:4433/api/auth/signin` + "?" + "API_KEY=2ef1e8c9-5c6f-4b40-873a-bcba13fc249b", body)
        .then(function(response) {
          obj.response = response;
        })
        .catch(function(error) {
          obj.response = error;
        });
    },
  },
};
</script>
