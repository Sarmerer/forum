<template>
  <div class="mr-2">
    <transition name="fade">
      <b-button-group v-if="hasPermission && !confirm && !compact" size="sm">
        <b-button
          class="controls-button"
          variant="light"
          v-b-tooltip.hover.bottom="'Edit'"
          :disabled="disabled"
          @click="$emit('edit-event')"
        >
          <b-icon-pencil-square color="white"></b-icon-pencil-square>
        </b-button>

        <b-button
          variant="outline-danger"
          class="controls-button"
          v-b-tooltip.hover.bottom="'Delete'"
          title="Delete"
          @click="confirm = true"
          :disabled="disabled"
        >
          <b-icon-trash color="red"></b-icon-trash>
        </b-button>
      </b-button-group>
      <b-button-group v-if="hasPermission && confirm" size="sm">
        <b-button
          variant="outline-success"
          class="confirm"
          @click="$emit('delete-event'), (confirm = false)"
          v-b-tooltip.hover.bottom="'Confirm'"
          :disabled="disabled"
        >
          <b-icon-check></b-icon-check>
        </b-button>
        <b-button
          variant="outline-danger"
          @click="confirm = false"
          class="confirm"
          v-b-tooltip.hover.bottom="'Cancel'"
          :disabled="disabled"
        >
          <b-icon-x></b-icon-x>
        </b-button>
      </b-button-group>
    </transition>
    <b-icon-three-dots
      v-if="hasPermission && compact"
      @click="$bvModal.show(modalID)"
    >
    </b-icon-three-dots>
  </div>
</template>
<script>
import { mapGetters } from "vuex";

export default {
  props: {
    hasPermission: { type: Boolean, required: true },
    disabled: { type: Boolean, required: true },
    compact: { type: Boolean },
    modalID: { type: String },
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
.controls {
  position: absolute;
  top: 5px;
  right: 10px;
}
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
