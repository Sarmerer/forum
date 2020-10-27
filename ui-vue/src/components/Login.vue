<template>
  <div>
    <b-modal
      id="auth-modal"
      :modal-class="authmodal"
      hide-footer
      hide-header
      centered
      no-stacking
      title="Sign In"
    >
      <h4>SIGN IN</h4>
      <b-form @submit.prevent="submit">
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
          <span class="secondary" v-b-modal.signup-modal>Sign up</span>
        </p></small
      >
    </b-modal>

    <b-modal
      id="signup-modal"
      title="Sign Up"
      :modal-class="authmodal"
      hide-footer
      hide-header
      centered
      no-stacking
    >
      <h4>SIGN UP</h4>

      <b-form @submit.prevent="submit">
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
          <span class="secondary" v-b-modal.auth-modal>Sign in</span>
        </p></small
      >
    </b-modal>
  </div>
</template>
<script>
import { mapActions } from "vuex";

export default {
  data() {
    return {
      form: {
        login: "",
        password: ""
      },
      response: "",
      authmodal: ["auth-modal"]
    };
  },
  methods: {
    ...mapActions({
      signIn: "auth/signIn"
    }),
    submit() {
      this.signIn(this.form);
    }
  }
};
</script>
<style>
.sidenav a {
  cursor: pointer;
  color: #818181 !important;
}

.sidenav a:hover {
  color: #f1f1f1 !important;
}

.auth-modal > .modal-dialog {
  max-width: 400px;
}

.auth-modal > .modal-dialog > .modal-content {
  background-color: rgba(255, 255, 255, 0.05);
  box-shadow: 5px 5px 6px 2px rgba(10, 10, 10, 0.3);
  text-align: center;
  background-image: url("../assets/img/home-hero.jpg");
  cursor: context-menu;
  padding: 30px 10px;
}

.modal-button {
  background-color: #278ea5;
  border: none;
  display: block;
  width: 100%;
}

.modal-button:hover {
  background-color: #278ea5;
  opacity: 0.8;
}
</style>
