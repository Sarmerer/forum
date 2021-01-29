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
    <div class="flex-container">
      <component
        class="auth"
        :is="activeComponent"
        v-on:changeActiveComponent="setActiveComponent"
        v-bind="componentProps"
      ></component>
    </div>
  </div>
</template>
<script>
export default {
  props: {
    prevRoute: { type: String, default: "/" },
    active: { type: String, default: "SignInForm" },
  },
  components: {
    SignInForm: () => import("@/components/forms/SignInForm"),
    SignUpForm: () => import("@/components/forms/SignUpForm"),
    Verify: () => import("@/components/auth/Verify"),
    Welcome: () => import("@/components/auth/Welcome"),
  },
  data() {
    return {
      componentProps: {},
      activeComponent: this.active,
    };
  },
  created() {
    this.componentProps.prevRoute = this.prevRoute;
  },
  methods: {
    setActiveComponent(event) {
      if (event.props) this.componentProps = event.props;
      if (typeof event === "string") return (this.activeComponent = event);
      if (event.component) return (this.activeComponent = event.component);
      this.activeComponent = "SignInForm";
    },
    successfulAuthHandler(event) {
      this.componentProps = { oldUser: event === "signIn" };
      this.activeComponent = "Welcome";
      this.redirect();
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
