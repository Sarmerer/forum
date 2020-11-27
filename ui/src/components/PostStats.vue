<template>
  <div :class="isMobile() ? 'card-m' : 'card'">
    <p>{{ pluralize(stats.commentsCount, "comment") }}</p>
    <p>{{ pluralize(stats.participantsCount, "participant") }}</p>
    <p v-if="stats.lastCommentFromName">
      Last comment from:
      <router-link :to="`/user/${stats.lastCommentFromID}`">{{
        stats.lastCommentFromName
      }}</router-link>
    </p>
    <p v-if="stats.lastCommentDate">
      Last activity:
      <time-ago :datetime="stats.lastCommentDate" :long="!isMobile()">
      </time-ago>
    </p>
  </div>
</template>
<script>
import TimeAgo from "vue2-timeago";

export default {
  props: {
    stats: { type: Object, required: true },
  },
  components: {
    TimeAgo,
  },
  methods: {
    pluralize(n, s) {
      return `${n} ${s}${n == 1 ? "" : "s"}`;
    },
  },
};
</script>
