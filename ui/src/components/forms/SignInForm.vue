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
      <b-input-group>
        <b-form-input
          autocomplete="off"
          v-model="form.password"
          type="password"
          required
          placeholder="Password"
        ></b-form-input>
        <template #append>
          <b-button type="submit" variant="dark" :disabled="!validform">
            <b-icon icon="arrow-right"></b-icon>
          </b-button>
        </template>
      </b-input-group>
      <span class="mb-1 d-flex justify-content-center"> or</span>
      <div class="oauths">
        <b-button
          variant="outline-light"
          class="w-100 text-white-50"
          href="https://github.com/login/oauth/authorize?client_id=df41c45c5f1e0a5b29fe"
        >
          <b-icon icon="github"></b-icon> GitHub
        </b-button>
        <b-button variant="outline-primary" class="w-100 text-white-50">
          <b-icon icon="google"></b-icon> Google
        </b-button>
      </div>
    </b-form-group>
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
  mounted() {
    if (this.$route?.query?.provider) {
      this.OAuth(this.$route.query).then(() => {
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
    }
  },
  methods: {
    ...mapActions({
      signIn: "auth/signIn",
      OAuth: "auth/OAuth",
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
<style lang="scss" scoped>
.oauths {
  display: flex;
  gap: 5px;
}
</style>
