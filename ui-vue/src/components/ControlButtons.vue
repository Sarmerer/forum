<template>
  <b-overlay
    :show="requestPending"
    rounded
    opacity="0.5"
    spinner-small
    variant="light"
    spinner-variant="primary"
    class="d-inline-block"
  >
    <transition name="fade">
      <b-button-group v-if="!confirm" size="sm">
        <b-button
          size="sm"
          lg="1"
          class="controls-button"
          variant="light"
          title="Edit"
          @click="editFunction()"
        >
          <img src="@/assets/svg/post/edit.svg" alt="edit" srcset="" />
        </b-button>
        <b-button
          size="sm"
          variant="outline-danger"
          lg="2"
          @click="confirm = !confirm"
          class="controls-button"
          title="Delete"
        >
          <img src="@/assets/svg/post/delete.svg" alt="delete" srcset="" />
        </b-button>
      </b-button-group>

      <b-button-group v-if="confirm" size="sm">
        <b-button
          size="sm"
          variant="outline-success"
          lg="2"
          class="confirm"
          @click="deleteFunction()"
          title="Confirm"
        >
          <img src="@/assets/svg/post/confirm.svg" alt="delete" srcset="" />
        </b-button>
        <b-button
          size="sm"
          lg="1"
          variant="outline-danger"
          @click="confirm = !confirm"
          class="confirm"
          title="Dismiss"
        >
          <img src="@/assets/svg/post/dismiss.svg" alt="edit" srcset="" />
        </b-button>
      </b-button-group>
    </transition>
  </b-overlay>
</template>
<script>
import { mapGetters } from "vuex";

export default {
  props: {
    deleteFunction: { type: Function, required: true },
    editFunction: { type: Function },
    requestPending: { type: Boolean, required: true },
  },
  computed: {
    ...mapGetters({
      user: "auth/user",
    }),
  },
  data() {
    return {
      confirm: false,
    };
  },
};
</script>
<style lang="scss" scoped>
.confirm {
  box-shadow: none;
  -moz-box-shadow: none;
  -webkit-box-shadow: none;
}
.controls-button {
  background-color: transparent;
  border-color: transparent;
  outline: none !important;
  outline-width: 0 !important;
  box-shadow: none;
  -moz-box-shadow: none;
  -webkit-box-shadow: none;
}

.fade-enter-active {
  transition: opacity 0.5s;
}
.fade-enter, .fade-leave-to /* .fade-leave-active below version 2.1.8 */ {
  opacity: 0;
}
</style>
