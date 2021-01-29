<template>
  <b-form @submit.prevent="submitSignUp">
    <h4>SIGN UP</h4>
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
        :state="emailState"
      ></b-form-input>

      <b-form-input
        autocomplete="off"
        v-model="form.password"
        :type="passwordInputType"
        required
        placeholder="Password"
        class="mb-2"
      ></b-form-input>

      <b-input-group class="mb-2">
        <b-form-input
          v-model="form.passwordConfirm"
          id="password"
          :type="passwordInputType"
          required
          placeholder="Confirm password"
          :state="passwordsState"
        ></b-form-input>

        <b-input-group-append>
          <b-button
            variant="darker text-white-50"
            v-b-tooltip.hover.right="'Show password'"
            size="sm"
            @click="showPassword"
          >
            <b-icon-eye> </b-icon-eye>
          </b-button>
        </b-input-group-append>
        <b-form-invalid-feedback
          id="password-feedback"
          v-if="passwordsState !== null"
        >
          {{ passwordsMatchFeedback }}
        </b-form-invalid-feedback>
      </b-input-group>

      <b-input-group>
        <b-form-input
          v-model="form.adminToken"
          type="password"
          placeholder="Admin token"
        ></b-form-input>
        <template #append>
          <b-button
            size="sm"
            variant="darker text-white-50"
            v-b-tooltip.hover.right="
              'You will be granted admin rights, if a valid admin token is entered here'
            "
          >
            <b-icon icon="info"></b-icon>
          </b-button>
        </template>
      </b-input-group>
    </b-form-group>
    <b-button
      type="submit"
      class="modal-button mb-4 w-100"
      :disabled="!validForm || requesting"
      variant="info"
      >Go</b-button
    >
    <small>
      <p class="text-white-50">
        Already have an account?
        <a
          class="secondary"
          @click="$emit('changeActiveComponent', 'SignInForm')"
          >Sign in</a
        >
      </p>
    </small>
  </b-form>
</template>
<script>
import { mapActions, mapGetters } from "vuex";

export default {
  props: { prevRoute: String },
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
    validEmail() {
      const re = /^(([^<>()[\]\\.,;:\s@"]+(\.[^<>()[\]\\.,;:\s@"]+)*)|(".+"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$/;
      return re.test(String(this.form.email).toLowerCase());
    },
    emailState() {
      return this.form.email.length === 0 ? null : this.validEmail;
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
      passwordInputType: "password",
      minLoginLength: 1,
      maxLoginLength: 15,
      minPasswordLength: 5,
      requesting: false,
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
      if (this.requesting) return;
      this.requesting = true;
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
          this.$emit("changeActiveComponent", {
            component: "Verify",
            props: { email: this.form.email, prevRoute: this.prevRoute },
          });
        }
        this.requesting = false;
      });
    },
    showPassword() {
      this.passwordInputType =
        this.passwordInputType === "password" ? "text" : "password";
    },
  },
};
</script>
