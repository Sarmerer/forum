<template>
  <transition name="fade">
    <b-button-group v-if="hasPermission && !confirm" size="sm">
      <b-button
        size="sm"
        lg="1"
        class="controls-button"
        variant="light"
        title="Edit"
        :disabled="disabled"
        @click="call(editCallback)"
      >
        <b-icon-pencil-square color="white"></b-icon-pencil-square>
      </b-button>
      <b-button
        size="sm"
        variant="outline-danger"
        lg="2"
        class="controls-button"
        title="Delete"
        @click="confirm = true"
        :disabled="disabled"
      >
        <b-icon-trash color="red"></b-icon-trash>
      </b-button>
    </b-button-group>
    <b-button-group v-if="hasPermission && confirm" size="sm">
      <b-button
        size="sm"
        variant="outline-success"
        lg="2"
        class="confirm"
        @click="call(deleteCallback).then((confirm = false))"
        title="Confirm"
        :disabled="disabled"
      >
        <b-icon-check></b-icon-check>
      </b-button>
      <b-button
        size="sm"
        lg="1"
        variant="outline-danger"
        @click="confirm = false"
        class="confirm"
        title="Dismiss"
        :disabled="disabled"
      >
        <b-icon-x></b-icon-x>
      </b-button>
    </b-button-group>
  </transition>
</template>
<script>
import { mapGetters } from "vuex";

export default {
  props: {
    hasPermission: { type: Boolean, required: true },
    deleteCallback: { type: Object, required: true },
    editCallback: { type: Object, required: true },
    disabled: { type: Boolean, required: true },
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
  methods: {
    call(prop) {
      return prop.args ? prop.callback(...prop.args) : prop.callback();
    },
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
