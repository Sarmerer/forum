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
          <template #badge>
            <div v-b-tooltip.left :title="userData.role == 2 ? 'Admin' : ''">
              <b-icon-code-slash></b-icon-code-slash>
            </div>
          </template>
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
            :datetime="userData.last_active"
            :long="!isMobile()"
            tooltip="right"
          >
          </time-ago>
        </small>
      </b-col>
    </b-row>
    <b-row class="mt-2" align-h="around">
      <div v-b-tooltip.hover.left title="Joined">
        <b-icon-calendar-event v-if="isMobile()"></b-icon-calendar-event>
        <small v-else>Joined</small>
        <br />
        <time-ago
          v-if="userData.created"
          tooltip="left"
          :datetime="userData.created"
        >
        </time-ago>
      </div>
      <div v-b-tooltip.hover.bottom title="Rating">
        <b-icon
          v-if="isMobile()"
          :icon="userData.rating >= 0 ? 'arrow-up' : 'arrow-down'"
        >
        </b-icon>
        <small v-else>Rating</small>
        <br />
        <span>{{ userData.rating }}</span>
      </div>

      <div v-b-tooltip.hover.bottom title="Posts">
        <b-icon-newspaper v-if="isMobile()"></b-icon-newspaper>
        <small v-else>Posts</small>
        <br />
        <span>{{ userData.posts }}</span>
      </div>

      <div v-b-tooltip.hover.right title="Comments">
        <b-icon-chat v-if="isMobile()"></b-icon-chat>
        <small v-else>Comments</small>
        <br />
        <span>{{ userData.comments }}</span>
      </div>
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