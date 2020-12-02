<template>
  <div :class="compact ? '' : 'rating-col'">
    <b-icon-arrow-up-short
      @click.prevent="authenticated ? callback('up', entity) : makeToast()"
      :class="`m-0 rating-opacity ${compact ? 'h4' : 'h3'}`"
      :style="`color: ${entity.your_reaction == 1 ? 'green' : 'white'}`"
    >
    </b-icon-arrow-up-short>
    <span class="rating-opacity">{{ entity.rating }}</span>
    <b-icon-arrow-down-short
      @click.prevent="authenticated ? callback('down', entity) : makeToast()"
      :class="`m-0 rating-opacity ${compact ? 'h4' : 'h3'}`"
      :style="`color: ${entity.your_reaction == -1 ? 'red' : 'white'}`"
    ></b-icon-arrow-down-short>
  </div>
</template>
<script>
import { mapGetters } from "vuex";
export default {
  props: {
    callback: {
      type: Function,
      required: true,
    },
    entity: {
      type: Object,
      required: true,
    },
    compact: {
      type: Boolean,
    },
  },
  computed: {
    ...mapGetters({
      authenticated: "auth/authenticated",
    }),
  },
  methods: {
    makeToast() {
      this.$bvToast.toast(
        "You need to be logged in, to rate posts and comments!",
        {
          title: "Oopsie!",
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

.rating-opacity {
  opacity: 87%;
}
</style>
q
