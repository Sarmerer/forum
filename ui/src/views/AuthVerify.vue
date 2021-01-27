<template>
  <div></div>
</template>
<script>
import { mapActions, mapGetters } from "vuex";
export default {
  data() {
    return {
      code: this.$route.query.code,
    };
  },
  computed: {
    ...mapGetters({ authError: "auth/authError" }),
  },
  created() {
    if (!this.code) return this.$router.push("/");
    this.verifyEmail(this.code).then(() => {
      let error = this.authError?.data || this.authError?.data;
      if (error) {
        this.$bvToast.toast(error.message, {
          title: "Oops!",
          variant: "danger",
          solid: true,
        });
        this.$store.commit("auth/setAuthError", null);
      }
      this.$router.push("/");
    });
  },
  methods: {
    ...mapActions({
      verifyEmail: "auth/verify",
    }),
  },
};
</script>
