<template>
  <b-navbar type="dark" fixed="bottom">
    <b-navbar-nav>
      <b-nav-item>
        <router-link to="/">
          <b-icon-house-door></b-icon-house-door>
        </router-link>
      </b-nav-item>
    </b-navbar-nav>
    <b-navbar-nav v-if="authenticated">
      <b-nav-item>
        <router-link to="/new-post"><b-icon-pen></b-icon-pen></router-link>
      </b-nav-item>
    </b-navbar-nav>
    <b-navbar-nav v-if="!authenticated">
      <b-nav-item>
        <a v-b-modal.signin-modal><b-icon-door-closed></b-icon-door-closed> </a>
        <AuthModals />
      </b-nav-item>
    </b-navbar-nav>
    <b-navbar-nav v-if="authenticated">
      <b-nav-item>
        <router-link :to="'/user/' + user.id">
          <b-img
            width="25px"
            v-b-tooltip.hover
            :title="user.display_name"
            :src="user.avatar"
          ></b-img>
        </router-link>
      </b-nav-item>
    </b-navbar-nav>
    <b-navbar-nav v-if="authenticated">
      <b-nav-item>
        <a @click.prevent="signOut"><b-icon icon="door-open"></b-icon></a>
      </b-nav-item>
    </b-navbar-nav>
    <b-navbar-nav v-if="authenticated && user.role > 0">
      <b-nav-item>
        <router-link to="/dashboard/admin">
          <b-icon-wrench></b-icon-wrench>
        </router-link>
      </b-nav-item>
    </b-navbar-nav>
  </b-navbar>
</template>
<script>
import { mapActions, mapGetters } from "vuex";
import AuthModals from "@/components/AuthModals";

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
    AuthModals,
  },
};
</script>
