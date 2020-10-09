<template>
  <div class="rating-column mr-2" style="text-align:center;">
    <svg
      style="display:block"
      @click.prevent="rate('up')"
      width="1.5em"
      height="1.5em"
      viewBox="0 0 16 16"
      class="bi bi-chevron-up"
      :fill="yourReaction == 1 ? 'green' : 'white'"
      xmlns="http://www.w3.org/2000/svg"
    >
      <path
        fill-rule="evenodd"
        d="M7.646 4.646a.5.5 0 0 1 .708 0l6 6a.5.5 0 0 1-.708.708L8 5.707l-5.646 5.647a.5.5 0 0 1-.708-.708l6-6z"
      />
    </svg>
    <span>{{ rating }}</span>
    <svg
      style="display:block"
      @click.prevent="rate('down')"
      width="1.5em"
      height="1.5em"
      viewBox="0 0 16 16"
      class="bi bi-chevron-down"
      :fill="yourReaction == -1 ? 'red' : 'white'"
      xmlns="http://www.w3.org/2000/svg"
    >
      <path
        fill-rule="evenodd"
        d="M1.646 4.646a.5.5 0 0 1 .708 0L8 10.293l5.646-5.647a.5.5 0 0 1 .708.708l-6 6a.5.5 0 0 1-.708 0l-6-6a.5.5 0 0 1 0-.708z"
      />
    </svg>
  </div>
</template>
<script>
import axios from "axios";

export default {
  props: {
    postID: {
      type: Number,
      required: true,
      default: -1,
    },
    rating: {
      type: Number,
      required: true,
      default: 0,
    },
    yourReaction: {
      type: Number,
      required: true,
      default: 0,
    },
  },
  data() {
    return {
      requesting: false,
    };
  },
  methods: {
    async rate(reaction) {
      if (this.requesting) return;
      this.requesting = true;
      let self = this;
      let r = reaction == "up" ? 1 : -1;
      if (
        (reaction == "up" && this.yourReaction == 1) ||
        (reaction == "down" && this.yourReaction == -1)
      ) {
        r = 0;
      }
      await axios
        .post("post/rate", { pid: self.postID, reaction: r })
        .then((response) => {
          // self.yourReaction = response.data.data.your_reaction;
          // self.rating = response.data.data.rating;
          this.$emit("update", {
            new_rating: response.data.data.rating,
            new_your_reaction: response.data.data.your_reaction,
          });
        })
        .catch((error) => {
          console.log(error);
          //TODO show alert saying that you need to be logged in to rate
        })
        .finally(() => {
          self.requesting = false;
        });
    },
  },
};
</script>
