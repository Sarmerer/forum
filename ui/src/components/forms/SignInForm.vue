<template>
  <b-form @submit.prevent="submitSignIn">
    <b-form-group>
      <b-form-input
        autocomplete="off"
        v-model="form.login"
        type="text"
        required
        placeholder="Login"
        class="mb-2 mt-4"
      ></b-form-input>
      <b-form-input
        autocomplete="off"
        v-model="form.password"
        type="password"
        required
        placeholder="Password"
        class="mb-2"
      ></b-form-input>
    </b-form-group>
    <b-button
      type="submit"
      class="modal-button mb-4 w-100"
      variant="info"
      :disabled="!validform"
      >Go</b-button
    >
  </b-form>
</template>
<script>
import { mapActions, mapGetters } from "vuex";

export default {
  computed: {
    ...mapGetters({ authError: "auth/authError" }),
    validform() {
      return this.form.login.length > 0 && this.form.password.length > 0;
    },
  },
  data() {
    return {
      form: {
        login: "",
        password: "",
      },
    };
  },
  methods: {
    ...mapActions({
      signIn: "auth/signIn",
    }),
    submitSignIn() {
      this.signIn(this.form).then(() => {
        let error = this.authError?.data?.message || this.authError?.data;
        if (error) {
          this.$bvToast.toast(error, {
            title: "Oops!",
            variant: "danger",
            solid: true,
          });
          this.$store.commit("auth/setAuthError", null);
        } else {
          this.$emit("success", "signin");
        }
      });
    },
  },
};
</script>
