<template>
  <b-navbar type="dark" fixed="bottom">
    <b-navbar-nav
      ><b-nav-item>
        <router-link to="/"
          ><span class="primary"
            ><b-icon
              icon="chat-left-dots"
            ></b-icon></span></router-link></b-nav-item
    ></b-navbar-nav>
    <b-navbar-nav
      ><b-nav-item>
        <router-link to="/post/new"
          ><b-icon icon="pen"></b-icon></router-link></b-nav-item
    ></b-navbar-nav>
    <b-navbar-nav v-if="!authenticated"
      ><b-nav-item
        ><a v-b-modal.auth-modal><b-icon icon="door-closed"></b-icon></a
        ><Login /></b-nav-item
    ></b-navbar-nav>
    <b-navbar-nav v-if="authenticated"
      ><b-nav-item>
        <router-link :to="'/user/' + user.id"
          ><b-icon icon="person"></b-icon
        ></router-link>
      </b-nav-item>
    </b-navbar-nav>
    <b-navbar-nav v-if="authenticated"
      ><b-nav-item>
        <router-link to="/"><b-icon icon="gear"></b-icon></router-link>
      </b-nav-item>
    </b-navbar-nav>
    <b-navbar-nav v-if="authenticated"
      ><b-nav-item>
        <a @click.prevent="signOut"><b-icon icon="door-open"></b-icon></a>
      </b-nav-item>
    </b-navbar-nav>
  </b-navbar>
</template>
<script>
import { mapActions, mapGetters } from "vuex";
import Login from "./Login";

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
  },
  components: {
    Login
  }
};
</script>