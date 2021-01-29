<template>
  <div>
    <div v-if="!isMobile()" class="card">
      <b-row>
        <b-col align="center">
          <b-avatar
            v-if="userData.avatar"
            :src="userData.avatar"
            alt="avatar"
            size="6rem"
            variant="dark"
          >
          </b-avatar>
          <router-link v-if="link" :to="`/user/${userData.id}`" class="primary">
            <h3 class="mb-0">
              {{ userData.alias }}
            </h3>
          </router-link>
          <h3 v-else class="primary mb-0">
            {{ userData.alias }}
          </h3>
          <span class="text-muted">active </span>
          <time-ago
            class="text-muted"
            :datetime="userData.last_active"
            tooltip="right"
          >
          </time-ago>
        </b-col>
      </b-row>
      <b-row class="mt-2">
        <b-col>
          <p class="m-0 px-2">
            <span v-if="userData.role === 2">
              <b-icon-code-slash></b-icon-code-slash>
              Admin
              <br />
            </span>
            <b-icon-calendar-event></b-icon-calendar-event>
            <small> Joined:</small>
            <time-ago :datetime="userData.created" tooltip="right"> </time-ago>
            <br />
            <b-icon :icon="userData.rating >= 0 ? 'arrow-up' : 'arrow-down'">
            </b-icon>
            <small> Rating:</small> {{ userData.rating }}<br />
            <b-icon-newspaper></b-icon-newspaper>
            <small> Posts:</small> {{ userData.posts }}<br />
            <b-icon-chat></b-icon-chat>
            <small> Comments: </small>{{ userData.comments }}
          </p>
        </b-col>
      </b-row>
    </div>
    <div v-if="isMobile()" class="card-m">
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
            :badge="userData.role === 2"
            badge-variant="dark"
          >
            <template v-if="userData.role === 2" #badge>
              <div v-b-tooltip.left :title="userData.role == 2 ? 'Admin' : ''">
                <b-icon-code-slash></b-icon-code-slash>
              </div>
            </template>
          </b-avatar>
        </b-col>
        <b-col>
          <router-link v-if="link" :to="`/user/${userData.id}`">
            <h3 class="primary text-break">
              {{ userData.alias }}
            </h3>
          </router-link>
          <h3 v-else class="primary mb-0 text-break">
            {{ userData.alias }}
          </h3>
          <small class="text-muted">
            active
            <time-ago :datetime="userData.last_active"> </time-ago>
          </small>
        </b-col>
      </b-row>
      <b-row class="mt-2">
        <b-col align="center">
          <b-icon-calendar-event></b-icon-calendar-event>
          <p class="m-0">
            <time-ago :datetime="userData.created"> </time-ago>
          </p>
        </b-col>
        <b-col align="center">
          <b-icon :icon="userData.rating >= 0 ? 'arrow-up' : 'arrow-down'">
          </b-icon>
          <p class="m-0">
            {{ userData.rating }}
          </p>
        </b-col>
        <b-col align="center">
          <b-icon-newspaper></b-icon-newspaper>
          <p class="m-0">
            {{ userData.posts }}
          </p>
        </b-col>
        <b-col align="center">
          <b-icon-chat></b-icon-chat>
          <p class="m-0">
            {{ userData.comments }}
          </p>
        </b-col>
      </b-row>
    </div>
  </div>
</template>
<script>
import TimeAgo from "@/components/TimeAgo";

export default {
  name: "user-card",
  props: { userData: Object, link: Boolean },
  components: {
    TimeAgo,
  },
};
</script>
