<template>
  <b-form @submit.prevent="submitSignUp">
    <b-form-group>
      <b-form-input
        autocomplete="off"
        id="login"
        v-model="form.login"
        type="text"
        required
        placeholder="Login"
        class="mb-2 mt-4"
        :state="loginState"
      ></b-form-input>
      <b-form-invalid-feedback
        class="mt-n2 mb-2"
        id="login-feedback"
        v-if="loginState !== null"
      >
        Login is too long
      </b-form-invalid-feedback>
      <b-form-input
        autocomplete="off"
        v-model="form.email"
        type="text"
        required
        placeholder="Email"
        class="mb-2"
        :state="!form.email.length ? null : true"
      ></b-form-input>

      <b-form-input
        autocomplete="off"
        v-model="form.password"
        type="password"
        required
        placeholder="Password"
        class="mb-2"
        :state="passwordsState"
      ></b-form-input>

      <b-form-input
        v-model="form.passwordConfirm"
        id="password"
        type="password"
        required
        placeholder="Confirm password"
        class="mb-2"
        :state="passwordsState"
      ></b-form-input>
      <b-form-invalid-feedback
        class="mt-n2 mb-2"
        id="password-feedback"
        v-if="passwordsState !== null"
      >
        {{ passwordsMatchFeedback }}
      </b-form-invalid-feedback>

      <b-form-checkbox v-model="form.admin">Sign up as admin</b-form-checkbox>
      <b-form-input
        v-if="form.admin"
        v-model="form.adminToken"
        type="password"
        required
        placeholder="Admin Token"
        class="mt-1"
      ></b-form-input>
    </b-form-group>
    <b-button
      type="submit"
      class="modal-button mb-4 w-100"
      :disabled="!validForm"
      variant="info"
      >Go</b-button
    >
  </b-form>
</template>
<script>
import { mapActions, mapGetters } from "vuex";

export default {
  computed: {
    ...mapGetters({ authError: "auth/authError", user: "auth/user" }),
    validForm() {
      return (
        this.passwordsMatch &&
        this.validPasswordLength &&
        this.form.email.length &&
        this.validLoginLength
      );
    },
    validLoginLength() {
      return (
        this.form.login.length >= this.minLoginLength &&
        this.form.login.length <= this.maxLoginLength
      );
    },
    loginState() {
      return this.form.login.length === 0 ? null : this.validLoginLength;
    },
    passwordsMatch() {
      return this.form.password === this.form.passwordConfirm;
    },
    validPasswordLength() {
      return this.form.password.length >= this.minPasswordLength;
    },
    passwordsState() {
      return !this.form.password.length && !this.form.passwordConfirm.length
        ? null
        : this.passwordsMatch && this.validPasswordLength;
    },
    passwordsMatchFeedback() {
      return !this.form.password.length || !this.form.passwordConfirm.length
        ? "Enter both passwords"
        : !this.passwordsMatch
        ? "Passwords don't match"
        : !this.validPasswordLength
        ? "Password is too short"
        : "";
    },
  },
  data() {
    return {
      minLoginLength: 1,
      maxLoginLength: 15,
      minPasswordLength: 5,
      form: {
        login: "",
        email: "",
        password: "",
        passwordConfirm: "",
        admin: false,
        adminToken: "",
      },
    };
  },
  methods: {
    ...mapActions({
      signUp: "auth/signUp",
    }),
    submitSignUp() {
      this.signUp(this.form).then(() => {
        let error = this.authError?.data?.message || this.authError?.data;
        if (error) {
          this.$bvToast.toast(error, {
            title: "Oops!",
            variant: "danger",
            solid: true,
          });
          this.$store.commit("auth/setAuthError", null);
        } else {
          this.$emit("success", "signup");
        }
      });
    },
  },
};
</script>
