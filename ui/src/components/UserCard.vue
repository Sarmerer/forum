<template>
  <div :class="isMobile() ? 'card-m' : 'card'">
    <b-row>
      <b-col cols="start">
        <b-avatar
          v-if="userData.avatar"
          :src="userData.avatar"
          alt="avatar"
          size="lg"
          class="ml-3"
          variant="dark"
          badge-offset="-0.2em"
          :badge="userData.role > 0"
          :badge-right="userData.role > 0"
          badge-variant="info"
        >
          <template #badge><b-icon-code-slash></b-icon-code-slash></template>
        </b-avatar>
      </b-col>
      <b-col>
        <router-link v-if="link" :to="`/user/${userData.id}`">
          <h3 class="primary">
            {{ userData.display_name }}
          </h3>
        </router-link>
        <h3 v-else class="primary mb-0">
          {{ userData.display_name }}
        </h3>
        <small>
          Active:
          <time-ago
            v-if="userData.last_active"
            v-b-tooltip.hover
            :title="userData.last_active"
            :datetime="userData.last_active"
            :long="!isMobile()"
          >
          </time-ago>
        </small>
      </b-col>
    </b-row>
    <b-row class="ml-2 mt-2">
      <b-col>
        <div v-b-tooltip.hover.left title="Joined">
          <b-icon-calendar-event></b-icon-calendar-event>
        </div>
        <time-ago
          v-if="userData.created"
          v-b-tooltip.hover.bottom
          :title="userData.created"
          :datetime="userData.created"
        >
        </time-ago>
      </b-col>
      <b-col>
        <div v-b-tooltip.hover.bottom title="Rating">
          <b-icon :icon="userData.rating >= 0 ? 'arrow-up' : 'arrow-down'">
          </b-icon>
        </div>
        {{ userData.rating }}
      </b-col>
      <b-col>
        <div v-b-tooltip.hover.bottom title="Posts">
          <b-icon-newspaper></b-icon-newspaper>
        </div>
        {{ userData.posts }}
      </b-col>
      <b-col>
        <div v-b-tooltip.hover.right title="Comments">
          <b-icon-chat></b-icon-chat>
        </div>
        {{ userData.comments }}
      </b-col>
    </b-row>
  </div>
</template>
<script>
import TimeAgo from "vue2-timeago";

export default {
  name: "user-card",
  props: { userData: Object, link: Boolean },
  components: {
    TimeAgo,
  },
};
</script>
