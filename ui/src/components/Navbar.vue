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
        <router-link to="/auth">
          <b-icon-door-closed> </b-icon-door-closed>
        </router-link>
      </b-nav-item>
    </b-navbar-nav>
    <b-navbar-nav v-if="authenticated">
      <b-nav-item>
        <router-link :to="'/user/' + user.id">
          <b-avatar
            size="25px"
            variant="dark"
            v-b-tooltip.hover.rightbottom="{ customClass: 'tooltip-sidenav' }"
            :title="user.display_name"
            :src="user.avatar"
          ></b-avatar>
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
        <router-link v-if="user.role === 2" to="/dashboard/admin">
          <b-icon-wrench></b-icon-wrench>
        </router-link>
      </b-nav-item>
    </b-navbar-nav>
  </b-navbar>
</template>
<script>
import { mapActions, mapGetters } from "vuex";

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
};
</script>
