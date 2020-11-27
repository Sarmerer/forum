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
      Last activity: {{ stats.lastCommentDate | formatDate }}
    </p>
  </div>
</template>
<script>
export default {
  props: {
    stats: { type: Object, required: true },
  },
  methods: {
    pluralize(n, s) {
      return `${n} ${s}${n == 1 ? "" : "s"}`;
    },
  },
};
</script>
