<template>
  <div class="sidenav">
    <router-link to="/"
      ><span class="primary"><b-icon icon="house-door"></b-icon></span
    ></router-link>
    <div v-if="!authenticated">
      <a v-b-modal.auth-modal><b-icon icon="door-closed"></b-icon></a>
      <Login />
    </div>
    <div v-if="authenticated">
      <router-link to="/new-post"><b-icon icon="pen"></b-icon></router-link>

      <router-link :to="'/user/' + user.id"
        ><b-icon icon="person"></b-icon
      ></router-link>

      <router-link to="/"><b-icon icon="gear"></b-icon></router-link>

      <a @click.prevent="signOut"><b-icon icon="door-open"></b-icon></a>
    </div>
    <router-link to="/dashboard/admin"
      ><b-icon icon="hammer"></b-icon
    ></router-link>
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
