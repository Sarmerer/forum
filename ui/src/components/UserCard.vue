<template>
  <div :class="isMobile() ? 'card-m' : 'card'">
    <div v-if="!isMobile()">
      <b-row >
        <b-col align="center">
          <b-img-lazy
            :src="userData.avatar"
            alt="avatar"
            rounded="circle"
            class="mb-1"
            width="100px"
          >
          </b-img-lazy>
          <h3 class="primary">
            {{ userData.display_name }}
            <b-badge v-if="userData.role == 2" class="background-variant">
              <b-icon-people></b-icon-people>
            </b-badge>
          </h3>
        </b-col>
      </b-row>
      <b-row>
        <b-col>
          <p class="m-0 px-2">
            Last active:
            <time-ago
              v-b-tooltip.hover
              :title="userData.last_online"
              :datetime="userData.last_online"
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
          <b-img-lazy
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
              <b-icon-people></b-icon-people>
            </b-badge>
          </h4>
        </b-col>
        <b-col>
          <p class="m-0">
            Last active:
            <time-ago
              v-b-tooltip.hover
              :title="userData.last_online"
              :datetime="userData.last_online"
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
  props: { userData: { type: Object } },
  components: {
    TimeAgo,
  },
};
</script>
<style lang="scss" scoped>
.user-card {
  cursor: pointer;
}

// .user-card img {
//   display: block;
//   margin-left: auto;
//   margin-right: auto;
//   margin-bottom: 15px;
//   width: 150px;
//   height: 150px;
//   border-radius: 150px;
// }
</style>
