<template>
  <div>
    <h4>
      Welcome{{ oldUser ? " back, " : ", " }}
      <span class="secondary">{{ user.alias }}</span
      >!
    </h4>
    <h6>You will be redirected back in {{ timeLeft }}</h6>
    <br />
    <b-button variant="outline-light" @click="redirect">
      Go back
    </b-button>
  </div>
</template>
<script>
import { mapGetters } from "vuex";

export default {
  props: {
    oldUser: Boolean,
    prevRoute: String,
  },
  data() {
    return {
      timeLeft: 5,
      interval: null,
      redirectParam: this.$route.query.redirect,
    };
  },
  beforeRouteLeave(from, to, next) {
    this.timeLeft = 5;
    clearInterval(this.interval);
    next();
  },
  computed: {
    ...mapGetters({
      user: "auth/user",
    }),
  },
  created() {
    this.interval = setInterval(() => {
      this.timeLeft <= 0 ? this.redirect() : this.timeLeft--;
    }, 1000);
  },
  methods: {
    emitRedirect() {
      this.$emit("redirect", this.$route?.query?.redirect);
    },
    redirect() {
      if (this.interval) clearInterval(this.interval);
      this.timeLeft = 5;
      let next = this.redirectParam || this.prevRoute;
      if (next !== this.$route.path) this.$router.push(next);
    },
  },
};
</script>
