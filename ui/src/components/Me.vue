<template>
  <div>
    <b-dropdown id="dropdown" right :text="user.display_name" class="m-md-2">
      <b-dropdown-item><router-link :to="'/user/' + user.id">Profile</router-link></b-dropdown-item>
      <b-dropdown-item><router-link to="/">Settings</router-link></b-dropdown-item>
      <b-dropdown-divider></b-dropdown-divider>
      <b-dropdown-item @click.prevent="signOut">Logout</b-dropdown-item>
    </b-dropdown>
  </div>
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
