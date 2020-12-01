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
      <div v-if="show" class="auth">
        <h4>SIGN IN</h4>
        <b-form @submit.prevent="submitSignIn">
          <b-input-group>
            <b-input
              v-model="form.login"
              type="text"
              required
              placeholder="Login"
              class="mb-2 mt-4"
            ></b-input>
            <b-form-invalid-feedback id="input-login-live-feedback"
              >Login is required.</b-form-invalid-feedback
            >
          </b-input-group>

          <b-input-group>
            <b-input
              v-model="form.password"
              type="password"
              required
              placeholder="Password"
              class="mb-2"
            ></b-input>
          </b-input-group>
          <b-button type="submit" class="modal-button mb-4">Go</b-button>
        </b-form>
        <small
          ><p>
            Don't have an account yet?
            <a class="secondary" v-on:click="show = false">Sign up</a>
          </p></small
        >
      </div>
      <div v-else class="auth">
        <h4>SIGN UP</h4>

        <b-form @submit.prevent="submitSignUp">
          <b-input-group>
            <b-input
              v-model="form.login"
              type="text"
              required
              placeholder="Login"
              class="mb-2 mt-4"
            ></b-input>
          </b-input-group>

          <b-input-group>
            <b-input
              v-model="form.email"
              type="text"
              required
              placeholder="Email"
              class="mb-2"
            ></b-input>
          </b-input-group>

          <b-input-group>
            <b-input
              v-model="form.password"
              type="password"
              required
              placeholder="Password"
              class="mb-2"
            ></b-input>
          </b-input-group>

          <b-input-group>
            <b-input
              v-model="form.repeatPassword"
              type="password"
              required
              placeholder="Re-enter Password"
              class="mb-2"
            ></b-input>
          </b-input-group>
          <b-button type="submit" class="modal-button mb-4">Go</b-button>
        </b-form>
        <small
          ><p>
            Already have an account?
            <a class="secondary" v-on:click="show = true">Sign in</a>
          </p></small
        >
      </div>
    </div>
  </div>
</template>



<script>
import { mapActions } from "vuex";
// import { validationMixin } from "vuelidate";
// import { required, minLength } from "vuelidate/lib/validators";

export default {
  // mixins: [validationMixin],
  data() {
    return {
      show: true,
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
      signUp: "auth/signUp",
    }),
    submitSignIn() {
      this.signIn(this.form);
      this.$router.push("/");
    },
    submitSignUp() {
      this.signUp(this.form);
      this.$router.push("/");
    },
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