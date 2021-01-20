<template>
  <div>
    <router-link
      :to="'/post/' + comment.post_id"
      :class="
        `user-card text-break ${isMobile() ? 'card-m' : 'card card-hover'}`
      "
      v-for="comment in comments"
      :key="comment.id"
      tag="div"
    >
      <h5>
        {{ comment.content }}
      </h5>
      <small>
        <span v-b-tooltip.hover title="Rating">
          <b-icon
            :icon="reactionIcon(comment.your_reaction)"
            :color="reactionColor(comment.your_reaction)"
          >
          </b-icon
          >{{ comment.rating }}
        </span>
        <time-ago :datetime="comment.created" tooltip="right"> </time-ago>
      </small>
    </router-link>
  </div>
</template>
<script>
// import ControlButtons from "@/components/ControlButtons";
import TimeAgo from "@/components/TimeAgo";
import api from "@/api/api";

export default {
  name: "CommentsTab",
  mounted() {
    this.getComments();
  },
  data() {
    return {
      comments: [],
    };
  },
  components: { TimeAgo },
  methods: {
    reactionColor(yourReaction) {
      return yourReaction === 1 ? "green" : yourReaction === -1 ? "red" : "";
    },
    reactionIcon(yourReaction) {
      return yourReaction === -1 ? "arrow-down" : "arrow-up";
    },
    async getComments() {
      if (this.activeTab === "Comments") return;
      return await api
        .post("comments/find", {
          by: "author",
          author: Number.parseInt(this.$route.params.id),
        })
        .then((response) => {
          this.comments = response.data.data || [];
        })
        .catch((error) => {
          console.log(error);
        })
        .then(() => (this.madeRequest = true));
    },
  },
};
</script>
