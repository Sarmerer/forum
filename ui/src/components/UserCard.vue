<template>
  <div :class="isMobile() ? 'card-m' : 'card'">
    <div v-if="!isMobile()">
      <b-row>
        <b-col align="center">
          <b-img-lazy
            v-if="userData.avatar"
            :src="userData.avatar"
            alt="avatar"
            rounded="circle"
            class="mb-1"
            width="100px"
          >
          </b-img-lazy>
          <router-link v-if="link" :to="`/user/${userData.id}`">
            <h3 class="primary">
              {{ userData.display_name }}
              <b-badge v-if="userData.role == 2" class="background-variant">
                <b-icon-code-slash></b-icon-code-slash>
              </b-badge>
            </h3>
          </router-link>
          <h3 v-else class="primary">
            {{ userData.display_name }}
            <b-badge v-if="userData.role == 2" class="background-variant">
              <b-icon-code-slash></b-icon-code-slash>
            </b-badge>
          </h3>
        </b-col>
      </b-row>
      <b-row>
        <b-col>
          <p class="m-0 px-2">
            Joined:
            <time-ago
              v-if="userData.created"
              v-b-tooltip.hover
              :title="userData.created"
              :datetime="userData.created"
              :long="!isMobile()"
            >
            </time-ago>
            <br />
            Active:
            <time-ago
              v-if="userData.last_active"
              v-b-tooltip.hover
              :title="userData.last_active"
              :datetime="userData.last_active"
              :long="!isMobile()"
            >
            </time-ago>
            <br />
            Rating: {{ userData.rating }}<br />
            Posts: {{ userData.posts }}<br />
            Comments: {{ userData.comments }}<br />
          </p>
        </b-col>
      </b-row>
    </div>

    <div v-else>
      <b-row align-v="center">
        <b-col align="center">
          <router-link v-if="link" :to="`/user/${userData.id}`">
            <b-img-lazy
              v-if="userData.avatar"
              :src="userData.avatar"
              alt="avatar"
              rounded="circle"
              class="mb-1"
              width="70px"
            >
            </b-img-lazy>
            <h4 class="primary m-0">
              {{ userData.display_name }}
              <b-badge v-if="userData.role == 2" class="background-variant">
                <b-icon-code-slash></b-icon-code-slash>
              </b-badge>
            </h4>
          </router-link>

          <div v-else>
            <b-img-lazy
              v-if="userData.avatar"
              :src="userData.avatar"
              alt="avatar"
              rounded="circle"
              class="mb-1"
              width="70px"
            >
            </b-img-lazy>
            <h4 class="primary m-0">
              {{ userData.display_name }}
              <b-badge v-if="userData.role == 2" class="background-variant">
                <b-icon-code-slash></b-icon-code-slash>
              </b-badge>
            </h4>
          </div>
        </b-col>
        <b-col>
          <p class="m-0">
            Active:
            <time-ago
              v-b-tooltip.hover
              :title="userData.last_active"
              :datetime="userData.last_active"
              :long="!isMobile()"
            >
            </time-ago>
            <br />
            Rating: {{ userData.rating }}<br />
            Posts: {{ userData.posts }}<br />
            Comments: {{ userData.comments }}<br />
          </p>
        </b-col>
      </b-row>
    </div>
  </div>
</template>
<script>
import TimeAgo from "vue2-timeago";

export default {
  name: "user-card",
  props: { userData: { type: Object }, link: { type: Boolean } },
  components: {
    TimeAgo,
  },
};
</script>
