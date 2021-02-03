<template>
  <div :class="size === 'sm' ? 'rating-col-inline' : 'rating-col'">
    <b-icon-arrow-up-short
      @click.prevent="authenticated ? rate('up') : makeToast()"
      :class="
        `${!isMobile() ? 'mb-n1' : 'm-0'} ${classGen(
          entity.your_reaction > 0,
          'positive'
        )}`
      "
    >
    </b-icon-arrow-up-short>
    <span
      v-if="!isMobile() && (size === 'sm' || size === 'lg')"
      class="rating-opacity rating-item"
      >{{ entity.rating }}
    </span>
    <small
      v-if="isMobile() && (size === 'sm' || size === 'lg')"
      class="rating-opacity rating-item"
      >{{ entity.rating }}
    </small>
    <b-icon-arrow-down-short
      @click.prevent="authenticated ? rate('down') : makeToast()"
      :class="
        `${!isMobile() ? 'mt-n1 ' : 'm-0 '} ${classGen(
          entity.your_reaction < 0,
          'negative'
        )}`
      "
    ></b-icon-arrow-down-short>
  </div>
</template>
<script>
import { mapGetters } from "vuex";
import api from "@/api/api";

export default {
  props: {
    entity: {
      type: Object,
      required: true,
    },
    endpoint: {
      type: String,
      required: true,
    },
    compact: Boolean,
    size: { type: String, default: "md" },
  },
  computed: {
    ...mapGetters({
      authenticated: "auth/authenticated",
    }),
  },
  data() {
    return {
      requesting: false,
    };
  },
  methods: {
    classGen(reaction, variant) {
      return `rating-opacity rating-item ${this.size === "sm" ? "h4" : "h3"} ${
        reaction ? `${variant}-active` : this.isMobile() ? "" : variant
      }`;
    },
    async rate(reaction) {
      if (this.requesting) return;
      let r = reaction == "up" ? 1 : -1;
      if (
        (reaction == "up" && this.entity.your_reaction == 1) ||
        (reaction == "down" && this.entity.your_reaction == -1)
      ) {
        r = 0;
      }
      await api
        .post(`${this.endpoint}/rate`, { id: this.entity.id, reaction: r })
        .then((response) => {
          this.entity.your_reaction = response.data.data.your_reaction;
          this.entity.rating = response.data.data.rating;
        })
        .catch((error) => {
          console.log(error);
        })
        .then(() => (this.requesting = false));
    },
    makeToast() {
      this.$bvToast.toast(
        "You need to be logged in, to rate entitys and comments!",
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

.positive:hover {
  color: #28a745;
}

.positive-active {
  color: #28a745;
}

.negative:hover {
  color: #dc3545;
}

.negative-active {
  color: #dc3545;
}
</style>
