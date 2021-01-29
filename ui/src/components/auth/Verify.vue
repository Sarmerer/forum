<template>
  <div>
    <b-form @submit.prevent="sendVerificationCode">
      <h4 class="mb-2 text-white-87">
        VERIFY
      </h4>
      <small class="text-white-87">Check your mail: {{ email }}</small>
      <br />
      <small class="text-white-87">Code will expire in 10 minutes</small>
      <b-form-group>
        <b-input-group size="sm">
          <b-form-input
            placeholder="Code"
            type="text"
            autocomplete="off"
            v-model="code"
          >
          </b-form-input>
          <b-input-group-append>
            <b-button
              :disabled="!code.length"
              size="sm"
              type="submit"
              variant="outline-dark"
              class="text-white-50"
            >
              <b-icon icon="arrow-right"></b-icon>
            </b-button>
          </b-input-group-append>
        </b-input-group>
        <small class="d-flex text-left text-white-87">
          <span
            v-if="timer <= 0"
            class="text-white-87 text-decoration-underline"
            @click="resendEmail"
            >Resend code
          </span>
          <span v-else>
            Resend code after: {{ timer }}
            {{ timer % 10 == 1 ? "second" : "seconds" }}
          </span>
        </small>
      </b-form-group>
    </b-form>
  </div>
</template>
<script>
import { mapActions, mapGetters } from "vuex";
import api from "@/api/api";

export default {
  props: {
    email: { type: String, default: "example@gmail.com" },
    prevRoute: String,
  },
  data() {
    return {
      code: "",
      cooldown: 60,
      timer: 0,
      interval: null,
      requesting: false,
    };
  },
  computed: {
    ...mapGetters({ authError: "auth/authError" }),
  },
  created() {
    this.startTimer();
  },
  beforeRouteLeave(to, from, next) {
    clearInterval(this.interval);
    next();
  },
  methods: {
    ...mapActions({
      verifyEmail: "auth/verify",
    }),
    startTimer() {
      this.timer = this.cooldown;
      this.interval = setInterval(() => {
        this.timer === 0 ? clearInterval(this.interval) : this.timer--;
      }, 1000);
    },
    async resendEmail() {
      if (this.timer > 0) return;
      this.startTimer();
      await api
        .post(`auth/send-verification?email=${this.email}`)
        .catch((error) => {
          if (error?.data?.status === 498) {
            this.expiredError();
            this.$emit("changeActiveComponent", "SignIn");
            return;
          }
          if (error?.data?.message)
            this.$bvToast.toast(error.data.message, {
              title: "Oops!",
              variant: "danger",
              solid: true,
            });
        });
    },
    sendVerificationCode() {
      if (this.requesting) return;
      this.requesting = true;
      this.verifyEmail({ code: this.code, email: this.email }).then(() => {
        let error = this.authError?.data || this.authError?.data;
        if (error) {
          if (error.status === 498) {
            this.expiredError();
            this.$emit("changeActiveComponent", "SignIn");
            return;
          }
          this.$bvToast.toast(error.message, {
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
    expiredError() {
      this.$bvToast.toast("Code has expired, you haven't used it in time", {
        title: "Oops!",
        variant: "danger",
        solid: true,
      });
    },
  },
};
</script>
