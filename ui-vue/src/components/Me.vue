<template>
  <div>
    <router-link :to="'/user/' + user.id"
      ><b-icon icon="person"></b-icon
    ></router-link>
    <router-link to="/"><b-icon icon="gear"></b-icon></router-link>
    <a @click.prevent="signOut"><b-icon icon="door-open"></b-icon></a>
    <!-- <b-dropdown id="dropdown-dropright" dropright class="m-md-2">
      <template #button-content>
        <b-icon icon="person"></b-icon>
      </template>
      <b-dropdown-item
        ><router-link :to="'/user/' + user.id"
          >Profile</router-link
        ></b-dropdown-item
      >
      <b-dropdown-item
        ><router-link to="/">Settings</router-link></b-dropdown-item
      >
      <b-dropdown-divider></b-dropdown-divider>
      <b-dropdown-item @click.prevent="signOut">Logout</b-dropdown-item>
    </b-dropdown> -->
  </div>
</template>
<script>
import { mapActions, mapGetters } from "vuex";
export default {
  computed: {
    ...mapGetters({
      authenticated: "auth/authenticated",
      user: "auth/user"
    })
  },
  methods: {
    ...mapActions({
      signOut: "auth/signOut"
    })
  }
};
</script>

