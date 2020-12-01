<template>
  <div class="sidenav">
    <router-link
      to="/"
      v-b-tooltip.hover.rightbottom="{ customClass: 'tooltip-sidenav' }"
      title="Home"
    >
      <b-icon-house-door></b-icon-house-door>
    </router-link>
    <div
      v-if="!authenticated"
      v-b-tooltip.hover.rightbottom="{ customClass: 'tooltip-sidenav' }"
      title="Sign in"
    >
      <a v-b-modal.auth-modal><b-icon-door-closed></b-icon-door-closed></a>
      <Login />
    </div>
    <div v-if="authenticated">
      <router-link
        to="/new-post"
        v-b-tooltip.hover.rightbottom="{ customClass: 'tooltip-sidenav' }"
        title="New post"
      >
        <b-icon-pen></b-icon-pen>
      </router-link>

      <router-link :to="'/user/' + user.id">
        <b-img
          width="30px"
          v-b-tooltip.hover.rightbottom="{ customClass: 'tooltip-sidenav' }"
          :title="user.display_name"
          :src="user.avatar"
        ></b-img>
      </router-link>

      <a
        @click.prevent="signOut"
        v-b-tooltip.hover.rightbottom="{ customClass: 'tooltip-sidenav' }"
        title="Sign out"
      >
        <b-icon-door-open></b-icon-door-open>
      </a>
      <router-link
        to="/dashboard/admin"
        v-b-tooltip.hover.rightbottom="{ customClass: 'tooltip-sidenav' }"
        title="Dashboard"
      >
        <b-icon-wrench></b-icon-wrench>
      </router-link>
    </div>
  </div>
</template>
<script>
import { mapActions, mapGetters } from "vuex";
import Login from "@/components/Login";

export default {
  computed: {
    ...mapGetters({
      authenticated: "auth/authenticated",
      user: "auth/user",
    }),
  },
  methods: {
    ...mapActions({
      signOut: "auth/signOut",
    }),
  },
  components: {
    Login,
  },
};
</script>
<style >
.tooltip-sidenav {
  margin-top: 22px;
}
</style>
