<template>
  <b-form inline @submit.prevent="submit">
    <b-input-group class="mb-2 mr-sm-2 mb-sm-0">
      <b-input style="width: 150px;height: 30px" v-model="form.login" type="text" required placeholder="Login"></b-input>
    </b-input-group>

    <b-input-group class="mb-2 mr-sm-2 mb-sm-0">
      <b-input style="width: 150px;height: 30px" v-model="form.password" type="password" required placeholder="Password"></b-input>
    </b-input-group>
    <b-button style="width: 50px;height: 30px" type="submit" variant="primary">Go</b-button>
  </b-form>
</template>
<script>
import { mapActions } from "vuex";

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
    ...mapActions({
      signIn: "auth/signIn",
    }),
    submit() {
      this.signIn(this.form);
    },
  },
};
</script>
