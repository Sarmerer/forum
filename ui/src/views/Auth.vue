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
      <div v-if="!signUpPageLocal" class="auth">
        <h4 align="center">SIGN IN</h4>
        <SignInForm :prevRoute="prevRoute" />
        <small
          ><p>
            Don't have an account yet?
            <a class="secondary" @click="signUpPageLocal = true">Sign up</a>
          </p></small
        >
      </div>
      <div v-else class="auth">
        <h4 align="center">SIGN UP</h4>
        <SignUpForm :prevRoute="prevRoute" />
        <small
          ><p>
            Already have an account?
            <a class="secondary" @click="signUpPageLocal = false">Sign in</a>
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
    prevRoute: { type: String },
    signUpPage: { type: Boolean, default: false },
  },
  components: {
    SignUpForm,
    SignInForm,
  },
  created() {
    if (this.signUpPage) this.signUpPageLocal = true;
  },
  data() {
    return {
      signUpPageLocal: false,
    };
  },
};
</script>

<style scoped>
.center {
  height: 400px;
  position: relative;
}

.center .auth {
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
