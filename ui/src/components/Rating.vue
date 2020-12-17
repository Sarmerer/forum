<template>
  <div :class="size === 'sm' ? 'rating-col-inline' : 'rating-col'">
    <b-icon-arrow-up-short
      @click.prevent="
        authenticated ? $emit('rate', ['up', entity]) : makeToast()
      "
      :class="`m-0 rating-opacity rating-item ${classUp}`"
      :style="`color: ${entity.your_reaction == 1 ? 'green' : 'white'}`"
    >
    </b-icon-arrow-up-short>
    <span
      v-if="size === 'sm' || size === 'lg'"
      class="rating-opacity rating-item"
      >{{ entity.rating }}</span
    >
    <b-icon-arrow-down-short
      @click.prevent="
        authenticated ? $emit('rate', ['down', entity]) : makeToast()
      "
      :class="`m-0 rating-opacity rating-item ${classDown}`"
      :style="`color: ${entity.your_reaction == -1 ? 'red' : 'white'}`"
    ></b-icon-arrow-down-short>
  </div>
</template>
<script>
import { mapGetters } from "vuex";
export default {
  props: {
    entity: {
      type: Object,
      required: true,
    },
    compact: Boolean,
    size: { type: String, default: "md" },
  },
  computed: {
    ...mapGetters({
      authenticated: "auth/authenticated",
    }),
    classUp() {
      return this.size === "sm" ? "h4" : "h3 mb-n1";
    },
    classDown() {
      return this.size === "sm" ? "h4" : "h3 mt-n1";
    },
  },
  methods: {
    makeToast() {
      this.$bvToast.toast(
        "You need to be logged in, to rate posts and comments!",
        {
          title: "Oops!",
          variant: "danger",
          solid: true,
        }
      );
    },
  },
};
</script>
<style lang="scss">
.rating-col {
  display: flex;
  flex-direction: column;
  text-align: left;
  align-items: center;
}

.rating-col-inline {
  display: flex;
  align-items: center;
}

.rating-item {
  vertical-align: middle;
}

.rating-opacity {
  opacity: 0.87;
}
</style>
