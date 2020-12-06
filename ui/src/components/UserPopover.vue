<template>
  <span :id="`user-info-popover${popoverID}`">
    <b-avatar
      v-if="!noAvatar"
      :src="userData.avatar"
      variant="dark"
      size="15px"
    ></b-avatar>
    <b-popover
      :target="`user-info-popover${popoverID}`"
      triggers="hover"
      :placement="popoverDirection || 'left'"
    >
      <span>
        <span v-if="userData.role === 2">
          <b-icon-code-slash></b-icon-code-slash>
          Admin
          <br />
        </span>
        <b-icon-calendar-event></b-icon-calendar-event>
        <small> Joined: </small>
        <time-ago :datetime="userData.created"> </time-ago>
        <br />
        <b-icon-clock></b-icon-clock>
        <small> Active: </small>
        <time-ago :datetime="userData.last_active"> </time-ago>
        <br />
        <b-icon :icon="userData.rating >= 0 ? 'arrow-up' : 'arrow-down'">
        </b-icon>
        <small> Rating:</small> {{ userData.rating }}<br />
        <b-icon-newspaper></b-icon-newspaper>
        <small> Posts:</small> {{ userData.posts }}<br />
        <b-icon-chat></b-icon-chat>
        <small> Comments: </small>{{ userData.comments }}
      </span>
    </b-popover>
    {{ userData.display_name }}
  </span>
</template>
<script>
import TimeAgo from "@/components/TimeAgo";
export default {
  name: "user-popover",
  props: {
    userData: { type: Object, required: true },
    popoverID: String,
    popoverDirection: String,
    noAvatar: Boolean,
  },
  components: {
    TimeAgo,
  },
};
</script>
<style lang="scss" scoped>
.popover {
  background-color: #282828;
  opacity: 0.87;
}

.popover span {
  color: white;
  opacity: 0.87;
}
</style>
