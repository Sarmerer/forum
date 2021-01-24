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
    <div class="flex-container" v-if="!showWelcomeMessage">
      <div v-if="!signUpPage && !signUpPageLocal" class="auth">
        <h4 align="center">SIGN IN</h4>
        <SignInForm v-on:success="successfulAuth" />
        <small>
          <p class="text-white-50">
            Don't have an account?
            <a class="secondary" @click="signUpPageLocal = true">Sign up</a>
          </p></small
        >
      </div>
      <div v-else class="auth">
        <h4 align="center">SIGN UP</h4>
        <SignUpForm v-on:success="successfulAuth" />
        <small>
          <p class="text-white-50">
            Already have an account?
            <a class="secondary" @click="signUpPageLocal = false">Sign in</a>
          </p></small
        >
      </div>
    </div>
    <div align="center" class="welcome-message" v-else>
      <h3>
        Welcome{{ oldUser ? " back, " : ", " }}
        <span class="secondary">{{ user.alias }}</span
        >!
      </h3>
      <h6>You will be redirected back in {{ timeLeft }}</h6>
      <br />
      <b-button variant="outline-light" @click="redirect">Go back</b-button>
    </div>
  </div>
</template>

<script>
import SignInForm from "@/components/forms/SignInForm";
import SignUpForm from "@/components/forms/SignUpForm";
import { mapGetters } from "vuex";

export default {
  props: {
    prevRoute: String,
    signUpPage: { type: Boolean, default: false },
  },
  components: {
    SignUpForm,
    SignInForm,
  },
  computed: {
    ...mapGetters({
      user: "auth/user",
    }),
  },
  data() {
    return {
      signUpPageLocal: false,
      showWelcomeMessage: false,
      oldUser: false,
      timeLeft: 5,
      interval: null,
    };
  },
  beforeRouteLeave(from, to, next) {
    this.timeLeft = 0;
    clearInterval(this.interval);
    next();
  },
  methods: {
    successfulAuth(event) {
      this.oldUser = event === "signin";
      this.showWelcomeMessage = true;
      this.interval = setInterval(() => {
        this.timeLeft--;
        if (this.timeLeft === 0) {
          clearInterval(this.interval);
          this.timeLeft = 5;
          this.prevRoute
            ? this.$router.push(this.prevRoute)
            : this.$router.back();
        }
      }, 1000);
    },
    redirect() {
      if (this.interval) clearInterval(this.interval);
      this.timeLeft = 5;
      this.prevRoute ? this.$router.push(this.prevRoute) : this.$router.back();
    },
  },
};
</script>

<style scoped>
.flex-container {
  display: flex;
  align-items: center;
  text-align: center;
  justify-content: center;
}

.auth {
  display: flex;
  justify-content: center;
  flex-direction: column;
  width: 300px;
}

.auth a {
  cursor: pointer;
}

@media (max-width: 768px) {
  .auth {
    width: 50%;
  }
}

@media (max-width: 596px) {
  .auth {
    width: 80%;
  }
}

.welcome-message {
  opacity: 0.87;
}
</style>
