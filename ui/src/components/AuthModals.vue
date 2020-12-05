<template>
  <div>
    <b-modal
      id="signin-modal"
      modal-class="auth-modal"
      hide-footer
      hide-header
      centered
      no-stacking
      title="Sign In"
    >
      <h4 align="center">SIGN IN</h4>
      <SignInForm />
      <small align="center">
        <p>
          Don't have an account yet?
          <span class="secondary" v-b-modal.signup-modal>Sign up</span>
        </p></small
      >
    </b-modal>
    <b-modal
      id="signup-modal"
      title="Sign Up"
      modal-class="auth-modal"
      hide-footer
      hide-header
      centered
      no-stacking
    >
      <h4 align="center">SIGN UP</h4>
      <SignUpForm />

      <small align="center">
        <p>
          Already have an account?
          <span class="secondary" v-b-modal.signin-modal>Sign in</span>
        </p>
      </small>
    </b-modal>
  </div>
</template>
<script>
import SignInForm from "@/components/forms/SignInForm";
import SignUpForm from "@/components/forms/SignUpForm";

export default {
  components: {
    SignInForm,
    SignUpForm,
  },
};
</script>
<style>
.auth-modal > .modal-dialog {
  max-width: 400px;
}

.auth-modal > .modal-dialog > .modal-content {
  background-color: rgba(255, 255, 255, 0.05);
  box-shadow: 5px 5px 6px 2px rgba(10, 10, 10, 0.3);
  /* text-align: center; */
  background-image: url("../assets/img/home-hero.jpg");
  cursor: context-menu;
  padding: 30px 10px;
}
</style>
