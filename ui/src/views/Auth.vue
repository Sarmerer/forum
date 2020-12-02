<template>
  <div>
    <div class="hero-image">
      <div class="hero-text">
        <h1>
          WELCOME<br />
          TO <span class="primary">FORUM</span>
        </h1>
      </div>
    </div>
    <div class="center">
      <div v-if="!signUpForm" class="auth">
        <h4>SIGN IN</h4>
        <SignInForm />
        <small
          ><p>
            Don't have an account yet?
            <a class="secondary" v-on:click="signUp = true">Sign up</a>
          </p></small
        >
      </div>
      <div v-else class="auth">
        <h4>SIGN UP</h4>
        <SignUpForm />
        <small
          ><p>
            Already have an account?
            <a class="secondary" v-on:click="signUp = false">Sign in</a>
          </p></small
        >
      </div>
    </div>
  </div>
</template>

<script>
import SignInForm from "@/components/SignInForm";
import SignUpForm from "@/components/SignUpForm";

export default {
  props: {
    signUpForm: { type: Boolean, default: false },
    prevRoute: { type: String },
  },
  components: {
    SignUpForm,
    SignInForm,
  },
};
</script>

<style scoped>
.center {
  height: 400px;
  position: relative;
}

.center .auth {
  text-align: center;
  margin: 0;
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
}

.center .auth a {
  cursor: pointer;
}

@media (max-width: 768px) {
  .center .auth {
    width: 50%;
  }
}

@media (max-width: 596px) {
  .center .auth {
    width: 80%;
  }
}
</style>
