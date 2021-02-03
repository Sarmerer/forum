<template>
  <b-form @submit.prevent="submitSignIn">
    <h4>SIGN IN</h4>
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
      <span class="mb-1">or</span>
      <div class="oauths">
        <b-button
          size="sm"
          variant="outline-dark"
          class="text-white-50"
          :href="
            `https://github.com/login/oauth/authorize?client_id=df41c45c5f1e0a5b29fe&redirect_uri=https://forum.sarmerer.ml/auth/github?redirect=${prevRoute}`
          "
          v-b-tooltip.hover.left="'GitHub'"
        >
          <b-icon size="sm" icon="github"></b-icon>
        </b-button>
        <b-button
          size="sm"
          variant="outline-dark"
          class="text-white-50"
          v-b-tooltip.hover.top="'Google'"
          href="https://accounts.google.com/o/oauth2/v2/auth?client_id=686483498224-2430dd197rbet1ck3tig05vtoub0ocvb.apps.googleusercontent.com&scope=https://www.googleapis.com/auth/userinfo.profile%20https://www.googleapis.com/auth/userinfo.email&access_type=offline&response_type=code&redirect_uri=https://forum.sarmerer.ml/auth/google"
        >
          <b-icon icon="google"></b-icon>
        </b-button>
        <b-button
          size="sm"
          variant="outline-dark"
          class="text-white-50"
          v-b-popover.hover.right="
            'Only your public data is being stored, because we care'
          "
        >
          <b-icon icon="info"></b-icon>
        </b-button>
      </div>
    </b-form-group>
    <small>
      <p class="text-white-50">
        Don't have an account?
        <a
          class="secondary"
          @click="$emit('changeActiveComponent', 'SignUpForm')"
          >Sign up</a
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
      provider: this.$route.params.provider || "",
      requesting: false,
    };
  },
  mounted() {
    if (this.provider) {
      if (this.requesting) return;
      this.requesting = true;
      let query = this.$route.query;
      query.provider = this.provider;
      this.OAuth(query).then(() => {
        let error = this.authError?.data || this.authError?.data;
        if (error) {
          if (error?.code === 409) {
            this.$router.push({
              name: "AuthMerge",
              params: {
                accounts: error?.data,
              },
            });
            this.$store.commit("auth/setAuthError", null);
          }
        } else {
          this.$emit("changeActiveComponent", {
            component: "Welcome",
            props: { prevRoute: this.$route?.query?.redirect },
          });
        }
        this.requesting = false;
      });
    }
  },
  methods: {
    ...mapActions({
      signIn: "auth/signIn",
      OAuth: "auth/OAuth",
    }),
    submitSignIn() {
      if (this.requesting) return;
      this.requesting = true;
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
          this.$emit("changeActiveComponent", {
            component: "Welcome",
            props: { prevRoute: this.prevRoute },
          });
        }
        this.requesting = false;
      });
    },
  },
};
</script>
<style lang="scss" scoped>
.oauths {
  display: flex;
  gap: 10px;
  flex-shrink: 1;
  justify-content: center;
}
</style>
