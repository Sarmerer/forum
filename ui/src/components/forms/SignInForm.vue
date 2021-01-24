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
      <span class="mb-1">or</span>
      <div class="oauths">
        <b-button
          size="sm"
          variant="outline-dark"
          class="text-white-50"
          href="https://github.com/login/oauth/authorize?client_id=df41c45c5f1e0a5b29fe"
          v-b-tooltip.hover.left="'GitHub'"
        >
          <b-icon size="sm" icon="github"></b-icon>
        </b-button>
        <b-button
          size="sm"
          variant="outline-dark"
          class="text-white-50"
          v-b-tooltip.hover.top="'Google'"
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
  gap: 10px;
  flex-shrink: 1;
  justify-content: center;
}
</style>
