<template>
  <b-form @submit.prevent="submitSignUp">
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
        v-model="form.email"
        type="text"
        required
        placeholder="Email"
        class="mb-2"
      ></b-form-input>
      <b-form-input
        autocomplete="off"
        v-model="form.password"
        type="password"
        required
        placeholder="Password"
        class="mb-2"
        :state="passwordsMatch && passwordLength"
      ></b-form-input>
      <b-form-input
        v-model="form.passwordConfirm"
        id="password"
        type="password"
        required
        placeholder="Confirm Password"
        class="mb-2"
        invalid-feedback="Passwords don't match"
        :state="passwordsMatch && passwordLength"
      ></b-form-input>
      <b-form-invalid-feedback id="password-feedback">
        {{ passwordsMatchFeedback }}
      </b-form-invalid-feedback>
    </b-form-group>
    <b-button
      type="submit"
      class="modal-button mb-4 w-100"
      :disabled="!validform"
      variant="info"
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
      return (
        this.passwordsMatch && this.form.login.length && this.form.email.length
      );
    },
    passwordsMatch() {
      return this.form.password === this.form.passwordConfirm;
    },
    passwordLength() {
      return (
        this.form.password.length > 0 && this.form.passwordConfirm.length > 0
      );
    },
    passwordsMatchFeedback() {
      return !this.passwordLength
        ? "Enter passwords"
        : !this.passwordsMatch
        ? "Passwords don't match"
        : "";
    },
  },
  data() {
    return {
      form: {
        login: "",
        email: "",
        password: "",
        passwordConfirm: "",
      },
    };
  },
  methods: {
    ...mapActions({
      signIn: "auth/signIn",
    }),
    submitSignUp() {
      console.log(this.form);
      this.signUp(this.form).then(() => {
        if (this.authError?.data?.message) {
          this.$bvToast.toast(this.authError.data.message, {
            title: "Oopsie!",
            variant: "danger",
            solid: true,
          });
          this.$store.commit("auth/setAuthError", null);
        }
      });
    },
  },
};
</script>
