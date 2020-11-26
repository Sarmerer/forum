<template>
  <div>
    <b-modal
      id="signout-modal"
      :modal-class="authmodal"
      hide-footer
      hide-header
      centered
      no-stacking
      title="Sign Out"
    >
      <h4>SIGN OUT</h4>
      <p>Are you sure? You might have to enter your looong password again...</p>
      <b-button
        class="modal-button mb-4"
        @click="$bvModal.hide('signout-modal')"
        >Oh no!</b-button
      >

      <small
        ><p>
          Yeah, I really need to log out...
          <a class="secondary" style="cursor: pointer" v-on:click="signOut"
            >Bye!</a
          >
        </p></small
      >
    </b-modal>
  </div>
</template>
<script>
import { mapActions } from "vuex";

export default {
  // mixins: [validationMixin],
  data() {
    return {
      authmodal: ["auth-modal"],
    };
  },

  methods: {
    ...mapActions({
      signOut: "auth/signOut",
    }),
  },
};
</script>

